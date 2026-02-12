<div align="center">

# 📄&nbsp;&nbsp;go-pandadoc

**Unofficial Go SDK for the [PandaDoc API](https://developers.pandadoc.com/reference/about).**

<br/>

<a href="https://github.com/mrz1836/go-pandadoc/releases"><img src="https://img.shields.io/github/release-pre/mrz1836/go-pandadoc?include_prereleases&style=flat-square&logo=github&color=black" alt="Release"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/mrz1836/go-pandadoc?style=flat-square&logo=go&color=00ADD8" alt="Go Version"></a>
<a href="https://github.com/mrz1836/go-pandadoc/blob/master/LICENSE"><img src="https://img.shields.io/github/license/mrz1836/go-pandadoc?style=flat-square&color=blue&v=1" alt="License"></a>

<br/>

<table align="center" border="0">
  <tr>
    <td align="right">
       <code>CI / CD</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/mrz1836/go-pandadoc/actions"><img src="https://img.shields.io/github/actions/workflow/status/mrz1836/go-pandadoc/fortress.yml?branch=master&label=build&logo=github&style=flat-square" alt="Build"></a>
       <a href="https://github.com/mrz1836/go-pandadoc/actions"><img src="https://img.shields.io/github/last-commit/mrz1836/go-pandadoc?style=flat-square&logo=git&logoColor=white&label=last%20update" alt="Last Commit"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Quality</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://goreportcard.com/report/github.com/mrz1836/go-pandadoc"><img src="https://goreportcard.com/badge/github.com/mrz1836/go-pandadoc?style=flat-square&v=1" alt="Go Report"></a>
       <a href="https://codecov.io/gh/mrz1836/go-pandadoc"><img src="https://codecov.io/gh/mrz1836/go-pandadoc/branch/master/graph/badge.svg?style=flat-square" alt="Coverage"></a>
    </td>
  </tr>

  <tr>
    <td align="right">
       <code>Security</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://scorecard.dev/viewer/?uri=github.com/mrz1836/go-pandadoc"><img src="https://api.scorecard.dev/projects/github.com/mrz1836/go-pandadoc/badge?style=flat-square" alt="Scorecard"></a>
       <a href=".github/SECURITY.md"><img src="https://img.shields.io/badge/policy-active-success?style=flat-square&logo=security&logoColor=white" alt="Security"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Community</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/mrz1836/go-pandadoc/graphs/contributors"><img src="https://img.shields.io/github/contributors/mrz1836/go-pandadoc?style=flat-square&color=orange" alt="Contributors"></a>
       <a href="https://mrz1818.com/"><img src="https://img.shields.io/badge/donate-bitcoin-ff9900?style=flat-square&logo=bitcoin" alt="Bitcoin"></a>
    </td>
  </tr>
</table>

</div>

<br/>
<br/>

<div align="center">

### <code>Project Navigation</code>

</div>

<table align="center">
  <tr>
    <td align="center" width="33%">
       📦&nbsp;<a href="#-installation"><code>Installation</code></a>
    </td>
    <td align="center" width="33%">
       🚀&nbsp;<a href="#-quick-start"><code>Quick&nbsp;Start</code></a>
    </td>
    <td align="center" width="33%">
       📚&nbsp;<a href="#-features"><code>Features</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       📊&nbsp;<a href="#-api-coverage"><code>API&nbsp;Coverage</code></a>
    </td>
    <td align="center">
       🧪&nbsp;<a href="#-examples--tests"><code>Examples&nbsp;&&nbsp;Tests</code></a>
    </td>
    <td align="center">
       📚&nbsp;<a href="#-documentation"><code>Documentation</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
      🛠️&nbsp;<a href="#%EF%B8%8F-code-standards"><code>Code&nbsp;Standards</code></a>
    </td>
    <td align="center">
      🤖&nbsp;<a href="#-ai-usage--assistant-guidelines"><code>AI&nbsp;Guidelines</code></a>
    </td>
    <td align="center">
       👥&nbsp;<a href="#-maintainers"><code>Maintainers</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🤝&nbsp;<a href="#-contributing"><code>Contributing</code></a>
    </td>
    <td align="center">
       📝&nbsp;<a href="#-license"><code>License</code></a>
    </td>
    <td align="center"></td>
  </tr>
</table>

<br/>

## 📦 Installation

```bash
go get github.com/mrz1836/go-pandadoc
```

<br/>

## 🚀 Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/mrz1836/go-pandadoc"
)

