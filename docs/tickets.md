# CloudMoor Ticket Backlog

Date: 2025-10-06

## Legend

- **Type:** Feature (F), Task (T), Bug (B, placeholder), Documentation (D), Research (R)
- **Priority:** P1 (Critical), P2 (High), P3 (Medium), P4 (Low)
- **Estimate:** Story points or ideal days (initial ballpark)
- **Dependencies:** Upstream tickets that must complete first
- **Acceptance Criteria:** Conditions required for completion
- **Owner:** Suggested role

---

## Milestone M0 – Foundations (2 weeks)

### TCK-001 · F · P1 · 5 pts

**Title:** Bootstrap repository structure and tooling

- **Description:** Initialize Go module, set up directory layout (`cmd/`, `internal/`, `pkg/`, `web/`, `deploy/`), add base configs.
- **Dependencies:** None
- **Acceptance Criteria:**
  - `go mod init` completed with agreed module path.
  - Base folders committed with README placeholders.
  - `.editorconfig`, `.gitignore`, and license file added.
- **Owner:** Lead Go engineer

### TCK-002 · T · P1 · 3 pts

**Title:** Configure CI pipeline skeleton

- **Description:** Setup GitHub Actions workflow for lint/test on push/PR, include Go version matrix.
- **Dependencies:** TCK-001
- **Acceptance Criteria:**
  - Workflow runs gofmt/go vet/golangci-lint on PRs.
  - Unit test job executes `go test ./...`.
  - Status badges documented in README.
- **Owner:** DevOps/SRE

### TCK-003 · F · P1 · 5 pts

**Title:** Implement connector interface contracts

- **Description:** Define Go interfaces for provider connectors (init, validate, mount, teardown) and scaffold plugin registry.
- **Dependencies:** TCK-001
- **Acceptance Criteria:**
  - Interface definitions reviewed and approved.
  - Registry supports loading connector metadata.
  - Placeholder connectors compile and register without errors.
- **Owner:** Lead Go engineer

### TCK-004 · F · P1 · 8 pts

**Title:** Build credential vault MVP

- **Description:** Implement encrypted storage for provider secrets using AES-GCM, master key management, CRUD API.
- **Dependencies:** TCK-001, TCK-003
- **Acceptance Criteria:**
  - Secrets stored encrypted at rest; round-trip tested.
  - unit tests cover encrypt/decrypt, key rotation.
  - CLI command `cloudmoor config vault test` validates setup.
- **Owner:** Cloud integrations engineer

### TCK-005 · F · P1 · 8 pts

**Title:** Create config persistence layer

- **Description:** Implement SQLite schema for providers, mounts, jobs, audit logs; include migrations.
- **Dependencies:** TCK-003
- **Acceptance Criteria:**
  - Schema generated via migrations; `go test` covers CRUD operations.
  - Database file path configurable; migration applied on startup.
  - Export/import CLI commands serialize configs to YAML.
- **Owner:** Lead Go engineer

### TCK-006 · F · P2 · 5 pts

**Title:** Skeleton for gRPC + REST gateway

- **Description:** Define protobuf for core services (MountService, ProviderService) and generate gRPC/REST stubs.
- **Dependencies:** TCK-003, TCK-005
- **Acceptance Criteria:**
  - Proto files linted; server compiles.
  - REST gateway accessible via `/api/v1/health`.
  - Example client integration test passes.
- **Owner:** Lead Go engineer

### TCK-007 · D · P2 · 3 pts

**Title:** Author contributor guidelines and coding standards

- **Description:** Document branching strategy, lint rules, testing expectations, PR checklist.
- **Dependencies:** TCK-001, TCK-002
- **Acceptance Criteria:**
  - `CONTRIBUTING.md` merged and referenced in README.
  - Templates for bugs, features, security issues added.
- **Owner:** Product/UX + Lead Go engineer

### TCK-008 · R · P2 · 2 pts

**Title:** Validate rclone library licensing and compatibility

