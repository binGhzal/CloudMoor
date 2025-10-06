# CloudMoor Remote Drive Mounting Application Plan

Date: 2025-10-06

## 1. Vision & Objectives

- Provide a unified application that mounts remote storage providers as local drives with persistent configuration across reboots.
- Support both power users (CLI) and always-on service deployments (headless daemon with optional Web UI), with an extensible path to a desktop GUI client.
- Deliver a plugin-based connector layer launching with Amazon S3/MinIO, WebDAV, and Dropbox, while maintaining a deferred backlog for additional providers.
- Focus on secure credential management, reliable reconnection on network hiccups, and minimal latency through smart caching.
- Maintain an architecture that can scale from single-user laptops to multi-tenant server deployments with minimal operational friction.

## 2. Target Platforms & Runtime

- **Operating systems:** macOS, Linux, Windows (WinFsp required for FUSE-like functionality).
- **Primary language:** Go (mature FUSE support via `bazil.org/fuse`, excellent cross-compilation, and ecosystem around `rclone` libraries for cloud connectors). Secondary components (Web UI) can use TypeScript/React.
- **Mount backend:** Use the `rclone` core libraries as a shared engine for remote providers, wrapped in a daemon that exposes a consistent API. Leverage `cgofuse` on macOS/Windows, and native FUSE on Linux.
- **Packaging philosophy:** Ship a single-statically linked binary per OS for the daemon/CLI, bundled Web UI assets, and provide optional container images for headless deployments.

## 3. Application Modes

1. **CLI Tool** (`cloudmoor`):
   - Commands: `config`, `mount`, `unmount`, `status`, `list`, `logs`, `service install/remove`, `cache purge`.
   - Support configuration import/export (YAML/TOML) and interactive setup wizards per provider.
   - Provide scripted workflows via `--json` and `--watch` flags for automation.
2. **Daemon Service** (`cloudmoord`):
   - Runs as a background process managing mounts, reconnection, scheduling, and exposure of a local gRPC/REST API.
   - Install scripts for systemd (Linux), launchd (macOS), and Windows Service.
   - Supports clustering mode (future) via shared metadata store.
3. **Optional Web UI** (`cloudmoor-console`):
   - Single-page app served by the daemon, providing dashboards, mount controls, activity logs, and credential management with RBAC.
   - Responsive design for mobile/tablet administration.
4. **Future Desktop GUI:**
   - Electron or Tauri app reusing daemon APIs for richer tray controls.
   - Offers OS-native notifications and quick actions.

### 4.1 Deployment Topologies

- **Local single-user:** CLI launches daemon on demand, persists configs locally.
- **Server/headless:** Daemon installed as service, Web UI enabled for remote management over HTTPS.
- **Multi-user teams:** RBAC via daemon, SSO integration (OIDC) scheduled for later roadmap.

### 4.2 Component Responsibility Matrix

| Component           | Responsibilities                                                                 | Technologies & Notes                                                |
| ------------------- | -------------------------------------------------------------------------------- | ------------------------------------------------------------------- |
| `cloudmoor` CLI     | User interaction, local config bootstrapping, daemon control, scripting hooks.   | Cobra-based CLI, communicates with daemon via gRPC/REST.            |
| `cloudmoord` daemon | Mount orchestration, scheduling, credential management, API hosting.             | Go services with worker pool, background scheduler.                 |
| Connector plugins   | Translate provider-specific config into rclone backend mounts; handle auth flow. | Go interfaces, compiled-in plugins, optional WASM sandbox (future). |
| Mount manager       | FUSE/WinFsp filesystem implementation, event hooks, metrics emission.            | rclone VFS, custom adapters.                                        |
| Persistence layer   | Store configs, secrets metadata, job history, audit logs.                        | SQLite + optional external Postgres (future).                       |
| Web UI              | Visualization, configuration wizards, admin controls.                            | React + Vite + Tailwind, served as static assets.                   |
| Telemetry stack     | Logging, metrics, tracing.                                                       | Zap, OpenTelemetry exporters, Prometheus endpoints.                 |

### 4.3 Control Plane vs Data Plane

- **Control plane:** Daemon services, APIs, scheduler, configuration management, credential vault, Web UI. Scales vertically; future horizontal scaling via distributed metadata store.
- **Data plane:** rclone VFS processes handling file operations between local mountpoint and remote provider. Runs per mount, isolated via worker goroutines and concurrency limits. Supports QoS controls and throttling policies.

### 4.4 Mount Lifecycle

