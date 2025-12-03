package ops

import (
	"testing"
)

func TestAnnotateOptions(t *testing.T) {
	t.Run("NewAnnotateOptions creates valid defaults", func(t *testing.T) {
		opts := NewAnnotateOptions()
		if err := opts.Validate(); err != nil {
			t.Errorf("default options should be valid: %v", err)
		}
	})

	t.Run("WithAnnotationTypes sets types", func(t *testing.T) {
		opts := NewAnnotateOptions().WithAnnotationTypes([]string{"entities", "sentiment"})
		if len(opts.AnnotationTypes) != 2 {
			t.Errorf("expected 2 types, got %d", len(opts.AnnotationTypes))
		}
	})

	t.Run("WithFormat sets format", func(t *testing.T) {
		opts := NewAnnotateOptions().WithFormat("inline")
		if opts.Format != "inline" {
			t.Errorf("expected inline, got %s", opts.Format)
		}
	})

	t.Run("WithIncludeConfidence sets flag", func(t *testing.T) {
		opts := NewAnnotateOptions().WithIncludeConfidence(true)
		if !opts.IncludeConfidence {
			t.Error("expected IncludeConfidence to be true")
		}
	})

	t.Run("WithMinConfidence sets confidence", func(t *testing.T) {
		opts := NewAnnotateOptions().WithMinConfidence(0.9)
		if opts.MinConfidence != 0.9 {
			t.Errorf("expected 0.9, got %f", opts.MinConfidence)
		}
	})

	t.Run("Validate rejects invalid confidence", func(t *testing.T) {
		opts := NewAnnotateOptions().WithMinConfidence(1.5)
		if err := opts.Validate(); err == nil {
			t.Error("expected error for confidence > 1")
		}
	})

	t.Run("Validate rejects invalid format", func(t *testing.T) {
		opts := NewAnnotateOptions().WithFormat("invalid_format")
		if err := opts.Validate(); err == nil {
			t.Error("expected error for invalid format")
		}
	})
}

func TestAnnotate(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("Annotate extracts entities", func(t *testing.T) {
		text := "John Smith visited Paris last week and met with Apple executives."
		opts := NewAnnotateOptions().
			WithAnnotationTypes([]string{"entities", "sentiment"}).
			WithFormat("structured")

		result, err := Annotate(text, opts)
		if err != nil {
			t.Fatalf("Annotate failed: %v", err)
		}

		if len(result.Annotations) == 0 {
			t.Error("expected annotations, got none")
		}
	})
}
