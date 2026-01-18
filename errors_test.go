// Package memu provides unit tests for error types.
// This file validates error types and Error() methods.
package memu

import (
	"strings"
	"testing"
)

// Helper functions for creating pointers.
func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

// TestClientError tests ClientError.
func TestClientError(t *testing.T) {
	statusCode := 500
	err := NewClientError("Test error", &statusCode, nil)

	if err.Message != "Test error" {
		t.Errorf("expected Message 'Test error', got '%s'", err.Message)
	}
	if err.StatusCode == nil || *err.StatusCode != 500 {
		t.Errorf("expected StatusCode 500, got %v", err.StatusCode)
	}
}

func TestClientError_Error(t *testing.T) {
	statusCode := 500
	err := NewClientError("Test error", &statusCode, nil)

	errStr := err.Error()
	if !strings.Contains(errStr, "500") {
		t.Errorf("expected error string to contain '500', got '%s'", errStr)
	}
	if !strings.Contains(errStr, "Test error") {
		t.Errorf("expected error string to contain 'Test error', got '%s'", errStr)
	}
}

func TestClientError_NilStatusCode(t *testing.T) {
	err := NewClientError("Test error", nil, nil)

	errStr := err.Error()
	if strings.Contains(errStr, "status") {
		t.Errorf("expected error string without status code, got '%s'", errStr)
	}
	if !strings.Contains(errStr, "Test error") {
		t.Errorf("expected error string to contain 'Test error', got '%s'", errStr)
	}
}

func TestClientError_WithResponse(t *testing.T) {
	statusCode := 400
	response := map[string]interface{}{
		"error":   "bad_request",
		"message": "Invalid parameters",
	}
	err := NewClientError("Bad request", &statusCode, response)

	if err.Response == nil {
		t.Fatal("expected Response to not be nil")
	}
	if err.Response["error"] != "bad_request" {
		t.Errorf("expected Response['error'] 'bad_request', got '%v'", err.Response["error"])
	}
}

// TestAuthenticationError tests AuthenticationError.
func TestAuthenticationError(t *testing.T) {
	statusCode := 401
	err := NewAuthenticationError(&statusCode, nil)

	if err.StatusCode == nil || *err.StatusCode != 401 {
		t.Errorf("expected StatusCode 401, got %v", err.StatusCode)
	}

	// Check that it's a ClientError
	var _ *ClientError = err.ClientError
}

func TestAuthenticationError_DefaultMessage(t *testing.T) {
	statusCode := 401
	err := NewAuthenticationError(&statusCode, nil)

	if !strings.Contains(err.Message, "Authentication") {
		t.Errorf("expected default message to contain 'Authentication', got '%s'", err.Message)
	}
}

func TestAuthenticationError_CustomMessage(t *testing.T) {
	statusCode := 401
	response := map[string]interface{}{
		"message": "Invalid API key provided",
	}
	err := NewAuthenticationError(&statusCode, response)

	if err.Message != "Invalid API key provided" {
		t.Errorf("expected custom message 'Invalid API key provided', got '%s'", err.Message)
	}
}

func TestAuthenticationError_TypeAssertion(t *testing.T) {
	statusCode := 401
	err := NewAuthenticationError(&statusCode, nil)

	// Test type assertion
	var genericErr error = err
	if _, ok := genericErr.(*AuthenticationError); !ok {
		t.Error("expected error to be *AuthenticationError")
	}
}

// TestRateLimitError tests RateLimitError.
func TestRateLimitError(t *testing.T) {
	statusCode := 429
	retryAfter := 30.0
	err := NewRateLimitError("Rate limit exceeded", &retryAfter, &statusCode, nil)

	if err.StatusCode == nil || *err.StatusCode != 429 {
		t.Errorf("expected StatusCode 429, got %v", err.StatusCode)
	}
	if err.RetryAfter == nil || *err.RetryAfter != 30.0 {
		t.Errorf("expected RetryAfter 30.0, got %v", err.RetryAfter)
	}
}

