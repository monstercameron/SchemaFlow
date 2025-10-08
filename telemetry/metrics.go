// package telemetry - placeholder for metrics recording
package telemetry

import (
	"fmt"

	"github.com/monstercameron/SchemaFlow/core"
)

// RecordMetric records a metric. This is a placeholder implementation.
func RecordMetric(name string, value int64, tags map[string]string) {
	// For now, we'll just log the metric call if debugging is enabled.
	// A proper implementation would send this to a metrics backend like Prometheus or Datadog.
	if core.IsMetricsEnabled() {
		if core.GetDebugMode() {
			core.GetLogger().Debug(fmt.Sprintf("Recording metric: %s = %d", name, value), "tags", tags)
		}
	}
}
