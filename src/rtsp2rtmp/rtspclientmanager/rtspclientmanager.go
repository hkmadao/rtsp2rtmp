package rtspclientmanager

import (
	"runtime/debug"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtspv2"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtspclient"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controllers"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service"
)

var rcmInstance *RtspClientManager

func init() {
	rcmInstance = &RtspClientManager{}
}

type RtspClientManager struct {
	rcs   sync.Map
	conns sync.Map
}

func GetSingleRtspClientManager() *RtspClientManager {
	return rcmInstance
}

func (rs *RtspClientManager) StartClient() {
	go rs.startConnections()
	go rs.stopConn(controllers.CodeStream())
}

func (rc *RtspClientManager) ExistsPublisher(code string) bool {
	exists := false
	rc.rcs.Range(func(key, value interface{}) bool {
		codeKey := key.(string)
		if codeKey == code {
			exists = true
			return false
		}
		return true
	})
	return exists
}

func (rs *RtspClientManager) stopConn(codeStream <-chan string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	for code := range codeStream {
		rs.rcs.Delete(code)
		v, b := rs.conns.Load(code)
		if b {
			r := v.(*rtspv2.RTSPClient)
			r.Close()
			close(r.OutgoingPacketQueue)
			logs.Info("camera [%s] close success", code)
			rs.conns.Delete(code)
		} else {
			logs.Info("RtspClient not exist, needn't close: %s", code)
		}
	}
}

func (s *RtspClientManager) startConnections() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("rtspManager panic %v", r)
		}
	}()
	es, err := service.CameraSelectAll()
	if err != nil {
		logs.Error("camera list query error: %s", err)
		return
	}
	timeTemp := time.Now()
	for {
		timeNow := time.Now()
		if timeNow.After(timeTemp.Add(30 * time.Second)) {
			es, err = service.CameraSelectAll()
			if err != nil {
				logs.Error("camera list query error: %s", err)
				return
			}
			timeTemp = timeNow
		}
		for _, camera := range es {
			if v, b := s.rcs.Load(camera.Code); b && v != nil {
				continue
			}
			if camera.Enabled != 1 {
				continue
			}
			go s.connRtsp(camera.Code)
		}
		<-time.After(1 * time.Second)
	}

}

func (s *RtspClientManager) connRtsp(code string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	defer func() {
		s.rcs.Delete(code)
		s.conns.Delete(code)
	}()
	//放置信息表示已经开始
	s.rcs.Store(code, struct{}{})
	q := entity.Camera{Code: code}
	c, err := service.CameraSelectOne(q)
	if err != nil {
		logs.Error("find camera [%s] error : %v", code, err)
		return
	}
	if c.Enabled != 1 {
		logs.Error("camera [%s] disabled : %v", code)
		return
	}
	logs.Info(c.Code, "connect", c.RtspUrl)
	rtspClientOptions := rtspv2.RTSPClientOptions{
		URL:              c.RtspUrl,
		Debug:            false,
		DialTimeout:      10 * time.Second,
		ReadWriteTimeout: 10 * time.Second,
		DisableAudio:     false,
	}
	session, err := rtspv2.Dial(rtspClientOptions)
	if err != nil {
		logs.Error("camera [%s] conn : %v", c.Code, err)
		c.OnlineStatus = 0
		if c.OnlineStatus == 1 {
			service.CameraUpdate(c)
		}
		return
	}
	codecs := session.CodecData
	// logs.Warn("camera: %s codecs: %v", code, session.CodecData)

	c.OnlineStatus = 1
	service.CameraUpdate(c)

	done := make(chan int)
	//添加缓冲，缓解前后速率不一致问题，但是如果收包平均速率大于消费平均速率，依然会导致丢包
	pktStream := make(chan av.Packet, 1024)
	defer func() {
		close(done)
		close(pktStream)
	}()

	rc := rtspclient.NewRtspClient(done, pktStream, code, codecs)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		for {
			select {
			case _, ok := <-session.Signals:
				if !ok {
					return
				}
				logs.Warn("camera: %s update codecs: %v", code, session.CodecData)
				rc.UpdateCodecs(session.CodecData)
			case <-done:
				return
			}
		}
	}()
	s.rcs.Store(code, rc)
	s.conns.Store(code, session)
	logs.Info("%s", string(session.SDPRaw))
	ticker := time.NewTicker(10 * time.Second)
	rtspStream := utils.OrDoneRefPacket(done, session.OutgoingPacketQueue)
Loop:
	for {
		select {
		case pkt, ok := <-rtspStream:
			if !ok {
				logs.Error("camera: %s rtsp packet stream is close", code)
				break Loop
			}
			//不能开goroutine,不能保证包的顺序
			select {
			case pktStream <- pkt:
			default:
				//添加缓冲，缓解前后速率不一致问题，但是如果收包平均速率大于消费平均速率，依然会导致丢包
				logs.Debug("rtspclient lose packet")
			}
			ticker.Reset(10 * time.Second)
		case <-ticker.C:
			logs.Error("camera: %s read packet from rtsp time out", code)
			break Loop
		}
	}

	//offline camera
	camera, err := service.CameraSelectOne(q)
	if err != nil {
		logs.Error("no camera error : %s", code)
	} else {
		camera.OnlineStatus = 0
		service.CameraUpdate(camera)
	}

	logs.Error("camera: %s session Close", code)
	session.Close()
}

func (r *RtspClientManager) Load(key interface{}) (interface{}, bool) {
	return r.rcs.Load(key)
}
func (r *RtspClientManager) Store(key, value interface{}) {
	r.rcs.Store(key, value)
}
func (r *RtspClientManager) Delete(key interface{}) {
	r.rcs.Delete(key)
}
