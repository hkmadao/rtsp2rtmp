package rtmpflv

import (
	"runtime/debug"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
)

type RtmpFlvManager struct {
	done      <-chan interface{}
	pktStream <-chan av.Packet
	code      string
	codecs    []av.CodecData
}

func NewRtmpFlvManager(done <-chan interface{}, pktStream <-chan av.Packet, code string, codecs []av.CodecData) *RtmpFlvManager {
	rfm := &RtmpFlvManager{
		done:      done,
		pktStream: pktStream,
		code:      code,
		codecs:    codecs,
	}
	go rfm.flvWrite()
	return rfm
}

func (rfm *RtmpFlvManager) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	NewRtmpFlvWriter(rfm.done, rfm.pktStream, rfm.code, rfm.codecs)
}
