# MemU Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/NevaMind-AI/memU-sdk-go.svg)](https://pkg.go.dev/github.com/NevaMind-AI/memU-sdk-go)
[![Go 1.21+](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for the [MemU Cloud API](https://api.memu.so) - Manage structured, long-term memory for AI agents.

## Features

- 🚀 **Full Cloud API v3 Coverage** - memorize, retrieve, categories, task status
- ⚡ **Context Support** - Native Go context for timeout and cancellation
- 🔄 **Automatic Retry** - Exponential backoff for failed requests
- ⏱️ **Rate Limit Handling** - Respects Retry-After headers
- 🛡️ **Type Safety** - Strongly typed models with full documentation
- 🎯 **Custom Errors** - Specific error types for different failure cases
- 📦 **Zero Dependencies** - Uses only Go standard library

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
        UserID:            "user_123",
        AgentID:           "my_assistant",
        WaitForCompletion: true,
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Task ID: %s\n", *result.TaskID)

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

Memorize a conversation and extract structured memory.

```go
func (c *Client) Memorize(ctx context.Context, req *MemorizeRequest) (*MemorizeResult, error)
```

**Request Fields:**
- `Conversation` - List of conversation messages (optional if ConversationText is provided)
- `ConversationText` - Alternative: raw conversation text (optional if Conversation is provided)
- `UserID` - User ID for scoping the memory (required)
- `AgentID` - Agent ID for scoping the memory (required)
- `UserName` - Display name for the user (default: "User")
- `AgentName` - Display name for the agent (default: "Assistant")
- `SessionDate` - Optional session date in ISO format
- `WaitForCompletion` - If true, poll until the task completes (default: false)
- `PollInterval` - Seconds between status checks when waiting (default: 2s)
- `Timeout` - Maximum time to wait for completion (default: 5 minutes)

**Example:**
```go
result, err := client.Memorize(ctx, &memu.MemorizeRequest{
    Conversation: []memu.ConversationMessage{
        {Role: "user", Content: "I love pasta"},
        {Role: "assistant", Content: "Great choice!"},
    },
    UserID:            "user_123",
    AgentID:           "agent_456",
    WaitForCompletion: true,
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
    fmt.Printf("%s: %s\n", *cat.Name, *cat.Summary)
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
    fmt.Printf("Task completed: %v\n", status.Result)
}
```

## Data Models

### MemorizeResult

```go
type MemorizeResult struct {
    TaskID     *string           // Task ID for async tracking
    Resource   *MemoryResource   // Created resource
    Items      []*MemoryItem     // Extracted memory items
    Categories []*MemoryCategory // Updated categories
}
```

### RetrieveResult

```go
type RetrieveResult struct {
    Categories    []*MemoryCategory // Relevant categories
    Items         []*MemoryItem     // Relevant memory items
    Resources     []*MemoryResource // Related raw resources
    NextStepQuery *string           // Rewritten query (if applicable)
}
```

### MemoryItem

```go
type MemoryItem struct {
    ID           *string    // Unique identifier
    Summary      *string    // Summary/description
    Content      *string    // Content text
    MemoryType   *string    // Type: profile, event, preference, etc.
    CategoryID   *string    // Category ID
    CategoryName *string    // Category name
    Score        *float64   // Relevance score (in retrieve)
    CreatedAt    *time.Time // Creation timestamp
    UpdatedAt    *time.Time // Last update timestamp
}
```

### MemoryCategory

```go
type MemoryCategory struct {
    ID          *string    // Unique identifier
    Name        *string    // Category name (e.g., 'personal info')
    Summary     *string    // Summary of content
    Content     *string    // Full content
    Description *string    // Description
    ItemCount   *int       // Number of items
    Score       *float64   // Relevance score (in retrieve)
    CreatedAt   *time.Time // Creation timestamp
    UpdatedAt   *time.Time // Last update timestamp
}
```

### TaskStatus

```go
type TaskStatus struct {
    TaskID    string                 // Task identifier
    Status    TaskStatusEnum         // PENDING, PROCESSING, COMPLETED, SUCCESS, FAILED
    Progress  *float64               // Progress percentage (0-100)
    Message   *string                // Status message or error
    Result    map[string]interface{} // Task result when completed
    CreatedAt *time.Time             // Task creation timestamp
    UpdatedAt *time.Time             // Last update timestamp
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

- 📚 [Full API Documentation](https://memu.pro/docs)
- 💬 [Discord Community](https://discord.gg/memu)
- 🐛 [Report Issues](https://github.com/NevaMind-AI/memU-sdk-go/issues)

## License

MIT License - see [LICENSE](./LICENSE) for details.
