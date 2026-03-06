package schemaflow

import (
	"context"
	"testing"
)

func TestNewCommonOptionsDefaults(t *testing.T) {
	opts := NewCommonOptions()

	if opts.Mode != TransformMode {
		t.Fatalf("expected default mode %v, got %v", TransformMode, opts.Mode)
	}
	if opts.Intelligence != Fast {
		t.Fatalf("expected default intelligence %v, got %v", Fast, opts.Intelligence)
	}
}

func TestExtractingBuilderAppliesCommonAndOperationOptions(t *testing.T) {
	ctx := context.Background()
	req := Extracting[struct {
		Name string `json:"name"`
	}]("John Doe").
		Strict().
		Smart().
		Steer("focus on the contact record").
		Threshold(0.9).
		Context(ctx).
		RequestID("req-123").
		Partial(false).
		SchemaHints(map[string]string{"name": "Full legal name"})

	if req.opts.CommonOptions.Mode != Strict {
		t.Fatalf("expected strict mode, got %v", req.opts.CommonOptions.Mode)
	}
	if req.opts.CommonOptions.Intelligence != Smart {
		t.Fatalf("expected smart intelligence, got %v", req.opts.CommonOptions.Intelligence)
	}
	if req.opts.CommonOptions.Steering != "focus on the contact record" {
		t.Fatalf("unexpected steering: %q", req.opts.CommonOptions.Steering)
	}
	if req.opts.CommonOptions.Threshold != 0.9 {
		t.Fatalf("expected threshold 0.9, got %v", req.opts.CommonOptions.Threshold)
	}
	if req.opts.CommonOptions.Context != ctx {
		t.Fatal("expected context to be preserved")
	}
	if req.opts.CommonOptions.RequestID != "req-123" {
		t.Fatalf("unexpected request id: %q", req.opts.CommonOptions.RequestID)
	}
	if req.opts.AllowPartial {
		t.Fatal("expected partial extraction to be disabled")
	}
	if req.opts.SchemaHints["name"] != "Full legal name" {
		t.Fatalf("unexpected schema hint: %#v", req.opts.SchemaHints)
	}
}

func TestGenerateAndTransformBuildersAreConfigurable(t *testing.T) {
	genReq := Generating[string]("write a release note").
		Creative().
		Quick().
		Style("concise").
		Count(2)

	if genReq.opts.CommonOptions.Mode != Creative {
		t.Fatalf("expected creative mode, got %v", genReq.opts.CommonOptions.Mode)
	}
	if genReq.opts.CommonOptions.Intelligence != Quick {
		t.Fatalf("expected quick intelligence, got %v", genReq.opts.CommonOptions.Intelligence)
	}
	if genReq.opts.Style != "concise" {
		t.Fatalf("unexpected style: %q", genReq.opts.Style)
	}
	if genReq.opts.Count != 2 {
		t.Fatalf("unexpected count: %d", genReq.opts.Count)
	}

	type source struct{ Name string }
	type target struct{ Label string }

	transformReq := Transforming[source, target](source{Name: "foo"}).
		Strict().
		Fast().
		Merge("merge").
		Steer("preserve identifiers")

	if transformReq.opts.CommonOptions.Mode != Strict {
		t.Fatalf("expected strict mode, got %v", transformReq.opts.CommonOptions.Mode)
	}
	if transformReq.opts.CommonOptions.Intelligence != Fast {
		t.Fatalf("expected fast intelligence, got %v", transformReq.opts.CommonOptions.Intelligence)
	}
	if transformReq.opts.MergeStrategy != "merge" {
		t.Fatalf("unexpected merge strategy: %q", transformReq.opts.MergeStrategy)
	}
	if transformReq.opts.CommonOptions.Steering != "preserve identifiers" {
		t.Fatalf("unexpected steering: %q", transformReq.opts.CommonOptions.Steering)
	}
}

func TestCollectionBuildersStayFluent(t *testing.T) {
	chooseReq := Choosing([]string{"a", "b"}).
		By("best quality", "lowest cost").
		Smart().
		Reasoning(false).
		Top(2).
		Steer("prefer durable options")

	if len(chooseReq.opts.Criteria) != 2 {
		t.Fatalf("expected 2 criteria, got %d", len(chooseReq.opts.Criteria))
	}
	if chooseReq.opts.CommonOptions.Intelligence != Smart {
		t.Fatalf("expected smart intelligence, got %v", chooseReq.opts.CommonOptions.Intelligence)
	}
	if chooseReq.opts.RequireReasoning {
		t.Fatal("expected reasoning to be disabled")
	}
	if chooseReq.opts.TopN != 2 {
		t.Fatalf("expected top 2, got %d", chooseReq.opts.TopN)
	}
	if chooseReq.opts.CommonOptions.Steering != "prefer durable options" {
		t.Fatalf("unexpected steering: %q", chooseReq.opts.CommonOptions.Steering)
	}

	filterReq := Filtering([]string{"a", "b"}).
		By("urgent only").
		Quick().
		KeepMatching(false).
		MinConfidence(0.95).
		Steer("drop backlog items")

	if filterReq.opts.Criteria != "urgent only" {
		t.Fatalf("unexpected criteria: %q", filterReq.opts.Criteria)
	}
	if filterReq.opts.CommonOptions.Intelligence != Quick {
		t.Fatalf("expected quick intelligence, got %v", filterReq.opts.CommonOptions.Intelligence)
	}
	if filterReq.opts.KeepMatching {
		t.Fatal("expected keep matching to be false")
	}
	if filterReq.opts.MinConfidence != 0.95 {
		t.Fatalf("unexpected confidence: %v", filterReq.opts.MinConfidence)
	}
	if filterReq.opts.CommonOptions.Steering != "drop backlog items" {
		t.Fatalf("unexpected steering: %q", filterReq.opts.CommonOptions.Steering)
	}

	sortReq := Sorting([]string{"a", "b"}).
		By("highest priority first").
		Fast().
		Desc().
		Steer("prioritize deadlines")

	if sortReq.opts.Criteria != "highest priority first" {
		t.Fatalf("unexpected criteria: %q", sortReq.opts.Criteria)
	}
	if sortReq.opts.CommonOptions.Intelligence != Fast {
		t.Fatalf("expected fast intelligence, got %v", sortReq.opts.CommonOptions.Intelligence)
	}
	if sortReq.opts.Direction != "descending" {
		t.Fatalf("unexpected direction: %q", sortReq.opts.Direction)
	}
	if sortReq.opts.CommonOptions.Steering != "prioritize deadlines" {
		t.Fatalf("unexpected steering: %q", sortReq.opts.CommonOptions.Steering)
	}
}

