package schemaflow

import (
	"context"
	"errors"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Initialize with test API key if available
	apiKey := os.Getenv("SCHEMAFLOW_API_KEY")
	if apiKey == "" {
		// Use a mock client for testing if no API key
		setupMockClient()
	} else {
		Init(apiKey)
	}
	
	// Run tests
	code := m.Run()
	os.Exit(code)
}

// setupMockClient creates a mock client for testing without API calls
func setupMockClient() {
	logger = NewLogger()
	logger.Info("Using mock client for testing")
	// Set the global callLLM to use our mock
	callLLM = mockLLMResponse
}

// mockLLMResponse provides canned responses for tests
func mockLLMResponse(ctx context.Context, systemPrompt, userPrompt string, opts OpOptions) (string, error) {
	// Return mock responses based on the operation type
	if strings.Contains(systemPrompt, "extraction expert") || strings.Contains(systemPrompt, "Extract structured data") {
		return `{"name": "John Doe", "age": 30, "email": "john@example.com"}`, nil
	}
	if strings.Contains(systemPrompt, "transformation expert") {
		return `{"id": 1, "fullName": "John Doe", "yearsOld": 30}`, nil
	}
	if strings.Contains(systemPrompt, "data generation expert") || strings.Contains(systemPrompt, "Generate structured data") {
		return `{"title": "Test Product", "price": 99.99, "inStock": true}`, nil
	}
	if strings.Contains(systemPrompt, "content generator") {
		return "Generated content based on prompt.", nil
	}
	if strings.Contains(systemPrompt, "classification expert") || strings.Contains(systemPrompt, "Classify") {
		// Check the categories to return an appropriate one
		if strings.Contains(systemPrompt, "junior") {
			return "mid", nil // Return mid for the test scenario
		}
		return "positive", nil
	}
	if strings.Contains(systemPrompt, "scoring expert") || strings.Contains(systemPrompt, "Score") {
		return "0.75", nil
	}
	if strings.Contains(systemPrompt, "summarization expert") || strings.Contains(systemPrompt, "Summarize") {
		return "This is a summary of the text.", nil
	}
	if strings.Contains(systemPrompt, "rewriting expert") || strings.Contains(systemPrompt, "Rewrite") {
		return "This is the rewritten text in a formal tone.", nil
	}
	if strings.Contains(systemPrompt, "translation expert") || strings.Contains(systemPrompt, "Translate") {
		return "Ceci est le texte traduit.", nil
	}
	if strings.Contains(systemPrompt, "expansion expert") || strings.Contains(systemPrompt, "Expand") {
		return "This is an expanded version with more details and context.", nil
	}
	if strings.Contains(systemPrompt, "comparison expert") || strings.Contains(systemPrompt, "Compare") {
		return "Item A is better than Item B in quality.", nil
	}
	if strings.Contains(systemPrompt, "similarity detection") || strings.Contains(systemPrompt, "Similar") {
		return "true", nil
	}
	if strings.Contains(systemPrompt, "selection expert") || strings.Contains(systemPrompt, "Choose") {
		return "1", nil
	}
	if strings.Contains(systemPrompt, "filtering expert") || strings.Contains(systemPrompt, "Filter") {
		return "[0, 2]", nil
	}
	if strings.Contains(systemPrompt, "sorting expert") || strings.Contains(systemPrompt, "Sort") {
		return "[2, 0, 1]", nil
	}
	if strings.Contains(systemPrompt, "pattern matching expert") {
		// Check if the input contains "test" and the condition is "test"
		// The userPrompt format is: "Does this input:\n<input>\n\nMatch this condition:\n<condition>"
		if testing.Verbose() {
			log.Printf("Pattern matching: userPrompt=%q", userPrompt)
		}
		if strings.Contains(userPrompt, "test input") && strings.Contains(userPrompt, "Match this condition:\ntest") {
			return "true", nil
		}
		return "false", nil
	}
	return "Mock response", nil
}

// Test structures
type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type Employee struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	YearsOld int    `json:"yearsOld"`
}

type Product struct {
	Title   string  `json:"title"`
	Price   float64 `json:"price"`
	InStock bool    `json:"inStock"`
}

// ============== Core Type Tests ==============

