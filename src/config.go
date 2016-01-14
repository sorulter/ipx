package main

type Config struct {
	NodeName     string `json:"node_name"`
	ParentServer struct {
		HostAndPort string `json:"host_and_port"`
		Method      string `json:"method"`
		Key         string `json:"key"`
	} `json:"parent_server"`
	FlowCounter struct {
		SSDB struct {
			Host    string `json:"host"`
			Port    uint16 `json:"port"`
			Auth    string `json:"auth"`
			MaxConn int    `json:"max_connect"`
		} `json:"ssdb"`
		DB struct {
			DSN string `json:"dsn"`
		} `json:"db"`
	} `json:"flow_counter"`
}
