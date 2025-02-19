package rtmppublisher

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
)

type RtmpServer interface {
	Load(key interface{}) (interface{}, bool)
	Store(key, value interface{})
	Delete(key interface{})
}

type Publisher struct {
	code         string
	codecs       []av.CodecData
	done         chan int
	connDone     <-chan int
	pktStream    <-chan av.Packet
	ffmPktStream <-chan av.Packet
	hfmPktStream <-chan av.Packet
	rfmPktStream <-chan av.Packet
	rtmpserver   RtmpServer
}

func NewPublisher(connDone <-chan int, pktStream <-chan av.Packet, code string, codecs []av.CodecData, rs RtmpServer) *Publisher {
	r := &Publisher{
		connDone:     connDone,
		pktStream:    pktStream,
		code:         code,
		codecs:       codecs,
		ffmPktStream: make(chan av.Packet, 1024),
		hfmPktStream: make(chan av.Packet, 1024),
		rtmpserver:   rs,
	}
	r.pktTransfer()
	return r
}

func (r *Publisher) Done() {
	<-r.connDone
}

func (rtmpClient *Publisher) pktTransfer() {
	done := utils.OrDoneInt(rtmpClient.done, rtmpClient.connDone)
	ffmPktStream, hfmPktStream, rfmPktStream := tee(done, rtmpClient.pktStream)
	rtmpClient.ffmPktStream = ffmPktStream
	rtmpClient.hfmPktStream = hfmPktStream
	rtmpClient.rfmPktStream = rfmPktStream
	logs.Debug("publisher [%s] create customer", rtmpClient.code)
	flvadmin.GetSingleFileFlvAdmin().FlvWrite(rtmpClient.ffmPktStream, rtmpClient.code, rtmpClient.codecs)
	flvadmin.GetSingleHttpFlvAdmin().AddHttpFlvManager(rtmpClient.hfmPktStream, rtmpClient.code, rtmpClient.codecs)
	flvadmin.GetSingleRtmpFlvAdmin().FlvWrite(rtmpClient.rfmPktStream, rtmpClient.code, rtmpClient.codecs)
}

func tee(done <-chan int, in <-chan av.Packet) (<-chan av.Packet, <-chan av.Packet, <-chan av.Packet) {
	//设置缓冲
	out1 := make(chan av.Packet, 1024)
	out2 := make(chan av.Packet, 1024)
	out3 := make(chan av.Packet, 1024)
	go func() {
		defer close(out1)
		defer close(out2)
		defer close(out3)
		for val := range in {
			var out1, out2, out3 = out1, out2, out3 // 私有变量覆盖
			for i := 0; i < 3; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					// logs.Debug("FileFlvManager write success")
					out1 = nil // 置空阻塞机制完成select轮询
				case out2 <- val:
					// logs.Debug("HttpflvAdmin write success")
					out2 = nil
				case out3 <- val:
					// logs.Debug("RtmpFlvManager write success")
					out3 = nil
				default:
					// logs.Debug("RtspClient tee lose packet")
				}
			}
		}
	}()
	return out1, out2, out3
}
