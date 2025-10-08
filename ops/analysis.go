package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/monstercameron/SchemaFlow/core"
)

// Classify categorizes text into predefined categories with specialized options.
//
// Examples:
//
//	// Basic classification
//	category, err := Classify("This product is amazing!", NewClassifyOptions().
//	    WithCategories([]string{"positive", "negative", "neutral"}))
//
//	// Multi-label classification
//	categories, err := Classify(text, NewClassifyOptions().
//	    WithCategories([]string{"tech", "business", "sports"}).
//	    WithMultiLabel(true).
//	    WithMaxCategories(2))
func Classify(input string, opts ClassifyOptions) (string, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return "", fmt.Errorf("invalid options: %w", err)
	}

	categories := opts.Categories
	opt := opts.toOpOptions()

	// Build classification instructions
	var instructions []string
	if opts.MultiLabel {
		instructions = append(instructions, "Allow multiple categories")
		if opts.MaxCategories > 0 {
			instructions = append(instructions, fmt.Sprintf("Return at most %d categories", opts.MaxCategories))
		}
	}

	if opts.MinConfidence > 0 {
		instructions = append(instructions, fmt.Sprintf("Minimum confidence: %.2f", opts.MinConfidence))
	}

	if len(opts.CategoryDescriptions) > 0 {
		for category, description := range opts.CategoryDescriptions {
			instructions = append(instructions, fmt.Sprintf("%s: %s", category, description))
		}
	}

	if len(instructions) > 0 {
		steering := strings.Join(instructions, ". ")
		if opts.OpOptions.Steering != "" {
			steering = opts.OpOptions.Steering + ". " + steering
		}
		opt.Steering = steering
	}

	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	categoriesJSON, _ := json.Marshal(categories)

	systemPrompt := fmt.Sprintf(`You are a text classification expert. Classify the input into one of the provided categories.

Categories: %s

Rules:
- Choose the most appropriate category
- Consider context and nuance
- Return ONLY the category name, nothing else`, string(categoriesJSON))

	userPrompt := fmt.Sprintf("Classify this text:\n%s", input)

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return "", core.ClassifyError{
			Input:      input,
			Categories: categories,
			Reason:     err.Error(),
		}
	}

	result := strings.TrimSpace(response)
	result = strings.Trim(result, "\"'")

	found := false
	for _, cat := range categories {
		if strings.EqualFold(result, cat) {
			found = true
			result = cat
			break
		}
	}

	if !found {
		return "", core.ClassifyError{
			Input:      input,
			Categories: categories,
			Reason:     fmt.Sprintf("invalid category returned: %s", result),
			Confidence: 0.5,
		}
	}

	return result, nil
}

// Score rates content based on specified criteria with specialized options.
//
// Examples:
//
//	// Basic scoring
//	score, err := Score(essay, NewScoreOptions().
//	    WithCriteria([]string{"clarity", "grammar", "relevance"}))
//
//	// Custom scale and rubric
//	score, err := Score(submission, NewScoreOptions().
//	    WithScaleMin(1).
//	    WithScaleMax(5).
//	    WithRubric(map[string]string{
//	        "quality": "Overall quality of work",
//	        "effort": "Evidence of effort and care",
//	    }))
func Score(input any, opts ScoreOptions) (float64, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return 0, fmt.Errorf("invalid options: %w", err)
	}

	opt := opts.toOpOptions()

	// Build scoring instructions
	var instructions []string
	instructions = append(instructions, fmt.Sprintf("Score range: %.1f to %.1f", opts.ScaleMin, opts.ScaleMax))

	if len(opts.Criteria) > 0 {
		instructions = append(instructions, fmt.Sprintf("Criteria: %s", strings.Join(opts.Criteria, ", ")))
	}

	if len(opts.Rubric) > 0 {
		for criterion, description := range opts.Rubric {
			instructions = append(instructions, fmt.Sprintf("%s: %s", criterion, description))
		}
	}

	if len(instructions) > 0 {
		steering := strings.Join(instructions, ". ")
		if opts.OpOptions.Steering != "" {
			steering = opts.OpOptions.Steering + ". " + steering
		}
		opt.Steering = steering
	}

	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	inputStr := fmt.Sprintf("%v", input)
	if str, ok := input.(string); ok {
		inputStr = str
	} else if bytes, err := json.Marshal(input); err == nil {
		inputStr = string(bytes)
	}

	systemPrompt := fmt.Sprintf(`You are a scoring expert. Evaluate the input and assign a numeric score.

Rules:
- Return a score between %.1f and %.1f
- Consider all relevant factors
- Be consistent in your scoring
- Return ONLY the numeric value, nothing else`, opts.ScaleMin, opts.ScaleMax)

	userPrompt := fmt.Sprintf("Score this input:\n%s", inputStr)

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return 0, core.ScoreError{
			Input:  input,
			Reason: err.Error(),
		}
	}

	response = strings.TrimSpace(response)
	response = strings.Trim(response, "\"'")

	score, err := strconv.ParseFloat(response, 64)
	if err != nil {
		return 0, core.ScoreError{
			Input:  input,
			Reason: fmt.Sprintf("failed to parse score: %v", err),
		}
	}

	// Normalize to scale
	if score < opts.ScaleMin {
		score = opts.ScaleMin
	} else if score > opts.ScaleMax {
		score = opts.ScaleMax
	}

	// Normalize if requested
	if opts.Normalize {
		score = (score - opts.ScaleMin) / (opts.ScaleMax - opts.ScaleMin)
	}

	return score, nil
}

