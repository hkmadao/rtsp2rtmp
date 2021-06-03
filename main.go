package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/yumrano/rtsp2rtmp/conf"
	"github.com/yumrano/rtsp2rtmp/writer/httpflv"
)

func main() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	// go serveHTTP()
	go httpflv.ServeHTTP()
	s := NewServer()
	go s.serveStreams()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println(sig)
		done <- true
	}()
	log.Println("Server Start Awaiting Signal")
	<-done
	log.Println("Exiting")
}
