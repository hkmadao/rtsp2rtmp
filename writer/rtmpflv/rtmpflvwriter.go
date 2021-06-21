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
	heartbeatStream chan interface{}
	endStream       chan interface{}
}

func NewRtmpFlvWriter(done <-chan interface{}, pktStream <-chan av.Packet, code string, codecs []av.CodecData) *RtmpFlvWriter {
	rfw := &RtmpFlvWriter{
		done:            done,
		pktStream:       pktStream,
		code:            code,
		codecs:          codecs,
		start:           false,
		heartbeatStream: make(chan interface{}),
		endStream:       make(chan interface{}),
	}
	go rfw.flvWrite(rfw.endStream, rfw.heartbeatStream)
	go rfw.monitor(rfw.endStream, rfw.heartbeatStream)
	return rfw
}

func (rfw *RtmpFlvWriter) monitor(endStream <-chan interface{}, heartbeatStream <-chan interface{}) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	for {
		select {
		case <-endStream:
			return
		case <-rfw.done:
			if rfw.conn == nil {
				return
			}
			err := rfw.conn.Close()
			if err != nil {
				logs.Error("close rtmp client connection error : %v", err)
			}
			return
		case <-time.After(10 * time.Second):
			logs.Error("rtmp cliect time out , close")
			if rfw.conn != nil {
				err := rfw.conn.Close()
				if err != nil {
					logs.Error("close rtmp client connection error : %v", err)
				}
			}
			rfw.start = false
			go rfw.flvWrite(rfw.endStream, rfw.heartbeatStream)
			continue
		case <-heartbeatStream:
			continue
		}

	}
}

func (rfw *RtmpFlvWriter) createConn() error {
	var camera models.Camera
	camera.Code = rfw.code
	camera, err := models.CameraSelectOne(camera)
	if err != nil {
		logs.Error("not found camera : %s", rfw.code)
		return err
	}
	rtmpConn, err := rtmp.Dial(camera.RtmpURL)
	if err != nil {
		logs.Error("rtmp client connection error : %v", err)
		return err
	}
	rfw.conn = rtmpConn
	return nil
}

//Write extends to writer.Writer
func (rfw *RtmpFlvWriter) flvWrite(endStream chan<- interface{}, heartbeatStream chan<- interface{}) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-rfw.done:
			if rfw.conn == nil {
				return
			}
			err := rfw.conn.Close()
			if err != nil {
				logs.Error("close rtmp client connection error : %v", err)
			}
			return
		case pkt := <-rfw.pktStream:
			if rfw.start {
				if err := rfw.conn.WritePacket(pkt); err != nil {
					logs.Error("writer packet to rtmp server error : %v", err)
					select {
					case rfw.heartbeatStream <- 0:
					case <-time.After(1 * time.Millisecond):
						logs.Error("send heartbeat timeout")
					}
					err := rfw.conn.Close()
					if err != nil {
						logs.Error("close rtmp client connection error : %v", err)
					}
					return
				}
				select {
				case rfw.heartbeatStream <- 1:
				case <-time.After(1 * time.Millisecond):
					logs.Error("send heartbeat timeout")
				}
				continue
			}
			if pkt.IsKeyFrame {
				if err := rfw.createConn(); err != nil {
					logs.Error("conn rtmp server error : %v", err)
					return
				}
				var err error
				err = rfw.conn.WriteHeader(rfw.codecs)
				if err != nil {
					logs.Error("writer header to rtmp server error : %v", err)
					select {
					case rfw.heartbeatStream <- 0:
					case <-time.After(1 * time.Millisecond):
						logs.Error("send heartbeat timeout")
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
					logs.Error("writer packet to rtmp server error : %v", err)
					select {
					case rfw.heartbeatStream <- 0:
					case <-time.After(1 * time.Millisecond):
						logs.Error("send heartbeat timeout")
					}
					err := rfw.conn.Close()
					if err != nil {
						logs.Error("close rtmp client connection error : %v", err)
					}
					return
				}
				select {
				case rfw.heartbeatStream <- 1:
				case <-time.After(1 * time.Millisecond):
					logs.Error("send heartbeat timeout")
				}
				continue
			}
		}
	}
}
