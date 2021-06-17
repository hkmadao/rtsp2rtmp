package conf

import "github.com/beego/beego/v2/core/logs"

func init() {
	// f := &logs.PatternLogFormatter{
	// 	Pattern:    "%w %F:%n %t %m",
	// 	WhenFormat: "2006-01-02 15:04:05.000",
	// }
	// logs.RegisterFormatter("pattern", f)

	// _ = logs.SetGlobalFormatter("pattern")
	// logs.SetLogger(logs.AdapterConsole)
	logs.SetLogger(logs.AdapterFile, `{"filename":"./output/log/rtsp2rtmp.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
}
