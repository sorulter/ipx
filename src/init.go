package main

import (
	"flag"
	"os"

	gconf "github.com/gocubes/config"
	"github.com/lessos/lessgo/data/hissdb"
	"github.com/lessos/lessgo/logger"
)

var (
	ssdb *hissdb.Connector
)

func init() {
	flag.Parse()

	provier, perr := gconf.New("etc/config.json", "json")
	if perr != nil {
		logger.Printf("fatal", "[init]create config provider error: %v\n", perr.Error())
		os.Exit(0)
	}

	gerr := provier.Get(&config)
	if gerr != nil {
		logger.Printf("fatal", "[init]get config data error: %v\n", gerr.Error())
		os.Exit(0)
	}

	initDB()

}

func initDB() {
	ssdb, err = hissdb.NewConnector(hissdb.Config{
		Host:    config.FlowCounter.SSDB.Host,
		Port:    config.FlowCounter.SSDB.Port,
		Auth:    config.FlowCounter.SSDB.Auth,
		MaxConn: config.FlowCounter.SSDB.MaxConn,
	})
	if err != nil {
		logger.Printf("fatal", "init ssdb error: %s", err.Error())
	}

}
