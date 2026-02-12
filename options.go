package pandadoc

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Option configures the client.
type Option func(*clientConfig) error

type clientConfig struct {
	baseURL     string
	timeout     time.Duration
	httpClient  *http.Client
	userAgent   string
	retryPolicy RetryPolicy
	apiKey      string
	accessToken string
	logger      Logger
}

// RetryPolicy controls transport-level retries.
type RetryPolicy struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	RetryOn429     bool
	RetryOn5xx     bool
}

// DefaultRetryPolicy returns a safe retry policy for API clients.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries:     2,
		InitialBackoff: 200 * time.Millisecond,
		MaxBackoff:     2 * time.Second,
		RetryOn429:     true,
		RetryOn5xx:     true,
	}
}

func (p RetryPolicy) normalize() RetryPolicy {
	if p.MaxRetries < 0 {
		p.MaxRetries = 0
	}
	if p.InitialBackoff <= 0 {
		p.InitialBackoff = 200 * time.Millisecond
	}
	if p.MaxBackoff <= 0 {
		p.MaxBackoff = 2 * time.Second
	}
	if p.MaxBackoff < p.InitialBackoff {
		p.MaxBackoff = p.InitialBackoff
	}

	return p
}

func defaultClientConfig() clientConfig {
	return clientConfig{
		baseURL:     DefaultBaseURL,
		timeout:     DefaultTimeout,
		userAgent:   DefaultUserAgent,
		retryPolicy: DefaultRetryPolicy(),
	}
}

// WithBaseURL sets a custom API base URL.
func WithBaseURL(baseURL string) Option {
	return func(cfg *clientConfig) error {
		cfg.baseURL = baseURL
		return nil
	}
}

var errInvalidTimeout = fmt.Errorf("timeout must be > 0")

// WithTimeout sets the HTTP timeout for requests.
func WithTimeout(timeout time.Duration) Option {
	return func(cfg *clientConfig) error {
		if timeout <= 0 {
			return errInvalidTimeout
		}
		cfg.timeout = timeout
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(cfg *clientConfig) error {
		if client == nil {
			return ErrNilHTTPClient
		}
		cfg.httpClient = client
		return nil
	}
}

// WithUserAgent sets a custom User-Agent.
func WithUserAgent(userAgent string) Option {
	return func(cfg *clientConfig) error {
		cfg.userAgent = strings.TrimSpace(userAgent)
		if cfg.userAgent == "" {
			cfg.userAgent = DefaultUserAgent
		}
		return nil
	}
}

// WithRetryPolicy sets a custom retry policy.
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(cfg *clientConfig) error {
		cfg.retryPolicy = policy.normalize()
		return nil
	}
}

// WithAPIKey sets API-Key auth.
func WithAPIKey(apiKey string) Option {
	return func(cfg *clientConfig) error {
		cfg.apiKey = strings.TrimSpace(apiKey)
		return nil
	}
}

// WithAccessToken sets OAuth Bearer auth.
func WithAccessToken(token string) Option {
	return func(cfg *clientConfig) error {
		cfg.accessToken = strings.TrimSpace(token)
		return nil
	}
}

// Logger defines the logging interface used by the client.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// WithLogger sets a custom logger.
func WithLogger(logger Logger) Option {
	return func(cfg *clientConfig) error {
		cfg.logger = logger
		return nil
	}
}
