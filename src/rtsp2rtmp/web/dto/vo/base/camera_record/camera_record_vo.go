package vo

import (
	"time"
)

// 摄像头记录
type CameraRecordVO struct {
	// 记录id
	IdCameraRecord string `json:"idCameraRecord"`
	// 创建时间:
	Created time.Time `json:"created"`
	// 临时文件名称:
	TempFileName string `json:"tempFileName"`
	// 临时文件标志:
	FgTemp bool `json:"fgTemp"`
	// 文件名称:
	FileName string `json:"fileName"`
	// 文件删除标志:
	FgRemove bool `json:"fgRemove"`
	// 文件时长
	Duration uint32 `json:"duration"`
	// 开始时间:
	StartTime time.Time `json:"startTime"`
	// 结束时间:
	EndTime time.Time `json:"endTime"`
	// 是否有音频
	HasAudio bool `json:"hasAudio"`
	// 摄像头主属性:
	IdCamera string `json:"idCamera"`
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
	// 加密标志:
	FgEncrypt bool `json:"fgEncrypt"`
	// 被动推送rtmp标志
	FgPassive bool `json:"fgPassive"`
	// rtmp识别码:
	RtmpAuthCode string `json:"rtmpAuthCode"`
	// 摄像头类型:
	CameraType string `json:"cameraType"`
}
