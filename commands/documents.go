// Package commands provides request DTOs for the PandaDoc API.
package commands

// ListDocumentsOptions contains options for listing documents.
type ListDocumentsOptions struct {
	// Page is the page number (1-indexed).
	Page int `json:"page,omitempty"`

	// Count is the number of results per page (max 100, default 50).
	Count int `json:"count,omitempty"`

	// Status filters documents by status.
	// Valid values: document.draft, document.sent, document.completed, document.uploaded, document.error, document.viewed, document.waiting_approval, document.approved, document.rejected, document.waiting_pay, document.paid, document.voided, document.declined, document.external_review
	Status string `json:"status,omitempty"`

	// OrderBy orders results by a field.
	// Valid values: name, date_created, date_modified, date_completed
	OrderBy string `json:"order_by,omitempty"`

	// Ascending orders results in ascending order when true.
	Ascending bool `json:"ascending,omitempty"`

	// TemplateID filters documents by template ID.
	TemplateID string `json:"template_id,omitempty"`

	// FolderUUID filters documents by folder UUID.
	FolderUUID string `json:"folder_uuid,omitempty"`

	// Tag filters documents by tag.
	Tag string `json:"tag,omitempty"`

	// Q is a search query string.
	Q string `json:"q,omitempty"`
}

// UpdateDocument contains fields for updating a document.
type UpdateDocument struct {
	// Name is the new document name.
	Name string `json:"name,omitempty"`

	// Recipients updates the document recipients.
	Recipients []Recipient `json:"recipients,omitempty"`

	// Fields updates the document fields.
	Fields map[string]interface{} `json:"fields,omitempty"`

	// Tokens updates the document tokens.
	Tokens []Token `json:"tokens,omitempty"`

	// Metadata updates the document metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Recipient represents a document recipient.
type Recipient struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role,omitempty"`
}

// Token represents a document token/variable.
type Token struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
