# Score Example - Code Quality Assessment

## What This Does

Demonstrates the **Score** operation: rating content based on specified criteria using numeric scores.

This example scores:
- **Input**: Code snippets
- **Criteria**: Readability, maintainability, documentation, error handling, best practices
- **Output**: Numeric score (1-10) with ranking

## Use Case

**Real-World Applications**:
- Code review automation
- Content quality assessment
- Essay grading
- Product quality ratings
- Performance evaluations

## How It Works

```go
score, err := ops.Score(
    codeSnippet,
    ops.NewScoreOptions().
        WithScaleMin(1).
        WithScaleMax(10).
        WithCriteria([]string{"readability", "maintainability", ...}),
)
```

The LLM intelligently:
1. Evaluates against multiple criteria
2. Balances different quality factors
3. Applies domain expertise
4. Generates consistent scores
5. Identifies strengths and weaknesses

## Running the Example

```bash
# Set your API key
export OPENAI_API_KEY='your-key-here'

# Run the example
go run main.go
```

## Expected Output

```
📊 Score Example - Code Quality Assessment
============================================================

🎯 Evaluation Criteria: [readability maintainability documentation error handling best practices]
📏 Scale: 1-10 (10 = excellent)

📝 Evaluating Submission #1 by Alice
---
func calculateTotal(prices []float64) float64 {
    total := 0.0
    for _, price := range prices {
        total += price
    }
    return total
}
---

✅ Score: 7.5/10 ⭐⭐⭐⭐⭐⭐⭐☆☆☆
👍 Good code quality. Minor improvements possible.

📝 Evaluating Submission #2 by Bob
---
func calculateTotal(p []float64) float64 {
    var t float64
    for i:=0;i<len(p);i++ {
        t=t+p[i]
    }
    return t
}
---

✅ Score: 4.5/10 ⭐⭐⭐⭐☆☆☆☆☆☆
⚠️  Acceptable but needs improvement.

📝 Evaluating Submission #3 by Carol
---
// CalculateTotal computes the sum of all prices in the slice.
// It returns 0.0 for an empty slice.
// Time complexity: O(n)
func CalculateTotal(prices []float64) float64 {
    if len(prices) == 0 {
        return 0.0
    }
    
    total := 0.0
    for _, price := range prices {
        if price < 0 {
            continue // Skip negative prices
        }
        total += price
    }
    return total
}
---

✅ Score: 9.2/10 ⭐⭐⭐⭐⭐⭐⭐⭐⭐☆
💎 Excellent! Production-ready code with best practices.

🏆 Final Rankings:
---
🥇 1. Carol - Score: 9.2/10
🥈 2. Alice - Score: 7.5/10
🥉 3. Bob - Score: 4.5/10

✨ Success! All code submissions evaluated
```

## Key Features Demonstrated

- ✅ **Multi-Criteria Evaluation**: Multiple quality factors
- ✅ **Numeric Scoring**: Consistent 1-10 scale
- ✅ **Comparative Analysis**: Ranking by quality
- ✅ **Visual Feedback**: Star ratings and medals

## Use Cases

1. **Code Review**: Automate quality assessment
2. **Education**: Grade assignments consistently
3. **Content**: Rate article/video quality
4. **Hiring**: Score technical interviews
5. **QA**: Evaluate test coverage and quality

## Learn More

- [SchemaFlow Documentation](../../README.md)
- [Score API Reference](../../docs/reference/API.md#score)
