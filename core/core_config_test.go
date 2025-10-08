package schemaflow

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
