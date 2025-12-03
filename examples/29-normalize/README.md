# Normalize Example

Demonstrates the `Normalize` operation for standardizing data formats.

## What is Normalize?

The `Normalize` operation standardizes data to consistent formats:
- Applies formatting rules (dates, addresses, names)
- Expands abbreviations
- Maps values to canonical forms
- Fixes typos and inconsistencies

## Key Features

- Custom normalization rules
- Canonical value mappings
- Change tracking
- Batch processing
- Typo correction

## Usage

```go
opts := ops.NewNormalizeOptions().
    WithRules(map[string]string{
        "date": "ISO 8601 format",
        "country": "Full country name",
    }).
    WithTrackChanges(true)

result, err := ops.Normalize(record, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithStandard()` | Standard to apply (ISO 8601, etc.) | "" |
| `WithRules()` | Field-specific normalization rules | {} |
| `WithCanonicalMappings()` | Value -> canonical mappings | {} |
| `WithFixTypos()` | Attempt to fix typos | false |
| `WithTrackChanges()` | Record what changed | false |

## Running

```bash
cd examples/29-normalize
go run main.go
```

## Example Output

```
Original:
  123 Main St., new york city, NY, USA 10001

Normalized:
  123 Main Street, New York City, NY, United States 10001

Changes made:
  - street: 'Main St.' -> 'Main Street' (abbreviation expanded)
  - city: 'new york city' -> 'New York City' (proper case)
  - country: 'USA' -> 'United States' (canonical form)
```

## Convenience Functions

```go
// Normalize a single text string
normalized, err := ops.NormalizeText(text, opts)

// Normalize a batch of records
results, err := ops.NormalizeBatch(records, opts)
```

## Use Cases

- **Data cleaning**: Standardize imported data
- **Address verification**: Normalize addresses for shipping
- **Date formatting**: Convert various date formats
- **Text processing**: Expand abbreviations, fix typos
