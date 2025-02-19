package tcpclientcommon

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
)

type CommandMessage struct {
	// "cameraAq" "historyVideoPage" "flvFileMediaInfo" "flvPlay" "flvFetchMoreData" "startPushRtmp" "stopPushRtmp" "getLiveMediaInfo"
	MessageType string `json:"messageType"`
	Param       string `json:"param"`
	MessageId   string `json:"messageId"`
}

// when connect to server, first send register packet to server
type RegisterInfo struct {
	ClientCode string `json:"clientCode"`
	DateStr    string `json:"dateStr"`
	Sign       string `json:"sign"`
	// "keepChannel" "cameraOnline" "cameraOffline" "cameraAq" "historyVideoPage" "flvFileMediaInfo" "flvPlay" "flvFetchMoreData" "startPushRtmp" "stopPushRtmp" "getLiveMediaInfo"
	ConnType   string `json:"connType"`
	MessageId  string `json:"messageId"`
	CameraCode string `json:"cameraCode"`
}

func newReisterInfo(connType string, messageId string, cameraCode string) (ri RegisterInfo, err error) {
	currentDateStr := time.Now().Format(time.RFC3339)
	clientCode, err := config.String("server.remote.client-code")
	if err != nil {
		logs.Error("get remote client-code error: %v\n", err)
		return
	}
	signSecret, err := config.String("server.remote.sign-secret")
	if err != nil {
		logs.Error("get remote sign-secret error: %v\n", err)
		return
	}
	planText := fmt.Sprintf("clientCode=%s&dateStr=%s&signSecret=%s", clientCode, currentDateStr, signSecret)
	signStr := utils.Md5(planText)

	ri = RegisterInfo{
		ClientCode: clientCode,
		ConnType:   connType,
		DateStr:    currentDateStr,
		Sign:       signStr,
		MessageId:  messageId,
		CameraCode: cameraCode,
	}
	return
}

func ConnectAndKeepChannelRegister(connType string) (conn net.Conn, err error) {
	return connectAndRegister(connType, "", "")
}

func ConnectAndCameraStatusRegister(connType string, cameraCode string) (conn net.Conn, err error) {
	return connectAndRegister(connType, "", cameraCode)
}

func ConnectAndResRegister(connType string, messageId string) (conn net.Conn, err error) {
	return connectAndRegister(connType, messageId, "")
}

func connectAndRegister(connType string, messageId string, cameraCode string) (conn net.Conn, err error) {
	serverIp, err := config.String("server.remote.server-ip")
	if err != nil {
		logs.Error("get remote server-ip error: %v. \n", err)
		return
	}
	port, err := config.Int("server.remote.port")
	if err != nil {
		logs.Error("get httpflv port error: %v. \n use default port : 9090", err)
		return
	}
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, port))
	if err != nil {
		logs.Error(err)
		return
	}

	// register to server
	ri, err := newReisterInfo(connType, messageId, cameraCode)
	if err != nil {
		logs.Error(err)
		return
	}
	registerBodyBytes, err := json.Marshal(ri)
	if err != nil {
		logs.Error(err)
		return
	}
	registerBodyLen := len(registerBodyBytes)
	registerBodyLenBytes := utils.Int32ToByteBigEndian(int32(registerBodyLen))
	messageBytes := append(registerBodyLenBytes, registerBodyBytes...)
	_, err = conn.Write(messageBytes)
	if err != nil {
		logs.Error("register error: %v", err)
		return
	}
	return
}

func WriteResult(result common.AppResult, writer io.Writer) (n int, err error) {
	messageBytes, err := json.Marshal(result)
	if err != nil {
		logs.Error(err)
		return
	}
	// messageLenBytes := utils.Int32ToByteBigEndian(int32(len(messageBytes)))
	// fullMessageBytes := append(messageLenBytes, messageBytes...)
	// n, err = writer.Write(fullMessageBytes)
	// if err != nil {
	// 	logs.Error("register error: %v", err)
	// 	return
	// }
	// return

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
	n, err = writer.Write(fullMessageBytes)
	if err != nil {
		logs.Error("register error: %v", err)
		return
	}
	return
}
