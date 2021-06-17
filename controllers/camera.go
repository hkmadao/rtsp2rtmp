package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/yumrano/rtsp2rtmp/models"
	"github.com/yumrano/rtsp2rtmp/result"
	"github.com/yumrano/rtsp2rtmp/utils"
)

func cros(w http.ResponseWriter, req *http.Request) {
	method := req.Method               //请求方法
	origin := req.Header.Get("Origin") //请求头部
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
		//  header的类型
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
		//              允许跨域设置                                                                                                      可以返回其他子段
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		w.Header().Set("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
		w.Header().Set("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
		w.Header().Set("content-type", "application/json")                                                                                                                                                           // 设置返回格式是json
	}

	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		r := result.Result{Code: 1, Msg: "Options Request!"}
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
	}
}

func CameraList(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	cros(w, req)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	cameras, err := models.CameraSelectAll()
	if err != nil {
		logs.Error("no camera found : %v", err)
		r.Code = 0
		r.Msg = "no camera found"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	page := result.Page{Total: len(cameras), Page: cameras}
	r.Data = page
	rbytes, _ := json.Marshal(r)
	w.Write(rbytes)
}

func CameraEdit(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	cros(w, req)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{
		Code: 1,
		Msg:  "",
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logs.Error("read body err, %v", err)
		return
	}
	logs.Info("json:", string(body))

	q := models.Camera{}
	if err = json.Unmarshal(body, &q); err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}

	if q.Id == "" || len(q.Id) == 0 {
		id, _ := utils.NextToke()
		count, err := models.CameraCountByCode(q.Code)
		if err != nil {
			logs.Error("check camera is exist error : %v", err)
			r.Code = 0
			r.Msg = "check camera is exist"
			rbytes, _ := json.Marshal(r)
			w.Write(rbytes)
			return
		}
		if count > 0 {
			logs.Error("camera code is exist error : %v", err)
			r.Code = 0
			r.Msg = "camera code is exist"
			rbytes, _ := json.Marshal(r)
			w.Write(rbytes)
			return
		}
		q.Id = id
		q.Created = time.Now()
		playAuthCode, _ := utils.NextToke()
		q.AuthCode = playAuthCode
		_, err = models.CameraInsert(q)
		if err != nil {
			logs.Error("camera insert error : %v", err)
			r.Code = 0
			r.Msg = "camera insert error"
			rbytes, _ := json.Marshal(r)
			w.Write(rbytes)
			return
		}
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	count, err := models.CameraCountByCode(q.Code)
	if err != nil {
		logs.Error("check camera is exist error : %v", err)
		r.Code = 0
		r.Msg = "check camera is exist"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	if count > 1 {
		logs.Error("camera code is exist error : %v", err)
		r.Code = 0
		r.Msg = "camera code is exist"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	camera, _ := models.CameraSelectById(q.Id)
	camera.Code = q.Code
	camera.AuthCode = q.AuthCode
	camera.Enabled = q.Enabled
	_, err = models.CameraUpdate(camera)
	if err != nil {
		logs.Error("camera insert error : %v", err)
		r.Code = 0
		r.Msg = "camera insert error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	rbytes, _ := json.Marshal(r)
	w.Write(rbytes)
}

func CameraDelete(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	cros(w, req)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logs.Error("read body err, %v", err)
		return
	}
	logs.Info("json:", string(body))

	q := models.Camera{}
	if err = json.Unmarshal(body, &q); err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	camera := models.Camera{Id: q.Id}
	models.CameraDelete(camera)

	if err != nil {
		logs.Error("delete camera error : %v", err)
		r.Code = 0
		r.Msg = "delete camera error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	//close camera conn
	select {
	case codeStream <- camera.Code:
	case <-time.After(1 * time.Second):
	}

	rbytes, _ := json.Marshal(r)
	w.Write(rbytes)
}

func CameraEnabled(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	cros(w, req)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logs.Error("read body err, %v", err)
		return
	}
	logs.Info("json:", string(body))

	q := models.Camera{}
	if err = json.Unmarshal(body, &q); err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}

	camera, err := models.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		r.Code = 0
		r.Msg = "query camera error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	camera.Enabled = q.Enabled
	_, err = models.CameraUpdate(camera)
	if err != nil {
		logs.Error("enabled camera status %d error : %v", camera.Enabled, err)
		r.Code = 0
		r.Msg = "enabled camera status %d error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	if q.Enabled != 1 {
		//close camera conn
		select {
		case codeStream <- camera.Code:
		case <-time.After(1 * time.Second):
		}
	}

	rbytes, _ := json.Marshal(r)
	w.Write(rbytes)
}

var codeStream = make(chan string)

func CodeStream() <-chan string {
	return codeStream
}
