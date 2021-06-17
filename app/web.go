package app

import (
	"net"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/yumrano/rtsp2rtmp/controllers"
)

func ServeHTTP() error {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("httpflv ServeHTTP panic %v", r)
		}
	}()
	port, err := config.Int("server.httpflv.port")
	if err != nil {
		logs.Error("get httpflv port error: %v. \n use default port : 8080", err)
		port = 8080
	}
	httpflvAddr := ":" + strconv.Itoa(port)
	flvListen, err := net.Listen("tcp", httpflvAddr)
	if err != nil {
		logs.Error("%v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/live/", controllers.Live)
	mux.HandleFunc("/camera/list", controllers.CameraList)
	mux.HandleFunc("/camera/edit", controllers.CameraEdit)
	mux.HandleFunc("/camera/delete", controllers.CameraDelete)
	mux.HandleFunc("/camera/enabled", controllers.CameraEnabled)
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	if err := http.Serve(flvListen, mux); err != nil {
		return err
	}
	return nil
}

// func cameraAll(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	result := make(map[string]interface{})
// 	result["code"] = 1
// 	result["msg"] = ""
// 	es, err := models.CameraSelectAll()
// 	if err != nil {
// 		logs.Error("get camera list error : %v", err)
// 		result["code"] = 0
// 		result["msg"] = ""
// 		result["data"] = "{}"
// 		rbytes, err := json.Marshal(result)
// 		if err != nil {
// 			logs.Error("parse json camera list error : %v", err)
// 			w.Write([]byte("{}"))
// 			return
// 		}
// 		w.Write(rbytes)
// 		return
// 	}

// 	data := make(map[string]interface{})
// 	data["total"] = len(es)
// 	data["page"] = es
// 	result["data"] = data
// 	rbytes, err := json.Marshal(result)
// 	if err != nil {
// 		logs.Error("parse json camera list error : %v", err)
// 		w.Write([]byte("{}"))
// 		return
// 	}
// 	w.Write(rbytes)
// }

// func cameraEdit(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	if r.Method == "OPTIONS" {
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
// 		w.Write([]byte("{}"))
// 		return
// 	}
// 	result := make(map[string]interface{})
// 	result["code"] = 1
// 	result["msg"] = ""

// 	var param map[string]interface{}
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		logs.Error("read param error : %v", err)
// 	}
// 	logs.Error("%s", r.Method)
// 	err = json.Unmarshal(body, &param)
// 	if err != nil {
// 		logs.Error("parse param to json error : %v", err)
// 	}
// 	if param["id"] == nil {
// 		param["id"] = ""
// 	}
// 	e := models.Camera{
// 		Id:       param["id"].(string),
// 		Code:     param["code"].(string),
// 		RtspURL:  param["rtspURL"].(string),
// 		RtmpURL:  param["rtmpURL"].(string),
// 		AuthCode: param["authCode"].(string),
// 	}
// 	var data int64
// 	if e.Id == "" {
// 		count, err := models.CameraCountByCode(e.Code)
// 		if err != nil || count > 0 {
// 			logs.Error("check code is exist camera error : %v", err)
// 			result["code"] = 0
// 			result["msg"] = "code exist !"
// 			result["data"] = "{}"
// 			rbytes, _ := json.Marshal(result)
// 			w.Write(rbytes)
// 			return
// 		}
// 		e.Id = time.Now().Format("20060102150405")
// 		_, err = models.CameraInsert(e)
// 		if err != nil {
// 			logs.Error("insert camera error : %v", err)
// 			result["code"] = 0
// 			result["msg"] = "insert camera error !"
// 			result["data"] = "{}"
// 			rbytes, _ := json.Marshal(result)
// 			w.Write(rbytes)
// 			return
// 		}
// 	} else {
// 		count, err := models.CameraCountByCode(e.Code)
// 		if err != nil || count > 1 {
// 			logs.Error("get camera list error : %v", err)
// 			result["code"] = 0
// 			result["msg"] = "code exist !"
// 			result["data"] = "{}"
// 			rbytes, _ := json.Marshal(result)
// 			w.Write(rbytes)
// 			return
// 		}
// 		_, err = models.CameraUpdate(e)
// 		if err != nil {
// 			logs.Error("update camera error : %v", err)
// 			result["code"] = 0
// 			result["msg"] = "update camera error !"
// 			result["data"] = "{}"
// 			rbytes, _ := json.Marshal(result)
// 			w.Write(rbytes)
// 			return
// 		}
// 	}

