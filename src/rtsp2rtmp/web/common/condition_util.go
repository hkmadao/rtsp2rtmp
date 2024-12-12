package common

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