func main() {
    apiKey := os.Getenv("PANDADOC_API_KEY")
    if apiKey == "" {
        log.Fatal("PANDADOC_API_KEY is required")
    }

    client, err := pandadoc.NewClientWithAPIKey(apiKey)
    if err != nil {
        log.Fatal(err)
    }

    status := pandadoc.DocumentStatusCompleted
    docs, err := client.Documents().List(context.Background(), &pandadoc.ListDocumentsOptions{
        Count:  10,
        Status: &status,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, doc := range docs.Results {
        fmt.Printf("Document: %s (%s)\n", doc.Name, doc.Status)
    }
}
```

<br/>

## 📚 Features

### Auth Modes

```go
// API key auth
client, err := pandadoc.NewClientWithAPIKey("api-key")
_ = client

// OAuth bearer auth
client, err = pandadoc.NewClientWithAccessToken("access-token")
_ = client
_ = err
```

### Client Options

```go
import "time"

client, err := pandadoc.NewClientWithAPIKey("api-key",
    pandadoc.WithTimeout(60 * time.Second),
    pandadoc.WithUserAgent("my-app/1.0"),
    pandadoc.WithBaseURL("https://api.pandadoc.com/"),
    pandadoc.WithRetryPolicy(pandadoc.DefaultRetryPolicy()),
)
```

### Documents Service

```go
// List documents
docs, err := client.Documents().List(ctx, &pandadoc.ListDocumentsOptions{Count: 25})
_ = docs

// Create a document from JSON payload
created, err := client.Documents().Create(ctx, pandadoc.DocumentCreateRequest{
    "name":          "Proposal",
    "template_uuid": "template-id",
    "recipients": []map[string]any{
        {"email": "jane@example.com", "first_name": "Jane", "last_name": "Doe", "role": "Signer"},
    },
})
_ = created

// Get status/details
status, err := client.Documents().Status(ctx, "document-id")
details, err := client.Documents().Details(ctx, "document-id")
_ = status
_ = details

// Update (204 no content)
err = client.Documents().Update(ctx, "document-id", pandadoc.DocumentUpdateRequest{
    "name": "Updated Name",
})
_ = err
```

### Product Catalog Service

```go
// Search catalog
items, err := client.ProductCatalog().Search(ctx, &pandadoc.SearchProductCatalogItemsOptions{
    Query:   "coffee",
    PerPage: 20,
})
_ = items

// Create/update/get/delete
createdItem, err := client.ProductCatalog().Create(ctx, pandadoc.CreateProductCatalogItemRequest{
    "type":  "regular",
    "title": "New Product",
})
item, err := client.ProductCatalog().Get(ctx, createdItem.UUID)
updated, err := client.ProductCatalog().Update(ctx, createdItem.UUID, pandadoc.UpdateProductCatalogItemRequest{
    "title": "Updated Product",
})
err = client.ProductCatalog().Delete(ctx, createdItem.UUID)
_ = item
_ = updated
_ = err
```

### OAuth Token Exchange

```go
oauthClient, err := pandadoc.NewClient()
token, err := oauthClient.OAuth().Token(ctx, &pandadoc.OAuthTokenRequest{
    GrantType:    "authorization_code",
    ClientID:     "client-id",
    ClientSecret: "client-secret",
    Code:         "authorization-code",
    Scope:        "read+write",
})
_ = token
_ = err
```

### Webhooks

```go
// Subscriptions
subs, err := client.WebhookSubscriptions().List(ctx, &pandadoc.ListWebhookSubscriptionsOptions{Count: 50, Page: 1})
createdSub, err := client.WebhookSubscriptions().Create(ctx, &pandadoc.WebhookSubscriptionRequest{
    Name: "my-subscription",
    URL:  "https://example.com/pandadoc/webhooks",
    Triggers: []pandadoc.WebhookTrigger{
        pandadoc.WebhookTriggerDocumentUpdated,
        pandadoc.WebhookTriggerDocumentStateChanged,
    },
})
sub, err := client.WebhookSubscriptions().Get(ctx, createdSub.UUID)
updatedSub, err := client.WebhookSubscriptions().Update(ctx, createdSub.UUID, &pandadoc.WebhookSubscriptionRequest{Name: "updated-name"})
key, err := client.WebhookSubscriptions().RegenerateSharedKey(ctx, createdSub.UUID)
err = client.WebhookSubscriptions().Delete(ctx, createdSub.UUID)
_ = subs
_ = sub
_ = updatedSub
_ = key

// Events
events, err := client.WebhookEvents().List(ctx, &pandadoc.ListWebhookEventsOptions{
    Type: "document_updated",
})
if len(events.Items) > 0 {
    event, err := client.WebhookEvents().Get(ctx, events.Items[0].UUID)
    _ = event
    _ = err
}
```

### Unit Testing & Mocking

The SDK now defines interfaces for all service interactions, making it easy to mock the client in your tests.

```go
// Create a mock struct that implements pandadoc.DocumentsService
type mockDocuments struct {
    pandadoc.DocumentsService // Embed interface to skip implementing all methods if not needed
}

func (m *mockDocuments) List(ctx context.Context, opts *pandadoc.ListDocumentsOptions) (*pandadoc.DocumentListResponse, error) {
    return &pandadoc.DocumentListResponse{
        Results: []pandadoc.DocumentSummary{
            {ID: "mock-doc-1", Name: "Mock Document"},
        },
    }, nil
}
```

### Observability

You can inject a custom logger to monitor SDK operations. The logger must implement the `pandadoc.Logger` interface.

```go
type myLogger struct{}

func (l *myLogger) Debugf(format string, args ...interface{}) { log.Printf("[DEBUG] "+format, args...) }
func (l *myLogger) Infof(format string, args ...interface{})  { log.Printf("[INFO] "+format, args...) }
func (l *myLogger) Errorf(format string, args ...interface{}) { log.Printf("[ERROR] "+format, args...) }

func main() {
    client, _ := pandadoc.NewClientWithAPIKey("key",
        pandadoc.WithLogger(&myLogger{}),
    )
}
```

<br/>

## 📊 API Coverage

This SDK implements core PandaDoc API functionality. Below is a comprehensive comparison of all available endpoints vs. what's currently supported.

<details>
<summary><strong>View Complete API Endpoint Coverage (115 endpoints)</strong></summary>

<br/>

### Coverage Summary
- ✅ **Implemented:** 5 services, 34 endpoints (~30% coverage)
- 📝 **Available in API:** 26 services, 115 endpoints
- 🎯 **Focus Areas:** Documents, Product Catalog, Webhooks, OAuth

---

### Legend
- ✅ = Fully implemented in SDK
- ⚠️ = Partially implemented
- ❌ = Not yet implemented
- 📄 [Docs] = Link to PandaDoc API reference

---

### 1. Documents ⚠️
*Core document lifecycle management - 20 of 22 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ✅ | GET | `/public/v1/documents` | `Documents().List()` | [📄](https://developers.pandadoc.com/reference/list-documents) |
| ✅ | POST | `/public/v1/documents` | `Documents().Create()` | [📄](https://developers.pandadoc.com/reference/create-document) |
| ✅ | POST | `/public/v1/documents?upload` | `Documents().CreateFromUpload()` | [📄](https://developers.pandadoc.com/reference/create-document-from-upload) |
| ✅ | GET | `/public/v1/documents/{id}` | `Documents().Status()` | [📄](https://developers.pandadoc.com/reference/status-document) |
| ✅ | PATCH | `/public/v1/documents/{id}` | `Documents().Update()` | [📄](https://developers.pandadoc.com/reference/update-document) |
| ✅ | DELETE | `/public/v1/documents/{id}` | `Documents().Delete()` | [📄](https://developers.pandadoc.com/reference/delete-document) |
| ✅ | GET | `/public/v1/documents/{id}/details` | `Documents().Details()` | [📄](https://developers.pandadoc.com/reference/details-document) |
| ✅ | GET | `/public/v1/documents/{id}/download` | `Documents().Download()` | [📄](https://developers.pandadoc.com/reference/download-document) |
| ✅ | GET | `/public/v1/documents/{id}/download-protected` | `Documents().DownloadProtected()` | [📄](https://developers.pandadoc.com/reference/download-protected-document) |
| ✅ | POST | `/public/v1/documents/{id}/send` | `Documents().Send()` | [📄](https://developers.pandadoc.com/reference/send-document) |
| ✅ | POST | `/public/v1/documents/{id}/session` | `Documents().CreateSession()` | [📄](https://developers.pandadoc.com/reference/create-document-link) |
| ✅ | POST | `/public/v1/documents/{id}/editing-sessions` | `Documents().CreateEditingSession()` | [📄](https://developers.pandadoc.com/reference/create-document-editing-session) |
| ✅ | PATCH | `/public/v1/documents/{id}/status` | `Documents().ChangeStatus()` | [📄](https://developers.pandadoc.com/reference/change-document-status) |
| ✅ | PATCH | `/public/v1/documents/{id}/status?upload` | `Documents().ChangeStatusWithUpload()` | [📄](https://developers.pandadoc.com/reference/change-document-status) |
| ✅ | POST | `/public/v1/documents/{id}/draft` | `Documents().RevertToDraft()` | [📄](https://developers.pandadoc.com/reference/document-revert-to-draft) |
| ✅ | GET | `/public/v1/documents/{document_id}/esign-disclosure` | `Documents().ESignDisclosure()` | [📄](https://developers.pandadoc.com/reference/document-esign-disclosure) |
| ✅ | POST | `/public/v1/documents/{id}/move-to-folder/{folder_id}` | `Documents().MoveToFolder()` | [📄](https://developers.pandadoc.com/reference/document-move-to-folder) |
| ✅ | PATCH | `/public/v1/documents/{id}/ownership` | `Documents().TransferOwnership()` | [📄](https://developers.pandadoc.com/reference/transfer-document-ownership) |
| ✅ | PATCH | `/public/v1/documents/ownership` | `Documents().TransferAllOwnership()` | [📄](https://developers.pandadoc.com/reference/transfer-all-documents-ownership) |
| ✅ | POST | `/public/v1/documents/{id}/append-content-library-item` | `Documents().AppendContentLibraryItem()` | [📄](https://developers.pandadoc.com/reference/append-content-library-item-to-document) |
| ❌ | POST | `/public/beta/documents/{document_id}/docx-export-tasks` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/createexportdocxtask) |
| ❌ | GET | `/public/beta/documents/{document_id}/docx-export-tasks/{task_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/getdocxexporttask) |

---

### 2. Product Catalog ✅
*Manage product catalog items - 5 of 5 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ✅ | GET | `/public/v2/product-catalog/items/search` | `ProductCatalog().Search()` | [📄](https://developers.pandadoc.com/reference/search-catalog-items) |
| ✅ | POST | `/public/v2/product-catalog/items` | `ProductCatalog().Create()` | [📄](https://developers.pandadoc.com/reference/create-catalog-item) |
| ✅ | GET | `/public/v2/product-catalog/items/{item_uuid}` | `ProductCatalog().Get()` | [📄](https://developers.pandadoc.com/reference/get-catalog-item) |
| ✅ | PATCH | `/public/v2/product-catalog/items/{item_uuid}` | `ProductCatalog().Update()` | [📄](https://developers.pandadoc.com/reference/update-catalog-item) |
| ✅ | DELETE | `/public/v2/product-catalog/items/{item_uuid}` | `ProductCatalog().Delete()` | [📄](https://developers.pandadoc.com/reference/delete-catalog-item) |

---

### 3. Webhook Subscriptions ✅
*Manage webhook endpoints and subscriptions - 6 of 6 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ✅ | GET | `/public/v1/webhook-subscriptions` | `WebhookSubscriptions().List()` | [📄](https://developers.pandadoc.com/reference/list-webhook-subscriptions) |
| ✅ | POST | `/public/v1/webhook-subscriptions` | `WebhookSubscriptions().Create()` | [📄](https://developers.pandadoc.com/reference/create-webhook-subscription) |
| ✅ | GET | `/public/v1/webhook-subscriptions/{id}` | `WebhookSubscriptions().Get()` | [📄](https://developers.pandadoc.com/reference/details-webhook-subscription) |
| ✅ | PATCH | `/public/v1/webhook-subscriptions/{id}` | `WebhookSubscriptions().Update()` | [📄](https://developers.pandadoc.com/reference/update-webhook-subscription) |
| ✅ | DELETE | `/public/v1/webhook-subscriptions/{id}` | `WebhookSubscriptions().Delete()` | [📄](https://developers.pandadoc.com/reference/delete-webhook-subscription) |
| ✅ | PATCH | `/public/v1/webhook-subscriptions/{id}/shared-key` | `WebhookSubscriptions().RegenerateSharedKey()` | [📄](https://developers.pandadoc.com/reference/update-webhook-subscription-shared-key) |

---

### 4. Webhook Events ✅
*Retrieve webhook event history - 2 of 2 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ✅ | GET | `/public/v1/webhook-events` | `WebhookEvents().List()` | [📄](https://developers.pandadoc.com/reference/list-webhook-event) |
| ✅ | GET | `/public/v1/webhook-events/{id}` | `WebhookEvents().Get()` | [📄](https://developers.pandadoc.com/reference/details-webhook-event) |

---

### 5. OAuth 2.0 Authentication ✅
*Handle OAuth token creation and refresh - 1 of 1 endpoint implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ✅ | POST | `/oauth2/access_token` | `OAuth().Token()` | [📄](https://developers.pandadoc.com/reference/access-token) |

---

### 6. Templates ❌
*Template management and operations - 0 of 8 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/templates` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-templates) |
| ❌ | POST | `/public/v1/templates` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-template) |
| ❌ | POST | `/public/v1/templates?upload` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-template-from-file) |
| ❌ | GET | `/public/v1/templates/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/template-status) |
| ❌ | GET | `/public/v1/templates/{id}/details` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/template-details) |
| ❌ | PATCH | `/public/v1/templates/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/template-update) |
| ❌ | DELETE | `/public/v1/templates/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/delete-template) |
| ❌ | POST | `/public/v1/templates/{id}/editing-sessions` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-template-editing-session) |

---

### 7. Document Recipients ❌
*Manage document recipients and signers - 0 of 4 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | POST | `/public/v1/documents/{id}/recipients` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/add-cc-recipient) |
| ❌ | PATCH | `/public/v1/documents/{id}/recipients/recipient/{recipient_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/edit-document-recipient) |
| ❌ | DELETE | `/public/v1/documents/{id}/recipients/{recipient_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/delete-document-recipient) |
| ❌ | POST | `/public/v1/documents/{id}/recipients/{recipient_id}/reassign` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/change-signer) |

---

### 8. Document Fields ❌
*Manage fillable fields in documents - 0 of 2 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/documents/{id}/fields` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-document-fields) |
| ❌ | POST | `/public/v1/documents/{id}/fields` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-document-fields) |

---

### 9. Document Attachments ❌
*Manage document attachments and uploads - 0 of 6 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/documents/{id}/attachments` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-attachment) |
| ❌ | POST | `/public/v1/documents/{id}/attachments` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-document-attachment) |
| ❌ | POST | `/public/v1/documents/{id}/attachments?upload` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-document-attachment-from-file-upload) |
| ❌ | GET | `/public/v1/documents/{id}/attachments/{attachment_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/attachment-details) |
| ❌ | DELETE | `/public/v1/documents/{id}/attachments/{attachment_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/delete-attachment) |
| ❌ | GET | `/public/v1/documents/{id}/attachments/{attachment_id}/download` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/download-attachment) |

---

### 10. Document Reminders ❌
*Auto-reminders and manual reminders for document recipients - 0 of 4 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/documents/{document_id}/auto-reminders` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/getdocumentautoremindersettings) |
| ❌ | PATCH | `/public/v1/documents/{document_id}/auto-reminders` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/updatedocumentautoremindersettings) |
| ❌ | GET | `/public/v1/documents/{document_id}/auto-reminders/status` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/statusdocumentautoreminder) |
| ❌ | POST | `/public/v1/documents/{document_id}/send-reminder` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/createmanualreminder) |

---

### 11. Document Sections (Bundles) ❌
*Manage document sections/bundles - 0 of 6 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/documents/{document_id}/sections` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-sections) |
| ❌ | POST | `/public/v1/documents/{document_id}/sections/uploads` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-document-section) |
| ❌ | POST | `/public/v1/documents/{document_id}/sections/uploads?upload` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-document-section-from-upload) |
| ❌ | GET | `/public/v1/documents/{document_id}/sections/uploads/{upload_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/document-section-upload-status) |
| ❌ | GET | `/public/v1/documents/{document_id}/sections/{section_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-document-section) |
| ❌ | DELETE | `/public/v1/documents/{document_id}/sections/{section_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/delete-section) |

---

### 12. Document Link to CRM ❌
*Link documents to CRM systems and objects - 0 of 4 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/documents/linked-objects` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-documents-by-linked-object) |
| ❌ | GET | `/public/v1/documents/{id}/linked-objects` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-linked-objects) |
| ❌ | POST | `/public/v1/documents/{id}/linked-objects` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/link-to-crm) |
| ❌ | DELETE | `/public/v1/documents/{id}/linked-objects/{linked_object_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/delete-linked-object) |

---

### 13. Contacts ❌
*Manage contact information - 0 of 5 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/contacts` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-contacts) |
| ❌ | POST | `/public/v1/contacts` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-contact) |
| ❌ | GET | `/public/v1/contacts/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/details-contact) |
| ❌ | PATCH | `/public/v1/contacts/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/update-contact) |
| ❌ | DELETE | `/public/v1/contacts/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/delete-contact) |

---

### 14. Content Library Items ❌
*Manage reusable content library items - 0 of 5 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/content-library-items` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-content-library-items) |
| ❌ | POST | `/public/v1/content-library-items` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-content-library-item) |
| ❌ | POST | `/public/v1/content-library-items?upload` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-content-library-item-from-file) |
| ❌ | GET | `/public/v1/content-library-items/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/content-library-item-details) |
| ❌ | GET | `/public/v1/content-library-items/{id}/details` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/content-library-item-details) |

---

### 15. Forms ❌
*Retrieve forms information - 0 of 1 endpoint implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/forms` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-forms) |