// Compare analyzes similarities and differences with specialized options.
//
// Examples:
//
//	// Basic comparison
//	comparison, err := Compare(product1, product2, NewCompareOptions())
//
//	// Detailed comparison with specific aspects
//	comparison, err := Compare(doc1, doc2, NewCompareOptions().
//	    WithComparisonAspects([]string{"content", "style", "accuracy"}).
//	    WithOutputFormat("table").
//	    WithFocusOn("differences"))
func Compare(itemA, itemB any, opts CompareOptions) (string, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return "", fmt.Errorf("invalid options: %w", err)
	}

	opt := opts.toOpOptions()

	// Build comparison instructions
	var instructions []string

	if len(opts.ComparisonAspects) > 0 {
		instructions = append(instructions, fmt.Sprintf("Compare these aspects: %s", strings.Join(opts.ComparisonAspects, ", ")))
	}

	if opts.OutputFormat != "" {
		instructions = append(instructions, fmt.Sprintf("Format as: %s", opts.OutputFormat))
	}

	if opts.FocusOn != "" {
		instructions = append(instructions, fmt.Sprintf("Focus on: %s", opts.FocusOn))
	}

	instructions = append(instructions, fmt.Sprintf("Depth level: %d/10", opts.Depth))

	if len(instructions) > 0 {
		steering := strings.Join(instructions, ". ")
		if opts.OpOptions.Steering != "" {
			steering = opts.OpOptions.Steering + ". " + steering
		}
		opt.Steering = steering
	}

	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	itemAString := fmt.Sprintf("%v", itemA)
	if s, ok := itemA.(string); ok {
		itemAString = s
	} else if bytes, err := json.Marshal(itemA); err == nil {
		itemAString = string(bytes)
	}

	itemBString := fmt.Sprintf("%v", itemB)
	if s, ok := itemB.(string); ok {
		itemBString = s
	} else if bytes, err := json.Marshal(itemB); err == nil {
		itemBString = string(bytes)
	}

	systemPrompt := `You are a comparison expert. Analyze and compare two items, highlighting similarities and differences.

Rules:
- Be objective and balanced
- Identify key similarities
- Identify key differences
- Provide clear, structured comparison`

	userPrompt := fmt.Sprintf("Compare these two items:\n\nItem A:\n%s\n\nItem B:\n%s", itemAString, itemBString)

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return "", core.CompareError{
			A:      itemA,
			B:      itemB,
			Reason: err.Error(),
		}
	}

	return strings.TrimSpace(response), nil
}

// SimilarOptions configures the Similar operation
type SimilarOptions struct {
	core.OpOptions
	SimilarityThreshold float64  // Threshold for similarity (0-1)
	Aspects             []string // Specific aspects to compare
}

// NewSimilarOptions creates SimilarOptions with defaults
func NewSimilarOptions() SimilarOptions {
	return SimilarOptions{
		OpOptions: core.OpOptions{
			Mode:         core.TransformMode,
			Intelligence: core.Fast,
		},
		SimilarityThreshold: 0.7,
	}
}

// Validate validates SimilarOptions
func (opts SimilarOptions) Validate() error {
	if opts.SimilarityThreshold < 0 || opts.SimilarityThreshold > 1 {
		return fmt.Errorf("similarity threshold must be between 0 and 1, got %f", opts.SimilarityThreshold)
	}
	return nil
}

// WithSimilarityThreshold sets the similarity threshold
func (opts SimilarOptions) WithSimilarityThreshold(threshold float64) SimilarOptions {
	opts.SimilarityThreshold = threshold
	return opts
}

// WithAspects sets specific aspects to compare
func (opts SimilarOptions) WithAspects(aspects []string) SimilarOptions {
	opts.Aspects = aspects
	return opts
}

// Similar checks semantic similarity with specialized options.
//
// Examples:
//
//	// Basic similarity check
//	similar, err := Similar("AI is great", "Artificial intelligence is wonderful",
//	    NewSimilarOptions())
//
//	// Custom threshold and aspects
//	similar, err := Similar(text1, text2, NewSimilarOptions().
