package telemetry

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger *Logger
	mu     sync.RWMutex
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// DebugLevel logs everything including debug information
	DebugLevel LogLevel = iota
	// InfoLevel logs informational messages and above
	InfoLevel
	// WarnLevel logs warnings and errors
	WarnLevel
	// ErrorLevel logs only errors
	ErrorLevel
	// FatalLevel logs fatal errors and exits
	FatalLevel
)

// Logger wraps slog.Logger to provide compatibility with the previous API.
type Logger struct {
	*slog.Logger
	level *slog.LevelVar
}

// NewLogger creates a new Logger backed by slog with optional JSON output.
func NewLogger() *Logger {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	handlerOpts := &slog.HandlerOptions{Level: lvl}
	var handler slog.Handler
	if os.Getenv("SCHEMAFLOW_LOG_FORMAT") == "json" {
		handler = slog.NewJSONHandler(os.Stderr, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, handlerOpts)
	}

	return &Logger{
		Logger: slog.New(handler),
		level:  lvl,
	}
}

// SetLevel updates the logging level by mapping the legacy levels to slog levels.
func (l *Logger) SetLevel(level LogLevel) {
	if l == nil || l.level == nil {
		return
	}

	switch level {
	case DebugLevel:
		l.level.Set(slog.LevelDebug)
	case InfoLevel:
		l.level.Set(slog.LevelInfo)
	case WarnLevel:
		l.level.Set(slog.LevelWarn)
	case ErrorLevel:
		l.level.Set(slog.LevelError)
	case FatalLevel:
		l.level.Set(slog.LevelError)
	default:
		l.level.Set(slog.LevelInfo)
	}
}

// WithFields returns a logger augmented with the provided structured attributes.
func (l *Logger) WithFields(fields map[string]any) *Logger {
	if l == nil {
		return nil
	}

	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}

	return &Logger{
		Logger: l.Logger.With(attrs...),
		level:  l.level,
	}
}

// Fatal logs an error message and exits the process.
func (l *Logger) Fatal(message string, args ...any) {
	if l == nil {
		return
	}
	l.Logger.Error(message, args...)
	os.Exit(1)
}

// GetLogger returns the global logger instance, creating one on-demand.
func GetLogger() *Logger {
	mu.RLock()
	if logger != nil {
		mu.RUnlock()
		return logger
	}
	mu.RUnlock()

	mu.Lock()
	if logger == nil {
		logger = NewLogger()
	}
	mu.Unlock()

	return logger
}
