package fluent

import (
	"context"
)

// ExtractRequest is a fluent builder for Extract.
type ExtractRequest[T any] struct {
	input any
	opts  ExtractOptions
}

// Extracting starts a fluent Extract request.
func Extracting[T any](input any) ExtractRequest[T] {
	return ExtractRequest[T]{input: input, opts: NewExtractOptions()}
}

func (r ExtractRequest[T]) WithOptions(opts ExtractOptions) ExtractRequest[T] {
	r.opts = opts
	return r
}

func (r ExtractRequest[T]) Configure(fn func(ExtractOptions) ExtractOptions) ExtractRequest[T] {
	r.opts = fn(r.opts)
	return r
}

func (r ExtractRequest[T]) Steer(steering string) ExtractRequest[T] {
	r.opts = r.opts.WithSteering(steering)
	return r
}

func (r ExtractRequest[T]) Threshold(threshold float64) ExtractRequest[T] {
	r.opts = r.opts.WithThreshold(threshold)
	return r
}

func (r ExtractRequest[T]) Strict() ExtractRequest[T] {
	r.opts = r.opts.WithMode(Strict)
	return r
}

func (r ExtractRequest[T]) Smart() ExtractRequest[T] {
	r.opts = r.opts.WithIntelligence(Smart)
	return r
}

func (r ExtractRequest[T]) Fast() ExtractRequest[T] {
	r.opts = r.opts.WithIntelligence(Fast)
	return r
}

func (r ExtractRequest[T]) Quick() ExtractRequest[T] {
	r.opts = r.opts.WithIntelligence(Quick)
	return r
}

func (r ExtractRequest[T]) Context(ctx context.Context) ExtractRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithContext(ctx)
	return r
}

func (r ExtractRequest[T]) RequestID(requestID string) ExtractRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithRequestID(requestID)
	return r
}

func (r ExtractRequest[T]) CorrelationID(correlationID string) ExtractRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithCorrelationID(correlationID)
	return r
}

func (r ExtractRequest[T]) Partial(allow bool) ExtractRequest[T] {
	r.opts = r.opts.WithAllowPartial(allow)
	return r
}

func (r ExtractRequest[T]) SchemaHints(hints map[string]string) ExtractRequest[T] {
	r.opts = r.opts.WithSchemaHints(hints)
	return r
}

func (r ExtractRequest[T]) Run() (T, error) {
	return Extract[T](r.input, r.opts)
}

// TransformRequest is a fluent builder for Transform.
type TransformRequest[T any, U any] struct {
	input T
	opts  TransformOptions
}

// Transforming starts a fluent Transform request.
func Transforming[T any, U any](input T) TransformRequest[T, U] {
	return TransformRequest[T, U]{input: input, opts: NewTransformOptions()}
}

func (r TransformRequest[T, U]) WithOptions(opts TransformOptions) TransformRequest[T, U] {
	r.opts = opts
	return r
}

func (r TransformRequest[T, U]) Configure(fn func(TransformOptions) TransformOptions) TransformRequest[T, U] {
	r.opts = fn(r.opts)
	return r
}

func (r TransformRequest[T, U]) Steer(steering string) TransformRequest[T, U] {
	r.opts = r.opts.WithSteering(steering)
	return r
}

func (r TransformRequest[T, U]) Strict() TransformRequest[T, U] {
	r.opts = r.opts.WithMode(Strict)
	return r
}

func (r TransformRequest[T, U]) Creative() TransformRequest[T, U] {
	r.opts = r.opts.WithMode(Creative)
	return r
}

func (r TransformRequest[T, U]) Smart() TransformRequest[T, U] {
	r.opts = r.opts.WithIntelligence(Smart)
	return r
}

func (r TransformRequest[T, U]) Fast() TransformRequest[T, U] {
	r.opts = r.opts.WithIntelligence(Fast)
	return r
}

func (r TransformRequest[T, U]) Quick() TransformRequest[T, U] {
	r.opts = r.opts.WithIntelligence(Quick)
	return r
}

func (r TransformRequest[T, U]) Merge(strategy string) TransformRequest[T, U] {
	r.opts = r.opts.WithMergeStrategy(strategy)
	return r
}

func (r TransformRequest[T, U]) Context(ctx context.Context) TransformRequest[T, U] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithContext(ctx)
	return r
}

func (r TransformRequest[T, U]) RequestID(requestID string) TransformRequest[T, U] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithRequestID(requestID)
	return r
}

func (r TransformRequest[T, U]) CorrelationID(correlationID string) TransformRequest[T, U] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithCorrelationID(correlationID)
	return r
}

func (r TransformRequest[T, U]) Run() (U, error) {
	return Transform[T, U](r.input, r.opts)
}

