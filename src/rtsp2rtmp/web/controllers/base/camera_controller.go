package base

import (
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func CameraAq(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := common.Result{Code: 1, Msg: ""}
	condition := common.AqCondition{}
	err := c.BindJSON(&condition)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}
	cameras, err := base_service.CameraFindCollectionByCondition(condition)
	if err != nil {
		logs.Error("no camera found : %v", err)
		r.Code = 0
		r.Msg = "no camera found"
		c.JSON(http.StatusOK, r)
		return
	}
	page := common.Page{Total: len(cameras), Page: cameras}
	r.Data = page
	c.JSON(http.StatusOK, r)
}
