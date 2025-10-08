# SchemaFlow API Reference

## Getting Started

### Installation

```bash
go get github.com/monstercameron/schemaflow
```

### Basic Setup

```go
package main

import "github.com/monstercameron/schemaflow"

func main() {
    // Option 1: Initialize globally
    schemaflow.Init("your-api-key")
    
    // Option 2: Create a client
    client := schemaflow.NewClient("your-api-key")
}
```

### Environment Variables

```bash
SCHEMAFLOW_API_KEY=your-api-key
SCHEMAFLOW_PROVIDER=openai       # or anthropic, local
SCHEMAFLOW_TIMEOUT=30s
SCHEMAFLOW_DEBUG=true
```

## Core Operations

### Extract
Extract structured data from unstructured text.

```go
type Person struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}

// Extract from text
person, err := schemaflow.Extract[Person]("John Doe, 30, john@example.com")

// With options
person, err := schemaflow.Extract[Person](text, schemaflow.OpOptions{
    Intelligence: schemaflow.Smart,  // or Fast, Quick
    Mode: schemaflow.Strict,         // or Creative
})

// With client
person, err := schemaflow.ClientExtract[Person](client, text)
```

### Transform
Transform data from one type to another.

```go
type Report struct {
    Title   string
    Content string
    Data    []float64
}

type Summary struct {
    Key     string
    Points  []string
    Average float64
}

summary, err := schemaflow.Transform[Report, Summary](report)
```

### Generate
Generate new data based on templates or examples.

```go
type Email struct {
    Subject string
    Body    string
    Tone    string
}

email, err := schemaflow.Generate[Email](schemaflow.OpOptions{
    Steering: "Professional thank you email for interview",
})
```

## Text Operations

### Summarize
Create concise summaries.

```go
summary, err := schemaflow.Summarize(longText, 100) // 100 words max

// With custom options
summary, err := schemaflow.Summarize(text, 200, schemaflow.OpOptions{
    Steering: "Focus on technical details",
})
```

### Rewrite
Rewrite text with different style or tone.

```go
formal, err := schemaflow.Rewrite(casualText, "formal business tone")
simple, err := schemaflow.Rewrite(technical, "explain like I'm five")
```

### Translate
Translate between languages.

```go
spanish, err := schemaflow.Translate(englishText, "Spanish")
japanese, err := schemaflow.Translate(text, "Japanese")
```

### Expand
Expand text with more detail.

```go
detailed, err := schemaflow.Expand(briefNotes, 500) // Expand to ~500 words
```

## Analysis Operations

### Classify
Categorize content into predefined groups.

```go
categories := []string{"spam", "important", "newsletter", "personal"}
category, confidence, err := schemaflow.Classify(emailContent, categories)

if confidence > 0.8 {
    fmt.Printf("Classified as: %s\n", category)
}
```

### Score
Score content based on criteria.

```go
score, err := schemaflow.Score(essay, "writing quality", 100) // Score out of 100
rating, err := schemaflow.Score(review, "sentiment", 5)       // 1-5 stars
```

### Compare
Compare two items and find differences.

```go
type Document struct {
    Title   string
    Content string
    Author  string
}

comparison, err := schemaflow.Compare(doc1, doc2)
// Returns differences, similarities, and analysis
```

### Similar
Find semantically similar items.

```go
type Product struct {
    Name        string
    Description string
}

// Find similar products
similar, scores, err := schemaflow.Similar(targetProduct, products, 0.8)
// Returns products with similarity >= 0.8
```

## Collection Operations

### Choose
Select the best option from choices.

```go
type Option struct {
    Name  string
    Pros  []string
    Cons  []string
    Cost  float64
}

best, reason, err := schemaflow.Choose(options, "best value for money")
```

### Filter
Filter items based on criteria.

```go
type Task struct {
    Title    string
    Priority int
    Due      time.Time
}

urgent, err := schemaflow.Filter(tasks, "priority > 8 and due today")
```

### Sort
Sort items by criteria.

```go
sorted, err := schemaflow.Sort(tasks, "priority desc, due asc")
```

## Extended Operations

### Validate
Validate data against rules.

```go
result, err := schemaflow.Validate(userData, `
    - Email must be valid
    - Age must be 18-100
    - Phone number must be US format
`)

if !result.Valid {
    fmt.Printf("Issues: %v\n", result.Issues)
    fmt.Printf("Suggestions: %v\n", result.Suggestions)
}
```

### Format
Format data into specific output formats.

