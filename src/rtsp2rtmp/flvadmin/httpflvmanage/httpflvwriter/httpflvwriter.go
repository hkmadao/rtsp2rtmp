package httpflvwriter

import (
	"io"
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
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

func (hfw *HttpFlvWriter) GetCode() string {
	return hfw.code
}

func (hfw *HttpFlvWriter) SetCodecs(codecs []av.CodecData) {
	logs.Warn("HttpFlvWriter: %s update codecs", hfw.code)
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
		hfw.ihfm.DeleteHFW(hfw.sessionId)
		close(hfw.done)
	}()
	pktStream := utils.OrDonePacket(hfw.playerDone, hfw.pktStream)

	notBlockStream := hfw.changeNotBlockStream(pktStream)
	timeNow := time.Now().Local()
	for {
		select {
		case <-ticker.C:
			return
		case pkt, ok := <-notBlockStream:
			if !ok {
				return
			}
			if err := hfw.writerPacket(pkt, &timeNow); err != nil {
				return
			}
			ticker.Reset(hfw.pulseInterval)
		}
	}

}

func (hfw *HttpFlvWriter) changeNotBlockStream(pktStream <-chan av.Packet) <-chan av.Packet {
	notBlockStream := make(chan av.Packet, 1024)
	go func() {
		defer close(notBlockStream)
		for {
			select {
			case pkt, ok := <-pktStream:
				if !ok {
					return
				}
				select {
				case notBlockStream <- pkt:
				default:
					logs.Error("camera: %s, session: %d lose package", hfw.code, hfw.sessionId)
				}
			}

		}
	}()
	return notBlockStream
}

func (hfw *HttpFlvWriter) writerPacket(pkt av.Packet, templateTime *time.Time) error {
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
	if time.Now().Local().After(templateTime.Add(1 * time.Minute)) {
		*templateTime = time.Now().Local()
		logs.Error("HttpFlvWriter ingrore package: %s", hfw.code)
	}
	return nil
}

// Write extends to io.Writer
func (hfw *HttpFlvWriter) Write(p []byte) (n int, err error) {
	// start := time.Now()
	defer func() {
		// logs.Debug(time.Since(start))
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
