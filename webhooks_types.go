package pandadoc

// WebhookPayloadOption controls additional payload sections in webhook deliveries.
type WebhookPayloadOption string

// Webhook payload option constants.
const (
	// WebhookPayloadMetadata includes metadata in webhook payload.
	WebhookPayloadMetadata WebhookPayloadOption = "metadata"
	// WebhookPayloadFields includes fields in webhook payload.
	WebhookPayloadFields WebhookPayloadOption = "fields"
	// WebhookPayloadProducts includes products in webhook payload.
	WebhookPayloadProducts WebhookPayloadOption = "products"
	// WebhookPayloadTokens includes tokens in webhook payload.
	WebhookPayloadTokens WebhookPayloadOption = "tokens"
	// WebhookPayloadPricing includes pricing in webhook payload.
	WebhookPayloadPricing WebhookPayloadOption = "pricing"
)

// WebhookTrigger controls which events trigger subscription deliveries.
type WebhookTrigger string

// Webhook trigger constants.
const (
	// WebhookTriggerRecipientCompleted triggers when a recipient completes.
	WebhookTriggerRecipientCompleted WebhookTrigger = "recipient_completed"
	// WebhookTriggerDocumentUpdated triggers when a document is updated.
	WebhookTriggerDocumentUpdated WebhookTrigger = "document_updated"
	// WebhookTriggerDocumentDeleted triggers when a document is deleted.
	WebhookTriggerDocumentDeleted WebhookTrigger = "document_deleted"
	// WebhookTriggerDocumentStateChanged triggers when document state changes.
	WebhookTriggerDocumentStateChanged WebhookTrigger = "document_state_changed"
	// WebhookTriggerDocumentCreationFailed triggers when document creation fails.
	WebhookTriggerDocumentCreationFailed WebhookTrigger = "document_creation_failed"
	// WebhookTriggerDocumentCompletedPDFReady triggers when completed PDF is ready.
	WebhookTriggerDocumentCompletedPDFReady WebhookTrigger = "document_completed_pdf_ready"
	// WebhookTriggerDocumentSectionAdded triggers when a document section is added.
	WebhookTriggerDocumentSectionAdded WebhookTrigger = "document_section_added"
	// WebhookTriggerQuoteUpdated triggers when a quote is updated.
	WebhookTriggerQuoteUpdated WebhookTrigger = "quote_updated"
	// WebhookTriggerTemplateCreated triggers when a template is created.
	WebhookTriggerTemplateCreated WebhookTrigger = "template_created"
	// WebhookTriggerTemplateUpdated triggers when a template is updated.
	WebhookTriggerTemplateUpdated WebhookTrigger = "template_updated"
	// WebhookTriggerTemplateDeleted triggers when a template is deleted.
	WebhookTriggerTemplateDeleted WebhookTrigger = "template_deleted"
	// WebhookTriggerContentLibraryItemCreated triggers when content library item is created.
	WebhookTriggerContentLibraryItemCreated WebhookTrigger = "content_library_item_created"
	// WebhookTriggerContentLibraryItemCreationFail triggers when content library creation fails.
	WebhookTriggerContentLibraryItemCreationFail WebhookTrigger = "content_library_item_creation_failed"
)

// ListWebhookSubscriptionsOptions configures webhook-subscription listing.
type ListWebhookSubscriptionsOptions struct {
	Count int
	Page  int
}

// WebhookSubscriptionRequest creates/updates webhook subscriptions.
type WebhookSubscriptionRequest struct {
	Name     string                 `json:"name,omitempty"`
	URL      string                 `json:"url,omitempty"`
	Active   *bool                  `json:"active,omitempty"`
	Triggers []WebhookTrigger       `json:"triggers,omitempty"`
	Payload  []WebhookPayloadOption `json:"payload,omitempty"`
}

// WebhookSubscription represents a webhook subscription item.
type WebhookSubscription struct {
	UUID      string                 `json:"uuid,omitempty"`
	Workspace string                 `json:"workspace_id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	URL       string                 `json:"url,omitempty"`
	Active    bool                   `json:"active"`
	Status    string                 `json:"status,omitempty"`
	SharedKey string                 `json:"shared_key,omitempty"`
	Triggers  []WebhookTrigger       `json:"triggers,omitempty"`
	Payload   []WebhookPayloadOption `json:"payload,omitempty"`
}

// WebhookSubscriptionListResponse lists webhook subscriptions.
type WebhookSubscriptionListResponse struct {
	Items []WebhookSubscription `json:"items"`
}

// UpdateWebhookSubscriptionSharedKeyResponse contains regenerated shared-key.
type UpdateWebhookSubscriptionSharedKeyResponse struct {
	SharedKey string `json:"shared_key"`
}

// ListWebhookEventsOptions configures webhook event listing.
type ListWebhookEventsOptions struct {
	Since          string
	To             string
	Type           string
	HTTPStatusCode int
	Error          *bool
}

// WebhookEventItem is a compact event-list entry.
type WebhookEventItem struct {
	UUID           string `json:"uuid,omitempty"`
	Name           string `json:"name,omitempty"`
	Type           string `json:"type,omitempty"`
	HTTPStatusCode int    `json:"http_status_code,omitempty"`
	DeliveryTime   string `json:"delivery_time,omitempty"`
	Error          bool   `json:"error"`
}

// WebhookEventListResponse represents webhook event pages.
type WebhookEventListResponse struct {
	Items []WebhookEventItem `json:"items"`
}

// WebhookEventDetailsResponse represents a single webhook event payload.
type WebhookEventDetailsResponse struct {
	UUID            string  `json:"uuid,omitempty"`
	Name            string  `json:"name,omitempty"`
	Type            string  `json:"type,omitempty"`
	EventTime       string  `json:"event_time,omitempty"`
	DeliveryTime    string  `json:"delivery_time,omitempty"`
	URL             string  `json:"url,omitempty"`
	HTTPStatusCode  int     `json:"http_status_code,omitempty"`
	Error           bool    `json:"error"`
	RequestBody     RawJSON `json:"request_body,omitempty"`
	ResponseBody    RawJSON `json:"response_body,omitempty"`
	ResponseHeaders RawJSON `json:"response_headers,omitempty"`
	Signature       string  `json:"signature,omitempty"`
}
