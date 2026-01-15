/**
 * [INPUT]: 依赖 testing 标准库; 依赖 errors.go 的所有错误类型
 * [OUTPUT]: 对外提供错误类型的单元测试，覆盖 ClientError, AuthenticationError, RateLimitError, NotFoundError, ValidationError
 * [POS]: SDK 根目录的测试层，验证 errors.go 的正确性
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */
package memu

import (
	"testing"
)

// ============================================================
// 测试: ClientError
// ============================================================

func TestClientError_Error(t *testing.T) {
	statusCode := 400
	err := &ClientError{
		Message:    "Test error",
		StatusCode: &statusCode,
		Response:   map[string]interface{}{"detail": "test"},
	}

	expected := "MemU API error (status 400): Test error"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestClientError_ErrorWithoutStatusCode(t *testing.T) {
	err := &ClientError{
		Message:  "Test error",
		Response: map[string]interface{}{"detail": "test"},
	}

	expected := "MemU API error: Test error"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// ============================================================
// 测试: AuthenticationError
// ============================================================

func TestNewAuthenticationError_WithMessage(t *testing.T) {
	statusCode := 401
	response := map[string]interface{}{
		"message": "Invalid API key",
	}

	err := NewAuthenticationError(&statusCode, response)

	if err.Message != "Invalid API key" {
		t.Errorf("Expected message 'Invalid API key', got '%s'", err.Message)
	}
	if err.StatusCode == nil || *err.StatusCode != 401 {
		t.Errorf("Expected status code 401, got %v", err.StatusCode)
	}
}

func TestNewAuthenticationError_DefaultMessage(t *testing.T) {
	statusCode := 401
	err := NewAuthenticationError(&statusCode, nil)

	expected := "Authentication failed. Please check your API key."
	if err.Message != expected {
		t.Errorf("Expected default message '%s', got '%s'", expected, err.Message)
	}
}

// ============================================================
// 测试: RateLimitError
// ============================================================

func TestNewRateLimitError_WithRetryAfter(t *testing.T) {
	statusCode := 429
	retryAfter := 60.0
	response := map[string]interface{}{
		"message": "Rate limit exceeded",
	}

	err := NewRateLimitError("Rate limit exceeded", &retryAfter, &statusCode, response)

	if err.Message != "Rate limit exceeded" {
		t.Errorf("Expected message 'Rate limit exceeded', got '%s'", err.Message)
	}
	if err.RetryAfter == nil || *err.RetryAfter != 60.0 {
		t.Errorf("Expected retry_after 60.0, got %v", err.RetryAfter)
	}
}

func TestNewRateLimitError_DefaultMessage(t *testing.T) {
	statusCode := 429
	err := NewRateLimitError("Rate limit exceeded. Please try again later.", nil, &statusCode, nil)

	expected := "Rate limit exceeded. Please try again later."
	if err.Message != expected {
		t.Errorf("Expected default message '%s', got '%s'", expected, err.Message)
	}
}

// ============================================================
// 测试: NotFoundError
// ============================================================

func TestNewNotFoundError_WithMessage(t *testing.T) {
	statusCode := 404
	response := map[string]interface{}{
		"message": "Task not found",
	}

	err := NewNotFoundError("/v3/task/123", &statusCode, response)

	if err.Message != "Task not found" {
		t.Errorf("Expected message 'Task not found', got '%s'", err.Message)
	}
	if err.StatusCode == nil || *err.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %v", err.StatusCode)
	}
}

func TestNewNotFoundError_DefaultMessage(t *testing.T) {
	statusCode := 404
	err := NewNotFoundError("/v3/task/123", &statusCode, nil)

	expected := "Resource not found: /v3/task/123"
	if err.Message != expected {
		t.Errorf("Expected default message '%s', got '%s'", expected, err.Message)
	}
}

// ============================================================
// 测试: ValidationError
// ============================================================

func TestNewValidationError_WithMessage(t *testing.T) {
	statusCode := 422
	response := map[string]interface{}{
		"message": "user_id is required",
	}

	err := NewValidationError(&statusCode, response)

	if err.Message != "user_id is required" {
		t.Errorf("Expected message 'user_id is required', got '%s'", err.Message)
	}
	if err.StatusCode == nil || *err.StatusCode != 422 {
		t.Errorf("Expected status code 422, got %v", err.StatusCode)
	}
}

func TestNewValidationError_DefaultMessage(t *testing.T) {
	statusCode := 422
	err := NewValidationError(&statusCode, nil)

	expected := "Request validation failed. Please check your request parameters."
	if err.Message != expected {
		t.Errorf("Expected default message '%s', got '%s'", expected, err.Message)
	}
}

// ============================================================
// 测试: 错误类型断言
// ============================================================

func TestErrorTypeAssertion(t *testing.T) {
	// 测试 AuthenticationError 可以被断言为 ClientError
	authErr := NewAuthenticationError(intPtr(401), nil)
	if _, ok := interface{}(authErr).(*AuthenticationError); !ok {
		t.Error("AuthenticationError should be assertable as *AuthenticationError")
	}

	// 测试 RateLimitError 可以被断言为 ClientError
	rateLimitErr := NewRateLimitError("Rate limit", nil, intPtr(429), nil)
	if _, ok := interface{}(rateLimitErr).(*RateLimitError); !ok {
		t.Error("RateLimitError should be assertable as *RateLimitError")
	}

	// 测试 NotFoundError 可以被断言为 ClientError
	notFoundErr := NewNotFoundError("/path", intPtr(404), nil)
	if _, ok := interface{}(notFoundErr).(*NotFoundError); !ok {
		t.Error("NotFoundError should be assertable as *NotFoundError")
	}

	// 测试 ValidationError 可以被断言为 ClientError
	validationErr := NewValidationError(intPtr(422), nil)
	if _, ok := interface{}(validationErr).(*ValidationError); !ok {
		t.Error("ValidationError should be assertable as *ValidationError")
	}
}

// ============================================================
// 辅助函数
// ============================================================

func intPtr(i int) *int {
	return &i
}
