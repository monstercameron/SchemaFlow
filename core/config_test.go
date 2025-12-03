package core

import (
	"testing"
	"time"
)

func TestGetMaxTokens(t *testing.T) {
	tests := []struct {
		speed    Speed
		expected int
	}{
		{Smart, 4000},
		{Fast, 2000},
		{Quick, 1000},
		{Speed(99), 2000},
	}

	for _, tt := range tests {
		t.Run(tt.speed.String(), func(t *testing.T) {
			if got := getMaxTokens(tt.speed); got != tt.expected {
				t.Errorf("getMaxTokens(%v) = %v, want %v", tt.speed, got, tt.expected)
			}
		})
	}
}

func TestGetModel(t *testing.T) {
	// Save original state
	originalProvider := provider
	defer func() { provider = originalProvider }()

	tests := []struct {
		name         string
		providerName string
		speed        Speed
		want         string
	}{
		{"OpenAI Smart", "openai", Smart, "gpt-5-2025-08-07"},
		{"OpenAI Fast", "openai", Fast, "gpt-5-nano-2025-08-07"},
		{"OpenRouter Smart", "openrouter", Smart, "openai/gpt-4o"},
		{"OpenRouter Fast", "openrouter", Fast, "openai/gpt-4o-mini"},
		{"Cerebras Smart", "cerebras", Smart, "llama-3.3-70b"},
		{"Cerebras Fast", "cerebras", Fast, "llama3.1-8b"},
		{"Anthropic Smart", "anthropic", Smart, "claude-3-5-sonnet-20240620"},
		{"Anthropic Fast", "anthropic", Fast, "claude-3-haiku-20240307"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider = tt.providerName
			if got := getModel(tt.speed); got != tt.want {
				t.Errorf("getModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetModelEnvOverride(t *testing.T) {
	// Save original state
	originalProvider := provider
	defer func() { provider = originalProvider }()

	t.Setenv("SCHEMAFLOW_MODEL", "custom-model-override")
	
	provider = "openai"
	if got := getModel(Smart); got != "custom-model-override" {
		t.Errorf("Expected env override 'custom-model-override', got %v", got)
	}

	provider = "anthropic"
	if got := getModel(Fast); got != "custom-model-override" {
		t.Errorf("Expected env override 'custom-model-override', got %v", got)
	}
}

func TestGetModelLevelOverride(t *testing.T) {
	// Save original state
	originalProvider := provider
	defer func() { provider = originalProvider }()

	t.Setenv("SCHEMAFLOW_MODEL_SMART", "smart-override")
	t.Setenv("SCHEMAFLOW_MODEL_FAST", "fast-override")
	
	provider = "openai"
	
	if got := getModel(Smart); got != "smart-override" {
		t.Errorf("Expected smart override 'smart-override', got %v", got)
	}

	if got := getModel(Fast); got != "fast-override" {
		t.Errorf("Expected fast override 'fast-override', got %v", got)
	}
	
	// Quick should still be default
	if got := getModel(Quick); got == "smart-override" || got == "fast-override" {
		t.Errorf("Expected default for Quick, got %v", got)
	}
}

func TestRecordMetric(t *testing.T) {
	original := metricsEnabled
	defer func() { metricsEnabled = original }()

	metricsEnabled = false
	recordMetric("test_metric", 100, map[string]string{"tag": "value"})

	metricsEnabled = true
	recordMetric("test_metric", 100, map[string]string{"tag": "value"})
}

func TestInitWithEnvironmentVariables(t *testing.T) {
	originalTimeout := timeout
	originalMaxRetries := maxRetries
	originalRetryBackoff := retryBackoff
	originalDebugMode := debugMode

	defer func() {
		timeout = originalTimeout
		maxRetries = originalMaxRetries
		retryBackoff = originalRetryBackoff
		debugMode = originalDebugMode
	}()

	t.Setenv("SCHEMAFLOW_TIMEOUT", "invalid")
	t.Setenv("SCHEMAFLOW_MAX_RETRIES", "invalid")
	t.Setenv("SCHEMAFLOW_RETRY_BACKOFF", "invalid")

	Init("")

	if timeout != 30*time.Second {
		t.Errorf("Expected default timeout, got %v", timeout)
	}
}
