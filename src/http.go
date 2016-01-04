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
