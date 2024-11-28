package flvadmin

import (
	"sync"

	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/rtmpflvwriter"
)

type RtmpFlvAdmin struct {
	rfms sync.Map
}

var rfmInstance *RtmpFlvAdmin

func init() {
	rfmInstance = &RtmpFlvAdmin{}
}

func GetSingleRtmpFlvAdmin() *RtmpFlvAdmin {
	return rfmInstance
}

func (rfm *RtmpFlvAdmin) FlvWrite(pktStream <-chan av.Packet, code string, codecs []av.CodecData) {
	rfw := rtmpflvwriter.NewRtmpFlvWriter(pktStream, code, codecs, rfm)
	rfm.rfms.Store(code, rfw)
}

func (rfm *RtmpFlvAdmin) StartWrite(code string) {
	v, ok := rfm.rfms.Load(code)
	if ok {
		rfw := v.(*rtmpflvwriter.RtmpFlvWriter)
		rfw.StopWrite()
		rfm.FlvWrite(rfw.GetPktStream(), code, rfw.GetCodecs())
	}
}

func (rfm *RtmpFlvAdmin) StopWrite(code string) {
	v, ok := rfm.rfms.Load(code)
	if ok {
		rfw := v.(*rtmpflvwriter.RtmpFlvWriter)
		rfw.StopWrite()
	}
}

func (rfm *RtmpFlvAdmin) UpdateFFWS(code string, rfw *rtmpflvwriter.RtmpFlvWriter) {
	_, ok := rfm.rfms.LoadAndDelete(code)
	if ok {
		rfm.rfms.Store(code, rfw)
	}
}

//更新sps、pps等信息
func (rfm *RtmpFlvAdmin) UpdateCodecs(code string, codecs []av.CodecData) {
	rfw, ok := rfm.rfms.Load(code)
	if ok {
		rfw := rfw.(*rtmpflvwriter.RtmpFlvWriter)
		rfw.SetCodecs(codecs)
		//sps、pps更新后，重新建立连接
		// logs.Warn("RtmpFlvAdmin: %s codecs change, restart RtmpFlvWriter", code)
		//这里只需要stop就可以，内部会重连
		// rfw.StopWrite()
	}
}
