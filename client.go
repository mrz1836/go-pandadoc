// Package pandadoc provides an unofficial Go SDK for the PandaDoc API.
//
// Example usage:
//
//	client, err := pandadoc.NewClient("your-api-key")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// List documents
//	docs, err := client.Documents().List(ctx, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get catalog items
//	items, err := client.Catalog().List(ctx, nil)
package pandadoc

import (
	"github.com/mrz1836/go-pandadoc/config"
	"github.com/mrz1836/go-pandadoc/internal/api/catalog"
	"github.com/mrz1836/go-pandadoc/internal/api/documents"
	"github.com/mrz1836/go-pandadoc/internal/httpclient"
)

// Client is the main PandaDoc API client.
type Client struct {
	config    *config.Config
	http      *httpclient.Client
	documents *documents.API
	catalog   *catalog.API
}

// NewClient creates a new PandaDoc API client with the given API key.
// Additional options can be provided to customize the client behavior.
//
// Example:
//
//	// Basic usage
//	client, err := pandadoc.NewClient("your-api-key")
//
//	// With options
//	client, err := pandadoc.NewClient("your-api-key",
//	    pandadoc.WithTimeout(60 * time.Second),
//	    pandadoc.WithUserAgent("my-app/1.0"),
//	)
func NewClient(apiKey string, options ...config.Option) (*Client, error) {
	cfg := config.New(apiKey, options...)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	httpClient := httpclient.New(
		cfg.BaseURL,
		cfg.APIKey,
		cfg.UserAgent,
		cfg.Timeout,
		cfg.Transport,
	)

	return &Client{
		config:    cfg,
		http:      httpClient,
		documents: documents.New(httpClient),
		catalog:   catalog.New(httpClient),
	}, nil
}

// Documents returns the Documents API for managing PandaDoc documents.
func (c *Client) Documents() *documents.API {
	return c.documents
}

// Catalog returns the Catalog API for managing product catalog items.
func (c *Client) Catalog() *catalog.API {
	return c.catalog
}
