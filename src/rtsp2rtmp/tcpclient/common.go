package tcpclient

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
	// "cameraAq" "historyVideoPage" "flvFileMediaInfo" "flvPlay" "flvFetchMoreData" "startPushRtmp" "stopPushRtmp"
	MessageType string `json:"messageType"`
	Param       string `json:"param"`
}

// when connect to server, first send register packet to server
type RegisterInfo struct {
	ClientId string `json:"clientId"`
	DateStr  string `json:"dateStr"`
	Sign     string `json:"sign"`
	// "keepChannel" "cameraAq" "historyVideoPage" "flvFileMediaInfo" "flvPlay" "flvFetchMoreData" "startPushRtmp" "stopPushRtmp"
	ConnType string `json:"connType"`
}

func newReisterInfo(connType string) (ri RegisterInfo, err error) {
	currentDateStr := time.Now().Format(time.RFC3339)
	clientId, err := config.String("server.remote.client-id")
	if err != nil {
		logs.Error("get remote client-id error: %v\n", err)
		return
	}
	signSecret, err := config.String("server.remote.sign-secret")
	if err != nil {
		logs.Error("get remote sign-secret error: %v\n", err)
		return
	}
	planText := fmt.Sprintf("cilentId=%s&dateStr=%s", clientId, currentDateStr)
	signStr, err := utils.EncryptAES([]byte(signSecret), planText)
	if err != nil {
		err = fmt.Errorf("buildSign error: %v", err)
		return
	}

	ri = RegisterInfo{
		ClientId: clientId,
		ConnType: connType,
		DateStr:  currentDateStr,
		Sign:     signStr,
	}
	return
}

func connectAndRegister(connType string) (conn net.Conn, err error) {
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
	defer func() {
		conn.Close()
	}()

	// register to server
	ri, err := newReisterInfo(connType)
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

func writeResult(result common.AppResult, writer io.Writer) (n int, err error) {
	messageBytes, err := json.Marshal(result)
	if err != nil {
		logs.Error(err)
		return
	}
	secretStr, err := config.String("server.remote.secret")
	if err != nil {
		logs.Error("get remote secret error: %v", err)
		return
	}
	encryptMessageStr, err := utils.EncryptAES([]byte(secretStr), string(messageBytes))
	if err != nil {
		logs.Error("EncryptAES error: %v", err)
		return
	}
	encryptMessageBytes := string(encryptMessageStr)
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
