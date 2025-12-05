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
	"context"

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
	ValidationResult      = ops.ValidationResult
	ValidateOptions       = ops.ValidateOptions
	ValidateResult[T any] = ops.ValidateResult[T]
	ValidationIssue       = ops.ValidationIssue
	QuestionOptions       = ops.QuestionOptions
	QuestionResult[A any] = ops.QuestionResult[A]
	// Procedural operations types
	Decision[T any] = ops.Decision[T]
	DecisionResult  = ops.DecisionResult
	GuardResult     = ops.GuardResult

	// New LLM operation types (v2)
	AnnotateOptions           = ops.AnnotateOptions
	Annotation                = ops.Annotation
	AnnotateResult            = ops.AnnotateResult
	ClusterOptions            = ops.ClusterOptions
	ClusterInfo[T any]        = ops.ClusterInfo[T]
	ClusterResult[T any]      = ops.ClusterResult[T]
	RankOptions               = ops.RankOptions
	RankedItem[T any]         = ops.RankedItem[T]
	RankResult[T any]         = ops.RankResult[T]
	CompressOptions           = ops.CompressOptions
	CompressResult[T any]     = ops.CompressResult[T]
	DecomposeOptions          = ops.DecomposeOptions
	DecomposedPart[T any]     = ops.DecomposedPart[T]
	DecomposeResult[T any]    = ops.DecomposeResult[T]
	EnrichOptions             = ops.EnrichOptions
	EnrichResult[T any]       = ops.EnrichResult[T]
	NormalizeOptions          = ops.NormalizeOptions
	NormalizeChange           = ops.NormalizeChange
	NormalizeResult[T any]    = ops.NormalizeResult[T]
	MatchOptions              = ops.MatchOptions
	MatchPair[S any, T any]   = ops.MatchPair[S, T]
	MatchResult[S any, T any] = ops.MatchResult[S, T]
	CritiqueOptions           = ops.CritiqueOptions
	CritiqueIssue             = ops.CritiqueIssue
	CritiquePositive          = ops.CritiquePositive
	CritiqueResult            = ops.CritiqueResult
	SynthesizeOptions         = ops.SynthesizeOptions
	SynthesisFact             = ops.SynthesisFact
	SynthesisInsight          = ops.SynthesisInsight
	SynthesisConflict         = ops.SynthesisConflict
	SynthesizeResult[T any]   = ops.SynthesizeResult[T]
	PredictOptions            = ops.PredictOptions
	PredictionInterval        = ops.PredictionInterval
	PredictionScenario        = ops.PredictionScenario
	PredictionFactor          = ops.PredictionFactor
	PredictResult[T any]      = ops.PredictResult[T]
	VerifyOptions             = ops.VerifyOptions
	ClaimVerification         = ops.ClaimVerification
	LogicIssue                = ops.LogicIssue
	ConsistencyIssue          = ops.ConsistencyIssue
	VerifyResult              = ops.VerifyResult

	// Analysis operation result types (refactored for Go-native generics)
	ClassifyResult[C any]      = ops.ClassifyResult[C]
	ClassifyAlternative[C any] = ops.ClassifyAlternative[C]
	ScoreResult                = ops.ScoreResult
	CompareResult[T any]       = ops.CompareResult[T]
	ComparisonPoint            = ops.ComparisonPoint
	SimilarResult              = ops.SimilarResult
	AspectMatch                = ops.AspectMatch

	// Data-centric LLM operations (v3)
	NegotiateOptions       = ops.NegotiateOptions
	Tradeoff               = ops.Tradeoff
	NegotiateResult[T any] = ops.NegotiateResult[T]

	ResolveOptions       = ops.ResolveOptions
	Conflict             = ops.Conflict
	ResolveResult[T any] = ops.ResolveResult[T]

	DeriveOptions       = ops.DeriveOptions
	Derivation          = ops.Derivation
	DeriveResult[U any] = ops.DeriveResult[U]

	ConformOptions       = ops.ConformOptions
	Adjustment           = ops.Adjustment
	ConformResult[T any] = ops.ConformResult[T]

	InterpolateOptions       = ops.InterpolateOptions
	FilledItem               = ops.FilledItem
	InterpolateResult[T any] = ops.InterpolateResult[T]

	ArbitrateOptions       = ops.ArbitrateOptions
	RuleEvaluation         = ops.RuleEvaluation
	OptionEvaluation       = ops.OptionEvaluation
	ArbitrateResult[T any] = ops.ArbitrateResult[T]

	ProjectOptions       = ops.ProjectOptions
	FieldMapping         = ops.FieldMapping
	ProjectResult[U any] = ops.ProjectResult[U]

	AuditOptions       = ops.AuditOptions
	AuditFinding       = ops.AuditFinding
	AuditSummary       = ops.AuditSummary
	AuditResult[T any] = ops.AuditResult[T]

	ComposeOptions       = ops.ComposeOptions
	ComposedField        = ops.ComposedField
	ComposeResult[T any] = ops.ComposeResult[T]

	PivotOptions       = ops.PivotOptions
	PivotMapping       = ops.PivotMapping
	PivotStats         = ops.PivotStats
	PivotResult[U any] = ops.PivotResult[U]

	// Text operation result types with metadata
	SummarizeResult        = ops.SummarizeResult
	RewriteResult          = ops.RewriteResult
	TranslateResult        = ops.TranslateResult
	TranslationAlternative = ops.TranslationAlternative
	ExpandResult           = ops.ExpandResult

	// Extended operation result types with metadata
	FormatResult       = ops.FormatResult
	MergeResult[T any] = ops.MergeResult[T]
	MergeConflict      = ops.MergeConflict
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

	// New LLM operation option constructors (v2)
	NewAnnotateOptions   = ops.NewAnnotateOptions
	NewClusterOptions    = ops.NewClusterOptions
	NewRankOptions       = ops.NewRankOptions
	NewCompressOptions   = ops.NewCompressOptions
	NewDecomposeOptions  = ops.NewDecomposeOptions
	NewEnrichOptions     = ops.NewEnrichOptions
	NewNormalizeOptions  = ops.NewNormalizeOptions
	NewMatchOptions      = ops.NewMatchOptions
	NewCritiqueOptions   = ops.NewCritiqueOptions
	NewSynthesizeOptions = ops.NewSynthesizeOptions
	NewPredictOptions    = ops.NewPredictOptions
	NewVerifyOptions     = ops.NewVerifyOptions
	NewValidateOptions   = ops.NewValidateOptions
	NewQuestionOptions   = ops.NewQuestionOptions
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

