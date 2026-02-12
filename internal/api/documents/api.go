// Package documents provides the Documents API implementation.
package documents

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/mrz1836/go-pandadoc/commands"
	"github.com/mrz1836/go-pandadoc/errors"
	"github.com/mrz1836/go-pandadoc/internal/httpclient"
	"github.com/mrz1836/go-pandadoc/models"
)

// API provides methods for interacting with the PandaDoc Documents API.
type API struct {
	client *httpclient.Client
}

// New creates a new Documents API instance.
func New(client *httpclient.Client) *API {
	return &API{client: client}
}

// List retrieves a list of documents with pagination support.
// API: GET /documents
func (a *API) List(ctx context.Context, opts *commands.ListDocumentsOptions) (*models.DocumentListResponse, error) {
	path := "documents"
	if opts != nil { //nolint:nestif // Query parameter building is straightforward
		params := url.Values{}
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Count > 0 {
			params.Set("count", strconv.Itoa(opts.Count))
		}
		if opts.Status != "" {
			params.Set("status", opts.Status)
		}
		if opts.OrderBy != "" {
			params.Set("order_by", opts.OrderBy)
		}
		if len(params) > 0 {
			path = path + "?" + params.Encode()
		}
	}

	resp, err := a.client.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result models.DocumentListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Get retrieves a document by ID.
// API: GET /documents/{id}
func (a *API) Get(ctx context.Context, id string) (*models.Document, error) {
	if id == "" {
		return nil, errors.ErrMissingDocumentID
	}

	path := fmt.Sprintf("documents/%s", id)
	resp, err := a.client.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result models.Document
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetStatus retrieves the status of a document by ID.
// API: GET /documents/{id}/status
func (a *API) GetStatus(ctx context.Context, id string) (*models.DocumentStatus, error) {
	if id == "" {
		return nil, errors.ErrMissingDocumentID
	}

	path := fmt.Sprintf("documents/%s/status", id)
	resp, err := a.client.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get document status: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result models.DocumentStatus
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetDetails retrieves the details/fields of a document by ID.
// API: GET /documents/{id}/details
func (a *API) GetDetails(ctx context.Context, id string) (*models.DocumentDetails, error) {
	if id == "" {
		return nil, errors.ErrMissingDocumentID
	}

	path := fmt.Sprintf("documents/%s/details", id)
	resp, err := a.client.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get document details: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result models.DocumentDetails
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Update updates a document by ID.
// API: PATCH /documents/{id}
func (a *API) Update(ctx context.Context, id string, update *commands.UpdateDocument) (*models.Document, error) {
	if id == "" {
		return nil, errors.ErrMissingDocumentID
	}
	if update == nil {
		return nil, errors.ErrMissingUpdateData
	}

	path := fmt.Sprintf("documents/%s", id)
	resp, err := a.client.DoRequest(ctx, "PATCH", path, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result models.Document
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
