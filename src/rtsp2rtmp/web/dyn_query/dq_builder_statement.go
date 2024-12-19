package dyn_query

type SelectStatement struct {
	Selects  []SelectExpr
	From     []FromExpr
	Join     []JoinExpr
	SqlWhere *ConditionExpression
	Having   *ConditionExpression
	Orders   []OrderExpr
}

type SelectExpr = ColumnRef

type ColumnRef struct {
	FgPrimary      bool
	TableAliasName string
	ColumnName     string
}

type TableRef struct {
	AliasName string
	TableName string
}

type FromExpr struct {
	TableRef
	Columns []ColumnRef
}

type JoinExpr struct {
	JoinType EJoinType
	TableRef
	On *ConditionExpression
}

type ConditionExpression struct {
	ConditionType EConditionType
	SimpleExprs   []SimpleExpr
	Child         *ConditionExpression
}

type SimpleExpr struct {
	ExprType EExprType
	ColumnRef
	ValueType EValueType
	Values    []interface{}
}

type OrderExpr struct {
	OrderType EOrderType
	ColumnRef
}

type EValueType = uint32

const (
	String = iota
	Number
	Bool
	DateTime
)

type EOrderType = uint32

const (
	ASC = iota
	DESC
)

type EExprType uint32

const (
	IsNull = iota
	NotNull
	Like
	LeftLike
	RightLike
	Equal
	NotEqual
	GTE
	GT
	LT
	LTE
	In
	NotIn
)

type EConditionType uint32

const (
	AND = iota
	OR
)

type EJoinType uint32

const (
	InnerJoin = iota
	LeftJoin
	RightJoin
)

func StartBuild() *SelectStatement {
	return new(SelectStatement)
}

func (selectStatemet *SelectStatement) BuildSelect(selectExprs ...SelectExpr) *SelectStatement {
	selectStatemet.Selects = append(selectStatemet.Selects, selectExprs...)
	return selectStatemet
}

func (selectStatemet *SelectStatement) BuildFrom(fromExprs ...FromExpr) *SelectStatement {
	selectStatemet.From = append(selectStatemet.From, fromExprs...)
	return selectStatemet
}

func (selectStatemet *SelectStatement) BuildJoin(joinExprs ...JoinExpr) *SelectStatement {
	selectStatemet.Join = append(selectStatemet.Join, joinExprs...)
	return selectStatemet
}

func (selectStatemet *SelectStatement) BuildWhere(conditionExpr ConditionExpression) *SelectStatement {
	selectStatemet.SqlWhere = &conditionExpr
	return selectStatemet
}

func (selectStatemet *SelectStatement) BuildHaving(conditionExpr ConditionExpression) *SelectStatement {
	selectStatemet.Having = &conditionExpr
	return selectStatemet
}

func (selectStatemet *SelectStatement) BuildOrder(orders ...OrderExpr) *SelectStatement {
	selectStatemet.Orders = append(selectStatemet.Orders, orders...)
	return selectStatemet
}
