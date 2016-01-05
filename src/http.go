package main

import (
	"log"
	"net"
	"net/http"
	"regexp"
)

type HttpServer struct {
	Uid             uint64
	Logger          *log.Logger
	NotFoundHandler http.Handler
	Tr              *http.Transport
	ConnectDial     func(network string, addr string) (net.Conn, error)
}

var (
	hasPort = regexp.MustCompile(`:\d+$`)
)

func (h *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func NewHttpServer(uid uint64) *HttpServer {
	return &HttpServer{
		Uid: uid,
		NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "Illegal requests.", 500)
		}),
		Tr: &http.Transport{
			Proxy: http.ProxyFromEnvironment},
	}
}
}
