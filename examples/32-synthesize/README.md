# Synthesize Example

Demonstrates the `Synthesize` operation for combining multiple sources with insights.

## What is Synthesize?

The `Synthesize` operation combines multiple sources into unified output:
- **Merge**: Combine sources into a cohesive whole
- **Compare**: Highlight similarities and differences
- **Integrate**: Create narrative with insights
- **Reconcile**: Resolve conflicting information

## Key Features

- Multiple synthesis strategies
- Conflict detection and resolution
- Source citation
- Insight generation
- Coverage tracking

## Usage

```go
opts := ops.NewSynthesizeOptions().
    WithStrategy("integrate").
    WithCiteSources(true).
    WithGenerateInsights(true)

result, err := ops.Synthesize[Report](sources, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithStrategy()` | Synthesis approach | "integrate" |
| `WithConflictResolution()` | How to resolve conflicts | "llm-decide" |
| `WithSourcePriorities()` | Priority order for conflicts | [] |
| `WithCiteSources()` | Include source references | true |
| `WithGenerateInsights()` | Generate additional insights | true |
| `WithFocusAreas()` | Areas to focus synthesis on | [] |

## Running

```bash
cd examples/32-synthesize
go run main.go
```

## Example Output

```
Synthesized Research Report:

Summary: Multiple studies show mixed results on remote work...

Key Findings:
  - Remote workers show higher individual productivity [Source 0]
  - Hybrid work offers better work-life balance [Source 1]
  - In-office collaboration remains stronger [Source 2]

Consensus: Optimal work arrangement depends on role and task type

Generated Insights:
  [pattern] Productivity gains may offset collaboration challenges
  [gap] Long-term effects on career progression not studied

Conflicts Identified:
  - Productivity metrics: Resolved by considering context
```

## Difference from Merge

- **Merge**: Simple combination of data
- **Synthesize**: Intelligent integration with insights and conflict handling

## Use Cases

- **Research synthesis**: Combine multiple studies
- **Data reconciliation**: Merge conflicting records
- **Document consolidation**: Combine related documents
- **Competitive analysis**: Compare multiple options
