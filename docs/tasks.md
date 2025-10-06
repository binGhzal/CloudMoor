# CloudMoor Execution TODOs

Date: 2025-10-06

<!-- markdownlint-disable MD007 -->

## How to Read This Document

- **Checkboxes** track completion. Update them during sprint rituals or milestone reviews.
- **Hint** lines offer quick-start guidance or reminders about tooling.
- **Comment** lines capture dependencies, owners, and validation notes. Align comments with the governance rules in `plan.md#21`.

## Milestone M0 – Foundations (3 weeks)

- [ ] **Task M0.1 – Repository & Tooling Setup** _(Tickets: TCK-001, TCK-007)_
  - _Hint:_ Start from a clean working tree and scaffold directories before wiring CI.
  - _Comment:_ Blocks every downstream engineering task; target completion in week 1.
  - [x] **Subtask M0.1.1 – Initialize repository scaffold**
    - _Hint:_ Run `go mod init github.com/binGhzal/cloudmoor` and commit baseline folders (`cmd/`, `internal/`, `pkg/`, `web/`, `deploy/`).
    - _Comment:_ Ensure module path matches release namespace to avoid future import churn.
  - [x] **Action:** Commit `.editorconfig`, `.gitignore`, LICENSE, and placeholder READMEs for each top-level directory.
    - _Hint:_ Mirror formatting settings from Golang + frontend conventions to reduce lint noise.
    - _Comment:_ Link README stubs to `plan.md` for traceability.
  - [x] **Subtask M0.1.2 – Establish contribution standards**
    - _Hint:_ Use existing OSS templates as inspiration for `CONTRIBUTING.md` and issue/PR templates.
    - _Comment:_ Coordinate with Product/UX to include security disclosure contact.
  - [x] **Action:** Publish contribution guide covering branching, linting, testing, and release cadence.
    - _Hint:_ Reference CI commands verbatim so developers can copy/paste locally.
    - _Comment:_ Mark required checks for PR merge in GitHub settings after document merges.
  - [x] **Action:** Add GitHub issue templates (bug, feature, security) and PR checklist.
  - _Hint:_ Reuse acceptance criteria boilerplate from the Ticket Backlog section below to stay consistent.
    - _Comment:_ Set default labels (e.g., `needs-triage`) to streamline intake.
  - [x] **Subtask M0.1.3 – Capture engineering standards**
    - _Hint:_ Synthesize best-practice research into `docs/spec.md` covering Go, React, testing, and tooling conventions.
    - _Comment:_ Keep the spec aligned with `docs/plan.md` §13 and refresh alongside milestone updates.
    - [x] **Action:** Publish `docs/spec.md` and reference it from README, plan, and tasks.
      - _Hint:_ Link spec wherever contributors look for process guidance to minimise drift.
      - _Comment:_ Update action items here whenever standards evolve materially.

- [ ] **Task M0.2 – CI/CD Skeleton & Quality Gates** _(Tickets: TCK-002)_
  - _Hint:_ Implement workflows incrementally—start with lint/test, then add build matrix.
  - _Comment:_ Align pipeline steps with governance quality gates (lint, unit tests, security scan).
  - [x] **Subtask M0.2.1 – Configure GitHub Actions pipeline**
    - _Hint:_ Base workflow on official Go + golangci-lint reusable actions.
    - _Comment:_ Use Go 1.22 and 1.21 to cover current + previous stable releases.
  - [x] **Action:** Add job running gofmt, govet, and golangci-lint on pull requests.
    - _Hint:_ Fail fast on formatting errors to encourage pre-commit hooks.
    - _Comment:_ Document command invocation in contributing guide.
  - [x] **Action:** Add test job executing `go test ./...` across OS matrix (linux, macos, windows).
    - _Hint:_ Use `actions/cache` to speed up module downloads.
    - _Comment:_ Gate merge on this job succeeding.
  - [x] **Subtask M0.2.2 – Expose pipeline status**
    - _Hint:_ Embed status badges in `README.md` once workflow names are finalized.
    - _Comment:_ Provide troubleshooting steps for common CI failures.
  - [x] **Action:** Update README with lint/test badge markdown.
    - _Hint:_ Use shields.io badge URLs tied to workflow file name.
    - _Comment:_ Keep badge section above fold for quick visibility.
  - [x] **Action:** Document local lint/test workflow mirroring CI steps.
    - _Hint:_ Add `make lint` / `make test` commands to support automation.
    - _Comment:_ Ensure docs mention required tool versions (golangci-lint, npm, etc.).

- [ ] **Task M0.3 – Core Domain Abstractions** _(Tickets: TCK-003, TCK-004, TCK-005)_
  - _Hint:_ Design interfaces and persistence schema together to avoid churn.
  - _Comment:_ Coordinate with integrations engineer to validate connector API surface.
  - [x] **Subtask M0.3.1 – Define connector interfaces**
    - _Hint:_ Capture lifecycle hooks (Init, Validate, Mount, Teardown) and metadata (schema, display name).
    - _Comment:_ Keep surface minimal; future providers should not require interface changes.
    - [x] **Action:** Implement connector registry loader with unit tests ensuring deterministic ordering.
      - _Hint:_ Use build tags or plugin metadata structs for compile-time registration.
      - _Comment:_ Store registry manifest in JSON to support Web UI discovery.
      - _Completed:_
        - ✅ Defined `internal/connectors` package with `Connector` interface (`Init`, `ValidateConfig`, `Open`, `Metadata`) and `Connection` interface (`Ping`, `ProviderID`, `Close`).
        - ✅ Implemented `Config` type (map[string]interface{}) with `GetString` and `GetBool` helper methods for type-safe access.
        - ✅ Created thread-safe registry with `RegisterProvider`, `GetProvider`, `ListProviders` (preserves registration order), `ProviderIDs` (alphabetical), and `ExportManifest` (JSON serialization).
        - ✅ Authored comprehensive table-driven unit tests covering successful registration, duplicate ID panics, empty ID validation, deterministic ordering, alphabetical ID sorting, JSON manifest generation, and config helpers.
        - ✅ Added testify dependency and verified all tests pass (go test ./... succeeds).
        - ✅ Removed `internal/placeholder` package now that real code exists.
  - [x] **Subtask M0.3.2 – Implement credential vault MVP**
    - _Hint:_ Leverage `crypto/aes` with envelope encryption and rotate master key via CLI.
    - _Comment:_ Provide secret abstraction that can swap to external stores later.
    - [x] **Action:** Add CRUD API and CLI command `cloudmoor config vault test` with unit coverage.
      - _Hint:_ Mock key ring during tests to avoid persisting secrets on disk.
      - _Comment:_ Emit structured audit logs on secret access.
      - _Completed:_
        - ✅ Created `internal/vault` package with `Store` interface defining `Put`, `Get`, `Delete`, `List`, and `HealthCheck` methods.
        - ✅ Defined `AuditEvent` struct and `AuditHook` function type for structured audit logging with timestamps, operation names, success flags, and metadata.
        - ✅ Implemented `KeyProvider` interface with `GetKey`, `RotateKey`, and `HealthCheck` methods to support pluggable key backends.
        - ✅ Built `aesgcmStore` implementing AES-256-GCM envelope encryption with thread-safe in-memory storage.
        - ✅ Created `InMemoryKeyProvider` for testing (generates random 32-byte keys) and `FileKeyProvider` for MVP (persists keys to disk with 0600 permissions).
        - ✅ Implemented comprehensive unit tests covering: round-trip encryption/decryption, empty key/value validation, not-found errors, list operations with sorting, health checks, key rotation, audit hook invocation, and GCM tampering detection.
        - ✅ Built Cobra-based CLI with `cloudmoor config vault test` command that performs all CRUD operations, validates results, and prints structured audit log in JSON format.
        - ✅ Added Cobra dependency; CLI successfully builds and executes all vault tests.
        - ✅ Verified all tests pass (go test ./... succeeds for connectors and vault packages).
  - [ ] **Subtask M0.3.3 – Create configuration persistence layer**
    - _Hint:_ Use `golang-migrate` for forward-only migrations and keep schema diagram in docs.
    - _Comment:_ Ensure `mounts` table stores semantic version for change detection.
    - [ ] **Action:** Implement SQLite migrations covering providers, mounts, jobs, audit logs.
      - _Hint:_ Back up DB prior to migration in integration tests.
      - _Comment:_ Validate migrations on macOS/Linux via CI.
      - _Plan:_
        - Introduce `internal/storage/sqlite` module housing migration runner (using `golang-migrate` with `embed` for migration files) and repository interfaces for providers, mounts, jobs, and audit entries.
        - Define initial schema: tables for `providers`, `mounts` (with semantic version + config JSON), `jobs`, `audit_logs`, and `settings`; include indexes for lookups and soft-delete columns for future expansion.
        - Provide bootstrap migration files `0001_initial.up.sql`/`.down.sql` plus Go helper to apply migrations at daemon startup with idempotent behavior.
        - Create config struct for database path/cache directories, defaulting to XDG-compliant locations; ensure tests use temp directories.
        - Add table-driven unit tests using `sqlite` in-memory DB to verify CRUD operations and versioned YAML export/import routines (Round-trip tests).
    - [ ] **Action:** Implement YAML export/import commands with round-trip tests.
      - _Hint:_ Use JSON schema validation before applying imports.
      - _Comment:_ Include version compatibility warnings in CLI output.
      - _Plan:_
        - Design serialization structs mirroring DB schema but decoupled from internal entities; include schema version header and metadata (timestamp, CLI version).
        - Implement exporter in `internal/configio` (or under `internal/storage`) that pulls data via repositories and writes YAML using `gopkg.in/yaml.v3` with deterministic ordering.
        - Implement importer that validates schema version, runs JSON-schema-based checks, and upserts into repositories within a transaction.
        - Create CLI commands `cloudmoor config export`/`import` with flags for output path, dry-run, and compatibility warnings; include docs referencing audit trail.
        - Write round-trip tests using golden fixtures, verifying idempotency and failure cases (unknown version, schema mismatch).