1. **Definition:** User creates/updates a mount definition (name, provider, credentials, local mount path, cache policy).
2. **Validation:** Connector validates credentials and connectivity; errors surfaced immediately.
3. **Preparation:** Daemon provisions cache directories, ensures dependencies (WinFsp, FUSE) are present, loads secrets.
4. **Mount Activation:** rclone VFS process starts, mountpoint registered with OS, telemetry streams begin.
5. **Monitoring:** Health checks monitor latency, error rates, cache utilization. Auto-retry with exponential backoff on failure.
6. **Graceful Unmount:** On user request or shutdown, daemon flushes caches, ensures remote sync completion, tears down mount.
7. **Persistence:** Final state recorded in database for resume and audit.

### 4.5 Caching & Sync Strategy

- Layered caching: metadata cache (SQLite), file chunk cache (configurable size/TTL), optional read-ahead/write-back policies.
- Consistency modes: `eventual` (default), `write-through`, `read-only`. Users pick per mount.
- Background sync workers reconcile remote changes via provider-specific change notifications (where available) or periodic listing.
- Cache eviction strategies: LRU for file blocks, TTL for metadata, manual purge command.
- Support offline mode: persisted cached files remain accessible; writes queued and replayed when connectivity resumes.

### 4.6 API Surface

| Service            | Key Methods                                                                                  | Notes                                                   |
| ------------------ | -------------------------------------------------------------------------------------------- | ------------------------------------------------------- |
| `MountService`     | `ListMounts`, `CreateMount`, `UpdateMount`, `DeleteMount`, `Mount`, `Unmount`, `GetMetrics`. | gRPC + REST; supports watch streams for status updates. |
| `ProviderService`  | `ListProviders`, `ValidateConfig`, `StartOAuthFlow`, `CompleteOAuthFlow`.                    | Device code flow support, secrets handled server-side.  |
| `ConfigService`    | `ExportConfig`, `ImportConfig`, `GetSettings`, `UpdateSettings`.                             | YAML/TOML export, versioned settings.                   |
| `AuthService`      | `Login`, `Refresh`, `ListUsers`, `CreateUser`, `AssignRole`.                                 | RBAC, optional external IdP integration.                |
| `TelemetryService` | `StreamLogs`, `StreamMetrics`, `GetAuditTrail`.                                              | Websocket support for Web UI dashboards.                |

- REST gateway follows `/api/v1/...` naming, with OpenAPI spec generated for client SDKs.
- CLI uses gRPC directly when running on same host, falling back to REST over HTTPS when remote.
  1.  **M0 – Foundations & rclone strategy (3 weeks):**
      - Bootstrap repo, choose tooling, set up CI/CD (GitHub Actions) and GoReleaser pipeline.
      - Implement config persistence, credential vault MVP, pluggable connector interface.
      - Run a comparative spike on embedding `librclone` versus orchestrating the rclone binary over RPC; capture the decision in `/docs/decisions/` (ticket TCK-009).
  2.  **M1 – Core Mount Experience (6 weeks):**
      - Ship S3/MinIO and WebDAV connectors with integration tests and cache MVP.
      - Deliver mount manager, daemon service (gRPC/REST), and CLI workflows for config/mount lifecycle.
      - Invest in cross-platform build tooling (macOS + Windows runners with MacFUSE/WinFsp) to unblock CGO releases.
  3.  **M2 – OAuth & Web UI Baseline (4 weeks):**
      - Harden the shared OAuth device flow service and complete Dropbox connector end-to-end.
      - Stand up Web UI alpha (config wizard, mount dashboard) and Prometheus/Zap observability improvements.
      - Close feedback loop on security (vault auditability, token refresh metrics).
  4.  **M3 – Reliability & Access Controls (4 weeks):**
      - Extend cache controls, add RBAC + audit logs, and refine Web UI ergonomics.
      - Prepare packaging (service installers, container image) and operations runbooks for beta.
      - Reassess deferred connectors based on performance and support signals.
  5.  **M4 – Beta Hardening & Release Prep (3 weeks):** - Complete security review, performance benchmarks, documentation suite, and beta program launch. - Finalise packaging artefacts, telemetry opt-in, and feedback triage cadence.
      | `jobs` | Background tasks (sync, cache purge, snapshot). | `id`, `job_type`, `payload`, `schedule`, `status`. |

### 5.2 Configuration Versioning & Migration

