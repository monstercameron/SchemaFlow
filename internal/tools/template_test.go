package tools

import (
	"context"
	"strings"
	"testing"
)

func TestTemplateToolBasic(t *testing.T) {
	result, _ := TemplateTool.Execute(context.Background(), map[string]any{
		"template": "Hello, {{.Name}}!",
		"data":     map[string]any{"Name": "World"},
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	output := result.Data.(string)
	if output != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got %q", output)
	}
}

func TestTemplateToolWithMap(t *testing.T) {
	result, _ := TemplateTool.Execute(context.Background(), map[string]any{
		"template": "{{.Name}} is {{.Age}} years old",
		"data": map[string]any{
			"Name": "Alice",
			"Age":  30,
		},
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	output := result.Data.(string)
	if output != "Alice is 30 years old" {
		t.Errorf("Expected 'Alice is 30 years old', got %q", output)
	}
}

func TestTemplateToolHTML(t *testing.T) {
	result, _ := TemplateTool.Execute(context.Background(), map[string]any{
		"template": "<p>{{.Content}}</p>",
		"data":     map[string]any{"Content": "<script>alert('xss')</script>"},
		"html":     true,
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	output := result.Data.(string)
	// HTML mode should escape the script tag
	if strings.Contains(output, "<script>") {
		t.Error("Expected HTML escaping")
	}
}

func TestTemplateToolRange(t *testing.T) {
	result, _ := TemplateTool.Execute(context.Background(), map[string]any{
		"template": "{{range .Items}}{{.}} {{end}}",
		"data":     map[string]any{"Items": []string{"a", "b", "c"}},
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	output := result.Data.(string)
	if output != "a b c " {
		t.Errorf("Expected 'a b c ', got %q", output)
	}
}

func TestTemplateToolInvalid(t *testing.T) {
	result, _ := TemplateTool.Execute(context.Background(), map[string]any{
		"template": "{{.Broken",
		"data":     map[string]any{},
	})

	if result.Success {
		t.Error("Expected failure for invalid template")
	}
}

func TestTemplateToolMissingTemplate(t *testing.T) {
	result, _ := TemplateTool.Execute(context.Background(), map[string]any{
		"data": map[string]any{},
	})

	if result.Success {
		t.Error("Expected failure for missing template")
	}
}

func TestStringTemplateTool(t *testing.T) {
	result, _ := StringTemplateTool.Execute(context.Background(), map[string]any{
		"template": "Hello, {{name}}! You have {{count}} messages.",
		"values":   map[string]any{"name": "Alice", "count": 5},
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	output := result.Data.(string)
	if output != "Hello, Alice! You have 5 messages." {
		t.Errorf("Expected interpolated string, got %q", output)
	}
}

func TestStringTemplateToolMissingKey(t *testing.T) {
	result, _ := StringTemplateTool.Execute(context.Background(), map[string]any{
		"template": "Hello, {{name}}! Status: {{status}}",
		"values":   map[string]any{"name": "Alice"},
	})

	// Should succeed but leave unmatched placeholders
	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	output := result.Data.(string)
	if !strings.Contains(output, "{{status}}") {
		t.Error("Expected unmatched placeholder to remain")
	}
}

func TestMarkdownTool(t *testing.T) {
	md := `# Title
## Subtitle
- Item 1
- Item 2

Regular paragraph.`

	result, _ := MarkdownTool.Execute(context.Background(), map[string]any{
		"markdown": md,
	})

	if !result.Success {
		t.Fatalf("Expected success: %s", result.Error)
	}

	html := result.Data.(string)
	if !strings.Contains(html, "<h1>Title</h1>") {
		t.Error("Expected h1 tag")
	}
	if !strings.Contains(html, "<h2>Subtitle</h2>") {
		t.Error("Expected h2 tag")
	}
	if !strings.Contains(html, "<li>") {
		t.Error("Expected list items")
	}
}

func TestMarkdownToHTMLHelpers(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{"# Header", "<h1>Header</h1>"},
		{"## Subheader", "<h2>Subheader</h2>"},
		{"### Small", "<h3>Small</h3>"},
		{"- Item", "<li>Item</li>"},
		{"* Item", "<li>Item</li>"},
		{"Text", "<p>Text</p>"},
	}

	for _, tt := range tests {
		result := markdownToHTML(tt.input)
		if !strings.Contains(result, tt.contains) {
			t.Errorf("markdownToHTML(%q) should contain %q, got %q", tt.input, tt.contains, result)
		}
	}
}
