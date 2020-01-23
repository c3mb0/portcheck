package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

type flags struct {
	wait    bool
	address string
	port    int
}

var opts = new(flags)

func init() {
	flag.StringVar(&opts.address, "a", "127.0.0.1", "address to check")
	flag.IntVar(&opts.port, "p", 0, "port to check")
	flag.BoolVar(&opts.wait, "w", false, "should wait until available")
	flag.Parse()
}

func main() {
	if opts.port == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if !opts.wait {
		os.Exit(isOpen())
	}
	for isOpen() != 0 {
		time.Sleep(time.Second)
	}
}

func isOpen() int {
	nc, err := net.Dial("tcp", fmt.Sprintf("%s:%d", opts.address, opts.port))
	if err != nil {
		return 1
	}
	defer nc.Close()
	nc.SetReadDeadline(time.Now().Add(time.Second))
	_, err = nc.Read(make([]byte, 1))
	if e, ok := err.(*net.OpError); err == nil || (ok && e.Timeout()) {
		return 0
	}
	return 1
}
