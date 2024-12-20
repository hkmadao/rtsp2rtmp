package ext

import (
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	dto_convert "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controller/convert"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func CameraShareEnabled(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.CameraShare{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	cameraShare, err := base_service.CameraShareSelectById(q.Id)
	if err != nil {
		logs.Error("query camerashare error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraShare.Enabled = q.Enabled
	_, err = base_service.CameraShareUpdateById(cameraShare)
	if err != nil {
		logs.Error("enabled camerashare status %d error : %v", cameraShare.Enabled, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertCameraShareToVO(cameraShare)
	if err != nil {
		logs.Error("CameraShareEnabled error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	result := common.SuccessResultWithMsg("succss", vo)
	ctx.JSON(http.StatusOK, result)
}

func CameraSharePlayAuthCodeReset(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.CameraShare{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	cameraShare, err := base_service.CameraShareSelectById(q.Id)
	if err != nil {
		logs.Error("query cameraShare error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	playAuthCode := utils.GenarateRandName()
	cameraShare.AuthCode = playAuthCode
	_, err = base_service.CameraShareUpdateById(cameraShare)
	if err != nil {
		logs.Error("CameraShare: %s PlayAuthCode reset error : %v", cameraShare.Id, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertCameraShareToVO(cameraShare)
	if err != nil {
		logs.Error("CameraSharePlayAuthCodeReset error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	result := common.SuccessResultWithMsg("succss", vo)
	ctx.JSON(http.StatusOK, result)
}
