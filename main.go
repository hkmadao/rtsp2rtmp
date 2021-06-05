package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/yumrano/rtsp2rtmp/conf"
	"github.com/yumrano/rtsp2rtmp/rlog"
	"github.com/yumrano/rtsp2rtmp/server"
	"github.com/yumrano/rtsp2rtmp/writer/httpflv"
)

func main() {
	go httpflv.ServeHTTP()
	s := server.NewServer()
	go s.ServeStreams()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	rlog.Log.Println("Server Start Awaiting Signal")
	select {
	case sig := <-sigs:
		rlog.Log.Println(sig)
	case <-done:
	}
	rlog.Log.Println("Exiting")
}
