# Annotate Example

Demonstrates the `Annotate` operation for adding metadata and labels to content.

## What is Annotate?

The `Annotate` operation adds structured metadata to text or data, including:
- **Named Entity Recognition (NER)**: Identify people, organizations, locations, etc.
- **Sentiment Analysis**: Detect positive, negative, or neutral sentiment
- **Topic Extraction**: Identify main topics and themes
- **Custom Labels**: Apply domain-specific labels

## Key Features

- Multiple annotation types: `ner`, `sentiment`, `topic`, `custom`
- Configurable confidence thresholds
- Optional span/position information
- Hierarchical annotations

## Usage

```go
opts := ops.NewAnnotateOptions().
    WithAnnotationType("ner").
    WithLabels([]string{"person", "organization", "location"}).
    WithIncludeSpans(true).
    WithMinConfidence(0.7)

result, err := ops.Annotate(text, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithAnnotationType()` | Type of annotation (ner, sentiment, topic, custom) | "custom" |
| `WithLabels()` | Labels to use for annotation | [] |
| `WithIncludeSpans()` | Include character positions | false |
| `WithMinConfidence()` | Minimum confidence threshold | 0.5 |
| `WithHierarchical()` | Enable nested annotations | false |
| `WithContext()` | Additional context for annotation | "" |

## Running

```bash
cd examples/23-annotate
go run main.go
```

## Example Output

```
Found 5 annotations:
  - Tim Cook (person): confidence 0.95
    Span: 10-18
  - Apple (organization): confidence 0.98
    Span: 0-5
  - San Jose (location): confidence 0.92
    Span: 65-73
  - California (location): confidence 0.94
    Span: 75-85
  - WWDC (event): confidence 0.88
    Span: 45-49
```
