package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadSpecAndResolveRequiredOperations(t *testing.T) {
	t.Parallel()

	specPath := filepath.Join("..", "..", "spec", "openapi", "pandadoc-openapi-7.18.0.json")
	doc, err := readSpec(specPath)
	if err != nil {
		t.Fatalf("readSpec failed: %v", err)
	}

	ops, err := resolveRequiredOperations(doc)
	if err != nil {
		t.Fatalf("resolveRequiredOperations failed: %v", err)
	}
	if len(ops) != len(requiredOperations) {
		t.Fatalf("unexpected operation count: got %d want %d", len(ops), len(requiredOperations))
	}
	if doc.Info.Version == "" {
		t.Fatalf("expected spec version")
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	out := filepath.Join(dir, "ops.go")
	specPath := filepath.Join("..", "..", "spec", "openapi", "pandadoc-openapi-7.18.0.json")

	if err := run([]string{"-spec", specPath, "-out", out}); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected generated manifest: %v", err)
	}
}

func TestRun_BadFlagAndBadSpec(t *testing.T) {
	t.Parallel()

	if err := run([]string{"-badflag"}); err == nil {
		t.Fatalf("expected flag parse error")
	}
	if err := run([]string{"-spec", filepath.Join(t.TempDir(), "missing.json")}); err == nil {
		t.Fatalf("expected missing spec error")
	}
}

func TestMain_Success(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "ops.go")
	specPath := filepath.Join("..", "..", "spec", "openapi", "pandadoc-openapi-7.18.0.json")

	origArgs := os.Args
	t.Cleanup(func() { os.Args = origArgs })
	os.Args = []string{"specsync", "-spec", specPath, "-out", out}

	main()
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected generated file after main: %v", err)
	}
}

func TestExitf_Subprocess(t *testing.T) {
	if os.Getenv("SPECSYNC_TEST_EXITF") == "1" {
		exitf("boom")
		return
	}

	//nolint:gosec // This is a test that deliberately runs a subprocess
	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestExitf_Subprocess")
	cmd.Env = append(os.Environ(), "SPECSYNC_TEST_EXITF=1")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected subprocess to exit with failure")
	}
}

func TestMain_ErrorSubprocess(t *testing.T) {
	if os.Getenv("SPECSYNC_TEST_MAIN_ERROR") == "1" {
		orig := os.Args
		os.Args = []string{"specsync", "-badflag"}
		main()
		os.Args = orig
		return
	}

	//nolint:gosec // This is a test that deliberately runs a subprocess
	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestMain_ErrorSubprocess")
	cmd.Env = append(os.Environ(), "SPECSYNC_TEST_MAIN_ERROR=1")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected subprocess failure from main error path")
	}
}

func TestWriteManifest(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	out := filepath.Join(dir, "ops.go")
	ops := []operation{{Method: "GET", Path: "/x", OperationID: "op", Tag: "Tag"}}

	if err := writeManifest(out, "v1", ops); err != nil {
		t.Fatalf("writeManifest failed: %v", err)
	}

	//nolint:gosec // Path is constructed from safe test directory
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read generated file failed: %v", err)
	}
	text := string(data)
	if !strings.Contains(text, `const PinnedSpecVersion = "v1"`) {
		t.Fatalf("missing version in generated file: %s", text)
	}
	if !strings.Contains(text, `OperationID: "op"`) {
		t.Fatalf("missing operation in generated file: %s", text)
	}
}

func TestReadSpec_Errors(t *testing.T) {
	t.Parallel()

	if _, err := readSpec(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatalf("expected missing file error")
	}

	bad := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(bad, []byte("{"), 0o600); err != nil {
		t.Fatalf("write bad file: %v", err)
	}
	if _, err := readSpec(bad); err == nil {
		t.Fatalf("expected invalid json error")
	}
}

func TestReadSpec_NoPaths(t *testing.T) {
	t.Parallel()

	noPaths := filepath.Join(t.TempDir(), "no_paths.json")
	if err := os.WriteFile(noPaths, []byte(`{"info":{"version":"v1"},"paths":{}}`), 0o600); err != nil {
		t.Fatalf("write no-path spec: %v", err)
	}

	if _, err := readSpec(noPaths); err == nil || !strings.Contains(err.Error(), errSpecContainsNoPaths.Error()) {
		t.Fatalf("expected no paths error, got %v", err)
	}
}

func TestResolveRequiredOperations_Errors(t *testing.T) {
	t.Parallel()

	first := requiredOperations[0]
	method := strings.ToLower(first.Method)

	t.Run("missing path", func(t *testing.T) {
		t.Parallel()
		_, err := resolveRequiredOperations(&openAPIDoc{Paths: map[string]map[string]json.RawMessage{"/other": {}}})
		if err == nil || !strings.Contains(err.Error(), "missing path in spec") {
			t.Fatalf("expected missing path error, got %v", err)
		}
	})

	t.Run("missing method", func(t *testing.T) {
		t.Parallel()
		_, err := resolveRequiredOperations(&openAPIDoc{Paths: map[string]map[string]json.RawMessage{first.Path: {}}})
		if err == nil || !strings.Contains(err.Error(), "missing operation in spec") {
			t.Fatalf("expected missing operation error, got %v", err)
		}
	})

	t.Run("invalid operation json", func(t *testing.T) {
		t.Parallel()
		_, err := resolveRequiredOperations(&openAPIDoc{Paths: map[string]map[string]json.RawMessage{
			first.Path: {method: json.RawMessage("{")},
		}})
		if err == nil || !strings.Contains(err.Error(), "decode operation") {
			t.Fatalf("expected decode operation error, got %v", err)
		}
	})

	t.Run("missing operation id", func(t *testing.T) {
		t.Parallel()
		_, err := resolveRequiredOperations(&openAPIDoc{Paths: map[string]map[string]json.RawMessage{
			first.Path: {method: json.RawMessage(`{"tags":["x"]}`)},
		}})
		if err == nil || !strings.Contains(err.Error(), "operation has no operationId") {
			t.Fatalf("expected missing operationId error, got %v", err)
		}
	})
}

func TestRun_ValidateOperationsAndWriteManifestErrors(t *testing.T) {
	t.Parallel()

	specMissingRequired := filepath.Join(t.TempDir(), "missing_required.json")
	if err := os.WriteFile(specMissingRequired, []byte(`{"info":{"version":"v1"},"paths":{"/not-required":{"get":{"operationId":"x"}}}}`), 0o600); err != nil {
		t.Fatalf("write spec fixture: %v", err)
	}
	if err := run([]string{"-spec", specMissingRequired}); err == nil || !strings.Contains(err.Error(), "validate operations") {
		t.Fatalf("expected validate operations error, got %v", err)
	}

	dir := t.TempDir()
	blockingFile := filepath.Join(dir, "blocked-parent")
	if err := os.WriteFile(blockingFile, []byte("x"), 0o600); err != nil {
		t.Fatalf("write blocking file: %v", err)
	}
	specPath := filepath.Join("..", "..", "spec", "openapi", "pandadoc-openapi-7.18.0.json")
	outPath := filepath.Join(blockingFile, "ops.go")
	if err := run([]string{"-spec", specPath, "-out", outPath}); err == nil || !strings.Contains(err.Error(), "write manifest") {
		t.Fatalf("expected write manifest error, got %v", err)
	}
}
