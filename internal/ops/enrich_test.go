package ops

import (
	"testing"
)

func TestEnrichOptions(t *testing.T) {
	t.Run("NewEnrichOptions creates valid defaults", func(t *testing.T) {
		opts := NewEnrichOptions()
		if err := opts.Validate(); err != nil {
			t.Errorf("default options should be valid: %v", err)
		}
	})

	t.Run("WithDeriveFields sets fields to derive", func(t *testing.T) {
		opts := NewEnrichOptions().WithDeriveFields([]string{"summary", "keywords"})
		if len(opts.DeriveFields) != 2 {
			t.Errorf("expected 2 fields, got %d", len(opts.DeriveFields))
		}
	})

	t.Run("WithDerivationRules sets rules", func(t *testing.T) {
		rules := map[string]string{
			"summary": "A brief summary of the content",
		}
		opts := NewEnrichOptions().WithDerivationRules(rules)
		if len(opts.DerivationRules) != 1 {
			t.Errorf("expected 1 rule, got %d", len(opts.DerivationRules))
		}
	})

	t.Run("WithDomain sets domain", func(t *testing.T) {
		opts := NewEnrichOptions().WithDomain("finance")
		if opts.Domain != "finance" {
			t.Errorf("expected finance, got %s", opts.Domain)
		}
	})

	t.Run("WithIncludeConfidence enables confidence", func(t *testing.T) {
		opts := NewEnrichOptions().WithIncludeConfidence(true)
		if !opts.IncludeConfidence {
			t.Error("expected IncludeConfidence to be true")
		}
	})
}

func TestEnrich(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("Enrich adds derived fields", func(t *testing.T) {
		type Article struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		type EnrichedArticle struct {
			Title    string   `json:"title"`
			Content  string   `json:"content"`
			Summary  string   `json:"summary"`
			Keywords []string `json:"keywords"`
		}

		article := Article{
			Title:   "Go Programming",
			Content: "Go is a statically typed language designed for simplicity.",
		}

		opts := NewEnrichOptions().
			WithDeriveFields([]string{"summary", "keywords"})

		result, err := Enrich[Article, EnrichedArticle](article, opts)
		if err != nil {
			t.Fatalf("Enrich failed: %v", err)
		}

		if result.Enriched.Summary == "" {
			t.Error("expected enriched summary")
		}
	})
}

func TestEnrichInPlace(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("EnrichInPlace modifies in place", func(t *testing.T) {
		type Document struct {
			Text     string   `json:"text"`
			Keywords []string `json:"keywords"`
		}

		doc := Document{
			Text: "Machine learning is transforming industries.",
		}

		opts := NewEnrichOptions().
			WithDeriveFields([]string{"keywords"})

		enriched, err := EnrichInPlace(doc, opts)
		if err != nil {
			t.Fatalf("EnrichInPlace failed: %v", err)
		}

		if len(enriched.Keywords) == 0 {
			t.Error("expected keywords to be added")
		}
	})
}
