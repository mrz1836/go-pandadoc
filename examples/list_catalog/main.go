// Example: List catalog items from PandaDoc
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mrz1836/go-pandadoc"
)

var errMissingAPIKey = errors.New("PANDADOC_API_KEY environment variable is required")

func run(ctx context.Context, getenv func(string) string, out io.Writer) error {
	apiKey := getenv("PANDADOC_API_KEY")
	if apiKey == "" {
		return errMissingAPIKey
	}

	opts := []pandadoc.Option{}
	if baseURL := getenv("PANDADOC_BASE_URL"); baseURL != "" {
		opts = append(opts, pandadoc.WithBaseURL(baseURL))
	}

	client, err := pandadoc.NewClientWithAPIKey(apiKey, opts...)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	items, err := client.ProductCatalog().Search(ctx, &pandadoc.SearchProductCatalogItemsOptions{PerPage: 25})
	if err != nil {
		return fmt.Errorf("search catalog: %w", err)
	}

	_, _ = fmt.Fprintf(out, "Found %d catalog items:\n", items.Total)
	for _, item := range items.Items {
		price := "N/A"
		if item.Price != nil {
			price = fmt.Sprintf("%.2f %s", *item.Price, item.Currency)
		}
		_, _ = fmt.Fprintf(out, "  - %s: %s (SKU: %s, Price: %s)\n", item.UUID, item.Title, item.SKU, price)
	}
	return nil
}

func main() {
	if err := run(context.Background(), os.Getenv, os.Stdout); err != nil {
		panic(err)
	}
}
