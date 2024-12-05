package fileflvwriter

import (
	"encoding/hex"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service"
)

type IFileFlvManager interface {
	DeleteFFW(sesessionId int64)
}

type FileFlvWriter struct {
	sessionId   int64
	done        chan int
	fgDoneClose bool
	tickerDone  chan int
	pktStream   chan av.Packet
	code        string
	codecs      []av.CodecData
	isStart     bool
	fd          *os.File
	fileName    string
	muxer       *flv.Muxer
	startTime   time.Time
	endTime     time.Time
	ffm         IFileFlvManager
}

func (ffw *FileFlvWriter) GetDone() <-chan int {
	return ffw.done
}

func (ffw *FileFlvWriter) GetCode() string {
	return ffw.code
}

func (ffw *FileFlvWriter) GetPktStream() chan av.Packet {
	return ffw.pktStream
}

func (ffw *FileFlvWriter) SetCodecs(codecs []av.CodecData) {
	ffw.codecs = codecs
}

func (ffw *FileFlvWriter) GetCodecs() []av.CodecData {
	return ffw.codecs
}

func (ffw *FileFlvWriter) GetSessionId() int64 {
	return ffw.sessionId
}

func NewFileFlvWriter(
	sessionId int64,
	pktStream chan av.Packet,
	code string,
	codecs []av.CodecData,
	ffm IFileFlvManager,
) *FileFlvWriter {

	ffw := &FileFlvWriter{
		sessionId:   sessionId,
		fgDoneClose: false,
		done:        make(chan int),
		tickerDone:  make(chan int),
		pktStream:   pktStream,
		code:        code,
		codecs:      codecs,
		isStart:     false,
		ffm:         ffm,
	}
	camera, err := service.CameraSelectOne(entity.Camera{Code: code})
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
				case _, ok := <-ffw.pktStream:
					if !ok {
						return
					}
				}
			}
		}()
		return ffw
	}
	go ffw.flvWrite()
	return ffw
}

func (ffw *FileFlvWriter) StopWrite() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		ffw.ffm.DeleteFFW(ffw.sessionId)
		ffw.fgDoneClose = true
		close(ffw.done)
	}()
}

func (ffw *FileFlvWriter) TickerStopWrite() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		select {
		case <-time.NewTicker(30 * time.Second).C: //等待30秒再关闭
			ffw.ffm.DeleteFFW(ffw.sessionId)
			ffw.fgDoneClose = true
			close(ffw.done)
		case <-ffw.GetDone():
		}
	}()
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
	fileName := getFileFlvPath() + "/" + ffw.code + "_" + time.Now().Format("20060102150405") + "_temp.flv"
	fd, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logs.Error("open file error :", err)
		return err
	}
	ffw.fd = fd
	ffw.fileName = fileName
	return nil
}

// Write extends to writer.Writer
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
		ffw.endTime = time.Now()
		ffw.muxer.WriteTrailer()
		ffw.fd.Close()

		//写入script tag data，主要补充视频的总时长，否则使用播放器播放看不到视频总时长
		ffw.writeScriptTagData()

		if !ffw.fgDoneClose {
			close(ffw.done)
		}
	}()

	muxer := flv.NewMuxer(ffw)
	ffw.muxer = muxer
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
			ffw.startTime = time.Now()
			continue
		}
		if time.Now().Local().After(timeNow.Add(1 * time.Minute)) {
			timeNow = time.Now().Local()
			logs.Error("FileFlvWriter ingrore package: %s", ffw.code)
		}
	}
}

func (ffw *FileFlvWriter) writeScriptTagData() {
	reverseFileName := utils.ReverseString(ffw.fileName)
	reverseNewFileName := strings.Replace(reverseFileName, utils.ReverseString("_temp.flv"), utils.ReverseString(".flv"), 1)
	newFileName := utils.ReverseString(reverseNewFileName)
	newflvFile, err := os.OpenFile(newFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logs.Error("create flv file error :", err)
		return
	}
	flvFile, err := os.OpenFile(ffw.fileName, os.O_RDWR, 0644)
	if err != nil {
		logs.Error("open file error :", err)
		return
	}
	buf := make([]byte, 10*1024)
	i := 1
	duration := float64(ffw.endTime.Sub(ffw.startTime).Seconds())
	durationBytes := utils.Float64ToByteBigEndian(duration)
	durationHexStr := hex.EncodeToString(durationBytes)
	scriptTagHexStr := "120000250000000000000002000A6F6E4D65746144617461080000000100086475726174696F6E00" + durationHexStr + "00000030"
	scriptTagBytes, err := hex.DecodeString(scriptTagHexStr)
	if err != nil {
		logs.Error("scriptTagHexStr: %s, DecodeString error : ", scriptTagHexStr, err)
		return
	}
	for {
		_, err := flvFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			logs.Error("read flv file error : %v", err)
		}

		if i == 1 {
			i = 2
			data1 := make([]byte, len(buf)+52)
			copy(data1, buf[:13])
			newData := append(data1[:13], scriptTagBytes...)
			newData = append(newData, buf[13:]...)
			newflvFile.Write(newData)
			continue
		}
		newflvFile.Write(buf)
	}
	err = flvFile.Close()
	if err != nil {
		logs.Error("close template flv file error :", err)
		return
	}
	err = os.Remove(ffw.fileName)
	if err != nil {
		logs.Error("remove template flv file error :", err)
		return
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
