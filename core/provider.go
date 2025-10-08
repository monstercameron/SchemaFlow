package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// Provider defines the interface for LLM providers
type Provider interface {
	// Complete sends a completion request to the LLM
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)

	// Name returns the provider name
	Name() string

	// EstimateCost estimates the cost for a request
	EstimateCost(req CompletionRequest) float64
}

// CompletionRequest represents a unified request format
type CompletionRequest struct {
	Model          string
	SystemPrompt   string
	UserPrompt     string
	Temperature    float64
	MaxTokens      int
	ResponseFormat string // "json" or "text"
}

// CompletionResponse represents a unified response format
type CompletionResponse struct {
	Content      string
	Usage        ProviderTokenUsage
	Model        string
	Provider     string
	FinishReason string
}

// ProviderTokenUsage tracks token consumption for provider responses
type ProviderTokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// ProviderConfig contains provider-specific configuration
type ProviderConfig struct {
	APIKey     string
	BaseURL    string
	OrgID      string
	Timeout    time.Duration
	MaxRetries int
	Debug      bool
}

// OpenAIProvider implements Provider for OpenAI
type OpenAIProvider struct {
	client *openai.Client
	config ProviderConfig
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config ProviderConfig) (*OpenAIProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	// Create client config
	clientConfig := openai.DefaultConfig(config.APIKey)

	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	if config.OrgID != "" {
		clientConfig.OrgID = config.OrgID
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &OpenAIProvider{
		client: client,
		config: config,
	}, nil
}

// Complete sends a completion request to OpenAI
func (provider *OpenAIProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: req.UserPrompt,
		},
	}

	chatRequest := openai.ChatCompletionRequest{
		Model:    req.Model,
		Messages: messages,
	}

	if req.Temperature > 0 {
		chatRequest.Temperature = float32(req.Temperature)
	}

	if req.MaxTokens > 0 {
		chatRequest.MaxTokens = req.MaxTokens
	}

	// Note: go-openai doesn't support response_format directly in the same way
	// We'll handle JSON formatting through prompt engineering

	completion, err := provider.client.CreateChatCompletion(ctx, chatRequest)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("OpenAI completion failed: %w", err)
	}

	if len(completion.Choices) == 0 {
		return CompletionResponse{}, fmt.Errorf("no completion choices returned")
	}

	return CompletionResponse{
		Content:      completion.Choices[0].Message.Content,
		Provider:     provider.Name(),
		Model:        completion.Model,
		FinishReason: string(completion.Choices[0].FinishReason),
		Usage: ProviderTokenUsage{
			PromptTokens:     completion.Usage.PromptTokens,
			CompletionTokens: completion.Usage.CompletionTokens,
			TotalTokens:      completion.Usage.TotalTokens,
		},
	}, nil
}

// Name returns the provider name
func (provider *OpenAIProvider) Name() string {
	return "openai"
}

// EstimateCost estimates the cost for OpenAI
func (provider *OpenAIProvider) EstimateCost(req CompletionRequest) float64 {
	// Rough estimation based on GPT-4 pricing
	// $0.03 per 1K prompt tokens, $0.06 per 1K completion tokens
	estimatedPromptTokens := len(req.SystemPrompt+req.UserPrompt) / 4
	estimatedCompletionTokens := 500 // Default estimate
	if req.MaxTokens > 0 {
		estimatedCompletionTokens = req.MaxTokens
	}

	promptCost := float64(estimatedPromptTokens) * 0.03 / 1000
	completionCost := float64(estimatedCompletionTokens) * 0.06 / 1000

	return promptCost + completionCost
}

