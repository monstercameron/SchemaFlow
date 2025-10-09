// package ops - Diff operation for intelligent difference detection
package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/monstercameron/SchemaFlow/core"
)

// DiffResult contains the results of a difference analysis
type DiffResult struct {
	Added    []string     `json:"added"`    // New fields/values that appeared
	Removed  []string     `json:"removed"`  // Fields/values that were removed
	Modified []DiffChange `json:"modified"` // Fields that changed with details
	Summary  string       `json:"summary"`  // LLM-generated explanation of changes
}

// DiffChange represents a single field modification
type DiffChange struct {
	Field      string `json:"field"`       // Field name that changed
	OldValue   any    `json:"old_value"`   // Previous value
	NewValue   any    `json:"new_value"`   // New value
	ChangeType string `json:"change_type"` // "modified", "type_changed", "structure_changed"
}

// DiffOptions configures the Diff operation
type DiffOptions struct {
	core.OpOptions
	Context      string   // Additional context about the data
	IgnoreFields []string // Fields to skip in comparison
	DeepCompare  bool     // Compare nested structures recursively
}

// NewDiffOptions creates DiffOptions with defaults
func NewDiffOptions() DiffOptions {
	return DiffOptions{
		OpOptions: core.OpOptions{
			Mode:         core.TransformMode,
			Intelligence: core.Fast,
		},
		DeepCompare: true,
	}
}

// WithContext sets additional context for difference analysis
func (opts DiffOptions) WithContext(context string) DiffOptions {
	opts.Context = context
	return opts
}

// WithIgnoreFields sets fields to skip in comparison
func (opts DiffOptions) WithIgnoreFields(fields []string) DiffOptions {
	opts.IgnoreFields = fields
	return opts
}

// WithDeepCompare enables/disables recursive struct comparison
func (opts DiffOptions) WithDeepCompare(deep bool) DiffOptions {
	opts.DeepCompare = deep
	return opts
}

// WithIntelligence sets the intelligence level
func (opts DiffOptions) WithIntelligence(intelligence core.Speed) DiffOptions {
	opts.OpOptions.Intelligence = intelligence
	return opts
}

// Validate validates DiffOptions
func (opts DiffOptions) Validate() error {
	return nil // No specific validation needed
}

// toOpOptions converts DiffOptions to core.OpOptions
func (opts DiffOptions) toOpOptions() core.OpOptions {
	return opts.OpOptions
}

// Diff intelligently compares two data instances and explains differences
//
// Examples:
//
//	// Compare customer records
//	result, err := Diff[Customer](oldCustomer, newCustomer,
//	    NewDiffOptions().WithContext("Customer management system"))
//
//	// Compare with ignored fields
//	result, err := Diff[Product](oldProduct, newProduct,
//	    NewDiffOptions().WithIgnoreFields([]string{"last_updated"}))
func Diff[T any](oldData, newData T, opts DiffOptions) (DiffResult, error) {
	return diffImpl(core.GetDefaultClient(), oldData, newData, opts)
}

// ClientDiff is the client-based version of Diff
func ClientDiff[T any](c *core.Client, oldData, newData T, opts DiffOptions) (DiffResult, error) {
	return diffImpl(c, oldData, newData, opts)
}

func diffImpl[T any](c *core.Client, oldData, newData T, opts DiffOptions) (DiffResult, error) {
	logger := core.GetLogger()
	logger.Debug("Starting diff operation", "requestID", opts.RequestID)

	result := DiffResult{}

	// Validate options
	if err := opts.Validate(); err != nil {
		logger.Error("Diff operation validation failed", "requestID", opts.RequestID, "error", err)
		return result, fmt.Errorf("invalid options: %w", err)
	}

	// Early exit for identical data (optimization)
	if reflect.DeepEqual(oldData, newData) {
		result.Summary = "No changes detected - data is identical"
		logger.Debug("Diff operation completed: no changes", "requestID", opts.RequestID)
		return result, nil
	}

	// Perform structural comparison
	changes, err := compareData(oldData, newData, opts)
	if err != nil {
		logger.Error("Diff operation comparison failed", "requestID", opts.RequestID, "error", err)
		return result, fmt.Errorf("comparison failed: %w", err)
	}

	result.Added = changes.Added
	result.Removed = changes.Removed
	result.Modified = changes.Modified

	// Generate intelligent summary if intelligence level requires it and there are changes
	hasChanges := len(result.Added) > 0 || len(result.Removed) > 0 || len(result.Modified) > 0
	if (opts.Intelligence != core.Fast || len(result.Modified) > 0) && hasChanges {
		summary, err := generateDiffSummary(c, oldData, newData, changes, opts)
		if err != nil {
			// Don't fail the whole operation if summary fails
			result.Summary = "Summary generation failed, but changes detected successfully"
			logger.Warn("Diff operation summary generation failed", "requestID", opts.RequestID, "error", err)
		} else {
			result.Summary = summary
		}
	} else if !hasChanges {
		result.Summary = "No changes detected between the data instances"
	}

	logger.Debug("Diff operation succeeded", "requestID", opts.RequestID, "added", len(result.Added), "removed", len(result.Removed), "modified", len(result.Modified))

	return result, nil
}

