package dyn_query

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"

	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/register"
)

func mergeMaps(map1, map2 map[string]uint32) map[string]uint32 {
	for k, v := range map2 {
		map1[k] = v
	}
	return map1
}

func Md5(unMd5Str string) (md5Str string) {
	md5Str = fmt.Sprintf("%x", md5.Sum([]byte(unMd5Str)))
	return
}

func getJoinNamesFromLogicNodes(
	logic_node *common.AqLogicNode,
) (result_map map[string]uint32) {
	result_map = make(map[string]uint32)
	if logic_node != nil {
		var filter_nodes = logic_node.FilterNodes
		for _, filter_node := range filter_nodes {
			var names []string = strings.Split(filter_node.Name, ".")
			if len(names) > 1 {
				names = names[:len(names)-1]
				result_map[strings.Join(names, ".")] = 1
			}
		}

		var child_result = getJoinNamesFromLogicNodes(logic_node.LogicNode)
		result_map = mergeMaps(result_map, child_result)
	}
	return
}

func getJoinNamesFromOrders(
	result_map map[string]uint32,
	orders []common.AqOrder,
) map[string]uint32 {
	for _, order := range orders {
		var names []string = strings.Split(order.Property, ",")
		if len(names) > 1 {
			names = names[:len(names)-1]
			result_map[strings.Join(names, ".")] = 1
		}
	}
	return result_map
}

func getDescAttr(
	attrs []string,
	parent_name string,
	index int,
) (attr_info common.AttributeInfo, err error) {
	var parent_desc = register.GetDescMapByKey(parent_name)
	if parent_desc == nil {
		err = fmt.Errorf("can not get description info: %s", parent_name)
		return
	}
	attr_info_pt := parent_desc.AttributeInfoMap[attrs[index]]
	if nil == attr_info_pt {
		err = fmt.Errorf("%s can not get attribute info: %s", parent_name, attrs[index])
		return
	}

	if len(attrs)-1 == index {
		attr_info = *attr_info_pt
		return
	} else {
		return getDescAttr(attrs, attr_info_pt.OutEntityName, index+1)
	}
}

func makeJoinSelect(
	attrs []string,
	refEntityName string,
	refEntityAliasName string,
	index int,
) (joinExprs []JoinExpr, err error) {
	if len(attrs) < 1 {
		return
	}
	var parent_desc = register.GetDescMapByKey(refEntityName)
	if parent_desc == nil {
		err = fmt.Errorf("can not get description info: %s", refEntityName)
		return
	}
	attr_info := parent_desc.AttributeInfoMap[attrs[index]]
	if nil == attr_info {
		err = fmt.Errorf("%s can not get attribute info: %s", refEntityName, attrs[index])
		return
	}
	var ref_desc = register.GetDescMapByKey(attr_info.OutEntityName)
	if ref_desc == nil {
		err = fmt.Errorf("can not get description info: %s", attr_info.OutEntityName)
		return
	}
	var joinExpr = new(JoinExpr)
	joinExpr.JoinType = LeftJoin
	joinExpr.TableName = ref_desc.EntityInfo.TableName
	var alias_name = strings.Join(attrs[:index+1], ".")
	var aliasNameMd5 = Md5(alias_name)
	joinExpr.AliasName = aliasNameMd5
	if attr_info.DataType == common.DATA_TYPE_REF || attr_info.DataType == common.DATA_TYPE_SINGLE_REF || attr_info.DataType == common.DATA_TYPE_AGG_REF || attr_info.DataType == common.DATA_TYPE_AGG_SINGLE_REF {
		fk_attr_info := ref_desc.AttributeInfoMap[attr_info.InnerAttributeName]
		if nil == fk_attr_info {
			err = fmt.Errorf("%s can not get out entity InnerAttributeName attribute info: %s", ref_desc.EntityInfo.ClassName, attr_info.InnerAttributeName)
			return
		}
		var on = new(ConditionExpression)
		on.ConditionType = AND
		var simpleExprs = make([]SimpleExpr, 2)
		simpleExprs[0].TableAliasName = refEntityAliasName
		simpleExprs[0].ColumnName = parent_desc.PkAttributeInfo.ColumnName
		simpleExprs[1].TableAliasName = aliasNameMd5
		simpleExprs[1].ColumnName = fk_attr_info.ColumnName
		on.SimpleExprs = simpleExprs
		joinExpr.On = on
	} else {
		out_entity_id_reversal_attr_info := ref_desc.AttributeInfoMap[attr_info.OutEntityIdReversalAttributeName]
		if nil == out_entity_id_reversal_attr_info {
			err = fmt.Errorf("entity desc: %s can not get out entity id reversal attribute info: %s", ref_desc.EntityInfo.ClassName, attr_info.OutEntityIdReversalAttributeName)
			return
		}
		var on = new(ConditionExpression)
		on.ConditionType = AND
		var simpleExprs = make([]SimpleExpr, 2)
		simpleExprs[0].TableAliasName = refEntityAliasName
		simpleExprs[0].ColumnName = parent_desc.PkAttributeInfo.ColumnName
		simpleExprs[1].TableAliasName = aliasNameMd5
		simpleExprs[1].ColumnName = out_entity_id_reversal_attr_info.ColumnName
		on.SimpleExprs = simpleExprs
		joinExpr.On = on
	}
	joinExprs = append(joinExprs, *joinExpr)
	var next_index = index + 1
	if len(attrs) > next_index {
		child_joinExprs, err_child := makeJoinSelect(attrs, attr_info.OutEntityName, aliasNameMd5, next_index)
		if err_child != nil {
			err = err_child
			return
		}
		joinExprs = append(joinExprs, child_joinExprs...)
	}
	return
}