- [ ] **Task M0.4 – API Surface & Service Contracts** _(Tickets: TCK-006, TCK-109)_
  - _Hint:_ Define protobuf schemas before coding server handlers.
  - _Comment:_ Align naming with REST paths to simplify auto-generated docs.
  - [ ] **Subtask M0.4.1 – Author core protobuf definitions**
    - _Hint:_ Use `buf` for lint/build; include health check service.
    - _Comment:_ Generate gRPC + grpc-gateway stubs into `internal/api`.
    - [ ] **Action:** Wire `/api/v1/health` endpoint and add smoke test.
      - _Hint:_ Use httptest server to validate HTTP gateway wiring.
      - _Comment:_ Document curl example in API docs.
  - [ ] **Subtask M0.4.2 – Publish architecture documentation**
    - _Hint:_ Export diagrams from modeling tool into `/docs/img` and embed in `architecture.md`.
    - _Comment:_ Review docs with engineering leads before merging.
    - [ ] **Action:** Convert plan’s architecture section into `/docs/architecture.md` including component + sequence diagrams.
      - _Hint:_ Use Mermaid to keep diagrams versionable.
      - _Comment:_ Link doc from root README and tasks checklist.

- [ ] **Task M0.5 – Governance & Licensing** _(Tickets: TCK-008)_
  - _Hint:_ Consult rclone upstream guidance on attribution and patch distribution.
  - _Comment:_ Store findings under `/docs/legal` for audit readiness.
  - [ ] **Subtask M0.5.1 – Complete rclone license review**
    - _Hint:_ Capture summary + decision matrix (bundle vs dependency) in markdown.
    - _Comment:_ Required before shipping binaries using rclone code.
    - [ ] **Action:** Publish `docs/legal/rclone.md` with obligations and compliance steps.
      - _Hint:_ Include release checklist updates referencing license requirements.
      - _Comment:_ Share decision log entry per governance section 21.

- [ ] **Task M0.6 – Rclone integration spike** _(Tickets: TCK-009)_
  - _Hint:_ Prototype direct librclone embedding to uncover CGO, packaging, and licensing risks.
  - _Comment:_ Required before finalizing daemon mount strategy (plan §11) and unblocking TCK-105.
  - [ ] **Subtask M0.6.1 – Evaluate librclone embedding**
    - _Hint:_ Build minimal Go wrapper linking against librclone to mount a local filesystem in CI and on macOS.
    - _Comment:_ Compare process-supervision trade-offs versus spawning the rclone binary.
    - [ ] **Action:** Publish spike summary with recommendation in `/docs/decisions/` and link from plan §11.
      - _Hint:_ Include dependency matrix, build flags, and debugging steps for each OS.
      - _Comment:_ Present findings during M0 exit review to set expectation for M1 implementation work.

## Milestone M1 – Core Providers & Daemon (6 weeks)

- [ ] **Task M1.1 – Connector Implementations (Launch set)** _(Tickets: TCK-102, TCK-104)_
  - _Hint:_ Share validation helpers and CLI wizard flows between S3/MinIO and WebDAV to minimize duplicate work.
  - _Comment:_ Merge each connector behind feature flags; FTP/SFTP (TCK-101) and Backblaze B2 (TCK-103) moved to deferred backlog.
  - [ ] **Subtask M1.1.1 – Amazon S3/MinIO connector (TCK-102)**
    - _Hint:_ Abstract credential chain to support env vars, IAM roles, and static keys.
    - _Comment:_ Validate with MinIO integration test part of CI nightly run.
    - [ ] **Action:** Implement multipart upload tuning and default cache policies.
      - _Hint:_ Expose cache configuration via provider metadata for Web UI.
      - _Comment:_ Track throughput metrics to inform benchmarks later.
  - [ ] **Subtask M1.1.2 – WebDAV connector (TCK-104)**
    - _Hint:_ Provide presets for Nextcloud/SharePoint to reduce user error.
    - _Comment:_ Document certificate handling options clearly.
    - [ ] **Action:** Validate connector against local WebDAV container with TLS toggle tests.
      - _Hint:_ Capture failing cases (self-signed certs) in regression tests.
      - _Comment:_ Surface warnings in CLI when TLS verification disabled.

- [ ] **Task M1.2 – Mount Manager & Runtime** _(Tickets: TCK-105)_
  - _Hint:_ Wrap rclone VFS process management in Go to control lifecycle and telemetry.
  - _Comment:_ Ensure concurrency-safe status tracking for multi-mount deployments.
  - [ ] **Subtask M1.2.1 – Build mount manager abstraction**
    - _Hint:_ Design state machine covering Initializing → Mounted → Reconciling → Unmounted.
    - _Comment:_ Provide hooks for logging/metrics instrumentation.
    - [ ] **Action:** Implement lifecycle API with backoff and jitter on retries.
      - _Hint:_ Reuse context cancellation to support graceful shutdowns.
      - _Comment:_ Add unit tests simulating transient failures.
  - [ ] **Subtask M1.2.2 – Verify cross-platform support**
    - _Hint:_ Add CI jobs for Linux/macos; schedule manual Windows validation using WinFsp.
    - _Comment:_ Track platform-specific quirks in `/docs/operations/known-issues.md`.
    - [ ] **Action:** Execute manual Windows smoke test checklist.
      - _Hint:_ Capture screenshots/logs for documentation.
      - _Comment:_ Log outcomes in decision register if deviations found.

