// Package schemaflow - Configuration and initialization
package schemaflow

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// Client represents a configured schemaflow client instance.
// This allows multiple clients with different configurations.
type Client struct {
	openaiClient *openai.Client  // Legacy OpenAI client (for backward compatibility)
	apiKey       string
	provider     Provider         // Provider interface for multiple LLM backends
	providerName string          // Name of the provider (openai, anthropic, local)
	timeout      time.Duration
	maxRetries   int
	retryBackoff time.Duration
	logger       *Logger
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
		logger:       NewLogger(),
	}
	
	// Initialize with OpenAI provider by default
	if apiKey != "" {
		client.openaiClient = openai.NewClient(apiKey)
		// Also create the new provider-based client
		provider, err := NewOpenAIProvider(ProviderConfig{
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
		localProvider, _ := NewLocalProvider(ProviderConfig{})
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
	var provider Provider
	var err error
	
	config := ProviderConfig{
		APIKey:     client.apiKey,
		Timeout:    client.timeout,
		MaxRetries: client.maxRetries,
		Debug:      client.debugMode,
	}
	
	switch providerName {
	case "openai":
		provider, err = NewOpenAIProvider(config)
	case "anthropic":
		provider, err = NewAnthropicProvider(config)
	case "local", "mock":
		provider, err = NewLocalProvider(config)
	default:
		// Try to get from global registry
		provider, err = GetProviderFromRegistry(providerName)
	}
	
	if err != nil {
		client.logger.Warn("Failed to create provider, using default", "provider", providerName, "error", err)
	} else {
		client.provider = provider
		client.logger.Info("Provider configured", "provider", providerName)
	}
	
	return client
}

// WithProviderInstance sets a custom provider instance directly
func (client *Client) WithProviderInstance(provider Provider) *Client {
	client.mu.Lock()
	defer client.mu.Unlock()
	
	if provider != nil {
		client.provider = provider
		client.providerName = provider.Name()
	}
	
	return client
}

// WithDebug enables debug mode for the client
func (client *Client) WithDebug(enabled bool) *Client {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.debugMode = enabled
	if enabled {
		client.logger.SetLevel(DebugLevel)
	}
	return client
}

// Global configuration and state management
var (
	client       *openai.Client      // OpenAI API client instance
	apiKey       string               // API key for authentication
	provider     string = "openai"   // LLM provider (openai, anthropic, etc.)
	timeout      time.Duration = 30 * time.Second // Default operation timeout
	maxRetries   int = 3             // Maximum retry attempts for failed operations
	retryBackoff time.Duration = 1 * time.Second  // Initial backoff duration for retries
	mu           sync.RWMutex        // Mutex for thread-safe configuration updates
	
	// Logging configuration
	logger       *Logger              // Structured logger instance
	debugMode    bool                 // Enable debug logging
	traceEnabled bool                 // Enable request tracing
	metricsEnabled bool               // Enable metrics collection
	
	// Default client for backward compatibility
	defaultClient *Client             // Default client instance
)

// Init initializes the schemaflow library with the provided API key.
// This must be called before using any operations.
// Configuration can be customized via environment variables:
//   - SCHEMAFLOW_API_KEY: API key (can be passed as parameter)
//   - SCHEMAFLOW_PROVIDER: LLM provider (default: "openai")
//   - SCHEMAFLOW_TIMEOUT: Operation timeout (default: "30s")
//   - SCHEMAFLOW_MAX_RETRIES: Retry attempts (default: 3)
//   - SCHEMAFLOW_RETRY_BACKOFF: Initial backoff duration (default: "1s")
//   - SCHEMAFLOW_DEBUG: Enable debug logging (default: false)
//   - SCHEMAFLOW_TRACE: Enable request tracing (default: false)
func Init(key string) {
	mu.Lock()
	defer mu.Unlock()
	
	// Initialize logger
	logger = NewLogger()
	
	// API key configuration
	apiKey = key
	if apiKey == "" {
		apiKey = os.Getenv("SCHEMAFLOW_API_KEY")
	}
	
	// Provider configuration
	if providerStr := os.Getenv("SCHEMAFLOW_PROVIDER"); providerStr != "" {
		provider = providerStr
		logger.Info("Provider configured", "provider", provider)
	}
	
	// Timeout configuration
	if timeoutStr := os.Getenv("SCHEMAFLOW_TIMEOUT"); timeoutStr != "" {
		if duration, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = duration
			logger.Info("Timeout configured", "timeout", timeout)
		} else {
			logger.Warn("Invalid timeout value, using default", "value", timeoutStr, "default", timeout)
		}
	}
	
	// Retry configuration
	if maxRetriesStr := os.Getenv("SCHEMAFLOW_MAX_RETRIES"); maxRetriesStr != "" {
		if _, err := fmt.Sscanf(maxRetriesStr, "%d", &maxRetries); err != nil {
			logger.Warn("Invalid max retries value, using default", "value", maxRetriesStr, "default", maxRetries)
		} else {
			logger.Info("Max retries configured", "maxRetries", maxRetries)
		}
	}
	
	// Retry backoff configuration
	if retryBackoffStr := os.Getenv("SCHEMAFLOW_RETRY_BACKOFF"); retryBackoffStr != "" {
		if duration, err := time.ParseDuration(retryBackoffStr); err == nil {
			retryBackoff = duration
			logger.Info("Retry backoff configured", "backoff", retryBackoff)
		} else {
			logger.Warn("Invalid retry backoff value, using default", "value", retryBackoffStr, "default", retryBackoff)
		}
	}
	
	// Debug mode configuration
	if debugStr := os.Getenv("SCHEMAFLOW_DEBUG"); debugStr == "true" || debugStr == "1" {
		debugMode = true
		logger.SetLevel(DebugLevel)
		logger.Debug("Debug mode enabled")
	}
	
	// Trace configuration
	if traceStr := os.Getenv("SCHEMAFLOW_TRACE"); traceStr == "true" || traceStr == "1" {
		traceEnabled = true
		logger.Info("Request tracing enabled")
	}
	
	// Metrics configuration
	if metricsStr := os.Getenv("SCHEMAFLOW_METRICS"); metricsStr == "true" || metricsStr == "1" {
		metricsEnabled = true
		logger.Info("Metrics collection enabled")
	}
	
	// Initialize OpenAI client if API key is present
	if apiKey != "" {
		client = openai.NewClient(apiKey)
		// Also create default client for new client-based API
		defaultClient = NewClient(apiKey).
			WithTimeout(timeout).
			WithProvider(provider).
			WithDebug(debugMode)
		logger.Info("schemaflow initialized successfully", "provider", provider)
	} else {
		logger.Warn("No API key provided, running in test/mock mode", "source", "Init")
		// Create a mock client for testing
		defaultClient = NewClient("")
	}
}

// applyDefaults applies default values to OpOptions if not specified
func applyDefaults(opts []OpOptions) OpOptions {
	if len(opts) == 0 {
		return OpOptions{
			Threshold:    0.7,
			Mode:         TransformMode,
			Intelligence: Smart,  // Changed from Fast to Smart
			context:      context.Background(),
			requestID:    generateRequestID(),
		}
	}
	
	opt := opts[0]
	
	// Apply defaults for zero values
	if opt.Threshold == 0 {
		opt.Threshold = 0.7
	}
	// Note: Mode and Intelligence can be 0 (Strict/Smart) which are valid values
	// Only apply defaults if struct was completely empty
	// We can't distinguish between explicit 0 and unset in Go
	// So we only apply Mode/Intelligence defaults for the empty options case above
	if opt.context == nil {
		opt.context = context.Background()
	}
	if opt.requestID == "" {
		opt.requestID = generateRequestID()
	}
	
	return opt
}

// getModel returns the appropriate OpenAI model based on intelligence level
func getModel(intelligence Speed) string {
	switch intelligence {
	case Smart:
		return "gpt-4-turbo-preview"
	case Fast:
		return "gpt-3.5-turbo"
	case Quick:
		return "gpt-3.5-turbo"
	default:
		return "gpt-3.5-turbo"
	}
}

// getMaxTokens returns the maximum token limit based on intelligence level
func getMaxTokens(intelligence Speed) int {
	switch intelligence {
	case Smart:
		return 4000
	case Fast:
		return 2000
	case Quick:
		return 1000
	default:
		return 2000
	}
}

// getTemperature returns the appropriate temperature setting based on mode
func getTemperature(mode Mode) float32 {
	switch mode {
	case Strict:
		return 0.1 // Very deterministic
	case TransformMode:
		return 0.3 // Balanced
	case Creative:
		return 0.7 // More creative/varied
	default:
		return 0.3
	}
}

// Helper functions for package access

// GetDebugMode returns the current debug mode status
func GetDebugMode() bool {
	mu.RLock()
	defer mu.RUnlock()
	return debugMode
}

// GetTraceEnabled returns whether tracing is enabled
func GetTraceEnabled() bool {
	mu.RLock()
	defer mu.RUnlock()
	return traceEnabled
}

// GenerateRequestID exports the generateRequestID function
func GenerateRequestID() string {
	return generateRequestID()
}

// GetDefaultClient returns the default client
func GetDefaultClient() *Client {
	return defaultClient
}

// GetClient returns the legacy client
func GetClient() interface{} {
	return client
}

// GetTimeout returns the global timeout
func GetTimeout() time.Duration {
	return timeout
}
