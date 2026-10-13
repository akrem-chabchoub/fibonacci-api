# fibonacci-api

HTTP service that returns the first n Fibonacci numbers.

## build

```sh
make binary
```


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

CI (`.github/workflows/ci.yml`): lint, commitlint on pull requests, test-race + test-cover, binary. Coverage floor is [go-test-coverage](https://github.com/vladopajic/go-test-coverage) (`.testcoverage.yml`): 80% on `internal/fibonacci` and `internal/api` once they exist.
