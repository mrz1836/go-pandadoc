package documents

import (
	"context"

	"github.com/mrz1836/go-pandadoc/commands"
	"github.com/mrz1836/go-pandadoc/models"
)

// Package documents provides the Documents API implementation.

// API provides methods for interacting with the PandaDoc Documents API.
type API struct {
	// TODO: Add HTTP client and configuration
}

// New creates a new Documents API instance.
func New() *API {
	return &API{}
}

// List retrieves a list of documents with pagination support.
func (a *API) List(ctx context.Context, opts *commands.ListDocumentsOptions) (*models.DocumentListResponse, error) {
	// TODO: Implement list documents endpoint
	return nil, nil
}

// Get retrieves a document by ID.
func (a *API) Get(ctx context.Context, id string) (*models.Document, error) {
	// TODO: Implement get document endpoint
	return nil, nil
}

// GetStatus retrieves the status of a document by ID.
func (a *API) GetStatus(ctx context.Context, id string) (*models.DocumentStatus, error) {
	// TODO: Implement get document status endpoint
	return nil, nil
}

// GetFields retrieves the fields of a document by ID.
func (a *API) GetFields(ctx context.Context, id string) ([]models.DocumentField, error) {
	// TODO: Implement get document fields endpoint
	return nil, nil
}

// Update updates a document by ID.
func (a *API) Update(ctx context.Context, id string, update *commands.UpdateDocument) (*models.Document, error) {
	// TODO: Implement update document endpoint
	return nil, nil
}
