package dyn_query

type DynQuery interface {
	BuildSql(selectStatement SelectStatement) (string, error)
}
