// Example: List documents from PandaDoc
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

	status := pandadoc.DocumentStatusCompleted
	docs, err := client.Documents().List(ctx, &pandadoc.ListDocumentsOptions{Count: 10, Status: &status})
	if err != nil {
		return fmt.Errorf("list documents: %w", err)
	}

	_, _ = fmt.Fprintf(out, "Found %d documents:\n", len(docs.Results))
	for _, doc := range docs.Results {
		_, _ = fmt.Fprintf(out, "  - %s: %s (status: %s)\n", doc.ID, doc.Name, doc.Status)
	}
	return nil
}

func main() {
	if err := run(context.Background(), os.Getenv, os.Stdout); err != nil {
		panic(err)
	}
}
