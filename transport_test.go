package pandadoc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestJoinPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		base string
		rel  string
		want string
	}{
		{"/", "documents", "/documents"},
		{"/public/v1/", "documents", "/public/v1/documents"},
		{"/public/v1", "documents", "/public/v1/documents"},
		{"public/v1", "/documents", "/public/v1/documents"},
		{"", "documents", "/documents"},
	}

	for _, tc := range tests {
		if got := joinPaths(tc.base, tc.rel); got != tc.want {
			t.Fatalf("joinPaths(%q,%q)=%q want %q", tc.base, tc.rel, got, tc.want)
		}
	}
}

func TestBuildURL_MergesQuery(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAPIKey("k", WithBaseURL("https://api.example.com/public/v1"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	q := url.Values{}
	q.Set("b", "2")
	full, err := c.buildURL("documents?a=1", q)
	if err != nil {
		t.Fatalf("buildURL failed: %v", err)
	}
	u, _ := url.Parse(full)
	if u.Path != "/public/v1/documents" {
		t.Fatalf("unexpected path: %s", u.Path)
	}
	if u.Query().Get("a") != "1" || u.Query().Get("b") != "2" {
		t.Fatalf("unexpected query: %s", u.RawQuery)
	}
}

func TestEncodeRequestBody(t *testing.T) {
	t.Parallel()

	payload, ct, err := encodeRequestBody(&request{jsonBody: map[string]string{"a": "b"}})
	if err != nil || ct != "application/json" || !strings.Contains(string(payload), "\"a\":\"b\"") {
		t.Fatalf("json body encoding failed: ct=%s err=%v payload=%s", ct, err, payload)
	}

	payload, ct, err = encodeRequestBody(&request{formBody: url.Values{"x": []string{"1"}}})
	if err != nil || ct != "application/x-www-form-urlencoded" || string(payload) != "x=1" {
		t.Fatalf("form encoding failed: ct=%s err=%v payload=%s", ct, err, payload)
	}

	payload, ct, err = encodeRequestBody(&request{multipart: &multipartPayload{
		Fields: map[string]string{"f": "v"},
		Files:  []multipartFile{{FieldName: "file", FileName: "a.txt", Reader: strings.NewReader("hello")}},
	}})
	if err != nil {
		t.Fatalf("multipart encoding failed: %v", err)
	}
	if !strings.HasPrefix(ct, "multipart/form-data;") {
		t.Fatalf("unexpected content-type: %s", ct)
	}

	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil || mediaType != "multipart/form-data" {
		t.Fatalf("failed to parse content-type: %v", err)
	}
	mr := multipart.NewReader(bytes.NewReader(payload), params["boundary"])
	part, err := mr.NextPart()
	if err != nil {
		t.Fatalf("failed to read multipart field: %v", err)
	}
	if part.FormName() != "f" {
		t.Fatalf("unexpected first part field: %s", part.FormName())
	}

	if _, _, err := encodeRequestBody(&request{jsonBody: map[string]any{}, formBody: url.Values{}}); err == nil {
		t.Fatalf("expected body-type conflict error")
	}

	if _, _, err := encodeRequestBody(&request{multipart: &multipartPayload{Files: []multipartFile{{FieldName: "f", Reader: nil}}}}); !errors.Is(err, ErrNilFileReader) {
		t.Fatalf("expected ErrNilFileReader, got %v", err)
	}
}

func TestDoAndDecodeJSON_SuccessAndAuth(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true}`)
	})

	var out map[string]bool
	err := client.decodeJSON(context.Background(), &request{method: http.MethodGet, path: "/x", requireAuth: true}, &out)
	if err != nil {
		t.Fatalf("decodeJSON failed: %v", err)
	}
	if !out["ok"] {
		t.Fatalf("expected decoded payload")
	}
}

func TestInjectAuth_BearerAndOptional(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAccessToken("tok", WithBaseURL("https://api.example.com"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.example.com", nil)
	if authErr := c.injectAuth(req, true); authErr != nil {
		t.Fatalf("injectAuth failed: %v", authErr)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer tok" {
		t.Fatalf("unexpected auth header: %s", got)
	}

	c2, err := NewClient(WithBaseURL("https://api.example.com"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.example.com", nil)
	if authErr := c2.injectAuth(req2, false); authErr != nil {
		t.Fatalf("injectAuth optional failed: %v", authErr)
	}
}

func TestDo_MissingAuthAndRetry(t *testing.T) {
	t.Parallel()

	c, err := NewClient(WithBaseURL("https://example.com"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if authErr := c.injectAuth((&http.Request{Header: make(http.Header)}), true); !errors.Is(authErr, ErrMissingAuthentication) {
		t.Fatalf("expected ErrMissingAuthentication")
	}

	attempts := 0
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = io.WriteString(w, `{"type":"throttled","detail":"Request was throttled."}`)
			return
		}
		_, _ = io.WriteString(w, `{"ok":true}`)
	}, WithRetryPolicy(RetryPolicy{MaxRetries: 1, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond, RetryOn429: true}))

	var out map[string]bool
	err = client.decodeJSON(context.Background(), &request{method: http.MethodGet, path: "/retry", requireAuth: true}, &out)
	if err != nil {
		t.Fatalf("expected retry success, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected two attempts, got %d", attempts)
	}
}

func TestDo_APIErrorParsing(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Request-Id", "rid-1")
		w.Header().Set("Retry-After", "7")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"code":"bad_request","detail":{"field":"name"}}`)
	})

	err := client.decodeJSON(context.Background(), &request{method: http.MethodGet, path: "/err", requireAuth: true}, &map[string]any{})
	if err == nil {
		t.Fatalf("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != "bad_request" || apiErr.RequestID != "rid-1" || apiErr.RetryAfter != "7" {
		t.Fatalf("unexpected parsed API error: %+v", apiErr)
	}
	if apiErr.Message == "" {
		t.Fatalf("expected synthesized message")
	}
}

