// Package schemaflow - LLM interaction and communication
package schemaflow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// callLLM is the actual LLM call function (can be replaced for testing)
var callLLM callLLMFunc

func init() {
	// Initialize callLLM with default implementation
	callLLM = defaultCallLLM
}

// defaultCallLLM makes a request to the LLM with the given prompts and options
func defaultCallLLM(ctx context.Context, systemPrompt, userPrompt string, opts OpOptions) (string, error) {
	// Check if we should use the new provider-based approach
	if defaultClient != nil && defaultClient.provider != nil {
		return providerCallLLM(ctx, systemPrompt, userPrompt, opts)
	}
	
	// Fallback to legacy OpenAI client
	if client == nil {
		return "", fmt.Errorf("schemaflow not initialized, call Init() first")
	}
	
	// Use operation context or create one
	if opts.context == nil {
		opts.context = context.Background()
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(opts.context, timeout)
	defer cancel()
	
	// Add request ID to context for tracing
	if opts.requestID != "" {
		ctx = context.WithValue(ctx, "requestID", opts.requestID)
	}
	
	model := GetModel(opts.Intelligence)
	maxTokens := getMaxTokens(opts.Intelligence)
	temperature := getTemperature(opts.Mode)
	
	// Log the request if debug is enabled
	if debugMode {
		logger.Debug("LLM request",
			"requestID", opts.requestID,
			"model", model,
			"temperature", temperature,
			"maxTokens", maxTokens,
			"mode", opts.Mode.String(),
			"intelligence", opts.Intelligence.String(),
		)
	}
	
	// Build messages
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userPrompt,
		},
	}
	
	// Add steering if provided
	if opts.Steering != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Additional guidance: " + opts.Steering,
		})
		
		if debugMode {
			logger.Debug("Steering applied", "requestID", opts.requestID, "steering", opts.Steering)
		}
	}
	
	// Retry logic with exponential backoff
	retries := maxRetries
	
	var lastErr error
	backoff := retryBackoff
	
	for attempt := 0; attempt <= retries; attempt++ {
		if attempt > 0 {
			logger.Warn("Retrying LLM request",
				"requestID", opts.requestID,
				"attempt", attempt,
				"maxRetries", retries,
				"backoff", backoff,
			)
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
		}
		
		startTime := time.Now()
		
		// Build request
		request := openai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
		}
		
		// GPT-5 models have specific requirements
		if strings.Contains(model, "gpt-5") {
			// GPT-5 only supports temperature = 1
			request.Temperature = 1.0
			// GPT-5 doesn't support max_tokens parameter
			// (would need MaxCompletionTokens which the library doesn't support yet)
		} else {
			// Use standard parameters for other models
			request.Temperature = temperature
			request.MaxTokens = maxTokens
		}
		
		resp, err := client.CreateChatCompletion(ctx, request)
		
		duration := time.Since(startTime)
		
		// Log metrics if enabled
		if metricsEnabled {
			recordMetric("llm_request_duration", duration.Milliseconds(), map[string]string{
				"model": model,
				"mode": opts.Mode.String(),
				"intelligence": opts.Intelligence.String(),
			})
		}
		
		if err != nil {
			lastErr = err
			logger.Error("LLM request failed",
				"requestID", opts.requestID,
				"attempt", attempt,
				"error", err,
				"duration", duration,
			)
			
			// Check if error is retryable
			if !isRetryableError(err) {
				break
			}
			continue
		}
		
		if len(resp.Choices) == 0 {
			lastErr = fmt.Errorf("no response from LLM")
			logger.Error("Empty LLM response",
				"requestID", opts.requestID,
				"attempt", attempt,
				"duration", duration,
			)
			continue
		}
		
		result := resp.Choices[0].Message.Content
		
		// Log successful response
		if debugMode {
			logger.Debug("LLM response received",
				"requestID", opts.requestID,
				"duration", duration,
				"responseLength", len(result),
				"tokensUsed", resp.Usage.TotalTokens,
			)
		}
		
		return result, nil
	}
	
	return "", fmt.Errorf("failed after %d retries: %w", retries, lastErr)
}

// parseJSON attempts to parse JSON from LLM response, handling common formatting issues
func parseJSON[T any](response string, target *T) error {
	response = strings.TrimSpace(response)
	
	// Remove markdown code blocks if present
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}
	
	// Try standard unmarshaling first
	if err := json.Unmarshal([]byte(response), target); err != nil {
		// Try with a decoder for better error messages
		decoder := json.NewDecoder(strings.NewReader(response))
		decoder.DisallowUnknownFields()
		if decodeErr := decoder.Decode(target); decodeErr != nil {
			logger.Error("JSON parsing failed",
				"error", decodeErr,
				"response", response[:min(len(response), 200)], // Log first 200 chars
			)
			return fmt.Errorf("JSON decode error: %w", decodeErr)
		}
	}
	
	return nil
}