- **Description:** Review rclone MIT license obligations, document attribution requirements.
- **Dependencies:** None
- **Acceptance Criteria:**
  - Legal note in docs with attribution plan.
  - Decision recorded on bundling approach vs extension link.
- **Owner:** Product/UX

---

## Milestone M1 – Core Providers & Daemon (4 weeks)

### TCK-101 · F · P1 · 13 pts

**Title:** Implement FTP/SFTP connector

- **Description:** Build connector supporting username/password and SSH key auth, passive mode, validation tests.
- **Dependencies:** M0 completion, TCK-004, TCK-005
- **Acceptance Criteria:**
  - Connector passes integration tests against vsftpd container.
  - CLI wizard prompts for credentials, stores in vault.
  - Mount/unmount operations succeed with smoke test.
- **Owner:** Cloud integrations engineer

### TCK-102 · F · P1 · 8 pts

**Title:** Implement Amazon S3/MinIO connector

- **Description:** Support IAM keys, STS assume role, custom endpoints; configure VFS cache defaults.
- **Dependencies:** TCK-004, TCK-005, TCK-006
- **Acceptance Criteria:**
  - Integration tests using MinIO container.
  - Bucket listing, upload/download validated.
  - Support for environment credential fallback.
- **Owner:** Cloud integrations engineer

### TCK-103 · F · P1 · 8 pts

**Title:** Implement Backblaze B2 connector

- **Description:** Integrate application keys, large file upload resume, throughput tuning.
- **Dependencies:** TCK-004, TCK-005, TCK-006
- **Acceptance Criteria:**
  - Integration tests using B2 sandbox or mock.
  - Large file (>2GB) upload resume verified.
  - Error retries with exponential backoff implemented.
- **Owner:** Cloud integrations engineer

### TCK-104 · F · P1 · 8 pts

**Title:** Implement WebDAV connector with presets

- **Description:** Provide generic WebDAV plus presets (Nextcloud, SharePoint); certificate options.
- **Dependencies:** TCK-004, TCK-005, TCK-006
- **Acceptance Criteria:**
  - Integration tests via local WebDAV container.
  - Presets selectable in CLI/Web UI; config saved.
  - TLS verification toggle documented.
- **Owner:** Cloud integrations engineer

### TCK-105 · F · P1 · 13 pts

**Title:** Develop mount manager & FUSE integration layer

- **Description:** Abstract platform-specific mounting via rclone VFS; manage lifecycle, retry policies.
- **Dependencies:** TCK-003, TCK-005, connectors (TCK-101..TCK-104)
- **Acceptance Criteria:**
  - Mount manager API supports create/start/stop/status.
  - Linux/macOS integration tests pass in CI.
  - Windows mount validated manually with WinFsp.
- **Owner:** Lead Go engineer

### TCK-106 · F · P1 · 8 pts

**Title:** Build daemon scheduler and health monitor

- **Description:** Background workers for mount health, retries, metrics collection, notification hooks.
- **Dependencies:** TCK-105
- **Acceptance Criteria:**
  - Health checks detect disconnect and trigger auto-reconnect.
  - Metrics exposed at `/metrics` Prometheus endpoint.
  - Alert hooks (log/Webhook stub) run on failure.
- **Owner:** Lead Go engineer

### TCK-107 · F · P1 · 8 pts

**Title:** Implement CLI workflows for config and mounts

- **Description:** Interactive commands for `config create`, `mount add`, status listing, log tailing.
- **Dependencies:** TCK-101..TCK-106
- **Acceptance Criteria:**
  - CLI UX reviewed, supports JSON output.
  - `cloudmoor mount <name>` works end-to-end.
  - Unit tests cover command parsing and API calls.
- **Owner:** Lead Go engineer

### TCK-108 · T · P2 · 5 pts

**Title:** Add integration test harness with Testcontainers

- **Description:** Create Go test suite spinning up FTP, MinIO, WebDAV containers for regression coverage.
- **Dependencies:** TCK-101..TCK-104
- **Acceptance Criteria:**
  - CI job runs integration suite nightly.
  - Tests validate mount lifecycle per connector.
  - Docs explain how to run locally.
