package main

import (
	"context"
	"fmt"
	"os"

	memu "github.com/NevaMind-AI/memU-sdk-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("MEMU_API_KEY")
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
	userID := "sdk_demo_user"
	agentID := "sdk_demo_agent"

	// Create client
	client, err := memu.NewClient(apiKey)
	if err != nil {
		fmt.Printf("❌ Failed to create client: %v\n", err)
		return
	}

	ctx := context.Background()

	// =========================================================
	// Step 1: Memorize a conversation
	// =========================================================
	fmt.Println("\n📝 Step 1: Memorizing conversation...")

	// Sample conversation to memorize
	conversation := []memu.ConversationMessage{
		{Role: "user", Content: "I really love Italian food, especially pasta."},
		{Role: "assistant", Content: "That's great! What's your favorite pasta dish?"},
		{Role: "user", Content: "I love carbonara! It's my absolute favorite."},
		{Role: "assistant", Content: "Carbonara is delicious! Do you cook it at home?"},
		{Role: "user", Content: "Sometimes, but I prefer dining out at authentic Italian restaurants."},
	}

	result, err := client.Memorize(ctx, &memu.MemorizeRequest{
		Conversation:      conversation,
		UserID:            userID,
		AgentID:           agentID,
		UserName:          "Demo User",
		AgentName:         "MemU Assistant",
		WaitForCompletion: false, // Don't wait, just get task ID
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
		fmt.Println("   Status: Memorization in progress...")

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
			if status.Progress != nil {
				fmt.Printf("   Progress: %.1f%%\n", *status.Progress)
			}
		}
	}

	// =========================================================
	// Step 3: List categories
	// =========================================================
	fmt.Println("\n📂 Step 3: Listing categories...")

	categories, err := client.ListCategories(ctx, &memu.ListCategoriesRequest{
		UserID: userID,
	})

	if err != nil {
		fmt.Printf("   Note: %v\n", err)
	} else {
		fmt.Printf("   Found %d categories:\n", len(categories))
		for i, cat := range categories {
			if i >= 5 {
				break
			}
			name := "Unknown"
			if cat.Name != nil {
				name = *cat.Name
			}
			summary := ""
			if cat.Summary != nil {
				summary = *cat.Summary
				if len(summary) > 50 {
					summary = summary[:50] + "..."
				}
			}
			fmt.Printf("      - %s: %s\n", name, summary)
		}
	}

	// =========================================================
	// Step 4: Retrieve memories
	// =========================================================
	fmt.Println("\n🔍 Step 4: Retrieving memories...")

	memories, err := client.Retrieve(ctx, &memu.RetrieveRequest{
		Query:   "What food does the user like?",
		UserID:  userID,
		AgentID: agentID,
	})

	if err != nil {
		fmt.Printf("   Note: %v\n", err)
	} else {
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
				summary := ""
				if item.Summary != nil {
					summary = *item.Summary
					if len(summary) > 60 {
						summary = summary[:60] + "..."
					}
				}
				fmt.Printf("      - [%s] %s\n", memType, summary)
			}
		}

		if len(memories.Categories) > 0 {
			fmt.Printf("   Related categories: %d\n", len(memories.Categories))
		}
	}

	fmt.Println("\n✨ Demo completed!")
	fmt.Println("\n============================================================")
	fmt.Println("📖 For more information, see README.md")
	fmt.Println("============================================================")
}
