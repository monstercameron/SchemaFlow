# SchemaFlow

SchemaFlow is a Go library for typed LLM operations.

It gives you a single public API built around fluent request builders so application code stays readable while retries, structured output contracts, logging, metrics, and cost tracking stay centralized.

## Install

```bash
go get github.com/monstercameron/schemaflow
```

Set an API key for your provider. OpenAI is the default.

```bash
export SCHEMAFLOW_API_KEY=your-api-key
```

## Quick Start

```go
package main

import (
    "fmt"

    schemaflow "github.com/monstercameron/schemaflow"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    if err := schemaflow.InitWithEnv(); err != nil {
        panic(err)
    }

    person, err := schemaflow.Extracting[Person]("John is 30 years old").
        Strict().
        Run()
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", person)
}
```

## API Shape

The public API is the fluent builder stack at the root package.

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

summary, err := schemaflow.Summarizing(longText).
    MaxLength(120).
    Run()
```

Builder conventions:
- `Run()` executes the operation
- `WithOptions(...)` replaces the full typed option object
- `Configure(func(...))` lets you mutate the typed option object directly
- common controls such as `Strict()`, `Smart()`, `Fast()`, `Quick()`, `Steer(...)`, `Context(...)`, and `RequestID(...)` are exposed where they make sense

## Core Builders

### Structured extraction and generation
- `Extracting[T](input)`
- `Transforming[T, U](input)`
- `Generating[T](prompt)`
- `Inferring[T](input)`
- `Parsing[T](input)`
- `Enriching[T, U](input)`
- `EnrichingInPlace[T](input)`
- `Normalizing[T](input)`
- `NormalizingText(input)`
- `NormalizingBatch[T](items)`
- `Projecting[T, U](input)`
- `Pivoting[T, U](input)`

### Text operations
- `Summarizing(input)`
- `Rewriting(input)`
- `Translating(input)`
- `Expanding(input)`
- `Completing(input)`
- `CompletingField[T](input, fieldName)`
- `Redacting[T](input)`
- `LLMRedacting(input)`

### Analysis and validation
- `Classifying[T, C](input)`
- `Scoring[T](input)`
- `Comparing[T](left, right)`
- `CheckingSimilarity[T](left, right)`
- `Validating[T](input)`
- `Asking[T, A](input, question)`
- `Explaining(input)`
- `Verifying(input)`
- `VerifyingClaim(claim)`

### Collection and reasoning
- `Choosing[T](items)`
- `Filtering[T](items)`
- `Sorting[T](items)`
- `Ranking[T](items)`
- `Clustering[T](items)`
- `Matching[S, T](sources, targets)`
- `MatchingOne[S, T](source, targets)`
- `Annotating[T](input)`
- `Compressing[T](input)`
- `CompressingText(input)`
- `Decomposing[T](input)`
- `DecomposingInto[T, U](input)`
- `Critiquing[T](input)`
- `Synthesizing[T](sources)`
- `Predicting[T](input)`
- `Negotiating[T](constraints)`
- `NegotiatingAdversarially[T](context)`
- `Resolving[T](sources)`
- `Deriving[T, U](input)`
- `Conforming[T](input, standard)`
- `Interpolating[T](items)`
- `Arbitrating[T](options)`
- `Auditing[T](input)`
- `Assembling[T](parts)`

### Compact helpers
- `ChooseBy(items, criteria...)`
- `FilterBy(items, criteria)`
- `SortBy(items, criteria)`

Full API reference: [docs/reference/API.md](docs/reference/API.md)

## Usage Examples

### Extract typed data

```go
type Invoice struct {
    Number string  `json:"number"`
    Total  float64 `json:"total"`
}

invoice, err := schemaflow.Extracting[Invoice](rawEmail).
    Strict().
    Fast().
    Run()
```

### Transform one type into another

```go
type Lead struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CRMContact struct {
    FullName string `json:"full_name"`
    Email    string `json:"email"`
}

contact, err := schemaflow.Transforming[Lead, CRMContact](lead).
    Strict().
    Steer("Preserve exact identifiers and emails").
    Run()
```

### Generate typed output

```go
type ReleaseNote struct {
    Title   string   `json:"title"`
    Bullets []string `json:"bullets"`
}

note, err := schemaflow.Generating[ReleaseNote]("Write release notes for version 2.3").
    Creative().
    Smart().
    Run()
