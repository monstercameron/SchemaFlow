package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/monstercameron/SchemaFlow/core"
)

// SuggestResult contains the result of a suggestion operation
type SuggestResult[T any] struct {
	Suggestions []T            `json:"suggestions"`        // The suggested items
	Scores      []float64      `json:"scores,omitempty"`   // Confidence scores (if ranked)
	Reasons     []string       `json:"reasons,omitempty"`  // Reasons for each suggestion
	Metadata    map[string]any `json:"metadata,omitempty"` // Additional metadata
}

// SuggestStrategy defines how suggestions are generated
type SuggestStrategy string

const (
	SuggestContextual SuggestStrategy = "contextual" // Based on current context/state
	SuggestPattern    SuggestStrategy = "pattern"    // Based on patterns in data
	SuggestGoal       SuggestStrategy = "goal"       // Based on stated goals
	SuggestHybrid     SuggestStrategy = "hybrid"     // Combination approach
)

// SuggestOptions configures the Suggest operation
type SuggestOptions struct {
	core.OpOptions
	CommonOptions

	// Suggestion strategy
	Strategy SuggestStrategy

	// Maximum number of suggestions to return
	TopN int

	// Whether to rank suggestions by relevance/confidence
	Ranked bool

	// Include confidence scores for each suggestion
	IncludeScores bool

	// Include reasoning for each suggestion
	IncludeReasons bool

	// Domain context for more relevant suggestions
	Domain string

	// Specific constraints or requirements
	Constraints []string

	// Custom suggestion categories
	Categories []string
}

// NewSuggestOptions creates SuggestOptions with defaults
func NewSuggestOptions() SuggestOptions {
	return SuggestOptions{
		CommonOptions: CommonOptions{
			Mode:         core.TransformMode,
			Intelligence: core.Smart,
		},
		Strategy:       SuggestContextual,
		TopN:           5,
		Ranked:         true,
		IncludeScores:  false,
		IncludeReasons: false,
	}
}

// Validate checks if the options are valid
func (opts SuggestOptions) Validate() error {
	if err := opts.CommonOptions.Validate(); err != nil {
		return err
	}

	if opts.TopN < 1 {
		return fmt.Errorf("topN must be at least 1, got %d", opts.TopN)
	}

	validStrategies := map[SuggestStrategy]bool{
		SuggestContextual: true,
		SuggestPattern:    true,
		SuggestGoal:       true,
		SuggestHybrid:     true,
	}
	if !validStrategies[opts.Strategy] {
		return fmt.Errorf("invalid strategy: %s", opts.Strategy)
	}

	return nil
}

// WithStrategy sets the suggestion strategy
func (opts SuggestOptions) WithStrategy(strategy SuggestStrategy) SuggestOptions {
	opts.Strategy = strategy
	return opts
}

// WithTopN sets the maximum number of suggestions to return
func (opts SuggestOptions) WithTopN(n int) SuggestOptions {
	opts.TopN = n
	return opts
}

// WithRanked enables/disables ranking of suggestions
func (opts SuggestOptions) WithRanked(ranked bool) SuggestOptions {
	opts.Ranked = ranked
	return opts
}

// WithIncludeScores includes confidence scores for suggestions
func (opts SuggestOptions) WithIncludeScores(include bool) SuggestOptions {
	opts.IncludeScores = include
	return opts
}

// WithIncludeReasons includes reasoning for each suggestion
func (opts SuggestOptions) WithIncludeReasons(include bool) SuggestOptions {
	opts.IncludeReasons = include
	return opts
}

// WithDomain sets the domain context for suggestions
func (opts SuggestOptions) WithDomain(domain string) SuggestOptions {
	opts.Domain = domain
	return opts
}

// WithConstraints sets specific constraints for suggestions
func (opts SuggestOptions) WithConstraints(constraints []string) SuggestOptions {
	opts.Constraints = constraints
	return opts
}

// WithCategories sets suggestion categories
func (opts SuggestOptions) WithCategories(categories []string) SuggestOptions {
	opts.Categories = categories
	return opts
}

