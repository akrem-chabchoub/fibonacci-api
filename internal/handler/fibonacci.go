// Package handler maps HTTP requests to the service layer.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/akrem-chabchoub/fibonacci-api/api"
	"github.com/akrem-chabchoub/fibonacci-api/internal/service"
)

const contentTypeJSON = "application/json"

type Fibonacci struct{}

func (h Fibonacci) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, api.InvalidArgument, "method not allowed")

		return
	}

	if !r.URL.Query().Has("n") {
		writeJSON(w, http.StatusOK, api.Limits{MinN: 0})

		return
	}

	n, err := service.ParseN(r.URL.Query().Get("n"))
	if err != nil {
		writeServiceError(w, err)

		return
	}

	etag := strconv.Quote("n=" + strconv.Itoa(n))
	w.Header().Set("ETag", etag)

	if sameETag(r.Header.Get("If-None-Match"), etag) {
		w.WriteHeader(http.StatusNotModified)

		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_ = service.WriteJSON(r.Context(), w, n)
}

func sameETag(ifNoneMatch, etag string) bool {
	if ifNoneMatch == "" {
		return false
	}

	if ifNoneMatch == "*" {
		return true
	}

	for _, candidate := range strings.Split(ifNoneMatch, ",") {
		if strings.TrimSpace(candidate) == etag {
			return true
		}
	}

	return false
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNTooBig):
		writeError(w, http.StatusBadRequest, api.TooLarge, service.MsgNTooBig)
	default:
		writeError(w, http.StatusBadRequest, api.InvalidArgument, service.MsgBadN)
	}
}

func writeError(w http.ResponseWriter, status int, code api.ErrorCode, message string) {
	writeJSON(w, status, api.Error{Code: code, Message: message})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
