package ops

import (
	"testing"
)

func TestRankOptions(t *testing.T) {
	t.Run("NewRankOptions creates valid defaults", func(t *testing.T) {
		opts := NewRankOptions()
		if err := opts.Validate(); err != nil {
			t.Errorf("default options should be valid: %v", err)
		}
	})

	t.Run("WithQuery sets query", func(t *testing.T) {
		opts := NewRankOptions().WithQuery("best programming languages")
		if opts.Query != "best programming languages" {
			t.Errorf("expected query to be set")
		}
	})

	t.Run("WithTopK sets limit", func(t *testing.T) {
		opts := NewRankOptions().WithTopK(10)
		if opts.TopK != 10 {
			t.Errorf("expected 10, got %d", opts.TopK)
		}
	})

	t.Run("WithBoostFields sets boost fields", func(t *testing.T) {
		opts := NewRankOptions().WithBoostFields(map[string]float64{"title": 2.0, "description": 1.5})
		if len(opts.BoostFields) != 2 {
			t.Errorf("expected 2 boost fields, got %d", len(opts.BoostFields))
		}
	})

	t.Run("WithPenalizeFields sets penalize fields", func(t *testing.T) {
		opts := NewRankOptions().WithPenalizeFields(map[string]float64{"spam_score": 0.5})
		if len(opts.PenalizeFields) != 1 {
			t.Errorf("expected 1 penalize field, got %d", len(opts.PenalizeFields))
		}
	})

	t.Run("WithIncludeExplanation enables explanation", func(t *testing.T) {
		opts := NewRankOptions().WithIncludeExplanation(true)
		if !opts.IncludeExplanation {
			t.Error("expected IncludeExplanation to be true")
		}
	})

	t.Run("Validate requires query", func(t *testing.T) {
		opts := RankOptions{}
		if err := opts.Validate(); err == nil {
			t.Error("expected error for empty query")
		}
	})

	t.Run("Validate rejects negative TopK", func(t *testing.T) {
		opts := NewRankOptions().WithQuery("test").WithTopK(-1)
		if err := opts.Validate(); err == nil {
			t.Error("expected error for negative TopK")
		}
	})
}

func TestRank(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("Rank orders by relevance", func(t *testing.T) {
		type Article struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		articles := []Article{
			{Title: "Go Programming", Content: "Learn Go basics"},
			{Title: "Python Tips", Content: "Python programming tips"},
			{Title: "Go Advanced", Content: "Advanced Go patterns"},
		}

		opts := NewRankOptions().
			WithQuery("Go programming").
			WithTopK(2)

		result, err := Rank(articles, opts)
		if err != nil {
			t.Fatalf("Rank failed: %v", err)
		}

		if len(result.Items) == 0 {
			t.Error("expected ranked items, got none")
		}
	})
}