```

### Filter and sort collections

```go
urgent, err := schemaflow.Filtering(tasks).
    By("high priority tasks due today").
    Run()

ordered, err := schemaflow.Sorting(urgent).
    By("most urgent operational risk first").
    Smart().
    Run()
```

### Summarize and rewrite text

```go
summary, err := schemaflow.Summarizing(article).
    MaxLength(160).
    Run()

rewrite, err := schemaflow.Rewriting(summary).
    Tone("executive").
    Run()
```

### Validate structured data

```go
result, err := schemaflow.Validating(customer).
    Rules("email must be valid, country must be ISO alpha-2, age must be at least 18").
    Run()

if err != nil {
    panic(err)
}
if !result.Valid {
    fmt.Println(result.Issues)
}
```

### Ask typed questions over context

```go
answer, err := schemaflow.Asking[string, string](report, "What changed from last quarter?").
    Strict().
    Run()
```

## Reliability

SchemaFlow treats the shared LLM path as infrastructure.

Built in:
- automatic request IDs when missing
- retries for transient provider failures and empty completions
- fail-fast behavior for non-retryable auth and request errors
- JSON response enforcement for structured operations
- timeout control through context and client configuration
- structured error and request logging

Retry-related environment variables:
- `SCHEMAFLOW_LLM_MAX_RETRIES`
- `SCHEMAFLOW_LLM_RETRY_BACKOFF`
- `SCHEMAFLOW_TIMEOUT`

Client tuning:

```go
client := schemaflow.NewClient(apiKey).
    WithRetries(3).
    WithRetryBackoff(500 * time.Millisecond).
    WithTimeout(30 * time.Second).
    WithProvider("openai")
```

## Logging

SchemaFlow uses structured logging backed by `slog`.

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
- prompt, completion, cached, reasoning, and total tokens when available
- prompt, completion, cached, reasoning, and total USD cost when available

Per-request cost history is tracked separately from low-cardinality aggregate metrics.

```go
import (
    "fmt"
    "time"

    "github.com/monstercameron/schemaflow/pricing"
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

## Providers

Default provider:
- `openai`

Built-in providers:
- `openai`
- `anthropic`
- `openrouter`
- `cerebras`
- `deepseek`
- `qwen`
- `zai`
- `local`

```go
client := schemaflow.NewClient(apiKey)

client.WithProvider("openai")
client.WithProvider("anthropic")
client.WithProvider("deepseek")
client.WithProvider("qwen")
client.WithProvider("zai")
client.WithProvider("local")
```

Provider notes:
- `anthropic` uses the native Anthropic Messages API
- `deepseek`, `qwen`, `zai`, `openrouter`, and `cerebras` use the shared OpenAI-compatible provider path
- provider-specific env vars are supported, including `ANTHROPIC_API_KEY`, `DEEPSEEK_API_KEY`, `DASHSCOPE_API_KEY`, and `ZAI_API_KEY`

Custom provider registration:

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

Overrides:
- `SCHEMAFLOW_MODEL`
- `SCHEMAFLOW_MODEL_SMART`
- `SCHEMAFLOW_MODEL_FAST`
- `SCHEMAFLOW_MODEL_QUICK`

## Environment

Common environment variables:
- `SCHEMAFLOW_API_KEY`
- `OPENAI_API_KEY`
- `SCHEMAFLOW_PROVIDER`
- `SCHEMAFLOW_TIMEOUT`
- `SCHEMAFLOW_MODEL`
- `SCHEMAFLOW_MODEL_SMART`
- `SCHEMAFLOW_MODEL_FAST`
- `SCHEMAFLOW_MODEL_QUICK`

If `SCHEMAFLOW_API_KEY` is unset and `OPENAI_API_KEY` is present, SchemaFlow will use `OPENAI_API_KEY`.

## Design Notes

- The root package is the public facade for downstream Go projects.
- Builder implementation lives under `internal/` so the exported API can evolve without exposing internal plumbing.
- Prefer collection-aware builders when they exist instead of scattering raw single-call logic across application code.
- The local provider is useful for tests and smoke runs, not for proving semantic correctness.

## Documentation

- [API reference](docs/reference/API.md)
- [Examples](examples/)
- [Production backlog](docs/engineering/backlog/PRODUCTION_TODO.md)

## Compatibility

The older direct-call function API still exists for existing consumers, but it is compatibility-only. New code should use the fluent builders shown here.
