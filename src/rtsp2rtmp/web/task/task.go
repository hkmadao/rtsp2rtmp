package task

import (
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/fileflvmanager/fileflvreader"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtmpserver"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/rtspclientmanager"
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
	go t.ClearHistoryVideo()
}

func (t *task) clearToken() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		<-time.After(24 * time.Hour)
	}
}

func (t *task) offlineCamera() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	fgUseFfmpeg, err := config.Bool("server.use-ffmpeg")
	if err != nil {
		logs.Error("get use-ffmpeg fail : %v", err)
		fgUseFfmpeg = false
	}
	for {
		condition := common.GetEqualCondition("onlineStatus", true)
		css, err := base_service.CameraFindCollectionByCondition(condition)
		if err != nil {
			logs.Error("query camera error : %v", err)
		}
		for _, cs := range css {
			if cs.CameraType == "rtmp" || fgUseFfmpeg {
				if exists := rtmpserver.GetSingleRtmpServer().ExistsPublisher(cs.Code); !exists {
					cs.OnlineStatus = false
					base_service.CameraUpdateById(cs)
				}
			} else {
				if exists := rtspclientmanager.GetSingleRtspClientManager().ExistsPublisher(cs.Code); !exists {
					cs.OnlineStatus = false
					base_service.CameraUpdateById(cs)
				}
			}
		}
		<-time.After(10 * time.Minute)
	}
}

func (t *task) ClearHistoryVideo() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		ltCreated := time.Now().Add(-7 * 24 * time.Hour)
		condition := common.GetLtCondition("created", ltCreated.Format(time.RFC3339))
		css, err := base_service.CameraRecordFindCollectionByCondition(condition)
		if err != nil {
			logs.Error("query CameraRecord error : %v", err)
		}
		for _, cs := range css {
			fileExists := false
			if cs.FgTemp {
				fileExists = fileflvreader.FlvFileExists(cs.TempFileName)
			} else {
				fileExists = fileflvreader.FlvFileExists(cs.FileName)
			}
			if !fileExists {
				logs.Info("flv file: %s not exists, clean", cs.FileName)
				base_service.CameraRecordDelete(cs)
			}
		}
		<-time.After(10 * time.Minute)
	}
}
