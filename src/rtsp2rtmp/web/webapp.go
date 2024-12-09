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
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controllers"
	base_controller "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/controllers/base"
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
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	router := gin.Default()
	router.Use(Cors())
	router.Use(Validate())

	router.POST("/system/login", login)

	router.GET("/live/:method/:code/:authCode.flv", controllers.HttpFlvPlay)

	router.POST("/camera/aq", base_controller.CameraAq)
	router.GET("/camera/list", controllers.CameraList)
	router.GET("/camera/detail", controllers.CameraDetail)
	router.POST("/camera/edit", controllers.CameraEdit)
	router.POST("/camera/delete/:id", controllers.CameraDelete)
	router.POST("/camera/enabled", controllers.CameraEnabled)
	router.POST("/camera/rtmppushchange", controllers.RtmpPushChange)
	router.POST("/camera/savevideochange", controllers.CameraSaveVideoChange)
	router.POST("/camera/livechange", controllers.CameraLiveChange)
	router.POST("/camera/playauthcodereset", controllers.CameraPlayAuthCodeReset)

	router.GET("/camerashare/list", controllers.CameraShareList)
	router.POST("/camerashare/edit", controllers.CameraShareEdit)
	router.POST("/camerashare/delete/:id", controllers.CameraShareDelete)
	router.POST("/camerashare/enabled", controllers.CameraShareEnabled)

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
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

// 验证token
func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		security, err := config.Bool("server.security")
		if err != nil {
			logs.Error("get server security error: %v. \n use default true", err)
			security = true
		}
		if !security {
			c.Next()
			return
		}
		if c.Request.URL.Path == "/system/login" || strings.HasPrefix(c.Request.URL.Path, "/live/") ||
			strings.HasPrefix(c.Request.URL.Path, "/rtsp2rtmp") {
			c.Next()
			return
		}
		r := common.Result{
			Code: 1,
			Msg:  "",
		}
		token := c.Request.Header.Get("token")
		if len(token) == 0 {
			logs.Error("token is null")
			r.Code = 0
			r.Msg = "token is null"
			c.JSON(http.StatusUnauthorized, r)
			c.Abort()
			return
		}
		tokenTime, b := tokens.Load(token)
		if !b {
			logs.Error("token error")
			r.Code = 0
			r.Msg = "token error"
			c.JSON(http.StatusUnauthorized, r)
			c.Abort()
			return
		}
		timeout := time.Now().After(tokenTime.(time.Time).Add(30 * time.Minute))
		if timeout {
			logs.Error("token is timeout")
			r.Code = 0
			r.Msg = "token is timeout"
			c.JSON(http.StatusUnauthorized, r)
			c.Abort()
			return
		}
		tokens.Store(token, time.Now())
		c.Next()
	}
}

func login(c *gin.Context) {
	r := common.Result{
		Code: 1,
		Msg:  "",
	}
	params := make(map[string]interface{})
	err := c.BindJSON(&params)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}
	userNameParam := params["userName"].(string)
	passwordParam := params["password"].(string)
	userName := config.DefaultString("server.user.name", "")
	password := config.DefaultString("server.user.password", "")
	if userNameParam == "" || passwordParam == "" || userNameParam != userName || passwordParam != password {
		logs.Error("userName : %s , password : %s error", userNameParam, passwordParam)
		r.Code = 0
		r.Msg = "userName or password error ! "
		c.JSON(http.StatusOK, r)
		return
	}
	logs.Info("用户[%s]登录成功！", userName)
	token, err := utils.NextToke()
	if err != nil {
		logs.Error("create token fail")
		r.Code = 0
		r.Msg = "create token fail"
		c.JSON(http.StatusOK, r)
		return
	}
	r.Data = map[string]string{"token": token}
	tokens.Store(token, time.Now())
	c.JSON(http.StatusOK, r)
}
