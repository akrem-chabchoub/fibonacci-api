package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFibonacci(t *testing.T) {
	t.Parallel()

	h := Fibonacci{}

	tests := []struct {
		name       string
		method     string
		target     string
		headers    http.Header
		wantStatus int
		wantBody   string
		wantAllow  string
		wantETag   string
	}{
		{
			name:       "sequence",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=5",
			wantStatus: http.StatusOK,
			wantBody:   "[0,1,1,2,3]",
			wantETag:   `"n=5"`,
		},
		{
			name:       "not modified",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=5",
			headers:    http.Header{"If-None-Match": {`"n=5"`}},
			wantStatus: http.StatusNotModified,
			wantETag:   `"n=5"`,
		},
		{
			name:       "stale etag",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=5",
			headers:    http.Header{"If-None-Match": {`"n=4"`}},
			wantStatus: http.StatusOK,
			wantBody:   "[0,1,1,2,3]",
			wantETag:   `"n=5"`,
		},
		{
			name:       "limits",
			method:     http.MethodGet,
			target:     "/v1/fibonacci",
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing value",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":"invalid_argument","message":"n must be a non-negative integer"}`,
		},
		{
			name:       "negative",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=-1",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":"invalid_argument","message":"n must be a non-negative integer"}`,
		},
		{
			name:       "too large",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=1001",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":"too_large","message":"n exceeds maximum"}`,
		},
		{
			name:       "past uint64",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=95",
			wantStatus: http.StatusOK,
			wantETag:   `"n=95"`,
		},
		{
			name:       "method not allowed",
			method:     http.MethodPost,
			target:     "/v1/fibonacci?n=5",
			wantStatus: http.StatusMethodNotAllowed,
			wantAllow:  "GET",
			wantBody:   `{"code":"method_not_allowed","message":"method not allowed"}`,
		},
		{
			name:       "head not allowed",
			method:     http.MethodHead,
			target:     "/v1/fibonacci?n=5",
			wantStatus: http.StatusMethodNotAllowed,
			wantAllow:  "GET",
			wantBody:   `{"code":"method_not_allowed","message":"method not allowed"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, tc.target, nil)
			for key, values := range tc.headers {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			res := rec.Result()
			defer func() { _ = res.Body.Close() }()

			if res.StatusCode != tc.wantStatus {
				t.Fatalf("status = %d, want %d", res.StatusCode, tc.wantStatus)
			}

			if tc.wantAllow != "" && res.Header.Get("Allow") != tc.wantAllow {
				t.Fatalf("Allow = %q, want %q", res.Header.Get("Allow"), tc.wantAllow)
			}

			if tc.wantETag != "" && res.Header.Get("ETag") != tc.wantETag {
				t.Fatalf("ETag = %q, want %q", res.Header.Get("ETag"), tc.wantETag)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			if tc.wantBody != "" || tc.wantStatus == http.StatusNotModified {
				if string(body) != tc.wantBody {
					t.Fatalf("body = %s, want %s", body, tc.wantBody)
				}
			}

			if tc.name == "limits" {
				var got struct {
					MinN int `json:"min_n"`
					MaxN int `json:"max_n"`
				}
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}
				if got.MinN != 0 || got.MaxN != 1000 {
					t.Fatalf("limits = %+v", got)
				}
			}

			if tc.name == "past uint64" {
				var seq []json.Number
				if err := json.Unmarshal(body, &seq); err != nil {
					t.Fatal(err)
				}
				if len(seq) != 95 || seq[94].String() != "19740274219868223167" {
					t.Fatalf("past uint64 = len %d last %v", len(seq), seq)
				}
			}
		})
	}
}
