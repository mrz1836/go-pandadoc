package catalog

import (
	"context"

	"github.com/mrz1836/go-pandadoc/commands"
	"github.com/mrz1836/go-pandadoc/models"
)

// Package catalog provides the Product Catalog API implementation.

// API provides methods for interacting with the PandaDoc Product Catalog API.
type API struct {
	// TODO: Add HTTP client and configuration
}

// New creates a new Catalog API instance.
func New() *API {
	return &API{}
}

// List retrieves a list of catalog items with pagination support.
func (a *API) List(ctx context.Context, opts *commands.ListCatalogOptions) (*models.CatalogListResponse, error) {
	// TODO: Implement list catalog items endpoint
	return nil, nil
}

// Get retrieves a catalog item by ID.
func (a *API) Get(ctx context.Context, id string) (*models.CatalogItem, error) {
	// TODO: Implement get catalog item endpoint
	return nil, nil
}