func TestModeString(t *testing.T) {
	tests := []struct {
		mode     Mode
		expected string
	}{
		{Strict, "strict"},
		{TransformMode, "transform"},
		{Creative, "creative"},
		{Mode(99), "unknown"},
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.expected {
				t.Errorf("Mode.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSpeedString(t *testing.T) {
	tests := []struct {
		speed    Speed
		expected string
	}{
		{Smart, "smart"},
		{Fast, "fast"},
		{Quick, "quick"},
		{Speed(99), "unknown"},
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.speed.String(); got != tt.expected {
				t.Errorf("Speed.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOpOptionsDefaults(t *testing.T) {
	opt := applyDefaults([]OpOptions{})
	
	if opt.Threshold != 0.7 {
		t.Errorf("Expected default threshold 0.7, got %f", opt.Threshold)
	}
	
	if opt.Mode != TransformMode {
		t.Errorf("Expected default mode TransformMode, got %v", opt.Mode)
	}
	
	if opt.Intelligence != Smart {
		t.Errorf("Expected default intelligence Smart, got %v", opt.Intelligence)
	}
	
	if opt.context == nil {
		t.Error("Expected context to be set")
	}
	
	if opt.requestID == "" {
		t.Error("Expected requestID to be generated")
	}
}

// ============== Data Operation Tests ==============

func TestExtract(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	tests := []struct {
		name  string
		input any
		opts  ExtractOptions
	}{
		{
			name:  "simple extraction",
			input: "John Doe, 30 years old, john@example.com",
			opts:  NewExtractOptions().WithMode(TransformMode),
		},
		{
			name:  "json input",
			input: `{"name": "John Doe", "age": 30}`,
			opts:  NewExtractOptions().WithMode(Strict),
		},
		{
			name:  "with steering",
			input: "Extract person data",
			opts:  NewExtractOptions().
				WithSteering("Focus on personal information").
				WithMode(TransformMode),
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			person, err := Extract[Person](tt.input, tt.opts)
			if err != nil {
				t.Fatalf("Extract failed: %v", err)
			}
			
			if person.Name != "John Doe" {
				t.Errorf("Expected name 'John Doe', got %s", person.Name)
			}
			
			if person.Age != 30 {
				t.Errorf("Expected age 30, got %d", person.Age)
			}
		})
	}
}

func TestTransform(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	person := Person{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}
	
	employee, err := Transform[Person, Employee](person, NewTransformOptions().
		WithMode(TransformMode).
		WithIntelligence(Fast))
	
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}
	
	if employee.FullName != "John Doe" {
		t.Errorf("Expected FullName 'John Doe', got %s", employee.FullName)
	}
	
	if employee.YearsOld != 30 {
		t.Errorf("Expected YearsOld 30, got %d", employee.YearsOld)
	}
}

func TestGenerate(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	tests := []struct {
		name   string
		prompt string
		opts   GenerateOptions
	}{
		{
			name:   "generate product",
			prompt: "Create a laptop product listing",
			opts:   NewGenerateOptions().WithMode(Creative),
		},
		{
			name:   "generate with constraints",
			prompt: "Generate a test product",
			opts:   NewGenerateOptions().
				WithSteering("Price should be under $100").
				WithMode(TransformMode),
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test structured generation
			product, err := Generate[Product](tt.prompt, tt.opts)
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}
			
			if product.Title == "" {
				t.Error("Expected product title to be set")
			}
			
			// Test string generation
			text, err := Generate[string](tt.prompt, tt.opts)
			if err != nil {
				t.Fatalf("Generate string failed: %v", err)
			}
			
			if text == "" {
				t.Error("Expected generated text to be non-empty")
			}
		})
	}
}

// ============== Text Operation Tests ==============

func TestSummarize(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	input := "This is a long text that needs to be summarized. It contains many details and information that should be condensed."
	
	summary, err := Summarize(input, NewSummarizeOptions().
		WithSteering("Keep it under 50 words").
		WithMode(TransformMode))
	
	if err != nil {
		t.Fatalf("Summarize failed: %v", err)
	}
	
	if summary == "" {
		t.Error("Expected summary to be non-empty")
	}
	
	if !strings.Contains(summary, "summary") {
		t.Errorf("Expected summary to contain 'summary', got: %s", summary)
	}
}

func TestRewrite(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	input := "This is informal text."
	
	rewritten, err := Rewrite(input, NewRewriteOptions().WithMode(TransformMode))
	
	if err != nil {
		t.Fatalf("Rewrite failed: %v", err)
	}
	
	if rewritten == "" {
		t.Error("Expected rewritten text to be non-empty")
	}
	
	if !strings.Contains(rewritten, "formal") {
		t.Errorf("Expected rewritten text to mention 'formal', got: %s", rewritten)
	}
}

func TestTranslate(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	input := "This is the text to translate."
	
	translated, err := Translate(input, NewTranslateOptions().WithTargetLanguage("Spanish").WithMode(TransformMode))
	
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}
	
	if translated == "" {
		t.Error("Expected translated text to be non-empty")
	}
}

func TestExpand(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	input := "Brief text"
	
	expanded, err := Expand(input, NewExpandOptions().WithMode(Creative))
	
	if err != nil {
		t.Fatalf("Expand failed: %v", err)
	}
	
	if expanded == "" {
		t.Error("Expected expanded text to be non-empty")
	}
	
	if len(expanded) <= len(input) {
		t.Error("Expected expanded text to be longer than input")
	}
}

// ============== Analysis Operation Tests ==============

func TestClassify(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	categories := []string{"positive", "negative", "neutral"}
	
	result, err := Classify("This is great!", NewClassifyOptions().WithCategories(categories).WithMode(TransformMode))
	
	if err != nil {
		t.Fatalf("Classify failed: %v", err)
	}
	
	if result != "positive" {
		t.Errorf("Expected classification 'positive', got %s", result)
	}
}

func TestScore(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	input := "Rate this text for quality"
	
	score, err := Score(input, NewScoreOptions().WithMode(TransformMode))
	
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}
	
	if score < 0 || score > 1 {
		t.Errorf("Expected score between 0-1, got %f", score)
	}
	
	if score != 0.75 {
		t.Errorf("Expected score 0.75 from mock, got %f", score)
	}
}

