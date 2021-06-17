package controllers

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/yumrano/rtsp2rtmp/models"
	"github.com/yumrano/rtsp2rtmp/result"
	"github.com/yumrano/rtsp2rtmp/writer/httpflv"
)

func Live(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	cros(w, req)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connection", "keep-alive")
	uri := strings.TrimSuffix(strings.TrimLeft(req.RequestURI, "/"), ".flv")
	uris := strings.Split(uri, "/")
	if len(uris) < 3 || uris[0] != "live" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	method := uris[1]
	code := uris[2]
	authCode := uris[3]
	r := result.Result{
		Code: 1,
		Msg:  "",
	}
	q := models.Camera{Code: code}
	camera, err := models.CameraSelectOne(q)
	if err != nil {
		logs.Error("camera query error : %v", err)
		r.Code = 0
		r.Msg = "camera query error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	if !(method == "temp" || method == "permanent") {
		logs.Error("method error : %s", method)
		r.Code = 0
		r.Msg = "method error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	// if method == "temp" {
	// 	csq := models.CameraShare{CameraId: camera.Id, AuthCode: authCode}
	// 	cs, err := models.CameraShareSelectOne(csq)
	// 	if err != nil {
	// 		logs.Error("CameraShareSelectOne error : %v", err)
	// 		r.Code = 0
	// 		r.Msg = "system error"
	// 		rbytes, _ := json.Marshal(r)
	// 		w.Write(rbytes)
	// 		return
	// 	}
	// 	if time.Now().After(cs.Created.Add(7 * 24 * time.Hour)) {
	// 		logs.Error("camera [%s] AuthCodeTemp expired : %s", camera.Code, authCode)
	// 		r.Code = 0
	// 		r.Msg = "authCode expired"
	// 		rbytes, _ := json.Marshal(r)
	// 		w.Write(rbytes)
	// 		return
	// 	}

	// }
	if method == "permanent" && authCode != camera.AuthCode {
		logs.Error("AuthCodePermanent error : %s", authCode)
		r.Code = 0
		r.Msg = "authCode error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	// if !server.ExistCamera(code) {
	// 	logs.Error("camera [%s] no connection", code)
	// 	r.Code = 0
	// 	r.Msg = "camera no connection"
	// 	rbytes, _ := json.Marshal(r)
	// 	w.Write(rbytes)
	// 	return
	// }
	logs.Info("player [%s] addr [%s] connecting", code, req.RemoteAddr)
	//管理员可以主动中断播放
	endStream, heartbeatStream, _, err := httpflv.AddHttpFlvPlayer(code, w)
	if err != nil {
		logs.Error("camera [%s] add player error : %s", code)
		r.Code = 0
		r.Msg = "add player error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
Loop:
	for {
		select {
		case <-endStream:
			logs.Info("player [%s] addr [%s] end", code, req.RemoteAddr)
			break Loop
		case <-heartbeatStream:
			// logs.Info("player [%s] addr [%s] continue", code, req.RemoteAddr)
			continue
		case <-time.After(10 * time.Second):
			logs.Info("player [%s] addr [%s] timeout", code, req.RemoteAddr)
			break Loop
		}
	}
	logs.Info("player [%s] addr [%s] exit", code, req.RemoteAddr)
}