- **Owner:** QA/Automation

### TCK-109 · D · P2 · 3 pts

**Title:** Draft architecture overview documentation

- **Description:** Convert sections of plan into `/docs/architecture.md`, include diagrams.
- **Dependencies:** TCK-105, TCK-106
- **Acceptance Criteria:**
  - Docs include component diagram, sequence for mount lifecycle.
  - Linked from README.
- **Owner:** Product/UX + Lead Go engineer

---

## Milestone M2 – OAuth Providers & Web UI (4 weeks)

### TCK-201 · F · P1 · 8 pts

**Title:** Implement OAuth device flow service

- **Description:** Build shared OAuth handler supporting device code, PKCE, token refresh persistence.
- **Dependencies:** M1 completion, TCK-004, TCK-006
- **Acceptance Criteria:**
  - Works with Dropbox sandbox; tokens stored securely.
  - CLI command `cloudmoor login` completes device flow.
  - Unit/integration tests cover expiration handling.
- **Owner:** Cloud integrations engineer

### TCK-202 · F · P1 · 5 pts

**Title:** Dropbox connector with OAuth

- **Description:** Use OAuth service to integrate Dropbox API, metadata caching.
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - OAuth flow completes via CLI & Web UI.
  - File operations succeed in integration tests.
  - Rate limit retries implemented.
- **Owner:** Cloud integrations engineer

### TCK-203 · F · P1 · 5 pts

**Title:** Google Drive connector (OAuth & service account)

- **Description:** Support OAuth client and service account auth, handle shared drives.
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - Integration tests cover both auth modes.
  - Change detection sync verified.
  - Scopes documented and minimal.
- **Owner:** Cloud integrations engineer

### TCK-204 · F · P1 · 5 pts

**Title:** OneDrive connector (personal & business)

- **Description:** Implement OneDrive graph API integration with tenant selection.
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - Device flow works for personal Microsoft account.
  - Business tenant selection documented.
  - Integration tests validate upload/download.
- **Owner:** Cloud integrations engineer

### TCK-205 · F · P1 · 5 pts

**Title:** Box connector with JWT + OAuth

- **Description:** Support enterprise JWT apps and standard OAuth; manage key rotation.
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - JWT auth path validated in integration test (mock or sandbox).
  - Token refresh + rotation automated.
  - Error handling for enterprise policies.
- **Owner:** Cloud integrations engineer

### TCK-206 · F · P1 · 5 pts

**Title:** pCloud connector with OAuth

- **Description:** Integrate pCloud API with folder selection, cache policies.
- **Dependencies:** TCK-201
- **Acceptance Criteria:**
  - OAuth flow tested with sandbox account.
  - Upload/download and metadata listing confirmed.
  - Documentation for API key setup.
- **Owner:** Cloud integrations engineer

### TCK-207 · F · P1 · 8 pts

**Title:** Build Web UI foundation (React + Vite)

- **Description:** Scaffold React app, integrate Tailwind, set up routing, connect to API gateway.
- **Dependencies:** TCK-006
- **Acceptance Criteria:**
  - Web UI builds via `npm run build` and served by daemon.
  - Auth flow with API tokens/OAuth login works.
  - Basic dashboard page renders mount list.
- **Owner:** Frontend engineer

### TCK-208 · F · P1 · 8 pts

**Title:** Implement Web UI configuration wizard

- **Description:** Multistep forms for provider setup, validation, credential entry, summary.
- **Dependencies:** TCK-207, connectors (TCK-202..TCK-206)
- **Acceptance Criteria:**
  - Wizard supports OAuth device initiation and completion.
  - Schema-driven forms render per provider.
  - Errors surfaced inline with accessible messaging.
- **Owner:** Frontend engineer

### TCK-209 · F · P2 · 5 pts

**Title:** Metrics & logging enhancements

