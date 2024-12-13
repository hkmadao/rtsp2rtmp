package web

import (
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_controller "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controllers/base"
	ext_controller "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controllers/ext"
)

var tokens sync.Map

func ClearExipresToken() {
	deleteTokens := []string{}
	// 遍历所有sync.Map中的键值对
	tokens.Range(func(k, v interface{}) bool {
		if time.Now().After(v.(time.Time).Add(30 * time.Minute)) {
			deleteTokens = append(deleteTokens, k.(string))
		}
		return true
	})
	for _, v := range deleteTokens {
		tokens.Delete(v)
	}
}

var webInstance *web

type web struct{}

func init() {
	webInstance = &web{}

}

func GetSingleWeb() *web {
	return webInstance
}

func (w *web) StartWeb() {
	go w.webRun()
}

func (w *web) webRun() {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()
	router := gin.Default()
	router.Use(Cors())
	router.Use(Validate())

	router.POST("/system/login", login)

	router.GET("/live/:method/:code/:authCode.flv", ext_controller.HttpFlvPlay)
	// camera
	router.POST("/camera/add", base_controller.CameraAdd)
	router.POST("/camera/update", base_controller.CameraUpdate)
	router.POST("/camera/remove", base_controller.CameraRemove)
	router.GET("/camera/getById/:id", base_controller.CameraGetById)
	router.GET("/camera/getByIds", base_controller.CameraGetByIds)
	router.POST("/camera/aq", base_controller.CameraAq)
	router.POST("/camera/aqPage", base_controller.CameraAqPage)
	// camerashare
	router.POST("/cameraShare/add", base_controller.CameraShareAdd)
	router.POST("/cameraShare/update", base_controller.CameraShareUpdate)
	router.POST("/cameraShare/remove", base_controller.CameraShareRemove)
	router.GET("/cameraShare/getById/:id", base_controller.CameraShareGetById)
	router.GET("/cameraShare/getByIds", base_controller.CameraShareGetByIds)
	router.POST("/cameraShare/aq", base_controller.CameraShareAq)
	router.POST("/cameraShare/aqPage", base_controller.CameraShareAqPage)

	// router.GET("/camera/list", ext_controller.CameraList)
	// router.GET("/camera/detail", ext_controller.CameraDetail)
	// router.POST("/camera/edit", ext_controller.CameraEdit)
	// router.POST("/camera/delete/:id", ext_controller.CameraDelete)
	router.POST("/camera/enabled", ext_controller.CameraEnabled)
	router.POST("/camera/rtmpPushChange", ext_controller.RtmpPushChange)
	router.POST("/camera/saveVideoChange", ext_controller.CameraSaveVideoChange)
	router.POST("/camera/liveChange", ext_controller.CameraLiveChange)
	router.POST("/camera/playAuthCodeReset", ext_controller.CameraPlayAuthCodeReset)

	// router.GET("/camerashare/list", ext_controller.CameraShareList)
	// router.POST("/camerashare/edit", ext_controller.CameraShareEdit)
	// router.POST("/camerashare/delete/:id", ext_controller.CameraShareDelete)
	router.POST("/cameraShare/enabled", ext_controller.CameraShareEnabled)

	staticPath, err := config.String("server.http.static.path")
	if err != nil {
		logs.Error("get httpflv staticPath error: %v. \n use default staticPath : ./resources/static", err)
		staticPath = "./resources/static"
	}

	router.StaticFS("/rtsp2rtmp", http.Dir(staticPath))

	port, err := config.Int("server.http.port")
	if err != nil {
		logs.Error("get httpflv port error: %v. \n use default port : 9090", err)
		port = 9090
	}
	err = router.Run(":" + strconv.Itoa(port))
	if err != nil {
		logs.Error("Start HTTP Server error", err)
	}
}

// 跨域
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//请求方法
		method := ctx.Request.Method
		//请求头部
		origin := ctx.Request.Header.Get("Origin")
		if origin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			// 这是允许访问所有域
			ctx.Header("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			// header的类型
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// 跨域关键设置 让浏览器可以解析
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			// 缓存请求信息 单位为秒
			ctx.Header("Access-Control-Max-Age", "172800")
			//  跨域请求是否需要带cookie信息 默认设置为true
			ctx.Header("Access-Control-Allow-Credentials", "false")
			// 设置返回格式是json
			ctx.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		ctx.Next()
	}
}

// 验证token
func Validate() gin.HandlerFunc {
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
		if ctx.Request.URL.Path == "/system/login" || strings.HasPrefix(ctx.Request.URL.Path, "/live/") ||
			strings.HasPrefix(ctx.Request.URL.Path, "/rtsp2rtmp") {
			ctx.Next()
			return
		}
		token := ctx.Request.Header.Get("token")
		if len(token) == 0 {
			logs.Error("token is empty")
			result := common.ErrorResult("token is empty")
			ctx.JSON(http.StatusUnauthorized, result)
			ctx.Abort()
			return
		}
		tokenTime, b := tokens.Load(token)
		if !b {
			logs.Error("token error")
			result := common.ErrorResult("token error")
			ctx.JSON(http.StatusUnauthorized, result)
			ctx.Abort()
			return
		}
		timeout := time.Now().After(tokenTime.(time.Time).Add(30 * time.Minute))
		if timeout {
			logs.Error("token is timeout")
			result := common.ErrorResult("token is timeout")
			ctx.JSON(http.StatusUnauthorized, result)
			ctx.Abort()
			return
		}
		tokens.Store(token, time.Now())
		ctx.Next()
	}
}

func login(ctx *gin.Context) {
	params := make(map[string]interface{})
	err := ctx.BindJSON(&params)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	userNameParam := params["userName"].(string)
	passwordParam := params["password"].(string)
	userName := config.DefaultString("server.user.name", "")
	password := config.DefaultString("server.user.password", "")
	if userNameParam == "" || passwordParam == "" || userNameParam != userName || passwordParam != password {
		logs.Error("userName : %s , password : %s error", userNameParam, passwordParam)
		result := common.ErrorResult("userName or password error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	logs.Info("用户[%s]登录成功！", userName)
	token, err := utils.NextToke()
	if err != nil {
		logs.Error("create token fail")
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	result := common.SuccessResultWithMsg("succss", map[string]string{"token": token})
	tokens.Store(token, time.Now())
	ctx.JSON(http.StatusOK, result)
}
