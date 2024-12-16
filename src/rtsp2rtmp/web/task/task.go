package task

import (
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtmpserver"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

var taskInstance *task

func init() {
	taskInstance = &task{}
}

type task struct {
}

func GetSingleTask() *task {
	return taskInstance
}

func (t *task) StartTask() {
	go t.clearToken()
	go t.offlineCamera()
}

func (t *task) clearToken() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		web.ClearExipresToken()
		<-time.After(24 * time.Hour)
	}
}

func (t *task) offlineCamera() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		condition := common.GetEmptyCondition()
		css, err := base_service.CameraFindCollectionByCondition(condition)
		if err != nil {
			logs.Error("query camera error : %v", err)
		}
		for _, cs := range css {
			if cs.OnlineStatus != 1 {
				continue
			}
			if exists := rtmpserver.GetSingleRtmpServer().ExistsPublisher(cs.Code); !exists {
				cs.OnlineStatus = 0
				base_service.CameraUpdateById(cs)
			}
		}
		<-time.After(10 * time.Minute)
	}
}
