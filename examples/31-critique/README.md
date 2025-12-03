# Critique Example

Demonstrates the `Critique` operation for evaluating with actionable feedback.

## What is Critique?

The `Critique` operation evaluates content with detailed, actionable feedback:
- Scores against specific criteria or rubrics
- Identifies issues with severity levels
- Provides concrete suggestions and fixes
- Includes positive feedback when requested

## Key Features

- Criteria-based evaluation
- Custom rubrics
- Multiple critique styles
- Severity-filtered issues
- Positive feedback inclusion

## Usage

```go
opts := ops.NewCritiqueOptions().
    WithCriteria([]string{"clarity", "accuracy", "completeness"}).
    WithIncludeSuggestions(true).
    WithIncludeFixes(true).
    WithStyle("constructive")

result, err := ops.Critique(content, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithCriteria()` | Evaluation criteria | [] |
| `WithRubric()` | Detailed rubric (criterion -> description) | {} |
| `WithStyle()` | Critique style (constructive, harsh, balanced) | "constructive" |
| `WithIncludeSuggestions()` | Add improvement suggestions | true |
| `WithIncludeFixes()` | Add specific fixes | true |
| `WithIncludePositives()` | Include positive feedback | true |
| `WithSeverityFilter()` | Filter by severity (all, major, minor, critical) | "all" |
| `WithMaxIssues()` | Maximum issues to report | 0 (unlimited) |

## Running

```bash
cd examples/31-critique
go run main.go
```

## Example Output

```
Overall Score: 0.45/1.00

Criteria Scores:
  argument_strength: 0.35
  evidence: 0.30
  clarity: 0.55
  structure: 0.60

Issues Found:

  1. [major] evidence
     Claims lack supporting data or citations
     Suggestion: Add specific statistics and cite scientific sources
     Fix: Replace "Scientists say it's getting warmer" with specific data

  2. [major] argument_strength
     Arguments are vague and lack specificity
     Suggestion: Provide concrete examples and specific proposals

Positive Feedback:
  - structure: Basic essay structure is present with intro and conclusion

Summary: The essay addresses an important topic but lacks depth...
```

## Difference from Score

- **Score**: Simple numeric rating
- **Critique**: Detailed feedback with issues, suggestions, and fixes

## Use Cases

- **Writing review**: Essays, articles, documentation
- **Code review**: Identify issues with suggestions
- **Content quality**: Marketing copy, product descriptions
- **Presentations**: Evaluate and improve slides/talks
