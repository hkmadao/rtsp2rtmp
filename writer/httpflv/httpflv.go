package httpflv

import (
	"fmt"
	"log"

	"net/http"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
)

var hms map[string]*HttpFlvManager

func init() {
	hms = make(map[string]*HttpFlvManager)
}

type HttpFlvManager struct {
	codecs []av.CodecData
	fws    map[string]*FlvResponseWriter
}

func NewHttpFlvManager() *HttpFlvManager {
	hm := &HttpFlvManager{
		fws: make(map[string]*FlvResponseWriter),
	}
	return hm
}

func (fm *HttpFlvManager) codec(code string, codecs []av.CodecData) {
	fm.codecs = codecs
	hms[code] = fm
}

//Write extends to writer.Writer
func (fm *HttpFlvManager) FlvWrite(code string, codecs []av.CodecData, done <-chan interface{}, pchan <-chan av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("HttpFlvManager FlvWrite pain %v", r)
		}
	}()
	fm.codec(code, codecs)
	for {
		select {
		case <-done:
			return
		case pkt := <-pchan:
			for _, fw := range fm.fws {
				if fw.close {
					fw.done <- nil
					delete(fm.fws, fw.sessionId)
					continue
				}
				if fw.isStart {
					if err := fw.muxer.WritePacket(pkt); err != nil {
						log.Printf("writer packet to httpflv error : %v\n", err)
						if fw.errTime > 20 {
							fw.close = true
							continue
						}
						fw.errTime = fw.errTime + 1
					} else {
						fw.errTime = 0
					}
					continue
				}
				if pkt.IsKeyFrame {
					muxer := flv.NewMuxer(fw)
					fw.muxer = muxer
					err := fw.muxer.WriteHeader(fm.codecs)
					if err != nil {
						log.Printf("writer header to httpflv error : %v\n", err)
						if fw.errTime > 20 {
							fw.close = true
							continue
						}
						fw.errTime = fw.errTime + 1
					}
					fw.isStart = true
					if err := fw.muxer.WritePacket(pkt); err != nil {
						log.Printf("writer packet to httpflv error : %v\n", err)
					}
				}
			}
		}
	}
}

type FlvResponseWriter struct {
	sessionId      string
	code           string
	isStart        bool
	responseWriter http.ResponseWriter
	codecs         []av.CodecData
	muxer          *flv.Muxer
	close          bool
	errTime        int
	done           chan<- interface{}
}

//Write extends to io.Writer
func (fw *FlvResponseWriter) Write(p []byte) (n int, err error) {
	n, err = fw.responseWriter.Write(p)
	if err != nil {
		fmt.Println("write httpflv error :", err)
	}
	return
}