func TestDownload(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/pdf" {
			t.Fatalf("expected pdf accept")
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=test.pdf")
		_, _ = io.WriteString(w, "PDF")
	})

	resp, err := client.download(context.Background(), &request{method: http.MethodGet, path: "/pdf", requireAuth: true, accept: "application/pdf"})
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}
	defer func() { _ = resp.Close() }()

	if resp.ContentType != "application/pdf" {
		t.Fatalf("unexpected content type: %s", resp.ContentType)
	}
	b, _ := io.ReadAll(resp.Body)
	if string(b) != "PDF" {
		t.Fatalf("unexpected body: %s", string(b))
	}
}

func TestParseRetryAfter(t *testing.T) {
	t.Parallel()

	if d, ok := parseRetryAfter("3"); !ok || d != 3*time.Second {
		t.Fatalf("unexpected delta parse: %v %v", d, ok)
	}
	if _, ok := parseRetryAfter("bogus"); ok {
		t.Fatalf("expected invalid retry-after to fail parse")
	}
	if _, ok := parseRetryAfter("-1"); ok {
		t.Fatalf("expected negative retry-after to fail parse")
	}

	h := time.Now().Add(2 * time.Second).UTC().Format(http.TimeFormat)
	if _, ok := parseRetryAfter(h); !ok {
		t.Fatalf("expected http-date retry-after parse")
	}
}

func TestPopulateAPIErrorFromBodyFallback(t *testing.T) {
	t.Parallel()

	e := &APIError{}
	populateAPIErrorFromBody(e, []byte("plain text"))
	if e.Message != "plain text" {
		t.Fatalf("unexpected fallback message: %s", e.Message)
	}

	e = &APIError{}
	populateAPIErrorFromBody(e, []byte(`{"error":"invalid_grant","error_description":"bad code"}`))
	if e.Code != "invalid_grant" || e.Message != "bad code" {
		t.Fatalf("unexpected oauth error parse: %+v", e)
	}

	e = &APIError{}
	populateAPIErrorFromBody(e, []byte(`{"type":"request_error","detail":"invalid"}`))
	if e.Code != "request_error" || e.Message != "invalid" {
		t.Fatalf("unexpected v1 parse: %+v", e)
	}
}

func TestSleepWithContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := sleepWithContext(ctx, time.Second); err == nil {
		t.Fatalf("expected context cancellation error")
	}

	if err := sleepWithContext(context.Background(), 0); err != nil {
		t.Fatalf("unexpected error for zero sleep: %v", err)
	}
}

func TestParseAPIError_EmptyBody(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("")),
	}
	err := parseAPIError(resp)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError")
	}
	if apiErr.Message != http.StatusText(http.StatusUnauthorized) {
		t.Fatalf("unexpected message: %s", apiErr.Message)
	}
}

func TestDecodeJSON_EmptyBodyAndNilOut(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.decodeJSON(context.Background(), &request{method: http.MethodGet, path: "/x", requireAuth: true, expectedStatus: []int{http.StatusNoContent}}, nil); err != nil {
		t.Fatalf("decodeJSON nil out failed: %v", err)
	}

	client2 := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "")
	})
	var out map[string]any
	if err := client2.decodeJSON(context.Background(), &request{method: http.MethodGet, path: "/x", requireAuth: true}, &out); err != nil {
		t.Fatalf("decodeJSON empty body failed: %v", err)
	}
}

