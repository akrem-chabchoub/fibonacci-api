package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/akrem-chabchoub/fibonacci-api/api"
	"github.com/akrem-chabchoub/fibonacci-api/internal/handler"
)

func newMux() http.Handler {
	router := httprouter.New()
	fib := handler.Fibonacci{}
	router.Handler(http.MethodGet, "/v1/fibonacci", fib)
	router.Handler(http.MethodGet, "/healthz", http.HandlerFunc(ready))
	router.Handler(http.MethodGet, "/readyz", http.HandlerFunc(ready))
	router.Handler(http.MethodGet, "/openapi.yaml", http.HandlerFunc(api.ServeOpenAPI))
	router.Handler(http.MethodGet, "/docs", http.HandlerFunc(api.ServeDocs))

	return withCORS(router)
}

func ready(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, If-None-Match, X-Request-Id")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(w, r)
	})
}
