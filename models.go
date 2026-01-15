package memu

import (
	"time"
)

// TaskStatusEnum represents the status of an asynchronous memorization task.
type TaskStatusEnum string

const (
	// TaskStatusPending indicates the task is pending.
	TaskStatusPending TaskStatusEnum = "PENDING"
	// TaskStatusProcessing indicates the task is being processed.
	TaskStatusProcessing TaskStatusEnum = "PROCESSING"
	// TaskStatusCompleted indicates the task has completed.
	TaskStatusCompleted TaskStatusEnum = "COMPLETED"
	// TaskStatusSuccess indicates the task succeeded.
	TaskStatusSuccess TaskStatusEnum = "SUCCESS"
	// TaskStatusFailed indicates the task failed.
	TaskStatusFailed TaskStatusEnum = "FAILED"
)

// MemoryResource represents a raw resource stored in MemU.
// Resources are the source materials (conversations, documents, images, etc.)
// from which memory items are extracted.
type MemoryResource struct {
	ID        *string                `json:"id,omitempty"`
	URL       *string                `json:"url,omitempty"`
	Modality  *string                `json:"modality,omitempty"`
	Caption   *string                `json:"caption,omitempty"`
	CreatedAt *time.Time             `json:"created_at,omitempty"`
	UpdatedAt *time.Time             `json:"updated_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// MemoryItem represents a discrete memory unit extracted from resources.
// Memory items are individual pieces of information such as preferences,
// skills, opinions, habits, relationships, etc.
type MemoryItem struct {
	ID           *string                `json:"id,omitempty"`
	Summary      *string                `json:"summary,omitempty"`
	Content      *string                `json:"content,omitempty"`
	MemoryType   *string                `json:"memory_type,omitempty"`
	CategoryID   *string                `json:"category_id,omitempty"`
	CategoryName *string                `json:"category_name,omitempty"`
	ResourceID   *string                `json:"resource_id,omitempty"`
	Score        *float64               `json:"score,omitempty"`
	CreatedAt    *time.Time             `json:"created_at,omitempty"`
	UpdatedAt    *time.Time             `json:"updated_at,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// MemoryCategory represents an aggregated memory category.
// Categories organize related memory items and provide summaries
// of clustered information (e.g., preferences.md, work_life.md).
type MemoryCategory struct {
	ID          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Summary     *string                `json:"summary,omitempty"`
	Description *string                `json:"description,omitempty"`
	Content     *string                `json:"content,omitempty"`
	ItemCount   *int                   `json:"item_count,omitempty"`
	Score       *float64               `json:"score,omitempty"`
	CreatedAt   *time.Time             `json:"created_at,omitempty"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TaskStatus represents status information for an asynchronous memorization task.
type TaskStatus struct {
	TaskID    string                 `json:"task_id"`
	Status    TaskStatusEnum         `json:"status"`
	Progress  *float64               `json:"progress,omitempty"`
	Message   *string                `json:"message,omitempty"`
	Result    map[string]interface{} `json:"result,omitempty"`
	CreatedAt *time.Time             `json:"created_at,omitempty"`
	UpdatedAt *time.Time             `json:"updated_at,omitempty"`
}

// MemorizeResult represents the result of a memorization operation.
type MemorizeResult struct {
	TaskID     *string           `json:"task_id,omitempty"`
	Resource   *MemoryResource   `json:"resource,omitempty"`
	Items      []*MemoryItem     `json:"items,omitempty"`
	Categories []*MemoryCategory `json:"categories,omitempty"`
}

// RetrieveResult represents the result of a memory retrieval operation.
type RetrieveResult struct {
	Categories    []*MemoryCategory `json:"categories,omitempty"`
	Items         []*MemoryItem     `json:"items,omitempty"`
	Resources     []*MemoryResource `json:"resources,omitempty"`
	NextStepQuery *string           `json:"next_step_query,omitempty"`
}

// ConversationMessage represents a single message in a conversation.
type ConversationMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MemorizeRequest represents a request to memorize a conversation.
type MemorizeRequest struct {
	// Conversation is a list of conversation messages.
	Conversation []ConversationMessage `json:"conversation,omitempty"`
	// ConversationText is an alternative to Conversation for raw text.
	ConversationText *string `json:"conversation_text,omitempty"`
	// UserID is the user ID for scoping the memory (required).
	UserID string `json:"user_id"`
	// AgentID is the agent ID for scoping the memory (required).
	AgentID string `json:"agent_id"`
	// UserName is the display name for the user (default: "User").
	UserName string `json:"user_name,omitempty"`
	// AgentName is the display name for the agent (default: "Assistant").
	AgentName string `json:"agent_name,omitempty"`
	// SessionDate is an optional session date in ISO format.
	SessionDate *string `json:"session_date,omitempty"`
	// WaitForCompletion indicates whether to poll until the task completes.
	WaitForCompletion bool `json:"-"`
	// PollInterval is the interval between status checks when waiting (default: 2s).
	PollInterval time.Duration `json:"-"`
	// Timeout is the maximum time to wait for completion (default: 5 minutes).
	Timeout time.Duration `json:"-"`
}

// RetrieveRequest represents a request to retrieve memories.
type RetrieveRequest struct {
	// Query can be a string or a list of conversation messages.
	Query interface{} `json:"query"`
	// UserID is the user ID for scoping (required).
	UserID string `json:"user_id"`
	// AgentID is the agent ID for scoping (required).
	AgentID string `json:"agent_id"`
}

// ListCategoriesRequest represents a request to list memory categories.
type ListCategoriesRequest struct {
	// UserID is the user ID for scoping (required).
	UserID string `json:"user_id"`
	// AgentID is the agent ID for scoping (optional).
	AgentID *string `json:"agent_id,omitempty"`
}
