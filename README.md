# SchemaFlow 🔄

**Production-ready LLM operations with Go's type safety**

## What 🎯

SchemaFlow brings **compile-time type safety** to LLM operations. No more parsing JSON strings and hoping for the best. Define your types, and let the LLM fill them.

```go
// Your types, LLM's data extraction
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

person, err := schemaflow.Extract[Person]("John is 30")
// person is typed, validated, and ready to use
```

## ⚠️ Critical Rule: Never Put LLMs in Loops!

**LLMs have 1-3 second latency. A loop with 100 items = 5 minutes. Batch processing = 5 seconds.**

```go
// ❌ NEVER DO THIS - Will take 30+ minutes for 1000 items
for _, item := range items {
    result, _ := schemaflow.Extract[Data](item)  // 2 sec * 1000 = 33 minutes!
}

// ✅ ALWAYS DO THIS - Takes 30 seconds for 1000 items
results := schemaflow.ExtractBatch[Data](
    schemaflow.Batch().WithMode(schemaflow.ParallelMode),
    items,
)
```

**Think collections, not iterations. LLMs are high-latency, high-throughput systems.**

## Why It Makes Sense 💡

**LLMs are unstructured → Go is structured → SchemaFlow is the bridge**

- **🔒 Type Safety First**: Generics ensure compile-time safety for all LLM operations
- **🚀 Production Ready**: Built-in retries, timeouts, error handling, and circuit breakers
- **💰 Cost Optimized**: Batch operations reduce API costs by up to 90%
- **🔄 Provider Agnostic**: Swap between OpenAI, Anthropic, or local models with one line
- **📊 Observable**: OpenTelemetry tracing, structured logging, and cost tracking built-in
- **🎮 Testable**: Mock provider for unit tests - no API calls needed

## The Batch-First Mental Model 🧠

**Traditional code thinking doesn't work with LLMs. Here's why:**

### ❌ The Loop Trap (What NOT to Do)
```go
// THIS IS AN ANTI-PATTERN - Each iteration = 2-3 seconds
for _, customer := range customers {
    enriched, _ := schemaflow.Extract[EnrichedCustomer](customer)
    // 1000 customers × 2 seconds = 33 minutes of waiting!
}
```

### ✅ The Batch Solution (What TO DO)
```go
// Process ALL customers in one go
enriched := schemaflow.ExtractBatch[EnrichedCustomer](
    schemaflow.Batch().WithMode(schemaflow.MergedMode),
    customers,
)
// 1000 customers = 30 seconds total!
```

### Why This Matters

1. **LLMs are high-latency systems** (1-3 seconds per call)
2. **LLMs are high-throughput systems** (can process many items at once)
3. **API pricing favors batching** (same cost for 1 or 50 items in merged mode)

**Think of LLMs like a cargo ship**: It takes the same time to cross the ocean whether carrying 1 container or 1000. Don't send 1000 ships when 1 will do!

## LLM Operations 📦

### Core Operations
- **`Extract[T]`** - Pull structured data from unstructured text
- **`Transform[T,U]`** - Convert between types using LLM intelligence
- **`Generate[T]`** - Create new data from templates and examples

### Text Operations
- **`Summarize`** - Intelligent compression to target length
- **`Rewrite`** - Change tone, style, or format
- **`Translate`** - Any language pair
- **`Expand`** - Add detail and context

### Analysis Operations
- **`Classify`** - Categorize into predefined groups
- **`Score`** - Rate based on custom criteria
- **`Compare`** - Semantic comparison of objects
- **`Similar`** - Find similar items by meaning

### Collection Operations
- **`Choose`** - Pick best option with reasoning
- **`Filter`** - Smart filtering with natural language
- **`Sort`** - Multi-criteria intelligent sorting

### Advanced Operations
- **`Validate`** - Check data against business rules
- **`Format`** - Transform to any output format
- **`Merge`** - Intelligently combine data sources
- **`Question`** - Ask questions about your data
- **`Deduplicate`** - Semantic deduplication

### Control Flow
- **`Decide`** - Multi-way decisions with LLM fallback
- **`Guard`** - Conditional execution with suggestions
- **`Workflow`** - Multi-step orchestration
- **`StateMachine`** - Type-safe state management

## Simple Example 🎬

```go
package main

import (
    "fmt"
    "github.com/monstercameron/schemaflow"
)

func main() {
    schemaflow.Init("your-api-key")
    
    // Extract structured data from messy text
    type Product struct {
        Name     string  `json:"name"`
        Price    float64 `json:"price"`
        InStock  bool    `json:"in_stock"`
    }
    
    product, _ := schemaflow.Extract[Product](
        "We have the iPhone 15 Pro available for $999 in stock!",
    )
    
    fmt.Printf("%+v\n", product)
    // Output: {Name:iPhone 15 Pro Price:999 InStock:true}
}
```

