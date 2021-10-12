package mintsql

import (
	"errors"
	"fmt"
	"net"
)

var DefaultConnPool *ConnPool

type Conn struct {
	net.Conn
}

func (c *Conn) Close() error {
	if DefaultConnPool == nil {
		return errors.New("pool not configured")
	}
	select {
	case DefaultConnPool.conn <- c:
		return nil
	default:
		return c.Conn.Close()
	}
}

type ConnPool struct {
	conn chan net.Conn
	addr *net.TCPAddr
}

func (cp *ConnPool) Get() (*Conn, error) {
	if cp.conn == nil {
		return nil, errors.New("pool closed")
	}
	poolConn := new(Conn)
	select {
	case c, ok := <-cp.conn:
		if !ok {
			return nil, errors.New("failed to get connection")
		}
		poolConn.Conn = c
		return poolConn, nil
	default:
		conn, err := net.DialTCP("tcp", nil, cp.addr)
		if err != nil {
			return nil, fmt.Errorf("failed to dial %s", cp.addr)
		}
		cp.conn <- conn
		poolConn.Conn = conn
		return poolConn, nil
	}
}

func (cp *ConnPool) Close() {
	cp.conn = nil
}

func (cp *ConnPool) Len() int {
	return len(cp.conn)
}

func NewConnPool(capacity int, addr *net.TCPAddr) error {
	if capacity <= 0 {
		return errors.New("invalid capacity")
	}
	DefaultConnPool = &ConnPool{conn: make(chan net.Conn, capacity), addr: addr}

	for i := 0; i < capacity; i++ {
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			return fmt.Errorf("failed to dial %s", addr)
		}
		DefaultConnPool.conn <- conn
	}
	return nil
}
