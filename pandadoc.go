// Package pandadoc provides an unofficial Go SDK for the PandaDoc API.
//
// This package re-exports commonly used types for convenience.
package pandadoc

import (
	"github.com/mrz1836/go-pandadoc/config"
)

// Re-export config options for convenience
var (
	// WithBaseURL sets a custom base URL for the API.
	WithBaseURL = config.WithBaseURL

	// WithTimeout sets the HTTP request timeout.
	WithTimeout = config.WithTimeout

	// WithUserAgent sets a custom User-Agent header.
	WithUserAgent = config.WithUserAgent

	// WithTransport sets a custom HTTP transport.
	WithTransport = config.WithTransport
)
