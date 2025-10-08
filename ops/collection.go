package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Choose selects the best option from a list with specialized options.
//
// Examples:
//
//	// Basic selection
//	best, err := Choose(options, NewChooseOptions().
//	    WithCriteria([]string{"quality", "price"}))
//
//	// Selection with reasoning
//	best, err := Choose(candidates, NewChooseOptions().
//	    WithRequireReasoning(true).
//	    WithTopN(3))
func Choose[T any](options []T, opts ChooseOptions) (T, error) {
	var result T
	
	// Validate options
	if err := opts.Validate(); err != nil {
		return result, fmt.Errorf("invalid options: %w", err)
	}
	
	if len(options) == 0 {
		return result, ChooseError{
			Options: []any{},
			Reason:  "no options provided",
		}
	}
	
	if len(options) == 1 {
		return options[0], nil
	}
	
	opOptions := opts.toOpOptions()
	
	// Build selection instructions
	var instructions []string
	
	if len(opts.Criteria) > 0 {
		instructions = append(instructions, fmt.Sprintf("Selection criteria: %s", strings.Join(opts.Criteria, ", ")))
	}
	
	if opts.RequireReasoning {
		instructions = append(instructions, "Provide reasoning for your choice")
	}
	
	if opts.TopN > 1 {
		instructions = append(instructions, fmt.Sprintf("Return top %d options", opts.TopN))
	}
	
	if opts.Strategy != "" {
		instructions = append(instructions, fmt.Sprintf("Use %s strategy", opts.Strategy))
	}
	
	if len(instructions) > 0 {
		steering := strings.Join(instructions, ". ")
		if opts.Steering != "" {
			steering = opts.Steering + ". " + steering
		}
		opOptions.Steering = steering
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return result, ChooseError{
			Options: interfaceSlice(options),
			Reason:  fmt.Sprintf("failed to marshal options: %v", err),
		}
	}
	
	systemPrompt := `You are a selection expert. Choose the best option from the provided list.

Rules:
- Evaluate all options carefully
- Select the most appropriate one based on the criteria
- Return ONLY the index number (0-based) of your choice
- Return a single integer, nothing else`
	
	userPrompt := fmt.Sprintf("Choose the best option from this list:\n%s", string(optionsJSON))
	
	response, err := callLLM(ctx, systemPrompt, userPrompt, opOptions)
	if err != nil {
		return result, ChooseError{
			Options: interfaceSlice(options),
			Reason:  err.Error(),
		}
	}
	
	response = strings.TrimSpace(response)
	var index int
	if _, err := fmt.Sscanf(response, "%d", &index); err != nil {
		return result, ChooseError{
			Options: interfaceSlice(options),
			Reason:  fmt.Sprintf("failed to parse index: %v", err),
		}
	}
	
	if index < 0 || index >= len(options) {
		return result, ChooseError{
			Options: interfaceSlice(options),
			Reason:  fmt.Sprintf("index %d out of range", index),
		}
	}
	
	return options[index], nil
}

// Filter semantically filters items with specialized options.
//
// Examples:
//
//	// Basic filtering
//	filtered, err := Filter(items, NewFilterOptions().
//	    WithCriteria("items with positive sentiment"))
//
//	// Complex filtering with confidence
//	filtered, err := Filter(products, NewFilterOptions().
//	    WithCriteria("electronics under $500").
//	    WithMinConfidence(0.8).
//	    WithIncludeReasons(true))
func Filter[T any](items []T, opts FilterOptions) ([]T, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}
	
	if len(items) == 0 {
		return items, nil
	}
	
	opOptions := opts.toOpOptions()
	
	// Build filter instructions
	var instructions []string
	
	instructions = append(instructions, fmt.Sprintf("Filter criteria: %s", opts.Criteria))
	
	if opts.KeepMatching {
		instructions = append(instructions, "Keep items that match the criteria")
	} else {
		instructions = append(instructions, "Remove items that match the criteria")
	}
	
	if opts.MinConfidence > 0 {
		instructions = append(instructions, fmt.Sprintf("Minimum confidence: %.2f", opts.MinConfidence))
	}
	
	if opts.IncludeReasons {
		instructions = append(instructions, "Include reasons for each decision")
	}
	
	steering := strings.Join(instructions, ". ")
	if opts.Steering != "" {
		steering = opts.Steering + ". " + steering
	}
	opOptions.Steering = steering
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, FilterError{
			Items:  interfaceSlice(items),
			Reason: fmt.Sprintf("failed to marshal items: %v", err),
		}
	}
	
	systemPrompt := `You are a filtering expert. Filter items based on the specified criteria.

Rules:
- Evaluate each item against the criteria
- Include items that match the criteria
- Return a JSON array of indices (0-based) of items to keep
- Return ONLY the JSON array, nothing else`
	
	userPrompt := fmt.Sprintf("Filter these items:\n%s", string(itemsJSON))

	response, err := callLLM(ctx, systemPrompt, userPrompt, opOptions)
	if err != nil {
		return nil, FilterError{
			Items:  interfaceSlice(items),
			Reason: err.Error(),
		}
	}
	
	var indices []int
	if err := parseJSON(response, &indices); err != nil {
		return nil, FilterError{
			Items:  interfaceSlice(items),
			Reason: fmt.Sprintf("failed to parse indices: %v", err),
		}
	}
	
	result := make([]T, 0, len(indices))
	for _, idx := range indices {
		if idx >= 0 && idx < len(items) {
			result = append(result, items[idx])
		}
	}
	
	return result, nil
}

