# Classify Example - Sentiment Analysis

## What This Does

Demonstrates the **Classify** operation: categorizing text into predefined categories.

This example classifies:
- **Input**: Customer review text
- **Categories**: positive, negative, neutral
- **Output**: Sentiment classification for each review

## Use Case

**Real-World Applications**:
- Customer feedback analysis
- Social media monitoring
- Email categorization (spam/not spam)
- Content moderation
- Support ticket classification

## How It Works

```go
sentiment, err := ops.Classify(
    review.Text,
    ops.NewClassifyOptions().WithCategories([]string{"positive", "negative", "neutral"}),
)
```

The LLM intelligently:
1. Analyzes text tone and language
2. Identifies emotional indicators
3. Considers context and nuance
4. Assigns to most appropriate category
5. Handles edge cases and sarcasm

## Running the Example

```bash
# Set your API key
export OPENAI_API_KEY='your-key-here'

# Run the example
go run main.go
```

## Expected Output

```
⭐ Classify Example - Sentiment Analysis
============================================================

📝 Customer Reviews to Classify:

Review #1 by Sarah M. 😊
  Product:    Wireless Headphones
  Sentiment:  positive
  Review:     "Amazing sound quality! The bass is incredible..."

Review #2 by Mike T. 😞
  Product:    Smart Watch
  Sentiment:  negative
  Review:     "Total disappointment. Stopped working after 2 weeks..."

Review #3 by Jessica L. 😐
  Product:    Coffee Maker
  Sentiment:  neutral
  Review:     "It's okay. Makes decent coffee. Nothing special..."

Review #4 by Tom K. 😊
  Product:    Standing Desk
  Sentiment:  positive
  Review:     "Life-changing for my back pain! Easy to assemble..."

Review #5 by Rachel P. 😞
  Product:    Bluetooth Speaker
  Sentiment:  negative
  Review:     "Arrived damaged and sound quality is terrible..."

📊 Classification Summary:
---
😊 Positive Reviews: 2
   - Review #1: Wireless Headphones ⭐⭐⭐⭐⭐
   - Review #4: Standing Desk ⭐⭐⭐⭐⭐

😞 Negative Reviews: 2
   - Review #2: Smart Watch ⭐
   - Review #5: Bluetooth Speaker ⭐

😐 Neutral Reviews: 1
   - Review #3: Coffee Maker ⭐⭐⭐

📈 Sentiment Distribution:
   Positive: 40%
   Negative: 40%
   Neutral:  20%

✨ Success! All reviews classified by sentiment
```

## Key Features Demonstrated

- ✅ **Sentiment Detection**: Positive/negative/neutral
- ✅ **Contextual Understanding**: Not just keyword matching
- ✅ **Batch Processing**: Multiple reviews
- ✅ **Analytics**: Summary statistics and distribution

## Use Cases

1. **E-commerce**: Analyze product reviews
2. **Social Media**: Monitor brand sentiment
3. **Customer Service**: Auto-categorize feedback
4. **Content Moderation**: Flag inappropriate content
5. **Market Research**: Track sentiment trends

## Learn More

- [SchemaFlow Documentation](../../README.md)
- [Classify API Reference](../../docs/reference/API.md#classify)
