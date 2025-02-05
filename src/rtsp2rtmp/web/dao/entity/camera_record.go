package entity

import (
	"time"
)
// 摄像头记录
type CameraRecord struct {
	// 记录id
	IdCameraRecord string  `orm:"pk;column(id_camera_record)" json:"idCameraRecord"`
	// 创建时间:
	Created time.Time `orm:"column(created)" json:"created"`
	// 临时文件名称:
	TempFileName string `orm:"column(temp_file_name)" json:"tempFileName"`
	// 临时文件标志:
	FgTemp bool `orm:"column(fg_temp)" json:"fgTemp"`
	// 文件名称:
	FileName string `orm:"column(file_name)" json:"fileName"`
	// 文件删除标志:
	FgRemove bool `orm:"column(fg_remove)" json:"fgRemove"`
	// 文件时长
	Duration uint32 `orm:"column(duration)" json:"duration"`
	// 开始时间:
	StartTime time.Time `orm:"column(start_time)" json:"startTime"`
	// 结束时间:
	EndTime time.Time `orm:"column(end_time)" json:"endTime"`
	// 是否有音频
	HasAudio bool `orm:"column(has_audio)" json:"hasAudio"`
	// 摄像头主属性:
	IdCamera string `orm:"column(id_camera)" json:"idCamera"`
	// 摄像头:
	Camera Camera `orm:"-" json:"camera"`
}
