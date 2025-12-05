// 20-complete: Complete partial text using LLM intelligence
// Intelligence: Fast (Cerebras gpt-oss-120b)
// Expectations:
// - Completes partial sentences/paragraphs naturally
// - Uses context for coherent completions
// - Supports different temperatures for creative vs. conservative output
// - CompleteField: completes a specific field in a struct using other fields as context

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	schemaflow "github.com/monstercameron/SchemaFlow"
)

// loadEnv loads environment variables from .env file
func loadEnv() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				log.Fatal("Error loading .env file")
			}
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	log.Fatal(".env file not found")
}

// BlogPost represents a blog post with partial content
type BlogPost struct {
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Body     string   `json:"body"` // This will be completed
}

// ProductDescription represents a product listing
type ProductDescription struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Description string  `json:"description"` // This will be completed
}

func main() {
	loadEnv()

	fmt.Println("✍️  Complete Example - Finish Partial Text with LLM")
	fmt.Println("=" + string(make([]byte, 60)))

	// Initialize SchemaFlow with Fast intelligence (Cerebras)
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	// Example 1: Basic sentence completion
	fmt.Println("\n1️⃣  Basic Sentence Completion")
	fmt.Println("─" + string(make([]byte, 60)))

	partial1 := "The benefits of using type-safe LLM operations include"
	fmt.Printf("   Input: %q\n", partial1)

	result1, err := schemaflow.Complete(partial1,
		schemaflow.NewCompleteOptions().
			WithMaxLength(100).
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✅ CompleteResult:")
		fmt.Printf("      Text:       %s\n", result1.Text)
		fmt.Printf("      Original:   %q\n", result1.Original)
		fmt.Printf("      Length:     %d characters added\n", result1.Length)
		fmt.Printf("      Confidence: %.2f\n", result1.Confidence)
	}

	// Example 2: Code completion
	fmt.Println("\n2️⃣  Code Completion")
	fmt.Println("─" + string(make([]byte, 60)))

	partial2 := "func validateEmail(email string) bool {\n    // Check if email is valid\n    "
	fmt.Printf("   Input:\n   %s\n", partial2)

	result2, err := schemaflow.Complete(partial2,
		schemaflow.NewCompleteOptions().
			WithMaxLength(150).
			WithTemperature(0.3). // Lower temp for code
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✅ CompleteResult:")
		fmt.Printf("      Text:\n%s\n", result2.Text)
		fmt.Printf("      Length:     %d characters added\n", result2.Length)
		fmt.Printf("      Confidence: %.2f\n", result2.Confidence)
	}

	// Example 3: Email completion with context
	fmt.Println("\n3️⃣  Email Completion with Context")
	fmt.Println("─" + string(make([]byte, 60)))

	partial3 := "Dear Customer Support,\n\nI am writing to request a refund because"
	context3 := []string{
		"Subject: Refund Request - Order #12345",
		"Customer purchased a laptop on Nov 15",
		"Product arrived damaged",
	}
	fmt.Printf("   Context: %v\n", context3)
	fmt.Printf("   Input: %q\n", partial3)

	result3, err := schemaflow.Complete(partial3,
		schemaflow.NewCompleteOptions().
			WithContext(context3).
			WithMaxLength(150).
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✅ CompleteResult:")
		fmt.Printf("      Text:       %s\n", result3.Text)
		fmt.Printf("      Length:     %d characters added\n", result3.Length)
		fmt.Printf("      Confidence: %.2f\n", result3.Confidence)
	}

	// Example 4: Creative story completion
	fmt.Println("\n4️⃣  Creative Story Completion (Higher Temperature)")
	fmt.Println("─" + string(make([]byte, 60)))

	partial4 := "The old lighthouse keeper had a secret that nobody in the village knew about. Every night at midnight, he would"
	fmt.Printf("   Input: %q\n", partial4)

	result4, err := schemaflow.Complete(partial4,
		schemaflow.NewCompleteOptions().
			WithMaxLength(100).
			WithTemperature(1.0). // Higher temp for creativity
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("   ✅ CompleteResult:")
		fmt.Printf("      Text:       %s\n", result4.Text)
		fmt.Printf("      Length:     %d characters added\n", result4.Length)
		fmt.Printf("      Confidence: %.2f\n", result4.Confidence)
	}

	// ========================================
	// CompleteField Examples - Complete a field in a struct
	// ========================================
	fmt.Println("\n" + "=" + string(make([]byte, 60)))
	fmt.Println("📦 CompleteField - Complete a Struct Field")
	fmt.Println("=" + string(make([]byte, 60)))

	// Example 5: Complete a blog post body using other fields as context
	fmt.Println("\n5️⃣  Complete Blog Post Body (using Title, Author, Tags as context)")
	fmt.Println("─" + string(make([]byte, 60)))

	blogPost := BlogPost{
		Title:    "The Future of Renewable Energy",
		Author:   "Dr. Sarah Chen",
		Category: "Technology",
		Tags:     []string{"solar", "wind", "sustainability", "climate"},
		Body:     "As the world faces increasing climate challenges, renewable energy sources are becoming",
	}

	fmt.Printf("   Input BlogPost:\n")
	fmt.Printf("      Title:    %s\n", blogPost.Title)
	fmt.Printf("      Author:   %s\n", blogPost.Author)
	fmt.Printf("      Category: %s\n", blogPost.Category)
	fmt.Printf("      Tags:     %v\n", blogPost.Tags)
	fmt.Printf("      Body:     %q\n", blogPost.Body)

	result5, err := schemaflow.CompleteField(blogPost,
		schemaflow.NewCompleteFieldOptions("Body").
			WithMaxLength(200).
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("\n   ✅ CompleteFieldResult:")
		fmt.Printf("      Field:      %s\n", result5.Field)
		fmt.Printf("      Original:   %q\n", result5.Original)
		fmt.Printf("      Completed:  %s\n", result5.Completed)
		fmt.Printf("      Length:     %d characters added\n", result5.Length)
		fmt.Printf("      Confidence: %.2f\n", result5.Confidence)
		fmt.Printf("\n   📄 Updated BlogPost.Body:\n      %s\n", result5.Data.Body)
	}

	// Example 6: Complete a product description
	fmt.Println("\n6️⃣  Complete Product Description (using Name, Price, Category as context)")
	fmt.Println("─" + string(make([]byte, 60)))

	product := ProductDescription{
		Name:        "EcoSmart Wireless Earbuds",
		Price:       79.99,
		Category:    "Electronics",
		Description: "These premium wireless earbuds feature",
	}

	fmt.Printf("   Input ProductDescription:\n")
	fmt.Printf("      Name:        %s\n", product.Name)
	fmt.Printf("      Price:       $%.2f\n", product.Price)
	fmt.Printf("      Category:    %s\n", product.Category)
	fmt.Printf("      Description: %q\n", product.Description)

	result6, err := schemaflow.CompleteField(product,
		schemaflow.NewCompleteFieldOptions("Description").
			WithMaxLength(150).
			WithTemperature(0.5). // Moderate creativity for product copy
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Println("\n   ✅ CompleteFieldResult:")
		fmt.Printf("      Field:      %s\n", result6.Field)
		fmt.Printf("      Original:   %q\n", result6.Original)
		fmt.Printf("      Completed:  %s\n", result6.Completed)
		fmt.Printf("      Length:     %d characters added\n", result6.Length)
		fmt.Printf("      Confidence: %.2f\n", result6.Confidence)
		fmt.Printf("\n   📄 Updated ProductDescription.Description:\n      %s\n", result6.Data.Description)
	}

	fmt.Println()
	fmt.Println("✨ Success! Partial text and struct fields completed with LLM intelligence")
}
