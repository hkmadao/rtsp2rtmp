package rtspclientmanager

import (
	"runtime/debug"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/controllers"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/models"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtspclient"
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
	for code := range codeStream {
		v, b := rs.conns.Load(code)
		if b {
			r := v.(*rtmp.Conn)
			err := r.Close()
			if err != nil {
				logs.Error("camera [%s] close error : %v", code, err)
				continue
			}
			logs.Info("camera [%s] close success", code)
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
	rtsp.DebugRtsp = false
	session, err := rtsp.Dial(c.RtspURL)
	if err != nil {
		logs.Error("camera [%s] conn : %v", c.Code, err)
		c.OnlineStatus = 0
		time.Sleep(5 * time.Second)
		if c.OnlineStatus == 1 {
			models.CameraUpdate(c)
		}
		return
	}
	session.RtpKeepAliveTimeout = 10 * time.Second
	codecs, err := session.Streams()
	if err != nil {
		logs.Error("camera [%s] get streams : %v", c.Code, err)
		return
	}

	c.OnlineStatus = 1
	models.CameraUpdate(c)

	done := make(chan interface{})
	//添加缓冲，缓解前后速率不一致问题，但是如果收包平均速率大于消费平均速率，依然会导致丢包
	pktStream := make(chan av.Packet, 50)
	defer func() {
		close(done)
		close(pktStream)
	}()

	rc := rtspclient.NewRtspClient(done, pktStream, code, codecs, s)
	s.rcs.Store(code, rc)
	s.conns.Store(code, session)
	for {
		pkt, err := session.ReadPacket()
		if err != nil {
			logs.Error("camera [%s] ReadPacket : %v", c.Code, err)
			break
		}
		//不能开goroutine,不能保证包的顺序
		select {
		case pktStream <- pkt:
		default:
			//添加缓冲，缓解前后速率不一致问题，但是如果收包平均速率大于消费平均速率，依然会导致丢包
			logs.Debug("rtmpserver lose packet")
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

	err = session.Close()
	if err != nil {
		logs.Error("close conn error : %v", err)
	}
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
