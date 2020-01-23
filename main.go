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
	var nc net.Conn
	var err error
	fn := func(code int) {
		nc.Close()
		os.Exit(code)
	}
	if !opts.wait {
		nc, err = dial()
		if err != nil {
			os.Exit(1)
		}
		if isOpen(nc) {
			fn(0)
		}
		fn(1)
	}
	for {
		nc, err = dial()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	for {
		if isOpen(nc) {
			fn(0)
		}
		time.Sleep(time.Second)
	}
}

func dial() (net.Conn, error) {
	nc, err := net.Dial("tcp", fmt.Sprintf("%s:%d", opts.address, opts.port))
	if err != nil {
		return nil, err
	}
	nc.SetReadDeadline(time.Now().Add(time.Second))
	return nc, nil
}

func isOpen(nc net.Conn) bool {
	_, err := nc.Read(make([]byte, 1))
	if e, ok := err.(*net.OpError); err == nil || (ok && e.Timeout()) {
		return true
	}
	return false
}
