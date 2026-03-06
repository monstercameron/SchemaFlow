# SchemaFlow

Production-ready LLM operations with Go's type safety.

SchemaFlow turns LLM calls into typed Go operations. You define structs and use semantic operations like `Extract`, `Transform`, `Generate`, `Rank`, `Enrich`, or `Verify` without manually juggling prompts, JSON parsing, retries, logging, and cost tracking in each call site.

## What It Is

SchemaFlow is a Go library, not an app.

It gives you:
- typed LLM operations with generics
- fluent request builders for common workflows
- provider abstraction for OpenAI, Anthropic, OpenRouter, Cerebras, and local mocks
- built-in retries, timeouts, structured logging, tracing, metrics, and cost tracking

## Quick Start

```bash
go get github.com/monstercameron/SchemaFlow
```

```bash
export SCHEMAFLOW_API_KEY=your-api-key
```

```go
package main

import (
    "fmt"

    schemaflow "github.com/monstercameron/SchemaFlow"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    schemaflow.Init("")

    person, err := schemaflow.Extracting[Person]("John is 30 years old").Run()
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", person)
}
```

## Fluent API

New integrations should prefer the fluent builders.

```go
person, err := schemaflow.Extracting[Person](rawText).
    Strict().
    Smart().
    Steer("Prefer explicit evidence over guesses").
    Run()

best, err := schemaflow.Choosing(products).
    By("lowest total cost", "best battery life").
    Fast().
    Run()

urgent, err := schemaflow.FilterBy(tasks, "high priority tasks due today")
```

The original function-based API remains supported:

```go
person, err := schemaflow.Extract[Person](rawText, schemaflow.NewExtractOptions())
summary, err := schemaflow.Summarize(text, schemaflow.NewSummarizeOptions())
result, err := schemaflow.Verify(claim, schemaflow.NewVerifyOptions())
```

## Core Operations

### Structured data
- `Extract[T]`
- `Transform[T, U]`
- `Generate[T]`
- `Infer[T]`
- `Parse[T]`
- `Enrich[T, U]`
- `Normalize[T]`
- `Project[T, U]`
- `Pivot[T, U]`

### Text and analysis
- `Summarize`
- `Rewrite`
- `Translate`
- `Expand`
- `Classify`
- `Score`
- `Compare`
- `Similar`
- `Critique`
- `Explain`
- `Verify`
- `Question`

### Collections and decisions
- `Choose`
- `Filter`
- `Sort`
- `Rank`
- `Cluster`
- `SemanticMatch`
- `MatchOne`
- `Decide`
- `Guard`
- `Resolve`
- `Negotiate`
- `Arbitrate`

See [docs/reference/API.md](docs/reference/API.md) for the full surface.

## Reliability Defaults

SchemaFlow now treats the shared LLM path as infrastructure, not a raw provider call.

Built in:
- request IDs generated automatically when missing
- provider-level retry policy support
- default retry budget for transient LLM failures
- fail-fast behavior for non-retryable request/auth errors
- validation for empty provider responses
- predictable timeout handling through context and `SCHEMAFLOW_TIMEOUT`

Retry-related environment variables:
- `SCHEMAFLOW_LLM_MAX_RETRIES`
- `SCHEMAFLOW_LLM_RETRY_BACKOFF`
- `SCHEMAFLOW_TIMEOUT`

Client-level tuning:

```go
client := schemaflow.NewClient(apiKey).
    WithRetries(3).
    WithRetryBackoff(500 * time.Millisecond).
    WithTimeout(30 * time.Second).
    WithProvider("openai")
```

## Logging And Review

SchemaFlow ships with structured logging backed by `slog`.

Features:
- text or JSON output
- level control
- optional file sink
- in-memory capture for review
- request-correlated LLM lifecycle logs

Environment variables:
- `SCHEMAFLOW_LOG_LEVEL=debug|info|warn|error`
- `SCHEMAFLOW_LOG_FORMAT=text|json`
- `SCHEMAFLOW_LOG_FILE=/path/to/schemaflow.log`
- `SCHEMAFLOW_LOG_BUFFER=1000`
- `SCHEMAFLOW_LOG_SOURCE=true`
- `SCHEMAFLOW_LOG_DISABLE_STDERR=true`
- `SCHEMAFLOW_LOG_DISABLE_CAPTURE=true`

Programmatic configuration:

```go
schemaflow.ConfigureLogging(schemaflow.LoggerConfig{
    Level:      schemaflow.LogDebug,
    Format:     "json",
    FilePath:   "schemaflow.log",
    BufferSize: 2000,
    Capture:    true,
})

entries := schemaflow.GetLogEntries()
fmt.Println("captured logs:", len(entries))
schemaflow.ResetLogEntries()
```

