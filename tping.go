package main

import (
	"flag"
	"fmt"
	"net"
        "os/signal"
	"os"
	"time"
)

var usage = `
Usage:

    tping [-c count] [-i interval] [-t timeout]  [-addr host ]

Examples:

    #i tping example continuously
   ./tping -addr 127.0.0.1:1234

    # tping default count  50 times
    ./tping -c 50 -addr 127.0.0.1:1234

    # tping default 1000 Period of check 
    ./tping -c 5 -i 1000 -addr 127.0.0.1:1234

    # tping default 500 for Connection timeout in ms. Should be less than check period
    ./tping -t 500 -addr 127.0.0.1:1234
`

func main() {
	var addrArg string
	var countArg int
	var timeoutArg int
	var periodArg int
	f := 0
	s := 0

	flag.StringVar(&addrArg, "addr", "", "provide the ip address with port ex: 127.0.0.1:1235")
	flag.IntVar(&countArg, "c", 50, "number of lines to read from the file")
	flag.IntVar(&periodArg, "i", 1000, "Period of check in ms")
	flag.IntVar(&timeoutArg, "t", 500, "Connection timeout in ms. Should be less than check period")

	flag.Parse()
	if addrArg == "" {
		fmt.Println("provide the ip address with port ex: 127.0.0.1:1235")
		fmt.Printf(usage)
		os.Exit(1)
	}

	var seqNumber uint64 = 0
	var timeout = time.Duration(timeoutArg) * time.Millisecond
	var period = time.Duration(periodArg) * time.Millisecond

	if timeout >= period {
		fmt.Println("timeout should be less than period")
		os.Exit(1)
	}

        e := make(chan os.Signal, 1)
        signal.Notify(e, os.Interrupt) 
        go func() {
                for _ = range e {
                 fmt.Printf("\n       ------ ping statistics ------\n")
                        fmt.Printf("%d packets transmitted, %d received, %d packet loss \n", s+f, s, f)
                os.Exit(1)
                  }
        }()

	ticker := time.NewTicker(period)
	quit := make(chan interface{})

	for c := 1; c <= countArg; c++ {
		seqNumber++
		select {
		case <-ticker.C:
			startTime := time.Now()
			conn, err := net.DialTimeout("tcp", addrArg, timeout)
			endTime := time.Now()
			if err != nil {
				os.Stdout.Write([]byte(startTime.Format("[2006-01-02T15:04:05]:") + " connection failed\n"))
				f++
			} else {
				defer conn.Close()
				var t = float64(endTime.Sub(startTime)) / float64(time.Millisecond)
				os.Stdout.Write([]byte(startTime.Format("[2006-01-02T15:04:05]:") + fmt.Sprintf(" addr=%s seq=%d time=%4.2fms\n", conn.RemoteAddr().String(), seqNumber, t)))
				s++
			}
		case <-quit:
			ticker.Stop()
			return
		}
		if c == countArg {
			fmt.Printf("\n       ------ ping statistics ------\n")
			fmt.Printf("%d packets transmitted, %d received, %d packet loss \n", s+f, s, f)
		}
	}
}
