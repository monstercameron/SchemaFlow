package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/monstercameron/SchemaFlow/internal/types"
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
	Usage        types.TokenUsage
	Model        string
	Provider     string
	FinishReason string
}

// ProviderConfig contains provider-specific configuration
type ProviderConfig struct {
	APIKey       string
	BaseURL      string
	OrgID        string
	Timeout      time.Duration
	MaxRetries   int
	Debug        bool
	ExtraHeaders map[string]string
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

// Complete sends a completion request to OpenAI using the Responses API
func (provider *OpenAIProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	// Use the new Responses API (POST /v1/responses)
	// Since go-openai v1.20.4 doesn't support this endpoint, we implement it manually.

	url := "https://api.openai.com/v1/responses"
	if provider.config.BaseURL != "" {
		url = strings.TrimRight(provider.config.BaseURL, "/") + "/responses"
	}

	// Construct request body
	requestBody := map[string]interface{}{
		"model":        req.Model,
		"input":        req.UserPrompt,
		"instructions": req.SystemPrompt,
	}

	if req.Temperature > 0 {
		requestBody["temperature"] = req.Temperature
	}

	if req.MaxTokens > 0 {
		requestBody["max_output_tokens"] = req.MaxTokens
	}

	// Handle response format
	if req.ResponseFormat == "json" {
		// The Responses API uses "text" object configuration
		requestBody["text"] = map[string]interface{}{
			"format": map[string]string{
				"type": "json_object", // Assuming json_object is supported or fallback to text with prompt
			},
		}
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+provider.config.APIKey)

	if provider.config.OrgID != "" {
		httpReq.Header.Set("OpenAI-Organization", provider.config.OrgID)
	}

	// Use a custom HTTP client or default
	client := &http.Client{
		Timeout: provider.config.Timeout,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("OpenAI request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return CompletionResponse{}, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var response struct {
		ID     string `json:"id"`
		Output []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
			TotalTokens  int `json:"total_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return CompletionResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Output) == 0 || len(response.Output[0].Content) == 0 {
		return CompletionResponse{}, fmt.Errorf("empty response from OpenAI")
	}

	// Extract text content
	content := ""
	for _, item := range response.Output[0].Content {
		if item.Type == "output_text" {
			content += item.Text
		}
	}

	return CompletionResponse{
		Content:      content,
		Provider:     provider.Name(),
		Model:        response.Model,
		FinishReason: "stop", // Responses API doesn't explicitly return finish_reason in the same way, assuming stop
		Usage: types.TokenUsage{
			PromptTokens:     response.Usage.InputTokens,
			CompletionTokens: response.Usage.OutputTokens,
			TotalTokens:      response.Usage.TotalTokens,
		},
	}, nil
}

// Name returns the provider name
func (provider *OpenAIProvider) Name() string {
	return "openai"
}

// EstimateCost estimates the cost for OpenAI
func (provider *OpenAIProvider) EstimateCost(req CompletionRequest) float64 {
	inputRate, outputRate := getModelRates("openai", req.Model)

	estimatedPromptTokens := len(req.SystemPrompt+req.UserPrompt) / 4
	estimatedCompletionTokens := 500 // Default estimate
	if req.MaxTokens > 0 {
		estimatedCompletionTokens = req.MaxTokens
	}

	promptCost := float64(estimatedPromptTokens) * inputRate
	completionCost := float64(estimatedCompletionTokens) * outputRate

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
	url := strings.TrimRight(provider.baseURL, "/") + "/v1/messages"

	model := req.Model
	if model == "" || strings.HasPrefix(model, "gpt") {
		// Default to Sonnet 3.5 if no valid model specified
		model = "claude-3-5-sonnet-20240620"
	}

	// Construct messages
	messages := []map[string]string{
		{
			"role":    "user",
			"content": req.UserPrompt,
		},
	}

	// Construct request body
	requestBody := map[string]interface{}{
		"model":      model,
		"messages":   messages,
		"max_tokens": 1024, // Default max tokens
	}

	if req.SystemPrompt != "" {
		requestBody["system"] = req.SystemPrompt
	}

	if req.Temperature > 0 {
		requestBody["temperature"] = req.Temperature
	}

	if req.MaxTokens > 0 {
		requestBody["max_tokens"] = req.MaxTokens
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", provider.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// Use a custom HTTP client or default
	client := &http.Client{
		Timeout: provider.config.Timeout,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("Anthropic request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return CompletionResponse{}, fmt.Errorf("Anthropic API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var response struct {
		ID      string `json:"id"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return CompletionResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Content) == 0 {
		return CompletionResponse{}, fmt.Errorf("empty response from Anthropic")
	}

	// Extract text content
	content := ""
	for _, item := range response.Content {
		if item.Type == "text" {
			content += item.Text
		}
	}

	return CompletionResponse{
		Content:      content,
		Provider:     provider.Name(),
		Model:        response.Model,
		FinishReason: "stop",
		Usage: types.TokenUsage{
			PromptTokens:     response.Usage.InputTokens,
			CompletionTokens: response.Usage.OutputTokens,
			TotalTokens:      response.Usage.InputTokens + response.Usage.OutputTokens,
		},
	}, nil
}

// Name returns the provider name
func (provider *AnthropicProvider) Name() string {
	return "anthropic"
}

// EstimateCost estimates the cost for Anthropic
func (provider *AnthropicProvider) EstimateCost(req CompletionRequest) float64 {
	inputRate, outputRate := getModelRates("anthropic", req.Model)

	estimatedPromptTokens := len(req.SystemPrompt+req.UserPrompt) / 4
	estimatedCompletionTokens := 500
	if req.MaxTokens > 0 {
		estimatedCompletionTokens = req.MaxTokens
	}

	promptCost := float64(estimatedPromptTokens) * inputRate
	completionCost := float64(estimatedCompletionTokens) * outputRate

	return promptCost + completionCost
}

// OpenRouterProvider implements Provider for OpenRouter
type OpenRouterProvider struct {
	client *openai.Client
	config ProviderConfig
}

// customTransport is a http.RoundTripper that adds custom headers
type customTransport struct {
	transport http.RoundTripper
	headers   map[string]string
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	return t.transport.RoundTrip(req)
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(config ProviderConfig) (*OpenRouterProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenRouter API key is required")
	}

	// Create client config
	clientConfig := openai.DefaultConfig(config.APIKey)

	// Set OpenRouter base URL
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	} else {
		clientConfig.BaseURL = "https://openrouter.ai/api/v1"
	}

	if config.OrgID != "" {
		clientConfig.OrgID = config.OrgID
	}

	// Add custom headers if provided
	if len(config.ExtraHeaders) > 0 {
		clientConfig.HTTPClient = &http.Client{
			Transport: &customTransport{
				transport: http.DefaultTransport,
				headers:   config.ExtraHeaders,
			},
		}
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &OpenRouterProvider{
		client: client,
		config: config,
	}, nil
}

// Complete sends a completion request to OpenRouter
func (provider *OpenRouterProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
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

	// OpenRouter supports response_format for some models
	if req.ResponseFormat == "json" {
		chatRequest.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		}
	}

	completion, err := provider.client.CreateChatCompletion(ctx, chatRequest)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("OpenRouter completion failed: %w", err)
	}

	if len(completion.Choices) == 0 {
		return CompletionResponse{}, fmt.Errorf("no completion choices returned")
	}

	return CompletionResponse{
		Content:      completion.Choices[0].Message.Content,
		Provider:     provider.Name(),
		Model:        completion.Model,
		FinishReason: string(completion.Choices[0].FinishReason),
		Usage: types.TokenUsage{
			PromptTokens:     completion.Usage.PromptTokens,
			CompletionTokens: completion.Usage.CompletionTokens,
			TotalTokens:      completion.Usage.TotalTokens,
		},
	}, nil
}

