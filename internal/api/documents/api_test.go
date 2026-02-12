package documents_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mrz1836/go-pandadoc/commands"
	"github.com/mrz1836/go-pandadoc/internal/api/documents"
	"github.com/mrz1836/go-pandadoc/internal/httpclient"
	"github.com/mrz1836/go-pandadoc/models"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*documents.API, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := httpclient.New(server.URL+"/", "test-api-key", "test-user-agent", 30*time.Second, nil)
	api := documents.New(client)
	return api, server.Close
}

func TestAPI_List(t *testing.T) {
	t.Parallel()

	expectedDocs := models.DocumentListResponse{
		Results: []models.Document{
			{
				ID:     "doc1",
				Name:   "Test Document 1",
				Status: "document.draft",
			},
			{
				ID:     "doc2",
				Name:   "Test Document 2",
				Status: "document.completed",
			},
		},
		Count: 2,
	}

	api, cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/documents" {
			t.Errorf("expected /documents, got %s", r.URL.Path)
		}
		// Check auth header
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedDocs); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	result, err := api.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 2 {
		t.Errorf("expected 2 documents, got %d", len(result.Results))
	}
	if result.Results[0].ID != "doc1" {
		t.Errorf("expected doc1, got %s", result.Results[0].ID)
	}
}

func TestAPI_List_WithOptions(t *testing.T) {
	t.Parallel()

	api, cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("page") != "2" {
			t.Errorf("expected page=2, got %s", query.Get("page"))
		}
		if query.Get("count") != "10" {
			t.Errorf("expected count=10, got %s", query.Get("count"))
		}
		if query.Get("status") != "document.completed" {
			t.Errorf("expected status=document.completed, got %s", query.Get("status"))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(models.DocumentListResponse{}); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	opts := &commands.ListDocumentsOptions{
		Page:   2,
		Count:  10,
		Status: "document.completed",
	}

	_, err := api.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPI_Get(t *testing.T) {
	t.Parallel()

	expectedDoc := models.Document{
		ID:     "doc123",
		Name:   "Test Document",
		Status: "document.draft",
	}

	api, cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/documents/doc123" {
			t.Errorf("expected /documents/doc123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedDoc); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	result, err := api.Get(context.Background(), "doc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "doc123" {
		t.Errorf("expected doc123, got %s", result.ID)
	}
}

func TestAPI_Get_EmptyID(t *testing.T) {
	t.Parallel()

	api := documents.New(nil)
	_, err := api.Get(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty ID")
	}
}

func TestAPI_GetStatus(t *testing.T) {
	t.Parallel()

	expectedStatus := models.DocumentStatus{
		ID:     "doc123",
		Status: "document.sent",
	}

	api, cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/documents/doc123/status" {
			t.Errorf("expected /documents/doc123/status, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedStatus); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	result, err := api.GetStatus(context.Background(), "doc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "document.sent" {
		t.Errorf("expected document.sent, got %s", result.Status)
	}
}

func TestAPI_GetDetails(t *testing.T) {
	t.Parallel()

	expectedDetails := models.DocumentDetails{
		ID:     "doc123",
		Name:   "Test Document",
		Status: "document.completed",
		Fields: map[string]interface{}{
			"field1": "value1",
		},
	}

	api, cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/documents/doc123/details" {
			t.Errorf("expected /documents/doc123/details, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedDetails); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	result, err := api.GetDetails(context.Background(), "doc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "doc123" {
		t.Errorf("expected doc123, got %s", result.ID)
	}
	if result.Fields["field1"] != "value1" {
		t.Errorf("expected field1=value1, got %v", result.Fields["field1"])
	}
}

func TestAPI_Update(t *testing.T) {
	t.Parallel()

	expectedDoc := models.Document{
		ID:     "doc123",
		Name:   "Updated Document",
		Status: "document.draft",
	}

	api, cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/documents/doc123" {
			t.Errorf("expected /documents/doc123, got %s", r.URL.Path)
		}

		// Verify request body
		var update commands.UpdateDocument
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if update.Name != "Updated Document" {
			t.Errorf("expected Updated Document, got %s", update.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedDoc); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})
	defer cleanup()

	update := &commands.UpdateDocument{
		Name: "Updated Document",
	}

	result, err := api.Update(context.Background(), "doc123", update)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "Updated Document" {
		t.Errorf("expected Updated Document, got %s", result.Name)
	}
}

func TestAPI_Update_NilUpdate(t *testing.T) {
	t.Parallel()

	api := documents.New(nil)
	_, err := api.Update(context.Background(), "doc123", nil)
	if err == nil {
		t.Error("expected error for nil update")
	}
}
