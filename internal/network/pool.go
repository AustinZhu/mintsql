package network

import (
	"net"
	"sync"
)

type ConnPool struct {
	mutex sync.RWMutex
	conn  chan net.Conn
}

func (*ConnPool) Get() (net.Conn, error) {
	panic("Not Implemented")
}

func (*ConnPool) Close() {
	panic("Not Implemented")
}

func (*ConnPool) Len() {
	panic("Not Implemented")
}
