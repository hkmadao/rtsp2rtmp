package main

import (
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/yumrano/rtsp2rtmp/conf"
	"github.com/yumrano/rtsp2rtmp/dao"
	"github.com/yumrano/rtsp2rtmp/rlog"
	"github.com/yumrano/rtsp2rtmp/writer/fileflv"
	"github.com/yumrano/rtsp2rtmp/writer/httpflv"
	"github.com/yumrano/rtsp2rtmp/writer/rtmpflv"
)

type Server struct {
	rms map[string]*RtspManager
}

func NewServer() *Server {
	return &Server{
		rms: make(map[string]*RtspManager),
	}
}

func (s *Server) serveStreams() {
	defer func() {
		if r := recover(); r != nil {
			rlog.Log.Printf("rtspManager pain %v", r)
		}
	}()
	es, err := dao.CameraSelectAll()
	if err != nil {
		rlog.Log.Printf("camera list is empty")
		return
	}
	for _, camera := range es {
		go func(c dao.Camera) {
			defer func() {
				if r := recover(); r != nil {
					rlog.Log.Printf("rtspManager pain %v", r)
				}
			}()
			for {
				rlog.Log.Println(c.Code, "connect", c.RtspURL)
				rtsp.DebugRtsp = true
				session, err := rtsp.Dial(c.RtspURL)
				if err != nil {
					rlog.Log.Println(c.Code, err)
					c.OnlineStatus = 0
					time.Sleep(5 * time.Second)
					if c.OnlineStatus == 1 {
						dao.CameraUpdate(c)
					}
					continue
				}
				session.RtpKeepAliveTimeout = 10 * time.Second
				if err != nil {
					rlog.Log.Println(c.Code, err)
					time.Sleep(5 * time.Second)
					continue
				}
				codec, err := session.Streams()
				if err != nil {
					rlog.Log.Println(c.Code, err)
					time.Sleep(5 * time.Second)
					continue
				}

				pRtmpFlvChan := make(chan av.Packet)
				pFlvFileChan := make(chan av.Packet)
				pHttpFlvChan := make(chan av.Packet)
				done := make(chan interface{})
				rm := &RtspManager{
					ffm:  fileflv.NewFileFlvManager(),
					hfm:  httpflv.NewHttpFlvManager(),
					frm:  rtmpflv.NewRtmpFlvManager(),
					done: done,
				}
				rm.pktTransfer(c.Code, codec, pRtmpFlvChan, pFlvFileChan, pHttpFlvChan)
				s.rms[c.Code] = rm

				c.OnlineStatus = 1
				dao.CameraUpdate(c)

				for {
					pkt, err := session.ReadPacket()
					if err != nil {
						rlog.Log.Println(c.Code, err)
						break
					}
					writeChan(pkt, pRtmpFlvChan, done)
					writeChan(pkt, pFlvFileChan, done)
					writeChan(pkt, pHttpFlvChan, done)
				}
				err = session.Close()
				if err != nil {
					rlog.Log.Println("session Close error", err)
				}
				rlog.Log.Println(c.Code, "reconnect wait 5s")
				time.Sleep(5 * time.Second)
			}
		}(camera)
	}
}

func writeChan(pkt av.Packet, c chan<- av.Packet, done <-chan interface{}) {
	select {
	case c <- pkt:
	case <-time.After(1 * time.Nanosecond):
	case <-done:
	}
}

type RtspManager struct {
	ffm  *fileflv.FileFlvManager
	hfm  *httpflv.HttpFlvManager
	frm  *rtmpflv.RtmpFlvManager
	done <-chan interface{}
}

func (r *RtspManager) pktTransfer(code string, codecs []av.CodecData, pRtmpFlvChan <-chan av.Packet, pFlvFileChan <-chan av.Packet, pHttpFlvChan <-chan av.Packet) {
	go r.frm.FlvWrite(code, codecs, r.done, pRtmpFlvChan)
	go r.hfm.FlvWrite(code, codecs, r.done, pHttpFlvChan)
	save, err := conf.GetBool("server.fileflv.save")
	if err == nil && save {
		go r.ffm.FlvWrite(code, codecs, r.done, pFlvFileChan)
	}
	rlog.Log.Printf("get server.fileflv.save error : %v", err)
}
