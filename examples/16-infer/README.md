# Infer Example - Smart Data Completion

## What This Does

Demonstrates the **Infer** operation: intelligently filling in missing fields in partial data using LLM reasoning.

This example takes incomplete person and product records and infers missing information based on:
- Available partial data
- Type schema information
- Optional context clues
- Logical reasoning

## Use Case

**Real-World Application**: Complete customer profiles, product catalogs, or any dataset with missing information using AI-powered inference.

## How It Works

```go
complete, err := schemaflow.Infer[Person](partialData, ops.NewInferOptions().
    WithContext("Additional facts to guide inference"))
```

The LLM intelligently:
1. Analyzes available fields
2. Uses context and logical reasoning
3. Fills in missing data with appropriate values
4. Maintains data consistency

## Running the Example

```bash
# Set your API key
export OPENAI_API_KEY='your-key-here'

# Run the example
go run main.go
```

## Expected Output

```
🧠 Smart Data Inference Example
==================================================

📥 Partial Person Data:
  Name: John
  Age:  30

✅ Inferred Complete Person:
  Name:  John
  Age:   30
  Email: john.doe@example.com
  City:  San Francisco

📦 Partial Product Data:
  Name: iPhone 15

✅ Inferred Complete Product:
  Name:     iPhone 15
  Price:    $999.00
  Category: smartphone
  Brand:    Apple

✨ Success! Partial data → Complete records
```

## Key Features Demonstrated

- ✅ **Smart Inference**: Uses available data to fill gaps logically
- ✅ **Context Awareness**: Incorporates additional facts for better results
- ✅ **Type Preservation**: Maintains original data types and structure
- ✅ **Reasonable Defaults**: Provides sensible values when inference is uncertain

## Learn More

- [SchemaFlow Documentation](../../README.md)
- [Infer API Reference](../../docs/reference/API.md#infer)