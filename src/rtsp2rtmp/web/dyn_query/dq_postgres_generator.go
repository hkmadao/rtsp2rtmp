package dyn_query

import (
	"fmt"
	"strings"
)

type DynQueryPostgres struct {
}

func (dynQueryPostgres *DynQueryPostgres) BuildSql(selectStatement SelectStatement) (sqlStr string, params []interface{}, err error) {
	if len(selectStatement.From) == 0 {
		err = fmt.Errorf("from is empty")
		return
	}
	var tokens = make([]string, 0)

	// select expr
	var selectTokens = make([]string, 0)
	if len(selectStatement.Selects) == 0 {
		var mainTable = selectStatement.From[0]
		selectTokens = append(selectTokens, fmt.Sprintf("\"%s\".*", mainTable.AliasName))
	} else {
		for _, selectExpr := range selectStatement.Selects {
			selectTokens = append(selectTokens, fmt.Sprintf("\"%s\".\"%s\"", selectExpr.TableAliasName, selectExpr.ColumnName))
		}
	}
	tokens = append(tokens, "SELECT", strings.Join(selectTokens, ","))

	// from expr
	var fromTokens = make([]string, 0)
	for _, tableRef := range selectStatement.From {
		fromTokens = append(fromTokens, fmt.Sprintf("\"%s\" \"%s\"", tableRef.TableName, tableRef.AliasName))
	}
	tokens = append(tokens, "FROM", strings.Join(fromTokens, ","))

	// join expr
	var fullJoinTokens = make([]string, 0)
	if len(selectStatement.Join) > 0 {
		for _, joinExpr := range selectStatement.Join {
			var joinTokens = make([]string, 0)
			if joinExpr.JoinType == InnerJoin {
				joinTokens = append(joinTokens, fmt.Sprintf("INNER JOIN \"%s\" \"%s\" ON", joinExpr.TableName, joinExpr.AliasName))
				var conditionExpr = joinExpr.On
				if nil == conditionExpr {
					err = fmt.Errorf("inner join condition is empty")
					return
				}
				if len(conditionExpr.SimpleExprs) != 2 {
					err = fmt.Errorf("inner join condition error")
					return
				}
				var simpleExprs = conditionExpr.SimpleExprs
				conditionToken := fmt.Sprintf("\"%s\".\"%s\" = \"%s\".\"%s\"", simpleExprs[0].TableAliasName, simpleExprs[0].ColumnName, simpleExprs[1].TableAliasName, simpleExprs[1].ColumnName)
				joinTokens = append(joinTokens, conditionToken)
			} else if joinExpr.JoinType == LeftJoin {
				joinTokens = append(joinTokens, fmt.Sprintf("LEFT JOIN \"%s\" \"%s\" ON", joinExpr.TableName, joinExpr.AliasName))
				var conditionExpr = joinExpr.On
				if nil == conditionExpr {
					err = fmt.Errorf("left join condition is empty")
					return
				}
				if len(conditionExpr.SimpleExprs) != 2 {
					err = fmt.Errorf("left join condition error")
					return
				}
				var simpleExprs = conditionExpr.SimpleExprs
				conditionToken := fmt.Sprintf("\"%s\".\"%s\" = \"%s\".\"%s\"", simpleExprs[0].TableAliasName, simpleExprs[0].ColumnName, simpleExprs[1].TableAliasName, simpleExprs[1].ColumnName)
				joinTokens = append(joinTokens, conditionToken)
			} else if joinExpr.JoinType == RightJoin {
				joinTokens = append(joinTokens, fmt.Sprintf("RIGHT JOIN \"%s\" \"%s\" ON", joinExpr.TableName, joinExpr.AliasName))
				var conditionExpr = joinExpr.On
				if nil == conditionExpr {
					err = fmt.Errorf("right join condition is empty")
					return
				}
				if len(conditionExpr.SimpleExprs) != 2 {
					err = fmt.Errorf("right join condition error")
					return
				}
				var simpleExprs = conditionExpr.SimpleExprs
				conditionToken := fmt.Sprintf("\"%s\".\"%s\" = \"%s\".\"%s\"", simpleExprs[0].TableAliasName, simpleExprs[0].ColumnName, simpleExprs[1].TableAliasName, simpleExprs[1].ColumnName)
				joinTokens = append(joinTokens, conditionToken)
			} else {
				err = fmt.Errorf("unsupport join: %d", joinExpr.JoinType)
				return
			}
			fullJoinTokens = append(fullJoinTokens, strings.Join(joinTokens, " "))
		}
	}
	tokens = append(tokens, strings.Join(fullJoinTokens, " "))

	// where expr
	var whereTokens = make([]string, 0)
	params = make([]interface{}, 0)
	if nil != selectStatement.SqlWhere {
		var conditionExpr = selectStatement.SqlWhere
		conditionToken, whereParams, makeErr := makeConditionsTokenPostgres(*conditionExpr)
		if makeErr != nil {
			err = fmt.Errorf("make where token error: %v", makeErr)
			return
		}
		params = append(params, whereParams...)
		whereTokens = append(whereTokens, conditionToken)
		tokens = append(tokens, "WHERE", strings.Join(whereTokens, " "))
	}

	// having expr
	var havingTokens = make([]string, 0)
	if nil != selectStatement.Having {
		var conditionExpr = selectStatement.Having
		conditionToken, havingParams, makeErr := makeConditionsTokenPostgres(*conditionExpr)
		if makeErr != nil {
			err = fmt.Errorf("make having token error: %v", makeErr)
			return
		}
		params = append(params, havingParams...)
		havingTokens = append(havingTokens, conditionToken)
		tokens = append(tokens, "HAVING", strings.Join(havingTokens, " "))
	}

	// order expr
	var orderTokens = make([]string, 0)
	if len(selectStatement.Orders) > 0 {
		for _, orderExpr := range selectStatement.Orders {
			if orderExpr.OrderType == ASC {
				orderTokens = append(orderTokens, fmt.Sprintf("\"%s\".\"%s\" ASC", orderExpr.TableAliasName, orderExpr.ColumnName))
			} else if orderExpr.OrderType == DESC {
				orderTokens = append(orderTokens, fmt.Sprintf("\"%s\".\"%s\" DESC", orderExpr.TableAliasName, orderExpr.ColumnName))
			} else {
				err = fmt.Errorf("unsupport OrderType: %d", orderExpr.OrderType)
				return
			}
		}
		tokens = append(tokens, "ORDER BY", strings.Join(orderTokens, ","))
	}

	sqlStr = strings.Join(tokens, " ")
	return
}

