// Package models provides data models for the PandaDoc API.
package models

import "time"

// CatalogItem represents a product in the PandaDoc catalog.
type CatalogItem struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	SKU         string                 `json:"sku,omitempty"`
	Price       *CatalogPrice          `json:"price,omitempty"`
	Currency    string                 `json:"currency,omitempty"`
	Quantity    *float64               `json:"quantity,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
	DateCreated time.Time              `json:"date_created,omitempty"`
	DateModified time.Time             `json:"date_modified,omitempty"`
}

// CatalogPrice represents pricing information for a catalog item.
type CatalogPrice struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency,omitempty"`
	Type     string  `json:"type,omitempty"` // e.g., "per_unit", "flat_rate"
}

// CatalogListResponse represents a paginated list of catalog items.
type CatalogListResponse struct {
	Results  []CatalogItem `json:"results"`
	Count    int           `json:"count"`
	Next     *string       `json:"next,omitempty"`
	Previous *string       `json:"previous,omitempty"`
}
