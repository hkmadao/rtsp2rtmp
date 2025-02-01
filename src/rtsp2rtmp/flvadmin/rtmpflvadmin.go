package flvadmin

import (
	"fmt"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/rtmpflvwriter"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
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
	condition := common.GetEqualCondition("code", code)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("FlvWrite found camera: %s error: %v, do painc", code, err)
		panic(fmt.Sprintf("FlvWrite found camera: %s error: %v", code, err))
	}
	rfw := rtmpflvwriter.NewRtmpFlvWriter(!camera.FgPassive, pktStream, code, codecs, rfm)
	rfm.rfms.Store(code, rfw)
}

func (rfm *RtmpFlvAdmin) StartWrite(code string, needPushRtmp bool) {
	v, ok := rfm.rfms.Load(code)
	if ok {
		rfw := v.(*rtmpflvwriter.RtmpFlvWriter)
		rfw.StopWrite()
		rfwNew := rtmpflvwriter.NewRtmpFlvWriter(needPushRtmp, rfw.GetPktStream(), code, rfw.GetCodecs(), rfm)
		rfm.rfms.Store(code, rfwNew)
	}
}

func (rfm *RtmpFlvAdmin) ReConntion(code string) {
	v, ok := rfm.rfms.Load(code)
	if ok {
		rfw := v.(*rtmpflvwriter.RtmpFlvWriter)
		rfw.StopWrite()
		rfwNew := rtmpflvwriter.NewRtmpFlvWriter(false, rfw.GetPktStream(), code, rfw.GetCodecs(), rfm)
		rfm.rfms.Store(code, rfwNew)
	}
}

func (rfm *RtmpFlvAdmin) RemoteStartWrite(code string) {
	v, ok := rfm.rfms.Load(code)
	if ok {
		rfw := v.(*rtmpflvwriter.RtmpFlvWriter)
		if !rfw.GetNeedPushRtmp() {
			rfwNew := rtmpflvwriter.NewRtmpFlvWriter(true, rfw.GetPktStream(), code, rfw.GetCodecs(), rfm)
			rfm.rfms.Store(code, rfwNew)
		}
	}
}

func (rfm *RtmpFlvAdmin) RemoteStopWrite(code string) {
	v, ok := rfm.rfms.Load(code)
	if ok {
		rfw := v.(*rtmpflvwriter.RtmpFlvWriter)
		rfw.StopWrite()
		rfwNew := rtmpflvwriter.NewRtmpFlvWriter(false, rfw.GetPktStream(), code, rfw.GetCodecs(), rfm)
		rfm.rfms.Store(code, rfwNew)
	}
}

func (rfm *RtmpFlvAdmin) UpdateFFWS(code string, rfw *rtmpflvwriter.RtmpFlvWriter) {
	_, ok := rfm.rfms.LoadAndDelete(code)
	if ok {
		rfm.rfms.Store(code, rfw)
	}
}

// 更新sps、pps等信息
func (rfm *RtmpFlvAdmin) UpdateCodecs(code string, codecs []av.CodecData) {
	rfw, ok := rfm.rfms.Load(code)
	if ok {
		rfw := rfw.(*rtmpflvwriter.RtmpFlvWriter)
		rfw.SetCodecs(codecs)
	}
}
