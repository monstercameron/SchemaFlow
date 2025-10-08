# Transform Example - Resume to CV

## What This Does

Demonstrates the **Transform** operation: converting data from one type to another with semantic understanding.

This example transforms:
- **Input**: Structured resume data (JSON/struct)
- **Output**: Professional markdown-formatted CV

## Use Case

**Real-World Application**: Generate customized CVs, reports, or documents from structured database records. Transform data between different formats (JSON → Markdown, CSV → HTML, etc.)

## How It Works

```go
cv, err := schemaflow.Transform[Resume, MarkdownCV](
    resume,
    ops.NewTransformOptions().WithSteering("Create professional CV..."),
)
```

The LLM intelligently:
1. Understands the input structure
2. Applies professional formatting rules
3. Reorganizes information for clarity
4. Adds appropriate markdown syntax
5. Enhances readability

## Running the Example

```bash
# Set your API key
export OPENAI_API_KEY='your-key-here'

# Run the example
go run main.go
```

## Expected Output

```
📄 Transform Example - Resume to CV
==================================================

📥 Input Resume (Structured):
  Name:   Jane Developer
  Email:  jane.dev@email.com
  Skills: [Go Python JavaScript Docker Kubernetes AWS]

✅ Transformed CV (Markdown):
---
# Jane Developer

📧 jane.dev@email.com | 📱 +1-555-0123

## Skills
- **Languages**: Go, Python, JavaScript
- **DevOps**: Docker, Kubernetes, AWS

## Experience

### Senior Engineer @ Tech Corp
*2020-2024*
...

✨ Success! Structured data → Formatted document
```

## Key Features Demonstrated

- ✅ **Type Transformation**: Resume → MarkdownCV
- ✅ **Intelligent Formatting**: Professional layout
- ✅ **Semantic Understanding**: Groups related info
- ✅ **Style Control**: Via steering parameter

## Learn More

- [SchemaFlow Documentation](../../README.md)
- [Transform API Reference](../../docs/reference/API.md#transform)