// Suggest generates context-aware suggestions based on input data and current state
//
// Examples:
//
//	// Basic suggestions
//	suggestions, err := Suggest[Action](currentState, NewSuggestOptions())
//
//	// Ranked suggestions with reasons
//	suggestions, err := Suggest[string](context, NewSuggestOptions().
//	    WithRanked(true).
//	    WithTopN(5).
//	    WithIncludeReasons(true))
//
//	// Domain-specific suggestions
//	suggestions, err := Suggest[ConfigOption](currentConfig,
//	    NewSuggestOptions().WithDomain("data-processing"))
func Suggest[T any](input any, opts SuggestOptions) ([]T, error) {
	logger := core.GetLogger()
	logger.Debug("Starting suggest operation", "requestID", opts.CommonOptions.RequestID, "inputType", fmt.Sprintf("%T", input), "topN", opts.TopN)

	// Validate options
	if err := opts.Validate(); err != nil {
		logger.Error("Suggest operation validation failed", "requestID", opts.CommonOptions.RequestID, "error", err)
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	opOptions := opts.toOpOptions()

	// Build suggestion instructions
	var instructions []string

	instructions = append(instructions, fmt.Sprintf("Generate suggestions using %s strategy", opts.Strategy))

	if opts.Domain != "" {
		instructions = append(instructions, fmt.Sprintf("Domain context: %s", opts.Domain))
	}

	if len(opts.Constraints) > 0 {
		instructions = append(instructions, fmt.Sprintf("Constraints: %s", strings.Join(opts.Constraints, ", ")))
	}

	if len(opts.Categories) > 0 {
		instructions = append(instructions, fmt.Sprintf("Categories: %s", strings.Join(opts.Categories, ", ")))
	}

	if opts.Ranked {
		instructions = append(instructions, "Rank suggestions by relevance")
	}

	if opts.IncludeScores {
		instructions = append(instructions, "Include confidence scores (0-1)")
	}

	if opts.IncludeReasons {
		instructions = append(instructions, "Provide reasoning for each suggestion")
	}

	instructions = append(instructions, fmt.Sprintf("Return top %d suggestions", opts.TopN))

	steering := strings.Join(instructions, ". ")
	if opts.CommonOptions.Steering != "" {
		steering = opts.CommonOptions.Steering + ". " + steering
	}
	opOptions.Steering = steering

	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	// Marshal input for LLM
	inputJSON, err := json.Marshal(input)
	if err != nil {
		logger.Error("Suggest operation marshal failed", "requestID", opts.CommonOptions.RequestID, "error", err)
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	systemPrompt := `You are an expert suggestion engine. Generate contextually relevant suggestions based on the provided input and requirements.

Rules:
- Analyze the input data and current context
- Generate practical, actionable suggestions
- Consider the specified constraints and domain
- Return suggestions as a JSON array
- Each suggestion should be relevant and helpful
- If ranking is requested, order by relevance (most relevant first)
- Include scores and reasons only if specifically requested`

	userPrompt := fmt.Sprintf("Generate suggestions based on this input:\n%s", string(inputJSON))

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opOptions)
	if err != nil {
		logger.Error("Suggest operation LLM call failed", "requestID", opts.CommonOptions.RequestID, "error", err)
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// Parse the response
	var suggestions []T
	if err := json.Unmarshal([]byte(response), &suggestions); err != nil {
		// Try to extract from a structured response
		var result struct {
			Suggestions []T `json:"suggestions"`
		}
		if err := json.Unmarshal([]byte(response), &result); err != nil {
			logger.Error("Suggest operation parse failed", "requestID", opts.CommonOptions.RequestID, "error", err)
			return nil, fmt.Errorf("failed to parse suggestions: %w", err)
		}
		suggestions = result.Suggestions
	}

	// Limit to TopN
	if len(suggestions) > opts.TopN {
		suggestions = suggestions[:opts.TopN]
	}

	logger.Debug("Suggest operation succeeded", "requestID", opts.CommonOptions.RequestID, "suggestionsCount", len(suggestions))

	return suggestions, nil
}

// SuggestWithResult provides detailed suggestion results with scores and reasons
func SuggestWithResult[T any](input any, opts SuggestOptions) (SuggestResult[T], error) {
	result := SuggestResult[T]{
		Metadata: make(map[string]any),
	}

	// Enable scores and reasons if not already set
	opts.IncludeScores = opts.IncludeScores || opts.Ranked
	opts.IncludeReasons = opts.IncludeReasons || true

	suggestions, err := Suggest[T](input, opts)
	if err != nil {
		return result, err
	}

	result.Suggestions = suggestions

	// For now, return basic result - full implementation would parse scores/reasons from LLM
	// This would require more sophisticated prompt engineering and response parsing

	return result, nil
}
