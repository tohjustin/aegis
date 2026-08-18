# Stdlib HTTP service refactor

Date: 2026-08-18

## Context

Aegis was a single-listener badge HTTP service using gorilla/mux, cobra, and zap, with application code in `service/` and a public badge library in `pkg/badge`. The process mixed CLI wiring, HTTP handlers, and git-forge API clients. Timeouts were 2s, shutdown had no deadline, and the Docker image was `scratch` with a `sh -c` entrypoint that cannot run.

## Design

Apply the Go HTTP service conventions:

- Keep `pkg/badge` as a public library (import path is already documented).
- Move the HTTP process to `cmd/aegis` + `internal/{config,observability,server,gitprovider,version}`.
- Dual listeners: app `:8080`, admin `:9090` (`/livez`, `/readyz`, `/metrics`, `/version`).
- Stdlib `flag` (defaults → env → flags), `log/slog`, ServeMux method+path patterns.
- Dedicated Prometheus registry; HTTP `path` label is the mux pattern.
- Domain fetches live in `internal/gitprovider` (no `http.Handler` types).
- App errors that users see as images stay SVG badges. JSON `{"error":...}` is used for admin/panic 5xx so internals are not leaked.
- Distroless static non-root image; Kubernetes manifests with probes on the admin port.

### Assumptions

- GitHub access token remains required at process start (previous fail-fast behavior).
- Default timeouts move from 2s to 15s read/write (2s was too tight for forge APIs).
- CLI flags are kebab-case stdlib flags (`-app-addr`, `-github-access-token`). `--port` and millisecond timeout integers are removed.
- NetworkPolicy allows HTTPS egress so GitHub/GitLab/Bitbucket fetches work, in addition to kube-dns.

## Verification

- `go test ./...`
- `go test -tags e2e -count=1 ./test/e2e/...`
- `golangci-lint run` when installed
- `make build` and `GET /readyz` after `make run` with `GITHUB_ACCESS_TOKEN` set
