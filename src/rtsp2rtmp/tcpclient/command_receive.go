package tcpclient

import (
	"encoding/json"
	"io"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
)

func StartCommandReceiveServer() {
	go func() {
		for {
			commandReceiveConnect()
			<-time.NewTicker(10 * time.Second).C
		}
	}()
}

func commandReceiveConnect() {
	conn, err := connectAndRegister("keepChannel", "")
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
		dataLen := utils.BigEndianToUint32(dataLenBytes)

		serverRepBytes := make([]byte, 0)
		for {
			buffer := make([]byte, dataLen-uint32(len(serverRepBytes)))
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					logs.Error("conn read message body error: %v", err)
					return
				}
				break
			}
			// 处理读取到的数据，n是实际读取的字节数
			serverRepBytes = append(serverRepBytes, buffer[:n]...)
			if uint32(len(serverRepBytes)) == dataLen {
				break
			}
		}

		// commandMessage := CommandMessage{}
		// err = json.Unmarshal(serverRepBytes, &commandMessage)
		// if err != nil {
		// 	logs.Error("message format error: %v", err)
		// 	continue
		// }

		secretStr, err := config.String("server.remote.secret")
		if err != nil {
			logs.Error("get remote secret error: %v", err)
			return
		}
		commandBytes, err := utils.DecryptAES([]byte(secretStr), serverRepBytes)
		if err != nil {
			logs.Error("message DecryptAES error: %v", err)
			continue
		}
		commandMessage := CommandMessage{}
		err = json.Unmarshal(commandBytes, &commandMessage)
		if err != nil {
			logs.Error("message format error: %v", err)
			continue
		}

		// do response
		go commandRes(commandMessage)
	}
}

func commandRes(commandMessage CommandMessage) {
	switch commandMessage.MessageType {
	case "cameraAq":
		cameraAq(commandMessage)
	case "historyVideoPage":
		historyVideoPage(commandMessage)
	case "flvFileMediaInfo":
		flvFileMediaInfo(commandMessage)
	case "flvPlay":
		flvPlay(commandMessage)
	case "flvFetchMoreData":
		flvFetchMoreData(commandMessage)
	default:
		logs.Error("unsupport commandType: %s", commandMessage.MessageType)
	}
}
