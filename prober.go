package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alouca/gosnmp"

	"flag"
	"fmt"
)

var errMessage string

func usage(msg string, code int) {
	errMessage = msg
	flag.Usage()
	os.Exit(code)
}
func main() {
	flag.Usage = func() {
		if errMessage != "" {
			fmt.Fprintf(os.Stderr, "%s\n", errMessage)
		}
		fmt.Fprintf(os.Stderr, "Usage of %s: <ip>:<community>:<port> <probing interval> <oid1> <oid2> ......<oidn>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "probing interval:\n  the delay between probes in seconds\n")
		fmt.Fprintf(os.Stderr, "oid[1 .... n]:\n  OIDS which you want to probe\n")
		flag.PrintDefaults()
	}

	counter := map[string][]float64{}
	timestamps := map[string][]float64{}
	argsWithoutProg := os.Args[1:]
	oids := os.Args[3:]
	addr := strings.Split(argsWithoutProg[0], ":")
	s, err := gosnmp.NewGoSNMP(addr[0], addr[1], gosnmp.Version2c, 5)
	if err != nil {
		log.Fatal(err)
	}
	dur, err := strconv.Atoi(argsWithoutProg[1])
	if err != nil {
		usage("Give me an valid time interval to probe", 2)
	}

	for {
		start := time.Now()
		for _, y := range oids {

			resp, err := s.Get(y)
			if err == nil {
				for _, v := range resp.Variables {
					switch v.Type {
					case gosnmp.Integer:

					default:
						i := float64(v.Value.(uint64))
						counter[v.Name] = append(counter[v.Name], i)
						timestamps[v.Name] = append(timestamps[v.Name], float64(time.Now().Unix()))

						ratecalculator(v.Name, len(counter[v.Name]), counter[v.Name], timestamps[v.Name])

					}
				}
			}

		}
		end := time.Now()
		delay := end.Sub(start)
		time.Sleep(time.Duration(dur)*time.Second - delay)

	}
}

func ratecalculator(oid string, i int, counter []float64, timestamps []float64) {
	if i != 1 {
		fmt.Println(fmt.Sprintf("rate of change of oid %s =  ", oid), (counter[i-1]-counter[i-2])/(timestamps[i-1]-timestamps[i-2]), fmt.Sprintf("for the interval = %d  ", i))
	} else {
		fmt.Println(fmt.Sprintf("rate of change of oid %s =  ", oid), (counter[i-1])/(timestamps[i-1]), fmt.Sprintf("for the interval = %d  ", i))
	}

}
