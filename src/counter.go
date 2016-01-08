package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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

// Upload local flow data to remote mysql server.
func upload() {
	now := time.Now().In(loc)
	before := now.Add(-6e10)

	uplimit := now.Format("ipx.200601021504")
	downlimit := before.Format("ipx.200601021504")

	if rs := ssdb.Cmd("scan", downlimit, uplimit, 9999999); rs.State == hissdb.ReplyOK {
		for _, v := range rs.Hash() {
			db.Table("flows").Where("user_id = ?", v.Key[17:]).Update("used", gorm.Expr("used + ?", v.Value))
		}
	} else {
		logger.Printf("warn", "[ssdb]cmd error: scan %s %s 9999999", downlimit, uplimit)
	}

}
