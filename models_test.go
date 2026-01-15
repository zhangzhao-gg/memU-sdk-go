/**
 * [INPUT]: 依赖 testing, time 标准库; 依赖 models.go 的所有数据模型
 * [OUTPUT]: 对外提供数据模型的单元测试，覆盖 JSON 序列化反序列化和字段验证
 * [POS]: SDK 根目录的测试层，验证 models.go 的正确性
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */
package memu

import (
	"encoding/json"
	"testing"
	"time"
)

// ============================================================
// 测试: TaskStatusEnum
// ============================================================

func TestTaskStatusEnum_Values(t *testing.T) {
	tests := []struct {
		status   TaskStatusEnum
		expected string
	}{
		{TaskStatusPending, "PENDING"},
		{TaskStatusProcessing, "PROCESSING"},
		{TaskStatusCompleted, "COMPLETED"},
		{TaskStatusSuccess, "SUCCESS"},
		{TaskStatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.expected {
			t.Errorf("Expected status '%s', got '%s'", tt.expected, string(tt.status))
		}
	}
}

// ============================================================
// 测试: MemoryResource
// ============================================================

func TestMemoryResource_JSONSerialization(t *testing.T) {
	now := time.Now()
	resource := &MemoryResource{
		ID:        stringPtr("res_123"),
		URL:       stringPtr("https://example.com/resource"),
		Modality:  stringPtr("text"),
		Caption:   stringPtr("Test caption"),
		CreatedAt: &now,
		UpdatedAt: &now,
		Metadata:  map[string]interface{}{"key": "value"},
	}

	// 序列化
	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal MemoryResource: %v", err)
	}

	// 反序列化
	var decoded MemoryResource
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal MemoryResource: %v", err)
	}

	// 验证字段
	if decoded.ID == nil || *decoded.ID != "res_123" {
		t.Errorf("Expected ID 'res_123', got %v", decoded.ID)
	}
	if decoded.URL == nil || *decoded.URL != "https://example.com/resource" {
		t.Errorf("Expected URL 'https://example.com/resource', got %v", decoded.URL)
	}
}

// ============================================================
// 测试: MemoryItem
// ============================================================

func TestMemoryItem_JSONSerialization(t *testing.T) {
	now := time.Now()
	score := 0.95
	item := &MemoryItem{
		ID:           stringPtr("item_123"),
		Summary:      stringPtr("Test summary"),
		Content:      stringPtr("Test content"),
		MemoryType:   stringPtr("preference"),
		CategoryID:   stringPtr("cat_123"),
		CategoryName: stringPtr("Test Category"),
		ResourceID:   stringPtr("res_123"),
		Score:        &score,
		CreatedAt:    &now,
		UpdatedAt:    &now,
		Metadata:     map[string]interface{}{"key": "value"},
	}

	// 序列化
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal MemoryItem: %v", err)
	}

	// 反序列化
	var decoded MemoryItem
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal MemoryItem: %v", err)
	}

	// 验证字段
	if decoded.ID == nil || *decoded.ID != "item_123" {
		t.Errorf("Expected ID 'item_123', got %v", decoded.ID)
	}
	if decoded.Score == nil || *decoded.Score != 0.95 {
		t.Errorf("Expected Score 0.95, got %v", decoded.Score)
	}
}

// ============================================================
// 测试: MemoryCategory
// ============================================================

func TestMemoryCategory_JSONSerialization(t *testing.T) {
	now := time.Now()
	score := 0.85
	itemCount := 10
	category := &MemoryCategory{
		ID:          stringPtr("cat_123"),
		Name:        stringPtr("Test Category"),
		Summary:     stringPtr("Test summary"),
		Description: stringPtr("Test description"),
		Content:     stringPtr("Test content"),
		ItemCount:   &itemCount,
		Score:       &score,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		Metadata:    map[string]interface{}{"key": "value"},
	}

	// 序列化
	data, err := json.Marshal(category)
	if err != nil {
		t.Fatalf("Failed to marshal MemoryCategory: %v", err)
	}

	// 反序列化
	var decoded MemoryCategory
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal MemoryCategory: %v", err)
	}

	// 验证字段
	if decoded.ID == nil || *decoded.ID != "cat_123" {
		t.Errorf("Expected ID 'cat_123', got %v", decoded.ID)
	}
	if decoded.ItemCount == nil || *decoded.ItemCount != 10 {
		t.Errorf("Expected ItemCount 10, got %v", decoded.ItemCount)
	}
}

// ============================================================
// 测试: TaskStatus
// ============================================================

