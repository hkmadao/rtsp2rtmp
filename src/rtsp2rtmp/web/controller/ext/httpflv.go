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
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/ext/flv_file"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

var playerMap sync.Map

func HttpFlvPlayMediaInfo(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	method, ok := ctx.Params.Get("method")
	if !ok || method == "" {
		logs.Error("path param method is rquired")
		result := common.ErrorResult("path param method is rquired")
		ctx.JSON(http.StatusOK, result)
		return
	}

	code, ok := ctx.Params.Get("code")
	if !ok || code == "" {
		logs.Error("path param code is rquired")
		result := common.ErrorResult("path param code is rquired")
		ctx.JSON(http.StatusOK, result)
		return
	}

	authCode, ok := ctx.Params.Get("authCode.flv")
	if !ok || authCode == "" {
		logs.Error("path param authCode is rquired")
		result := common.ErrorResult("path param authCode is rquired")
		ctx.JSON(http.StatusOK, result)
		return
	}
	authCode = utils.ReverseString(strings.Replace(utils.ReverseString(authCode), utils.ReverseString(".flv"), "", 1))

	conditon := common.GetEqualCondition("code", code)
	camera, err := base_service.CameraFindOneByCondition(conditon)
	if err != nil {
		logs.Error("camera query error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	if !(method == "temp" || method == "permanent") {
		logs.Error("method error : %s", method)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	if method == "temp" {
		var filters = []common.EqualFilter{{Name: "cameraId", Value: camera.Id}, {Name: "authCode", Value: authCode}}
		condition := common.GetEqualConditions(filters)

		exists, err := base_service.CameraShareExistsByCondition(condition)
		if err != nil {
			logs.Error("cameraShareExistsByCondition error : %v", err)
			result := common.ErrorResult("internal error")
			ctx.JSON(http.StatusOK, result)
			return
		}
		if !exists {
			logs.Error("camera [%s] AuthCodeTemp expired : %s", camera.Code, authCode)
			result := common.ErrorResult("auth error")
			ctx.JSON(http.StatusOK, result)
			return
		}

		cs, err := base_service.CameraShareFindOneByCondition(condition)
		if err != nil {
			logs.Error("CameraShareSelectOne error : %v", err)
			result := common.ErrorResult("internal error")
			ctx.JSON(http.StatusOK, result)
			return
		}
		if time.Now().Before(cs.StartTime) || time.Now().After(cs.Deadline) {
			logs.Error("camera [%s] AuthCodeTemp expired : %s", camera.Code, authCode)
			result := common.ErrorResult("auth error")
			ctx.JSON(http.StatusOK, result)
			return
		}

	}
	if method == "permanent" && authCode != camera.PlayAuthCode {
		logs.Error("AuthCodePermanent error : %s", authCode)
		result := common.ErrorResult("auth error")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	liveMediaInfo, err := flvadmin.GetSingleHttpFlvAdmin().GetLiveInfo(camera.Code)
	if err != nil {
		logs.Error("getLiveInfo error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	result := common.SuccessResultData(liveMediaInfo)
	ctx.JSON(http.StatusOK, result)
}

func HttpFlvPlay(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	method, ok := ctx.Params.Get("method")
	if !ok || method == "" {
		logs.Error("path param method is rquired")
		http.Error(ctx.Writer, "path param method is rquired", http.StatusBadRequest)
		return
	}

	code, ok := ctx.Params.Get("code")
	if !ok || code == "" {
		logs.Error("path param code is rquired")
		http.Error(ctx.Writer, "path param code is rquired", http.StatusBadRequest)
		return
	}

	authCode, ok := ctx.Params.Get("authCode.flv")
	if !ok || authCode == "" {
		logs.Error("path param authCode is rquired")
		http.Error(ctx.Writer, "path param authCode is rquired", http.StatusBadRequest)
		return
	}
	authCode = utils.ReverseString(strings.Replace(utils.ReverseString(authCode), utils.ReverseString(".flv"), "", 1))

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

		exists, err := base_service.CameraShareExistsByCondition(condition)
		if err != nil {
			logs.Error("cameraShareExistsByCondition error : %v", err)
			result := common.ErrorResult("internal error")
			ctx.JSON(http.StatusBadRequest, result)
			return
		}
		if !exists {
			logs.Error("camera [%s] AuthCodeTemp expired : %s", camera.Code, authCode)
			result := common.ErrorResult("auth error")
			ctx.JSON(http.StatusBadRequest, result)
			return
		}

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
	flvPlayerDone, addHttpFlvPlayerErr := flvadmin.GetSingleHttpFlvAdmin().AddHttpFlvPlayer(playerDone, timeout/2, code, ctx.Writer)
	if addHttpFlvPlayerErr != nil {
		logs.Error("camera [%s] add player error : %v", code, addHttpFlvPlayerErr)
		if addHttpFlvPlayerErr.IsCustomError() {
			result := common.ErrorResult(addHttpFlvPlayerErr.Error())
			ctx.JSON(http.StatusBadRequest, result)
			return
		}

		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	<-flvPlayerDone
	logs.Info("player [%s] addr [%s] exit", code, ctx.Request.RemoteAddr)
}

func HttpFlvVODFileMediaInfo(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()

	fileName, ok := ctx.Params.Get("fileName")
	if !ok {
		logs.Error("path param fileName is rquired")
		http.Error(ctx.Writer, "path param fileName is rquired", http.StatusBadRequest)
		return
	}

	if strings.Contains(fileName, "..") {
		logs.Error("fileName: %s illegal", fileName)
		http.Error(ctx.Writer, fmt.Sprintf("fileName: %s illegal", fileName), http.StatusBadRequest)
		return
	}

	durationInt, err := fileflvreader.FlvDurationReadUntilErr(fileName)
	mediaInfo := flv_file.FlvMediaInfo{
		Duration: uint32(durationInt),
		HasAudio: true,
	}
	if err != nil {
		logs.Error("file: %s get mediaInfo error", fileName)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	result := common.SuccessResultData(mediaInfo)
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
		logs.Error("path param fileName is rquired")
		http.Error(ctx.Writer, "path param fileName is rquired", http.StatusBadRequest)
		return
	}

	playerId := ctx.Query("playerId")
	if playerId == "" {
		logs.Error("query param playerId is rquired")
		http.Error(ctx.Writer, "query param playerId is rquired", http.StatusBadRequest)
		return
	}

	if strings.Contains(fileName, "..") {
		logs.Error("fileName: %s illegal", fileName)
		http.Error(ctx.Writer, fmt.Sprintf("fileName: %s illegal", fileName), http.StatusBadRequest)
		return
	}

	seekSecond := ctx.Query("seekSecond")
	if seekSecond == "" {
		logs.Error("query param seekSecond is rquired")
		http.Error(ctx.Writer, "query param seekSecond is rquired", http.StatusBadRequest)
		return
	}
	seekSecondUint, err := strconv.ParseUint(seekSecond, 10, 64)
	if err != nil {
		logs.Error("query param seekSecond need uint")
		http.Error(ctx.Writer, "query param seekSecond need uint", http.StatusBadRequest)
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
		logs.Error("query param playerId is rquired")
		http.Error(ctx.Writer, "query param playerId is rquired", http.StatusBadRequest)
		return
	}

	seekSecond := ctx.Query("seekSecond")
	if playerId == "" {
		logs.Error("query param seekSecond is rquired")
		http.Error(ctx.Writer, "query param seekSecond is rquired", http.StatusBadRequest)
		return
	}
	seekSecondUint, err := strconv.ParseUint(seekSecond, 10, 64)
	if err != nil {
		logs.Error("query param seekSecond need uint")
		http.Error(ctx.Writer, "query param seekSecond need uint", http.StatusBadRequest)
		return
	}

	value, ok := playerMap.Load(playerId)
	if !ok {
		logs.Error("playerId: %s not exists or complate", playerId)
		result := common.SuccessResultMsg(fmt.Sprintf("playerId: %s not exists or complate, skip this request", playerId))
		ctx.JSON(http.StatusOK, result)
		return
	}
	loadFfr := (value.(*fileflvreader.FileFlvReader))
	loadFfr.SetSeekSecond(seekSecondUint)

	logs.Info("vod player [%s] fetch data, addr [%s]", playerId, ctx.Request.RemoteAddr)
	result := common.SuccessResultMsg("fetch sccess")
	ctx.JSON(http.StatusOK, result)
}
