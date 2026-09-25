# Copilot Instructions for sso-dashboard

This repo aggregates Keycloak event logs for the BC Gov SSO service. It has one
main Go service (`aggregator/`) plus supporting infra/config directories
(`helm/`, `grafana/`, `loki-authorizer/`, `service-account-generator/`,
`scripts/`). Most day-to-day code work happens in `aggregator/`.

## Architecture

Log flow: **Grafana Alloy** (log collector on Keycloak) → **Loki** (log
store, backed by S3) → **aggregator** (Go HTTP server) → **Postgres** →
**Grafana** (dashboards, queries both Loki and the aggregator's Postgres DB).

`aggregator/` contains two independently-deployed Go binaries sharing the same
module and `model` package:

- **`cmd/aggregator`** — HTTP server. Receives log push requests from Alloy at
  `/api/promtail/push` (`promtail/index.go`), which relies on Grafana Loki's
  own `github.com/grafana/loki/pkg/loghttp/push.ParseRequest` to parse the
  push payload instead of a custom parser. It extracts labels
  (`environment`, `realmId`, `clientId`, `eventType`, `username`, `timestamp`)
  from each log stream, derives the IDP from the username
  (`idpFromUsername`), and upserts one row per unique label combination into
  `client_events_with_idp` via `model.UpsertClientEventWithIDP`. Streams
  older than 24h or missing required labels are skipped.
- **`cmd/compactor`** — no HTTP server; runs scheduled jobs via
  `go-co-op/gocron` (`model.RunEventsJob`, `model.RunSessionsJob`): deleting
  old client events (retention controlled by `RETENTION_PERIOD` env var) and
  collecting client session counts from Keycloak (notifying RocketChat on
  failures via `webhooks/`).
- **`model/`** — all Postgres access (via `github.com/go-pg/pg/v10`, raw SQL
  through `pgdb.Query(nil, query, params...)`, no ORM/model-struct usage).
  `db.go` owns the shared `*pg.DB` singleton (`GetDB()`), configured from
  `config.LoadDatabaseConfig()` — both binaries defer `model.GetDB().Close()`.
- **`config/`** — env-driven config loading (DB connection, timezone/location).
- **`keycloak/`** — Keycloak REST client + OAuth token manager
  (`tokenManager.go`) used by the compactor to fetch session counts.
- Python (`alembic`, `models.py`, `database.py`) is used **only** for DB
  schema migrations, not application logic — the Go code never uses an ORM,
  so schema changes must be kept in sync manually between `models.py`
  (SQLAlchemy, used to autogenerate migrations) and the raw SQL in
  `model/*.go`.

## Build, run, and test

All commands below are run from `aggregator/` (the Go module root):

```sh
go build -o build/aggregator ./cmd/aggregator/main.go   # or: make build
go build -o build/compactor ./cmd/compactor/main.go
gofmt -w -s .                                            # or: make format
```

Local dev loops (uses `reflex` to rebuild on `.go` changes):

```sh
make aggregator-dev   # runs ./cmd/aggregator
make compactor-dev    # runs ./cmd/compactor
```

Tests live next to the code they test (`keycloak/*_test.go`,
`model/*_test.go`) and are run per-package, not with a repo-wide `go test
./...` (see `.github/workflows/test.yaml`):

```sh
cd keycloak && go test -v
cd model && go test -v

# single test:
go test -run TestName ./keycloak/...
```

CI (`test.yaml`) triggers only on changes under `aggregator/**` and pins Go
via `go-version: ^1.21.0` in the workflow — keep that in sync if the `go`
directive in `go.mod` changes.

Local Postgres/full-stack dev environment (Keycloak, Postgres, Alloy, Loki,
Grafana) is provisioned via `aggregator/localdev/docker-compose.yaml` — see
`aggregator/README.md` for the full setup sequence and required `DB_PORT`
env var when using docker-compose (`5433` vs local Postgres default).

## Key conventions

- Do not ignore the /plan command when running copilot queries.

- **Dependency pinning quirks**: `github.com/grafana/loki` is imported at a
  pre-v1-module-path pseudo-version (Loki's own go.mod doesn't use semantic
  import versioning even past v1, so bumping it means resolving a specific
  commit SHA via `go get github.com/grafana/loki@<sha>`, not a `vX.Y.Z` tag
  directly — the target tag's commit must be looked up first). The
  `github.com/grafana/loki/pkg/push` submodule must be pinned to the exact
  same commit as the main `loki` module or the build breaks
  (`undefined: push.LabelAdapter` is the symptom of a mismatch).
  `github.com/go-pg/pg`, by contrast, *does* use standard semantic import
  versioning — v10+ requires the `/v10` path suffix
  (`github.com/go-pg/pg/v10`).
- `go.mod` mirrors select `replace` directives from `grafana/loki`'s own
  `go.mod` (e.g. `hashicorp/memberlist` fork, `grafana/regexp` pin, pinned
  `otelhttp` version) to keep the resolved dependency graph consistent with
  what Loki itself is tested against — check Loki's go.mod for updates when
  bumping the Loki dependency.
- The Go version in `go.mod`, the builder base image in `aggregator/Dockerfile`
  (`FROM golang:X-...`), and the `go-version` in
  `.github/workflows/test.yaml` should be kept in lockstep.
- Deployment is directory-driven: GitHub Actions workflows filter on
  `paths:` for `aggregator/**`, `helm/aggregator/**`, etc., so a change only
  triggers the relevant deploy workflow (`deploy-aggregator.yaml`,
  `deploy-alloy.yaml`, `deploy-dashboard.yaml`, `terraform.yaml`). Aggregator
  images are built via `docker/build-push-action` and deployed with Helm to
  OpenShift (`oc-login` + `helm upgrade --install`), not raw `kubectl`.