func TestRateLimitError_NilRetryAfter(t *testing.T) {
	statusCode := 429
	err := NewRateLimitError("Rate limit exceeded", nil, &statusCode, nil)

	if err.RetryAfter != nil {
		t.Errorf("expected RetryAfter to be nil, got %v", err.RetryAfter)
	}
}

func TestRateLimitError_TypeAssertion(t *testing.T) {
	statusCode := 429
	err := NewRateLimitError("Rate limit", nil, &statusCode, nil)

	var genericErr error = err
	if _, ok := genericErr.(*RateLimitError); !ok {
		t.Error("expected error to be *RateLimitError")
	}
}

// TestNotFoundError tests NotFoundError.
func TestNotFoundError(t *testing.T) {
	statusCode := 404
	err := NewNotFoundError("/api/v3/memory/task/123", &statusCode, nil)

	if err.StatusCode == nil || *err.StatusCode != 404 {
		t.Errorf("expected StatusCode 404, got %v", err.StatusCode)
	}
	if !strings.Contains(err.Message, "not found") {
		t.Errorf("expected message to contain 'not found', got '%s'", err.Message)
	}
}

func TestNotFoundError_CustomMessage(t *testing.T) {
	statusCode := 404
	response := map[string]interface{}{
		"message": "Task not found",
	}
	err := NewNotFoundError("/api/v3/memory/task/123", &statusCode, response)

	if err.Message != "Task not found" {
		t.Errorf("expected custom message 'Task not found', got '%s'", err.Message)
	}
}

func TestNotFoundError_TypeAssertion(t *testing.T) {
	statusCode := 404
	err := NewNotFoundError("/path", &statusCode, nil)

	var genericErr error = err
	if _, ok := genericErr.(*NotFoundError); !ok {
		t.Error("expected error to be *NotFoundError")
	}
}

// TestValidationError tests ValidationError.
func TestValidationError(t *testing.T) {
	statusCode := 422
	err := NewValidationError(&statusCode, nil)

	if err.StatusCode == nil || *err.StatusCode != 422 {
		t.Errorf("expected StatusCode 422, got %v", err.StatusCode)
	}
}

func TestValidationError_DefaultMessage(t *testing.T) {
	statusCode := 422
	err := NewValidationError(&statusCode, nil)

	if !strings.Contains(err.Message, "validation") || !strings.Contains(err.Message, "failed") {
		t.Errorf("expected default message to contain 'validation' and 'failed', got '%s'", err.Message)
	}
}

func TestValidationError_CustomMessage(t *testing.T) {
	statusCode := 422
	response := map[string]interface{}{
		"message": "user_id is required",
	}
	err := NewValidationError(&statusCode, response)

	if err.Message != "user_id is required" {
		t.Errorf("expected custom message 'user_id is required', got '%s'", err.Message)
	}
}

func TestValidationError_TypeAssertion(t *testing.T) {
	statusCode := 422
	err := NewValidationError(&statusCode, nil)

	var genericErr error = err
	if _, ok := genericErr.(*ValidationError); !ok {
		t.Error("expected error to be *ValidationError")
	}
}

// TestErrorHierarchy tests error hierarchy.
func TestErrorHierarchy(t *testing.T) {
	statusCode := 401

	tests := []struct {
		name     string
		err      error
		isClient bool
	}{
		{"AuthenticationError", NewAuthenticationError(&statusCode, nil), true},
		{"RateLimitError", NewRateLimitError("rate limit", nil, &statusCode, nil), true},
		{"NotFoundError", NewNotFoundError("/path", &statusCode, nil), true},
		{"ValidationError", NewValidationError(&statusCode, nil), true},
		{"ClientError", NewClientError("error", &statusCode, nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// All errors should implement error interface
			if tt.err.Error() == "" {
				t.Error("expected non-empty error string")
			}
		})
	}
}