## Complex Example 🚀

```go
// Multi-provider pipeline with cost optimization and observability
func ProcessOrders(orders []Order) error {
    // Use different models for different tasks
    extractClient := schemaflow.NewClient(apiKey).
        WithProvider("openai").          // GPT-4 for extraction
        WithDebug(true)
    
    validateClient := schemaflow.NewClient(apiKey).
        WithProvider("anthropic")        // Claude for validation
    
    // Cost-optimized batch processing pipeline
    pipeline := schemaflow.NewPipeline("order-processing").
        Add("extract", func(ctx context.Context, in any) (any, error) {
            return schemaflow.ClientExtract[OrderData](extractClient, in)
        }).
        Add("validate", func(ctx context.Context, in any) (any, error) {
            data := in.(OrderData)
            return schemaflow.ClientValidate(validateClient, data, 
                "must have valid email, items > 0, total > 0")
        }).
        Add("enrich", func(ctx context.Context, in any) (any, error) {
            // Add steering for specific behavior
            return schemaflow.Generate[EnrichedOrder](schemaflow.OpOptions{
                Steering: "Add shipping estimates and tax calculations",
                Intelligence: schemaflow.Quick,  // Use fast model for simple tasks
            })
        })
    
    // Process in batches with 90% cost reduction
    batch := schemaflow.Batch().
        WithMode(schemaflow.MergedMode).
        WithBatchSize(50).
        WithTimeout(5 * time.Minute)
    
    results := schemaflow.ProcessBatch(batch, orders, pipeline)
    
    // Full observability
    log.Printf("Processed: %d orders", results.Metadata.Succeeded)
    log.Printf("Cost: $%.4f", results.Metadata.EstimatedCost)
    log.Printf("Tokens saved: %d", results.Metadata.TokensSaved)
    
    return results.Error
}
```

## Key Features ⚡

### 🎯 **Steering**
Guide LLM behavior without changing code:
```go
result, _ := schemaflow.Extract[Data](input, schemaflow.OpOptions{
    Steering: "Focus on financial data only",
})
```

### 🔄 **Provider Layer**
Mix and match models for optimal performance/cost:
```go
smartClient := client.WithProvider("openai")      // Complex tasks
fastClient := client.WithProvider("anthropic")    // Quick responses  
testClient := client.WithProvider("local")        // Free testing
```

### 📊 **Observability**
Built-in logging, tracing, and cost tracking:
```go
// OpenTelemetry tracing
ctx = otel.Start(ctx, "process-batch")

// Structured logging
logger.Info("Processing", "items", len(items), "cost", batch.EstimatedCost())

// Cost tracking
fmt.Printf("This operation cost: $%.4f\n", result.Metadata.Cost)
```

### 🚀 **Batch Processing**
```go
// Process 1000 items with 10 parallel workers
results := schemaflow.ExtractBatch[Invoice](
    schemaflow.Batch().WithMode(schemaflow.ParallelMode).WithConcurrency(10),
    invoices,
)

// Or merge into 20 API calls instead of 1000 (95% cost reduction)
results := schemaflow.ExtractBatch[Invoice](
    schemaflow.Batch().WithMode(schemaflow.MergedMode).WithBatchSize(50),
    invoices,
)
```

## Get Started 🏃

```bash
go get github.com/monstercameron/schemaflow
```

```go
// Set your API key
export SCHEMAFLOW_API_KEY=your-api-key

// Initialize and go
schemaflow.Init(os.Getenv("SCHEMAFLOW_API_KEY"))
```

## Why SchemaFlow? 🤔

1. **It's Just Go** - No DSL, no magic strings, just strongly-typed Go
2. **Battle-Tested** - Production-ready with retries, timeouts, and error handling
3. **Cost Efficient** - Batch processing saves 90% on API costs
4. **Observable** - Know what's happening with built-in logging and tracing
5. **Testable** - Mock provider means no API calls in tests
6. **Flexible** - Swap providers, add steering, customize everything

## Ready to Make Your LLMs Type-Safe? 🎯

Stop wrestling with JSON strings. Start writing type-safe LLM operations.

[**📚 Full API Documentation →**](API.md)

[**⭐ Star on GitHub**](https://github.com/monstercameron/schemaflow) | [**🐛 Report Issues**](https://github.com/monstercameron/schemaflow/issues) | [**💬 Discussions**](https://github.com/monstercameron/schemaflow/discussions)

---

*Built for Go developers who want LLMs without the chaos* 🚀