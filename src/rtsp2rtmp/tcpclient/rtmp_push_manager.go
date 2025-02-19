package tcpclient

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/tcpclient/tcpclientcommon"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

type RtmpPushParam struct {
	CameraCode string `json:"cameraCode"`
}

func startRtmpPush(commandMessage tcpclientcommon.CommandMessage) {
	conn, err := tcpclientcommon.ConnectAndResRegister("startPushRtmp", commandMessage.MessageId)
	if err != nil {
		logs.Error("startPushRtmp connect to server error: %v", err)
		return
	}
	defer conn.Close()
	paramStr := commandMessage.Param
	param := RtmpPushParam{}
	err = json.Unmarshal([]byte(paramStr), &param)
	if err != nil {
		logs.Error("startPushRtmp message format error: %v", err)
		result := common.ErrorResult(fmt.Sprintf("startPushRtmp message format error: %v", err))
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	condition := common.GetEqualCondition("code", param.CameraCode)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("startPushRtmp: find camera: %s error: %v", param.CameraCode, err)
		result := common.ErrorResult(fmt.Sprintf("startPushRtmp: find camera: %s error: %v", param.CameraCode, err))
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
	}
	if !camera.FgPassive {
		logs.Error("startPushRtmp: camera: %s fgPassive is false", param.CameraCode)
		result := common.ErrorResult(fmt.Sprintf("startPushRtmp: camera: %s fgPassive is false", param.CameraCode))
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
	}
	flvadmin.GetSingleRtmpFlvAdmin().RemoteStartWrite(param.CameraCode)

	result := common.SuccessResultData("startPushRtmp success")
	_, err = tcpclientcommon.WriteResult(result, conn)
	if err != nil {
		logs.Error(err)
		return
	}
}

func stopRtmpPush(commandMessage tcpclientcommon.CommandMessage) {
	conn, err := tcpclientcommon.ConnectAndResRegister("stopPushRtmp", commandMessage.MessageId)
	if err != nil {
		logs.Error("stopPushRtmp connect to server error: %v", err)
		return
	}
	defer conn.Close()
	paramStr := commandMessage.Param
	param := RtmpPushParam{}
	err = json.Unmarshal([]byte(paramStr), &param)
	if err != nil {
		logs.Error("stopPushRtmp message format error: %v", err)
		result := common.ErrorResult(fmt.Sprintf("stopPushRtmp message format error: %v", err))
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	condition := common.GetEqualCondition("code", param.CameraCode)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("stopPushRtmp: find camera: %s error: %v", param.CameraCode, err)
		result := common.ErrorResult(fmt.Sprintf("stopPushRtmp: find camera: %s error: %v", param.CameraCode, err))
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
	}
	if !camera.FgPassive {
		logs.Error("stopPushRtmp: camera: %s fgPassive is false", param.CameraCode)
		result := common.ErrorResult(fmt.Sprintf("stopPushRtmp: camera: %s fgPassive is false", param.CameraCode))
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
	}
	flvadmin.GetSingleRtmpFlvAdmin().RemoteStopWrite(param.CameraCode)

	defer conn.Close()
	result := common.SuccessResultData("stopPushRtmp success")
	_, err = tcpclientcommon.WriteResult(result, conn)
	if err != nil {
		logs.Error(err)
		return
	}
}
