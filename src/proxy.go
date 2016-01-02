package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/kofj/goproxy"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

type Proxy struct {
	conn     net.Conn
	listener net.Listener
}

type ProxydManager struct {
	sync.Mutex
	proxy map[uint64]*Proxy
}

var (
)

func (pm *ProxydManager) add(uid uint64, listener net.Listener, conn net.Conn) {
	pm.Lock()
	pm.proxy[uid] = &Proxy{conn, listener}
	pm.Unlock()
}

func newHttpServer() {
	proxy.Logger.Fatal(http.ListenAndServe(":10010", proxy))
func (pm *ProxydManager) updateConn(uid uint64, conn net.Conn) {
	pm.Lock()
	if pm.proxy[uid].conn != nil {
		pm.proxy[uid].conn.Close()
	}
	pm.proxy[uid].conn = conn
	pm.Unlock()

}

func (pm *ProxydManager) get(uid uint64) (prx *Proxy, ok bool) {
	pm.Lock()
	prx, ok = pm.proxy[uid]
	pm.Unlock()
	return
}

func (pm *ProxydManager) del(uid uint64) {
	if prx, ok := pm.get(uid); ok {
		if prx.conn != nil {
			prx.conn.Close()
		}
		prx.listener.Close()
		pm.Lock()
		delete(pm.proxy, uid)
		pm.Unlock()
	}
}
