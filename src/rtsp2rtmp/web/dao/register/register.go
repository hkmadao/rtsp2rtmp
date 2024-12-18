package register

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/desc"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	_ "github.com/lib/pq"
)

func init() {

	driver, err := config.String("server.database.driver")
	if err != nil {
		logs.Error("database driver param is null")
		return
	}
	url, err := config.String("server.database.url")
	if err != nil {
		logs.Error("database url param is null")
		return
	}
	driveType, err := config.Int("server.database.driver-type")
	if err != nil {
		logs.Error("database driver-type param is null")
		return
	}
	showSql, err := config.Bool("server.database.show-sql")
	if err != nil {
		logs.Error("database show-sql param error : %v", err)
	}
	if showSql {
		orm.Debug = showSql
	}
	logs.Info("user database %v", driver)
	orm.RegisterDriver(driver, orm.DriverType(driveType))
	orm.RegisterDataBase("default", driver, url)
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(entity.Camera))
	orm.RegisterModel(new(entity.CameraShare))
	orm.RegisterModelWithPrefix("sys_", new(entity.User))
	orm.RegisterModelWithPrefix("sys_", new(entity.Token))
}

var descMap = make(map[string]*common.EntityDesc)

func init() {
	descMap["Camera"] = desc.GetCameraDesc()
	descMap["CameraShare"] = desc.GetCameraShareDesc()
	descMap["User"] = desc.GetUserDesc()
	descMap["Token"] = desc.GetTokenDesc()
}

func GetDescMapByKey(key string) *common.EntityDesc {
	return descMap[key]
}
