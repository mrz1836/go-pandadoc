package pandadoc

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// productCatalogService implements ProductCatalogService.
type productCatalogService struct {
	client *Client
}

// Search searches catalog items.
func (s *productCatalogService) Search(ctx context.Context, opts *SearchProductCatalogItemsOptions) (*SearchProductCatalogItemsResponse, error) {
	query := url.Values{}
	if opts == nil {
		opts = &SearchProductCatalogItemsOptions{}
	}
	if opts.Page > 0 {
		query.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		query.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if opts.Query != "" {
		query.Set("query", opts.Query)
	}
	if opts.OrderBy != "" {
		query.Set("order_by", opts.OrderBy)
	}
	for _, v := range opts.Types {
		query.Add("types", string(v))
	}
	for _, v := range opts.BillingTypes {
		query.Add("billing_types", string(v))
	}
	for _, v := range opts.ExcludeUUIDs {
		query.Add("exclude_uuids", v)
	}
	if opts.CategoryID != "" {
		query.Set("category_id", opts.CategoryID)
	}
	if opts.NoCategory != nil {
		query.Set("no_category", strconv.FormatBool(*opts.NoCategory))
	}

	var out SearchProductCatalogItemsResponse
	err := s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v2/product-catalog/items/search",
		requireAuth: true,
		query:       query,
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Create creates a catalog item.
func (s *productCatalogService) Create(ctx context.Context, reqBody CreateProductCatalogItemRequest) (*ProductCatalogItemResponse, error) {
	if reqBody == nil {
		return nil, ErrNilRequest
	}

	var out ProductCatalogItemResponse
	err := s.client.decodeJSON(ctx, &request{
		method:      http.MethodPost,
		path:        "/public/v2/product-catalog/items",
		requireAuth: true,
		jsonBody:    reqBody,
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Get gets a catalog item by UUID.
func (s *productCatalogService) Get(ctx context.Context, itemUUID string) (*ProductCatalogItemResponse, error) {
	escapedID, err := escapePathParam(itemUUID)
	if err != nil {
		return nil, err
	}

	var out ProductCatalogItemResponse
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v2/product-catalog/items/" + escapedID,
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Update updates a catalog item by UUID.
func (s *productCatalogService) Update(ctx context.Context, itemUUID string, reqBody UpdateProductCatalogItemRequest) (*ProductCatalogItemResponse, error) {
	escapedID, err := escapePathParam(itemUUID)
	if err != nil {
		return nil, err
	}
	if reqBody == nil {
		return nil, ErrNilRequest
	}

	var out ProductCatalogItemResponse
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodPatch,
		path:        "/public/v2/product-catalog/items/" + escapedID,
		requireAuth: true,
		jsonBody:    reqBody,
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Delete deletes a catalog item by UUID.
func (s *productCatalogService) Delete(ctx context.Context, itemUUID string) error {
	escapedID, err := escapePathParam(itemUUID)
	if err != nil {
		return err
	}

	return s.client.decodeJSON(ctx, &request{
		method:      http.MethodDelete,
		path:        "/public/v2/product-catalog/items/" + escapedID,
		requireAuth: true,
	}, nil)
}