- [ ] **Task M1.3 – Daemon Scheduler & Health Monitoring** _(Tickets: TCK-106)_
  - _Hint:_ Build background workers using Go contexts and ticker-based loops.
  - _Comment:_ Provide instrumentation for reconnect attempts and failure counts.
  - [ ] **Subtask M1.3.1 – Implement health workers**
    - _Hint:_ Use Prometheus client to expose mount status metrics.
    - _Comment:_ Add structured logging with correlation IDs for incidents.
    - [ ] **Action:** Detect disconnects, trigger auto-reconnect, emit alerts to event bus.
      - _Hint:_ Simulate network flaky scenarios in integration tests.
      - _Comment:_ Document default retry intervals in docs.
  - [ ] **Subtask M1.3.2 – Create alerting hooks**
    - _Hint:_ Stub webhook/email channels now; implement full notifications later (Task X.4).
    - _Comment:_ Ensure alert payload includes mount ID and recent error log.
    - [ ] **Action:** Log structured alert events and provide CLI command to list recent alerts.
      - _Hint:_ Store alerts in `audit_logs` table for traceability.
      - _Comment:_ Confirm with Product what severity levels to expose.

- [ ] **Task M1.4 – CLI Workflow Enablement** _(Tickets: TCK-107)_
  - _Hint:_ Build Cobra commands modularly to share REST/gRPC clients.
  - _Comment:_ Provide JSON output for automation from day one.
  - [ ] **Subtask M1.4.1 – Build interactive CLI commands**
    - _Hint:_ Use promptui/survey for interactive flows with fallback to flags for scripting.
    - _Comment:_ Ensure commands degrade gracefully when daemon unreachable.
    - [ ] **Action:** Implement `config create`, `mount add`, `status`, and `logs` commands.
      - _Hint:_ Add `--json` and `--watch` options consistent with plan.
      - _Comment:_ Document CLI examples in `/docs/cli.md`.
  - [ ] **Subtask M1.4.2 – Add automated CLI tests**
    - _Hint:_ Use golden files for CLI outputs to simplify assertions.
    - _Comment:_ Run tests in CI to prevent regressions before release.
    - [ ] **Action:** Add unit tests for command parsing plus integration test for mount lifecycle.
      - _Hint:_ Mock daemon responses via gRPC test server.
      - _Comment:_ Tie tests to acceptance criteria from TCK-107.

- [ ] **Task M1.5 – Integration Test Harness** _(Tickets: TCK-108)_
  - _Hint:_ Centralize Testcontainers setup to reuse across connector suites.
  - _Comment:_ Schedule nightly job separate from PR gating if runtime heavy.
  - [ ] **Subtask M1.5.1 – Stand up regression harness**
    - _Hint:_ Provide make targets (`make test-integration`) to launch containers locally.
    - _Comment:_ Document required Docker resources to avoid CI failures.
  - [ ] **Action:** Automate container lifecycle and add sample tests for each launch connector (S3/MinIO, WebDAV).
    - _Hint:_ Employ context timeouts to prevent hanging tests.
    - _Comment:_ Publish how-to guide in `/docs/testing/integration.md`.

## Milestone M2 – OAuth Providers & Web UI (4 weeks)

- [ ] **Task M2.1 – OAuth Device Flow Platform** _(Tickets: TCK-201)_
  - _Hint:_ Design token persistence to reuse vault primitives from M0.
  - _Comment:_ Align UX copy between CLI and Web UI flows for consistency.
  - [ ] **Subtask M2.1.1 – Shared OAuth service implementation**
    - _Hint:_ Embed device code polling loop with context cancellation.
    - _Comment:_ Add metrics for OAuth success/failure counts.
    - [ ] **Action:** Implement CLI `cloudmoor login` flow with automated expiry tests.
      - _Hint:_ Leverage httptest server to mock provider responses.
      - _Comment:_ Document provider-specific scope differences.
  - [ ] **Subtask M2.1.2 – Web UI + API integration**
    - _Hint:_ Use WebSocket or SSE to notify UI when device auth completes.
    - _Comment:_ Audit log should capture device code issuance and completion.
    - [ ] **Action:** Expose `/api/v1/oauth/device/start` and `/complete` endpoints consumed by UI wizard.
      - _Hint:_ Secure endpoints via short-lived tokens to prevent reuse.
      - _Comment:_ Sync API reference with OpenAPI spec.

- [ ] **Task M2.2 – Dropbox Connector & OAuth Hardening** _(Tickets: TCK-202)_
  - _Hint:_ Reuse the shared OAuth service to minimize bespoke Dropbox logic; lean on caching helpers from M1.
  - _Comment:_ Google Drive, OneDrive, Box, and pCloud connectors are deferred to the backlog.
  - [ ] **Subtask M2.2.1 – Dropbox connector (TCK-202)**
    - _Hint:_ Use incremental sync endpoints for efficiency.
    - _Comment:_ Capture refresh token storage process in provider docs.
    - [ ] **Action:** Complete OAuth flow in CLI & Web UI with retry/backoff.
      - _Hint:_ Mock Dropbox API via `httptest` to avoid flakiness.
      - _Comment:_ Add acceptance tests ensuring metadata caching works.

- [ ] **Task M2.3 – Web UI Experience (Phase 1)** _(Tickets: TCK-207, TCK-208)_
  - _Hint:_ Treat Web UI as SPA served by daemon; support local dev via Vite proxy.
  - _Comment:_ Accessibility from start reduces rework during M3 polish.
  - [ ] **Subtask M2.3.1 – Scaffold React + Vite app**
    - _Hint:_ Use TypeScript + Tailwind; configure absolute imports.
    - _Comment:_ Add eslint/prettier config matching frontend standards.
    - [ ] **Action:** Implement authentication guard using API tokens/OAuth session.
      - _Hint:_ Store tokens securely (httpOnly cookies or encrypted storage).
      - _Comment:_ Document dev proxy configuration.
  - [ ] **Subtask M2.3.2 – Build configuration wizard**
    - _Hint:_ Render forms dynamically from connector schemas.
    - _Comment:_ Ensure wizard persists progress in case of browser refresh.
    - [ ] **Action:** Handle OAuth initiation + completion with inline status feedback.
      - _Hint:_ Use toasts/snackbars for async states; fallback for CLI flows.
      - _Comment:_ Write Cypress tests covering wizard steps.
  - [ ] **Subtask M2.3.3 – Dashboard & telemetry views**
    - _Hint:_ Use charts sparingly; focus on clarity for mount status.
    - _Comment:_ Pull metrics from `/metrics` endpoint via backend proxy to avoid CORS issues.
    - [ ] **Action:** Implement log streaming viewer via WebSocket.
      - _Hint:_ Add reconnect logic for websocket drops.
      - _Comment:_ Provide download/export option for support.

- [ ] **Task M2.4 – Observability Enhancements** _(Tickets: TCK-209)_
  - _Hint:_ Standardize logging schema before propagating to connectors.
  - _Comment:_ Provide sample Grafana dashboard JSON alongside docs.
  - [ ] **Subtask M2.4.1 – Structured logging improvements**
    - _Hint:_ Adopt Zap with fields for mount ID, request ID, connector.
    - _Comment:_ Align log levels with incident response playbooks.
    - [ ] **Action:** Add correlation IDs to CLI and daemon interactions.
      - _Hint:_ Use context propagation to avoid manual wiring.
      - _Comment:_ Update runbooks with log query examples.
  - [ ] **Subtask M2.4.2 – Metrics instrumentation**
    - _Hint:_ Expose counters/gauges for latency, throughput, cache hits.
    - _Comment:_ Validate format with Prometheus lint tool.
    - [ ] **Action:** Publish sample Grafana dashboard in `/docs/observability/dashboard.json`.
      - _Hint:_ Include panels for error rate, reconnect attempts, cache utilization.
      - _Comment:_ Keep JSON small; document import steps.

