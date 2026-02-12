// Example: List catalog items from PandaDoc
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

	// List catalog items
	opts := &commands.ListCatalogOptions{
		Count: 25,
	}

	items, err := client.Catalog().List(context.Background(), opts)
	if err != nil {
		log.Fatalf("Failed to list catalog items: %v", err)
	}

	// Print results
	fmt.Printf("Found %d catalog items:\n", items.Count)
	for _, item := range items.Results {
		price := "N/A"
		if item.Price != nil {
			price = fmt.Sprintf("%.2f %s", item.Price.Value, item.Price.Currency)
		}
		fmt.Printf("  - %s: %s (SKU: %s, Price: %s)\n", item.ID, item.Name, item.SKU, price)
	}
}
