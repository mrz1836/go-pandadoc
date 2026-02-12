package pandadoc

import (
	stderrors "errors"
	"fmt"
	"net/http"
)

var (
	// ErrInvalidBaseURL indicates the configured API base URL is invalid.
	ErrInvalidBaseURL = stderrors.New("invalid base URL")

	// ErrNilHTTPClient indicates a nil HTTP client was provided.
	ErrNilHTTPClient = stderrors.New("http client cannot be nil")

	// ErrMissingAuthentication indicates no auth method is configured for an authenticated call.
	ErrMissingAuthentication = stderrors.New("missing authentication credentials")

	// ErrMultipleAuthenticationMethods indicates both API-Key and Bearer auth were configured.
	ErrMultipleAuthenticationMethods = stderrors.New("only one authentication method can be configured")

	// ErrEmptyPathParameter indicates a required path parameter was empty.
	ErrEmptyPathParameter = stderrors.New("path parameter cannot be empty")

	// ErrNilRequest indicates a required request payload is nil.
	ErrNilRequest = stderrors.New("request payload cannot be nil")

	// ErrNilFileReader indicates an upload request has no file reader.
	ErrNilFileReader = stderrors.New("file reader is required")
)

// APIError represents a non-2xx response from PandaDoc.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    any

	RequestID  string
	RetryAfter string
	RawBody    string
	Headers    http.Header
}

// Error implements error.
func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code != "" {
		return fmt.Sprintf("pandadoc API error: status=%d code=%s message=%s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("pandadoc API error: status=%d message=%s", e.StatusCode, e.Message)
}

// IsRateLimited returns true if err is a 429 API error.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if !stderrors.As(err, &apiErr) {
		return false
	}
	return apiErr.StatusCode == http.StatusTooManyRequests
}

// IsUnauthorized returns true if err is a 401 API error.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if !stderrors.As(err, &apiErr) {
		return false
	}
	return apiErr.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if err is a 403 API error.
func IsForbidden(err error) bool {
	var apiErr *APIError
	if !stderrors.As(err, &apiErr) {
		return false
	}
	return apiErr.StatusCode == http.StatusForbidden
}

// IsNotFound returns true if err is a 404 API error.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if !stderrors.As(err, &apiErr) {
		return false
	}
	return apiErr.StatusCode == http.StatusNotFound
}