// Classify categorizes any Go type into typed categories.
//
// Type parameter T specifies the input type.
// Type parameter C specifies the category type (typically string or custom enum).
//
// Example:
//
//	result, err := schemaflow.Classify[string, string]("Great product!",
//	    schemaflow.NewClassifyOptions().WithCategories([]string{"positive", "negative", "neutral"}))
//	fmt.Printf("Category: %s (%.0f%% confidence)\n", result.Category, result.Confidence*100)
func Classify[T any, C any](input T, opts ClassifyOptions) (ClassifyResult[C], error) {
	return ops.Classify[T, C](input, opts)
}

// Score rates any Go type based on specified criteria.
//
// Type parameter T specifies the input type.
//
// Example:
//
//	result, err := schemaflow.Score[Essay](essay,
//	    schemaflow.NewScoreOptions().WithCriteria([]string{"clarity", "grammar"}))
//	fmt.Printf("Score: %.1f/10\n", result.Value)
func Score[T any](input T, opts ScoreOptions) (ScoreResult, error) {
	return ops.Score(input, opts)
}

// Compare analyzes similarities and differences between two items of the same type.
//
// Type parameter T specifies the type of items being compared.
//
// Example:
//
//	result, err := schemaflow.Compare[Product](product1, product2, schemaflow.NewCompareOptions())
//	fmt.Printf("Similarity: %.0f%%\n", result.SimilarityScore*100)
func Compare[T any](itemA, itemB T, opts CompareOptions) (CompareResult[T], error) {
	return ops.Compare(itemA, itemB, opts)
}

