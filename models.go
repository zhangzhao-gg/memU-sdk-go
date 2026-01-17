/**
 * [INPUT]: 依赖 time 的 Time 类型
 * [OUTPUT]: 对外提供 TaskStatusEnum, MemoryResource, MemoryItem, MemoryCategory, TaskStatus, MemorizeResult, RetrieveResult, MemorizeRequest, RetrieveRequest, ListCategoriesRequest, ConversationMessage, Validator 接口
 * [POS]: SDK 根目录的数据层，定义所有请求响应模型和验证接口，被 client.go 和用户代码消费
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */
package memu

import (
	"fmt"
)

// Validator 定义参数验证接口
// ============================================================
// 消除重复: 统一的验证逻辑，避免在每个方法中重复验证代码
// ============================================================
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
	Modality    *string                `json:"modality,omitempty"`
	ResourceURL *string                `json:"resource_url,omitempty"`
	Caption     *string                `json:"caption,omitempty"`
	Content     map[string]interface{} `json:"content,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// MemoryItem represents a discrete memory unit extracted from resources.
// Memory items are individual pieces of information such as preferences,
// skills, opinions, habits, relationships, etc.
type MemoryItem struct {
	Content    *string `json:"content,omitempty"`
	MemoryType *string `json:"memory_type,omitempty"`
}

// MemoryCategory represents an aggregated memory category.
// Categories organize related memory items and provide summaries
// of clustered information (e.g., preferences.md, work_life.md).
type MemoryCategory struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Summary     *string `json:"summary,omitempty"`
	UserID      *string `json:"user_id,omitempty"`
	AgentID     *string `json:"agent_id,omitempty"`
}

// TaskStatus represents status information for an asynchronous memorization task.
type TaskStatus struct {
	TaskID  string         `json:"task_id"`
	Status  TaskStatusEnum `json:"status"`
	Message *string        `json:"message,omitempty"`
}

// RetrieveResult represents the result of a memory retrieval operation.
type RetrieveResult struct {
	RewrittenQuery *string           `json:"rewritten_query,omitempty"`
	Categories     []*MemoryCategory `json:"categories,omitempty"`
	Items          []*MemoryItem     `json:"items,omitempty"`
	Resources      []*MemoryResource `json:"resources,omitempty"`
}

// ConversationMessage represents a single message in a conversation.
type ConversationMessage struct {
	Role      string  `json:"role"`
	Content   string  `json:"content"`
	Name      *string `json:"name,omitempty"`
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
	TaskID  *string `json:"task_id,omitempty"`
	Status  *string `json:"status,omitempty"`
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

// ============================================================
// 参数验证方法
// ============================================================

// Validate 验证 MemorizeRequest 参数
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

// Validate 验证 RetrieveRequest 参数
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

// Validate 验证 ListCategoriesRequest 参数
func (r *ListCategoriesRequest) Validate() error {
	if r.UserID == "" {
		return fmt.Errorf("ListCategories: UserID is required")
	}
	return nil
}
