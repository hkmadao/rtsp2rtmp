package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/yumrano/rtsp2rtmp/conf"
	"github.com/yumrano/rtsp2rtmp/rlog"
	"github.com/yumrano/rtsp2rtmp/writer/httpflv"
)

func main() {
	go httpflv.ServeHTTP()
	s := NewServer()
	go s.serveStreams()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		rlog.Log.Println(sig)
		done <- true
	}()
	rlog.Log.Println("Server Start Awaiting Signal")
	<-done
	rlog.Log.Println("Exiting")
}
