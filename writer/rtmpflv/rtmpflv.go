package rtmpflv

import (
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/yumrano/rtsp2rtmp/dao"
	"github.com/yumrano/rtsp2rtmp/rlog"
)

type RtmpFlvManager struct {
	fw *RtmpFlvWriter
}

func NewRtmpFlvManager() *RtmpFlvManager {
	return &RtmpFlvManager{}
}

func (fm *RtmpFlvManager) IsStop() bool {
	if fm.fw == nil {
		return true
	}
	return fm.fw.stop
}

func (fm *RtmpFlvManager) codec(code string, codecs []av.CodecData) {
	var camera dao.Camera
	camera.Code = code
	camera, err := dao.CameraSelectOne(camera)
	if err != nil {
		rlog.Log.Printf("not found camera : %s", code)
		return
	}
	rtmpConn, err := rtmp.Dial(camera.RtmpURL)
	if err != nil {
		rlog.Log.Printf("rtmp client connection error : %v", err)
		return
	}
	fm.fw = &RtmpFlvWriter{
		code:   code,
		conn:   rtmpConn,
		codecs: codecs,
	}
}

//Write extends to writer.Writer
func (fm *RtmpFlvManager) FlvWrite(code string, codecs []av.CodecData, done <-chan interface{}, pchan <-chan av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			rlog.Log.Printf("RtmpFlvManager FlvWrite panic %v", r)
		}
	}()
	fm.codec(code, codecs)
	for {
		select {
		case <-done:
			fm.fw.conn.Close()
			return
		case pkt := <-pchan:
			if fm.fw.start {
				if err := fm.fw.conn.WritePacket(pkt); err != nil {
					rlog.Log.Printf("writer packet to rtmp server error : %v\n", err)
					fm.fw.stop = true
					fm.fw.conn.Close()
				}
				continue
			}
			if pkt.IsKeyFrame {
				err := fm.fw.conn.WriteHeader(fm.fw.codecs)
				if err != nil {
					rlog.Log.Printf("writer header to rtmp server error : %v\n", err)
					fm.fw.stop = true
					fm.fw.conn.Close()
					continue
				}
				fm.fw.start = true
				err = fm.fw.conn.WritePacket(pkt)
				if err != nil {
					rlog.Log.Printf("writer packet to rtmp server error : %v\n", err)
					fm.fw.stop = true
					fm.fw.conn.Close()
					continue
				}
			}
		}
	}
}

type RtmpFlvWriter struct {
	code   string
	start  bool
	stop   bool
	conn   *rtmp.Conn
	codecs []av.CodecData
}
