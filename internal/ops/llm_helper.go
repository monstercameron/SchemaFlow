package ops

import (
	"context"
	"fmt"
	"strings"

	"github.com/monstercameron/SchemaFlow/internal/config"
	"github.com/monstercameron/SchemaFlow/internal/llm"
	"github.com/monstercameron/SchemaFlow/internal/types"
)

var defaultProvider llm.Provider

// LLMCaller is the function type for calling the LLM
type LLMCaller func(ctx context.Context, system, user string, opts types.OpOptions) (string, error)

// Custom LLM caller for testing
var customLLMCaller LLMCaller

// setLLMCaller sets a custom LLM caller (for testing)
func setLLMCaller(caller LLMCaller) {
	customLLMCaller = caller
}

// SetDefaultProvider sets the default LLM provider for operations
func SetDefaultProvider(p llm.Provider) {
	defaultProvider = p
}

// callLLM executes an LLM request using the default provider
func callLLM(ctx context.Context, systemPrompt, userPrompt string, opts types.OpOptions) (string, error) {
	// Use custom caller if set (for testing)
	if customLLMCaller != nil {
		return customLLMCaller(ctx, systemPrompt, userPrompt, opts)
	}

	if defaultProvider == nil {
		// Try to initialize a default provider (e.g. OpenAI from env)
		// For now, just return error if not set
		return "", fmt.Errorf("no LLM provider configured")
	}
	return CallLLM(ctx, defaultProvider, systemPrompt, userPrompt, opts)
}

// CallLLM executes an LLM request using the provided provider
func CallLLM(ctx context.Context, provider llm.Provider, systemPrompt, userPrompt string, opts types.OpOptions) (string, error) {
	// Determine model
	model := config.GetModel(opts.Intelligence, provider.Name())
	maxTokens := config.GetMaxTokens(opts.Intelligence)
	temperature := config.GetTemperature(opts.Mode)
	effectiveSystemPrompt := applySteering(systemPrompt, opts.Steering)
	responseFormat := inferResponseFormat(effectiveSystemPrompt, userPrompt)

	req := llm.CompletionRequest{
		Model:          model,
		SystemPrompt:   strengthenSystemPrompt(effectiveSystemPrompt, responseFormat),
		UserPrompt:     userPrompt,
		Temperature:    float64(temperature),
		MaxTokens:      maxTokens,
		ResponseFormat: responseFormat,
	}

	resp, err := provider.Complete(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

func applySteering(systemPrompt, steering string) string {
	steering = strings.TrimSpace(steering)
	if steering == "" {
		return systemPrompt
	}
	return strings.TrimSpace(systemPrompt + "\n\nAdditional instructions:\n" + steering)
}

func inferResponseFormat(systemPrompt, userPrompt string) string {
	combined := strings.ToLower(systemPrompt + "\n" + userPrompt)
	jsonSignals := []string{
		"return a json object",
		"return a json array",
		"return only valid json",
		"return only json",
		"valid json",
		"json object",
		"json array",
		"matching the schema",
	}
	for _, signal := range jsonSignals {
		if strings.Contains(combined, signal) {
			return "json"
		}
	}
	return "text"
}

func strengthenSystemPrompt(systemPrompt, responseFormat string) string {
	baseRules := strings.TrimSpace(`Perform the semantic task faithfully using the provided input.
Do not merely restate schemas, field names, or type descriptions.
Infer, compare, rank, validate, transform, or summarize based on the actual content.`)

	if responseFormat != "json" {
		return strings.TrimSpace(baseRules + "\n\n" + systemPrompt)
	}

	jsonRules := strings.TrimSpace(`After reasoning about the task, return only the final JSON answer.
Do not include markdown fences, prose, placeholders, or schema descriptions.
Every field must be populated with task-relevant values supported by the input or clearly inferred from it.`)

	return strings.TrimSpace(baseRules + "\n\n" + jsonRules + "\n\n" + systemPrompt)
}
