# fibonacci-api

HTTP service that returns the first n Fibonacci numbers.

## build

```sh
make binary
./_out/fibonacci-api --listen :8080
```

Ctrl+C (or SIGTERM) finishes in-flight requests, then exits. Default shutdown wait is 10s (`--shutdown-timeout`).

```sh
curl 'http://127.0.0.1:8080/v1/fibonacci?n=5'
# contract: http://127.0.0.1:8080/docs  and  /openapi.yaml
```

`make generate` refreshes `api/types.gen.go` from `api/openapi.yaml`.

## docker

```sh
make docker-image
docker compose up --build
```

## test

```sh
make test            # go test (TDD loop)
make test-race       # go test -race (what CI runs)
make test-cover      # writes _out/coverage.out
make test-cover-html # open the profile in a browser
```

## lint

```sh
make lint
```

Commit messages follow conventional commits (`feat:`, `fix:`, `chore:`, `ci:`, `test:`, `docs:`, `refactor:`), max 98 characters. GitHub Actions checks them with [commitlint](https://github.com/conventional-changelog/commitlint) (`.commitlintrc.yaml`). In the editor, the Conventional Commits extension can help write them.

CI (`.github/workflows/ci.yml`): lint, commitlint on pull requests, test-race + test-cover, binary. Coverage floor is [go-test-coverage](https://github.com/vladopajic/go-test-coverage) (`.testcoverage.yml`): 80% on `internal/service` and `internal/handler`.
