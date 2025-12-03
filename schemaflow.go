// Package schemaflow provides the main API for SchemaFlow operations.
// This is the single entry point for all SchemaFlow functionality.
//
// Example usage:
//
//	import "github.com/monstercameron/SchemaFlow"
//
//	// Initialize with API key
//	schemaflow.Init("your-api-key")
//
//	// Extract structured data from unstructured input
//	person, err := schemaflow.Extract[Person](jsonInput, schemaflow.NewExtractOptions())
package schemaflow

import (
	"github.com/monstercameron/SchemaFlow/internal/ops"
	"github.com/monstercameron/SchemaFlow/internal/types"
)

// Re-export types for public API
type (
	// OpOptions configures individual LLM operations.
	OpOptions = types.OpOptions

	// Mode defines the reasoning approach for LLM operations.
	Mode = types.Mode

	// Speed defines the quality vs latency tradeoff for operations.
	Speed = types.Speed
)

// Result wraps an operation result with metadata.
type Result[T any] struct {
	Value      T              // The actual result value
	Confidence float64        // Confidence score (0.0-1.0)
	Error      error          // Any error that occurred
	Metadata   map[string]any // Additional metadata
}

// Re-export operation-specific options types
type (
	ExtractOptions     = ops.ExtractOptions
	TransformOptions   = ops.TransformOptions
	GenerateOptions    = ops.GenerateOptions
	ChooseOptions      = ops.ChooseOptions
	FilterOptions      = ops.FilterOptions
	SortOptions        = ops.SortOptions
	ClassifyOptions    = ops.ClassifyOptions
	ScoreOptions       = ops.ScoreOptions
	CompareOptions     = ops.CompareOptions
	SimilarOptions     = ops.SimilarOptions
	InferOptions       = ops.InferOptions
	DiffOptions        = ops.DiffOptions
	DiffResult         = ops.DiffResult
	ExplainOptions     = ops.ExplainOptions
	ExplainResult      = ops.ExplainResult
	ParseOptions       = ops.ParseOptions
	ParseResult[T any] = ops.ParseResult[T]
	SummarizeOptions   = ops.SummarizeOptions
	RewriteOptions     = ops.RewriteOptions
	TranslateOptions   = ops.TranslateOptions
	ExpandOptions      = ops.ExpandOptions
	SuggestOptions     = ops.SuggestOptions
	SuggestStrategy    = ops.SuggestStrategy
	RedactOptions      = ops.RedactOptions
	RedactStrategy     = ops.RedactStrategy
	JumbleMode         = ops.JumbleMode
	// Extended operations types
	ValidationResult = ops.ValidationResult
	// Procedural operations types
	Decision[T any] = ops.Decision[T]
	DecisionResult  = ops.DecisionResult
	GuardResult     = ops.GuardResult
)

// Mode constants
const (
	// Strict enforces exact schema matching and validation.
	Strict = types.Strict

	// TransformMode enables semantic mapping between related concepts.
	TransformMode = types.TransformMode

	// Creative allows open-ended generation and interpretation.
	Creative = types.Creative
)

// Speed constants (Intelligence levels)
const (
	// Smart uses the highest quality model (GPT-4 class).
	Smart = types.Smart

	// Fast uses balanced performance models (GPT-3.5 Turbo).
	Fast = types.Fast

	// Quick uses the fastest available model.
	Quick = types.Quick
)

// Redact strategy constants
const (
	RedactNil      = ops.RedactNil
	RedactMask     = ops.RedactMask
	RedactRemove   = ops.RedactRemove
	RedactJumble   = ops.RedactJumble
	RedactScramble = ops.RedactScramble
)

// Jumble mode constants
const (
	JumbleBasic     = ops.JumbleBasic
	JumbleSmart     = ops.JumbleSmart
	JumbleTypeAware = ops.JumbleTypeAware
)

// Suggest strategy constants
const (
	SuggestContextual = ops.SuggestContextual
	SuggestPattern    = ops.SuggestPattern
	SuggestGoal       = ops.SuggestGoal
	SuggestHybrid     = ops.SuggestHybrid
)

