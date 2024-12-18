package base

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	dto_convert "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controller/convert"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	camera_share_po "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/po/base/camera_share"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func CameraShareAdd(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := camera_share_po.CameraSharePO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraShare, err := dto_convert.ConvertPOToCameraShare(po)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	id, _ := utils.UUID()
	cameraShare.Id = id
	_, err = base_service.CameraShareCreate(cameraShare)
	if err != nil {
		logs.Error("insert error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	cameraShareAfterSave, err := base_service.CameraShareSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertCameraShareToVO(cameraShareAfterSave)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func CameraShareUpdate(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := camera_share_po.CameraSharePO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraShare, err := dto_convert.ConvertPOToCameraShare(po)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	var id = cameraShare.Id

	_, err = base_service.CameraShareSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	_, err = base_service.CameraShareUpdateById(cameraShare)
	if err != nil {
		logs.Error("insert error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	cameraShareAfterSave, err := base_service.CameraShareSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertCameraShareToVO(cameraShareAfterSave)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func CameraShareRemove(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := camera_share_po.CameraSharePO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	var id = po.Id

	cameraShareGetById, err := base_service.CameraShareSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	_, err = base_service.CameraShareDelete(cameraShareGetById)
	if err != nil {
		logs.Error("delete error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(cameraShareGetById)
	ctx.JSON(http.StatusOK, result)
}

func CameraShareGetById(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	id, ok := ctx.Params.Get("id")
	if !ok {
		logs.Error("get param id failed")
		result := common.ErrorResult("get param id failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraShare, err := base_service.CameraShareSelectById(id)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	vo, err := dto_convert.ConvertCameraShareToVO(cameraShare)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func CameraShareGetByIds(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	idsStr := ctx.Query("ids")
	idList := strings.Split(idsStr, ",")
	if len(idList) == 0 {
		logs.Error("get param ids failed")
		result := common.ErrorResult("get param ids failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraShares, err := base_service.CameraShareSelectByIds(idList)
	if err != nil {
		logs.Error("getByIds error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	voList, err := dto_convert.ConvertCameraShareToVOList(cameraShares)
	if err != nil {
		logs.Error("getByIds error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(voList)
	ctx.JSON(http.StatusOK, result)
}

func CameraShareAq(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	condition := common.AqCondition{}
	err := ctx.BindJSON(&condition)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraShares, err := base_service.CameraShareFindCollectionByCondition(condition)
	if err != nil {
		logs.Error("aq error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	voList, err := dto_convert.ConvertCameraShareToVOList(cameraShares)
	if err != nil {
		logs.Error("aq error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(voList)
	ctx.JSON(http.StatusOK, result)
}

func CameraShareAqPage(ctx *gin.Context) {
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
	pageInfo, err := base_service.CameraShareFindPageByCondition(pageInfoInput)
	if err != nil {
		logs.Error("aqPage error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	var cameraShares = make([]entity.CameraShare, 0)
	for _, data := range pageInfo.DataList {
		cameraShares = append(cameraShares, data.(entity.CameraShare))
	}
	voList, err := dto_convert.ConvertCameraShareToVOList(cameraShares)
	if err != nil {
		logs.Error("aqPage error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	var dataList = make([]interface{}, 0)
	for _, vo := range voList {
		dataList = append(dataList, vo)
	}
	pageInfo.DataList = dataList
	result := common.SuccessResultData(pageInfo)
	ctx.JSON(http.StatusOK, result)
}
