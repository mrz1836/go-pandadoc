package catalog_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mrz1836/go-pandadoc/commands"
	"github.com/mrz1836/go-pandadoc/internal/api/catalog"
	"github.com/mrz1836/go-pandadoc/internal/httpclient"
	"github.com/mrz1836/go-pandadoc/models"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*catalog.API, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := httpclient.New(server.URL+"/", "test-api-key", "test-user-agent", 30*time.Second, nil)
	api := catalog.New(client)
	return api, server.Close
}

func TestAPI_List(t *testing.T) {
	t.Parallel()

	expectedItems := models.CatalogListResponse{
		Results: []models.CatalogItem{
			{
				ID:   "item1",
				Name: "Product 1",
				SKU:  "SKU001",
			},
			{
				ID:   "item2",
				Name: "Product 2",
				SKU:  "SKU002",
			},
		},
		Count: 2,
	}

	api, cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/catalog" {
			t.Errorf("expected /catalog, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedItems); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	result, err := api.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Results))
	}
	if result.Results[0].ID != "item1" {
		t.Errorf("expected item1, got %s", result.Results[0].ID)
	}
}

func TestAPI_List_WithOptions(t *testing.T) {
	t.Parallel()

	api, cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("page") != "2" {
			t.Errorf("expected page=2, got %s", query.Get("page"))
		}
		if query.Get("count") != "25" {
			t.Errorf("expected count=25, got %s", query.Get("count"))
		}
		if query.Get("q") != "widget" {
			t.Errorf("expected q=widget, got %s", query.Get("q"))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(models.CatalogListResponse{}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	opts := &commands.ListCatalogOptions{
		Page:  2,
		Count: 25,
		Q:     "widget",
	}

	_, err := api.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPI_Get(t *testing.T) {
	t.Parallel()

	price := 99.99
	expectedItem := models.CatalogItem{
		ID:          "item123",
		Name:        "Test Product",
		Description: "A test product",
		SKU:         "TEST-SKU",
		Price: &models.CatalogPrice{
			Value:    price,
			Currency: "USD",
		},
	}

	api, cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/catalog/item123" {
			t.Errorf("expected /catalog/item123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedItem); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	result, err := api.Get(context.Background(), "item123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "item123" {
		t.Errorf("expected item123, got %s", result.ID)
	}
	if result.Name != "Test Product" {
		t.Errorf("expected Test Product, got %s", result.Name)
	}
	if result.Price == nil || result.Price.Value != price {
		t.Errorf("expected price %.2f, got %v", price, result.Price)
	}
}

func TestAPI_Get_EmptyID(t *testing.T) {
	t.Parallel()

	api := catalog.New(nil)
	_, err := api.Get(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty ID")
	}
}