- **Description:** Add structured logging with Zap fields, Prometheus metrics (latency, throughput, cache), expose `/metrics`.
- **Dependencies:** TCK-106
- **Acceptance Criteria:**
  - Metrics visible in Prometheus scrape.
  - Log entries include correlation IDs.
  - Sample Grafana dashboard provided.
- **Owner:** DevOps/SRE

### TCK-210 · D · P2 · 3 pts

**Title:** Update documentation with OAuth provider guides

- **Description:** Write guides for Dropbox, Google Drive, OneDrive, Box, pCloud setup.
- **Dependencies:** TCK-202..TCK-206
- **Acceptance Criteria:**
  - `/docs/providers/*.md` created with step-by-step instructions.
  - Screenshots or references for OAuth console configuration.
- **Owner:** Product/UX

---

## Milestone M3 – Advanced Providers & UX (3 weeks)

### TCK-301 · F · P1 · 8 pts

**Title:** Mega connector with client-side encryption

- **Description:** Implement Mega API integration, handle key derivation, throttling limits.
- **Dependencies:** M2 completion, TCK-004
- **Acceptance Criteria:**
  - Integration tests using Mega SDK or mock.
  - Encryption key handling validated and documented.
  - Throttle management prevents API bans during stress test.
- **Owner:** Cloud integrations engineer

### TCK-302 · F · P1 · 5 pts

**Title:** Advanced cache policy controls

- **Description:** Implement per-mount cache size/TTL, offline mode, manual purge command.
- **Dependencies:** TCK-105, TCK-106
- **Acceptance Criteria:**
  - Cache policies configurable via CLI/Web UI.
  - Offline mode retains cached data and queues writes.
  - `cloudmoor cache purge` clears caches and logs outcome.
- **Owner:** Lead Go engineer

### TCK-303 · F · P1 · 5 pts

**Title:** RBAC & multi-user support in daemon

- **Description:** Add user management, roles (Admin, Operator, Read-only), auth tokens.
- **Dependencies:** TCK-006, TCK-207
- **Acceptance Criteria:**
  - RBAC enforced at API layer; tests cover role restrictions.
  - Web UI user management page functional.
  - Audit logs capture role changes.
- **Owner:** Lead Go engineer + Frontend engineer

### TCK-304 · F · P2 · 5 pts

**Title:** Web UI polish & accessibility pass

- **Description:** Apply design system, dark mode, keyboard navigation, localization scaffolding.
- **Dependencies:** TCK-207, TCK-208
- **Acceptance Criteria:**
  - Lighthouse accessibility score ≥90.
  - i18n strings externalized; locale switcher prototype.
  - Keyboard navigation validated in QA checklist.
- **Owner:** Frontend engineer

### TCK-305 · F · P2 · 5 pts

**Title:** Desktop tray prototype (Tauri/Electron)

- **Description:** Build lightweight desktop app showing mount status, start/stop controls.
- **Dependencies:** TCK-107, TCK-303
- **Acceptance Criteria:**
  - Prototype runs on macOS & Windows.
  - Communicates with daemon via localhost API.
  - Basic notifications on mount failure.
- **Owner:** Frontend engineer

### TCK-306 · D · P2 · 3 pts

**Title:** Publish operations runbooks

- **Description:** Document mount failure remediation, credential rotation, disaster recovery.
- **Dependencies:** TCK-302, TCK-303
- **Acceptance Criteria:**
  - Runbooks stored under `/docs/operations/`.
  - Includes flowcharts or checklists for common incidents.
- **Owner:** DevOps/SRE

---

## Milestone M4 – Hardening & Release (2 weeks)

### TCK-401 · T · P1 · 5 pts

**Title:** Security review & dependency audit

- **Description:** Perform dependency scan, SBOM generation, penetration test checklist.
- **Dependencies:** Previous milestones
- **Acceptance Criteria:**
  - SBOM uploaded with release assets.
  - Critical findings resolved or tracked.
  - Security review report stored internally.
- **Owner:** DevOps/SRE + Lead Go engineer

### TCK-402 · F · P1 · 5 pts

