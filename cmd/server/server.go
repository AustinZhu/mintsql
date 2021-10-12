package server

import (
	"context"
	"log"
	"mintsql/internal/backend"
	"net"
)

const (
	HOST     = "127.0.0.1"
	PROTOCOL = "tcp"
	PORT     = 7384
)

type Server struct {
	Addr   *net.TCPAddr
	Engine *backend.Engine
}

func New(host string, port uint16, poolSize int) (s *Server) {
	// TODO server network options
	s = &Server{}
	s.Addr = &net.TCPAddr{
		IP:   net.ParseIP(HOST),
		Port: PORT,
	}
	return s
}

func (s *Server) Run() {
	log.Printf("Welcome to mintsql Server.")
	l, err := net.ListenTCP(PROTOCOL, s.Addr)
	if err != nil {
		log.Fatal("error listening: ", err)
	}
	defer func(l *net.TCPListener) {
		err := l.Close()
		if err != nil {
			log.Println(err)
		}
	}(l)
	log.Printf("Listening on %s", s.Addr)
	for {
		conn, err := l.AcceptTCP()
		log.Printf("Accepted incoming connection on %s", conn.RemoteAddr())
		if err != nil {
			log.Fatalln(err)
		}
		go s.Handle(context.TODO(), conn)
	}
}

func (s *Server) Handle(ctx context.Context, conn *net.TCPConn) {
	panic("Not Implemented")
}
