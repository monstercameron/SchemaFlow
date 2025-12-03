package ops

import (
	"context"
	"fmt"

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

	req := llm.CompletionRequest{
		Model:        model,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  float64(temperature),
		MaxTokens:    maxTokens,
	}

	// Set response format if needed (e.g. for strict mode or extraction)
	// This logic might need to be more sophisticated based on the operation
	if opts.Mode == types.Strict {
		// req.ResponseFormat = "json" // Maybe?
	}

	resp, err := provider.Complete(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}
