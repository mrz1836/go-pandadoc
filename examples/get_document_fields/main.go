// Example: Get document details and fields from PandaDoc
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mrz1836/go-pandadoc"
)

var (
	errMissingAPIKey = errors.New("PANDADOC_API_KEY environment variable is required")
	errUsage         = errors.New("usage: get_document_fields <document_id>")
)

func run(ctx context.Context, args []string, getenv func(string) string, out io.Writer) error {
	apiKey := getenv("PANDADOC_API_KEY")
	if apiKey == "" {
		return errMissingAPIKey
	}
	if len(args) < 2 {
		return errUsage
	}
	docID := args[1]

	opts := []pandadoc.Option{}
	if baseURL := getenv("PANDADOC_BASE_URL"); baseURL != "" {
		opts = append(opts, pandadoc.WithBaseURL(baseURL))
	}

	client, err := pandadoc.NewClientWithAPIKey(apiKey, opts...)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	details, err := client.Documents().Details(ctx, docID)
	if err != nil {
		return fmt.Errorf("get details: %w", err)
	}

	_, _ = fmt.Fprintf(out, "Document: %s\n", details.Name)
	_, _ = fmt.Fprintf(out, "Status: %s\n", details.Status)
	_, _ = fmt.Fprintf(out, "ID: %s\n", details.ID)

	if len(details.Fields) > 0 {
		_, _ = fmt.Fprintln(out, "\nFields:")
		for _, field := range details.Fields {
			_, _ = fmt.Fprintf(out, "  %s (%s): %v\n", field.Name, field.Type, field.Value)
		}
	}
	if len(details.Tokens) > 0 {
		_, _ = fmt.Fprintln(out, "\nTokens:")
		for _, token := range details.Tokens {
			_, _ = fmt.Fprintf(out, "  %s: %v\n", token.Name, token.Value)
		}
	}

	return nil
}

func main() {
	if err := run(context.Background(), os.Args, os.Getenv, os.Stdout); err != nil {
		panic(err)
	}
}
