# Verify Example

Demonstrates the `Verify` operation for fact-checking and consistency validation.

## What is Verify?

The `Verify` operation checks claims for accuracy and consistency:
- Fact-checks statements against knowledge
- Validates logical reasoning
- Checks internal consistency
- Provides evidence and corrections

## Key Features

- Multiple verification types
- Source-based verification
- Logic checking
- Consistency validation
- Severity-graded verdicts

## Usage

```go
opts := ops.NewVerifyOptions().
    WithCheckFacts(true).
    WithCheckLogic(true).
    WithIncludeEvidence(true)

result, err := ops.Verify(content, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithSources()` | Knowledge sources to verify against | [] |
| `WithStrictness()` | Verification strictness (strict, moderate, lenient) | "moderate" |
| `WithCheckFacts()` | Verify factual accuracy | true |
| `WithCheckLogic()` | Check logical consistency | true |
| `WithCheckConsistency()` | Check internal consistency | true |
| `WithIncludeEvidence()` | Include supporting evidence | true |
| `WithExplainReasoning()` | Explain verification logic | true |
| `WithMinConfidence()` | Minimum confidence for "verified" | 0.7 |

## Verdicts

| Verdict | Meaning |
|---------|---------|
| `verified` | Accurate and supported by evidence |
| `false` | Demonstrably incorrect |
| `partially_true` | Contains some truth but incomplete/misleading |
| `misleading` | Technically true but presented deceptively |
| `unverifiable` | Cannot be verified with available information |

## Running

```bash
cd examples/34-verify
go run main.go
```

## Example Output

```
Overall Verdict: mixed (Trust Score: 0.72)

Individual Claims:

  ✓ verified (95% confidence)
  Claim: The Earth is approximately 4.5 billion years old.
  Reasoning: Consistent with radiometric dating evidence

  ✓ verified (98% confidence)
  Claim: Water boils at 100 degrees Celsius at sea level.
  Reasoning: Standard physics at 1 atm pressure

  ✗ false (92% confidence)
  Claim: The capital of Australia is Sydney.
  Reasoning: Common misconception
  Correction: The capital of Australia is Canberra
```

## Difference from Validate

- **Validate**: Checks structure/schema compliance
- **Verify**: Checks factual accuracy and truth

## Use Cases

- **Content moderation**: Fact-check articles
- **Research validation**: Verify citations and claims
- **Document review**: Check consistency
- **Argument analysis**: Evaluate logical validity
