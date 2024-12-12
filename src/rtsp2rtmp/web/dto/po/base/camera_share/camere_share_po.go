package camerashare

import "time"

// 摄像头分享
type CameraSharePO struct {
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
}
