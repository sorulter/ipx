package main

import (
	"errors"
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
	proxyManager = ProxydManager{proxy: map[uint64]*Proxy{}}
)

func newHttpServer(uid uint64, port uint16) (ok bool, err error) {
	// proxy
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.Tr.Dial = func(network, addr string) (conn net.Conn, err error) {
		// use existed conn
		prx, ok := proxyManager.get(uid)
		if !ok {
			return prx.conn, nil
		}

		// new conn
		cipher, err := ss.NewCipher(config.ParentServer.Method, config.ParentServer.Key)
		if err != nil {
			proxy.Logger.Fatal("create shadowsocks cipher:", err)
			return nil, err
		}
		conn, err = dial(addr, config.ParentServer.HostAndPort, cipher.Copy(), uid)
		if err != nil {
			proxy.Logger.Fatalf("can't connect to shadowsocks parent %s for %s: %v\n",
				config.ParentServer.HostAndPort, addr, err)
			return nil, err
		}
		fmt.Println("Dial to parent", conn)
		proxyManager.updateConn(uid, conn)

		return
	}

	// http server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return ok, errors.New(fmt.Sprintf("Listen new http server port(%d) error: %s", port, err.Error()))

	}
	proxyManager.add(uid, listener, nil)
	server := http.Server{
		Handler: proxy,
	}
	go server.Serve(listener)
	return true, err

}

func stop(uid uint64) {
	proxyManager.del(uid)
}

func (pm *ProxydManager) add(uid uint64, listener net.Listener, conn net.Conn) {
	pm.Lock()
	pm.proxy[uid] = &Proxy{conn, listener}
	pm.Unlock()
}

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