- [ ] **Task M2.5 – Documentation Expansion** _(Tickets: TCK-210)_
  - _Hint:_ Capture screenshots during connector integration to avoid rework later.
  - _Comment:_ Ensure docs reference security guidance for OAuth credentials.
  - [ ] **Subtask M2.5.1 – Provider setup guides**
    - _Hint:_ Maintain consistent structure (prerequisites, steps, troubleshooting).
    - _Comment:_ Pair review with Product/UX for clarity and localization readiness.
  - [ ] **Action:** Publish guides for Amazon S3/MinIO, WebDAV, and Dropbox under `/docs/providers/`.
    - _Hint:_ Use relative image paths to keep repo portable.
    - _Comment:_ Link each guide from Web UI tooltips/help menu.

## Milestone M3 – Advanced Providers & UX (4 weeks)

- [ ] **Task M3.1 – Mega Connector & Security Enhancements** _(Tickets: TCK-301)_
  - _Hint:_ Reference Mega SDK docs for cryptography specifics.
  - _Comment:_ Highlight throttling rules to avoid account suspension.
  - [ ] **Subtask M3.1.1 – Implement Mega API integration**
    - _Hint:_ Derive crypto keys locally and persist hashed salts only.
    - _Comment:_ Add stress tests for API rate limiting.
    - [ ] **Action:** Validate encryption/decryption via regression tests with fixture data.
      - _Hint:_ Keep fixture data small to stay within repo limits.
      - _Comment:_ Document key recovery steps in provider guide.
  - [ ] **Subtask M3.1.2 – Document encryption handling**
    - _Hint:_ Use tables/flowcharts to explain offline recovery.
    - _Comment:_ Include security review sign-off in decision log.
    - [ ] **Action:** Publish `/docs/providers/mega-security.md` outlining key lifecycle.
      - _Hint:_ Add FAQ for lost password scenarios.
      - _Comment:_ Link doc to runbooks and onboarding guides.

- [ ] **Task M3.2 – Advanced Caching & Offline Mode** _(Tickets: TCK-302)_
  - _Hint:_ Build policy engine with clear states to avoid cache corruption.
  - _Comment:_ Provide CLI and UI controls for manual purge.
  - [ ] **Subtask M3.2.1 – Extend cache policy engine**
    - _Hint:_ Store cache configuration per mount in DB with TTL + size columns.
    - _Comment:_ Monitor disk usage metrics for capacity planning.
    - [ ] **Action:** Implement offline queue replay with durable persistence.
      - _Hint:_ Use WAL or job queue to guarantee delivery on reconnect.
      - _Comment:_ Write chaos test simulating prolonged offline period.
  - [ ] **Subtask M3.2.2 – Administrative tooling**
    - _Hint:_ Provide CLI + Web UI entry points for cache controls.
    - _Comment:_ Audit log all purge operations for compliance.
    - [ ] **Action:** Ship `cloudmoor cache purge` command with success/failure telemetry.
      - _Hint:_ Provide dry-run mode for change previews.
      - _Comment:_ Document command in CLI manual.
    - [ ] **Action:** Surface cache settings in Web UI with real-time metrics.
      - _Hint:_ Use gauge/slider UI for size limits.
      - _Comment:_ Ensure forms validate lower/upper bounds.

- [ ] **Task M3.3 – RBAC & Multi-User Support** _(Tickets: TCK-303)_
  - _Hint:_ Implement token issuance with refresh + revoke capabilities.
  - _Comment:_ Enforce least privilege roles across API endpoints.
  - [ ] **Subtask M3.3.1 – Daemon-side RBAC**
    - _Hint:_ Store users/roles in DB with passwordless auth ready for SSO.
    - _Comment:_ Add middleware enforcing role-based permissions.
    - [ ] **Action:** Write API tests verifying Admin/Operator/Read-only matrices.
      - _Hint:_ Use table-driven tests for coverage.
      - _Comment:_ Document results in security review artifacts.
  - [ ] **Subtask M3.3.2 – Web UI management experience**
    - _Hint:_ Provide role assignment UI with confirmation modals.
    - _Comment:_ Reflect audit logs inline for transparency.
    - [ ] **Action:** Build user management screen with search/filter + invite flow.
      - _Hint:_ Reuse table components from dashboard.
      - _Comment:_ Add Cypress tests covering role changes.

- [ ] **Task M3.4 – UX Polish & Accessibility** _(Tickets: TCK-304, TCK-502)_
  - _Hint:_ Conduct accessibility audit early; fix keyboard navigation issues.
  - _Comment:_ Prepare localization scaffold for future languages.
  - [ ] **Subtask M3.4.1 – Design system and theming**
    - _Hint:_ Create Tailwind theme tokens for light/dark modes.
    - _Comment:_ Hit Lighthouse accessibility score ≥90.
    - [ ] **Action:** Implement responsive layouts and color contrast adjustments.
      - _Hint:_ Test on mobile breakpoints using BrowserStack or similar.
      - _Comment:_ Capture before/after visuals in design docs.
  - [ ] **Subtask M3.4.2 – Localization scaffold**
    - _Hint:_ Externalize strings using ICU message format.
    - _Comment:_ Provide English + Spanish translations as baseline.
    - [ ] **Action:** Add locale switcher prototype and CLI `--lang` flag support.
      - _Hint:_ Store language preference in persisted settings.
      - _Comment:_ Plan QA pass with bilingual reviewer.

- [ ] **Task M3.5 – Desktop Tray Prototype** _(Tickets: TCK-305)_
  - _Hint:_ Start with Tauri for shared Rust backend; fall back to Electron if blockers arise.
  - _Comment:_ Keep footprint small; prototype only.
  - [ ] **Subtask M3.5.1 – Build tray application**
    - _Hint:_ Expose minimal UI: mount list, start/stop buttons, status indicator.
    - _Comment:_ Use daemon REST API over localhost with auth token secured.
    - [ ] **Action:** Produce macOS & Windows builds with auto-update disabled (prototype stage).
      - _Hint:_ Document build steps in `/docs/desktop/README.md`.
      - _Comment:_ Collect OS-specific issues in backlog.
  - [ ] **Subtask M3.5.2 – Feedback loop**
    - _Hint:_ Onboard internal testers first; gather notes via issue template.
    - _Comment:_ Convert validated feedback into backlog items or tasks.
    - [ ] **Action:** Summarize findings in decision log entry referencing TCK-305.
      - _Hint:_ Use template from `/docs/decisions/`.
      - _Comment:_ Present summary during M3 exit review.

- [ ] **Task M3.6 – Operational Runbooks** _(Tickets: TCK-306)_
  - _Hint:_ Collaborate with DevOps to capture real-world remediation steps.
  - _Comment:_ Ensure runbooks align with alert payloads from M1.3.
  - [ ] **Subtask M3.6.1 – Author operations documentation**
    - _Hint:_ Structure runbooks as actionable checklists with escalation paths.
    - _Comment:_ Include RACI table for incident response.
    - [ ] **Action:** Publish mount failure, credential rotation, and DR playbooks under `/docs/operations/`.
      - _Hint:_ Add diagrams/flowcharts for complex procedures.
      - _Comment:_ Review with support leads prior to release.

## Milestone M4 – Hardening & Release (3 weeks)

