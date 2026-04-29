package errors

import "net/http"

// ErrorType 错误类型
type ErrorType string

const (
	// ErrorTypeValidation 验证错误
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeDatabase 数据库错误
	ErrorTypeDatabase ErrorType = "database"
	// ErrorTypeAuth 认证错误
	ErrorTypeAuth ErrorType = "auth"
	// ErrorTypeInternal 内部错误
	ErrorTypeInternal ErrorType = "internal"
	// ErrorTypeNotFound 资源不存在错误
	ErrorTypeNotFound ErrorType = "not_found"
)

// TypedError 带类型的错误
type TypedError struct {
	*AppError
	Type ErrorType `json:"type"`
}

// NewTypedError 创建带类型的错误
func NewTypedError(code int, errType ErrorType, message string, err error) *TypedError {
	return &TypedError{
		AppError: New(code, message, err),
		Type:     errType,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(message string, err error) *TypedError {
	return NewTypedError(http.StatusBadRequest, ErrorTypeValidation, message, err)
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(message string, err error) *TypedError {
	return NewTypedError(http.StatusInternalServerError, ErrorTypeDatabase, message, err)
}

// NewAuthError 创建认证错误
func NewAuthError(message string, err error) *TypedError {
	return NewTypedError(http.StatusUnauthorized, ErrorTypeAuth, message, err)
}

// NewInternalError 创建内部错误
func NewInternalError(message string, err error) *TypedError {
	return NewTypedError(http.StatusInternalServerError, ErrorTypeInternal, message, err)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(message string, err error) *TypedError {
	return NewTypedError(http.StatusNotFound, ErrorTypeNotFound, message, err)
}
