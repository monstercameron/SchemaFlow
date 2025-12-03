# Enrich Example

Demonstrates the `Enrich` operation for adding derived/inferred fields.

## What is Enrich?

The `Enrich` operation adds new fields to data by deriving or inferring values:
- Adds metadata like categories, tags, summaries
- Infers properties from existing fields
- Fills in missing information intelligently
- Domain-specific enrichment

## Key Features

- Type-safe input/output transformation
- Custom derivation rules
- Domain-aware enrichment
- Confidence scores for inferred values

## Usage

```go
type Product struct {
    Name  string `json:"name"`
    Price float64 `json:"price"`
}

type EnrichedProduct struct {
    Name       string   `json:"name"`
    Price      float64  `json:"price"`
    Category   string   `json:"category"`
    Keywords   []string `json:"keywords"`
}

opts := ops.NewEnrichOptions().
    WithDeriveFields([]string{"category", "keywords"}).
    WithDomain("e-commerce")

result, err := ops.Enrich[Product, EnrichedProduct](product, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithDeriveFields()` | Fields to add/derive | [] |
| `WithDerivationRules()` | Custom rules per field | {} |
| `WithDomain()` | Domain context for inference | "" |
| `WithIncludeConfidence()` | Add confidence scores | false |
| `WithOverwriteExisting()` | Overwrite non-empty fields | false |

## Running

```bash
cd examples/28-enrich
go run main.go
```

## Example Output

```
Original product: Pro Gaming Keyboard RGB
Enriched with:
  Category: Electronics > Computer Peripherals > Keyboards
  Keywords: [gaming, mechanical, RGB, keyboard, Cherry MX]
  Target Market: Gamers, Enthusiasts
  Price Range: Premium
```

## Difference from Transform

- **Transform**: Changes data structure/format
- **Enrich**: Adds new derived fields while preserving originals

## Use Cases

- **Product catalogs**: Add categories, tags, SEO keywords
- **Contact data**: Infer company size, industry, region
- **Documents**: Add summaries, reading time, difficulty
- **Events**: Add time zones, weather, related events
