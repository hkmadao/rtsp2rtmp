package tcpclient

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/tcpclient/tcpclientcommon"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/ext/live"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func getLiveMediaInfo(commandMessage tcpclientcommon.CommandMessage) {
	conn, err := tcpclientcommon.ConnectAndResRegister("getLiveMediaInfo", commandMessage.MessageId)
	if err != nil {
		logs.Error("getLiveMediaInfo connect to server error: %v", err)
		return
	}
	defer conn.Close()

	paramStr := commandMessage.Param
	rtmpPushParam := RtmpPushParam{}
	err = json.Unmarshal([]byte(paramStr), &rtmpPushParam)
	if err != nil {
		logs.Error("getLiveMediaInfo message format error: %v", err)
		result := common.ErrorResult(fmt.Sprintf("getLiveMediaInfo message format error: %v", err))
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}

	condition := common.GetEqualCondition("code", rtmpPushParam.CameraCode)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("CameraFindOneByCondition error: %v", err)
		result := common.ErrorResult("CameraFindOneByCondition error")
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}

	if !camera.RtmpPushStatus {
		liveMediaInfo := live.LiveMediaInfo{
			HasAudio:     false,
			OnlineStatus: false,
			AnchorName:   camera.Code,
		}
		result := common.SuccessResultData(liveMediaInfo)
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
	}

	liveMediaInfo, err := flvadmin.GetSingleHttpFlvAdmin().GetLiveInfo(camera.Code)
	if err != nil {
		logs.Error("getLiveInfo error : %v", err)
		result := common.ErrorResult("getLiveInfo error")
		_, err = tcpclientcommon.WriteResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}

	result := common.SuccessResultData(liveMediaInfo)
	_, err = tcpclientcommon.WriteResult(result, conn)
	if err != nil {
		logs.Error(err)
		return
	}
}
