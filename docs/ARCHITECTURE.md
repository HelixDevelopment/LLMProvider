# LLMProvider Architecture

## Overview

Standalone Go module (`digital.vasic.llmprovider`) providing a unified interface for 43+ LLM providers with retry logic, circuit breaker fault tolerance, health monitoring, and dynamic model discovery.

## Package Structure

```
pkg/
  provider/    -- LLMProvider interface (core abstraction)
  models/      -- LLMRequest, LLMResponse, ProviderCapabilities, Tool types
  providers/   -- 43 provider implementations (one package per provider)
  retry/       -- RetryConfig with exponential backoff and jitter
  circuit/     -- CircuitBreaker wrapping providers for fault tolerance
  health/      -- HealthMonitor tracking provider availability
  http/        -- HTTP client with retry
  discovery/   -- 3-tier dynamic model discovery
```

## Provider Adapter Pattern

All providers implement the `LLMProvider` interface:

```go
type LLMProvider interface {
    Complete(ctx context.Context, req *LLMRequest) (*LLMResponse, error)
    CompleteStream(ctx context.Context, req *LLMRequest) (<-chan *LLMResponse, error)
    HealthCheck() error
    GetCapabilities() *ProviderCapabilities
    ValidateConfig(config map[string]interface{}) (bool, []string)
}
```

Each provider translates the common `LLMRequest` into provider-specific API calls and normalizes the response back to `LLMResponse`. Capabilities advertise supported features: streaming, function calling, vision, tools, search, reasoning, code completion.

## Circuit Breaker

Three states: **Closed** (normal), **Open** (failing, requests short-circuited), **Half-Open** (testing recovery).

- `FailureThreshold`: 5 consecutive failures opens the circuit
- `SuccessThreshold`: 2 successes in half-open state closes it
- `Timeout`: 30 seconds before transitioning from open to half-open
- `HalfOpenMaxRequests`: 3 probe requests allowed while half-open

`CircuitBreakerManager` manages breakers for multiple providers. State change listeners can be registered for monitoring.

## Health Monitoring

`HealthMonitor` periodically checks all registered providers:

- **Check interval**: 30 seconds (configurable)
- **States**: Healthy, Unhealthy, Degraded, Unknown
- **Thresholds**: 2 consecutive successes to mark healthy, 3 failures for unhealthy
- Tracks latency, consecutive failures, total check/success/failure counts per provider

## 3-Tier Model Discovery

Dynamic model listing without hardcoded catalogs:

1. **Tier 1**: Query the provider's own API (`/v1/models` or equivalent)
2. **Tier 2**: Query `models.dev` API for the provider's model catalog
3. **Tier 3**: Fall back to hardcoded known models

Results are cached with configurable TTL (default: 1 hour). Custom `ModelFilter` and `ResponseParser` functions handle non-standard APIs.

## Retry Logic

`ExecuteWithRetry` wraps any operation with:
- Configurable max retries
- Exponential backoff with jitter
- Context cancellation support
- Retry-worthy error classification

## Key Design Decisions

- **One package per provider**: Each provider is isolated with its own HTTP client logic, model mappings, and API quirks.
- **Interface-first**: The `LLMProvider` interface allows consumers to swap providers without code changes.
- **Defensive resilience**: Circuit breaker + health monitor + retry ensures graceful degradation when providers go down.
- **Minimal dependencies**: Only `logrus` for logging and `testify` for testing. Everything else is standard library.
