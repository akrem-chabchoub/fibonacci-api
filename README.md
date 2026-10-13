# fibonacci-api

HTTP service that returns the first n Fibonacci numbers.

## build

```sh
make binary
```

## test

```sh
make test       # go test -race (default, what CI runs)
make test-fast  # no race, for a tight TDD loop
```