func TestDecodeJSON_InvalidJSON(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "{")
	})

	var out map[string]any
	err := client.decodeJSON(context.Background(), &request{method: http.MethodGet, path: "/bad", requireAuth: true}, &out)
	if err == nil {
		t.Fatalf("expected decode error")
	}
}

func TestDo_ContextCancelledBeforeRetry(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, `{}`)
	}, WithRetryPolicy(RetryPolicy{MaxRetries: 2, InitialBackoff: 100 * time.Millisecond, MaxBackoff: 100 * time.Millisecond}))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out map[string]any
	err := client.decodeJSON(ctx, &request{method: http.MethodGet, path: "/retry", requireAuth: true}, &out)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDoRequest_BuildRequestFailure(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAPIKey("k", WithBaseURL("https://api.example.com"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	_, err = c.do(context.Background(), &request{method: "\n", path: "/x", requireAuth: true}) //nolint:bodyclose // request build fails before response
	if err == nil {
		t.Fatalf("expected request build error")
	}
}

func TestEncodeRequestBody_JSONMarshalFailure(t *testing.T) {
	t.Parallel()

	_, _, err := encodeRequestBody(&request{jsonBody: map[string]any{"x": func() {}}})
	if err == nil {
		t.Fatalf("expected marshal failure")
	}
}

func TestPopulateAPIErrorFromBody_DetailObject(t *testing.T) {
	t.Parallel()

	e := &APIError{}
	populateAPIErrorFromBody(e, mustJSON(t, map[string]any{"detail": map[string]any{"x": "y"}}))
	if e.Details == nil {
		t.Fatalf("expected details object")
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}
	return b
}

func TestDrainAndClose(t *testing.T) {
	t.Parallel()

	if err := drainAndClose(nil); err != nil {
		t.Fatalf("drainAndClose nil failed: %v", err)
	}
	if err := drainAndClose(io.NopCloser(strings.NewReader("x"))); err != nil {
		t.Fatalf("drainAndClose failed: %v", err)
	}
}

func TestHeaderFirst(t *testing.T) {
	t.Parallel()

	h := make(http.Header)
	h.Set("B", "2")
	if v := headerFirst(h, "A", "B"); v != "2" {
		t.Fatalf("unexpected header value: %s", v)
	}
}

func TestStringValue(t *testing.T) {
	t.Parallel()

	if v, ok := stringValue("x"); !ok || v != "x" {
		t.Fatalf("unexpected string conversion")
	}
	if _, ok := stringValue(1); ok {
		t.Fatalf("expected false for non-string")
	}
}

func TestEscapePathParam(t *testing.T) {
	t.Parallel()

	if _, err := escapePathParam(" "); !errors.Is(err, ErrEmptyPathParameter) {
		t.Fatalf("expected ErrEmptyPathParameter")
	}
	v, err := escapePathParam("a/b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "a%2Fb" {
		t.Fatalf("unexpected escaped value: %s", v)
	}
}

func TestRetryBackoff(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAPIKey("k")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	c.retryPolicy = RetryPolicy{MaxRetries: 2, InitialBackoff: 100 * time.Millisecond, MaxBackoff: 300 * time.Millisecond}.normalize()

	if b := c.backoff(0); b != 100*time.Millisecond {
		t.Fatalf("unexpected backoff: %v", b)
	}
	if b := c.backoff(10); b != 300*time.Millisecond {
		t.Fatalf("unexpected capped backoff: %v", b)
	}
}

func TestShouldRetryHelpers(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAPIKey("k")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	c.retryPolicy = RetryPolicy{MaxRetries: 1, RetryOn429: true, RetryOn5xx: true}.normalize()

	if !c.shouldRetryOnError(0, errTestDummy) {
		t.Fatalf("expected retry on transport error")
	}
	if c.shouldRetryOnError(1, errTestDummy) {
		t.Fatalf("expected no retry beyond max")
	}

	resp := &http.Response{StatusCode: http.StatusTooManyRequests}
	if !c.shouldRetryOnStatus(0, resp) {
		t.Fatalf("expected retry for 429")
	}
	resp.StatusCode = http.StatusInternalServerError
	if !c.shouldRetryOnStatus(0, resp) {
		t.Fatalf("expected retry for 5xx")
	}
	resp.StatusCode = http.StatusBadRequest
	if c.shouldRetryOnStatus(0, resp) {
		t.Fatalf("did not expect retry for 4xx")
	}
}

func TestDo_NoReq(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAPIKey("k")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if _, err := c.do(context.Background(), nil); !errors.Is(err, ErrNilRequest) { //nolint:bodyclose // early validation error
		t.Fatalf("expected ErrNilRequest, got %v", err)
	}
}

func TestBuildURL_InvalidPath(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAPIKey("k", WithBaseURL("https://api.example.com"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if _, err := c.buildURL("", nil); err == nil {
		t.Fatalf("expected missing path error")
	}
}

func TestDownloadCloseNil(t *testing.T) {
	t.Parallel()

	var d *DownloadResponse
	if err := d.Close(); err != nil {
		t.Fatalf("unexpected nil close error: %v", err)
	}

	d = &DownloadResponse{}
	if err := d.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}

func TestDoTransportError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatalf("expected hijacker")
		}
		conn, _, _ := hj.Hijack()
		_ = conn.Close()
		_ = r
	}))
	defer srv.Close()

	client, err := NewClientWithAPIKey("k", WithBaseURL(srv.URL), WithRetryPolicy(RetryPolicy{MaxRetries: 0, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond}))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	_, err = client.do(context.Background(), &request{method: http.MethodGet, path: "/x", requireAuth: true}) //nolint:bodyclose // connection closed before response body
	if err == nil {
		t.Fatalf("expected transport error")
	}
}

