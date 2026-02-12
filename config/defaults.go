package config

import "time"

const (
	// DefaultBaseURL is the default PandaDoc API base URL.
	DefaultBaseURL = "https://api.pandadoc.com/public/v1/"

	// DefaultTimeout is the default HTTP request timeout.
	DefaultTimeout = 30 * time.Second
)

// setDefaultValues applies default values to unset configuration fields.
func (c *Config) setDefaultValues() {
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}
}
