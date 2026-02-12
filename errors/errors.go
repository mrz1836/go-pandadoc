// Package errors defines error types for the PandaDoc API client.
//
//nolint:revive // Package name intentionally matches error domain
package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrMissingAPIKey is returned when no API key is provided.
	ErrMissingAPIKey = errors.New("missing API key")

	// ErrMissingBaseURL is returned when no base URL is provided.
	ErrMissingBaseURL = errors.New("missing base URL")

	// ErrRateLimitExceeded is returned when the API rate limit is exceeded (429).
	ErrRateLimitExceeded = errors.New("API rate limit exceeded")

	// ErrMissingDocumentID is returned when a document ID is required but not provided.
	ErrMissingDocumentID = errors.New("document ID is required")

	// ErrMissingCatalogItemID is returned when a catalog item ID is required but not provided.
	ErrMissingCatalogItemID = errors.New("catalog item ID is required")

	// ErrMissingUpdateData is returned when update data is required but not provided.
	ErrMissingUpdateData = errors.New("update data is required")
)

// APIError represents an error response from the PandaDoc API.
type APIError struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Details    map[string]string `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("PandaDoc API error (status %d): %s - %v", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("PandaDoc API error (status %d): %s", e.StatusCode, e.Message)
}