func TestCompare(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	comparison, err := Compare("Item A", "Item B", NewCompareOptions().WithMode(TransformMode))
	
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}
	
	if comparison == "" {
		t.Error("Expected comparison to be non-empty")
	}
	
	if !strings.Contains(comparison, "Item A") || !strings.Contains(comparison, "Item B") {
		t.Error("Expected comparison to mention both items")
	}
}

func TestSimilar(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	similar, err := Similar("Hello world", "Hi world", NewSimilarOptions().WithSimilarityThreshold(0.7))
	
	if err != nil {
		t.Fatalf("Similar failed: %v", err)
	}
	
	if !similar {
		t.Error("Expected texts to be similar")
	}
}

// ============== Collection Operation Tests ==============

func TestChoose(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	options := []string{"Option A", "Option B", "Option C"}
	
	choice, err := Choose(options, NewChooseOptions().WithCriteria([]string{"best"}))
	
	if err != nil {
		t.Fatalf("Choose failed: %v", err)
	}
	
	if choice != options[1] {
		t.Errorf("Expected choice to be Option B, got %s", choice)
	}
}

func TestFilter(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	items := []string{"Important", "Not important", "Critical"}
	
	filtered, err := Filter(items, NewFilterOptions().WithCriteria("keep positive"))
	
	if err != nil {
		t.Fatalf("Filter failed: %v", err)
	}
	
	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered items, got %d", len(filtered))
	}
}

func TestSort(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	items := []string{"Low priority", "High priority", "Medium priority"}
	
	sorted, err := Sort(items, NewSortOptions().WithCriteria("by value"))
	
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}
	
	if len(sorted) != len(items) {
		t.Errorf("Expected %d sorted items, got %d", len(items), len(sorted))
	}
	
	// Mock returns [2, 0, 1] as indices
	if sorted[0] != items[2] {
		t.Errorf("Expected first item to be %s, got %s", items[2], sorted[0])
	}
}

// ============== Control Flow Tests ==============

func TestMatch(t *testing.T) {
	executed := false
	
	// Debug: log what happens
	t.Logf("Running Match with input 'test input' and condition 'test'")
	
	Match("test input",
		When("test", func() {
			t.Log("Matched 'test' condition")
			executed = true
		}),
		When("other", func() {
			t.Error("Should not match 'other'")
		}),
		Otherwise(func() {
			t.Error("Should not reach Otherwise")
		}),
	)
	
	if !executed {
		t.Error("Expected 'test' case to be executed")
	}
}

func TestMatchWithType(t *testing.T) {
	err := ExtractError{Reason: "test error"}
	matched := false
	
	Match(err,
		When(ExtractError{}, func() {
			matched = true
		}),
		Otherwise(func() {
			t.Error("Should match ExtractError type")
		}),
	)
	
	if !matched {
		t.Error("Expected ExtractError to be matched")
	}
}

func TestMatchOtherwise(t *testing.T) {
	executed := false
	
	Match("unknown",
		When("test", func() {
			t.Error("Should not match 'test'")
		}),
		Otherwise(func() {
			executed = true
		}),
	)
	
	if !executed {
		t.Error("Expected Otherwise to be executed")
	}
}

// ============== Error Type Tests ==============

