// Example: List documents from PandaDoc
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mrz1836/go-pandadoc"
	"github.com/mrz1836/go-pandadoc/commands"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("PANDADOC_API_KEY")
	if apiKey == "" {
		log.Fatal("PANDADOC_API_KEY environment variable is required")
	}

	// Create client
	client, err := pandadoc.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List documents with options
	opts := &commands.ListDocumentsOptions{
		Count:  10,
		Status: "document.completed",
	}

	docs, err := client.Documents().List(context.Background(), opts)
	if err != nil {
		log.Fatalf("Failed to list documents: %v", err)
	}

	// Print results
	fmt.Printf("Found %d documents:\n", docs.Count) //nolint:forbidigo // CLI output
	for _, doc := range docs.Results {
		fmt.Printf("  - %s: %s (status: %s)\n", doc.ID, doc.Name, doc.Status) //nolint:forbidigo // CLI output
	}
}
