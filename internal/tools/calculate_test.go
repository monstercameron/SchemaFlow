package tools

import (
	"context"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
		hasError bool
	}{
		// Basic arithmetic
		{"2 + 2", 4, false},
		{"10 - 3", 7, false},
		{"5 * 6", 30, false},
		{"20 / 4", 5, false},
		{"10 % 3", 1, false},

		// Decimals
		{"3.14 + 2.86", 6, false},
		{"10.5 / 2", 5.25, false},

		// Parentheses
		{"(2 + 3) * 4", 20, false},
		{"2 + (3 * 4)", 14, false},
		{"((2 + 3) * 4) / 2", 10, false},

		// Negative numbers
		{"-5 + 3", -2, false},
		{"5 + -3", 2, false},

		// Percentages
		{"15% of 200", 30, false},
		{"20% of 50", 10, false},
		{"50%", 0.5, false},
		{"100%", 1, false},

		// Functions
		{"sqrt(16)", 4, false},
		{"sqrt(2)", math.Sqrt(2), false},
		{"abs(-5)", 5, false},
		{"floor(3.7)", 3, false},
		{"ceil(3.2)", 4, false},
		{"round(3.5)", 4, false},
		{"pow(2, 3)", 8, false},
		{"pow(10, 2)", 100, false},

		// Complex expressions
		{"sqrt(16) + pow(2, 3)", 12, false},
		{"2 * sqrt(25) + 5", 15, false},

		// Errors
		{"10 / 0", 0, true},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := Calculate(tt.expr)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for %q", tt.expr)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %q: %v", tt.expr, err)
				return
			}

			// Use approximate comparison for floating point
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("Calculate(%q) = %v, expected %v", tt.expr, result, tt.expected)
			}
		})
	}
}

func TestCalculateTool(t *testing.T) {
	result, err := CalculateTool.Execute(context.Background(), map[string]any{
		"expression": "2 + 2",
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !result.Success {
		t.Errorf("Expected success, got error: %s", result.Error)
	}
	if result.Data.(float64) != 4 {
		t.Errorf("Expected 4, got %v", result.Data)
	}

	// Test with invalid expression
	result, _ = CalculateTool.Execute(context.Background(), map[string]any{
		"expression": "invalid",
	})
	if result.Success {
		t.Error("Expected failure for invalid expression")
	}

	// Test with missing expression
	result, _ = CalculateTool.Execute(context.Background(), map[string]any{})
	if result.Success {
		t.Error("Expected failure for missing expression")
	}
}

func TestCalculateStats(t *testing.T) {
	numbers := []float64{1, 2, 3, 4, 5}
	stats := CalculateStats(numbers)

	if stats.Count != 5 {
		t.Errorf("Count = %d, expected 5", stats.Count)
	}
	if stats.Sum != 15 {
		t.Errorf("Sum = %f, expected 15", stats.Sum)
	}
	if stats.Mean != 3 {
		t.Errorf("Mean = %f, expected 3", stats.Mean)
	}
	if stats.Min != 1 {
		t.Errorf("Min = %f, expected 1", stats.Min)
	}
	if stats.Max != 5 {
		t.Errorf("Max = %f, expected 5", stats.Max)
	}

	// Test empty slice
	emptyStats := CalculateStats([]float64{})
	if emptyStats.Count != 0 {
		t.Error("Expected 0 count for empty slice")
	}
}

func TestCalculateToolRegistered(t *testing.T) {
	// Use a fresh registry to test registration
	registry := NewRegistry()
	err := registry.Register(CalculateTool)
	if err != nil {
		t.Fatalf("Failed to register tool: %v", err)
	}

	tool, ok := registry.Get("calculate")
	if !ok {
		t.Fatal("Calculate tool not registered")
	}
	if tool.Category != CategoryComputation {
		t.Errorf("Expected category %s, got %s", CategoryComputation, tool.Category)
	}
}
