package ops

import (
	"testing"
)

func TestClusterOptions(t *testing.T) {
	t.Run("NewClusterOptions creates valid defaults", func(t *testing.T) {
		opts := NewClusterOptions()
		if err := opts.Validate(); err != nil {
			t.Errorf("default options should be valid: %v", err)
		}
	})

	t.Run("WithNumClusters sets cluster count", func(t *testing.T) {
		opts := NewClusterOptions().WithNumClusters(5)
		if opts.NumClusters != 5 {
			t.Errorf("expected 5 clusters, got %d", opts.NumClusters)
		}
	})

	t.Run("WithNamingStrategy sets strategy", func(t *testing.T) {
		opts := NewClusterOptions().WithNamingStrategy("descriptive")
		if opts.NamingStrategy != "descriptive" {
			t.Errorf("expected descriptive, got %s", opts.NamingStrategy)
		}
	})

	t.Run("WithSimilarityThreshold sets threshold", func(t *testing.T) {
		opts := NewClusterOptions().WithSimilarityThreshold(0.8)
		if opts.SimilarityThreshold != 0.8 {
			t.Errorf("expected 0.8, got %f", opts.SimilarityThreshold)
		}
	})

	t.Run("WithIncludeOutliers enables outliers", func(t *testing.T) {
		opts := NewClusterOptions().WithIncludeOutliers(true)
		if !opts.IncludeOutliers {
			t.Error("expected IncludeOutliers to be true")
		}
	})

	t.Run("Validate rejects invalid naming strategy", func(t *testing.T) {
		opts := NewClusterOptions()
		opts.NamingStrategy = "invalid"
		if err := opts.Validate(); err == nil {
			t.Error("expected error for invalid naming strategy")
		}
	})

	t.Run("Validate rejects invalid similarity threshold", func(t *testing.T) {
		opts := NewClusterOptions().WithSimilarityThreshold(1.5)
		if err := opts.Validate(); err == nil {
			t.Error("expected error for threshold > 1")
		}
	})
}

func TestCluster(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("Cluster groups similar items", func(t *testing.T) {
		items := []string{
			"Python programming tutorial",
			"JavaScript for beginners",
			"Machine learning basics",
			"Deep learning with PyTorch",
			"React web development",
			"Vue.js framework guide",
		}

		opts := NewClusterOptions().
			WithNumClusters(2).
			WithNamingStrategy("descriptive")

		result, err := Cluster(items, opts)
		if err != nil {
			t.Fatalf("Cluster failed: %v", err)
		}

		if len(result.Clusters) == 0 {
			t.Error("expected clusters, got none")
		}
	})
}
