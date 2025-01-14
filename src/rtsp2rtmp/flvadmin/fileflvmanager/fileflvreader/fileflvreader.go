package fileflvreader

import (
	"io"
	"os"
	"runtime/debug"
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
}

func (ffw *FileFlvReader) GetDone() <-chan int {
	return ffw.done
}

func (ffw *FileFlvReader) GetCode() string {
	return ffw.code
}

func (ffw *FileFlvReader) SetCodecs(codecs []av.CodecData) {
	ffw.codecs = codecs
}

func (ffw *FileFlvReader) GetCodecs() []av.CodecData {
	return ffw.codecs
}

func (ffw *FileFlvReader) SetSeekSecond(seekSecond uint64) {
	if seekSecond > ffw.seekSecond {
		ffw.seekSecond = seekSecond
	}
}

func NewFileFlvReader(
	seekSecond uint64,
	writer io.Writer,
	fileName string,
) *FileFlvReader {
	myHttpWriter := newMyHttpWriter(writer)
	ffw := &FileFlvReader{
		fgDoneClose: false,
		done:        make(chan int),
		httpWrite:   myHttpWriter,
		fileName:    fileName,
		codecs:      make([]av.CodecData, 0),
		seekSecond:  seekSecond,
		fgStart:     false,
	}
	go ffw.flvRead()
	return ffw
}

func (ffw *FileFlvReader) StopRead() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		ffw.fgDoneClose = true
		close(ffw.done)
	}()
}

func (ffw *FileFlvReader) TickerStopRead() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		select {
		case <-time.NewTicker(30 * time.Second).C: //等待30秒再关闭
			ffw.fgDoneClose = true
			close(ffw.done)
		case <-ffw.GetDone():
		}
	}()
}

func (ffw *FileFlvReader) Read(p []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	n, err = ffw.fd.Read(p)
	if err != nil {
		logs.Error("Read file error : %v", err)
	}

	return
}

func (ffw *FileFlvReader) openFlvFile() error {
	fullFileName := getFileFlvPath() + "/" + ffw.fileName + ".flv"
	fd, err := os.OpenFile(fullFileName, os.O_RDWR, 0644)
	if err != nil {
		logs.Error("open file: %s error : %v", fullFileName, err)
		return err
	}
	ffw.fd = fd
	ffw.fullFileName = fullFileName
	return nil
}

func (ffw *FileFlvReader) flvRead() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	if err := ffw.openFlvFile(); err != nil {
		logs.Error("open file flv error : %v", err)
		return
	}
	defer func() {
		ffw.endTime = time.Now()
		ffw.fd.Close()

		if !ffw.fgDoneClose {
			close(ffw.done)
		}
	}()

	ffw.startTime = time.Now()
	httpWriteMuxer := flv.NewMuxer(ffw.httpWrite)
	demuxer := flv.NewDemuxer(ffw)
	ffw.muxer = demuxer
	codecs, err := demuxer.Streams()
	if err != nil {
		logs.Error("read codecs from flv file error : %v", err)
		return
	}
	ffw.codecs = codecs
	httpWriteMuxer.WriteHeader(codecs)

	if ffw.seekSecond > 0 {
		for {
			pkt, err := demuxer.ReadPacket()
			if err != nil {
				logs.Error("read file %s ReadPacket error : %v", ffw.fileName, err)
				break
			}
			if pkt.Time >= time.Duration(ffw.seekSecond)*time.Second {
				break
			}
		}
	}

	timeStart := time.Now()

Loop:
	for {
		pkt, err := demuxer.ReadPacket()
		if err != nil {
			logs.Error("read file %s ReadPacket error : %v", ffw.fileName, err)
			break
		}
		if !ffw.fgStart && pkt.IsKeyFrame {
			continue
		}
		ffw.fgStart = true
		err = httpWriteMuxer.WritePacket(pkt)
		if err != nil {
			logs.Error("read file %s WritePacket error : %v", ffw.fileName, err)
			break
		}
		// fmt.Printf("pkt.Time=%s\n", pkt.Time)
		sinceTime := time.Since(timeStart) + time.Duration(ffw.seekSecond)*time.Second
		// fmt.Printf("sinceTime=%s\n", sinceTime)
		// fmt.Printf("result=%t\n", pkt.Time > (sinceTime+60*time.Second))
		if pkt.Time > (sinceTime + 3*time.Minute) {
			select {
			case <-ffw.done:
				break Loop
			case <-time.NewTicker(5 * time.Second).C:
			}
		}
		select {
		case <-ffw.done:
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

func FlvDurationRead(fileName string) (n int, err error) {
	fullFileName := getFileFlvPath() + "/" + fileName + ".flv"
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
	for {
		pkt, readErr := demuxer.ReadPacket()
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			logs.Error("read file %s ReadPacket error : %v", fileName, err)
			err = readErr
			break
		}

		totalTime = pkt.Time
	}
	n = int(totalTime / time.Millisecond)
	return
}
