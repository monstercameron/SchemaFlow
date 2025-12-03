package main

import (
	"fmt"
	"log"
	"os"

	"github.com/monstercameron/SchemaFlow/internal/ops"
)

func main() {
	// Ensure environment is configured
	if os.Getenv("SCHEMAFLOW_API_KEY") == "" {
		log.Fatal("SCHEMAFLOW_API_KEY environment variable not set")
	}

	fmt.Println("=== Enrich Example ===")

	// Example 1: Enrich a product with derived fields
	fmt.Println("--- Example 1: Product Enrichment ---")
	type Product struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	type EnrichedProduct struct {
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		Price        float64  `json:"price"`
		Category     string   `json:"category"`
		Keywords     []string `json:"keywords"`
		TargetMarket string   `json:"target_market"`
		PriceRange   string   `json:"price_range"`
	}

	product := Product{
		Name:        "Pro Gaming Keyboard RGB",
		Description: "Mechanical keyboard with Cherry MX switches, customizable RGB lighting, and programmable macros",
		Price:       149.99,
	}

	opts := ops.NewEnrichOptions().
		WithDeriveFields([]string{"category", "keywords", "target_market", "price_range"}).
		WithDomain("e-commerce")

	result, err := ops.Enrich[Product, EnrichedProduct](product, opts)
	if err != nil {
		log.Fatalf("Enrichment failed: %v", err)
	}

	fmt.Printf("Original product: %s\n", product.Name)
	fmt.Printf("Enriched with:\n")
	fmt.Printf("  Category: %s\n", result.Enriched.Category)
	fmt.Printf("  Keywords: %v\n", result.Enriched.Keywords)
	fmt.Printf("  Target Market: %s\n", result.Enriched.TargetMarket)
	fmt.Printf("  Price Range: %s\n", result.Enriched.PriceRange)
	fmt.Println()

	// Example 2: Enrich a contact with inferred information
	fmt.Println("--- Example 2: Contact Enrichment ---")
	type Contact struct {
		Email   string `json:"email"`
		Company string `json:"company"`
		Title   string `json:"title"`
	}

	type EnrichedContact struct {
		Email      string `json:"email"`
		Company    string `json:"company"`
		Title      string `json:"title"`
		Department string `json:"department"`
		Seniority  string `json:"seniority"`
		Industry   string `json:"industry"`
		Region     string `json:"region"`
	}

	contact := Contact{
		Email:   "john.smith@techcorp.io",
		Company: "TechCorp Inc",
		Title:   "VP of Engineering",
	}

	contactOpts := ops.NewEnrichOptions().
		WithDeriveFields([]string{"department", "seniority", "industry", "region"}).
		WithDerivationRules(map[string]string{
			"seniority": "Infer from job title (executive, senior, mid, junior)",
			"region":    "Infer from email domain if possible",
		})

	contactResult, err := ops.Enrich[Contact, EnrichedContact](contact, contactOpts)
	if err != nil {
		log.Fatalf("Contact enrichment failed: %v", err)
	}

	fmt.Printf("Original: %s at %s\n", contact.Title, contact.Company)
	fmt.Printf("Enriched:\n")
	fmt.Printf("  Department: %s\n", contactResult.Enriched.Department)
	fmt.Printf("  Seniority: %s\n", contactResult.Enriched.Seniority)
	fmt.Printf("  Industry: %s\n", contactResult.Enriched.Industry)
	fmt.Printf("  Region: %s\n", contactResult.Enriched.Region)
	fmt.Println()

	// Example 3: Enrich text content
	fmt.Println("--- Example 3: Article Enrichment ---")
	type Article struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	type EnrichedArticle struct {
		Title        string   `json:"title"`
		Content      string   `json:"content"`
		Author       string   `json:"author"`
		Summary      string   `json:"summary"`
		Topics       []string `json:"topics"`
		ReadingTime  string   `json:"reading_time"`
		Difficulty   string   `json:"difficulty"`
		KeyTakeaways []string `json:"key_takeaways"`
	}

	article := Article{
		Title:   "Understanding Kubernetes Operators",
		Content: "Kubernetes Operators extend the Kubernetes API to manage complex applications...",
		Author:  "Jane Developer",
	}

	articleOpts := ops.NewEnrichOptions().
		WithDeriveFields([]string{"summary", "topics", "reading_time", "difficulty", "key_takeaways"}).
		WithDomain("technology").
		WithIncludeConfidence(true)

	articleResult, err := ops.Enrich[Article, EnrichedArticle](article, articleOpts)
	if err != nil {
		log.Fatalf("Article enrichment failed: %v", err)
	}

	fmt.Printf("Article: %s\n", article.Title)
	fmt.Printf("Enriched:\n")
	fmt.Printf("  Summary: %s\n", articleResult.Enriched.Summary)
	fmt.Printf("  Topics: %v\n", articleResult.Enriched.Topics)
	fmt.Printf("  Reading Time: %s\n", articleResult.Enriched.ReadingTime)
	fmt.Printf("  Difficulty: %s\n", articleResult.Enriched.Difficulty)
	fmt.Printf("  Key Takeaways: %v\n", articleResult.Enriched.KeyTakeaways)

	// Example 4: In-place enrichment
	fmt.Println("\n--- Example 4: In-Place Enrichment ---")
	type Document struct {
		Text     string   `json:"text"`
		Keywords []string `json:"keywords"`
		Language string   `json:"language"`
	}

	doc := Document{
		Text: "Machine learning is transforming healthcare diagnostics.",
	}

	inPlaceOpts := ops.NewEnrichOptions().
		WithDeriveFields([]string{"keywords", "language"})

	enrichedDoc, err := ops.EnrichInPlace(doc, inPlaceOpts)
	if err != nil {
		log.Fatalf("In-place enrichment failed: %v", err)
	}

	fmt.Printf("Original text: %s\n", doc.Text)
	fmt.Printf("Added keywords: %v\n", enrichedDoc.Keywords)
	fmt.Printf("Detected language: %s\n", enrichedDoc.Language)

	fmt.Println("\n=== Enrich Example Complete ===")
}