- [ ] **Task M4.1 – Security & Compliance** _(Tickets: TCK-401)_
  - _Hint:_ Integrate security scanning into CI early in the milestone.
  - _Comment:_ Capture findings in security backlog with owner + due date.
  - [ ] **Subtask M4.1.1 – Dependency and security audit**
    - _Hint:_ Use Syft/Grype or Trivy for SBOM + vulnerability reports.
    - _Comment:_ Attach reports to release candidate artifacts.
    - [ ] **Action:** Automate SBOM generation in release workflow and archive outputs.
      - _Hint:_ Store SBOM alongside binaries in GitHub Releases.
      - _Comment:_ Update compliance checklist with storage location.
  - [ ] **Subtask M4.1.2 – Penetration test readiness**
    - _Hint:_ Build internal checklist covering auth, encryption, logging controls.
    - _Comment:_ Share sanitized report with beta customers if requested.
    - [ ] **Action:** Produce security review report and log in `/docs/security/reports/`.
      - _Hint:_ Tie recommendations to follow-up tasks.
      - _Comment:_ Secure approvals from security advisor.

- [ ] **Task M4.2 – Packaging & Distribution** _(Tickets: TCK-402)_
  - _Hint:_ Use Goreleaser to generate installers/binaries consistently.
  - _Comment:_ Test on clean VMs to ensure dependencies bundled.
  - [ ] **Subtask M4.2.1 – Build installers per platform**
    - _Hint:_ Automate Homebrew tap, Debian/RPM packaging, and Windows MSI/NSIS.
    - _Comment:_ Document uninstall steps for each platform.
    - [ ] **Action:** Validate `cloudmoor service install` post-install on macOS, Linux, Windows.
      - _Hint:_ Script smoke tests that mount dummy storage post-install.
      - _Comment:_ Capture logs/screenshots for release notes.
  - [ ] **Subtask M4.2.2 – Container & Helm delivery**
    - _Hint:_ Build multi-arch images and publish to GHCR.
    - _Comment:_ Provide Helm chart values for TLS + persistent storage.
    - [ ] **Action:** Publish Helm chart and automation templates (Terraform, Ansible).
      - _Hint:_ Validate chart via `helm lint` and kind cluster deploy.
      - _Comment:_ Keep versioning in sync with main release.

- [ ] **Task M4.3 – Documentation & Onboarding** _(Tickets: TCK-403, TCK-404)_
  - _Hint:_ Drive docs updates in parallel with packaging to avoid bottleneck.
  - _Comment:_ Include telemetry opt-in flow per privacy policy.
  - [ ] **Subtask M4.3.1 – Getting started experience**
    - _Hint:_ Produce quickstart guide with CLI + Web UI paths.
    - _Comment:_ Add troubleshooting FAQ covering top 5 install issues.
    - [ ] **Action:** Publish `/docs/getting-started.md`, FAQ, and walkthrough assets.
      - _Hint:_ Record short Loom/videos for onboarding.
      - _Comment:_ Embed links in README and Web UI help drawer.
  - [ ] **Subtask M4.3.2 – Beta program launch**
    - _Hint:_ Use static landing page + form (or GitHub Discussion) for beta intake.
    - _Comment:_ Ensure telemetry opt-in respects regional laws.
    - [ ] **Action:** Onboard pilot testers, set up feedback triage cadence, and log learnings.
      - _Hint:_ Track feedback in shared spreadsheet or issue label.
      - _Comment:_ Review metrics during weekly steering sync as per plan section 21.

- [ ] **Task M4.4 – Performance Benchmarking** _(Tickets: TCK-405)_
  - _Hint:_ Run benchmarks on consistent hardware/VM sizes for reproducibility.
  - _Comment:_ Store baseline numbers to detect regressions in future releases.
  - [ ] **Subtask M4.4.1 – Execute benchmark suite**
    - _Hint:_ Include sequential + parallel read/write scenarios per provider.
    - _Comment:_ Publish results and regression thresholds in docs.
    - [ ] **Action:** Produce `/docs/performance.md` with charts/tables summarizing results.
      - _Hint:_ Use lightweight charting library (e.g., Mermaid or embedded images).
      - _Comment:_ Record benchmark configuration in appendix for repeatability.

## Deferred Connector Tasks (Post-v1)

- [ ] **Task D.1 – FTP/SFTP connector (TCK-101)** _(Deferred)_
  - _Hint:_ Base implementation on rclone SFTP backend, ensuring passive mode toggle.
  - _Comment:_ Target after launch once vault and registry patterns are validated in production.
  - [ ] **Action:** Support username/password and SSH key auth with vsftpd container tests.
    - _Hint:_ Use Testcontainers to spin up vsftpd with seeded directories.
    - _Comment:_ Document known SSH cipher constraints in provider guide.
  - [ ] **Action:** Wire CLI wizard to capture credentials securely.
    - _Hint:_ Offer JSON output for automation; reuse validation prompts across providers.
    - _Comment:_ Add unit tests for wizard flows in headless mode.

- [ ] **Task D.2 – Backblaze B2 connector (TCK-103)** _(Deferred)_
  - _Hint:_ Mock large file uploads locally if sandbox unavailable; ensure resume logic.
  - _Comment:_ Schedule once sustained bandwidth benchmarks complete for S3/MinIO.
  - [ ] **Action:** Implement exponential backoff and resume for >2GB uploads.
    - _Hint:_ Store upload state in cache DB to survive restarts.
    - _Comment:_ Add metrics for retry counts.

- [ ] **Task D.3 – Google Drive connector (TCK-203)** _(Deferred)_
  - _Hint:_ Implement service-account impersonation path for enterprises.
  - _Comment:_ Reassess priority based on customer demand after beta feedback.
  - [ ] **Action:** Validate shared drive support with integration tests.
    - _Hint:_ Use Google-provided test data set where possible.
    - _Comment:_ Document prerequisites (client secrets JSON) in provider guide.

- [ ] **Task D.4 – OneDrive connector (TCK-204)** _(Deferred)_
  - _Hint:_ Support both personal and business endpoints via config flag.
  - _Comment:_ Track Graph API throttling learnings from Dropbox rollout before revisiting.
  - [ ] **Action:** Add tests covering tenant selection and file operations.
    - _Hint:_ Provide sample tenant config JSON for QA team.
    - _Comment:_ Capture known limitations (e.g., SharePoint site support) in docs.

- [ ] **Task D.5 – Box connector (TCK-205)** _(Deferred)_
  - _Hint:_ Use JWT app config for enterprise integration; secure private key storage.
  - _Comment:_ Revisit alongside enterprise feedback on governance requirements.
  - [ ] **Action:** Implement file operations with policy-aware error handling.
    - _Hint:_ Surface admin consent instructions in UI wizard.
    - _Comment:_ Add audit logging for Box-specific retention policies.

- [ ] **Task D.6 – pCloud connector (TCK-206)** _(Deferred)_
  - _Hint:_ Provide folder picker UI to limit scope.
  - _Comment:_ Evaluate after Dropbox customer onboarding to understand OAuth coverage.
  - [ ] **Action:** Validate upload/download flows and metadata listing.
    - _Hint:_ Include offline cache scenario in tests.
    - _Comment:_ Ensure connector handles 2FA-enabled accounts gracefully.

## Cross-Cutting Backlog (Post-M4 or As Capacity Allows)

- [ ] **Task X.1 – External Secret Store Integrations** _(Tickets: TCK-501)_
  - _Hint:_ Evaluate HashiCorp Vault + AWS Secrets Manager first; prioritize based on user demand.
  - _Comment:_ Target follow-up RFC with phased rollout plan.
  - [ ] **Subtask X.1.1 – Research & spike**
    - _Hint:_ Prototype minimal Vault integration for credential fetch.
    - _Comment:_ Document API changes required for pluggable backends.
    - [ ] **Action:** Publish feasibility report with recommended next steps.
      - _Hint:_ Include security considerations (auth, rotation).
      - _Comment:_ Tag stakeholders in decision log entry.

