package api_err

import "github.com/zeromicro/x/errors"

var (
	ErrCodeInternalErrorCode = 10001
	ErrCodeParamErrorCode    = 10002
)

var (
	ErrCodeInvalidChainNotSupport = errors.New(20002, "invalid chain not support")
)

func NewErrorWithCodeAndMsg(code int, msg string) error {
	return errors.New(code, msg)
}
