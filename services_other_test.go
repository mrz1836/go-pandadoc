package pandadoc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

//nolint:gocognit // Test function that validates all product catalog methods
func TestProductCatalogService_AllMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/public/v2/product-catalog/items/search":
			if r.URL.Query().Get("query") != "coffee" {
				t.Fatalf("missing query filter")
			}
			_, _ = io.WriteString(w, `{"items":[{"uuid":"i1","title":"Coffee"}],"has_more_items":false,"total":1}`)
		case r.Method == http.MethodPost && r.URL.Path == "/public/v2/product-catalog/items":
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode create payload: %v", err)
			}
			_, _ = io.WriteString(w, `{"uuid":"i1","title":"Coffee"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/public/v2/product-catalog/items/i1":
			_, _ = io.WriteString(w, `{"uuid":"i1","title":"Coffee"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/public/v2/product-catalog/items/i1":
			_, _ = io.WriteString(w, `{"uuid":"i1","title":"Coffee Updated"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/public/v2/product-catalog/items/i1":
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	searchResp, err := client.ProductCatalog().Search(context.Background(), &SearchProductCatalogItemsOptions{Query: "coffee", Page: 1, PerPage: 10})
	if err != nil || searchResp.Total != 1 || len(searchResp.Items) != 1 {
		t.Fatalf("Search failed: %v %+v", err, searchResp)
	}

	created, err := client.ProductCatalog().Create(context.Background(), CreateProductCatalogItemRequest{"title": "Coffee"})
	if err != nil || created.UUID != "i1" {
		t.Fatalf("Create failed: %v %+v", err, created)
	}

	got, err := client.ProductCatalog().Get(context.Background(), "i1")
	if err != nil || got.UUID != "i1" {
		t.Fatalf("Get failed: %v %+v", err, got)
	}

	updated, err := client.ProductCatalog().Update(context.Background(), "i1", UpdateProductCatalogItemRequest{"title": "Coffee Updated"})
	if err != nil || updated.Title != "Coffee Updated" {
		t.Fatalf("Update failed: %v %+v", err, updated)
	}

	if err := client.ProductCatalog().Delete(context.Background(), "i1"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestProductCatalogService_SearchQuerySerialization(t *testing.T) {
	t.Parallel()

	noCategory := true
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/public/v2/product-catalog/items/search" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("page") != "3" || q.Get("per_page") != "20" || q.Get("order_by") != "-date_modified" || q.Get("category_id") != "cat1" || q.Get("no_category") != "true" {
			t.Fatalf("unexpected query params: %s", r.URL.RawQuery)
		}
		if len(q["types"]) != 1 || q["types"][0] != string(ProductCatalogItemTypeRegular) {
			t.Fatalf("unexpected types query: %v", q["types"])
		}
		if len(q["billing_types"]) != 1 || q["billing_types"][0] != string(ProductCatalogBillingTypeRecurring) {
			t.Fatalf("unexpected billing types query: %v", q["billing_types"])
		}
		if len(q["exclude_uuids"]) != 2 {
			t.Fatalf("unexpected exclude_uuids: %v", q["exclude_uuids"])
		}
		_, _ = io.WriteString(w, `{"items":[],"has_more_items":false,"total":0}`)
	})

	_, err := client.ProductCatalog().Search(context.Background(), &SearchProductCatalogItemsOptions{
		Page:         3,
		PerPage:      20,
		Query:        "q",
		OrderBy:      "-date_modified",
		Types:        []ProductCatalogItemType{ProductCatalogItemTypeRegular},
		BillingTypes: []ProductCatalogBillingType{ProductCatalogBillingTypeRecurring},
		ExcludeUUIDs: []string{"u1", "u2"},
		CategoryID:   "cat1",
		NoCategory:   &noCategory,
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
}

func TestProductCatalogService_Validation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = w
		_ = r
		t.Fatalf("handler should not be called")
	})

	if _, err := client.ProductCatalog().Create(context.Background(), nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if _, err := client.ProductCatalog().Update(context.Background(), "i1", nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if _, err := client.ProductCatalog().Get(context.Background(), ""); err == nil {
		t.Fatalf("expected empty path parameter error")
	}
	if err := client.ProductCatalog().Delete(context.Background(), ""); err == nil {
		t.Fatalf("expected empty path parameter error")
	}
}

func TestOAuthService_Token(t *testing.T) {
	t.Parallel()

	c, err := NewClient(WithBaseURL("http://example.invalid"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if _, tokenErr := c.OAuth().Token(context.Background(), nil); tokenErr == nil {
		t.Fatalf("expected nil request error")
	}

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/oauth2/access_token" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			t.Fatalf("unexpected content-type: %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "grant_type=authorization_code") {
			t.Fatalf("unexpected oauth form body: %s", string(body))
		}
		_, _ = io.WriteString(w, `{"access_token":"a","refresh_token":"r","token_type":"Bearer","expires_in":3600}`)
	})

	resp, err := client.OAuth().Token(context.Background(), &OAuthTokenRequest{
		GrantType:    "authorization_code",
		ClientID:     "cid",
		ClientSecret: "sec",
		Code:         "code",
		Scope:        "read+write",
	})
	if err != nil {
		t.Fatalf("Token failed: %v", err)
	}
	if resp.AccessToken != "a" || resp.TokenType != "Bearer" {
		t.Fatalf("unexpected token response: %+v", resp)
	}
}

func TestWebhookSubscriptionsService_AllMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/public/v1/webhook-subscriptions":
			_, _ = io.WriteString(w, `{"items":[{"uuid":"w1","name":"sub"}]}`)
		case r.Method == http.MethodPost && r.URL.Path == "/public/v1/webhook-subscriptions":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"uuid":"w1","name":"sub"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/public/v1/webhook-subscriptions/w1":
			_, _ = io.WriteString(w, `{"uuid":"w1","name":"sub"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/public/v1/webhook-subscriptions/w1":
			_, _ = io.WriteString(w, `{"uuid":"w1","name":"updated"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/public/v1/webhook-subscriptions/w1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPatch && r.URL.Path == "/public/v1/webhook-subscriptions/w1/shared-key":
			_, _ = io.WriteString(w, `{"shared_key":"new-key"}`)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	list, err := client.WebhookSubscriptions().List(context.Background(), &ListWebhookSubscriptionsOptions{Count: 10, Page: 2})
	if err != nil || len(list.Items) != 1 {
		t.Fatalf("List failed: %v %+v", err, list)
	}
	created, err := client.WebhookSubscriptions().Create(context.Background(), &WebhookSubscriptionRequest{Name: "sub", URL: "https://example.com"})
	if err != nil || created.UUID != "w1" {
		t.Fatalf("Create failed: %v %+v", err, created)
	}
	got, err := client.WebhookSubscriptions().Get(context.Background(), "w1")
	if err != nil || got.UUID != "w1" {
		t.Fatalf("Get failed: %v %+v", err, got)
	}
	upd, err := client.WebhookSubscriptions().Update(context.Background(), "w1", &WebhookSubscriptionRequest{Name: "updated"})
	if err != nil || upd.Name != "updated" {
		t.Fatalf("Update failed: %v %+v", err, upd)
	}
	if delErr := client.WebhookSubscriptions().Delete(context.Background(), "w1"); delErr != nil {
		t.Fatalf("Delete failed: %v", delErr)
	}
	keyResp, err := client.WebhookSubscriptions().RegenerateSharedKey(context.Background(), "w1")
	if err != nil || keyResp.SharedKey != "new-key" {
		t.Fatalf("RegenerateSharedKey failed: %v %+v", err, keyResp)
	}
}

func TestWebhookSubscriptionsService_Validation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = w
		_ = r
		t.Fatalf("handler should not be called")
	})

	if _, err := client.WebhookSubscriptions().Create(context.Background(), nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if _, err := client.WebhookSubscriptions().Update(context.Background(), "w1", nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if _, err := client.WebhookSubscriptions().Get(context.Background(), ""); err == nil {
		t.Fatalf("expected empty id error")
	}
	if err := client.WebhookSubscriptions().Delete(context.Background(), ""); err == nil {
		t.Fatalf("expected empty id error")
	}
	if _, err := client.WebhookSubscriptions().RegenerateSharedKey(context.Background(), ""); err == nil {
		t.Fatalf("expected empty id error")
	}
}

func TestWebhookEventsService_AllMethods(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/public/v1/webhook-events":
			if r.URL.Query().Get("type") != "document_updated" {
				t.Fatalf("expected type filter")
			}
			_, _ = io.WriteString(w, `{"items":[{"uuid":"e1","name":"event"}]}`)
		case r.Method == http.MethodGet && r.URL.Path == "/public/v1/webhook-events/e1":
			_, _ = io.WriteString(w, `{"uuid":"e1","name":"event","error":false}`)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	list, err := client.WebhookEvents().List(context.Background(), &ListWebhookEventsOptions{Type: "document_updated", HTTPStatusCode: 200, Error: ptrBool(false)})
	if err != nil || len(list.Items) != 1 {
		t.Fatalf("List failed: %v %+v", err, list)
	}

	details, err := client.WebhookEvents().Get(context.Background(), "e1")
	if err != nil || details.UUID != "e1" {
		t.Fatalf("Get failed: %v %+v", err, details)
	}
}

func TestWebhookEventsService_Validation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = w
		_ = r
		t.Fatalf("handler should not be called")
	})

	if _, err := client.WebhookEvents().Get(context.Background(), ""); err == nil {
		t.Fatalf("expected empty id error")
	}
}
