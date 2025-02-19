package ext

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

func getTokenName() string {
	return "Token"
}

func TokenDeleteByUsername(username string) (i int64, err error) {
	o := orm.NewOrm()
	rowResult, err := o.Raw("DELETE sys_token WHERE username = ?", username).Exec()
	if err != nil {
		logs.Error("delete user: %s tokens error : %v", username, err)
		return 0, err
	}
	return rowResult.RowsAffected()
}
