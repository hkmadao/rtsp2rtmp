package ext

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/fileflvmanager/fileflvreader"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

var playerMap sync.Map

func HttpFlvPlay(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	uri := strings.TrimSuffix(strings.TrimLeft(ctx.Request.RequestURI, "/"), ".flv")
	uris := strings.Split(uri, "/")
	if len(uris) < 3 || uris[0] != "live" {
		http.Error(ctx.Writer, "invalid path", http.StatusBadRequest)
		return
	}
	method := uris[1]
	code := uris[2]
	authCode := uris[3]

	conditon := common.GetEqualCondition("code", code)
	camera, err := base_service.CameraFindOneByCondition(conditon)
	if err != nil {
		logs.Error("camera query error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	if !(method == "temp" || method == "permanent") {
		logs.Error("method error : %s", method)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	if method == "temp" {
		var filters = []common.EqualFilter{{Name: "cameraId", Value: camera.Id}, {Name: "authCode", Value: authCode}}
		condition := common.GetEqualConditions(filters)
		cs, err := base_service.CameraShareFindOneByCondition(condition)

		if err != nil {
			logs.Error("CameraShareSelectOne error : %v", err)
			result := common.ErrorResult("internal error")
			ctx.JSON(http.StatusBadRequest, result)
			return
		}
		if time.Now().Before(cs.StartTime) || time.Now().After(cs.Deadline) {
			logs.Error("camera [%s] AuthCodeTemp expired : %s", camera.Code, authCode)
			result := common.ErrorResult(fmt.Sprintf("auth error"))
			ctx.JSON(http.StatusBadRequest, result)
			return
		}

	}
	if method == "permanent" && authCode != camera.PlayAuthCode {
		logs.Error("AuthCodePermanent error : %s", authCode)
		result := common.ErrorResult("auth error")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	logs.Info("player [%s] addr [%s] connecting", code, ctx.Request.RemoteAddr)
	//管理员可以主动中断播放
	playerDone := make(chan int)
	defer close(playerDone)
	const timeout = 10 * time.Second
	flvPlayerDone, err := flvadmin.GetSingleHttpFlvAdmin().AddHttpFlvPlayer(playerDone, timeout/2, code, ctx.Writer)
	if err != nil {
		logs.Error("camera [%s] add player error : %s", code, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	<-flvPlayerDone
	logs.Info("player [%s] addr [%s] exit", code, ctx.Request.RemoteAddr)
}

func HttpFlvVODFileDuration(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()

	fileName, ok := ctx.Params.Get("fileName")
	if !ok {
		logs.Error("get param fileName failed")
		http.Error(ctx.Writer, "invalid path", http.StatusBadRequest)
		return
	}

	if strings.Contains(fileName, "..") {
		logs.Error("fileName: %s illegal", fileName)
		http.Error(ctx.Writer, fmt.Sprintf("fileName: %s illegal", fileName), http.StatusBadRequest)
		return
	}

	duration, err := fileflvreader.FlvDurationReadUntilErr(fileName)
	if err != nil {
		logs.Error("file: %s get duration error", fileName)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	result := common.SuccessResultData(duration)
	ctx.JSON(http.StatusOK, result)
}

func HttpFlvVODStart(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	fileName, ok := ctx.Params.Get("fileName")
	if !ok {
		logs.Error("get param fileName failed")
		http.Error(ctx.Writer, "invalid path", http.StatusBadRequest)
		return
	}

	playerId := ctx.Query("playerId")
	if playerId == "" {
		logs.Error("get param playerId failed")
		http.Error(ctx.Writer, "invalid path", http.StatusBadRequest)
		return
	}

	if strings.Contains(fileName, "..") {
		logs.Error("fileName: %s illegal", fileName)
		http.Error(ctx.Writer, fmt.Sprintf("fileName: %s illegal", fileName), http.StatusBadRequest)
		return
	}

	seekSecond := ctx.Query("seekSecond")
	if seekSecond == "" {
		logs.Error("get param seekSecond failed")
		http.Error(ctx.Writer, "invalid path", http.StatusBadRequest)
		return
	}
	seekSecondUint, err := strconv.ParseUint(seekSecond, 10, 64)
	if err != nil {
		logs.Error("get param seekSecond failed")
		http.Error(ctx.Writer, "invalid path", http.StatusBadRequest)
		return
	}

	ffr := fileflvreader.NewFileFlvReader(seekSecondUint, ctx.Writer, fileName)
	_, ok = playerMap.Load(playerId)
	if ok {
		logs.Error("playerId: %s exists", playerId)
		http.Error(ctx.Writer, fmt.Sprintf("playerId: %s exists", playerId), http.StatusBadRequest)
		return
	}
	playerMap.Store(playerId, ffr)
	<-ffr.GetDone()
	playerMap.Delete(playerId)
	logs.Info("vod player [%s] addr [%s] exit", fileName, ctx.Request.RemoteAddr)
}

func HttpFlvVODFetch(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	playerId := ctx.Query("playerId")
	if playerId == "" {
		logs.Error("get param playerId failed")
		http.Error(ctx.Writer, "invalid path", http.StatusBadRequest)
		return
	}

	seekSecond := ctx.Query("seekSecond")
	if playerId == "" {
		logs.Error("get param seekSecond failed")
		http.Error(ctx.Writer, "invalid path", http.StatusBadRequest)
		return
	}
	seekSecondUint, err := strconv.ParseUint(seekSecond, 10, 64)
	if err != nil {
		logs.Error("get param seekSecond failed")
		http.Error(ctx.Writer, "invalid path", http.StatusBadRequest)
		return
	}

	value, ok := playerMap.Load(playerId)
	if !ok {
		logs.Error("playerId: %s not exists", playerId)
		http.Error(ctx.Writer, fmt.Sprintf("playerId: %s not exists", playerId), http.StatusBadRequest)
		return
	}
	loadFfr := (value.(*fileflvreader.FileFlvReader))
	loadFfr.SetSeekSecond(seekSecondUint)

	logs.Info("vod player [%s] fetch data, addr [%s]", playerId, ctx.Request.RemoteAddr)
	result := common.SuccessResultMsg("fetch sccess")
	ctx.JSON(http.StatusOK, result)
}
