// package ops - Extended operations for data validation, formatting, and analysis
package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/monstercameron/SchemaFlow/core"
)

// ValidationResult contains the results of a validation operation
type ValidationResult struct {
	Valid       bool     `json:"valid"`
	Issues      []string `json:"issues"`
	Confidence  float64  `json:"confidence"`
	Suggestions []string `json:"suggestions"`
}

// Validate checks if data meets specified criteria using LLM interpretation
//
// Examples:
//
//	result, err := Validate(person, "age must be 18-100, email must be valid")
//	if !result.Valid {
//	    fmt.Printf("Validation issues: %v\n", result.Issues)
//	}
func Validate[T any](data T, rules string, opts ...core.OpOptions) (ValidationResult, error) {
	return validateImpl(core.GetDefaultClient(), data, rules, opts...)
}

// ClientValidate is the client-based version of Validate
func ClientValidate[T any](c *core.Client, data T, rules string, opts ...core.OpOptions) (ValidationResult, error) {
	// Pass client to implementation
	return validateImpl(c, data, rules, opts...)
}

func validateImpl[T any](c *core.Client, data T, rules string, opts ...core.OpOptions) (ValidationResult, error) {
	opt := core.ApplyDefaults(opts...)
	opt.Client = c // Set client from parameter

	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	// Convert data to JSON for validation
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("failed to marshal data: %w", err)
	}

	systemPrompt := `You are a data validation expert. Validate the provided data against the given rules.

Return a JSON object with:
{
  "valid": boolean,
  "issues": ["list of validation issues, empty if valid"],
  "confidence": 0.0-1.0,
  "suggestions": ["list of suggestions to fix issues, empty if valid"]
}`

	userPrompt := fmt.Sprintf(`Validate this data:
%s

Against these rules:
%s`, string(dataJSON), rules)

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("validation failed: %w", err)
	}

	var result ValidationResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// Try to parse as plain text if JSON parsing fails
		result.Valid = strings.Contains(strings.ToLower(response), "valid")
		result.Confidence = 0.5
		if !result.Valid {
			result.Issues = []string{response}
		}
	}

	return result, nil
}

// Format converts data to a specific output format using LLM interpretation
//
// Examples:
//
//	// Format as markdown table
//	formatted, err := Format(data, "markdown table with headers")
//
//	// Format as professional bio
//	bio, err := Format(person, "professional bio in third person")
func Format(data any, template string, opts ...core.OpOptions) (string, error) {
	return formatImpl(core.GetDefaultClient(), data, template, opts...)
}

// ClientFormat is the client-based version of Format
func ClientFormat(c *core.Client, data any, template string, opts ...core.OpOptions) (string, error) {
	return formatImpl(c, data, template, opts...)
}

func formatImpl(c *core.Client, data any, template string, opts ...core.OpOptions) (string, error) {
	opt := core.ApplyDefaults(opts...)
	opt.Client = c
	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	// Convert data to string representation
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	default:
		dataJSON, err := json.Marshal(data)
		if err != nil {
			dataStr = fmt.Sprintf("%v", data)
		} else {
			dataStr = string(dataJSON)
		}
	}

	systemPrompt := `You are a formatting expert. Convert the provided data into the requested format.
Follow the template instructions precisely and produce clean, well-formatted output.`

	userPrompt := fmt.Sprintf(`Format this data:
%s

Into this format:
%s`, dataStr, template)

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return "", fmt.Errorf("formatting failed: %w", err)
	}

	return strings.TrimSpace(response), nil
}

// Merge intelligently combines multiple data sources into a single result
//
// Examples:
//
//	// Merge customer records preferring newest data
//	merged, err := Merge([]Customer{dbRecord, apiResponse, csvRow}, "prefer newest")
func Merge[T any](sources []T, strategy string, opts ...core.OpOptions) (T, error) {
	return mergeImpl(core.GetDefaultClient(), sources, strategy, opts...)
}

// ClientMerge is the client-based version of Merge
func ClientMerge[T any](c *core.Client, sources []T, strategy string, opts ...core.OpOptions) (T, error) {
	return mergeImpl(c, sources, strategy, opts...)
}

func mergeImpl[T any](c *core.Client, sources []T, strategy string, opts ...core.OpOptions) (T, error) {
	var result T

	if len(sources) == 0 {
		return result, fmt.Errorf("no sources to merge")
	}

	if len(sources) == 1 {
		return sources[0], nil
	}

	opt := core.ApplyDefaults(opts...)
	opt.Client = c
	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	// Convert sources to JSON
	var sourcesJSON []string
	for i, source := range sources {
		sourceJSON, err := json.Marshal(source)
		if err != nil {
			return result, fmt.Errorf("failed to marshal source %d: %w", i, err)
		}
		sourcesJSON = append(sourcesJSON, string(sourceJSON))
	}

	// Get type information
	typeInfo := GenerateTypeSchema(reflect.TypeOf(result))

	systemPrompt := fmt.Sprintf(`You are a data merging expert. Merge multiple data sources into a single result.
Follow the merging strategy and produce a result matching this schema:
%s

Return only the merged JSON object.`, typeInfo)

	userPrompt := fmt.Sprintf(`Merge these sources:
%s

Using strategy: %s`, strings.Join(sourcesJSON, "\n"), strategy)

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return result, fmt.Errorf("merge failed: %w", err)
	}

	// Parse the merged result
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return result, fmt.Errorf("failed to parse merged result: %w", err)
	}

	return result, nil
}

