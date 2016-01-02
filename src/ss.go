package main

import (
	"net"
	"time"

	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

type Conn struct {
	net.Conn
	Uid uint64
}

func dial(addr, server string, cipher *ss.Cipher, uid uint64) (c *Conn, err error) {
	var conn *ss.Conn
	conn, err = ss.Dial(addr, server, cipher)
	c = newConn(conn, uid)
	return
}

func newConn(c net.Conn, uid uint64) *Conn {
	return &Conn{
		Conn: c,
		Uid:  uid,
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	return
}
