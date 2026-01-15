/**
 * [INPUT]: 依赖 testing, net/http/httptest, context, time; 依赖 retry.go 的 RetryPolicy 相关类型
 * [OUTPUT]: 对 retry.go 的测试覆盖
 * [POS]: SDK 根目录的测试层，测试重试策略功能
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */
package memu

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestDefaultRetryPolicy_ShouldRetry(t *testing.T) {
	policy := NewDefaultRetryPolicy(nil)

	tests := []struct {
		name       string
		attempt    int
		statusCode int
		err        error
		want       bool
	}{
		{
			name:       "retry on network error",
			attempt:    0,
			statusCode: 0,
			err:        errors.New("network error"),
			want:       true,
		},
		{
			name:       "retry on 429",
			attempt:    0,
			statusCode: http.StatusTooManyRequests,
			err:        nil,
			want:       true,
		},
		{
			name:       "retry on 500",
			attempt:    0,
			statusCode: http.StatusInternalServerError,
			err:        nil,
			want:       true,
		},
		{
			name:       "retry on 502",
			attempt:    0,
			statusCode: http.StatusBadGateway,
			err:        nil,
			want:       true,
		},
		{
			name:       "retry on 503",
			attempt:    0,
			statusCode: http.StatusServiceUnavailable,
			err:        nil,
			want:       true,
		},
		{
			name:       "retry on 504",
			attempt:    0,
			statusCode: http.StatusGatewayTimeout,
			err:        nil,
			want:       true,
		},
		{
			name:       "don't retry on 400",
			attempt:    0,
			statusCode: http.StatusBadRequest,
			err:        nil,
			want:       false,
		},
		{
			name:       "don't retry after max attempts",
			attempt:    3,
			statusCode: http.StatusInternalServerError,
			err:        nil,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.ShouldRetry(tt.attempt, tt.statusCode, tt.err)
			if got != tt.want {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultRetryPolicy_GetBackoff(t *testing.T) {
	policy := NewDefaultRetryPolicy(nil)

	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{
			name:    "first attempt",
			attempt: 0,
			want:    1 * time.Second,
		},
		{
			name:    "second attempt",
			attempt: 1,
			want:    2 * time.Second,
		},
		{
			name:    "third attempt",
			attempt: 2,
			want:    4 * time.Second,
		},
		{
			name:    "fourth attempt",
			attempt: 3,
			want:    8 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.GetBackoff(tt.attempt)
			if got != tt.want {
				t.Errorf("GetBackoff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultRetryPolicy_MaxDelay(t *testing.T) {
	config := &RetryConfig{
		MaxRetries: 10,
		BaseDelay:  1 * time.Second,
		MaxDelay:   10 * time.Second,
		RetryableStatusCodes: map[int]bool{
			http.StatusInternalServerError: true,
		},
	}
	policy := NewDefaultRetryPolicy(config)

	// Attempt 10 should give 1024 seconds, but capped at 10
	backoff := policy.GetBackoff(10)
	if backoff != 10*time.Second {
		t.Errorf("GetBackoff(10) = %v, want %v (should be capped)", backoff, 10*time.Second)
	}
}

func TestNoRetryPolicy(t *testing.T) {
	policy := NewNoRetryPolicy()

	// Should never retry
	if policy.ShouldRetry(0, http.StatusInternalServerError, nil) {
		t.Error("NoRetryPolicy should never retry")
	}

	if policy.ShouldRetry(0, 0, errors.New("error")) {
		t.Error("NoRetryPolicy should never retry on error")
	}

	// Should always return 0 backoff
	if backoff := policy.GetBackoff(0); backoff != 0 {
		t.Errorf("NoRetryPolicy.GetBackoff() = %v, want 0", backoff)
	}
}

func TestCustomRetryPolicy(t *testing.T) {
	shouldRetryFunc := func(attempt int, statusCode int, err error) bool {
		// Only retry on 503
		return statusCode == http.StatusServiceUnavailable
	}

	getBackoffFunc := func(attempt int) time.Duration {
		// Fixed 5 second backoff
		return 5 * time.Second
	}

	policy := NewCustomRetryPolicy(5, shouldRetryFunc, getBackoffFunc)

	// Should retry on 503
	if !policy.ShouldRetry(0, http.StatusServiceUnavailable, nil) {
		t.Error("CustomRetryPolicy should retry on 503")
	}

	// Should not retry on 500
	if policy.ShouldRetry(0, http.StatusInternalServerError, nil) {
		t.Error("CustomRetryPolicy should not retry on 500")
	}

	// Should not retry after max attempts
	if policy.ShouldRetry(5, http.StatusServiceUnavailable, nil) {
		t.Error("CustomRetryPolicy should not retry after max attempts")
	}

	// Should return fixed backoff
	if backoff := policy.GetBackoff(0); backoff != 5*time.Second {
		t.Errorf("CustomRetryPolicy.GetBackoff() = %v, want 5s", backoff)
	}
}

func TestWithRetryPolicy(t *testing.T) {
	customPolicy := NewNoRetryPolicy()
	client, err := NewClient("test-key", WithRetryPolicy(customPolicy))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.retryPolicy != customPolicy {
		t.Error("WithRetryPolicy did not set custom policy")
	}
}

func TestWithHTTPClient(t *testing.T) {
	customHTTPClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	client, err := NewClient("test-key", WithHTTPClient(customHTTPClient))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.httpClient != customHTTPClient {
		t.Error("WithHTTPClient did not set custom HTTP client")
	}
}
