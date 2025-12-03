package ops

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/monstercameron/SchemaFlow/internal/llm"
	"github.com/monstercameron/SchemaFlow/internal/telemetry"
	"github.com/monstercameron/SchemaFlow/internal/types"
)

// CompleteResult contains the completion result and metadata
type CompleteResult struct {
	Text       string         `json:"text"`       // The completed text
	Original   string         `json:"original"`   // The original partial text
	Length     int            `json:"length"`     // Length of completion
	Confidence float64        `json:"confidence"` // Confidence score (0.0-1.0)
	Metadata   map[string]any `json:"metadata"`   // Additional metadata
}

// CompleteOptions configures the Complete operation
type CompleteOptions struct {
	types.OpOptions
	Context       []string // Previous context messages/text
	MaxLength     int      // Maximum length of completion
	StopSequences []string // Sequences that stop generation
	Temperature   float32  // Creativity level (0.0-2.0)
	TopP          float32  // Nucleus sampling (0.0-1.0)
	TopK          int      // Top-k sampling
}

// NewCompleteOptions creates CompleteOptions with defaults
func NewCompleteOptions() CompleteOptions {
	return CompleteOptions{
		OpOptions: types.OpOptions{
			Mode:         types.Creative,
			Intelligence: types.Smart,
		},
		MaxLength:   100,
		Temperature: 0.7,
		TopP:        0.9,
		TopK:        50,
	}
}

// WithContext provides previous messages/text for context
func (opts CompleteOptions) WithContext(context []string) CompleteOptions {
	opts.Context = context
	return opts
}

// WithMaxLength sets the maximum completion length
func (opts CompleteOptions) WithMaxLength(maxLength int) CompleteOptions {
	opts.MaxLength = maxLength
	return opts
}

// WithStopSequences sets sequences that stop generation
func (opts CompleteOptions) WithStopSequences(sequences []string) CompleteOptions {
	opts.StopSequences = sequences
	return opts
}

// WithTemperature sets the creativity level
func (opts CompleteOptions) WithTemperature(temperature float32) CompleteOptions {
	opts.Temperature = temperature
	return opts
}

// WithTopP sets nucleus sampling parameter
func (opts CompleteOptions) WithTopP(topP float32) CompleteOptions {
	opts.TopP = topP
	return opts
}

// WithTopK sets top-k sampling parameter
func (opts CompleteOptions) WithTopK(topK int) CompleteOptions {
	opts.TopK = topK
	return opts
}

// WithIntelligence sets the intelligence level
func (opts CompleteOptions) WithIntelligence(intelligence types.Speed) CompleteOptions {
	opts.OpOptions.Intelligence = intelligence
	return opts
}

// WithMode sets the reasoning mode
func (opts CompleteOptions) WithMode(mode types.Mode) CompleteOptions {
	opts.OpOptions.Mode = mode
	return opts
}

// Validate validates CompleteOptions
func (opts CompleteOptions) Validate() error {
	if opts.MaxLength <= 0 {
		return fmt.Errorf("MaxLength must be positive")
	}
	if opts.Temperature < 0 || opts.Temperature > 2.0 {
		return fmt.Errorf("Temperature must be between 0.0 and 2.0")
	}
	if opts.TopP <= 0 || opts.TopP > 1.0 {
		return fmt.Errorf("TopP must be between 0.0 and 1.0")
	}
	if opts.TopK <= 0 {
		return fmt.Errorf("TopK must be positive")
	}
	return nil
}

// toOpOptions converts CompleteOptions to types.OpOptions
func (opts CompleteOptions) toOpOptions() types.OpOptions {
	return opts.OpOptions
}

// Complete intelligently completes partial text using LLM intelligence.
// It uses context from previous messages/text to generate coherent completions.
// Note: This function requires a provider to be passed.
func Complete(ctx context.Context, provider llm.Provider, partialText string, opts CompleteOptions) (CompleteResult, error) {
	return completeImpl(ctx, provider, partialText, opts)
}

