package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	var server *Server
	if len(os.Args) > 2 {
		host, port := os.Args[1], os.Args[2]
		p, err := strconv.Atoi(port)
		if err != nil {
			log.Fatalln("not a valid port number", err)
		}
		server = New(host, uint16(p))
	}
	server = New(HOST, PORT)
	server.Run()
}
