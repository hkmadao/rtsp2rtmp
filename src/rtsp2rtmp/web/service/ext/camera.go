package ext

import (
	"fmt"
	"os"
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	ext_vo "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/ext/camera_record"
)

func CameraFindRecordFiles() (fileInfoList []ext_vo.RecordFileInfo, err error) {
	fileFlvPath, err := getFileFlvPath()
	if err != nil {
		return
	}
	logs.Info("fileFlvPath : %s", fileFlvPath)
	entries, err := os.ReadDir(fileFlvPath)
	if err != nil {
		return
	}
	fileInfoList = make([]ext_vo.RecordFileInfo, 0)
	for _, entry := range entries {
		info, errEntry := entry.Info()
		if errEntry != nil {
			err = errEntry
			return
		}
		if info.IsDir() {
			continue
		}
		fileName := info.Name()
		if !strings.HasSuffix(fileName, ".flv") {
			continue
		}
		fileInfo := ext_vo.RecordFileInfo{
			FileName: fileName,
			ModTime:  info.ModTime(),
			Size:     info.Size(),
		}
		fileInfoList = append(fileInfoList, fileInfo)
	}
	fmt.Printf("%v\n", fileInfoList)
	return
}

func getFileFlvPath() (fileFlvPath string, err error) {
	fileFlvPath, err = config.String("server.fileflv.path")
	if err != nil {
		logs.Error("get fileflv path error : %v", err)
	}
	return
}
