// Package tools registry provides helper functions for tool management.
// Individual tool files register themselves via init() functions.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// ToolHandler is a handler function for tool calls from LLMs
type ToolHandler func(ctx context.Context, name string, args map[string]any) (string, error)

// CreateToolHandler creates a handler for processing tool calls using the default registry
func CreateToolHandler() ToolHandler {
	return CreateToolHandlerWithRegistry(DefaultRegistry)
}

// CreateToolHandlerWithRegistry creates a handler for processing tool calls using a specific registry
func CreateToolHandlerWithRegistry(registry *Registry) ToolHandler {
	return func(ctx context.Context, name string, args map[string]any) (string, error) {
		result, err := registry.Execute(ctx, name, args)
		if err != nil {
			return "", err
		}

		// Convert result to JSON for LLM consumption
		output, err := json.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("failed to marshal result: %w", err)
		}

		return string(output), nil
	}
}

// ToolSpec represents a tool specification for OpenAI function calling
type ToolSpec struct {
	Type     string       `json:"type"`
	Function FunctionSpec `json:"function"`
}

// FunctionSpec represents a function specification
type FunctionSpec struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// GetOpenAITools returns tool specifications in OpenAI function calling format
func GetOpenAITools() []ToolSpec {
	return GetOpenAIToolsFromRegistry(DefaultRegistry)
}

// GetOpenAIToolsFromRegistry returns tool specifications from a specific registry
func GetOpenAIToolsFromRegistry(registry *Registry) []ToolSpec {
	tools := registry.List()
	specs := make([]ToolSpec, len(tools))

	for i, tool := range tools {
		specs[i] = ToolSpec{
			Type: "function",
			Function: FunctionSpec{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
		}
	}

	return specs
}

// GetAnthropicTools returns tool specifications in Anthropic format
func GetAnthropicTools() []map[string]any {
	return GetAnthropicToolsFromRegistry(DefaultRegistry)
}

// GetAnthropicToolsFromRegistry returns tool specifications from a specific registry
func GetAnthropicToolsFromRegistry(registry *Registry) []map[string]any {
	tools := registry.List()
	specs := make([]map[string]any, len(tools))

	for i, tool := range tools {
		specs[i] = map[string]any{
			"name":         tool.Name,
			"description":  tool.Description,
			"input_schema": tool.Parameters,
		}
	}

	return specs
}

// ToolCategories returns a map of tool categories to tool names
func ToolCategories() map[string][]string {
	return map[string][]string{
		"computation": {"calculate", "convert", "regex"},
		"http":        {"fetch", "post", "web_search", "scrape", "browser", "webhook", "encode_url", "build_url"},
		"file":        {"read_file", "write_file", "list_dir", "copy_file", "move_file", "delete_file", "file_exists", "file_info", "watch_file", "search_files"},
		"database":    {"sqlite", "migrate", "seed", "backup", "vector_db"},
		"cache":       {"cache", "memoize"},
		"security":    {"hash", "base64", "token", "encrypt", "decrypt"},
		"time":        {"now", "parse_time", "duration", "schedule", "holiday", "geo", "weather"},
		"data":        {"csv", "json", "xml", "table", "diff"},
		"finance":     {"chart", "currency", "stock", "tax", "interest"},
		"messaging":   {"email", "sms", "push", "slack", "discord", "webhook_notify"},
		"image":       {"vision", "ocr", "image_info", "image_resize", "image_crop", "image_convert", "image_base64", "thumbnail"},
		"audio":       {"tts", "stt", "audio_info", "audio_convert", "audio_trim", "audio_analyze"},
		"template":    {"template", "string_template", "markdown"},
		"archive":     {"pdf", "zip", "tar", "qrcode", "barcode"},
		"execution":   {"shell", "run_code"},
		"ai":          {"embed", "similarity", "semantic_search", "classify", "sentiment", "translate"},
	}
}

// ToolsByCategory returns tools for a specific category
func ToolsByCategory(category string) []*Tool {
	names, ok := ToolCategories()[category]
	if !ok {
		return nil
	}

	var tools []*Tool
	for _, name := range names {
		if tool, ok := DefaultRegistry.Get(name); ok {
			tools = append(tools, tool)
		}
	}
	return tools
}

// CreateSubRegistry creates a new registry with only selected tools
func CreateSubRegistry(toolNames ...string) *Registry {
	registry := NewRegistry()
	for _, name := range toolNames {
		if tool, ok := DefaultRegistry.Get(name); ok {
			registry.Register(tool)
		}
	}
	return registry
}

// CreateCategoryRegistry creates a new registry with tools from specific categories
func CreateCategoryRegistry(categories ...string) *Registry {
	registry := NewRegistry()
	for _, category := range categories {
		for _, tool := range ToolsByCategory(category) {
			registry.Register(tool)
		}
	}
	return registry
}
