package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	config Config
	err    error
)

func main() {
	// Upload flow data.
	go func() {
		for {
			if time.Now().In(loc).Second() == 1 {
				upload()
			}
			time.Sleep(1e9)
		}
	}()

	go crtl()

	// Exit main loop
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
forever:

	for {
		select {
		case s := <-sig:
			fmt.Printf("\nSignal (%d) received, stopping\n", s)
			for uid, _ := range proxyManager.proxy {
				fmt.Printf("Stop connection of user %d\n", uid)
				proxyManager.del(uid)

			}
			break forever
		}
	}

}
