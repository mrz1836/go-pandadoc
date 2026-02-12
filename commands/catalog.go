// Package commands provides request DTOs for the PandaDoc API.
package commands

// ListCatalogOptions contains options for listing catalog items.
type ListCatalogOptions struct {
	// Page is the page number (1-indexed).
	Page int `json:"page,omitempty"`

	// Count is the number of results per page (max 100, default 50).
	Count int `json:"count,omitempty"`

	// Q is a search query string.
	Q string `json:"q,omitempty"`
}
