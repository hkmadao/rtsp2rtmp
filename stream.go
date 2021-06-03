package main

import (
	"log"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/yumrano/rtsp2rtmp/conf"
	"github.com/yumrano/rtsp2rtmp/dao"
	"github.com/yumrano/rtsp2rtmp/writer/flvfile"
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

type RtspManager struct {
	ffm  *flvfile.FlvFileManager
	hfm  *httpflv.HttpFlvManager
	frm  *rtmpflv.FlvRtmpManager
	done <-chan interface{}
}

func (r *RtspManager) pktTransfer(code string, codecs []av.CodecData, pRtmpFlvChan <-chan av.Packet, pFlvFileChan <-chan av.Packet, pHttpFlvChan <-chan av.Packet) {
	// go r.frm.FlvWrite(code, codecs, r.done, pRtmpFlvChan)
	go r.hfm.FlvWrite(code, codecs, r.done, pHttpFlvChan)
	save, err := conf.GetBool("server.flvfile.save")
	if err == nil && save {
		go r.ffm.FlvWrite(code, codecs, r.done, pFlvFileChan)
	}
}

func (s *Server) serveStreams() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("rtspManager pain %v", r)
		}
	}()
	es, err := dao.CameraSelectAll()
	if err != nil {
		log.Printf("camera list is empty")
		return
	}
	for _, camera := range es {
		go func(c dao.Camera) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("rtspManager pain %v", r)
				}
			}()
			for {
				log.Println(c.Code, "connect", c.RtspURL)
				rtsp.DebugRtsp = true
				session, err := rtsp.Dial(c.RtspURL)
				if err != nil {
					log.Println(c.Code, err)
					c.OnlineStatus = 0
					time.Sleep(5 * time.Second)
					if c.OnlineStatus == 1 {
						dao.CameraUpdate(c)
					}
					continue
				}
				session.RtpKeepAliveTimeout = 10 * time.Second
				if err != nil {
					log.Println(c.Code, err)
					time.Sleep(5 * time.Second)
					continue
				}
				codec, err := session.Streams()
				if err != nil {
					log.Println(c.Code, err)
					time.Sleep(5 * time.Second)
					continue
				}

				pRtmpFlvChan := make(chan av.Packet)
				pFlvFileChan := make(chan av.Packet)
				pHttpFlvChan := make(chan av.Packet)
				done := make(chan interface{})
				rm := &RtspManager{
					ffm:  flvfile.NewFlvFileManager(),
					hfm:  httpflv.NewHttpFlvManager(),
					frm:  rtmpflv.NewFlvRtmpManager(),
					done: done,
				}
				rm.pktTransfer(c.Code, codec, pRtmpFlvChan, pFlvFileChan, pHttpFlvChan)
				s.rms[c.Code] = rm

				c.OnlineStatus = 1
				dao.CameraUpdate(c)

				for {
					pkt, err := session.ReadPacket()
					if err != nil {
						log.Println(c.Code, err)
						break
					}
					select {
					case pRtmpFlvChan <- pkt:
					default:
					}
					select {
					case pFlvFileChan <- pkt:
					default:
					}
					select {
					case pHttpFlvChan <- pkt:
					default:
					}
				}
				err = session.Close()
				if err != nil {
					log.Println("session Close error", err)
				}
				log.Println(c.Code, "reconnect wait 5s")
				time.Sleep(5 * time.Second)
			}
		}(camera)
	}
}
