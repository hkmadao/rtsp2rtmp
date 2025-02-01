package rtmpflvwriter

import (
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

type IRtmpFlvManager interface {
	UpdateFFWS(string, *RtmpFlvWriter)
}

type RtmpFlvWriter struct {
	needPushRtmp  bool
	stop          bool
	done          chan int
	pktStream     <-chan av.Packet
	code          string
	codecs        []av.CodecData
	start         bool
	startTime     time.Time
	conn          *rtmp.Conn
	pulseInterval time.Duration
	irfm          IRtmpFlvManager
}

func (rfw *RtmpFlvWriter) GetDone() <-chan int {
	return rfw.done
}

func (rfw *RtmpFlvWriter) GetPktStream() <-chan av.Packet {
	return rfw.pktStream
}

func (rfw *RtmpFlvWriter) GetCode() string {
	return rfw.code
}

func (rfw *RtmpFlvWriter) SetCodecs(codecs []av.CodecData) {
	logs.Warn("RtmpFlvWriter: %s update codecs", rfw.code)
	rfw.codecs = codecs
}

func (rfw *RtmpFlvWriter) GetCodecs() []av.CodecData {
	return rfw.codecs
}

func (rfw *RtmpFlvWriter) GetNeedPushRtmp() bool {
	return rfw.needPushRtmp
}

func NewRtmpFlvWriter(needPushRtmp bool, pktStream <-chan av.Packet, code string, codecs []av.CodecData, irfm IRtmpFlvManager) *RtmpFlvWriter {
	rfw := &RtmpFlvWriter{
		needPushRtmp:  needPushRtmp,
		stop:          false,
		done:          make(chan int),
		pktStream:     pktStream,
		code:          code,
		codecs:        codecs,
		start:         false,
		pulseInterval: 5 * time.Second,
		irfm:          irfm,
	}
	go rfw.flvWrite()
	return rfw
}

func (rfw *RtmpFlvWriter) StopWrite() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		rfw.stop = true
		close(rfw.done)
	}()
}

func (rfw *RtmpFlvWriter) createConn() error {
	condition := common.GetEqualCondition("code", rfw.code)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("not found camera : %s", rfw.code)
		return err
	}
	rtmpConn, err := rtmp.Dial(camera.RtmpUrl)
	if err != nil {
		logs.Error("rtmp client connection error : %v", err)
		return err
	}
	logs.Info("rtmp client connection success : %s", rfw.code)
	rfw.conn = rtmpConn
	return nil
}

// Write extends to writer.Writer
func (rfw *RtmpFlvWriter) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	condition := common.GetEqualCondition("code", rfw.code)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("not found camera : %s", rfw.code)
		return
	}
	if camera.OnlineStatus != true {
		return
	}
	if camera.RtmpPushStatus != true || !rfw.needPushRtmp {
		go func() {
			for {
				select {
				case <-rfw.GetDone():
					return
				case <-rfw.pktStream:
				}
			}
		}()
		return
	}

	ticker := time.NewTicker(rfw.pulseInterval)
	defer func() {
		if rfw.conn != nil {
			logs.Info("disconnect old RtmpFlvWriter : %s", rfw.code)
			rfw.conn.Close()
		}
		if rfw.stop {
			logs.Info("stop RtmpFlvWriter : %s", rfw.code)
			return
		}
		_, pktStreamOk := <-rfw.pktStream
		if pktStreamOk {
			logs.Info("to create NewRtmpFlvWriter : %s", rfw.code)
			rfwn := NewRtmpFlvWriter(true, rfw.pktStream, rfw.code, rfw.codecs, rfw.irfm)
			rfwn.irfm.UpdateFFWS(rfwn.code, rfwn)
		} else {
			logs.Info("RtmpFlvWriter pktStream is closed : %s", rfw.code)
		}
	}()
	pktStream := utils.OrDonePacket(rfw.done, rfw.pktStream)
	timeNow := time.Now().Local()
	for {
		select {
		case <-ticker.C:
			logs.Error("flvWrite ticker time out: %s", rfw.code)
			return
		case pkt, ok := <-pktStream:
			if !ok {
				return
			}
			// if rfw.start {
			// 	pktTime := time.Now().Sub(rfw.startTime)
			// 	pkt.Time = pktTime
			// }
			if err := rfw.writerPacket(pkt, &timeNow); err != nil {
				logs.Error("flvWrite error : %v", err)
				return
			}
			ticker.Reset(rfw.pulseInterval)
		}
	}
}

func (rfw *RtmpFlvWriter) writerPacket(pkt av.Packet, templateTime *time.Time) error {
	if rfw.start {
		if err := rfw.conn.WritePacket(pkt); err != nil {
			logs.Error("writer packet to rtmp server error : %v", err)
			return err
		}
		return nil
	}
	if pkt.IsKeyFrame {
		if err := rfw.createConn(); err != nil {
			logs.Error("conn rtmp server error : %v", err)
			return err
		}
		var err error
		err = rfw.conn.WriteHeader(rfw.codecs)
		rfw.startTime = time.Now()
		logs.Info("KeyFrame WriteHeader to rtmp server : %s, codesc: %v", rfw.code, rfw.codecs)
		if err != nil {
			logs.Error("writer header to rtmp server error : %v", err)
			return err
		}
		logs.Info("KeyFrame WriteHeader to rtmp server success : %s", rfw.code)
		err = rfw.conn.WritePacket(pkt)
		if err != nil {
			logs.Error("writer packet to rtmp server error : %v", err)
			return err
		}
		logs.Info("writer KeyFrame to rtmp server success : %s", rfw.code)
		rfw.start = true
		return nil
	}
	if time.Now().Local().After(templateTime.Add(1 * time.Minute)) {
		*templateTime = time.Now().Local()
		logs.Error("RtmpFlvWriter ingrore package: %s", rfw.code)
	}
	return nil
}
