// Package handler maps HTTP requests to the service layer.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/akrem-chabchoub/fibonacci-api/internal/service"
)

const contentTypeJSON = "application/json"

type ErrorCode string

const (
	InvalidArgument  ErrorCode = "invalid_argument"
	TooLarge         ErrorCode = "too_large"
	MethodNotAllowed ErrorCode = "method_not_allowed"
	NotFound         ErrorCode = "not_found"
)

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type Limits struct {
	MinN int `json:"min_n"`
	MaxN int `json:"max_n"`
}

type Fibonacci struct{}

func (h Fibonacci) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		WriteError(w, http.StatusMethodNotAllowed, MethodNotAllowed, "method not allowed")

		return
	}

	if !r.URL.Query().Has("n") {
		writeJSON(w, http.StatusOK, Limits{MinN: 0, MaxN: service.MaxN})

		return
	}

	n, err := service.ParseN(r.URL.Query().Get("n"))
	if err != nil {
		writeServiceError(w, err)

		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_ = service.WriteJSON(r.Context(), w, n)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNTooBig):
		WriteError(w, http.StatusBadRequest, TooLarge, service.MsgNTooBig)
	default:
		WriteError(w, http.StatusBadRequest, InvalidArgument, service.MsgBadN)
	}
}

func WriteError(w http.ResponseWriter, status int, code ErrorCode, message string) {
	writeJSON(w, status, Error{Code: code, Message: message})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
