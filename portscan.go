package main

//just for learning purpose..

import (
	"flag"
	"fmt"
	"net"
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
	if !verbose {
		fmt.Printf("Scanning -- %s\n", p)
	}
	log(fmt.Sprintf("Connecting %s", this))
	conn, err := net.DialTimeout("tcp", p, time.Duration(timeout)*time.Second)
	if err != nil {
		log(fmt.Sprintf("Error %s", err.Error()))
		return &Result{TCPLocation: this, Err: err.Error(), IsOpen: false}
	}
	defer func() {
		conn.Close()
	}()
	return &Result{TCPLocation: this, IsOpen: true}
}

func log(log string) {
	if verbose {
		fmt.Printf("[ %s ]: %s\n", time.Now().Format(time.RFC3339Nano), log)
	}
}

func main() {
	flag.StringVar(&host, "host", "localhost", "host address")
	flag.IntVar(&start, "start", 20, "start")
	flag.IntVar(&end, "end", 25, "end")
	flag.IntVar(&timeout, "w", 5, "tcp timeout in seconds")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.IntVar(&port, "p", -1, "specific port, this option ignores start/end flag")

	flag.Parse()

	log(fmt.Sprintf("scanning from %s:%d, %s:%d", host, start, host, end))
	log(fmt.Sprintf("timeout - %d", timeout))

	if port != -1 {
		start = port
		end = start + 1
	}

	for i := start; i < end; i++ {
		result := (&TCPLocation{host, i}).Scan()
		log(fmt.Sprintf("result - %s", result))
		if result.IsOpen {
			fmt.Printf("Is Open - %s\n", result)
		}
	}
}
