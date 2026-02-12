package pandadoc

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) *Client {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	baseOpts := make([]Option, 0, 2+len(opts))
	baseOpts = append(baseOpts,
		WithBaseURL(srv.URL),
		WithRetryPolicy(RetryPolicy{MaxRetries: 0, InitialBackoff: 1, MaxBackoff: 1}),
	)
	baseOpts = append(baseOpts, opts...)

	client, err := NewClientWithAPIKey("test-api-key", baseOpts...)
	if err != nil {
		t.Fatalf("NewClientWithAPIKey failed: %v", err)
	}

	return client
}

func ptrBool(v bool) *bool {
	return &v
}
