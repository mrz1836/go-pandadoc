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
    <td align="center" width="33%">
       🚀&nbsp;<a href="#-installation"><code>Installation</code></a>
    </td>
    <td align="center" width="33%">
       🧪&nbsp;<a href="#-examples--tests"><code>Examples&nbsp;&&nbsp;Tests</code></a>
    </td>
    <td align="center" width="33%">
       📚&nbsp;<a href="#-documentation"><code>Documentation</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🤝&nbsp;<a href="#-contributing"><code>Contributing</code></a>
    </td>
    <td align="center">
      🛠️&nbsp;<a href="#-code-standards"><code>Code&nbsp;Standards</code></a>
    </td>
    <td align="center">
      ⚖️&nbsp;<a href="#-license"><code>License</code></a>
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

    "github.com/mrz1836/go-pandadoc"
    "github.com/mrz1836/go-pandadoc/commands"
)

func main() {
    // Create a new client
    client, err := pandadoc.NewClient("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    // List documents
    docs, err := client.Documents().List(context.Background(), &commands.ListDocumentsOptions{
        Count:  10,
        Status: "document.completed",
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

### Documents API

```go
// List documents with pagination and filters
docs, err := client.Documents().List(ctx, &commands.ListDocumentsOptions{
    Page:   1,
    Count:  25,
    Status: "document.completed",
})

// Get a document by ID
doc, err := client.Documents().Get(ctx, "document-id")

// Get document status
status, err := client.Documents().GetStatus(ctx, "document-id")

// Get document details (including fields)
details, err := client.Documents().GetDetails(ctx, "document-id")

// Update a document
updated, err := client.Documents().Update(ctx, "document-id", &commands.UpdateDocument{
    Name: "Updated Document Name",
})
```

### Product Catalog API

```go
// List catalog items
items, err := client.Catalog().List(ctx, &commands.ListCatalogOptions{
    Count: 50,
    Q:     "widget",
})

// Get a catalog item by ID
item, err := client.Catalog().Get(ctx, "item-id")
```

### Client Options

```go
import "time"

client, err := pandadoc.NewClient("your-api-key",
    pandadoc.WithTimeout(60 * time.Second),
    pandadoc.WithUserAgent("my-app/1.0"),
    pandadoc.WithBaseURL("https://api.pandadoc.com/public/v1/"),
)
```

<br/>

## 🧪 Examples & Tests

### Running Examples

```bash
# Set your API key
export PANDADOC_API_KEY="your-api-key"

# List documents
go run examples/list_documents/main.go

# Get document fields
go run examples/get_document_fields/main.go <document-id>

# List catalog items
go run examples/list_catalog/main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

<br/>

## 📚 Documentation

- [PandaDoc API Documentation](https://developers.pandadoc.com/reference/about)
- [Go Package Documentation](https://pkg.go.dev/github.com/mrz1836/go-pandadoc)

<br/>

## 🤝 Contributing

Contributions are welcome! Please read the [contributing guidelines](CONTRIBUTING.md) first.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

<br/>

## 🛠️ Code Standards

This project follows:

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)
- [Conventional Commits](https://www.conventionalcommits.org/)

<br/>

## ⚖️ License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

<br/>

## 👥 Maintainers

| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:------------------------------------------------------------------------------------------------:|
| [MrZ](https://github.com/mrz1836)                                                                |

<br/>

---

<div align="center">

**This is an unofficial SDK and is not affiliated with PandaDoc.**

</div>