func TestTaskStatus_JSONSerialization(t *testing.T) {
	now := time.Now()
	progress := 75.0
	status := &TaskStatus{
		TaskID:    "task_123",
		Status:    TaskStatusProcessing,
		Progress:  &progress,
		Message:   stringPtr("Processing..."),
		Result:    map[string]interface{}{"items": 5},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// 序列化
	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal TaskStatus: %v", err)
	}

	// 反序列化
	var decoded TaskStatus
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal TaskStatus: %v", err)
	}

	// 验证字段
	if decoded.TaskID != "task_123" {
		t.Errorf("Expected TaskID 'task_123', got %v", decoded.TaskID)
	}
	if decoded.Status != TaskStatusProcessing {
		t.Errorf("Expected Status 'PROCESSING', got %v", decoded.Status)
	}
	if decoded.Progress == nil || *decoded.Progress != 75.0 {
		t.Errorf("Expected Progress 75.0, got %v", decoded.Progress)
	}
}

// ============================================================
// 测试: ConversationMessage
// ============================================================

func TestConversationMessage_JSONSerialization(t *testing.T) {
	message := &ConversationMessage{
		Role:    "user",
		Content: "Hello, world!",
	}

	// 序列化
	data, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal ConversationMessage: %v", err)
	}

	// 反序列化
	var decoded ConversationMessage
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ConversationMessage: %v", err)
	}

	// 验证字段
	if decoded.Role != "user" {
		t.Errorf("Expected Role 'user', got %v", decoded.Role)
	}
	if decoded.Content != "Hello, world!" {
		t.Errorf("Expected Content 'Hello, world!', got %v", decoded.Content)
	}
}

// ============================================================
// 测试: MemorizeRequest
// ============================================================

func TestMemorizeRequest_JSONSerialization(t *testing.T) {
	waitForCompletion := true
	pollInterval := 2 * time.Second
	timeout := 300 * time.Second
	sessionDate := "2026-01-15"

	req := &MemorizeRequest{
		UserID:            "user_123",
		AgentID:           "agent_123",
		UserName:          "Test User",
		AgentName:         "Test Agent",
		SessionDate:       &sessionDate,
		WaitForCompletion: waitForCompletion,
		PollInterval:      pollInterval,
		Timeout:           timeout,
		Conversation: []ConversationMessage{
			{Role: "user", Content: "Hello"},
			{Role: "assistant", Content: "Hi there!"},
		},
	}

	// 序列化
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal MemorizeRequest: %v", err)
	}

	// 反序列化
	var decoded MemorizeRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal MemorizeRequest: %v", err)
	}

	// 验证字段
	if decoded.UserID != "user_123" {
		t.Errorf("Expected UserID 'user_123', got %v", decoded.UserID)
	}
	if len(decoded.Conversation) != 2 {
		t.Errorf("Expected 2 conversation messages, got %v", len(decoded.Conversation))
	}
}

// ============================================================
// 测试: RetrieveRequest
// ============================================================

func TestRetrieveRequest_JSONSerialization_StringQuery(t *testing.T) {
	req := &RetrieveRequest{
		UserID:  "user_123",
		AgentID: "agent_123",
		Query:   "test query",
	}

	// 序列化
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal RetrieveRequest: %v", err)
	}

	// 反序列化
	var decoded RetrieveRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal RetrieveRequest: %v", err)
	}

	// 验证字段
	if decoded.UserID != "user_123" {
		t.Errorf("Expected UserID 'user_123', got %v", decoded.UserID)
	}
	if queryStr, ok := decoded.Query.(string); !ok || queryStr != "test query" {
		t.Errorf("Expected Query 'test query', got %v", decoded.Query)
	}
}

func TestRetrieveRequest_JSONSerialization_MessagesQuery(t *testing.T) {
	req := &RetrieveRequest{
		UserID:  "user_123",
		AgentID: "agent_123",
		Query: []ConversationMessage{
			{Role: "user", Content: "What do you know about me?"},
		},
	}

	// 序列化
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal RetrieveRequest: %v", err)
	}

	// 反序列化
	var decoded RetrieveRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal RetrieveRequest: %v", err)
	}

	// 验证字段
	if decoded.UserID != "user_123" {
		t.Errorf("Expected UserID 'user_123', got %v", decoded.UserID)
	}
	// Query 应该是一个数组
	if decoded.Query == nil {
		t.Error("Expected Query to be non-nil")
	}
}

// ============================================================
// 测试: ListCategoriesRequest
// ============================================================

func TestListCategoriesRequest_JSONSerialization(t *testing.T) {
	agentID := "agent_123"
	req := &ListCategoriesRequest{
		UserID:  "user_123",
		AgentID: &agentID,
	}

	// 序列化
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ListCategoriesRequest: %v", err)
	}

	// 反序列化
	var decoded ListCategoriesRequest
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ListCategoriesRequest: %v", err)
	}

	// 验证字段
	if decoded.UserID != "user_123" {
		t.Errorf("Expected UserID 'user_123', got %v", decoded.UserID)
	}
	if decoded.AgentID == nil || *decoded.AgentID != "agent_123" {
		t.Errorf("Expected AgentID 'agent_123', got %v", decoded.AgentID)
	}
}

// ============================================================
// 辅助函数
// ============================================================

func stringPtr(s string) *string {
	return &s
}
