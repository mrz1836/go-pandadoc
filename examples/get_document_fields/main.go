// Example: Get document details and fields from PandaDoc
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mrz1836/go-pandadoc"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("PANDADOC_API_KEY")
	if apiKey == "" {
		log.Fatal("PANDADOC_API_KEY environment variable is required")
	}

	// Get document ID from args
	if len(os.Args) < 2 {
		log.Fatal("Usage: get_document_fields <document_id>")
	}
	docID := os.Args[1]

	// Create client
	client, err := pandadoc.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get document details
	details, err := client.Documents().GetDetails(context.Background(), docID)
	if err != nil {
		log.Fatalf("Failed to get document details: %v", err)
	}

	// Print document info
	fmt.Printf("Document: %s\n", details.Name) //nolint:forbidigo // CLI output
	fmt.Printf("Status: %s\n", details.Status) //nolint:forbidigo // CLI output
	fmt.Printf("ID: %s\n", details.ID)         //nolint:forbidigo // CLI output

	// Print fields
	if len(details.Fields) > 0 {
		fmt.Println("\nFields:") //nolint:forbidigo // CLI output
		for name, value := range details.Fields {
			fmt.Printf("  %s: %v\n", name, value) //nolint:forbidigo // CLI output
		}
	}

	// Print tokens
	if len(details.Tokens) > 0 {
		fmt.Println("\nTokens:") //nolint:forbidigo // CLI output
		for _, token := range details.Tokens {
			fmt.Printf("  %s: %s\n", token.Name, token.Value) //nolint:forbidigo // CLI output
		}
	}
}
