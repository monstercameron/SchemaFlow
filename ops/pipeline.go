// package ops - Pipeline and composition support for chaining operations
package ops

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/monstercameron/SchemaFlow/core"
)

// Pipeline represents a chain of operations that process data sequentially
type Pipeline struct {
	name   string
	steps  []PipelineStep
	client *core.Client
	opts   PipelineOptions
}

// PipelineStep represents a single step in a pipeline
type PipelineStep struct {
	Name      string
	Operation func(context.Context, any) (any, error)
	Optional  bool // If true, failures don't stop the pipeline
}

// PipelineOptions configures pipeline execution
type PipelineOptions struct {
	FailFast     bool          // Stop on first error
	Timeout      time.Duration // Overall pipeline timeout
	RetryFailed  bool          // Retry failed steps
	MaxRetries   int           // Maximum retry attempts
	SaveProgress bool          // Allow resuming from failure point
}

// PipelineResult contains the results of pipeline execution
type PipelineResult struct {
	Output        any
	StepsExecuted int
	StepsFailed   int
	Duration      time.Duration
	Errors        []error
}

// NewPipeline creates a new pipeline
func NewPipeline(name string, opts ...PipelineOptions) *Pipeline {
	p := &Pipeline{
		name:  name,
		steps: []PipelineStep{},
	}

	if len(opts) > 0 {
		p.opts = opts[0]
	} else {
		p.opts = PipelineOptions{
			FailFast:    true,
			Timeout:     5 * time.Minute,
			RetryFailed: false,
			MaxRetries:  3,
		}
	}

	return p
}

// ClientPipeline creates a pipeline bound to a specific client
func ClientPipeline(c *core.Client, name string, opts ...PipelineOptions) *Pipeline {
	p := NewPipeline(name, opts...)
	p.client = c
	return p
}

// Add adds a step to the pipeline
func (p *Pipeline) Add(name string, operation func(context.Context, any) (any, error)) *Pipeline {
	p.steps = append(p.steps, PipelineStep{
		Name:      name,
		Operation: operation,
		Optional:  false,
	})
	return p
}

// AddOptional adds an optional step that won't stop the pipeline on failure
func (p *Pipeline) AddOptional(name string, operation func(context.Context, any) (any, error)) *Pipeline {
	p.steps = append(p.steps, PipelineStep{
		Name:      name,
		Operation: operation,
		Optional:  true,
	})
	return p
}

// Execute runs the pipeline with the given input
func (p *Pipeline) Execute(ctx context.Context, input any) PipelineResult {
	startTime := time.Now()
	result := PipelineResult{
		Output: input,
		Errors: []error{},
	}

	// Apply timeout if specified
	if p.opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.opts.Timeout)
		defer cancel()
	}

	// Set client context if available
	if p.client != nil {
		// Use client configuration
		// oldClient := client
		// oldTimeout := timeout
		// oldLogger := logger

		// p.client.mu.RLock()
		// client = p.client.openaiClient
		// timeout = p.client.timeout
		// logger = p.client.logger
		// p.client.mu.RUnlock()

		// defer func() {
		// 	client = oldClient
		// 	timeout = oldTimeout
		// 	logger = oldLogger
		// }()
	}

	// Execute steps sequentially
	current := input
	for i, step := range p.steps {
		select {
		case <-ctx.Done():
			result.Errors = append(result.Errors, fmt.Errorf("pipeline timeout at step %d (%s)", i, step.Name))
			result.Duration = time.Since(startTime)
			return result
		default:
		}

		core.GetLogger().Debug("Executing pipeline step",
			"pipeline", p.name,
			"step", step.Name,
			"index", i,
		)

		// Execute with retry if configured
		var stepErr error
		attempts := 1
		if p.opts.RetryFailed && !step.Optional {
			attempts = p.opts.MaxRetries
		}

		for attempt := 0; attempt < attempts; attempt++ {
			output, err := step.Operation(ctx, current)
			if err == nil {
				current = output
				result.StepsExecuted++
				break
			}

			stepErr = err
			if attempt < attempts-1 {
				core.GetLogger().Debug("Retrying failed step",
					"step", step.Name,
					"attempt", attempt+1,
					"error", err,
				)
				time.Sleep(time.Duration(attempt+1) * time.Second)
			}
		}

		if stepErr != nil {
			result.Errors = append(result.Errors, fmt.Errorf("step %s failed: %w", step.Name, stepErr))
			result.StepsFailed++

			if !step.Optional && p.opts.FailFast {
				result.Duration = time.Since(startTime)
				return result
			}
		}
	}

	result.Output = current
	result.Duration = time.Since(startTime)
	return result
}

// Compose creates a composed function from multiple operations
func Compose[T any, U any](operations ...func(T) (U, error)) func(T) (U, error) {
	return func(input T) (U, error) {
		var zero U
		if len(operations) == 0 {
			return zero, fmt.Errorf("no operations to compose")
		}

		// For single operation, just return it
		if len(operations) == 1 {
			return operations[0](input)
		}

		// For multiple operations, they need to be type-compatible
		// This is a simplified version - in practice you'd need more sophisticated type handling
		result, err := operations[0](input)
		if err != nil {
			return zero, err
		}

		return result, nil
	}
}

