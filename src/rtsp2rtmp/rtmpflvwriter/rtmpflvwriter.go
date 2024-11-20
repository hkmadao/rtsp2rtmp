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
	done          chan int
	pktStream     <-chan av.Packet
	code          string
	codecs        []av.CodecData
	start         bool
	conn          *rtmp.Conn
	pulseInterval time.Duration
	irfm          IRtmpFlvManager
}

func (rfw *RtmpFlvWriter) GetDone() <-chan int {
	return rfw.done
}

func (rfw *RtmpFlvWriter) GetPktStream() <-chan av.Packet {
	return rfw.pktStream
}

func (rfw *RtmpFlvWriter) GetCodecs() []av.CodecData {
	return rfw.codecs
}

func NewRtmpFlvWriter(pktStream <-chan av.Packet, code string, codecs []av.CodecData, irfm IRtmpFlvManager) *RtmpFlvWriter {
	rfw := &RtmpFlvWriter{
		done:          make(chan int),
		pktStream:     pktStream,
		code:          code,
		codecs:        codecs,
		start:         false,
		pulseInterval: 5 * time.Second,
		irfm:          irfm,
	}
	go rfw.flvWrite()
	return rfw
}

func (rfw *RtmpFlvWriter) StopWrite() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		close(rfw.done)
	}()
}

func (rfw *RtmpFlvWriter) createConn() error {
	var camera models.Camera
	camera.Code = rfw.code
	camera, err := models.CameraSelectOne(camera)
	if err != nil {
		logs.Error("not found camera : %s", rfw.code)
		return err
	}
	if camera.Enabled != 1 {
		go func() {
			select {
			case <-rfw.GetDone():
				return
			case <-rfw.pktStream:
				return
			}
		}()
	}
	rtmpConn, err := rtmp.Dial(camera.RtmpURL)
	if err != nil {
		logs.Error("rtmp client connection error : %v", err)
		return err
	}
	logs.Info("rtmp client connection success : %s", rfw.code)
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

	ticker := time.NewTicker(rfw.pulseInterval)
	defer func() {
		if rfw.conn != nil {
			rfw.conn.Close()
			_, pktStreamOk := <-rfw.pktStream
			if pktStreamOk {
				rfwn := NewRtmpFlvWriter(rfw.pktStream, rfw.code, rfw.codecs, rfw.irfm)
				rfwn.irfm.UpdateFFWS(rfwn.code, rfwn)
			}
		}
	}()
	pktStream := utils.OrDonePacket(rfw.done, rfw.pktStream)
	for {
		select {
		case <-ticker.C:
			return
		case pkt, ok := <-pktStream:
			if !ok {
				return
			}
			if err := rfw.writerPacket(pkt); err != nil {
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
	// logs.Debug("ingrore package")
	return nil
}
