# CloudMoor Remote Drive Mounting Application Plan

Date: 2025-10-06

## 1. Vision & Objectives

- Provide a unified application that mounts remote storage providers as local drives with persistent configuration across reboots.
- Support both power users (CLI) and always-on service deployments (headless daemon with optional Web UI), with an extensible path to a desktop GUI client.
- Deliver a plugin-based connector layer covering the following services at launch: FTP, SFTP, Amazon S3, MinIO, Backblaze B2, OpenStack Swift, Dropbox, Google Drive, OneDrive, Box, Mega, WebDAV, and pCloud.
- Focus on secure credential management, reliable reconnection on network hiccups, and minimal latency through smart caching.

## 2. Target Platforms & Runtime

- **Operating systems:** macOS, Linux, Windows (WinFsp required for FUSE-like functionality).
- **Primary language:** Go (mature FUSE support via `bazil.org/fuse`, excellent cross-compilation, and ecosystem around `rclone` libraries for cloud connectors). Secondary components (Web UI) can use TypeScript/React.
- **Mount backend:** Use the `rclone` core libraries as a shared engine for remote providers, wrapped in a daemon that exposes a consistent API. Leverage `cgofuse` on macOS/Windows, and native FUSE on Linux.

## 3. Application Modes

1. **CLI Tool** (`cloudmoor`):
   - Commands: `config`, `mount`, `unmount`, `status`, `list`, `logs`, `service install/remove`, `cache purge`.
   - Support configuration import/export (YAML/TOML) and interactive setup wizards per provider.
2. **Daemon Service** (`cloudmoord`):
   - Runs as a background process managing mounts, reconnection, scheduling, and exposure of a local gRPC/REST API.
   - Install scripts for systemd (Linux), launchd (macOS), and Windows Service.
3. **Optional Web UI** (`cloudmoor-console`):
   - Single-page app served by the daemon, providing dashboards, mount controls, activity logs, and credential management with RBAC.
4. **Future Desktop GUI:**
   - Electron or Tauri app reusing daemon APIs for richer tray controls.

## 4. High-Level Architecture

- **Core Daemon:** Orchestrates mounts, retries, monitoring. Written in Go, exposing gRPC and REST endpoints.
- **Connector Plugins:** Thin wrappers around rclone remotes with uniform lifecycle (init, validate, mount, sync, cleanup). Configurable via JSON Schema definitions stored in the database.
- **Mount Manager:** Abstracts FUSE operations; handles platform-specific adapters (FUSE, WinFsp).
- **Persistence Layer:** SQLite (via `modernc.org/sqlite` for pure Go) storing profiles, credentials (encrypted), mount states, cache metadata.
- **Secrets Vault:** Pluggable backends: file-based encrypted keystore (default), with future support for OS keychains, HashiCorp Vault.
- **API Gateway:** gRPC endpoints (for internal clients) with auto-generated REST gateway for CLI/Web UI integrations.
- **Task & Event Bus:** Lightweight pub/sub (e.g., `nats.go` embedded JetStream or Go channels) for broadcasting mount events to UI/logging subsystems.
- **Observability:** Structured logging (Zap), metrics via Prometheus, health endpoints.

### 4.1 Deployment Topologies

- **Local single-user:** CLI launches daemon on demand, persists configs locally.
- **Server/headless:** Daemon installed as service, Web UI enabled for remote management over HTTPS.
- **Multi-user teams:** RBAC via daemon, SSO integration (OIDC) scheduled for later roadmap.

## 5. Persistence & State Management

- **Configuration DB:** Stores provider configs, mount definitions, schedules, user roles, and audit logs.
- **Cache Directory:** Per-mount configurable location with size/TTL policies; optional persistent caching using `rclone` VFS cache.
- **Snapshotting:** Periodic export of configs to encrypted archive for backup/disaster recovery.
- **Resume Strategy:** On startup, daemon reads persisted mounts and attempts reconnection respecting backoff policies.

## 6. Security Considerations

- Credentials encrypted using AES-GCM with master key stored in OS keychain where available; fallback to passphrase-protected file.
- Option to integrate with cloud secret managers (AWS Secrets Manager, Azure Key Vault, GCP Secret Manager) via provider-specific config.
- TLS everywhere: self-signed bootstrap with user-provided certs, mutual TLS for remote admin.
- Role-based access in daemon API: Admin, Operator, Read-only.
- Audit logging for credential access, mount lifecycle events, and UI actions.

## 7. Provider Integration Strategy

