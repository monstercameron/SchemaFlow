# Audit Example

This example demonstrates the `Audit` operation, which performs deep inspection for issues, anomalies, and policy violations.

## What it does

The `Audit` operation analyzes data for problems including security issues, compliance violations, data quality issues, and internal inconsistencies.

## Categories of Findings

- **Security**: Sensitive data exposure, weak validation
- **Compliance**: Policy violations, regulatory issues
- **Quality**: Data accuracy, format problems
- **Consistency**: Internal contradictions
- **Completeness**: Missing required data

## Key Features

- Configurable audit policies
- Severity levels (info to critical)
- Full findings with recommendations
- Summary statistics

## Running the Example

```bash
cd examples/43-audit
go run main.go
```
