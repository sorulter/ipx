package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/kofj/goproxy"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

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
		c, err = ss.Dial(addr, config.ParentServer.HostAndPort, cipher.Copy())
		if err != nil {
			proxy.Logger.Fatalf("can't connect to shadowsocks parent %s for %s: %v\n",
				config.ParentServer.HostAndPort, addr, err)
			return nil, err
		}
		return
	}

}

func newHttpServer() {
	proxy.Logger.Fatal(http.ListenAndServe(":10010", proxy))
}
