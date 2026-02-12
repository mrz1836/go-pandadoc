package pagination

import (
	"testing"

	"github.com/mrz1836/go-pandadoc/models"
)

func TestNextPage(t *testing.T) {
	t.Run("nil current options", func(t *testing.T) {
		next := NextPage(nil)

		if next.Page != 1 {
			t.Errorf("expected Page 1, got %d", next.Page)
		}

		if next.Count != 50 {
			t.Errorf("expected Count 50, got %d", next.Count)
		}
	})

	t.Run("existing options", func(t *testing.T) {
		current := &models.ListOptions{
			Page:  2,
			Count: 25,
		}

		next := NextPage(current)

		if next.Page != 3 {
			t.Errorf("expected Page 3, got %d", next.Page)
		}

		if next.Count != 25 {
			t.Errorf("expected Count 25, got %d", next.Count)
		}
	})
}

func TestHasMore(t *testing.T) {
	t.Run("nil metadata", func(t *testing.T) {
		if HasMore(nil) {
			t.Error("expected false for nil metadata")
		}
	})

	t.Run("next is nil", func(t *testing.T) {
		meta := &models.PaginationMeta{
			Count:    50,
			Next:     nil,
			Previous: nil,
		}

		if HasMore(meta) {
			t.Error("expected false when Next is nil")
		}
	})

	t.Run("next is not nil", func(t *testing.T) {
		nextURL := "https://api.pandadoc.com/public/v1/documents?page=2"
		meta := &models.PaginationMeta{
			Count:    50,
			Next:     &nextURL,
			Previous: nil,
		}

		if !HasMore(meta) {
			t.Error("expected true when Next is not nil")
		}
	})
}