- [ ] **Task X.2 – Enhanced Localization** _(Tickets: TCK-502 continuation)_
  - _Hint:_ Coordinate with translation vendor/community once scaffold ready.
  - _Comment:_ Align release timeline with marketing commitments.
  - [ ] **Subtask X.2.1 – Expand translation coverage**
    - _Hint:_ Prioritize languages based on customer pipeline.
    - _Comment:_ Update QA checklist with locale-specific tests.
    - [ ] **Action:** Add additional language packs and CLI localization support.
      - _Hint:_ Reuse ICU message files; ensure fallback to English.
      - _Comment:_ Validate formatting (dates/numbers) per locale.

- [ ] **Task X.3 – Scheduled Sync Jobs** _(Tickets: TCK-503)_
  - _Hint:_ Extend daemon scheduler with cron parser (robfig/cron v3).
  - _Comment:_ Provide conflict resolution policies aligning with user preferences.
  - [ ] **Subtask X.3.1 – Implement cron-based sync engine**
    - _Hint:_ Persist schedules with next run time to avoid drift.
    - _Comment:_ Add metrics for job success/failure counts.
    - [ ] **Action:** Write tests for overlapping schedules and manual triggers.
      - _Hint:_ Use simulated clock to test edge cases.
      - _Comment:_ Document configuration examples in `/docs/scheduling.md`.

- [ ] **Task X.4 – Notification Channels** _(Tickets: TCK-504)_
  - _Hint:_ Reuse alert hooks from M1.3; add transport-specific senders.
  - _Comment:_ Provide per-mount configuration UI + CLI flags.
  - [ ] **Subtask X.4.1 – Webhook & email notifications**
    - _Hint:_ Store delivery attempts with exponential backoff.
    - _Comment:_ Allow templating for email subject/body.
    - [ ] **Action:** Implement webhook sender with signing secret + SMTP email integration.
      - _Hint:_ Use background workers to avoid blocking health loops.
      - _Comment:_ Update runbooks with notification troubleshooting tips.

- [ ] **Task X.5 – Web UI MFA Evaluation** _(Tickets: TCK-505)_
  - _Hint:_ Interview security-sensitive customers to gauge urgency.
  - _Comment:_ Outcome should feed roadmap prioritization.
  - [ ] **Subtask X.5.1 – Assess MFA requirements**
    - _Hint:_ Compare TOTP, WebAuthn, and external IdP delegation options.
    - _Comment:_ Capture effort estimate and dependencies.
    - [ ] **Action:** Publish recommendation memo with proposed milestone placement.
      - _Hint:_ Include risk analysis if MFA deferred.
      - _Comment:_ Attach memo to steering committee meeting notes.

---

**Usage Tips:**

- Update checkboxes during standups, sprint reviews, and milestone exit reviews.
- When a task’s scope changes, create/modify the corresponding ticket entry in this document and reflect edits per Governance Section 21.
- Log clarifications or deviations in `/docs/decisions/` to maintain traceability.

---

## Ticket Backlog (Detailed Acceptance Criteria)

### Legend

- **Type:** Feature (F), Task (T), Bug (B, placeholder), Documentation (D), Research (R)
- **Priority:** P1 (Critical), P2 (High), P3 (Medium), P4 (Low)
- **Estimate:** Story points or ideal days (initial ballpark)
- **Dependencies:** Upstream tickets that must complete first
- **Acceptance Criteria:** Conditions required for completion
- **Owner:** Suggested role

### Ticket Backlog – Milestone M0 (3 weeks)

#### TCK-001 · F · P1 · 5 pts — Bootstrap repository structure and tooling

- **Description:** Initialize Go module, set up directory layout (`cmd/`, `internal/`, `pkg/`, `web/`, `deploy/`), add base configs.
- **Dependencies:** None
- **Acceptance Criteria:**
  - `go mod init` completed with agreed module path.
  - Base folders committed with README placeholders.
  - `.editorconfig`, `.gitignore`, and license file added.

#### TCK-002 · T · P1 · 3 pts — Configure CI pipeline skeleton

- **Description:** Setup GitHub Actions workflow for lint/test on push/PR, include Go version matrix.
- **Dependencies:** TCK-001
- **Acceptance Criteria:**
  - Workflow runs gofmt/go vet/golangci-lint on PRs.
  - Unit test job executes `go test ./...`.
  - Status badges documented in README.
- **Owner:** DevOps/SRE

#### TCK-003 · F · P1 · 5 pts — Implement connector interface contracts

- **Description:** Define Go interfaces for provider connectors (init, validate, mount, teardown) and scaffold plugin registry.
- **Dependencies:** TCK-001
- **Acceptance Criteria:**
  - Interface definitions reviewed and approved.
  - Registry supports loading connector metadata.
  - Placeholder connectors compile and register without errors.
- **Owner:** Lead Go engineer

#### TCK-004 · F · P1 · 8 pts — Build credential vault MVP

- **Description:** Implement encrypted storage for provider secrets using AES-GCM, master key management, CRUD API.
- **Dependencies:** TCK-001, TCK-003
- **Acceptance Criteria:**
  - Secrets stored encrypted at rest; round-trip tested.
  - Unit tests cover encrypt/decrypt, key rotation.
  - CLI command `cloudmoor config vault test` validates setup.
- **Owner:** Cloud integrations engineer

#### TCK-005 · F · P1 · 8 pts — Create config persistence layer

- **Description:** Implement SQLite schema for providers, mounts, jobs, audit logs; include migrations.
- **Dependencies:** TCK-003
- **Acceptance Criteria:**
  - Schema generated via migrations; `go test` covers CRUD operations.
  - Database file path configurable; migration applied on startup.
  - Export/import CLI commands serialize configs to YAML.
- **Owner:** Lead Go engineer

#### TCK-006 · F · P2 · 5 pts — Skeleton for gRPC + REST gateway

- **Description:** Define protobuf for core services (MountService, ProviderService) and generate gRPC/REST stubs.
- **Dependencies:** TCK-003, TCK-005
- **Acceptance Criteria:**
  - Proto files linted; server compiles.
  - REST gateway accessible via `/api/v1/health`.
  - Example client integration test passes.
- **Owner:** Lead Go engineer

#### TCK-007 · D · P2 · 3 pts — Author contributor guidelines and coding standards

- **Description:** Document branching strategy, lint rules, testing expectations, PR checklist.
- **Dependencies:** TCK-001, TCK-002
- **Acceptance Criteria:**
  - `CONTRIBUTING.md` merged and referenced in README.
  - Templates for bugs, features, security issues added.
- **Owner:** Product/UX + Lead Go engineer

#### TCK-008 · R · P2 · 2 pts — Validate rclone library licensing and compatibility

- **Description:** Review rclone MIT license obligations, document attribution requirements.
- **Dependencies:** None
- **Acceptance Criteria:**
  - Legal note in docs with attribution plan.
  - Decision recorded on bundling approach vs extension link.
- **Owner:** Product/UX

#### TCK-009 · R · P1 · 5 pts — Evaluate librclone embedding feasibility

- **Description:** Prototype direct librclone embedding via CGO to understand build, packaging, and telemetry implications relative to shelling out to the binary.
- **Dependencies:** TCK-003
- **Acceptance Criteria:**
  - Minimal Go harness linking librclone builds and mounts a test remote on macOS and Linux.
  - Pros/cons documented comparing embedded vs subprocess execution (performance, footprint, upgrade story).
  - Decision log entry created with recommendation and next steps.
- **Owner:** Lead Go engineer

### Ticket Backlog – Milestone M1 (6 weeks)

#### TCK-102 · F · P1 · 8 pts — Implement Amazon S3/MinIO connector

- **Description:** Support IAM keys, STS assume role, custom endpoints; configure VFS cache defaults.
- **Dependencies:** TCK-004, TCK-005, TCK-006
- **Acceptance Criteria:**
  - Integration tests using MinIO container.
  - Bucket listing, upload/download validated.
  - Support for environment credential fallback.
