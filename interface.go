/**
 * [INPUT]: 依赖 context 的 Context; 依赖 models.go 的 MemorizeRequest, RetrieveRequest, ListCategoriesRequest, MemorizeResult, RetrieveResult, MemoryCategory, TaskStatus
 * [OUTPUT]: 对外提供 MemUClient 接口定义
 * [POS]: SDK 根目录的接口层，定义核心接口契约，被 client.go 实现，被用户代码和测试代码消费
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */
package memu

import (
	"context"
)

// MemUClient defines the interface for interacting with the MemU API.
// This interface allows for easy mocking and testing.
type MemUClient interface {
	// Memorize memorizes a conversation and extracts structured memory.
	Memorize(ctx context.Context, req *MemorizeRequest) (*MemorizeResult, error)

	// GetTaskStatus gets the status of a memorization task.
	GetTaskStatus(ctx context.Context, taskID string) (*TaskStatus, error)

	// Retrieve retrieves relevant memories based on a query.
	Retrieve(ctx context.Context, req *RetrieveRequest) (*RetrieveResult, error)

	// ListCategories lists all memory categories.
	ListCategories(ctx context.Context, req *ListCategoriesRequest) ([]*MemoryCategory, error)
}

// Ensure Client implements MemUClient interface
var _ MemUClient = (*Client)(nil)
