package web

import (
	"net/http"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	base_controller "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controller/base"
	ext_controller "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controller/ext"
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
	router.Use(ext_controller.TokenValidate())

	router.POST("/login", ext_controller.Login)
	router.POST("/logout", ext_controller.Logout)

	router.GET("/live/:method/:code/:authCode.flv", ext_controller.HttpFlvPlay)
	// user
	router.POST("/user/updatePw", ext_controller.ChangePassword)
	router.POST("/user/add", base_controller.UserAdd)
	router.POST("/user/update", base_controller.UserUpdate)
	router.POST("/user/remove", base_controller.UserRemove)
	router.POST("/user/batchRemove", base_controller.UserBatchRemove)
	router.GET("/user/getById/:id", base_controller.UserGetById)
	router.GET("/user/getByIds", base_controller.UserGetByIds)
	router.POST("/user/aq", base_controller.UserAq)
	router.POST("/user/aqPage", base_controller.UserAqPage)
	// toke
	router.POST("/token/add", base_controller.TokenAdd)
	router.POST("/token/update", base_controller.TokenUpdate)
	router.POST("/token/remove", base_controller.TokenRemove)
	router.POST("/token/batchRemove", base_controller.TokenBatchRemove)
	router.GET("/token/getById/:id", base_controller.TokenGetById)
	router.GET("/token/getByIds", base_controller.TokenGetByIds)
	router.POST("/token/aq", base_controller.TokenAq)
	router.POST("/token/aqPage", base_controller.TokenAqPage)
	// camera
	router.POST("/camera/add", base_controller.CameraAdd)
	router.POST("/camera/update", base_controller.CameraUpdate)
	router.POST("/camera/remove", base_controller.CameraRemove)
	router.POST("/camera/batchRemove", base_controller.CameraBatchRemove)
	router.GET("/camera/getById/:id", base_controller.CameraGetById)
	router.GET("/camera/getByIds", base_controller.CameraGetByIds)
	router.POST("/camera/aq", base_controller.CameraAq)
	router.POST("/camera/aqPage", base_controller.CameraAqPage)
	router.POST("/camera/enabled", ext_controller.CameraEnabled)
	router.POST("/camera/rtmpPushChange", ext_controller.RtmpPushChange)
	router.POST("/camera/saveVideoChange", ext_controller.CameraSaveVideoChange)
	router.POST("/camera/liveChange", ext_controller.CameraLiveChange)
	router.POST("/camera/playAuthCodeReset", ext_controller.CameraPlayAuthCodeReset)
	// camerashare
	router.POST("/cameraShare/add", base_controller.CameraShareAdd)
	router.POST("/cameraShare/update", base_controller.CameraShareUpdate)
	router.POST("/cameraShare/remove", base_controller.CameraShareRemove)
	router.POST("/cameraShare/batchRemove", base_controller.CameraShareBatchRemove)
	router.GET("/cameraShare/getById/:id", base_controller.CameraShareGetById)
	router.GET("/cameraShare/getByIds", base_controller.CameraShareGetByIds)
	router.POST("/cameraShare/aq", base_controller.CameraShareAq)
	router.POST("/cameraShare/aqPage", base_controller.CameraShareAqPage)
	router.POST("/cameraShare/enabled", ext_controller.CameraShareEnabled)
	router.POST("/cameraShare/playAuthCodeReset", ext_controller.CameraSharePlayAuthCodeReset)

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
