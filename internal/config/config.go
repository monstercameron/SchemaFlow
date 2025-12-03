package config

import (
	"os"
	"sync"
	"time"

	"github.com/monstercameron/SchemaFlow/internal/types"
)

var (
	debugMode      bool
	traceEnabled   bool
	metricsEnabled bool
	apiKey         string
	mu             sync.RWMutex
)

// Init initializes the configuration with an API key
func Init(key string) {
	mu.Lock()
	defer mu.Unlock()

	apiKey = key
	if apiKey == "" {
		apiKey = os.Getenv("SCHEMAFLOW_API_KEY")
	}
}

// GetAPIKey returns the configured API key
func GetAPIKey() string {
	mu.RLock()
	defer mu.RUnlock()
	if apiKey != "" {
		return apiKey
	}
	return os.Getenv("SCHEMAFLOW_API_KEY")
}

// GetTimeout returns the default timeout for operations
func GetTimeout() time.Duration {
	if timeoutStr := os.Getenv("SCHEMAFLOW_TIMEOUT"); timeoutStr != "" {
		if duration, err := time.ParseDuration(timeoutStr); err == nil {
			return duration
		}
	}
	return 30 * time.Second
}

// GetDebugMode returns whether debug mode is enabled
func GetDebugMode() bool {
	mu.RLock()
	defer mu.RUnlock()
	if debugMode {
		return true
	}
	return os.Getenv("SCHEMAFLOW_DEBUG") == "true"
}

// SetDebugMode sets the debug mode
func SetDebugMode(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	debugMode = enabled
}

// GetTraceEnabled returns whether tracing is enabled
func GetTraceEnabled() bool {
	mu.RLock()
	defer mu.RUnlock()
	if traceEnabled {
		return true
	}
	return os.Getenv("SCHEMAFLOW_TRACE") == "true"
}

// SetTraceEnabled sets tracing mode
func SetTraceEnabled(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	traceEnabled = enabled
}

// IsMetricsEnabled returns whether metrics collection is enabled
func IsMetricsEnabled() bool {
	mu.RLock()
	defer mu.RUnlock()
	if metricsEnabled {
		return true
	}
	return os.Getenv("SCHEMAFLOW_METRICS") == "true"
}

// SetMetricsEnabled sets whether metrics collection is enabled
func SetMetricsEnabled(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	metricsEnabled = enabled
}

// GetModel returns the appropriate model based on intelligence level
func GetModel(intelligence types.Speed, provider string) string {
	// Check for global model override via environment variable
	if envModel := os.Getenv("SCHEMAFLOW_MODEL"); envModel != "" {
		return envModel
	}

	// Check for intelligence-level specific overrides
	var envLevelModel string
	switch intelligence {
	case types.Smart:
		envLevelModel = os.Getenv("SCHEMAFLOW_MODEL_SMART")
	case types.Fast:
		envLevelModel = os.Getenv("SCHEMAFLOW_MODEL_FAST")
	case types.Quick:
		envLevelModel = os.Getenv("SCHEMAFLOW_MODEL_QUICK")
	}

	if envLevelModel != "" {
		return envLevelModel
	}

	// Check global provider
	if provider == "openrouter" {
		switch intelligence {
		case types.Smart:
			return "openai/gpt-4o"
		case types.Fast:
			return "openai/gpt-4o-mini"
		case types.Quick:
			return "openai/gpt-4o-mini"
		default:
			return "openai/gpt-4o-mini"
		}
	}

	if provider == "cerebras" {
		switch intelligence {
		case types.Smart:
			return "llama-3.3-70b"
		case types.Fast:
			return "llama3.1-8b"
		case types.Quick:
			return "llama3.1-8b"
		default:
			return "llama3.1-8b"
		}
	}

	if provider == "anthropic" {
		switch intelligence {
		case types.Smart:
			return "claude-3-5-sonnet-20240620"
		case types.Fast:
			return "claude-3-haiku-20240307"
		case types.Quick:
			return "claude-3-haiku-20240307"
		default:
			return "claude-3-haiku-20240307"
		}
	}

	switch intelligence {
	case types.Smart:
		return "gpt-5-2025-08-07"
	case types.Fast:
		return "gpt-5-nano-2025-08-07"
	case types.Quick:
		return "gpt-5-mini-2025-08-07"
	default:
		return "gpt-5-nano-2025-08-07"
	}
}

// GetMaxTokens returns the maximum token limit based on intelligence level
func GetMaxTokens(intelligence types.Speed) int {
	switch intelligence {
	case types.Smart:
		return 4000
	case types.Fast:
		return 2000
	case types.Quick:
		return 1000
	default:
		return 2000
	}
}

// GetTemperature returns the appropriate temperature setting based on mode
func GetTemperature(mode types.Mode) float32 {
	switch mode {
	case types.Strict:
		return 0.1 // Very deterministic
	case types.TransformMode:
		return 0.3 // Balanced
	case types.Creative:
		return 0.7 // More creative/varied
	default:
		return 0.3
	}
}
