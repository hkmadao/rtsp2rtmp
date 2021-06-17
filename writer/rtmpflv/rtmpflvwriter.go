package rtmpflv

import (
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/yumrano/rtsp2rtmp/models"
)

type RtmpFlvWriter struct {
	done            <-chan interface{}
	pktStream       <-chan av.Packet
	code            string
	codecs          []av.CodecData
	start           bool
	conn            *rtmp.Conn
	heartbeatStream chan int
}

func NewRtmpFlvWriter(done <-chan interface{}, pktStream <-chan av.Packet, code string, codecs []av.CodecData) *RtmpFlvWriter {
	rfw := &RtmpFlvWriter{
		done:            done,
		pktStream:       pktStream,
		code:            code,
		codecs:          codecs,
		start:           false,
		heartbeatStream: make(chan int),
	}
	go rfw.flvWrite()
	return rfw
}

func (rfw *RtmpFlvWriter) Monitor() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	for {
		select {
		case heartbeat := <-rfw.heartbeatStream:
			if heartbeat != 1 {
				go rfw.flvWrite()
				return
			}
		case <-rfw.done:
			return
		case <-time.After(10 * time.Second):
			logs.Error("rtmp cliect time out , close")
			err := rfw.conn.Close()
			if err != nil {
				logs.Error("close rtmp client connection error : %v", err)
			}
			go rfw.flvWrite()
		}

	}
}

func (rfw *RtmpFlvWriter) createConn() {
	var camera models.Camera
	camera.Code = rfw.code
	camera, err := models.CameraSelectOne(camera)
	if err != nil {
		logs.Error("not found camera : %s", rfw.code)
		return
	}
	rtmpConn, err := rtmp.Dial(camera.RtmpURL)
	if err != nil {
		logs.Error("rtmp client connection error : %v", err)
		return
	}
	rfw.conn = rtmpConn
}

//Write extends to writer.Writer
func (rfw *RtmpFlvWriter) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	rfw.createConn()
	for {
		select {
		case <-rfw.done:
			err := rfw.conn.Close()
			if err != nil {
				logs.Error("close rtmp client connection error : %v", err)
			}
			return
		case pkt := <-rfw.pktStream:
			if rfw.start {
				if err := rfw.conn.WritePacket(pkt); err != nil {
					logs.Error("writer packet to rtmp server error : %v\n", err)
					select {
					case rfw.heartbeatStream <- 0:
					case <-time.After(1 * time.Microsecond):
					}
					err := rfw.conn.Close()
					if err != nil {
						logs.Error("close rtmp client connection error : %v", err)
					}
					return
				}
				select {
				case rfw.heartbeatStream <- 1:
				case <-time.After(1 * time.Microsecond):
				}
				return
			}
			if pkt.IsKeyFrame {
				err := rfw.conn.WriteHeader(rfw.codecs)
				if err != nil {
					logs.Error("writer header to rtmp server error : %v\n", err)
					select {
					case rfw.heartbeatStream <- 0:
					case <-time.After(1 * time.Microsecond):
					}
					err := rfw.conn.Close()
					if err != nil {
						logs.Error("close rtmp client connection error : %v", err)
					}
					return
				}
				rfw.start = true
				err = rfw.conn.WritePacket(pkt)
				if err != nil {
					logs.Error("writer packet to rtmp server error : %v\n", err)
					select {
					case rfw.heartbeatStream <- 0:
					case <-time.After(1 * time.Microsecond):
					}
					err := rfw.conn.Close()
					if err != nil {
						logs.Error("close rtmp client connection error : %v", err)
					}
					return
				}
				select {
				case rfw.heartbeatStream <- 1:
				case <-time.After(1 * time.Microsecond):
				}
			}
		}
	}
}