- Every change to provider or mount definitions increments a semantic version stored in `mounts.version`.
- Migrations handled via `golang-migrate`, with pre-flight backups and post-migration smoke tests.
- Configuration export includes schema version; import process performs compatibility validation and prompts for remediation when required.

### 5.3 Backup & Restore Workflow

1. Nightly cron job triggers `cloudmoor backup create`.
2. Daemon snapshots SQLite DB, encrypts archive, and uploads to user-selected storage (S3, local path, etc.).
3. Restoration performed via CLI/web wizard with integrity checks and dry-run validation.
4. Supports point-in-time recovery using WAL files retained for X days (configurable).

## 6. Security Considerations

- Credentials encrypted using AES-GCM with master key stored in OS keychain where available; fallback to passphrase-protected file.
- Option to integrate with cloud secret managers (AWS Secrets Manager, Azure Key Vault, GCP Secret Manager) via provider-specific config.
- TLS everywhere: self-signed bootstrap with user-provided certs, mutual TLS for remote admin.
- Role-based access in daemon API: Admin, Operator, Read-only.
- Audit logging for credential access, mount lifecycle events, and UI actions.
- Support for FIPS-compliant crypto modules when running in regulated environments.

### 6.1 Threat Model Snapshot

- **Assets:** Credentials, cached data, configuration exports, audit logs.
- **Attack vectors:** Local privilege escalation, intercepted OAuth flows, compromised API tokens, supply-chain tampering.
- **Mitigations:** Least-privilege daemon service accounts, short-lived OAuth device codes, signed release artifacts, SBOM generation, dependency scanning, tamper-evident audit logs.

### 6.2 Compliance & Audit Readiness

- Map controls to SOC2 (security, availability), ISO 27001 Annex A, and GDPR data subject rights.
- Provide data retention configuration, audit log export, and DPIA guidance for EU customers.
- Document secure deployment checklist (network segmentation, TLS cert rotation, secret rotation cadence).

## 7. Provider Integration Strategy

| Provider          | Protocol/SDK                          | Notes                                                                               |
| ----------------- | ------------------------------------- | ----------------------------------------------------------------------------------- |
| Amazon S3 / MinIO | `rclone/backend/s3` (custom endpoint) | Launch connector; supports IAM roles, static keys, and S3-compatible endpoints.     |
| WebDAV            | `rclone/backend/webdav`               | Launch connector; presets for Nextcloud and SharePoint, TLS toggle for self-hosted. |
| Dropbox           | `rclone/backend/dropbox`              | Launch connector; exercises shared OAuth device flow service and token lifecycle.   |

- Each launch connector defines validation schema, environment variable defaults, and UI form layout metadata.
- Deferred providers remain part of the product vision and are documented in §7.2 for future milestones.

### 7.1 Launch Providers (v1 scope)

#### Amazon S3 & MinIO

- Support IAM role assumption via AWS STS, web identity tokens, and static access keys.
- Allow custom region/endpoint configuration for MinIO, DigitalOcean Spaces, and other S3-compatible services.
- Expose sensible defaults for multipart uploads and optional SSE-KMS integration; advanced tuning lands post-GA.

#### WebDAV

- Provide presets for Nextcloud, SharePoint, and ownCloud with pre-filled endpoints and default TLS settings.
- Offer certificate pinning options for self-hosted instances while surfacing warnings when TLS verification is disabled.
- Document optimistic vs. pessimistic locking semantics to set user expectations.

#### Dropbox

- Implements the shared OAuth device flow service, including polling cadence, user messaging, and token refresh.
- Supports scope minimisation, incremental backoff for rate limits, and secure credential storage via the vault.
- Serves as the reference implementation for future OAuth-centric providers.

### 7.2 Deferred Connectors (post-v1 backlog)

- FTP/SFTP, Backblaze B2, OpenStack Swift, Google Drive, OneDrive, Box, Mega, and pCloud remain strategic targets but are intentionally deferred until after the beta launch.
- Each deferred connector retains lightweight research notes and acceptance criteria in the Ticket Backlog section of `docs/tasks.md` so they can be reintroduced when resourcing allows.
- Device-flow providers (Google Drive, OneDrive, Box, pCloud) will reuse the Dropbox patterns once the shared service proves stable.

### 7.3 Connector Certification Checklist

- Connectivity & auth validation (success/failure).
- Large file upload/download (>5 GB) throughput measurement.
- Metadata consistency tests (rename, move, permission changes).
- Offline/online reconnection scenarios.
- OAuth token revocation and recovery workflow.

## 8. User Experience

### CLI

