package tcpclient

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/fileflvmanager/fileflvreader"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/ext/flv_file"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

type FlvFileMediaInfoParam struct {
	IdCameraRecord string `json:"idCameraRecord"`
}

func flvFileMediaInfo(commandMessage CommandMessage) {
	paramStr := commandMessage.Param
	param := FlvFileMediaInfoParam{}
	err := json.Unmarshal([]byte(paramStr), &param)
	if err != nil {
		logs.Error("flvFileMediaInfo message format error: %v", err)
		return
	}
	idCameraRecord := param.IdCameraRecord
	conn, err := connectAndRegister("flvFileMediaInfo", commandMessage.MessageId)
	if err != nil {
		logs.Error("flvFileMediaInfo connect to server error: %v", err)
		return
	}
	camera_record, err := base_service.CameraRecordSelectById(idCameraRecord)
	if err != nil {
		logs.Error("idCameraRecord: %s CameraRecordSelectById error: %v", idCameraRecord, err)
		result := common.ErrorResult(fmt.Sprintf("idCameraRecord: %s CameraRecordSelectById error", idCameraRecord))
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}

	mediaInfo := flv_file.FlvMediaInfo{
		Duration: camera_record.Duration,
		HasAudio: true,
	}

	if camera_record.FgTemp {
		durationInt, err := fileflvreader.FlvDurationReadUntilErr(camera_record.TempFileName)
		mediaInfo = flv_file.FlvMediaInfo{
			Duration: uint32(durationInt),
			HasAudio: true,
		}
		if err != nil {
			logs.Error("file: %s get mediaInfo error", camera_record.TempFileName)
			result := common.ErrorResult(fmt.Sprintf("file: %s get mediaInfo error", camera_record.TempFileName))
			_, err = writeResult(result, conn)
			if err != nil {
				logs.Error(err)
				return
			}
			return
		}
	}

	result := common.SuccessResultData(mediaInfo)
	_, err = writeResult(result, conn)
	if err != nil {
		logs.Error(err)
		return
	}
}
