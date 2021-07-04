package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
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
	orm.RegisterModel(new(Camera))
	orm.RegisterModel(new(CameraShare))
}
