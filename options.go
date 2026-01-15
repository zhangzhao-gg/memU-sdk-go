/**
 * [INPUT]: 依赖 net/http 的 Client, time 的 Duration; 依赖 retry.go 的 RetryPolicy
 * [OUTPUT]: 对外提供 Option 类型, WithBaseURL, WithTimeout, WithMaxRetries, WithHTTPClient, WithRetryPolicy 函数
 * [POS]: SDK 根目录的配置层，被 client.go 的 NewClient 消费，实现选项模式
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */
package memu

import (
	"net/http"
	"time"
)

// Option is a function that configures a Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithTimeout sets the request timeout duration.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retry attempts for failed requests.
func WithMaxRetries(retries int) Option {
	return func(c *Client) {
		c.maxRetries = retries
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithRetryPolicy sets a custom retry policy.
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(c *Client) {
		c.retryPolicy = policy
	}
}
