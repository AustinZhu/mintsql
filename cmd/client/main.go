package main

import (
	"fmt"
	"os"
)

func init() {
	if len(os.Args) > 2 {
		host, port = os.Args[1], os.Args[2]
		sqlClient = New(host, port)
	} else {
		sqlClient = New(HOST, PORT)
	}
	fmt.Println("Welcome to MintSQL Client.")
}

func main() {
	sqlClient.Run()
}
