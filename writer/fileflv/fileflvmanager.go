package fileflv

import (
	"runtime/debug"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
)

type FileFlvManager struct {
	done      <-chan interface{}
	pktStream <-chan av.Packet
	code      string
	codecs    []av.CodecData
}

func NewFileFlvManager(done <-chan interface{}, pktStream <-chan av.Packet, code string, codecs []av.CodecData) *FileFlvManager {
	ffm := &FileFlvManager{
		done:      done,
		pktStream: pktStream,
		code:      code,
		codecs:    codecs,
	}
	go ffm.flvWrite()
	return ffm
}

func (ffm *FileFlvManager) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	NewFileFlvWriter(ffm.done, ffm.pktStream, ffm.code, ffm.codecs)
}
