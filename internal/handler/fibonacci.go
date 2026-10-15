// Package handler maps HTTP requests to the service layer.
package handler

import (
	"errors"
	"net/http"

	"github.com/akrem-chabchoub/fibonacci-api/internal/service"
)

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
		WriteJSON(w, http.StatusOK, Limits{MinN: 0, MaxN: service.MaxN})

		return
	}

	n, err := service.ParseN(r.URL.Query().Get("n"))
	if err != nil {
		if errors.Is(err, service.ErrNTooBig) {
			WriteError(w, http.StatusBadRequest, TooLarge, service.MsgNTooBig)
		} else {
			WriteError(w, http.StatusBadRequest, InvalidArgument, service.MsgBadN)
		}

		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_ = service.WriteJSON(r.Context(), w, n)
}
