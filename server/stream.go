package server

import (
	"runtime/debug"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/yumrano/rtsp2rtmp/controllers"
	"github.com/yumrano/rtsp2rtmp/models"
	"github.com/yumrano/rtsp2rtmp/writer/fileflv"
	"github.com/yumrano/rtsp2rtmp/writer/httpflv"
	"github.com/yumrano/rtsp2rtmp/writer/rtmpflv"
)

var rms sync.Map

type Server struct {
	codeStream <-chan string //管理员关闭

}

func NewServer() *Server {
	codeStream := controllers.CodeStream()
	s := &Server{
		codeStream: codeStream,
	}
	go s.serveStreams()
	go s.stopConn()
	return s
}

func ExistCamera(code string) bool {
	_, b := rms.Load(code)
	return b
}

func (s *Server) stopConn() {
	codeStream := controllers.CodeStream()
	for {
		code := <-codeStream
		v, b := rms.Load(code)
		if b {
			r := v.(*RtspManager)
			err := r.conn.Close()
			if err != nil {
				logs.Error("camera [%s] close error : %v", code, err)
				return
			}
			logs.Info("camera [%s] close success", code)
		}
	}

}

func (s *Server) serveStreams() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("rtspManager panic %v", r)
		}
	}()
	for {
		es, err := models.CameraSelectAll()
		if err != nil {
			logs.Error("camera list is empty")
			return
		}
		for _, camera := range es {
			if v, b := rms.Load(camera.Code); b && v != nil {
				continue
			}
			if camera.Enabled != 1 {
				continue
			}
			go s.connRtsp(camera.Code)
		}
		<-time.After(30 * time.Second)
	}

}

func (s *Server) connRtsp(code string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		q := models.Camera{Code: code}
		c, err := models.CameraSelectOne(q)
		if err != nil {
			logs.Error("find camera [%s] error : %v", code, err)
			return
		}
		if c.Enabled != 1 {
			logs.Error("camera [%s] disabled : %v", code)
			return
		}
		logs.Info(c.Code, "connect", c.RtspURL)
		rtsp.DebugRtsp = true
		session, err := rtsp.Dial(c.RtspURL)
		if err != nil {
			logs.Error("camera [%s] conn : %v", c.Code, err)
			c.OnlineStatus = 0
			time.Sleep(5 * time.Second)
			if c.OnlineStatus == 1 {
				models.CameraUpdate(c)
			}
			return
		}
		session.RtpKeepAliveTimeout = 10 * time.Second
		codec, err := session.Streams()
		if err != nil {
			logs.Error("camera [%s] get streams : %v", c.Code, err)
			time.Sleep(5 * time.Second)
			return
		}

		rm := NewRtmpManager(session, c.Code, codec)
		rms.Store(c.Code, rm)

		c.OnlineStatus = 1
		models.CameraUpdate(c)
		for {
			pkt, err := session.ReadPacket()
			if err != nil {
				logs.Error("camera [%s] ReadPacket : %v", c.Code, err)
				break
			}
			//不能开goroutine,不能保证包的顺序
			writeChan(pkt, rm.rfPktStream, rm.rfPktDone, "rfp")
			writeChan(pkt, rm.ffPktStream, rm.ffPktDone, "ffp")
			writeChan(pkt, rm.hfPktStream, rm.hfPktDone, "hfp")
		}
		close(rm.rfPktDone)
		close(rm.ffPktDone)
		close(rm.hfPktDone)
		err = session.Close()
		if err != nil {
			logs.Error("session Close error : %v", err)
		}
		//offline camera
		q = models.Camera{Code: code}
		c, err = models.CameraSelectOne(q)
		if err != nil {
			logs.Error("find camera [%s] error : %v", code, err)
			return
		}
		c.OnlineStatus = 0
		models.CameraUpdate(c)

		rms.Delete(c.Code)
		logs.Info("camera [%s] reconnect wait 5s", c.Code)
		time.Sleep(5 * time.Second)
	}
}

func writeChan(pkt av.Packet, c chan<- av.Packet, done <-chan interface{}, t string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("writeChan panic : %v", r)
		}
	}()
	select {
	case c <- pkt:
	case <-time.After(1 * time.Millisecond):
		// if t == "hfp" {
		// 	logs.Info("lose pkt %s", t)
		// }
	case <-done:
	}
}

type RtspManager struct {
	conn        *rtsp.Client
	code        string
	codecs      []av.CodecData
	rfPktDone   chan interface{}
	ffPktDone   chan interface{}
	hfPktDone   chan interface{}
	rfPktStream chan av.Packet
	ffPktStream chan av.Packet
	hfPktStream chan av.Packet
}

func NewRtmpManager(conn *rtsp.Client, code string, codecs []av.CodecData) *RtspManager {
	rfPktDone := make(chan interface{})
	ffPktDone := make(chan interface{})
	hfPktDone := make(chan interface{})
	rfPktStream := make(chan av.Packet, 10)
	ffPktStream := make(chan av.Packet, 10)
	hfPktStream := make(chan av.Packet, 10)
	rm := &RtspManager{
		conn:        conn,
		code:        code,
		codecs:      codecs,
		rfPktDone:   rfPktDone,
		ffPktDone:   ffPktDone,
		hfPktDone:   hfPktDone,
		rfPktStream: rfPktStream,
		ffPktStream: ffPktStream,
		hfPktStream: hfPktStream,
	}
	go rm.pktTransfer()
	return rm
}

func (rm *RtspManager) pktTransfer() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	rtmpflv.NewRtmpFlvManager(rm.rfPktDone, rm.rfPktStream, rm.code, rm.codecs)
	httpflv.NewHttpFlvManager(rm.hfPktDone, rm.hfPktStream, rm.code, rm.codecs)
	save, err := config.Bool("server.fileflv.save")
	if err != nil {
		logs.Error("get server.fileflv.save error : %v", err)
		return
	}
	if save {
		fileflv.NewFileFlvManager(rm.ffPktDone, rm.ffPktStream, rm.code, rm.codecs)
	}
}
