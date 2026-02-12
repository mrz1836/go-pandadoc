package errors

import (
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	t.Run("error without details", func(t *testing.T) {
		err := &APIError{
			StatusCode: 404,
			Message:    "Document not found",
		}

		expected := "PandaDoc API error (status 404): Document not found"
		if err.Error() != expected {
			t.Errorf("expected '%s', got '%s'", expected, err.Error())
		}
	})

	t.Run("error with details", func(t *testing.T) {
		err := &APIError{
			StatusCode: 400,
			Message:    "Invalid request",
			Details: map[string]string{
				"field": "email",
				"error": "invalid format",
			},
		}

		result := err.Error()
		// Should contain all components
		if result == "" {
			t.Error("expected non-empty error message")
		}

		// Check for key components (exact order may vary due to map iteration)
		expectedParts := []string{"status 400", "Invalid request"}
		for _, part := range expectedParts {
			if !contains(result, part) {
				t.Errorf("expected error message to contain '%s', got '%s'", part, result)
			}
		}
	})
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrMissingAPIKey", ErrMissingAPIKey, "missing API key"},
		{"ErrMissingBaseURL", ErrMissingBaseURL, "missing base URL"},
		{"ErrRateLimitExceeded", ErrRateLimitExceeded, "API rate limit exceeded"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, tt.err.Error())
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInMiddle(s, substr)))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
