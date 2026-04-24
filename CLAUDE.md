## MANDATORY: No CI/CD Pipelines

**NO GitHub Actions, GitLab CI/CD, or any automated pipeline may exist in this repository!**

- No `.github/workflows/` directory
- No `.gitlab-ci.yml` file
- No Jenkinsfile, .travis.yml, .circleci, or any other CI configuration
- All builds and tests are run manually or via Makefile targets
- This rule is permanent and non-negotiable

## Project Overview

**LLMProvider** is a standalone Go module providing a shared LLM provider interface, 40+ provider adapters, retry logic, circuit breaker, and health monitoring.

**Module:** `digital.vasic.llmprovider`

## Build Commands

```bash
# Build
go build ./...

# Run all tests
go test ./... -race -count=1

# Run core tests only (no network calls)
go test ./pkg/models/... ./pkg/retry/... ./pkg/circuit/... ./pkg/health/... ./pkg/provider/... -race -count=1

# Run specific provider tests
go test ./pkg/providers/claude/... -race -count=1

# Vet
go vet ./...
```

## Architecture

```
pkg/
  provider/    - LLMProvider interface
  models/      - LLMRequest, LLMResponse, ProviderCapabilities
  retry/       - RetryConfig, ExecuteWithRetry, backoff
  circuit/     - CircuitBreaker, CircuitBreakerManager
  health/      - HealthMonitor, ProviderHealth
  http/        - HTTP client with retry
  discovery/   - 3-tier model discovery (API, models.dev, fallback)
  providers/   - 40+ provider implementations
```

## Key Patterns

- All providers implement `provider.LLMProvider` interface
- Circuit breaker wraps providers for fault tolerance
- Health monitor tracks provider availability
- Retry logic with exponential backoff and jitter
- Discovery caches model lists with configurable TTL

## Dependencies

- `github.com/sirupsen/logrus` - Logging
- `github.com/stretchr/testify` - Testing
- Standard library for everything else

## Environment Variables

Provider API keys are loaded from `.env` file. See `.env.example`.

## Definition of Done

A change is NOT done because code compiles and tests pass. "Done" requires pasted
terminal output from a real run of the real system, produced in the same session as
the change. Coverage and passing suites measure the LLM's model of the product, not
the product.

1. **No self-certification.** *Verified, tested, working, complete, fixed, passing*
   are forbidden in commits, PRs, and agent replies without accompanying pasted
   output from a same-session real-system run.
2. **Demo before code.** Every task begins with the runnable acceptance demo below.
3. **Real system.** Demos run against real artifacts — built binaries, live
   databases, instrumented devices — not mocks/stubs/in-memory fakes.
4. **Skips are loud.** `t.Skip` / `@Ignore` / `xit` / `it.skip` without a trailing
   `SKIP-OK: #<ticket>` annotation fails `make ci-validate-all`.
5. **Contract tests on every seam.** Any change touching a module↔module boundary
   runs one roundtrip test asserting the wire format on both sides.
6. **Evidence in the PR.** PR body contains a fenced `## Demo` block with exact
   command(s) + output.

### Acceptance demo for this module

```bash
# TODO — replace with a 10-line real-system demo. See examples in
# HelixAgent/docs/development/dod-dropin/templates/CLAUDE_md_clause.md
```
