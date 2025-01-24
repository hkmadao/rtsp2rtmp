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

func flvFetchMoreData(paramStr string) {
	param := FetchMoreDataParam{}
	err := json.Unmarshal([]byte(paramStr), &param)
	if err != nil {
		logs.Error("flvPlay message format error: %v", err)
		return
	}
	conn, err := connectAndRegister("flvFetchMoreData")
	if err != nil {
		logs.Error("flvFetchMoreData connect to server error: %v", err)
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
