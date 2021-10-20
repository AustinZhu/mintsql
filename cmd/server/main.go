package main

import (
	"log"
	"os"
)

func init() {
	if len(os.Args) > 1 {
		port = os.Args[1]
		sqlServer = New(port)
	} else {
		sqlServer = New(PORT)
	}
	log.Printf("Welcome to MintSQL Server.")
}

func main() {
	sqlServer.Run()
}
