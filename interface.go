// Package memu defines the core interface for the MemU SDK.
// This interface allows for easy mocking and testing.
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