func completeImpl(ctx context.Context, provider llm.Provider, partialText string, opts CompleteOptions) (CompleteResult, error) {
	logger := telemetry.GetLogger()
	logger.Debug("Starting complete operation", "requestID", opts.RequestID, "partialLength", len(partialText))

	result := CompleteResult{
		Original: partialText,
		Metadata: make(map[string]any),
	}

	// Validate options
	if err := opts.Validate(); err != nil {
		logger.Error("Complete operation validation failed", "requestID", opts.RequestID, "error", err)
		return result, fmt.Errorf("invalid options: %w", err)
	}

	// Validate input
	if strings.TrimSpace(partialText) == "" {
		logger.Error("Complete operation failed: empty input", "requestID", opts.RequestID)
		return result, fmt.Errorf("partial text cannot be empty")
	}

	// Use provided context or create one with timeout
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Build prompt
	systemPrompt := buildCompleteSystemPrompt(opts)
	userPrompt := buildCompleteUserPrompt(partialText, opts)

	// Call LLM
	opOpts := opts.toOpOptions()
	response, err := CallLLM(ctx, provider, systemPrompt, userPrompt, opOpts)
	if err != nil {
		logger.Error("Complete operation LLM call failed", "requestID", opts.RequestID, "error", err)
		return result, fmt.Errorf("completion failed: %w", err)
	}

	// Process response
	completedText := processCompletionResponse(response, partialText, opts)

	result.Text = completedText
	result.Length = len(completedText) - len(partialText)

	// Basic confidence estimation
	result.Confidence = estimateCompletionConfidence(completedText, partialText)

	// Add metadata
	result.Metadata["model"] = "llm-completion"
	result.Metadata["temperature"] = opts.Temperature
	result.Metadata["max_length"] = opts.MaxLength
	result.Metadata["context_messages"] = len(opts.Context)

	logger.Debug("Complete operation succeeded", "requestID", opts.RequestID, "completionLength", result.Length, "confidence", result.Confidence)

	return result, nil
}

// buildCompleteSystemPrompt creates the system prompt for completion
func buildCompleteSystemPrompt(opts CompleteOptions) string {
	prompt := "You are an intelligent text completion assistant. Complete the given partial text naturally and coherently.\n\n"

	switch opts.OpOptions.Mode {
	case types.Strict:
		prompt += "Complete the text following strict grammatical and logical rules. Maintain formal tone and precision.\n"
	case types.TransformMode:
		prompt += "Complete the text in a balanced, natural way. Adapt to the existing style and context.\n"
	case types.Creative:
		prompt += "Complete the text creatively and engagingly. Feel free to be imaginative while staying relevant.\n"
	}

	if len(opts.Context) > 0 {
		prompt += "Use the provided context to understand the conversation flow and maintain coherence.\n"
	}

	if len(opts.StopSequences) > 0 {
		prompt += fmt.Sprintf("Stop generation when you encounter any of these sequences: %s\n",
			strings.Join(opts.StopSequences, ", "))
	}

	prompt += fmt.Sprintf("Limit your completion to approximately %d characters.\n", opts.MaxLength)
	prompt += "Return only the completion text, no explanations or metadata."

	return prompt
}

// buildCompleteUserPrompt creates the user prompt for completion
func buildCompleteUserPrompt(partialText string, opts CompleteOptions) string {
	prompt := ""

	// Add context if provided
	if len(opts.Context) > 0 {
		prompt += "Context:\n"
		for i, ctx := range opts.Context {
			prompt += fmt.Sprintf("%d. %s\n", i+1, ctx)
		}
		prompt += "\n"
	}

	prompt += fmt.Sprintf("Complete this text: %s", partialText)

	return prompt
}

// processCompletionResponse processes the LLM response
func processCompletionResponse(response, originalText string, opts CompleteOptions) string {
	// Clean up the response
	response = strings.TrimSpace(response)

	// Remove any unwanted prefixes that might repeat the original text
	if strings.HasPrefix(response, originalText) {
		response = strings.TrimPrefix(response, originalText)
		response = strings.TrimSpace(response)
	}

	// Apply stop sequences
	for _, stopSeq := range opts.StopSequences {
		if idx := strings.Index(response, stopSeq); idx != -1 {
			response = response[:idx]
			break
		}
	}

	// Combine original with completion
	completedText := originalText
	if response != "" {
		// Add space if needed
		if !strings.HasSuffix(originalText, " ") && !strings.HasPrefix(response, " ") {
			completedText += " "
		}
		completedText += response
	}

	// Limit total length
	maxTotalLength := len(originalText) + opts.MaxLength
	if len(completedText) > maxTotalLength {
		completedText = completedText[:maxTotalLength]
		// Try to cut at word boundary
		if lastSpace := strings.LastIndex(completedText, " "); lastSpace > len(originalText) {
			completedText = completedText[:lastSpace]
		}
	}

	return completedText
}

// estimateCompletionConfidence provides a basic confidence estimate
func estimateCompletionConfidence(completed, original string) float64 {
	if len(completed) <= len(original) {
		return 0.0 // No completion generated
	}

	completion := strings.TrimSpace(completed[len(original):])
	if completion == "" {
		return 0.0
	}

	// Basic heuristics for confidence
	confidence := 0.5 // Base confidence

	// Longer completions tend to be more confident
	if len(completion) > 20 {
		confidence += 0.2
	}

	// Completions with proper punctuation
	if strings.Contains(completion, ".") || strings.Contains(completion, "!") || strings.Contains(completion, "?") {
		confidence += 0.1
	}

	// Avoid overconfidence
	if confidence > 0.9 {
		confidence = 0.9
	}

	return confidence
}