// Name returns the provider name
func (provider *OpenRouterProvider) Name() string {
	return "openrouter"
}

// EstimateCost estimates the cost for OpenRouter
func (provider *OpenRouterProvider) EstimateCost(req CompletionRequest) float64 {
	inputRate, outputRate := getModelRates("openrouter", req.Model)

	estimatedPromptTokens := len(req.SystemPrompt+req.UserPrompt) / 4
	estimatedCompletionTokens := 500
	if req.MaxTokens > 0 {
		estimatedCompletionTokens = req.MaxTokens
	}

	promptCost := float64(estimatedPromptTokens) * inputRate
	completionCost := float64(estimatedCompletionTokens) * outputRate

	return promptCost + completionCost
}

// CerebrasProvider implements Provider for Cerebras
type CerebrasProvider struct {
	client *openai.Client
	config ProviderConfig
}

// NewCerebrasProvider creates a new Cerebras provider
func NewCerebrasProvider(config ProviderConfig) (*CerebrasProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Cerebras API key is required")
	}

	// Create client config
	clientConfig := openai.DefaultConfig(config.APIKey)

	// Set Cerebras base URL
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	} else {
		clientConfig.BaseURL = "https://api.cerebras.ai/v1"
	}

	if config.OrgID != "" {
		clientConfig.OrgID = config.OrgID
	}

	// Add custom headers if provided
	if len(config.ExtraHeaders) > 0 {
		clientConfig.HTTPClient = &http.Client{
			Transport: &customTransport{
				transport: http.DefaultTransport,
				headers:   config.ExtraHeaders,
			},
		}
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &CerebrasProvider{
		client: client,
		config: config,
	}, nil
}

