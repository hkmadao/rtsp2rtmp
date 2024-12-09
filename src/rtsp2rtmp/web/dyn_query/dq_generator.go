package dyn_query

import (
	"errors"

	"github.com/beego/beego/v2/core/config"
)

type DynQuery interface {
	BuildSql(selectStatement SelectStatement) (string, []interface{}, error)
}

func NewDynQuery() (qb DynQuery, err error) {
	driver, err := config.String("server.database.driver")
	if err != nil {
		err = errors.New("database driver param is null")
		return
	}
	if driver == "mysql" {
		qb = new(DynQueryMysql)
	} else if driver == "tidb" {
		qb = new(DynQueryTidb)
	} else if driver == "postgres" {
		qb = new(DynQueryPostgres)
	} else if driver == "sqlite" {
		qb = new(DynQuerySqlite)
	} else {
		err = errors.New("unknown driver for query builder")
	}
	return
}

func getQuestionMarkArr(arrLen int) []string {
	var strArr = make([]string, arrLen)
	for i := 0; i < arrLen; i++ {
		strArr[i] = "?"
	}
	return strArr
}
