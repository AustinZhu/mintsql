package server

import (
	"bufio"
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

func New(host string, port uint16) (s *Server) {
	s = &Server{
		Engine: new(backend.Engine),
	}
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
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	str, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err)
	}
	res, err := s.Engine.Execute(ctx, str)
	resp := res.String()
	_, err = conn.Write([]byte(resp))
	if err != nil {
		log.Fatalln(err)
	}
}
