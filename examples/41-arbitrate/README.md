# Arbitrate Example

This example demonstrates the `Arbitrate` operation, which makes rule-based decisions with full audit trail.

## What it does

The `Arbitrate` operation evaluates multiple options against a set of rules and selects the best option, providing a complete audit trail of how each rule was applied to each option.

## Use Cases

- **Candidate selection**: Evaluate job applicants against criteria
- **Policy enforcement**: Apply business rules to decisions
- **Compliance checking**: Ensure choices meet regulatory requirements
- **Automated approvals**: Route requests based on defined rules

## Key Features

- Full audit trail of rule evaluations
- Per-rule pass/fail with reasoning
- Overall scoring and winner selection
- Configurable strictness (all must pass vs. best effort)

## Running the Example

```bash
cd examples/41-arbitrate
go run main.go
```