func TestErrorTypes(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{
			name:     "ExtractError",
			err:      ExtractError{Reason: "failed", TargetType: "Person", Confidence: 0.5},
			contains: "extract error",
		},
		{
			name:     "GenerateError",
			err:      GenerateError{Reason: "failed", TargetType: "Product"},
			contains: "generate error",
		},
		{
			name:     "TransformError",
			err:      TransformError{Reason: "failed", FromType: "A", ToType: "B"},
			contains: "transform error",
		},
		{
			name:     "ClassifyError",
			err:      ClassifyError{Reason: "failed", Categories: []string{"a", "b"}},
			contains: "classify error",
		},
		{
			name:     "SummarizeError",
			err:      SummarizeError{Reason: "failed", Length: 100},
			contains: "summarize error",
		},
		{
			name:     "ScoreError",
			err:      ScoreError{Reason: "failed"},
			contains: "score error",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()
			if !strings.Contains(errStr, tt.contains) {
				t.Errorf("Expected error to contain '%s', got: %s", tt.contains, errStr)
			}
		})
	}
}

// ============== Configuration Tests ==============

func TestInit(t *testing.T) {
	// Save original values
	origTimeout := timeout
	origMaxRetries := maxRetries
	origDebugMode := debugMode
	defer func() {
		timeout = origTimeout
		maxRetries = origMaxRetries
		debugMode = origDebugMode
	}()
	
	// Test with environment variables
	os.Setenv("SCHEMAFLOW_TIMEOUT", "10s")
	os.Setenv("SCHEMAFLOW_MAX_RETRIES", "5")
	os.Setenv("SCHEMAFLOW_DEBUG", "true")
	defer func() {
		os.Unsetenv("SCHEMAFLOW_TIMEOUT")
		os.Unsetenv("SCHEMAFLOW_MAX_RETRIES")
		os.Unsetenv("SCHEMAFLOW_DEBUG")
	}()
	
	Init("")
	
	if timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10s, got %v", timeout)
	}
	
	if maxRetries != 5 {
		t.Errorf("Expected maxRetries to be 5, got %d", maxRetries)
	}
	
	if !debugMode {
		t.Error("Expected debug mode to be enabled")
	}
}

func TestGetModel(t *testing.T) {
	tests := []struct {
		speed    Speed
		expected string
	}{
		{Smart, "gpt-4-turbo-preview"},
		{Fast, "gpt-3.5-turbo"},
		{Quick, "gpt-3.5-turbo"},
	}
	
	for _, tt := range tests {
		t.Run(tt.speed.String(), func(t *testing.T) {
			if got := getModel(tt.speed); got != tt.expected {
				t.Errorf("getModel(%v) = %v, want %v", tt.speed, got, tt.expected)
			}
		})
	}
}

func TestGetTemperature(t *testing.T) {
	tests := []struct {
		mode     Mode
		expected float32
	}{
		{Strict, 0.1},
		{TransformMode, 0.3},
		{Creative, 0.7},
	}
	
	for _, tt := range tests {
		t.Run(tt.mode.String(), func(t *testing.T) {
			if got := getTemperature(tt.mode); got != tt.expected {
				t.Errorf("getTemperature(%v) = %v, want %v", tt.mode, got, tt.expected)
			}
		})
	}
}

