package auth

import "net/http"

// Package auth provides authentication utilities for the PandaDoc API.

// InjectAPIKey adds the API Key authentication header to the request.
// Format: Authorization: API-Key {apiKey}
func InjectAPIKey(req *http.Request, apiKey string) {
	req.Header.Set("Authorization", "API-Key "+apiKey)
}

// ValidateAPIKey validates that an API key is not empty.
func ValidateAPIKey(apiKey string) bool {
	return apiKey != ""
}
