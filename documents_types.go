package pandadoc

import "io"

// DocumentStatusCode is the numeric status code used in some document requests.
type DocumentStatusCode int

// Document status code constants.
const (
	// DocumentStatusDraft represents the draft status.
	DocumentStatusDraft DocumentStatusCode = 0
	// DocumentStatusSent represents the sent status.
	DocumentStatusSent DocumentStatusCode = 1
	// DocumentStatusCompleted represents the completed status.
	DocumentStatusCompleted DocumentStatusCode = 2
	// DocumentStatusUploaded represents the uploaded status.
	DocumentStatusUploaded DocumentStatusCode = 3
	// DocumentStatusError represents the error status.
	DocumentStatusError DocumentStatusCode = 4
	// DocumentStatusViewed represents the viewed status.
	DocumentStatusViewed DocumentStatusCode = 5
	// DocumentStatusWaitingApproval represents the waiting approval status.
	DocumentStatusWaitingApproval DocumentStatusCode = 6
	// DocumentStatusApproved represents the approved status.
	DocumentStatusApproved DocumentStatusCode = 7
	// DocumentStatusRejected represents the rejected status.
	DocumentStatusRejected DocumentStatusCode = 8
	// DocumentStatusWaitingPay represents the waiting pay status.
	DocumentStatusWaitingPay DocumentStatusCode = 9
	// DocumentStatusPaid represents the paid status.
	DocumentStatusPaid DocumentStatusCode = 10
	// DocumentStatusVoided represents the voided status.
	DocumentStatusVoided DocumentStatusCode = 11
	// DocumentStatusDeclined represents the declined status.
	DocumentStatusDeclined DocumentStatusCode = 12
	// DocumentStatusExternalReview represents the external review status.
	DocumentStatusExternalReview DocumentStatusCode = 13
)

// DocumentOrderBy controls document list ordering.
type DocumentOrderBy string

// Document order by constants.
const (
	// DocumentOrderByName orders by document name.
	DocumentOrderByName DocumentOrderBy = "name"
	// DocumentOrderByDateCreated orders by creation date.
	DocumentOrderByDateCreated DocumentOrderBy = "date_created"
	// DocumentOrderByDateStatusChange orders by status change date.
	DocumentOrderByDateStatusChange DocumentOrderBy = "date_status_changed"
	// DocumentOrderByDateLastAction orders by last action date.
	DocumentOrderByDateLastAction DocumentOrderBy = "date_of_last_action"
	// DocumentOrderByDateModified orders by modification date.
	DocumentOrderByDateModified DocumentOrderBy = "date_modified"
	// DocumentOrderByDateSent orders by sent date.
	DocumentOrderByDateSent DocumentOrderBy = "date_sent"
	// DocumentOrderByDateCompleted orders by completion date.
	DocumentOrderByDateCompleted DocumentOrderBy = "date_completed"
	// DocumentOrderByDateExpiration orders by expiration date.
	DocumentOrderByDateExpiration DocumentOrderBy = "date_expiration"
	// DocumentOrderByDateDeclined orders by declined date.
	DocumentOrderByDateDeclined DocumentOrderBy = "date_declined"
	// DocumentOrderByStatus orders by status.
	DocumentOrderByStatus DocumentOrderBy = "status"
	// DocumentOrderByNameDesc orders by document name descending.
	DocumentOrderByNameDesc DocumentOrderBy = "-name"
	// DocumentOrderByDateCreatedDesc orders by creation date descending.
	DocumentOrderByDateCreatedDesc DocumentOrderBy = "-date_created"
	// DocumentOrderByStatusChangeDesc orders by status change date descending.
	DocumentOrderByStatusChangeDesc DocumentOrderBy = "-date_status_changed"
	// DocumentOrderByLastActionDesc orders by last action date descending.
	DocumentOrderByLastActionDesc DocumentOrderBy = "-date_of_last_action"
	// DocumentOrderByDateModifiedDesc orders by modification date descending.
	DocumentOrderByDateModifiedDesc DocumentOrderBy = "-date_modified"
	// DocumentOrderByDateSentDesc orders by sent date descending.
	DocumentOrderByDateSentDesc DocumentOrderBy = "-date_sent"
	// DocumentOrderByDateCompletedDes orders by completion date descending.
	DocumentOrderByDateCompletedDes DocumentOrderBy = "-date_completed"
	// DocumentOrderByDateExpirationDe orders by expiration date descending.
	DocumentOrderByDateExpirationDe DocumentOrderBy = "-date_expiration"
	// DocumentOrderByDateDeclinedDesc orders by declined date descending.
	DocumentOrderByDateDeclinedDesc DocumentOrderBy = "-date_declined"
	// DocumentOrderByStatusDesc orders by status descending.
	DocumentOrderByStatusDesc DocumentOrderBy = "-status"
)

