package main

import (
	"fmt"
	"os"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// Review represents a customer review
type Review struct {
	ID      int
	Author  string
	Product string
	Text    string
}

// ClassifiedReview extends Review with classification results
type ClassifiedReview struct {
	Review
	Sentiment  string
	Confidence float64
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		os.Exit(1)
	}

	// Customer reviews to classify
	reviews := []Review{
		{
			ID:      1,
			Author:  "Sarah M.",
			Product: "Wireless Headphones",
			Text:    "Amazing sound quality! The bass is incredible and battery lasts for days. Best purchase I've made this year!",
		},
		{
			ID:      2,
			Author:  "Mike T.",
			Product: "Smart Watch",
			Text:    "Total disappointment. Stopped working after 2 weeks. Customer service was unhelpful. Don't waste your money.",
		},
		{
			ID:      3,
			Author:  "Jessica L.",
			Product: "Coffee Maker",
			Text:    "It's okay. Makes decent coffee. Nothing special but gets the job done. Average quality for the price.",
		},
		{
			ID:      4,
			Author:  "Tom K.",
			Product: "Standing Desk",
			Text:    "Life-changing for my back pain! Easy to assemble, smooth adjustment. Highly recommend for anyone working from home.",
		},
		{
			ID:      5,
			Author:  "Rachel P.",
			Product: "Bluetooth Speaker",
			Text:    "Arrived damaged and sound quality is terrible. Bass is non-existent. Returning immediately.",
		},
	}

	// Categories for classification
	categories := []string{"positive", "negative", "neutral"}

	fmt.Println("⭐ Classify Example - Sentiment Analysis")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println("\n📝 Customer Reviews to Classify:")
	fmt.Println()

	// Classify each review - now with confidence scores and alternatives!
	results := make(map[string][]ClassifiedReview)
	for _, category := range categories {
		results[category] = []ClassifiedReview{}
	}

	for _, review := range reviews {
		classifyOpts := schemaflow.NewClassifyOptions().WithCategories(categories)
		classifyOpts.OpOptions.Intelligence = schemaflow.Fast

		// Use the new generic Classify with typed result
		result, err := schemaflow.Classify[string, string](review.Text, classifyOpts)
		if err != nil {
			schemaflow.GetLogger().Warn("Failed to classify review", "review_id", review.ID, "error", err)
			continue
		}

		classified := ClassifiedReview{
			Review:     review,
			Sentiment:  result.Category,
			Confidence: result.Confidence,
		}
		results[result.Category] = append(results[result.Category], classified)

		// Show individual classification with confidence
		emoji := "😐"
		if result.Category == "positive" {
			emoji = "😊"
		} else if result.Category == "negative" {
			emoji = "😞"
		}

		fmt.Printf("Review #%d by %s %s\n", review.ID, review.Author, emoji)
		fmt.Printf("  Product:    %s\n", review.Product)
		fmt.Printf("  Sentiment:  %s (%.0f%% confidence)\n", result.Category, result.Confidence*100)
		if result.Reasoning != "" {
			fmt.Printf("  Reasoning:  %s\n", result.Reasoning)
		}
		if len(result.Alternatives) > 0 {
			fmt.Printf("  Alternatives:\n")
			for _, alt := range result.Alternatives {
				fmt.Printf("    - %s (%.0f%%)\n", alt.Category, alt.Confidence*100)
			}
		}
		fmt.Printf("  Review:     \"%s\"\n\n", review.Text)
	}

	// Show summary
	fmt.Println("📊 Classification Summary:")
	fmt.Println("---")
	fmt.Printf("😊 Positive Reviews: %d\n", len(results["positive"]))
	for _, r := range results["positive"] {
		fmt.Printf("   - Review #%d: %s ⭐⭐⭐⭐⭐ (%.0f%% confident)\n", r.ID, r.Product, r.Confidence*100)
	}

	fmt.Printf("\n😞 Negative Reviews: %d\n", len(results["negative"]))
	for _, r := range results["negative"] {
		fmt.Printf("   - Review #%d: %s ⭐ (%.0f%% confident)\n", r.ID, r.Product, r.Confidence*100)
	}

	fmt.Printf("\n😐 Neutral Reviews: %d\n", len(results["neutral"]))
	for _, r := range results["neutral"] {
		fmt.Printf("   - Review #%d: %s ⭐⭐⭐ (%.0f%% confident)\n", r.ID, r.Product, r.Confidence*100)
	}

	// Calculate sentiment distribution
	total := len(reviews)
	positivePercent := float64(len(results["positive"])) / float64(total) * 100
	negativePercent := float64(len(results["negative"])) / float64(total) * 100
	neutralPercent := float64(len(results["neutral"])) / float64(total) * 100

	fmt.Println("\n📈 Sentiment Distribution:")
	fmt.Printf("   Positive: %.0f%%\n", positivePercent)
	fmt.Printf("   Negative: %.0f%%\n", negativePercent)
	fmt.Printf("   Neutral:  %.0f%%\n", neutralPercent)

	fmt.Println("\n✨ Success! All reviews classified with confidence scores")
}
