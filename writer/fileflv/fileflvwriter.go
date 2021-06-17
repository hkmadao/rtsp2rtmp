package fileflv

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
)

type FileFlvWriter struct {
	done      <-chan interface{}
	pktStream <-chan av.Packet
	code      string
	codecs    []av.CodecData
	start     bool
	prepare   bool
	fd        *os.File
}

func NewFileFlvWriter(done <-chan interface{}, pktStream <-chan av.Packet, code string, codecs []av.CodecData) *FileFlvWriter {
	ffw := &FileFlvWriter{
		done:      done,
		pktStream: pktStream,
		code:      code,
		codecs:    codecs,
	}
	go ffw.flvWrite()
	return ffw
}

func (ffw *FileFlvWriter) Write(p []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	if ffw.prepare {
		return
	}
	n, err = ffw.fd.Write(p)
	if err != nil {
		logs.Error("write file error : %v", err)
	}
	return
}

func (ffw *FileFlvWriter) createFlvFile() {
	fd, err := os.OpenFile(getFileFlvPath()+"/"+ffw.code+"_"+time.Now().Format("20060102150405")+".flv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logs.Error("open file error :", err)
	}
	if ffw.fd == nil {
		ffw.fd = fd
		return
	}
	fdOld := ffw.fd
	ffw.prepare = true
	ffw.start = false
	ffw.fd = fd
	ffw.prepare = false
	fdOld.Close()
}

//Write extends to writer.Writer
func (ffw *FileFlvWriter) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	ffw.createFlvFile()
	muxer := flv.NewMuxer(ffw)
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-ffw.done:
			ffw.fd.Close()
			return
		case <-ticker.C: //split flvFile
			ffw.createFlvFile()
		case pkt := <-ffw.pktStream:
			if ffw.start {
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
				ffw.start = true
			}
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
