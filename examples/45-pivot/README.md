# Pivot Example

This example demonstrates the `Pivot` operation, which restructures data relationships between typed objects.

## What it does

The `Pivot` operation transforms data structure - pivoting rows to columns, flattening nested objects, or grouping data in new ways.

## Use Cases

- **Reporting**: Pivot time series data for charts
- **Data flattening**: Convert nested structures to flat records
- **Aggregation**: Group and summarize data
- **Schema transformation**: Restructure for different systems

## Key Features

- Row-to-column and column-to-row pivoting
- Configurable aggregation methods
- Nesting and flattening
- Tracks transformations and data loss

## Running the Example

```bash
cd examples/45-pivot
go run main.go
```
