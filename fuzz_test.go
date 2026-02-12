package pandadoc

import (
	"encoding/json"
	"strings"
	"testing"
)

func FuzzJoinPaths(f *testing.F) {
	f.Add("/public/v1/", "documents")
	f.Add("/", "x")
	f.Add("public", "/v1")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, base, rel string) {
		got := joinPaths(base, rel)
		if got == "" {
			t.Fatalf("joinPaths returned empty path")
		}
		if !strings.HasPrefix(got, "/") {
			t.Fatalf("joinPaths must return absolute path: %q", got)
		}
	})
}

func FuzzPopulateAPIErrorFromBody(f *testing.F) {
	f.Add(`{"code":"bad","detail":"invalid"}`)
	f.Add(`{"type":"request_error","detail":{"x":"y"}}`)
	f.Add("plain")
	f.Add("")

	f.Fuzz(func(t *testing.T, body string) {
		e := &APIError{}
		populateAPIErrorFromBody(e, []byte(body))
		if e.Message == "" && strings.TrimSpace(body) != "" {
			var obj map[string]any
			if err := json.Unmarshal([]byte(body), &obj); err == nil {
				if _, hasKnown := obj["message"]; hasKnown {
					t.Fatalf("expected message for JSON with message field")
				}
			}
		}
	})
}