- **Owner:** Cloud integrations engineer

#### TCK-104 · F · P1 · 8 pts — Implement WebDAV connector with presets

- **Description:** Provide generic WebDAV plus presets (Nextcloud, SharePoint); certificate options.
- **Dependencies:** TCK-004, TCK-005, TCK-006
- **Acceptance Criteria:**
  - Integration tests via local WebDAV container.
  - Presets selectable in CLI/Web UI; config saved.
  - TLS verification toggle documented.
- **Owner:** Cloud integrations engineer

#### TCK-105 · F · P1 · 13 pts — Develop mount manager & FUSE integration layer

- **Description:** Abstract platform-specific mounting via rclone VFS; manage lifecycle, retry policies.
- **Dependencies:** TCK-003, TCK-005, connectors (TCK-102, TCK-104)
- **Acceptance Criteria:**
  - Mount manager API supports create/start/stop/status.
  - Linux/macOS integration tests pass in CI.
  - Windows mount validated manually with WinFsp.
- **Owner:** Lead Go engineer

#### TCK-106 · F · P1 · 8 pts — Build daemon scheduler and health monitor

- **Description:** Background workers for mount health, retries, metrics collection, notification hooks.
- **Dependencies:** TCK-105
- **Acceptance Criteria:**
  - Health checks detect disconnect and trigger auto-reconnect.
  - Metrics exposed at `/metrics` Prometheus endpoint.
  - Alert hooks (log/Webhook stub) run on failure.
- **Owner:** Lead Go engineer

#### TCK-107 · F · P1 · 8 pts — Implement CLI workflows for config and mounts

- **Description:** Interactive commands for `config create`, `mount add`, status listing, log tailing.
- **Dependencies:** TCK-102, TCK-104, TCK-105, TCK-106
- **Acceptance Criteria:**
  - CLI UX reviewed, supports JSON output.
  - `cloudmoor mount <name>` works end-to-end.
  - Unit tests cover command parsing and API calls.
- **Owner:** Lead Go engineer

#### TCK-108 · T · P2 · 5 pts — Add integration test harness with Testcontainers

- **Description:** Create Go test suite spinning up MinIO (S3) and WebDAV containers for regression coverage.
- **Dependencies:** TCK-102, TCK-104
- **Acceptance Criteria:**
  - CI job runs integration suite nightly.
  - Tests validate mount lifecycle for S3/MinIO and WebDAV connectors.
  - Docs explain how to run locally.
- **Owner:** QA/Automation

#### TCK-109 · D · P2 · 3 pts — Draft architecture overview documentation

- **Description:** Convert sections of plan into `/docs/architecture.md`, include diagrams.
- **Dependencies:** TCK-105, TCK-106
- **Acceptance Criteria:**
  - Docs include component diagram, sequence for mount lifecycle.
  - Linked from README.
- **Owner:** Product/UX + Lead Go engineer

### Ticket Backlog – Milestone M2 (4 weeks)

#### TCK-201 · F · P1 · 8 pts — Implement OAuth device flow service

- **Description:** Build shared OAuth handler supporting device code, PKCE, token refresh persistence.
- **Dependencies:** M1 completion, TCK-004, TCK-006
- **Acceptance Criteria:**
  - Works with Dropbox sandbox; tokens stored securely.
  - CLI command `cloudmoor login` completes device flow.
  - Unit/integration tests cover expiration handling.
- **Owner:** Cloud integrations engineer

#### TCK-202 · F · P1 · 5 pts — Dropbox connector with OAuth

- **Description:** Use OAuth service to integrate Dropbox API, metadata caching.
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - OAuth flow completes via CLI & Web UI.
  - File operations succeed in integration tests.
  - Rate limit retries implemented.
- **Owner:** Cloud integrations engineer

#### TCK-207 · F · P1 · 8 pts — Build Web UI foundation (React + Vite)

- **Description:** Scaffold React app, integrate Tailwind, set up routing, connect to API gateway.
- **Dependencies:** TCK-006
- **Acceptance Criteria:**
  - Web UI builds via `npm run build` and served by daemon.
  - Auth flow with API tokens/OAuth login works.
  - Basic dashboard page renders mount list.
- **Owner:** Frontend engineer

#### TCK-208 · F · P1 · 8 pts — Implement Web UI configuration wizard

- **Description:** Multistep forms for provider setup, validation, credential entry, summary.
- **Dependencies:** TCK-207, TCK-202
- **Acceptance Criteria:**
  - Wizard supports OAuth device initiation and completion.
  - Schema-driven forms render per provider.
  - Errors surfaced inline with accessible messaging.
- **Owner:** Frontend engineer

#### TCK-209 · F · P2 · 5 pts — Metrics & logging enhancements

- **Description:** Add structured logging with Zap fields, Prometheus metrics (latency, throughput, cache), expose `/metrics`.
- **Dependencies:** TCK-106
- **Acceptance Criteria:**
  - Metrics visible in Prometheus scrape.
  - Log entries include correlation IDs.
  - Sample Grafana dashboard provided.
- **Owner:** DevOps/SRE

#### TCK-210 · D · P2 · 3 pts — Update documentation with launch provider guides

- **Description:** Write guides for Amazon S3/MinIO, WebDAV, and Dropbox setup.
- **Dependencies:** TCK-102, TCK-104, TCK-202
- **Acceptance Criteria:**
  - `/docs/providers/*.md` created with step-by-step instructions for the three launch connectors.
  - Screenshots or references for relevant console configuration (AWS, WebDAV preset, Dropbox OAuth).
- **Owner:** Product/UX

### Ticket Backlog – Milestone M3 (4 weeks)

#### TCK-301 · F · P1 · 8 pts — Mega connector with client-side encryption

- **Description:** Implement Mega API integration, handle key derivation, throttling limits.
- **Dependencies:** M2 completion, TCK-004
- **Acceptance Criteria:**
  - Integration tests using Mega SDK or mock.
  - Encryption key handling validated and documented.
  - Throttle management prevents API bans during stress test.
- **Owner:** Cloud integrations engineer

#### TCK-302 · F · P1 · 5 pts — Advanced cache policy controls

- **Description:** Implement per-mount cache size/TTL, offline mode, manual purge command.
- **Dependencies:** TCK-105, TCK-106
- **Acceptance Criteria:**
  - Cache policies configurable via CLI/Web UI.
  - Offline mode retains cached data and queues writes.
  - `cloudmoor cache purge` clears caches and logs outcome.
- **Owner:** Lead Go engineer

#### TCK-303 · F · P1 · 5 pts — RBAC & multi-user support in daemon

- **Description:** Add user management, roles (Admin, Operator, Read-only), auth tokens.
- **Dependencies:** TCK-006, TCK-207
- **Acceptance Criteria:**
  - RBAC enforced at API layer; tests cover role restrictions.
  - Web UI user management page functional.
  - Audit logs capture role changes.
- **Owner:** Lead Go engineer + Frontend engineer

#### TCK-304 · F · P2 · 5 pts — Web UI polish & accessibility pass

- **Description:** Apply design system, dark mode, keyboard navigation, localization scaffolding.
- **Dependencies:** TCK-207, TCK-208
- **Acceptance Criteria:**
  - Lighthouse accessibility score ≥90.
  - I18n strings externalized; locale switcher prototype.
  - Keyboard navigation validated in QA checklist.
- **Owner:** Frontend engineer

#### TCK-305 · F · P2 · 5 pts — Desktop tray prototype (Tauri/Electron)

- **Description:** Build lightweight desktop app showing mount status, start/stop controls.
- **Dependencies:** TCK-107, TCK-303
- **Acceptance Criteria:**
  - Prototype runs on macOS & Windows.
  - Communicates with daemon via localhost API.
  - Basic notifications on mount failure.
- **Owner:** Frontend engineer

