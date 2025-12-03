package schemaflow

import (
	"os"
	"sync"
	"time"

	"github.com/monstercameron/SchemaFlow/internal/llm"
	"github.com/monstercameron/SchemaFlow/internal/telemetry"
	openai "github.com/sashabaranov/go-openai"
)

// Client represents a configured schemaflow client instance.
type Client struct {
	openaiClient *openai.Client // Legacy OpenAI client
	apiKey       string
	provider     llm.Provider
	providerName string
	timeout      time.Duration
	maxRetries   int
	retryBackoff time.Duration
	logger       *telemetry.Logger
	debugMode    bool
	mu           sync.RWMutex
}

// NewClient creates a new client with custom configuration
func NewClient(apiKey string) *Client {
	client := &Client{
		apiKey:       apiKey,
		providerName: "openai",
		timeout:      30 * time.Second,
		maxRetries:   3,
		retryBackoff: 1 * time.Second,
		logger:       telemetry.NewLogger(),
	}

	// Initialize with OpenAI provider by default
	if apiKey != "" {
		client.openaiClient = openai.NewClient(apiKey)
		// Also create the new provider-based client
		provider, err := llm.NewOpenAIProvider(llm.ProviderConfig{
			APIKey:     apiKey,
			Timeout:    client.timeout,
			MaxRetries: client.maxRetries,
			Debug:      client.debugMode,
		})
		if err == nil {
			client.provider = provider
		} else {
			client.logger.Warn("Failed to create OpenAI provider", "error", err)
		}
	} else {
		// Use local provider for testing
		localProvider, _ := llm.NewLocalProvider(llm.ProviderConfig{})
		client.provider = localProvider
	}

	return client
}

// WithTimeout sets a custom timeout for the client
func (client *Client) WithTimeout(timeout time.Duration) *Client {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.timeout = timeout
	return client
}

// WithProvider sets a custom provider for the client
func (client *Client) WithProvider(providerName string) *Client {
	client.mu.Lock()
	defer client.mu.Unlock()

	client.providerName = providerName

	// Create the appropriate provider based on name
	var provider llm.Provider
	var err error

	config := llm.ProviderConfig{
		APIKey:     client.apiKey,
		Timeout:    client.timeout,
		MaxRetries: client.maxRetries,
		Debug:      client.debugMode,
	}

	switch providerName {
	case "openai":
		provider, err = llm.NewOpenAIProvider(config)
	case "anthropic":
		provider, err = llm.NewAnthropicProvider(config)
	case "openrouter":
		provider, err = llm.NewOpenRouterProvider(config)
	case "cerebras":
		provider, err = llm.NewCerebrasProvider(config)
	case "local", "mock":
		provider, err = llm.NewLocalProvider(config)
	default:
		// Try to get from global registry
		provider, err = llm.GetProviderFromRegistry(providerName)
	}

	if err != nil {
		client.logger.Warn("Failed to create provider, using default", "provider", providerName, "error", err)
	} else {
		client.provider = provider
		client.logger.Info("Provider configured", "provider", providerName)
	}

	return client
}

// WithDebug enables debug mode for the client
func (client *Client) WithDebug(enabled bool) *Client {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.debugMode = enabled
	if enabled {
		client.logger.SetLevel(telemetry.DebugLevel)
	}
	return client
}

// Global configuration
var (
	defaultClient *Client
	mu            sync.RWMutex
)

// Init initializes the schemaflow library with the provided API key.
func Init(key string) {
	mu.Lock()
	defer mu.Unlock()

	apiKey := key
	if apiKey == "" {
		apiKey = os.Getenv("SCHEMAFLOW_API_KEY")
	}

	provider := "openai"
	if p := os.Getenv("SCHEMAFLOW_PROVIDER"); p != "" {
		provider = p
	}

	timeout := 30 * time.Second
	if t := os.Getenv("SCHEMAFLOW_TIMEOUT"); t != "" {
		if d, err := time.ParseDuration(t); err == nil {
			timeout = d
		}
	}

	debugMode := false
	if d := os.Getenv("SCHEMAFLOW_DEBUG"); d == "true" || d == "1" {
		debugMode = true
	}

	if apiKey != "" {
		defaultClient = NewClient(apiKey).
			WithTimeout(timeout).
			WithProvider(provider).
			WithDebug(debugMode)
	} else {
		defaultClient = NewClient("")
	}
}

// GetDefaultClient returns the default client
func GetDefaultClient() *Client {
	return defaultClient
}

// InitWithEnv initializes SchemaFlow from environment variables.
// It reads configuration from a .env file if path is provided.
func InitWithEnv(paths ...string) error {
	// Load .env file if path provided (optional)
	// For now, just use environment variables directly
	apiKey := os.Getenv("SCHEMAFLOW_API_KEY")
	if apiKey == "" {
		// Not an error - might be using local provider
	}

	Init(apiKey)
	return nil
}

// GetLogger returns the default logger for the schemaflow package.
func GetLogger() *telemetry.Logger {
	if defaultClient != nil {
		return defaultClient.logger
	}
	return telemetry.NewLogger()
}
