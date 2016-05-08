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

		ReadTimeout: time.Duration(config.KeepAliveTimeout) * time.Second,
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

func (pm *ProxydManager) clean(uids []uint64) {
	pm.Lock()
	defer pm.Unlock()
	for id := range pm.proxy {
		var hasid bool
		for _, uid := range uids {
			if uid == id {
				hasid = true
			}
		}
		if !hasid {
			pm.del(id)
		}
	}

}

type Port struct {
	NodeName     string
	Port         uint16
	UserId       uint64
	Activate     int
	Used         int
	Forever      int
	Combo        int
	Extra        int
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
		ids   []uint64
	)

	db.Table("ports").Select(
		"node_name,`port`,ports.user_id,ports.updated_at as port_update_at,activate,used,forever,combo,extra,combo_end_date").Joins(
		"JOIN flows ON ports.user_id = flows.user_id JOIN users ON ports.user_id = users.id").Where(
		" node_name = ?", config.NodeName,
	).Find(&ports)

	now := time.Now().In(loc)
	for _, port := range ports {
		ids = append(ids, port.UserId)

		// No any flows.
		if port.Combo+port.Forever+port.Extra == 0 {
			continue
		}

		_, isRunning := proxyManager.get(port.UserId)

		// Not running and have enough flows.
		if !isRunning && !now.After(port.ComboEndDate) && port.Combo+port.Forever+port.Extra > port.Used && port.Activate == 1 {
			// fmt.Printf("[Start] user %d, port %d, used: %d, flows: %d\n", port.UserId, port.Used, port.Combo, port.Forever)
			start(port.UserId, port.Port)
		}

		// After combo flows end time, but have enough forever flows.
		if !isRunning && now.After(port.ComboEndDate) && port.Forever+port.Extra > port.Used && port.Activate == 1 {
			start(port.UserId, port.Port)
		}

		// Is running but have not enough flows.
		if isRunning && port.Used >= port.Combo+port.Forever+port.Extra {
			// fmt.Printf("Stop user %d, port %d\n", port.UserId, port.Port)
			stop(port.UserId)
		}

		// Is running but now is after combo end time.
		if isRunning && now.After(port.ComboEndDate) && port.Forever+port.Extra <= port.Used {
			fmt.Printf("Combo is after the time,Stop user %d, port %d\n", port.UserId, port.Port)
			stop(port.UserId)
		}
	}
	proxyManager.clean(ids)
}
