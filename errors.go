/**
 * [INPUT]: 依赖 fmt 的格式化能力
 * [OUTPUT]: 对外提供 ClientError, AuthenticationError, RateLimitError, NotFoundError, ValidationError 类型及其构造函数
 * [POS]: SDK 根目录的错误处理层，定义错误类型层次，被 client.go 和用户代码消费
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */
package memu

import (
	"fmt"
)

// ClientError is the base error type for all MemU SDK errors.
type ClientError struct {
	Message    string
	StatusCode *int
	Response   map[string]interface{}
}

// Error implements the error interface.
func (e *ClientError) Error() string {
	if e.StatusCode != nil {
		return fmt.Sprintf("MemU API error (status %d): %s", *e.StatusCode, e.Message)
	}
	return fmt.Sprintf("MemU API error: %s", e.Message)
}

// AuthenticationError is raised when API authentication fails (401).
type AuthenticationError struct {
	*ClientError
}

// NewAuthenticationError creates a new AuthenticationError.
func NewAuthenticationError(statusCode *int, response map[string]interface{}) *AuthenticationError {
	message := "Authentication failed. Please check your API key."
	if response != nil {
		if msg, ok := response["message"].(string); ok && msg != "" {
			message = msg
		}
	}
	return &AuthenticationError{
		ClientError: &ClientError{
			Message:    message,
			StatusCode: statusCode,
			Response:   response,
		},
	}
}

// RateLimitError is raised when API rate limit is exceeded (429).
type RateLimitError struct {
	*ClientError
	RetryAfter *float64
}

// NewRateLimitError creates a new RateLimitError.
func NewRateLimitError(message string, retryAfter *float64, statusCode *int, response map[string]interface{}) *RateLimitError {
	return &RateLimitError{
		ClientError: &ClientError{
			Message:    message,
			StatusCode: statusCode,
			Response:   response,
		},
		RetryAfter: retryAfter,
	}
}

// NotFoundError is raised when a requested resource is not found (404).
type NotFoundError struct {
	*ClientError
}

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(path string, statusCode *int, response map[string]interface{}) *NotFoundError {
	message := fmt.Sprintf("Resource not found: %s", path)
	if response != nil {
		if msg, ok := response["message"].(string); ok && msg != "" {
			message = msg
		}
	}
	return &NotFoundError{
		ClientError: &ClientError{
			Message:    message,
			StatusCode: statusCode,
			Response:   response,
		},
	}
}

// ValidationError is raised when request validation fails (422).
type ValidationError struct {
	*ClientError
}

// NewValidationError creates a new ValidationError.
func NewValidationError(statusCode *int, response map[string]interface{}) *ValidationError {
	message := "Request validation failed. Please check your request parameters."
	if response != nil {
		if msg, ok := response["message"].(string); ok && msg != "" {
			message = msg
		}
	}
	return &ValidationError{
		ClientError: &ClientError{
			Message:    message,
			StatusCode: statusCode,
			Response:   response,
		},
	}
}

// NewClientError creates a new ClientError.
func NewClientError(message string, statusCode *int, response map[string]interface{}) *ClientError {
	return &ClientError{
		Message:    message,
		StatusCode: statusCode,
		Response:   response,
	}
}