// Complete sends a completion request to Cerebras
func (provider *CerebrasProvider) Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
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

	// Cerebras supports response_format for some models (json_object)
	if req.ResponseFormat == "json" {
		chatRequest.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		}
	}

	completion, err := provider.client.CreateChatCompletion(ctx, chatRequest)
	if err != nil {
		return CompletionResponse{}, fmt.Errorf("Cerebras completion failed: %w", err)
	}

	if len(completion.Choices) == 0 {
		return CompletionResponse{}, fmt.Errorf("no completion choices returned")
	}

	return CompletionResponse{
		Content:      completion.Choices[0].Message.Content,
		Provider:     provider.Name(),
		Model:        completion.Model,
		FinishReason: string(completion.Choices[0].FinishReason),
		Usage: types.TokenUsage{
			PromptTokens:     completion.Usage.PromptTokens,
			CompletionTokens: completion.Usage.CompletionTokens,
			TotalTokens:      completion.Usage.TotalTokens,
		},
	}, nil
}

// Name returns the provider name
func (provider *CerebrasProvider) Name() string {
	return "cerebras"
}

// EstimateCost estimates the cost for Cerebras
func (provider *CerebrasProvider) EstimateCost(req CompletionRequest) float64 {
	inputRate, outputRate := getModelRates("cerebras", req.Model)

	estimatedPromptTokens := len(req.SystemPrompt+req.UserPrompt) / 4
	estimatedCompletionTokens := 500
	if req.MaxTokens > 0 {
		estimatedCompletionTokens = req.MaxTokens
	}

	promptCost := float64(estimatedPromptTokens) * inputRate
	completionCost := float64(estimatedCompletionTokens) * outputRate

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
		Usage: types.TokenUsage{
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

// getModelRates returns the input and output cost per token for a given model
// It checks environment variables first, then falls back to provider defaults.
// Environment variables format: SCHEMAFLOW_COST_INPUT_<MODEL> and SCHEMAFLOW_COST_OUTPUT_<MODEL>
// where <MODEL> is the sanitized model name (uppercase, dots/dashes replaced with underscores).
// The value should be the cost per 1,000,000 (1M) tokens.
func getModelRates(provider, model string) (float64, float64) {
	// Check for intelligence level overrides first if the model matches a known level env var
	// Note: This is tricky because we don't know the intelligence level here, only the model name.
	// However, if the user set SCHEMAFLOW_MODEL_SMART="gpt-4", they might also set SCHEMAFLOW_COST_INPUT_SMART.
	// But the request was "map the models and cost/token at the levels in the env vars".
	// This implies we should check SCHEMAFLOW_COST_INPUT_SMART if the current model IS the smart model.
	// Since we don't have the level context here easily without circular dependency or API change,
	// we will stick to model-based lookup which is robust.
	// If the user sets SCHEMAFLOW_MODEL_SMART="my-model", they can set SCHEMAFLOW_COST_INPUT_MY_MODEL.

	// Wait, the user specifically asked: "map the models and cost/token at the levels in the env vars"
	// This could mean: SCHEMAFLOW_COST_INPUT_SMART=30.0
	// If we can't know if "model" is "Smart", we can't apply "Smart" pricing.
	// But wait, the caller of EstimateCost knows the model, but not necessarily the level used to select it.
	// Let's support looking up by the exact model name first (as implemented),
	// AND also support looking up by "SMART", "FAST", "QUICK" if the model name matches one of those keywords (unlikely)
	// OR, we can just rely on the user setting the cost for the *model* they mapped to the level.

	// However, if the user wants to say "Smart level costs X", and they mapped Smart -> GPT-4,
	// they probably want to set cost for GPT-4.
	// BUT, if they want to abstract it:
	// SCHEMAFLOW_MODEL_SMART=gpt-4
	// SCHEMAFLOW_COST_INPUT_SMART=30.0

	// To support this, we'd need to know if 'model' == os.Getenv("SCHEMAFLOW_MODEL_SMART").

	smartModel := os.Getenv("SCHEMAFLOW_MODEL_SMART")
	fastModel := os.Getenv("SCHEMAFLOW_MODEL_FAST")
	quickModel := os.Getenv("SCHEMAFLOW_MODEL_QUICK")

	var levelSuffix string
	if model == smartModel && smartModel != "" {
		levelSuffix = "SMART"
	} else if model == fastModel && fastModel != "" {
		levelSuffix = "FAST"
	} else if model == quickModel && quickModel != "" {
		levelSuffix = "QUICK"
	}

	if levelSuffix != "" {
		inputEnv := os.Getenv(fmt.Sprintf("SCHEMAFLOW_COST_INPUT_%s", levelSuffix))
		outputEnv := os.Getenv(fmt.Sprintf("SCHEMAFLOW_COST_OUTPUT_%s", levelSuffix))

		if inputEnv != "" && outputEnv != "" {
			inputVal, err1 := strconv.ParseFloat(inputEnv, 64)
			outputVal, err2 := strconv.ParseFloat(outputEnv, 64)
			if err1 == nil && err2 == nil {
				return inputVal / 1_000_000, outputVal / 1_000_000
			}
		}
	}

	sanitizedModel := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(model, "-", "_"), ".", "_"))

	inputEnv := os.Getenv(fmt.Sprintf("SCHEMAFLOW_COST_INPUT_%s", sanitizedModel))
	outputEnv := os.Getenv(fmt.Sprintf("SCHEMAFLOW_COST_OUTPUT_%s", sanitizedModel))

	var inputRate, outputRate float64
	var inputFound, outputFound bool

	if inputEnv != "" {
		if val, err := strconv.ParseFloat(inputEnv, 64); err == nil {
			inputRate = val / 1_000_000
			inputFound = true
		}
	}

	if outputEnv != "" {
		if val, err := strconv.ParseFloat(outputEnv, 64); err == nil {
			outputRate = val / 1_000_000
			outputFound = true
		}
	}

	if inputFound && outputFound {
		return inputRate, outputRate
	}

	// Defaults if not fully specified in env vars
	defaultInput, defaultOutput := getDefaultRates(provider, model)

	if !inputFound {
		inputRate = defaultInput
	}
	if !outputFound {
		outputRate = defaultOutput
	}

	return inputRate, outputRate
}

func getDefaultRates(provider, _ string) (float64, float64) {
	switch provider {
	case "openai":
		// Default to GPT-4 pricing ($30/1M input, $60/1M output)
		return 30.0 / 1_000_000, 60.0 / 1_000_000
	case "anthropic":
		// Claude 3 Opus: $15/1M, $75/1M
		return 15.0 / 1_000_000, 75.0 / 1_000_000
	case "openrouter":
		// Generic: $1/1M, $2/1M
		return 1.0 / 1_000_000, 2.0 / 1_000_000
	case "cerebras":
		// Cheap: $0.10/1M
		return 0.10 / 1_000_000, 0.10 / 1_000_000
	default:
		return 0, 0
	}
}
