package main

import (
	"flag"
)

func main() {
	var (
		addr string = "localhost:3000"
	)

	flag.StringVar(&addr, "addr", "localhost:3000", "")
	flag.Parse()

	server := NewServer(addr)
	StartServer(server)
}
