package ext

import (
	"net/http"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

type LoginResult struct {
	NickName    string    `json:"nickName"`
	Username    string    `json:"username"`
	Token       string    `json:"token"`
	ExpiredTime time.Time `json:"expiredTime"`
}

func Login(ctx *gin.Context) {
	params := make(map[string]interface{})
	err := ctx.BindJSON(&params)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	userNameParam := params["username"].(string)
	passwordParam := params["password"].(string)
	if userNameParam == "" || passwordParam == "" {
		logs.Error("username or passowrd is empty")
		result := common.ErrorResult("username or passowrd is empty")
		ctx.JSON(http.StatusOK, result)
		return
	}
	condition := common.GetEqualCondition("account", userNameParam)
	user, err := base_service.UserFindOneByCondition(condition)
	if err != nil {
		logs.Error("find user error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	if user.UserPwd == "" {
		logs.Error("user: %s password is empty", user.Account)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	password := user.UserPwd
	if utils.Md5(passwordParam) != password {
		logs.Error("userName : %s , password error", user.Account)
		result := common.ErrorResult("userName or password error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	token, err := utils.GenarateRandStr(32)
	if err != nil {
		logs.Error("user: %s create token fail: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	sysToken := entity.Token{}
	sysToken.CreateTime = time.Now()
	idToken, _ := utils.GenerateId()
	sysToken.IdToken = idToken
	sysToken.ExpiredTime = time.Now().Add(1 * time.Hour)
	sysToken.NickName = user.NickName
	sysToken.Username = user.Account
	sysToken.Token = token
	_, err = base_service.TokenCreate(sysToken)
	if err != nil {
		logs.Error("user: %s store token fail", user.Account)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	logs.Info("user: %s login success", user.Account)
	loginResult := LoginResult{}
	loginResult.ExpiredTime = sysToken.ExpiredTime
	loginResult.NickName = sysToken.NickName
	loginResult.Username = sysToken.Username
	loginResult.Token = sysToken.Token
	result := common.SuccessResultWithMsg("succss", loginResult)
	ctx.JSON(http.StatusOK, result)
}

func Logout(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	if len(token) == 0 {
		logs.Error("token is empty")
		result := common.ErrorResult("token is empty")
		ctx.JSON(http.StatusUnauthorized, result)
		return
	}
	condition := common.GetEqualCondition("token", token)
	sysToken, err := base_service.TokenFindOneByCondition(condition)
	if err != nil {
		logs.Error("find sysToken error : %v", err)
		result := common.ErrorResult("token error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	_, err = base_service.TokenDelete(sysToken)
	if err != nil {
		logs.Error("user: %s token delete error : %v", sysToken.Username, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	logs.Info("user: %s logout success", sysToken.Username)
	result := common.SuccessResultMsg("logout success, please relogin")
	ctx.JSON(http.StatusOK, result)
}

// 验证token
func TokenValidate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		security, err := config.Bool("server.security")
		if err != nil {
			logs.Error("get server security error: %v. \n use default true", err)
			security = true
		}
		if !security {
			ctx.Next()
			return
		}
		if ctx.Request.URL.Path == "/login" || strings.HasPrefix(ctx.Request.URL.Path, "/live/") ||
			strings.HasPrefix(ctx.Request.URL.Path, "/rtsp2rtmp") {
			ctx.Next()
			return
		}
		token := ctx.Request.Header.Get("Authorization")
		if len(token) == 0 {
			logs.Error("token is empty")
			result := common.ErrorResult("token is empty")
			ctx.JSON(http.StatusUnauthorized, result)
			ctx.Abort()
			return
		}
		condition := common.GetEqualCondition("token", token)
		sysToken, err := base_service.TokenFindOneByCondition(condition)
		if err != nil {
			logs.Error("find sysToken error : %v", err)
			result := common.ErrorResult("token error")
			ctx.JSON(http.StatusOK, result)
			ctx.Abort()
			return
		}

		timeout := time.Now().After(sysToken.ExpiredTime)
		if timeout {
			logs.Error("token is expired")
			result := common.ErrorResult("token is expired")
			ctx.JSON(http.StatusUnauthorized, result)
			ctx.Abort()
			return
		}
		sysToken.ExpiredTime = time.Now().Add(1 * time.Hour)
		_, err = base_service.TokenUpdateById(sysToken)
		if err != nil {
			logs.Error("user: %s update token fail", sysToken.Username)
			result := common.ErrorResult("internal error")
			ctx.JSON(http.StatusOK, result)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
