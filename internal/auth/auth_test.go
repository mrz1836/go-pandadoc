package auth

import (
	"net/http"
	"testing"
)

func TestInjectAPIKey(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.example.com/test", nil)

	InjectAPIKey(req, "test-api-key-123")

	authHeader := req.Header.Get("Authorization")
	expected := "API-Key test-api-key-123"

	if authHeader != expected {
		t.Errorf("expected Authorization header '%s', got '%s'", expected, authHeader)
	}
}

func TestValidateAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		expected bool
	}{
		{"valid key", "test-api-key", true},
		{"empty key", "", false},
		{"whitespace only", "   ", true}, // Only checks if not empty
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAPIKey(tt.apiKey)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
