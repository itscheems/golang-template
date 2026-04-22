# golang-template

<div align="left">

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)
[![CI](https://github.com/itscheems/golang-template/actions/workflows/ci.yml/badge.svg)](https://github.com/itscheems/golang-template/actions/workflows/ci.yml)

</div>

A modern, **multi-module Go workspace** template for monorepos with shared libraries and one or more deployable services.

## Overview

This template uses a [`go.work`](https://go.dev/ref/mod#workspaces) workspace so that several Go modules — typically a handful of deployable services and a few shared libraries — live side by side and import each other through their canonical module paths, with no published versions or `replace` gymnastics in user-facing code.

It ships pre-wired with:

- Pinned Go toolchain and tooling via [**mise**](https://mise.jdx.dev)
- [**Task**](https://taskfile.dev) recipes that iterate every module in the workspace
- [**golangci-lint v2**](https://golangci-lint.run) (per-module) with integrated formatters (gofumpt + goimports)
- [**GoReleaser v2**](https://goreleaser.com) with multi-platform binaries, SBOMs, and per-service builds
- [**govulncheck**](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) running per module on every CI build
- Distroless `nonroot` Docker image that builds any service via `--build-arg SERVICE=…`
- Conventional-commits driven CHANGELOG ([**git-cliff**](https://git-cliff.org))
- License header enforcement ([**hawkeye**](https://github.com/korandoru/hawkeye)) — `.go` files use `//` line comments
- Dependabot for every `go.mod`, GitHub Actions, and Docker

> **Note**: This template targets **multi-module monorepos**. For a small single-module project, use [`simple-golang-template`](https://github.com/itscheems/simple-golang-template) instead.

## Layout

```text
.
├── go.work                                  # workspace declaration
├── services/                                # deployable binaries; each is its own module
│   └── app/
│       ├── go.mod                           # module .../services/app
│       ├── cmd/app/main.go                  # entry point
│       └── internal/
│           └── version/                     # build-time stamps via -ldflags
├── libs/                                    # shared libraries; each is its own module
│   └── shared/
│       ├── go.mod                           # module .../libs/shared
│       └── log/                             # example: slog wrapper imported by services
├── deploy/
│   ├── Dockerfile                           # distroless, non-root, ARG SERVICE selects which one
│   └── .dockerignore
├── docs/
├── .github/
│   ├── dependabot.yml
│   └── workflows/{ci,cd,pr-check-ci}.yml
├── .mise.toml                               # pinned Go + tool versions
├── Taskfile.yml                             # workspace-aware recipes
├── .golangci.yml                            # shared lint config (run per module)
├── .goreleaser.yaml                         # release pipeline
├── cliff.toml                               # CHANGELOG generation
├── licenserc.toml                           # license header rules
├── .editorconfig
├── .gitignore
├── LICENSE                                  # Apache-2.0
├── CHANGELOG.md
└── README.md
```

### Adding a new module

1. Create the directory, e.g. `services/worker/cmd/worker/main.go`.
2. `cd services/worker && go mod init github.com/itscheems/golang-template/services/worker`.
3. Add it to the workspace: `go work use ./services/worker`.
4. If it imports another local module, add a `require … v0.0.0-…` line **and** a matching `replace` directive in its `go.mod` (see `services/app/go.mod` for the pattern). The replace lets standalone tooling (`go mod tidy`, single-module builds) resolve cross-module imports the same way `go.work` does.
5. Update `MODULES` in `Taskfile.yml` and the `gomod` entries in `.github/dependabot.yml`.
6. For releases, add a new `builds:` entry to `.goreleaser.yaml` pointing at the new `dir:`.

## Prerequisites

All tool versions are pinned in `.mise.toml`. The fast path:

```bash
# Install mise once: https://mise.jdx.dev/getting-started.html
mise install     # installs Go, golangci-lint, goreleaser, task, gofumpt, ...
mise activate    # or add `eval "$(mise activate zsh)"` to your shell rc
```

Manual installation needs at least: [Go](https://go.dev/dl/) 1.26.x, [Task](https://taskfile.dev/installation/) 3.x, [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/) v2.x, and [hawkeye](https://github.com/korandoru/hawkeye).

## Quick Start

```bash
# 1. Use this template on GitHub ("Use this template" → "Create a new repository")
# 2. Clone your new repo
git clone https://github.com/<you>/<your-repo>.git
cd <your-repo>

# 3. Rewrite the module path (one editor pass)
OLD=github.com/itscheems/golang-template
NEW=github.com/<you>/<your-repo>
git ls-files | xargs sed -i '' "s|$OLD|$NEW|g"   # macOS; drop '' on Linux
go work sync && (cd services/app && go mod tidy)

# 4. Install dev tools
mise install

# 5. Run everything CI runs, locally
task ci
```

## Development

```bash
task                       # list available tasks
task run                   # go run services/app/cmd/app
task run SERVICE=worker    # … any other service
task build                 # build dist/app
task fmt                   # gofumpt + goimports + license headers
task lint                  # golangci-lint (per module) + hawkeye + go vet
task test                  # race-detector + per-module coverage
task audit                 # govulncheck + gosec across every module
task ci                    # lint + test + build (everything CI runs)
task release-snapshot      # dry-run a GoReleaser build
```

Aliases: `task l` (lint), `task t` (test).

### Without `task`

```bash
# Per-module
cd services/app
go run ./cmd/app
go build -o ../../dist/app ./cmd/app
go test -race ./...
golangci-lint run

# Workspace-wide
go work sync
```

## Docker

Build any service by passing `--build-arg SERVICE=…`. Build context **must** be the repo root so the workspace and shared libs are visible.

```bash
docker build -f deploy/Dockerfile --build-arg SERVICE=app -t golang-template:dev .
docker run --rm golang-template:dev
```

The image is `gcr.io/distroless/static-debian12:nonroot` — no shell, no package manager, runs as UID `nonroot` (65532).

## Releasing

Push a semver tag and the `cd.yml` workflow will:

1. Generate `CHANGELOG.md` with git-cliff and commit it back to `main`.
2. Run GoReleaser to build binaries for `linux|darwin|windows` × `amd64|arm64` for every service in `.goreleaser.yaml`.
3. Generate a CycloneDX SBOM per archive (Syft).
4. Publish a GitHub Release with checksums, archives, and notes.

```bash
git tag v0.1.0
git push origin v0.1.0
```

For per-service release tags (e.g. `app/v0.1.0`, `worker/v0.2.0`), enable [GoReleaser's monorepo mode](https://goreleaser.com/customization/monorepo/).

## Conventions

- **Commits**: [Conventional Commits](https://www.conventionalcommits.org) — drives changelog grouping.
- **Go style**: `gofumpt` (superset of `gofmt`) + `goimports` with `-local github.com/itscheems/golang-template`.
- **Layout**: Go's [`internal/`](https://go.dev/ref/spec#Implementation_restriction:_Programs) rule is honoured — `services/app/internal/...` is unimportable from outside `services/app`.
- **Cross-module imports**: shared code goes under `libs/`, requires both a `require` line **and** a `replace` directive in the consumer's `go.mod` (see "Adding a new module").
- **License headers**: every `.go` file carries the Apache-2.0 header; enforced by `hawkeye check` in CI.

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
