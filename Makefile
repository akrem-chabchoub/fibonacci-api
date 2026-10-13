ARTIFACTS ?= _out
CGO_ENABLED ?= 0
GO ?= go

.PHONY: binary test test-race lint

binary:
	@mkdir -p $(ARTIFACTS)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(ARTIFACTS)/fibonacci-api ./cmd/fibonacci-api

test:
	$(GO) test ./...

test-race:
	CGO_ENABLED=1 $(GO) test -race ./...

lint:
	golangci-lint run ./...
