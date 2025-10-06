# CloudMoor Engineering Standards

> Last updated: 2025-10-06

## 1. Purpose

This specification codifies the engineering practices, project layout, and quality bars CloudMoor follows across the Go backend, React + Tailwind frontend, documentation, and automation. It supplements the product roadmap in `docs/plan.md` and the execution backlog in `docs/tasks.md` by defining **how** we build, test, and document the system.

## 2. Scope & Guiding Principles

- Uphold **security-first** defaults for credential handling, transport security, and observability, as outlined in `docs/plan.md` §6.
- Embrace incremental delivery via the milestone plan while keeping the codebase **easy to navigate** and **safe to change**.
- Prefer **convention over configuration** and reuse standard ecosystem tooling to reduce onboarding friction.
- Every feature must ship with automated tests, documentation updates, and CI enforcement.

## 3. Repository Layout & Naming

CloudMoor adopts the community "Standard Go Project Layout" guidelines for backend code and a feature-oriented structure for the frontend.

### 3.1 Top-level folders

| Path        | Standard contents & notes                                                                                  |
| ----------- | ---------------------------------------------------------------------------------------------------------- |
| `cmd/`      | One package per executable (`cloudmoor`, `cloudmoord`, future tooling). Minimal `main.go` that wires deps. |
| `internal/` | Non-exported Go packages (core services, connectors, persistence).                                         |
| `pkg/`      | Reusable Go packages that may be imported by external tools (keep surface area intentional).               |
| `web/`      | React + Vite application (`src/`, `public/`, config). Tailwind and TypeScript live here.                   |
| `deploy/`   | Packaging artefacts (Docker, Helm, Terraform, service scripts).                                            |
| `docs/`     | Product plan, backlog, this spec, runbooks, provider guides, diagrams.                                     |
| `tools/`    | (Future) Developer tooling such as code generation helpers or pre-commit scripts.                          |
| `.github/`  | Workflows, issue templates, release configuration.                                                         |

Future directories (e.g., `build/`, `scripts/`) should be documented here once added.

### 3.2 Go package conventions

- Use **lowercase, no underscores** for package names (`mountmanager`, `credentials`).
- Group domain logic under `internal/<domain>` (e.g., `internal/mounts`, `internal/vault`). Shared primitives that may be imported by connectors go in `pkg/`.
- Keep executable wiring inside `cmd/<binary>/main.go`. Avoid business logic in `main` packages.
- Configuration files live under `configs/` (to be added during M0.1) and should never be imported as Go packages.

### 3.3 Frontend structure (`web/`)

Follow a hybrid feature + layer layout inspired by industry React TypeScript best practices:

```text
web/
├── src/
│   ├── app/              # App shell, router, providers, layout frames
│   ├── features/         # Feature folders (e.g., mounts, auth, settings)
│   │   └── <feature>/
│   │       ├── components/
│   │       ├── hooks/
│   │       ├── api/      # API client wrappers typed with generated SDKs
│   │       ├── routes.tsx
│   │       └── index.ts
│   ├── components/       # Truly shared UI elements (Button, Modal, Table)
│   ├── hooks/            # Reusable cross-feature hooks (useMediaQuery)
│   ├── lib/              # Utilities, formatting helpers
│   ├── services/         # Data-access helpers (REST/gRPC bridges)
│   ├── styles/           # Tailwind config extensions, global CSS
│   └── types/            # Global TypeScript types & API response interfaces
├── tests/                # UI integration tests (Cypress/Playwright) if not colocated
├── public/
├── package.json
└── vite.config.ts
```

Tests for components/hooks live alongside their subjects (`Component.test.tsx`, `useThing.test.ts`).

## 4. Go Backend Standards

### 4.1 Language & toolchain

- Target Go **1.22.x** (match `go.mod`). Maintain compatibility with one previous minor release in CI matrix.
- Enforce formatting with `gofmt` + `gofumpt`; imports via `goimports`.
- Static analysis through `golangci-lint` (enable `gosimple`, `staticcheck`, `gocyclo`, `misspell`, `errcheck`, `prealloc`, `revive`).
- Use `buf` for protobuf linting/builds.

### 4.2 Coding guidelines

- **Context first**: Pass `context.Context` as the first argument to functions that perform I/O or long-running work.
- **Error handling**: Return wrapped errors using `fmt.Errorf("...: %w", err)`; avoid using panics outside program start-up.
- **Logging**: Use `zap.Logger` (structured logging) from a shared `internal/logging` package. No `fmt.Println` in production code.
- **Configuration**: Centralize config parsing in `internal/config`, reading from environment + YAML, with validation using `go-playground/validator`.
- **Dependencies**: Inject via constructors; prefer interfaces close to the consumer (`internal/mounts.Manager`, `internal/vault.Store`).
- **Concurrency**: Guard shared state with channels or `sync` primitives; ensure goroutines respect context cancellation.

