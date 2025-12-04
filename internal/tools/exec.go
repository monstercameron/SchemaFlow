package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// ShellTool executes shell commands
var ShellTool = &Tool{
	Name:        "shell",
	Description: "Execute shell commands with timeout and output capture.",
	Category:    CategoryComputation,
	Parameters: ObjectSchema(map[string]ParameterSchema{
		"command": StringParam("Shell command to execute"),
		"args":    {Type: "array", Description: "Command arguments"},
		"dir":     StringParam("Working directory"),
		"timeout": NumberParam("Timeout in seconds (default: 30)"),
	}, []string{"command"}),
	Execute: executeShell,
}

func executeShell(ctx context.Context, params map[string]any) (Result, error) {
	command, _ := params["command"].(string)
	if command == "" {
		return ErrorResultFromError(fmt.Errorf("command is required")), nil
	}

	timeout := 30.0
	if t, ok := params["timeout"].(float64); ok && t > 0 {
		timeout = t
	}

	dir, _ := params["dir"].(string)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Prepare command based on OS
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else if ctx.Err() == context.DeadlineExceeded {
			return NewResultWithMeta(map[string]any{
				"stdout":    stdout.String(),
				"stderr":    stderr.String(),
				"exit_code": -1,
				"error":     "command timed out",
				"duration":  duration.Seconds(),
			}, map[string]any{"timed_out": true}), nil
		} else {
			return ErrorResultFromError(fmt.Errorf("command failed: %w", err)), nil
		}
	}

	return NewResultWithMeta(map[string]any{
		"stdout":    stdout.String(),
		"stderr":    stderr.String(),
		"exit_code": exitCode,
		"duration":  duration.Seconds(),
	}, nil), nil
}

// RunCodeTool executes code snippets (stub - requires sandboxed environment)
var RunCodeTool = &Tool{
	Name:        "run_code",
	Description: "Execute code snippets in sandboxed environment (stub - security considerations)",
	Category:    CategoryComputation,
	IsStub:      true,
	Parameters: ObjectSchema(map[string]ParameterSchema{
		"language": EnumParam("Programming language", []string{"python", "javascript", "go", "ruby", "php"}),
		"code":     StringParam("Code to execute"),
		"timeout":  NumberParam("Timeout in seconds"),
	}, []string{"language", "code"}),
	Execute: func(ctx context.Context, params map[string]any) (Result, error) {
		language, _ := params["language"].(string)
		code, _ := params["code"].(string)
		return NewResultWithMeta(map[string]any{
			"stub":       true,
			"language":   language,
			"code_lines": len(strings.Split(code, "\n")),
			"message":    "Code execution requires sandboxed environment for security",
		}, map[string]any{"stubbed": true}), nil
	},
}

// EmbedTool generates text embeddings (stub - requires embedding model)
var EmbedTool = &Tool{
	Name:        "embed",
	Description: "Generate text embeddings (stub - requires OpenAI or similar API)",
	Category:    CategoryAI,
	IsStub:      true,
	RequiresAuth: true,
	Parameters: ObjectSchema(map[string]ParameterSchema{
		"text":  StringParam("Text to embed"),
		"model": StringParam("Embedding model name"),
	}, []string{"text"}),
	Execute: func(ctx context.Context, params map[string]any) (Result, error) {
		text, _ := params["text"].(string)
		return NewResultWithMeta(map[string]any{
			"stub":        true,
			"text_length": len(text),
			"message":     "Embedding generation requires AI API integration",
		}, map[string]any{"stubbed": true}), nil
	},
}

// SimilarityTool calculates text similarity (stub)
var SimilarityTool = &Tool{
	Name:        "similarity",
	Description: "Calculate semantic similarity between texts (stub - requires embeddings)",
	Category:    CategoryAI,
	IsStub:      true,
	RequiresAuth: true,
	Parameters: ObjectSchema(map[string]ParameterSchema{
		"text1": StringParam("First text"),
		"text2": StringParam("Second text"),
	}, []string{"text1", "text2"}),
	Execute: func(ctx context.Context, params map[string]any) (Result, error) {
		return NewResultWithMeta(map[string]any{
			"stub":    true,
			"message": "Semantic similarity requires embedding model",
		}, map[string]any{"stubbed": true}), nil
	},
}

