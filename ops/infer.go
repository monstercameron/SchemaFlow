// package ops - Infer operation for smart missing data inference
package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/monstercameron/SchemaFlow/core"
)

// InferOptions configures the Infer operation
type InferOptions struct {
	core.OpOptions
	Context string // Additional context or known facts to aid inference
}

// NewInferOptions creates InferOptions with defaults
func NewInferOptions() InferOptions {
	return InferOptions{
		OpOptions: core.OpOptions{
			Mode:         core.TransformMode,
			Intelligence: core.Fast,
		},
	}
}

// WithContext sets additional context for inference
func (opts InferOptions) WithContext(context string) InferOptions {
	opts.Context = context
	return opts
}

// WithIntelligence sets the intelligence level
func (opts InferOptions) WithIntelligence(intelligence core.Speed) InferOptions {
	opts.OpOptions.Intelligence = intelligence
	return opts
}

// Validate validates InferOptions
func (opts InferOptions) Validate() error {
	return nil // No specific validation needed
}

// toOpOptions converts InferOptions to core.OpOptions
func (opts InferOptions) toOpOptions() core.OpOptions {
	return opts.OpOptions
}

// Infer smartly fills in missing fields in partial data using LLM inference
//
// Examples:
//
//	// Infer missing fields in a person record
//	complete, err := Infer[Person](Person{Name: "John", Age: 30},
//	    NewInferOptions().WithContext("This person works in tech"))
//
//	// Infer product details from partial information
//	product, err := Infer[Product](Product{Name: "iPhone 15"},
//	    NewInferOptions().WithContext("Latest Apple smartphone released in 2023"))
func Infer[T any](partialData T, opts InferOptions) (T, error) {
	return inferImpl(core.GetDefaultClient(), partialData, opts)
}

// ClientInfer is the client-based version of Infer
func ClientInfer[T any](c *core.Client, partialData T, opts InferOptions) (T, error) {
	return inferImpl(c, partialData, opts)
}

func inferImpl[T any](c *core.Client, partialData T, opts InferOptions) (T, error) {
	logger := core.GetLogger()
	logger.Debug("Starting infer operation", "requestID", opts.RequestID, "dataType", fmt.Sprintf("%T", partialData))

	var result T

	// Validate options
	if err := opts.Validate(); err != nil {
		logger.Error("Infer operation validation failed", "requestID", opts.RequestID, "error", err)
		return result, fmt.Errorf("invalid options: %w", err)
	}

	opt := opts.toOpOptions()
	opt.Client = c

	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	// Get type information
	targetType := reflect.TypeOf(result)
	typeSchema := GenerateTypeSchema(targetType)

	// Marshal partial data to JSON
	partialJSON, err := json.Marshal(partialData)
	if err != nil {
		logger.Error("Infer operation marshal failed", "requestID", opts.RequestID, "error", err)
		return result, fmt.Errorf("failed to marshal partial data: %w", err)
	}

	// Build system prompt
	systemPrompt := fmt.Sprintf(`You are a data inference expert. Given partial data and context, intelligently infer missing fields.

Target schema:
%s

Rules:
- Infer missing fields based on available data and logical reasoning
- Use reasonable defaults when inference is uncertain
- Maintain data consistency and coherence
- Return ONLY valid JSON matching the complete schema`, typeSchema)

	// Build user prompt
	userPrompt := fmt.Sprintf("Complete this partial data by inferring missing fields:\n%s", string(partialJSON))

	if opts.Context != "" {
		userPrompt += fmt.Sprintf("\n\nAdditional context: %s", opts.Context)
	}

	// Call LLM for inference
	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		logger.Error("Infer operation LLM call failed", "requestID", opts.RequestID, "error", err)
		return result, fmt.Errorf("inference failed: %w", err)
	}

	// Parse inferred data
	if err := core.ParseJSON(response, &result); err != nil {
		logger.Error("Infer operation parse failed", "requestID", opts.RequestID, "error", err)
		return result, fmt.Errorf("failed to parse inferred result: %w", err)
	}

	logger.Debug("Infer operation succeeded", "requestID", opts.RequestID)

	return result, nil
}
