// Package httpclient provides HTTP client utilities for making API requests.
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mrz1836/go-pandadoc/errors"
	"github.com/mrz1836/go-pandadoc/internal/auth"
)

// Client wraps an HTTP client with configuration.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	userAgent  string
}

// New creates a new HTTP client wrapper.
func New(baseURL, apiKey, userAgent string, timeout time.Duration, transport http.RoundTripper) *Client {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		baseURL:   baseURL,
		apiKey:    apiKey,
		userAgent: userAgent,
	}
}

// DoRequest performs an HTTP request with authentication and error handling.
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// Build full URL
	fullURL, err := c.buildURL(path)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Marshal body if provided
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", marshalErr)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	// Inject API Key authentication
	auth.InjectAPIKey(req, c.apiKey)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for API errors
	if err := c.checkResponse(resp); err != nil {
		_ = resp.Body.Close()
		return nil, err
	}

	return resp, nil
}

// buildURL constructs the full API URL from base URL and path.
func (c *Client) buildURL(path string) (string, error) {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	rel, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(rel).String(), nil
}

// checkResponse checks the HTTP response for errors.
func (c *Client) checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// Handle rate limiting specifically
	if resp.StatusCode == http.StatusTooManyRequests {
		return errors.ErrRateLimitExceeded
	}

	// Parse API error response
	apiErr := &errors.APIError{
		StatusCode: resp.StatusCode,
	}

	// Try to decode error body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err == nil && len(bodyBytes) > 0 {
		var errorResponse struct {
			Message string            `json:"message"`
			Details map[string]string `json:"details"`
		}
		if json.Unmarshal(bodyBytes, &errorResponse) == nil {
			apiErr.Message = errorResponse.Message
			apiErr.Details = errorResponse.Details
		} else {
			apiErr.Message = string(bodyBytes)
		}
	}

	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}

	return apiErr
}
