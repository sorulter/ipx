package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/lessos/lessgo/logger"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

type Proxy struct {
	listener net.Listener
}

type ProxydManager struct {
	sync.Mutex
	proxy map[uint64]*Proxy
}

var (
	proxyManager = ProxydManager{proxy: map[uint64]*Proxy{}}
)

func start(uid uint64, port uint16) (ok bool, err error) {
	// proxy
	proxy := NewHttpServer(uid)
	logger.Printf("info", "[ctrl]start listen port %d as proxy http server for user %d.", port, uid)

	proxy.Tr.Dial = func(network, addr string) (conn net.Conn, err error) {
		// new conn
		cipher, err := ss.NewCipher(config.ParentServer.Method, config.ParentServer.Key)
		if err != nil {
			return nil, err
		}
		conn, err = dial(addr, config.ParentServer.HostAndPort, cipher.Copy(), uid)
		if err != nil {
			return nil, err
		}

		return
	}

	proxy.Counter = counter
	proxy.FailFlowCounter = failFlowCounter

	// http server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return ok, errors.New(fmt.Sprintf("Listen new http server port(%d) error: %s", port, err.Error()))

	}
	go proxyManager.add(uid, listener, nil)
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
	pm.proxy[uid] = &Proxy{listener}
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
		prx.listener.Close()
		pm.Lock()
		delete(pm.proxy, uid)
		pm.Unlock()
	}
}

type Port struct {
	NodeName     string
	Port         uint16
	UserId       uint64
	Used         int
	Free         int
	ComboFlows   int
	ComboEndDate time.Time
	PortUpdateAt time.Time
}

func crtl() {
	getAndListenPorts()

	for {
		time.Sleep(15e9)
		getAndListenPorts()
	}
}

func getAndListenPorts() {
	var (
		ports []Port
	)

	db.Table("ports").Select(
		"node_name,`port`,ports.user_id,ports.updated_at as port_update_at,used,free,combo_flows,combo_end_date").Joins(
		"JOIN flows ON ports.user_id = flows.user_id").Where(
		" node_name = ?", config.NodeName,
	).Find(&ports)

	now := time.Now().In(loc)
	for _, port := range ports {
		// No any flows.
		if port.ComboFlows+port.Free == 0 {
			return
		}

		_, isRunning := proxyManager.get(port.UserId)

		fmt.Printf("[check]user %d (port %d) is running: %v\n", port.UserId, port.Port, isRunning)

		// Not running and have enough flows.
		if !isRunning && port.Used < port.ComboFlows+port.Free {
			// fmt.Printf("[Start] user %d, port %d, used: %d, flows: %d\n", port.UserId, port.Used, port.ComboFlows, port.Free)
			start(port.UserId, port.Port)
		}

		// Is running but have not enough flows.
		if isRunning && port.Used >= port.ComboFlows+port.Free {
			// fmt.Printf("Stop user %d, port %d\n", port.UserId, port.Port)
			stop(port.UserId)
		}

		// Is running but now is after combo end time.
		if isRunning && now.After(port.ComboEndDate) && port.Free <= port.Used {
			fmt.Printf("Combo is after the time,Stop user %d, port %d\n", port.UserId, port.Port)
			stop(port.UserId)
		}
	}
}
