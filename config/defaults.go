package config

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default PandaDoc API base URL.
	DefaultBaseURL = "https://api.pandadoc.com/public/v1/"

	// DefaultTimeout is the default HTTP request timeout.
	DefaultTimeout = 30 * time.Second

	// DefaultUserAgent is the default User-Agent header.
	DefaultUserAgent = "go-pandadoc/1.0.0"
)

// setDefaultValues applies default values to unset configuration fields.
func (c *Config) setDefaultValues() {
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}
	if c.UserAgent == "" {
		c.UserAgent = DefaultUserAgent
	}
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}
}
