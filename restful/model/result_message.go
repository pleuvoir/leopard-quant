package model

const (
	SuccessCode = 0
	FailCode    = 1
)

func NewResultMessage(code int, msg string) *ResultMessage {
	return &ResultMessage{Code: code, Msg: msg}
}

func Success(msg string) *ResultMessage {
	return NewResultMessage(SuccessCode, msg)
}

func Fail(msg string) *ResultMessage {
	return NewResultMessage(FailCode, msg)
}

type ResultMessage struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *ResultMessage) SetCode(code int) *ResultMessage {
	r.Code = code
	return r
}

func (r *ResultMessage) SetMsg(msg string) *ResultMessage {
	r.Msg = msg
	return r
}

func (r *ResultMessage) SetData(data interface{}) *ResultMessage {
	r.Data = data
	return r
}
