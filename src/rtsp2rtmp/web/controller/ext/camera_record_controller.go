package ext

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/flvadmin/fileflvmanager/fileflvreader"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/ext/flv_file"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func CameraRecordFileDuration(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()

	idCameraRecord, ok := ctx.Params.Get("idCameraRecord")
	if !ok || idCameraRecord == "" {
		logs.Error("path param idCameraRecord is rquired")
		http.Error(ctx.Writer, "path param idCameraRecord is rquired", http.StatusBadRequest)
		return
	}

	camera_record, err := base_service.CameraRecordSelectById(idCameraRecord)
	if err != nil {
		logs.Error("CameraRecordSelectById error: %v", err)
		http.Error(ctx.Writer, fmt.Sprintf("idCameraRecord: %s not found", idCameraRecord), http.StatusBadRequest)
		return
	}

	mediaInfo := flv_file.FlvMediaInfo{
		Duration: camera_record.Duration,
		HasAudio: true,
	}

	if camera_record.FgTemp {
		durationInt, err := fileflvreader.FlvDurationReadUntilErr(camera_record.TempFileName)
		mediaInfo = flv_file.FlvMediaInfo{
			Duration: uint32(durationInt),
			HasAudio: true,
		}
		if err != nil {
			logs.Error("file: %s get mediaInfo error", camera_record.TempFileName)
			http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	result := common.SuccessResultData(mediaInfo)
	ctx.JSON(http.StatusOK, result)
}

func CameraRecordFilePlay(ctx *gin.Context) {
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

	idCameraRecord, ok := ctx.Params.Get("idCameraRecord")
	if !ok || idCameraRecord == "" {
		logs.Error("path param idCameraRecord is rquired")
		http.Error(ctx.Writer, "path param idCameraRecord is rquired", http.StatusBadRequest)
		return
	}

	camera_record, err := base_service.CameraRecordSelectById(idCameraRecord)
	if err != nil {
		logs.Error("CameraRecordSelectById error: %v", err)
		http.Error(ctx.Writer, fmt.Sprintf("idCameraRecords: %s not found", idCameraRecord), http.StatusBadRequest)
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

	fileName := camera_record.FileName
	if camera_record.FgTemp {
		fileName = camera_record.TempFileName
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

func CameraRecordFileFetch(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	playerId := ctx.Query("playerId")
	if playerId == "" {
		logs.Error("query param playerId failed")
		http.Error(ctx.Writer, "query param playerId is rquired", http.StatusBadRequest)
		return
	}

	seekSecond := ctx.Query("seekSecond")
	if playerId == "" {
		logs.Error("get param seekSecond failed")
		http.Error(ctx.Writer, "query param seekSecond is rquired", http.StatusBadRequest)
		return
	}
	seekSecondUint, err := strconv.ParseUint(seekSecond, 10, 64)
	if err != nil {
		logs.Error("get param seekSecond failed")
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
