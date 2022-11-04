package response

type ResponseMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const FailureCode = -1
const SuccessCode = 0
const SuccessMessage = "SUCCESS"
const FailureMsgToken = "token验证失败"
const FailureNullToken = "token 为空"
const FailureParserToken = "token 解析失败"

func SuccessMsg(data interface{}) *ResponseMsg {
	msg := &ResponseMsg{
		Code: SuccessCode,
		Msg:  SuccessMessage,
		Data: data,
	}
	return msg
}

func FailMsg(msg string) *ResponseMsg {
	msgObj := &ResponseMsg{
		Code: FailureCode,
		Msg:  msg,
	}
	return msgObj
}
