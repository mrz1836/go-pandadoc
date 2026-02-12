package pandadoc

import (
	"errors"
	"io"
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

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newRoundTripperClient(t *testing.T, rt roundTripperFunc, opts ...Option) *Client {
	t.Helper()

	httpClient := &http.Client{Transport: rt}
	baseOpts := make([]Option, 0, 3+len(opts))
	baseOpts = append(baseOpts,
		WithBaseURL("https://api.example.com"),
		WithHTTPClient(httpClient),
		WithRetryPolicy(RetryPolicy{MaxRetries: 0, InitialBackoff: 1, MaxBackoff: 1}),
	)
	baseOpts = append(baseOpts, opts...)

	client, err := NewClientWithAPIKey("test-api-key", baseOpts...)
	if err != nil {
		t.Fatalf("NewClientWithAPIKey failed: %v", err)
	}

	return client
}

type errorReader struct {
	err error
}

func (r *errorReader) Read([]byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return 0, io.EOF
}

func newErrorReader() *errorReader {
	return &errorReader{err: errors.New("forced read error")} //nolint:err113 // static error for testing
}
