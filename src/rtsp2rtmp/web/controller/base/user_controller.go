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
	user_po "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/po/base/user"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func UserAdd(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := user_po.UserPO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	user, err := dto_convert.ConvertPOToUser(po)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	id, _ := utils.GenerateId()
	user.IdUser = id
	_, err = base_service.UserCreate(user)
	if err != nil {
		logs.Error("insert error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	userAfterSave, err := base_service.UserSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertUserToVO(userAfterSave)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func UserUpdate(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := user_po.UserPO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	user, err := dto_convert.ConvertPOToUser(po)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	var id = user.IdUser

	userById, err := base_service.UserSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	user.UserPwd = userById.UserPwd
	_, err = base_service.UserUpdateById(user)
	if err != nil {
		logs.Error("insert error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	userAfterSave, err := base_service.UserSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	vo, err := dto_convert.ConvertUserToVO(userAfterSave)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func UserRemove(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	po := user_po.UserPO{}
	err := ctx.BindJSON(&po)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	var id = po.IdUser

	userGetById, err := base_service.UserSelectById(id)
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	_, err = base_service.UserDelete(userGetById)
	if err != nil {
		logs.Error("delete error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(userGetById)
	ctx.JSON(http.StatusOK, result)
}

func UserBatchRemove(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	poes := []user_po.UserPO{}
	err := ctx.BindJSON(&poes)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	users, err := dto_convert.ConvertPOListToUser(poes)
	_, err = base_service.UserBatchDelete(users)
	if err != nil {
		logs.Error("delete error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultMsg("remove success")
	ctx.JSON(http.StatusOK, result)
}

func UserGetById(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	id, ok := ctx.Params.Get("id")
	if !ok {
		logs.Error("get param id failed")
		result := common.ErrorResult("get param id failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	user, err := base_service.UserSelectById(id)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	vo, err := dto_convert.ConvertUserToVO(user)
	if err != nil {
		logs.Error("getById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(vo)
	ctx.JSON(http.StatusOK, result)
}

func UserGetByIds(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	idsStr := ctx.Query("ids")
	idList := strings.Split(idsStr, ",")
	if len(idList) == 0 {
		logs.Error("get param ids failed")
		result := common.ErrorResult("get param ids failed")
		ctx.JSON(http.StatusOK, result)
		return
	}
	users, err := base_service.UserSelectByIds(idList)
	if err != nil {
		logs.Error("getByIds error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	voList, err := dto_convert.ConvertUserToVOList(users)
	if err != nil {
		logs.Error("getByIds error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(voList)
	ctx.JSON(http.StatusOK, result)
}

func UserAq(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	condition := common.AqCondition{}
	err := ctx.BindJSON(&condition)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	users, err := base_service.UserFindCollectionByCondition(condition)
	if err != nil {
		logs.Error("aq error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	voList, err := dto_convert.ConvertUserToVOList(users)
	if err != nil {
		logs.Error("aq error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(voList)
	ctx.JSON(http.StatusOK, result)
}

func UserAqPage(ctx *gin.Context) {
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
	pageInfo, err := base_service.UserFindPageByCondition(pageInfoInput)
	if err != nil {
		logs.Error("aqPage error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	var users = make([]entity.User, 0)
	for _, data := range pageInfo.DataList {
		users = append(users, data.(entity.User))
	}
	voList, err := dto_convert.ConvertUserToVOList(users)
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