---

### 16. Folders ❌
*Organize documents and templates into folders - 0 of 6 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/documents/folders` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-documents-folders) |
| ❌ | POST | `/public/v1/documents/folders` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-document-folder) |
| ❌ | PUT | `/public/v1/documents/folders/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/rename-document-folder) |
| ❌ | GET | `/public/v1/templates/folders` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-templates-folders) |
| ❌ | POST | `/public/v1/templates/folders` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-templates-folder) |
| ❌ | PUT | `/public/v1/templates/folders/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/rename-template-folder) |

---

### 17. Quotes ❌
*Manage quotes within documents - 0 of 1 endpoint implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | PUT | `/public/v1/documents/{document_id}/quotes/{quote_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/update-quote) |

---

### 18. Members ❌
*Manage workspace members and users - 0 of 4 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/members` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-members) |
| ❌ | GET | `/public/v1/members/current` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/member-details) |
| ❌ | GET | `/public/v1/members/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/member-details) |
| ❌ | POST | `/public/v1/members/{member_id}/token` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-member-token) |

---

### 19. User and Workspace Management ❌
*Manage users, workspaces, and organizational settings - 0 of 8 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/users` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-users) |
| ❌ | POST | `/public/v1/users` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-user) |
| ❌ | GET | `/public/v1/workspaces` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/get-workspaces-list) |
| ❌ | POST | `/public/v1/workspaces` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-workspace) |
| ❌ | POST | `/public/v1/workspaces/{workspace_id}/api-keys` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-api-key) |
| ❌ | POST | `/public/v1/workspaces/{workspace_id}/deactivate` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/deactivate-workspace) |
| ❌ | POST | `/public/v1/workspaces/{workspace_id}/members` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/add-member-to-workspace) |
| ❌ | DELETE | `/public/v1/workspaces/{workspace_id}/members/{member_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/remove-member) |

