package ops

import (
	"context"
	"strings"
	"testing"

	"github.com/monstercameron/SchemaFlow/internal/llm"
	"github.com/monstercameron/SchemaFlow/internal/types"
)

type captureProvider struct {
	req llm.CompletionRequest
}

func (p *captureProvider) Complete(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
	p.req = req
	return llm.CompletionResponse{Content: `{"ok":true}`}, nil
}

func (p *captureProvider) Name() string {
	return "local"
}

func (p *captureProvider) EstimateCost(req llm.CompletionRequest) float64 {
	return 0
}

func TestInferResponseFormat(t *testing.T) {
	tests := []struct {
		name       string
		system     string
		user       string
		wantFormat string
	}{
		{
			name:       "json object contract",
			system:     "Return a JSON object with fields name and age.",
			wantFormat: "json",
		},
		{
			name:       "schema contract",
			system:     "Return ONLY valid JSON matching the schema.",
			wantFormat: "json",
		},
		{
			name:       "plain text summary",
			system:     "Summarize the text in two sentences.",
			wantFormat: "text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferResponseFormat(tt.system, tt.user); got != tt.wantFormat {
				t.Fatalf("inferResponseFormat() = %q, want %q", got, tt.wantFormat)
			}
		})
	}
}

func TestCallLLMUsesStructuredContracts(t *testing.T) {
	provider := &captureProvider{}

	_, err := CallLLM(
		context.Background(),
		provider,
		`You are a ranking expert. Return a JSON object with rankings.`,
		`Rank these items.`,
		types.OpOptions{Intelligence: types.Fast, Mode: types.TransformMode},
	)
	if err != nil {
		t.Fatalf("CallLLM() error = %v", err)
	}

	if provider.req.ResponseFormat != "json" {
		t.Fatalf("ResponseFormat = %q, want json", provider.req.ResponseFormat)
	}
	if !strings.Contains(provider.req.SystemPrompt, "Perform the semantic task faithfully") {
		t.Fatalf("system prompt missing semantic grounding: %q", provider.req.SystemPrompt)
	}
	if !strings.Contains(provider.req.SystemPrompt, "return only the final JSON answer") {
		t.Fatalf("system prompt missing JSON grounding: %q", provider.req.SystemPrompt)
	}
}

func TestCallLLMLeavesTextOpsAsText(t *testing.T) {
	provider := &captureProvider{}

	_, err := CallLLM(
		context.Background(),
		provider,
		`You are a text summarization expert. Create concise summaries that preserve key information.`,
		`Summarize this article.`,
		types.OpOptions{Intelligence: types.Fast, Mode: types.TransformMode},
	)
	if err != nil {
		t.Fatalf("CallLLM() error = %v", err)
	}

	if provider.req.ResponseFormat != "text" {
		t.Fatalf("ResponseFormat = %q, want text", provider.req.ResponseFormat)
	}
	if strings.Contains(provider.req.SystemPrompt, "return only the final JSON answer") {
		t.Fatalf("text prompt unexpectedly forced JSON rules: %q", provider.req.SystemPrompt)
	}
}

func TestCallLLMAppliesSteeringToSystemPrompt(t *testing.T) {
	provider := &captureProvider{}

	_, err := CallLLM(
		context.Background(),
		provider,
		`You are a filtering expert.`,
		`Filter these items.`,
		types.OpOptions{
			Intelligence: types.Fast,
			Mode:         types.TransformMode,
			Steering:     "Return only a JSON array of matching strings.",
		},
	)
	if err != nil {
		t.Fatalf("CallLLM() error = %v", err)
	}

	if !strings.Contains(provider.req.SystemPrompt, "Additional instructions:") {
		t.Fatalf("system prompt missing steering section: %q", provider.req.SystemPrompt)
	}
	if !strings.Contains(provider.req.SystemPrompt, "Return only a JSON array of matching strings.") {
		t.Fatalf("system prompt missing steering content: %q", provider.req.SystemPrompt)
	}
}
