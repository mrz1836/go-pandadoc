package config

import (
	"net/http"
	"time"
)

// Option is a functional option for configuring the Config.
type Option func(*Config)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// WithTimeout sets a custom timeout for HTTP requests.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(userAgent string) Option {
	return func(c *Config) {
		c.UserAgent = userAgent
	}
}

// WithTransport sets a custom HTTP transport.
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Config) {
		c.Transport = transport
	}
}
