package common

type EValueType = string

const (
	// 字符串
	VALUE_TYPE_STRING = "string"
	// 数值
	VALUE_TYPE_NUMBER = "number"
	// 布尔
	VALUE_TYPE_BOOL = "bool"
	// 日期
	VALUE_TYPE_DATE_TIME = "DateTime"
)

type EDataType = string

const (
	// 主键
	DATA_TYPE_INTERNAL_PK = "InternalPK"
	// 外键
	DATA_TYPE_INTERNAL_FK = "InternalFK"
	// 1对多关系 子实体的引用属性
	DATA_TYPE_REF = "InternalRef"
	// 1对1关系 子实体的引用属性
	DATA_TYPE_SINGLE_REF = "InternalSingleRef"
	// 1对1关系 主实体的子属性
	DATA_TYPE_SINGLE = "InternalSingle"
	// 1对多关系 主实体的子属性
	DATA_TYPE_ARRAY = "InternalArray"
	// agg 外键
	DATA_TYPE_AGG_FK = "InternalAggFK"
	// 1对多关系 子实体的引用属性
	DATA_TYPE_AGG_REF = "InternalAggRef"
	// agg 1对1关系 子实体的引用属性
	DATA_TYPE_AGG_SINGLE_REF = "InternalAggSingleRef"
	// agg 1对1关系 主实体的子属性
	DATA_TYPE_AGG_SINGLE = "InternalAggSingle"
	// agg 1对多关系 主实体的子属性
	DATA_TYPE_AGG_ARRAY = "InternalAggArray"
)

type EntityInfo struct {
	// 实体名称
	Name string
	// 实体显示名称
	DisplayName string
	// 类名
	ClassName string
	// 表名
	TableName string
	// 服务路径
	BasePath string
}

type AttributeInfo struct {
	// 属性名称
	Name string
	// 属性显示名称
	DisplayName string
	// 数据库字段名
	ColumnName string
	// 值类型
	ValueType EValueType
	// 数据类型
	DataType EDataType
	// 关联的内部属性名称（外键属性，外键引用属性）
	InnerAttributeName string
	// 外部实体名
	OutEntityName string
	// 外部实体主属性名
	OutEntityPkAttributeName string
	// 外部实体引用本实体的属性名称
	OutEntityReversalAttributeName string
	// 外部实体引用本实体的外键属性名称
	OutEntityIdReversalAttributeName string
}

type EntityDesc struct {
	// 实体信息
	EntityInfo EntityInfo
	// 属性信息
	AttributeInfoMap map[string]*AttributeInfo
	// 获取主键属性描述
	PkAttributeInfo *AttributeInfo
	// 获取不在同一个聚合根下的外键Id属性描述
	NormalFkIdAttributeInfos []*AttributeInfo
	// 获取不在同一个聚合根下的外键属性描述
	NormalFkAttributeInfos []*AttributeInfo
	// 获取不在同一个聚合根下子属性描述
	NormalChildren []*AttributeInfo
	// 1对1情况的子属性
	NormalOne2OneChildren []*AttributeInfo
}
