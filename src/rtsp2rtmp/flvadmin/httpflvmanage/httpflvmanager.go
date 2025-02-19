package httpflvmanage

import (
	"fmt"
	"io"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/httpflvmanage/httpflvwriter"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

type SyncMap struct {
	sync.Map
	size int32 // 原子计数器，用于跟踪map的大小
}

func (sm *SyncMap) Store(key, value interface{}) {
	sm.Map.Store(key, value)
	atomic.AddInt32(&sm.size, 1) // 每次存储时增加计数器
}

func (sm *SyncMap) Delete(key interface{}) {
	sm.Map.Delete(key)
	atomic.AddInt32(&sm.size, -1) // 每次删除时减少计数器
}

func (sm *SyncMap) IsEmpty() bool {
	return atomic.LoadInt32(&sm.size) == 0 // 加载计数器的值并检查是否为0
}

type HttpFlvManager struct {
	fgDoneClose bool
	done        chan int
	pktStream   <-chan av.Packet
	code        string
	codecs      []av.CodecData
	hfws        SyncMap
	mutex       sync.Mutex
}

func (hfm *HttpFlvManager) GetCode() string {
	return hfm.code
}

func (hfm *HttpFlvManager) SetCodecs(codecs []av.CodecData) {
	logs.Warn("HttpFlvManager: %s update codecs", hfm.code)
	hfm.codecs = codecs
	hfm.hfws.Range(func(key, value interface{}) bool {
		wi := value.(*httpflvwriter.HttpFlvWriter)
		wi.SetCodecs(hfm.codecs)
		return true
	})
}

func (hfm *HttpFlvManager) GetDone() <-chan int {
	return hfm.done
}

func (hfm *HttpFlvManager) GetPktStream() <-chan av.Packet {
	return hfm.pktStream
}

func (hfm *HttpFlvManager) GetCodecs() []av.CodecData {
	return hfm.codecs
}

func NewHttpFlvManager(pktStream <-chan av.Packet, code string, codecs []av.CodecData) *HttpFlvManager {
	hfm := &HttpFlvManager{
		fgDoneClose: false,
		done:        make(chan int),
		pktStream:   pktStream,
		code:        code,
		codecs:      codecs,
	}

	go hfm.flvWrite()
	return hfm
}

func (hfm *HttpFlvManager) StopWrite() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		hfm.CloseDone()
	}()
}

func (hfm *HttpFlvManager) CloseDone() {
	hfm.mutex.Lock()
	if !hfm.fgDoneClose {
		hfm.fgDoneClose = true
		close(hfm.done)
	}
	hfm.mutex.Unlock()
}

// Write extends to writer.Writer
func (hfm *HttpFlvManager) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	defer hfm.CloseDone()
	condition := common.GetEqualCondition("code", hfm.code)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("query camera error : %v", err)
		return
	}
	if !camera.OnlineStatus {
		return
	}
	if !camera.Live {
		for {
			select {
			case <-hfm.GetDone():
				return
			case _, ok := <-hfm.pktStream:
				if !ok {
					return
				}
			}
		}
	}
	for pkt := range utils.OrDonePacket(hfm.done, hfm.pktStream) {
		hfm.hfws.Range(func(key, value interface{}) bool {
			wi := value.(*httpflvwriter.HttpFlvWriter)
			select {
			case wi.GetPktStream() <- pkt:
				// logs.Debug("flvWrite pkt")
			default:
				//当播放者速率跟不上时，会发生丢包
				logs.Debug("camera [%s] http flv sessionId [%d] write fail", hfm.code, wi.GetSessionId())
			}
			return true
		})
	}
}

// 添加播放者
func (hfm *HttpFlvManager) AddHttpFlvPlayer(
	playerDone <-chan int,
	pulseInterval time.Duration,
	writer io.Writer,
) (<-chan int, error) {
	condition := common.GetEqualCondition("code", hfm.code)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("query camera error : %v", err)
		return nil, err
	}
	if !camera.OnlineStatus {
		return nil, fmt.Errorf("camera offline")
	}
	if !camera.Live {
		return nil, fmt.Errorf("camera live disabled")
	}
	sessionId := utils.NextValSnowflakeID()
	//添加缓冲
	pktStream := make(chan av.Packet, 1024)
	hfw := httpflvwriter.NewHttpFlvWriter(hfm.GetDone(), playerDone, pulseInterval, pktStream, hfm.code, hfm.codecs, writer, sessionId, hfm)
	hfm.hfws.Store(sessionId, hfw)
	go func() {
		<-hfw.GetDone()
		hfm.hfws.Delete(sessionId)
	}()
	return hfw.GetDone(), nil
}

func (hfm *HttpFlvManager) DeleteHFW(sesessionId int64) {
	hfm.hfws.LoadAndDelete(sesessionId)
}

func (hfm *HttpFlvManager) IsCameraExistsPlayer() bool {
	return hfm.hfws.IsEmpty()
}