- `cloudmoor login <provider>` launches OAuth device code where applicable.
- `cloudmoor mount <mount-name>` handles on-demand mount with interactive prompts if config missing.
- `cloudmoor service enable --web-ui` installs daemon and exposes port configuration.
- Rich `--json` output for scripting.
- `cloudmoor watch mounts` streams status updates to the terminal with color-coded output.

### Web UI

- React + Tailwind, served by daemon (or standalone static bundle).
- Features: Dashboard of mounts (status, throughput), configuration wizard per provider, activity log, user management, settings.
- Notifications: toast + webhooks/email for mount failures.
- Supports dark mode, responsive layout, keyboard shortcuts, and localization (initially English, Spanish, German).

### GUI (Future)

- Desktop tray with start/stop controls, quick shortcuts to open mounted volumes, status indicator.
- Native auto-update channel leveraging platform-specific installers.

### 8.1 User Journeys

| Persona         | Goal                                     | Journey Steps                                                                                                        |
| --------------- | ---------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| DevOps engineer | Mount S3 bucket on Linux server.         | Install CLI ➝ `cloudmoor service enable` ➝ Configure S3 mount via CLI wizard ➝ Validate with `cloudmoor status`.     |
| Power user      | Sync Dropbox files locally with caching. | Run daemon in background ➝ Complete OAuth via Web UI ➝ Enable offline cache mode ➝ Monitor throughput via dashboard. |
| IT admin        | Manage team access to mounts.            | Invite users via Web UI ➝ Assign roles ➝ Configure audit log exports ➝ Set up alerts for mount failures.             |

### 8.2 Accessibility & Localization

- WCAG 2.1 AA compliance target for Web UI (color contrast, focus management, ARIA labels).
- Provide keyboard navigable CLI prompts and optional high-contrast output mode.
- Localization framework is deferred until after the beta launch; the initial release ships with English-only UI strings while we gather feedback on terminology and tone.

## 9. Deployment & Packaging

- **Build tooling:** Go modules, Goreleaser for multi-platform binaries, Docker images for headless deployments.
- **Installers:** Homebrew tap (macOS), Debian/RPM packages, Chocolatey/Scoop for Windows.
- **Service setup:** Provide scripts (`cloudmoor service install`) to register daemon with systemd/launchd/WinFsp.
- **TLS management:** ACME integration for public endpoints; self-signed generator for LAN use.
- Deliver container images with non-root user, health checks, and optional sidecar exporter.

### 9.1 Distribution Artifacts

- `cloudmoor` & `cloudmoord` binaries (per OS/arch).
- Docker/OCI image (amd64, arm64) published to GHCR.
- Helm chart for Kubernetes deployments (daemonset per node, optional Web UI ingress).
- Terraform module snippets for cloud VM deployments.

### 9.2 Infrastructure-as-Code Targets

- Provide Ansible role for on-prem installs.
- Publish systemd unit, launchd plist, Windows service template in `/deploy` directory.
- Optional Pulumi program for managed environments.

## 10. Testing & QA Strategy

- Unit tests for mount manager, config parser, credential vault.
- Integration tests per provider using mocked endpoints or official sandboxes.
- End-to-end smoke tests that mount and validate filesystem semantics using temp directories.
- Performance benchmarks measuring throughput, cache hit rates, reconnection times.
- Static analysis: golangci-lint, SAST for credential handling.
- Chaos testing (fault injection) to verify resilience to network outages.

### 10.1 Automation Pyramid

- **Unit (70%)**: Pure Go tests, mocks for provider interfaces.
- **Integration (20%)**: Spin up containers (MinIO, WebDAV) via Testcontainers-go and mock Dropbox OAuth device responses.
- **E2E (10%)**: CLI-driven scenarios executed in GitHub Actions matrix (macOS, Linux, Windows).

### 10.2 Test Environments

| Environment | Purpose                            | Notes                                                       |
| ----------- | ---------------------------------- | ----------------------------------------------------------- |
| `dev`       | Rapid iteration, feature branches. | Uses local SQLite, minimal providers.                       |
| `qa`        | Pre-release validation.            | Realistic data sets, nightly test suites, metrics captured. |
| `staging`   | Release candidate verification.    | Mirrors production packaging, used for dogfooding.          |

### 10.3 Performance Benchmarks

- Measure sequential and parallel read/write throughput per provider with various cache modes.
- Track mount initialization time, reconnect latency, memory footprint under load.
- Publish benchmark dashboards and compare across releases.

## 11. Roadmap & Milestones

