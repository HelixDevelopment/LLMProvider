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
