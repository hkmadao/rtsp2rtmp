package flvadmin

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/httpflvmanage"
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

//添加播放者
func (hfa *HttpFlvAdmin) AddHttpFlvPlayer(
	playerDone <-chan int,
	pulseInterval time.Duration,
	code string,
	writer io.Writer,
) (<-chan int, error) {
	v, b := hfa.hfms.Load(code)
	if b {
		hfm := v.(*httpflvmanage.HttpFlvManager)
		return hfm.AddHttpFlvPlayer(playerDone, pulseInterval, writer)
	}
	return nil, errors.New("camera no connection")
}

//更新sps、pps等信息
func (hfa *HttpFlvAdmin) UpdateCodecs(code string, codecs []av.CodecData) {
	rfw, ok := hfa.hfms.Load(code)
	if ok {
		rfw := rfw.(*httpflvmanage.HttpFlvManager)
		rfw.SetCodecs(codecs)
	}
}
