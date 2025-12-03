package ops

import (
	"testing"
)

func TestNormalizeOptions(t *testing.T) {
	t.Run("NewNormalizeOptions creates valid defaults", func(t *testing.T) {
		opts := NewNormalizeOptions()
		if err := opts.Validate(); err != nil {
			t.Errorf("default options should be valid: %v", err)
		}
	})

	t.Run("WithStandard sets standard", func(t *testing.T) {
		opts := NewNormalizeOptions().WithStandard("ISO 8601")
		if opts.Standard != "ISO 8601" {
			t.Errorf("expected ISO 8601, got %s", opts.Standard)
		}
	})

	t.Run("WithRules sets rules", func(t *testing.T) {
		rules := map[string]string{
			"phone": "E.164 format",
		}
		opts := NewNormalizeOptions().WithRules(rules)
		if len(opts.Rules) != 1 {
			t.Errorf("expected 1 rule, got %d", len(opts.Rules))
		}
	})

	t.Run("WithFixTypos enables typo fixing", func(t *testing.T) {
		opts := NewNormalizeOptions().WithFixTypos(true)
		if !opts.FixTypos {
			t.Error("expected FixTypos to be true")
		}
	})

	t.Run("WithCanonicalMappings sets mappings", func(t *testing.T) {
		mappings := map[string]string{
			"US": "United States",
		}
		opts := NewNormalizeOptions().WithCanonicalMappings(mappings)
		if len(opts.CanonicalMappings) != 1 {
			t.Errorf("expected 1 mapping, got %d", len(opts.CanonicalMappings))
		}
	})

	t.Run("WithStrict enables strict mode", func(t *testing.T) {
		opts := NewNormalizeOptions().WithStrict(true)
		if !opts.Strict {
			t.Error("expected Strict to be true")
		}
	})
}

func TestNormalize(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("Normalize standardizes data", func(t *testing.T) {
		type Address struct {
			Street  string `json:"street"`
			City    string `json:"city"`
			Country string `json:"country"`
		}

		address := Address{
			Street:  "123 Main St.",
			City:    "new york",
			Country: "USA",
		}

		opts := NewNormalizeOptions().
			WithRules(map[string]string{
				"city":    "Proper case",
				"country": "Full country name",
			}).
			WithStrict(true)

		result, err := Normalize(address, opts)
		if err != nil {
			t.Fatalf("Normalize failed: %v", err)
		}

		if result.Normalized.City != "New York" {
			t.Errorf("expected 'New York', got '%s'", result.Normalized.City)
		}
	})
}

func TestNormalizeText(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("NormalizeText standardizes text", func(t *testing.T) {
		text := "The mtg is tmrw @ 3pm. Pls confirm ASAP"
		opts := NewNormalizeOptions().
			WithRules(map[string]string{
				"abbreviations": "Expand to full words",
			})

		result, err := NormalizeText(text, opts)
		if err != nil {
			t.Fatalf("NormalizeText failed: %v", err)
		}

		if result == text {
			t.Error("expected text to be normalized")
		}
	})
}

func TestNormalizeBatch(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("NormalizeBatch processes multiple items", func(t *testing.T) {
		type Record struct {
			Date string `json:"date"`
		}

		records := []Record{
			{Date: "Jan 1, 2024"},
			{Date: "2024/01/02"},
			{Date: "03-01-2024"},
		}

		opts := NewNormalizeOptions().
			WithStandard("ISO 8601").
			WithRules(map[string]string{
				"date": "YYYY-MM-DD format",
			})

		results, err := NormalizeBatch(records, opts)
		if err != nil {
			t.Fatalf("NormalizeBatch failed: %v", err)
		}

		if len(results) != 3 {
			t.Errorf("expected 3 results, got %d", len(results))
		}
	})
}
