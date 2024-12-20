package ext

import (
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
	ext_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/ext"
)

type ModifyPasswordParams struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
}

func ChangePassword(ctx *gin.Context) {
	modifyPwdParam := ModifyPasswordParams{}
	err := ctx.BindJSON(&modifyPwdParam)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	if modifyPwdParam.Password == "" {
		logs.Error("new password is empty")
		result := common.ErrorResult(fmt.Sprintf("new password is empty"))
		ctx.JSON(http.StatusOK, result)
		return
	}

	condition := common.GetEqualCondition("account", modifyPwdParam.Username)
	user, err := base_service.UserFindOneByCondition(condition)
	if err != nil {
		logs.Error("find user error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	if err != nil {
		logs.Error("query by id error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	security, err := config.Bool("server.security")
	if err != nil {
		logs.Error("get server security error: %v. \n use default true", err)
		security = true
	}
	if security {
		if utils.Md5(modifyPwdParam.OldPassword) != user.UserPwd {
			logs.Error("userName : %s , password error", user.Account)
			result := common.ErrorResult("old password error")
			ctx.JSON(http.StatusOK, result)
			return
		}
	}

	newPwd := utils.Md5(modifyPwdParam.Password)
	user.UserPwd = newPwd
	_, err = base_service.UserUpdateById(user)
	if err != nil {
		logs.Error("insert error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	ext_service.TokenDeleteByUsername(user.Account)

	result := common.SuccessResultMsg("user password change success, please relogin")
	ctx.JSON(http.StatusOK, result)
}
