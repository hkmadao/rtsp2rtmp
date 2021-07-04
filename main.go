package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/conf" // 必须先导入配置文件
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtspclientmanager"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/task"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web"

	"github.com/beego/beego/v2/core/logs"
)

func main() {
	rtspclientmanager.GetSingleRtspClientManager().StartClient()
	task.GetSingleTask().StartTask()
	web.GetSingleWeb().StartWeb()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	logs.Info("Server Start Awaiting Signal")
	select {
	case sig := <-sigs:
		logs.Info(sig)
	case <-done:
	}
	logs.Info("Exiting")
}
