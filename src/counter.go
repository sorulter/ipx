package main

import (
	"fmt"
	"time"

	"github.com/lessos/lessgo/data/hissdb"
	"github.com/lessos/lessgo/logger"
)

var (
	loc, _ = time.LoadLocation("Asia/Shanghai")
)

func counter(uid uint64, bytes int64) {
	key := fmt.Sprintf("ipx.%s.%d", time.Now().In(loc).Format("200601021504"), uid)
	if rs := ssdb.Cmd("incr", key, fmt.Sprint(bytes)); rs.State != hissdb.ReplyOK {
		logger.Printf("warn", "Log flow error: %v, incr %s %d\n", rs.State, key, bytes)
	}
}