// GenerateRequest is a fluent builder for Generate.
type GenerateRequest[T any] struct {
	prompt string
	opts   GenerateOptions
}

// Generating starts a fluent Generate request.
func Generating[T any](prompt string) GenerateRequest[T] {
	return GenerateRequest[T]{prompt: prompt, opts: NewGenerateOptions()}
}

func (r GenerateRequest[T]) WithOptions(opts GenerateOptions) GenerateRequest[T] {
	r.opts = opts
	return r
}

func (r GenerateRequest[T]) Configure(fn func(GenerateOptions) GenerateOptions) GenerateRequest[T] {
	r.opts = fn(r.opts)
	return r
}

func (r GenerateRequest[T]) Steer(steering string) GenerateRequest[T] {
	r.opts = r.opts.WithSteering(steering)
	return r
}

func (r GenerateRequest[T]) Strict() GenerateRequest[T] {
	r.opts = r.opts.WithMode(Strict)
	return r
}

func (r GenerateRequest[T]) Creative() GenerateRequest[T] {
	r.opts = r.opts.WithMode(Creative)
	return r
}

func (r GenerateRequest[T]) Smart() GenerateRequest[T] {
	r.opts = r.opts.WithIntelligence(Smart)
	return r
}

func (r GenerateRequest[T]) Fast() GenerateRequest[T] {
	r.opts = r.opts.WithIntelligence(Fast)
	return r
}

func (r GenerateRequest[T]) Quick() GenerateRequest[T] {
	r.opts = r.opts.WithIntelligence(Quick)
	return r
}

func (r GenerateRequest[T]) Count(count int) GenerateRequest[T] {
	r.opts.Count = count
	return r
}

func (r GenerateRequest[T]) Style(style string) GenerateRequest[T] {
	r.opts.Style = style
	return r
}

func (r GenerateRequest[T]) Context(ctx context.Context) GenerateRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithContext(ctx)
	return r
}

func (r GenerateRequest[T]) RequestID(requestID string) GenerateRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithRequestID(requestID)
	return r
}

func (r GenerateRequest[T]) CorrelationID(correlationID string) GenerateRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithCorrelationID(correlationID)
	return r
}

func (r GenerateRequest[T]) Run() (T, error) {
	return Generate[T](r.prompt, r.opts)
}

// ChooseRequest is a fluent builder for Choose.
type ChooseRequest[T any] struct {
	options []T
	opts    ChooseOptions
}

// Choosing starts a fluent Choose request.
func Choosing[T any](options []T) ChooseRequest[T] {
	return ChooseRequest[T]{options: options, opts: NewChooseOptions()}
}

func (r ChooseRequest[T]) WithOptions(opts ChooseOptions) ChooseRequest[T] {
	r.opts = opts
	return r
}

func (r ChooseRequest[T]) Configure(fn func(ChooseOptions) ChooseOptions) ChooseRequest[T] {
	r.opts = fn(r.opts)
	return r
}

func (r ChooseRequest[T]) By(criteria ...string) ChooseRequest[T] {
	r.opts = r.opts.WithCriteria(criteria)
	return r
}

func (r ChooseRequest[T]) Steer(steering string) ChooseRequest[T] {
	r.opts = r.opts.WithSteering(steering)
	return r
}

func (r ChooseRequest[T]) Smart() ChooseRequest[T] {
	r.opts = r.opts.WithIntelligence(Smart)
	return r
}

func (r ChooseRequest[T]) Fast() ChooseRequest[T] {
	r.opts = r.opts.WithIntelligence(Fast)
	return r
}

func (r ChooseRequest[T]) Quick() ChooseRequest[T] {
	r.opts = r.opts.WithIntelligence(Quick)
	return r
}

func (r ChooseRequest[T]) Context(ctx context.Context) ChooseRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithContext(ctx)
	return r
}

func (r ChooseRequest[T]) RequestID(requestID string) ChooseRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithRequestID(requestID)
	return r
}

func (r ChooseRequest[T]) CorrelationID(correlationID string) ChooseRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithCorrelationID(correlationID)
	return r
}

func (r ChooseRequest[T]) Top(n int) ChooseRequest[T] {
	r.opts = r.opts.WithTopN(n)
	return r
}

func (r ChooseRequest[T]) Reasoning(enabled bool) ChooseRequest[T] {
	r.opts = r.opts.WithRequireReasoning(enabled)
	return r
}

func (r ChooseRequest[T]) Run() (T, error) {
	return Choose[T](r.options, r.opts)
}

// FilterRequest is a fluent builder for Filter.
type FilterRequest[T any] struct {
	items []T
	opts  FilterOptions
}

