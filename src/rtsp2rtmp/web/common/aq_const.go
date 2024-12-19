package common

type EDoStatus = uint32

// do_status const
const (
	DO_UNCHANGE = 1
	DO_UPDATE   = 2
	DO_NEW      = 3
	DO_DELETE   = 4
)

type ELogicOperatorCode = string

// logic operator code const
const (
	LOGIC_OPERATOR_CODE_AND = "and"
	LOGIC_OPERATOR_CODE_OR  = "or"
)

type EOperatorCode = string

// operator code const
const (
	OPERATOR_CODE_EQUAL      = "equal"
	OPERATOR_CODE_NOT_EQUAL  = "notEqual"
	OPERATOR_CODE_IN         = "in"
	OPERATOR_CODE_NOT_IN     = "notIn"
	OPERATOR_CODE_LT         = "lessThan"
	OPERATOR_CODE_LTE        = "lessThanEqual"
	OPERATOR_CODE_GT         = "greaterThan"
	OPERATOR_CODE_GTE        = "greaterThanEqual"
	OPERATOR_CODE_LIKE       = "like"
	OPERATOR_CODE_LEFT_LIKE  = "leftLike"
	OPERATOR_CODE_RIGHT_LIKE = "rightLike"
)

type EOrderDirection = string

// order direction const
const (
	ORDER_DIRECTION_ASC  = "ASC"
	ORDER_DIRECTION_DESC = "DESC"
)