| Provider        | Protocol/SDK                             | Notes                                                                       |
| --------------- | ---------------------------------------- | --------------------------------------------------------------------------- |
| FTP/SFTP        | `rclone/backend/sftp`                    | Combine FTP & SFTP connectors, support key-based auth, passive mode.        |
| Amazon S3       | `rclone/backend/s3`                      | Support multiple credential sources (IAM roles, access keys, web identity). |
| MinIO           | `rclone/backend/s3` with custom endpoint | Allow custom region/endpoint, TLS settings.                                 |
| Backblaze B2    | `rclone/backend/b2`                      | Handle application keys; tune upload concurrency.                           |
| OpenStack Swift | `rclone/backend/swift`                   | Support Keystone v2/v3 auth.                                                |
| Dropbox         | `rclone/backend/dropbox`                 | OAuth 2.0 flow; refresh tokens stored securely.                             |
| Google Drive    | `rclone/backend/drive`                   | Service account & OAuth options; file metadata caching.                     |
| OneDrive        | `rclone/backend/onedrive`                | Support personal and business variants.                                     |
| Box             | `rclone/backend/box`                     | JWT-based enterprise auth and OAuth for individuals.                        |
| Mega            | `rclone/backend/mega`                    | Client-side encryption support; watch throttling.                           |
| WebDAV          | `rclone/backend/webdav`                  | Generic connector; preset profiles (Nextcloud, SharePoint).                 |
| pCloud          | `rclone/backend/pcloud`                  | OAuth flow; ensure support for custom directories.                          |

- Each connector defines validation schema, environment variable defaults, and UI form layout metadata.
- Provide integration tests per connector using local test containers or mocked APIs where feasible.

## 8. User Experience

### CLI

- `cloudmoor login <provider>` launches OAuth device code where applicable.
- `cloudmoor mount <mount-name>` handles on-demand mount with interactive prompts if config missing.
- `cloudmoor service enable --web-ui` installs daemon and exposes port configuration.
- Rich `--json` output for scripting.

### Web UI

- React + Tailwind, served by daemon (or standalone static bundle).
- Features: Dashboard of mounts (status, throughput), configuration wizard per provider, activity log, user management, settings.
- Notifications: toast + webhooks/email for mount failures.

### GUI (Future)

- Desktop tray with start/stop controls, quick shortcuts to open mounted volumes, status indicator.

## 9. Deployment & Packaging

- **Build tooling:** Go modules, Goreleaser for multi-platform binaries, Docker images for headless deployments.
- **Installers:** Homebrew tap (macOS), Debian/RPM packages, Chocolatey/Scoop for Windows.
- **Service setup:** Provide scripts (`cloudmoor service install`) to register daemon with systemd/launchd/WinFsp.
- **TLS management:** ACME integration for public endpoints; self-signed generator for LAN use.

## 10. Testing & QA Strategy

- Unit tests for mount manager, config parser, credential vault.
- Integration tests per provider using mocked endpoints or official sandboxes.
- End-to-end smoke tests that mount and validate filesystem semantics using temp directories.
- Performance benchmarks measuring throughput, cache hit rates, reconnection times.
- Static analysis: golangci-lint, SAST for credential handling.

## 11. Roadmap & Milestones

1. **M0 - Foundations (2 weeks):**
   - Bootstrap repo, choose tooling, set up CI/CD (GitHub Actions) and goreleaser pipeline.
   - Implement config persistence, credential vault, pluggable connector interface.
2. **M1 - Core Providers (4 weeks):**
   - FTP/SFTP, S3/MinIO, Backblaze, WebDAV connectors.
   - CLI workflows for config/mount, daemon service with gRPC/REST API.
   - Basic caching and reconnection logic.
3. **M2 - OAuth-heavy Providers (4 weeks):**
   - Dropbox, Google Drive, OneDrive, Box, pCloud connectors with OAuth device flow.
   - Web UI for configuration and monitoring.
   - Metrics and logging enhancements.
4. **M3 - Advanced Providers & UX (3 weeks):**
   - Mega integration, advanced cache policies, mount profiles.
   - RBAC, multi-user support, audit logs.
   - Polished Web UI, optional desktop tray prototype.
5. **M4 - Hardening & Release (2 weeks):**
   - Security review, packaging, documentation, beta release program.

## 12. Risks & Mitigations

- **Cross-platform FUSE complexity:** Mitigate by leveraging battle-tested rclone mounting logic and WinFsp adapters; add integration test matrix in CI.
- **OAuth token management pitfalls:** Centralize through secure vault, refresh tokens proactively, allow external secret stores.
- **Provider API rate limits:** Implement adaptive throttling, exponential backoff, and cache metadata to reduce chatter.
- **Local cache growth:** Provide configurable quotas and eviction policies; add monitoring alerts.
- **User trust & security:** Offer transparency logs, optional remote wipe, and strong docs around credential handling.

## 13. Documentation & Developer Experience

- **Docs site:** mdBook or Docusaurus under `/docs` with tutorials, provider-specific guides, troubleshooting.
- **API references:** Auto-generate from protobuf definitions using `buf` or `grpc-gateway` swagger.
- **Examples:** Recipes for common mounts (e.g., MinIO in Kubernetes, S3 cross-account).
- **Contribution guide:** Coding standards, branching, testing requirements.

## 14. Next Steps

- Validate Go + rclone library licensing alignment (rclone is MIT; ensure compliance).
- Prototype core daemon skeleton with one provider (S3) to verify mount performance.
- Design protobuf/gRPC contracts and generate stubs.
- Draft UX wireframes for Web UI and CLI flows.
- Set up CI pipeline skeleton (lint, tests, cross-build) and container packaging.
