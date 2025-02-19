package dyn_query

import (
	"fmt"
)

func (builder QuerySqlBuilder) GetCountSql() (sqlStr string, params []interface{}, err error) {
	sqlStr, params, err = builder.dynQuery.BuildCountSql()
	if err != nil {
		err = fmt.Errorf("makeSqlByCondition error: %v", err)
		return
	}
	return
}

func (builder QuerySqlBuilder) GetSql() (sqlStr string, params []interface{}, err error) {
	sqlStr, params, err = builder.dynQuery.BuildSql(false)
	if err != nil {
		err = fmt.Errorf("makeSqlByCondition error: %v", err)
		return
	}
	return
}

func (builder QuerySqlBuilder) GetPageSql(
	pageIndex uint64, pageSize uint64,
) (sqlStr string, params []interface{}, err error) {
	sqlStr, params, err = builder.dynQuery.BuildPageSql(pageIndex, pageSize)
	if err != nil {
		err = fmt.Errorf("makeSqlByCondition error: %v", err)
		return
	}
	return
}
