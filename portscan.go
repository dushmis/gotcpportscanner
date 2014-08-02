package main

//just for learning purpose..

import (
	"flag"
	"fmt"
	"net"
	"time"
)

var host string
var start, end int

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
	fmt.Printf("Scanning -- %s\n", p)
	conn, err := net.DialTimeout("tcp", p, 1*time.Second)
	if err != nil {
		return &Result{TCPLocation: this, Err: err.Error(), IsOpen: false}
	}
	defer func() {
		conn.Close()
	}()
	return &Result{TCPLocation: this, IsOpen: true}
}

func main() {
	flag.StringVar(&host, "host", "localhost", "host address")
	flag.IntVar(&start, "start", 20, "start")
	flag.IntVar(&end, "end", 25, "end")
	flag.Parse()

	for i := start; i < end; i++ {
		result := (&TCPLocation{host, i}).Scan()
		if result.IsOpen {
			fmt.Printf("Is Open - %s\n", result)
		}
	}
}
