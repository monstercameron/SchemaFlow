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

	fmt.Println("=== Rank Example ===")

	// Example 1: Rank articles by relevance
	fmt.Println("--- Example 1: Article Relevance Ranking ---")
	type Article struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	articles := []Article{
		{Title: "Introduction to Python", Content: "Python basics for beginners", Author: "Alice"},
		{Title: "Advanced Go Patterns", Content: "Design patterns in Go programming", Author: "Bob"},
		{Title: "Go Concurrency Deep Dive", Content: "Goroutines and channels explained", Author: "Charlie"},
		{Title: "Web Development with React", Content: "Building modern web apps", Author: "Diana"},
		{Title: "Go Performance Optimization", Content: "Making Go programs faster", Author: "Eve"},
		{Title: "Python Data Science", Content: "Using pandas and numpy", Author: "Frank"},
	}

	query := "Go programming patterns and best practices"
	opts := ops.NewRankOptions().
		WithQuery(query).
		WithTopK(3).
		WithIncludeExplanation(true)

	result, err := ops.Rank(articles, opts)
	if err != nil {
		log.Fatalf("Ranking failed: %v", err)
	}

	fmt.Printf("Query: '%s'\n", query)
	fmt.Printf("Top %d results:\n\n", len(result.Items))
	for i, item := range result.Items {
		fmt.Printf("  %d. %s (Score: %.2f)\n", i+1, item.Item.Title, item.Score)
		if item.Explanation != "" {
			fmt.Printf("     Reason: %s\n", item.Explanation)
		}
	}
	fmt.Println()

	// Example 2: Rank candidates for a job
	fmt.Println("--- Example 2: Candidate Ranking ---")
	type Candidate struct {
		Name       string   `json:"name"`
		Experience int      `json:"experience_years"`
		Skills     []string `json:"skills"`
		Education  string   `json:"education"`
	}

	candidates := []Candidate{
		{Name: "Alice", Experience: 5, Skills: []string{"Go", "Python", "Kubernetes"}, Education: "MS Computer Science"},
		{Name: "Bob", Experience: 3, Skills: []string{"Java", "Spring", "AWS"}, Education: "BS Software Engineering"},
		{Name: "Charlie", Experience: 7, Skills: []string{"Go", "Docker", "Terraform"}, Education: "PhD Distributed Systems"},
		{Name: "Diana", Experience: 2, Skills: []string{"JavaScript", "React", "Node"}, Education: "BS Computer Science"},
		{Name: "Eve", Experience: 4, Skills: []string{"Go", "gRPC", "Microservices"}, Education: "MS Software Engineering"},
	}

	jobQuery := "Senior Go developer with cloud infrastructure experience"
	candidateOpts := ops.NewRankOptions().
		WithQuery(jobQuery).
		WithTopK(3).
		WithBoostFields(map[string]float64{"skills": 2.0, "experience_years": 1.5}).
		WithIncludeExplanation(true)

	candidateResult, err := ops.Rank(candidates, candidateOpts)
	if err != nil {
		log.Fatalf("Candidate ranking failed: %v", err)
	}

	fmt.Printf("Job: '%s'\n", jobQuery)
	fmt.Printf("Top candidates:\n\n")
	for i, item := range candidateResult.Items {
		fmt.Printf("  %d. %s\n", i+1, item.Item.Name)
		fmt.Printf("     Experience: %d years\n", item.Item.Experience)
		fmt.Printf("     Skills: %v\n", item.Item.Skills)
		fmt.Printf("     Match Score: %.2f\n", item.Score)
		if item.Explanation != "" {
			fmt.Printf("     Why: %s\n", item.Explanation)
		}
		fmt.Println()
	}

	// Example 3: Rank with penalization
	fmt.Println("--- Example 3: Ranking with Boost/Penalize ---")
	type Product struct {
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		Rating    float64 `json:"rating"`
		InStock   bool    `json:"in_stock"`
		SpamScore float64 `json:"spam_score"` // Higher = more suspicious
	}

	products := []Product{
		{Name: "Premium Widget A", Price: 99.99, Rating: 4.8, InStock: true, SpamScore: 0.1},
		{Name: "Budget Widget B", Price: 29.99, Rating: 4.2, InStock: true, SpamScore: 0.8},
		{Name: "Luxury Widget C", Price: 199.99, Rating: 4.9, InStock: false, SpamScore: 0.05},
		{Name: "Standard Widget D", Price: 49.99, Rating: 4.5, InStock: true, SpamScore: 0.2},
	}

	productQuery := "Best value widget for everyday use"
	productOpts := ops.NewRankOptions().
		WithQuery(productQuery).
		WithBoostFields(map[string]float64{"rating": 2.0, "in_stock": 1.5}).
		WithPenalizeFields(map[string]float64{"spam_score": 0.5})

	productResult, err := ops.Rank(products, productOpts)
	if err != nil {
		log.Fatalf("Product ranking failed: %v", err)
	}

	fmt.Printf("Query: '%s'\n", productQuery)
	fmt.Printf("Ranked products:\n")
	for i, item := range productResult.Items {
		fmt.Printf("  %d. %s - $%.2f (Rating: %.1f, Score: %.2f)\n",
			i+1, item.Item.Name, item.Item.Price, item.Item.Rating, item.Score)
	}

	fmt.Println("\n=== Rank Example Complete ===")
}
