package pandadoc

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient_ConstructorsAndServices(t *testing.T) {
	t.Parallel()

	c, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if c.Documents() == nil || c.ProductCatalog() == nil || c.OAuth() == nil || c.WebhookSubscriptions() == nil || c.WebhookEvents() == nil {
		t.Fatal("expected all service accessors to be initialized")
	}
	if c.baseURL.String() != "https://api.pandadoc.com/" {
		t.Fatalf("unexpected default baseURL: %s", c.baseURL.String())
	}

	c2, err := NewClientWithAPIKey("k")
	if err != nil {
		t.Fatalf("NewClientWithAPIKey failed: %v", err)
	}
	if c2.apiKey != "k" {
		t.Fatalf("expected api key to be set")
	}

	c3, err := NewClientWithAccessToken("tok")
	if err != nil {
		t.Fatalf("NewClientWithAccessToken failed: %v", err)
	}
	if c3.accessToken != "tok" {
		t.Fatalf("expected access token to be set")
	}
}

func TestNewClient_ValidationAndOptions(t *testing.T) {
	t.Parallel()

	if _, err := NewClient(WithAPIKey("k"), WithAccessToken("t")); err == nil {
		t.Fatal("expected auth conflict error")
	}

	if _, err := NewClient(WithBaseURL("::bad")); err == nil {
		t.Fatal("expected invalid base url error")
	}

	if _, err := NewClient(WithHTTPClient(nil)); err == nil {
		t.Fatal("expected nil http client error")
	}

	if _, err := NewClient(WithTimeout(0)); err == nil {
		t.Fatal("expected invalid timeout error")
	}

	hc := &http.Client{Timeout: 5 * time.Second}
	c, err := NewClient(WithHTTPClient(hc), WithTimeout(9*time.Second), WithUserAgent("custom/1"), WithRetryPolicy(RetryPolicy{MaxRetries: 3, InitialBackoff: time.Second, MaxBackoff: 3 * time.Second}))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if c.httpClient == hc {
		t.Fatalf("expected client clone when timeout override is provided")
	}
	if c.httpClient.Timeout != 9*time.Second {
		t.Fatalf("unexpected timeout: %v", c.httpClient.Timeout)
	}
	if c.userAgent != "custom/1" {
		t.Fatalf("unexpected user agent: %s", c.userAgent)
	}
	if c.retryPolicy.MaxRetries != 3 {
		t.Fatalf("unexpected retry policy: %+v", c.retryPolicy)
	}
}

func TestNormalizeBaseURL(t *testing.T) {
	t.Parallel()

	u, err := normalizeBaseURL("https://api.pandadoc.com/public/v1")
	if err != nil {
		t.Fatalf("normalizeBaseURL failed: %v", err)
	}
	if u.String() != "https://api.pandadoc.com/public/v1/" {
		t.Fatalf("unexpected normalized URL: %s", u.String())
	}

	if _, err := normalizeBaseURL("  "); err != nil {
		t.Fatalf("expected default URL for blank input: %v", err)
	}

	if _, err := normalizeBaseURL("/relative"); err == nil {
		t.Fatalf("expected relative URL rejection")
	}
}

func TestRetryPolicyNormalize(t *testing.T) {
	t.Parallel()

	p := RetryPolicy{MaxRetries: -1, InitialBackoff: -1, MaxBackoff: 0}
	n := p.normalize()
	if n.MaxRetries != 0 {
		t.Fatalf("unexpected max retries: %d", n.MaxRetries)
	}
	if n.InitialBackoff <= 0 || n.MaxBackoff <= 0 {
		t.Fatalf("expected normalized positive backoffs: %+v", n)
	}

	p2 := RetryPolicy{MaxRetries: 1, InitialBackoff: 5 * time.Second, MaxBackoff: time.Second}
	n2 := p2.normalize()
	if n2.MaxBackoff != n2.InitialBackoff {
		t.Fatalf("expected max backoff to be clamped to initial backoff: %+v", n2)
	}
}

func TestNewClient_NilOptionAndUserAgentFallback(t *testing.T) {
	t.Parallel()

	c, err := NewClient(nil, WithUserAgent("   "))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if c.userAgent != DefaultUserAgent {
		t.Fatalf("expected default user-agent, got %q", c.userAgent)
	}
}