// comparisonResult holds intermediate comparison data
type comparisonResult struct {
	Added    []string
	Removed  []string
	Modified []DiffChange
}

// compareData performs the actual data comparison using reflection
func compareData(oldData, newData any, opts DiffOptions) (comparisonResult, error) {
	result := comparisonResult{}

	oldValue := reflect.ValueOf(oldData)
	newValue := reflect.ValueOf(newData)

	// Handle pointer types
	if oldValue.Kind() == reflect.Ptr {
		if oldValue.IsNil() && newValue.IsNil() {
			return result, nil
		}
		if oldValue.IsNil() || newValue.IsNil() {
			return result, fmt.Errorf("cannot compare nil and non-nil pointers")
		}
		oldValue = oldValue.Elem()
		newValue = newValue.Elem()
	}

	// Ensure both are structs
	if oldValue.Kind() != reflect.Struct || newValue.Kind() != reflect.Struct {
		return result, fmt.Errorf("diff operation requires struct types")
	}

	oldType := oldValue.Type()
	newType := newValue.Type()

	// Create ignored fields map for O(1) lookup (optimization)
	ignoredFields := make(map[string]bool, len(opts.IgnoreFields))
	for _, field := range opts.IgnoreFields {
		ignoredFields[field] = true
	}

	// Pre-allocate slices with estimated capacity (optimization)
	numFields := oldValue.NumField()
	if newValue.NumField() > numFields {
		numFields = newValue.NumField()
	}
	result.Added = make([]string, 0, numFields/4)        // Estimate 25% added
	result.Removed = make([]string, 0, numFields/4)      // Estimate 25% removed
	result.Modified = make([]DiffChange, 0, numFields/2) // Estimate 50% modified

	// Compare each field in old struct
	for i := 0; i < oldValue.NumField(); i++ {
		field := oldType.Field(i)
		fieldName := field.Name

		// Skip ignored fields (O(1) lookup)
		if ignoredFields[fieldName] {
			continue
		}

		oldFieldValue := oldValue.Field(i)
		newFieldValue := newValue.FieldByName(fieldName)

		// Check if field exists in new struct
		if !newFieldValue.IsValid() {
			// Field removed in new struct
			result.Removed = append(result.Removed, fieldName)
			continue
		}

		// Compare field values
		change := compareField(fieldName, oldFieldValue, newFieldValue, opts.DeepCompare)
		if change != nil {
			result.Modified = append(result.Modified, *change)
		}
	}

	// Check for added fields in new struct
	for i := 0; i < newValue.NumField(); i++ {
		field := newType.Field(i)
		fieldName := field.Name

		// Skip ignored fields (O(1) lookup)
		if ignoredFields[fieldName] {
			continue
		}

		// Check if field exists in old struct
		if !oldValue.FieldByName(fieldName).IsValid() {
			result.Added = append(result.Added, fieldName)
		}
	}

	return result, nil
}

