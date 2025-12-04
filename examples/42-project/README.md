# Project Example

This example demonstrates the `Project` operation, which transforms data structure while preserving semantics.

## What it does

The `Project` operation transforms data from one type to another, similar to SQL SELECT with column mapping. It maps fields semantically, handling field renames, exclusions, and type conversions.

## Use Cases

- **API responses**: Project internal models to public DTOs
- **Data views**: Create different views of the same data
- **Privacy filtering**: Exclude sensitive fields
- **Schema migration**: Transform old formats to new ones

## Key Features

- Field mapping (explicit or inferred)
- Field exclusion for privacy/security
- Documents all mappings and transformations
- Tracks lost fields and inferred values

## Running the Example

```bash
cd examples/42-project
go run main.go
```
