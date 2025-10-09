// Package schemaflow provides the main API for SchemaFlow operations.
// This file re-exports the core package for easier imports.
package schemaflow

import (
	"github.com/monstercameron/SchemaFlow/core"
	"github.com/monstercameron/SchemaFlow/ops"
)

// Re-export core types for convenience
type (
	Client    = core.Client
	OpOptions = core.OpOptions
	Speed     = core.Speed
	Mode      = core.Mode
	Logger    = core.Logger
)

// Result is a generic type that must be used with a type parameter
type Result[T any] = core.Result[T]

// Intelligence levels (Speed)
const (
	Quick = core.Quick
	Fast  = core.Fast
	Smart = core.Smart
)

// Modes
const (
	Strict        = core.Strict
	TransformMode = core.TransformMode
	Creative      = core.Creative
)

// Configuration functions
var (
	Init         = core.Init
	InitWithEnv  = core.InitWithEnv
	NewClient    = core.NewClient
	SetDebugMode = core.SetDebugMode
	GetDebugMode = core.GetDebugMode
)

// Operation functions - these are generic and must be used with type parameters
// Example: schemaflow.Extract[Person](input, opts)

// Extract converts unstructured data into strongly-typed Go structs
// This function accepts either core.OpOptions or ops.ExtractOptions for backward compatibility
func Extract[T any](input any, opts ...interface{}) (T, error) {
	if len(opts) == 0 {
		return ops.Extract[T](input, ops.NewExtractOptions())
	}

	switch opt := opts[0].(type) {
	case ops.ExtractOptions:
		return ops.Extract[T](input, opt)
	case OpOptions:
		// Convert OpOptions to ExtractOptions for backward compatibility
		extractOpts := ops.NewExtractOptions()
		// Just set the OpOptions field directly - it has precedence
		extractOpts.OpOptions = opt
		return ops.Extract[T](input, extractOpts)
	default:
		return ops.Extract[T](input, ops.NewExtractOptions())
	}
}

// Transform converts data from one type to another using LLM intelligence
// This function accepts either core.OpOptions or ops.TransformOptions for backward compatibility
func Transform[T any, U any](input T, opts ...interface{}) (U, error) {
	if len(opts) == 0 {
		return ops.Transform[T, U](input, ops.NewTransformOptions())
	}

	switch opt := opts[0].(type) {
	case ops.TransformOptions:
		return ops.Transform[T, U](input, opt)
	case OpOptions:
		// Convert OpOptions to TransformOptions for backward compatibility
		transformOpts := ops.NewTransformOptions()
		transformOpts.OpOptions = opt
		return ops.Transform[T, U](input, transformOpts)
	default:
		return ops.Transform[T, U](input, ops.NewTransformOptions())
	}
}

// Generate creates new data based on templates and examples
// This function accepts either core.OpOptions or ops.GenerateOptions for backward compatibility
func Generate[T any](prompt string, opts ...interface{}) (T, error) {
	if len(opts) == 0 {
		return ops.Generate[T](prompt, ops.NewGenerateOptions())
	}

	switch opt := opts[0].(type) {
	case ops.GenerateOptions:
		return ops.Generate[T](prompt, opt)
	case OpOptions:
		// Convert OpOptions to GenerateOptions for backward compatibility
		generateOpts := ops.NewGenerateOptions()
		generateOpts.OpOptions = opt
		return ops.Generate[T](prompt, generateOpts)
	default:
		return ops.Generate[T](prompt, ops.NewGenerateOptions())
	}
}

// Choose selects the best option from a list based on criteria
// This function accepts either core.OpOptions or ops.ChooseOptions for backward compatibility
func Choose[T any](options []T, opts ...interface{}) (T, error) {
	if len(opts) == 0 {
		return ops.Choose(options, ops.NewChooseOptions())
	}

	switch opt := opts[0].(type) {
	case ops.ChooseOptions:
		return ops.Choose(options, opt)
	case OpOptions:
		// Convert OpOptions to ChooseOptions for backward compatibility
		chooseOpts := ops.NewChooseOptions()
		chooseOpts.OpOptions = opt
		return ops.Choose(options, chooseOpts)
	default:
		return ops.Choose(options, ops.NewChooseOptions())
	}
}

// Filter filters items based on natural language criteria
// This function accepts either core.OpOptions or ops.FilterOptions for backward compatibility
func Filter[T any](items []T, opts ...interface{}) ([]T, error) {
	if len(opts) == 0 {
		return ops.Filter(items, ops.NewFilterOptions())
	}

	switch opt := opts[0].(type) {
	case ops.FilterOptions:
		return ops.Filter(items, opt)
	case OpOptions:
		// Convert OpOptions to FilterOptions for backward compatibility
		filterOpts := ops.NewFilterOptions()
		filterOpts.OpOptions = opt
		return ops.Filter(items, filterOpts)
	default:
		return ops.Filter(items, ops.NewFilterOptions())
	}
}

