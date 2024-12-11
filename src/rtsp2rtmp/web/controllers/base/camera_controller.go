package base

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func CameraGetById(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	id, ok := ctx.Params.Get("id")
	if !ok {
		logs.Error("get param id failed")
		result := common.ErrorResult("get param id failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	camera, err := base_service.CameraSelectById(id)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(camera)
	ctx.JSON(http.StatusOK, result)
}

func CameraGetByIds(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	idsStr := ctx.Query("ids")
	idList := strings.Split(idsStr, ",")
	if len(idList) == 0 {
		logs.Error("get param ids failed")
		result := common.ErrorResult("get param ids failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	camera, err := base_service.CameraSelectByIds(idList)
	if err != nil {
		logs.Error("getByIds error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(camera)
	ctx.JSON(http.StatusOK, result)
}

func CameraAq(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	condition := common.AqCondition{}
	err := ctx.BindJSON(&condition)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameras, err := base_service.CameraFindCollectionByCondition(condition)
	if err != nil {
		logs.Error("aq error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(cameras)
	ctx.JSON(http.StatusOK, result)
}

func CameraAqPage(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	pageInfoInput := common.AqPageInfoInput{}
	err := ctx.BindJSON(&pageInfoInput)
	if err != nil {
		ctx.AbortWithError(500, err)
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	pageInfo, err := base_service.CameraFindPageByCondition(pageInfoInput)
	if err != nil {
		logs.Error("aqPage error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(pageInfo)
	ctx.JSON(http.StatusOK, result)
}
