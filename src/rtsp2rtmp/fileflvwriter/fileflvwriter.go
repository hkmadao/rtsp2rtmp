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
	done      chan int
	pktStream <-chan av.Packet
	code      string
	codecs    []av.CodecData
	isStart   bool
	fd        *os.File
	iffm      IFileFlvManager
}

func (ffw *FileFlvWriter) GetDone() <-chan int {
	return ffw.done
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
		done:      make(chan int),
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
		go func() {
			for {
				select {
				case <-ffw.GetDone():
					return
				case <-ffw.pktStream:
				}
			}
		}()
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
		close(ffw.done)
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
		case <-ffw.done:
			return
		case <-time.After(1 * time.Hour):
			ffw.StopWrite()
			_, pktStreamOk := <-ffw.pktStream
			if pktStreamOk {
				logs.Info("to create NewFileFlvWriter : %s", ffw.code)
				ffwn := NewFileFlvWriter(ffw.pktStream, ffw.code, ffw.codecs, ffw.iffm)
				ffwn.iffm.UpdateFFWS(ffwn.code, ffwn)
			} else {
				logs.Info("FileFlvWriter pktStream is closed : %s", ffw.code)
			}
			return
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
	defer func() {
		close(ffw.done)
		ffw.fd.Close()
	}()
	muxer := flv.NewMuxer(ffw)
	timeNow := time.Now().Local()
	for pkt := range utils.OrDonePacket(ffw.done, ffw.pktStream) {
		if ffw.isStart {
			if err := muxer.WritePacket(pkt); err != nil {
				logs.Error("writer packet to flv file error : %v", err)
			}
			continue
		}
		if pkt.IsKeyFrame {
			ffw.isStart = true
			err := muxer.WriteHeader(ffw.codecs)
			if err != nil {
				logs.Error("writer header to flv file error : %v", err)
				ffw.isStart = false
			}
			if err := muxer.WritePacket(pkt); err != nil {
				logs.Error("writer packet to flv file error : %v", err)
				ffw.isStart = false
			}
			continue
		}
		if time.Now().Local().After(timeNow.Add(1 * time.Minute)) {
			timeNow = time.Now().Local()
			logs.Error("FileFlvWriter ingrore package: %s", ffw.code)
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
