package pandadoc

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

func TestDocumentsService_List_AllFilters(t *testing.T) {
	t.Parallel()

	status := DocumentStatusCompleted
	statusNe := DocumentStatusDeclined
	deleted := true

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/public/v1/documents" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		q := r.URL.Query()
		assertQueryEq(t, q, "template_id", "tpl")
		assertQueryEq(t, q, "form_id", "form")
		assertQueryEq(t, q, "folder_uuid", "folder")
		assertQueryEq(t, q, "contact_id", "contact")
		assertQueryEq(t, q, "count", "50")
		assertQueryEq(t, q, "page", "2")
		assertQueryEq(t, q, "order_by", string(DocumentOrderByDateCreatedDesc))
		assertQueryEq(t, q, "created_from", "2024-01-01T00:00:00Z")
		assertQueryEq(t, q, "created_to", "2024-01-02T00:00:00Z")
		assertQueryEq(t, q, "deleted", "true")
		assertQueryEq(t, q, "id", "doc-id")
		assertQueryEq(t, q, "completed_from", "2024-01-03T00:00:00Z")
		assertQueryEq(t, q, "completed_to", "2024-01-04T00:00:00Z")
		assertQueryEq(t, q, "membership_id", "mem")
		assertQueryEq(t, q, "modified_from", "2024-01-05T00:00:00Z")
		assertQueryEq(t, q, "modified_to", "2024-01-06T00:00:00Z")
		assertQueryEq(t, q, "q", "query")
		assertQueryEq(t, q, "status", "2")
		assertQueryEq(t, q, "status__ne", "12")
		assertQueryEq(t, q, "tag", "tag")
		meta := q["metadata"]
		if len(meta) != 2 {
			t.Fatalf("expected 2 metadata values, got %v", meta)
		}

		_, _ = io.WriteString(w, `{"results":[{"id":"d1","name":"n1"}]}`)
	}, WithRetryPolicy(RetryPolicy{MaxRetries: 0, InitialBackoff: 1, MaxBackoff: 1}))

	resp, err := client.Documents().List(context.Background(), &ListDocumentsOptions{
		TemplateID:    "tpl",
		FormID:        "form",
		FolderUUID:    "folder",
		ContactID:     "contact",
		Count:         50,
		Page:          2,
		OrderBy:       DocumentOrderByDateCreatedDesc,
		CreatedFrom:   "2024-01-01T00:00:00Z",
		CreatedTo:     "2024-01-02T00:00:00Z",
		Deleted:       &deleted,
		ID:            "doc-id",
		CompletedFrom: "2024-01-03T00:00:00Z",
		CompletedTo:   "2024-01-04T00:00:00Z",
		MembershipID:  "mem",
		Metadata:      map[string]string{"a": "1", "b": "2"},
		ModifiedFrom:  "2024-01-05T00:00:00Z",
		ModifiedTo:    "2024-01-06T00:00:00Z",
		Q:             "query",
		Status:        &status,
		StatusNot:     &statusNe,
		Tag:           "tag",
	})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(resp.Results) != 1 || resp.Results[0].ID != "d1" {
		t.Fatalf("unexpected list response: %+v", resp)
	}
}

func TestDocumentsService_CreateAndStatusAndDetails(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/public/v1/documents":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"doc1","name":"Created"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/public/v1/documents/doc1":
			_, _ = io.WriteString(w, `{"id":"doc1","status":"document.draft"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/public/v1/documents/doc1/details":
			_, _ = io.WriteString(w, `{"id":"doc1","fields":[{"name":"Field1","type":"text","value":"v"}]}`)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	created, err := client.Documents().Create(context.Background(), DocumentCreateRequest{"name": "Created"})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.ID != "doc1" {
		t.Fatalf("unexpected create response: %+v", created)
	}

	statusResp, err := client.Documents().Status(context.Background(), "doc1")
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if statusResp.Status != "document.draft" {
		t.Fatalf("unexpected status: %+v", statusResp)
	}

	details, err := client.Documents().Details(context.Background(), "doc1")
	if err != nil {
		t.Fatalf("Details failed: %v", err)
	}
	if len(details.Fields) != 1 || details.Fields[0].Name != "Field1" {
		t.Fatalf("unexpected details: %+v", details)
	}
}

//nolint:gocognit // Test function that validates multiple upload scenarios
func TestDocumentsService_CreateFromUpload_AndChangeStatusWithUpload(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Fatalf("parse content-type: %v", err)
		}
		if mediaType != "multipart/form-data" {
			t.Fatalf("expected multipart, got %s", mediaType)
		}
		mr := multipart.NewReader(r.Body, params["boundary"])
		seenFile := false
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("read multipart: %v", err)
			}
			if part.FormName() == "file" {
				seenFile = true
			}
		}
		if !seenFile {
			t.Fatalf("expected file part")
		}

		if strings.Contains(r.URL.RawQuery, "upload") && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"u1"}`)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	created, err := client.Documents().CreateFromUpload(context.Background(), &CreateDocumentFromUploadRequest{
		FileField: "file",
		FileName:  "doc.pdf",
		File:      strings.NewReader("pdf-data"),
		Fields:    map[string]string{"name": "Upload"},
	})
	if err != nil {
		t.Fatalf("CreateFromUpload failed: %v", err)
	}
	if created.ID != "u1" {
		t.Fatalf("unexpected create upload response: %+v", created)
	}

	notify := true
	if err := client.Documents().ChangeStatusWithUpload(context.Background(), "u1", &ChangeDocumentStatusWithUploadRequest{
		Status:           DocumentStatusCompleted,
		NotifyRecipients: &notify,
		FileField:        "file",
		FileName:         "evidence.pdf",
		File:             strings.NewReader("pdf-data-2"),
	}); err != nil {
		t.Fatalf("ChangeStatusWithUpload failed: %v", err)
	}
}

