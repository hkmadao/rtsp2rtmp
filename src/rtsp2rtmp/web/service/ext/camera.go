package ext

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
)

func CameraFindRecordFiles(e entity.Camera) (i int64, err error) {
	fileFlvPath, err := getFileFlvPath()
	if err != nil {
		return
	}
	logs.Info("fileFlvPath : %s", fileFlvPath)
	entries, err := os.ReadDir(fileFlvPath)
	if err != nil {
		return
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, errEntry := entry.Info()
		if errEntry != nil {
			err = errEntry
			return
		}
		if !info.IsDir() {
			infos = append(infos, info)
		}
	}
	fmt.Printf("%v\n", infos)
	return
}

func getFileFlvPath() (fileFlvPath string, err error) {
	fileFlvPath, err = config.String("server.fileflv.path")
	if err != nil {
		logs.Error("get fileflv path error : %v", err)
	}
	return
}
