package main

import (
	"flag"
	"fmt"

	"github.com/bikefrivolously/go-to-aws/internal/server"
)

func main() {
	var port = flag.Int("port", 8000, "Listen port for the web server.")
	var address = flag.String("address", "localhost", "Listen address for the web server.")
	flag.Parse()
	fmt.Printf("Starting server on %s:%d...\n", *address, *port)
	s := server.Server{Port: *port, Address: *address}
	s.Run()
}
