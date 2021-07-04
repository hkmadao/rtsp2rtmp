package fileflvwriter

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/models"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
)

type IFileFlvManager interface {
	UpdateFFWS(string, *FileFlvWriter)
}

type FileFlvWriter struct {
	selfDone  chan interface{}
	pktStream <-chan av.Packet
	code      string
	codecs    []av.CodecData
	isStart   bool
	fd        *os.File
	iffm      IFileFlvManager
}

func (ffw *FileFlvWriter) GetPktStream() <-chan av.Packet {
	return ffw.pktStream
}

func (ffw *FileFlvWriter) GetCodecs() []av.CodecData {
	return ffw.codecs
}

func NewFileFlvWriter(
	pktStream <-chan av.Packet,
	code string,
	codecs []av.CodecData,
	iffm IFileFlvManager,
) *FileFlvWriter {

	ffw := &FileFlvWriter{
		selfDone:  make(chan interface{}, 10),
		pktStream: pktStream,
		code:      code,
		codecs:    codecs,
		iffm:      iffm,
		isStart:   false,
	}
	camera, err := models.CameraSelectOne(models.Camera{Code: code})
	if err != nil {
		logs.Error("query camera error : %v", err)
		return ffw
	}
	if camera.OnlineStatus != 1 {
		return ffw
	}
	if camera.SaveVideo != 1 {
		return ffw
	}
	go ffw.flvWrite()
	go ffw.splitFile()
	return ffw
}

func (ffw *FileFlvWriter) StopWrite() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		//有多个地方监听seleDone,需要写入多次才能退出多个goroutine
		for i := 0; i < 10; i++ {
			select {
			case ffw.selfDone <- struct{}{}:
			default:
			}
		}
	}()
}

func (ffw *FileFlvWriter) splitFile() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-ffw.selfDone:
			return
		case <-time.After(1 * time.Hour):
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
					}
				}()
				//有多个地方监听seleDone,需要写入多次才能退出多个goroutine
				for i := 0; i < 10; i++ {
					select {
					case ffw.selfDone <- struct{}{}:
					default:
					}
				}
				ffw.fd.Close()
			}()
			ffwn := NewFileFlvWriter(ffw.pktStream, ffw.code, ffw.codecs, ffw.iffm)
			ffwn.iffm.UpdateFFWS(ffwn.code, ffwn)
			continue
		}
	}
}

func (ffw *FileFlvWriter) Write(p []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	n, err = ffw.fd.Write(p)
	if err != nil {
		logs.Error("write file error : %v", err)
	}
	return
}

func (ffw *FileFlvWriter) createFlvFile() error {
	fd, err := os.OpenFile(getFileFlvPath()+"/"+ffw.code+"_"+time.Now().Format("20060102150405")+".flv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logs.Error("open file error :", err)
		return err
	}
	ffw.fd = fd
	return nil
}

//Write extends to writer.Writer
func (ffw *FileFlvWriter) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	if err := ffw.createFlvFile(); err != nil {
		logs.Error("create file flv error : %v", err)
		return
	}
	defer close(ffw.selfDone)
	defer ffw.fd.Close()
	muxer := flv.NewMuxer(ffw)
	for pkt := range utils.OrDonePacket(ffw.selfDone, ffw.pktStream) {
		if ffw.isStart {
			if err := muxer.WritePacket(pkt); err != nil {
				logs.Error("writer packet to flv file error : %v", err)
			}
			continue
		}
		if pkt.IsKeyFrame {
			err := muxer.WriteHeader(ffw.codecs)
			if err != nil {
				logs.Error("writer header to flv file error : %v", err)
			}
			if err := muxer.WritePacket(pkt); err != nil {
				logs.Error("writer packet to flv file error : %v", err)
			}
			ffw.isStart = true
		}
	}
}

func getFileFlvPath() string {
	fileFlvPath, err := config.String("server.fileflv.path")
	if err != nil {
		logs.Error("get fileflv path error :", err)
		return ""
	}
	return fileFlvPath
}
