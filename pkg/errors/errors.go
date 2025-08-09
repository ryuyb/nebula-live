package errors

import "errors"

// 通用错误
var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrInternalServer = errors.New("internal server error")
)

// API错误响应
type APIError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// NewAPIError 创建API错误
func NewAPIError(code int, error, message string) *APIError {
	return &APIError{
		Code:    code,
		Error:   error,
		Message: message,
	}
}
