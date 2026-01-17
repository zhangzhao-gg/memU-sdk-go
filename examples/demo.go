/**
 * [INPUT]: 依赖 context, fmt, os 标准库; 依赖 github.com/NevaMind-AI/memU-sdk-go 的 Client 和所有 API 方法
 * [OUTPUT]: 可执行的 main 程序，演示 SDK 完整用法
 * [POS]: examples/ 的唯一示例文件，展示四大 API 方法的使用和错误处理
 * [PROTOCOL]: 变更时更新此头部，然后检查 CLAUDE.md
 */
package main

import (
	"context"
	"fmt"

	memu "github.com/NevaMind-AI/memU-sdk-go"
)

func main() {
	// Get API key from environment variable
	apiKey := "mu_dWYYNR03digW-m293RTomJnKsApamLBA9hnD7_dPtW6BpiPzxVxttgUDdrOkVyWu3d7lXTXPuSErW0x7LMpKl-6mhUh6msv1QY87FA"
	if apiKey == "" {
		fmt.Println("❌ MEMU_API_KEY environment variable not set")
		fmt.Println("   Please get your API key from https://memu.so")
		fmt.Println("   Then run: export MEMU_API_KEY=your_api_key")
		return
	}

	fmt.Println("============================================================")
	fmt.Println("🚀 MemU SDK Demo - Cloud API")
	fmt.Println("============================================================")

	// Demo user and agent IDs
	userID := "user_123"
	agentID := "agent456"

	// Create client
	client, err := memu.NewClient(apiKey)
	if err != nil {
		fmt.Printf("❌ Failed to create client: %v\n", err)
		return
	}

	ctx := context.Background()

	// =========================================================
	// Step 1: Memorize a conversation (with optional metadata)
	// =========================================================
	fmt.Println("\n📝 Step 1: Memorizing conversation...")

	// Optional: Add speaker names and timestamps
	userName := "John"
	assistantName := "Coach"
	time1 := "2024-01-15T10:30:00Z"
	time2 := "2024-01-15T10:30:15Z"
	time3 := "2024-01-15T10:31:00Z"

	// Sample conversation to memorize
	conversation := []memu.ConversationMessage{
		{
			Role:      "user",
			Content:   "I love playing tennis on weekends",
			Name:      &userName,
			CreatedAt: &time1,
		},
		{
			Role:      "assistant",
			Content:   "That's great! Tennis is an excellent way to stay active.",
			Name:      &assistantName,
			CreatedAt: &time2,
		},
		{
			Role:      "user",
			Content:   "I usually play at the local club every Saturday morning.",
			Name:      &userName,
			CreatedAt: &time3,
		},
	}

	sessionDate := "2024-01-15T10:30:00Z"
	result, err := client.Memorize(ctx, &memu.MemorizeRequest{
		Conversation: conversation,
		UserID:       userID,
		AgentID:      agentID,
		UserName:     "John Doe",
		AgentName:    "Tennis Coach AI",
		SessionDate:  &sessionDate,
	})

	if err != nil {
		if authErr, ok := err.(*memu.AuthenticationError); ok {
			fmt.Printf("❌ Authentication failed: %v\n", authErr)
			return
		}
		fmt.Printf("❌ Failed to memorize: %v\n", err)
		return
	}

	if result.TaskID != nil {
		fmt.Printf("   ✅ Task submitted: %s\n", *result.TaskID)
		if result.Status != nil {
			fmt.Printf("   Status: %s\n", *result.Status)
		}
		if result.Message != nil {
			fmt.Printf("   Message: %s\n", *result.Message)
		}

		// =========================================================
		// Step 2: Check task status
		// =========================================================
		fmt.Println("\n⏳ Step 2: Checking task status...")

		status, err := client.GetTaskStatus(ctx, *result.TaskID)
		if err != nil {
			fmt.Printf("   Note: Failed to get task status: %v\n", err)
		} else {
			fmt.Printf("   Task ID: %s\n", status.TaskID)
			fmt.Printf("   Status: %s\n", status.Status)
		}
	}

	// =========================================================
	// Step 3: List categories
	// =========================================================
	fmt.Println("\n📂 Step 3: Listing categories...")

	categories, err := client.ListCategories(ctx, &memu.ListCategoriesRequest{
		UserID:  userID,
		AgentID: &agentID,
	})

	if err != nil {
		fmt.Printf("   Note: %v\n", err)
	} else {
		fmt.Printf("   Found %d categories:\n", len(categories))
		for i, cat := range categories {
			if i >= 5 {
				break
			}

			// Print category name
			name := "Unknown"
			if cat.Name != nil {
				name = *cat.Name
			}
			fmt.Printf("\n      Category %d: %s\n", i+1, name)

			// Print description
			if cat.Description != nil && *cat.Description != "" {
				desc := *cat.Description
				if len(desc) > 80 {
					desc = desc[:80] + "..."
				}
				fmt.Printf("         Description: %s\n", desc)
			}

			// Print summary
			if cat.Summary != nil && *cat.Summary != "" {
				summary := *cat.Summary
				if len(summary) > 100 {
					summary = summary[:100] + "..."
				}
				fmt.Printf("         Summary: %s\n", summary)
			}

			// Print user_id and agent_id
			if cat.UserID != nil {
				fmt.Printf("         User ID: %s\n", *cat.UserID)
			}
			if cat.AgentID != nil {
				fmt.Printf("         Agent ID: %s\n", *cat.AgentID)
			}
		}
	}

	// =========================================================
	// Step 4a: Retrieve memories with string query
	// =========================================================
	fmt.Println("\n🔍 Step 4a: Retrieving memories (string query)...")

	memories, err := client.Retrieve(ctx, &memu.RetrieveRequest{
		Query:   "What are the user's hobbies and interests?",
		UserID:  userID,
		AgentID: agentID,
	})

	if err != nil {
		fmt.Printf("   Note: %v\n", err)
	} else {
		if memories.RewrittenQuery != nil {
			fmt.Printf("   Rewritten Query: %s\n", *memories.RewrittenQuery)
		}
		fmt.Printf("   Found %d memory items\n", len(memories.Items))
		if len(memories.Items) > 0 {
			for i, item := range memories.Items {
				if i >= 5 {
					break
				}
				memType := "unknown"
				if item.MemoryType != nil {
					memType = *item.MemoryType
				}
				content := ""
				if item.Content != nil {
					content = *item.Content
					if len(content) > 100 {
						content = content[:100] + "..."
					}
				}
				fmt.Printf("      - [%s] %s\n", memType, content)
			}
		}

		if len(memories.Categories) > 0 {
			fmt.Printf("   Related categories: %d\n", len(memories.Categories))
		}
	}

	// =========================================================
	// Step 4b: Retrieve memories with conversation array query
	// =========================================================
	fmt.Println("\n🔍 Step 4b: Retrieving memories (conversation array query)...")

	conversationQuery := []memu.ConversationMessage{
		{
			Role:    "user",
			Content: "What sports do I usually enjoy?",
		},
	}

	memories2, err := client.Retrieve(ctx, &memu.RetrieveRequest{
		Query:   conversationQuery,
		UserID:  userID,
		AgentID: agentID,
	})

	if err != nil {
		fmt.Printf("   Note: %v\n", err)
	} else {
		if memories2.RewrittenQuery != nil {
			fmt.Printf("   Rewritten Query: %s\n", *memories2.RewrittenQuery)
		}
		fmt.Printf("   Found %d memory items\n", len(memories2.Items))
		if len(memories2.Items) > 0 {
			for i, item := range memories2.Items {
				if i >= 5 {
					break
				}
				memType := "unknown"
				if item.MemoryType != nil {
					memType = *item.MemoryType
				}
				content := ""
				if item.Content != nil {
					content = *item.Content
					if len(content) > 100 {
						content = content[:100] + "..."
					}
				}
				fmt.Printf("      - [%s] %s\n", memType, content)
			}
		}

		if len(memories2.Categories) > 0 {
			fmt.Printf("   Related categories: %d\n", len(memories2.Categories))
		}
	}

	fmt.Println("\n✨ Demo completed!")
	fmt.Println("\n============================================================")
	fmt.Println("📖 For more information, see README.md")
	fmt.Println("============================================================")
}
