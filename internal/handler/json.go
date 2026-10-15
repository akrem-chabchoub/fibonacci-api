package handler

import (
	"encoding/json"
	"net/http"
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

func WriteError(w http.ResponseWriter, status int, code ErrorCode, message string) {
	WriteJSON(w, status, Error{Code: code, Message: message})
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
