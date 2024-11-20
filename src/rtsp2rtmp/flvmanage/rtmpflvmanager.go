package flvmanage

import (
	"sync"

	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtmpflvwriter"
)

type rtmpFlvManager struct {
	rfms sync.Map
}

var rfmInstance *rtmpFlvManager

func init() {
	rfmInstance = &rtmpFlvManager{}
}

func GetSingleRtmpFlvManager() *rtmpFlvManager {
	return rfmInstance
}

func (rfm *rtmpFlvManager) FlvWrite(pktStream <-chan av.Packet, code string, codecs []av.CodecData) {
	ffw := rtmpflvwriter.NewRtmpFlvWriter(pktStream, code, codecs, rfm)
	rfm.rfms.Store(code, ffw)
}

func (rfm *rtmpFlvManager) StopWrite(code string) {
	v, ok := rfm.rfms.Load(code)
	if ok {
		ffw := v.(*rtmpflvwriter.RtmpFlvWriter)
		ffw.StopWrite()
	}
}

func (rfm *rtmpFlvManager) UpdateFFWS(code string, rfw *rtmpflvwriter.RtmpFlvWriter) {
	_, ok := rfm.rfms.LoadAndDelete(code)
	if ok {
		rfm.rfms.Store(code, rfw)
	}
}