// Similar checks semantic similarity between two items of the same type.
//
// Type parameter T specifies the type of items being compared.
//
// Example:
//
//	result, err := schemaflow.Similar[string]("AI is great", "Artificial intelligence is wonderful",
//	    schemaflow.NewSimilarOptions())
//	fmt.Printf("Similar: %v (score: %.0f%%)\n", result.IsSimilar, result.Score*100)
func Similar[T any](itemA, itemB T, opts SimilarOptions) (SimilarResult, error) {
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

// SummarizeWithMetadata summarizes text and returns rich metadata including
// compression ratio, key points, and confidence score.
//
// Example:
//
//	result, err := schemaflow.SummarizeWithMetadata(longText, schemaflow.NewSummarizeOptions())
//	fmt.Printf("Summary: %s\nKey points: %v\nCompression: %.0f%%\n",
//	    result.Text, result.KeyPoints, result.CompressionRatio*100)
func SummarizeWithMetadata(input string, opts SummarizeOptions) (SummarizeResult, error) {
	return ops.SummarizeWithMetadata(input, opts)
}

// Rewrite rewrites text according to specified instructions.
//
// Example:
//
//	rewritten, err := schemaflow.Rewrite(text, schemaflow.NewRewriteOptions().WithStyle("formal"))
func Rewrite(input string, opts RewriteOptions) (string, error) {
	return ops.Rewrite(input, opts)
}

// RewriteWithMetadata rewrites text and returns metadata including
// what changes were made, tone achieved, and confidence score.
//
// Example:
//
//	result, err := schemaflow.RewriteWithMetadata(text, schemaflow.NewRewriteOptions().WithTargetTone("professional"))
//	fmt.Printf("Rewritten: %s\nChanges: %v\nTone: %s\n",
//	    result.Text, result.ChangesMade, result.ToneAchieved)
func RewriteWithMetadata(input string, opts RewriteOptions) (RewriteResult, error) {
	return ops.RewriteWithMetadata(input, opts)
}

// Translate translates text to a target language.
//
// Example:
//
//	translated, err := schemaflow.Translate(text, schemaflow.NewTranslateOptions().WithTargetLanguage("Spanish"))
func Translate(input string, opts TranslateOptions) (string, error) {
	return ops.Translate(input, opts)
}

// TranslateWithMetadata translates text and returns metadata including
// detected source language, confidence, and alternative translations.
//
// Example:
//
//	result, err := schemaflow.TranslateWithMetadata(text, schemaflow.NewTranslateOptions().WithTargetLanguage("French"))
//	fmt.Printf("Translation: %s\nDetected language: %s\nConfidence: %.0f%%\n",
//	    result.Text, result.SourceLanguageDetected, result.Confidence*100)
func TranslateWithMetadata(input string, opts TranslateOptions) (TranslateResult, error) {
	return ops.TranslateWithMetadata(input, opts)
}

// Expand expands brief text into more detailed content.
//
// Example:
//
//	expanded, err := schemaflow.Expand(briefText, schemaflow.NewExpandOptions().WithTargetLength(500))
func Expand(input string, opts ExpandOptions) (string, error) {
	return ops.Expand(input, opts)
}

// ExpandWithMetadata expands text and returns metadata including
// expansion ratio, what content was added, and confidence score.
//
// Example:
//
//	result, err := schemaflow.ExpandWithMetadata(briefText, schemaflow.NewExpandOptions())
//	fmt.Printf("Expanded: %s\nExpansion ratio: %.1fx\nAdded: %v\n",
//	    result.Text, result.ExpansionRatio, result.AddedContent)
func ExpandWithMetadata(input string, opts ExpandOptions) (ExpandResult, error) {
	return ops.ExpandWithMetadata(input, opts)
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

// Re-export LLM-powered Redact types and functions
type (
	RedactLLMOptions = ops.RedactLLMOptions
	RedactLLMResult  = ops.RedactLLMResult
	RedactSpan       = ops.RedactSpan
)

var NewRedactLLMOptions = ops.NewRedactLLMOptions

// RedactLLM uses LLM to intelligently identify and mask sensitive data at character level.
// Unlike the regex-based Redact, this uses AI to understand context and detect sensitive
// information that patterns might miss.
//
// Example:
//
//	result, err := schemaflow.RedactLLM("Contact john@email.com for support",
//	    schemaflow.NewRedactLLMOptions().
//	        WithCategories([]string{"email", "phone"}).
//	        WithMaskChar('*').
//	        WithShowFirst(2).
//	        WithShowLast(2))
//	// result.Text = "Contact jo**********om for support"
//	// result.Spans = [{Start: 8, End: 22, Category: "email", Original: "john@email.com"}]
func RedactLLM(text string, opts RedactLLMOptions) (RedactLLMResult, error) {
	return ops.RedactLLM(context.Background(), text, opts)
}

// Re-export Complete types and functions
type (
	CompleteOptions = ops.CompleteOptions
	CompleteResult  = ops.CompleteResult
)

var NewCompleteOptions = ops.NewCompleteOptions

// NewCompleteFieldOptions creates options for completing a specific field in a struct.
var NewCompleteFieldOptions = ops.NewCompleteFieldOptions

// CompleteFieldResult is the result of completing a field in a struct.
type CompleteFieldResult[T any] = ops.CompleteFieldResult[T]

// CompleteFieldOptions configures the CompleteField operation.
type CompleteFieldOptions = ops.CompleteFieldOptions

// Complete intelligently completes partial text using LLM intelligence.
// Note: This is a simplified wrapper that uses the default provider.
//
// Example:
//
//	result, err := schemaflow.Complete(partialText, schemaflow.NewCompleteOptions())
func Complete(partialText string, opts CompleteOptions) (CompleteResult, error) {
	return ops.Complete(context.Background(), nil, partialText, opts)
}

// CompleteField completes a specific string field in a struct and returns a new copy
// with the completed field. The struct's other fields provide context for completion.
//
// Example:
//
//	type BlogPost struct {
//	    Title string `json:"title"`
//	    Body  string `json:"body"`
//	}
//	post := BlogPost{Title: "AI in Healthcare", Body: "Artificial intelligence is transforming"}
//	result, err := schemaflow.CompleteField[BlogPost](post, schemaflow.NewCompleteFieldOptions("Body"))
//	// result.Data.Body now contains the completed text
func CompleteField[T any](data T, opts CompleteFieldOptions) (CompleteFieldResult[T], error) {
	return ops.CompleteField[T](context.Background(), nil, data, opts)
}

// Validate checks if data meets specified validation rules with rich result.
//
// Example:
//
//	result, err := schemaflow.Validate[Person](person, schemaflow.NewValidateOptions().
//	    WithRules("age must be 18-100, email must be valid"))
//	if !result.Valid {
//	    for _, err := range result.Errors {
//	        fmt.Printf("Error: %s\n", err.Message)
//	    }
//	}
func Validate[T any](data T, opts ValidateOptions) (ValidateResult[T], error) {
	return ops.Validate(data, opts)
}

// ValidateLegacy is the legacy validation function for backward compatibility.
//
// Example:
//
//	result, err := schemaflow.ValidateLegacy(person, "age must be 18-100")
func ValidateLegacy[T any](data T, rules string, opts ...OpOptions) (ValidationResult, error) {
	return ops.ValidateLegacy(data, rules, opts...)
}

// Question answers questions about data and returns a typed answer.
//
// Example:
//
//	result, err := schemaflow.Question[Report, string](report, schemaflow.NewQuestionOptions("What is the main finding?"))
//	fmt.Println(result.Answer, "confidence:", result.Confidence)
func Question[T any, A any](data T, opts QuestionOptions) (QuestionResult[A], error) {
	return ops.Question[T, A](data, opts)
}

// QuestionLegacy answers questions about data using the legacy string interface.
//
// Example:
//
//	answer, err := schemaflow.QuestionLegacy(report, "What are the top 3 risks?")
func QuestionLegacy(data any, question string, opts ...OpOptions) (string, error) {
	return ops.QuestionLegacy(data, question, opts...)
}

// Merge combines multiple data sources using a specified strategy.
//
// Example:
//
//	merged, err := schemaflow.Merge(sources, "first-wins")
func Merge[T any](sources []T, strategy string, opts ...OpOptions) (T, error) {
	return ops.Merge(sources, strategy, opts...)
}

// MergeWithMetadata combines data sources and returns metadata including
// which sources were used, any conflicts, and how they were resolved.
//
// Example:
//
//	result, err := schemaflow.MergeWithMetadata(sources, "prefer-newest")
//	fmt.Printf("Merged: %+v\nConflicts: %d\nConfidence: %.0f%%\n",
//	    result.Merged, len(result.Conflicts), result.Confidence*100)
func MergeWithMetadata[T any](sources []T, strategy string, opts ...OpOptions) (MergeResult[T], error) {
	return ops.MergeWithMetadata(sources, strategy, opts...)
}

// Format converts data to a specific output format.
//
// Example:
//
//	formatted, err := schemaflow.Format(data, "markdown table")
func Format(data any, template string, opts ...OpOptions) (string, error) {
	return ops.Format(data, template, opts...)
}

// FormatWithMetadata formats data and returns metadata including
// what format was applied, transformation notes, and confidence score.
//
// Example:
//
//	result, err := schemaflow.FormatWithMetadata(data, "professional bio in third person")
//	fmt.Printf("Formatted: %s\nFormat applied: %s\n", result.Text, result.FormatApplied)
func FormatWithMetadata(data any, template string, opts ...OpOptions) (FormatResult, error) {
	return ops.FormatWithMetadata(data, template, opts...)
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

// === New LLM Operations (v2) ===

// Annotate extracts semantic annotations (entities, sentiments, topics) from text.
//
// Example:
//
//	result, err := schemaflow.Annotate(document, schemaflow.NewAnnotateOptions().WithTypes([]string{"entities", "sentiment"}))
func Annotate[T any](input T, opts AnnotateOptions) (AnnotateResult, error) {
	return ops.Annotate(input, opts)
}

// Cluster groups items by semantic similarity into natural clusters.
//
// Example:
//
//	result, err := schemaflow.Cluster(documents, schemaflow.NewClusterOptions().WithMaxClusters(5))
func Cluster[T any](items []T, opts ClusterOptions) (ClusterResult[T], error) {
	return ops.Cluster(items, opts)
}

// Rank orders items by relevance to a query using semantic understanding.
//
// Example:
//
//	result, err := schemaflow.Rank(documents, schemaflow.NewRankOptions().WithQuery("machine learning"))
func Rank[T any](items []T, opts RankOptions) (RankResult[T], error) {
	return ops.Rank(items, opts)
}

// Compress reduces content while preserving essential meaning.
//
// Example:
//
//	result, err := schemaflow.Compress(document, schemaflow.NewCompressOptions().WithRatio(0.3))
func Compress[T any](input T, opts CompressOptions) (CompressResult[T], error) {
	return ops.Compress(input, opts)
}

// CompressText compresses text content using semantic compression.
//
// Example:
//
//	compressed, err := schemaflow.CompressText("long text here...", schemaflow.NewCompressOptions().WithRatio(0.5))
func CompressText(input string, opts CompressOptions) (string, error) {
	return ops.CompressText(input, opts)
}

// Decompose breaks down complex items into atomic parts using LLM intelligence.
//
// Example:
//
//	result, err := schemaflow.Decompose(complexTask, schemaflow.NewDecomposeOptions().WithMode("hierarchical"))
func Decompose[T any](input T, opts DecomposeOptions) (DecomposeResult[T], error) {
	return ops.Decompose(input, opts)
}

// DecomposeToSlice breaks down complex items and returns just the parts.
//
// Example:
//
//	parts, err := schemaflow.DecomposeToSlice[Task, SubTask](complexTask, schemaflow.NewDecomposeOptions())
func DecomposeToSlice[T any, U any](input T, opts DecomposeOptions) ([]U, error) {
	return ops.DecomposeToSlice[T, U](input, opts)
}

// Enrich adds derived or inferred fields to existing data.
//
// Example:
//
//	result, err := schemaflow.Enrich[Person, EnrichedPerson](person, schemaflow.NewEnrichOptions().WithFields([]string{"age_bracket", "generation"}))
func Enrich[T any, U any](input T, opts EnrichOptions) (EnrichResult[U], error) {
	return ops.Enrich[T, U](input, opts)
}

// EnrichInPlace enriches data in place, returning the same type.
//
// Example:
//
//	enriched, err := schemaflow.EnrichInPlace(person, schemaflow.NewEnrichOptions())
func EnrichInPlace[T any](input T, opts EnrichOptions) (T, error) {
	return ops.EnrichInPlace(input, opts)
}

// Normalize standardizes data formats and values for consistency.
//
// Example:
//
//	result, err := schemaflow.Normalize(data, schemaflow.NewNormalizeOptions().WithRules([]string{"dates", "phone_numbers"}))
func Normalize[T any](input T, opts NormalizeOptions) (NormalizeResult[T], error) {
	return ops.Normalize(input, opts)
}

// NormalizeText normalizes text content.
//
// Example:
//
//	normalized, err := schemaflow.NormalizeText("messy text...", schemaflow.NewNormalizeOptions())
func NormalizeText(input string, opts NormalizeOptions) (string, error) {
	return ops.NormalizeText(input, opts)
}

// NormalizeBatch normalizes multiple items at once.
//
// Example:
//
//	results, err := schemaflow.NormalizeBatch(items, schemaflow.NewNormalizeOptions())
func NormalizeBatch[T any](items []T, opts NormalizeOptions) ([]NormalizeResult[T], error) {
	return ops.NormalizeBatch(items, opts)
}

// SemanticMatch pairs items from two sets based on semantic similarity.
// Note: This is different from the control-flow Match operation.
//
// Example:
//
//	result, err := schemaflow.SemanticMatch(resumes, jobs, schemaflow.NewMatchOptions().WithThreshold(0.7))
func SemanticMatch[S any, T any](sources []S, targets []T, opts MatchOptions) (MatchResult[S, T], error) {
	return ops.SemanticMatch(sources, targets, opts)
}

// MatchOne finds the best matching targets for a single source item.
//
// Example:
//
//	matches, err := schemaflow.MatchOne(resume, jobs, schemaflow.NewMatchOptions())
func MatchOne[S any, T any](source S, targets []T, opts MatchOptions) ([]MatchPair[S, T], error) {
	return ops.MatchOne(source, targets, opts)
}

// Critique provides constructive feedback on content quality.
//
// Example:
//
//	result, err := schemaflow.Critique(essay, schemaflow.NewCritiqueOptions().WithAspects([]string{"clarity", "argument_strength"}))
func Critique[T any](input T, opts CritiqueOptions) (CritiqueResult, error) {
	return ops.Critique(input, opts)
}

// Synthesize combines information from multiple sources into a coherent whole.
//
// Example:
//
//	result, err := schemaflow.Synthesize[Summary](articles, schemaflow.NewSynthesizeOptions().WithPerspective("balanced"))
func Synthesize[T any](sources []any, opts SynthesizeOptions) (SynthesizeResult[T], error) {
	return ops.Synthesize[T](sources, opts)
}

// Predict forecasts future states or outcomes based on patterns.
//
// Example:
//
//	result, err := schemaflow.Predict[Forecast](salesData, schemaflow.NewPredictOptions().WithHorizon("next_quarter"))
func Predict[T any](historicalData any, opts PredictOptions) (PredictResult[T], error) {
	return ops.Predict[T](historicalData, opts)
}

// Verify checks facts, logic, and consistency in content.
//
// Example:
//
//	result, err := schemaflow.Verify("The Earth is flat.", schemaflow.NewVerifyOptions().WithMode("factual"))
func Verify(input string, opts VerifyOptions) (VerifyResult, error) {
	return ops.Verify(input, opts)
}

// VerifyClaim checks a specific claim.
//
// Example:
//
//	result, err := schemaflow.VerifyClaim("GDP grew 5%", schemaflow.NewVerifyOptions())
func VerifyClaim(claim string, opts VerifyOptions) (ClaimVerification, error) {
	return ops.VerifyClaim(claim, opts)
}

// === Data-Centric LLM Operations (v3) ===

// Negotiate reconciles competing constraints to find an optimal solution.
//
// Type parameter T specifies the output solution type.
//
// Example:
//
//	result, err := schemaflow.Negotiate[Schedule](constraints)
//	result, err := schemaflow.Negotiate[Schedule](constraints, schemaflow.NegotiateOptions{
//	    Strategy: "balanced",
//	})
func Negotiate[T any](constraints any, opts ...NegotiateOptions) (NegotiateResult[T], error) {
	return ops.Negotiate[T](constraints, opts...)
}

// AdversarialPosition represents one party's position in a negotiation.
type AdversarialPosition[T any] = ops.AdversarialPosition[T]

// AdversarialContext provides the negotiation dynamics between two parties.
type AdversarialContext[T any] = ops.AdversarialContext[T]

// TermMovement tracks how a specific term moved during negotiation.
type TermMovement = ops.TermMovement

// AdversarialResult contains the outcome of an adversarial negotiation.
type AdversarialResult[T any] = ops.AdversarialResult[T]

// AdversarialOptions configures the adversarial negotiation.
type AdversarialOptions = ops.AdversarialOptions

// NegotiateAdversarial conducts a two-party adversarial negotiation.
//
// This models real-world negotiations where two parties with opposing interests
// must find common ground. The leverage parameter determines who has more power
// and thus who should concede more.
//
// Type parameter T specifies the structure of positions and the final deal.
//
// Example:
//
//	type SalaryTerms struct {
//	    BaseSalary int `json:"base_salary"`
//	    RemoteDays int `json:"remote_days"`
//	}
//	ctx := schemaflow.AdversarialContext[SalaryTerms]{
//	    Ours:        schemaflow.AdversarialPosition[SalaryTerms]{Position: SalaryTerms{BaseSalary: 160000, RemoteDays: 5}},
//	    Theirs:      schemaflow.AdversarialPosition[SalaryTerms]{Position: SalaryTerms{BaseSalary: 130000, RemoteDays: 2}},
//	    OurLeverage: "strong",
//	}
//	result, err := schemaflow.NegotiateAdversarial[SalaryTerms](ctx)
//	// result.Deal has the final terms
//	// result.TermMovements shows who moved on each term
//	// result.WhoConcededMore indicates "they" since we had strong leverage
func NegotiateAdversarial[T any](context AdversarialContext[T], opts ...AdversarialOptions) (AdversarialResult[T], error) {
	return ops.NegotiateAdversarial[T](context, opts...)
}

// Resolve resolves conflicts when multiple typed sources disagree.
//
// Type parameter T specifies the type of sources and the resolved output.
//
// Example:
//
//	result, err := schemaflow.Resolve(conflictingSources)
//	result, err := schemaflow.Resolve(conflictingSources, schemaflow.ResolveOptions{
//	    Strategy: "most-complete",
//	})
func Resolve[T any](sources []T, opts ...ResolveOptions) (ResolveResult[T], error) {
	return ops.Resolve[T](sources, opts...)
}

// Derive infers new typed fields from existing data.
//
// Type parameter T specifies the input type.
// Type parameter U specifies the output type with derived fields.
//
// Example:
//
//	result, err := schemaflow.Derive[Person, EnrichedPerson](person)
//	result, err := schemaflow.Derive[Person, EnrichedPerson](person, schemaflow.DeriveOptions{
//	    TargetFields: []string{"age_category", "generation"},
//	})
func Derive[T any, U any](input T, opts ...DeriveOptions) (DeriveResult[U], error) {
	return ops.Derive[T, U](input, opts...)
}

// Conform transforms data to match specific standards (USPS, ISO8601, E164, etc.).
//
// Type parameter T specifies the data type to conform.
//
// Example:
//
//	result, err := schemaflow.Conform(address, "USPS")
//	result, err := schemaflow.Conform(phoneData, "E164", schemaflow.ConformOptions{
//	    Strict: true,
//	})
func Conform[T any](input T, standard string, opts ...ConformOptions) (ConformResult[T], error) {
	return ops.Conform[T](input, standard, opts...)
}

// Interpolate fills gaps in typed sequences intelligently.
//
// Type parameter T specifies the type of sequence items.
//
// Example:
//
//	result, err := schemaflow.Interpolate(timeSeriesData)
//	result, err := schemaflow.Interpolate(sparseRecords, schemaflow.InterpolateOptions{
//	    Method: "contextual",
//	})
func Interpolate[T any](items []T, opts ...InterpolateOptions) (InterpolateResult[T], error) {
	return ops.Interpolate[T](items, opts...)
}

// Arbitrate makes rule-based decisions with full audit trail.
//
// Type parameter T specifies the type of options to choose from.
//
// Example:
//
//	result, err := schemaflow.Arbitrate(candidates)
//	result, err := schemaflow.Arbitrate(candidates, schemaflow.ArbitrateOptions{
//	    Rules: []string{"must have 3+ years experience", "prefer local candidates"},
//	})
func Arbitrate[T any](options []T, opts ...ArbitrateOptions) (ArbitrateResult[T], error) {
	return ops.Arbitrate[T](options, opts...)
}

// Project transforms structure while preserving semantics.
//
// Type parameter T specifies the input type.
// Type parameter U specifies the projected output type.
//
// Example:
//
//	result, err := schemaflow.Project[Order, OrderSummary](order)
//	result, err := schemaflow.Project[UserProfile, PublicProfile](profile, schemaflow.ProjectOptions{
//	    Mappings: map[string]string{"full_name": "display_name"},
//	    Exclude:  []string{"password_hash", "ssn"},
//	})
func Project[T any, U any](input T, opts ...ProjectOptions) (ProjectResult[U], error) {
	return ops.Project[T, U](input, opts...)
}

// Audit performs deep inspection for issues, anomalies, and policy violations.
//
// Type parameter T specifies the type of data to audit.
//
// Example:
//
//	result, err := schemaflow.Audit(customerRecord)
//	result, err := schemaflow.Audit(financialData, schemaflow.AuditOptions{
//	    Policies:   []string{"PII must be encrypted", "Amounts must balance"},
//	    Categories: []string{"security", "compliance"},
//	})
func Audit[T any](data T, opts ...AuditOptions) (AuditResult[T], error) {
	return ops.Audit[T](data, opts...)
}

// Assemble builds a complex typed object from multiple parts.
//
// Type parameter T specifies the target type to compose.
//
// Example:
//
//	result, err := schemaflow.Assemble[UserProfile]([]any{basicInfo, addressData, preferences})
//	result, err := schemaflow.Assemble[Document](parts, schemaflow.ComposeOptions{
//	    MergeStrategy: "smart",
//	    FillGaps:      true,
//	})
func Assemble[T any](parts []any, opts ...ComposeOptions) (ComposeResult[T], error) {
	return ops.Assemble[T](parts, opts...)
}

// Pivot restructures data relationships between typed objects.
//
// Type parameter T specifies the input type.
// Type parameter U specifies the pivoted output type.
//
// Example:
//
//	result, err := schemaflow.Pivot[[]SalesRow, []SalesPivot](sales, schemaflow.PivotOptions{
//	    PivotOn:   []string{"Month"},
//	    Aggregate: "sum",
//	})
//	result, err := schemaflow.Pivot[Nested, Flat](data, schemaflow.PivotOptions{
//	    Flatten: true,
//	})
func Pivot[T any, U any](input T, opts ...PivotOptions) (PivotResult[U], error) {
	return ops.Pivot[T, U](input, opts...)
}
