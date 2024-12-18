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
	token_po "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/po/base/token"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func TokenAdd(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := token_po.TokenPO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	token, err := dto_convert.ConvertPOToToken(po)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	id, _ := utils.UUID()
	token.IdToken = id
	_, err = base_service.TokenCreate(token)
	if err != nil {
		logs.Error("insert error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	tokenAfterSave, err := base_service.TokenSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertTokenToVO(tokenAfterSave)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func TokenUpdate(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := token_po.TokenPO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	token, err := dto_convert.ConvertPOToToken(po)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	var id = token.IdToken

	_, err = base_service.TokenSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	_, err = base_service.TokenUpdateById(token)
	if err != nil {
		logs.Error("insert error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	tokenAfterSave, err := base_service.TokenSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertTokenToVO(tokenAfterSave)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func TokenRemove(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := token_po.TokenPO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	var id = po.IdToken

	tokenGetById, err := base_service.TokenSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	_, err = base_service.TokenDelete(tokenGetById)
	if err != nil {
		logs.Error("delete error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(tokenGetById)
	ctx.JSON(http.StatusOK, result)
}

func TokenGetById(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	id, ok := ctx.Params.Get("id")
	if !ok {
		logs.Error("get param id failed")
		result := common.ErrorResult("get param id failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	token, err := base_service.TokenSelectById(id)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	vo, err := dto_convert.ConvertTokenToVO(token)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func TokenGetByIds(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	idsStr := ctx.Query("ids")
	idList := strings.Split(idsStr, ",")
	if len(idList) == 0 {
		logs.Error("get param ids failed")
		result := common.ErrorResult("get param ids failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	tokens, err := base_service.TokenSelectByIds(idList)
	if err != nil {
		logs.Error("getByIds error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	voList, err := dto_convert.ConvertTokenToVOList(tokens)
	if err != nil {
		logs.Error("getByIds error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(voList)
	ctx.JSON(http.StatusOK, result)
}

func TokenAq(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	condition := common.AqCondition{}
	err := ctx.BindJSON(&condition)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	tokens, err := base_service.TokenFindCollectionByCondition(condition)
	if err != nil {
		logs.Error("aq error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	voList, err := dto_convert.ConvertTokenToVOList(tokens)
	if err != nil {
		logs.Error("aq error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(voList)
	ctx.JSON(http.StatusOK, result)
}

func TokenAqPage(ctx *gin.Context) {
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
	pageInfo, err := base_service.TokenFindPageByCondition(pageInfoInput)
	if err != nil {
		logs.Error("aqPage error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	var tokens = make([]entity.Token, 0)
	for _, data := range pageInfo.DataList {
		tokens = append(tokens, data.(entity.Token))
	}
	voList, err := dto_convert.ConvertTokenToVOList(tokens)
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
