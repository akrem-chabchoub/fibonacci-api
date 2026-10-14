// Package handler maps HTTP requests to the service layer.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/akrem-chabchoub/fibonacci-api/internal/service"
)

const contentTypeJSON = "application/json"

// Fibonacci serves GET /v1/fibonacci.
type Fibonacci struct{}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type limits struct {
	MinN int `json:"min_n"`
	MaxN int `json:"max_n"`
}

// ServeHTTP handles the Fibonacci HTTP API.
func (h Fibonacci) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, service.CodeBadN, "method not allowed")

		return
	}

	if !r.URL.Query().Has("n") {
		writeJSON(w, http.StatusOK, limits{MinN: 0, MaxN: service.MaxN})

		return
	}

	n, err := service.ParseN(r.URL.Query().Get("n"))
	if err != nil {
		writeServiceError(w, err)

		return
	}

	body, err := service.JSON(n)
	if err != nil {
		writeServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.Header().Set("ETag", strconv.Quote("n="+strconv.Itoa(n)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNTooBig):
		writeError(w, http.StatusBadRequest, service.CodeNTooBig, service.MsgNTooBig)
	default:
		writeError(w, http.StatusBadRequest, service.CodeBadN, service.MsgBadN)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, apiError{Code: code, Message: message})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
