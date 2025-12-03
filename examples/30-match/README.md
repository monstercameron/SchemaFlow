# SemanticMatch Example

Demonstrates the `SemanticMatch` operation for finding best matches between collections.

## What is SemanticMatch?

The `SemanticMatch` operation finds optimal pairings between two collections:
- Semantic matching based on meaning
- Best-fit, one-to-one, or all-matches strategies
- Field-weighted matching
- Entity resolution and deduplication

Note: Named `SemanticMatch` to distinguish from the control-flow `Match` operation.

## Key Features

- Multiple matching strategies
- Configurable similarity threshold
- Field-level weighting
- Partial match support
- Unmatched item tracking

## Usage

```go
opts := ops.NewMatchOptions().
    WithStrategy("best-fit").
    WithThreshold(0.5).
    WithMatchFields([]string{"name", "description"})

result, err := ops.SemanticMatch[Query, Product](queries, products, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithStrategy()` | Match method (best-fit, all-matches, one-to-one, one-to-many) | "best-fit" |
| `WithThreshold()` | Minimum score to consider a match | 0.5 |
| `WithMatchFields()` | Fields to use for matching | all |
| `WithFieldWeights()` | Importance of each field | equal |
| `WithAllowPartial()` | Allow partial matches | true |
| `WithMaxMatches()` | Max matches per source item | 1 |

## Running

```bash
cd examples/30-match
go run main.go
```

## Example Output

```
Match results:

  Query: 'bluetooth audio device'
  Match: Wireless Bluetooth Headphones (Score: 0.92)
  Why: Matches on bluetooth, audio, and wireless keywords

  Query: 'charging accessories'
  Match: USB-C Charging Cable (Score: 0.85)
  Why: Direct match on charging and accessories category

  Query: 'gaming equipment'
  Match: Mechanical Gaming Keyboard (Score: 0.88)
  Why: Strong match on gaming category and peripherals
```

## Convenience Functions

```go
// Match a single item to a collection
matches, err := ops.MatchOne(query, candidates, opts)
```

## Use Cases

- **Search**: Match queries to products/content
- **Recruiting**: Match candidates to job requirements
- **Entity resolution**: Deduplicate records
- **Recommendation**: Match users to items