// Sort orders items semantically with specialized options.
//
// Examples:
//
//	// Basic sorting
//	sorted, err := Sort(items, NewSortOptions().
//	    WithCriteria("by importance"))
//
//	// Multi-criteria sorting
//	sorted, err := Sort(products, NewSortOptions().
//	    WithCriteria("by quality").
//	    WithSecondaryCriteria([]string{"by price", "by popularity"}).
//	    WithDirection("descending"))
func Sort[T any](items []T, opts SortOptions) ([]T, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}
	
	if len(items) <= 1 {
		return items, nil
	}
	
	opOptions := opts.toOpOptions()
	
	// Build sort instructions
	var instructions []string
	
	instructions = append(instructions, fmt.Sprintf("Sort criteria: %s", opts.Criteria))
	
	if opts.Direction != "" {
		instructions = append(instructions, fmt.Sprintf("Direction: %s", opts.Direction))
	}
	
	if len(opts.SecondaryCriteria) > 0 {
		instructions = append(instructions, fmt.Sprintf("Secondary criteria: %s", strings.Join(opts.SecondaryCriteria, ", ")))
	}
	
	if opts.ComparisonLogic != "" {
		instructions = append(instructions, fmt.Sprintf("Comparison logic: %s", opts.ComparisonLogic))
	}
	
	if opts.Stable {
		instructions = append(instructions, "Maintain relative order of equal elements")
	}
	
	if opts.IncludeScores {
		instructions = append(instructions, "Include sort scores")
	}
	
	steering := strings.Join(instructions, ". ")
	if opts.Steering != "" {
		steering = opts.Steering + ". " + steering
	}
	opOptions.Steering = steering
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, SortError{
			Items:  interfaceSlice(items),
			Reason: fmt.Sprintf("failed to marshal items: %v", err),
		}
	}
	
	systemPrompt := `You are a sorting expert. Sort items based on the specified criteria.

Rules:
- Evaluate all items according to the sorting criteria
- Arrange them in the proper order
- Return a JSON array of indices (0-based) representing the sorted order
- Return ONLY the JSON array, nothing else`
	
	userPrompt := fmt.Sprintf("Sort these items:\n%s", string(itemsJSON))
	
	response, err := callLLM(ctx, systemPrompt, userPrompt, opOptions)
	if err != nil {
		return nil, SortError{
			Items:  interfaceSlice(items),
			Reason: err.Error(),
		}
	}
	
	var indices []int
	if err := parseJSON(response, &indices); err != nil {
		return nil, SortError{
			Items:  interfaceSlice(items),
			Reason: fmt.Sprintf("failed to parse indices: %v", err),
		}
	}
	
	if len(indices) != len(items) {
		return nil, SortError{
			Items:  interfaceSlice(items),
			Reason: fmt.Sprintf("received %d indices for %d items", len(indices), len(items)),
		}
	}
	
	result := make([]T, len(items))
	for i, idx := range indices {
		if idx < 0 || idx >= len(items) {
			return nil, SortError{
				Items:  interfaceSlice(items),
				Reason: fmt.Sprintf("index %d out of range", idx),
			}
		}
		result[i] = items[idx]
	}
	
	return result, nil
}

func interfaceSlice[T any](items []T) []any {
	result := make([]any, len(items))
	for i, item := range items {
		result[i] = item
	}
	return result
}
