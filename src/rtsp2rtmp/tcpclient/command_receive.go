package tcpclient

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
)

func StartCommandReceiveServer() {
	go func() {
		for {
			commandReceiveConnect()
			<-time.NewTicker(1 * time.Minute).C
		}
	}()
}

func commandReceiveConnect() {
	conn, err := connectAndRegister("keepChannel")
	if err != nil {
		logs.Error("keepChannel connect to server error: %v", err)
		return
	}
	logs.Info("keepChannel connect to server successful, remote: %s", conn.RemoteAddr().String())
	// read command
	for {
		dataLenBytes := make([]byte, 4)
		_, err := conn.Read(dataLenBytes)
		if err != nil {
			logs.Error("read len error: %v", err)
			break
		}
		var dataLen int32
		buffer := bytes.NewBuffer(dataLenBytes)
		err = binary.Read(buffer, binary.LittleEndian, &dataLen)
		if err != nil {
			logs.Error(err)
			break
		}
		logs.Info("receive message dataLen: %d", dataLen)
		serverRepBytes := make([]byte, dataLen)
		_, err = conn.Read(serverRepBytes)
		if err != nil {
			logs.Error("read message body error: %v", err)
			break
		}
		secretCommandStr := string(serverRepBytes)
		logs.Info("receive message: %s", secretCommandStr)
		secretStr, err := config.String("server.remote.secret")
		if err != nil {
			logs.Error("get remote secret error: %v", err)
			return
		}
		commandStr, err := utils.DecryptAES([]byte(secretStr), secretCommandStr)
		if err != nil {
			logs.Error("message DecryptAES error: %v", err)
			continue
		}
		commandMessage := CommandMessage{}
		err = json.Unmarshal([]byte(commandStr), &commandMessage)
		if err != nil {
			logs.Error("message format error: %v", err)
			continue
		}

		switch commandMessage.MessageType {
		case "cameraAq":
			cameraAq(commandMessage.Param)
		case "historyVideoPage":
			historyVideoPage(commandMessage.Param)
		case "flvFileMediaInfo":
			flvFileMediaInfo(commandMessage.Param)
		case "flvPlay":
			flvPlay(commandMessage.Param)
		case "flvFetchMoreData":
			flvFetchMoreData(commandMessage.Param)
		default:
			logs.Error("unsupport commandType: %s", commandMessage.MessageType)
		}
	}
}
