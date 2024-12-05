package controllers

import (
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/result"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service"
)

func HttpFlvPlay(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Connection", "keep-alive")
	uri := strings.TrimSuffix(strings.TrimLeft(c.Request.RequestURI, "/"), ".flv")
	uris := strings.Split(uri, "/")
	if len(uris) < 3 || uris[0] != "live" {
		http.Error(c.Writer, "invalid path", http.StatusBadRequest)
		return
	}
	method := uris[1]
	code := uris[2]
	authCode := uris[3]
	r := result.Result{
		Code: 1,
		Msg:  "",
	}
	q := entity.Camera{Code: code}
	camera, err := service.CameraSelectOne(q)
	if err != nil {
		logs.Error("camera query error : %v", err)
		r.Code = 0
		r.Msg = "camera query error"
		c.JSON(http.StatusBadRequest, r)
		return
	}
	if !(method == "temp" || method == "permanent") {
		logs.Error("method error : %s", method)
		r.Code = 0
		r.Msg = "method error"
		c.JSON(http.StatusBadRequest, r)
		return
	}
	if method == "temp" {
		csq := entity.CameraShare{CameraId: camera.Id, AuthCode: authCode}
		cs, err := service.CameraShareSelectOne(csq)
		if err != nil {
			logs.Error("CameraShareSelectOne error : %v", err)
			r.Code = 0
			r.Msg = "system error"
			c.JSON(http.StatusBadRequest, r)
			return
		}
		if time.Now().Before(cs.StartTime) || time.Now().After(cs.Deadline) {
			logs.Error("camera [%s] AuthCodeTemp expired : %s", camera.Code, authCode)
			r.Code = 0
			r.Msg = "authCode expired"
			c.JSON(http.StatusBadRequest, r)
			return
		}

	}
	if method == "permanent" && authCode != camera.PlayAuthCode {
		logs.Error("AuthCodePermanent error : %s", authCode)
		r.Code = 0
		r.Msg = "authCode error"
		c.JSON(http.StatusBadRequest, r)
		return
	}

	logs.Info("player [%s] addr [%s] connecting", code, c.Request.RemoteAddr)
	//管理员可以主动中断播放
	playerDone := make(chan int)
	defer close(playerDone)
	const timeout = 10 * time.Second
	flvPlayerDone, err := flvadmin.GetSingleHttpFlvAdmin().AddHttpFlvPlayer(playerDone, timeout/2, code, c.Writer)
	if err != nil {
		logs.Error("camera [%s] add player error : %s", code, err)
		r.Code = 0
		r.Msg = "add player error"
		c.JSON(http.StatusBadRequest, r)
		return
	}
	<-flvPlayerDone
	logs.Info("player [%s] addr [%s] exit", code, c.Request.RemoteAddr)
}
