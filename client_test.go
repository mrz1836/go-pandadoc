package pandadoc_test

import (
	"testing"
	"time"

	"github.com/mrz1836/go-pandadoc"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	client, err := pandadoc.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_EmptyAPIKey(t *testing.T) {
	t.Parallel()

	_, err := pandadoc.NewClient("")
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	t.Parallel()

	client, err := pandadoc.NewClient("test-api-key",
		pandadoc.WithTimeout(60*time.Second),
		pandadoc.WithUserAgent("test-agent/1.0"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestClient_Documents(t *testing.T) {
	t.Parallel()

	client, err := pandadoc.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	docs := client.Documents()
	if docs == nil {
		t.Error("expected non-nil Documents API")
	}
}

func TestClient_Catalog(t *testing.T) {
	t.Parallel()

	client, err := pandadoc.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	catalog := client.Catalog()
	if catalog == nil {
		t.Error("expected non-nil Catalog API")
	}
}
