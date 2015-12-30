package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/lessos/lessgo/logger"
)

type User struct {
	Uid    uint64 `json:"uid"`
	Port   uint16 `json:"port"`
	Action string `json:"action"`
}

var (
	user          User
	userFormValue string
)

func startApiServer() {
	http.HandleFunc("/api/v1/do", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		userFormValue = req.FormValue("user")
		err = json.Unmarshal([]byte(userFormValue), &user)
		if err != nil {
			resp := fmt.Sprintf("parser user data error:%v\r\n", err)
			logger.Printf("warn", "[API]parser user data error:%v\r\n", err)
			w.Write([]byte(resp))
			return
		}

		if user.Uid == 0 || user.Port < 1024 || user.Action == "" {
			logger.Print("warn", "[api request]user data invalid.")
			w.Write([]byte("user data invalid.\r\n"))
			return
		}

	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ApiServerPort), nil); err != nil {
		logger.Printf("fatal", "API Server error: %v\r\n", err.Error())
		os.Exit(0)
	}
}