func TestCollectionOptionTypesExposeCommonBuilders(t *testing.T) {
	ctx := context.Background()

	chooseOpts := NewChooseOptions().
		WithMode(Strict).
		WithIntelligence(Smart).
		WithSteering("best fit").
		WithThreshold(0.8).
		WithContext(ctx).
		WithRequestID("choose-1")
	if chooseOpts.CommonOptions.Mode != Strict || chooseOpts.CommonOptions.Intelligence != Smart || chooseOpts.CommonOptions.Steering != "best fit" || chooseOpts.CommonOptions.Threshold != 0.8 || chooseOpts.CommonOptions.Context != ctx || chooseOpts.CommonOptions.RequestID != "choose-1" {
		t.Fatalf("choose options lost common builder state: %#v", chooseOpts)
	}

	filterOpts := NewFilterOptions().
		WithMode(Strict).
		WithIntelligence(Quick).
		WithSteering("keep only compliant").
		WithThreshold(0.7).
		WithContext(ctx).
		WithRequestID("filter-1")
	if filterOpts.CommonOptions.Mode != Strict || filterOpts.CommonOptions.Intelligence != Quick || filterOpts.CommonOptions.Steering != "keep only compliant" || filterOpts.CommonOptions.Threshold != 0.7 || filterOpts.CommonOptions.Context != ctx || filterOpts.CommonOptions.RequestID != "filter-1" {
		t.Fatalf("filter options lost common builder state: %#v", filterOpts)
	}

	sortOpts := NewSortOptions().
		WithMode(Strict).
		WithIntelligence(Fast).
		WithSteering("latest first").
		WithThreshold(0.6).
		WithContext(ctx).
		WithRequestID("sort-1")
	if sortOpts.CommonOptions.Mode != Strict || sortOpts.CommonOptions.Intelligence != Fast || sortOpts.CommonOptions.Steering != "latest first" || sortOpts.CommonOptions.Threshold != 0.6 || sortOpts.CommonOptions.Context != ctx || sortOpts.CommonOptions.RequestID != "sort-1" {
		t.Fatalf("sort options lost common builder state: %#v", sortOpts)
	}
}

func TestExtendedBuildersExposeFluentConfiguration(t *testing.T) {
	classifyReq := Classifying[string, string]("refund requested").
		Categories("billing", "support").
		Smart().
		Steer("prefer the most actionable label")
	if len(classifyReq.opts.Categories) != 2 {
		t.Fatalf("expected categories to be captured, got %#v", classifyReq.opts.Categories)
	}
	if classifyReq.opts.CommonOptions.Intelligence != Smart {
		t.Fatalf("expected smart intelligence, got %v", classifyReq.opts.CommonOptions.Intelligence)
	}

	parseReq := Parsing[map[string]any]("name|john").
		AllowLLMFallback(true).
		AutoFix(true).
		Quick().
		RequestID("parse-1")
	if !parseReq.opts.AllowLLMFallback || !parseReq.opts.AutoFix {
		t.Fatalf("expected parse flags to be enabled: %#v", parseReq.opts)
	}
	if parseReq.opts.OpOptions.Intelligence != Quick || parseReq.opts.OpOptions.RequestID != "parse-1" {
		t.Fatalf("expected parse op options to be updated: %#v", parseReq.opts.OpOptions)
	}

	questionReq := Asking[string, string]("quarterly report", "What changed?").
		Strict().
		Steer("answer briefly")
	if questionReq.opts.Question != "What changed?" {
		t.Fatalf("unexpected question: %q", questionReq.opts.Question)
	}
	if questionReq.opts.CommonOptions.Mode != Strict || questionReq.opts.CommonOptions.Steering != "answer briefly" {
		t.Fatalf("unexpected question common options: %#v", questionReq.opts.CommonOptions)
	}
}

func TestDirectStyleBuildersExposeUnifiedControls(t *testing.T) {
	resolveReq := Resolving([]string{"a", "b"}).
		Strategy("merge").
		Smart().
		Steer("prefer the most complete record")
	if resolveReq.opts.Strategy != "merge" {
		t.Fatalf("unexpected resolve strategy: %q", resolveReq.opts.Strategy)
	}
	if resolveReq.opts.Intelligence != Smart || resolveReq.opts.Steering != "prefer the most complete record" {
		t.Fatalf("unexpected resolve direct options: %#v", resolveReq.opts)
	}

	projectReq := Projecting[map[string]any, map[string]any](map[string]any{"id": 1}).
		Exclude("secret", "token").
		Fast()
	if len(projectReq.opts.Exclude) != 2 {
		t.Fatalf("expected projected exclude fields, got %#v", projectReq.opts.Exclude)
	}
	if projectReq.opts.Intelligence != Fast {
		t.Fatalf("expected fast intelligence, got %v", projectReq.opts.Intelligence)
	}
}