// ListDocumentsOptions controls list/search behavior for documents.
type ListDocumentsOptions struct {
	TemplateID    string
	FormID        string
	FolderUUID    string
	ContactID     string
	Count         int
	Page          int
	OrderBy       DocumentOrderBy
	CreatedFrom   string
	CreatedTo     string
	Deleted       *bool
	ID            string
	CompletedFrom string
	CompletedTo   string
	MembershipID  string
	Metadata      map[string]string
	ModifiedFrom  string
	ModifiedTo    string
	Q             string
	Status        *DocumentStatusCode
	StatusNot     *DocumentStatusCode
	Tag           string
}

// DocumentSummary represents core document fields used by multiple endpoints.
type DocumentSummary struct {
	ID             string `json:"id,omitempty"`
	UUID           string `json:"uuid,omitempty"`
	Name           string `json:"name,omitempty"`
	Status         string `json:"status,omitempty"`
	DateCreated    string `json:"date_created,omitempty"`
	DateModified   string `json:"date_modified,omitempty"`
	DateCompleted  string `json:"date_completed,omitempty"`
	ExpirationDate string `json:"expiration_date,omitempty"`
	Version        string `json:"version,omitempty"`
}

// DocumentListResponse is returned by list/search documents endpoint.
type DocumentListResponse struct {
	Results []DocumentSummary `json:"results"`
}

// DocumentCreateRequest is a flexible create-document payload.
type DocumentCreateRequest map[string]any

// DocumentUpdateRequest is a flexible update-document payload.
type DocumentUpdateRequest map[string]any

// DocumentCreateResponse is returned when creating a document.
type DocumentCreateResponse struct {
	DocumentSummary

	Links       []DocumentLink `json:"links,omitempty"`
	InfoMessage string         `json:"info_message,omitempty"`
}

// DocumentLink is a link object from create responses.
type DocumentLink struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
	Type string `json:"type,omitempty"`
}

// DocumentStatusResponse is returned by GET /documents/{id}.
type DocumentStatusResponse = DocumentSummary

// DocumentRevertToDraftResponse is returned by revert-to-draft endpoint.
type DocumentRevertToDraftResponse = DocumentSummary

// DocumentField represents a field entry in document details.
type DocumentField struct {
	UUID        string  `json:"uuid,omitempty"`
	Name        string  `json:"name,omitempty"`
	Title       string  `json:"title,omitempty"`
	MergeField  string  `json:"merge_field,omitempty"`
	Placeholder string  `json:"placeholder,omitempty"`
	FieldID     string  `json:"field_id,omitempty"`
	Type        string  `json:"type,omitempty"`
	Value       any     `json:"value,omitempty"`
	AssignedTo  RawJSON `json:"assigned_to,omitempty"`
}

// DocumentToken is a token/value pair from document details.
type DocumentToken struct {
	Name  string `json:"name,omitempty"`
	Value any    `json:"value,omitempty"`
}

// DocumentRecipient is a compact recipient payload.
type DocumentRecipient struct {
	ID            string `json:"id,omitempty"`
	ContactID     string `json:"contact_id,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Email         string `json:"email,omitempty"`
	Role          string `json:"role,omitempty"`
	RecipientType string `json:"recipient_type,omitempty"`
	HasCompleted  bool   `json:"has_completed,omitempty"`
}

// LinkedObject represents a linked CRM object reference.
type LinkedObject struct {
	ID         string `json:"id,omitempty"`
	Provider   string `json:"provider,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityID   string `json:"entity_id,omitempty"`
}

