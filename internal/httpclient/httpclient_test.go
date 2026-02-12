package httpclient

import (
	"context"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mrz1836/go-pandadoc/errors"
)

func TestNew(t *testing.T) {
	client := New("https://api.example.com", "test-key", "test-agent/1.0", 30*time.Second, nil)

	if client == nil {
		t.Fatal("expected client to be created")
	}

	if client.baseURL != "https://api.example.com" {
		t.Errorf("expected baseURL 'https://api.example.com', got '%s'", client.baseURL)
	}

	if client.apiKey != "test-key" {
		t.Errorf("expected apiKey 'test-key', got '%s'", client.apiKey)
	}

	if client.userAgent != "test-agent/1.0" {
		t.Errorf("expected userAgent 'test-agent/1.0', got '%s'", client.userAgent)
	}
}

func TestDoRequest_SuccessfulRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Authorization") != "API-Key test-key" {
			t.Errorf("expected Authorization header 'API-Key test-key', got '%s'", r.Header.Get("Authorization"))
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("User-Agent") != "test-agent/1.0" {
			t.Errorf("expected User-Agent 'test-agent/1.0', got '%s'", r.Header.Get("User-Agent"))
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status": "ok"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := New(server.URL, "test-key", "test-agent/1.0", 30*time.Second, nil)
	resp, err := client.DoRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	_ = resp.Body.Close()
}

func TestDoRequest_WithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL, "test-key", "test-agent/1.0", 30*time.Second, nil)

	body := map[string]string{"key": "value"}
	resp, err := client.DoRequest(context.Background(), http.MethodPost, "/test", body)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_ = resp.Body.Close()
}

func TestDoRequest_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		if _, err := w.Write([]byte(`{"message": "Rate limit exceeded"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := New(server.URL, "test-key", "test-agent/1.0", 30*time.Second, nil)
	resp, err := client.DoRequest(context.Background(), http.MethodGet, "/test", nil)
	if resp != nil {
		_ = resp.Body.Close()
	}

	if !stderrors.Is(err, errors.ErrRateLimitExceeded) {
		t.Errorf("expected ErrRateLimitExceeded, got %v", err)
	}
}

func TestDoRequest_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"message": "Invalid request", "details": {"field": "error"}}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := New(server.URL, "test-key", "test-agent/1.0", 30*time.Second, nil)
	resp, err := client.DoRequest(context.Background(), http.MethodGet, "/test", nil)
	if resp != nil {
		_ = resp.Body.Close()
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *errors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code 400, got %d", apiErr.StatusCode)
	}

	if apiErr.Message != "Invalid request" {
		t.Errorf("expected message 'Invalid request', got '%s'", apiErr.Message)
	}
}

func TestDoRequest_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL, "test-key", "test-agent/1.0", 30*time.Second, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	resp, err := client.DoRequest(ctx, http.MethodGet, "/test", nil)
	if resp != nil {
		_ = resp.Body.Close()
	}

	if err == nil {
		t.Fatal("expected error due to context cancellation")
	}
}

func TestBuildURL(t *testing.T) {
	client := New("https://api.example.com/v1/", "test-key", "test-agent/1.0", 30*time.Second, nil)

	t.Run("simple path", func(t *testing.T) {
		url, err := client.buildURL("documents")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expected := "https://api.example.com/v1/documents"
		if url != expected {
			t.Errorf("expected URL '%s', got '%s'", expected, url)
		}
	})

	t.Run("path with leading slash", func(t *testing.T) {
		url, err := client.buildURL("/documents")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expected := "https://api.example.com/documents"
		if url != expected {
			t.Errorf("expected URL '%s', got '%s'", expected, url)
		}
	})
}