func TestJoinPaths_EmptyRelative(t *testing.T) {
	t.Parallel()

	if got := joinPaths("/public/v1/", ""); got != "/public/v1" {
		t.Fatalf("unexpected joined path: %s", got)
	}
	if got := joinPaths("/public/v1", ""); got != "/public/v1" {
		t.Fatalf("unexpected joined path: %s", got)
	}
}

func TestBuildURL_ParseEndpointPathError(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAPIKey("k", WithBaseURL("https://api.example.com"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if _, err := c.buildURL("http://[::1", nil); err == nil {
		t.Fatalf("expected parse endpoint path error")
	}
}

func TestDo_BuildURLAndBodyEncodeFailures(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAPIKey("k", WithBaseURL("https://api.example.com"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if resp, err2 := c.do(context.Background(), &request{method: http.MethodGet, path: "http://[::1", requireAuth: true}); err2 == nil {
		_ = resp.Body.Close()
		t.Fatalf("expected buildURL error")
	}

	resp, bodyErr := c.do(context.Background(), &request{
		method:      http.MethodGet,
		path:        "/x",
		requireAuth: true,
		jsonBody:    map[string]any{"x": 1},
		formBody:    url.Values{"a": []string{"1"}},
	})
	if resp != nil {
		_ = resp.Body.Close()
	}
	if !errors.Is(bodyErr, errOnlyOneBodyType) {
		t.Fatalf("expected body-type conflict error, got %v", bodyErr)
	}
}

func TestDo_CustomHeadersPropagation(t *testing.T) {
	t.Parallel()

	var seen []string
	client := newRoundTripperClient(t, func(req *http.Request) (*http.Response, error) {
		seen = req.Header.Values("X-Test")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		}, nil
	})

	var out map[string]bool
	err := client.decodeJSON(context.Background(), &request{
		method:      http.MethodGet,
		path:        "/x",
		requireAuth: true,
		headers:     http.Header{"X-Test": []string{"a", "b"}},
	}, &out)
	if err != nil {
		t.Fatalf("decodeJSON failed: %v", err)
	}
	if len(seen) != 2 || seen[0] != "a" || seen[1] != "b" {
		t.Fatalf("unexpected propagated headers: %v", seen)
	}
}

func TestDoAttempt_InjectAuthError(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{Transport: roundTripperFunc(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})}
	c, err := NewClient(WithBaseURL("https://api.example.com"), WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	resp, retryable, err := c.doAttempt(context.Background(), &request{method: http.MethodGet, requireAuth: true}, "https://api.example.com/x", nil, "")
	if resp != nil {
		_ = resp.Body.Close()
	}
	if !errors.Is(err, ErrMissingAuthentication) {
		t.Fatalf("expected ErrMissingAuthentication, got %v", err)
	}
	if retryable {
		t.Fatalf("did not expect auth error to be retryable")
	}
}

func TestDecodeJSON_ReadBodyError(t *testing.T) {
	t.Parallel()

	client := newRoundTripperClient(t, func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(newErrorReader()),
		}, nil
	})

	var out map[string]any
	err := client.decodeJSON(context.Background(), &request{method: http.MethodGet, path: "/x", requireAuth: true}, &out)
	if err == nil || !strings.Contains(err.Error(), "decode response body") {
		t.Fatalf("expected read response body error, got %v", err)
	}
}

func TestDo_RetryableTransportErrorThenFail(t *testing.T) {
	t.Parallel()

	attempts := 0
	client := newRoundTripperClient(t, func(_ *http.Request) (*http.Response, error) {
		attempts++
		return nil, errTestDummy
	}, WithRetryPolicy(RetryPolicy{
		MaxRetries:     1,
		InitialBackoff: time.Nanosecond,
		MaxBackoff:     time.Nanosecond,
	}))

	_, err := client.do(context.Background(), &request{method: http.MethodGet, path: "/retry", requireAuth: true}) //nolint:bodyclose // expected transport error
	if !errors.Is(err, errTestDummy) {
		t.Fatalf("expected transport error, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected two attempts, got %d", attempts)
	}
}

func TestDo_RetryOnStatus_CanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	client := newRoundTripperClient(t, func(_ *http.Request) (*http.Response, error) {
		cancel()
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Header:     http.Header{"Retry-After": []string{"1"}},
			Body:       io.NopCloser(strings.NewReader(`{}`)),
		}, nil
	}, WithRetryPolicy(RetryPolicy{
		MaxRetries:     1,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     time.Millisecond,
		RetryOn429:     true,
	}))

	_, err := client.do(ctx, &request{method: http.MethodGet, path: "/retry", requireAuth: true}) //nolint:bodyclose // expected context cancellation
	if err == nil {
		t.Fatalf("expected context cancellation error")
	}
}

