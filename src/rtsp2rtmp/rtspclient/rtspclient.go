package rtspclient

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
)

type RtspClient struct {
	code         string
	codecs       []av.CodecData
	connDone     <-chan int
	done         chan int
	pktStream    <-chan av.Packet
	ffmPktStream <-chan av.Packet
	hfmPktStream <-chan av.Packet
	rfmPktStream <-chan av.Packet
}

func NewRtspClient(connDone <-chan int, pktStream <-chan av.Packet, code string, codecs []av.CodecData) *RtspClient {
	done := make(chan int)
	rc := &RtspClient{
		connDone:     connDone,
		done:         done,
		pktStream:    pktStream,
		code:         code,
		codecs:       codecs,
		ffmPktStream: make(chan av.Packet, 1024),
		hfmPktStream: make(chan av.Packet, 1024),
		rfmPktStream: make(chan av.Packet, 1024),
	}
	rc.pktTransfer()
	return rc
}

func (rtspClient *RtspClient) Done() {
	<-rtspClient.connDone
}

// 主动关闭
func (rtspClient *RtspClient) Close() {
	close(rtspClient.done)
}

// 更新sps、pps等信息
func (rtspClient *RtspClient) UpdateCodecs(codecs []av.CodecData) {
	rtspClient.codecs = codecs
	logs.Warn("RtspClient: %s update codecs", rtspClient.code)
	flvadmin.GetSingleFileFlvAdmin().UpdateCodecs(rtspClient.code, codecs)
	flvadmin.GetSingleHttpFlvAdmin().UpdateCodecs(rtspClient.code, codecs)
	flvadmin.GetSingleRtmpFlvAdmin().UpdateCodecs(rtspClient.code, codecs)
}

func (rtspClient *RtspClient) pktTransfer() {
	done := utils.OrDoneInt(rtspClient.done, rtspClient.connDone)
	ffmPktStream, hfmPktStream, rfmPktStream := tee(done, rtspClient.pktStream)
	rtspClient.ffmPktStream = ffmPktStream
	rtspClient.hfmPktStream = hfmPktStream
	rtspClient.rfmPktStream = rfmPktStream
	logs.Info("publisher [%s] create customer", rtspClient.code)
	flvadmin.GetSingleFileFlvAdmin().FlvWrite(rtspClient.ffmPktStream, rtspClient.code, rtspClient.codecs)
	flvadmin.GetSingleHttpFlvAdmin().AddHttpFlvManager(rtspClient.hfmPktStream, rtspClient.code, rtspClient.codecs)
	flvadmin.GetSingleRtmpFlvAdmin().FlvWrite(rtspClient.rfmPktStream, rtspClient.code, rtspClient.codecs)
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
