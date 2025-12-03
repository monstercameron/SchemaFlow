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

	fmt.Println("=== SemanticMatch Example ===")

	// Example 1: Match queries to products
	fmt.Println("--- Example 1: Product Matching ---")
	type Product struct {
		ID          int      `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Category    string   `json:"category"`
		Tags        []string `json:"tags"`
	}

	type Query struct {
		SearchTerm string `json:"search_term"`
	}

	products := []Product{
		{ID: 1, Name: "Wireless Bluetooth Headphones", Description: "Over-ear noise cancelling headphones", Category: "Audio", Tags: []string{"wireless", "bluetooth", "audio"}},
		{ID: 2, Name: "USB-C Charging Cable", Description: "Fast charging cable for phones and laptops", Category: "Accessories", Tags: []string{"cable", "charging", "usb-c"}},
		{ID: 3, Name: "Mechanical Gaming Keyboard", Description: "RGB keyboard with Cherry MX switches", Category: "Peripherals", Tags: []string{"keyboard", "gaming", "mechanical"}},
		{ID: 4, Name: "Wireless Mouse", Description: "Ergonomic mouse with adjustable DPI", Category: "Peripherals", Tags: []string{"mouse", "wireless", "ergonomic"}},
		{ID: 5, Name: "Portable Power Bank", Description: "20000mAh battery pack with fast charging", Category: "Accessories", Tags: []string{"battery", "charging", "portable"}},
	}

	queries := []Query{
		{SearchTerm: "bluetooth audio device"},
		{SearchTerm: "charging accessories"},
		{SearchTerm: "gaming equipment"},
	}

	opts := ops.NewMatchOptions().
		WithStrategy("best-fit").
		WithThreshold(0.5).
		WithMatchFields([]string{"name", "description", "tags"})

	result, err := ops.SemanticMatch(queries, products, opts)
	if err != nil {
		log.Fatalf("Matching failed: %v", err)
	}

	fmt.Println("Match results:")
	for _, match := range result.Matches {
		fmt.Printf("\n  Query: '%s'\n", queries[match.SourceIndex].SearchTerm)
		fmt.Printf("  Match: %s (Score: %.2f)\n", products[match.TargetIndex].Name, match.Score)
		if match.Explanation != "" {
			fmt.Printf("  Why: %s\n", match.Explanation)
		}
	}
	fmt.Println()

	// Example 2: Match candidates to job requirements
	fmt.Println("--- Example 2: Candidate-Job Matching ---")
	type JobReq struct {
		Title      string   `json:"title"`
		Skills     []string `json:"skills"`
		Experience int      `json:"experience_years"`
	}

	type Candidate struct {
		Name       string   `json:"name"`
		Skills     []string `json:"skills"`
		Experience int      `json:"experience"`
	}

	jobs := []JobReq{
		{Title: "Senior Go Developer", Skills: []string{"Go", "Kubernetes", "gRPC"}, Experience: 5},
		{Title: "Frontend Engineer", Skills: []string{"React", "TypeScript", "CSS"}, Experience: 3},
	}

	candidates := []Candidate{
		{Name: "Alice", Skills: []string{"Go", "Docker", "Kubernetes"}, Experience: 6},
		{Name: "Bob", Skills: []string{"JavaScript", "React", "Node.js"}, Experience: 4},
		{Name: "Charlie", Skills: []string{"Python", "Django", "PostgreSQL"}, Experience: 5},
		{Name: "Diana", Skills: []string{"Go", "gRPC", "Microservices"}, Experience: 4},
	}

	jobOpts := ops.NewMatchOptions().
		WithStrategy("one-to-one").
		WithThreshold(0.4).
		WithFieldWeights(map[string]float64{
			"skills":     2.0,
			"experience": 1.0,
		})

	jobResult, err := ops.SemanticMatch(jobs, candidates, jobOpts)
	if err != nil {
		log.Fatalf("Job matching failed: %v", err)
	}

	fmt.Println("Job-Candidate matches:")
	for _, match := range jobResult.Matches {
		fmt.Printf("\n  %s -> %s (Score: %.2f)\n", jobs[match.SourceIndex].Title, candidates[match.TargetIndex].Name, match.Score)
	}

	fmt.Printf("\nUnmatched candidates: %d\n", len(jobResult.UnmatchedTargets))
	fmt.Println()

	// Example 3: Match a single query using MatchOne
	fmt.Println("--- Example 3: Single Item Matching ---")
	type SearchQuery struct {
		Term string `json:"term"`
	}
	singleQuery := SearchQuery{Term: "noise cancelling earbuds for travel"}

	matches, err := ops.MatchOne(singleQuery, products, ops.NewMatchOptions().WithStrategy("best-fit"))
	if err != nil {
		log.Fatalf("Single matching failed: %v", err)
	}

	fmt.Printf("Query: '%s'\n", singleQuery.Term)
	if len(matches) > 0 {
		fmt.Printf("Best match: %s (Score: %.2f)\n", matches[0].Target.Name, matches[0].Score)
	}
	fmt.Println()

	// Example 4: Entity resolution matching
	fmt.Println("--- Example 4: Entity Resolution ---")
	type Record struct {
		Name    string `json:"name"`
		Company string `json:"company"`
		Email   string `json:"email"`
	}

	sourceRecords := []Record{
		{Name: "John Smith", Company: "Acme Corp", Email: "john.s@acme.com"},
		{Name: "Jane Doe", Company: "Tech Inc", Email: "jdoe@tech.io"},
	}

	targetRecords := []Record{
		{Name: "J. Smith", Company: "Acme Corporation", Email: "john.smith@acme.com"},
		{Name: "Jane D.", Company: "Tech Incorporated", Email: "jane.doe@tech.io"},
		{Name: "Bob Wilson", Company: "Other LLC", Email: "bob@other.com"},
	}

	entityOpts := ops.NewMatchOptions().
		WithStrategy("all-matches").
		WithThreshold(0.6).
		WithAllowPartial(true)

	entityResult, err := ops.SemanticMatch(sourceRecords, targetRecords, entityOpts)
	if err != nil {
		log.Fatalf("Entity matching failed: %v", err)
	}

	fmt.Println("Entity resolution matches:")
	for _, match := range entityResult.Matches {
		source := sourceRecords[match.SourceIndex]
		target := targetRecords[match.TargetIndex]
		fmt.Printf("  %s (%s) <-> %s (%s) [%.2f]\n",
			source.Name, source.Company,
			target.Name, target.Company,
			match.Score)
	}

	fmt.Println("\n=== SemanticMatch Example Complete ===")
}
