// Package main demonstrates the complete usage of the MemU SDK.
// This example shows how to use all four major API methods with proper error handling.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	memu "github.com/NevaMind-AI/memU-sdk-go"
)

func main() {
	// Get API key from environment variable
	apiKey := "mu_xxx"
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
	userID := "user_test"
	agentID := "agent_test"

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
		// Step 2: Check task status (poll until SUCCESS)
		// =========================================================
		fmt.Println("\n⏳ Step 2: Checking task status...")

		maxWait := 60 * time.Second
		pollInterval := 2 * time.Second
		startTime := time.Now()
		completed := false

		for time.Since(startTime) < maxWait {
			status, err := client.GetTaskStatus(ctx, *result.TaskID)
			if err != nil {
				fmt.Printf("   ❌ Failed to get task status: %v\n", err)
				break
			}

			fmt.Printf("   Status: %s", status.Status)
			if status.DetailInfo != "" {
				fmt.Printf(" - %s", status.DetailInfo)
			}
			fmt.Println()

			// Check if task is completed
			if status.Status == memu.TaskStatusSuccess {
				completed = true
				fmt.Printf("\n   ✅ Task completed successfully in %.1fs\n", time.Since(startTime).Seconds())
				fmt.Printf("   Task ID: %s\n", status.TaskID)
				break
			} else if status.Status == memu.TaskStatusFailed {
				fmt.Printf("\n   ❌ Task failed: %s\n", status.DetailInfo)
				break
			} else if status.Status == memu.TaskStatusCompleted {
				// COMPLETED is also a success state
				completed = true
				fmt.Printf("\n   ✅ Task completed in %.1fs\n", time.Since(startTime).Seconds())
				fmt.Printf("   Task ID: %s\n", status.TaskID)
				break
			}

			// Wait before next poll
			time.Sleep(pollInterval)
		}

		if !completed {
			fmt.Printf("\n   ⚠️  Timeout: Task did not complete within %.0fs\n", maxWait.Seconds())
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
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("\n   📂 Categories: %d found\n", len(categories))
		for i, cat := range categories {
			if i >= 5 {
				fmt.Printf("      ... and %d more categories\n", len(categories)-5)
				break
			}

			// Separator line
			fmt.Printf("      ─────────────────────────────────────\n")

			// Name
			name := "(unnamed)"
			if cat.Name != nil {
				name = *cat.Name
			}
			fmt.Printf("      📁 %s\n", name)

			// Description
			if cat.Description != nil && *cat.Description != "" {
				desc := *cat.Description
				if len(desc) > 80 {
					desc = desc[:80] + "..."
				}
				fmt.Printf("         📝 %s\n", desc)
			}

			// Summary
			if cat.Summary != nil && *cat.Summary != "" {
				summary := *cat.Summary
				if len(summary) > 100 {
					summary = summary[:100] + "..."
				}
				fmt.Printf("         📄 %s\n", summary)
			}

			// ID information (displayed on same line)
			ids := ""
			if cat.UserID != nil {
				ids += fmt.Sprintf("user:%s ", *cat.UserID)
			}
			if cat.AgentID != nil {
				ids += fmt.Sprintf("agent:%s", *cat.AgentID)
			}
			if ids != "" {
				fmt.Printf("         🔑 %s\n", ids)
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
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		// Display rewritten query
		if memories.RewrittenQuery != nil {
			fmt.Printf("   📝 Rewritten Query: %s\n", *memories.RewrittenQuery)
		}

		// Display memory items
		fmt.Printf("\n   📦 Memory Items: %d found\n", len(memories.Items))
		for i, item := range memories.Items {
			if i >= 5 {
				fmt.Printf("      ... and %d more items\n", len(memories.Items)-5)
				break
			}
			memType := "unknown"
			if item.MemoryType != nil {
				memType = *item.MemoryType
			}
			content := "(empty)"
			if item.Content != nil {
				content = *item.Content
				if len(content) > 80 {
					content = content[:80] + "..."
				}
			}
			fmt.Printf("      %d. [%s] %s\n", i+1, memType, content)
		}

		// Display related categories
		if len(memories.Categories) > 0 {
			fmt.Printf("\n   📂 Related Categories: %d found\n", len(memories.Categories))
			for i, cat := range memories.Categories {
				if i >= 3 {
					break
				}
				name := "(unnamed)"
				if cat.Name != nil {
					name = *cat.Name
				}
				summary := ""
				if cat.Summary != nil {
					summary = *cat.Summary
					if len(summary) > 60 {
						summary = summary[:60] + "..."
					}
				}
				fmt.Printf("      - %s\n", name)
				if summary != "" {
					fmt.Printf("        Summary: %s\n", summary)
				}
			}
		}

		// Display resources
		if len(memories.Resources) > 0 {
			fmt.Printf("\n   🖼️  Resources: %d found\n", len(memories.Resources))
			for i, res := range memories.Resources {
				if i >= 3 {
					break
				}
				modality := "unknown"
				if res.Modality != nil {
					modality = *res.Modality
				}
				url := "(no url)"
				if res.ResourceURL != nil {
					url = *res.ResourceURL
					if len(url) > 50 {
						url = url[:50] + "..."
					}
				}
				fmt.Printf("      - [%s] %s\n", modality, url)
				if res.Caption != nil && *res.Caption != "" {
					fmt.Printf("        Caption: %s\n", *res.Caption)
				}
				if len(res.Content) > 0 {
					fmt.Printf("        Content: %v\n", res.Content)
				}
				if len(res.Metadata) > 0 {
					fmt.Printf("        Metadata: %v\n", res.Metadata)
				}
			}
		}
	}

	// // =========================================================
	// // Step 4b: Retrieve memories with conversation array query
	// // =========================================================
	fmt.Println("\n🔍 Step 4b: Retrieving memories (conversation array query)...")

	// Simulate multi-turn conversation context query
	conversationQuery := []memu.ConversationMessage{
		{
			Role:    "user",
			Content: "I want to be more active this year",
		},
		{
			Role:    "assistant",
			Content: "That's a great goal! What kind of activities interest you?",
		},
		{
			Role:    "user",
			Content: "What sports do I usually enjoy?",
		},
	}

	// Print the actual query structure being sent
	queryJSON, _ := json.MarshalIndent(conversationQuery, "   ", "  ")
	fmt.Printf("\n   📤 Sending query:\n   %s\n", string(queryJSON))

	// Build complete request payload for debugging
	requestPayload := map[string]interface{}{
		"user_id":  userID,
		"agent_id": agentID,
		"query":    conversationQuery,
	}
	payloadJSON, _ := json.MarshalIndent(requestPayload, "   ", "  ")
	fmt.Printf("\n   📤 Full request payload:\n   %s\n", string(payloadJSON))

	memories2, err := client.Retrieve(ctx, &memu.RetrieveRequest{
		Query:   conversationQuery,
		UserID:  userID,
		AgentID: agentID,
	})

	if err != nil {
		fmt.Printf("\n   ❌ Error: %v\n", err)
	} else {
		fmt.Println("\n   ✅ Success! Response received:")
		fmt.Println("   ═══════════════════════════════════════════════════════════")

		// 1. Display rewritten query
		if memories2.RewrittenQuery != nil {
			fmt.Printf("\n   📝 Rewritten Query:\n")
			fmt.Printf("      \"%s\"\n", *memories2.RewrittenQuery)
		}

		// 2. Display categories
		if len(memories2.Categories) > 0 {
			fmt.Printf("\n   📂 Categories (%d):\n", len(memories2.Categories))
			for i, cat := range memories2.Categories {
				fmt.Printf("\n      ┌─ Category %d\n", i+1)

				if cat.Name != nil {
					fmt.Printf("      │  Name: %s\n", *cat.Name)
				}

				if cat.Description != nil && *cat.Description != "" {
					fmt.Printf("      │  Description: %s\n", *cat.Description)
				}

				if cat.Summary != nil && *cat.Summary != "" {
					summary := *cat.Summary
					if len(summary) > 100 {
						summary = summary[:100] + "..."
					}
					fmt.Printf("      │  Summary: %s\n", summary)
				}

				fmt.Printf("      └─\n")
			}
		} else {
			fmt.Printf("\n   📂 Categories: (none)\n")
		}

		// 3. Display memory items
		if len(memories2.Items) > 0 {
			fmt.Printf("\n   📦 Memory Items (%d):\n", len(memories2.Items))
			for i, item := range memories2.Items {
				fmt.Printf("\n      ┌─ Item %d\n", i+1)

				if item.MemoryType != nil {
					fmt.Printf("      │  Type: %s\n", *item.MemoryType)
				}

				if item.Content != nil {
					content := *item.Content
					if len(content) > 150 {
						content = content[:150] + "..."
					}
					fmt.Printf("      │  Content: %s\n", content)
				}

				fmt.Printf("      └─\n")
			}
		} else {
			fmt.Printf("\n   📦 Memory Items: (none)\n")
		}

		// 4. Display resources
		if len(memories2.Resources) > 0 {
			fmt.Printf("\n   🖼️  Resources (%d):\n", len(memories2.Resources))
			for i, res := range memories2.Resources {
				fmt.Printf("\n      ┌─ Resource %d\n", i+1)

				if res.Modality != nil {
					fmt.Printf("      │  Modality: %s\n", *res.Modality)
				}

				if res.ResourceURL != nil {
					fmt.Printf("      │  URL: %s\n", *res.ResourceURL)
				}

				if res.Caption != nil && *res.Caption != "" {
					fmt.Printf("      │  Caption: %s\n", *res.Caption)
				}

				if len(res.Content) > 0 {
					contentJSON, _ := json.MarshalIndent(res.Content, "      │  ", "  ")
					fmt.Printf("      │  Content: %s\n", string(contentJSON))
				}

				if len(res.Metadata) > 0 {
					metadataJSON, _ := json.MarshalIndent(res.Metadata, "      │  ", "  ")
					fmt.Printf("      │  Metadata: %s\n", string(metadataJSON))
				}

				fmt.Printf("      └─\n")
			}
		} else {
			fmt.Printf("\n   🖼️  Resources: (none)\n")
		}

		fmt.Println("\n   ═══════════════════════════════════════════════════════════")
	}

	fmt.Println("\n✨ Demo completed!")
	fmt.Println("\n============================================================")
	fmt.Println("📖 For more information, see README.md")
	fmt.Println("============================================================")
}