func makeColumnOrderBy(
	orders []common.AqOrder,
	main_entity_name string,
	main_table_alias string,
) (orderExprs []OrderExpr, err error) {
	for _, order := range orders {
		if strings.ToUpper(order.Direction) == common.ORDER_DIRECTION_ASC {
			var names = strings.Split(order.Property, ".")
			if len(names) < 2 {
				var attr_info, err_get = getDescAttr(names, main_entity_name, 0)
				if err_get != nil {
					err = fmt.Errorf("make_column_order_by error: %v", err_get)
				}
				var orderExpr = OrderExpr{
					OrderType: ASC,
				}
				orderExpr.TableAliasName = main_table_alias
				orderExpr.ColumnName = attr_info.ColumnName
				orderExprs = append(orderExprs, orderExpr)
			} else {
				var alias_name = strings.Join(names[:len(names)-1], ".")
				var attr_info, err_get = getDescAttr(names, main_entity_name, 0)
				if err_get != nil {
					err = fmt.Errorf("make_column_order_by error: %v", err_get)
				}
				var orderExpr = OrderExpr{
					OrderType: ASC,
				}
				orderExpr.TableAliasName = Md5(alias_name)
				orderExpr.ColumnName = attr_info.ColumnName
				orderExprs = append(orderExprs, orderExpr)
			}
		} else {
			var names = strings.Split(order.Property, ".")
			if len(names) < 2 {
				var attr_info, err_get = getDescAttr(names, main_entity_name, 0)
				if err_get != nil {
					err = fmt.Errorf("make_column_order_by error: %v", err_get)
				}
				var orderExpr = OrderExpr{
					OrderType: DESC,
				}
				orderExpr.TableAliasName = main_table_alias
				orderExpr.ColumnName = attr_info.ColumnName
				orderExprs = append(orderExprs, orderExpr)
			} else {
				var alias_name = strings.Join(names[:len(names)-1], ".")
				var attr_info, err_get = getDescAttr(names, main_entity_name, 0)
				if err_get != nil {
					err = fmt.Errorf("make_column_order_by error: %v", err_get)
				}
				var orderExpr = OrderExpr{
					OrderType: DESC,
				}
				orderExpr.TableAliasName = Md5(alias_name)
				orderExpr.ColumnName = attr_info.ColumnName
				orderExprs = append(orderExprs, orderExpr)
			}
		}
	}
	return
}

