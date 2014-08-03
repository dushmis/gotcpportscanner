package main

//just for learning purpose..

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var host string
var start, end, port, timeout int

var verbose bool

type TCPLocation struct {
	Host string
	Port int
}

type Result struct {
	TCPLocation *TCPLocation
	Err         string
	IsOpen      bool
}

func (this *Result) String() string {
	return fmt.Sprintf("{ TCPLocation:'%s', Err:'%s', IsOpen:'%t' }", this.TCPLocation, this.Err, this.IsOpen)
}

func (this *TCPLocation) String() string {
	return fmt.Sprintf("{ host=%s, port=%d }", this.Host, this.Port)
}

func (this *TCPLocation) Scan() *Result {
	var p string
	p = fmt.Sprintf("%s:%d", this.Host, this.Port)
	// if !verbose {
	// fmt.Printf("Scanning -- %s\n", p)
	// }
	_log_(fmt.Sprintf("Connecting %s", this))
	conn, err := net.DialTimeout("tcp", p, time.Duration(timeout)*time.Second)
	if err != nil {
		_log_(fmt.Sprintf("Error %s", err.Error()))
		return &Result{TCPLocation: this, Err: err.Error(), IsOpen: false}
	}
	defer func() {
		conn.Close()
	}()
	return &Result{TCPLocation: this, IsOpen: true}
}

func _log_(logs string) {
	if verbose {
		log.Printf("\033[92m%s\033[0m\n", logs)
	}
}

func main() {
	flag.StringVar(&host, "host", "localhost", "IP Address or host")
	flag.IntVar(&start, "start", 20, "start port")
	flag.IntVar(&end, "end", 25, "end port")
	flag.IntVar(&timeout, "w", 5, "tcp timeout in seconds")
	flag.IntVar(&port, "p", -1, "specific port, this option ignores start/end flag")

	flag.BoolVar(&verbose, "v", false, "verbose")

	var so = flag.Bool("s", false, "success only logs")

	flag.Parse()

	_log_(fmt.Sprintf("scanning from %s:%d, %s:%d", host, start, host, end))
	_log_(fmt.Sprintf("timeout - %d", timeout))

	singlePort := false

	if port != -1 {
		start = port
		end = start + 1
		singlePort = true
	}

	IPAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(2)
	}

	host = IPAddr.IP.String()

	var wg sync.WaitGroup

	var portRange int

	if portRange = end - start; portRange < 0 {
		log.Fatal("Invalid port range...")
		os.Exit(3)
	}

	TCPLocations := make([]*TCPLocation, portRange)

	j := 0
	for i := start; i < end; i++ {
		TCPLocations[j] = (&TCPLocation{host, i})
		j++
	}

	for i := range TCPLocations {
		wg.Add(1)
		go func(tcpLocation *TCPLocation) {
			result := tcpLocation.Scan()
			_log_(fmt.Sprintf("result - %s", result))
			if result.IsOpen {
				fmt.Printf("\033[92mSUCCESS\033[0m - %s:%d\n", result.TCPLocation.Host, result.TCPLocation.Port)
			} else {
				if !*so {
					fmt.Printf("\033[91mFAILURE\033[0m - %s:%d\n", result.TCPLocation.Host, result.TCPLocation.Port)
				}
				if singlePort {
					os.Exit(1)
				}
			}
			defer wg.Done()
		}(TCPLocations[i])
	}
	wg.Wait()
}
