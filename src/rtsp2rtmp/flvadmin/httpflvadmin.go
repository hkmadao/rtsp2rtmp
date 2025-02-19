package flvadmin

import (
	"io"
	"sync"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/httpflvmanage"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/ext/live"
)

var hfas *HttpFlvAdmin

type HttpFlvAdmin struct {
	hfms sync.Map
}

func init() {
	hfas = &HttpFlvAdmin{}
}

func GetSingleHttpFlvAdmin() *HttpFlvAdmin {
	return hfas
}

func (hfa *HttpFlvAdmin) AddHttpFlvManager(
	pktStream <-chan av.Packet,
	code string,
	codecs []av.CodecData,
) {
	hfm := httpflvmanage.NewHttpFlvManager(pktStream, code, codecs)
	hfa.hfms.Store(code, hfm)
	go func() {
		<-hfm.GetDone()
		hfa.hfms.Delete(code)
	}()
}

func (hfa *HttpFlvAdmin) StopWrite(code string) {
	v, ok := hfa.hfms.Load(code)
	if ok {
		ffw := v.(*httpflvmanage.HttpFlvManager)
		ffw.StopWrite()
	}
}

func (hfa *HttpFlvAdmin) StartWrite(code string) {
	v, ok := hfa.hfms.Load(code)
	if ok {
		ffw := v.(*httpflvmanage.HttpFlvManager)
		ffw.StopWrite()
		hfa.AddHttpFlvManager(ffw.GetPktStream(), code, ffw.GetCodecs())
	}
}

// 添加播放者
func (hfa *HttpFlvAdmin) AddHttpFlvPlayer(
	playerDone <-chan int,
	pulseInterval time.Duration,
	code string,
	writer io.Writer,
) (<-chan int, *common.Rtmp2FlvCustomError) {
	v, b := hfa.hfms.Load(code)
	if b {
		hfm := v.(*httpflvmanage.HttpFlvManager)
		flvPlayerDone, err := hfm.AddHttpFlvPlayer(playerDone, pulseInterval, writer)
		if err != nil {
			return flvPlayerDone, common.InternalError(err)
		}
		return flvPlayerDone, nil
	}
	return nil, common.CustomError("camera no connection")
}

// 更新sps、pps等信息
func (hfa *HttpFlvAdmin) UpdateCodecs(code string, codecs []av.CodecData) {
	rfw, ok := hfa.hfms.Load(code)
	if ok {
		rfw := rfw.(*httpflvmanage.HttpFlvManager)
		rfw.SetCodecs(codecs)
	}
}

func (hfa *HttpFlvAdmin) GetLiveInfo(code string) (*live.LiveMediaInfo, error) {
	liveMediaInfo := &live.LiveMediaInfo{
		HasAudio:     false,
		OnlineStatus: false,
		AnchorName:   code,
	}
	rfw, ok := hfa.hfms.Load(code)
	if ok {
		rfw := rfw.(*httpflvmanage.HttpFlvManager)

		liveMediaInfo.HasAudio = hasAudio(rfw.GetCodecs())
		liveMediaInfo.OnlineStatus = true
	}
	return liveMediaInfo, nil
}

func hasAudio(streams []av.CodecData) bool {
	for _, stream := range streams {
		if stream.Type().IsAudio() {
			return true
		}
	}
	return false
}