1. **M0 – Foundations & rclone strategy (3 weeks):**
   - Bootstrap repo, choose tooling, set up CI/CD (GitHub Actions) and GoReleaser pipeline.
   - Implement config persistence, credential vault MVP, pluggable connector interface.
   - Run a comparative spike on embedding `librclone` versus orchestrating the rclone binary over RPC; capture the decision in `/docs/decisions/` (ticket TCK-009).
2. **M1 – Core Mount Experience (6 weeks):**
   - Ship S3/MinIO and WebDAV connectors with integration tests and cache MVP.
   - Deliver mount manager, daemon service (gRPC/REST), and CLI workflows for config/mount lifecycle.
   - Invest in cross-platform build tooling (macOS + Windows runners with MacFUSE/WinFsp) to unblock CGO releases.
3. **M2 – OAuth & Web UI Baseline (4 weeks):**
   - Harden the shared OAuth device flow service and complete Dropbox connector end-to-end.
   - Stand up Web UI alpha (config wizard, mount dashboard) and Prometheus/Zap observability improvements.
   - Close feedback loop on security (vault auditability, token refresh metrics).
4. **M3 – Reliability & Access Controls (4 weeks):**
   - Extend cache controls, add RBAC + audit logs, and refine Web UI ergonomics.
   - Prepare packaging (service installers, container image) and operations runbooks for beta.
   - Reassess deferred connectors based on performance and support signals.
5. **M4 – Beta Hardening & Release Prep (3 weeks):**
   - Complete security review, performance benchmarks, documentation suite, and beta program launch.
   - Finalise packaging artefacts, telemetry opt-in, and feedback triage cadence.

### 11.1 Milestone Deliverables

| Milestone | Key Deliverables                                                            | Acceptance Criteria                                                        |
| --------- | --------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| M0        | Repo scaffolding, CI pipeline, connector interface, rclone integration RFC. | CI green on lint/test, documented decision on rclone embedding approach.   |
| M1        | S3/MinIO + WebDAV connectors, daemon + CLI, cache MVP.                      | Mount/unmount S3 and WebDAV in integration tests across macOS/Linux/Win.   |
| M2        | OAuth device flow service, Dropbox connector, Web UI alpha, observability.  | Dropbox device flow succeeds, Web UI shows mount status, metrics exported. |
| M3        | Cache controls, RBAC, packaging groundwork, runbooks.                       | RBAC enforced, cache tuning exposed in CLI/UI, install scripts verified.   |
| M4        | Security/packaging docs, beta launch readiness.                             | Security review signed off, docs published, beta cohort onboarded.         |

## 12. Risks & Mitigations

- **Cross-platform FUSE complexity:** Mitigate by leveraging battle-tested rclone mounting logic and WinFsp adapters; add integration test matrix in CI.
- **OAuth token management pitfalls:** Centralize through secure vault, refresh tokens proactively, allow external secret stores.
- **Provider API rate limits:** Implement adaptive throttling, exponential backoff, and cache metadata to reduce chatter.
- **Local cache growth:** Provide configurable quotas and eviction policies; add monitoring alerts.
- **User trust & security:** Offer transparency logs, optional remote wipe, and strong docs around credential handling.
- **Dependency drift:** Automate dependency scanning (Dependabot), maintain SBOM, pin versions.

### 12.1 Risk Monitoring Triggers

- Mount failure rate >2% per hour.
- OAuth token refresh failures >5% per day.
- Cache eviction time > target threshold (configurable SLA).
- Security advisories affecting FUSE/rclone dependencies.

## 13. Documentation & Developer Experience

- **Docs site:** mdBook or Docusaurus under `/docs` with tutorials, provider-specific guides, troubleshooting.
- **API references:** Auto-generate from protobuf definitions using `buf` or `grpc-gateway` swagger.
- **Examples:** Recipes for common mounts (e.g., MinIO in Kubernetes, S3 cross-account).
- **Contribution guide:** Coding standards, branching, testing requirements.
- Provide quickstart templates (Docker Compose, Terraform sample) and CLI cheat sheet.

### 13.1 Developer Tooling

- Pre-commit hooks for formatting, linting, license headers.
- Makefile tasks for build/test/release; `justfile` optional for convenience.
- Live reload for Web UI during development via Vite proxying to daemon.

### 13.2 Community & Support

- GitHub Discussions for Q&A, roadmap transparency.
- Issue templates for bug reports, feature requests, security disclosures.
- Monthly community call (post v1.0) to gather feedback.

