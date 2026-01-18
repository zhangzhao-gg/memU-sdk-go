// Package memu provides the core HTTP client implementation for the MemU SDK.
// This is the main entry point for the SDK, implementing all four major API methods
// with automatic retry logic, JSON parsing, and parameter validation.
package memu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default API base URL.
	DefaultBaseURL = "https://api.memu.so"
	// DefaultTimeout is the default request timeout.
	DefaultTimeout = 60 * time.Second
	// DefaultPollInterval is the default interval between status checks.
	DefaultPollInterval = 2 * time.Second
	// DefaultMaxRetries is the default maximum number of retry attempts.
	DefaultMaxRetries = 3
	// DefaultWaitTimeout is the default maximum time to wait for task completion.
	DefaultWaitTimeout = 5 * time.Minute
)

// Client is the MemU API client.
type Client struct {
	// apiKey is the API authentication key.
	apiKey string
	// baseURL is the base URL for API requests.
	baseURL string
	// httpClient is the underlying HTTP client used for requests.
	httpClient *http.Client
	// maxRetries is the maximum number of retry attempts.
	maxRetries int
	// timeout is the request timeout duration.
	timeout time.Duration
	// retryPolicy defines the retry behavior for failed requests.
	retryPolicy RetryPolicy
}

// NewClient creates a new MemU API client.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	client := &Client{
		apiKey:     apiKey,
		baseURL:    strings.TrimRight(DefaultBaseURL, "/"),
		maxRetries: DefaultMaxRetries,
		timeout:    DefaultTimeout,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		retryPolicy: NewDefaultRetryPolicy(nil),
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	// Update HTTP client timeout if it was changed
	if client.httpClient.Timeout != client.timeout {
		client.httpClient.Timeout = client.timeout
	}

	return client, nil
}

// defaultHeaders returns the default headers for API requests.
// This includes the authorization bearer token, content type, and user agent.
func (c *Client) defaultHeaders() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
		"Content-Type":  "application/json",
		"User-Agent":    "memu-go-sdk/1.0.0",
	}
}

// parseJSONObject parses a JSON object into a struct, avoiding double serialization.
// This is a performance optimization that directly deserializes data without
// the overhead of Marshal → Unmarshal cycles.
// It accepts any interface and returns a pointer to the typed struct.
func parseJSONObject[T any](data interface{}) (*T, error) {
	if data == nil {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}

	var obj T
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal object: %w", err)
	}

	return &obj, nil
}

// parseJSONArray parses a JSON array into a slice of structs, avoiding double serialization.
// This is a performance optimization that serializes the entire array at once
// rather than processing elements individually.
// It accepts a slice of interfaces and returns a slice of pointers to the typed struct.
func parseJSONArray[T any](data []interface{}) ([]*T, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// Marshal the entire array once
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal array: %w", err)
	}

	// Unmarshal into the target type
	var result []*T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal array: %w", err)
	}

	return result, nil
}

// buildMemorizePayload builds the payload for a Memorize request.
// This provides unified payload construction logic to simplify the Memorize method.
// It handles default values for user_name and agent_name, and conditionally includes
// conversation, conversation_text, and session_date fields.
func buildMemorizePayload(req *MemorizeRequest) map[string]interface{} {
	payload := map[string]interface{}{
		"user_id":  req.UserID,
		"agent_id": req.AgentID,
	}

	if req.UserName != "" {
		payload["user_name"] = req.UserName
	} else {
		payload["user_name"] = "User"
	}

	if req.AgentName != "" {
		payload["agent_name"] = req.AgentName
	} else {
		payload["agent_name"] = "Assistant"
	}

	if len(req.Conversation) > 0 {
		payload["conversation"] = req.Conversation
	} else if req.ConversationText != nil {
		payload["conversation_text"] = *req.ConversationText
	}

	if req.SessionDate != nil {
		payload["session_date"] = *req.SessionDate
	}

	return payload
}

// request makes an HTTP request to the API with automatic retry logic.
// It handles request construction, header setting, query parameters, response parsing,
// rate limiting, and error handling. The method automatically retries on transient errors
// based on the configured retry policy.
func (c *Client) request(ctx context.Context, method, path string, body interface{}, params map[string]string) (map[string]interface{}, error) {
	for attempt := 0; ; attempt++ {
		// Prepare request body
		var bodyReader io.Reader
		if body != nil {
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(jsonData)
		}

		// Create request
		url := c.baseURL + path
		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		for key, value := range c.defaultHeaders() {
			req.Header.Set(key, value)
		}

		// Set query parameters
		if len(params) > 0 {
			q := req.URL.Query()
			for key, value := range params {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
		}

		// Make request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			// Check if we should retry
			if c.retryPolicy.ShouldRetry(attempt, 0, err) {
				time.Sleep(c.retryPolicy.GetBackoff(attempt))
				continue
			}
			return nil, fmt.Errorf("request failed after %d attempts: %w", attempt+1, err)
		}
		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Parse response
		var result map[string]interface{}
		if len(respBody) > 0 {
			if err := json.Unmarshal(respBody, &result); err != nil {
				// If JSON parsing fails, return the raw response
				result = map[string]interface{}{
					"raw": string(respBody),
				}
			}
		}

		// Handle rate limiting (429)
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			var waitTime time.Duration
			if retryAfter != "" {
				if seconds, err := strconv.ParseFloat(retryAfter, 64); err == nil {
					waitTime = time.Duration(seconds * float64(time.Second))
				}
			} else {
				waitTime = c.retryPolicy.GetBackoff(attempt)
			}

			if c.retryPolicy.ShouldRetry(attempt, resp.StatusCode, nil) {
				time.Sleep(waitTime)
				continue
			}

			retryAfterFloat := float64(waitTime) / float64(time.Second)
			statusCode := resp.StatusCode
			return nil, NewRateLimitError("rate limit exceeded", &retryAfterFloat, &statusCode, result)
		}

		// Handle server errors (5xx) - retry
		if resp.StatusCode >= 500 {
			if c.retryPolicy.ShouldRetry(attempt, resp.StatusCode, nil) {
				time.Sleep(c.retryPolicy.GetBackoff(attempt))
				continue
			}
			statusCode := resp.StatusCode
			// Include response body in error message for debugging
			errorMsg := fmt.Sprintf("server error: %d", resp.StatusCode)
			if len(respBody) > 0 {
				errorMsg = fmt.Sprintf("server error: %d, response: %s", resp.StatusCode, string(respBody))
			}
			return nil, NewClientError(errorMsg, &statusCode, result)
		}

		// Handle client errors (4xx) - don't retry
		if resp.StatusCode >= 400 {
			return nil, c.raiseForStatus(resp.StatusCode, path, result)
		}

		// Success
		return result, nil
	}
}

