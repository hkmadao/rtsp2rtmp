package ext

import "time"

// 记录
type RecordFileInfo struct {
	// 文件名称
	FileName string `json:"fileName"`
	// 文件大小
	Size int64 `json:"size"`
	// 最后修改时间
	ModTime time.Time `json:"modTime"`
}