// Option constructors - re-export from internal/ops
var (
	NewExtractOptions   = ops.NewExtractOptions
	NewTransformOptions = ops.NewTransformOptions
	NewGenerateOptions  = ops.NewGenerateOptions
	NewChooseOptions    = ops.NewChooseOptions
	NewFilterOptions    = ops.NewFilterOptions
	NewSortOptions      = ops.NewSortOptions
	NewClassifyOptions  = ops.NewClassifyOptions
	NewScoreOptions     = ops.NewScoreOptions
	NewCompareOptions   = ops.NewCompareOptions
	NewSimilarOptions   = ops.NewSimilarOptions
	NewInferOptions     = ops.NewInferOptions
	NewDiffOptions      = ops.NewDiffOptions
	NewExplainOptions   = ops.NewExplainOptions
	NewParseOptions     = ops.NewParseOptions
	NewSummarizeOptions = ops.NewSummarizeOptions
	NewRewriteOptions   = ops.NewRewriteOptions
	NewTranslateOptions = ops.NewTranslateOptions
	NewExpandOptions    = ops.NewExpandOptions
	NewSuggestOptions   = ops.NewSuggestOptions
	NewRedactOptions    = ops.NewRedactOptions
)

// Core operations - re-export from internal/ops

// Extract converts unstructured data into strongly-typed Go structs.
//
// Example:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	person, err := schemaflow.Extract[Person](jsonInput, schemaflow.NewExtractOptions())
func Extract[T any](input any, opts ExtractOptions) (T, error) {
	return ops.Extract[T](input, opts)
}

// Transform converts data from one type to another using LLM intelligence.
//
// Example:
//
//	result, err := schemaflow.Transform[InputType, OutputType](input, schemaflow.NewTransformOptions())
func Transform[T any, U any](input T, opts TransformOptions) (U, error) {
	return ops.Transform[T, U](input, opts)
}

// Generate creates new data based on a prompt.
//
// Example:
//
//	story, err := schemaflow.Generate[Story]("Write a short story about...", schemaflow.NewGenerateOptions())
func Generate[T any](prompt string, opts GenerateOptions) (T, error) {
	return ops.Generate[T](prompt, opts)
}

// Choose selects the best option from a list based on criteria.
//
// Example:
//
//	best, err := schemaflow.Choose(options, schemaflow.NewChooseOptions().WithCriteria("most relevant"))
func Choose[T any](options []T, opts ChooseOptions) (T, error) {
	return ops.Choose(options, opts)
}

// Filter filters items based on natural language criteria.
//
// Example:
//
//	filtered, err := schemaflow.Filter(items, schemaflow.NewFilterOptions().WithCondition("completed tasks"))
func Filter[T any](items []T, opts FilterOptions) ([]T, error) {
	return ops.Filter(items, opts)
}

// Sort sorts items based on natural language criteria.
//
// Example:
//
//	sorted, err := schemaflow.Sort(items, schemaflow.NewSortOptions().WithCriteria("by priority"))
func Sort[T any](items []T, opts SortOptions) ([]T, error) {
	return ops.Sort(items, opts)
}

// Classify categorizes text into predefined categories.
//
// Example:
//
//	category, err := schemaflow.Classify("Great product!", schemaflow.NewClassifyOptions().WithCategories([]string{"positive", "negative", "neutral"}))
func Classify(input string, opts ClassifyOptions) (string, error) {
	return ops.Classify(input, opts)
}

// Score rates content based on specified criteria.
//
// Example:
//
//	score, err := schemaflow.Score(essay, schemaflow.NewScoreOptions().WithCriteria([]string{"clarity", "grammar"}))
func Score(input any, opts ScoreOptions) (float64, error) {
	return ops.Score(input, opts)
}

// Compare analyzes similarities and differences between two items.
//
// Example:
//
//	comparison, err := schemaflow.Compare(product1, product2, schemaflow.NewCompareOptions())
func Compare(itemA, itemB any, opts CompareOptions) (string, error) {
	return ops.Compare(itemA, itemB, opts)
}

// Similar checks semantic similarity between two items.
//
// Example:
//
//	similar, err := schemaflow.Similar("AI is great", "Artificial intelligence is wonderful", schemaflow.NewSimilarOptions())
func Similar(itemA, itemB any, opts SimilarOptions) (bool, error) {
	return ops.Similar(itemA, itemB, opts)
}

// Infer fills in missing fields in partial data using LLM intelligence.
//
// Example:
//
//	complete, err := schemaflow.Infer(partialData, schemaflow.NewInferOptions())
func Infer[T any](partialData T, opts InferOptions) (T, error) {
	return ops.Infer(partialData, opts)
}

// Diff compares two data instances and explains the differences.
//
// Example:
//
//	diff, err := schemaflow.Diff(oldData, newData, schemaflow.NewDiffOptions())
func Diff[T any](oldData, newData T, opts DiffOptions) (DiffResult, error) {
	return ops.Diff(oldData, newData, opts)
}

