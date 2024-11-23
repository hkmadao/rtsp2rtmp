package rtspclient

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvmanage"
)

type IRtspClientManager interface {
	Load(key interface{}) (interface{}, bool)
	Store(key, value interface{})
	Delete(key interface{})
}

type RtspClient struct {
	code         string
	codecs       []av.CodecData
	connDone     <-chan int
	pktStream    <-chan av.Packet
	ffmPktStream <-chan av.Packet
	hfmPktStream <-chan av.Packet
	rfmPktStream <-chan av.Packet
	ircm         IRtspClientManager
}

func NewRtspClient(connDone <-chan int, pktStream <-chan av.Packet, code string, codecs []av.CodecData, ircm IRtspClientManager) *RtspClient {
	r := &RtspClient{
		connDone:     connDone,
		pktStream:    pktStream,
		code:         code,
		codecs:       codecs,
		ffmPktStream: make(chan av.Packet, 1024),
		hfmPktStream: make(chan av.Packet, 1024),
		rfmPktStream: make(chan av.Packet, 1024),
		ircm:         ircm,
	}
	r.pktTransfer()
	return r
}

func (r *RtspClient) Done() {
	<-r.connDone
}

func (r *RtspClient) pktTransfer() {
	ffmPktStream, hfmPktStream, rfmPktStream := tee(r.connDone, r.pktStream)
	r.ffmPktStream = ffmPktStream
	r.hfmPktStream = hfmPktStream
	r.rfmPktStream = rfmPktStream
	logs.Debug("publisher [%s] create customer", r.code)
	flvmanage.GetSingleFileFlvManager().FlvWrite(r.ffmPktStream, r.code, r.codecs)
	flvmanage.GetSingleHttpflvAdmin().AddHttpFlvManager(r.hfmPktStream, r.code, r.codecs)
	flvmanage.GetSingleRtmpFlvManager().FlvWrite(r.rfmPktStream, r.code, r.codecs)
}

func tee(done <-chan int, in <-chan av.Packet) (<-chan av.Packet, <-chan av.Packet, <-chan av.Packet) {
	//设置缓冲，调节前后速率
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
