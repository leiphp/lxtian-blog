package response

import (
	"fmt"
	"net/http"
)

// HttpError 用于包含状态码和消息的错误类型
type HttpError struct {
	Message    string
	StatusCode int
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// NewHttpError 构造器
func NewHttpError(message string, status int) *HttpError {
	return &HttpError{
		Message:    message,
		StatusCode: status,
	}
}

// 常用错误定义
var (
	ErrTokenInvalid = NewHttpError("登录已过期，请重新登录", http.StatusUnauthorized)
	ErrForbidden    = NewHttpError("无权限访问", http.StatusForbidden)
	ErrBadRequest   = NewHttpError("请求参数错误", http.StatusBadRequest)
	ErrServerError  = NewHttpError("服务器错误", http.StatusInternalServerError)
)
