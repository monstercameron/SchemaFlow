# Conform Example

This example demonstrates the `Conform` operation, which transforms data to match specific standards.

## What it does

The `Conform` operation takes data and transforms it to match a specific standard format (like USPS addresses, ISO8601 dates, E164 phone numbers, etc.), documenting all adjustments made.

## Supported Standards

- **USPS**: US Postal Service address format
- **ISO8601**: Date/time format
- **E164**: International phone number format
- **ISO3166**: Country codes
- **Custom standards**: Define your own rules

## Use Cases

- **Address standardization**: Normalize addresses for shipping
- **Data quality**: Ensure dates, phones, IDs match expected formats
- **API integration**: Transform data to match external API requirements
- **Compliance**: Ensure data meets regulatory format requirements

## Running the Example

```bash
cd examples/39-conform
go run main.go
```