func makeCondition(
	logic_node *common.AqLogicNode,
	conditionExpr *ConditionExpression,
	main_table_alias string,
	main_entity_name string,
) (err error) {
	conditionExpr.ConditionType = AND
	if logic_node.LogicOperatorCode == common.LOGIC_OPERATOR_CODE_OR {
		conditionExpr.ConditionType = OR
	}

	// parent level node
	var simple_exprs, err_make = makeSimpleExpr(
		logic_node,
		main_table_alias,
		main_entity_name,
	)
	if err_make != nil {
		err = fmt.Errorf("makeCondition error: %v", err_make)
		return
	}

	conditionExpr.SimpleExprs = append(conditionExpr.SimpleExprs, simple_exprs...)
	// children level node
	var sub_logic_node = logic_node.LogicNode
	if nil != sub_logic_node && len(sub_logic_node.FilterNodes) > 0 {
		conditionExpr.Child = new(ConditionExpression)
		recursionMakeSimpleExpr(
			sub_logic_node,
			main_table_alias,
			main_entity_name,
			conditionExpr.Child,
		)
	}
	return
}

func recursionMakeSimpleExpr(
	sub_logic_node *common.AqLogicNode,
	main_table_alias string,
	main_entity_name string,
	conditionExpr *ConditionExpression,
) (err error) {
	//has children logic node
	if nil != sub_logic_node {
		if len(sub_logic_node.FilterNodes) == 0 {
			return
		}
		conditionExpr.ConditionType = AND
		if sub_logic_node.LogicOperatorCode == common.LOGIC_OPERATOR_CODE_OR {
			conditionExpr.ConditionType = OR
		}
		var sub_simple_exprs, err_make = makeSimpleExpr(
			sub_logic_node,
			main_table_alias,
			main_entity_name,
		)
		if err_make != nil {
			err = fmt.Errorf("recursion_make_simple_expr error: %v", err_make)
			return
		}

		conditionExpr.SimpleExprs = append(conditionExpr.SimpleExprs, sub_simple_exprs...)
		recursionMakeSimpleExpr(
			sub_logic_node.LogicNode,
			main_table_alias,
			main_entity_name,
			conditionExpr.Child,
		)
	}
	return
}

func makeSimpleExpr(
	logic_node *common.AqLogicNode,
	main_table_alias string,
	main_entity_name string,
) (simpleExprs []SimpleExpr, err error) {
	simpleExprs = make([]SimpleExpr, 0)
	for _, filter_node := range logic_node.FilterNodes {
		var simpleExpr, err_build = buildFilter(
			filter_node,
			main_entity_name,
			main_table_alias,
			filter_node.OperatorCode,
		)
		if err_build != nil {
			err = fmt.Errorf("make_simple_expr error: %v", err_build)
			return
		}
		simpleExprs = append(simpleExprs, simpleExpr)
	}
	return
}

func buildFilter(
	filter_node common.AqFilterNode,
	main_entity_name string,
	main_table_alias string,
	operator_code common.EOperatorCode,
) (simpleExpr SimpleExpr, err error) {
	if strings.Contains(filter_node.Name, ".") {
		var names = strings.Split(filter_node.Name, ".")
		var alias_name = strings.Join(names[:len(names)-1], ".")
		var attr_info, err_get = getDescAttr(names, main_entity_name, 0)
		if err_get != nil {
			err = fmt.Errorf("build_filter error: %v", err_get)
			return
		}
		var expr_temp, err_build = buildFilterByOperatorCode(
			operator_code,
			filter_node,
			alias_name,
			attr_info,
			true,
		)
		if err_build != nil {
			err = fmt.Errorf("build_filter error: %v", err_build)
			return
		}
		simpleExpr = expr_temp
	} else {
		attr_info, err_get := getDescAttr(strings.Split(filter_node.Name, "."), main_entity_name, 0)
		if err_get != nil {
			err = fmt.Errorf("build_filter error: %v", err_get)
			return
		}
		expr_temp, err_build := buildFilterByOperatorCode(
			operator_code,
			filter_node,
			main_table_alias,
			attr_info,
			false,
		)
		if err_build != nil {
			err = fmt.Errorf("build_filter error: %v", err_build)
			return
		}
		simpleExpr = expr_temp
	}
	return
}

