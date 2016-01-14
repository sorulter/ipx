package main

import (
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
	gconf "github.com/gocubes/config"
	"github.com/jinzhu/gorm"
	"github.com/lessos/lessgo/data/hissdb"
)

var (
	db   gorm.DB
	ssdb *hissdb.Connector
)

func init() {
	flag.Parse()

	provier, perr := gconf.New("etc/config.json", "json")
	if perr != nil {
		log.Fatalf("[init]create config provider error: %v\n", perr.Error())
	}

	gerr := provier.Get(&config)
	if gerr != nil {
		log.Fatalf("[init]get config data error: %v\n", gerr.Error())
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
		log.Fatalf("[init db]Connect to ssdb error: %s", err.Error())
	}

	db, err = gorm.Open("mysql", config.FlowCounter.DB.DSN)
	if err != nil {
		log.Fatalf("[init db]MySQL Connect error: %v\n", err.Error())
	}
}