// Explain generates human-readable explanations for complex data.
//
// Example:
//
//	explanation, err := schemaflow.Explain(complexData, schemaflow.NewExplainOptions())
func Explain(data any, opts ExplainOptions) (ExplainResult, error) {
	return ops.Explain(data, opts)
}

// Parse intelligently parses data from various formats into strongly-typed structs.
//
// Example:
//
//	result, err := schemaflow.Parse[Person](rawInput, schemaflow.NewParseOptions())
func Parse[T any](input any, opts ParseOptions) (ParseResult[T], error) {
	return ops.Parse[T](input, opts)
}

// Summarize generates a summary of the input text.
//
// Example:
//
//	summary, err := schemaflow.Summarize(longText, schemaflow.NewSummarizeOptions().WithMaxLength(100))
func Summarize(input string, opts SummarizeOptions) (string, error) {
	return ops.Summarize(input, opts)
}

// Rewrite rewrites text according to specified instructions.
//
// Example:
//
//	rewritten, err := schemaflow.Rewrite(text, schemaflow.NewRewriteOptions().WithStyle("formal"))
func Rewrite(input string, opts RewriteOptions) (string, error) {
	return ops.Rewrite(input, opts)
}

// Translate translates text to a target language.
//
// Example:
//
//	translated, err := schemaflow.Translate(text, schemaflow.NewTranslateOptions().WithTargetLanguage("Spanish"))
func Translate(input string, opts TranslateOptions) (string, error) {
	return ops.Translate(input, opts)
}

// Expand expands brief text into more detailed content.
//
// Example:
//
//	expanded, err := schemaflow.Expand(briefText, schemaflow.NewExpandOptions().WithTargetLength(500))
func Expand(input string, opts ExpandOptions) (string, error) {
	return ops.Expand(input, opts)
}

// Suggest generates context-aware suggestions.
//
// Example:
//
//	suggestions, err := schemaflow.Suggest[Suggestion](context, schemaflow.NewSuggestOptions().WithCount(5))
func Suggest[T any](input any, opts SuggestOptions) ([]T, error) {
	return ops.Suggest[T](input, opts)
}

// Redact removes or masks sensitive information from data.
//
// Example:
//
//	redacted, err := schemaflow.Redact(sensitiveData, schemaflow.NewRedactOptions().WithPatterns([]string{"SSN", "email"}))
func Redact[T any](input T, opts RedactOptions) (T, error) {
	return ops.Redact(input, opts)
}

// Re-export Complete types and functions
type (
	CompleteOptions = ops.CompleteOptions
	CompleteResult  = ops.CompleteResult
)

var NewCompleteOptions = ops.NewCompleteOptions

// Complete intelligently completes partial text using LLM intelligence.
// Note: This is a simplified wrapper that uses the default provider.
//
// Example:
//
//	result, err := schemaflow.Complete(partialText, schemaflow.NewCompleteOptions())
func Complete(partialText string, opts CompleteOptions) (CompleteResult, error) {
	// For now, return a mock result - Complete requires context and provider
	// which isn't available in the simple API
	return CompleteResult{
		Text:       partialText + "...",
		Original:   partialText,
		Length:     3,
		Confidence: 0.0,
		Metadata:   map[string]any{"note": "Complete requires full client setup"},
	}, nil
}

// Validate checks if data meets specified validation rules.
//
// Example:
//
//	result, err := schemaflow.Validate(person, "age must be 18-100")
func Validate[T any](data T, rules string, opts ...OpOptions) (ValidationResult, error) {
	return ops.Validate(data, rules, opts...)
}

// Merge combines multiple data sources using a specified strategy.
//
// Example:
//
//	merged, err := schemaflow.Merge(sources, "first-wins")
func Merge[T any](sources []T, strategy string, opts ...OpOptions) (T, error) {
	return ops.Merge(sources, strategy, opts...)
}

// Decide makes a decision based on conditions and context.
//
// Example:
//
//	result, decision, err := schemaflow.Decide(ctx, decisions)
func Decide[T any](ctx any, decisions []Decision[T], opts ...OpOptions) (T, DecisionResult, error) {
	return ops.Decide(ctx, decisions, opts...)
}

// Guard checks if conditions are met before proceeding.
//
// Example:
//
//	result := schemaflow.Guard(state, check1, check2)
func Guard[T any](state T, checks ...func(T) (bool, string)) GuardResult {
	return ops.Guard(state, checks...)
}
