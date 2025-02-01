package tcpclient

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	dto_convert "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controller/convert"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func cameraAq(commandMessage CommandMessage) {
	conn, err := connectAndRegister("cameraAq", commandMessage.MessageId)
	if err != nil {
		logs.Error("cameraAq connect to server error: %v", err)
		return
	}
	defer conn.Close()

	paramStr := commandMessage.Param
	condition := common.AqCondition{}
	err = json.Unmarshal([]byte(paramStr), &condition)
	if err != nil {
		logs.Error("cameraAq message format error: %v", err)
		result := common.ErrorResult(fmt.Sprintf("cameraAq message format error: %v", err))
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}

	cameras, err := base_service.CameraFindCollectionByCondition(condition)
	if err != nil {
		logs.Error("CameraFindCollectionByCondition error: %v", err)
		result := common.ErrorResult("CameraFindCollectionByCondition error")
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	voList, err := dto_convert.ConvertCameraToVOList(cameras)
	if err != nil {
		logs.Error("ConvertCameraToVOList error: %v", err)
		result := common.ErrorResult("ConvertCameraToVOList error")
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	result := common.SuccessResultData(voList)
	_, err = writeResult(result, conn)
	if err != nil {
		logs.Error(err)
		return
	}
}
