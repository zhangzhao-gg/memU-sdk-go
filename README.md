# MemU Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/NevaMind-AI/memU-sdk-go.svg)](https://pkg.go.dev/github.com/NevaMind-AI/memU-sdk-go)
[![Go 1.21+](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for the [MemU Cloud API](https://api.memu.so) - Manage structured, long-term memory for AI agents.

## Features

- Full Cloud API v3 Coverage - memorize, retrieve, categories, task status
- Context Support - Native Go context for timeout and cancellation
- Automatic Retry - Exponential backoff with configurable retry policies
- Rate Limit Handling - Respects Retry-After headers
- Type Safety - Strongly typed models with full documentation
- Custom Errors - Specific error types for different failure cases
- Zero Dependencies - Uses only Go standard library

## Installation

```bash
go get github.com/NevaMind-AI/memU-sdk-go
```

## Quick Start

### Get Your API Key

1. Sign up at [memu.so](https://memu.so)
2. Navigate to your dashboard to obtain your API key

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    memu "github.com/NevaMind-AI/memU-sdk-go"
)

func main() {
    // Initialize the client
    client, err := memu.NewClient("your_api_key")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Memorize a conversation
    result, err := client.Memorize(ctx, &memu.MemorizeRequest{
        Conversation: []memu.ConversationMessage{
            {Role: "user", Content: "I love Italian food, especially pasta."},
            {Role: "assistant", Content: "That's great! What's your favorite dish?"},
            {Role: "user", Content: "Carbonara is my absolute favorite!"},
        },
        UserID:    "user_123",
        AgentID:   "my_assistant",
        UserName:  "User",
        AgentName: "Assistant",
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Task ID: %s\n", *result.TaskID)
    fmt.Printf("Status: %s\n", *result.Status)

    // Poll for task completion
    for {
        status, err := client.GetTaskStatus(ctx, *result.TaskID)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Status: %s\n", status.Status)

        if status.Status == memu.TaskStatusSuccess || status.Status == memu.TaskStatusCompleted {
            fmt.Println("Task completed successfully!")
            break
        } else if status.Status == memu.TaskStatusFailed {
            fmt.Printf("Task failed: %s\n", status.Message)
            break
        }

        time.Sleep(2 * time.Second)
    }

    // Retrieve memories
    memories, err := client.Retrieve(ctx, &memu.RetrieveRequest{
        Query:   "What food does the user like?",
        UserID:  "user_123",
        AgentID: "my_assistant",
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d relevant memories\n", len(memories.Items))
    for _, item := range memories.Items {
        if item.Content != nil {
            fmt.Printf("  - [%s] %s\n", *item.MemoryType, *item.Content)
        }
    }
}
```

## API Reference

### Client

#### NewClient

```go
func NewClient(apiKey string, opts ...Option) (*Client, error)
```

Creates a new MemU API client.

**Parameters:**
- `apiKey`: Your MemU API key (required)
- `opts`: Optional configuration options

**Options:**
- `WithBaseURL(url string)` - Set custom base URL (default: https://api.memu.so)
- `WithTimeout(timeout time.Duration)` - Set request timeout (default: 60s)
- `WithMaxRetries(retries int)` - Set max retry attempts (default: 3)
- `WithHTTPClient(client *http.Client)` - Use custom HTTP client
- `WithRetryPolicy(policy RetryPolicy)` - Set custom retry policy

**Example:**
```go
client, err := memu.NewClient(
    "your_api_key",
    memu.WithTimeout(30 * time.Second),
    memu.WithMaxRetries(5),
)
```

### Methods

#### Memorize

Memorize a conversation and extract structured memory. This is an **asynchronous operation** that returns immediately with a task ID.

```go
func (c *Client) Memorize(ctx context.Context, req *MemorizeRequest) (*MemorizeResult, error)
```

**Request Fields:**
- `Conversation` - List of conversation messages (optional if ConversationText is provided)
  - **Minimum 3 messages required**
  - Each message includes:
    - `Role` - "user" or "assistant" (required)
    - `Content` - Message content (required)
    - `Name` - Speaker name (optional)
    - `CreatedAt` - Timestamp in ISO format (optional)
- `ConversationText` - Alternative: raw conversation text (optional if Conversation is provided)
- `UserID` - User ID for scoping the memory (required)
- `AgentID` - Agent ID for scoping the memory (required)
- `UserName` - Display name for the user (default: "User")
- `AgentName` - Display name for the agent (default: "Assistant")
- `SessionDate` - Optional session date in ISO format

**Response Fields:**
- `TaskID` - Task ID for async tracking
- `Status` - Task status (typically "PENDING")
- `Message` - Descriptive message

**Example:**
```go
// Basic example
result, err := client.Memorize(ctx, &memu.MemorizeRequest{
    Conversation: []memu.ConversationMessage{
        {Role: "user", Content: "I love pasta"},
        {Role: "assistant", Content: "Great choice!"},
        {Role: "user", Content: "Especially carbonara"},
    },
    UserID:  "user_123",
    AgentID: "agent_456",
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Task ID: %s\n", *result.TaskID)
fmt.Printf("Status: %s\n", *result.Status)

// Poll for completion
status, err := client.GetTaskStatus(ctx, *result.TaskID)

// Advanced example with full metadata
name1 := "John"
name2 := "Coach"
time1 := "2024-01-15T10:30:00Z"
time2 := "2024-01-15T10:30:15Z"
time3 := "2024-01-15T10:31:00Z"
sessionDate := "2024-01-15T10:30:00Z"

result, err := client.Memorize(ctx, &memu.MemorizeRequest{
    Conversation: []memu.ConversationMessage{
        {
            Role:      "user",
            Content:   "I love playing tennis on weekends",
            Name:      &name1,
            CreatedAt: &time1,
        },
        {
            Role:      "assistant",
            Content:   "That's great! Tennis is excellent exercise.",
            Name:      &name2,
            CreatedAt: &time2,
        },
        {
            Role:      "user",
            Content:   "I play every Saturday morning",
            Name:      &name1,
            CreatedAt: &time3,
        },
    },
    UserID:      "user_123",
    AgentID:     "agent_456",
    UserName:    "John Doe",
    AgentName:   "Tennis Coach AI",
    SessionDate: &sessionDate,
})
```


#### Retrieve

Retrieve relevant memories based on a query.

```go
func (c *Client) Retrieve(ctx context.Context, req *RetrieveRequest) (*RetrieveResult, error)
```

**Request Fields:**
- `Query` - Query string or list of conversation messages (required)
- `UserID` - User ID for scoping (required)
- `AgentID` - Agent ID for scoping (required)

**Example:**
```go
// Simple text query
memories, err := client.Retrieve(ctx, &memu.RetrieveRequest{
    Query:   "What are the user's food preferences?",
    UserID:  "user_123",
    AgentID: "agent_456",
})

// Conversation-aware query
memories, err := client.Retrieve(ctx, &memu.RetrieveRequest{
    Query: []memu.ConversationMessage{
        {Role: "user", Content: "What do they like?"},
        {Role: "assistant", Content: "They have several preferences."},
        {Role: "user", Content: "Tell me about food specifically"},
    },
    UserID:  "user_123",
    AgentID: "agent_456",
})
```

#### ListCategories

List all memory categories for a user.

```go
func (c *Client) ListCategories(ctx context.Context, req *ListCategoriesRequest) ([]*MemoryCategory, error)
```

**Request Fields:**
- `UserID` - User ID for scoping (required)
- `AgentID` - Agent ID for scoping (optional)

**Example:**
```go
categories, err := client.ListCategories(ctx, &memu.ListCategoriesRequest{
    UserID: "user_123",
})

for _, cat := range categories {
    if cat.Name != nil {
        fmt.Printf("%s", *cat.Name)
    }
    if cat.Summary != nil {
        fmt.Printf(": %s\n", *cat.Summary)
    }
}
```

#### GetTaskStatus

Get the status of an asynchronous memorization task.

```go
func (c *Client) GetTaskStatus(ctx context.Context, taskID string) (*TaskStatus, error)
```

**Example:**
```go
status, err := client.GetTaskStatus(ctx, "task_abc123")
if status.Status == memu.TaskStatusCompleted {
    fmt.Printf("Task completed: %s\n", status.Message)
}
```

## Data Models

### MemorizeResult

```go
type MemorizeResult struct {
    TaskID  *string // Task ID for async tracking
    Status  *string // Task status (typically "PENDING")
    Message *string // Descriptive message
}
```

### RetrieveResult

```go
type RetrieveResult struct {
    RewrittenQuery *string           // Rewritten query (if applicable)
    Categories     []*MemoryCategory // Relevant categories
    Items          []*MemoryItem     // Relevant memory items
    Resources      []*MemoryResource // Related raw resources
}
```

### MemoryItem

```go
type MemoryItem struct {
    Content    *string // Content text
    MemoryType *string // Type: profile, event, preference, etc.
}
```

### MemoryCategory

```go
type MemoryCategory struct {
    Name        *string // Category name (e.g., 'personal info')
    Description *string // Description
    Summary     *string // Summary of content
    UserID      *string // User ID
    AgentID     *string // Agent ID
}
```

### MemoryResource

```go
type MemoryResource struct {
    Modality    *string                // Resource modality
    ResourceURL *string                // Resource URL
    Caption     *string                // Caption
    Content     map[string]interface{} // Content data
    Metadata    map[string]interface{} // Metadata
}
```

### TaskStatus

```go
type TaskStatus struct {
    TaskID     string         // Task identifier
    Status     TaskStatusEnum // PENDING, PROCESSING, COMPLETED, SUCCESS, FAILED
    Message    string         // Status message or error
    DetailInfo string         // Detailed information
}
```

### TaskStatusEnum

```go
const (
    TaskStatusPending    TaskStatusEnum = "PENDING"
    TaskStatusProcessing TaskStatusEnum = "PROCESSING"
    TaskStatusCompleted  TaskStatusEnum = "COMPLETED"
    TaskStatusSuccess    TaskStatusEnum = "SUCCESS"
    TaskStatusFailed     TaskStatusEnum = "FAILED"
)
```

### ConversationMessage

```go
type ConversationMessage struct {
    Role      string  // "user" or "assistant"
    Content   string  // Message content
    Name      *string // Speaker name (optional)
    CreatedAt *string // Timestamp in ISO format (optional)
}
```

## Retry Policy

The SDK includes a flexible retry policy system for handling transient failures.

### Default Retry Policy

By default, the SDK retries on:
- Network errors
- HTTP 429 (Too Many Requests)
- HTTP 500 (Internal Server Error)
- HTTP 502 (Bad Gateway)
- HTTP 503 (Service Unavailable)
- HTTP 504 (Gateway Timeout)

With exponential backoff: `baseDelay * 2^attempt` (capped at 32 seconds).

### Custom Retry Policy

```go
// Create a custom retry policy
customPolicy := memu.NewCustomRetryPolicy(
    5, // max retries
    func(attempt int, statusCode int, err error) bool {
        // Custom retry logic
        return attempt < 5 && (statusCode >= 500 || err != nil)
    },
    func(attempt int) time.Duration {
        // Custom backoff logic
        return time.Duration(attempt) * time.Second
    },
)

client, err := memu.NewClient(
    "your_api_key",
    memu.WithRetryPolicy(customPolicy),
)
```

### No Retry Policy

```go
// Disable retries
client, err := memu.NewClient(
    "your_api_key",
    memu.WithRetryPolicy(memu.NewNoRetryPolicy()),
)
```

## Error Handling

The SDK provides specific error types for different error cases:

```go
import memu "github.com/NevaMind-AI/memU-sdk-go"

result, err := client.Memorize(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *memu.AuthenticationError:
        // Invalid API key (401)
        fmt.Printf("Authentication failed: %v\n", e)
    case *memu.RateLimitError:
        // Rate limit exceeded (429)
        fmt.Printf("Rate limited. Retry after %.1f seconds\n", *e.RetryAfter)
    case *memu.NotFoundError:
        // Resource not found (404)
        fmt.Printf("Not found: %v\n", e)
    case *memu.ValidationError:
        // Request validation failed (422)
        fmt.Printf("Validation error: %v\n", e.Response)
    case *memu.ClientError:
        // Other API errors
        fmt.Printf("API error: %v\n", e)
    default:
        // Network or other errors
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Error Types

- `ClientError` - Base error type with status code and response
- `AuthenticationError` - Invalid API key (401)
- `RateLimitError` - Rate limit exceeded (429), includes RetryAfter field
- `NotFoundError` - Resource not found (404)
- `ValidationError` - Request validation failed (422)

## Examples

See the [examples](./examples/) directory for complete working examples:

- [`demo.go`](./examples/demo.go) - Complete workflow demonstration

To run the example:

```bash
export MEMU_API_KEY=your_api_key
go run examples/demo.go
```

## Context and Timeouts

All API methods accept a `context.Context` parameter for timeout and cancellation control:

```go
// Set a timeout for a single request
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.Memorize(ctx, req)
```

## Development

### Building

```bash
# Clone the repository
git clone https://github.com/NevaMind-AI/memU-sdk-go.git
cd memU-sdk-go

# Build
go build ./...

# Run tests
go test ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint
golangci-lint run

# Vet
go vet ./...
```

## Support

- [Full API Documentation](https://memu.pro/docs)
- [Discord Community](https://discord.gg/memu)
- [Report Issues](https://github.com/NevaMind-AI/memU-sdk-go/issues)

## License

MIT License - see [LICENSE](./LICENSE) for details.
