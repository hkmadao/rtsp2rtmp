package camera

import "time"

// 摄像头
type CameraVO struct {
	// 摄像头主属性
	Id string `json:"id"`
	// code:
	Code string `json:"code"`
	// rtsp地址:
	RtspUrl string `json:"rtspUrl"`
	// rtmp推送地址:
	RtmpUrl string `json:"rtmpUrl"`
	// 播放权限码:
	PlayAuthCode string `json:"playAuthCode"`
	// 在线状态:
	OnlineStatus int `json:"onlineStatus"`
	// 启用状态:
	Enabled int `json:"enabled"`
	// rtmp推送状态:
	RtmpPushStatus int `json:"rtmpPushStatus"`
	// 保存录像状态:
	SaveVideo int `json:"saveVideo"`
	// 直播状态:
	Live int `json:"live"`
	// 创建时间:
	Created time.Time `json:"created"`
	// 摄像头分享
	CameraShares []CameraShareVO `vo:"ignore" json:"cameraShares"`
}

// 摄像头分享
type CameraShareVO struct {
	// 摄像头分享主属性
	Id string `json:"id"`
	// 名称:
	Name string `json:"name"`
	// 权限码:
	AuthCode string `json:"authCode"`
	// 启用状态:
	Enabled int `json:"enabled"`
	// 创建时间:
	Created time.Time `json:"created"`
	// 开始时间:
	StartTime time.Time `json:"startTime"`
	// 结束时间:
	Deadline time.Time `json:"deadline"`
	// 摄像头id:
	CameraId string `json:"cameraId"`
}
