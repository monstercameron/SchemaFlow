package tools

import (
	"context"
	"testing"
)

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	// Test registering a tool
	testTool := &Tool{
		Name:        "test_tool",
		Description: "A test tool",
		Category:    CategoryComputation,
		Parameters: ObjectSchema(map[string]ParameterSchema{
			"input": StringParam("Test input"),
		}, []string{"input"}),
		Execute: func(ctx context.Context, params map[string]any) (Result, error) {
			input := params["input"].(string)
			return NewResult("processed: " + input), nil
		},
	}

	err := registry.Register(testTool)
	if err != nil {
		t.Fatalf("Failed to register tool: %v", err)
	}

	// Test duplicate registration
	err = registry.Register(testTool)
	if err == nil {
		t.Error("Expected error for duplicate registration")
	}

	// Test get
	tool, ok := registry.Get("test_tool")
	if !ok {
		t.Fatal("Tool not found")
	}
	if tool.Name != "test_tool" {
		t.Errorf("Expected name 'test_tool', got '%s'", tool.Name)
	}

	// Test execute
	result, err := registry.Execute(context.Background(), "test_tool", map[string]any{
		"input": "hello",
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !result.Success {
		t.Error("Expected success")
	}
	if result.Data != "processed: hello" {
		t.Errorf("Expected 'processed: hello', got '%v'", result.Data)
	}

	// Test list
	tools := registry.List()
	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}

	// Test list by category
	compTools := registry.ListByCategory(CategoryComputation)
	if len(compTools) != 1 {
		t.Errorf("Expected 1 computation tool, got %d", len(compTools))
	}

	httpTools := registry.ListByCategory(CategoryHTTP)
	if len(httpTools) != 0 {
		t.Errorf("Expected 0 HTTP tools, got %d", len(httpTools))
	}
}

func TestResult(t *testing.T) {
	// Test NewResult
	r := NewResult("data")
	if !r.Success {
		t.Error("Expected success")
	}
	if r.Data != "data" {
		t.Errorf("Expected 'data', got '%v'", r.Data)
	}

	// Test NewResultWithMeta
	r = NewResultWithMeta("data", map[string]any{"key": "value"})
	if r.Metadata["key"] != "value" {
		t.Error("Expected metadata")
	}

	// Test ErrorResult
	r = ErrorResult(context.DeadlineExceeded)
	if r.Success {
		t.Error("Expected failure")
	}
	if r.Error == "" {
		t.Error("Expected error message")
	}

	// Test StubResult
	r = StubResult("test stub")
	if !r.Success {
		t.Error("Expected success for stub")
	}
	if r.Metadata["stubbed"] != true {
		t.Error("Expected stubbed metadata")
	}
}

func TestToOpenAIFormat(t *testing.T) {
	registry := NewRegistry()
	_ = registry.Register(&Tool{
		Name:        "get_weather",
		Description: "Get weather for a location",
		Category:    CategoryHTTP,
		Parameters: ObjectSchema(map[string]ParameterSchema{
			"location": StringParam("City name"),
		}, []string{"location"}),
	})

	format := registry.ToOpenAIFormat()
	if len(format) != 1 {
		t.Fatalf("Expected 1 function, got %d", len(format))
	}

	fn := format[0]
	if fn["type"] != "function" {
		t.Errorf("Expected type 'function', got '%v'", fn["type"])
	}

	fnDef := fn["function"].(map[string]any)
	if fnDef["name"] != "get_weather" {
		t.Errorf("Expected name 'get_weather', got '%v'", fnDef["name"])
	}
}

func TestParameterSchemas(t *testing.T) {
	// Test StringParam
	s := StringParam("test")
	if s.Type != "string" {
		t.Errorf("Expected 'string', got '%s'", s.Type)
	}

	// Test NumberParam
	n := NumberParam("test")
	if n.Type != "number" {
		t.Errorf("Expected 'number', got '%s'", n.Type)
	}

	// Test BoolParam
	b := BoolParam("test")
	if b.Type != "boolean" {
		t.Errorf("Expected 'boolean', got '%s'", b.Type)
	}

	// Test EnumParam
	e := EnumParam("test", []string{"a", "b"})
	if len(e.Enum) != 2 {
		t.Errorf("Expected 2 enum values, got %d", len(e.Enum))
	}

	// Test ObjectSchema
	o := ObjectSchema(map[string]ParameterSchema{
		"field": StringParam("test"),
	}, []string{"field"})
	if o["type"] != "object" {
		t.Errorf("Expected 'object', got '%v'", o["type"])
	}
}

func TestEmptyToolName(t *testing.T) {
	registry := NewRegistry()
	err := registry.Register(&Tool{
		Name: "",
	})
	if err == nil {
		t.Error("Expected error for empty tool name")
	}
}

func TestExecuteNonExistent(t *testing.T) {
	registry := NewRegistry()
	_, err := registry.Execute(context.Background(), "nonexistent", nil)
	if err == nil {
		t.Error("Expected error for non-existent tool")
	}
}

func TestDefaultRegistryFunctions(t *testing.T) {
	// Save original registry
	original := DefaultRegistry

	// Create a new registry for testing
	DefaultRegistry = NewRegistry()

	// Cleanup: restore original registry
	defer func() {
		DefaultRegistry = original
	}()

	err := Register(&Tool{
		Name:        "default_test",
		Description: "Test tool",
		Execute: func(ctx context.Context, params map[string]any) (Result, error) {
			return NewResult("ok"), nil
		},
	})
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	tool, ok := Get("default_test")
	if !ok {
		t.Fatal("Tool not found in default registry")
	}
	if tool.Name != "default_test" {
		t.Errorf("Expected 'default_test', got '%s'", tool.Name)
	}

	result, err := Execute(context.Background(), "default_test", nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Data != "ok" {
		t.Errorf("Expected 'ok', got '%v'", result.Data)
	}

	tools := List()
	if len(tools) == 0 {
		t.Error("Expected at least one tool")
	}
}
