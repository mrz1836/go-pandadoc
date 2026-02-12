// Package pandadoc provides a pure-stdlib Go SDK for the PandaDoc Public API.
package pandadoc

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default PandaDoc API base URL.
	DefaultBaseURL = "https://api.pandadoc.com/"

	// DefaultTimeout is the default HTTP timeout.
	DefaultTimeout = 30 * time.Second

	// DefaultUserAgent is the default user-agent sent to PandaDoc.
	DefaultUserAgent = "go-pandadoc/0.2.0"
)

// Client is the root PandaDoc SDK client.
type Client struct {
	baseURL     *url.URL
	httpClient  *http.Client
	userAgent   string
	retryPolicy RetryPolicy

	apiKey      string
	accessToken string

	documents            *DocumentsService
	productCatalog       *ProductCatalogService
	oauth                *OAuthService
	webhookSubscriptions *WebhookSubscriptionsService
	webhookEvents        *WebhookEventsService
}

// NewClient creates a new PandaDoc client.
//
// If no auth option is configured, only OAuth token exchange calls are available.
func NewClient(opts ...Option) (*Client, error) {
	cfg := defaultClientConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	if cfg.apiKey != "" && cfg.accessToken != "" {
		return nil, ErrMultipleAuthenticationMethods
	}

	baseURL, err := normalizeBaseURL(cfg.baseURL)
	if err != nil {
		return nil, err
	}

	httpClient := cfg.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.timeout}
	} else if cfg.timeout > 0 && httpClient.Timeout != cfg.timeout {
		clone := *httpClient
		clone.Timeout = cfg.timeout
		httpClient = &clone
	}

	client := &Client{
		baseURL:     baseURL,
		httpClient:  httpClient,
		userAgent:   cfg.userAgent,
		retryPolicy: cfg.retryPolicy.normalize(),
		apiKey:      cfg.apiKey,
		accessToken: cfg.accessToken,
	}

	client.documents = &DocumentsService{client: client}
	client.productCatalog = &ProductCatalogService{client: client}
	client.oauth = &OAuthService{client: client}
	client.webhookSubscriptions = &WebhookSubscriptionsService{client: client}
	client.webhookEvents = &WebhookEventsService{client: client}

	return client, nil
}

// NewClientWithAPIKey creates a new PandaDoc client using API-Key auth.
func NewClientWithAPIKey(apiKey string, opts ...Option) (*Client, error) {
	return NewClient(append([]Option{WithAPIKey(apiKey)}, opts...)...)
}

// NewClientWithAccessToken creates a new PandaDoc client using Bearer auth.
func NewClientWithAccessToken(token string, opts ...Option) (*Client, error) {
	return NewClient(append([]Option{WithAccessToken(token)}, opts...)...)
}

// Documents exposes document-related endpoints.
func (c *Client) Documents() *DocumentsService {
	return c.documents
}

// ProductCatalog exposes product-catalog endpoints.
func (c *Client) ProductCatalog() *ProductCatalogService {
	return c.productCatalog
}

// OAuth exposes OAuth token operations.
func (c *Client) OAuth() *OAuthService {
	return c.oauth
}

// WebhookSubscriptions exposes webhook-subscription endpoints.
func (c *Client) WebhookSubscriptions() *WebhookSubscriptionsService {
	return c.webhookSubscriptions
}

// WebhookEvents exposes webhook-event endpoints.
func (c *Client) WebhookEvents() *WebhookEventsService {
	return c.webhookEvents
}

func normalizeBaseURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		raw = DefaultBaseURL
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, ErrInvalidBaseURL
	}

	if u.Path == "" {
		u.Path = "/"
	}
	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	return u, nil
}
