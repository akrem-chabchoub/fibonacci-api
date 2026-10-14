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
		wantStatus int
		wantBody   string
		wantAllow  string
	}{
		{
			name:       "sequence",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=5",
			wantStatus: http.StatusOK,
			wantBody:   "[0,1,1,2,3]",
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
		},
		{
			name:       "negative",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=-1",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "too large",
			method:     http.MethodGet,
			target:     "/v1/fibonacci?n=95",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "method not allowed",
			method:     http.MethodPost,
			target:     "/v1/fibonacci?n=5",
			wantStatus: http.StatusMethodNotAllowed,
			wantAllow:  "GET",
		},
		{
			name:       "head not allowed",
			method:     http.MethodHead,
			target:     "/v1/fibonacci?n=5",
			wantStatus: http.StatusMethodNotAllowed,
			wantAllow:  "GET",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, tc.target, nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.wantStatus {
				t.Fatalf("status = %d, want %d", res.StatusCode, tc.wantStatus)
			}

			if tc.wantAllow != "" && res.Header.Get("Allow") != tc.wantAllow {
				t.Fatalf("Allow = %q, want %q", res.Header.Get("Allow"), tc.wantAllow)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			if tc.wantBody != "" && string(body) != tc.wantBody {
				t.Fatalf("body = %s, want %s", body, tc.wantBody)
			}

			if tc.name == "limits" {
				var got struct {
					MinN int `json:"min_n"`
					MaxN int `json:"max_n"`
				}
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}
				if got.MinN != 0 || got.MaxN != 94 {
					t.Fatalf("limits = %+v", got)
				}
			}
		})
	}
}
