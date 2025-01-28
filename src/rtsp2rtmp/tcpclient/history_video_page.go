package tcpclient

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	dto_convert "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controller/convert"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func historyVideoPage(commandMessage CommandMessage) {
	paramStr := commandMessage.Param
	pageInfoInput := common.AqPageInfoInput{}
	err := json.Unmarshal([]byte(paramStr), &pageInfoInput)
	if err != nil {
		logs.Error("historyVideoPage message format error: %v", err)
		return
	}
	conn, err := connectAndRegister("historyVideoPage", commandMessage.MessageId)
	if err != nil {
		logs.Error("historyVideoPage connect to server error: %v", err)
		return
	}
	defer conn.Close()

	pageInfo, err := base_service.CameraRecordFindPageByCondition(pageInfoInput)
	if err != nil {
		logs.Error("aqPage error : %v", err)
		result := common.ErrorResult("CameraRecordFindPageByCondition error")
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	var cameraRecords = make([]entity.CameraRecord, 0)
	for _, data := range pageInfo.DataList {
		cameraRecords = append(cameraRecords, data.(entity.CameraRecord))
	}
	voList, err := dto_convert.ConvertCameraRecordToVOList(cameraRecords)
	if err != nil {
		logs.Error("aqPage error: %v", err)
		result := common.ErrorResult(fmt.Sprintf("ConvertCameraRecordToVOList error"))
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	var dataList = make([]interface{}, 0)
	for _, vo := range voList {
		dataList = append(dataList, vo)
	}
	pageInfo.DataList = dataList
	result := common.SuccessResultData(pageInfo)
	_, err = writeResult(result, conn)
	if err != nil {
		logs.Error(err)
		return
	}
}
