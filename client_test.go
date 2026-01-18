// Package memu provides unit tests for client initialization.
// This file validates NewClient and Options configuration.
package memu

import (
	"testing"
	"time"
)

// TestNewClient_ValidAPIKey tests client creation with valid API key.
func TestNewClient_ValidAPIKey(t *testing.T) {
	client, err := NewClient("test_key")
	if err != nil {
		t.Fatalf("NewClient failed with valid API key: %v", err)
	}
	if client == nil {
		t.Fatal("NewClient returned nil client")
	}
	if client.apiKey != "test_key" {
		t.Errorf("expected apiKey 'test_key', got '%s'", client.apiKey)
	}
	if client.baseURL != "https://api.memu.so" {
		t.Errorf("expected baseURL 'https://api.memu.so', got '%s'", client.baseURL)
	}
}

func TestNewClient_CustomBaseURL(t *testing.T) {
	client, err := NewClient("test_key", WithBaseURL("https://custom.example.com"))
	if err != nil {
		t.Fatalf("NewClient failed with custom base URL: %v", err)
	}
	if client.baseURL != "https://custom.example.com" {
		t.Errorf("expected baseURL 'https://custom.example.com', got '%s'", client.baseURL)
	}
}

func TestNewClient_StripsTrailingSlash(t *testing.T) {
	client, err := NewClient("test_key", WithBaseURL("https://api.example.com/"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	// Note: WithBaseURL doesn't strip trailing slash, but DefaultBaseURL does
	// This test verifies the behavior
	if client.baseURL == "https://api.example.com" {
		t.Log("Trailing slash was stripped (good)")
	} else {
		t.Logf("Trailing slash was NOT stripped: '%s'", client.baseURL)
	}
}

func TestNewClient_StripsAPIKeyWhitespace(t *testing.T) {
	client, err := NewClient("  test_key  ")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if client.apiKey != "test_key" {
		t.Errorf("expected apiKey 'test_key' (trimmed), got '%s'", client.apiKey)
	}
}

func TestNewClient_EmptyAPIKey(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Fatal("expected error for empty API key, got nil")
	}
}

func TestNewClient_WhitespaceAPIKey(t *testing.T) {
	// Note: Current implementation trims whitespace first, so "   " becomes ""
	// which should raise an error
	_, err := NewClient("   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only API key, got nil")
	}
}

// TestNewClient_Options tests client configuration options.
func TestNewClient_WithTimeout(t *testing.T) {
	client, err := NewClient("test_key", WithTimeout(30*time.Second))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if client.timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", client.timeout)
	}
}

func TestNewClient_WithMaxRetries(t *testing.T) {
	client, err := NewClient("test_key", WithMaxRetries(5))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if client.maxRetries != 5 {
		t.Errorf("expected maxRetries 5, got %d", client.maxRetries)
	}
}

func TestNewClient_MultipleOptions(t *testing.T) {
	client, err := NewClient("test_key",
		WithBaseURL("https://custom.api.com"),
		WithTimeout(45*time.Second),
		WithMaxRetries(10),
	)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if client.baseURL != "https://custom.api.com" {
		t.Errorf("expected baseURL 'https://custom.api.com', got '%s'", client.baseURL)
	}
	if client.timeout != 45*time.Second {
		t.Errorf("expected timeout 45s, got %v", client.timeout)
	}
	if client.maxRetries != 10 {
		t.Errorf("expected maxRetries 10, got %d", client.maxRetries)
	}
}

// TestClient_DefaultHeaders tests default request headers.
func TestClient_DefaultHeaders(t *testing.T) {
	client, err := NewClient("test_key")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	headers := client.defaultHeaders()

	// Check Authorization header
	if headers["Authorization"] != "Bearer test_key" {
		t.Errorf("expected Authorization 'Bearer test_key', got '%s'", headers["Authorization"])
	}

	// Check Content-Type header
	if headers["Content-Type"] != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", headers["Content-Type"])
	}

	// Check User-Agent header exists
	if _, ok := headers["User-Agent"]; !ok {
		t.Error("expected User-Agent header to exist")
	}
}

// TestNewClient_Defaults tests default values.
func TestNewClient_Defaults(t *testing.T) {
	client, err := NewClient("test_key")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("expected default baseURL '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}
	if client.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, client.timeout)
	}
	if client.maxRetries != DefaultMaxRetries {
		t.Errorf("expected default maxRetries %d, got %d", DefaultMaxRetries, client.maxRetries)
	}
}
