package main

import (
	"flag"
	"github.com/initialed85/frinkconv-api/internal/helpers"
	"github.com/initialed85/frinkconv-api/pkg/frinkconv_server"
	"log"
)

func main() {
	port := flag.Int("port", 8080, "HTTP port to listen on")
	processes := flag.Int("processes", 4, "Number of frinkconv REPL processes to spawn")

	flag.Parse()

	server, err := frinkconv_server.New(*port, *processes)
	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	helpers.WaitForSigInt()
}
