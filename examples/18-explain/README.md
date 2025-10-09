# 18. Explain - Generate Human Explanations

The Explain operation generates human-readable explanations for complex data or code in simple terms. It can adapt explanations for different audiences (technical, non-technical, children, executives) and provide various output formats.

## Features

- **Multi-Audience Support**: Tailored explanations for different audiences
- **Flexible Depth**: Control explanation detail level (1-4)
- **Multiple Formats**: Paragraph, bullet-points, step-by-step, Q&A, structured
- **Context Awareness**: Include additional context for better explanations
- **Focus Areas**: Target specific aspects (overview, usage, implementation, etc.)

## Usage

```go
// Basic explanation for non-technical audience
explanation, err := ops.Explain(complexData,
    ops.NewExplainOptions().
        WithAudience("non-technical").
        WithDepth(3))

// Technical explanation with implementation details
explanation, err := ops.Explain(codeStructure,
    ops.NewExplainOptions().
        WithAudience("technical").
        WithFocus("implementation"))

// Executive summary in bullet points
explanation, err := ops.Explain(metrics,
    ops.NewExplainOptions().
        WithAudience("executive").
        WithFormat("bullet-points"))
```

## Options

### Audience Types
- `"technical"` - Technical audience with detailed terminology
- `"non-technical"` - General audience, avoid jargon
- `"children"` - Simple explanations for kids (ages 8-12)
- `"executive"` - Business-focused for executives
- `"beginner"` - Complete beginners, start with basics
- `"expert"` - Advanced technical details for domain experts

### Depth Levels
- `1` - High-level overview only
- `2` - Moderate detail with key concepts
- `3` - Good detail with examples
- `4` - Comprehensive with full technical depth

### Output Formats
- `"paragraph"` - Cohesive paragraph explanation
- `"bullet-points"` - Clear bullet point structure
- `"step-by-step"` - Numbered step-by-step breakdown
- `"qa"` - Q&A format with common questions
- `"structured"` - Structured sections with headings

### Focus Areas
- `"overview"` - General overview of the data/code
- `"usage"` - How to use or interact with it
- `"implementation"` - Technical implementation details
- `"benefits"` - Benefits and advantages
- `"limitations"` - Limitations and constraints
- `"examples"` - Usage examples and scenarios

## Result Structure

```go
type ExplainResult struct {
    Explanation string            // The human-readable explanation
    Summary     string            // Brief overview
    KeyPoints   []string          // Important points to remember
    Audience    string            // Target audience for the explanation
    Complexity  string            // "simple", "intermediate", "detailed", "comprehensive"
    Metadata    map[string]any    // Additional explanation metadata
}
```

## Running the Example

```bash
cd examples/18-explain
go run main.go
```

This will demonstrate:
1. Non-technical audience explanation
2. Technical implementation-focused explanation
3. Executive summary in bullet points
4. Beginner-friendly step-by-step explanation
5. Simple data structure explanation for children
6. Explanation metadata display

## Use Cases

- **API Documentation**: Explain complex API responses
- **Data Analysis**: Make data insights accessible to non-technical stakeholders
- **Code Review**: Generate explanations of code functionality
- **User Interfaces**: Create help text and tooltips
- **Educational Content**: Generate learning materials
- **Business Intelligence**: Explain metrics and KPIs to executives
- **Debugging**: Understand complex data structures during development

## Performance Notes

- Uses the configured intelligence level (Smart/Fast/Quick) for explanation generation
- Larger data structures may take longer to analyze and explain
- Consider using lower depth levels for faster results with simpler data
- The operation automatically analyzes data structure complexity