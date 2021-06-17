package httpflv

import (
	"errors"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/yumrano/rtsp2rtmp/utils"
)

var hfms sync.Map

type writerInfo struct {
	sessionId       int64
	code            string
	heartbeatStream <-chan interface{}
	endStream       <-chan interface{}
	pktStream       chan<- av.Packet
}

//添加播放者
func AddHttpFlvPlayer(code string, writer http.ResponseWriter) (endStream <-chan interface{}, heartbeatStream <-chan interface{}, playerDone chan<- interface{}, err error) {
	v, b := hfms.Load(code)
	if !b {
		err = errors.New("camera no connection")
		return
	}
	sessionId := utils.NextValSnowflakeID()
	hfm := v.(*HttpFlvManager)
	hfw := NewHttpFlvWriter(hfm.done, hfm.code, hfm.codecs, writer, sessionId)
	endStream = hfw.GetEndStream()
	//one2two chan
	heartbeatStream1, heartbeatStream2 := utils.Tee(endStream, hfw.GetHeartbeatStream(), 1*time.Millisecond)
	wi := &writerInfo{
		sessionId:       sessionId,
		code:            code,
		heartbeatStream: heartbeatStream1,
		endStream:       endStream,
		pktStream:       hfw.GetPktStream(),
	}
	hfm.wis.Store(sessionId, wi)
	heartbeatStream = heartbeatStream2
	playerDone = hfw.GetPlayerDone()
	go monitor(wi)
	return
}

func monitor(wi *writerInfo) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-wi.heartbeatStream:
			continue
		case <-wi.endStream:
			//end info
			if v, b := hfms.Load(wi.code); b {
				hfm := v.(*HttpFlvManager)
				hfm.wis.Delete(wi.sessionId)
			}
			return
		case <-time.After(10 * time.Second):
			//time out
			if v, b := hfms.Load(wi.code); b {
				hfm := v.(*HttpFlvManager)
				hfm.wis.Delete(wi.sessionId)
			}
			return
		}
	}
}

func ExistsHttpFlvManager(code string) bool {
	_, b := hfms.Load(code)
	return b
}

type HttpFlvManager struct {
	done      <-chan interface{}
	pktStream <-chan av.Packet
	code      string
	codecs    []av.CodecData
	wis       sync.Map
}

func NewHttpFlvManager(done <-chan interface{}, pktStream <-chan av.Packet, code string, codecs []av.CodecData) *HttpFlvManager {
	hfm := &HttpFlvManager{
		done:      done,
		pktStream: pktStream,
		code:      code,
		codecs:    codecs,
	}
	go hfm.flvWrite()
	hfms.Store(code, hfm)
	return hfm
}

//Write extends to writer.Writer
func (hfm *HttpFlvManager) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-hfm.done:
			return
		case pkt := <-hfm.pktStream:
			hfm.wis.Range(func(key, value interface{}) bool {
				go func(pkt1 av.Packet) {
					defer func() {
						if r := recover(); r != nil {
							logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
						}
					}()
					wi := value.(*writerInfo)
					select {
					case wi.pktStream <- pkt1:
					case <-time.After(1 * time.Millisecond):
						// logs.Info("lose pkt")
					}
				}(pkt)
				return true
			})
		}
	}
}
