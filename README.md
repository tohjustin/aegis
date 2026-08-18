<p align="center">
  <img src="./docs/title.png">
</p>

## About

<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/tohjustin/aegis)](https://goreportcard.com/report/github.com/tohjustin/aegis) -->
<!-- [![CodeCov](https://codecov.io/gh/tohjustin/aegis/branch/master/graph/badge.svg?token=HRJhI2JVS0)](https://codecov.io/gh/tohjustin/aegis) -->
<!-- [![GoDoc](https://godoc.org/github.com/tohjustin/aegis/pkg/badge?status.svg)](http://godoc.org/github.com/tohjustin/aegis/pkg/badge/) -->
[![CI](https://github.com/tohjustin/aegis/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/tohjustin/aegis/actions/workflows/ci.yml)
[![Font Awesome](https://aegisbadges.appspot.com/static?icon=brands/font-awesome&subject=Font%20Awesome&status=v7.2.0&color=318FE0)](https://fontawesome.com/)
[![Release](https://aegisbadges.appspot.com/static?subject=release&status=v1.1.0)](https://github.com/tohjustin/aegis/releases)
[![License](https://aegisbadges.appspot.com/static?subject=license&status=MIT)](https://opensource.org/licenses/MIT)

Aegis is a SVG badge generation service, icons powered by [Font Awesome](https://fontawesome.com/). See the [changelog](./CHANGELOG.md) for notable changes.

## Usage

### Query Parameters

> NOTE: You can use these query parameters to modify any badges generated from badge services below.

Query parameters for configuring your badge appearance.

| Query Parameter | Description                  | Input Format                                                                                       | Example                                       |
| --------------- | ---------------------------- | -------------------------------------------------------------------------------------------------- | --------------------------------------------- |
| color           | Sets the badge primary color | RGB Hex Values, [CSS Color Keywords](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value) | "fff", "1BACBF", "mediumturquoise"            |
| icon            | Sets the badge icon          | Any one of the available [Font Awesome Icons](https://fontawesome.com/icons): `<STYLE>/<NAME>`     | "brands/github", "regular/star", "solid/star" |
| status          | Sets the badge status text   | Any URL-encoded string                                                                             | "Build%20Status", "ビルド状態"                           |
| style           | Sets the badge style         | Any one of the 4 available badge styles (classic, flat, plastic, semaphoreci)                      | "classic", "flat", "plastic", "semaphoreci"   |
| subject         | Sets the badge subject text  | Any URL-encoded string                                                                             | "Failed", "失敗"                                  |

### Static Badge Service

| Path                             | Description            | Example                                                                                                           |
| -------------------------------- | ---------------------- | ----------------------------------------------------------------------------------------------------------------- |
| [/static?subject=release&status=v2.0.0](https://aegisbadges.appspot.com/static?subject=release&status=v2.0.0)                          | Static badge           | ![static](https://aegisbadges.appspot.com/static?subject=release&status=v2.0.0)                                  |
| [/static?subject=style&status=classic&style=classic](https://aegisbadges.appspot.com/static?subject=style&status=classic&style=classic)<br>[/static?subject=style&status=flat&style=flat](https://aegisbadges.appspot.com/static?subject=style&status=flat&style=flat)<br>[/static?subject=style&status=plastic&style=plastic](https://aegisbadges.appspot.com/static?subject=style&status=plastic&style=plastic)<br>[/static?subject=style&status=semaphoreci&style=semaphoreci](https://aegisbadges.appspot.com/static?subject=style&status=semaphoreci&style=semaphoreci) | With various badge styles | ![static](https://aegisbadges.appspot.com/static?subject=style&status=classic&style=classic)<br>![static](https://aegisbadges.appspot.com/static?subject=style&status=flat&style=flat)<br>![static](https://aegisbadges.appspot.com/static?subject=style&status=plastic&style=plastic)<br>![static](https://aegisbadges.appspot.com/static?subject=style&status=semaphoreci&style=semaphoreci) |
| [/static?subject=license&status=AGPL%20v3&icon=solid/balance-scale](https://aegisbadges.appspot.com/static?subject=license&status=AGPL%20v3&icon=solid/balance-scale) | With icon | ![static](https://aegisbadges.appspot.com/static?subject=license&status=AGPL%20v3&icon=solid/balance-scale) |
| [/static?subject=ビルド状態&status=成功&color=26A876](https://aegisbadges.appspot.com/static?subject=ビルド状態&status=成功&color=26A876) | With non-english characters | ![static](https://aegisbadges.appspot.com/static?subject=ビルド状態&status=成功&color=26A876) |

### Bitbucket Badge Service

[![Bitbucket Cloud REST API](https://aegisbadges.appspot.com/static?icon=brands/bitbucket&subject=Bitbucket%20Cloud%20REST%20API&status=v2.0)](https://developer.atlassian.com/bitbucket/api/2/reference/)

| Path                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | Description        | Example                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| /bitbucket/forks/`<USERNAME>`/`<REPO_SLUG>`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | Fork count         | ![bitbucket/forks](https://aegisbadges.appspot.com/bitbucket/forks/atlassian/aui-react?)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| /bitbucket/issues/`<USERNAME>`/`<REPO_SLUG>`<br>/bitbucket/issues/`<USERNAME>`/`<REPO_SLUG>`?state=new<br>/bitbucket/issues/`<USERNAME>`/`<REPO_SLUG>`?state=open<br>/bitbucket/issues/`<USERNAME>`/`<REPO_SLUG>`?state=resolved<br>/bitbucket/issues/`<USERNAME>`/`<REPO_SLUG>`?state=on-hold<br>/bitbucket/issues/`<USERNAME>`/`<REPO_SLUG>`?state=invalid<br>/bitbucket/issues/`<USERNAME>`/`<REPO_SLUG>`?state=duplicate<br>/bitbucket/issues/`<USERNAME>`/`<REPO_SLUG>`?state=wontfix<br>/bitbucket/issues/`<USERNAME>`/`<REPO_SLUG>`?state=closed<br> | Issue count        | ![bitbucket/issues](https://aegisbadges.appspot.com/bitbucket/issues/atlassian/aui-react)<br>![bitbucket/new-issues](https://aegisbadges.appspot.com/bitbucket/issues/atlassian/aui-react?state=new)<br>![bitbucket/open-issues](https://aegisbadges.appspot.com/bitbucket/issues/atlassian/aui-react?state=open)<br>![bitbucket/resolved-issues](https://aegisbadges.appspot.com/bitbucket/issues/atlassian/aui-react?state=resolved)<br>![bitbucket/on-hold-issues](https://aegisbadges.appspot.com/bitbucket/issues/atlassian/aui-react?state=on-hold)<br>![bitbucket/invalid-issues](https://aegisbadges.appspot.com/bitbucket/issues/atlassian/aui-react?state=invalid)<br>![bitbucket/duplicate-issues](https://aegisbadges.appspot.com/bitbucket/issues/atlassian/aui-react?state=duplicate)<br>![bitbucket/wontfix-issues](https://aegisbadges.appspot.com/bitbucket/issues/atlassian/aui-react?state=wontfix)<br>![bitbucket/closed-issues](https://aegisbadges.appspot.com/bitbucket/issues/atlassian/aui-react?state=closed)<br> |
| /bitbucket/pull-requests/`<USERNAME>`/`<REPO_SLUG>`<br>/bitbucket/pull-requests/`<USERNAME>`/`<REPO_SLUG>`?state=open<br>/bitbucket/pull-requests/`<USERNAME>`/`<REPO_SLUG>`?state=declined<br>/bitbucket/pull-requests/`<USERNAME>`/`<REPO_SLUG>`?state=merged<br>/bitbucket/pull-requests/`<USERNAME>`/`<REPO_SLUG>`?state=superseded<br>                                                                                                                                                                                                                 | Pull Request count | ![bitbucket/pull-requests](https://aegisbadges.appspot.com/bitbucket/pull-requests/atlassian/aui-react)<br>![bitbucket/open-pull-requests](https://aegisbadges.appspot.com/bitbucket/pull-requests/atlassian/aui-react?state=open)<br>![bitbucket/declined-pull-requests](https://aegisbadges.appspot.com/bitbucket/pull-requests/atlassian/aui-react?state=declined)<br>![bitbucket/merged-pull-requests](https://aegisbadges.appspot.com/bitbucket/pull-requests/atlassian/aui-react?state=merged)<br>![bitbucket/superseded-pull-requests](https://aegisbadges.appspot.com/bitbucket/pull-requests/atlassian/aui-react?state=superseded)                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |

### GitHub Badge Service

[![GitHub GraphQL API](https://aegisbadges.appspot.com/static?icon=brands/github&subject=GitHub%20GraphQL%20API&status=v4)](https://developer.github.com/v4/)

| Path                                                                                                                                                                                                                                          | Description        | Example                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| /github/forks/`<OWNER>`/`<REPOSITORY>`                                                                                                                                                                                                        | Fork count         | ![github/forks](https://aegisbadges.appspot.com/github/forks/google/gopacket)                                                                                                                                                                                                                                                                                                                                                                                 |
| /github/issues/`<OWNER>`/`<REPOSITORY>`<br>/github/issues/`<OWNER>`/`<REPOSITORY>`?state=open<br>/github/issues/`<OWNER>`/`<REPOSITORY>`?state=closed<br>                                                                                     | Issue count        | ![github/issues](https://aegisbadges.appspot.com/github/issues/google/gopacket)<br>![github/open-issues](https://aegisbadges.appspot.com/github/issues/google/gopacket?state=open)<br>![github/closed-issues](https://aegisbadges.appspot.com/github/issues/google/gopacket?state=closed)                                                                                                                                                                 |
| /github/pull-requests/`<OWNER>`/`<REPOSITORY>`<br>/github/pull-requests/`<OWNER>`/`<REPOSITORY>`?state=open<br>/github/pull-requests/`<OWNER>`/`<REPOSITORY>`?state=closed<br>/github/pull-requests/`<OWNER>`/`<REPOSITORY>`?state=merged<br> | Pull Request count | ![github/pull-requests](https://aegisbadges.appspot.com/github/pull-requests/google/gopacket)<br>![github/open-pull-requests](https://aegisbadges.appspot.com/github/pull-requests/google/gopacket?state=open)<br>![github/closed-pull-requests](https://aegisbadges.appspot.com/github/pull-requests/google/gopacket?state=closed)<br>![github/merged-pull-requests](https://aegisbadges.appspot.com/github/pull-requests/google/gopacket?state=merged) |
| /github/stars/`<OWNER>`/`<REPOSITORY>`                                                                                                                                                                                                        | Star count         | ![github/stars](https://aegisbadges.appspot.com/github/stars/google/gopacket)                                                                                                                                                                                                                                                                                                                                                                                 |

### GitLab Badge Service

[![GitLab API](https://aegisbadges.appspot.com/static?icon=brands/gitlab&subject=GitLab%20API&status=v4)](https://docs.gitlab.com/ee/api/)

| Path                                                                                                                                                                                                                                                                                                                                              | Description         | Example                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| /gitlab/forks/`<NAMESPACE>`/`<PROJECT_NAME>`                                                                                                                                                                                                                                                                                                      | Fork count          | ![gitlab/forks](https://aegisbadges.appspot.com/gitlab/forks/gitlab-org/gitaly)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| /gitlab/issues/`<NAMESPACE>`/`<PROJECT_NAME>`<br>/gitlab/issues/`<NAMESPACE>`/`<PROJECT_NAME>`?state=opened<br>/gitlab/issues/`<NAMESPACE>`/`<PROJECT_NAME>`?state=closed<br>                                                                                                                                                                     | Issue count         | ![gitlab/issues](https://aegisbadges.appspot.com/gitlab/issues/gitlab-org/gitaly)<br>![gitlab/opened-issues](https://aegisbadges.appspot.com/gitlab/issues/gitlab-org/gitaly?state=opened)<br>![gitlab/closed-issues](https://aegisbadges.appspot.com/gitlab/issues/gitlab-org/gitaly?state=closed)<br>                                                                                                                                                                                                                                                                                                      |
| /gitlab/merge-requests/`<NAMESPACE>`/`<PROJECT_NAME>`<br>/gitlab/merge-requests/`<NAMESPACE>`/`<PROJECT_NAME>`?state=opened<br>/gitlab/merge-requests/`<NAMESPACE>`/`<PROJECT_NAME>`?state=closed<br>/gitlab/merge-requests/`<NAMESPACE>`/`<PROJECT_NAME>`?state=locked<br>/gitlab/merge-requests/`<NAMESPACE>`/`<PROJECT_NAME>`?state=merged<br> | Merge Request count | ![gitlab/merge-requests](https://aegisbadges.appspot.com/gitlab/merge-requests/gitlab-org/gitaly)<br>![gitlab/opened-merge-requests](https://aegisbadges.appspot.com/gitlab/merge-requests/gitlab-org/gitaly?state=opened)<br>![gitlab/closed-merge-requests](https://aegisbadges.appspot.com/gitlab/merge-requests/gitlab-org/gitaly?state=closed)<br>![gitlab/locked-merge-requests](https://aegisbadges.appspot.com/gitlab/merge-requests/gitlab-org/gitaly?state=locked)<br>![gitlab/merged-merge-requests](https://aegisbadges.appspot.com/gitlab/merge-requests/gitlab-org/gitaly?state=merged)<br> |
| /gitlab/stars/`<NAMESPACE>`/`<PROJECT_NAME>`<br>                                                                                                                                                                                                                                                                                                  | Star count          | ![gitlab/stars](https://aegisbadges.appspot.com/gitlab/stars/gitlab-org/gitaly)<br>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |

## Getting Started

A GitHub access token is required at process start (`GITHUB_ACCESS_TOKEN` or `-github-access-token`).

```shell
make build
./bin/aegis -github-access-token "$GITHUB_ACCESS_TOKEN"
```

Or with Docker:

```shell
make docker-run GITHUB_ACCESS_TOKEN="$GITHUB_ACCESS_TOKEN"
```

### Make targets

| Target | Description |
| ------ | ----------- |
| `check-deps` | Verify `go` and `golangci-lint` are on `PATH` |
| `build` | Build `bin/aegis` (`CGO_ENABLED=0`, trimmed, version ldflags) |
| `build-all` | Cross-compile linux/darwin amd64/arm64 |
| `run` | Build and run |
| `vet` | `go vet ./...` |
| `lint` | `golangci-lint run` |
| `tidy` | `go mod tidy` |
| `test` | Unit tests |
| `test-e2e` | Blackbox tests against a real binary |
| `test-all` | Unit + e2e |
| `docker-build` | Distroless static non-root image |
| `docker-run` | Build and run, publishing 8080 and 9090 |

### Endpoints

| Listener | Path | Description |
| -------- | ---- | ----------- |
| app `:8080` | `GET /static` | Static badge |
| app `:8080` | `GET /github/{method}/{owner}/{repo}` | GitHub badge |
| app `:8080` | `GET /gitlab/{method}/{owner}/{repo}` | GitLab badge |
| app `:8080` | `GET /bitbucket/{method}/{owner}/{repo}` | Bitbucket badge |
| admin `:9090` | `GET /livez` | Liveness |
| admin `:9090` | `GET /readyz` | Readiness |
| admin `:9090` | `GET /metrics` | Prometheus metrics |
| admin `:9090` | `GET /version` | Build metadata |

### Configuration

Precedence: built-in defaults → environment variables → CLI flags.

| Flag | Env | Default | Description |
| ---- | --- | ------- | ----------- |
| `-app-addr` | `APP_ADDR` | `:8080` | Application listen address |
| `-admin-addr` | `ADMIN_ADDR` | `:9090` | Admin listen address |
| `-log-level` | `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `-log-format` | `LOG_FORMAT` | `json` | `json` or `text` |
| `-read-timeout` | `READ_TIMEOUT` | `15s` | Request read timeout |
| `-read-header-timeout` | `READ_HEADER_TIMEOUT` | `5s` | Header read timeout |
| `-write-timeout` | `WRITE_TIMEOUT` | `15s` | Response write timeout |
| `-idle-timeout` | `IDLE_TIMEOUT` | `60s` | Keep-alive idle timeout |
| `-shutdown-timeout` | `SHUTDOWN_TIMEOUT` | `15s` | Graceful shutdown timeout |
| `-histogram-buckets` | `HISTOGRAM_BUCKETS` | Prometheus defaults | HTTP duration buckets |
| `-exclude-cache-control-headers` | `EXCLUDE_CACHE_CONTROL_HEADERS` | `false` | Omit `Cache-Control` on badges |
| `-root-redirect-url` | `ROOT_REDIRECT_URL` | empty | Redirect target for `GET /` |
| `-github-access-token` | `GITHUB_ACCESS_TOKEN` | required | GitHub API token |
| `-v` / `--version` | | | Print version and exit |

### Layout

```
cmd/aegis/                 # process wiring and lifecycle
internal/config/           # defaults → env → flags
internal/observability/    # slog + Prometheus
internal/server/           # app + admin HTTP
internal/gitprovider/      # GitHub / GitLab / Bitbucket clients
internal/version/          # ldflags-injected build info
pkg/badge/                 # public SVG badge library
test/e2e/                  # blackbox tests (//go:build e2e)
deploy/kubernetes/         # Deployment, PDB, NetworkPolicy, ServiceMonitor
```

Kubernetes manifests live in [`deploy/kubernetes`](./deploy/kubernetes). The container image is `gcr.io/distroless/static-debian12:nonroot` (no shell, no Docker `HEALTHCHECK`; probes belong on `:9090`).

## License

Aegis is [MIT licensed](./LICENSE).