// providerCallLLM makes a request using the provider abstraction
func providerCallLLM(ctx context.Context, systemPrompt, userPrompt string, opts OpOptions) (string, error) {
	provider := defaultClient.provider
	if provider == nil {
		return "", fmt.Errorf("no provider configured")
	}
	
	// Use operation context or create one
	if opts.context == nil {
		opts.context = context.Background()
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(opts.context, timeout)
	defer cancel()
	
	// Add request ID to context for tracing
	if opts.requestID != "" {
		ctx = context.WithValue(ctx, "requestID", opts.requestID)
	}
	
	model := GetModel(opts.Intelligence)
	maxTokens := getMaxTokens(opts.Intelligence)
	temperature := float64(getTemperature(opts.Mode))
	
	// Log the request if debug is enabled
	if debugMode {
		logger.Debug("Provider LLM request",
			"requestID", opts.requestID,
			"provider", provider.Name(),
			"model", model,
			"temperature", temperature,
			"maxTokens", maxTokens,
			"mode", opts.Mode.String(),
			"intelligence", opts.Intelligence.String(),
		)
	}
	
	// Create provider request
	request := CompletionRequest{
		Model:        model,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  temperature,
		MaxTokens:    maxTokens,
	}
	
	// Add response format hint if needed
	if strings.Contains(systemPrompt, "JSON") || strings.Contains(systemPrompt, "json") {
		request.ResponseFormat = "json"
	}
	
	// Add steering if provided
	if opts.Steering != "" {
		request.SystemPrompt += "\n\nAdditional guidance: " + opts.Steering
		
		if debugMode {
			logger.Debug("Steering applied", "requestID", opts.requestID, "steering", opts.Steering)
		}
	}
	
	// Retry logic with exponential backoff
	retries := maxRetries
	var lastErr error
	backoff := retryBackoff
	
	for attempt := 0; attempt <= retries; attempt++ {
		if attempt > 0 {
			logger.Warn("Retrying provider request",
				"requestID", opts.requestID,
				"provider", provider.Name(),
				"attempt", attempt,
				"maxRetries", retries,
				"backoff", backoff,
			)
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
		}
		
		startTime := time.Now()
		
		// Make the provider request
		resp, err := provider.Complete(ctx, request)
		
		duration := time.Since(startTime)
		
		// Log metrics if enabled
		if metricsEnabled {
			recordMetric("provider_request_duration", duration.Milliseconds(), map[string]string{
				"provider": provider.Name(),
				"model": model,
				"mode": opts.Mode.String(),
				"intelligence": opts.Intelligence.String(),
			})
		}
		
		if err != nil {
			lastErr = err
			logger.Error("Provider request failed",
				"requestID", opts.requestID,
				"provider", provider.Name(),
				"attempt", attempt,
				"error", err,
				"duration", duration,
			)
			
			// Check if error is retryable
			if !isRetryableError(err) {
				break
			}
			continue
		}
		
		if resp.Content == "" {
			lastErr = fmt.Errorf("empty response from provider")
			logger.Error("Empty provider response",
				"requestID", opts.requestID,
				"provider", provider.Name(),
				"attempt", attempt,
				"duration", duration,
			)
			continue
		}
		
		// Log successful response
		if debugMode {
			logger.Debug("Provider response received",
				"requestID", opts.requestID,
				"provider", provider.Name(),
				"duration", duration,
				"responseLength", len(resp.Content),
				"tokensUsed", resp.Usage.TotalTokens,
			)
		}
		
		return resp.Content, nil
	}
	
	return "", fmt.Errorf("failed after %d retries: %w", retries, lastErr)
}

// isRetryableError determines if an error should trigger a retry
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	errorString := err.Error()
	
	// Retry on rate limits, timeouts, and temporary failures
	retryablePatterns := []string{
		"rate limit",
		"timeout",
		"temporary",
		"unavailable",
		"connection",
		"429", // Rate limit status code
		"503", // Service unavailable
		"504", // Gateway timeout
	}
	
	for _, pattern := range retryablePatterns {
		if strings.Contains(strings.ToLower(errorString), pattern) {
			return true
		}
	}
	
	return false
}
// CallLLM is the exported function for making LLM calls from subpackages
func CallLLM(ctx context.Context, systemPrompt, userPrompt string, opts OpOptions) (string, error) {
	return callLLM(ctx, systemPrompt, userPrompt, opts)
}
