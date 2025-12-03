package logger

import (
	"log"
	"os"
	"sync"
)

// Logger defines the interface for logging operations
type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

// DefaultLogger is a simple logger implementation using standard log package
type DefaultLogger struct {
	debugEnabled bool
}

var (
	globalLogger Logger
	mu           sync.RWMutex
)

func init() {
	// Initialize with default logger
	globalLogger = &DefaultLogger{
		debugEnabled: os.Getenv("SCHEMAFLOW_DEBUG") == "true",
	}
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return globalLogger
}

// SetLogger sets the global logger instance
func SetLogger(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = l
}

// GetDebugMode returns whether debug mode is enabled
func GetDebugMode() bool {
	return os.Getenv("SCHEMAFLOW_DEBUG") == "true"
}

// IsMetricsEnabled returns whether metrics collection is enabled
func IsMetricsEnabled() bool {
	return os.Getenv("SCHEMAFLOW_METRICS") == "true"
}

// Debug logs a debug message
func (l *DefaultLogger) Debug(msg string, keysAndValues ...any) {
	if l.debugEnabled {
		log.Printf("[DEBUG] %s %v", msg, keysAndValues)
	}
}

// Info logs an info message
func (l *DefaultLogger) Info(msg string, keysAndValues ...any) {
	log.Printf("[INFO] %s %v", msg, keysAndValues)
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(msg string, keysAndValues ...any) {
	log.Printf("[WARN] %s %v", msg, keysAndValues)
}

// Error logs an error message
func (l *DefaultLogger) Error(msg string, keysAndValues ...any) {
	log.Printf("[ERROR] %s %v", msg, keysAndValues)
}