## Metrics And Cost Tracking

SchemaFlow records:
- request counts
- request durations
- prompt/completion/total tokens
- prompt/completion/total cost in USD
- cached and reasoning token/cost fields when providers return them

Low-cardinality metrics are recorded through the telemetry registry, while exact per-request cost history is tracked separately for review and reporting.

Examples:

```go
import (
    "fmt"
    "time"

    "github.com/monstercameron/SchemaFlow/pricing"
)

summary := pricing.GetCostSummary(time.Now().Add(-1*time.Hour), map[string]string{
    "provider": "openai",
})
fmt.Printf("requests=%d avg_cost=%.6f avg_tokens=%.1f\n",
    summary.RequestCount,
    summary.AverageCostPerRequest,
    summary.AverageTokensPerRequest,
)

record, ok := pricing.GetRequestCost("req-123")
if ok {
    fmt.Printf("request cost: %.6f total tokens: %d\n",
        record.Cost.TotalCost,
        record.TokenUsage.TotalTokens,
    )
}
```

## Provider Layer

```go
client := schemaflow.NewClient(apiKey)

client.WithProvider("openai")
client.WithProvider("anthropic")
client.WithProvider("openrouter")
client.WithProvider("cerebras")
client.WithProvider("deepseek")
client.WithProvider("qwen")
client.WithProvider("zai")
client.WithProvider("local")
```

Built-in provider notes:
- `anthropic` uses Anthropic's native Messages API
- `deepseek`, `qwen`, `zai`, `openrouter`, and `cerebras` use the shared OpenAI-compatible provider path
- provider-specific env vars are supported, for example `DEEPSEEK_API_KEY`, `DASHSCOPE_API_KEY`, `ZAI_API_KEY`, `ANTHROPIC_API_KEY`

Custom provider integration:

```go
schemaflow.RegisterProviderFactory("myvendor", func(cfg schemaflow.ProviderConfig) (schemaflow.Provider, error) {
    cfg.BaseURL = "https://vendor.example.com/v1"
    return schemaflow.NewOpenAICompatibleProvider("myvendor", cfg)
})

client := schemaflow.NewClient("").
    WithProviderConfig("myvendor", schemaflow.ProviderConfig{
        APIKey: "vendor-key",
    })
```

Default OpenAI intelligence mapping:
- `Smart -> gpt-5.4`
- `Fast -> gpt-5-mini`
- `Quick -> gpt-5-nano`

You can override model mapping with:
- `SCHEMAFLOW_MODEL`
- `SCHEMAFLOW_MODEL_SMART`
- `SCHEMAFLOW_MODEL_FAST`
- `SCHEMAFLOW_MODEL_QUICK`

## Example Workflow

```go
package main

import (
    "fmt"

    schemaflow "github.com/monstercameron/SchemaFlow"
)

type Ticket struct {
    Title    string `json:"title"`
    Priority string `json:"priority"`
    Team     string `json:"team"`
}

func ProcessTickets(raw []string) error {
    parsed := make([]Ticket, 0, len(raw))
    for _, item := range raw {
        ticket, err := schemaflow.Extracting[Ticket](item).
            Fast().
            Steer("Infer priority conservatively").
            Run()
        if err != nil {
            return err
        }
        parsed = append(parsed, ticket)
    }

    urgent, err := schemaflow.FilterBy(parsed, "priority is high and team is platform")
    if err != nil {
        return err
    }

    ordered, err := schemaflow.Sorting(urgent).
        By("most urgent operational risk first").
        Smart().
        Run()
    if err != nil {
        return err
    }

    fmt.Println("tickets to handle:", len(ordered))
    return nil
}
```

## Notes

- Prefer collection-aware APIs when they exist; avoid spraying raw LLM calls through business logic.
- Use the local provider for tests and smoke runs, not as proof of semantic correctness.
- Use JSON output only for operations that truly need structure; SchemaFlow now infers and enforces that automatically for shared LLM calls.

## Documentation

- [API reference](docs/reference/API.md)
- [`examples/`](examples/)
- [Production backlog](docs/notes/PRODUCTION_TODO.md)

## Why SchemaFlow

1. It is just Go with generics, not a separate DSL.
2. It centralizes reliability concerns instead of scattering prompt and parse glue.
3. It gives you typed outputs for semantic operations.
4. It includes observability and cost analysis primitives out of the box.
5. It keeps provider choice and intelligence mapping configurable.
