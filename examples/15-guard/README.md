# Guard Example - Order State Validation

## What This Does

Demonstrates the **Guard** operation: validating state and conditions before proceeding with an action.

This example validates:
- **Input**: E-commerce orders ready for shipping
- **Checks**: 5 guard conditions (payment, address, items, status, total)
- **Output**: Pass/fail with detailed check results

## Use Case

**Real-World Applications**:
- Pre-condition validation (before actions)
- State machine transitions
- Workflow gate checks
- Safety validations
- Business rule enforcement
- Access control verification

## How It Works

```go
result := ops.Guard(order, checks...)

if result.Passed {
    // Safe to proceed with shipping
    shipOrder(order)
} else {
    // Show which checks failed
    for _, check := range result.Checks {
        if !check.Passed {
            log.Printf("Failed: %s", check.Message)
        }
    }
}
```

The guard system:
1. Runs all validation checks
2. Collects results and messages
3. Fails fast on first failure (optional)
4. Returns detailed pass/fail status
5. Provides actionable feedback

## Running the Example

```bash
# Set your API key
export OPENAI_API_KEY='your-key-here'

# Run the example
go run main.go
```

## Expected Output

```
🛡️ Guard Example - Order State Validation
============================================================

1. Valid Order - Ready to Ship
---
   Order ID: ORD-1001
   Status: processing
   Items: 3
   Total: $149.99
   Payment: true
   Address: 123 Main St, City, State 12345

   🛡️ Running guard checks...

   ✅ ALL GUARDS PASSED - Safe to ship

   Check Results:
      ✓ Payment confirmed
      ✓ Shipping address present
      ✓ 3 items ready
      ✓ Order status valid for shipping
      ✓ Order total: $149.99

2. Invalid - No Payment
---
   Order ID: ORD-1002
   Status: processing
   Items: 2
   Total: $89.50
   Payment: false
   Address: 456 Oak Ave, Town, State 54321

   🛡️ Running guard checks...

   ❌ GUARDS FAILED - Cannot proceed

   Check Results:
      ✗ Payment not received
      ✓ Shipping address present
      ✓ 2 items ready
      ✓ Order status valid for shipping
      ✓ Order total: $89.50

3. Invalid - Missing Address
---
   Order ID: ORD-1003
   Status: processing
   Items: 1
   Total: $29.99
   Payment: true
   Address: (missing)

   🛡️ Running guard checks...

   ❌ GUARDS FAILED - Cannot proceed

   Check Results:
      ✓ Payment confirmed
      ✗ Shipping address missing
      ✓ 1 items ready
      ✓ Order status valid for shipping
      ✓ Order total: $29.99

4. Invalid - Already Shipped
---
   Order ID: ORD-1004
   Status: shipped
   Items: 5
   Total: $299.99
   Payment: true
   Address: 789 Pine Rd, Village, State 98765

   🛡️ Running guard checks...

   ❌ GUARDS FAILED - Cannot proceed

   Check Results:
      ✓ Payment confirmed
      ✓ Shipping address present
      ✓ 5 items ready
      ✗ Order already shipped
      ✓ Order total: $299.99

📊 Guard Summary:
   Total orders checked: 4
   Passed all guards: 1
   Failed guards: 3

✨ Success! Guard checks complete
```

## Key Features Demonstrated

- ✅ **Multiple Checks**: Run several validations at once
- ✅ **Clear Feedback**: Know exactly what passed/failed
- ✅ **Type-Safe**: Guards work with any data type
- ✅ **Actionable**: Easy to show user what's missing

## Guard Patterns

Common guard check patterns:
1. **Pre-conditions**: Verify state before action
2. **Authorization**: Check permissions
3. **Data Validation**: Ensure data integrity
4. **Resource Availability**: Check limits/quotas
5. **Business Rules**: Enforce policies
6. **Safety Checks**: Prevent dangerous operations

## Learn More

- [SchemaFlow Documentation](../../README.md)
- [Guard API Reference](../../docs/reference/API.md#guard)
