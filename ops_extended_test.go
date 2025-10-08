package schemaflow

import (
	"context"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	setupMockClient()
	
	// Update mock to return validation results
	callLLM = func(ctx context.Context, system, user string, opts OpOptions) (string, error) {
		if strings.Contains(system, "validation expert") {
			if strings.Contains(user, "age must be 18-100") {
				return `{
					"valid": true,
					"issues": [],
					"confidence": 0.95,
					"suggestions": []
				}`, nil
			}
			return `{
				"valid": false,
				"issues": ["Age is outside valid range"],
				"confidence": 0.9,
				"suggestions": ["Set age between 18 and 100"]
			}`, nil
		}
		return mockLLMResponse(ctx, system, user, opts)
	}
	
	t.Run("ValidData", func(t *testing.T) {
		person := Person{Name: "John", Age: 30, Email: "john@example.com"}
		result, err := Validate(person, "age must be 18-100, email must be valid")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !result.Valid {
			t.Errorf("Expected valid=true, got false")
		}
		
		if len(result.Issues) > 0 {
			t.Errorf("Expected no issues, got %v", result.Issues)
		}
		
		if result.Confidence < 0.9 {
			t.Errorf("Expected high confidence, got %.2f", result.Confidence)
		}
	})
	
	t.Run("InvalidData", func(t *testing.T) {
		person := Person{Name: "Jane", Age: 150, Email: "invalid"}
		result, err := Validate(person, "age must be reasonable")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if result.Valid {
			t.Errorf("Expected valid=false, got true")
		}
		
		if len(result.Issues) == 0 {
			t.Errorf("Expected validation issues, got none")
		}
		
		if len(result.Suggestions) == 0 {
			t.Errorf("Expected suggestions, got none")
		}
	})
}

func TestFormat(t *testing.T) {
	setupMockClient()
	
	// Update mock for formatting
	callLLM = func(ctx context.Context, system, user string, opts OpOptions) (string, error) {
		if strings.Contains(system, "formatting expert") {
			if strings.Contains(user, "markdown table") {
				return "| Name | Age | Email |\n|------|-----|-------|\n| John | 30 | john@example.com |", nil
			}
			if strings.Contains(user, "professional bio") {
				return "John is a 30-year-old professional with extensive experience.", nil
			}
		}
		return mockLLMResponse(ctx, system, user, opts)
	}
	
	t.Run("FormatAsTable", func(t *testing.T) {
		person := Person{Name: "John", Age: 30, Email: "john@example.com"}
		formatted, err := Format(person, "markdown table with headers")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.Contains(formatted, "|") {
			t.Errorf("Expected markdown table format, got: %s", formatted)
		}
		
		if !strings.Contains(formatted, "Name") || !strings.Contains(formatted, "Age") {
			t.Errorf("Expected headers in table, got: %s", formatted)
		}
	})
	
	t.Run("FormatAsBio", func(t *testing.T) {
		person := Person{Name: "John", Age: 30, Email: "john@example.com"}
		bio, err := Format(person, "professional bio in third person")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.Contains(bio, "John") {
			t.Errorf("Expected name in bio, got: %s", bio)
		}
		
		if strings.Contains(bio, "|") {
			t.Errorf("Bio should not be in table format, got: %s", bio)
		}
	})
}

func TestMerge(t *testing.T) {
	setupMockClient()
	
	// Update mock for merging
	callLLM = func(ctx context.Context, system, user string, opts OpOptions) (string, error) {
		if strings.Contains(system, "merging expert") {
			// Return a merged person with combined data
			return `{"name": "John Doe", "age": 30, "email": "john.doe@example.com"}`, nil
		}
		return mockLLMResponse(ctx, system, user, opts)
	}
	
	t.Run("MergeMultipleSources", func(t *testing.T) {
		sources := []Person{
			{Name: "John", Age: 30, Email: "john@old.com"},
			{Name: "John Doe", Age: 30, Email: "john@new.com"},
			{Name: "John Doe", Age: 30, Email: "john.doe@example.com"},
		}
		
		merged, err := Merge(sources, "prefer newest, combine names")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if merged.Name == "" {
			t.Errorf("Expected merged name, got empty")
		}
		
		if merged.Email == "" {
			t.Errorf("Expected merged email, got empty")
		}
	})
	
	t.Run("MergeSingleSource", func(t *testing.T) {
		sources := []Person{
			{Name: "Jane", Age: 25, Email: "jane@example.com"},
		}
		
		merged, err := Merge(sources, "any strategy")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		// Should return the single source unchanged
		if merged.Name != "Jane" || merged.Age != 25 {
			t.Errorf("Single source should be returned unchanged, got %+v", merged)
		}
	})
	
	t.Run("MergeEmptySources", func(t *testing.T) {
		sources := []Person{}
		
		_, err := Merge(sources, "any strategy")
		
		if err == nil {
			t.Errorf("Expected error for empty sources")
		}
	})
}

