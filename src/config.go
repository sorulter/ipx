package main

type Config struct {
	ApiServerPort uint16 `json:"api_server_port"`
	ParentServer  struct {
		HostAndPort string `json:"host_and_port"`
		Method      string `json:"method"`
		Key         string `json:"key"`
	} `json:"parent_server"`
}
