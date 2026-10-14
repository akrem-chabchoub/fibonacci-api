ARTIFACTS ?= _out
CGO_ENABLED ?= 0
GO ?= go
IMAGE ?= fibonacci-api:local

.PHONY: binary dev test test-integration test-race test-cover test-cover-html lint docker-image generate

binary:
	@mkdir -p $(ARTIFACTS)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(ARTIFACTS)/fibonacci-api ./cmd/fibonacci-api

dev:
	@command -v air >/dev/null || { echo "go install github.com/air-verse/air@v1.67.4"; exit 1; }
	air

test:
	$(GO) test ./...

test-integration:
	$(GO) test ./internal/server -count=1

test-race:
	CGO_ENABLED=1 $(GO) test -race ./...

lint:
	golangci-lint run ./...

test-cover:
	@mkdir -p $(ARTIFACTS)
	$(GO) test -coverprofile=$(ARTIFACTS)/coverage.out ./...
	$(GO) tool cover -func=$(ARTIFACTS)/coverage.out

test-cover-html: test-cover
	$(GO) tool cover -html=$(ARTIFACTS)/coverage.out -o $(ARTIFACTS)/coverage.html

docker-image:
	docker build -t $(IMAGE) .

generate:
	$(GO) generate ./api/
