# LLMProvider

Standalone Go module providing a unified interface for 43+ LLM providers with built-in fault tolerance.

**Module:** `digital.vasic.llmprovider`

## Features

- **43 provider adapters** -- OpenAI, Anthropic, Gemini, Groq, Mistral, DeepSeek, Ollama, and many more
- **Circuit breaker** -- Automatic fault isolation with closed/open/half-open state machine
- **Health monitoring** -- Periodic health checks with latency tracking and status history
- **Retry with backoff** -- Exponential backoff with jitter for transient failures
- **Model discovery** -- 3-tier dynamic model listing (provider API, models.dev, fallback)
- **Streaming support** -- Channel-based streaming responses
- **Tool/function calling** -- OpenAI-compatible tool use across supporting providers

## Quick Start

```go
import (
    "digital.vasic.llmprovider/pkg/models"
    "digital.vasic.llmprovider/pkg/providers/openai"
)

provider := openai.New("sk-...", "gpt-4o")

resp, err := provider.Complete(ctx, &models.LLMRequest{
    Prompt: "Explain circuit breakers in distributed systems",
    ModelParams: models.ModelParameters{
        Temperature: 0.7,
        MaxTokens:   1024,
    },
})
```

## Build

```bash
go build ./...
go test ./... -race -count=1
go vet ./...
```

## Documentation

- [Architecture](docs/ARCHITECTURE.md) -- Package structure, patterns, design decisions
- [Providers](docs/PROVIDERS.md) -- Complete list of supported providers with capabilities

## License

Apache-2.0
