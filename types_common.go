package pandadoc

import "encoding/json"

// RawObject is a flexible JSON object shape for endpoints with wide payload variance.
type RawObject map[string]any

// APIListResponse models list responses that return items.
type APIListResponse[T any] struct {
	Items []T `json:"items"`
}

// UserReference is a compact user descriptor used across multiple endpoints.
type UserReference struct {
	ID           string `json:"id,omitempty"`
	MembershipID string `json:"membership_id,omitempty"`
	Email        string `json:"email,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
}

// MoneyAmount represents an amount/currency pair.
type MoneyAmount struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// NamedContentBlock is used by document detail payloads for image/table/text blocks.
type NamedContentBlock struct {
	Name string `json:"name,omitempty"`
}

// RawJSON is a helper alias for nested payloads that are intentionally untyped.
type RawJSON = json.RawMessage