### 4.3 Testing guidelines

- Place tests in the same package with `_test.go` suffix. Use external test packages (`package foo_test`) when verifying public APIs.
- Default to **table-driven tests** with descriptive names. Skeleton:

  ```go
  func TestValidateMount(t *testing.T) {
      t.Parallel()

      cases := map[string]struct {
          input   mounts.Config
          wantErr string
      }{
          "valid config": {...},
          "missing bucket": {...},
      }

      for name, tc := range cases {
          tc := tc
          t.Run(name, func(t *testing.T) {
              t.Parallel()

              err := ValidateMount(tc.input)
              if tc.wantErr != "" {
                  require.ErrorContains(t, err, tc.wantErr)
                  return
              }
              require.NoError(t, err)
          })
      }
  }
  ```

- Prefer `testify/require` for assertions to abort on failure; `assert` is acceptable when subsequent checks add value.
- Cover error paths, edge cases (empty configs, large payloads, timeouts).
- Use `testing/iotest`, `httptest`, and fake connectors to avoid hitting real services in unit tests.
- Benchmarks belong in `_test.go` files using `testing.B` (e.g., cache eviction, mount lifecycle hot paths).

### 4.4 Integration tests

- Use `testcontainers-go` for provider simulations (MinIO, WebDAV, Dropbox mock).
- Keep integration suites under `internal/<domain>/integration_test.go` or `tests/integration/` if spanning packages.
- Guard long-running suites behind `go test -tags=integration` to control CI scope.

### 4.5 Code generation

- Store generated code under `internal/api` and mark files with `// Code generated by ... DO NOT EDIT.`
- Use `mage` or `make` targets (`make generate`) to run protobuf and mock generators (`mockery` for interfaces).

## 5. React + TypeScript Standards

### 5.1 Toolchain & configuration

- Vite + TypeScript in **strict** mode (`"strict": true`, `"noImplicitAny": true`).
- ESLint with `typescript-eslint`, `eslint-plugin-react`, `eslint-plugin-react-hooks`, `eslint-plugin-jsx-a11y`, and Tailwind recommended rules.
- Prettier formatting with repo-wide `.prettierrc` (2-space indent, single quotes, semicolons).
- Path aliases managed via `tsconfig.json` (`@/features`, `@/components`).

### 5.2 Component guidelines

- Use functional components with hooks; avoid class components.
- Derive UI from data – no side-effects inside render. Fetch data via hooks (React Query or custom) that live in `features/<name>/hooks/`.
- Co-locate component styles (Tailwind utility classes) or `*.module.css` when utilising CSS modules for complex cases.
- Ensure accessibility: labelled inputs, ARIA attributes when needed, maintain focus order.
- Use Storybook (to be introduced in M2) for visual regression and component documentation.

### 5.3 State management & data access

- Local state: React `useState` / `useReducer`.
- Shared cross-feature state: React Query (preferred) or Context. Global contexts declared under `src/app/providers`.
- API access: Generate OpenAPI client or gRPC-web stubs; wrap in thin service modules inside each feature.

### 5.4 Testing

- Unit/component tests with React Testing Library + Jest. Follow AAA (Arrange, Act, Assert).
- Mock network requests with `msw` (Mock Service Worker).
- Accessibility smoke tests using `@axe-core/react` or `jest-axe`.
- Snapshot tests only for stable, structural output (avoid brittle snapshots).
- End-to-end tests (Cypress or Playwright) per milestone M2.3; keep under `web/tests/`.

## 6. Cross-cutting Testing Strategy

| Layer                  | Backend tooling                              | Frontend tooling                                 | Notes                                                                              |
| ---------------------- | -------------------------------------------- | ------------------------------------------------ | ---------------------------------------------------------------------------------- |
| Formatting/linting     | `gofmt`, `gofumpt`, `golangci-lint`          | ESLint, Prettier                                 | Enforced in pre-commit and CI.                                                     |
| Unit                   | `go test ./...`, table-driven + testify      | Jest + React Testing Library                     | Aim for ≥70% coverage of mount manager, vault, connectors.                         |
| Integration            | `go test -tags=integration` (Testcontainers) | Cypress/Playwright (UI), API contract tests      | Nightly CI job; provide fixtures for deterministic results.                        |
| Accessibility          | `axe-core` via CLI (future)                  | `jest-axe`, manual keyboard audits               | Accessibility is part of Definition of Done for UI tasks.                          |
| Performance/benchmarks | `testing.B`, `benchstat`, k6 (API)           | Lighthouse CI (Web UI)                           | Baselines captured in `/docs/performance.md` (M4 deliverable).                     |
| Security               | `gosec`, `trivy`, dependency scanning        | `npm audit`, `depcheck`, ESLint security plugins | Fail builds on critical findings; document remediation in security review reports. |

