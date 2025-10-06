# CloudMoor

_Unified remote storage mounting across every platform._

> **Project status:** Pre-alpha (Milestone M0 in progress). Expect active development and frequent changes.

## Table of Contents

- [CloudMoor](#cloudmoor)
  - [Table of Contents](#table-of-contents)
  - [Why CloudMoor?](#why-cloudmoor)
  - [Key Features](#key-features)
  - [Architecture at a Glance](#architecture-at-a-glance)
  - [Project Structure](#project-structure)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Clone the Repository](#clone-the-repository)
    - [Bootstrap Tooling (upcoming in M0)](#bootstrap-tooling-upcoming-in-m0)
    - [Running the Daemon \& CLI (planned)](#running-the-daemon--cli-planned)
    - [Web UI Development (planned)](#web-ui-development-planned)
    - [Running Tests (planned)](#running-tests-planned)
  - [Operational Workflows](#operational-workflows)
  - [Roadmap](#roadmap)
  - [Documentation \& Backlog](#documentation--backlog)
  - [Contributing](#contributing)
  - [Security](#security)
  - [Governance](#governance)
  - [License](#license)
  - [Acknowledgements](#acknowledgements)

## Why CloudMoor?

CloudMoor is a Go-based platform that mounts remote storage providers (S3/MinIO, WebDAV, Dropbox—with a backlog of additional providers) as native drives on macOS, Linux, and Windows. It builds on the battle-tested rclone ecosystem while adding:

- A daemon with resilient reconnects, caching, and observability built in.
- A CLI for scripting and headless automation.
- A Web UI for secure multi-user administration.
- A future desktop tray app for quick controls.

## Key Features

- **Connector ecosystem:** Launch connectors include S3/MinIO, WebDAV, and Dropbox, with deferred backlog items such as FTP/SFTP, Backblaze B2, Google Drive, OneDrive, Box, pCloud, and more.
- **Smart caching:** Layered metadata and file caching with offline write-back support.
- **Security-first design:** AES-GCM credential vault, audit logging, TLS everywhere, RBAC-ready APIs.
- **Unified APIs:** gRPC core with REST gateway, OpenAPI generation, and CLI/Web UI parity.
- **Observability:** Structured logging with Zap, Prometheus metrics, webhook hooks, optional Grafana dashboards.
- **Cross-platform builds:** Single binary per OS, optional Docker/Helm packaging, and service install scripts for systemd, launchd, and Windows Service.

## Architecture at a Glance

CloudMoor separates the **control plane** (daemon, APIs, scheduler, credential vault, Web UI) from the **data plane** (rclone-powered VFS processes handling file IO). A high-level overview:

- `cloudmoord` daemon orchestrates mounts, retries, scheduling, and observability.
- `cloudmoor` CLI manages configuration, mount lifecycle, and service installation.
- Connector plugins translate provider-specific configuration into a common interface.
- SQLite persists configuration, credentials (encrypted), mount state, and audit events.
- Prometheus metrics, structured logs, and alert hooks power operations.
- Optional Web UI (React + Vite + Tailwind) provides dashboards, wizards, and RBAC.

Refer to [`docs/plan.md`](docs/plan.md) for the full architecture deep dive, component responsibility matrix, and security model.

## Project Structure

| Path            | Purpose                                                                           |
| --------------- | --------------------------------------------------------------------------------- |
| `cmd/`          | Go `main` packages for executables (`cloudmoor`, `cloudmoord`, supporting tools). |
| `internal/`     | Core domain packages (mount manager, connectors, vault, persistence, services).   |
| `pkg/`          | Intentional public Go packages for limited external reuse.                        |
| `web/`          | React + Vite console application scaffold.                                        |
| `deploy/`       | Packaging artefacts (Docker, Helm, installers, IaC).                              |
| `docs/`         | Strategic plan, specs, backlog, and operational runbooks.                         |
| `.editorconfig` | Shared formatting rules for Go, TypeScript, Markdown, and Makefiles.              |
| `.gitignore`    | Ignore rules covering Go builds, frontend artefacts, tooling caches, and secrets. |
| `LICENSE`       | Project licensing (MIT).                                                          |

Refer to [`docs/spec.md`](docs/spec.md#31-top-level-folders) for detailed layout conventions and update this table as new folders land.

## Getting Started

> **Note:** Code scaffolding and build scripts are being laid down during Milestone M0. This section outlines the intended workflow once the bootstrap tasks land.

### Prerequisites

- Go 1.22+
- Node.js 20+ (for the React-based Web UI)
- Make (optional but recommended)
- Docker (required for integration tests with Testcontainers)

### Clone the Repository

```bash
git clone https://github.com/binGhzal/CloudMoor.git
cd CloudMoor
```

### Bootstrap Tooling (upcoming in M0)

Once `Task M0.1` completes, the repository will include:

- `make lint` and `make test` helpers mirroring CI.
- `go mod tidy` and `npm install` scripts for backend/frontend dependencies.
- Pre-commit hooks for formatting and linting.

Until then, keep an eye on the checkbox progress in [`docs/tasks.md`](docs/tasks.md).

### Running the Daemon & CLI (planned)

```bash
# Build CLI and daemon binaries (placeholder commands)
make build

# Start daemon locally
./bin/cloudmoord --config ./config/cloudmoor.yaml

# Configure and mount a provider via CLI
./bin/cloudmoor config create --provider s3
./bin/cloudmoor mount add my-s3-mount
./bin/cloudmoor mount start my-s3-mount
```

These commands will be finalized as the CLI and daemon land in Milestones M0–M1.

### Web UI Development (planned)

```bash
cd web
npm install
npm run dev
```

The daemon will proxy Web UI assets in production, while `npm run dev` provides a Vite-powered local workflow.

### Running Tests (planned)

```bash
# Unit tests
make test

# Integration tests (spins up MinIO (S3) and WebDAV containers)
make test-integration
```

Refer to upcoming docs under `docs/testing/` for detailed guidance once the harness is in place (Task M1.5).

## Operational Workflows

- **Credential vault management:** `cloudmoor config vault test` verifies encryption setup and key rotation (Task M0.3.2).
- **Mount monitoring:** Prometheus metrics and structured logs (Tasks M1.3 & M2.4) feed dashboards and alerts.
- **Cache controls:** CLI/Web UI expose cache tuning, offline mode, and purge commands (Task M3.2).
- **Runbooks:** Operational procedures (mount failures, credential rotation, DR) will live under `docs/operations/` (Task M3.6).

## Roadmap

| Milestone                         | Focus                                                                                                      | Duration |
| --------------------------------- | ---------------------------------------------------------------------------------------------------------- | -------- |
| **M0 – Foundations**              | Repo scaffold, CI, connector interfaces, credential vault, persistence, proto contracts, librclone spike.  | 3 weeks  |
| **M1 – Core Providers & Daemon**  | S3/MinIO & WebDAV connectors, mount manager, CLI workflows, Testcontainers harness, daemon runtime.        | 6 weeks  |
| **M2 – OAuth Providers & Web UI** | Dropbox connector, OAuth device flow service, Web UI foundation, observability enhancements.               | 4 weeks  |
| **M3 – Advanced Providers & UX**  | Mega connector evaluation, advanced caching, RBAC, UX polish, desktop tray prototype, operations runbooks. | 4 weeks  |
| **M4 – Hardening & Release**      | Security review, packaging, documentation, performance benchmarks, beta program launch.                    | 3 weeks  |

Detailed acceptance criteria and dependencies are tracked in [`docs/tasks.md`](docs/tasks.md#ticket-backlog-detailed-acceptance-criteria).

## Documentation & Backlog

- 📘 [Strategic Plan](docs/plan.md): Product vision, architecture, KPIs, staffing, governance.
- ✅ [Execution TODOs + Ticket Backlog](docs/tasks.md): One consolidated checklist with hints, governance notes, and detailed acceptance criteria.
- 🧭 [Engineering Standards](docs/spec.md): Repository layout, coding conventions, testing strategy, and automation expectations.

When scope changes, update the relevant document and log decisions per the governance rules below.

## Contributing

CloudMoor is pre-release, but we welcome early feedback and contributions:

1. Review the roadmap and open tickets to identify work in flight.
2. Coordinate via issues or discussions before large changes.
3. Follow the upcoming `CONTRIBUTING.md` (Task M0.1.2) for branching, linting, and testing expectations.
4. Reference ticket IDs in commits to maintain plan → task traceability.

> **Note:** Security-sensitive disclosures should go directly to the maintainers (see [Security](#security)).

## Security

- Report vulnerabilities privately to `security@cloudmoor.dev` (placeholder). Avoid filing public issues for sensitive reports.
- Credentials are always encrypted at rest; see the [Security section of the plan](docs/plan.md#6-security-considerations) for threat model details.
- Security review gates are enforced at each milestone exit (see [Governance](#governance)).

## Governance

CloudMoor follows the delivery governance process described in [`docs/plan.md#21-delivery-governance--tracking`](docs/plan.md#21-delivery-governance--tracking):

- Weekly steering sync to review burn-down, risks, and blockers.
- Milestone exit reviews covering scope, defects, security, and documentation.
- Change control via lightweight RFCs; updates reflected in plan, tickets, and tasks within 24 hours.
- Decision log stored under `docs/decisions/` (to be created) for cross-team traceability.

## License

CloudMoor is released under the [MIT License](LICENSE), aligning with rclone’s permissive requirements while enabling broad community adoption.

## Acknowledgements

- [rclone](https://rclone.org/) for the foundational remote storage backends.
- [bazil.org/fuse](https://github.com/bazil/fuse) and [WinFsp](https://github.com/billziss-gh/winfsp) for cross-platform filesystem support.
- Open-source communities behind React, Vite, Tailwind, Prometheus, Zap, and Testcontainers.

---

**Need something that isn’t documented yet?** Open an issue or start a discussion so we can fold it into the plan, tickets, and task backlog.
