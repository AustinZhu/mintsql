package main

import (
	"context"
	"github.com/google/uuid"
	"log"
	"mintsql/internal/backend"
	"mintsql/internal/store/table"
	"net"
	"strconv"
)

const (
	HOST     = "0.0.0.0"
	PROTOCOL = "tcp"
	PORT     = "7384"
)

var (
	sqlServer *Server
	port      string
)

type Server struct {
	Addr   *net.TCPAddr
	Engine *backend.Engine
}

func New(port string) (s *Server) {
	port_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalln("not a valid port number", err)
	}

	s = &Server{
		Engine: backend.Setup(),
		Addr: &net.TCPAddr{
			IP:   net.ParseIP(HOST),
			Port: port_,
		},
	}
	return s
}

func (s *Server) Run() {
	l, err := net.ListenTCP(PROTOCOL, s.Addr)
	if err != nil {
		log.Fatal("error listening: ", err)
	}

	defer func(l *net.TCPListener) {
		err := l.Close()
		if err != nil {
			return
		}
	}(l)

	log.Printf("Listening on %s", s.Addr)
	for {
		conn, err := l.AcceptTCP()
		log.Printf("Accepted incoming connection on %s", conn.RemoteAddr())
		if err != nil {
			log.Println(err)
			continue
		}

		go s.HandleRepl(context.Background(), conn)
	}
}

func (s *Server) HandleRepl(ctx context.Context, conn *net.TCPConn) {
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	for {
		var n int
		raw := make([]byte, 1024)
		n, err := conn.Read(raw)
		if err != nil {
			log.Println(err)
			return
		}

		ctx = context.WithValue(context.Background(), "uuid", uuid.New())
		ctx = context.WithValue(ctx, "addr", conn.RemoteAddr())

		var res *table.Result
		var resp string
		query := string(raw[:n])

		res, err = s.Engine.Execute(ctx, query)
		if err != nil {
			resp = err.Error()
		} else if res == nil {
			resp = "ok"
		} else {
			resp = res.String()
		}

		_, err = conn.Write([]byte(resp))
		if err != nil {
			log.Println(err)
			return
		}
	}
}
