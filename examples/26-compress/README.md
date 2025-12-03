# Compress Example

Demonstrates the `Compress` operation for reducing content while preserving meaning.

## What is Compress?

The `Compress` operation reduces content length while maintaining key information:
- **Extractive**: Selects and combines key sentences
- **Abstractive**: Rewrites content more concisely
- **Hybrid**: Combines both approaches
- Configurable compression ratio

## Key Features

- Multiple compression strategies
- Priority preservation for key topics
- Structure preservation option
- Configurable compression ratio

## Usage

```go
opts := ops.NewCompressOptions().
    WithCompressionRatio(0.3).
    WithStrategy("extractive").
    WithPreserveStructure(true)

result, err := ops.Compress(longText, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithCompressionRatio()` | Target size as fraction of original | 0.5 |
| `WithStrategy()` | Compression approach (extractive, abstractive, hybrid) | "extractive" |
| `WithPriority()` | Topics/elements to preserve | [] |
| `WithPreserveStructure()` | Maintain document structure | false |
| `WithMaxLength()` | Absolute maximum length | 0 |

## Running

```bash
cd examples/26-compress
go run main.go
```

## Example Output

```
Original length: 1247 characters
Compressed length: 374 characters
Actual compression ratio: 0.30

Compressed text:
AI has transformed technology, powering search and autonomous vehicles through
deep learning neural networks. Healthcare benefits from diagnostic AI and
accelerated drug discovery. Ethical concerns about bias and privacy require
attention. AI integration will continue in homes, education, and transportation.
```

## Difference from Summarize

- **Summarize**: Creates a structured summary with sections
- **Compress**: Focuses purely on size reduction

## Use Cases

- **Document compression**: Reduce storage/transmission size
- **Tweet generation**: Fit content into character limits
- **Brief creation**: Create executive summaries
- **Mobile optimization**: Reduce content for mobile displays
