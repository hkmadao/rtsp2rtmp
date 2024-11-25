package flvadmin

import (
	"sync"

	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/fileflvmanager"
)

type FileFlvAdmin struct {
	ffws sync.Map
}

var ffmInstance *FileFlvAdmin

func init() {
	ffmInstance = &FileFlvAdmin{}
}

func GetSingleFileFlvAdmin() *FileFlvAdmin {
	return ffmInstance
}

func (ffm *FileFlvAdmin) FlvWrite(pktStream <-chan av.Packet, code string, codecs []av.CodecData) {
	ffw := fileflvmanager.NewFileFlvManager(pktStream, code, codecs)
	ffm.ffws.Store(code, ffw)
}

func (ffm *FileFlvAdmin) StopWrite(code string) {
	v, ok := ffm.ffws.Load(code)
	if ok {
		ffw := v.(*fileflvmanager.FileFlvManager)
		ffw.StopWrite()
	}
}

func (ffm *FileFlvAdmin) StartWrite(code string) {
	v, ok := ffm.ffws.Load(code)
	if ok {
		ffw := v.(*fileflvmanager.FileFlvManager)
		ffw.StopWrite()
		ffm.FlvWrite(ffw.GetPktStream(), code, ffw.GetCodecs())
	}
}

func (ffm *FileFlvAdmin) UpdateFFWS(code string, ffw *fileflvmanager.FileFlvManager) {
	_, ok := ffm.ffws.LoadAndDelete(code)
	if ok {
		ffm.ffws.Store(code, ffw)
	}
}

//更新sps、pps等信息
func (ffm *FileFlvAdmin) UpdateCodecs(code string, codecs []av.CodecData) {
	rfw, ok := ffm.ffws.Load(code)
	if ok {
		rfw := rfw.(*fileflvmanager.FileFlvManager)
		rfw.SetCodecs(codecs)
	}
}
