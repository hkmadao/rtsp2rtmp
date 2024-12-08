package entity

import (
	"time"
)

// 摄像头
type Camera struct {
	// 摄像头主属性
	Id string `orm:"pk;column(id)" json:"id"`
	// code:
	Code string `orm:"column(code)" json:"code"`
	// rtsp地址:
	RtspUrl string `orm:"column(rtsp_url)" json:"rtspUrl"`
	// rtmp推送地址:
	RtmpUrl string `orm:"column(rtmp_url)" json:"rtmpUrl"`
	// 播放权限码:
	PlayAuthCode string `orm:"column(play_auth_code)" json:"playAuthCode"`
	// 在线状态:
	OnlineStatus int `orm:"column(online_status)" json:"onlineStatus"`
	// 启用状态:
	Enabled int `orm:"column(enabled)" json:"enabled"`
	// rtmp推送状态:
	RtmpPushStatus int `orm:"column(rtmp_push_status)" json:"rtmpPushStatus"`
	// 保存录像状态:
	SaveVideo int `orm:"column(save_video)" json:"saveVideo"`
	// 直播状态:
	Live int `orm:"column(live)" json:"live"`
	// 创建时间:
	Created time.Time `orm:"column(created)" json:"created"`
	// 摄像头分享
	CameraShares []CameraShare `orm:"-" json:"cameraShares"`
}
