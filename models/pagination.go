package models

// ListOptions represents common pagination options for list operations.
type ListOptions struct {
	Page  int // Page number (1-indexed)
	Count int // Number of results per page (max 100, default 50)
}

// PaginationMeta contains pagination metadata from API responses.
type PaginationMeta struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}
