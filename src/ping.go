package main

import (
	"log"
	"net/http"
)

func ping() {
	http.HandleFunc("/ping", PongServer)
	err := http.ListenAndServeTLS(config.Ping, config.PingCert, config.PingKey, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func PongServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("pong"))
}