// Question answers questions about data using LLM analysis
//
// Examples:
//
//	answer, err := Question(report, "What are the top 3 risks?")
//	summary, err := Question(data, "Summarize the key findings")
func Question(data any, question string, opts ...core.OpOptions) (string, error) {
	return questionImpl(core.GetDefaultClient(), data, question, opts...)
}

// ClientQuestion is the client-based version of Question
func ClientQuestion(c *core.Client, data any, question string, opts ...core.OpOptions) (string, error) {
	return questionImpl(c, data, question, opts...)
}

func questionImpl(c *core.Client, data any, question string, opts ...core.OpOptions) (string, error) {
	opt := core.ApplyDefaults(opts...)
	opt.Client = c
	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	// Convert data to string representation
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	default:
		dataJSON, err := json.Marshal(data)
		if err != nil {
			dataStr = fmt.Sprintf("%v", data)
		} else {
			dataStr = string(dataJSON)
		}
	}

	systemPrompt := `You are a data analysis expert. Answer questions about the provided data accurately and concisely.
Base your answers only on the information provided.`

	userPrompt := fmt.Sprintf(`Data:
%s

Question: %s`, dataStr, question)

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return "", fmt.Errorf("question answering failed: %w", err)
	}

	return strings.TrimSpace(response), nil
}

// DeduplicateResult contains the results of deduplication
type DeduplicateResult[T any] struct {
	Unique       []T
	Duplicates   [][]T // Groups of duplicates
	TotalRemoved int
}

// Deduplicate removes duplicates using semantic similarity
//
// Examples:
//
//	result, err := Deduplicate(customers, 0.85) // 85% similarity threshold
//	fmt.Printf("Removed %d duplicates\n", result.TotalRemoved)
func Deduplicate[T any](items []T, threshold float64, opts ...core.OpOptions) (DeduplicateResult[T], error) {
	return deduplicateImpl(core.GetDefaultClient(), items, threshold, opts...)
}

// ClientDeduplicate is the client-based version of Deduplicate
func ClientDeduplicate[T any](c *core.Client, items []T, threshold float64, opts ...core.OpOptions) (DeduplicateResult[T], error) {
	return deduplicateImpl(c, items, threshold, opts...)
}

func deduplicateImpl[T any](c *core.Client, items []T, threshold float64, opts ...core.OpOptions) (DeduplicateResult[T], error) {
	result := DeduplicateResult[T]{
		Unique:     []T{},
		Duplicates: [][]T{},
	}

	if len(items) <= 1 {
		result.Unique = items
		return result, nil
	}

	opt := core.ApplyDefaults(opts...)
	opt.Client = c
	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	// Convert items to JSON for comparison
	itemsJSON := make([]string, len(items))
	for i, item := range items {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return result, fmt.Errorf("failed to marshal item %d: %w", i, err)
		}
		itemsJSON[i] = string(itemJSON)
	}

	systemPrompt := fmt.Sprintf(`You are a deduplication expert. Identify duplicate items based on semantic similarity.
Items with similarity >= %.2f should be considered duplicates.

Return a JSON object with:
{
  "groups": [
    [0, 5, 8],  // indices of items that are duplicates of each other
    [2, 7],     // another group of duplicates
    [1],        // unique item
    [3],        // unique item
    ...
  ]
}`, threshold)

	userPrompt := fmt.Sprintf(`Find duplicates in these items:
%s`, strings.Join(itemsJSON, "\n"))

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return result, fmt.Errorf("deduplication failed: %w", err)
	}

	// Parse the grouping response
	var grouping struct {
		Groups [][]int `json:"groups"`
	}

	if err := json.Unmarshal([]byte(response), &grouping); err != nil {
		// Fallback: treat all items as unique if parsing fails
		result.Unique = items
		return result, nil
	}

	// Process groups
	seen := make(map[int]bool)
	for _, group := range grouping.Groups {
		if len(group) == 0 {
			continue
		}

		// Mark all indices as seen
		for _, idx := range group {
			if idx >= 0 && idx < len(items) {
				seen[idx] = true
			}
		}

		if len(group) == 1 {
			// Unique item
			if idx := group[0]; idx >= 0 && idx < len(items) {
				result.Unique = append(result.Unique, items[idx])
			}
		} else {
			// Group of duplicates - keep first, track others
			if idx := group[0]; idx >= 0 && idx < len(items) {
				result.Unique = append(result.Unique, items[idx])
			}

			// Track the duplicate group
			var dupGroup []T
			for _, idx := range group {
				if idx >= 0 && idx < len(items) {
					dupGroup = append(dupGroup, items[idx])
				}
			}
			if len(dupGroup) > 1 {
				result.Duplicates = append(result.Duplicates, dupGroup)
				result.TotalRemoved += len(dupGroup) - 1
			}
		}
	}

	// Add any items not mentioned in groups as unique
	for i, item := range items {
		if !seen[i] {
			result.Unique = append(result.Unique, item)
		}
	}

	return result, nil
}