// compareField compares individual field values
func compareField(fieldName string, oldValue, newValue reflect.Value, deepCompare bool) *DiffChange {
	// Handle nil values
	if !oldValue.IsValid() && !newValue.IsValid() {
		return nil
	}
	if !oldValue.IsValid() || !newValue.IsValid() {
		return &DiffChange{
			Field:      fieldName,
			OldValue:   nil,
			NewValue:   nil,
			ChangeType: "structure_changed",
		}
	}

	// Handle pointer types
	if oldValue.Kind() == reflect.Ptr && newValue.Kind() == reflect.Ptr {
		if oldValue.IsNil() && newValue.IsNil() {
			return nil
		}
		if oldValue.IsNil() || newValue.IsNil() {
			return &DiffChange{
				Field:      fieldName,
				OldValue:   formatValue(oldValue),
				NewValue:   formatValue(newValue),
				ChangeType: "modified",
			}
		}
		oldValue = oldValue.Elem()
		newValue = newValue.Elem()
	}

	// Handle different types
	if oldValue.Kind() != newValue.Kind() {
		return &DiffChange{
			Field:      fieldName,
			OldValue:   formatValue(oldValue),
			NewValue:   formatValue(newValue),
			ChangeType: "type_changed",
		}
	}

	// Compare based on type
	switch oldValue.Kind() {
	case reflect.Struct:
		if deepCompare {
			// For structs, compare if they're different (simplified check)
			if !reflect.DeepEqual(oldValue.Interface(), newValue.Interface()) {
				return &DiffChange{
					Field:      fieldName,
					OldValue:   formatValue(oldValue),
					NewValue:   formatValue(newValue),
					ChangeType: "structure_changed",
				}
			}
		}
	case reflect.Slice, reflect.Array:
		if !reflect.DeepEqual(oldValue.Interface(), newValue.Interface()) {
			return &DiffChange{
				Field:      fieldName,
				OldValue:   formatValue(oldValue),
				NewValue:   formatValue(newValue),
				ChangeType: "modified",
			}
		}
	case reflect.Map:
		if !reflect.DeepEqual(oldValue.Interface(), newValue.Interface()) {
			return &DiffChange{
				Field:      fieldName,
				OldValue:   formatValue(oldValue),
				NewValue:   formatValue(newValue),
				ChangeType: "modified",
			}
		}
	default:
		// Primitive types: direct comparison
		if !reflect.DeepEqual(oldValue.Interface(), newValue.Interface()) {
			return &DiffChange{
				Field:      fieldName,
				OldValue:   oldValue.Interface(),
				NewValue:   newValue.Interface(),
				ChangeType: "modified",
			}
		}
	}

	return nil
}

// formatValue formats a reflect.Value for display
func formatValue(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	return v.Interface()
}

// generateDiffSummary uses LLM to create an intelligent summary of changes
func generateDiffSummary(c *core.Client, oldData, newData any, changes comparisonResult, opts DiffOptions) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), core.GetTimeout())
	defer cancel()

	// Marshal data for prompt (only when needed)
	oldJSON, err := json.MarshalIndent(oldData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal old data: %w", err)
	}

	newJSON, err := json.MarshalIndent(newData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal new data: %w", err)
	}

	changesJSON, err := json.MarshalIndent(changes, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal changes: %w", err)
	}

	// Build system prompt
	systemPrompt := `You are an expert data analyst specializing in change detection and impact analysis.
Given two versions of data and the detected changes, provide a concise but insightful summary explaining:
1. What the changes represent in practical terms
2. Any patterns or significance of the changes
3. Potential implications or business impact
Keep the summary under 200 words and focus on actionable insights.`

	// Build user prompt using strings.Builder for efficiency
	var promptBuilder strings.Builder
	promptBuilder.Grow(1024) // Pre-allocate capacity
	promptBuilder.WriteString("Analyze these data changes:\n\nOLD DATA:\n")
	promptBuilder.Write(oldJSON)
	promptBuilder.WriteString("\n\nNEW DATA:\n")
	promptBuilder.Write(newJSON)
	promptBuilder.WriteString("\n\nDETECTED CHANGES:\n")
	promptBuilder.Write(changesJSON)

	if opts.Context != "" {
		promptBuilder.WriteString("\n\nCONTEXT: ")
		promptBuilder.WriteString(opts.Context)
	}

	userPrompt := promptBuilder.String()

	// Call LLM for summary
	opt := opts.toOpOptions()
	opt.Client = c

	response, err := core.CallLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return "", fmt.Errorf("summary generation failed: %w", err)
	}

	return strings.TrimSpace(response), nil
}
