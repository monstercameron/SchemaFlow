package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnv loads environment variables from a .env file
func LoadEnv(filePath ...string) error {
	envFile := ".env"
	if len(filePath) > 0 {
		envFile = filePath[0]
	}

	// Try to find .env file in current directory first
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		// Try parent directories
		dir, _ := os.Getwd()
		for dir != "/" && dir != "" {
			testPath := filepath.Join(dir, ".env")
			if _, err := os.Stat(testPath); err == nil {
				envFile = testPath
				break
			}
			dir = filepath.Dir(dir)
		}
	}

	file, err := os.Open(envFile)
	if err != nil {
		// .env file is optional
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		// Only set if not already set (env vars take precedence)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// InitWithEnv initializes SchemaFlow with environment variables, optionally loading from .env file
func InitWithEnv(envFile ...string) error {
	// Load .env file if it exists
	if len(envFile) > 0 {
		if err := LoadEnv(envFile[0]); err != nil {
			return fmt.Errorf("failed to load env file: %w", err)
		}
	} else {
		LoadEnv() // Try to load default .env
	}

	// Get API key from environment
	apiKey := os.Getenv("SCHEMAFLOW_API_KEY")
	if apiKey == "" {
		// Fallback to OPENAI_API_KEY
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	if apiKey == "" {
		return fmt.Errorf("no API key found in environment (SCHEMAFLOW_API_KEY or OPENAI_API_KEY)")
	}

	// Initialize with the API key
	Init(apiKey)

	// Override model names if specified
	// Check both SCHEMAFLOW_ and OPENAI_ prefixes
	if model := os.Getenv("SCHEMAFLOW_MODEL_SMART"); model != "" {
		SetModel(Smart, model)
	} else if model := os.Getenv("OPENAI_SMART_MODEL"); model != "" {
		SetModel(Smart, model)
	}

	if model := os.Getenv("SCHEMAFLOW_MODEL_FAST"); model != "" {
		SetModel(Fast, model)
	} else if model := os.Getenv("OPENAI_FAST_MODEL"); model != "" {
		SetModel(Fast, model)
	}

	if model := os.Getenv("SCHEMAFLOW_MODEL_QUICK"); model != "" {
		SetModel(Quick, model)
	} else if model := os.Getenv("OPENAI_QUICK_MODEL"); model != "" {
		SetModel(Quick, model)
	}

	// Log configuration if debug is enabled
	if debugMode {
		logger.Info("SchemaFlow initialized with environment configuration",
			"provider", provider,
			"timeout", timeout,
			"maxRetries", maxRetries,
			"debugMode", debugMode,
		)
	}

	return nil
}

// Model configuration
var modelOverrides = make(map[Speed]string)

// SetModel allows overriding the model for a specific intelligence level
func SetModel(intelligence Speed, model string) {
	mu.Lock()
	defer mu.Unlock()
	modelOverrides[intelligence] = model

	if debugMode {
		logger.Debug("Model override set",
			"intelligence", intelligence.String(),
			"model", model,
		)
	}
}

// GetModel returns the model name for the given intelligence level
func GetModel(intelligence Speed) string {
	mu.RLock()
	defer mu.RUnlock()

	// Check for override first
	if model, ok := modelOverrides[intelligence]; ok {
		return model
	}

	// Return default
	return getModel(intelligence)
}

// SetProvider sets the LLM provider (openai, anthropic, etc.)
func SetProvider(p string) {
	mu.Lock()
	defer mu.Unlock()
	provider = p

	if debugMode {
		logger.Info("Provider changed", "provider", provider)
	}
}

// GetProvider returns the current LLM provider
func GetProvider() string {
	mu.RLock()
	defer mu.RUnlock()
	return provider
}

// SetDebugMode enables or disables debug logging
func SetDebugMode(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	debugMode = enabled

	if logger != nil {
		if enabled {
			logger.SetLevel(DebugLevel)
			logger.Debug("Debug mode enabled")
		} else {
			logger.SetLevel(InfoLevel)
			logger.Info("Debug mode disabled")
		}
	}
}

// IsDebugMode returns whether debug mode is enabled
func IsDebugMode() bool {
	mu.RLock()
	defer mu.RUnlock()
	return debugMode
}

// SetMetricsEnabled enables or disables metrics collection
func SetMetricsEnabled(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	metricsEnabled = enabled

	if debugMode {
		logger.Info("Metrics collection changed", "enabled", enabled)
	}
}

// IsMetricsEnabled returns whether metrics collection is enabled
func IsMetricsEnabled() bool {
	mu.RLock()
	defer mu.RUnlock()
	return metricsEnabled
}
