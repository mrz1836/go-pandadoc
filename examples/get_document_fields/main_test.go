package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRun_Validation(t *testing.T) {
	t.Parallel()

	if err := run(context.Background(), []string{"cmd"}, func(string) string { return "" }, io.Discard); err == nil || !strings.Contains(err.Error(), "PANDADOC_API_KEY") {
		t.Fatalf("expected missing API key error, got %v", err)
	}

	env := map[string]string{"PANDADOC_API_KEY": "k"}
	getenv := func(k string) string { return env[k] }
	if err := run(context.Background(), []string{"cmd"}, getenv, io.Discard); err == nil || !strings.Contains(strings.ToLower(err.Error()), "usage") {
		t.Fatalf("expected usage error, got %v", err)
	}
}

func TestRun_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/v1/documents/doc1/details" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"id":"doc1","name":"Doc","status":"document.draft","fields":[{"name":"F","type":"text","value":"v"}],"tokens":[{"name":"T","value":"X"}]}`)
	}))
	defer srv.Close()

	env := map[string]string{"PANDADOC_API_KEY": "k", "PANDADOC_BASE_URL": srv.URL}
	getenv := func(k string) string { return env[k] }

	var out bytes.Buffer
	if err := run(context.Background(), []string{"cmd", "doc1"}, getenv, &out); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if !strings.Contains(out.String(), "Document: Doc") || !strings.Contains(out.String(), "Fields:") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestRun_CreateClientError(t *testing.T) {
	t.Parallel()

	env := map[string]string{"PANDADOC_API_KEY": "k", "PANDADOC_BASE_URL": "::bad"}
	getenv := func(k string) string { return env[k] }

	err := run(context.Background(), []string{"cmd", "doc1"}, getenv, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "create client") {
		t.Fatalf("expected create client error, got %v", err)
	}
}

func TestRun_PropagatesDetailsError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"detail":"bad auth"}`)
	}))
	defer srv.Close()

	env := map[string]string{"PANDADOC_API_KEY": "k", "PANDADOC_BASE_URL": srv.URL}
	getenv := func(k string) string { return env[k] }

	err := run(context.Background(), []string{"cmd", "doc1"}, getenv, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "get details") {
		t.Fatalf("expected get details error, got %v", err)
	}
}

func TestMain_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"id":"doc1","name":"Doc","status":"document.draft","fields":[],"tokens":[]}`)
	}))
	defer srv.Close()

	t.Setenv("PANDADOC_API_KEY", "k")
	t.Setenv("PANDADOC_BASE_URL", srv.URL)

	origArgs := os.Args
	t.Cleanup(func() { os.Args = origArgs })
	os.Args = []string{"cmd", "doc1"}

	main()
}

func TestMain_PanicsOnError(t *testing.T) {
	t.Setenv("PANDADOC_API_KEY", "")
	t.Setenv("PANDADOC_BASE_URL", "")

	origArgs := os.Args
	t.Cleanup(func() { os.Args = origArgs })
	os.Args = []string{"cmd"}

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic from main on missing env/args")
		}
	}()
	main()
}
