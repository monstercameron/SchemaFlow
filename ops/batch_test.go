package ops

import (
	"testing"
	"time"
)

func TestBatchOperations(t *testing.T) {
	// Setup
	setupMockClient()
	
	t.Run("ParallelMode", func(t *testing.T) {
		inputs := []interface{}{
			"John Doe, 30 years old",
			"Jane Smith, 25 years old",
			"Bob Johnson, 35 years old",
		}
		
		batch := Batch().
			WithMode(ParallelMode).
			WithConcurrency(2).
			WithTimeout(10 * time.Second)
		
		results := ExtractBatch[Person](batch, inputs)
		
		// Check metadata
		if results.Metadata.Mode != ParallelMode {
			t.Errorf("Expected ParallelMode, got %v", results.Metadata.Mode)
		}
		
		if results.Metadata.TotalItems != 3 {
			t.Errorf("Expected 3 total items, got %d", results.Metadata.TotalItems)
		}
		
		// Check results
		for i, err := range results.Errors {
			if err != nil {
				t.Errorf("Unexpected error at index %d: %v", i, err)
			}
		}
	})
	
	t.Run("MergedMode", func(t *testing.T) {
		inputs := []interface{}{
			"Alice Brown, 28 years old",
			"Charlie Davis, 32 years old",
			"Eve Wilson, 29 years old",
			"Frank Miller, 31 years old",
			"Grace Lee, 27 years old",
		}
		
		batch := Batch().
			WithMode(MergedMode).
			WithBatchSize(3).
			WithTimeout(10 * time.Second)
		
		results := ExtractBatch[Person](batch, inputs)
		
		// Check metadata
		if results.Metadata.Mode != MergedMode {
			t.Errorf("Expected MergedMode, got %v", results.Metadata.Mode)
		}
		
		if results.Metadata.TotalItems != 5 {
			t.Errorf("Expected 5 total items, got %d", results.Metadata.TotalItems)
		}
		
		// In merged mode, we should have fewer API calls than items
		expectedCalls := 2 // 5 items with batch size 3 = 2 calls
		if results.Metadata.APICallsMade > expectedCalls {
			t.Logf("Warning: More API calls than expected. Got %d, expected <= %d",
				results.Metadata.APICallsMade, expectedCalls)
		}
	})
	
	t.Run("SmartBatch", func(t *testing.T) {
		// Small batch - should use ParallelMode
		smallInputs := []interface{}{
			"Item 1",
			"Item 2",
		}
		
		smartBatch := NewClient("").SmartBatch()
		results := ExtractSmart[Person](smartBatch, smallInputs)
		
		if results.Metadata.Mode != ParallelMode {
			t.Errorf("Expected ParallelMode for small batch, got %v (value=%d)", results.Metadata.Mode, results.Metadata.Mode)
		}
		
		// Large batch - should use MergedMode
		largeInputs := make([]interface{}, 25)
		for i := range largeInputs {
			largeInputs[i] = "Test input"
		}
		
		results = ExtractSmart[Person](smartBatch, largeInputs)
		
		if results.Metadata.Mode != MergedMode {
			t.Errorf("Expected MergedMode for large batch, got %v", results.Metadata.Mode)
		}
	})
}

func TestBatchChunking(t *testing.T) {
	batch := Batch().WithBatchSize(3)
	
	inputs := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	chunks := batch.createChunks(inputs, 3)
	
	if len(chunks) != 4 {
		t.Errorf("Expected 4 chunks, got %d", len(chunks))
	}
	
	// Check chunk sizes
	expectedSizes := []int{3, 3, 3, 1}
	for i, chunk := range chunks {
		if len(chunk) != expectedSizes[i] {
			t.Errorf("Chunk %d: expected size %d, got %d", i, expectedSizes[i], len(chunk))
		}
	}
}

func TestBatchMetadata(t *testing.T) {
	setupMockClient()
	
	inputs := []interface{}{
		"Test 1",
		"Test 2",
		"Test 3",
	}
	
	batch := Batch().WithMode(ParallelMode)
	results := ExtractBatch[Person](batch, inputs)
	
	// Check metadata calculations
	metadata := results.Metadata
	
	if metadata.TotalItems != len(inputs) {
		t.Errorf("TotalItems mismatch: got %d, want %d", metadata.TotalItems, len(inputs))
	}
	
	if metadata.Succeeded+metadata.Failed != metadata.TotalItems {
		t.Errorf("Succeeded+Failed should equal TotalItems")
	}
	
	if metadata.Duration == 0 {
		t.Error("Duration should be non-zero")
	}
}

// Benchmark batch operations
func BenchmarkBatchParallel(b *testing.B) {
	setupMockClient()
	inputs := make([]interface{}, 100)
	for i := range inputs {
		inputs[i] = "Test input"
	}
	
	batch := Batch().WithMode(ParallelMode).WithConcurrency(10)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ExtractBatch[Person](batch, inputs)
	}
}

func BenchmarkBatchMerged(b *testing.B) {
	setupMockClient()
	inputs := make([]interface{}, 100)
	for i := range inputs {
		inputs[i] = "Test input"
	}
	
	batch := Batch().WithMode(MergedMode).WithBatchSize(20)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ExtractBatch[Person](batch, inputs)
	}
}