// AnthropicProvider implements Provider for Anthropic Claude
type AnthropicProvider struct {
	config  ProviderConfig
	apiKey  string
	baseURL string
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(config ProviderConfig) (*AnthropicProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("anthropic API key is required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}

	return &AnthropicProvider{
		config:  config,
		apiKey:  config.APIKey,
		baseURL: baseURL,
	}, nil
}

// Complete sends a completion request to Anthropic
func (provider *AnthropicProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	// Note: This is a simplified implementation
	// In production, you would use the official Anthropic SDK

	model := req.Model
	if model == "" || strings.HasPrefix(model, "gpt") {
		// Map OpenAI models to Claude models
		model = "claude-3-opus-20240229"
	}

	// Anthropic uses a different message format
	prompt := fmt.Sprintf("%s\n\nHuman: %s\n\nAssistant:", req.SystemPrompt, req.UserPrompt)

	// This would make an actual API call to Anthropic
	// For now, return a mock response
	return CompletionResponse{
		Content:      fmt.Sprintf("Anthropic response to: %s", req.UserPrompt),
		Provider:     provider.Name(),
		Model:        model,
		FinishReason: "stop",
		Usage: ProviderTokenUsage{
			PromptTokens:     len(prompt) / 4,
			CompletionTokens: 100,
			TotalTokens:      len(prompt)/4 + 100,
		},
	}, nil
}

// Name returns the provider name
func (provider *AnthropicProvider) Name() string {
	return "anthropic"
}

// EstimateCost estimates the cost for Anthropic
func (provider *AnthropicProvider) EstimateCost(req CompletionRequest) float64 {
	// Claude 3 Opus pricing: $15 per million input tokens, $75 per million output tokens
	estimatedPromptTokens := len(req.SystemPrompt+req.UserPrompt) / 4
	estimatedCompletionTokens := 500
	if req.MaxTokens > 0 {
		estimatedCompletionTokens = req.MaxTokens
	}

	promptCost := float64(estimatedPromptTokens) * 15.0 / 1_000_000
	completionCost := float64(estimatedCompletionTokens) * 75.0 / 1_000_000

	return promptCost + completionCost
}

// LocalProvider implements Provider for local/mock models
type LocalProvider struct {
	config  ProviderConfig
	handler func(context.Context, CompletionRequest) (string, error)
}

// NewLocalProvider creates a new local/mock provider
func NewLocalProvider(config ProviderConfig) (*LocalProvider, error) {
	return &LocalProvider{
		config: config,
	}, nil
}

// WithHandler sets a custom handler for the local provider
func (provider *LocalProvider) WithHandler(handler func(context.Context, CompletionRequest) (string, error)) *LocalProvider {
	provider.handler = handler
	return provider
}

// Complete processes a request locally
func (provider *LocalProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	var content string
	var err error

	if provider.handler != nil {
		content, err = provider.handler(ctx, req)
		if err != nil {
			return CompletionResponse{}, err
		}
	} else {
		// Default mock behavior for testing
		content = provider.mockResponse(req)
	}

	return CompletionResponse{
		Content:      content,
		Provider:     provider.Name(),
		Model:        "local-model",
		FinishReason: "stop",
		Usage: ProviderTokenUsage{
			PromptTokens:     len(req.SystemPrompt+req.UserPrompt) / 4,
			CompletionTokens: len(content) / 4,
			TotalTokens:      (len(req.SystemPrompt+req.UserPrompt) + len(content)) / 4,
		},
	}, nil
}

// mockResponse generates a mock response for testing
func (provider *LocalProvider) mockResponse(req CompletionRequest) string {
	// Analyze the prompt to generate appropriate mock responses
	userPrompt := strings.ToLower(req.UserPrompt)

	// Check for extraction patterns
	if strings.Contains(userPrompt, "extract") || strings.Contains(req.SystemPrompt, "extraction") {
		if req.ResponseFormat == "json" {
			return `{"name": "John Doe", "age": 30, "email": "john@example.com"}`
		}
		return "Extracted data: John Doe, 30 years old"
	}

	// Check for validation patterns
	if strings.Contains(userPrompt, "validate") || strings.Contains(req.SystemPrompt, "validation") {
		return `{"valid": true, "issues": [], "confidence": 0.95}`
	}

	// Check for transformation patterns
	if strings.Contains(userPrompt, "transform") || strings.Contains(req.SystemPrompt, "transform") {
		return `{"result": "transformed data", "success": true}`
	}

	// Default response
	if req.ResponseFormat == "json" {
		return `{"response": "Mock response for testing", "success": true}`
	}

	return "Mock response for: " + req.UserPrompt
}

// Name returns the provider name
func (provider *LocalProvider) Name() string {
	return "local"
}

// EstimateCost returns 0 for local provider
func (provider *LocalProvider) EstimateCost(req CompletionRequest) float64 {
	return 0.0
}

// ProviderRegistry manages available providers
type ProviderRegistry struct {
	providers       map[string]Provider
	defaultProvider string
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry
func (registry *ProviderRegistry) Register(name string, provider Provider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}
	registry.providers[name] = provider
	if registry.defaultProvider == "" {
		registry.defaultProvider = name
	}
	return nil
}

// Get retrieves a provider by name
func (registry *ProviderRegistry) Get(name string) (Provider, error) {
	if name == "" {
		name = registry.defaultProvider
	}

	provider, ok := registry.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return provider, nil
}

// SetDefault sets the default provider
func (registry *ProviderRegistry) SetDefault(name string) error {
	if _, ok := registry.providers[name]; !ok {
		return fmt.Errorf("provider %s not found", name)
	}
	registry.defaultProvider = name
	return nil
}

// List returns all registered provider names
func (registry *ProviderRegistry) List() []string {
	names := make([]string, 0, len(registry.providers))
	for name := range registry.providers {
		names = append(names, name)
	}
	return names
}

// Global provider registry
var globalRegistry = NewProviderRegistry()

// RegisterProvider registers a provider globally
func RegisterProvider(name string, provider Provider) error {
	return globalRegistry.Register(name, provider)
}

// GetProviderFromRegistry gets a provider from the global registry
func GetProviderFromRegistry(name string) (Provider, error) {
	return globalRegistry.Get(name)
}

// SetDefaultProvider sets the global default provider
func SetDefaultProvider(name string) error {
	return globalRegistry.SetDefault(name)
}