// ============== Utility Tests ==============

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()
	
	if id1 == "" {
		t.Error("Expected request ID to be non-empty")
	}
	
	if id1 == id2 {
		t.Error("Expected unique request IDs")
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"rate limit error", errors.New("rate limit exceeded"), true},
		{"timeout error", errors.New("request timeout"), true},
		{"429 error", errors.New("status 429"), true},
		{"503 error", errors.New("status 503"), true},
		{"connection error", errors.New("connection refused"), true},
		{"non-retryable", errors.New("invalid input"), false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRetryableError(tt.err); got != tt.expected {
				t.Errorf("isRetryableError(%v) = %v, want %v", tt.err, got, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{1, 2, 1},
		{2, 1, 1},
		{1, 1, 1},
		{-1, 0, -1},
	}
	
	for _, tt := range tests {
		result := min(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

// ============== Steering Preset Tests ==============

func TestSteeringPresets(t *testing.T) {
	tests := []struct {
		name     string
		preset   func(...string) string
		contains string
	}{
		{"BusinessTone", Steering.BusinessTone, "professional"},
		{"CasualTone", Steering.CasualTone, "friendly"},
		{"TechnicalTone", Steering.TechnicalTone, "technical"},
		{"UrgencyScore", Steering.UrgencyScore, "urgency"},
		{"ImportanceScore", Steering.ImportanceScore, "importance"},
		{"QualityScore", Steering.QualityScore, "quality"},
		{"PrioritySort", Steering.PrioritySort, "priority"},
		{"WorkContext", Steering.WorkContext, "work"},
		{"HomeContext", Steering.HomeContext, "home"},
		{"MobileContext", Steering.MobileContext, "mobile"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.preset()
			if !strings.Contains(strings.ToLower(result), tt.contains) {
				t.Errorf("Expected %s preset to contain '%s', got: %s", tt.name, tt.contains, result)
			}
			
			// Test with additional context
			resultWithContext := tt.preset("Additional context")
			if !strings.Contains(resultWithContext, "Additional Context") {
				t.Error("Expected preset to include additional context")
			}
		})
	}
}

// ============== Logger Tests ==============

func TestLogger(t *testing.T) {
	// Create a test logger
	logger := NewLogger()
	
	// Test different log levels
	logger.SetLevel(DebugLevel)
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	
	// Test with fields
	logger.WithFields(map[string]any{
		"requestID": "test-123",
	}).Info("message with request ID")
	
	// Test JSON output
	logger.SetJSONMode(true)
	logger.Info("json message", "key", "value")
	
	// Verify logger doesn't panic
	logger.WithFields(map[string]any{
		"operation": "test",
		"duration":  time.Second,
	}).Info("structured log")
}

// ============== Parsing Tests ==============

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Person
		wantErr  bool
	}{
		{
			name:     "valid json",
			input:    `{"name": "John", "age": 30, "email": "john@example.com"}`,
			expected: Person{Name: "John", Age: 30, Email: "john@example.com"},
			wantErr:  false,
		},
		{
			name:     "json with markdown",
			input:    "```json\n{\"name\": \"John\", \"age\": 30, \"email\": \"john@example.com\"}\n```",
			expected: Person{Name: "John", Age: 30, Email: "john@example.com"},
			wantErr:  false,
		},
		{
			name:     "json with generic markdown",
			input:    "```\n{\"name\": \"John\", \"age\": 30, \"email\": \"john@example.com\"}\n```",
			expected: Person{Name: "John", Age: 30, Email: "john@example.com"},
			wantErr:  false,
		},
		{
			name:     "invalid json",
			input:    "not json",
			expected: Person{},
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Person
			err := parseJSON(tt.input, &result)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ============== Context Tests ==============

func TestContextTimeout(t *testing.T) {
	// Replace LLM call with a slow mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	
	callLLM = func(ctx context.Context, systemPrompt, userPrompt string, opts OpOptions) (string, error) {
		select {
		case <-time.After(100 * time.Millisecond):
			return "response", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	
	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	
	opts := OpOptions{
		context: ctx,
		Mode: TransformMode,
	}
	
	// This should timeout
	_, err := callLLM(ctx, "test", "test", opts)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

// ============== Benchmark Tests ==============

func BenchmarkExtract(b *testing.B) {
	// Use mock for consistent benchmarking
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	input := "John Doe, 30 years old"
	opts := NewExtractOptions().WithMode(TransformMode)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Extract[Person](input, opts)
	}
}

func BenchmarkGenerateRequestID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateRequestID()
	}
}

func BenchmarkParseJSON(b *testing.B) {
	input := `{"name": "John", "age": 30, "email": "john@example.com"}`
	var result Person
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseJSON(input, &result)
	}
}

// ============== Integration Tests ==============

func TestOperationChaining(t *testing.T) {
	// Replace LLM call with mock
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	callLLM = mockLLMResponse
	
	// Chain: Extract -> Transform -> Score
	input := "John Doe, 30 years old, john@example.com"
	
	// Step 1: Extract
	person, err := Extract[Person](input, NewExtractOptions().WithMode(TransformMode))
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}
	
	// Step 2: Transform
	employee, err := Transform[Person, Employee](person, NewTransformOptions().WithMode(TransformMode))
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}
	
	// Step 3: Score
	score, err := Score(employee, NewScoreOptions().WithMode(TransformMode))
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}
	
	if score <= 0 {
		t.Error("Expected positive score for chained operations")
	}
}

func TestErrorRecovery(t *testing.T) {
	// Test error recovery with Match
	oldCallLLM := callLLM
	defer func() { callLLM = oldCallLLM }()
	
	// Mock that returns an error
	callLLM = func(ctx context.Context, systemPrompt, userPrompt string, opts OpOptions) (string, error) {
		return "", errors.New("LLM error")
	}
	
	recovered := false
	_, err := Extract[Person]("test input", NewExtractOptions().WithMode(TransformMode))
	
	Match(err,
		When(ExtractError{}, func() {
			recovered = true
		}),
		Otherwise(func() {
			t.Error("Expected ExtractError")
		}),
	)
	
	if !recovered {
		t.Error("Expected error recovery through Match")
	}
}