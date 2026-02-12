package models

// Package models contains response data structures returned by the PandaDoc API.

// Document represents a PandaDoc document.
type Document struct {
	// TODO: Define fields based on PandaDoc API documentation
}

// DocumentStatus represents the status of a document.
type DocumentStatus struct {
	// TODO: Define fields based on PandaDoc API documentation
}

// DocumentField represents a field within a document.
type DocumentField struct {
	// TODO: Define fields based on PandaDoc API documentation
}

// DocumentListResponse represents the response from listing documents.
type DocumentListResponse struct {
	Results []Document `json:"results"`
	Count   int        `json:"count"`
	Next    *string    `json:"next"`
	// TODO: Add pagination fields
}
