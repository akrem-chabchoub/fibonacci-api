package api

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.yaml
var openAPI []byte

//go:embed docs.html
var docs []byte

func ServeOpenAPI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	_, _ = w.Write(openAPI)
}

func ServeDocs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(docs)
}
