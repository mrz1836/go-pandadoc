package models

// CatalogItem represents an item in the product catalog.
type CatalogItem struct {
	// TODO: Define fields based on PandaDoc API documentation
}

// CatalogListResponse represents the response from listing catalog items.
type CatalogListResponse struct {
	Results []CatalogItem `json:"results"`
	Count   int           `json:"count"`
	Next    *string       `json:"next"`
	// TODO: Add pagination fields
}
