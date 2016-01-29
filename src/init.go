package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	gconf "github.com/gocubes/config"
	"github.com/jinzhu/gorm"
	"github.com/lessos/lessgo/data/hissdb"
)

var (
	db   gorm.DB
	ssdb *hissdb.Connector
	pre  string
	v    bool
	Git  string
	At   string
)

func init() {
	flag.StringVar(&pre, "prefix", ".", "config file prefix path")
	flag.BoolVar(&v, "v", false, "show version info")
	flag.Parse()

	if v {
		fmt.Printf("Version: %s,\nBuild At: %s UTC.\n", Git, At)
		os.Exit(0)
	}

	provier, perr := gconf.New(pre+"/etc/config.json", "json")
	if perr != nil {
		log.Fatalf("[init]create config provider error: %v\n", perr.Error())
	}

	gerr := provier.Get(&config)
	if gerr != nil {
		log.Fatalf("[init]get config data error: %v\n", gerr.Error())
	}
	provier.ReloadOn(&config, syscall.SIGUSR1)

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
	// db = *db.Debug()
	db.LogMode(true)
	db.DB().SetMaxIdleConns(50)
	db.DB().SetMaxOpenConns(200)

}
