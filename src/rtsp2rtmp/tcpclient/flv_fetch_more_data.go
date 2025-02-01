package tcpclient

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/fileflvmanager/fileflvreader"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
)

type FetchMoreDataParam struct {
	PlayerId   string `json:"playerId"`
	SeekSecond uint64 `json:"seekSecond"`
}

func flvFetchMoreData(commandMessage CommandMessage) {
	conn, err := connectAndRegister("flvFetchMoreData", commandMessage.MessageId)
	if err != nil {
		logs.Error("flvFetchMoreData connect to server error: %v", err)
		return
	}
	defer conn.Close()

	paramStr := commandMessage.Param
	param := FetchMoreDataParam{}
	err = json.Unmarshal([]byte(paramStr), &param)
	if err != nil {
		logs.Error("flvFetchMoreData message format error: %v", err)
		result := common.ErrorResult(fmt.Sprintf("flvFetchMoreData message format error: %v", err))
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}

	value, ok := playerMap.Load(param.PlayerId)
	if !ok {
		logs.Error("playerId: %s not exists or complate", param.PlayerId)
		result := common.SuccessResultMsg(fmt.Sprintf("playerId: %s not exists or complate, skip this request", param.PlayerId))
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	loadFfr := (value.(*fileflvreader.FileFlvReader))
	loadFfr.SetSeekSecond(param.SeekSecond)

	logs.Info("vod player [%s] fetch data, addr [%s]", param.PlayerId, conn.LocalAddr().String())
	result := common.SuccessResultMsg("fetch sccess")
	_, err = writeResult(result, conn)
	if err != nil {
		logs.Error(err)
		return
	}
}
