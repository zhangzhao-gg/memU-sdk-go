// Package memu defines all request and response models for the MemU SDK.
// This file contains data structures and validation interfaces used throughout the SDK.
package memu

import (
	"fmt"
)

// Validator defines the parameter validation interface.
// This provides unified validation logic to avoid code duplication.
type Validator interface {
	Validate() error
}

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
	// Modality specifies the type of resource (e.g., "text", "image", "audio").
	Modality *string `json:"modality,omitempty"`
	// ResourceURL is the URL where the resource is stored.
	ResourceURL *string `json:"resource_url,omitempty"`
	// Caption is a textual description of the resource.
	Caption *string `json:"caption,omitempty"`
	// Content contains the actual resource data as a flexible map.
	Content map[string]interface{} `json:"content,omitempty"`
	// Metadata contains additional metadata about the resource.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// MemoryItem represents a discrete memory unit extracted from resources.
// Memory items are individual pieces of information such as preferences,
// skills, opinions, habits, relationships, etc.
type MemoryItem struct {
	// Content is the textual content of the memory item.
	Content *string `json:"content,omitempty"`
	// MemoryType categorizes the type of memory (e.g., "preference", "skill", "fact").
	MemoryType *string `json:"memory_type,omitempty"`
}

// MemoryCategory represents an aggregated memory category.
// Categories organize related memory items and provide summaries
// of clustered information (e.g., preferences.md, work_life.md).
type MemoryCategory struct {
	// Name is the category name (e.g., "preferences", "work_life").
	Name *string `json:"name,omitempty"`
	// Description provides a brief description of the category.
	Description *string `json:"description,omitempty"`
	// Summary contains a summary of all memories in this category.
	Summary *string `json:"summary,omitempty"`
	// UserID is the user ID this category belongs to.
	UserID *string `json:"user_id,omitempty"`
	// AgentID is the agent ID this category is associated with.
	AgentID *string `json:"agent_id,omitempty"`
}

// TaskStatus represents status information for an asynchronous memorization task.
type TaskStatus struct {
	// TaskID is the unique identifier for the task.
	TaskID string `json:"task_id"`
	// Status indicates the current status of the task.
	Status TaskStatusEnum `json:"status"`
	// Message provides a human-readable status message.
	Message string `json:"message,omitempty"`
	// DetailInfo contains additional detailed information about the task.
	DetailInfo string `json:"detail_info,omitempty"`
}

// RetrieveResult represents the result of a memory retrieval operation.
type RetrieveResult struct {
	// RewrittenQuery is the query after being rewritten by the system for better retrieval.
	RewrittenQuery *string `json:"rewritten_query,omitempty"`
	// Categories contains the retrieved memory categories.
	Categories []*MemoryCategory `json:"categories,omitempty"`
	// Items contains the retrieved memory items.
	Items []*MemoryItem `json:"items,omitempty"`
	// Resources contains the retrieved memory resources.
	Resources []*MemoryResource `json:"resources,omitempty"`
}

// ConversationMessage represents a single message in a conversation.
type ConversationMessage struct {
	// Role is the role of the message sender (e.g., "user", "assistant", "system").
	Role string `json:"role"`
	// Content is the textual content of the message.
	Content string `json:"content"`
	// Name is an optional name for the message sender.
	Name *string `json:"name,omitempty"`
	// CreatedAt is an optional timestamp for when the message was created.
	CreatedAt *string `json:"created_at,omitempty"`
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
}

// MemorizeResult represents the result of a memorization operation.
// The API returns only task_id, status, and message.
// To get the extracted memories, use GetTaskStatus or Retrieve API.
type MemorizeResult struct {
	// TaskID is the unique identifier for the memorization task.
	TaskID *string `json:"task_id,omitempty"`
	// Status indicates the current status of the memorization task.
	Status *string `json:"status,omitempty"`
	// Message provides a human-readable message about the task.
	Message *string `json:"message,omitempty"`
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

// Validate validates MemorizeRequest parameters.
func (r *MemorizeRequest) Validate() error {
	if r.UserID == "" {
		return fmt.Errorf("Memorize: UserID is required")
	}
	if r.AgentID == "" {
		return fmt.Errorf("Memorize: AgentID is required")
	}
	if len(r.Conversation) == 0 && r.ConversationText == nil {
		return fmt.Errorf("Memorize: either Conversation or ConversationText must be provided")
	}
	if len(r.Conversation) > 0 && len(r.Conversation) < 3 {
		return fmt.Errorf("Memorize: Conversation must contain at least 3 messages")
	}
	return nil
}

// Validate validates RetrieveRequest parameters.
func (r *RetrieveRequest) Validate() error {
	if r.Query == nil {
		return fmt.Errorf("Retrieve: Query is required")
	}
	if r.UserID == "" {
		return fmt.Errorf("Retrieve: UserID is required")
	}
	if r.AgentID == "" {
		return fmt.Errorf("Retrieve: AgentID is required")
	}
	return nil
}

// Validate validates ListCategoriesRequest parameters.
func (r *ListCategoriesRequest) Validate() error {
	if r.UserID == "" {
		return fmt.Errorf("ListCategories: UserID is required")
	}
	return nil
}