---

### 20. API Logs ❌
*Retrieve API activity logs - 0 of 4 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/logs` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-api-logs) |
| ❌ | GET | `/public/v1/logs/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/api-log-details) |
| ❌ | GET | `/public/v2/logs` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-api-logs) |
| ❌ | GET | `/public/v2/logs/{id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/api-log-details) |

---

### 21. Document Settings (v2) ❌
*Manage document-specific settings - 0 of 2 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v2/documents/{document_id}/settings` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/document-settings-get) |
| ❌ | PATCH | `/public/v2/documents/{document_id}/settings` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/document-settings-update) |

---

### 22. Template Settings (v2) ❌
*Manage template-specific settings - 0 of 2 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v2/templates/{template_id}/settings` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/template-settings-get) |
| ❌ | PATCH | `/public/v2/templates/{template_id}/settings` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/template-settings-update) |

---

### 23. Document Audit Trail (v2) ❌
*Retrieve document audit logs - 0 of 1 endpoint implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v2/documents/{document_id}/audit-trail` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-document-audit-trail) |

---

### 24. Document Structure View (v2) ❌
*Add named items to document structure - 0 of 1 endpoint implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | POST | `/public/v2/dsv/{document_id}/add-named-items` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/add-dsv-named-items) |

---

### 25. Notary (v2) ❌
*Manage notarization requests and services - 0 of 4 endpoints implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v2/notary/notaries` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/list-notaries) |
| ❌ | POST | `/public/v2/notary/notarization-requests` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/create-notarization-request) |
| ❌ | GET | `/public/v2/notary/notarization-requests/{session_request_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/notarization-request-details) |
| ❌ | DELETE | `/public/v2/notary/notarization-requests/{session_request_id}` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/delete-notarization-request) |

