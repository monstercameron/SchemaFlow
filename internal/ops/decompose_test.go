package ops

import (
	"testing"
)

func TestDecomposeOptions(t *testing.T) {
	t.Run("NewDecomposeOptions creates valid defaults", func(t *testing.T) {
		opts := NewDecomposeOptions()
		if err := opts.Validate(); err != nil {
			t.Errorf("default options should be valid: %v", err)
		}
	})

	t.Run("WithStrategy sets strategy", func(t *testing.T) {
		opts := NewDecomposeOptions().WithStrategy("hierarchical")
		if opts.Strategy != "hierarchical" {
			t.Errorf("expected hierarchical, got %s", opts.Strategy)
		}
	})

	t.Run("WithMaxDepth sets max depth", func(t *testing.T) {
		opts := NewDecomposeOptions().WithMaxDepth(5)
		if opts.MaxDepth != 5 {
			t.Errorf("expected 5, got %d", opts.MaxDepth)
		}
	})

	t.Run("WithTargetParts sets target parts", func(t *testing.T) {
		opts := NewDecomposeOptions().WithTargetParts(4)
		if opts.TargetParts != 4 {
			t.Errorf("expected 4, got %d", opts.TargetParts)
		}
	})

	t.Run("WithIncludeDependencies enables dependencies", func(t *testing.T) {
		opts := NewDecomposeOptions().WithIncludeDependencies(true)
		if !opts.IncludeDependencies {
			t.Error("expected IncludeDependencies to be true")
		}
	})

	t.Run("Validate rejects invalid strategy", func(t *testing.T) {
		opts := NewDecomposeOptions()
		opts.Strategy = "invalid"
		if err := opts.Validate(); err == nil {
			t.Error("expected error for invalid strategy")
		}
	})

	t.Run("Validate rejects negative depth", func(t *testing.T) {
		opts := NewDecomposeOptions().WithMaxDepth(-1)
		if err := opts.Validate(); err == nil {
			t.Error("expected error for negative depth")
		}
	})
}

func TestDecompose(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("Decompose breaks down complex task", func(t *testing.T) {
		task := "Build a web application with user authentication and database"
		opts := NewDecomposeOptions().
			WithStrategy("sequential").
			WithIncludeDependencies(true)

		result, err := Decompose(task, opts)
		if err != nil {
			t.Fatalf("Decompose failed: %v", err)
		}

		if len(result.Parts) == 0 {
			t.Error("expected decomposed parts, got none")
		}
	})
}

func TestDecomposeToSlice(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("DecomposeToSlice returns typed slice", func(t *testing.T) {
		type Step struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		task := "Deploy an application to production"
		opts := NewDecomposeOptions().WithStrategy("sequential")

		result, err := DecomposeToSlice[string, Step](task, opts)
		if err != nil {
			t.Fatalf("DecomposeToSlice failed: %v", err)
		}

		if len(result) == 0 {
			t.Error("expected parts, got none")
		}
	})
}
