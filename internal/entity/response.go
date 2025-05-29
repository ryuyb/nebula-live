package entity

// Response 是一个通用的响应结构，包含HTTP状态码、消息和数据字段，其中data是泛型类型
type Response[T any] struct {
	Code    int    `json:"code"`           // HTTP状态码
	Message string `json:"message"`        // 响应消息
	Data    T      `json:"data,omitempty"` // 响应数据，可选
}

// NewResponse 创建一个新的Response实例
func NewResponse[T any](code int, message string, data T) *Response[T] {
	return &Response[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// SuccessResponse 创建一个成功的响应，默认使用HTTP 200状态码
func SuccessResponse[T any](data T) *Response[T] {
	return NewResponse(200, "OK", data)
}

// ErrorResponse 创建一个错误的响应，使用指定的HTTP状态码和消息
func ErrorResponse(code int, message string) *Response[any] {
	var zero any
	return NewResponse(code, message, zero)
}

// BadRequestResponse 创建一个HTTP 400 Bad Request错误响应
func BadRequestResponse() *Response[any] {
	return ErrorResponse(400, "Bad Request")
}

// UnauthorizedResponse 创建一个HTTP 401 Unauthorized错误响应
func UnauthorizedResponse() *Response[any] {
	return ErrorResponse(401, "Unauthorized")
}

// ForbiddenResponse 创建一个HTTP 403 Forbidden错误响应
func ForbiddenResponse() *Response[any] {
	return ErrorResponse(403, "Forbidden")
}

// NotFoundResponse 创建一个HTTP 404 Not Found错误响应
func NotFoundResponse() *Response[any] {
	return ErrorResponse(404, "Not Found")
}

// InternalServerErrorResponse 创建一个HTTP 500 Internal Server Error错误响应
func InternalServerErrorResponse() *Response[any] {
	return ErrorResponse(500, "Internal Server Error")
}
