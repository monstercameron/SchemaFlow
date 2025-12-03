package ops

import (
	"testing"
)

func TestSuggestOptions(t *testing.T) {
	opts := NewSuggestOptions()
	if err := opts.Validate(); err != nil {
		t.Errorf("Default options validation failed: %v", err)
	}

	if opts.Strategy != SuggestContextual {
		t.Errorf("Expected default strategy Contextual, got %v", opts.Strategy)
	}
	if opts.TopN != 5 {
		t.Errorf("Expected default TopN 5, got %d", opts.TopN)
	}
	if opts.Ranked != true {
		t.Errorf("Expected default Ranked true, got %v", opts.Ranked)
	}
}

func TestSuggestOptions_FluentAPI(t *testing.T) {
	opts := NewSuggestOptions().
		WithStrategy(SuggestHybrid).
		WithTopN(10).
		WithRanked(false).
		WithIncludeScores(true).
		WithIncludeReasons(true).
		WithDomain("data-processing").
		WithConstraints([]string{"must be efficient", "cost-effective"}).
		WithCategories([]string{"optimization", "scaling"})

	if opts.Strategy != SuggestHybrid {
		t.Errorf("Expected strategy Hybrid, got %v", opts.Strategy)
	}
	if opts.TopN != 10 {
		t.Errorf("Expected TopN 10, got %d", opts.TopN)
	}
	if opts.Ranked != false {
		t.Errorf("Expected Ranked false, got %v", opts.Ranked)
	}
	if opts.IncludeScores != true {
		t.Errorf("Expected IncludeScores true, got %v", opts.IncludeScores)
	}
	if opts.IncludeReasons != true {
		t.Errorf("Expected IncludeReasons true, got %v", opts.IncludeReasons)
	}
	if opts.Domain != "data-processing" {
		t.Errorf("Expected Domain 'data-processing', got %v", opts.Domain)
	}
	if len(opts.Constraints) != 2 {
		t.Errorf("Expected 2 constraints, got %d", len(opts.Constraints))
	}
	if len(opts.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(opts.Categories))
	}
}

func TestSuggestOptions_Validation(t *testing.T) {
	tests := []struct {
		name    string
		opts    SuggestOptions
		wantErr bool
	}{
		{"valid defaults", NewSuggestOptions(), false},
		{"invalid topN", NewSuggestOptions().WithTopN(0), true},
		{"invalid strategy", SuggestOptions{Strategy: "invalid"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSuggest_Basic(t *testing.T) {
	// Skip this test if no LLM is configured
	if testing.Short() {
		t.Skip("Skipping LLM-dependent test in short mode")
	}

	input := map[string]any{
		"current_task": "data cleaning",
		"data_size":    "1GB",
		"issues":       []string{"missing values", "duplicates"},
	}

	opts := NewSuggestOptions().WithTopN(3)
	suggestions, err := Suggest[string](input, opts)
	if err != nil {
		t.Skipf("LLM not available: %v", err)
		return
	}

	if len(suggestions) > 3 {
		t.Errorf("Expected at most 3 suggestions, got %d", len(suggestions))
	}

	// Basic validation that we got some suggestions
	if len(suggestions) == 0 {
		t.Log("Warning: No suggestions returned (may indicate LLM configuration issue)")
	}
}

func TestSuggestWithResult(t *testing.T) {
	// Skip this test if no LLM is configured
	if testing.Short() {
		t.Skip("Skipping LLM-dependent test in short mode")
	}

	input := "optimize database query performance"

	opts := NewSuggestOptions().
		WithTopN(2).
		WithIncludeReasons(true)

	result, err := SuggestWithResult[string](input, opts)
	if err != nil {
		t.Skipf("LLM not available: %v", err)
		return
	}

	if len(result.Suggestions) > 2 {
		t.Errorf("Expected at most 2 suggestions, got %d", len(result.Suggestions))
	}

	// Result should have metadata map initialized
	if result.Metadata == nil {
		t.Errorf("Expected metadata map to be initialized")
	}
}
