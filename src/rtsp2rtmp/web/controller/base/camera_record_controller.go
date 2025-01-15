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
	camera_record_po "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/po/base/camera_record"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func CameraRecordAdd(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := camera_record_po.CameraRecordPO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraRecord, err := dto_convert.ConvertPOToCameraRecord(po)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	id, _ := utils.GenerateId()
	cameraRecord.IdCameraRecord = id
	_, err = base_service.CameraRecordCreate(cameraRecord)
	if err != nil {
		logs.Error("insert error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	cameraRecordAfterSave, err := base_service.CameraRecordSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertCameraRecordToVO(cameraRecordAfterSave)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func CameraRecordUpdate(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := camera_record_po.CameraRecordPO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraRecord, err := dto_convert.ConvertPOToCameraRecord(po)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	var id = cameraRecord.IdCameraRecord

	_, err = base_service.CameraRecordSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	_, err = base_service.CameraRecordUpdateById(cameraRecord)
	if err != nil {
		logs.Error("insert error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	cameraRecordAfterSave, err := base_service.CameraRecordSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertCameraRecordToVO(cameraRecordAfterSave)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func CameraRecordRemove(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := camera_record_po.CameraRecordPO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	var id = po.IdCameraRecord

	cameraRecordGetById, err := base_service.CameraRecordSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	_, err = base_service.CameraRecordDelete(cameraRecordGetById)
	if err != nil {
		logs.Error("delete error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(cameraRecordGetById)
	ctx.JSON(http.StatusOK, result)
}

func CameraRecordBatchRemove(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	poes := []camera_record_po.CameraRecordPO{}
	err := ctx.BindJSON(&poes)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraRecords, err := dto_convert.ConvertPOListToCameraRecord(poes)
	_, err = base_service.CameraRecordBatchDelete(cameraRecords)
	if err != nil {
		logs.Error("delete error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultMsg("remove success")
	ctx.JSON(http.StatusOK, result)
}

func CameraRecordGetById(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	id, ok := ctx.Params.Get("id")
	if !ok {
		logs.Error("get param id failed")
		result := common.ErrorResult("get param id failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraRecord, err := base_service.CameraRecordSelectById(id)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	vo, err := dto_convert.ConvertCameraRecordToVO(cameraRecord)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func CameraRecordGetByIds(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	idsStr := ctx.Query("ids")
	idList := strings.Split(idsStr, ",")
	if len(idList) == 0 {
		logs.Error("get param ids failed")
		result := common.ErrorResult("get param ids failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraRecords, err := base_service.CameraRecordSelectByIds(idList)
	if err != nil {
		logs.Error("getByIds error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	voList, err := dto_convert.ConvertCameraRecordToVOList(cameraRecords)
	if err != nil {
		logs.Error("getByIds error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(voList)
	ctx.JSON(http.StatusOK, result)
}

func CameraRecordAq(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	condition := common.AqCondition{}
	err := ctx.BindJSON(&condition)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	cameraRecords, err := base_service.CameraRecordFindCollectionByCondition(condition)
	if err != nil {
		logs.Error("aq error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	voList, err := dto_convert.ConvertCameraRecordToVOList(cameraRecords)
	if err != nil {
		logs.Error("aq error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(voList)
	ctx.JSON(http.StatusOK, result)
}

func CameraRecordAqPage(ctx *gin.Context) {
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
	pageInfo, err := base_service.CameraRecordFindPageByCondition(pageInfoInput)
	if err != nil {
		logs.Error("aqPage error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	var cameraRecords = make([]entity.CameraRecord, 0)
	for _, data := range pageInfo.DataList {
		cameraRecords = append(cameraRecords, data.(entity.CameraRecord))
	}
	voList, err := dto_convert.ConvertCameraRecordToVOList(cameraRecords)
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
