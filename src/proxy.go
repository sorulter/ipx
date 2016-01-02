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
	proxy *goproxy.ProxyHttpServer
)

func init() {
	proxy = goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	// set transport
	fmt.Println("[main.info]Set Tr.Dial")

	proxy.Tr.Dial = func(network, addr string) (c net.Conn, err error) {
		fmt.Println("Dial to parent")
		cipher, err := ss.NewCipher(config.ParentServer.Method, config.ParentServer.Key)
		if err != nil {
			proxy.Logger.Fatal("create shadowsocks cipher:", err)
			return nil, err
		}
		c, err = dial(addr, config.ParentServer.HostAndPort, cipher.Copy(), 121)
		if err != nil {
			proxy.Logger.Fatalf("can't connect to shadowsocks parent %s for %s: %v\n",
				config.ParentServer.HostAndPort, addr, err)
			return nil, err
		}
		return
	}

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
