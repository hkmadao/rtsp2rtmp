package vo

import (
	"time"
)

// 摄像头分享
type CameraShareVO struct {
	// 摄像头分享主属性
	Id string `json:"id"`
	// 名称:
	Name string `json:"name"`
	// 权限码:
	AuthCode string `json:"authCode"`
	// 启用状态:
	Enabled bool `json:"enabled"`
	// 创建时间:
	Created time.Time `json:"created"`
	// 开始时间:
	StartTime time.Time `json:"startTime"`
	// 结束时间:
	Deadline time.Time `json:"deadline"`
	// 摄像头id:
	CameraId string `json:"cameraId"`
	// 摄像头:
	Camera CameraVO `vo:"ignore" json:"camera"`
}
type CameraVO struct {
	// 摄像头主属性
	Id string `json:"id"`
	// 编号:
	Code string `json:"code"`
	// rtsp地址:
	RtspUrl string `json:"rtspUrl"`
	// rtmp推送地址:
	RtmpUrl string `json:"rtmpUrl"`
	// 播放权限码:
	PlayAuthCode string `json:"playAuthCode"`
	// 在线状态:
	OnlineStatus bool `json:"onlineStatus"`
	// 启用状态:
	Enabled bool `json:"enabled"`
	// rtmp推送状态:
	RtmpPushStatus bool `json:"rtmpPushStatus"`
	// 保存录像状态:
	SaveVideo bool `json:"saveVideo"`
	// 直播状态:
	Live bool `json:"live"`
	// 创建时间:
	Created time.Time `json:"created"`
}
