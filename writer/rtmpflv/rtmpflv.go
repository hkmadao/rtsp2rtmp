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
	if fm.fw.errTime > 9 {
		if fm.fw.conn != nil {
			err := fm.fw.conn.Close()
			if err != nil {
				rlog.Log.Printf("conn close error : %v", err)
			}
		}
		return true
	}
	return false
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
			rlog.Log.Printf("RtmpFlvManager FlvWrite pain %v", r)
			if fm.fw != nil {
				if fm.fw.conn != nil {
					fm.fw.conn.Close()
				}
				fm.fw.errTime = 99
			}
		}
	}()
	fm.codec(code, codecs)
	for {
		select {
		case <-done:
			fm.fw.conn.Close()
			return
		case pkt := <-pchan:
			if fm.fw.isStart && fm.fw.errTime < 10 {
				if err := fm.fw.conn.WritePacket(pkt); err != nil {
					rlog.Log.Printf("writer packet to rtmp server error : %v\n", err)
					fm.fw.errTime = fm.fw.errTime + 1
					continue
				}
				fm.fw.errTime = 0
				continue
			}
			if pkt.IsKeyFrame && fm.fw.errTime < 10 {
				err := fm.fw.conn.WriteHeader(fm.fw.codecs)
				if err != nil {
					rlog.Log.Printf("writer header to rtmp server error : %v\n", err)
				}
				fm.fw.isStart = true
				err = fm.fw.conn.WritePacket(pkt)
				if err != nil {
					rlog.Log.Printf("writer packet to rtmp server error : %v\n", err)
					fm.fw.errTime = fm.fw.errTime + 1
					continue
				}
				fm.fw.errTime = 0
			}
		}
	}
}

type RtmpFlvWriter struct {
	code    string
	isStart bool
	errTime int //发送数据包连续失败次数
	conn    *rtmp.Conn
	codecs  []av.CodecData
}
