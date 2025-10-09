// package ops - Complete operation for intelligent text completion
package ops

import (
	"context"
	"fmt"
	"strings"

	"github.com/monstercameron/SchemaFlow/core"
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
	core.OpOptions
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
		OpOptions: core.OpOptions{
			Mode:         core.Creative,
			Intelligence: core.Smart,
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
func (opts CompleteOptions) WithIntelligence(intelligence core.Speed) CompleteOptions {
	opts.OpOptions.Intelligence = intelligence
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

// toOpOptions converts CompleteOptions to core.OpOptions
func (opts CompleteOptions) toOpOptions() core.OpOptions {
	return opts.OpOptions
}

// Complete intelligently completes partial text using LLM intelligence.
// It uses context from previous messages/text to generate coherent completions.
//
// Examples:
//
//	// Complete a sentence
//	result, err := Complete("The weather today is", NewCompleteOptions().WithMaxLength(50))
//
//	// Complete with context
//	result, err := Complete("Please send me the", NewCompleteOptions().
//	    WithContext([]string{"User: I need help with my order", "Assistant: I'd be happy to help!"}).
//	    WithMaxLength(100))
//
//	// Complete code with stop sequences
//	result, err := Complete("function calculateTotal(items) {", NewCompleteOptions().
//	    WithStopSequences([]string{"}", "\n\n"}).
//	    WithMaxLength(200))
//
// The operation automatically handles different content types and maintains coherence.
func Complete(partialText string, opts CompleteOptions) (CompleteResult, error) {
	return completeImpl(partialText, opts)
}

// ClientComplete is the client-based version of Complete
func ClientComplete(c *core.Client, partialText string, opts CompleteOptions) (CompleteResult, error) {
	return completeImpl(partialText, opts)
}

func completeImpl(partialText string, opts CompleteOptions) (CompleteResult, error) {
	result := CompleteResult{
		Original: partialText,
		Metadata: make(map[string]any),
	}

	// Validate options
	if err := opts.Validate(); err != nil {
		return result, fmt.Errorf("invalid options: %w", err)
	}

	// Validate input
	if strings.TrimSpace(partialText) == "" {
		return result, fmt.Errorf("partial text cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	// Build prompt
	systemPrompt := buildCompleteSystemPrompt(opts)
	userPrompt := buildCompleteUserPrompt(partialText, opts)

	// Call LLM
	opOpts := opts.toOpOptions()
	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opOpts)
	if err != nil {
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

	return result, nil
}

// buildCompleteSystemPrompt creates the system prompt for completion
func buildCompleteSystemPrompt(opts CompleteOptions) string {
	prompt := "You are an intelligent text completion assistant. Complete the given partial text naturally and coherently.\n\n"

	switch opts.OpOptions.Mode {
	case core.Strict:
		prompt += "Complete the text following strict grammatical and logical rules. Maintain formal tone and precision.\n"
	case core.TransformMode:
		prompt += "Complete the text in a balanced, natural way. Adapt to the existing style and context.\n"
	case core.Creative:
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

	// Limit length
	if len(response) > opts.MaxLength {
		response = response[:opts.MaxLength]
		// Try to cut at word boundary
		if lastSpace := strings.LastIndex(response, " "); lastSpace > len(response)/2 {
			response = response[:lastSpace]
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
