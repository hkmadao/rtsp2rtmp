package controllers

import (
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/utils"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service"
)

func CameraShareList(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := common.Result{Code: 1, Msg: ""}
	cameraId := c.Query("cameraId")
	if cameraId == "" {
		logs.Error("no cameraId found")
		r.Code = 0
		r.Msg = "no cameraId found"
		c.JSON(http.StatusOK, r)
		return
	}
	cameraShares, err := service.CameraShareSelectByCameraId(cameraId)
	if err != nil {
		logs.Error("no camerashare found : %v", err)
		r.Code = 0
		r.Msg = "no camerashare found"
		c.JSON(http.StatusOK, r)
		return
	}
	page := common.Page{Total: len(cameraShares), Page: cameraShares}
	r.Data = page
	c.JSON(http.StatusOK, r)
}

func CameraShareEdit(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := common.Result{
		Code: 1,
		Msg:  "",
	}
	q := entity.CameraShare{}
	err := c.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}

	_, err = service.CameraSelectById(q.CameraId)
	if err != nil {
		logs.Error("not fount camera : %v", err)
		r.Code = 0
		r.Msg = "not fount camera"
		c.JSON(http.StatusOK, r)
		return
	}

	if q.Id == "" || len(q.Id) == 0 {
		id, _ := utils.UUID()
		q.Id = id
		q.Created = time.Now()
		playAuthCode, _ := utils.UUID()
		q.AuthCode = playAuthCode
		_, err = service.CameraShareInsert(q)
		if err != nil {
			logs.Error("camerashare insert error : %v", err)
			r.Code = 0
			r.Msg = "camerashare insert error"
			c.JSON(http.StatusOK, r)
			return
		}
		c.JSON(http.StatusOK, r)
		return
	}
	cameraShare, _ := service.CameraShareSelectById(q.Id)
	cameraShare.Name = q.Name
	cameraShare.StartTime = q.StartTime
	cameraShare.Deadline = q.Deadline
	// camera.Enabled = q.Enabled
	_, err = service.CameraShareUpdate(cameraShare)
	if err != nil {
		logs.Error("camerashare insert error : %v", err)
		r.Code = 0
		r.Msg = "camerashare insert error"
		c.JSON(http.StatusOK, r)
		return
	}
	c.JSON(http.StatusOK, r)
}

func CameraShareDelete(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := common.Result{Code: 1, Msg: ""}
	id, b := c.Params.Get("id")
	if !b {
		r.Code = 0
		r.Msg = "id is null"
		c.JSON(http.StatusOK, r)
		return
	}
	camera := entity.CameraShare{Id: id}
	_, err := service.CameraShareDelete(camera)

	if err != nil {
		logs.Error("delete camerashare error : %v", err)
		r.Code = 0
		r.Msg = "delete camerashare error"
		c.JSON(http.StatusOK, r)
		return
	}

	c.JSON(http.StatusOK, r)
}

func CameraShareEnabled(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := common.Result{Code: 1, Msg: ""}
	q := entity.CameraShare{}
	err := c.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}

	camera, err := service.CameraShareSelectById(q.Id)
	if err != nil {
		logs.Error("query camerashare error : %v", err)
		r.Code = 0
		r.Msg = "query camerashare error"
		c.JSON(http.StatusOK, r)
		return
	}
	camera.Enabled = q.Enabled
	_, err = service.CameraShareUpdate(camera)
	if err != nil {
		logs.Error("enabled camerashare status %d error : %v", camera.Enabled, err)
		r.Code = 0
		r.Msg = "enabled camerashare status %d error"
		c.JSON(http.StatusOK, r)
		return
	}

	c.JSON(http.StatusOK, r)
}