#### TCK-306 · D · P2 · 3 pts — Publish operations runbooks

- **Description:** Document mount failure remediation, credential rotation, disaster recovery.
- **Dependencies:** TCK-302, TCK-303
- **Acceptance Criteria:**
  - Runbooks stored under `/docs/operations/`.
  - Includes flowcharts or checklists for common incidents.
- **Owner:** DevOps/SRE

### Ticket Backlog – Milestone M4 (3 weeks)

#### TCK-401 · T · P1 · 5 pts — Security review & dependency audit

- **Description:** Perform dependency scan, SBOM generation, penetration test checklist.
- **Dependencies:** Previous milestones
- **Acceptance Criteria:**
  - SBOM uploaded with release assets.
  - Critical findings resolved or tracked.
  - Security review report stored internally.
- **Owner:** DevOps/SRE + Lead Go engineer

#### TCK-402 · F · P1 · 5 pts — Packaging & installer tooling

- **Description:** Create Homebrew tap, Debian/RPM packages, Windows installer (MSI) via WiX or NSIS.
- **Dependencies:** TCK-107, TCK-207
- **Acceptance Criteria:**
  - Installers tested on target platforms.
  - `cloudmoor service install` works post-install.
  - Docs updated with install instructions.
- **Owner:** DevOps/SRE

#### TCK-403 · F · P1 · 5 pts — User documentation & onboarding guides

- **Description:** Produce getting started guide, provider-specific tutorials, FAQ.
- **Dependencies:** TCK-109, TCK-210, TCK-306
- **Acceptance Criteria:**
  - `/docs/getting-started.md` + tutorial videos/screenshots.
  - FAQ covers top support scenarios.
  - Documentation reviewed by Product/UX.
- **Owner:** Product/UX

#### TCK-404 · T · P2 · 3 pts — Beta release program launch

- **Description:** Set up feedback channels, telemetry opt-in, announce beta to pilot users.
- **Dependencies:** TCK-401..TCK-403
- **Acceptance Criteria:**
  - Landing page for beta signups live.
  - Feedback triage process documented.
  - First batch of pilot testers onboarded.
- **Owner:** Product/UX

#### TCK-405 · T · P2 · 3 pts — Performance benchmarking report

- **Description:** Execute benchmark suite, document throughput, latency, resource usage across providers.
- **Dependencies:** TCK-302, TCK-209
- **Acceptance Criteria:**
  - Report published in `/docs/performance.md`.
  - Regression thresholds defined for future releases.
  - Findings reviewed with engineering team.
- **Owner:** QA/Automation + DevOps/SRE

### Deferred Backlog – Post-v1 Connectors

#### TCK-101 · F · P2 · 13 pts — Implement FTP/SFTP connector _(Deferred)_

- **Description:** Build connector supporting username/password and SSH key auth, passive mode, validation tests. Deferred until after launch to focus on the narrower beta scope.
- **Status:** Deferred to post-v1 (plan §7.2).
- **Dependencies:** TCK-004, TCK-005
- **Acceptance Criteria:**
  - Connector passes integration tests against vsftpd container.
  - CLI wizard prompts for credentials, stores in vault.
  - Mount/unmount operations succeed with smoke test.
- **Owner:** Cloud integrations engineer

#### TCK-103 · F · P2 · 8 pts — Implement Backblaze B2 connector _(Deferred)_

- **Description:** Integrate application keys, large file upload resume, throughput tuning. Deferred until post-v1 to keep launch footprint manageable.
- **Status:** Deferred to post-v1 (plan §7.2).
- **Dependencies:** TCK-004, TCK-005, TCK-006
- **Acceptance Criteria:**
  - Integration tests using B2 sandbox or mock.
  - Large file (>2GB) upload resume verified.
  - Error retries with exponential backoff implemented.
- **Owner:** Cloud integrations engineer

#### TCK-203 · F · P2 · 5 pts — Google Drive connector _(Deferred)_

- **Description:** Support OAuth client and service account auth, handle shared drives. Deferred until beta feedback confirms demand.
- **Status:** Deferred to post-v1 (plan §7.2).
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - Integration tests cover both auth modes.
  - Change detection sync verified.
  - Scopes documented and minimal.
- **Owner:** Cloud integrations engineer

#### TCK-204 · F · P2 · 5 pts — OneDrive connector _(Deferred)_

- **Description:** Implement OneDrive Graph API integration with tenant selection. Deferred until after Dropbox learnings are incorporated.
- **Status:** Deferred to post-v1 (plan §7.2).
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - Device flow works for personal Microsoft account.
  - Business tenant selection documented.
  - Integration tests validate upload/download.
- **Owner:** Cloud integrations engineer

#### TCK-205 · F · P2 · 5 pts — Box connector _(Deferred)_

- **Description:** Support enterprise JWT apps and standard OAuth; manage key rotation. Deferred to avoid overextending beta milestone.
- **Status:** Deferred to post-v1 (plan §7.2).
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - JWT auth path validated in integration test (mock or sandbox).
  - Token refresh + rotation automated.
  - Error handling for enterprise policies documented.
- **Owner:** Cloud integrations engineer

#### TCK-206 · F · P2 · 5 pts — pCloud connector _(Deferred)_

- **Description:** Integrate pCloud API with folder selection and cache policies. Deferred until post-v1 once OAuth platform is hardened.
- **Status:** Deferred to post-v1 (plan §7.2).
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - OAuth flow tested with sandbox account.
  - Upload/download and metadata listing confirmed.
  - Documentation for API key setup.
- **Owner:** Cloud integrations engineer

### Cross-Cutting & Backlog Items

#### TCK-501 · R · P3 · 3 pts — Investigate external secret manager integrations

- **Description:** Research architecture for HashiCorp Vault, AWS Secrets Manager extensions.
- **Dependencies:** TCK-004
- **Acceptance Criteria:**
  - Proposal document with recommended approach.
  - Spike code proving feasibility.
- **Owner:** Cloud integrations engineer

#### TCK-502 · F · P3 · 5 pts — Implement localization framework

- **Description:** Add ICU message format support, locale detection, translation files.
- **Dependencies:** TCK-304
- **Acceptance Criteria:**
  - Web UI builds with language packs.
  - CLI supports `--lang` flag for prompts.
  - English + Spanish translations verified.
- **Owner:** Frontend engineer

#### TCK-503 · F · P3 · 5 pts — Add scheduled sync jobs with cron expressions

- **Description:** Allow mounts to run sync tasks on schedule with conflict resolution options.
- **Dependencies:** TCK-302, TCK-106
- **Acceptance Criteria:**
  - Cron parser integrated; jobs persisted in DB.
  - Sync conflicts resolved via policy settings.
  - Tests cover overlapping schedules.
- **Owner:** Lead Go engineer

#### TCK-504 · F · P3 · 5 pts — Implement webhook & email notifications

- **Description:** Notify users on mount failures, token expiry via configurable channels.
- **Dependencies:** TCK-209, TCK-303
- **Acceptance Criteria:**
  - Webhooks configurable per mount.
  - Email notifications sent via SMTP integration (mock in dev).
  - Audited in logs.
- **Owner:** DevOps/SRE

#### TCK-505 · R · P3 · 3 pts — Evaluate MFA requirements for Web UI

- **Description:** Gather stakeholder needs, review implementation options, scope effort.
- **Dependencies:** TCK-303
- **Acceptance Criteria:**
  - Recommendation doc with effort estimate.
  - Decision captured in roadmap backlog.
- **Owner:** Product/UX

### Ticket Notes

- Estimates are preliminary and should be refined during sprint planning.
- Tickets marked P1 are blockers for milestone completion.
- Backlog items can be pulled into milestones after foundational work stabilizes.

<!-- markdownlint-enable MD007 -->
