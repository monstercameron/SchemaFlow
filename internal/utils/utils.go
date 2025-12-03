package utils

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/monstercameron/SchemaFlow/internal/config"
	"github.com/monstercameron/SchemaFlow/internal/logger"
)

var requestIDCounter uint64

// GenerateRequestID generates a unique request ID for tracing
func GenerateRequestID() string {
	// Use atomic counter to ensure uniqueness even when called in quick succession
	counter := atomic.AddUint64(&requestIDCounter, 1)
	return fmt.Sprintf("%d-%d-%d", time.Now().UnixNano(), os.Getpid(), counter)
}

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RecordMetric records a metric for observability (placeholder for actual implementation)
func RecordMetric(name string, value int64, tags map[string]string) {
	if !config.IsMetricsEnabled() {
		return
	}

	// This is a placeholder - in production, integrate with your metrics system
	// (Prometheus, DataDog, CloudWatch, etc.)
	if config.GetDebugMode() {
		logger.GetLogger().Debug("Metric recorded",
			"name", name,
			"value", value,
			"tags", tags,
		)
	}
}
