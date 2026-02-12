package pandadoc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

// documentsService implements DocumentsService.
type documentsService struct {
	client *Client
}

// List lists/searches documents.
func (s *documentsService) List(ctx context.Context, opts *ListDocumentsOptions) (*DocumentListResponse, error) {
	query := url.Values{}
	if opts != nil {
		buildDocumentListQuery(query, opts)
	}

	var out DocumentListResponse
	err := s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/documents",
		query:       query,
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func buildDocumentListQuery(query url.Values, opts *ListDocumentsOptions) {
	setIfNotEmpty(query, "template_id", opts.TemplateID)
	setIfNotEmpty(query, "form_id", opts.FormID)
	setIfNotEmpty(query, "folder_uuid", opts.FolderUUID)
	setIfNotEmpty(query, "contact_id", opts.ContactID)
	setIfPositive(query, "count", opts.Count)
	setIfPositive(query, "page", opts.Page)
	setIfNotEmpty(query, "order_by", string(opts.OrderBy))
	setIfNotEmpty(query, "created_from", opts.CreatedFrom)
	setIfNotEmpty(query, "created_to", opts.CreatedTo)
	setIfNotNil(query, "deleted", opts.Deleted)
	setIfNotEmpty(query, "id", opts.ID)
	setIfNotEmpty(query, "completed_from", opts.CompletedFrom)
	setIfNotEmpty(query, "completed_to", opts.CompletedTo)
	setIfNotEmpty(query, "membership_id", opts.MembershipID)
	setMetadataIfNotEmpty(query, opts.Metadata)
	setIfNotEmpty(query, "modified_from", opts.ModifiedFrom)
	setIfNotEmpty(query, "modified_to", opts.ModifiedTo)
	setIfNotEmpty(query, "q", opts.Q)
	setStatusIfNotNil(query, "status", opts.Status)
	setStatusIfNotNil(query, "status__ne", opts.StatusNot)
	setIfNotEmpty(query, "tag", opts.Tag)
}

func setIfNotEmpty(query url.Values, key, value string) {
	if value != "" {
		query.Set(key, value)
	}
}

func setIfPositive(query url.Values, key string, value int) {
	if value > 0 {
		query.Set(key, strconv.Itoa(value))
	}
}

func setIfNotNil(query url.Values, key string, value *bool) {
	if value != nil {
		query.Set(key, strconv.FormatBool(*value))
	}
}

func setStatusIfNotNil(query url.Values, key string, value *DocumentStatusCode) {
	if value != nil {
		query.Set(key, strconv.Itoa(int(*value)))
	}
}

func setMetadataIfNotEmpty(query url.Values, metadata map[string]string) {
	if len(metadata) == 0 {
		return
	}
	keys := make([]string, 0, len(metadata))
	for k := range metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		query.Add("metadata", "metadata_"+k+"="+metadata[k])
	}
}

