package ffmpegmanager

import (
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/go-cmd/cmd"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	ext_controller "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controllers/ext"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

type FFmpegInfo struct {
	startTime      time.Time
	templatePasswd string
	ffMpegCmd      *cmd.Cmd
}

var rcmInstance *FFmpegManager

func init() {
	rcmInstance = &FFmpegManager{}
}

type FFmpegManager struct {
	rcs sync.Map
}

func GetSingleFFmpegManager() *FFmpegManager {
	return rcmInstance
}

func (rs *FFmpegManager) StartClient() {
	go rs.startConnections()
	go rs.stopConn(ext_controller.CodeStream())
	go rs.checkBlockFFmpegProcess()
}

func (rc *FFmpegManager) ExistsPublisher(code string) bool {
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

func (rs *FFmpegManager) checkBlockFFmpegProcess() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		condition := common.GetEmptyCondition()
		es, err := base_service.CameraFindCollectionByCondition(condition)
		if err != nil {
			logs.Error("camera list query error: %s", err)
			return
		}
		for _, camera := range es {
			if camera.OnlineStatus != 1 {
				v, b := rs.rcs.Load(camera.Code)
				if b {
					ffmpegInfo := v.(FFmpegInfo)
					if time.Since(ffmpegInfo.startTime) > 30*time.Second {
						logs.Info("camera [%s] rtsp block, close it", camera.Code)
						err := ffmpegInfo.ffMpegCmd.Stop()
						if err != nil {
							logs.Error("camera [%s] rtsp block, close it error: %v", err)
						}
						rs.rcs.Delete(camera.Code)
					}
				}
			}
			<-time.NewTicker(10 * time.Second).C
		}
	}
}

func (rs *FFmpegManager) stopConn(codeStream <-chan string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()

	for code := range codeStream {
		v, b := rs.rcs.Load(code)
		if b {
			ffmpegInfo := v.(FFmpegInfo)
			ffmpegInfo.ffMpegCmd.Stop()
			logs.Info("camera [%s] close success", code)
		} else {
			logs.Info("ffmpeg proccess not exist, needn't close: %s", code)
		}
	}
}

func (s *FFmpegManager) startConnections() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("rtspManager panic %v", r)
		}
	}()
	for {
		condition := common.GetEmptyCondition()
		es, err := base_service.CameraFindCollectionByCondition(condition)
		if err != nil {
			logs.Error("camera list query error: %s", err)
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

func (ffmpegManager *FFmpegManager) connRtsp(code string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	defer func() {
		ffmpegManager.rcs.Delete(code)
	}()
	condition := common.GetEqualCondition("code", code)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("find camera [%s] error : %v", code, err)
		return
	}
	if camera.Enabled != 1 {
		logs.Error("camera [%s] disabled : %v", code)
		return
	}
	rtmpPort, err := config.Int("server.rtmp.port")
	if err != nil {
		logs.Error("get rtmp port fail : %v", err)
		return
	}
	templatePasswd, _ := utils.UUID()
	rtmpUrl := fmt.Sprintf("rtmp://127.0.0.1:%d/%s/%s", rtmpPort, code, templatePasswd)
	portOpen := checkTargetPortStatus(camera.RtspUrl)
	if !portOpen {
		logs.Error("rtspUrl: %s port not open", camera.RtspUrl)
		return
	}
	// 只支持h264编码, 使用"-c:v copy", 不要使用其他选项, 出发视频转码会导致cpu很高
	ffmpegCmd := cmd.NewCmd("ffmpeg", "-i", camera.RtspUrl, "-c:v", "copy", "-c:a", "aac", "-f", "flv", rtmpUrl)
	ffmpegManager.rcs.Store(code, FFmpegInfo{startTime: time.Now(), templatePasswd: templatePasswd, ffMpegCmd: ffmpegCmd})
	logs.Info("ffmpeg start connect rtsp : command : %s %s", ffmpegCmd.Name, strings.Join(ffmpegCmd.Args, " "))
	statusChan := ffmpegCmd.Start()
	finalStatus := <-statusChan
	if finalStatus.Error != nil {
		logs.Error("ffmpeg start connect rtsp failed:", finalStatus.Error)
	} else {
		logs.Info("ffmpeg complate connect rtsp : %s", code)
	}

}

func (r *FFmpegManager) Load(key interface{}) (interface{}, bool) {
	return r.rcs.Load(key)
}
func (r *FFmpegManager) Store(key, value interface{}) {
	r.rcs.Store(key, value)
}
func (r *FFmpegManager) Delete(key interface{}) {
	r.rcs.Delete(key)
}

func ValiadRtmpInfo(code string, passwd string) (fgSuccess bool) {
	v, ok := GetSingleFFmpegManager().rcs.Load(code)
	if !ok {
		fgSuccess = false
		return
	}
	fFmpegInfo, ok := v.(FFmpegInfo)
	if !ok {
		fgSuccess = false
		return
	}
	if fFmpegInfo.templatePasswd != passwd {
		fgSuccess = false
		return
	}
	fgSuccess = true
	return
}

func checkTargetPortStatus(rtspUrl string) (open bool) {
	// rtsp://127.0.0.1:554/2
	urlSplits := strings.Split(rtspUrl, "/")
	if len(urlSplits) < 3 {
		logs.Warn("rtspUrl: %s error", rtspUrl)
		open = false
		return
	}
	address := urlSplits[2]
	timeout := 5 * time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		logs.Warn("rtspUrl: %s port not open", rtspUrl)
		open = false
		return
	}
	defer conn.Close()

	open = true
	return
}
