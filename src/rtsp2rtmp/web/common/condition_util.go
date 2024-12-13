package common

func GetEmptyCondition() (condition AqCondition) {
	condition = AqCondition{}
	return
}

func GetEqualCondition(fieldName string, value interface{}) (condition AqCondition) {
	filterNode := AqFilterNode{}
	filterNode.OperatorCode = OPERATOR_CODE_EQUAL
	filterNode.Name = fieldName
	filterNode.FilterParams = []interface{}{value}

	logicNode := AqLogicNode{}
	logicNode.LogicOperatorCode = LOGIC_OPERATOR_CODE_AND
	logicNode.FilterNodes = []AqFilterNode{filterNode}
	condition = AqCondition{LogicNode: &logicNode}
	return
}

func GetInCondition(fieldName string, values []interface{}) (condition AqCondition) {
	filterNode := AqFilterNode{}
	filterNode.OperatorCode = OPERATOR_CODE_IN
	filterNode.Name = fieldName
	filterNode.FilterParams = values

	logicNode := AqLogicNode{}
	logicNode.LogicOperatorCode = LOGIC_OPERATOR_CODE_AND
	logicNode.FilterNodes = []AqFilterNode{filterNode}
	condition = AqCondition{LogicNode: &logicNode}
	return
}

type EqualFilter struct {
	Name  string
	Value interface{}
}

func GetEqualConditions(equalFilters []EqualFilter) (condition AqCondition) {
	var filters = []AqFilterNode{}
	for _, equalFilter := range equalFilters {
		filterNode := AqFilterNode{}
		filterNode.OperatorCode = OPERATOR_CODE_EQUAL
		filterNode.Name = equalFilter.Name
		filterNode.FilterParams = []interface{}{equalFilter.Value}
		filters = append(filters, filterNode)
	}

	logicNode := AqLogicNode{}
	logicNode.LogicOperatorCode = LOGIC_OPERATOR_CODE_AND
	logicNode.FilterNodes = filters
	condition = AqCondition{LogicNode: &logicNode}
	return
}
