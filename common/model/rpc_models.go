package model

const success = 0

type ResultMessage struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Method  string `json:"method"`
}

type RequestMessage struct {
	Method string `json:"method"`
	Param  any    `json:"param"`
}

func NewSuccess(method string, data any) ResultMessage {
	return ResultMessage{
		Code:    success,
		Message: "success",
		Data:    data,
		Method:  method,
	}
}

func NewFail(code int32, method string, message string) ResultMessage {
	return ResultMessage{
		Code:    code,
		Message: message,
		Data:    nil,
		Method:  method,
	}
}
