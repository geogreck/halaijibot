FROM golang:1.20.4-alpine3.17 AS base

WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS build
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/halajibot cmd/halajibot/*.go

FROM golangci/golangci-lint:v1.51.2-alpine AS lint-base

FROM base AS lint
RUN --mount=target=. \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    golangci-lint run --timeout 10m0s ./...

FROM golang:1.20.4-alpine3.17
COPY --from=build /out/halajibot /
ENTRYPOINT "/halajibot"