func makeConditionsTokenPostgres(conditionExpr ConditionExpression) (conditionToken string, params []interface{}, err error) {
	var logicCode = "AND"
	if conditionExpr.ConditionType == OR {
		logicCode = "OR"
	}
	var conditionTokens = make([]string, 0)
	if len(conditionExpr.SimpleExprs) > 0 {
		for _, simpleExpr := range conditionExpr.SimpleExprs {
			var simpleExprTokens = make([]string, 0)
			simpleExprTokens = append(simpleExprTokens, fmt.Sprintf("\"%s\".\"%s\"", simpleExpr.TableAliasName, simpleExpr.ColumnName))
			if simpleExpr.ExprType == IsNull {
				simpleExprTokens = append(simpleExprTokens, "IS NULL")
			} else if simpleExpr.ExprType == NotNull {
				simpleExprTokens = append(simpleExprTokens, "NOT NULL")
			} else if simpleExpr.ExprType == Like {
				simpleExprTokens = append(simpleExprTokens, "LIKE", "?")
				params = append(params, "%"+fmt.Sprintf("%s", simpleExpr.Values[0])+"%")
			} else if simpleExpr.ExprType == LeftLike {
				simpleExprTokens = append(simpleExprTokens, "LIKE", "?")
				params = append(params, "%"+fmt.Sprintf("%s", simpleExpr.Values[0]))
			} else if simpleExpr.ExprType == RightJoin {
				simpleExprTokens = append(simpleExprTokens, "LIKE", "?")
				params = append(params, fmt.Sprintf("%s", simpleExpr.Values[0])+"%")
			} else if simpleExpr.ExprType == Equal {
				simpleExprTokens = append(simpleExprTokens, "=", "?")
				params = append(params, simpleExpr.Values[0])
			} else if simpleExpr.ExprType == NotEqual {
				simpleExprTokens = append(simpleExprTokens, "!=", "?")
				params = append(params, simpleExpr.Values[0])
			} else if simpleExpr.ExprType == GT {
				simpleExprTokens = append(simpleExprTokens, ">", "?")
				params = append(params, simpleExpr.Values[0])
			} else if simpleExpr.ExprType == GTE {
				simpleExprTokens = append(simpleExprTokens, ">=", "?")
				params = append(params, simpleExpr.Values[0])
			} else if simpleExpr.ExprType == LT {
				simpleExprTokens = append(simpleExprTokens, "<", "?")
				params = append(params, simpleExpr.Values[0])
			} else if simpleExpr.ExprType == LTE {
				simpleExprTokens = append(simpleExprTokens, "<=", "?")
				params = append(params, simpleExpr.Values[0])
			} else if simpleExpr.ExprType == In {
				qustionMarkArr := getQuestionMarkArr(len(simpleExpr.Values))
				simpleExprTokens = append(simpleExprTokens, "IN(", strings.Join(qustionMarkArr, ","), ")")
				params = append(params, simpleExpr.Values...)
			} else if simpleExpr.ExprType == NotIn {
				qustionMarkArr := getQuestionMarkArr(len(simpleExpr.Values))
				simpleExprTokens = append(simpleExprTokens, "NOT IN(", strings.Join(qustionMarkArr, ","), ")")
				params = append(params, simpleExpr.Values...)
			} else {
				err = fmt.Errorf("unsupport expr type: %d", simpleExpr.ExprType)
				return
			}

			conditionTokens = append(conditionTokens, strings.Join(simpleExprTokens, " "))
		}
		var child = conditionExpr.Child
		if nil != child && len(child.SimpleExprs) > 0 {
			childConditionToken, childParams, child_err := makeConditionsTokenPostgres(*child)
			if child_err != nil {
				err = child_err
				return
			}
			params = append(params, childParams...)
			conditionToken = strings.Join(conditionTokens, " "+logicCode+" ") + " " + logicCode + " ( " + childConditionToken + " )"
		} else {
			conditionToken = strings.Join(conditionTokens, " "+logicCode+" ")
		}
	}
	return
}
