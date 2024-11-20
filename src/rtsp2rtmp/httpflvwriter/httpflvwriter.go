package httpflvwriter

import (
	"io"
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/models"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
)

type IHttpFlvManager interface {
	DeleteHFW(sesessionId int64)
}

type HttpFlvWriter struct {
	sessionId         int64
	pulseInterval     time.Duration
	pktStream         chan av.Packet
	code              string
	codecs            []av.CodecData
	start             bool
	writer            io.Writer
	muxer             *flv.Muxer
	done              chan int
	httpflvManageDone <-chan int //来自管理者的关闭
	playerDone        <-chan int //来自播放者的关闭
	ihfm              IHttpFlvManager
}

func (hfw *HttpFlvWriter) GetDone() <-chan int {
	return utils.OrDoneInt(hfw.done, hfw.httpflvManageDone)
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

func (hfw *HttpFlvWriter) GetSessionId() int64 {
	return hfw.sessionId
}

func NewHttpFlvWriter(
	httpflvManageDone <-chan int,
	playerDone <-chan int,
	pulseInterval time.Duration,
	pktStream chan av.Packet,
	code string,
	codecs []av.CodecData,
	writer io.Writer,
	sessionId int64,
	ihfm IHttpFlvManager,
) *HttpFlvWriter {
	hfw := &HttpFlvWriter{
		sessionId:         sessionId,
		pulseInterval:     pulseInterval,
		pktStream:         pktStream,
		code:              code,
		codecs:            codecs,
		writer:            writer,
		done:              make(chan int),
		playerDone:        playerDone,
		httpflvManageDone: httpflvManageDone,
		ihfm:              ihfm,
		start:             false,
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
		go func() {
			select {
			case <-hfw.GetDone():
				return
			case <-hfw.pktStream:
				return
			}
		}()
		return hfw
	}

	go hfw.httpWrite()
	return hfw
}

func (hfw *HttpFlvWriter) httpWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	ticker := time.NewTicker(hfw.pulseInterval)
	defer func() {
		close(hfw.done)
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	pktStream := utils.OrDonePacket(hfw.playerDone, hfw.pktStream)
	for {
		select {
		case <-ticker.C:
			return
		case pkt, ok := <-pktStream:
			if !ok {
				return
			}
			if err := hfw.writerPacket(pkt); err != nil {
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
	logs.Debug("ingrore package")
	return nil
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
