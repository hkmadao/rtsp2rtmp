package dyn_query

import (
	"errors"

	"github.com/beego/beego/v2/core/config"
)

type DynQuery interface {
	newDynQuery(selectStatement SelectStatement) DynQuery
	//return row sql string
	BuildCountSql() (rowSql string, params []interface{}, err error)
	//return row sql string
	BuildPageSql(pageIndex uint64, pageSize uint64) (rowSql string, params []interface{}, err error)
	//return row sql string
	BuildSql(fgCount bool) (rowSql string, params []interface{}, err error)
	// return condition sql string and param list.
	// for example:
	//
	//	`user_name` = ? AND `code` = ? AND (`enabled` = ? OR `status` = ?)
	makeConditionsToken(conditionExpr ConditionExpression) (conditionSql string, params []interface{}, err error)
}

func NewDynQuery(selectStatement SelectStatement) (qb DynQuery, err error) {
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
		qb = new(DynQueryPostgres).newDynQuery(selectStatement)
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
