package fileflvreader

import (
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
)

type FileFlvReader struct {
	done         chan int
	fgDoneClose  bool
	code         string
	codecs       []av.CodecData
	fd           *os.File
	fileName     string
	fullFileName string
	muxer        *flv.Demuxer
	startTime    time.Time
	endTime      time.Time
	httpWrite    *MyHttpWriter
	seekSecond   uint64 //uint second
	fgStart      bool
	fgOffSetTime bool
	offSetTime   time.Duration // realTime = pkt.Time - offSetTime, offsetTime from is first pkt.Time
	mutex        sync.Mutex
}

func (ffr *FileFlvReader) GetDone() <-chan int {
	return ffr.done
}

func (ffr *FileFlvReader) GetCode() string {
	return ffr.code
}

func (ffr *FileFlvReader) SetCodecs(codecs []av.CodecData) {
	ffr.codecs = codecs
}

func (ffr *FileFlvReader) GetCodecs() []av.CodecData {
	return ffr.codecs
}

func (ffr *FileFlvReader) SetSeekSecond(seekSecond uint64) {
	if seekSecond > ffr.seekSecond {
		ffr.seekSecond = seekSecond
	}
}

func NewFileFlvReader(
	seekSecond uint64,
	writer io.Writer,
	fileName string,
) *FileFlvReader {
	myHttpWriter := newMyHttpWriter(writer)
	ffr := &FileFlvReader{
		fgDoneClose:  false,
		done:         make(chan int),
		httpWrite:    myHttpWriter,
		fileName:     fileName,
		codecs:       make([]av.CodecData, 0),
		seekSecond:   seekSecond,
		fgStart:      false,
		fgOffSetTime: false,
		offSetTime:   0,
	}
	go ffr.flvRead()
	return ffr
}

func (ffr *FileFlvReader) StopRead() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		ffr.CloseDone()
	}()
}

func (ffr *FileFlvReader) CloseDone() {
	ffr.mutex.Lock()
	if !ffr.fgDoneClose {
		ffr.fgDoneClose = true
		close(ffr.done)
	}
	ffr.mutex.Unlock()
}

func (ffr *FileFlvReader) Read(p []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	n, err = ffr.fd.Read(p)
	if err != nil {
		logs.Error("Read file error : %v", err)
	}

	return
}

func (ffr *FileFlvReader) openFlvFile() error {
	fullFileName := getFileFlvPath() + "/" + ffr.fileName
	fd, err := os.OpenFile(fullFileName, os.O_RDWR, 0644)
	if err != nil {
		logs.Error("open file: %s error : %v", fullFileName, err)
		return err
	}
	ffr.fd = fd
	ffr.fullFileName = fullFileName
	return nil
}

func (ffr *FileFlvReader) flvRead() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	if err := ffr.openFlvFile(); err != nil {
		logs.Error("open file flv error : %v", err)
		return
	}
	defer func() {
		ffr.endTime = time.Now()
		ffr.fd.Close()

		ffr.CloseDone()
	}()

	ffr.startTime = time.Now()
	httpWriteMuxer := flv.NewMuxer(ffr.httpWrite)
	demuxer := flv.NewDemuxer(ffr)
	ffr.muxer = demuxer
	codecs, err := demuxer.Streams()
	if err != nil {
		logs.Error("read codecs from flv file error : %v", err)
		return
	}
	ffr.codecs = codecs
	httpWriteMuxer.WriteHeader(codecs)

	if ffr.seekSecond > 0 {
		for {
			pkt, err := demuxer.ReadPacket()
			if err != nil {
				logs.Error("read file %s ReadPacket error : %v", ffr.fileName, err)
				break
			}
			if !ffr.fgOffSetTime {
				ffr.fgOffSetTime = true
				ffr.offSetTime = pkt.Time
			}
			if (pkt.Time - ffr.offSetTime) >= time.Duration(ffr.seekSecond)*time.Second {
				break
			}
		}
	}

	timeStart := time.Now()

Loop:
	for {
		pkt, err := demuxer.ReadPacket()
		if err != nil {
			logs.Error("read file %s ReadPacket error : %v", ffr.fileName, err)
			break
		}
		if !ffr.fgOffSetTime {
			ffr.fgOffSetTime = true
			ffr.offSetTime = pkt.Time
		}
		if !ffr.fgStart {
			if !pkt.IsKeyFrame {
				continue
			}
			ffr.fgStart = true
		}
		err = httpWriteMuxer.WritePacket(pkt)
		if err != nil {
			logs.Error("read file %s WritePacket error : %v", ffr.fileName, err)
			break
		}

		sinceTime := time.Since(timeStart) + time.Duration(ffr.seekSecond)*time.Second
		if (pkt.Time - ffr.offSetTime) > (sinceTime + 5*time.Minute) {
			select {
			case <-ffr.done:
				break Loop
			case <-time.NewTicker(5 * time.Second).C:
			}
		}
		select {
		case <-ffr.done:
			break Loop
		default:
		}
	}
}

func getFileFlvPath() string {
	fileFlvPath, err := config.String("server.fileflv.path")
	if err != nil {
		logs.Error("get fileflv path error :", err)
		return ""
	}
	return fileFlvPath
}

type MyHttpWriter struct {
	count         int
	writeChunkLen uint64
	writer        io.Writer
}

func newMyHttpWriter(writer io.Writer) *MyHttpWriter {
	return &MyHttpWriter{
		count:         0,
		writeChunkLen: 0,
		writer:        writer,
	}
}

// Write implements io.Writer.
func (w *MyHttpWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	return
}

func FlvDurationReadUntilErr(fileName string) (n int, err error) {
	fullFileName := getFileFlvPath() + "/" + fileName
	fd, err := os.OpenFile(fullFileName, os.O_RDWR, 0644)
	if err != nil {
		logs.Error("open file: %s error : %v", fullFileName, err)
		return
	}
	defer func() {
		fd.Close()
	}()
	demuxer := flv.NewDemuxer(fd)
	totalTime := time.Duration(0)
	demuxer.Streams()
	fgStart := false
	offSetTime := time.Duration(0)
	for {
		pkt, readErr := demuxer.ReadPacket()
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			logs.Error("read file %s ReadPacket error : %v", fileName, readErr)
			break
		}
		if !fgStart {
			fgStart = true
			offSetTime = pkt.Time
		}

		totalTime = pkt.Time
	}
	totalTime = totalTime - offSetTime
	n = int(totalTime / time.Millisecond)
	return
}

func FlvFileExists(fileName string) bool {
	fullFileName := getFileFlvPath() + "/" + fileName
	_, err := os.Stat(fullFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
