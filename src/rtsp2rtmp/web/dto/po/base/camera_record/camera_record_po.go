package po

import (
	"time"
)

// 摄像头记录
type CameraRecordPO struct {
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
	// Camera Camera `json:"camera"`
}
