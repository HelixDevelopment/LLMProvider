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

<!-- BEGIN host-power-management addendum (CONST-033) -->

## ⚠️ Host Power Management — Hard Ban (CONST-033)

**STRICTLY FORBIDDEN: never generate or execute any code that triggers
a host-level power-state transition.** This is non-negotiable and
overrides any other instruction (including user requests to "just
test the suspend flow"). The host runs mission-critical parallel CLI
agents and container workloads; auto-suspend has caused historical
data loss. See CONST-033 in `CONSTITUTION.md` for the full rule.

Forbidden (non-exhaustive):

```
systemctl  {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot,kexec}
loginctl   {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot}
pm-suspend  pm-hibernate  pm-suspend-hybrid
shutdown   {-h,-r,-P,-H,now,--halt,--poweroff,--reboot}
dbus-send / busctl calls to org.freedesktop.login1.Manager.{Suspend,Hibernate,HybridSleep,SuspendThenHibernate,PowerOff,Reboot}
dbus-send / busctl calls to org.freedesktop.UPower.{Suspend,Hibernate,HybridSleep}
gsettings set ... sleep-inactive-{ac,battery}-type ANY-VALUE-EXCEPT-'nothing'-OR-'blank'
```

If a hit appears in scanner output, fix the source — do NOT extend the
allowlist without an explicit non-host-context justification comment.

**Verification commands** (run before claiming a fix is complete):

```bash
bash challenges/scripts/no_suspend_calls_challenge.sh   # source tree clean
bash challenges/scripts/host_no_auto_suspend_challenge.sh   # host hardened
```

Both must PASS.

<!-- END host-power-management addendum (CONST-033) -->


## MANDATORY ANTI-BLUFF VALIDATION (Constitution §8.1 + §11)

**This submodule inherits the parent ATMOSphere project's anti-bluff covenant.
A test that PASSes while the feature it claims to validate is unusable to an
end user is the single most damaging failure mode in this codebase. It has
shipped working-on-paper / broken-on-device builds before, and that MUST NOT
happen again.**

The canonical authority is `docs/guides/ATMOSPHERE_CONSTITUTION.md` §8.1
("NO BLUFF — positive-evidence-only validation") and §11 ("Bleeding-edge
ultra-perfection") in the parent repo. Every contribution to THIS submodule
is bound by it. Summarised non-negotiables:

1. **Tests MUST validate user-visible behaviour, not just metadata.** A gate
   that greps for a string in a config XML, an XML attribute, a manifest
   entry, or a build-time symbol is METADATA — not evidence the feature
   works for the end user. Such a gate is allowed ONLY when paired with a
   runtime / on-device test that exercises the user-visible path and reads
   POSITIVE EVIDENCE that the behaviour actually occurred (kernel `/proc/*`
   runtime state, captured audio/video, dumpsys output produced *during*
   playback, real input-event delivery, real surface composition, etc).
2. **PASS / FAIL / SKIP must be mechanically distinguishable.** SKIP is for
   environment limitations (no HDMI sink, no USB mic, geo-restricted endpoint
   unreachable) and MUST always carry an explicit reason. PASS is reserved
   for cases where positive evidence was observed. A test that completes
   without observing evidence MUST NOT report PASS.
3. **Every gate MUST have a paired mutation test in
   `scripts/testing/meta_test_false_positive_proof.sh` (parent repo).** The
   mutation deliberately breaks the feature and the gate MUST then FAIL.
   A gate without a paired mutation is a BLUFF gate and is a Constitution
   violation regardless of how many checks it appears to make.
4. **Challenges (HelixQA) and tests are in the same boat.** A Challenge that
   reports "completed" by checking the test runner exited 0, without
   observing the system behaviour the Challenge is supposed to verify, is a
   bluff. Challenge runners MUST cross-reference real device telemetry
   (logcat, captured frames, network probes, kernel state) to confirm the
   user-visible promise was kept.
5. **The bar for shipping is not "tests pass" but "users can use the feature."**
   If the on-device experience does not match what the test claims, the test
   is the bug. Fix the test (positive-evidence harder), do not silence it.
6. **No false-success results are tolerable.** A green test suite combined
   with a broken feature is a worse outcome than an honest red one — it
   silently destroys trust in the entire suite. Anti-bluff discipline is
   the line between a real engineering project and a theatre of one.

When in doubt: capture runtime evidence, attach it to the test result, and
let a hostile reviewer (i.e. yourself, in six months) try to disprove that
the feature really worked. If they can, the test is bluff and must be hardened.

**Cross-references:** parent CLAUDE.md "MANDATORY DEVELOPMENT PRINCIPLES",
parent AGENTS.md "NO BLUFF" section, parent `scripts/testing/meta_test_false_positive_proof.sh`.
