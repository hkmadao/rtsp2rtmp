package rlog

import (
	"log"
	"os"
	"time"

	"github.com/yumrano/rtsp2rtmp/conf"
)

var Log *log.Logger

func init() {
	// log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	logPath, err := conf.GetString("server.log.path")
	if err != nil {
		log.Println("get logPath error :", err)
	}
	lfile, err := os.OpenFile(logPath+"/"+time.Now().Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("open file error :", err)
	}
	Log = log.New(lfile, "", log.Llongfile|log.Lmicroseconds|log.Ldate)
	// Log = log.New(log.Default().Writer(), "", log.Llongfile|log.Lmicroseconds|log.Ldate)
	Log.Printf("test message : %s", "message")
}
