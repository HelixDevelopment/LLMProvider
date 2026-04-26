# CLAUDE.md - LLMProvider Module


## Definition of Done

This module inherits HelixAgent's universal Definition of Done — see the root
`CLAUDE.md` and `docs/development/definition-of-done.md`. In one line: **no
task is done without pasted output from a real run of the real system in the
same session as the change.** Coverage and green suites are not evidence.

### Acceptance demo for this module

```bash
# Circuit breaker + health monitor + retry policy for provider fault tolerance
cd LLMProvider && GOMAXPROCS=2 nice -n 19 go test -count=1 -race -v \
  -run 'TestDefaultCircuitBreakerConfig|TestHealthMonitor_|TestDefaultRetryConfig' ./pkg/...
```
Expect: PASS; breaker opens after 3 consecutive failures, recovers after cooldown. `LLMProvider/README.md` shows the full `LLMProvider` interface.


## Overview

`digital.vasic.llmprovider` is a generic, reusable Go module providing LLM provider abstractions and utilities. It defines the core `LLMProvider` interface and common patterns for building LLM provider implementations, including circuit breakers, health monitoring, retry logic, and lazy loading. The module is designed for AI/LLM applications that need to integrate multiple LLM providers with fault tolerance and observability.

**Module**: `digital.vasic.llmprovider` (Go 1.25+)
**Dependencies**: `digital.vasic.models`, `github.com/sirupsen/logrus`
**Test Dependencies**: `github.com/stretchr/testify`

## Build & Test

```bash
go build ./...
go test ./... -count=1 -race
go test ./... -short              # Unit tests only
```

## Code Style

- Standard Go conventions, `gofmt` formatting
- Imports grouped: stdlib, third-party, internal (blank line separated)
- Line length ≤ 100 characters
- Naming: `camelCase` private, `PascalCase` exported, acronyms all-caps
- Errors: always check, wrap with `fmt.Errorf("...: %w", err)`
- Tests: table-driven, `testify`, naming `Test<Struct>_<Method>_<Scenario>`

## Package Structure

| Package | Purpose |
|---------|---------|
| `llmprovider` (root) | Core types: `LLMProvider` interface, circuit breaker, health monitor, retry config, lazy provider, and associated utilities |

## Key Interfaces

- `LLMProvider`: Interface for LLM provider implementations with `Complete`, `CompleteStream`, `HealthCheck`, `GetCapabilities`, `ValidateConfig`
- `CircuitBreaker`: Wraps an `LLMProvider` with fault tolerance (closed/open/half-open states)
- `HealthMonitor`: Tracks provider health with configurable thresholds and intervals
- `RetryConfig`: Configurable retry logic with exponential backoff and jitter
- `LazyProvider`: Lazy initialization of providers with optional event publishing

## Core Components

### LLMProvider Interface

The foundational interface that all LLM provider implementations must satisfy:

```go
type LLMProvider interface {
    Complete(ctx context.Context, req *models.LLMRequest) (*models.LLMResponse, error)
    CompleteStream(ctx context.Context, req *models.LLMRequest) (<-chan *models.LLMResponse, error)
    HealthCheck() error
    GetCapabilities() *models.ProviderCapabilities
    ValidateConfig(config map[string]interface{}) (bool, []string)
}
```

### Circuit Breaker

Prevents cascading failures when providers are unhealthy:
- **Closed**: Normal operation, requests pass through
- **Open**: Provider is failing, requests are short-circuited
- **Half-Open**: Testing if provider has recovered

### Health Monitor

Tracks provider health with:
- Configurable check intervals and timeouts
- Consecutive failure/success thresholds
- Health status transitions (healthy, degraded, unhealthy, unknown)
- Listener support for health status changes

### Retry Logic

Configurable retry with:
- Exponential backoff with configurable multiplier
- Jitter to prevent thundering herd
- HTTP status code detection (429, 500, 502, 503, 504)
- Context cancellation support

