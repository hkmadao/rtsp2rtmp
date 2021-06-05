package httpflv

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yumrano/rtsp2rtmp/conf"
	"github.com/yumrano/rtsp2rtmp/dao"
	"github.com/yumrano/rtsp2rtmp/rlog"
)

func ServeHTTP() error {
	defer func() {
		if r := recover(); r != nil {
			rlog.Log.Printf("httpflv ServeHTTP pain %v", r)
		}
	}()
	port, err := conf.GetInt("server.httpflv.port")
	if err != nil {
		rlog.Log.Printf("get httpflv port error: %v. \n use default port : 8080", err)
		port = 8080
	}
	httpflvAddr := ":" + strconv.Itoa(port)
	flvListen, err := net.Listen("tcp", httpflvAddr)
	if err != nil {
		rlog.Log.Printf("%v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/live/", reciver)
	mux.HandleFunc("/camera/all", cameraAll)
	mux.HandleFunc("/camera/edit", cameraEdit)
	mux.HandleFunc("/camera/delete", cameraDelete)
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("statics"))))
	if err := http.Serve(flvListen, mux); err != nil {
		return err
	}
	return nil
}

func cameraAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	result := make(map[string]interface{})
	result["code"] = 1
	result["msg"] = ""
	es, err := dao.CameraSelectAll()
	if err != nil {
		rlog.Log.Printf("get camera list error : %v", err)
		result["code"] = 0
		result["msg"] = ""
		result["data"] = "{}"
		rbytes, err := json.Marshal(result)
		if err != nil {
			rlog.Log.Printf("parse json camera list error : %v", err)
			w.Write([]byte("{}"))
			return
		}
		w.Write(rbytes)
		return
	}

	data := make(map[string]interface{})
	data["total"] = len(es)
	data["page"] = es
	result["data"] = data
	rbytes, err := json.Marshal(result)
	if err != nil {
		rlog.Log.Printf("parse json camera list error : %v", err)
		w.Write([]byte("{}"))
		return
	}
	w.Write(rbytes)
}

func cameraEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
		w.Write([]byte("{}"))
		return
	}
	result := make(map[string]interface{})
	result["code"] = 1
	result["msg"] = ""

	var param map[string]interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rlog.Log.Printf("read param error : %v", err)
	}
	rlog.Log.Printf("%s", r.Method)
	err = json.Unmarshal(body, &param)
	if err != nil {
		rlog.Log.Printf("parse param to json error : %v", err)
	}
	if param["id"] == nil {
		param["id"] = ""
	}
	e := dao.Camera{
		Id:                param["id"].(string),
		Code:              param["code"].(string),
		RtspURL:           param["rtspURL"].(string),
		RtmpURL:           param["rtmpURL"].(string),
		AuthCodeTemp:      param["authCodeTemp"].(string),
		AuthCodePermanent: param["authCodePermanent"].(string),
	}
	var data int64
	if e.Id == "" {
		count, err := dao.CameraCountByCode(e.Code)
		if err != nil || count > 0 {
			rlog.Log.Printf("check code is exist camera error : %v", err)
			result["code"] = 0
			result["msg"] = "code exist !"
			result["data"] = "{}"
			rbytes, _ := json.Marshal(result)
			w.Write(rbytes)
			return
		}
		e.Id = time.Now().Format("20060102150405")
		_, err = dao.CameraInsert(e)
		if err != nil {
			rlog.Log.Printf("insert camera error : %v", err)
			result["code"] = 0
			result["msg"] = "insert camera error !"
			result["data"] = "{}"
			rbytes, _ := json.Marshal(result)
			w.Write(rbytes)
			return
		}
	} else {
		count, err := dao.CameraCountByCode(e.Code)
		if err != nil || count > 1 {
			rlog.Log.Printf("get camera list error : %v", err)
			result["code"] = 0
			result["msg"] = "code exist !"
			result["data"] = "{}"
			rbytes, _ := json.Marshal(result)
			w.Write(rbytes)
			return
		}
		_, err = dao.CameraUpdate(e)
		if err != nil {
			rlog.Log.Printf("update camera error : %v", err)
			result["code"] = 0
			result["msg"] = "update camera error !"
			result["data"] = "{}"
			rbytes, _ := json.Marshal(result)
			w.Write(rbytes)
			return
		}
	}

	result["data"] = data
	rbytes, _ := json.Marshal(result)
	w.Write(rbytes)
}

func cameraDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
		w.Write([]byte("{}"))
		return
	}
	result := make(map[string]interface{})
	result["code"] = 1
	result["msg"] = ""

	var param map[string]interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rlog.Log.Printf("read param error : %v", err)
	}
	rlog.Log.Printf("%s", r.Method)
	err = json.Unmarshal(body, &param)
	if err != nil {
		rlog.Log.Printf("parse param to json error : %v", err)
	}
	e := dao.Camera{Id: param["id"].(string)}
	data, err := dao.CameraDelete(e)
	if err != nil {
		rlog.Log.Printf("delete camera error : %v", err)
		result["code"] = 0
		result["msg"] = "delete camera error !"
		result["data"] = "{}"
		rbytes, _ := json.Marshal(result)
		w.Write(rbytes)
		return
	}

	result["data"] = data
	rbytes, _ := json.Marshal(result)
	w.Write(rbytes)
}

func reciver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	sessionId := time.Now().Format("20060102150405")
	fw := &HttpFlvWriter{
		sessionId:      sessionId,
		isStart:        false,
		responseWriter: w,
	}
	uri := strings.TrimSuffix(strings.TrimLeft(r.RequestURI, "/"), ".flv")
	uris := strings.Split(uri, "/")
	if len(uris) < 2 || uris[0] != "live" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	hms[uris[1]].fws[sessionId] = fw
	done := make(chan interface{})
	hms[uris[1]].fws[sessionId].done = done
	<-done
	rlog.Log.Printf("session %s exit", sessionId)
}

// func ServeHTTP() {
// 	router := gin.Default()
// 	router.GET("/recive", reciver)
// 	port, err := conf.GetInt("server.httpflv.port")
// 	if err != nil {
// 		rlog.Log.Printf("get httpflv port error: %v. \n use default port : 8080", err)
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
