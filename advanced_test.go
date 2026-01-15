/**
 * [INPUT]: 依赖 testing, net/http/httptest, context, time, encoding/json
 * [OUTPUT]: 对 client.go 中高级功能的测试覆盖
 * [POS]: SDK 根目录的测试层，测试异步任务和复杂场景
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */
package memu

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMemorize_WithWaitForCompletion_Success(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v3/memory/memorize" {
			// Return task ID
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"task_id": "test-task-123",
			})
		} else if r.URL.Path == "/api/v3/memory/memorize/status/test-task-123" {
			callCount++
			w.WriteHeader(http.StatusOK)
			if callCount >= 2 {
				// Task completed
				json.NewEncoder(w).Encode(map[string]interface{}{
					"task_id": "test-task-123",
					"status":  "COMPLETED",
					"result": map[string]interface{}{
						"resource": map[string]interface{}{
							"id": "resource-1",
						},
						"items": []interface{}{
							map[string]interface{}{
								"id":      "item-1",
								"content": "test content",
							},
						},
						"categories": []interface{}{
							map[string]interface{}{
								"id":   "cat-1",
								"name": "test category",
							},
						},
					},
				})
			} else {
				// Task still processing
				json.NewEncoder(w).Encode(map[string]interface{}{
					"task_id": "test-task-123",
					"status":  "PROCESSING",
				})
			}
		}
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	req := &MemorizeRequest{
		UserID:            "user-1",
		AgentID:           "agent-1",
		Conversation:      []ConversationMessage{{Role: "user", Content: "test"}},
		WaitForCompletion: true,
		PollInterval:      10 * time.Millisecond,
		Timeout:           5 * time.Second,
	}

	result, err := client.Memorize(context.Background(), req)
	if err != nil {
		t.Fatalf("Memorize() error = %v", err)
	}

	if result.Resource == nil {
		t.Error("Expected resource to be parsed")
	}

	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result.Items))
	}

	if len(result.Categories) != 1 {
		t.Errorf("Expected 1 category, got %d", len(result.Categories))
	}
}

func TestMemorize_WithWaitForCompletion_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v3/memory/memorize" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"task_id": "test-task-123",
			})
		} else if r.URL.Path == "/api/v3/memory/memorize/status/test-task-123" {
			// Always return processing
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"task_id": "test-task-123",
				"status":  "PROCESSING",
			})
		}
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	req := &MemorizeRequest{
		UserID:            "user-1",
		AgentID:           "agent-1",
		Conversation:      []ConversationMessage{{Role: "user", Content: "test"}},
		WaitForCompletion: true,
		PollInterval:      10 * time.Millisecond,
		Timeout:           50 * time.Millisecond,
	}

	_, err := client.Memorize(context.Background(), req)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestMemorize_WithWaitForCompletion_Failed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v3/memory/memorize" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"task_id": "test-task-123",
			})
		} else if r.URL.Path == "/api/v3/memory/memorize/status/test-task-123" {
			// Return failed status
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"task_id": "test-task-123",
				"status":  "FAILED",
				"message": "processing failed",
			})
		}
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	req := &MemorizeRequest{
		UserID:            "user-1",
		AgentID:           "agent-1",
		Conversation:      []ConversationMessage{{Role: "user", Content: "test"}},
		WaitForCompletion: true,
		PollInterval:      10 * time.Millisecond,
		Timeout:           5 * time.Second,
	}

	_, err := client.Memorize(context.Background(), req)
	if err == nil {
		t.Error("Expected error for failed task")
	}
}

func TestMemorize_WithConversationText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		if _, ok := payload["conversation_text"]; !ok {
			t.Error("Expected conversation_text in payload")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"resource": map[string]interface{}{
				"id": "resource-1",
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	conversationText := "This is a test conversation"
	req := &MemorizeRequest{
		UserID:           "user-1",
		AgentID:          "agent-1",
		ConversationText: &conversationText,
	}

	_, err := client.Memorize(context.Background(), req)
	if err != nil {
		t.Fatalf("Memorize() error = %v", err)
	}
}

func TestMemorize_WithCustomNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		if payload["user_name"] != "Alice" {
			t.Errorf("Expected user_name=Alice, got %v", payload["user_name"])
		}

		if payload["agent_name"] != "Bob" {
			t.Errorf("Expected agent_name=Bob, got %v", payload["agent_name"])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	req := &MemorizeRequest{
		UserID:       "user-1",
		AgentID:      "agent-1",
		Conversation: []ConversationMessage{{Role: "user", Content: "test"}},
		UserName:     "Alice",
		AgentName:    "Bob",
	}

	_, err := client.Memorize(context.Background(), req)
	if err != nil {
		t.Fatalf("Memorize() error = %v", err)
	}
}

func TestMemorize_WithSessionDate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		if _, ok := payload["session_date"]; !ok {
			t.Error("Expected session_date in payload")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	sessionDate := "2024-01-01"
	req := &MemorizeRequest{
		UserID:       "user-1",
		AgentID:      "agent-1",
		Conversation: []ConversationMessage{{Role: "user", Content: "test"}},
		SessionDate:  &sessionDate,
	}

	_, err := client.Memorize(context.Background(), req)
	if err != nil {
		t.Fatalf("Memorize() error = %v", err)
	}
}

func TestValidate_MemorizeRequest_MissingConversation(t *testing.T) {
	req := &MemorizeRequest{
		UserID:  "user-1",
		AgentID: "agent-1",
	}

	err := req.Validate()
	if err == nil {
		t.Error("Expected validation error for missing conversation")
	}
}

func TestValidate_RetrieveRequest_MissingQuery(t *testing.T) {
	req := &RetrieveRequest{
		UserID:  "user-1",
		AgentID: "agent-1",
	}

	err := req.Validate()
	if err == nil {
		t.Error("Expected validation error for missing query")
	}
}

func TestValidate_ListCategoriesRequest_MissingUserID(t *testing.T) {
	req := &ListCategoriesRequest{}

	err := req.Validate()
	if err == nil {
		t.Error("Expected validation error for missing user_id")
	}
}

func TestParseJSONObject_NilData(t *testing.T) {
	result, err := parseJSONObject[MemoryResource](nil)
	if err != nil {
		t.Errorf("parseJSONObject(nil) error = %v", err)
	}
	if result != nil {
		t.Error("Expected nil result for nil data")
	}
}

func TestParseJSONArray_EmptyArray(t *testing.T) {
	result, err := parseJSONArray[MemoryItem]([]interface{}{})
	if err != nil {
		t.Errorf("parseJSONArray([]) error = %v", err)
	}
	if result != nil {
		t.Error("Expected nil result for empty array")
	}
}

func TestListCategories_WithAgentID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		if _, ok := payload["agent_id"]; !ok {
			t.Error("Expected agent_id in payload")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"categories": []interface{}{
				map[string]interface{}{
					"id":   "cat-1",
					"name": "test",
				},
			},
		})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	agentID := "agent-1"
	req := &ListCategoriesRequest{
		UserID:  "user-1",
		AgentID: &agentID,
	}

	_, err := client.ListCategories(context.Background(), req)
	if err != nil {
		t.Fatalf("ListCategories() error = %v", err)
	}
}
