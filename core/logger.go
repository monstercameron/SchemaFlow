// Package schemaflow - Structured logging implementation
package schemaflow

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
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

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging with levels and JSON output
type Logger struct {
	mu       sync.RWMutex
	level    LogLevel
	output   io.Writer
	jsonMode bool
	fields   map[string]any
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Fields    map[string]any `json:"fields,omitempty"`
	Caller    string         `json:"caller,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		level:    InfoLevel,
		output:   os.Stderr,
		jsonMode: os.Getenv("SCHEMAFLOW_LOG_FORMAT") == "json",
		fields:   make(map[string]any),
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput sets the output writer for logs
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// SetJSONMode enables or disables JSON logging
func (l *Logger) SetJSONMode(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.jsonMode = enabled
}

// WithFields returns a new logger with additional fields
func (l *Logger) WithFields(fields map[string]any) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	newFields := make(map[string]any)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	
	return &Logger{
		level:    l.level,
		output:   l.output,
		jsonMode: l.jsonMode,
		fields:   newFields,
	}
}

// log writes a log entry at the specified level
func (l *Logger) log(level LogLevel, message string, fields ...any) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	// Check if we should log this level
	if level < l.level {
		return
	}
	
	// Create log entry
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level.String(),
		Message:   message,
		Fields:    make(map[string]any),
	}
	
	// Add persistent fields
	for k, v := range l.fields {
		entry.Fields[k] = v
	}
	
	// Parse variadic fields
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			entry.Fields[key] = fields[i+1]
		}
	}
	
	// Add request ID if present
	if reqID, ok := entry.Fields["requestID"].(string); ok && reqID != "" {
		entry.RequestID = reqID
		delete(entry.Fields, "requestID")
	}
	
	// Add caller information for errors and debug
	if level >= ErrorLevel || level == DebugLevel {
		if caller := getCaller(3); caller != "" {
			entry.Caller = caller
		}
	}
	
	// Format and write the log
	var output string
	if l.jsonMode {
		if data, err := json.Marshal(entry); err == nil {
			output = string(data) + "\n"
		} else {
			output = fmt.Sprintf("LOG_ERROR: %v\n", err)
		}
	} else {
		output = formatTextLog(entry)
	}
	
	fmt.Fprint(l.output, output)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...any) {
	l.log(DebugLevel, message, fields...)
}

// Info logs an informational message
func (l *Logger) Info(message string, fields ...any) {
	l.log(InfoLevel, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...any) {
	l.log(WarnLevel, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...any) {
	l.log(ErrorLevel, message, fields...)
}

// Fatal logs a fatal error and exits
func (l *Logger) Fatal(message string, fields ...any) {
	l.log(FatalLevel, message, fields...)
	os.Exit(1)
}

// formatTextLog formats a log entry as human-readable text
func formatTextLog(entry LogEntry) string {
	var sb strings.Builder
	
	// Timestamp and level
	sb.WriteString(entry.Timestamp)
	sb.WriteString(" [")
	sb.WriteString(entry.Level)
	sb.WriteString("] ")
	
	// Request ID if present
	if entry.RequestID != "" {
		sb.WriteString("[")
		sb.WriteString(entry.RequestID)
		sb.WriteString("] ")
	}
	
	// Message
	sb.WriteString(entry.Message)
	
	// Fields
	if len(entry.Fields) > 0 {
		sb.WriteString(" | ")
		first := true
		for k, v := range entry.Fields {
			if !first {
				sb.WriteString(", ")
			}
			first = false
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(formatValue(v))
		}
	}
	
	// Caller
	if entry.Caller != "" {
		sb.WriteString(" | caller=")
		sb.WriteString(entry.Caller)
	}
	
	sb.WriteString("\n")
	return sb.String()
}

// formatValue formats a value for text logging
func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		if strings.Contains(val, " ") || strings.Contains(val, "=") {
			return fmt.Sprintf("%q", val)
		}
		return val
	case error:
		return val.Error()
	default:
		return fmt.Sprintf("%v", val)
	}
}

// getCaller returns the caller information
func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	
	// Extract just the filename
	parts := strings.Split(file, "/")
	if len(parts) > 0 {
		file = parts[len(parts)-1]
	}
	
	
	return fmt.Sprintf("%s:%d", file, line)
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	if logger == nil {
		return NewLogger()
	}
	return logger
}