### Lazy Provider

Lazy initialization pattern:
- Deferred provider initialization until first use
- Configurable timeout and retry attempts
- Optional event bus integration for provider lifecycle events

## Dependencies

- **digital.vasic.models**: For `LLMRequest`, `LLMResponse`, `ProviderCapabilities` types
- **github.com/sirupsen/logrus**: For structured logging in circuit breaker
- **Standard library**: `context`, `sync`, `time`, `net/http`, etc.

## Thread Safety

- `CircuitBreaker`, `HealthMonitor`, and `CircuitBreakerManager` are thread-safe using `sync.RWMutex`
- `RetryConfig` is immutable after creation
- `LazyProvider` is thread-safe for concurrent initialization
- All exported methods are safe for concurrent use unless otherwise documented

## Example Usage

```go
import (
    "context"
    "digital.vasic.llmprovider"
    "digital.vasic.models"
)

func main() {
    provider := // create your provider implementation
    cb := llmprovider.NewDefaultCircuitBreaker("my-provider", provider)
    
    req := &models.LLMRequest{
        Prompt: "Hello, world!",
        MaxTokens: 100,
    }
    
    resp, err := cb.Complete(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(resp.Text)
}
```

## Integration with HelixAgent

This module is extracted from HelixAgent's `internal/llm` package. In HelixAgent, provider implementations (Claude, DeepSeek, Gemini, etc.) implement the `LLMProvider` interface and use these utilities for fault tolerance and observability.

## Integration Seams

| Direction | Sibling modules |
|-----------|-----------------|
| Upstream (this module imports) | Models |
| Downstream (these import this module) | DebateOrchestrator, HelixLLM |

*Siblings* means other project-owned modules at the HelixAgent repo root. The root HelixAgent app and external systems are not listed here — the list above is intentionally scoped to module-to-module seams, because drift *between* sibling modules is where the "tests pass, product broken" class of bug most often lives. See root `CLAUDE.md` for the rules that keep these seams contract-tested.



---

## Universal Mandatory Constraints

> Cascaded from the HelixAgent root `CLAUDE.md` via `/tmp/UNIVERSAL_MANDATORY_RULES.md`.
> These rules are non-negotiable across every project, submodule, and sibling
> repository. Project-specific addenda are welcome but cannot weaken or
> override these.

### Hard Stops (permanent, non-negotiable)

1. **NO CI/CD pipelines.** No `.github/workflows/`, `.gitlab-ci.yml`,
   `Jenkinsfile`, `.travis.yml`, `.circleci/`, or any automated pipeline.
   No Git hooks either. All builds and tests run manually or via
   Makefile/script targets.
2. **NO HTTPS for Git.** SSH URLs only (`git@github.com:…`,
   `git@gitlab.com:…`, etc.) for clones, fetches, pushes, and submodule
   updates. Including for public repos. SSH keys are configured on every
   service.
3. **NO manual container commands.** Container orchestration is owned by
   the project's binary/orchestrator (e.g. `make build` → `./bin/<app>`).
   Direct `docker`/`podman start|stop|rm` and `docker-compose up|down`
   are prohibited as workflows. The orchestrator reads its configured
   `.env` and brings up everything.

### Mandatory Development Standards

1. **100% Test Coverage.** Every component MUST have unit, integration,
   E2E, automation, security/penetration, and benchmark tests. No false
   positives. Mocks/stubs ONLY in unit tests; all other test types use
   real data and live services.
2. **Challenge Coverage.** Every component MUST have Challenge scripts
   (`./challenges/scripts/`) validating real-life use cases. No false
   success — validate actual behavior, not return codes.
3. **Real Data.** Beyond unit tests, all components MUST use actual API
   calls, real databases, live services. No simulated success. Fallback
   chains tested with actual failures.
