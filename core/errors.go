package core

import (
	"errors"
	"fmt"
	"time"
)

// ExtractError indicates a failure in the Extract operation.
// Contains context about what was being extracted and why it failed.
type ExtractError struct {
	Input      any       // The input that failed to extract
	TargetType string    // The Go type we tried to extract to
	Reason     string    // Human-readable failure reason
	Confidence float64   // How confident the model was (0.0-1.0)
	RequestID  string    // Trace ID for debugging
	Timestamp  time.Time // When the error occurred
}

// Error implements the error interface for ExtractError
func (e ExtractError) Error() string {
	return fmt.Sprintf("[%s] extract error: %s (type: %s, confidence: %.2f, request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.TargetType, e.Confidence, e.RequestID)
}

// Unwrap returns the underlying error if any
func (e ExtractError) Unwrap() error {
	return errors.New(e.Reason)
}

// GenerateError indicates a failure in the Generate operation.
type GenerateError struct {
	Prompt     string // The generation prompt that failed
	TargetType string // The type we tried to generate
	Reason     string // Human-readable failure reason
	RequestID  string // Trace ID for debugging
	Timestamp  time.Time
}

// Error implements the error interface for GenerateError
func (e GenerateError) Error() string {
	return fmt.Sprintf("[%s] generate error: %s (type: %s, request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.TargetType, e.RequestID)
}

// TransformError indicates a failure in the Transform operation.
type TransformError struct {
	Input      any     // The input value that failed to transform
	FromType   string  // Source type name
	ToType     string  // Target type name
	Reason     string  // Human-readable failure reason
	Confidence float64 // Model confidence
	RequestID  string  // Trace ID
	Timestamp  time.Time
}

// Error implements the error interface for TransformError
func (e TransformError) Error() string {
	return fmt.Sprintf("[%s] transform error: %s (from: %s, to: %s, request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.FromType, e.ToType, e.RequestID)
}

// ClassifyError indicates a failure in the Classify operation.
type ClassifyError struct {
	Input      string   // The text that failed to classify
	Categories []string // Available categories
	Reason     string   // Failure reason
	Confidence float64  // Model confidence
	RequestID  string   // Trace ID
	Timestamp  time.Time
}

// Error implements the error interface for ClassifyError
func (e ClassifyError) Error() string {
	return fmt.Sprintf("[%s] classify error: %s (categories: %v, request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.Categories, e.RequestID)
}

// SummarizeError indicates a failure in the Summarize operation.
type SummarizeError struct {
	Input     string // The text that failed to summarize
	Length    int    // Input text length
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for SummarizeError
func (e SummarizeError) Error() string {
	return fmt.Sprintf("[%s] summarize error: %s (length: %d, request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.Length, e.RequestID)
}

// TranslateError indicates a failure in the Translate operation.
type TranslateError struct {
	Input     string // The text that failed to translate
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for TranslateError
func (e TranslateError) Error() string {
	return fmt.Sprintf("[%s] translate error: %s (request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.RequestID)
}

// RewriteError indicates a failure in the Rewrite operation.
type RewriteError struct {
	Input     string // The text that failed to rewrite
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for RewriteError
func (e RewriteError) Error() string {
	return fmt.Sprintf("[%s] rewrite error: %s (request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.RequestID)
}

// ExpandError indicates a failure in the Expand operation.
type ExpandError struct {
	Input     string // The text that failed to expand
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for ExpandError
func (e ExpandError) Error() string {
	return fmt.Sprintf("[%s] expand error: %s (request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.RequestID)
}

// CompareError indicates a failure in the Compare operation.
type CompareError struct {
	A, B      any    // The items that failed to compare
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for CompareError
func (e CompareError) Error() string {
	return fmt.Sprintf("[%s] compare error: %s (request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.RequestID)
}

// SimilarError indicates a failure in the Similar operation.
type SimilarError struct {
	Input, Target string // The texts that failed similarity check
	Reason        string // Failure reason
	RequestID     string // Trace ID
	Timestamp     time.Time
}

// Error implements the error interface for SimilarError
func (e SimilarError) Error() string {
	return fmt.Sprintf("[%s] similar error: %s (request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.RequestID)
}

// ScoreError indicates a failure in the Score operation.
type ScoreError struct {
	Input     any    // The input that failed to score
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for ScoreError
func (e ScoreError) Error() string {
	return fmt.Sprintf("[%s] score error: %s (request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.RequestID)
}

// ChooseError indicates a failure in the Choose operation.
type ChooseError struct {
	Options   []any  // The options that failed to choose from
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for ChooseError
func (e ChooseError) Error() string {
	return fmt.Sprintf("[%s] choose error: %s (options: %d, request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, len(e.Options), e.RequestID)
}

// FilterError indicates a failure in the Filter operation.
type FilterError struct {
	Items     []any  // The items that failed to filter
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for FilterError
func (e FilterError) Error() string {
	return fmt.Sprintf("[%s] filter error: %s (items: %d, request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, len(e.Items), e.RequestID)
}

// SortError indicates a failure in the Sort operation.
type SortError struct {
	Items     []any  // The items that failed to sort
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for SortError
func (e SortError) Error() string {
	return fmt.Sprintf("[%s] sort error: %s (items: %d, request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, len(e.Items), e.RequestID)
}

// MatchError indicates a failure in the Match operation.
type MatchError struct {
	Input     any    // The input that failed to match
	Cases     int    // Number of cases evaluated
	Reason    string // Failure reason
	RequestID string // Trace ID
	Timestamp time.Time
}

// Error implements the error interface for MatchError
func (e MatchError) Error() string {
	return fmt.Sprintf("[%s] match error: %s (cases: %d, request: %s)",
		e.Timestamp.Format(time.RFC3339), e.Reason, e.Cases, e.RequestID)
}
