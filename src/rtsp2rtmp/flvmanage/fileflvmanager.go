package flvmanage

import (
	"sync"

	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/fileflvwriter"
)

type fileFlvManager struct {
	ffws sync.Map
}

var ffmInstance *fileFlvManager

func init() {
	ffmInstance = &fileFlvManager{}
}

func GetSingleFileFlvManager() *fileFlvManager {
	return ffmInstance
}

func (ffm *fileFlvManager) FlvWrite(pktStream <-chan av.Packet, code string, codecs []av.CodecData) {
	ffw := fileflvwriter.NewFileFlvWriter(pktStream, code, codecs, ffm)
	ffm.ffws.Store(code, ffw)
}

func (ffm *fileFlvManager) StopWrite(code string) {
	v, ok := ffm.ffws.Load(code)
	if ok {
		ffw := v.(*fileflvwriter.FileFlvWriter)
		ffw.StopWrite()
	}
}

func (ffm *fileFlvManager) StartWrite(code string) {
	v, ok := ffm.ffws.Load(code)
	if ok {
		ffw := v.(*fileflvwriter.FileFlvWriter)
		ffw.StopWrite()
		ffm.FlvWrite(ffw.GetPktStream(), code, ffw.GetCodecs())
	}
}

func (ffm *fileFlvManager) UpdateFFWS(code string, ffw *fileflvwriter.FileFlvWriter) {
	_, ok := ffm.ffws.LoadAndDelete(code)
	if ok {
		ffm.ffws.Store(code, ffw)
	}
}