## 14. Next Steps

- Validate Go + rclone library licensing alignment (rclone is MIT; ensure compliance).
- Prototype core daemon skeleton with one provider (S3) to verify mount performance.
- Complete the `librclone` vs. RPC orchestration spike (TCK-009) and log the outcome in `/docs/decisions/`.
- Design protobuf/gRPC contracts and generate stubs.
- Draft UX wireframes for Web UI and CLI flows.
- Set up CI pipeline skeleton (lint, tests, cross-build) and container packaging.
- Begin authoring developer handbook and user onboarding guides.

## 15. DevOps & CI/CD Pipeline

- GitHub Actions workflows for lint, unit tests, integration tests, packaging, release.
- Matrix builds across macOS, Linux, Windows, amd64/arm64.
- Release pipeline signs binaries (Cosign), publishes SBOM (Syft), and uploads artifacts to GH Releases.
- Canary channel: nightly builds with automated smoke tests; promote to beta after 3 consecutive green runs.

## 16. Operations & Monitoring Playbooks

- Runbook for common incidents (mount failure, credential expiry, cache saturation) with step-by-step remediation.
- Alerting via Prometheus Alertmanager → Slack/Email/Webhooks.
- Log retention guidance (e.g., 30 days local, 90 days centralized).
- Disaster recovery plan covering database restore, re-seeding secrets, redeploying daemon.

## 17. Resource & Staffing Plan

| Role                        | FTE  | Responsibilities                                            |
| --------------------------- | ---- | ----------------------------------------------------------- |
| Lead Go engineer            | 1    | Core daemon, mount manager, connector framework.            |
| Cloud integrations engineer | 1    | Implement/maintain provider connectors, OAuth flows.        |
| Frontend engineer           | 0.5  | Web UI, design system, accessibility.                       |
| DevOps/SRE                  | 0.5  | CI/CD, packaging, observability, deployments.               |
| QA/Automation               | 0.5  | Test frameworks, integration suites, release validation.    |
| Product/UX                  | 0.25 | User research, roadmap management, documentation oversight. |

## 18. Success Metrics & KPIs

| Metric                        | Target                  | Measurement                              |
| ----------------------------- | ----------------------- | ---------------------------------------- |
| Mount success rate            | ≥99.5% per 7-day window | Telemetry counters, alert if below.      |
| Reconnect time after outage   | <60 seconds median      | Daemon metrics, chaos test verification. |
| OAuth flow completion         | ≥95% success            | AuthService logs, user feedback.         |
| Web UI performance            | Time-to-interactive <3s | Lighthouse CI, synthetic monitoring.     |
| Support tickets per 100 users | <5 per month            | Helpdesk integration (future).           |

## 19. Future Enhancements Backlog

- Multi-tenant RBAC with LDAP/AD integration.
- Scheduled sync jobs with cron-like syntax and conflict resolution policies.
- Differential sync (rsync-like delta transfers) where providers support checksums.
- Built-in file sharing links per provider (Dropbox, Box) from Web UI.
- Mobile admin app leveraging REST API.
- Secrets backend integrations (Vault, AWS KMS) GA.

## 20. Assumptions & Open Questions

- Assume rclone licensing remains compatible with CloudMoor distribution; verify attribution requirements.
- Clarify whether offline writeback conflicts should favor local or remote changes by default.
- Determine minimal desktop OS versions supported (Windows 10+, macOS 12+, Ubuntu 22.04 LTS?).
- Confirm need for HIPAA/PCI compliance in early releases.
- Evaluate demand for multi-factor authentication in Web UI for v1.0 or post-launch.

## 21. Delivery Governance & Tracking

- **Milestone exit reviews:** Each milestone (M0–M4) must pass a formal review covering scope completion, outstanding defects, security checklist, and documentation updates before progression.
- **Weekly steering sync:** Cross-functional leads meet weekly to review burn-down charts, risk register updates, and unblock critical dependencies.
- **Change control:** Material roadmap or scope adjustments require a written RFC reviewed by engineering, product, and security owners; accepted changes are reflected in both `plan.md` and `tasks.md` within 24 hours.
- **Quality gates:** No feature advances to the next environment without passing automated lint/test suites, integration smoke tests, and relevant security scanning jobs captured in CI pipelines.
- **Traceability:** Each task in `tasks.md` references its originating ticket identifier to ensure traceable linkage from strategic plan → backlog → execution. Deviations are logged in a lightweight decision log stored under `/docs/decisions/`.