func buildFilterByOperatorCode(
	operator_code common.EOperatorCode,
	filter_node common.AqFilterNode,
	alias_name string,
	attr_info common.AttributeInfo,
	fg_md5_alias_name bool,
) (simpleExpr SimpleExpr, err error) {
	if operator_code == common.OPERATOR_CODE_EQUAL {
		var params = filter_node.FilterParams
		simpleExpr, err = makeParamCondition(
			Equal,
			alias_name,
			attr_info,
			params,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_NOT_EQUAL {
		var params = filter_node.FilterParams
		simpleExpr, err = makeParamCondition(
			NotEqual,
			alias_name,
			attr_info,
			params,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_LT {
		var params = filter_node.FilterParams
		simpleExpr, err = makeParamCondition(
			LT,
			alias_name,
			attr_info,
			params,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_LTE {
		var params = filter_node.FilterParams
		simpleExpr, err = makeParamCondition(
			LTE,
			alias_name,
			attr_info,
			params,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_GT {
		var params = filter_node.FilterParams
		simpleExpr, err = makeParamCondition(
			GT,
			alias_name,
			attr_info,
			params,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_GTE {
		var params = filter_node.FilterParams
		simpleExpr, err = makeParamCondition(
			GTE,
			alias_name,
			attr_info,
			params,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_IN {
		simpleExpr, err = makeParamCondition(
			In,
			alias_name,
			attr_info,
			filter_node.FilterParams,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_NOT_IN {
		simpleExpr, err = makeParamCondition(
			NotIn,
			alias_name,
			attr_info,
			filter_node.FilterParams,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_LIKE {
		var params = filter_node.FilterParams
		simpleExpr, err = makeParamCondition(
			Like,
			alias_name,
			attr_info,
			params,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_LEFT_LIKE {
		var params = filter_node.FilterParams
		simpleExpr, err = makeParamCondition(
			LeftLike,
			alias_name,
			attr_info,
			params,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else if operator_code == common.OPERATOR_CODE_RIGHT_LIKE {
		var params = filter_node.FilterParams
		simpleExpr, err = makeParamCondition(
			RightLike,
			alias_name,
			attr_info,
			params,
			fg_md5_alias_name,
		)
		if err != nil {
			err = fmt.Errorf("build_filter_by_operator_code error: %v", err)
			return
		}
	} else {
		err = fmt.Errorf("unsupport operator code: %s", operator_code)
		return
	}
	return
}

func makeParamCondition(
	exprType EExprType,
	alias_name string,
	attr_info common.AttributeInfo,
	values []interface{},
	fg_md5_alias_name bool,
) (expr_temp SimpleExpr, err error) {
	if fg_md5_alias_name {
		alias_name = Md5(alias_name)
	}
	valueType, err := getValueType(attr_info)
	if err != nil {
		err = fmt.Errorf("make condition error: attr: %s unsupport value type: %s", attr_info.Name, attr_info.ValueType)
		return
	}

	expr_temp = SimpleExpr{
		ExprType:  exprType,
		ValueType: valueType,
		Values:    values,
	}
	expr_temp.TableAliasName = alias_name
	expr_temp.ColumnName = attr_info.ColumnName
	return
}

func getValueType(attr_info common.AttributeInfo) (valueType EValueType, err error) {
	if attr_info.ValueType == common.VALUE_TYPE_STRING {
		valueType = String
		return
	} else if attr_info.ValueType == common.VALUE_TYPE_NUMBER {
		valueType = Number
		return
	} else if attr_info.ValueType == common.VALUE_TYPE_BOOL {
		valueType = Bool
		return
	} else {
		err = fmt.Errorf("attr: %s unsupport value type: %s", attr_info.Name, attr_info.ValueType)
		return
	}
}

type QuerySqlBuilder struct {
	dynQuery DynQuery
}

func NewQuerySqlBuilder(
	aq_condition common.AqCondition,
	main_entity_name string,
) (builder *QuerySqlBuilder, err error) {
	var mainDesc = register.GetDescMapByKey(main_entity_name)
	var main_table_name = mainDesc.EntityInfo.TableName
	var main_table_alias = mainDesc.EntityInfo.TableName
	var selectStatemet = SelectStatement{}
	var root_logic_node = aq_condition.LogicNode
	var orders = aq_condition.Orders
	var join_name_map = getJoinNamesFromLogicNodes(root_logic_node)
	var join_name_map1 = getJoinNamesFromOrders(join_name_map, orders)
	var join_names = make([]string, 0)
	for key := range join_name_map1 {
		join_names = append(join_names, key)
	}
	sort.Sort(ByLength(join_names))

	var parent_desc = register.GetDescMapByKey(main_entity_name)
	if parent_desc == nil {
		err = fmt.Errorf("can not get description info: %s", main_entity_name)
		return
	}
	var mainColumns = make([]ColumnRef, 0)
	for _, value := range parent_desc.AttributeInfoMap {
		if value.DataType == common.DATA_TYPE_AGG_REF ||
			value.DataType == common.DATA_TYPE_AGG_SINGLE_REF ||
			value.DataType == common.DATA_TYPE_ARRAY ||
			value.DataType == common.DATA_TYPE_SINGLE {
			continue
		}
		var column = ColumnRef{}
		if value.DataType == common.DATA_TYPE_INTERNAL_PK {
			column.FgPrimary = true
		}
		column.ColumnName = value.ColumnName
		column.TableAliasName = main_table_alias
		mainColumns = append(mainColumns, column)
	}

	selectStatemet.Selects = mainColumns

	var from = make([]FromExpr, 1)
	from[0] = FromExpr{Columns: mainColumns}
	from[0].TableName = main_table_name
	from[0].AliasName = main_table_alias
	selectStatemet.From = from

	jonExprs, err := makeJoinSelect(join_names, main_entity_name, main_table_alias, 0)
	if err != nil {
		err = fmt.Errorf("makeSqlByCondition error: %v", err)
		return
	}
	selectStatemet.Join = jonExprs

	var conditionExpr = ConditionExpression{}
	if root_logic_node != nil && len(root_logic_node.FilterNodes) > 0 {
		err = makeCondition(
			root_logic_node,
			&conditionExpr,
			main_table_alias,
			main_entity_name,
		)
		if err != nil {
			err = fmt.Errorf("makeSqlByCondition error: %v", err)
			return
		}
	}
	selectStatemet.SqlWhere = &conditionExpr

	var orderExprs = make([]OrderExpr, 0)
	orderExprs, err = makeColumnOrderBy(orders, main_entity_name, main_table_alias)
	if err != nil {
		err = fmt.Errorf("makeSqlByCondition error: %v", err)
		return
	}
	selectStatemet.Orders = orderExprs

	var dynQuery, err_dyn_query = NewDynQuery(selectStatemet)
	if nil != err_dyn_query {
		err = fmt.Errorf("makeSqlByCondition error: %v", err_dyn_query)
		return
	}

	builder = new(QuerySqlBuilder)
	builder.dynQuery = dynQuery

	return
}

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

type ByLength []string

func (s ByLength) Len() int { return len(s) }
func (s ByLength) Less(i, j int) bool {
	return len(strings.Split(s[i], ".")) < len(strings.Split(s[j], "."))
}
func (s ByLength) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
