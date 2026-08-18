# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/2.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Refactored the HTTP process onto stdlib `net/http`, `flag`, and `log/slog`, with a dedicated Prometheus registry and dual listeners (`:8080` app, `:9090` admin).
- **Breaking:** replaced cobra flags. `-app-addr` / `-admin-addr` replace `--port`. Timeouts are Go durations (default 15s), not milliseconds. `-v` / `--version` replace the `version` subcommand.
- Docker image is now distroless static non-root; GitHub token is read from `GITHUB_ACCESS_TOKEN` instead of a shell-wrapped `CMD`.
- Make interface is now `build`, `test`, `test-e2e`, `lint`, `docker-build`, and related targets. Coverage/`codecov` CI upload was dropped in favor of vet, unit, e2e, lint, and docker buildx.

## [1.1.0] - 2026-04-16

### Added

- Default handler for unmatched HTTP routes.
- Option to redirect requests to the root path.
- Documentation examples for additional badge styles.

### Changed

- **Breaking:** upgraded Font Awesome from v5.12.0 to v7.2.0. Existing badges may need updated `icon` query values if identifiers changed between those versions.
- Raised the Go runtime baseline to 1.26.2.
- Migrated CI from CircleCI to GitHub Actions.

### Removed

- Dropped the deprecated `bump-version` Make target.
- Dropped `freebsd` and `solaris` from official cross-compile targets.
- Removed the FOSSA badge from documentation.

### Fixed

- Log errors only when badge generation fails.
- Corrected letter casing on error badges.
- Fixed Bitbucket badge links in documentation.

## [1.0.0] - 2020-01-29

Initial public release.

[Unreleased]: https://github.com/tohjustin/aegis/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/tohjustin/aegis/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/tohjustin/aegis/releases/tag/v1.0.0
