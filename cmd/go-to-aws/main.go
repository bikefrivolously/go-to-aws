package main

import (
	"fmt"

	"github.com/bikefrivolously/go-to-aws/internal/server"
)

func main() {
	fmt.Println("Starting server...")
	s := server.Server{}
	s.Run()
}
