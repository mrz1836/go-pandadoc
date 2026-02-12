<div align="center">

# 📄&nbsp;&nbsp;go-pandadoc

**Unofficial Go SDK for the PandaDoc API.**

<br/>

<a href="https://github.com/mrz1836/go-pandadoc/releases"><img src="https://img.shields.io/github/release-pre/mrz1836/go-pandadoc?include_prereleases&style=flat-square&logo=github&color=black" alt="Release"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/mrz1836/go-pandadoc?style=flat-square&logo=go&color=00ADD8" alt="Go Version"></a>
<a href="https://github.com/mrz1836/go-pandadoc/blob/master/LICENSE"><img src="https://img.shields.io/github/license/mrz1836/go-pandadoc?style=flat-square&color=blue" alt="License"></a>

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
       <a href="https://goreportcard.com/report/github.com/mrz1836/go-pandadoc"><img src="https://goreportcard.com/badge/github.com/mrz1836/go-pandadoc?style=flat-square" alt="Go Report"></a>
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
    <td align="center" width="20%">
       📦&nbsp;<a href="#-installation"><code>Installation</code></a>
    </td>
    <td align="center" width="20%">
       🚀&nbsp;<a href="#-quick-start"><code>Quick&nbsp;Start</code></a>
    </td>
    <td align="center" width="20%">
       📚&nbsp;<a href="#-features"><code>Features</code></a>
    </td>
    <td align="center" width="20%">
       🧪&nbsp;<a href="#-examples--tests"><code>Examples&nbsp;&&nbsp;Tests</code></a>
    </td>
    <td align="center" width="20%">
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
    <td align="center">
       🤝&nbsp;<a href="#-contributing"><code>Contributing</code></a>
    </td>
    <td align="center">
      📝&nbsp;<a href="#-license"><code>License</code></a>
    </td>
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
