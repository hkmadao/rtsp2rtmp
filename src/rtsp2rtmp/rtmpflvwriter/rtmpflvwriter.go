package rtmpflvwriter

import (
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/models"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
)

type IRtmpFlvManager interface {
	UpdateFFWS(string, *RtmpFlvWriter)
}

type RtmpFlvWriter struct {
	selfDone        chan interface{}
	pktStream       <-chan av.Packet
	code            string
	codecs          []av.CodecData
	start           bool
	conn            *rtmp.Conn
	pulseInterval   time.Duration
	heartbeatStream chan interface{}
	irfm            IRtmpFlvManager
}

func (rfw *RtmpFlvWriter) GetPktStream() <-chan av.Packet {
	return rfw.pktStream
}

func (rfw *RtmpFlvWriter) GetCodecs() []av.CodecData {
	return rfw.codecs
}

func NewRtmpFlvWriter(pktStream <-chan av.Packet, code string, codecs []av.CodecData, irfm IRtmpFlvManager) *RtmpFlvWriter {
	rfw := &RtmpFlvWriter{
		selfDone:        make(chan interface{}, 10),
		pktStream:       pktStream,
		code:            code,
		codecs:          codecs,
		start:           false,
		pulseInterval:   5 * time.Second,
		heartbeatStream: make(chan interface{}),
		irfm:            irfm,
	}
	go rfw.flvWrite()
	go rfw.monitor()
	return rfw
}

func (rfw *RtmpFlvWriter) StopWrite() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		//有多个地方监听seleDone,需要写入多次才能退出多个goroutine
		for i := 0; i < 10; i++ {
			select {
			case rfw.selfDone <- struct{}{}:
			default:
			}
		}
	}()
}

func (rfw *RtmpFlvWriter) monitor() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-rfw.selfDone:
			return
		case _, ok := <-rfw.heartbeatStream:
			if ok {
				logs.Debug("heartbeat")
				continue
			}
			rfw.StopWrite()
			rfwn := NewRtmpFlvWriter(rfw.pktStream, rfw.code, rfw.codecs, rfw.irfm)
			rfwn.irfm.UpdateFFWS(rfwn.code, rfwn)
			return
		case <-time.After(2 * rfw.pulseInterval):
			//time out
			rfw.StopWrite()
			rfwn := NewRtmpFlvWriter(rfw.pktStream, rfw.code, rfw.codecs, rfw.irfm)
			rfwn.irfm.UpdateFFWS(rfwn.code, rfwn)
			return
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
func (rfw *RtmpFlvWriter) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	pulse := time.NewTicker(rfw.pulseInterval).C

	sendPulse := func() {
		select {
		case rfw.heartbeatStream <- struct{}{}:
		default:
		}
	}

	sendPulse()
	done := make(chan interface{}, 10)
	go func(done <-chan interface{}) {
		defer func() {
			close(rfw.heartbeatStream)
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		for {
			select {
			case <-pulse:
				sendPulse()
			case <-done:
				return
			}
		}
	}(done)

	ticker := time.NewTicker(rfw.pulseInterval)
	defer func() {
		if rfw.conn != nil {
			rfw.conn.Close()
		}
	}()
	pktStream := utils.OrDonePacket(rfw.selfDone, rfw.pktStream)
	for {
		select {
		case <-ticker.C:
			close(done)
			return
		case pkt := <-pktStream:
			if err := rfw.writerPacket(pkt); err != nil {
				close(done)
				return
			}
			ticker.Reset(rfw.pulseInterval)
		}
	}
}

func (rfw *RtmpFlvWriter) writerPacket(pkt av.Packet) error {
	if rfw.start {
		if err := rfw.conn.WritePacket(pkt); err != nil {
			logs.Error("writer packet to rtmp server error : %v", err)
			return err
		}
		return nil
	}
	if pkt.IsKeyFrame {
		if err := rfw.createConn(); err != nil {
			logs.Error("conn rtmp server error : %v", err)
			return err
		}
		var err error
		err = rfw.conn.WriteHeader(rfw.codecs)
		if err != nil {
			logs.Error("writer header to rtmp server error : %v", err)
			return err
		}
		rfw.start = true
		err = rfw.conn.WritePacket(pkt)
		if err != nil {
			logs.Error("writer packet to rtmp server error : %v", err)
			return err
		}
		return nil
	}
	return nil
}