## 7. Documentation Standards

- **Single source of truth**: `docs/tasks.md` (execution checklist), `docs/plan.md` (strategy), `docs/spec.md` (this document), plus domain-specific guides.
- Use [Mermaid](https://mermaid.js.org/) diagrams committed as text when possible. Exported images go in `docs/img/` with descriptive names.
- Every feature change must update related docs, including CLI help, API references, and provider guides.
- Maintain a decision log under `docs/decisions/` following the template introduced in Task M0.6.1.
- Run `markdownlint` locally prior to commits (configure rule overrides in `.markdownlint.json`).

## 8. Automation & Developer Experience

### 8.1 Makefile/Justfile targets (to be added in M0.1)

| Target           | Description                                                                |
| ---------------- | -------------------------------------------------------------------------- |
| `make bootstrap` | Install Go tools, Node dependencies, pre-commit hooks.                     |
| `make lint`      | Run gofmt, gofumpt, golangci-lint, eslint, prettier check.                 |
| `make test`      | Execute Go unit tests and Jest suite.                                      |
| `make test-int`  | Execute integration suites (Go + Cypress) with necessary environment vars. |
| `make generate`  | Run protobuf/gRPC/GraphQL (future) code generation.                        |
| `make build`     | Compile Go binaries and web bundle (production).                           |

### 8.2 Pre-commit hooks

Adopt `lefthook` or `pre-commit` to run formatting, linting, and unit tests on staged files. Enforce gofmt/gofumpt, golangci-lint, eslint, prettier, and jest (focused).

### 8.3 Git workflow

- Branch naming: `feat/<ticket-id>-short-slug`, `fix/<ticket-id>-slug`, `docs/<ticket-id>-slug`.
- Conventional commits strongly encouraged (`type(scope): summary`).
- Pull requests must:
  - Reference ticket(s) from `docs/tasks.md`.
  - Include checklist verifying tests, lint, docs.
  - Attach screenshots for UI changes; paste metrics for performance updates.

## 9. CI/CD Expectations

- GitHub Actions workflows: `lint.yml` (Go + Node lint), `test.yml` (Go test matrix Linux/macos/windows + Jest), `integration.yml` (nightly heavy suites), `release.yaml` (GoReleaser + SBOM + cosign).
- Mandatory checks: formatting, go test, golangci-lint, eslint, jest.
- Optional but recommended: coverage upload (`codecov`), dependency diff alerts.
- Releases must pass security scan jobs (`trivy`, `gosec`) before artifacts publish.

## 10. Security & Secrets Handling

- Follow vault design in `docs/plan.md` §6. Key handling lives in `internal/vault` with AES-GCM envelope encryption.
- No secrets in git. Use `.env.example` for developer defaults; instruct use of OS keychains or `.envrc` with direnv.
- Ensure TLS is enforced in all daemon network listeners; self-signed certificates only for local dev with explicit flag.
- Keep third-party dependency list current via Dependabot and `go list -m -json all` snapshots stored in release artifacts.

## 11. Release Readiness Checklist

Before exiting a milestone, confirm:

1. All relevant `docs/tasks.md` checkboxes for the milestone are marked complete.
2. Code passes lint/test pipelines across supported OSes/architectures.
3. Documentation (plan/spec/task/backlog, provider guides) reflects new capabilities.
4. Security review items are captured in `/docs/security/`.
5. Decision logs updated for major architectural choices.

## 12. References & Further Reading

- Standard Go Project Layout (golang-standards/project-layout) — repository structure inspiration.
- Dave Cheney, "Prefer table-driven tests" — canonical guidance on Go testing patterns.
- Golangci-lint documentation — recommended linter suite and configuration examples.
- React TypeScript project structure guides (Medium, 2024) — feature-oriented folder layout with hooks/services/types directories.
- React Testing Library + jest-axe articles — accessibility-focused component testing strategies.

All future updates to this spec must include a changelog entry in `docs/plan.md` §13 and reference the relevant governance decision if standards change materially.
