package config

import (
	"net/http"
	"time"
)

// Package config provides configuration management for the PandaDoc API client.

// Config holds the configuration for the PandaDoc API client.
type Config struct {
	APIKey    string            // PandaDoc API Key
	BaseURL   string            // Base URL for the PandaDoc API (default: https://api.pandadoc.com/public/v1/)
	Timeout   time.Duration     // HTTP request timeout
	Transport http.RoundTripper // Custom HTTP transport (optional)
}

// New creates a new Config with the provided options.
func New(apiKey string, options ...Option) *Config {
	cfg := &Config{
		APIKey: apiKey,
	}

	for _, opt := range options {
		opt(cfg)
	}

	cfg.setDefaultValues()
	return cfg
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}
	if c.BaseURL == "" {
		return ErrMissingBaseURL
	}
	return nil
}
