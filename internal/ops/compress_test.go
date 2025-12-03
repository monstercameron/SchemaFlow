package ops

import (
	"testing"
)

func TestCompressOptions(t *testing.T) {
	t.Run("NewCompressOptions creates valid defaults", func(t *testing.T) {
		opts := NewCompressOptions()
		if err := opts.Validate(); err != nil {
			t.Errorf("default options should be valid: %v", err)
		}
	})

	t.Run("WithCompressionRatio sets ratio", func(t *testing.T) {
		opts := NewCompressOptions().WithCompressionRatio(0.5)
		if opts.CompressionRatio != 0.5 {
			t.Errorf("expected 0.5, got %f", opts.CompressionRatio)
		}
	})

	t.Run("WithStrategy sets strategy", func(t *testing.T) {
		opts := NewCompressOptions().WithStrategy("abstractive")
		if opts.Strategy != "abstractive" {
			t.Errorf("expected abstractive, got %s", opts.Strategy)
		}
	})

	t.Run("WithPriority sets priority", func(t *testing.T) {
		opts := NewCompressOptions().WithPriority("key_facts")
		if opts.Priority != "key_facts" {
			t.Errorf("expected key_facts, got %s", opts.Priority)
		}
	})

	t.Run("WithPreserveInfo sets info to preserve", func(t *testing.T) {
		opts := NewCompressOptions().WithPreserveInfo([]string{"dates", "names"})
		if len(opts.PreserveInfo) != 2 {
			t.Errorf("expected 2 preserve items, got %d", len(opts.PreserveInfo))
		}
	})

	t.Run("Validate rejects invalid compression ratio", func(t *testing.T) {
		opts := NewCompressOptions().WithCompressionRatio(1.5)
		if err := opts.Validate(); err == nil {
			t.Error("expected error for ratio > 1")
		}
	})

	t.Run("Validate rejects invalid strategy", func(t *testing.T) {
		opts := NewCompressOptions()
		opts.Strategy = "invalid"
		if err := opts.Validate(); err == nil {
			t.Error("expected error for invalid strategy")
		}
	})
}

func TestCompress(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("Compress reduces text length", func(t *testing.T) {
		text := `This is a long document that contains many important details 
		about various topics. It discusses programming, software development, 
		and best practices for writing clean code. The document also covers 
		testing strategies and deployment procedures.`

		opts := NewCompressOptions().
			WithCompressionRatio(0.5).
			WithStrategy("extractive")

		result, err := Compress(text, opts)
		if err != nil {
			t.Fatalf("Compress failed: %v", err)
		}

		if result.ActualRatio >= 1.0 {
			t.Error("expected compression to reduce size")
		}
	})
}

func TestCompressText(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("CompressText returns compressed string", func(t *testing.T) {
		text := "This is a long text that should be compressed to a shorter version."
		opts := NewCompressOptions().WithCompressionRatio(0.5)

		result, err := CompressText(text, opts)
		if err != nil {
			t.Fatalf("CompressText failed: %v", err)
		}

		if result == "" {
			t.Error("expected non-empty compressed text")
		}
	})
}
