# syntax=docker/dockerfile:1

FROM golang:1.27.2-bookworm AS build
WORKDIR /src

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/fibonacci-api ./cmd/fibonacci-api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/fibonacci-api /fibonacci-api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/fibonacci-api"]
