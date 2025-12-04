# Interpolate Example

This example demonstrates the `Interpolate` operation, which fills gaps in typed sequences intelligently.

## What it does

The `Interpolate` operation takes a sequence with missing data points and fills them intelligently using context and patterns from the existing data.

## Use Cases

- **Time series**: Fill missing data points in metrics
- **Data cleaning**: Handle null/missing values intelligently
- **Sequence completion**: Complete partially recorded sequences
- **Historical reconstruction**: Infer missing historical records

## Key Features

- Understands sequence patterns (linear, cyclical, etc.)
- Context-aware gap filling
- Documents which items were interpolated
- Provides confidence for each filled value

## Running the Example

```bash
cd examples/40-interpolate
go run main.go
```
