package main

import (
	"os"
	"os/signal"
	"syscall"

	// _ "net/http/pprof"

	"github.com/beego/beego/v2/core/logs"
	"github.com/yumrano/rtsp2rtmp/app"
	_ "github.com/yumrano/rtsp2rtmp/conf"
	"github.com/yumrano/rtsp2rtmp/server"
)

func main() {
	go app.ServeHTTP()
	server.NewServer()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	http.ListenAndServe("0.0.0.0:6060", nil)
	// }()
	logs.Info("Server Start Awaiting Signal")
	select {
	case sig := <-sigs:
		logs.Error(sig)
	case <-done:
	}
	logs.Error("Exiting")
}
