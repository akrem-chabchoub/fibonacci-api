package server

import (
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
	res.Body.Close()
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
	health.Body.Close()

	if health.StatusCode != http.StatusOK {
		t.Fatalf("healthz status = %d", health.StatusCode)
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
