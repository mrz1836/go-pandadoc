package pandadoc

import (
	"errors"
	"net/http"
	"testing"
)

var errTestDummy = errors.New("test error")

func TestAPIErrorHelpers(t *testing.T) {
	t.Parallel()

	err401 := &APIError{StatusCode: http.StatusUnauthorized, Message: "unauthorized"}
	err403 := &APIError{StatusCode: http.StatusForbidden, Message: "forbidden"}
	err404 := &APIError{StatusCode: http.StatusNotFound, Message: "not found"}
	err429 := &APIError{StatusCode: http.StatusTooManyRequests, Code: "too_many_requests", Message: "Request was throttled"}

	if !IsUnauthorized(err401) || IsUnauthorized(errTestDummy) {
		t.Fatalf("unauthorized helper mismatch")
	}
	if !IsForbidden(err403) {
		t.Fatalf("forbidden helper mismatch")
	}
	if !IsNotFound(err404) {
		t.Fatalf("not found helper mismatch")
	}
	if !IsRateLimited(err429) {
		t.Fatalf("rate limit helper mismatch")
	}
	if IsForbidden(errTestDummy) {
		t.Fatalf("expected false for non-API forbidden check")
	}
	if IsNotFound(errTestDummy) {
		t.Fatalf("expected false for non-API not-found check")
	}
	if IsRateLimited(errTestDummy) {
		t.Fatalf("expected false for non-API rate-limit check")
	}
	if err429.Error() == "" {
		t.Fatalf("expected non-empty error string")
	}
	plain := (&APIError{StatusCode: 400, Message: "bad"}).Error()
	if plain == "" {
		t.Fatalf("expected non-empty plain error")
	}

	var nilErr *APIError
	if nilErr.Error() != "" {
		t.Fatalf("expected nil receiver error string to be empty")
	}
}