func TestDocumentsService_NoContentMutations(t *testing.T) {
	t.Parallel()

	seen := map[string]int{}
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seen[r.Method+" "+r.URL.Path]++
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/public/v1/documents/doc1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPatch && r.URL.Path == "/public/v1/documents/doc1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPatch && r.URL.Path == "/public/v1/documents/doc1/status":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPatch && r.URL.Path == "/public/v1/documents/doc1/ownership":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPatch && r.URL.Path == "/public/v1/documents/ownership":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/public/v1/documents/doc1/move-to-folder/f1":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	if err := client.Documents().Delete(context.Background(), "doc1"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if err := client.Documents().Update(context.Background(), "doc1", DocumentUpdateRequest{"name": "x"}); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if err := client.Documents().ChangeStatus(context.Background(), "doc1", &ChangeDocumentStatusRequest{Status: DocumentStatusCompleted}); err != nil {
		t.Fatalf("ChangeStatus failed: %v", err)
	}
	if err := client.Documents().TransferOwnership(context.Background(), "doc1", TransferDocumentOwnershipRequest{"to": "user"}); err != nil {
		t.Fatalf("TransferOwnership failed: %v", err)
	}
	if err := client.Documents().TransferAllOwnership(context.Background(), TransferAllDocumentsOwnershipRequest{"to": "user"}); err != nil {
		t.Fatalf("TransferAllOwnership failed: %v", err)
	}
	if err := client.Documents().MoveToFolder(context.Background(), "doc1", "f1"); err != nil {
		t.Fatalf("MoveToFolder failed: %v", err)
	}
}

func TestDocumentsService_OtherJSONOperations(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/public/v1/documents/doc1/esign-disclosure":
			_, _ = io.WriteString(w, `{"result":{"is_enabled":true}}`)
		case r.Method == http.MethodPost && r.URL.Path == "/public/v1/documents/doc1/draft":
			_, _ = io.WriteString(w, `{"id":"doc1","status":"document.draft"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/public/v1/documents/doc1/send":
			_, _ = io.WriteString(w, `{"id":"doc1","status":"document.sent"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/public/v1/documents/doc1/editing-sessions":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"s1","token":"t"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/public/v1/documents/doc1/session":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"sess1"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/public/v1/documents/doc1/append-content-library-item":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"block_mapping":{"a":"b"}}`)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	esign, err := client.Documents().ESignDisclosure(context.Background(), "doc1")
	if err != nil || esign.Result == nil || !esign.Result.IsEnabled {
		t.Fatalf("ESignDisclosure failed: %v %+v", err, esign)
	}

	revert, err := client.Documents().RevertToDraft(context.Background(), "doc1")
	if err != nil || revert.Status != "document.draft" {
		t.Fatalf("RevertToDraft failed: %v %+v", err, revert)
	}

	sendResp, err := client.Documents().Send(context.Background(), "doc1", DocumentSendRequest{"silent": true})
	if err != nil || sendResp.Status != "document.sent" {
		t.Fatalf("Send failed: %v %+v", err, sendResp)
	}

	edit, err := client.Documents().CreateEditingSession(context.Background(), "doc1", CreateDocumentEditingSessionRequest{"member": "x"})
	if err != nil || edit.ID != "s1" {
		t.Fatalf("CreateEditingSession failed: %v %+v", err, edit)
	}

	session, err := client.Documents().CreateSession(context.Background(), "doc1", CreateDocumentSessionRequest{"recipient": "x"})
	if err != nil || session.ID != "sess1" {
		t.Fatalf("CreateSession failed: %v %+v", err, session)
	}

	appendResp, err := client.Documents().AppendContentLibraryItem(context.Background(), "doc1", AppendContentLibraryItemRequest{"id": "cli1"})
	if err != nil {
		t.Fatalf("AppendContentLibraryItem failed: %v", err)
	}
	if appendResp.BlockMapping["a"] != "b" {
		t.Fatalf("unexpected append response: %+v", appendResp)
	}
}

func TestDocumentsService_Downloads(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/download") || strings.HasSuffix(r.URL.Path, "/download-protected") {
			w.Header().Set("Content-Type", "application/pdf")
			_, _ = io.WriteString(w, "PDF")
			return
		}
		t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
	})

	d1, err := client.Documents().Download(context.Background(), "doc1")
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	b1, _ := io.ReadAll(d1.Body)
	_ = d1.Close()
	if string(b1) != "PDF" {
		t.Fatalf("unexpected download body")
	}

	d2, err := client.Documents().DownloadProtected(context.Background(), "doc1")
	if err != nil {
		t.Fatalf("DownloadProtected failed: %v", err)
	}
	b2, _ := io.ReadAll(d2.Body)
	_ = d2.Close()
	if string(b2) != "PDF" {
		t.Fatalf("unexpected protected download body")
	}
}