// Create creates a new document from JSON payload.
func (s *documentsService) Create(ctx context.Context, reqBody DocumentCreateRequest) (*DocumentCreateResponse, error) {
	if reqBody == nil {
		return nil, ErrNilRequest
	}

	var out DocumentCreateResponse
	err := s.client.decodeJSON(ctx, &request{
		method:         http.MethodPost,
		path:           "/public/v1/documents",
		requireAuth:    true,
		jsonBody:       reqBody,
		expectedStatus: []int{http.StatusCreated},
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// CreateFromUpload creates a new document from multipart upload payload.
func (s *documentsService) CreateFromUpload(ctx context.Context, reqBody *CreateDocumentFromUploadRequest) (*DocumentCreateResponse, error) {
	if reqBody == nil {
		return nil, ErrNilRequest
	}
	if reqBody.File == nil {
		return nil, ErrNilFileReader
	}

	fieldName := reqBody.FileField
	if fieldName == "" {
		fieldName = "file"
	}

	var out DocumentCreateResponse
	err := s.client.decodeJSON(ctx, &request{
		method:      http.MethodPost,
		path:        "/public/v1/documents?upload",
		requireAuth: true,
		multipart: &multipartPayload{
			Fields: reqBody.Fields,
			Files: []multipartFile{{
				FieldName: fieldName,
				FileName:  reqBody.FileName,
				Reader:    reqBody.File,
			}},
		},
		expectedStatus: []int{http.StatusCreated},
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Status returns a document status payload.
func (s *documentsService) Status(ctx context.Context, id string) (*DocumentStatusResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}

	var out DocumentStatusResponse
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/documents/" + escapedID,
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete deletes a document.
func (s *documentsService) Delete(ctx context.Context, id string) error {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return err
	}

	return s.client.decodeJSON(ctx, &request{
		method:         http.MethodDelete,
		path:           "/public/v1/documents/" + escapedID,
		requireAuth:    true,
		expectedStatus: []int{http.StatusNoContent},
	}, nil)
}

// Update updates a document and returns no payload on success.
func (s *documentsService) Update(ctx context.Context, id string, reqBody DocumentUpdateRequest) error {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return err
	}
	if reqBody == nil {
		return ErrNilRequest
	}

	return s.client.decodeJSON(ctx, &request{
		method:         http.MethodPatch,
		path:           "/public/v1/documents/" + escapedID,
		requireAuth:    true,
		jsonBody:       reqBody,
		expectedStatus: []int{http.StatusNoContent},
	}, nil)
}

// ESignDisclosure gets e-sign disclosure settings for a document.
func (s *documentsService) ESignDisclosure(ctx context.Context, documentID string) (*DocumentESignDisclosureResponse, error) {
	escapedID, err := escapePathParam(documentID)
	if err != nil {
		return nil, err
	}

	var out DocumentESignDisclosureResponse
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/documents/" + escapedID + "/esign-disclosure",
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// ChangeStatus changes a document status.
func (s *documentsService) ChangeStatus(ctx context.Context, id string, reqBody *ChangeDocumentStatusRequest) error {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return err
	}
	if reqBody == nil {
		return ErrNilRequest
	}

	return s.client.decodeJSON(ctx, &request{
		method:         http.MethodPatch,
		path:           "/public/v1/documents/" + escapedID + "/status",
		requireAuth:    true,
		jsonBody:       reqBody,
		expectedStatus: []int{http.StatusNoContent},
	}, nil)
}

// ChangeStatusWithUpload changes status with multipart payload.
func (s *documentsService) ChangeStatusWithUpload(ctx context.Context, id string, reqBody *ChangeDocumentStatusWithUploadRequest) error {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return err
	}
	if reqBody == nil {
		return ErrNilRequest
	}
	if reqBody.File == nil {
		return ErrNilFileReader
	}

	fieldName := reqBody.FileField
	if fieldName == "" {
		fieldName = "file"
	}

	fields := map[string]string{
		"status": strconv.Itoa(int(reqBody.Status)),
	}
	for k, v := range reqBody.Fields {
		fields[k] = v
	}
	if reqBody.Note != "" {
		fields["note"] = reqBody.Note
	}
	if reqBody.NotifyRecipients != nil {
		fields["notify_recipients"] = strconv.FormatBool(*reqBody.NotifyRecipients)
	}

	return s.client.decodeJSON(ctx, &request{
		method:      http.MethodPatch,
		path:        "/public/v1/documents/" + escapedID + "/status?upload",
		requireAuth: true,
		multipart: &multipartPayload{
			Fields: fields,
			Files: []multipartFile{{
				FieldName: fieldName,
				FileName:  reqBody.FileName,
				Reader:    reqBody.File,
			}},
		},
		expectedStatus: []int{http.StatusNoContent},
	}, nil)
}

// RevertToDraft reverts a document to draft.
func (s *documentsService) RevertToDraft(ctx context.Context, id string) (*DocumentRevertToDraftResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}

	var out DocumentRevertToDraftResponse
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodPost,
		path:        "/public/v1/documents/" + escapedID + "/draft",
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Details returns document details.
func (s *documentsService) Details(ctx context.Context, id string) (*DocumentDetailsResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}

	var out DocumentDetailsResponse
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/documents/" + escapedID + "/details",
		requireAuth: true,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Send sends a document.
func (s *documentsService) Send(ctx context.Context, id string, reqBody DocumentSendRequest) (*DocumentSendResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}
	if reqBody == nil {
		return nil, ErrNilRequest
	}

	var out DocumentSendResponse
	err = s.client.decodeJSON(ctx, &request{
		method:      http.MethodPost,
		path:        "/public/v1/documents/" + escapedID + "/send",
		requireAuth: true,
		jsonBody:    reqBody,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateEditingSession creates an editing session for a document.
func (s *documentsService) CreateEditingSession(ctx context.Context, id string, reqBody CreateDocumentEditingSessionRequest) (*CreateDocumentEditingSessionResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}
	if reqBody == nil {
		return nil, ErrNilRequest
	}

	var out CreateDocumentEditingSessionResponse
	err = s.client.decodeJSON(ctx, &request{
		method:         http.MethodPost,
		path:           "/public/v1/documents/" + escapedID + "/editing-sessions",
		requireAuth:    true,
		jsonBody:       reqBody,
		expectedStatus: []int{http.StatusCreated},
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateSession creates an embedded signing session for a document.
func (s *documentsService) CreateSession(ctx context.Context, id string, reqBody CreateDocumentSessionRequest) (*CreateDocumentSessionResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}
	if reqBody == nil {
		return nil, ErrNilRequest
	}

	var out CreateDocumentSessionResponse
	err = s.client.decodeJSON(ctx, &request{
		method:         http.MethodPost,
		path:           "/public/v1/documents/" + escapedID + "/session",
		requireAuth:    true,
		jsonBody:       reqBody,
		expectedStatus: []int{http.StatusCreated},
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Download downloads a completed document PDF.
func (s *documentsService) Download(ctx context.Context, id string) (*DownloadResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}

	return s.client.download(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/documents/" + escapedID + "/download",
		requireAuth: true,
		accept:      "application/pdf",
	})
}

// DownloadProtected downloads a password-protected document PDF.
func (s *documentsService) DownloadProtected(ctx context.Context, id string) (*DownloadResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}

	return s.client.download(ctx, &request{
		method:      http.MethodGet,
		path:        "/public/v1/documents/" + escapedID + "/download-protected",
		requireAuth: true,
		accept:      "application/pdf",
	})
}

// TransferOwnership transfers ownership of a single document.
func (s *documentsService) TransferOwnership(ctx context.Context, id string, reqBody TransferDocumentOwnershipRequest) error {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return err
	}
	if reqBody == nil {
		return ErrNilRequest
	}

	return s.client.decodeJSON(ctx, &request{
		method:         http.MethodPatch,
		path:           "/public/v1/documents/" + escapedID + "/ownership",
		requireAuth:    true,
		jsonBody:       reqBody,
		expectedStatus: []int{http.StatusNoContent},
	}, nil)
}

// TransferAllOwnership transfers ownership for all documents.
func (s *documentsService) TransferAllOwnership(ctx context.Context, reqBody TransferAllDocumentsOwnershipRequest) error {
	if reqBody == nil {
		return ErrNilRequest
	}

	return s.client.decodeJSON(ctx, &request{
		method:         http.MethodPatch,
		path:           "/public/v1/documents/ownership",
		requireAuth:    true,
		jsonBody:       reqBody,
		expectedStatus: []int{http.StatusNoContent},
	}, nil)
}

// MoveToFolder moves a document to a folder.
func (s *documentsService) MoveToFolder(ctx context.Context, id, folderID string) error {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return err
	}
	escapedFolderID, err := escapePathParam(folderID)
	if err != nil {
		return fmt.Errorf("folder id: %w", err)
	}

	return s.client.decodeJSON(ctx, &request{
		method:         http.MethodPost,
		path:           "/public/v1/documents/" + escapedID + "/move-to-folder/" + escapedFolderID,
		requireAuth:    true,
		expectedStatus: []int{http.StatusNoContent},
	}, nil)
}

// AppendContentLibraryItem appends a content library item to a document.
func (s *documentsService) AppendContentLibraryItem(ctx context.Context, id string, reqBody AppendContentLibraryItemRequest) (*AppendContentLibraryItemResponse, error) {
	escapedID, err := escapePathParam(id)
	if err != nil {
		return nil, err
	}
	if reqBody == nil {
		return nil, ErrNilRequest
	}

	var out AppendContentLibraryItemResponse
	err = s.client.decodeJSON(ctx, &request{
		method:         http.MethodPost,
		path:           "/public/v1/documents/" + escapedID + "/append-content-library-item",
		requireAuth:    true,
		jsonBody:       reqBody,
		expectedStatus: []int{http.StatusCreated},
	}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
