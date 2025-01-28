package tcpclient

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/fileflvmanager/fileflvreader"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

var playerMap sync.Map

type PlayParam struct {
	PlayerId       string `json:"playerId"`
	IdCameraRecord string `json:"idCameraRecord"`
	SeekSecond     uint64 `json:"seekSecond"`
}

type FlvPush struct {
	conn net.Conn
}

// override io.Writer
func (flvPush FlvPush) Write(p []byte) (n int, err error) {
	messageBytes := p
	// messageLenBytes := utils.Int32ToByteBigEndian(int32(len(messageBytes)))
	// fullMessageBytes := append(messageLenBytes, messageBytes...)
	// n, err = flvPush.conn.Write(fullMessageBytes)
	// if err != nil {
	// 	logs.Error("register error: %v", err)
	// 	return
	// }
	// return len(p), nil
	secretStr, err := config.String("server.remote.secret")
	if err != nil {
		logs.Error("get remote secret error: %v", err)
		return
	}
	encryptMessageBytes, err := utils.EncryptAES([]byte(secretStr), messageBytes)
	if err != nil {
		logs.Error("EncryptAES error: %v", err)
		return
	}

	encryptMessageLen := len(encryptMessageBytes)
	encryptMessageLenBytes := utils.Int32ToByteBigEndian(int32(encryptMessageLen))
	fullMessageBytes := append(encryptMessageLenBytes, encryptMessageBytes...)
	_, err = flvPush.conn.Write(fullMessageBytes)
	if err != nil {
		logs.Error("flvPush Write error: %v", err)
		return
	}
	return len(p), nil
}

func flvPlay(commandMessage CommandMessage) {
	paramStr := commandMessage.Param
	playParam := PlayParam{}
	err := json.Unmarshal([]byte(paramStr), &playParam)
	if err != nil {
		logs.Error("flvPlay message format error: %v", err)
		return
	}
	conn, err := connectAndRegister("flvPlay", commandMessage.MessageId)
	if err != nil {
		logs.Error("flvPlay connect to server error: %v", err)
		return
	}
	defer conn.Close()

	camera_record, err := base_service.CameraRecordSelectById(playParam.IdCameraRecord)
	if err != nil {
		logs.Error("CameraRecordSelectById error: %v", err)
		result := common.ErrorResult(fmt.Sprintf("idCameraRecord: %s CameraRecordSelectById error", playParam.IdCameraRecord))
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	fileName := camera_record.FileName
	if camera_record.FgTemp {
		fileName = camera_record.TempFileName
	}

	flvPush := FlvPush{conn: conn}

	ffr := fileflvreader.NewFileFlvReader(playParam.SeekSecond, flvPush, fileName)
	_, ok := playerMap.Load(playParam.PlayerId)
	if ok {
		logs.Error("playerId: %s exists", playParam.PlayerId)
		result := common.ErrorResult(fmt.Sprintf("playerId: %s exists", playParam.PlayerId))
		_, err = writeResult(result, conn)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	playerMap.Store(playParam.PlayerId, ffr)
	<-ffr.GetDone()
	playerMap.Delete(playParam.PlayerId)
	logs.Info("vod player [%s] addr [%s] exit", fileName, conn.LocalAddr().String())
}
