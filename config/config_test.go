package config

import (
	stderrors "errors"
	"net/http"
	"testing"
	"time"

	"github.com/mrz1836/go-pandadoc/errors"
)

func TestNew(t *testing.T) {
	t.Run("creates config with defaults", func(t *testing.T) {
		cfg := New("test-api-key")

		if cfg.APIKey != "test-api-key" {
			t.Errorf("expected APIKey 'test-api-key', got '%s'", cfg.APIKey)
		}

		if cfg.BaseURL != DefaultBaseURL {
			t.Errorf("expected BaseURL '%s', got '%s'", DefaultBaseURL, cfg.BaseURL)
		}

		if cfg.Timeout != DefaultTimeout {
			t.Errorf("expected Timeout %v, got %v", DefaultTimeout, cfg.Timeout)
		}

		if cfg.UserAgent != DefaultUserAgent {
			t.Errorf("expected UserAgent '%s', got '%s'", DefaultUserAgent, cfg.UserAgent)
		}

		if cfg.Transport == nil {
			t.Error("expected Transport to be set to default")
		}
	})

	t.Run("applies custom options", func(t *testing.T) {
		customURL := "https://custom.pandadoc.com/api/"
		customTimeout := 60 * time.Second
		customUA := "custom-agent/1.0"
		customTransport := &http.Transport{}

		cfg := New("test-api-key",
			WithBaseURL(customURL),
			WithTimeout(customTimeout),
			WithUserAgent(customUA),
			WithTransport(customTransport),
		)

		if cfg.BaseURL != customURL {
			t.Errorf("expected BaseURL '%s', got '%s'", customURL, cfg.BaseURL)
		}

		if cfg.Timeout != customTimeout {
			t.Errorf("expected Timeout %v, got %v", customTimeout, cfg.Timeout)
		}

		if cfg.UserAgent != customUA {
			t.Errorf("expected UserAgent '%s', got '%s'", customUA, cfg.UserAgent)
		}

		if cfg.Transport != customTransport {
			t.Error("expected Transport to be custom transport")
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := New("test-api-key")
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("missing API key", func(t *testing.T) {
		cfg := &Config{
			BaseURL: DefaultBaseURL,
		}
		cfg.setDefaultValues()

		err := cfg.Validate()
		if !stderrors.Is(err, errors.ErrMissingAPIKey) {
			t.Errorf("expected ErrMissingAPIKey, got %v", err)
		}
	})

	t.Run("missing base URL", func(t *testing.T) {
		cfg := &Config{
			APIKey: "test-api-key",
		}

		err := cfg.Validate()
		if !stderrors.Is(err, errors.ErrMissingBaseURL) {
			t.Errorf("expected ErrMissingBaseURL, got %v", err)
		}
	})
}
