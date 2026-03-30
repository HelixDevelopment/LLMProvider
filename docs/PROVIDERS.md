# LLMProvider Supported Providers

43 provider implementations in `pkg/providers/`. Each implements the `LLMProvider` interface.

## Commercial Providers

| Provider | Package | Notes |
|----------|---------|-------|
| OpenAI | `openai` | GPT-4o, GPT-4, GPT-3.5 series. Streaming, function calling, vision. |
| Anthropic | `anthropic` | Claude 4 Opus/Sonnet, Claude 3.5 series. Streaming, tools, vision. |
| Claude | `claude` | Alternate Claude adapter with different auth flow. |
| Google Gemini | `gemini` | Gemini 2.0 Flash, Gemini Pro. Streaming, vision. |
| Mistral | `mistral` | Mistral Large, Medium, Small. Code-focused models. |
| Codestral | `codestral` | Mistral's code-generation endpoint. |
| Cohere | `cohere` | Command R+, Command R. RAG-optimized. |
| AI21 | `ai21` | Jamba models. |
| Perplexity | `perplexity` | Search-augmented generation. |
| xAI | `xai` | Grok models. |

## Cloud/Infrastructure Providers

| Provider | Package | Notes |
|----------|---------|-------|
| NVIDIA | `nvidia` | NIM endpoints. GPU-optimized inference. |
| Cloudflare | `cloudflare` | Workers AI. Edge inference. |
| Replicate | `replicate` | Model hosting platform. Async prediction API. |
| Modal | `modal` | Serverless GPU inference. |
| HuggingFace | `huggingface` | Inference API. Thousands of open models. |
| GitHub Models | `githubmodels` | GitHub-hosted model inference. |

## Inference Providers

| Provider | Package | Notes |
|----------|---------|-------|
| Groq | `groq` | Ultra-fast inference (LPU). Llama, Mixtral. |
| Together | `together` | Open model hosting. Llama, CodeLlama, Mixtral. |
| Fireworks | `fireworks` | Fast open model inference. |
| SambaNova | `sambanova` | Custom silicon inference. |
| Cerebras | `cerebras` | Wafer-scale inference. |
| Hyperbolic | `hyperbolic` | GPU marketplace inference. |
| Novita | `novita` | Multi-model inference API. |
| SiliconFlow | `siliconflow` | Chinese model inference platform. |

## Chinese/Asian Providers

| Provider | Package | Notes |
|----------|---------|-------|
| Qwen | `qwen` | Alibaba Qwen series. Vision support. |
| DeepSeek | `deepseek` | DeepSeek Coder, DeepSeek Chat. |
| Kimi | `kimi` | Moonshot AI. Long-context models. |
| Zhipu | `zhipu` | GLM-4 series. Chinese language focus. |

## Specialized Providers

| Provider | Package | Notes |
|----------|---------|-------|
| Ollama | `ollama` | Local model inference. Self-hosted. |
| OpenRouter | `openrouter` | Multi-provider routing. 100+ models. |
| NLP Cloud | `nlpcloud` | Hosted NLP models. |
| Upstage | `upstage` | Solar models. Document AI. |
| Venice | `venice` | Privacy-focused inference. |
| Sarvam | `sarvam` | Indian language models. |
| VulaVula | `vulavula` | African language models. |
| PublicAI | `publicai` | Community inference network. |
| Chutes | `chutes` | Decentralized GPU inference. |

## Agent Providers

| Provider | Package | Notes |
|----------|---------|-------|
| Junie | `junie` | JetBrains Junie coding agent. |
| Kilo | `kilo` | Kilo Code coding agent. |
| Nia | `nia` | Nia AI assistant. |
| Zai | `zai` | Zai coding assistant. |
| Zen | `zen` | Zen AI coding agent. |

## Generic Provider

| Provider | Package | Notes |
|----------|---------|-------|
| Generic | `generic` | OpenAI-compatible adapter for any `/v1/chat/completions` endpoint. |

## Capability Matrix

| Capability | Providers Supporting It |
|------------|------------------------|
| Streaming | OpenAI, Anthropic, Gemini, Groq, Mistral, DeepSeek, Together, Fireworks, Ollama, and most others |
| Function Calling | OpenAI, Anthropic, Gemini, Mistral, Groq |
| Vision | OpenAI (GPT-4o), Anthropic (Claude), Gemini, Qwen |
| Tool Use | OpenAI, Anthropic, Gemini |
| Search | Perplexity |
| Code Completion | Codestral, DeepSeek, Junie, Kilo |
