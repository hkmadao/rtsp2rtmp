package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/conf" // 必须先导入配置文件
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/ffmpegmanager"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtmpserver"

	// "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtspclientmanager"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web"
	_ "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/register"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/task"

	// "net/http"
	// _ "net/http/pprof"

	"github.com/beego/beego/v2/core/logs"
)

func main() {
	// rtspclientmanager.GetSingleRtspClientManager().StartClient()
	rtmpserver.GetSingleRtmpServer().StartRtmpServer()
	ffmpegmanager.GetSingleFFmpegManager().StartClient()
	task.GetSingleTask().StartTask()
	web.GetSingleWeb().StartWeb()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	logs.Info("Server Start Awaiting Signal")
	// http.ListenAndServe("0.0.0.0:6060", nil)
	select {
	case sig := <-sigs:
		logs.Info(sig)
	case <-done:
	}
	logs.Info("Exiting")
}
