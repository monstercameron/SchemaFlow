package core

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var requestIDCounter uint64

// generateRequestID generates a unique request ID for tracing
func generateRequestID() string {
	// Use atomic counter to ensure uniqueness even when called in quick succession
	counter := atomic.AddUint64(&requestIDCounter, 1)
	return fmt.Sprintf("%d-%d-%d", time.Now().UnixNano(), os.Getpid(), counter)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// recordMetric records a metric for observability (placeholder for actual implementation)
func recordMetric(name string, value int64, tags map[string]string) {
	if !metricsEnabled {
		return
	}

	// This is a placeholder - in production, integrate with your metrics system
	// (Prometheus, DataDog, CloudWatch, etc.)
	if debugMode {
		logger.Debug("Metric recorded",
			"name", name,
			"value", value,
			"tags", tags,
		)
	}
}