```go
// Format as markdown table
table, err := schemaflow.Format(data, "markdown table with headers")

// Format as HTML
html, err := schemaflow.Format(data, "HTML with Bootstrap classes")

// Format as report
report, err := schemaflow.Format(data, "executive summary report")
```

### Merge
Intelligently merge multiple data sources.

```go
type Customer struct {
    ID      string
    Name    string
    Email   string
    Updated time.Time
}

// Merge with strategy
merged, err := schemaflow.Merge([]Customer{dbRecord, apiData, csvRow}, 
    "prefer newest, combine addresses")
```

### Question
Answer questions about data.

```go
answer, err := schemaflow.Question(financialData, "What was the Q3 revenue?")
insights, err := schemaflow.Question(metrics, "What are the top 3 trends?")
```

### Deduplicate
Remove duplicates using semantic similarity.

```go
result, err := schemaflow.Deduplicate(customers, 0.85) // 85% similarity threshold

fmt.Printf("Removed %d duplicates\n", result.TotalRemoved)
fmt.Printf("Unique items: %d\n", len(result.Unique))
```

## Batch Operations

### Parallel Mode
Process items concurrently for speed.

```go
batch := schemaflow.Batch().
    WithMode(schemaflow.ParallelMode).
    WithConcurrency(10).
    WithTimeout(5 * time.Minute)

results := schemaflow.ExtractBatch[Person](batch, inputs)

fmt.Printf("Processed: %d/%d\n", results.Metadata.Succeeded, results.Metadata.TotalItems)
fmt.Printf("Duration: %v\n", results.Metadata.Duration)
```

### Merged Mode
Combine items into single API calls for cost savings.

```go
batch := schemaflow.Batch().
    WithMode(schemaflow.MergedMode).
    WithBatchSize(50) // Items per API call

results := schemaflow.ExtractBatch[Invoice](batch, invoices)

fmt.Printf("API calls: %d (saved %d%%)\n", 
    results.Metadata.APICallsMade,
    (1000 - results.Metadata.APICallsMade) * 100 / 1000)
```

### Smart Batch
Automatic mode selection based on data.

```go
smartBatch := client.SmartBatch()
results := schemaflow.ExtractSmart[Product](smartBatch, products)
// Automatically chooses best mode
```

## Pipeline Operations

### Creating Pipelines

```go
pipeline := schemaflow.NewPipeline("data-processing").
    Add("extract", func(ctx context.Context, input any) (any, error) {
        return schemaflow.Extract[Data](input.(string))
    }).
    Add("validate", func(ctx context.Context, input any) (any, error) {
        return schemaflow.Validate(input, "rules")
    }).
    Add("transform", func(ctx context.Context, input any) (any, error) {
        return schemaflow.Transform[Data, Output](input.(Data))
    })

result := pipeline.Execute(context.Background(), rawInput)
```

### Composition Functions

```go
// Chain operations
process := schemaflow.Then(
    func(s string) (Data, error) { return schemaflow.Extract[Data](s) },
    func(d Data) (Output, error) { return schemaflow.Transform[Data, Output](d) },
)

output, err := process(input)

// Map over collections
results, err := schemaflow.Map(items, func(item Item) (Result, error) {
    return processItem(item)
})

// Parallel map
results, err := schemaflow.MapConcurrent(items, processItem, 10)

// Reduce
total, err := schemaflow.Reduce(numbers, func(a, b int) int {
    return a + b
})
```

## Procedural Operations

### Decide
Make decisions with LLM fallback.

```go
decisions := []schemaflow.Decision[string]{
    {
        Value: "urgent",
        Condition: func(ctx any) bool {
            return ctx.(Task).Priority > 8
        },
        Description: "High priority tasks",
    },
    {
        Value: "normal",
        Description: "Regular tasks",
    },
}

action, result, err := schemaflow.Decide(task, decisions)
```

### State Machine

```go
type State string
const (
    Pending State = "pending"
    Active  State = "active"
    Done    State = "done"
)

sm := schemaflow.NewStateMachine[State, string](Pending)

sm.AddTransition(Pending, "start", Active)
sm.AddTransition(Active, "complete", Done)

newState, err := sm.Transition("start")
```

### Workflow

```go
workflow := schemaflow.NewWorkflow("order-processing")

workflow.AddStep(schemaflow.WorkflowStep{
    Name: "validate",
    Execute: func(ctx context.Context, state map[string]any) error {
        // Validation logic
        return nil
    },
})

workflow.AddStep(schemaflow.WorkflowStep{
    Name: "process",
    Execute: func(ctx context.Context, state map[string]any) error {
        // Processing logic
        return nil
    },
    Dependencies: []string{"validate"},
})

err := workflow.Execute(context.Background())
```