4. **Health & Observability.** Every service MUST expose health
   endpoints. Circuit breakers for all external dependencies.
   Prometheus / OpenTelemetry integration where applicable.
5. **Documentation & Quality.** Update `CLAUDE.md`, `AGENTS.md`, and
   relevant docs alongside code changes. Pass language-appropriate
   format/lint/security gates. Conventional Commits:
   `<type>(<scope>): <description>`.
6. **Validation Before Release.** Pass the project's full validation
   suite (`make ci-validate-all`-equivalent) plus all challenges
   (`./challenges/scripts/run_all_challenges.sh`).
7. **No Mocks or Stubs in Production.** Mocks, stubs, fakes,
   placeholder classes, TODO implementations are STRICTLY FORBIDDEN in
   production code. All production code is fully functional with real
   integrations. Only unit tests may use mocks/stubs.
8. **Comprehensive Verification.** Every fix MUST be verified from all
   angles: runtime testing (actual HTTP requests / real CLI
   invocations), compile verification, code structure checks,
   dependency existence checks, backward compatibility, and no false
   positives in tests or challenges. Grep-only validation is NEVER
   sufficient.
9. **Resource Limits for Tests & Challenges (CRITICAL).** ALL test and
   challenge execution MUST be strictly limited to 30-40% of host
   system resources. Use `GOMAXPROCS=2`, `nice -n 19`, `ionice -c 3`,
   `-p 1` for `go test`. Container limits required. The host runs
   mission-critical processes — exceeding limits causes system crashes.
10. **Bugfix Documentation.** All bug fixes MUST be documented in
    `docs/issues/fixed/BUGFIXES.md` (or the project's equivalent) with
    root cause analysis, affected files, fix description, and a link to
    the verification test/challenge.
11. **Real Infrastructure for All Non-Unit Tests.** Mocks/fakes/stubs/
    placeholders MAY be used ONLY in unit tests (files ending
    `_test.go` run under `go test -short`, equivalent for other
    languages). ALL other test types — integration, E2E, functional,
    security, stress, chaos, challenge, benchmark, runtime
    verification — MUST execute against the REAL running system with
    REAL containers, REAL databases, REAL services, and REAL HTTP
    calls. Non-unit tests that cannot connect to real services MUST
    skip (not fail).
12. **Reproduction-Before-Fix (CONST-032 — MANDATORY).** Every reported
    error, defect, or unexpected behavior MUST be reproduced by a
    Challenge script BEFORE any fix is attempted. Sequence:
    (1) Write the Challenge first. (2) Run it; confirm fail (it
    reproduces the bug). (3) Then write the fix. (4) Re-run; confirm
    pass. (5) Commit Challenge + fix together. The Challenge becomes
    the regression guard for that bug forever.
13. **Concurrent-Safe Containers (Go-specific, where applicable).** Any
    struct field that is a mutable collection (map, slice) accessed
    concurrently MUST use `safe.Store[K,V]` / `safe.Slice[T]` from
    `digital.vasic.concurrency/pkg/safe` (or the project's equivalent
    primitives). Bare `sync.Mutex + map/slice` combinations are
    prohibited for new code.

### Definition of Done (universal)

A change is NOT done because code compiles and tests pass. "Done"
requires pasted terminal output from a real run, produced in the same
session as the change.

- **No self-certification.** Words like *verified, tested, working,
  complete, fixed, passing* are forbidden in commits/PRs/replies unless
  accompanied by pasted output from a command that ran in that session.
- **Demo before code.** Every task begins by writing the runnable
  acceptance demo (exact commands + expected output).
- **Real system, every time.** Demos run against real artifacts.
- **Skips are loud.** `t.Skip` / `@Ignore` / `xit` / `describe.skip`
  without a trailing `SKIP-OK: #<ticket>` comment break validation.
- **Evidence in the PR.** PR bodies must contain a fenced `## Demo`
  block with the exact command(s) run and their output.
