package fileflvmanager

import (
	"runtime/debug"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/fileflvmanager/fileflvwriter"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/models"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
)

type FileFlvManager struct {
	done      chan int
	pktStream <-chan av.Packet
	code      string
	codecs    []av.CodecData
	ffws      sync.Map
}

func (ffm *FileFlvManager) GetCode() string {
	return ffm.code
}

func (ffm *FileFlvManager) SetCodecs(codecs []av.CodecData) {
	ffm.codecs = codecs
	ffm.ffws.Range(func(key, value interface{}) bool {
		wi := value.(*fileflvwriter.FileFlvWriter)
		wi.SetCodecs(ffm.codecs)
		return true
	})
}

func (ffm *FileFlvManager) GetDone() <-chan int {
	return ffm.done
}

func (ffm *FileFlvManager) GetPktStream() <-chan av.Packet {
	return ffm.pktStream
}

func (ffm *FileFlvManager) GetCodecs() []av.CodecData {
	return ffm.codecs
}

func NewFileFlvManager(pktStream <-chan av.Packet, code string, codecs []av.CodecData) *FileFlvManager {
	ffm := &FileFlvManager{
		done:      make(chan int),
		pktStream: pktStream,
		code:      code,
		codecs:    codecs,
	}
	camera, err := models.CameraSelectOne(models.Camera{Code: code})
	if err != nil {
		logs.Error("query camera error : %v", err)
		return ffm
	}
	if camera.OnlineStatus != 1 {
		return ffm
	}
	if camera.SaveVideo != 1 {
		go func() {
			for {
				select {
				case <-ffm.GetDone():
					return
				case _, ok := <-ffm.pktStream:
					if !ok {
						return
					}
				}
			}
		}()
		return ffm
	}
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for {
			select {
			case <-ffm.done:
				return
			case <-ticker.C:
				ffm.ffws.Range(func(key, value interface{}) bool {
					ffw := value.(*fileflvwriter.FileFlvWriter)
					if ffw.GetCode() == code {
						ffw.TickerStopWrite()
					}
					return true
				})
				sessionId := utils.NextValSnowflakeID()
				//添加缓冲，减少包到达速率震荡导致丢包
				pktStream := make(chan av.Packet, 1024)
				newFfw := fileflvwriter.NewFileFlvWriter(sessionId, pktStream, code, ffm.codecs, ffm)
				ffm.ffws.Store(sessionId, newFfw)
			}
		}
	}()
	sessionId := utils.NextValSnowflakeID()
	//添加缓冲，减少包到达速率震荡导致丢包
	ffwPktStream := make(chan av.Packet, 1024)
	newFfw := fileflvwriter.NewFileFlvWriter(sessionId, ffwPktStream, code, codecs, ffm)
	ffm.ffws.Store(sessionId, newFfw)
	go ffm.flvWrite()
	return ffm
}

func (ffm *FileFlvManager) StopWrite() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		close(ffm.done)
	}()
}

//Write extends to writer.Writer
func (ffm *FileFlvManager) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for pkt := range utils.OrDonePacket(ffm.done, ffm.pktStream) {
		ffm.ffws.Range(func(key, value interface{}) bool {
			ffw := value.(*fileflvwriter.FileFlvWriter)
			select {
			case ffw.GetPktStream() <- pkt:
				// logs.Debug("flvWrite pkt")
			default:
				//当播放者速率跟不上时，会发生丢包
				logs.Debug("camera [%s] file flv write fail", ffm.code)
			}
			return true
		})
	}
}

func (ffm *FileFlvManager) DeleteFFW(sesessionId int64) {
	ffm.ffws.LoadAndDelete(sesessionId)
}
