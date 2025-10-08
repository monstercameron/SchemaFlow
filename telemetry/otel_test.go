package schemaflow

import (
	"context"
	"testing"
	"time"
)

func TestStartSpan(t *testing.T) {
	// Test span creation
	ctx := context.Background()
	opts := OpOptions{
		Mode:         TransformMode,
		Intelligence: Fast,
	}
	
	newCtx, span := StartSpan(ctx, "test-operation", opts)
	if span == nil {
		t.Error("Expected span to be created")
	}
	defer span.End()

	if newCtx == nil {
		t.Error("Expected context with span")
	}
}

func TestRecordLLMCall(t *testing.T) {
	ctx := context.Background()
	opts := OpOptions{
		Mode:         TransformMode,
		Intelligence: Fast,
	}
	
	_, span := StartSpan(ctx, "test-operation", opts)
	defer span.End()

	// Test recording LLM call
	usage := &TokenUsage{
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
	}
	
	cost := &CostInfo{
		TotalCost: 0.05,
		Currency:  "USD",
	}
	
	RecordLLMCall(span, "gpt-3.5-turbo", "openai", usage, cost, 1*time.Second, nil)
	
	// No error expected, just ensuring it doesn't panic
}

func TestAddSpanTags(t *testing.T) {
	ctx := context.Background()
	opts := OpOptions{
		Mode:         TransformMode,
		Intelligence: Fast,
	}
	
	newCtx, span := StartSpan(ctx, "test-operation", opts)
	defer span.End()

	// Test adding span tags
	tags := map[string]string{
		"user":      "test-user",
		"operation": "test",
	}
	
	AddSpanTags(newCtx, tags)
	
	// No error expected, just ensuring it doesn't panic
}

func TestGetSpanID(t *testing.T) {
	ctx := context.Background()
	opts := OpOptions{
		Mode:         TransformMode,
		Intelligence: Fast,
	}
	
	newCtx, span := StartSpan(ctx, "test-operation", opts)
	defer span.End()

	// Test getting span ID
	spanID := GetSpanID(newCtx)
	if spanID == "" {
		t.Error("Expected span ID to be set")
	}
}