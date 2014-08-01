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

type Host struct {
	Ip string
}

func (this *Host) IsUp() bool {
	return true
}

type TCPLocation struct {
	Host *Host
	Port int
}

func (this *TCPLocation) String() string {
	return fmt.Sprintf("{ host=%s, port=%d }", this.Host, this.Port)
}

func (this *TCPLocation) IsOpen() bool {
	var p string
	p = fmt.Sprintf("%s:%d", this.Host.Ip, this.Port)
	fmt.Printf("%s ", p)
	conn, err := net.DialTimeout("tcp", p, 3*time.Second)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return false
	}
	defer func() {
		conn.Close()
	}()
	return true
}

func main() {
	flag.StringVar(&host, "host", "localhost", "host address")
	flag.IntVar(&start, "start", 20, "start")
	flag.IntVar(&end, "end", 25, "end")
	flag.Parse()

	for i := start; i < end; i++ {
		fmt.Printf("Is Open - %t\n", (&TCPLocation{&Host{host}, i}).IsOpen())
	}
}
