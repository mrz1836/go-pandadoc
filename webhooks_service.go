package pandadoc

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// WebhookSubscriptionsService handles webhook subscription APIs.
type WebhookSubscriptionsService struct {
	client *Client
}

// List lists webhook subscriptions.
func (s *WebhookSubscriptionsService) List(ctx context.Context, opts *ListWebhookSubscriptionsOptions) (*WebhookSubscriptionListResponse, error) {
	query := url.Values{}
	if opts == nil {
		opts = &ListWebhookSubscriptionsOptions{}
	}
	if opts.Count > 0 {
		query.Set("count", strconv.Itoa(opts.Count))
	}
	if opts.Page > 0 {
		query.Set("page", strconv.Itoa(opts.Page))
	}

	var out WebhookSubscriptionListResponse
	err := s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/webhook-subscriptions",
		query:       query,
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a webhook subscription.
func (s *WebhookSubscriptionsService) Create(ctx context.Context, reqBody *WebhookSubscriptionRequest) (*WebhookSubscription, error) {
	if reqBody == nil {
		return nil, ErrNilRequest
	}

	var out WebhookSubscription
	err := s.client.decodeJSON(ctx, &request{
		method:         http.MethodPost,
		path:           "/public/v1/webhook-subscriptions",
		requireAuth:    true,
		jsonBody:       reqBody,
		expectedStatus: []int{http.StatusCreated},
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets a webhook subscription by ID.
func (s *WebhookSubscriptionsService) Get(ctx context.Context, id string) (*WebhookSubscription, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}

	var out WebhookSubscription
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/webhook-subscriptions/" + escapedID,
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates a webhook subscription.
func (s *WebhookSubscriptionsService) Update(ctx context.Context, id string, reqBody *WebhookSubscriptionRequest) (*WebhookSubscription, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}
	if reqBody == nil {
		return nil, ErrNilRequest
	}

	var out WebhookSubscription
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodPatch,
		path:        "/public/v1/webhook-subscriptions/" + escapedID,
		requireAuth: true,
		jsonBody:    reqBody,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a webhook subscription.
func (s *WebhookSubscriptionsService) Delete(ctx context.Context, id string) error {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return err
	}

	return s.client.decodeJSON(ctx, &request{
		method:         http.MethodDelete,
		path:           "/public/v1/webhook-subscriptions/" + escapedID,
		requireAuth:    true,
		expectedStatus: []int{http.StatusNoContent},
	}, nil)
}

// RegenerateSharedKey regenerates a webhook subscription shared key.
func (s *WebhookSubscriptionsService) RegenerateSharedKey(ctx context.Context, id string) (*UpdateWebhookSubscriptionSharedKeyResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}

	var out UpdateWebhookSubscriptionSharedKeyResponse
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodPatch,
		path:        "/public/v1/webhook-subscriptions/" + escapedID + "/shared-key",
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// WebhookEventsService handles webhook event APIs.
type WebhookEventsService struct {
	client *Client
}

// List lists webhook events.
func (s *WebhookEventsService) List(ctx context.Context, opts *ListWebhookEventsOptions) (*WebhookEventListResponse, error) {
	query := url.Values{}
	if opts == nil {
		opts = &ListWebhookEventsOptions{}
	}
	if opts.Since != "" {
		query.Set("since", opts.Since)
	}
	if opts.To != "" {
		query.Set("to", opts.To)
	}
	if opts.Type != "" {
		query.Set("type", opts.Type)
	}
	if opts.HTTPStatusCode > 0 {
		query.Set("http_status_code", strconv.Itoa(opts.HTTPStatusCode))
	}
	if opts.Error != nil {
		query.Set("error", strconv.FormatBool(*opts.Error))
	}

	var out WebhookEventListResponse
	err := s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/webhook-events",
		query:       query,
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Get gets webhook event details.
func (s *WebhookEventsService) Get(ctx context.Context, id string) (*WebhookEventDetailsResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}

	var out WebhookEventDetailsResponse
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/webhook-events/" + escapedID,
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