// 	result["data"] = data
// 	rbytes, _ := json.Marshal(result)
// 	w.Write(rbytes)
// }

// func cameraDelete(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	if r.Method == "OPTIONS" {
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
// 		w.Write([]byte("{}"))
// 		return
// 	}
// 	result := make(map[string]interface{})
// 	result["code"] = 1
// 	result["msg"] = ""

// 	var param map[string]interface{}
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		logs.Error("read param error : %v", err)
// 	}
// 	logs.Error("%s", r.Method)
// 	err = json.Unmarshal(body, &param)
// 	if err != nil {
// 		logs.Error("parse param to json error : %v", err)
// 	}
// 	e := models.Camera{Id: param["id"].(string)}
// 	data, err := models.CameraDelete(e)
// 	if err != nil {
// 		logs.Error("delete camera error : %v", err)
// 		result["code"] = 0
// 		result["msg"] = "delete camera error !"
// 		result["data"] = "{}"
// 		rbytes, _ := json.Marshal(result)
// 		w.Write(rbytes)
// 		return
// 	}

// 	result["data"] = data
// 	rbytes, _ := json.Marshal(result)
// 	w.Write(rbytes)
// }

// func reciver(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	r := result.Result{
// 		Code: 1,
// 		Msg:  "",
// 	}
// 	uri := strings.TrimSuffix(strings.TrimLeft(req.RequestURI, "/"), ".flv")
// 	uris := strings.Split(uri, "/")
// 	if len(uris) < 2 || uris[0] != "live" {
// 		http.Error(w, "invalid path", http.StatusBadRequest)
// 		return
// 	}
// 	code := uris[1]
// 	logs.Info("player [%s] addr [%s] connecting", code, req.RemoteAddr)
// 	//管理员可以主动中断播放
// 	endStream, heartbeatStream, _, err := httpflv.AddHttpFlvPlayer(code, w)
// 	if err != nil {
// 		logs.Error("camera [%s] add player error : %s", code)
// 		r.Code = 0
// 		r.Msg = "add player error"
// 		rbytes, _ := json.Marshal(r)
// 		w.Write(rbytes)
// 		return
// 	}
// Loop:
// 	for {
// 		select {
// 		case <-endStream:
// 			break Loop
// 		case <-heartbeatStream:
// 			continue
// 		case <-time.After(10 * time.Second):
// 			logs.Info("player [%s] addr [%s] timeout exit", code, req.RemoteAddr)
// 			break Loop
// 		}
// 	}
// 	logs.Info("player [%s] addr [%s] exit", code, req.RemoteAddr)
// }

// func ServeHTTP() {
// 	router := gin.Default()
// 	router.GET("/recive", reciver)
// 	port, err := conf.GetInt("server.httpflv.port")
// 	if err != nil {
// 		logs.Error("get httpflv port error: %v. \n use default port : 8080", err)
// 		port = 8080
// 	}
// 	err = router.Run(":" + strconv.Itoa(port))
// 	if err != nil {
// 		rlog.Log.Fatalln("Start HTTP Server error", err)
// 	}
// }
// func reciver(c *gin.Context) {
// 	c.Header("Access-Control-Allow-Origin", "*")
// 	fw := &FlvResponseWriter{
// 		Key:     "1",
// 		IsStart: false,
// 		writer:  c.Writer,
// 	}
// 	httpFlvWriter.FlvResponseWriters["1"] = fw
// }
