package conf

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/beego/beego/v2/core/config"
	_ "github.com/beego/beego/v2/core/config/yaml"
	"github.com/beego/beego/v2/core/logs"
)

func init() {
	var mode string
	var lout string
	flag.StringVar(&lout, "lout", logs.AdapterFile, "日志输出方式，默认file")
	flag.StringVar(&mode, "mode", "prod", "启动环境，默认prod")
	// 解析命令行参数写入注册的flag里
	flag.Parse()

	loadConf(mode)
	logConfig(lout)
}

func loadConf(mode string) {
	if mode == "dev" {
		filePath := "./resources/conf/conf-dev.yml"
		err := config.InitGlobalInstance("yaml", filePath)
		if err != nil {
			fmt.Printf("read conf file [%s] error : %v", filePath, err)
			return
		}
		return
	}

	filePath := "./resources/conf/conf-prod.yml"
	err := config.InitGlobalInstance("yaml", filePath)
	if err != nil {
		fmt.Printf("read conf file [%s] error : %v", filePath, err)
		return
	}
}

//日志配置
func logConfig(lout string) {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if lout == logs.AdapterConsole {
		f := &logs.PatternLogFormatter{
			Pattern:    "%w %F:%n %t %m",
			WhenFormat: "2006-01-02 15:04:05.000",
		}
		logs.RegisterFormatter("pattern", f)

		_ = logs.SetGlobalFormatter("pattern")
		level, err := config.Int("server.log.level")
		if err != nil {
			fmt.Printf("can not get log level : %v , use default level logs.LevelInformational", err)
			level = logs.LevelInformational
		}
		logs.SetLogger(logs.AdapterConsole, `{"level":`+strconv.Itoa(level)+`,"color":true}`)
		return
	}
	logPath, err := config.String("server.log.path")
	if err != nil {
		fmt.Printf("can not get log path : %v , use default path ./resources/output/log", err)
		logPath = "./output/log"
	}
	level, err := config.Int("server.log.level")
	if err != nil {
		fmt.Printf("can not get log level : %v , use default level logs.LevelInformational", err)
		level = logs.LevelInformational
	}
	logs.SetLogger(logs.AdapterFile, `{"filename":"`+logPath+`/rtsp2rtmp.log","level":`+strconv.Itoa(level)+`,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
}
