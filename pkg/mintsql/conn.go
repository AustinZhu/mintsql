package mintsql

import (
	"errors"
	"fmt"
	"net"
)

type Conn struct {
	net.Conn
	pool *ConnPool
}

func (c *Conn) Close() error {
	select {
	case c.pool.conn <- c:
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
	poolConn := &Conn{
		pool: cp,
	}
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

func NewConnPool(capacity int, addr *net.TCPAddr) (*ConnPool, error) {
	if capacity <= 0 {
		return nil, errors.New("invalid capacity")
	}
	cp := &ConnPool{conn: make(chan net.Conn, capacity), addr: addr}

	for i := 0; i < capacity; i++ {
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			return nil, fmt.Errorf("failed to dial %s", addr)
		}
		cp.conn <- conn
	}
	return cp, nil
}
