package errors

import "net/http"

// AppError 应用错误类型
type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// New 创建新的应用错误
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// BadRequest 创建 400 错误
func BadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, message, err)
}

// InternalServerError 创建 500 错误
func InternalServerError(message string, err error) *AppError {
	return New(http.StatusInternalServerError, message, err)
}

// NotFound 创建 404 错误
func NotFound(message string, err error) *AppError {
	return New(http.StatusNotFound, message, err)
}

// Unauthorized 创建 401 错误
func Unauthorized(message string, err error) *AppError {
	return New(http.StatusUnauthorized, message, err)
}

// Forbidden 创建 403 错误
func Forbidden(message string, err error) *AppError {
	return New(http.StatusForbidden, message, err)
}
