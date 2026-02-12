package pandadoc_test

import (
	"context"
	"testing"

	"github.com/mrz1836/go-pandadoc"
)

// Ensure a mock can be created as described in README
type mockDocuments struct {
	pandadoc.DocumentsService
}

func (m *mockDocuments) List(_ context.Context, _ *pandadoc.ListDocumentsOptions) (*pandadoc.DocumentListResponse, error) {
	return &pandadoc.DocumentListResponse{
		Results: []pandadoc.DocumentSummary{
			{ID: "mock-doc-1", Name: "Mock Document"},
		},
	}, nil
}

func TestMockability(t *testing.T) {
	// This test simply verifies that the mock struct satisfies the interface
	// and acts as a compile-time check for the README example.
	var svc pandadoc.DocumentsService = &mockDocuments{}

	resp, err := svc.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("mock list failed: %v", err)
	}
	if len(resp.Results) != 1 || resp.Results[0].ID != "mock-doc-1" {
		t.Errorf("unexpected mock response: %+v", resp)
	}
}
