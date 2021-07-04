package httpflvwriter

import (
	"io"
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/models"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
)

type IHttpFlvManager interface {
	DeleteHFW(sesessionId int64)
}

type HttpFlvWriter struct {
	sessionId             int64
	pulseInterval         time.Duration
	pktStream             chan av.Packet
	code                  string
	codecs                []av.CodecData
	start                 bool
	writer                io.Writer
	muxer                 *flv.Muxer
	playerDone            <-chan interface{} //来自播放者的关闭
	heartbeatStream       chan interface{}   //心跳包通道
	playerHeartbeatStream <-chan interface{}
	selfHeartbeatStream   <-chan interface{}
	ihfm                  IHttpFlvManager
}

func (hfw *HttpFlvWriter) SetCode(code string) {
	hfw.code = code
}

func (hfw *HttpFlvWriter) SetCodecs(codecs []av.CodecData) {
	hfw.codecs = codecs
}

func (hfw *HttpFlvWriter) GetPktStream() chan<- av.Packet {
	return hfw.pktStream
}

// func (hfw *HttpFlvWriter) GetReadPktStream() <-chan av.Packet {
// 	return hfw.pktStream
// }

func (hfw *HttpFlvWriter) GetSessionId() int64 {
	return hfw.sessionId
}

func NewHttpFlvWriter(
	playerDone <-chan interface{},
	pulseInterval time.Duration,
	pktStream chan av.Packet,
	code string,
	codecs []av.CodecData,
	writer io.Writer,
	sessionId int64,
	ihfm IHttpFlvManager,
) *HttpFlvWriter {
	heartbeatStream := make(chan interface{})
	playerHeartbeatStream, selfHeartbeatStream := utils.Tee(playerDone, heartbeatStream)
	hfw := &HttpFlvWriter{
		sessionId:             sessionId,
		pulseInterval:         pulseInterval,
		pktStream:             pktStream,
		code:                  code,
		codecs:                codecs,
		writer:                writer,
		playerDone:            playerDone,
		heartbeatStream:       heartbeatStream,
		playerHeartbeatStream: playerHeartbeatStream,
		selfHeartbeatStream:   selfHeartbeatStream,
		ihfm:                  ihfm,
		start:                 false,
	}

	camera, err := models.CameraSelectOne(models.Camera{Code: code})
	if err != nil {
		logs.Error("query camera error : %v", err)
		return hfw
	}
	if camera.OnlineStatus != 1 {
		return hfw
	}
	if camera.Live != 1 {
		return hfw
	}

	go hfw.httpWrite()
	go hfw.monitor()
	return hfw
}

func (hfw *HttpFlvWriter) GetPlayerHeartbeatStream() <-chan interface{} {
	return hfw.playerHeartbeatStream
}

func (hfw *HttpFlvWriter) httpWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	pulse := time.NewTicker(hfw.pulseInterval).C

	sendPulse := func() {
		select {
		case hfw.heartbeatStream <- struct{}{}:
		default:
		}
	}

	sendPulse()
	done := make(chan interface{})
	go func(done <-chan interface{}) {
		defer func() {
			close(hfw.heartbeatStream)
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		writer := hfw.writer.(gin.ResponseWriter)
		for {
			select {
			case <-writer.CloseNotify():
				return
			case <-pulse:
				sendPulse()
			case <-done:
				return
			case <-hfw.playerDone:
				return
			}
		}
	}(done)

	ticker := time.NewTicker(hfw.pulseInterval)
	pktStream := utils.OrDonePacket(hfw.playerDone, hfw.pktStream)
	for {
		select {
		case <-ticker.C:
			close(done)
			return
		case pkt := <-pktStream:
			if err := hfw.writerPacket(pkt); err != nil {
				close(done)
				return
			}
			ticker.Reset(hfw.pulseInterval)
		}
	}

}
func (hfw *HttpFlvWriter) writerPacket(pkt av.Packet) error {
	if hfw.start {
		if err := hfw.muxer.WritePacket(pkt); err != nil {
			logs.Error("writer packet to httpflv error : %v", err)
			return err
		}
		// logs.Debug("httpWrite")
		return nil
	}
	if pkt.IsKeyFrame {
		muxer := flv.NewMuxer(hfw)
		hfw.muxer = muxer
		err := hfw.muxer.WriteHeader(hfw.codecs)
		if err != nil {
			logs.Error("writer header to httpflv error : %v", err)
			return err
		}
		hfw.start = true
		if err := hfw.muxer.WritePacket(pkt); err != nil {
			logs.Error("writer packet to httpflv error : %v", err)
			return err
		}
		return nil
	}
	return nil
}

func (hfw *HttpFlvWriter) monitor() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		select {
		case _, ok := <-hfw.selfHeartbeatStream:
			if ok {
				logs.Debug("heartbeat")
				continue
			}
			hfw.ihfm.DeleteHFW(hfw.sessionId)
			return
		case <-time.After(2 * hfw.pulseInterval):
			//time out
			hfw.ihfm.DeleteHFW(hfw.sessionId)
			return
		}
	}
}

//Write extends to io.Writer
func (hfw *HttpFlvWriter) Write(p []byte) (n int, err error) {
	start := time.Now()
	defer func() {
		logs.Debug(time.Since(start))
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	// logs.Debug("write pkt")
	n, err = hfw.writer.Write(p)
	if err != nil {
		logs.Error("write httpflv error : %v", err)
	}
	return
}
