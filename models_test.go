// Package memu provides unit tests for data models and validation.
// This file validates model structures and Validate methods.
package memu

import (
	"strings"
	"testing"
)

// Helper function for creating string pointers.
func strPtr(s string) *string {
	return &s
}

// TestMemoryItem tests MemoryItem model.
func TestMemoryItem(t *testing.T) {
	content := "User prefers Italian food"
	memType := "preference"

	item := MemoryItem{
		Content:    &content,
		MemoryType: &memType,
	}

	if item.Content == nil || *item.Content != content {
		t.Errorf("expected Content '%s', got '%v'", content, item.Content)
	}
	if item.MemoryType == nil || *item.MemoryType != memType {
		t.Errorf("expected MemoryType '%s', got '%v'", memType, item.MemoryType)
	}
}

func TestMemoryItem_NilFields(t *testing.T) {
	item := MemoryItem{}

	if item.Content != nil {
		t.Error("expected Content to be nil")
	}
	if item.MemoryType != nil {
		t.Error("expected MemoryType to be nil")
	}
}

// TestMemoryCategory tests MemoryCategory model.
func TestMemoryCategory(t *testing.T) {
	name := "preferences"
	summary := "User preferences"
	desc := "Category for user preferences"

	category := MemoryCategory{
		Name:        &name,
		Summary:     &summary,
		Description: &desc,
	}

	if category.Name == nil || *category.Name != name {
		t.Errorf("expected Name '%s', got '%v'", name, category.Name)
	}
	if category.Summary == nil || *category.Summary != summary {
		t.Errorf("expected Summary '%s', got '%v'", summary, category.Summary)
	}
	if category.Description == nil || *category.Description != desc {
		t.Errorf("expected Description '%s', got '%v'", desc, category.Description)
	}
}

// TestMemoryResource tests MemoryResource model.
func TestMemoryResource(t *testing.T) {
	url := "https://example.com/chat.json"
	modality := "conversation"
	caption := "A conversation"

	resource := MemoryResource{
		ResourceURL: &url,
		Modality:    &modality,
		Caption:     &caption,
	}

	if resource.ResourceURL == nil || *resource.ResourceURL != url {
		t.Errorf("expected ResourceURL '%s', got '%v'", url, resource.ResourceURL)
	}
	if resource.Modality == nil || *resource.Modality != modality {
		t.Errorf("expected Modality '%s', got '%v'", modality, resource.Modality)
	}
}

func TestMemoryResource_WithMetadata(t *testing.T) {
	resource := MemoryResource{
		Metadata: map[string]interface{}{
			"date":   "2024-01-15",
			"source": "chat",
		},
	}

	if resource.Metadata == nil {
		t.Fatal("expected Metadata to not be nil")
	}
	if resource.Metadata["date"] != "2024-01-15" {
		t.Errorf("expected Metadata['date'] '2024-01-15', got '%v'", resource.Metadata["date"])
	}
}

// TestTaskStatus tests TaskStatus model.
func TestTaskStatus(t *testing.T) {
	status := TaskStatus{
		TaskID:  "task_123",
		Status:  TaskStatusCompleted,
		Message: "Task completed successfully",
	}

	if status.TaskID != "task_123" {
		t.Errorf("expected TaskID 'task_123', got '%s'", status.TaskID)
	}
	if status.Status != TaskStatusCompleted {
		t.Errorf("expected Status COMPLETED, got '%s'", status.Status)
	}
}

func TestTaskStatusEnum_Values(t *testing.T) {
	tests := []struct {
		status TaskStatusEnum
		value  string
	}{
		{TaskStatusPending, "PENDING"},
		{TaskStatusProcessing, "PROCESSING"},
		{TaskStatusCompleted, "COMPLETED"},
		{TaskStatusSuccess, "SUCCESS"},
		{TaskStatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.value {
			t.Errorf("expected TaskStatusEnum '%s', got '%s'", tt.value, tt.status)
		}
	}
}

// TestMemorizeResult tests MemorizeResult model.
func TestMemorizeResult(t *testing.T) {
	taskID := "task_123"
	status := "PENDING"
	message := "Task submitted"

	result := MemorizeResult{
		TaskID:  &taskID,
		Status:  &status,
		Message: &message,
	}

	if result.TaskID == nil || *result.TaskID != taskID {
		t.Errorf("expected TaskID '%s', got '%v'", taskID, result.TaskID)
	}
	if result.Status == nil || *result.Status != status {
		t.Errorf("expected Status '%s', got '%v'", status, result.Status)
	}
}

// TestRetrieveResult tests RetrieveResult model.
func TestRetrieveResult(t *testing.T) {
	result := RetrieveResult{
		Categories: []*MemoryCategory{},
		Items:      []*MemoryItem{},
		Resources:  []*MemoryResource{},
	}

	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}
	if len(result.Categories) != 0 {
		t.Errorf("expected 0 categories, got %d", len(result.Categories))
	}
	if len(result.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(result.Resources))
	}
}

