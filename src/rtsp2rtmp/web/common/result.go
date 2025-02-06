package common

import "errors"

type AppResult struct {
	Status  uint32      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageInfoInput struct {
	// 当前页码
	PageIndex uint64 `json:"pageIndex"`
	// 分页大小
	PageSize uint64 `json:"pageSize"`
	// 总记录数
	TotalCount uint64 `json:"totalCount"`
}

type PageInfo struct {
	PageInfoInput PageInfoInput `json:"pageInfoInput"`
	DataList      []interface{} `json:"dataList"`
}

type DeleteRefErrorMessageVO struct {
	// 被引用id
	IdData string `json:"idData"`
	// 错误信息
	Message string `json:"Message"`
	// 被引用类名
	SourceClassName string `json:"sourceClassName"`
	// 引用类名
	RefClassName string `json:"refClassName"`
}

func (appResult *AppResult) IsFailed() bool {
	return appResult.Status != 0
}

func ErrorResult(msg string) AppResult {
	return AppResult{Status: 1, Message: msg}
}

func SuccessResultData(data interface{}) AppResult {
	return AppResult{Status: 0, Data: data}
}

func SuccessResultMsg(msg string) AppResult {
	return AppResult{Status: 0, Message: msg}
}

func SuccessResultWithMsg(msg string, data interface{}) AppResult {
	return AppResult{Status: 0, Message: msg, Data: data}
}

type Rtmp2FlvCustomError struct {
	kind int
	err  error
}

func CustomError(errMsg string) *Rtmp2FlvCustomError {
	return &Rtmp2FlvCustomError{kind: 0, err: errors.New(errMsg)}
}

func InternalError(err error) *Rtmp2FlvCustomError {
	if err == nil {
		panic("err param is nil")
	}
	return &Rtmp2FlvCustomError{kind: 1, err: err}
}

func (ce *Rtmp2FlvCustomError) IsCustomError() bool {
	return ce.kind == 0
}

func (ce *Rtmp2FlvCustomError) Error() string {
	return ce.err.Error()
}
