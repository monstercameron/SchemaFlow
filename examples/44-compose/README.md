# Compose Example

This example demonstrates the `Compose` operation, which builds complex typed objects from multiple parts.

## What it does

The `Compose` operation takes multiple partial data sources and intelligently assembles them into a single coherent typed object, handling conflicts and gaps.

## Use Cases

- **Data assembly**: Combine data from multiple APIs
- **Form aggregation**: Merge multi-step form submissions
- **Document composition**: Build documents from sections
- **Profile building**: Assemble user profiles from various sources

## Key Features

- Tracks source contribution per field
- Resolves conflicts between sources
- Optional gap filling
- Completeness scoring

## Running the Example

```bash
cd examples/44-compose
go run main.go
```
