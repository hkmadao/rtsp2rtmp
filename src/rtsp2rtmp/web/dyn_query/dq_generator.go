package dyn_query

import "errors"

type DynQuery interface {
	BuildSql(selectStatement SelectStatement) (string, error)
}

func NewDynQuery(driver string) (qb DynQuery, err error) {
	if driver == "mysql" {
		qb = new(DynQueryMysql)
		// } else if driver == "tidb" {
		// 	qb = new(TiDBQueryBuilder)
		// } else if driver == "postgres" {
		// 	qb = new(PostgresQueryBuilder)
		// } else if driver == "sqlite" {
		// 	err = errors.New("sqlite query builder is not supported yet")
	} else {
		err = errors.New("unknown driver for query builder")
	}
	return
}
