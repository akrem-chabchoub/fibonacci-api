package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	WriteError(rec, http.StatusBadRequest, InvalidArgument, "n must be a non-negative integer")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != contentTypeJSON {
		t.Fatalf("content-type = %q", rec.Header().Get("Content-Type"))
	}

	want := `{"code":"invalid_argument","message":"n must be a non-negative integer"}`
	if rec.Body.String() != want {
		t.Fatalf("body = %s, want %s", rec.Body.String(), want)
	}
}
