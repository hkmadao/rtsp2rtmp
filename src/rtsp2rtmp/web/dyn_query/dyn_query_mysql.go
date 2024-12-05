package dyn_query

import (
	"fmt"
	"strings"
)

type DynQueryMysql struct {
}

func (dynQueryMysql *DynQueryMysql) BuildSql(selectStatement SelectStatement) (sqlStr string, err error) {
	if len(selectStatement.From) == 0 {
		err = fmt.Errorf("from is empty")
		return
	}
	var tokens = make([]string, 0)

	// select expr
	var selectTokens = make([]string, 0)
	if len(selectStatement.Selects) == 0 {
		var mainTable = selectStatement.From[0]
		selectTokens = append(selectTokens, "`"+mainTable.AliasName+"`.*")
	} else {
		for _, selectExpr := range selectStatement.Selects {
			selectTokens = append(selectTokens, "`"+selectExpr.TableAliasName+"`.`"+selectExpr.ColumnName+"`")
		}
	}
	tokens = append(tokens, "SELECT", strings.Join(selectTokens, ","))

	// from expr
	var fromTokens = make([]string, 0)
	for _, tableRef := range selectStatement.From {
		fromTokens = append(fromTokens, "`"+tableRef.TableName+"` `"+tableRef.AliasName+"`")
	}
	tokens = append(tokens, "FROM", strings.Join(fromTokens, ","))

	// join expr
	var fullJoinTokens = make([]string, 0)
	if len(selectStatement.Join) > 0 {
		for _, joinExpr := range selectStatement.Join {
			var joinTokens = make([]string, 0)
			if joinExpr.JoinType == InnerJoin {
				joinTokens = append(joinTokens, "INNER JION `"+joinExpr.TableName+"` `"+joinExpr.AliasName+"`")
				var conditions = joinExpr.On
				conditionTokens, makeErr := makeConditionsToken(conditions)
				if makeErr != nil {
					err = fmt.Errorf("make inner join token error: %v", makeErr)
					return
				}
				joinTokens = append(joinTokens, strings.Join(conditionTokens, " "))
			} else if joinExpr.JoinType == LeftJoin {
				joinTokens = append(joinTokens, "LEFT JION `"+joinExpr.TableName+"` `"+joinExpr.AliasName+"`")
				var conditions = joinExpr.On
				conditionTokens, makeErr := makeConditionsToken(conditions)
				if makeErr != nil {
					err = fmt.Errorf("make left join token error: %v", makeErr)
					return
				}
				joinTokens = append(joinTokens, strings.Join(conditionTokens, " "))
			} else if joinExpr.JoinType == RightJoin {
				joinTokens = append(joinTokens, "RIGHT JION `"+joinExpr.TableName+"` `"+joinExpr.AliasName+"`")
				var conditions = joinExpr.On
				conditionTokens, makeErr := makeConditionsToken(conditions)
				if makeErr != nil {
					err = fmt.Errorf("make right join token error: %v", makeErr)
					return
				}
				joinTokens = append(joinTokens, strings.Join(conditionTokens, " "))
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
	if len(selectStatement.SqlWhere) > 0 {
		var conditions = selectStatement.SqlWhere
		conditionTokens, makeErr := makeConditionsToken(conditions)
		if makeErr != nil {
			err = fmt.Errorf("make where token error: %v", makeErr)
			return
		}
		whereTokens = append(whereTokens, strings.Join(conditionTokens, " "))
		tokens = append(tokens, "WHERE", strings.Join(whereTokens, " "))
	}

	// having expr
	var havingTokens = make([]string, 0)
	if len(selectStatement.Having) > 0 {
		var conditions = selectStatement.Having
		conditionTokens, makeErr := makeConditionsToken(conditions)
		if makeErr != nil {
			err = fmt.Errorf("make having token error: %v", makeErr)
			return
		}
		havingTokens = append(havingTokens, strings.Join(conditionTokens, " "))
		tokens = append(tokens, "HAVING", strings.Join(havingTokens, " "))
	}

	// order expr
	var orderTokens = make([]string, 0)
	if len(selectStatement.Orders) > 0 {
		for _, orderExpr := range selectStatement.Orders {
			if orderExpr.OrderType == ASC {
				orderTokens = append(orderTokens, "`"+orderExpr.TableAliasName+"`.`"+orderExpr.ColumnName+"` ASC")
			} else if orderExpr.OrderType == DESC {
				orderTokens = append(orderTokens, "`"+orderExpr.TableAliasName+"`.`"+orderExpr.ColumnName+"` DESC")
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

func makeConditionsToken(conditionExprs []ConditionExpression) (conditionTokens []string, err error) {
	conditionTokens = make([]string, 0)
	if len(conditionExprs) > 0 {
		for _, conditionExpr := range conditionExprs {
			var logicCode = "AND"
			if conditionExpr.ConditionType == OR {
				logicCode = "OR"
			}
			var andTokens = make([]string, 0)
			if len(conditionExpr.SimpleExprs) > 0 {
				for _, simpleExpr := range conditionExpr.SimpleExprs {
					var simpleExprTokens = make([]string, 0)
					simpleExprTokens = append(simpleExprTokens, "`"+simpleExpr.TableAliasName+"` `"+simpleExpr.ColumnName+"`")
					if simpleExpr.ExprType == IsNull {
						simpleExprTokens = append(simpleExprTokens, "IS NULL")
					} else if simpleExpr.ExprType == NotNull {
						simpleExprTokens = append(simpleExprTokens, "NOT NULL")
					} else if simpleExpr.ExprType == Like {
						simpleExprTokens = append(simpleExprTokens, "LIKE", "'%"+simpleExpr.Values[0]+"%'")
					} else if simpleExpr.ExprType == LeftLike {
						simpleExprTokens = append(simpleExprTokens, "LIKE", "'%"+simpleExpr.Values[0]+"'")
					} else if simpleExpr.ExprType == RightJoin {
						simpleExprTokens = append(simpleExprTokens, "LIKE", "'"+simpleExpr.Values[0]+"%'")
					} else if simpleExpr.ExprType == Equal {
						if simpleExpr.ValueType == String {
							simpleExprTokens = append(simpleExprTokens, "=", "'"+simpleExpr.Values[0]+"'")
						} else {
							simpleExprTokens = append(simpleExprTokens, "=", simpleExpr.Values[0])
						}
					} else if simpleExpr.ExprType == NotEqual {
						if simpleExpr.ValueType == String {
							simpleExprTokens = append(simpleExprTokens, "!=", "'"+simpleExpr.Values[0]+"'")
						} else {
							simpleExprTokens = append(simpleExprTokens, "!=", simpleExpr.Values[0])
						}
					} else if simpleExpr.ExprType == GT {
						if simpleExpr.ValueType == String {
							simpleExprTokens = append(simpleExprTokens, ">", "'"+simpleExpr.Values[0]+"'")
						} else {
							simpleExprTokens = append(simpleExprTokens, ">", simpleExpr.Values[0])
						}
					} else if simpleExpr.ExprType == GTE {
						if simpleExpr.ValueType == String {
							simpleExprTokens = append(simpleExprTokens, ">=", "'"+simpleExpr.Values[0]+"'")
						} else {
							simpleExprTokens = append(simpleExprTokens, ">=", simpleExpr.Values[0])
						}
					} else if simpleExpr.ExprType == LT {
						if simpleExpr.ValueType == String {
							simpleExprTokens = append(simpleExprTokens, "<", "'"+simpleExpr.Values[0]+"'")
						} else {
							simpleExprTokens = append(simpleExprTokens, "<", simpleExpr.Values[0])
						}
					} else if simpleExpr.ExprType == LTE {
						if simpleExpr.ValueType == String {
							simpleExprTokens = append(simpleExprTokens, "<=", "'"+simpleExpr.Values[0]+"'")
						} else {
							simpleExprTokens = append(simpleExprTokens, "<=", simpleExpr.Values[0])
						}
					} else if simpleExpr.ExprType == In {
						if simpleExpr.ValueType == String {
							simpleExprTokens = append(simpleExprTokens, "IN(", "'"+strings.Join(simpleExpr.Values, "','")+"')")
						} else {
							simpleExprTokens = append(simpleExprTokens, "IN(", strings.Join(simpleExpr.Values, ","), ")")
						}
					} else if simpleExpr.ExprType == NotIn {
						if simpleExpr.ValueType == String {
							simpleExprTokens = append(simpleExprTokens, "NOT IN(", "'"+strings.Join(simpleExpr.Values, "','")+"')")
						} else {
							simpleExprTokens = append(simpleExprTokens, "NOT IN(", strings.Join(simpleExpr.Values, ","), ")")
						}
					} else {
						err = fmt.Errorf("unsupport expr type: %d", simpleExpr.ExprType)
						return
					}

					andTokens = append(andTokens, strings.Join(simpleExprTokens, " "))
				}
			}
			conditionTokens = append(conditionTokens, strings.Join(andTokens, logicCode))
		}
	}
	return
}
