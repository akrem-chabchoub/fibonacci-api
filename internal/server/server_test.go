package server

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRunServesFibonacciAndShutsDown(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	srv := New(Options{Addr: "127.0.0.1:0"})
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Run(ctx)
	}()

	select {
	case <-srv.Started():
	case err := <-errCh:
		t.Fatalf("server exited early: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("server did not start")
	}

	base := "http://" + srv.Addr()

	res, err := http.Get(base + "/v1/fibonacci?n=5")
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}

	if string(body) != "[0,1,1,2,3]" {
		t.Fatalf("body = %s", body)
	}

	health, err := http.Get(base + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	_ = health.Body.Close()

	if health.StatusCode != http.StatusOK {
		t.Fatalf("healthz status = %d", health.StatusCode)
	}

	spec, err := http.Get(base + "/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}

	specBody, err := io.ReadAll(spec.Body)
	_ = spec.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if spec.StatusCode != http.StatusOK {
		t.Fatalf("openapi status = %d", spec.StatusCode)
	}

	if !bytes.Contains(specBody, []byte("/v1/fibonacci")) {
		t.Fatal("openapi.yaml missing /v1/fibonacci")
	}

	docs, err := http.Get(base + "/docs")
	if err != nil {
		t.Fatal(err)
	}
	_ = docs.Body.Close()

	if docs.StatusCode != http.StatusOK {
		t.Fatalf("docs status = %d", docs.StatusCode)
	}

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("graceful shutdown timed out")
	}
}
