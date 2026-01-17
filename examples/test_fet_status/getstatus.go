package main

import (
	"context"
	"fmt"

	memu "github.com/NevaMind-AI/memU-sdk-go"
)

func main() {
	apiKey := "mu_dWYYNR03digW-m293RTomJnKsApamLBA9hnD7_dPtW6BpiPzxVxttgUDdrOkVyWu3d7lXTXPuSErW0x7LMpKl-6mhUh6msv1QY87FA"
	client, err := memu.NewClient(apiKey)
	if err != nil {
		fmt.Printf("❌ Failed to create client: %v\n", err)
		return
	}

	// =========================================================
	// Step 2: Check task status
	// =========================================================
	fmt.Println("\n⏳ Step 2: Checking task status...")
	ctx := context.Background()
	status, err := client.GetTaskStatus(ctx, "38OIfcmC1JAonnokV5Y4V4Dx1VB")
	if err != nil {
		fmt.Printf("   Note: Failed to get task status: %v\n", err)
	} else {
		fmt.Printf("   Task ID: %s\n", status.TaskID)
		fmt.Printf("   Status: %s\n", status.Status)
		if status.Message != nil {
			fmt.Printf("   Message: %s\n", *status.Message)
		}
	}
}
