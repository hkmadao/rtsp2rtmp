package rtspclientmanager

import (
	"runtime/debug"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtspv2"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/controllers"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/models"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtspclient"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
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
	go rs.serveStreams()
	done := make(chan interface{})
	go rs.stopConn(done, controllers.CodeStream())
}

func (rc *RtspClientManager) ExistsPublisher(code string) bool {
	exists := false
	rc.rcs.Range(func(key, value interface{}) bool {
		codeKey := key.(string)
		if code == codeKey {
			exists = true
			return false
		}
		exists = codeKey == code
		return true
	})
	return exists
}

func (rs *RtspClientManager) stopConn(done <-chan interface{}, codeStream <-chan string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	for code := range codeStream {
		v, b := rs.conns.Load(code)
		if b {
			r := v.(*rtspv2.RTSPClient)
			r.Close()
			logs.Info("camera [%s] close success", code)
			rs.rcs.Delete(code)
		} else {
			logs.Error("codeStream error")
		}
	}
}

func (s *RtspClientManager) serveStreams() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("rtspManager panic %v", r)
		}
	}()
	for {
		es, err := models.CameraSelectAll()
		if err != nil {
			logs.Error("camera list is empty")
			return
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
		<-time.After(30 * time.Second)
	}

}

func (s *RtspClientManager) connRtsp(code string) {
	defer func() {
		s.rcs.Delete(code)
		s.conns.Delete(code)
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	//放置信息表示已经开始
	s.rcs.Store(code, struct{}{})
	q := models.Camera{Code: code}
	c, err := models.CameraSelectOne(q)
	if err != nil {
		logs.Error("find camera [%s] error : %v", code, err)
		return
	}
	if c.Enabled != 1 {
		logs.Error("camera [%s] disabled : %v", code)
		return
	}
	logs.Info(c.Code, "connect", c.RtspURL)
	ro := rtspv2.RTSPClientOptions{
		URL:              c.RtspURL,
		Debug:            false,
		DialTimeout:      10 * time.Second,
		ReadWriteTimeout: 10 * time.Second,
		DisableAudio:     true,
	}
	session, err := rtspv2.Dial(ro)
	if err != nil {
		logs.Error("camera [%s] conn : %v", c.Code, err)
		c.OnlineStatus = 0
		time.Sleep(5 * time.Second)
		if c.OnlineStatus == 1 {
			models.CameraUpdate(c)
		}
		return
	}
	codecs := session.CodecData

	c.OnlineStatus = 1
	models.CameraUpdate(c)

	done := make(chan int)
	//添加缓冲，缓解前后速率不一致问题，但是如果收包平均速率大于消费平均速率，依然会导致丢包
	pktStream := make(chan av.Packet, 50)
	defer func() {
		close(done)
		close(pktStream)
	}()

	rc := rtspclient.NewRtspClient(done, pktStream, code, codecs, s)
	s.rcs.Store(code, rc)
	s.conns.Store(code, session)
	logs.Info("%s", string(session.SDPRaw))
	for pkt := range utils.OrDoneRefPacket(done, session.OutgoingPacketQueue) {
		//不能开goroutine,不能保证包的顺序
		select {
		case pktStream <- pkt:
		default:
			//添加缓冲，缓解前后速率不一致问题，但是如果收包平均速率大于消费平均速率，依然会导致丢包
			logs.Debug("rtspclient lose packet")
		}
	}

	if err != nil {
		logs.Error("session Close error : %v", err)
	}
	//offline camera
	camera, err := models.CameraSelectOne(q)
	if err != nil {
		logs.Error("no camera error : %s", code)
	} else {
		camera.OnlineStatus = 0
		models.CameraUpdate(camera)
	}

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
