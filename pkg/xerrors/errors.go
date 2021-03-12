package xerrors

import "fmt"

type BizError struct {
	Code ErrorCode `json:"code"`
	Msg  string    `json:"msg"`
}

func NewBizError(code ErrorCode, msg string, args ...interface{}) *BizError {
	return &BizError{
		Code: code,
		Msg:  fmt.Sprintf(msg, args...),
	}
}

func NewInvalidInputError(msg string, args ...interface{}) *BizError {
	return NewBizError(InvalidInput, msg, args...)
}

func (e *BizError) Error() string {
	return e.Msg
}

type ErrorCode int

const (
	InvalidInput   ErrorCode = 10000 + iota // 输入异常
	InvalidToken                            // 无效的Token
	ValidateFailed                          // 数据校验异常
	RecordNotFound                          // 数据不存在

	InvalidCredential // 用户名或密码错误

	Unknown = 99999
)

var (
	InvalidInputError      = NewBizError(InvalidInput, "无效的输入参数")
	InvalidCredentialError = NewBizError(InvalidCredential, "用户名或密码错误")
	InvalidTokenError      = NewBizError(InvalidToken, "无效的token")
	RecordNotFoundError    = NewBizError(RecordNotFound, "数据不存在")
)
