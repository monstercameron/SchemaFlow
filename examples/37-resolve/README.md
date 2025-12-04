# Resolve Example

This example demonstrates the `Resolve` operation, which resolves conflicts when multiple typed sources disagree.

## What it does

The `Resolve` operation takes multiple sources of the same type that may contain conflicting information and produces a single reconciled result, documenting which conflicts were found and how they were resolved.

## Use Cases

- **Data deduplication**: Merge duplicate records with different values
- **Multi-source aggregation**: Combine data from APIs, databases, user input
- **Version reconciliation**: Merge changes from different versions
- **Truth finding**: Determine the most reliable value from conflicting sources

## Key Features

- Identifies specific conflicts between sources
- Documents resolution strategy for each conflict
- Returns confidence scores
- Preserves source attribution

## Running the Example

```bash
cd examples/37-resolve
go run main.go
```
