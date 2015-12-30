package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	config Config
	err    error
)

func main() {

	// Exit main loop
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
forever:

	for {
		select {
		case s := <-sig:
			fmt.Printf("\nSignal (%d) received, stopping\n", s)
			break forever
		}
	}

}
