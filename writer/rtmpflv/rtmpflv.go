package rtmpflv

import (
	"fmt"
	"log"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/yumrano/rtsp2rtmp/dao"
)

type FlvRtmpManager struct {
	fw *FlvRtmpWriter
}

func NewFlvRtmpManager() *FlvRtmpManager {
	return &FlvRtmpManager{}
}

func (fm *FlvRtmpManager) codec(code string, codecs []av.CodecData) {
	var camera dao.Camera
	camera.Code = code
	camera, err := dao.CameraSelectOne(camera)
	if err != nil {
		log.Printf("not found camera : %s", code)
		return
	}
	rtmpConn, err := rtmp.Dial(camera.RtmpURL)
	if err != nil {
		fmt.Println("rtmp client connection error ", err)
	}
	fm.fw = &FlvRtmpWriter{
		code: code,
		conn: rtmpConn,
	}
}

//Write extends to writer.Writer
func (fm *FlvRtmpManager) FlvWrite(code string, codecs []av.CodecData, done <-chan interface{}, pchan <-chan av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("FlvRtmpManager FlvWrite pain %v", r)
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
					log.Printf("writer packet to flv file error : %v\n", err)
				}
				continue
			}
			if pkt.IsKeyFrame {
				err := fm.fw.conn.WriteHeader(fm.fw.codecs)
				if err != nil {
					log.Printf("writer header to flv file error : %v\n", err)
				}
				fm.fw.isStart = true
			}
		}
	}
}

type FlvRtmpWriter struct {
	code    string
	isStart bool
	conn    *rtmp.Conn
	codecs  []av.CodecData
}
