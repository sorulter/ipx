package main

type Config struct {
	ApiServerPort uint16 `json:"api_server_port"`
	ParentServer  struct {
		HostAndPort string `json:"host_and_port"`
		Method      string `json:"method"`
		Key         string `json:"key"`
	} `json:"parent_server"`
	FlowCounter struct {
		DB struct {
			Host    string `json:"host"`
			Port    uint16 `json:"port"`
			Auth    string `json:"auth"`
			MaxConn int    `json:"max_connect"`
		} `json:"db"`
	} `json:"flow_counter"`
}