func TestQuestion(t *testing.T) {
	setupMockClient()
	
	// Update mock for questions
	callLLM = func(ctx context.Context, system, user string, opts OpOptions) (string, error) {
		if strings.Contains(system, "data analysis expert") {
			if strings.Contains(user, "What is the average age") {
				return "The average age is 30 years.", nil
			}
			if strings.Contains(user, "How many people") {
				return "There are 3 people in the data.", nil
			}
		}
		return mockLLMResponse(ctx, system, user, opts)
	}
	
	t.Run("QuestionAboutData", func(t *testing.T) {
		data := []Person{
			{Name: "John", Age: 30},
			{Name: "Jane", Age: 25},
			{Name: "Bob", Age: 35},
		}
		
		answer, err := Question(data, "What is the average age?")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.Contains(answer, "30") || !strings.Contains(answer, "average") {
			t.Errorf("Expected answer about average age, got: %s", answer)
		}
	})
	
	t.Run("CountQuestion", func(t *testing.T) {
		data := []Person{
			{Name: "John", Age: 30},
			{Name: "Jane", Age: 25},
			{Name: "Bob", Age: 35},
		}
		
		answer, err := Question(data, "How many people are there?")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.Contains(answer, "3") || !strings.Contains(answer, "people") {
			t.Errorf("Expected answer about count, got: %s", answer)
		}
	})
}

func TestDeduplicate(t *testing.T) {
	setupMockClient()
	
	// Update mock for deduplication
	callLLM = func(ctx context.Context, system, user string, opts OpOptions) (string, error) {
		if strings.Contains(system, "deduplication expert") {
			// Return groups where items 0 and 2 are duplicates
			return `{
				"groups": [
					[0, 2],
					[1],
					[3]
				]
			}`, nil
		}
		return mockLLMResponse(ctx, system, user, opts)
	}
	
	t.Run("DeduplicateWithDuplicates", func(t *testing.T) {
		items := []Person{
			{Name: "John Doe", Age: 30, Email: "john@example.com"},
			{Name: "Jane Smith", Age: 25, Email: "jane@example.com"},
			{Name: "John D.", Age: 30, Email: "john.d@example.com"}, // Duplicate of first
			{Name: "Bob Johnson", Age: 35, Email: "bob@example.com"},
		}
		
		result, err := Deduplicate(items, 0.85)
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		// Should have 3 unique items (removed 1 duplicate)
		if len(result.Unique) != 3 {
			t.Errorf("Expected 3 unique items, got %d", len(result.Unique))
		}
		
		// Should have 1 removed
		if result.TotalRemoved != 1 {
			t.Errorf("Expected 1 removed, got %d", result.TotalRemoved)
		}
		
		// Should have 1 duplicate group
		if len(result.Duplicates) != 1 {
			t.Errorf("Expected 1 duplicate group, got %d", len(result.Duplicates))
		}
	})
	
	t.Run("DeduplicateNoDuplicates", func(t *testing.T) {
		// Update mock for no duplicates case
		callLLM = func(ctx context.Context, system, user string, opts OpOptions) (string, error) {
			if strings.Contains(system, "deduplication expert") {
				return `{
					"groups": [
						[0],
						[1]
					]
				}`, nil
			}
			return mockLLMResponse(ctx, system, user, opts)
		}
		
		items := []Person{
			{Name: "Alice", Age: 20, Email: "alice@example.com"},
			{Name: "Bob", Age: 30, Email: "bob@example.com"},
		}
		
		result, err := Deduplicate(items, 0.85)
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		// All items should be unique
		if len(result.Unique) != 2 {
			t.Errorf("Expected 2 unique items, got %d", len(result.Unique))
		}
		
		if result.TotalRemoved != 0 {
			t.Errorf("Expected 0 removed, got %d", result.TotalRemoved)
		}
	})
	
	t.Run("DeduplicateEmptyList", func(t *testing.T) {
		items := []Person{}
		
		result, err := Deduplicate(items, 0.85)
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if len(result.Unique) != 0 {
			t.Errorf("Expected 0 unique items for empty input, got %d", len(result.Unique))
		}
	})
}

// Test with client
func TestExtendedOperationsWithClient(t *testing.T) {
	setupMockClient()
	
	client := NewClient("")
	
	t.Run("ClientValidate", func(t *testing.T) {
		person := Person{Name: "Test", Age: 25}
		result, err := ClientValidate(client, person, "age > 0")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		// Just check it runs without error
		_ = result
	})
	
	t.Run("ClientFormat", func(t *testing.T) {
		data := "test data"
		formatted, err := ClientFormat(client, data, "uppercase")
		
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		// Just check it runs without error
		_ = formatted
	})
}