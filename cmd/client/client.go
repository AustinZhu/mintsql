package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	HOST     = "127.0.0.1"
	PROTOCOL = "tcp"
	PORT     = "7384"
)

var (
	host      string
	port      string
	sqlClient *Client
)

type Client struct {
	Addr *net.TCPAddr
	Conn *net.TCPConn
}

func New(host string, port string) *Client {
	port_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalln("invalid port number", err)
	}
	return &Client{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP(host),
			Port: port_,
		},
		Conn: nil,
	}
}

func (c *Client) Run() {
	var err error
	c.Conn, err = net.DialTCP(PROTOCOL, nil, c.Addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		p := make([]byte, 1024)
		fmt.Print("> ")
		n, err := os.Stdin.Read(p)
		if err != nil {
			fmt.Println(err)
			break
		}
		input := p[:n]

		_, err = c.Conn.Write(input)
		if err != nil {
			fmt.Println(err)
			break
		}

		if strings.TrimSpace(string(input)) == "\\quit" {
			break
		}

		raw := make([]byte, 1024)
		n, err = c.Conn.Read(raw)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(string(raw[:n]))
	}
}
