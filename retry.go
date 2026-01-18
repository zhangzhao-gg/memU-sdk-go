// Package memu provides retry policy implementations for the MemU SDK.
// This file defines retry behavior configuration and various retry strategies.
package memu

import (
	"net/http"
	"time"
)

// RetryPolicy defines the interface for retry behavior.
type RetryPolicy interface {
	// ShouldRetry determines if a request should be retried based on the attempt number and error.
	ShouldRetry(attempt int, statusCode int, err error) bool

	// GetBackoff returns the backoff duration for a given attempt.
	GetBackoff(attempt int) time.Duration
}

// RetryConfig holds retry configuration.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts.
	MaxRetries int

	// BaseDelay is the base delay for exponential backoff.
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries.
	MaxDelay time.Duration

	// RetryableStatusCodes are HTTP status codes that should trigger a retry.
	RetryableStatusCodes map[int]bool
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries: 3,
		BaseDelay:  1 * time.Second,
		MaxDelay:   32 * time.Second,
		RetryableStatusCodes: map[int]bool{
			http.StatusTooManyRequests:     true, // 429
			http.StatusInternalServerError: true, // 500
			http.StatusBadGateway:          true, // 502
			http.StatusServiceUnavailable:  true, // 503
			http.StatusGatewayTimeout:      true, // 504
		},
	}
}

// defaultRetryPolicy is the default retry policy implementation.
type defaultRetryPolicy struct {
	// config holds the retry configuration.
	config *RetryConfig
}

// NewDefaultRetryPolicy creates a new default retry policy.
func NewDefaultRetryPolicy(config *RetryConfig) RetryPolicy {
	if config == nil {
		config = DefaultRetryConfig()
	}
	return &defaultRetryPolicy{config: config}
}

// ShouldRetry implements RetryPolicy.
func (p *defaultRetryPolicy) ShouldRetry(attempt int, statusCode int, err error) bool {
	// Don't retry if we've exceeded max retries
	if attempt >= p.config.MaxRetries {
		return false
	}

	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on specific status codes
	if statusCode > 0 {
		return p.config.RetryableStatusCodes[statusCode]
	}

	return false
}

// GetBackoff implements RetryPolicy.
func (p *defaultRetryPolicy) GetBackoff(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	backoff := p.config.BaseDelay * (1 << uint(attempt))

	// Cap at max delay
	if backoff > p.config.MaxDelay {
		backoff = p.config.MaxDelay
	}

	return backoff
}

// noRetryPolicy never retries.
type noRetryPolicy struct{}

// NewNoRetryPolicy creates a policy that never retries.
func NewNoRetryPolicy() RetryPolicy {
	return &noRetryPolicy{}
}

// ShouldRetry implements RetryPolicy.
func (p *noRetryPolicy) ShouldRetry(attempt int, statusCode int, err error) bool {
	return false
}

// GetBackoff implements RetryPolicy.
func (p *noRetryPolicy) GetBackoff(attempt int) time.Duration {
	return 0
}

// CustomRetryFunc is a function type for custom retry logic.
type CustomRetryFunc func(attempt int, statusCode int, err error) bool

// CustomBackoffFunc is a function type for custom backoff logic.
type CustomBackoffFunc func(attempt int) time.Duration

// customRetryPolicy allows custom retry logic.
type customRetryPolicy struct {
	// shouldRetry is the custom function to determine if a retry should occur.
	shouldRetry CustomRetryFunc
	// getBackoff is the custom function to calculate backoff duration.
	getBackoff CustomBackoffFunc
	// maxRetries is the maximum number of retry attempts.
	maxRetries int
}

// NewCustomRetryPolicy creates a custom retry policy.
func NewCustomRetryPolicy(maxRetries int, shouldRetry CustomRetryFunc, getBackoff CustomBackoffFunc) RetryPolicy {
	return &customRetryPolicy{
		shouldRetry: shouldRetry,
		getBackoff:  getBackoff,
		maxRetries:  maxRetries,
	}
}

// ShouldRetry implements RetryPolicy.
func (p *customRetryPolicy) ShouldRetry(attempt int, statusCode int, err error) bool {
	if attempt >= p.maxRetries {
		return false
	}
	return p.shouldRetry(attempt, statusCode, err)
}

// GetBackoff implements RetryPolicy.
func (p *customRetryPolicy) GetBackoff(attempt int) time.Duration {
	return p.getBackoff(attempt)
}