**Title:** Packaging & installer tooling

- **Description:** Create Homebrew tap, Debian/RPM packages, Windows installer (MSI) via WiX or NSIS.
- **Dependencies:** TCK-107, TCK-207
- **Acceptance Criteria:**
  - Installers tested on target platforms.
  - `cloudmoor service install` works post-install.
  - Docs updated with install instructions.
- **Owner:** DevOps/SRE

### TCK-403 · F · P1 · 5 pts

**Title:** User documentation & onboarding guides

- **Description:** Produce getting started guide, provider-specific tutorials, FAQ.
- **Dependencies:** TCK-109, TCK-210, TCK-306
- **Acceptance Criteria:**
  - `/docs/getting-started.md` + tutorial videos/screenshots.
  - FAQ covers top support scenarios.
  - Documentation reviewed by Product/UX.
- **Owner:** Product/UX

### TCK-404 · T · P2 · 3 pts

**Title:** Beta release program launch

- **Description:** Set up feedback channels, telemetry opt-in, announce beta to pilot users.
- **Dependencies:** TCK-401..TCK-403
- **Acceptance Criteria:**
  - Landing page for beta signups live.
  - Feedback triage process documented.
  - First batch of pilot testers onboarded.
- **Owner:** Product/UX

### TCK-405 · T · P2 · 3 pts

**Title:** Performance benchmarking report

- **Description:** Execute benchmark suite, document throughput, latency, resource usage across providers.
- **Dependencies:** TCK-302, TCK-209
- **Acceptance Criteria:**
  - Report published in `/docs/performance.md`.
  - Regression thresholds defined for future releases.
  - Findings reviewed with engineering team.
- **Owner:** QA/Automation + DevOps/SRE

---

## Cross-Cutting & Backlog Items

### TCK-501 · R · P3 · 3 pts

**Title:** Investigate external secret manager integrations

- **Description:** Research architecture for HashiCorp Vault, AWS Secrets Manager extensions.
- **Dependencies:** TCK-004
- **Acceptance Criteria:**
  - Proposal document with recommended approach.
  - Spike code proving feasibility.
- **Owner:** Cloud integrations engineer

### TCK-502 · F · P3 · 5 pts

**Title:** Implement localization framework

- **Description:** Add ICU message format support, locale detection, translation files.
- **Dependencies:** TCK-304
- **Acceptance Criteria:**
  - Web UI builds with language packs.
  - CLI supports `--lang` flag for prompts.
  - English + Spanish translations verified.
- **Owner:** Frontend engineer

### TCK-503 · F · P3 · 5 pts

**Title:** Add scheduled sync jobs with cron expressions

- **Description:** Allow mounts to run sync tasks on schedule with conflict resolution options.
- **Dependencies:** TCK-302, TCK-106
- **Acceptance Criteria:**
  - Cron parser integrated; jobs persisted in DB.
  - Sync conflicts resolved via policy settings.
  - Tests cover overlapping schedules.
- **Owner:** Lead Go engineer

### TCK-504 · F · P3 · 5 pts

**Title:** Implement webhook & email notifications

- **Description:** Notify users on mount failures, token expiry via configurable channels.
- **Dependencies:** TCK-209, TCK-303
- **Acceptance Criteria:**
  - Webhooks configurable per mount.
  - Email notifications sent via SMTP integration (mock in dev).
  - Audited in logs.
- **Owner:** DevOps/SRE

### TCK-505 · R · P3 · 3 pts

**Title:** Evaluate MFA requirements for Web UI

- **Description:** Gather stakeholder needs, review implementation options, scope effort.
- **Dependencies:** TCK-303
- **Acceptance Criteria:**
  - Recommendation doc with effort estimate.
  - Decision captured in roadmap backlog.
- **Owner:** Product/UX

---

## Notes

- Estimates are preliminary and should be refined during sprint planning.
- Tickets marked P1 are blockers for milestone completion.
- Backlog items can be pulled into milestones after foundational work stabilizes.
