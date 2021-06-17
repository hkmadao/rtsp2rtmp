package app

import (
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/yumrano/rtsp2rtmp/models"
	"github.com/yumrano/rtsp2rtmp/server"
)

type task struct {
}

func NewTask() *task {
	t := &task{}
	go t.offlineCamera()
	return t
}

func (t *task) offlineCamera() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		css, err := models.CameraSelectAll()
		if err != nil {
			logs.Error("query camera error : %v", err)
		}
		for _, cs := range css {
			if cs.OnlineStatus != 1 {
				continue
			}
			exist := server.ExistCamera(cs.Code)
			if !exist {
				cs.OnlineStatus = 0
				models.CameraUpdate(cs)
			}
		}
		<-time.After(10 * time.Minute)
	}
}