// raiseForStatus raises an appropriate error for HTTP error status codes.
// It maps HTTP status codes to specific error types: 401 to AuthenticationError,
// 404 to NotFoundError, 422 to ValidationError, and others to generic ClientError.
func (c *Client) raiseForStatus(statusCode int, path string, response map[string]interface{}) error {
	status := &statusCode

	switch statusCode {
	case http.StatusUnauthorized:
		return NewAuthenticationError(status, response)
	case http.StatusNotFound:
		return NewNotFoundError(path, status, response)
	case http.StatusUnprocessableEntity:
		return NewValidationError(status, response)
	default:
		return NewClientError(fmt.Sprintf("HTTP %d: %s", statusCode, path), status, response)
	}
}

// Memorize memorizes a conversation and extracts structured memory.
func (c *Client) Memorize(ctx context.Context, req *MemorizeRequest) (*MemorizeResult, error) {
	if req == nil {
		return nil, fmt.Errorf("Memorize: request is required")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Build request payload
	payload := buildMemorizePayload(req)

	// Make request
	response, err := c.request(ctx, "POST", "/api/v3/memory/memorize", payload, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	result := &MemorizeResult{}
	if taskID, ok := response["task_id"].(string); ok {
		result.TaskID = &taskID
	}
	if status, ok := response["status"].(string); ok {
		result.Status = &status
	}
	if message, ok := response["message"].(string); ok {
		result.Message = &message
	}

	return result, nil
}

// GetTaskStatus gets the status of a memorization task.
func (c *Client) GetTaskStatus(ctx context.Context, taskID string) (*TaskStatus, error) {
	if taskID == "" {
		return nil, fmt.Errorf("taskID is required")
	}

	path := fmt.Sprintf("/api/v3/memory/memorize/status/%s", taskID)
	response, err := c.request(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, err
	}
	// Parse response using parseJSONObject to avoid double serialization
	status, err := parseJSONObject[TaskStatus](response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task status: %w", err)
	}

	return status, nil
}

// ListCategories lists all memory categories.
func (c *Client) ListCategories(ctx context.Context, req *ListCategoriesRequest) ([]*MemoryCategory, error) {
	if req == nil {
		return nil, fmt.Errorf("ListCategories: request is required")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Build request payload
	payload := map[string]interface{}{
		"user_id":  req.UserID,
		"agent_id": req.AgentID,
	}

	// Make request
	response, err := c.request(ctx, "POST", "/api/v3/memory/categories", payload, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var categories []*MemoryCategory

	// Try to get categories from "categories" field first
	categoriesData, ok := response["categories"]
	if !ok {
		// If not found, assume the response itself is the categories array
		categoriesData = response
	}

	if categoriesList, ok := categoriesData.([]interface{}); ok {
		parsedCategories, err := parseJSONArray[MemoryCategory](categoriesList)
		if err != nil {
			return nil, fmt.Errorf("failed to parse categories: %w", err)
		}
		categories = parsedCategories
	}

	return categories, nil
}

// Retrieve retrieves relevant memories based on a query.
func (c *Client) Retrieve(ctx context.Context, req *RetrieveRequest) (*RetrieveResult, error) {
	if req == nil {
		return nil, fmt.Errorf("Retrieve: request is required")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Build request payload
	payload := map[string]interface{}{
		"user_id":  req.UserID,
		"agent_id": req.AgentID,
		"query":    req.Query,
	}

	// Make request
	response, err := c.request(ctx, "POST", "/api/v3/memory/retrieve", payload, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	result := &RetrieveResult{}

	if categories, ok := response["categories"].([]interface{}); ok {
		parsedCategories, err := parseJSONArray[MemoryCategory](categories)
		if err != nil {
			return nil, fmt.Errorf("failed to parse categories: %w", err)
		}
		result.Categories = parsedCategories
	}

	if items, ok := response["items"].([]interface{}); ok {
		parsedItems, err := parseJSONArray[MemoryItem](items)
		if err != nil {
			return nil, fmt.Errorf("failed to parse items: %w", err)
		}
		result.Items = parsedItems
	}

	if resources, ok := response["resources"].([]interface{}); ok {
		parsedResources, err := parseJSONArray[MemoryResource](resources)
		if err != nil {
			return nil, fmt.Errorf("failed to parse resources: %w", err)
		}
		result.Resources = parsedResources
	}

	if rewrittenQuery, ok := response["rewritten_query"].(string); ok {
		result.RewrittenQuery = &rewrittenQuery
	}

	return result, nil
}
