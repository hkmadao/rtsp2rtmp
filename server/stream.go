package server

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

func (s *Server) ServeStreams() {
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

				pRtmpFlvDone := make(chan interface{})
				pFlvFileDone := make(chan interface{})
				pHttpFlvDone := make(chan interface{})
				pRtmpFlvChan := make(chan av.Packet, 10)
				pFlvFileChan := make(chan av.Packet, 10)
				pHttpFlvChan := make(chan av.Packet, 10)
				rm := &RtspManager{
					ffm:          fileflv.NewFileFlvManager(),
					hfm:          httpflv.NewHttpFlvManager(),
					frm:          rtmpflv.NewRtmpFlvManager(),
					pRtmpFlvDone: pRtmpFlvDone,
					pFlvFileDone: pFlvFileDone,
					pHttpFlvDone: pHttpFlvDone,
					pRtmpFlvChan: pRtmpFlvChan,
					pFlvFileChan: pFlvFileChan,
					pHttpFlvChan: pHttpFlvChan,
				}
				rm.pktTransfer(c.Code, codec)
				s.rms[c.Code] = rm

				c.OnlineStatus = 1
				dao.CameraUpdate(c)

				for {
					pkt, err := session.ReadPacket()
					if err != nil {
						rlog.Log.Println(c.Code, err)
						break
					}
					writeChan(pkt, pRtmpFlvChan, pRtmpFlvDone)
					writeChan(pkt, pFlvFileChan, pFlvFileDone)
					writeChan(pkt, pHttpFlvChan, pHttpFlvDone)
				}
				close(pRtmpFlvDone)
				close(pFlvFileDone)
				close(pHttpFlvDone)
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
	defer func() {
		if r := recover(); r != nil {
			rlog.Log.Printf("writeChan painc : %v", r)
		}
	}()
	select {
	case c <- pkt:
	case <-time.After(1 * time.Millisecond):
	case <-done:
	}
}

type RtspManager struct {
	ffm          *fileflv.FileFlvManager
	hfm          *httpflv.HttpFlvManager
	frm          *rtmpflv.RtmpFlvManager
	pRtmpFlvDone chan interface{}
	pFlvFileDone chan interface{}
	pHttpFlvDone chan interface{}
	pRtmpFlvChan chan av.Packet
	pFlvFileChan chan av.Packet
	pHttpFlvChan chan av.Packet
}

func (r *RtspManager) pktTransfer(code string, codecs []av.CodecData) {
	go r.startRtmpClient(code, codecs)
	go r.hfm.FlvWrite(code, codecs, r.pHttpFlvDone, r.pHttpFlvChan)
	save, err := conf.GetBool("server.fileflv.save")
	if err != nil && save {
		rlog.Log.Printf("get server.fileflv.save error : %v", err)
		return
	}
	if save {
		go r.ffm.FlvWrite(code, codecs, r.pFlvFileDone, r.pFlvFileChan)
	}
}

func (r *RtspManager) startRtmpClient(code string, codecs []av.CodecData) {
	go r.frm.FlvWrite(code, codecs, r.pRtmpFlvDone, r.pRtmpFlvChan)
	for {
		if r.frm.IsStop() {
			select {
			case r.pFlvFileDone <- nil: //结束上一个goroutine
			case <-time.After(1 * time.Millisecond):
			}
			go r.frm.FlvWrite(code, codecs, r.pRtmpFlvDone, r.pRtmpFlvChan)
		}
		select {
		case <-r.pRtmpFlvDone:
			return
		case <-time.After(30 * time.Second):
		}
	}
}
