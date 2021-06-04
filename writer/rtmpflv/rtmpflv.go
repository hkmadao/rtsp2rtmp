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
	}
	fm.fw = &RtmpFlvWriter{
		code: code,
		conn: rtmpConn,
	}
}

//Write extends to writer.Writer
func (fm *RtmpFlvManager) FlvWrite(code string, codecs []av.CodecData, done <-chan interface{}, pchan <-chan av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			rlog.Log.Printf("RtmpFlvManager FlvWrite pain %v", r)
		}
	}()
	fm.codec(code, codecs)
	for {
		select {
		case <-done:
			fm.fw.conn.Close()
			return
		case pkt := <-pchan:
			if fm.fw.isStart {
				if err := fm.fw.conn.WritePacket(pkt); err != nil {
					rlog.Log.Printf("writer packet to flv file error : %v\n", err)
				}
				continue
			}
			if pkt.IsKeyFrame {
				err := fm.fw.conn.WriteHeader(fm.fw.codecs)
				if err != nil {
					rlog.Log.Printf("writer header to flv file error : %v\n", err)
				}
				fm.fw.isStart = true
			}
		}
	}
}

type RtmpFlvWriter struct {
	code    string
	isStart bool
	conn    *rtmp.Conn
	codecs  []av.CodecData
}
