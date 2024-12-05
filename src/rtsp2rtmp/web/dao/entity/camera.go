package entity

import (
	"time"
)

type Camera struct {
	Id             string        `orm:"pk;column(id)" json:"id"`
	Code           string        `orm:"column(code)" json:"code"`
	RtspURL        string        `orm:"column(rtsp_url)" json:"rtspURL"`
	RtmpURL        string        `orm:"column(rtmp_url)" json:"rtmpURL"`
	PlayAuthCode   string        `orm:"column(play_auth_code)" json:"playAuthCode"`
	OnlineStatus   int           `orm:"column(online_status)" json:"onlineStatus"`
	Enabled        int           `orm:"column(enabled)" json:"enabled"`
	RtmpPushStatus int           `orm:"column(rtmp_push_status)" json:"rtmpPushStatus"`
	SaveVideo      int           `orm:"column(save_video)" json:"saveVideo"`
	Live           int           `orm:"column(live)" json:"live"`
	Created        time.Time     `orm:"column(created)" json:"created"`
	CameraShares   []CameraShare `orm:"-" json:"cameraShares"`
}
