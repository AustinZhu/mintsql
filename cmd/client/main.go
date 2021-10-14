package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Welcome to mintsql Client.")
	serverAddr, err := net.ResolveTCPAddr("tcp", ":7384")
	if err != nil {
		log.Fatalln("unable to resolve address:", err)
	}
	tcp, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Fatalln("unable to connect to server:", err)
	}
	reader := bufio.NewReader(os.Stdin)
	line := ""
	for line != "quit" {
		line, err = reader.ReadString('\n')
		_, err := tcp.Write([]byte(line))
		if err != nil {
			log.Println("unable to write:", err)
		}
		res, err := bufio.NewReader(tcp).ReadString('\n')
		fmt.Println(res)
	}
}
