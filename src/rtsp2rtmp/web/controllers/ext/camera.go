package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	bas_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func CameraEnabled(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.Camera{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	camera, err := bas_service.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	camera.Enabled = q.Enabled
	if q.Enabled != 1 {
		camera.OnlineStatus = 0
	}
	_, err = bas_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("enabled camera status %d error : %v", camera.Enabled, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	if q.Enabled != 1 {
		logs.Debug("close camera conn: %s", camera.Code)
		select {
		case codeStream <- camera.Code:
		case <-time.After(1 * time.Second):
		}
	}

	result := common.SuccessResultWithMsg("succss", camera)
	ctx.JSON(http.StatusOK, result)
}

func RtmpPushChange(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.Camera{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	camera, err := bas_service.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	camera.RtmpPushStatus = q.RtmpPushStatus
	_, err = bas_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("RtmpPushEnabled camera status %d error : %v", camera.Enabled, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	switch {
	case q.RtmpPushStatus != 1:
		logs.Info("camera [%s] stop push rtmp", q.Code)
		flvadmin.GetSingleRtmpFlvAdmin().StopWrite(q.Code)
	case q.RtmpPushStatus == 1:
		flvadmin.GetSingleRtmpFlvAdmin().StartWrite(q.Code)
		logs.Info("camera [%s] start push rtmp", q.Code)
	}

	result := common.SuccessResultWithMsg("succss", camera)
	ctx.JSON(http.StatusOK, result)
}

func CameraSaveVideoChange(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.Camera{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	camera, err := bas_service.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	camera.SaveVideo = q.SaveVideo
	_, err = bas_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("SaveVideo camera status %d error : %v", camera.SaveVideo, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	switch {
	case q.SaveVideo != 1:
		logs.Info("camera [%s] stop save video", q.Code)
		flvadmin.GetSingleFileFlvAdmin().StopWrite(q.Code)
	case q.SaveVideo == 1:
		flvadmin.GetSingleFileFlvAdmin().StartWrite(q.Code)
		logs.Info("camera [%s] start save video", q.Code)
	}

	result := common.SuccessResultWithMsg("succss", camera)
	ctx.JSON(http.StatusOK, result)
}

func CameraLiveChange(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.Camera{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	camera, err := bas_service.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	camera.Live = q.Live
	_, err = bas_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("Live camera status %d error : %v", camera.Live, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	switch {
	case q.Live != 1:
		flvadmin.GetSingleHttpFlvAdmin().StopWrite(q.Code)
	case q.Live == 1:
		flvadmin.GetSingleHttpFlvAdmin().StartWrite(q.Code)
	}

	result := common.SuccessResultWithMsg("succss", camera)
	ctx.JSON(http.StatusOK, result)
}

func CameraPlayAuthCodeReset(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.Camera{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	camera, err := bas_service.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	playAuthCode, _ := utils.UUID()
	camera.PlayAuthCode = playAuthCode
	_, err = bas_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("PlayAuthCode camera status %d error : %v", camera.PlayAuthCode, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	flvadmin.GetSingleHttpFlvAdmin().StopWrite(q.Code)
	flvadmin.GetSingleHttpFlvAdmin().StartWrite(q.Code)

	result := common.SuccessResultWithMsg("succss", camera)
	ctx.JSON(http.StatusOK, result)
}

var codeStream = make(chan string)

func CodeStream() <-chan string {
	return codeStream
}