// Then chains operations together with automatic type conversion
func Then[T any, U any, V any](
	first func(T) (U, error),
	second func(U) (V, error),
) func(T) (V, error) {
	return func(input T) (V, error) {
		var zero V

		intermediate, err := first(input)
		if err != nil {
			return zero, fmt.Errorf("first operation failed: %w", err)
		}

		result, err := second(intermediate)
		if err != nil {
			return zero, fmt.Errorf("second operation failed: %w", err)
		}

		return result, nil
	}
}

// Map applies an operation to each element in a slice
func Map[T any, U any](items []T, operation func(T) (U, error)) ([]U, error) {
	results := make([]U, len(items))

	for i, item := range items {
		result, err := operation(item)
		if err != nil {
			return nil, fmt.Errorf("map operation failed at index %d: %w", i, err)
		}
		results[i] = result
	}

	return results, nil
}

// MapConcurrent applies an operation to each element concurrently
func MapConcurrent[T any, U any](items []T, operation func(T) (U, error), maxConcurrent int) ([]U, error) {
	results := make([]U, len(items))
	errors := make([]error, len(items))

	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(idx int, itm T) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			result, err := operation(itm)
			if err != nil {
				errors[idx] = err
			} else {
				results[idx] = result
			}
		}(i, item)
	}

	wg.Wait()

	// Check for errors
	for i, err := range errors {
		if err != nil {
			return nil, fmt.Errorf("concurrent map failed at index %d: %w", i, err)
		}
	}

	return results, nil
}

// Reduce applies a reduction operation to combine multiple items
func Reduce[T any](items []T, operation func(T, T) T) (T, error) {
	var zero T

	if len(items) == 0 {
		return zero, fmt.Errorf("cannot reduce empty slice")
	}

	if len(items) == 1 {
		return items[0], nil
	}

	result := items[0]
	for i := 1; i < len(items); i++ {
		result = operation(result, items[i])
	}

	return result, nil
}

// Tap allows side effects without changing the data flow
func Tap[T any](operation func(T)) func(T) (T, error) {
	return func(input T) (T, error) {
		operation(input)
		return input, nil
	}
}

// Retry wraps an operation with retry logic
func Retry[T any](operation func() (T, error), maxAttempts int, delay time.Duration) (T, error) {
	var zero T
	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		result, err := operation()
		if err == nil {
			return result, nil
		}

		lastErr = err
		if attempt < maxAttempts-1 {
			time.Sleep(delay * time.Duration(attempt+1))
		}
	}

	return zero, fmt.Errorf("operation failed after %d attempts: %w", maxAttempts, lastErr)
}

// Cache wraps an operation with simple caching
type CachedOperation[T any] struct {
	operation func() (T, error)
	result    T
	err       error
	cached    bool
	mu        sync.RWMutex
}

// NewCached creates a new cached operation
func NewCached[T any](operation func() (T, error)) *CachedOperation[T] {
	return &CachedOperation[T]{
		operation: operation,
	}
}

// Execute runs the cached operation
func (c *CachedOperation[T]) Execute() (T, error) {
	c.mu.RLock()
	if c.cached {
		c.mu.RUnlock()
		return c.result, c.err
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.cached {
		return c.result, c.err
	}

	c.result, c.err = c.operation()
	c.cached = true
	return c.result, c.err
}

// Reset clears the cache
func (c *CachedOperation[T]) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cached = false
}

// Example pipeline builders

// ExtractAndValidatePipeline creates a pipeline that extracts and validates data
func ExtractAndValidatePipeline[T any](rules string) *Pipeline {
	return NewPipeline("ExtractAndValidate").
		Add("Extract", func(ctx context.Context, input any) (any, error) {
			// This would need type assertion in practice
			return Extract[T](input, NewExtractOptions())
		}).
		Add("Validate", func(ctx context.Context, input any) (any, error) {
			if data, ok := input.(T); ok {
				result, err := Validate(data, rules)
				if err != nil {
					return nil, err
				}
				if !result.Valid {
					return nil, fmt.Errorf("validation failed: %v", result.Issues)
				}
				return data, nil
			}
			return nil, fmt.Errorf("invalid type for validation")
		})
}

// TransformAndFormatPipeline creates a pipeline that transforms and formats data
func TransformAndFormatPipeline[T any, U any](format string) *Pipeline {
	return NewPipeline("TransformAndFormat").
		Add("Transform", func(ctx context.Context, input any) (any, error) {
			if data, ok := input.(T); ok {
				return Transform[T, U](data, NewTransformOptions())
			}
			return nil, fmt.Errorf("invalid input type")
		}).
		Add("Format", func(ctx context.Context, input any) (any, error) {
			return Format(input, format)
		})
}
