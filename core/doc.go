/*
Package schemaflow provides type-safe, production-ready LLM operations for Go applications.

SchemaFlow simplifies working with Large Language Models by providing strongly-typed operations
that handle the complexity of prompt engineering, response parsing, and error handling.

# Features

  - Type-safe operations with Go generics
  - Multiple LLM provider support (OpenAI, Anthropic, local)
  - Batch processing with cost optimization
  - Pipeline and composition support
  - Procedural programming operations
  - Full backward compatibility

# Installation

	go get github.com/monstercameron/schemaflow

# Quick Start

Initialize the library with your API key:

	import "github.com/monstercameron/schemaflow"

	func main() {
	    schemaflow.Init("your-api-key")
	    
	    // Extract structured data
	    type Person struct {
	        Name string `json:"name"`
	        Age  int    `json:"age"`
	    }
	    
	    person, err := schemaflow.Extract[Person]("John Doe, 30 years old")
	}

# Client-Based Usage

For multiple configurations or providers:

	client := schemaflow.NewClient("api-key").
	    WithTimeout(30 * time.Second).
	    WithProvider("anthropic")

	person, err := schemaflow.ClientExtract[Person](client, input)

# Operation Categories

The library is organized into logical operation categories:

Core Operations (data_operations.go):
  - Extract: Extract structured data from text
  - Transform: Transform between data types
  - Generate: Generate new data from templates

Text Operations (text_operations.go):
  - Summarize: Create concise summaries
  - Rewrite: Rewrite with different styles
  - Translate: Translate between languages
  - Expand: Expand text with more detail

Analysis Operations (analysis_operations.go):
  - Classify: Categorize content
  - Score: Score based on criteria
  - Compare: Compare items
  - Similar: Find similar items

Collection Operations (collection_operations.go):
  - Choose: Select best option
  - Filter: Filter by criteria
  - Sort: Sort by criteria

Extended Operations (extended_operations.go):
  - Validate: Validate against rules
  - Format: Format data
  - Merge: Merge multiple sources
  - Question: Answer questions about data
  - Deduplicate: Remove duplicates

Batch Operations (batch_operations.go):
  - ParallelMode: Concurrent processing (5-10x faster)
  - MergedMode: Combined API calls (70-90% cost reduction)
  - SmartBatch: Automatic optimization

Pipeline Operations (pipeline.go):
  - Pipeline: Chain operations
  - Compose: Function composition
  - Map/Reduce: Collection processing

Procedural Operations (procedural_ops.go):
  - Decide: Multi-way decisions
  - StateMachine: State management
  - Workflow: Multi-step workflows
  - Guard: Condition checking

# Batch Processing

Process multiple items efficiently:

	// Parallel processing for speed
	batch := schemaflow.Batch().
	    WithMode(schemaflow.ParallelMode).
	    WithConcurrency(10)
	
	results := schemaflow.ExtractBatch[Person](batch, inputs)

	// Merged processing for cost savings
	batch := schemaflow.Batch().
	    WithMode(schemaflow.MergedMode).
	    WithBatchSize(50)
	
	results := schemaflow.ExtractBatch[Invoice](batch, invoices)

# Pipelines

Chain operations together:

	pipeline := schemaflow.NewPipeline("process").
	    Add("extract", extractOp).
	    Add("validate", validateOp).
	    Add("transform", transformOp)
	
	result := pipeline.Execute(ctx, input)

# Provider Support

Switch between different LLM providers:

	// OpenAI (default)
	client := schemaflow.NewClient(apiKey)

	// Anthropic
	client := schemaflow.NewClient(apiKey).WithProvider("anthropic")

	// Local/Mock for testing
	testClient := schemaflow.NewClient("").WithProvider("local")

# Configuration

Configure via environment variables:

	SCHEMAFLOW_API_KEY=your-api-key
	SCHEMAFLOW_PROVIDER=openai
	SCHEMAFLOW_TIMEOUT=30s
	SCHEMAFLOW_DEBUG=true

Or programmatically:

	client := schemaflow.NewClient(apiKey).
	    WithTimeout(60 * time.Second).
	    WithDebug(true)

# Error Handling

All operations return errors that should be checked:

	result, err := schemaflow.Extract[Data](input)
	if err != nil {
	    // Handle error
	    log.Printf("Extraction failed: %v", err)
	    return err
	}

# Performance

Different modes for different needs:

  - Single operations: Standard performance
  - Parallel batch: 5-10x throughput improvement
  - Merged batch: 70-90% cost reduction
  - Local provider: Instant, free testing

# Testing

Use the local provider for testing:

	testClient := schemaflow.NewClient("").WithProvider("local")
	
	// Configure custom responses
	provider := schemaflow.NewLocalProvider(schemaflow.ProviderConfig{})
	provider.WithHandler(func(ctx context.Context, req schemaflow.CompletionRequest) (string, error) {
	    return "test response", nil
	})
	
	testClient.WithProviderInstance(provider)

# Examples

See the examples/ directory for complete applications:
  - SmartTodo: AI-powered task management
  - Data processing pipelines
  - Document analysis
  - Multi-client configurations

# Best Practices

1. Use typed operations instead of raw LLM calls
2. Batch process when handling multiple items
3. Use appropriate intelligence levels (Quick/Fast/Smart)
4. Test with the local provider
5. Handle all errors appropriately
6. Use pipelines for complex workflows
7. Configure timeouts appropriately

# Thread Safety

All operations are thread-safe. Clients can be shared across goroutines.

# Links

Documentation: https://github.com/monstercameron/schemaflow
Issues: https://github.com/monstercameron/schemaflow/issues
Examples: https://github.com/monstercameron/schemaflow/examples
*/
package schemaflow