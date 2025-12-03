package ops

import (
	"testing"
)

func TestMatchOptions(t *testing.T) {
	t.Run("NewMatchOptions creates valid defaults", func(t *testing.T) {
		opts := NewMatchOptions()
		if err := opts.Validate(); err != nil {
			t.Errorf("default options should be valid: %v", err)
		}
	})

	t.Run("WithStrategy sets strategy", func(t *testing.T) {
		opts := NewMatchOptions().WithStrategy("semantic")
		if opts.Strategy != "semantic" {
			t.Errorf("expected semantic, got %s", opts.Strategy)
		}
	})

	t.Run("WithThreshold sets threshold", func(t *testing.T) {
		opts := NewMatchOptions().WithThreshold(0.9)
		if opts.Threshold != 0.9 {
			t.Errorf("expected 0.9, got %f", opts.Threshold)
		}
	})

	t.Run("WithMatchFields sets fields", func(t *testing.T) {
		opts := NewMatchOptions().WithMatchFields([]string{"name", "description"})
		if len(opts.MatchFields) != 2 {
			t.Errorf("expected 2 fields, got %d", len(opts.MatchFields))
		}
	})

	t.Run("WithFieldWeights sets weights", func(t *testing.T) {
		weights := map[string]float64{"name": 2.0}
		opts := NewMatchOptions().WithFieldWeights(weights)
		if opts.FieldWeights["name"] != 2.0 {
			t.Errorf("expected weight 2.0, got %f", opts.FieldWeights["name"])
		}
	})

	t.Run("WithAllowPartial enables partial matching", func(t *testing.T) {
		opts := NewMatchOptions().WithAllowPartial(true)
		if !opts.AllowPartial {
			t.Error("expected AllowPartial to be true")
		}
	})

	t.Run("Validate rejects invalid strategy", func(t *testing.T) {
		opts := NewMatchOptions()
		opts.Strategy = "invalid"
		if err := opts.Validate(); err == nil {
			t.Error("expected error for invalid strategy")
		}
	})

	t.Run("Validate rejects invalid threshold", func(t *testing.T) {
		opts := NewMatchOptions().WithThreshold(1.5)
		if err := opts.Validate(); err == nil {
			t.Error("expected error for threshold > 1")
		}
	})
}

func TestSemanticMatch(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("SemanticMatch finds best matches", func(t *testing.T) {
		type Product struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		type Query struct {
			SearchTerm string `json:"search_term"`
		}

		products := []Product{
			{Name: "Blue Widget", Description: "A blue colored widget"},
			{Name: "Red Gadget", Description: "A red colored gadget"},
			{Name: "Green Widget", Description: "A green colored widget"},
		}

		queries := []Query{
			{SearchTerm: "blue thing"},
			{SearchTerm: "red device"},
		}

		opts := NewMatchOptions().
			WithStrategy("semantic").
			WithThreshold(0.5)

		result, err := SemanticMatch(queries, products, opts)
		if err != nil {
			t.Fatalf("SemanticMatch failed: %v", err)
		}

		if len(result.Matches) == 0 {
			t.Error("expected matched pairs, got none")
		}
	})
}

func TestMatchOne(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("MatchOne finds single best match", func(t *testing.T) {
		type Item struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		type Query struct {
			Term string `json:"term"`
		}

		query := Query{Term: "wireless headphones"}
		candidates := []Item{
			{ID: 1, Name: "Wired Earbuds"},
			{ID: 2, Name: "Bluetooth Headphones"},
			{ID: 3, Name: "USB Speaker"},
		}

		opts := NewMatchOptions().WithStrategy("best-fit")

		matches, err := MatchOne(query, candidates, opts)
		if err != nil {
			t.Fatalf("MatchOne failed: %v", err)
		}

		if len(matches) == 0 {
			t.Error("expected at least one match")
		}
	})
}