// Sort sorts items based on natural language criteria
// This function accepts either core.OpOptions or ops.SortOptions for backward compatibility
func Sort[T any](items []T, opts ...interface{}) ([]T, error) {
	if len(opts) == 0 {
		return ops.Sort(items, ops.NewSortOptions())
	}

	switch opt := opts[0].(type) {
	case ops.SortOptions:
		return ops.Sort(items, opt)
	case OpOptions:
		// Convert OpOptions to SortOptions for backward compatibility
		sortOpts := ops.NewSortOptions()
		sortOpts.OpOptions = opt
		return ops.Sort(items, sortOpts)
	default:
		return ops.Sort(items, ops.NewSortOptions())
	}
}

// Infer fills in missing fields in partial data using LLM intelligence
// This function accepts either core.OpOptions or ops.InferOptions for backward compatibility
func Infer[T any](partialData T, opts ...interface{}) (T, error) {
	if len(opts) == 0 {
		return ops.Infer[T](partialData, ops.NewInferOptions())
	}

	switch opt := opts[0].(type) {
	case ops.InferOptions:
		return ops.Infer[T](partialData, opt)
	case OpOptions:
		// Convert OpOptions to InferOptions for backward compatibility
		inferOpts := ops.NewInferOptions()
		inferOpts.OpOptions = opt
		return ops.Infer[T](partialData, inferOpts)
	default:
		return ops.Infer[T](partialData, ops.NewInferOptions())
	}
}

// Diff compares two data instances and explains the differences intelligently
// This function accepts either core.OpOptions or ops.DiffOptions for backward compatibility
func Diff[T any](oldData, newData T, opts ...interface{}) (ops.DiffResult, error) {
	if len(opts) == 0 {
		return ops.Diff[T](oldData, newData, ops.NewDiffOptions())
	}

	switch opt := opts[0].(type) {
	case ops.DiffOptions:
		return ops.Diff[T](oldData, newData, opt)
	case OpOptions:
		// Convert OpOptions to DiffOptions for backward compatibility
		diffOpts := ops.NewDiffOptions()
		diffOpts.OpOptions = opt
		return ops.Diff[T](oldData, newData, diffOpts)
	default:
		return ops.Diff[T](oldData, newData, ops.NewDiffOptions())
	}
}

// Explain generates human explanations for complex data or code in simple terms
// This function accepts either core.OpOptions or ops.ExplainOptions for backward compatibility
func Explain(data any, opts ...interface{}) (ops.ExplainResult, error) {
	if len(opts) == 0 {
		return ops.Explain(data, ops.NewExplainOptions())
	}

	switch opt := opts[0].(type) {
	case ops.ExplainOptions:
		return ops.Explain(data, opt)
	case OpOptions:
		// Convert OpOptions to ExplainOptions for backward compatibility
		explainOpts := ops.NewExplainOptions()
		explainOpts.OpOptions = opt
		return ops.Explain(data, explainOpts)
	default:
		return ops.Explain(data, ops.NewExplainOptions())
	}
}

// Parse intelligently parses data from various formats into strongly-typed Go structs
// This function accepts either core.OpOptions or ops.ParseOptions for backward compatibility
func Parse[T any](input any, opts ...interface{}) (ops.ParseResult[T], error) {
	if len(opts) == 0 {
		return ops.Parse[T](input, ops.NewParseOptions())
	}

	switch opt := opts[0].(type) {
	case ops.ParseOptions:
		return ops.Parse[T](input, opt)
	case OpOptions:
		// Convert OpOptions to ParseOptions for backward compatibility
		parseOpts := ops.NewParseOptions()
		parseOpts.OpOptions = opt
		return ops.Parse[T](input, parseOpts)
	default:
		return ops.Parse[T](input, ops.NewParseOptions())
	}
}

// Complete intelligently completes partial text using LLM intelligence
// This function accepts either core.OpOptions or ops.CompleteOptions for backward compatibility
func Complete(partialText string, opts ...interface{}) (ops.CompleteResult, error) {
	if len(opts) == 0 {
		return ops.Complete(partialText, ops.NewCompleteOptions())
	}

	switch opt := opts[0].(type) {
	case ops.CompleteOptions:
		return ops.Complete(partialText, opt)
	case OpOptions:
		// Convert OpOptions to CompleteOptions for backward compatibility
		completeOpts := ops.NewCompleteOptions()
		completeOpts.OpOptions = opt
		return ops.Complete(partialText, completeOpts)
	default:
		return ops.Complete(partialText, ops.NewCompleteOptions())
	}
}