---

### 26. Communication Preferences ❌
*Manage SMS opt-outs - 0 of 1 endpoint implemented*

| Status | Method | Endpoint | SDK Method | API Docs |
|--------|--------|----------|------------|----------|
| ❌ | GET | `/public/v1/sms-opt-outs` | *Not implemented* | [📄](https://developers.pandadoc.com/reference/listrecentsmsoptouts) |

---

</details>

<br/>

## 🧪 Examples & Tests

All unit tests and fuzz tests run via [GitHub Actions](https://github.com/mrz1836/go-pandadoc/actions) and use [Go version 1.24.x](https://go.dev/doc/go1.24). View the [configuration file](.github/workflows/fortress.yml).

Run all tests (fast):

```bash script
magex test
```

Run all tests with race detector (slower):
```bash script
magex test:race
```

<br/>

## 📚 Documentation

- [PandaDoc API Documentation](https://developers.pandadoc.com/reference/about)
- [Go Package Documentation](https://pkg.go.dev/github.com/mrz1836/go-pandadoc)

<br/>

## 🛠️ Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## 🤖 AI Usage & Assistant Guidelines
Read the [AI Usage & Assistant Guidelines](.github/tech-conventions/ai-compliance.md) for details on how AI is used in this project and how to interact with the AI assistants.

<br/>

## 👥 Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:------------------------------------------------------------------------------------------------:|
|                                [MrZ](https://github.com/mrz1836)                                 |

<br/>

## 🤝 Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap:
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-pandadoc&utm_term=go-pandadoc&utm_content=go-pandadoc) to ensure this journey continues indefinitely! :rocket:


[![Stars](https://img.shields.io/github/stars/mrz1836/go-pandadoc?label=Please%20like%20us&style=social)](https://github.com/mrz1836/go-pandadoc/stargazers)

<br/>

## 📝 License

[![License](https://img.shields.io/github/license/mrz1836/go-pandadoc.svg?style=flat)](LICENSE)
