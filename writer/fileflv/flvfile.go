package fileflv

import (
	"os"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
	"github.com/yumrano/rtsp2rtmp/conf"
	"github.com/yumrano/rtsp2rtmp/rlog"
)

type FileFlvManager struct {
	fw *FileFlvWriter
}

func NewFileFlvManager() *FileFlvManager {
	return &FileFlvManager{}
}

func (fm *FileFlvManager) codec(code string, codecs []av.CodecData) {
	fd, err := os.OpenFile(getFileFlvPath()+"/"+code+"_"+time.Now().Format("20060102150405")+".flv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		rlog.Log.Println("open file error :", err)
	}
	fm.fw = &FileFlvWriter{
		codecs: codecs,
		code:   code,
		fd:     fd,
	}
}

//Write extends to writer.Writer
func (fm *FileFlvManager) FlvWrite(code string, codecs []av.CodecData, done <-chan interface{}, pchan <-chan av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			rlog.Log.Printf("FlvFileManager FlvWrite pain %v", r)
		}
	}()
	fm.codec(code, codecs)
	muxer := flv.NewMuxer(fm.fw)
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-done:
			fm.fw.fd.Close()
			return
		case <-ticker.C: //split flvFile
			fd, err := os.OpenFile(getFileFlvPath()+"/"+fm.fw.code+"_"+time.Now().Format("20060102150405")+".flv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				rlog.Log.Println("open file error :", err)
			}
			fdOld := fm.fw.fd
			fm.fw.prepare = true
			fm.fw.isStart = false
			fm.fw.fd = fd
			fm.fw.prepare = false
			fdOld.Close()
		case pkt := <-pchan:
			if fm.fw.isStart {
				if err := muxer.WritePacket(pkt); err != nil {
					rlog.Log.Printf("writer packet to flv file error : %v\n", err)
				}
				continue
			}
			if pkt.IsKeyFrame {
				err := muxer.WriteHeader(fm.fw.codecs)
				if err != nil {
					rlog.Log.Printf("writer header to flv file error : %v\n", err)
				}
				if err := muxer.WritePacket(pkt); err != nil {
					rlog.Log.Printf("writer packet to flv file error : %v\n", err)
				}
				fm.fw.isStart = true
			}
		}
	}
}

type FileFlvWriter struct {
	code    string
	isStart bool
	prepare bool
	fd      *os.File
	codecs  []av.CodecData
}

//Write extends to io.Writer
func (fw *FileFlvWriter) Write(p []byte) (n int, err error) {
	if fw.prepare {
		return
	}
	n, err = fw.fd.Write(p)
	if err != nil {
		rlog.Log.Println("write file error :", err)
	}
	return
}

func getFileFlvPath() string {
	fileFlvPath, err := conf.GetString("server.fileflv.path")
	if err != nil {
		rlog.Log.Println("get fileflv path error :", err)
		return ""
	}
	return fileFlvPath
}
