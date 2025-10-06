# CloudMoor AI Contributor Guide

## Project snapshot
- Pre-alpha repo focused on designing a Go-based remote storage mounting platform (see `README.md`).
- Primary components: `cloudmoord` daemon + gRPC/REST APIs, `cloudmoor` CLI, optional React/Vite Web UI, future desktop tray.
- Strategic context lives in `docs/plan.md`; engineering conventions in `docs/spec.md`; execution backlog in `docs/tasks.md`.

## Architecture cues
- Separate control plane (daemon, scheduler, vault, APIs, UI) from data plane (rclone-backed VFS workers per mount); changes must respect that split.
- Connector plugins translate provider configs into a shared interface—initial focus on S3/MinIO, WebDAV, Dropbox (plan §§3–7).
- Persistence via SQLite (config, secrets, mounts, audit); future Postgres optional. Credential vault uses AES-GCM with OS keychain fallback (plan §6).

## Repository expectations
- Planned layout follows Standard Go project structure (`cmd/`, `internal/`, `pkg/`, `web/`, `deploy/`, `docs/`). Document new folders in `docs/spec.md` §3.
- When scaffolding Go module, use `go mod init github.com/binGhzal/cloudmoor` (tasks M0.1). Keep executable wiring in `cmd/<binary>/main.go` and domain logic under `internal/`.
- Frontend lives in `web/` with feature-first directories (`src/features/<name>/...`) as outlined in `docs/spec.md` §3.3.

## Implementation patterns
- Go services pass `context.Context` first, wrap errors with `%w`, log via shared `zap.Logger` (spec §4). Depend on constructor injection for pluggable pieces.
- CLI commands built with Cobra should expose JSON output (`--json`, `--watch`) and degrade gracefully if the daemon is unreachable (plan §3, tasks TCK-107).
- API contracts originate from protobuf definitions using `buf`; generate stubs into `internal/api` and expose REST via `grpc-gateway` under `/api/v1/...` (plan §4.6, spec §4.5).
- Web UI consumes generated clients per feature folder, enforces accessibility and Tailwind styling, and should proxy through the daemon for metrics/log streams (spec §5, plan §8).

## Workflows & tooling
- Target Go 1.22.x, Node 20+, golangci-lint, React Testing Library/Jest, and Testcontainers-go (spec §§4–6). Plan to wire Make targets (`make lint`, `make test`, `make build`) once Task M0.1 lands—mirror those commands in CI.
- Until Makefiles exist, run `gofmt/gofumpt`, `golangci-lint run`, and `go test ./...` manually; frontend uses `npm run lint` / `npm test` per spec.
- Integration tests spin up MinIO and WebDAV containers behind the `integration` build tag; prefer `testcontainers-go` helpers defined under `tests/integration/` (tasks TCK-108).

## Testing checklist
- Default to table-driven Go tests with `testify/require`; keep benchmarks in `_test.go` (spec §4.3).
- Mock connectors, HTTP servers, and key stores to avoid touching real services; record chaos cases (network drops, cache exhaustion) per plan §4.4/§5.3.
- Frontend tests live next to components (`Component.test.tsx`), mock network via `msw`, and run accessibility assertions with `jest-axe` (spec §5.4).

## Documentation & governance
- Every code change must sync `docs/plan.md`, `docs/spec.md`, and `docs/tasks.md` when scope, standards, or backlog evolve (plan §21, spec §7).
- Reference ticket IDs (TCK-###) in branches, commits, and PRs; keep roadmap alignment by updating checkboxes in `docs/tasks.md`.
- Record major decisions in `/docs/decisions/` using the milestone governance template once that directory exists (plan §21).

## Security expectations
- Encrypt credentials via the vault layer; never log secrets. Ensure TLS is on by default for daemon listeners with optional mTLS for remote admin (plan §6).
- Audit logging is mandatory for credential access, mount lifecycle, and UI actions—surface structured fields (`mount_id`, `request_id`) per spec §4.4 and plan §6.

## Collaboration guardrails
- Validate new work against milestone goals in `docs/tasks.md`; defer connectors or features explicitly tagged as post-v1 (plan §7.2, tasks backlog).
- Favor incremental, test-backed commits that keep CI green across macOS/Linux/Windows; note any OS-specific TODOs in `docs/operations/known-issues.md` when it appears.
- If introducing tooling or workflows, add corresponding lint/test badges or usage notes in `README.md` and reference them from the contributing guide once created.
