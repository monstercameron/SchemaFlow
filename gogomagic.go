// Package schemaflow provides LLM-powered fuzzy programming operations for Go.
// It enables semantic reasoning and natural language processing while maintaining
// Go's type safety and idiomatic patterns.
//
// Core Features:
//   - Type-safe operations with Go generics
//   - Natural language steering for fine-tuned control
//   - Multiple intelligence levels (Smart/Fast/Quick)
//   - Three reasoning modes (Strict/Transform/Creative)
//   - Comprehensive error handling with context
//   - Request tracing and metrics collection
//   - Retry logic with exponential backoff
//
// Basic Usage:
//
//	import "github.com/monstercameron/schemaflow"
//
//	func main() {
//	    // Initialize with API key
//	    schemaflow.Init("your-api-key")
//
//	    // Extract structured data from text
//	    person, err := schemaflow.Extract[Person]("John Doe, 30 years old")
//
//	    // Generate data from prompts
//	    product, err := schemaflow.Generate[Product]("Create a laptop listing")
//
//	    // Transform between types
//	    employee, err := schemaflow.Transform[Person, Employee](person)
//	}
//
// Configuration:
//
// The library can be configured through environment variables:
//   - SCHEMAFLOW_API_KEY: OpenAI API key
//   - SCHEMAFLOW_PROVIDER: LLM provider (default: "openai")
//   - SCHEMAFLOW_TIMEOUT: Operation timeout (default: "30s")
//   - SCHEMAFLOW_MAX_RETRIES: Number of retries (default: 3)
//   - SCHEMAFLOW_DEBUG: Enable debug logging (default: false)
//   - SCHEMAFLOW_TRACE: Enable request tracing (default: false)
//   - SCHEMAFLOW_METRICS: Enable metrics collection (default: false)
//
// Operations:
//
// Data Operations:
//   - Extract[T]: Convert unstructured data to typed structs
//   - Transform[T,U]: Convert between different types
//   - Generate[T]: Create data from natural language
//
// Text Operations:
//   - Summarize: Create concise summaries
//   - Rewrite: Transform text style and tone
//   - Translate: Natural language translation
//   - Expand: Elaborate on concepts
//
// Analysis Operations:
//   - Classify: Categorize text into groups
//   - Score: Rate quality or relevance
//   - Compare: Semantic comparison
//   - Similar: Check semantic similarity
//
// Collection Operations:
//   - Choose: Select best option from list
//   - Filter: Semantic filtering of items
//   - Sort: Order by semantic criteria
//
// Control Flow:
//   - Match: Pattern matching with fuzzy logic
//
// Error Handling:
//
// All operations return typed errors with context:
//
//	person, err := schemaflow.Extract[Person](input)
//	if err != nil {
//	    switch e := err.(type) {
//	    case schemaflow.ExtractError:
//	        log.Printf("Extraction failed: %s (confidence: %.2f)",
//	            e.Reason, e.Confidence)
//	    default:
//	        log.Printf("Unexpected error: %v", err)
//	    }
//	}
//
// Testing:
//
// The library supports mocking for testing without API calls:
//
//	// In your test file
//	oldCallLLM := schemaflow.callLLM
//	defer func() { schemaflow.callLLM = oldCallLLM }()
//
//	schemaflow.callLLM = func(ctx context.Context, system, user string, opts OpOptions) (string, error) {
//	    return `{"name": "Test User", "age": 30}`, nil
//	}
//
// Thread Safety:
//
// All operations are thread-safe and can be called concurrently.
// The Init function should be called once at program startup.
//
// For more examples and detailed documentation, see:
// https://github.com/monstercameron/schemaflow
package schemaflow

// This file serves as the main entry point and documentation for the schemaflow package.
// All implementation has been organized into separate files:
//
// - types.go: Core type definitions (Mode, Speed, OpOptions, Result, Case)
// - errors.go: Custom error types for all operations
// - config.go: Initialization and configuration management
// - llm.go: LLM interaction and communication layer
// - utils.go: Helper functions and utilities
// - data_operations.go: Extract, Transform, Generate operations
// - text_operations.go: Text manipulation operations
// - analysis_operations.go: Analysis and comparison operations
// - collection_operations.go: Collection manipulation operations
// - control_operations.go: Control flow operations
// - steering.go: Steering presets and shortcuts
// - logger.go: Structured logging implementation
// - debug.go: Debugging utilities and tracing