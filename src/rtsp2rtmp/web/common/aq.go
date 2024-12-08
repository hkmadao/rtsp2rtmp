package common

type AqPageInfoInput struct {
	/// 当前页码
	PageIndex uint64 `json:"pageIndex"`
	/// 分页大小
	PageSize uint64 `json:"pageSize"`
	/// 总记录数
	TotalCount uint64 `json:"totalCount"`
	/// 查询条件
	LogicNode *AqLogicNode `json:"logicNode"`
	/// 排序设置
	Orders []AqOrder `json:"orders"`
}

type AqCondition struct {
	// 查询条件
	LogicNode *AqLogicNode `json:"logicNode"`
	// 分页信息
	Orders []AqOrder `json:"orders"`
}

type AqLogicNode struct {
	// 逻辑操作编码
	LogicOperatorCode ELogicOperatorCode `json:"logicOperatorCode"`
	//子节点
	LogicNode *AqLogicNode `json:"logicNode"`
	//查询条件集合
	FilterNodes []AqFilterNode `json:"filterNodes"`
}

type AqOrder struct {
	// 排序方向
	Direction EOrderDirection `json:"direction"`
	// 排序属性
	Property string `json:"property"`
	// 是否忽略
	IgnoreCase bool `json:"ignoreCase"`
}

type AqFilterNode struct {
	// 查询条件名称
	Name string `json:"name"`
	// 比较操作符编码
	OperatorCode EOperatorCode `json:"operatorCode"`
	// 查询参数
	FilterParams []interface{} `json:"filterParams"`
}
