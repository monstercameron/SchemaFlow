# Negotiate Example

This example demonstrates the `Negotiate` operation, which reconciles competing constraints to find an optimal solution.

## What it does

The `Negotiate` operation takes a set of constraints (which may conflict or compete) and finds a balanced solution that satisfies as many constraints as possible while documenting the tradeoffs made.

## Use Cases

- **Resource allocation**: Balance budget vs. features vs. timeline
- **Scheduling**: Find meeting times that work for multiple parties
- **Configuration**: Find settings that satisfy multiple requirements
- **Trade-off analysis**: When you can't have everything, find the best compromise

## Key Features

- Returns typed result matching your specified output schema
- Documents all tradeoffs made during negotiation
- Provides satisfaction scores for each constraint
- Explains reasoning behind the solution

## Running the Example

```bash
cd examples/36-negotiate
go run main.go
```