// Filtering starts a fluent Filter request.
func Filtering[T any](items []T) FilterRequest[T] {
	return FilterRequest[T]{items: items, opts: NewFilterOptions()}
}

func (r FilterRequest[T]) WithOptions(opts FilterOptions) FilterRequest[T] {
	r.opts = opts
	return r
}

func (r FilterRequest[T]) Configure(fn func(FilterOptions) FilterOptions) FilterRequest[T] {
	r.opts = fn(r.opts)
	return r
}

func (r FilterRequest[T]) By(criteria string) FilterRequest[T] {
	r.opts = r.opts.WithCriteria(criteria)
	return r
}

func (r FilterRequest[T]) Steer(steering string) FilterRequest[T] {
	r.opts = r.opts.WithSteering(steering)
	return r
}

func (r FilterRequest[T]) Smart() FilterRequest[T] {
	r.opts = r.opts.WithIntelligence(Smart)
	return r
}

func (r FilterRequest[T]) Fast() FilterRequest[T] {
	r.opts = r.opts.WithIntelligence(Fast)
	return r
}

func (r FilterRequest[T]) Quick() FilterRequest[T] {
	r.opts = r.opts.WithIntelligence(Quick)
	return r
}

func (r FilterRequest[T]) Context(ctx context.Context) FilterRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithContext(ctx)
	return r
}

func (r FilterRequest[T]) RequestID(requestID string) FilterRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithRequestID(requestID)
	return r
}

func (r FilterRequest[T]) CorrelationID(correlationID string) FilterRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithCorrelationID(correlationID)
	return r
}

func (r FilterRequest[T]) KeepMatching(keep bool) FilterRequest[T] {
	r.opts.KeepMatching = keep
	return r
}

func (r FilterRequest[T]) MinConfidence(confidence float64) FilterRequest[T] {
	r.opts = r.opts.WithMinConfidence(confidence)
	return r
}

func (r FilterRequest[T]) Run() ([]T, error) {
	return Filter[T](r.items, r.opts)
}

// SortRequest is a fluent builder for Sort.
type SortRequest[T any] struct {
	items []T
	opts  SortOptions
}

// Sorting starts a fluent Sort request.
func Sorting[T any](items []T) SortRequest[T] {
	return SortRequest[T]{items: items, opts: NewSortOptions()}
}

func (r SortRequest[T]) WithOptions(opts SortOptions) SortRequest[T] {
	r.opts = opts
	return r
}

func (r SortRequest[T]) Configure(fn func(SortOptions) SortOptions) SortRequest[T] {
	r.opts = fn(r.opts)
	return r
}

func (r SortRequest[T]) By(criteria string) SortRequest[T] {
	r.opts = r.opts.WithCriteria(criteria)
	return r
}

func (r SortRequest[T]) Steer(steering string) SortRequest[T] {
	r.opts = r.opts.WithSteering(steering)
	return r
}

func (r SortRequest[T]) Smart() SortRequest[T] {
	r.opts = r.opts.WithIntelligence(Smart)
	return r
}

func (r SortRequest[T]) Fast() SortRequest[T] {
	r.opts = r.opts.WithIntelligence(Fast)
	return r
}

func (r SortRequest[T]) Quick() SortRequest[T] {
	r.opts = r.opts.WithIntelligence(Quick)
	return r
}

func (r SortRequest[T]) Context(ctx context.Context) SortRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithContext(ctx)
	return r
}

func (r SortRequest[T]) RequestID(requestID string) SortRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithRequestID(requestID)
	return r
}

func (r SortRequest[T]) CorrelationID(correlationID string) SortRequest[T] {
	r.opts.CommonOptions = r.opts.CommonOptions.WithCorrelationID(correlationID)
	return r
}

func (r SortRequest[T]) Asc() SortRequest[T] {
	r.opts = r.opts.WithDirection("ascending")
	return r
}

func (r SortRequest[T]) Desc() SortRequest[T] {
	r.opts = r.opts.WithDirection("descending")
	return r
}

func (r SortRequest[T]) Run() ([]T, error) {
	return Sort[T](r.items, r.opts)
}

// ChooseBy is a compact entrypoint for selection-style tasks.
func ChooseBy[T any](options []T, criteria ...string) (T, error) {
	return Choosing(options).By(criteria...).Run()
}

// FilterBy is a compact entrypoint for natural-language filtering.
func FilterBy[T any](items []T, criteria string) ([]T, error) {
	return Filtering(items).By(criteria).Run()
}

// SortBy is a compact entrypoint for natural-language sorting.
func SortBy[T any](items []T, criteria string) ([]T, error) {
	return Sorting(items).By(criteria).Run()
}
