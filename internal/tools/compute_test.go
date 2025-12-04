package tools

import (
	"context"
	"math"
	"testing"
)

func TestRegexMatch(t *testing.T) {
	tests := []struct {
		pattern  string
		text     string
		expected bool
	}{
		{`\d+`, "abc123def", true},
		{`\d+`, "abcdef", false},
		{`^hello`, "hello world", true},
		{`^hello`, "world hello", false},
		{`[a-z]+@[a-z]+\.[a-z]+`, "test@example.com", true},
	}

	for _, tt := range tests {
		result, err := RegexMatch(tt.pattern, tt.text)
		if err != nil {
			t.Errorf("RegexMatch(%q, %q) error: %v", tt.pattern, tt.text, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("RegexMatch(%q, %q) = %v, expected %v", tt.pattern, tt.text, result, tt.expected)
		}
	}
}

func TestRegexFind(t *testing.T) {
	result, err := RegexFind(`\d+`, "abc123def456")
	if err != nil {
		t.Fatalf("RegexFind error: %v", err)
	}
	if result != "123" {
		t.Errorf("RegexFind = %q, expected %q", result, "123")
	}
}

func TestRegexFindAll(t *testing.T) {
	results, err := RegexFindAll(`\d+`, "abc123def456ghi789")
	if err != nil {
		t.Fatalf("RegexFindAll error: %v", err)
	}
	expected := []string{"123", "456", "789"}
	if len(results) != len(expected) {
		t.Fatalf("RegexFindAll returned %d results, expected %d", len(results), len(expected))
	}
	for i, r := range results {
		if r != expected[i] {
			t.Errorf("Result[%d] = %q, expected %q", i, r, expected[i])
		}
	}
}

func TestRegexReplace(t *testing.T) {
	result, err := RegexReplace(`\d+`, "abc123def456", "X")
	if err != nil {
		t.Fatalf("RegexReplace error: %v", err)
	}
	if result != "abcXdefX" {
		t.Errorf("RegexReplace = %q, expected %q", result, "abcXdefX")
	}
}

func TestRegexExtract(t *testing.T) {
	pattern := `(?P<area>\d{3})-(?P<exchange>\d{3})-(?P<number>\d{4})`
	text := "Call me at 555-123-4567"

	result, err := RegexExtract(pattern, text)
	if err != nil {
		t.Fatalf("RegexExtract error: %v", err)
	}

	if result["area"] != "555" {
		t.Errorf("area = %q, expected %q", result["area"], "555")
	}
	if result["exchange"] != "123" {
		t.Errorf("exchange = %q, expected %q", result["exchange"], "123")
	}
	if result["number"] != "4567" {
		t.Errorf("number = %q, expected %q", result["number"], "4567")
	}
}

func TestRegexTool(t *testing.T) {
	// Test match action
	result, err := RegexTool.Execute(context.Background(), map[string]any{
		"action":  "match",
		"pattern": `\d+`,
		"text":    "abc123",
	})
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Success {
		t.Errorf("Expected success, got error: %s", result.Error)
	}
	if result.Data != true {
		t.Errorf("Expected true, got %v", result.Data)
	}

	// Test findall action
	result, _ = RegexTool.Execute(context.Background(), map[string]any{
		"action":  "findall",
		"pattern": `\d+`,
		"text":    "a1b2c3",
	})
	matches := result.Data.([]string)
	if len(matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(matches))
	}

	// Test invalid pattern
	result, _ = RegexTool.Execute(context.Background(), map[string]any{
		"action":  "match",
		"pattern": `[`,
		"text":    "test",
	})
	if result.Success {
		t.Error("Expected failure for invalid pattern")
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		value    float64
		from     string
		to       string
		expected float64
	}{
		// Length
		{1, "km", "m", 1000},
		{1000, "m", "km", 1},
		{1, "mi", "km", 1.609344},
		{1, "ft", "m", 0.3048},
		{12, "in", "ft", 1},

		// Weight
		{1, "kg", "g", 1000},
		{1000, "g", "kg", 1},
		{1, "lb", "kg", 0.453592},

		// Time
		{1, "h", "min", 60},
		{60, "min", "h", 1},
		{1, "day", "h", 24},
		{24, "hr", "day", 1},

		// Data
		{1, "kb", "b", 1024},
		{1024, "mb", "gb", 1},
		{1, "gb", "mb", 1024},

		// Temperature
		{0, "c", "f", 32},
		{100, "c", "f", 212},
		{32, "f", "c", 0},
		{0, "c", "k", 273.15},
		{273.15, "k", "c", 0},
	}

	for _, tt := range tests {
		t.Run(tt.from+"_to_"+tt.to, func(t *testing.T) {
			result, err := Convert(tt.value, tt.from, tt.to)
			if err != nil {
				t.Fatalf("Convert(%v, %s, %s) error: %v", tt.value, tt.from, tt.to, err)
			}
			if math.Abs(result-tt.expected) > 0.001 {
				t.Errorf("Convert(%v, %s, %s) = %v, expected %v", tt.value, tt.from, tt.to, result, tt.expected)
			}
		})
	}
}

func TestConvertTool(t *testing.T) {
	result, err := ConvertTool.Execute(context.Background(), map[string]any{
		"value": 100.0,
		"from":  "km",
		"to":    "m",
	})
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Success {
		t.Errorf("Expected success, got error: %s", result.Error)
	}
	if result.Data.(float64) != 100000 {
		t.Errorf("Expected 100000, got %v", result.Data)
	}

	// Test invalid conversion
	result, _ = ConvertTool.Execute(context.Background(), map[string]any{
		"value": 100.0,
		"from":  "km",
		"to":    "kg",
	})
	if result.Success {
		t.Error("Expected failure for incompatible units")
	}
}

func TestConvertTemperatureFullNames(t *testing.T) {
	result, err := Convert(0, "celsius", "fahrenheit")
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}
	if result != 32 {
		t.Errorf("Expected 32, got %v", result)
	}

	result, err = Convert(0, "celsius", "kelvin")
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}
	if result != 273.15 {
		t.Errorf("Expected 273.15, got %v", result)
	}
}
