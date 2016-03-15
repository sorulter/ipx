package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lessos/lessgo/data/hissdb"
	"github.com/lessos/lessgo/logger"
)

var (
	loc, _ = time.LoadLocation("Asia/Shanghai")
)

func counter(uid uint64, bytes int64, ipstr string) {
	key := fmt.Sprintf("ipx.%s.%d.%d", time.Now().In(loc).Format("200601021504"), uid, ip2long(ipstr))
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
			id, ip := getIdIp(v.Key)
			if dberr := db.Table("flows").Where("user_id = ?", id).Updates(map[string]interface{}{
				"used":       gorm.Expr("used + ? * ?", v.Value, config.Multiple),
				"updated_at": time.Now(),
			}).Error; dberr != nil {
				logger.Printf("warn", "[upload]update mysql error: %v", dberr.Error())
			}

			// Logs
			used, _ := strconv.ParseFloat(v.Value, 0)
			db.Exec("INSERT INTO `logs_"+logshash(id)+"` (`user_id`, `flows`, `node`, `client_ip`, `used_at`) VALUES (?, ?, ?, ?, ?)", id, float32(used)*config.Multiple, config.NodeName, ip, now.Format("2006/01/02 15:04:05"))

		}
	} else {
		logger.Printf("warn", "[ssdb]cmd error: scan %s %s 9999999", downlimit, uplimit)
	}

}

func failFlowCounter(target string, bytes int64) {
	key := fmt.Sprintf("ipx.fail.flow.%s", time.Now().In(loc).Format("2006010215"))
	if rs := ssdb.Cmd("incr", key, fmt.Sprint(bytes)); rs.State != hissdb.ReplyOK {
		logger.Printf("warn", "Log fial flow error: %v, incr %s %d\n", rs.State, key, bytes)
	}
}

func getIdIp(key string) (id uint64, ip string) {
	slice := strings.Split(key, ".")
	id, _ = strconv.ParseUint(slice[2], 0, 32)
	ip = slice[3]
	return
}
