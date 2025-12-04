# Derive Example

This example demonstrates the `Derive` operation, which infers new typed fields from existing data.

## What it does

The `Derive` operation takes data of one type and produces a new type with additional derived/inferred fields. Unlike simple transformation, it uses LLM intelligence to infer values that aren't directly present in the source.

## Use Cases

- **Data enrichment**: Add computed fields like age brackets, generations
- **Feature engineering**: Derive ML features from raw data
- **Business logic**: Infer customer segments, risk scores
- **Data augmentation**: Add semantic information to raw records

## Key Features

- Type-safe input and output
- Documents how each field was derived
- Tracks derivation confidence per field
- Preserves original data while adding new fields

## Running the Example

```bash
cd examples/38-derive
go run main.go
```
