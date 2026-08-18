# build stage
ARG GO_VERSION=1.26.2
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG VERSION=dev
ARG BUILD_DATE=unknown
ARG GIT_COMMIT=unknown
ARG GIT_TAG=unknown
ARG GIT_TREE_STATE=unknown

RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags="-s -w \
      -X github.com/tohjustin/aegis/internal/version.Version=${VERSION} \
      -X github.com/tohjustin/aegis/internal/version.BuildDate=${BUILD_DATE} \
      -X github.com/tohjustin/aegis/internal/version.GitCommit=${GIT_COMMIT} \
      -X github.com/tohjustin/aegis/internal/version.GitTag=${GIT_TAG} \
      -X github.com/tohjustin/aegis/internal/version.GitTreeState=${GIT_TREE_STATE}" \
    -o /out/aegis ./cmd/aegis

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/aegis /aegis
USER nonroot:nonroot
EXPOSE 8080 9090
ENTRYPOINT ["/aegis"]
