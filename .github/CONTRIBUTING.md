# Contributing to CloudMoor

Thanks for helping build CloudMoor! This guide explains how to plan work, set up your environment, and submit high-quality changes that keep our milestone roadmap on track.

## Before You Start

- **Read the docs:** Familiarize yourself with the [project plan](../docs/plan.md), [engineering spec](../docs/spec.md), and the active [execution checklist](../docs/tasks.md). Every change should trace back to a roadmap item or ticket.
- **Follow the Code of Conduct:** We adhere to the [Contributor Covenant](CODE_OF_CONDUCT.md). Report concerns privately to `conduct@cloudmoor.dev`.
- **Security first:** Vulnerabilities should be disclosed via the private channel documented in [SECURITY.md](SECURITY.md). Do not open public issues for security problems.

## Planning Work

1. Choose a task from `docs/tasks.md` or open a new issue describing the problem and proposed scope.
2. Coordinate on the issue before starting large or ambiguous work. Mention which milestone and ticket ID (e.g., `TCK-001`) you are tackling.
3. Create a feature branch using the format described below and keep scope tight to ease reviews.

## Branching & Commit Guidelines

- **Branch naming:** `feat/<ticket-id>-short-slug`, `fix/<ticket-id>-slug`, or `docs/<ticket-id>-slug`. Example: `feat/TCK-003-connector-registry`.
- **Commits:** Prefer [Conventional Commits](https://www.conventionalcommits.org/) (`feat: add connector registry`). Include the ticket ID in the commit message body if it isn’t in the summary.
- **Traceability:** Reference the relevant checklist item or ticket in commit bodies and pull requests so updates flow back to `docs/tasks.md`.

## Environment Setup

| Tool            | Version | Notes                                                     |
| --------------- | ------- | --------------------------------------------------------- |
| Go              | 1.22.x  | Required for backend code. Install via `golang.org/dl`.   |
| Node.js         | 20.x    | Used by the React/Vite Web UI under `web/`.               |
| npm             | 10.x    | Ships with Node 20. Replace with pnpm/yarn only if agreed |
| Docker          | Latest  | Needed for Testcontainers-based integration suites.       |
| Make (optional) | Latest  | Make targets will arrive later in M0 (`make lint`, etc.). |

Clone the repository and install dependencies:

```bash
# Clone
git clone https://github.com/binGhzal/CloudMoor.git
cd CloudMoor

# Backend tools (run as needed)
go mod tidy

# Frontend tools
cd web
npm install
```

## Formatting, Linting & Tests

Until automated Make targets land, run the following commands locally before opening a pull request. These mirror the GitHub Actions lint and test workflows.

```bash
# From repo root
gofmt -w $(go list -f '{{.Dir}}' ./...)
gofumpt -w ./...
go vet ./...
golangci-lint run
go test ./...

# Frontend (if you touched web/)
cd web
npm run lint
npm test
```

> **Tip:** Configure your editor to run `gofmt`/`gofumpt` on save and respect settings from `.editorconfig`.

## Documentation Expectations

- Update `docs/spec.md`, `docs/plan.md`, or `docs/tasks.md` when your change alters standards, roadmap scope, or checklist status.
- Add or refresh inline documentation, CLI help, and API comments impacted by your change.
- For new features, consider adding a short note or guide under `docs/` so future contributors can follow your approach.

## Pull Request Checklist

Every pull request should:

1. Target the `main` branch and keep a narrow, reviewable scope.
2. Reference the ticket ID or checklist item in the PR title or description.
3. Include a summary of changes, testing performed, and screenshots/logs when applicable.
4. Confirm that linting, formatting, and tests above were executed (or explain why a step is skipped).
5. Ensure docs and task checkboxes are updated in the same PR when scope or status changes.
6. Request review from the relevant domain owners (Go, frontend, DevOps) when touching cross-cutting code.

CI will rerun the lint/test matrix on every pull request—keep it green to unblock reviewers.

## Release & Governance Notes

- Milestone exit reviews require clean checklists in `docs/tasks.md` and up-to-date documentation. Plan documentation updates alongside code.
- Major architectural or roadmap changes should start with an RFC or decision log entry under `docs/decisions/` (task TCK-009 introduces the template).
- Security and licensing considerations are tracked in `docs/plan.md` §6 and §9; loop in maintainers early if your change touches either area.

## Getting Help

- Open a discussion or issue for questions that benefit the wider community.
- Ping maintainers on the issue or PR when you need timely feedback.
- For urgent security matters, email `security@cloudmoor.dev` (see [SECURITY.md](SECURITY.md)).

Thank you for contributing! Together we’ll deliver a reliable, secure remote storage platform.
