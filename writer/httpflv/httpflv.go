package httpflv

import (
	"net/http"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
	"github.com/yumrano/rtsp2rtmp/rlog"
)

var hms map[string]*HttpFlvManager

func init() {
	hms = make(map[string]*HttpFlvManager)
}

func AddHttpFlvPlayer(done <-chan interface{}, code string, writer http.ResponseWriter) <-chan int {
	heartChan := make(chan int)
	sessionId := time.Now().Format("20060102150405")
	fw := &HttpFlvWriter{
		sessionId: sessionId,
		writer:    writer,
		heartChan: heartChan,
		codecs:    hms[code].codecs,
		code:      code,
	}
	hms[code].fws[sessionId] = fw
	return heartChan
}

type HttpFlvManager struct {
	codecs []av.CodecData
	fws    map[string]*HttpFlvWriter
}

func NewHttpFlvManager() *HttpFlvManager {
	hm := &HttpFlvManager{}
	return hm
}

func (fm *HttpFlvManager) codec(code string, codecs []av.CodecData) {
	fm.codecs = codecs
	fm.fws = make(map[string]*HttpFlvWriter)
	hms[code] = fm
}

//Write extends to writer.Writer
func (fm *HttpFlvManager) FlvWrite(code string, codecs []av.CodecData, done <-chan interface{}, pchan <-chan av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			rlog.Log.Printf("HttpFlvManager FlvWrite panic %v", r)
		}
	}()
	fm.codec(code, codecs)
	for {
		select {
		case <-done:
			return
		case pkt := <-pchan:
			deleteKeys := make([]string, 2)
			for _, fw := range fm.fws {
				if fw.IsClose() {
					deleteKeys = append(deleteKeys, fw.sessionId)
				}
				go fw.HttpWrite(pkt)
			}
			for _, sessionId := range deleteKeys {
				delete(fm.fws, sessionId)
			}
		}
	}
}

type HttpFlvWriter struct {
	sessionId string
	code      string
	start     bool
	writer    http.ResponseWriter
	codecs    []av.CodecData
	muxer     *flv.Muxer
	close     bool
	done      <-chan interface{}
	heartChan chan<- int
}

func (fw *HttpFlvWriter) IsClose() bool {
	return fw.close
}

func (fw *HttpFlvWriter) HttpWrite(pkt av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			rlog.Log.Printf("httpWrite panic : %v", r)
		}
	}()
	if fw.start {
		if err := fw.muxer.WritePacket(pkt); err != nil {
			rlog.Log.Printf("writer packet to httpflv error : %v\n", err)
			close(fw.heartChan)
			fw.close = true
			return
		}
		return
	}
	if pkt.IsKeyFrame {
		muxer := flv.NewMuxer(fw)
		fw.muxer = muxer
		err := fw.muxer.WriteHeader(fw.codecs)
		if err != nil {
			rlog.Log.Printf("writer header to httpflv error : %v\n", err)
			close(fw.heartChan)
			fw.close = true
			return
		}
		fw.start = true
		if err := fw.muxer.WritePacket(pkt); err != nil {
			rlog.Log.Printf("writer packet to httpflv error : %v\n", err)
		}
	}

}

//Write extends to io.Writer
func (fw *HttpFlvWriter) Write(p []byte) (n int, err error) {
	n, err = fw.writer.Write(p)
	if err != nil {
		rlog.Log.Println("write httpflv error :", err)
	}
	for {
		select {
		case fw.heartChan <- 1:
			return
		case <-time.After(1 * time.Millisecond):
			return
		case <-fw.done:
			return
		}
	}
}