// SemanticSearchTool performs semantic search (stub)
var SemanticSearchTool = &Tool{
	Name:        "semantic_search",
	Description: "Search documents by semantic meaning (stub - requires vector database)",
	Category:    CategoryAI,
	IsStub:      true,
	RequiresAuth: true,
	Parameters: ObjectSchema(map[string]ParameterSchema{
		"query":     StringParam("Search query"),
		"documents": {Type: "array", Description: "Documents to search"},
		"top_k":     NumberParam("Number of results to return"),
	}, []string{"query", "documents"}),
	Execute: func(ctx context.Context, params map[string]any) (Result, error) {
		query, _ := params["query"].(string)
		return NewResultWithMeta(map[string]any{
			"stub":    true,
			"query":   query,
			"message": "Semantic search requires embedding model and vector database",
		}, map[string]any{"stubbed": true}), nil
	},
}

// ClassifyTool classifies text into categories (stub)
var ClassifyTool = &Tool{
	Name:        "classify",
	Description: "Classify text into predefined categories (stub - requires LLM)",
	Category:    CategoryAI,
	IsStub:      true,
	RequiresAuth: true,
	Parameters: ObjectSchema(map[string]ParameterSchema{
		"text":       StringParam("Text to classify"),
		"categories": {Type: "array", Description: "Possible categories"},
	}, []string{"text", "categories"}),
	Execute: func(ctx context.Context, params map[string]any) (Result, error) {
		text, _ := params["text"].(string)
		return NewResultWithMeta(map[string]any{
			"stub":        true,
			"text_length": len(text),
			"message":     "Classification requires LLM integration",
		}, map[string]any{"stubbed": true}), nil
	},
}

// SentimentTool analyzes text sentiment (stub)
var SentimentTool = &Tool{
	Name:        "sentiment",
	Description: "Analyze text sentiment (stub - requires sentiment model or LLM)",
	Category:    CategoryAI,
	IsStub:      true,
	RequiresAuth: true,
	Parameters: ObjectSchema(map[string]ParameterSchema{
		"text": StringParam("Text to analyze"),
	}, []string{"text"}),
	Execute: func(ctx context.Context, params map[string]any) (Result, error) {
		text, _ := params["text"].(string)
		return NewResultWithMeta(map[string]any{
			"stub":        true,
			"text_length": len(text),
			"message":     "Sentiment analysis requires ML model or LLM",
		}, map[string]any{"stubbed": true}), nil
	},
}

// TranslateTool translates text (stub)
var TranslateTool = &Tool{
	Name:        "translate",
	Description: "Translate text between languages (stub - requires translation API)",
	Category:    CategoryAI,
	IsStub:      true,
	RequiresAuth: true,
	Parameters: ObjectSchema(map[string]ParameterSchema{
		"text":   StringParam("Text to translate"),
		"from":   StringParam("Source language code"),
		"to":     StringParam("Target language code"),
	}, []string{"text", "to"}),
	Execute: func(ctx context.Context, params map[string]any) (Result, error) {
		text, _ := params["text"].(string)
		to, _ := params["to"].(string)
		return NewResultWithMeta(map[string]any{
			"stub":        true,
			"text_length": len(text),
			"target":      to,
			"message":     "Translation requires translation API (Google, DeepL, etc.)",
		}, map[string]any{"stubbed": true}), nil
	},
}

func init() {
	_ = Register(ShellTool)
	_ = Register(RunCodeTool)
	_ = Register(EmbedTool)
	_ = Register(SimilarityTool)
	_ = Register(SemanticSearchTool)
	_ = Register(ClassifyTool)
	_ = Register(SentimentTool)
	_ = Register(TranslateTool)
}
