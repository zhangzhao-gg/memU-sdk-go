package memu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
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
	apiKey     string
	baseURL    string
	httpClient *http.Client
	maxRetries int
	timeout    time.Duration
}

// NewClient creates a new MemU API client.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	client := &Client{
		apiKey:     strings.TrimSpace(apiKey),
		baseURL:    strings.TrimRight(DefaultBaseURL, "/"),
		maxRetries: DefaultMaxRetries,
		timeout:    DefaultTimeout,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
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
func (c *Client) defaultHeaders() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
		"Content-Type":  "application/json",
		"User-Agent":    "memu-go-sdk/1.0.0",
	}
}

// request makes an HTTP request to the API with automatic retry logic.
func (c *Client) request(ctx context.Context, method, path string, body interface{}, params map[string]string) (map[string]interface{}, error) {
	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
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
			lastErr = err
			if attempt < c.maxRetries-1 {
				waitTime := time.Duration(math.Pow(2, float64(attempt))) * time.Second
				time.Sleep(waitTime)
				continue
			}
			return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries, err)
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
				waitTime = time.Duration(math.Pow(2, float64(attempt))) * time.Second
			}

			if attempt < c.maxRetries-1 {
				time.Sleep(waitTime)
				continue
			}

			retryAfterFloat := float64(waitTime) / float64(time.Second)
			statusCode := resp.StatusCode
			return nil, NewRateLimitError("rate limit exceeded", &retryAfterFloat, &statusCode, result)
		}

		// Handle server errors (5xx) - retry
		if resp.StatusCode >= 500 {
			if attempt < c.maxRetries-1 {
				waitTime := time.Duration(math.Pow(2, float64(attempt))) * time.Second
				time.Sleep(waitTime)
				continue
			}
			statusCode := resp.StatusCode
			return nil, NewClientError(fmt.Sprintf("server error: %d", resp.StatusCode), &statusCode, result)
		}

		// Handle client errors (4xx) - don't retry
		if resp.StatusCode >= 400 {
			return nil, c.raiseForStatus(resp.StatusCode, path, result)
		}

		// Success
		return result, nil
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries, lastErr)
}

// raiseForStatus raises an appropriate error for HTTP error status codes.
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
		return nil, fmt.Errorf("request is required")
	}

	if len(req.Conversation) == 0 && req.ConversationText == nil {
		return nil, fmt.Errorf("either Conversation or ConversationText must be provided")
	}

	if req.UserID == "" {
		return nil, fmt.Errorf("UserID is required")
	}

	if req.AgentID == "" {
		return nil, fmt.Errorf("AgentID is required")
	}

	// Build request payload
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

	// Wait for completion if requested
	if req.WaitForCompletion && result.TaskID != nil {
		pollInterval := req.PollInterval
		if pollInterval == 0 {
			pollInterval = DefaultPollInterval
		}

		timeout := req.Timeout
		if timeout == 0 {
			timeout = DefaultWaitTimeout
		}

		// Create a context with timeout
		waitCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-waitCtx.Done():
				return nil, fmt.Errorf("memorization task timed out after %v", timeout)
			case <-ticker.C:
				status, err := c.GetTaskStatus(ctx, *result.TaskID)
				if err != nil {
					return nil, err
				}

				if status.Status == TaskStatusCompleted || status.Status == TaskStatusSuccess {
					// Parse result
					if status.Result != nil {
						if resource, ok := status.Result["resource"].(map[string]interface{}); ok {
							resourceData, _ := json.Marshal(resource)
							var r MemoryResource
							if err := json.Unmarshal(resourceData, &r); err == nil {
								result.Resource = &r
							}
						}

						if items, ok := status.Result["items"].([]interface{}); ok {
							for _, item := range items {
								itemData, _ := json.Marshal(item)
								var i MemoryItem
								if err := json.Unmarshal(itemData, &i); err == nil {
									result.Items = append(result.Items, &i)
								}
							}
						}

						if categories, ok := status.Result["categories"].([]interface{}); ok {
							for _, cat := range categories {
								catData, _ := json.Marshal(cat)
								var c MemoryCategory
								if err := json.Unmarshal(catData, &c); err == nil {
									result.Categories = append(result.Categories, &c)
								}
							}
						}
					}
					return result, nil
				}

				if status.Status == TaskStatusFailed {
					message := "memorization task failed"
					if status.Message != nil {
						message = *status.Message
					}
					return nil, fmt.Errorf(message)
				}
			}
		}
	}

	// Parse immediate result
	if resource, ok := response["resource"].(map[string]interface{}); ok {
		resourceData, _ := json.Marshal(resource)
		var r MemoryResource
		if err := json.Unmarshal(resourceData, &r); err == nil {
			result.Resource = &r
		}
	}

	if items, ok := response["items"].([]interface{}); ok {
		for _, item := range items {
			itemData, _ := json.Marshal(item)
			var i MemoryItem
			if err := json.Unmarshal(itemData, &i); err == nil {
				result.Items = append(result.Items, &i)
			}
		}
	}

	if categories, ok := response["categories"].([]interface{}); ok {
		for _, cat := range categories {
			catData, _ := json.Marshal(cat)
			var c MemoryCategory
			if err := json.Unmarshal(catData, &c); err == nil {
				result.Categories = append(result.Categories, &c)
			}
		}
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

	// Parse response
	responseData, _ := json.Marshal(response)
	var status TaskStatus
	if err := json.Unmarshal(responseData, &status); err != nil {
		return nil, fmt.Errorf("failed to parse task status: %w", err)
	}

	return &status, nil
}

// Retrieve retrieves relevant memories based on a query.
func (c *Client) Retrieve(ctx context.Context, req *RetrieveRequest) (*RetrieveResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	if req.Query == nil {
		return nil, fmt.Errorf("Query is required")
	}

	if req.UserID == "" {
		return nil, fmt.Errorf("UserID is required")
	}

	if req.AgentID == "" {
		return nil, fmt.Errorf("AgentID is required")
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
		for _, cat := range categories {
			catData, _ := json.Marshal(cat)
			var c MemoryCategory
			if err := json.Unmarshal(catData, &c); err == nil {
				result.Categories = append(result.Categories, &c)
			}
		}
	}

	if items, ok := response["items"].([]interface{}); ok {
		for _, item := range items {
			itemData, _ := json.Marshal(item)
			var i MemoryItem
			if err := json.Unmarshal(itemData, &i); err == nil {
				result.Items = append(result.Items, &i)
			}
		}
	}

	if resources, ok := response["resources"].([]interface{}); ok {
		for _, res := range resources {
			resData, _ := json.Marshal(res)
			var r MemoryResource
			if err := json.Unmarshal(resData, &r); err == nil {
				result.Resources = append(result.Resources, &r)
			}
		}
	}

	if nextStepQuery, ok := response["next_step_query"].(string); ok {
		result.NextStepQuery = &nextStepQuery
	}

	return result, nil
}

// ListCategories lists all memory categories.
func (c *Client) ListCategories(ctx context.Context, req *ListCategoriesRequest) ([]*MemoryCategory, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	if req.UserID == "" {
		return nil, fmt.Errorf("UserID is required")
	}

	// Build request payload
	payload := map[string]interface{}{
		"user_id": req.UserID,
	}

	if req.AgentID != nil {
		payload["agent_id"] = *req.AgentID
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
		for _, cat := range categoriesList {
			catData, _ := json.Marshal(cat)
			var c MemoryCategory
			if err := json.Unmarshal(catData, &c); err == nil {
				categories = append(categories, &c)
			}
		}
	}

	return categories, nil
}
