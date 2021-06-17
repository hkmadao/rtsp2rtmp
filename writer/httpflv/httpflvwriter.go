package httpflv

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
)

type HttpFlvWriter struct {
	sessionId       int64
	done            <-chan interface{}
	pktStream       chan av.Packet
	code            string
	codecs          []av.CodecData
	start           bool
	writer          http.ResponseWriter
	muxer           *flv.Muxer
	close           bool
	playerDone      chan interface{} //来自播放者的关闭
	heartbeatStream chan interface{} //心跳包通道
	endStream       chan interface{} //结束通道，告知父进程结束
}

func NewHttpFlvWriter(done <-chan interface{}, code string, codecs []av.CodecData, writer http.ResponseWriter, sessionId int64) *HttpFlvWriter {
	playerDone := make(chan interface{})
	heartbeatStream := make(chan interface{})
	endStream := make(chan interface{})
	pktStream := make(chan av.Packet, 1*1024)
	hfw := &HttpFlvWriter{
		sessionId:       sessionId,
		done:            done,
		pktStream:       pktStream,
		code:            code,
		codecs:          codecs,
		writer:          writer,
		playerDone:      playerDone,
		heartbeatStream: heartbeatStream,
		endStream:       endStream,
		start:           false,
		close:           false,
	}
	go hfw.httpWrite()
	return hfw
}

func (hfw *HttpFlvWriter) GetEndStream() <-chan interface{} {
	return hfw.endStream
}

func (hfw *HttpFlvWriter) GetHeartbeatStream() <-chan interface{} {
	return hfw.heartbeatStream
}

func (hfw *HttpFlvWriter) GetPlayerDone() chan<- interface{} {
	return hfw.playerDone
}

func (hfw *HttpFlvWriter) GetPktStream() chan<- av.Packet {
	return hfw.pktStream
}

func (hfw *HttpFlvWriter) httpWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-hfw.done:
			close(hfw.endStream)
			return
		case <-hfw.playerDone:
			close(hfw.endStream)
			return
		case pkt := <-hfw.pktStream:
			if hfw.start {
				if err := hfw.muxer.WritePacket(pkt); err != nil {
					logs.Error("writer packet to httpflv error : %v", err)
					close(hfw.endStream)
					return
				}
				// logs.Info("start send heartbeat ")
				select {
				case hfw.heartbeatStream <- 1:
					// logs.Info("send heartbeat sucessful")
					continue
				case <-hfw.done:
					logs.Info("send heartbeat done")
					continue
				case <-hfw.playerDone:
					logs.Info("send heartbeat playerDone")
					continue
				case <-time.After(1 * time.Millisecond):
					logs.Info("send heartbeat time out")
				}
				continue
			}
			if pkt.IsKeyFrame {
				muxer := flv.NewMuxer(hfw)
				hfw.muxer = muxer
				err := hfw.muxer.WriteHeader(hfw.codecs)
				if err != nil {
					logs.Error("writer header to httpflv error : %v", err)
					close(hfw.endStream)
					return
				}
				hfw.start = true
				if err := hfw.muxer.WritePacket(pkt); err != nil {
					logs.Error("writer packet to httpflv error : %v", err)
					close(hfw.endStream)
					return
				}
				// logs.Info("start send heartbeat ")
				select {
				case hfw.heartbeatStream <- 1:
					// logs.Info("send heartbeat sucessful")
					continue
				case <-hfw.done:
					continue
				case <-hfw.playerDone:
					continue
				case <-time.After(1 * time.Millisecond):
					logs.Info("send heartbeat time out")
				}
				continue
			}
		}
	}

}

//Write extends to io.Writer
func (hfw *HttpFlvWriter) Write(p []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	n, err = hfw.writer.Write(p)
	if err != nil {
		logs.Error("write httpflv error : %v", err)
	}
	return
}
