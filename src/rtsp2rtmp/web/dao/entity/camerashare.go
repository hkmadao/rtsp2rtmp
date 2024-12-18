package entity

import (
	"time"
)

// 摄像头分享
type CameraShare struct {
	// 摄像头分享主属性
	Id string `orm:"pk;column(id)" json:"id"`
	// 名称:
	Name string `orm:"column(name)" json:"name"`
	// 权限码:
	AuthCode string `orm:"column(auth_code)" json:"authCode"`
	// 启用状态:
	Enabled int `orm:"column(enabled)" json:"enabled"`
	// 创建时间:
	Created time.Time `orm:"column(created)" json:"created"`
	// 开始时间:
	StartTime time.Time `orm:"column(start_time)" json:"startTime"`
	// 结束时间:
	Deadline time.Time `orm:"column(deadline)" json:"deadline"`
	// 摄像头id:
	CameraId string `orm:"column(camera_id)" json:"cameraId"`
	// 摄像头:
	Camera Camera `orm:"-" json:"camera"`
}
