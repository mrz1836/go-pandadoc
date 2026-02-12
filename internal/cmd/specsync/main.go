// Package main is the specsync command that synchronizes the SDK against OpenAPI spec changes.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var errSpecContainsNoPaths = fmt.Errorf("spec contains no paths")

type openAPIDoc struct {
	Info struct {
		Version string `json:"version"`
	} `json:"info"`
	Paths map[string]map[string]json.RawMessage `json:"paths"`
}

type openAPIOperation struct {
	OperationID string   `json:"operationId"`
	Tags        []string `json:"tags"`
}

type operation struct {
	Method      string
	Path        string
	OperationID string
	Tag         string
}

var requiredOperations = []operation{ //nolint:gochecknoglobals // needed for test validation
	// Documents (20)
	{Method: "GET", Path: "/public/v1/documents"},
	{Method: "POST", Path: "/public/v1/documents"},
	{Method: "POST", Path: "/public/v1/documents?upload"},
	{Method: "GET", Path: "/public/v1/documents/{id}"},
	{Method: "DELETE", Path: "/public/v1/documents/{id}"},
	{Method: "PATCH", Path: "/public/v1/documents/{id}"},
	{Method: "GET", Path: "/public/v1/documents/{document_id}/esign-disclosure"},
	{Method: "PATCH", Path: "/public/v1/documents/{id}/status"},
	{Method: "PATCH", Path: "/public/v1/documents/{id}/status?upload"},
	{Method: "POST", Path: "/public/v1/documents/{id}/draft"},
	{Method: "GET", Path: "/public/v1/documents/{id}/details"},
	{Method: "POST", Path: "/public/v1/documents/{id}/send"},
	{Method: "POST", Path: "/public/v1/documents/{id}/editing-sessions"},
	{Method: "POST", Path: "/public/v1/documents/{id}/session"},
	{Method: "GET", Path: "/public/v1/documents/{id}/download"},
	{Method: "GET", Path: "/public/v1/documents/{id}/download-protected"},
	{Method: "PATCH", Path: "/public/v1/documents/{id}/ownership"},
	{Method: "PATCH", Path: "/public/v1/documents/ownership"},
	{Method: "POST", Path: "/public/v1/documents/{id}/move-to-folder/{folder_id}"},
	{Method: "POST", Path: "/public/v1/documents/{id}/append-content-library-item"},

	// Product catalog (5)
	{Method: "GET", Path: "/public/v2/product-catalog/items/search"},
	{Method: "POST", Path: "/public/v2/product-catalog/items"},
	{Method: "GET", Path: "/public/v2/product-catalog/items/{item_uuid}"},
	{Method: "PATCH", Path: "/public/v2/product-catalog/items/{item_uuid}"},
	{Method: "DELETE", Path: "/public/v2/product-catalog/items/{item_uuid}"},

	// OAuth (1)
	{Method: "POST", Path: "/oauth2/access_token"},

	// Webhook subscriptions/events (8)
	{Method: "GET", Path: "/public/v1/webhook-subscriptions"},
	{Method: "POST", Path: "/public/v1/webhook-subscriptions"},
	{Method: "GET", Path: "/public/v1/webhook-subscriptions/{id}"},
	{Method: "PATCH", Path: "/public/v1/webhook-subscriptions/{id}"},
	{Method: "DELETE", Path: "/public/v1/webhook-subscriptions/{id}"},
	{Method: "PATCH", Path: "/public/v1/webhook-subscriptions/{id}/shared-key"},
	{Method: "GET", Path: "/public/v1/webhook-events"},
	{Method: "GET", Path: "/public/v1/webhook-events/{id}"},
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		exitf("%v", err)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("specsync", flag.ContinueOnError)
	specPath := fs.String("spec", "internal/spec/openapi/pandadoc-openapi-7.18.0.json", "path to pinned OpenAPI spec")
	outPath := fs.String("out", "internal/spec/operations_gen.go", "path to generated operation manifest")
	if err := fs.Parse(args); err != nil {
		return err
	}

	doc, err := readSpec(*specPath)
	if err != nil {
		return fmt.Errorf("read spec: %w", err)
	}

	resolved, err := resolveRequiredOperations(doc)
	if err != nil {
		return fmt.Errorf("validate operations: %w", err)
	}

	if err := writeManifest(*outPath, doc.Info.Version, resolved); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

func readSpec(path string) (*openAPIDoc, error) {
	//nolint:gosec // Path is from command-line flag, not direct user input
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var doc openAPIDoc
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	if len(doc.Paths) == 0 {
		return nil, errSpecContainsNoPaths
	}

	return &doc, nil
}

func resolveRequiredOperations(doc *openAPIDoc) ([]operation, error) {
	resolved := make([]operation, 0, len(requiredOperations))

	for _, want := range requiredOperations {
		pathOps, ok := doc.Paths[want.Path]
		if !ok {
			return nil, fmt.Errorf("missing path in spec: %s: %w", want.Path, errSpecContainsNoPaths)
		}

		method := strings.ToLower(want.Method)
		rawOp, ok := pathOps[method]
		if !ok {
			return nil, fmt.Errorf("missing operation in spec: %s %s: %w", want.Method, want.Path, errSpecContainsNoPaths)
		}
		var op openAPIOperation
		if opErr := json.Unmarshal(rawOp, &op); opErr != nil {
			return nil, fmt.Errorf("decode operation %s %s: %w", want.Method, want.Path, opErr)
		}
		if op.OperationID == "" {
			return nil, fmt.Errorf("operation has no operationId: %s %s: %w", want.Method, want.Path, errSpecContainsNoPaths)
		}
		tag := ""
		if len(op.Tags) > 0 {
			tag = op.Tags[0]
		}

		resolved = append(resolved, operation{
			Method:      want.Method,
			Path:        want.Path,
			OperationID: op.OperationID,
			Tag:         tag,
		})
	}

	sort.Slice(resolved, func(i, j int) bool {
		if resolved[i].Tag != resolved[j].Tag {
			return resolved[i].Tag < resolved[j].Tag
		}
		if resolved[i].Path != resolved[j].Path {
			return resolved[i].Path < resolved[j].Path
		}
		return resolved[i].Method < resolved[j].Method
	})

	return resolved, nil
}

func writeManifest(path, version string, ops []operation) error {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by internal/cmd/specsync; DO NOT EDIT.\n")
	buf.WriteString("\n")
	buf.WriteString("package spec\n\n")
	buf.WriteString("// PinnedSpecVersion is the PandaDoc OpenAPI version this SDK is synced against.\n")
	buf.WriteString(fmt.Sprintf("const PinnedSpecVersion = %q\n\n", version))
	buf.WriteString("// CoveredOperations is the exact operation manifest supported by this SDK milestone.\n")
	buf.WriteString("var CoveredOperations = []Operation{\n")
	for _, op := range ops {
		buf.WriteString(fmt.Sprintf("\t{Method: %q, Path: %q, OperationID: %q, Tag: %q},\n", op.Method, op.Path, op.OperationID, op.Tag))
	}
	buf.WriteString("}\n")

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o600)
}

func exitf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