func TestRetrieveResult_WithData(t *testing.T) {
	content := "User likes pizza"
	memType := "preference"
	rewrittenQuery := "What food does the user like?"

	result := RetrieveResult{
		RewrittenQuery: &rewrittenQuery,
		Items: []*MemoryItem{
			{Content: &content, MemoryType: &memType},
		},
	}

	if result.RewrittenQuery == nil || *result.RewrittenQuery != rewrittenQuery {
		t.Errorf("expected RewrittenQuery '%s', got '%v'", rewrittenQuery, result.RewrittenQuery)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	if *result.Items[0].Content != content {
		t.Errorf("expected item content '%s', got '%s'", content, *result.Items[0].Content)
	}
}

// TestConversationMessage tests ConversationMessage model.
func TestConversationMessage(t *testing.T) {
	name := "John"
	createdAt := "2024-01-15T10:30:00Z"

	msg := ConversationMessage{
		Role:      "user",
		Content:   "Hello, world!",
		Name:      &name,
		CreatedAt: &createdAt,
	}

	if msg.Role != "user" {
		t.Errorf("expected Role 'user', got '%s'", msg.Role)
	}
	if msg.Content != "Hello, world!" {
		t.Errorf("expected Content 'Hello, world!', got '%s'", msg.Content)
	}
	if msg.Name == nil || *msg.Name != name {
		t.Errorf("expected Name '%s', got '%v'", name, msg.Name)
	}
}

// TestMemorizeRequest_Validate tests MemorizeRequest validation.
func TestMemorizeRequest_Validate_Valid(t *testing.T) {
	req := &MemorizeRequest{
		UserID:  "user_123",
		AgentID: "agent_456",
		Conversation: []ConversationMessage{
			{Role: "user", Content: "Message 1"},
			{Role: "assistant", Content: "Message 2"},
			{Role: "user", Content: "Message 3"},
		},
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestMemorizeRequest_Validate_MissingUserID(t *testing.T) {
	req := &MemorizeRequest{
		AgentID: "agent_456",
		Conversation: []ConversationMessage{
			{Role: "user", Content: "Message 1"},
			{Role: "assistant", Content: "Message 2"},
			{Role: "user", Content: "Message 3"},
		},
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for missing UserID")
	}
	if !strings.Contains(err.Error(), "UserID") {
		t.Errorf("expected error message to contain 'UserID', got: %v", err)
	}
}

func TestMemorizeRequest_Validate_MissingAgentID(t *testing.T) {
	req := &MemorizeRequest{
		UserID: "user_123",
		Conversation: []ConversationMessage{
			{Role: "user", Content: "Message 1"},
			{Role: "assistant", Content: "Message 2"},
			{Role: "user", Content: "Message 3"},
		},
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for missing AgentID")
	}
	if !strings.Contains(err.Error(), "AgentID") {
		t.Errorf("expected error message to contain 'AgentID', got: %v", err)
	}
}

func TestMemorizeRequest_Validate_MissingConversation(t *testing.T) {
	req := &MemorizeRequest{
		UserID:  "user_123",
		AgentID: "agent_456",
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for missing conversation")
	}
	if !strings.Contains(err.Error(), "Conversation") {
		t.Errorf("expected error message to contain 'Conversation', got: %v", err)
	}
}

func TestMemorizeRequest_Validate_TooFewMessages(t *testing.T) {
	req := &MemorizeRequest{
		UserID:  "user_123",
		AgentID: "agent_456",
		Conversation: []ConversationMessage{
			{Role: "user", Content: "Message 1"},
			{Role: "assistant", Content: "Message 2"},
		},
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for too few messages")
	}
	if !strings.Contains(err.Error(), "at least 3") {
		t.Errorf("expected error message to contain 'at least 3', got: %v", err)
	}
}

func TestMemorizeRequest_Validate_WithConversationText(t *testing.T) {
	text := "User: Hello\nAssistant: Hi there!"
	req := &MemorizeRequest{
		UserID:           "user_123",
		AgentID:          "agent_456",
		ConversationText: &text,
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error with ConversationText, got: %v", err)
	}
}

// TestRetrieveRequest_Validate tests RetrieveRequest validation.
func TestRetrieveRequest_Validate_Valid(t *testing.T) {
	req := &RetrieveRequest{
		Query:   "What are the user's hobbies?",
		UserID:  "user_123",
		AgentID: "agent_456",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRetrieveRequest_Validate_MissingQuery(t *testing.T) {
	req := &RetrieveRequest{
		UserID:  "user_123",
		AgentID: "agent_456",
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for missing Query")
	}
	if !strings.Contains(err.Error(), "Query") {
		t.Errorf("expected error message to contain 'Query', got: %v", err)
	}
}

func TestRetrieveRequest_Validate_MissingUserID(t *testing.T) {
	req := &RetrieveRequest{
		Query:   "What are the user's hobbies?",
		AgentID: "agent_456",
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for missing UserID")
	}
}

func TestRetrieveRequest_Validate_MissingAgentID(t *testing.T) {
	req := &RetrieveRequest{
		Query:  "What are the user's hobbies?",
		UserID: "user_123",
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for missing AgentID")
	}
}

func TestRetrieveRequest_Validate_ConversationQuery(t *testing.T) {
	req := &RetrieveRequest{
		Query: []ConversationMessage{
			{Role: "user", Content: "Tell me about their hobbies"},
		},
		UserID:  "user_123",
		AgentID: "agent_456",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error with conversation query, got: %v", err)
	}
}

// TestListCategoriesRequest_Validate tests ListCategoriesRequest validation.
func TestListCategoriesRequest_Validate_Valid(t *testing.T) {
	req := &ListCategoriesRequest{
		UserID: "user_123",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestListCategoriesRequest_Validate_WithAgentID(t *testing.T) {
	agentID := "agent_456"
	req := &ListCategoriesRequest{
		UserID:  "user_123",
		AgentID: &agentID,
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestListCategoriesRequest_Validate_MissingUserID(t *testing.T) {
	req := &ListCategoriesRequest{}

	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for missing UserID")
	}
	if !strings.Contains(err.Error(), "UserID") {
		t.Errorf("expected error message to contain 'UserID', got: %v", err)
	}
}
