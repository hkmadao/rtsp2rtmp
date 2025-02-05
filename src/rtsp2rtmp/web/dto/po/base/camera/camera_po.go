package po

import (
	"time"
)

// 摄像头
type CameraPO struct {
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
	// 加密标志:
	FgEncrypt bool `json:"fgEncrypt"`
	// 被动推送rtmp标志
	FgPassive bool `json:"fgPassive"`
	// rtmp识别码:
	RtmpAuthCode string `json:"rtmpAuthCode"`
	// 摄像头类型:
	CameraType string `json:"cameraType"`
	// 摄像头分享
	// CameraShares []CameraShare `json:"cameraShares"`
}
