package entity

import (
	"time"
)
// 摄像头
type Camera struct {
	// 摄像头主属性
	Id string  `orm:"pk;column(id)" json:"id"`
	// 编号:
	Code string `orm:"column(code)" json:"code"`
	// rtsp地址:
	RtspUrl string `orm:"column(rtsp_url)" json:"rtspUrl"`
	// rtmp推送地址:
	RtmpUrl string `orm:"column(rtmp_url)" json:"rtmpUrl"`
	// 播放权限码:
	PlayAuthCode string `orm:"column(play_auth_code)" json:"playAuthCode"`
	// 在线状态:
	OnlineStatus bool `orm:"column(online_status)" json:"onlineStatus"`
	// 启用状态:
	Enabled bool `orm:"column(enabled)" json:"enabled"`
	// rtmp推送状态:
	RtmpPushStatus bool `orm:"column(rtmp_push_status)" json:"rtmpPushStatus"`
	// 保存录像状态:
	SaveVideo bool `orm:"column(save_video)" json:"saveVideo"`
	// 直播状态:
	Live bool `orm:"column(live)" json:"live"`
	// 创建时间:
	Created time.Time `orm:"column(created)" json:"created"`
	// 摄像头分享
	CameraShares []CameraShare `orm:"-" json:"cameraShares"`
}