func TestDownload_ErrorPath(t *testing.T) {
	t.Parallel()

	client := newRoundTripperClient(t, func(_ *http.Request) (*http.Response, error) {
		return nil, errTestDummy
	})

	if _, err := client.download(context.Background(), &request{method: http.MethodGet, path: "/x", requireAuth: true}); !errors.Is(err, errTestDummy) {
		t.Fatalf("expected transport error from download, got %v", err)
	}
}

func TestEncodeMultipart_DefaultsAndErrors(t *testing.T) {
	t.Parallel()

	payload, contentType, err := encodeMultipart(&multipartPayload{
		Files: []multipartFile{{
			Reader: strings.NewReader("data"),
		}},
	})
	if err != nil {
		t.Fatalf("encodeMultipart defaults failed: %v", err)
	}
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatalf("parse media type failed: %v", err)
	}
	mr := multipart.NewReader(bytes.NewReader(payload), params["boundary"])
	part, err := mr.NextPart()
	if err != nil {
		t.Fatalf("read multipart part failed: %v", err)
	}
	if part.FormName() != "file" || part.FileName() != "upload.bin" {
		t.Fatalf("unexpected multipart defaults: form=%s file=%s", part.FormName(), part.FileName())
	}

	if _, _, err := encodeMultipart(&multipartPayload{
		Files: []multipartFile{{
			FieldName: "file",
			FileName:  "x.txt",
			Reader:    newErrorReader(),
		}},
	}); err == nil {
		t.Fatalf("expected multipart copy error")
	}
}

func TestStatusExpectedAndExtractErrorDetails(t *testing.T) {
	t.Parallel()

	if statusExpected(http.StatusOK, []int{http.StatusCreated}) {
		t.Fatalf("did not expect status to match explicit expected list")
	}

	e := &APIError{}
	extractErrorDetails(e, map[string]any{"details": map[string]any{"field": "name"}})
	if e.Details == nil {
		t.Fatalf("expected details to be extracted")
	}
}

func TestParseRetryAfter_EmptyAndPastDate(t *testing.T) {
	t.Parallel()

	if _, ok := parseRetryAfter(""); ok {
		t.Fatalf("expected empty retry-after to fail parse")
	}

	past := time.Now().Add(-time.Minute).UTC().Format(http.TimeFormat)
	d, ok := parseRetryAfter(past)
	if !ok || d != 0 {
		t.Fatalf("expected past retry-after to parse as zero delay: ok=%v d=%v", ok, d)
	}
}

func TestSleepWithContext_TimerCompletion(t *testing.T) {
	t.Parallel()

	if err := sleepWithContext(context.Background(), time.Millisecond); err != nil {
		t.Fatalf("expected sleepWithContext timer completion, got %v", err)
	}
}

func TestRetryBackoff_UncappedBranch(t *testing.T) {
	t.Parallel()

	c, err := NewClientWithAPIKey("k")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	c.retryPolicy = RetryPolicy{
		MaxRetries:     5,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     5 * time.Second,
	}.normalize()

	if got := c.backoff(1); got != 200*time.Millisecond {
		t.Fatalf("unexpected backoff for uncapped branch: %v", got)
	}
}
