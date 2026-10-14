package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeOpenAPI(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	ServeOpenAPI(rec, httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/yaml" {
		t.Fatalf("content-type = %q", rec.Header().Get("Content-Type"))
	}

	if !bytes.Contains(rec.Body.Bytes(), []byte("/v1/fibonacci")) {
		t.Fatal("spec missing /v1/fibonacci")
	}
}

func TestServeDocs(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	ServeDocs(rec, httptest.NewRequest(http.MethodGet, "/docs", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("content-type = %q", rec.Header().Get("Content-Type"))
	}

	if !bytes.Contains(rec.Body.Bytes(), []byte("swagger-ui")) {
		t.Fatal("docs missing swagger-ui")
	}
}
