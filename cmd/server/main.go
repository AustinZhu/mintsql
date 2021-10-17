package main

import (
	"log"
	"os"
)

func init() {
	if len(os.Args) > 2 {
		host, port = os.Args[1], os.Args[2]
		sqlServer = New(host, port)
	} else {
		sqlServer = New(HOST, PORT)
	}
	log.Printf("Welcome to mintsql Server.")
}

func main() {
	sqlServer.Run()
}
