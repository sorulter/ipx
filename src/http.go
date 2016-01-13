package main

import (
	"io"
	"net"
	"net/http"
	"regexp"

	"github.com/lessos/lessgo/logger"
)

type HttpServer struct {
	Uid             uint64
	NotFoundHandler http.Handler
	Tr              *http.Transport
	ConnectDial     func(network string, addr string) (net.Conn, error)
	Counter         func(uid uint64, bytes int64)
	FailFlowCounter func(target string, bytes int64)
}

var (
	hasPort = regexp.MustCompile(`:\d+$`)
)

func (h *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		h.handleHttps(w, r)
	} else {
		h.handleHttp(w, r)
	}
}

func NewHttpServer(uid uint64) *HttpServer {
	return &HttpServer{
		Uid: uid,
		NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "Illegal requests.", 500)
		}),
		Tr:      &http.Transport{Proxy: http.ProxyFromEnvironment},
		Counter: func(uid uint64, bytes int64) {},

		FailFlowCounter: func(target string, bytes int64) {},
	}
}

func httpError(w io.WriteCloser, err error) {
	if _, err2 := io.WriteString(w, "HTTP/1.1 502 Bad Gateway\r\n\r\n"); err2 != nil {
		logger.Printf("warn", "[proxy]Error responding to client: %s", err2.Error())
	}
	if err3 := w.Close(); err3 != nil {
		logger.Printf("warn", "[proxy]Error closing client connection: %s", err3.Error())
	}
	logger.Printf("warn", "[proxy]Error dail to remote server,response 502 and closed request: %v", err.Error())
}

func (h *HttpServer) handleHttps(w http.ResponseWriter, r *http.Request) {
	hij, ok := w.(http.Hijacker)
	if !ok {
		logger.Print("warn", "[proxy]httpserver does not support hijacking")
		return
	}

	proxyClient, _, e := hij.Hijack()
	if e != nil {
		logger.Printf("warn", "[proxy]Cannot hijack connection %v", e.Error())
		return
	}

	host := r.URL.Host
	if !hasPort.MatchString(host) {
		host += ":80"
	}
	targetSiteCon, err := h.connectDial("tcp", host)
	if err != nil {
		httpError(proxyClient, err)
		return
	}
	// log.Printf("Accepting CONNECT to %s", host)
	proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	go h.copyAndClose(targetSiteCon, proxyClient, r, "to server")
	go h.copyAndClose(proxyClient, targetSiteCon, r, "to client")

}

func (h *HttpServer) connectDial(network, addr string) (c net.Conn, err error) {
	if h.ConnectDial == nil {
		return h.dial(network, addr)
	}
	return h.ConnectDial(network, addr)
}

func (h *HttpServer) dial(network, addr string) (c net.Conn, err error) {
	if h.Tr.Dial != nil {
		return h.Tr.Dial(network, addr)
	}
	return net.Dial(network, addr)
}

func (h *HttpServer) copyAndClose(w, r net.Conn, req *http.Request, do string) {
	connOk := true
	n, err := io.Copy(w, r)
	if err != nil {
		connOk = false
		logger.Printf("warn", "[proxy]Error %s: %s, %d bytes,host is: %v", do, err.Error(), n, req.URL.Host)
		go h.FailFlowCounter(req.URL.Host, n)
	} else {
		go h.Counter(h.Uid, n)
	}

	if err := r.Close(); err != nil && connOk {
		logger.Printf("warn", "[proxy]Error closing client connection: %s", err.Error())
	}
}

func (h *HttpServer) handleHttp(w http.ResponseWriter, r *http.Request) {
	if !r.URL.IsAbs() {
		h.NotFoundHandler.ServeHTTP(w, r)
		return
	}

	removeProxyHeaders(r)

	resp, err := h.Tr.RoundTrip(r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	nh := copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	nr, _ := io.Copy(w, resp.Body)
	if err := resp.Body.Close(); err != nil {
		logger.Printf("warn", "[proxy]close response body error:%s", err.Error())
	}
	go h.Counter(h.Uid, nr+nh)
}

func removeProxyHeaders(r *http.Request) {
	r.RequestURI = "" // this must be reset when serving a request with the client
	// If no Accept-Encoding header exists, Transport will add the headers it can accept
	// and would wrap the response body with the relevant reader.
	r.Header.Del("Accept-Encoding")
	// curl can add that, see
	// http://homepage.ntlworld.com/jonathan.deboynepollard/FGA/web-proxy-connection-header.html
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Proxy-Authenticate")
	r.Header.Del("Proxy-Authorization")
	// Connection, Authenticate and Authorization are single hop Header:
	// http://www.w3.org/Protocols/rfc2616/rfc2616.txt
	// 14.10 Connection
	//   The Connection general-header field allows the sender to specify
	//   options that are desired for that particular connection and MUST NOT
	//   be communicated by proxies over further connections.
	r.Header.Del("Connection")
}

func copyHeaders(dst, src http.Header) (n int64) {
	for k, _ := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			n += int64(len(k) + len(v))
			dst.Add(k, v)
		}
	}
	return
}
