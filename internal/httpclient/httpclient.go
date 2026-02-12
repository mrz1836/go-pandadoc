package httpclient

import (
	"context"
	"net/http"
	"time"
)

// Package httpclient provides HTTP client utilities for making API requests.

// Client wraps an HTTP client with configuration.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// New creates a new HTTP client wrapper.
func New(baseURL, apiKey string, timeout time.Duration, transport http.RoundTripper) *Client {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// DoRequest performs an HTTP request with authentication.
func (c *Client) DoRequest(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	// TODO: Implement request execution with auth injection and error handling
	return nil, nil
}