### Control Flow

```go
// Loop with condition
final, err := schemaflow.LoopWhile(
    initialState,
    func(s State) bool { return !s.Complete },
    func(s State) (State, error) { return processState(s) },
    100, // max iterations
)

// Switch with typed returns
result := schemaflow.Switch(value, map[string]func() Result{
    "a": func() Result { return ResultA() },
    "b": func() Result { return ResultB() },
}, func() Result { return DefaultResult() })

// If-else with typed returns
result := schemaflow.IfElse(condition,
    func() string { return "true case" },
    func() string { return "false case" },
)

// Try with panic recovery
result, err := schemaflow.Try(func() (Result, error) {
    return riskyOperation()
})
```

## Provider Configuration

### Switching Providers

```go
// OpenAI (default)
client := schemaflow.NewClient(apiKey)

// Anthropic
client := schemaflow.NewClient(apiKey).WithProvider("anthropic")

// Local/Mock
testClient := schemaflow.NewClient("").WithProvider("local")

// Custom provider
customProvider := &MyCustomProvider{}
client.WithProviderInstance(customProvider)
```

### Provider Interface

```go
type Provider interface {
    Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
    Name() string
    EstimateCost(req CompletionRequest) float64
}
```

### Local Provider for Testing

```go
provider := schemaflow.NewLocalProvider(schemaflow.ProviderConfig{})

provider.WithHandler(func(ctx context.Context, req schemaflow.CompletionRequest) (string, error) {
    // Custom mock response logic
    if strings.Contains(req.UserPrompt, "extract") {
        return `{"name": "Test", "value": 123}`, nil
    }
    return "mock response", nil
})

testClient := schemaflow.NewClient("").WithProviderInstance(provider)
```

## Options and Configuration

### Operation Options

```go
type OpOptions struct {
    Intelligence Speed      // Quick, Fast, Smart
    Mode         Mode       // Strict, TransformMode, Creative
    Threshold    float64    // Confidence threshold (0-1)
    Steering     string     // Additional instructions
}

// Example usage
result, err := schemaflow.Extract[Data](input, schemaflow.OpOptions{
    Intelligence: schemaflow.Smart,
    Mode: schemaflow.Strict,
    Threshold: 0.9,
    Steering: "Focus on numerical data",
})
```

### Client Configuration

```go
client := schemaflow.NewClient(apiKey).
    WithTimeout(60 * time.Second).
    WithDebug(true).
    WithProvider("anthropic")
```

## Error Handling

All operations return errors that should be checked:

```go
result, err := schemaflow.Extract[Data](input)
if err != nil {
    switch {
    case errors.Is(err, schemaflow.ErrInvalidInput):
        // Handle invalid input
    case errors.Is(err, schemaflow.ErrTimeout):
        // Handle timeout
    case errors.Is(err, schemaflow.ErrRateLimit):
        // Handle rate limiting
    default:
        // Handle other errors
    }
}
```

## Testing

```go
func TestMyFunction(t *testing.T) {
    // Create test client with local provider
    testClient := schemaflow.NewClient("").WithProvider("local")
    
    // Use in tests
    result, err := schemaflow.ClientExtract[MyType](testClient, "test input")
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## Performance Tips

1. **Use batch operations** for multiple items
2. **Choose appropriate intelligence level** - Quick for simple, Smart for complex
3. **Use Merged mode** for similar items to save costs
4. **Cache results** when possible using `CachedOperation`
5. **Set appropriate timeouts** based on operation complexity
6. **Use local provider** for development and testing

## Common Patterns

### ETL Pipeline
```go
pipeline := schemaflow.NewPipeline("etl").
    Add("extract", extractOp).
    Add("transform", transformOp).
    Add("load", loadOp)
```

### Validation Flow
```go
valid, err := schemaflow.Validate(data, rules)
if !valid.Valid {
    formatted, _ := schemaflow.Format(valid.Issues, "user-friendly errors")
    return errors.New(formatted)
}
```

### Batch Processing with Error Handling
```go
batch := schemaflow.Batch().WithMode(schemaflow.ParallelMode)
results := schemaflow.ExtractBatch[Type](batch, items)

for i, result := range results.Results {
    if result.Error != nil {
        log.Printf("Item %d failed: %v", i, result.Error)
        continue
    }
    // Process successful result
}
```