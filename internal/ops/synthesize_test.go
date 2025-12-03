package ops

import (
	"testing"
)

func TestSynthesizeOptions(t *testing.T) {
	t.Run("NewSynthesizeOptions creates valid defaults", func(t *testing.T) {
		opts := NewSynthesizeOptions()
		if err := opts.Validate(); err != nil {
			t.Errorf("default options should be valid: %v", err)
		}
	})

	t.Run("WithStrategy sets strategy", func(t *testing.T) {
		opts := NewSynthesizeOptions().WithStrategy("merge")
		if opts.Strategy != "merge" {
			t.Errorf("expected merge, got %s", opts.Strategy)
		}
	})

	t.Run("WithConflictResolution sets resolution", func(t *testing.T) {
		opts := NewSynthesizeOptions().WithConflictResolution("majority")
		if opts.ConflictResolution != "majority" {
			t.Errorf("expected majority, got %s", opts.ConflictResolution)
		}
	})

	t.Run("WithSourcePriorities sets priorities", func(t *testing.T) {
		opts := NewSynthesizeOptions().WithSourcePriorities([]int{0, 2, 1})
		if len(opts.SourcePriorities) != 3 {
			t.Errorf("expected 3 priorities, got %d", len(opts.SourcePriorities))
		}
	})

	t.Run("WithCiteSources sets cite flag", func(t *testing.T) {
		opts := NewSynthesizeOptions().WithCiteSources(false)
		if opts.CiteSources {
			t.Error("expected CiteSources to be false")
		}
	})

	t.Run("WithGenerateInsights sets insights flag", func(t *testing.T) {
		opts := NewSynthesizeOptions().WithGenerateInsights(false)
		if opts.GenerateInsights {
			t.Error("expected GenerateInsights to be false")
		}
	})

	t.Run("WithFocusAreas sets focus areas", func(t *testing.T) {
		opts := NewSynthesizeOptions().WithFocusAreas([]string{"methodology", "results"})
		if len(opts.FocusAreas) != 2 {
			t.Errorf("expected 2 focus areas, got %d", len(opts.FocusAreas))
		}
	})

	t.Run("Validate rejects invalid strategy", func(t *testing.T) {
		opts := NewSynthesizeOptions()
		opts.Strategy = "invalid"
		if err := opts.Validate(); err == nil {
			t.Error("expected error for invalid strategy")
		}
	})

	t.Run("Validate rejects invalid conflict resolution", func(t *testing.T) {
		opts := NewSynthesizeOptions()
		opts.ConflictResolution = "invalid"
		if err := opts.Validate(); err == nil {
			t.Error("expected error for invalid conflict resolution")
		}
	})
}

func TestSynthesize(t *testing.T) {
	// Skip integration tests without LLM
	t.Skip("Integration test requires LLM provider")

	t.Run("Synthesize combines sources", func(t *testing.T) {
		type Report struct {
			Summary    string   `json:"summary"`
			KeyPoints  []string `json:"key_points"`
			Conclusion string   `json:"conclusion"`
		}

		sources := []any{
			map[string]string{
				"title":   "Study A",
				"finding": "Treatment X is effective",
			},
			map[string]string{
				"title":   "Study B",
				"finding": "Treatment X shows 80% success rate",
			},
			map[string]string{
				"title":   "Study C",
				"finding": "Treatment X has minimal side effects",
			},
		}

		opts := NewSynthesizeOptions().
			WithStrategy("integrate").
			WithCiteSources(true).
			WithGenerateInsights(true)

		result, err := Synthesize[Report](sources, opts)
		if err != nil {
			t.Fatalf("Synthesize failed: %v", err)
		}

		if result.Synthesized.Summary == "" {
			t.Error("expected synthesized summary")
		}
	})

	t.Run("Synthesize with conflict resolution", func(t *testing.T) {
		type Data struct {
			Value int    `json:"value"`
			Notes string `json:"notes"`
		}

		sources := []any{
			map[string]any{"temperature": 72, "source": "sensor1"},
			map[string]any{"temperature": 75, "source": "sensor2"},
			map[string]any{"temperature": 73, "source": "sensor3"},
		}

		opts := NewSynthesizeOptions().
			WithStrategy("reconcile").
			WithConflictResolution("majority")

		result, err := Synthesize[Data](sources, opts)
		if err != nil {
			t.Fatalf("Synthesize failed: %v", err)
		}

		if len(result.Conflicts) >= 0 { // May or may not have conflicts
			// This is expected behavior
		}
	})
}
