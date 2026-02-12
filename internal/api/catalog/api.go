// Package catalog provides the Product Catalog API implementation.
package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/mrz1836/go-pandadoc/commands"
	"github.com/mrz1836/go-pandadoc/errors"
	"github.com/mrz1836/go-pandadoc/internal/httpclient"
	"github.com/mrz1836/go-pandadoc/models"
)

// API provides methods for interacting with the PandaDoc Product Catalog API.
type API struct {
	client *httpclient.Client
}

// New creates a new Catalog API instance.
func New(client *httpclient.Client) *API {
	return &API{client: client}
}

// List retrieves a list of catalog items with pagination support.
// API: GET /catalog
func (a *API) List(ctx context.Context, opts *commands.ListCatalogOptions) (*models.CatalogListResponse, error) {
	path := "catalog"
	if opts != nil { //nolint:nestif // Query parameter building is straightforward
		params := url.Values{}
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Count > 0 {
			params.Set("count", strconv.Itoa(opts.Count))
		}
		if opts.Q != "" {
			params.Set("q", opts.Q)
		}
		if len(params) > 0 {
			path = path + "?" + params.Encode()
		}
	}

	resp, err := a.client.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list catalog items: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result models.CatalogListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// Get retrieves a catalog item by ID.
// API: GET /catalog/{id}
func (a *API) Get(ctx context.Context, id string) (*models.CatalogItem, error) {
	if id == "" {
		return nil, errors.ErrMissingCatalogItemID
	}

	path := fmt.Sprintf("catalog/%s", id)
	resp, err := a.client.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get catalog item: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result models.CatalogItem
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
