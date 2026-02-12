package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRun_RequiresAPIKey(t *testing.T) {
	t.Parallel()

	err := run(context.Background(), func(string) string { return "" }, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "PANDADOC_API_KEY") {
		t.Fatalf("expected missing API key error, got %v", err)
	}
}

func TestRun_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/v2/product-catalog/items/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"items":[{"uuid":"i1","title":"Item 1","sku":"SKU1","price":12.5,"currency":"USD"}],"has_more_items":false,"total":1}`)
	}))
	defer srv.Close()

	env := map[string]string{"PANDADOC_API_KEY": "k", "PANDADOC_BASE_URL": srv.URL}
	getenv := func(k string) string { return env[k] }

	var out bytes.Buffer
	if err := run(context.Background(), getenv, &out); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !strings.Contains(out.String(), "Found 1 catalog items") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestRun_CreateClientError(t *testing.T) {
	t.Parallel()

	env := map[string]string{"PANDADOC_API_KEY": "k", "PANDADOC_BASE_URL": "::bad"}
	getenv := func(k string) string { return env[k] }

	err := run(context.Background(), getenv, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "create client") {
		t.Fatalf("expected create client error, got %v", err)
	}
}

func TestRun_PropagatesAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, `{"detail":"server failure"}`)
	}))
	defer srv.Close()

	env := map[string]string{"PANDADOC_API_KEY": "k", "PANDADOC_BASE_URL": srv.URL}
	getenv := func(k string) string { return env[k] }

	err := run(context.Background(), getenv, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "search catalog") {
		t.Fatalf("expected search catalog error, got %v", err)
	}
}

func TestMain_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"items":[],"has_more_items":false,"total":0}`)
	}))
	defer srv.Close()

	t.Setenv("PANDADOC_API_KEY", "k")
	t.Setenv("PANDADOC_BASE_URL", srv.URL)
	main()
}

func TestMain_PanicsOnError(t *testing.T) {
	t.Setenv("PANDADOC_API_KEY", "")
	t.Setenv("PANDADOC_BASE_URL", "")

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic from main on missing env")
		}
	}()
	main()
}
