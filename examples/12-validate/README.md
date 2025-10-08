# Validate Example - User Registration Validation

## What This Does

Demonstrates the **Validate** operation: intelligent data validation using LLM reasoning.

This example validates:
- **Input**: User registration data
- **Rules**: Complex business rules (email format, password strength, age restrictions)
- **Output**: ValidationResult with detailed issues

## Use Case

**Real-World Applications**:
- Form validation with complex rules
- Data quality checks
- Business rule validation
- Contract compliance checking
- Document verification
- Configuration validation

## How It Works

```go
result, err := ops.Validate(registration, validationRules)

// Check result
if result.Valid {
    // Process registration
} else {
    // Show errors to user
    for _, issue := range result.Issues {
        fmt.Printf("%s: %s\n", issue.Field, issue.Message)
    }
}
```

The LLM intelligently:
1. Understands complex validation rules
2. Applies context-aware validation
3. Provides detailed error messages
4. Handles edge cases gracefully
5. Supports custom business logic

## Running the Example

```bash
# Set your API key
export OPENAI_API_KEY='your-key-here'

# Run the example
go run main.go
```

## Expected Output

```
✅ Validate Example - User Registration Validation
============================================================

1. Valid Registration
---
   Username: johndoe123
   Email: john.doe@example.com
   Password: Se************ (length: 15)
   Age: 25
   Country: USA

   ✅ VALID - Registration accepted

2. Invalid Email
---
   Username: janedoe
   Email: not-an-email
   Password: Go************ (length: 16)
   Age: 30
   Country: Canada

   ❌ INVALID - Errors found:
      • Email: Invalid email format
      
3. Weak Password
---
   Username: bobsmith
   Email: bob@example.com
   Password: 12*** (length: 5)
   Age: 22
   Country: UK

   ❌ INVALID - Errors found:
      • Password: Too short (minimum 8 characters)
      • Password: Missing uppercase letter
      • Password: Missing special character
      
4. Underage User
---
   Username: younguser
   Email: young@example.com
   Password: St************ (length: 14)
   Age: 15
   Country: Germany

   ❌ INVALID - Errors found:
      • Age: Must be 18 or older

📊 Validation Summary:
   Total tested: 4 registrations
   Expected: 1 valid, 3 invalid

✨ Success! Validation complete
```

## Key Features Demonstrated

- ✅ **Intelligent Validation**: Beyond regex patterns
- ✅ **Detailed Feedback**: Field-specific error messages
- ✅ **Context-Aware**: Understands business rules
- ✅ **Multiple Issues**: Reports all validation problems

## Advantages Over Traditional Validation

1. **Natural Language Rules**: Write rules in plain English
2. **Context Understanding**: Handles nuanced validation
3. **Better Error Messages**: Human-friendly explanations
4. **Adaptive**: Learns from examples
5. **Less Code**: No complex regex or validation libraries

## Learn More

- [SchemaFlow Documentation](../../README.md)
- [Validate API Reference](../../docs/reference/API.md#validate)