//nolint:gocognit // Test function that validates multiple error scenarios
func TestDocumentsService_ValidationErrors(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = w
		_ = r
		t.Fatalf("handler should not be called")
	})

	if _, err := client.Documents().Status(context.Background(), ""); err == nil {
		t.Fatalf("expected path param error")
	}
	if _, err := client.Documents().ESignDisclosure(context.Background(), ""); err == nil {
		t.Fatalf("expected path param error")
	}
	if err := client.Documents().Delete(context.Background(), ""); err == nil {
		t.Fatalf("expected path param error")
	}
	if _, err := client.Documents().RevertToDraft(context.Background(), ""); err == nil {
		t.Fatalf("expected path param error")
	}
	if _, err := client.Documents().Details(context.Background(), ""); err == nil {
		t.Fatalf("expected path param error")
	}
	if _, err := client.Documents().Download(context.Background(), ""); err == nil {
		t.Fatalf("expected path param error")
	}
	if _, err := client.Documents().DownloadProtected(context.Background(), ""); err == nil {
		t.Fatalf("expected path param error")
	}
	if err := client.Documents().Update(context.Background(), "doc", nil); err == nil {
		t.Fatalf("expected nil request error")
	}
	if err := client.Documents().ChangeStatus(context.Background(), "doc", nil); err == nil {
		t.Fatalf("expected nil request error")
	}
	if _, err := client.Documents().Create(context.Background(), nil); err == nil {
		t.Fatalf("expected nil request error")
	}
	if _, err := client.Documents().CreateFromUpload(context.Background(), &CreateDocumentFromUploadRequest{}); err == nil {
		t.Fatalf("expected nil file reader error")
	}
	if err := client.Documents().ChangeStatusWithUpload(context.Background(), "doc", &ChangeDocumentStatusWithUploadRequest{}); err == nil {
		t.Fatalf("expected nil file reader error")
	}
	if _, err := client.Documents().Send(context.Background(), "doc", nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if _, err := client.Documents().Send(context.Background(), "", DocumentSendRequest{}); err == nil {
		t.Fatalf("expected path param error")
	}
	if _, err := client.Documents().CreateEditingSession(context.Background(), "doc", nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if _, err := client.Documents().CreateEditingSession(context.Background(), "", CreateDocumentEditingSessionRequest{}); err == nil {
		t.Fatalf("expected path param error")
	}
	if _, err := client.Documents().CreateSession(context.Background(), "doc", nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if _, err := client.Documents().CreateSession(context.Background(), "", CreateDocumentSessionRequest{}); err == nil {
		t.Fatalf("expected path param error")
	}
	if err := client.Documents().TransferOwnership(context.Background(), "doc", nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if err := client.Documents().TransferOwnership(context.Background(), "", TransferDocumentOwnershipRequest{}); err == nil {
		t.Fatalf("expected path param error")
	}
	if err := client.Documents().TransferAllOwnership(context.Background(), nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if err := client.Documents().MoveToFolder(context.Background(), "", "folder"); err == nil {
		t.Fatalf("expected path param error")
	}
	if err := client.Documents().MoveToFolder(context.Background(), "doc", ""); err == nil {
		t.Fatalf("expected folder path param error")
	}
	if _, err := client.Documents().AppendContentLibraryItem(context.Background(), "doc", nil); err == nil {
		t.Fatalf("expected nil request")
	}
	if _, err := client.Documents().AppendContentLibraryItem(context.Background(), "", AppendContentLibraryItemRequest{}); err == nil {
		t.Fatalf("expected path param error")
	}
	if err := client.Documents().ChangeStatus(context.Background(), "", &ChangeDocumentStatusRequest{}); err == nil {
		t.Fatalf("expected path param error")
	}
	if err := client.Documents().ChangeStatusWithUpload(context.Background(), "", &ChangeDocumentStatusWithUploadRequest{File: strings.NewReader("x")}); err == nil {
		t.Fatalf("expected path param error")
	}
}

func assertQueryEq(t *testing.T, q map[string][]string, key, want string) {
	t.Helper()
	vals := q[key]
	got := ""
	if len(vals) > 0 {
		got = vals[0]
	}
	if got != want {
		t.Fatalf("query[%s]=%q want %q", key, got, want)
	}
}

func TestDocumentsService_CreateRequestBodyJSON(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["name"] != "doc" {
			t.Fatalf("unexpected payload: %+v", payload)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{"id":"doc1"}`)
	})

	if _, err := client.Documents().Create(context.Background(), DocumentCreateRequest{"name": "doc"}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
}