// DocumentTemplateReference is a template descriptor in details responses.
type DocumentTemplateReference struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// DocumentDetailsResponse is returned by GET /documents/{id}/details.
type DocumentDetailsResponse struct {
	ApprovalExecution               any                        `json:"approval_execution,omitempty"`
	AutonumberingSequenceNamePrefix string                     `json:"autonumbering_sequence_name_prefix,omitempty"`
	ContentDateModified             string                     `json:"content_date_modified,omitempty"`
	CreatedBy                       *UserReference             `json:"created_by,omitempty"`
	DateCompleted                   string                     `json:"date_completed,omitempty"`
	DateCreated                     string                     `json:"date_created,omitempty"`
	DateModified                    string                     `json:"date_modified,omitempty"`
	DateSent                        string                     `json:"date_sent,omitempty"`
	ExpirationDate                  string                     `json:"expiration_date,omitempty"`
	Fields                          []DocumentField            `json:"fields,omitempty"`
	FolderUUID                      string                     `json:"folder_uuid,omitempty"`
	GrandTotal                      *MoneyAmount               `json:"grand_total,omitempty"`
	ID                              string                     `json:"id,omitempty"`
	Images                          []NamedContentBlock        `json:"images,omitempty"`
	LinkedObjects                   []LinkedObject             `json:"linked_objects,omitempty"`
	Metadata                        map[string]any             `json:"metadata,omitempty"`
	Name                            string                     `json:"name,omitempty"`
	Pricing                         RawJSON                    `json:"pricing,omitempty"`
	Recipients                      []DocumentRecipient        `json:"recipients,omitempty"`
	RefNumber                       string                     `json:"ref_number,omitempty"`
	SentBy                          *UserReference             `json:"sent_by,omitempty"`
	Status                          string                     `json:"status,omitempty"`
	Tables                          []NamedContentBlock        `json:"tables,omitempty"`
	Tags                            []string                   `json:"tags,omitempty"`
	Template                        *DocumentTemplateReference `json:"template,omitempty"`
	Texts                           []NamedContentBlock        `json:"texts,omitempty"`
	Tokens                          []DocumentToken            `json:"tokens,omitempty"`
	Version                         string                     `json:"version,omitempty"`
}

// DocumentESignDisclosureResponse models e-sign disclosure settings for a document.
type DocumentESignDisclosureResponse struct {
	Result *DocumentESignDisclosure `json:"result,omitempty"`
}

// DocumentESignDisclosure contains disclosure details.
type DocumentESignDisclosure struct {
	IsEnabled           bool   `json:"is_enabled"`
	CompanyName         string `json:"company_name,omitempty"`
	ESignDisclosureText string `json:"esign_disclosure_text,omitempty"`
}

// ChangeDocumentStatusRequest changes a document status.
type ChangeDocumentStatusRequest struct {
	Status           DocumentStatusCode `json:"status"`
	Note             string             `json:"note,omitempty"`
	NotifyRecipients *bool              `json:"notify_recipients,omitempty"`
}

// ChangeDocumentStatusWithUploadRequest changes status with multipart payload.
type ChangeDocumentStatusWithUploadRequest struct {
	Status           DocumentStatusCode
	Note             string
	NotifyRecipients *bool
	FileField        string
	FileName         string
	File             io.Reader
	Fields           map[string]string
}

// CreateDocumentFromUploadRequest uploads a file and creates a document.
type CreateDocumentFromUploadRequest struct {
	FileField string
	FileName  string
	File      io.Reader
	Fields    map[string]string
}

// DocumentSendRequest is a flexible send payload.
type DocumentSendRequest map[string]any

// DocumentSendResponse is returned by send endpoint.
type DocumentSendResponse struct {
	DocumentSummary

	Recipients []DocumentRecipient `json:"recipients,omitempty"`
}

// CreateDocumentEditingSessionRequest is a flexible editing-session payload.
type CreateDocumentEditingSessionRequest map[string]any

// CreateDocumentEditingSessionResponse models editing-session response.
type CreateDocumentEditingSessionResponse struct {
	ID         string `json:"id,omitempty"`
	Token      string `json:"token,omitempty"`
	Key        string `json:"key,omitempty"`
	Email      string `json:"email,omitempty"`
	ExpiresAt  string `json:"expires_at,omitempty"`
	DocumentID string `json:"document_id,omitempty"`
}

// CreateDocumentSessionRequest is a flexible embedded-session payload.
type CreateDocumentSessionRequest map[string]any

// CreateDocumentSessionResponse is returned by document session endpoint.
type CreateDocumentSessionResponse struct {
	ID        string `json:"id,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// TransferDocumentOwnershipRequest is a flexible ownership transfer payload.
type TransferDocumentOwnershipRequest map[string]any

// TransferAllDocumentsOwnershipRequest is a flexible bulk ownership transfer payload.
type TransferAllDocumentsOwnershipRequest map[string]any

// AppendContentLibraryItemRequest is a flexible append-content payload.
type AppendContentLibraryItemRequest map[string]any

// AppendContentLibraryItemResponse models append-content response.
type AppendContentLibraryItemResponse struct {
	BlockMapping map[string]string `json:"block_mapping,omitempty"`
	CLI          RawObject         `json:"cli,omitempty"`
}
