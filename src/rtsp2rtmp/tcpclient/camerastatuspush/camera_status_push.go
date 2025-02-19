package camerastatuspush

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/tcpclient/tcpclientcommon"
)

func CameraOnlinePush(cameraCode string) {
	conn, err := tcpclientcommon.ConnectAndCameraStatusRegister("cameraOnline", cameraCode)
	if err != nil {
		logs.Error("cameraAq connect to server error: %v", err)
		return
	}
	defer conn.Close()
}

func CameraOfflinePush(cameraCode string) {
	conn, err := tcpclientcommon.ConnectAndCameraStatusRegister("cameraOffline", cameraCode)
	if err != nil {
		logs.Error("cameraAq connect to server error: %v", err)
		return
	}
	defer conn.Close()
}
