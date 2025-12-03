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

	fmt.Println("=== Synthesize Example ===")

	// Example 1: Synthesize research from multiple sources
	fmt.Println("--- Example 1: Research Synthesis ---")
	type ResearchReport struct {
		Summary     string   `json:"summary"`
		KeyFindings []string `json:"key_findings"`
		Consensus   string   `json:"consensus"`
		Gaps        []string `json:"gaps"`
	}

	sources := []any{
		map[string]any{
			"title":   "Study A: Remote Work Productivity",
			"year":    2023,
			"finding": "Remote workers showed 13% higher productivity",
			"sample":  "500 tech employees",
		},
		map[string]any{
			"title":   "Study B: Hybrid Work Analysis",
			"year":    2023,
			"finding": "Hybrid workers reported better work-life balance but 5% lower productivity",
			"sample":  "1200 corporate employees",
		},
		map[string]any{
			"title":   "Study C: Office vs Remote Collaboration",
			"year":    2024,
			"finding": "In-office teams showed stronger collaboration metrics, remote teams better focus time",
			"sample":  "300 teams across industries",
		},
	}

	opts := ops.NewSynthesizeOptions().
		WithStrategy("integrate").
		WithCiteSources(true).
		WithGenerateInsights(true)

	result, err := ops.Synthesize[ResearchReport](sources, opts)
	if err != nil {
		log.Fatalf("Synthesis failed: %v", err)
	}

	fmt.Println("Synthesized Research Report:")
	fmt.Printf("\nSummary: %s\n", result.Synthesized.Summary)
	fmt.Println("\nKey Findings:")
	for _, finding := range result.Synthesized.KeyFindings {
		fmt.Printf("  - %s\n", finding)
	}
	fmt.Printf("\nConsensus: %s\n", result.Synthesized.Consensus)

	if len(result.Insights) > 0 {
		fmt.Println("\nGenerated Insights:")
		for _, insight := range result.Insights {
			fmt.Printf("  [%s] %s\n", insight.Type, insight.Insight)
		}
	}

	if len(result.Conflicts) > 0 {
		fmt.Println("\nConflicts Identified:")
		for _, conflict := range result.Conflicts {
			fmt.Printf("  - %s: %s\n", conflict.Topic, conflict.Resolution)
		}
	}
	fmt.Println()

	// Example 2: Reconcile conflicting data
	fmt.Println("--- Example 2: Data Reconciliation ---")
	type CustomerProfile struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}

	customerSources := []any{
		map[string]any{
			"source": "CRM",
			"name":   "John Smith",
			"email":  "john@email.com",
			"phone":  "555-1234",
		},
		map[string]any{
			"source":  "Billing",
			"name":    "J. Smith",
			"email":   "john.smith@email.com",
			"phone":   "555-1234",
			"address": "123 Main St",
		},
		map[string]any{
			"source": "Support",
			"name":   "John Smith",
			"phone":  "555-5678", // Different phone
		},
	}

	reconcileOpts := ops.NewSynthesizeOptions().
		WithStrategy("reconcile").
		WithConflictResolution("source-priority").
		WithSourcePriorities([]int{1, 0, 2}) // Billing > CRM > Support

	customerResult, err := ops.Synthesize[CustomerProfile](customerSources, reconcileOpts)
	if err != nil {
		log.Fatalf("Reconciliation failed: %v", err)
	}

	fmt.Println("Reconciled Customer Profile:")
	fmt.Printf("  Name: %s\n", customerResult.Synthesized.Name)
	fmt.Printf("  Email: %s\n", customerResult.Synthesized.Email)
	fmt.Printf("  Phone: %s\n", customerResult.Synthesized.Phone)
	fmt.Printf("  Address: %s\n", customerResult.Synthesized.Address)

	if len(customerResult.Conflicts) > 0 {
		fmt.Println("\nResolved Conflicts:")
		for _, c := range customerResult.Conflicts {
			fmt.Printf("  - %s: Chose source %d\n", c.Topic, c.Chosen)
		}
	}
	fmt.Println()

	// Example 3: Compare and contrast
	fmt.Println("--- Example 3: Compare and Contrast ---")
	type ComparisonReport struct {
		Similarities   []string `json:"similarities"`
		Differences    []string `json:"differences"`
		Recommendation string   `json:"recommendation"`
	}

	products := []any{
		map[string]any{
			"name":     "Product A",
			"price":    99.99,
			"features": []string{"Feature 1", "Feature 2", "Feature 3"},
			"rating":   4.5,
		},
		map[string]any{
			"name":     "Product B",
			"price":    79.99,
			"features": []string{"Feature 1", "Feature 4", "Feature 5"},
			"rating":   4.2,
		},
	}

	compareOpts := ops.NewSynthesizeOptions().
		WithStrategy("compare").
		WithFocusAreas([]string{"price", "features", "value"})

	compareResult, err := ops.Synthesize[ComparisonReport](products, compareOpts)
	if err != nil {
		log.Fatalf("Comparison failed: %v", err)
	}

	fmt.Println("Product Comparison:")
	fmt.Println("\nSimilarities:")
	for _, s := range compareResult.Synthesized.Similarities {
		fmt.Printf("  - %s\n", s)
	}
	fmt.Println("\nDifferences:")
	for _, d := range compareResult.Synthesized.Differences {
		fmt.Printf("  - %s\n", d)
	}
	fmt.Printf("\nRecommendation: %s\n", compareResult.Synthesized.Recommendation)

	// Example 4: Merge with insights
	fmt.Println("\n--- Example 4: Document Merge ---")
	type MergedDoc struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Authors  []string `json:"authors"`
		Sections []string `json:"sections"`
	}

	docs := []any{
		map[string]any{
			"title":   "Architecture Overview",
			"author":  "Alice",
			"content": "The system uses microservices architecture...",
		},
		map[string]any{
			"title":   "Technical Implementation",
			"author":  "Bob",
			"content": "Each service is containerized using Docker...",
		},
		map[string]any{
			"title":   "Deployment Guide",
			"author":  "Charlie",
			"content": "Deploy to Kubernetes using helm charts...",
		},
	}

	mergeOpts := ops.NewSynthesizeOptions().
		WithStrategy("merge").
		WithCiteSources(true)

	mergeResult, err := ops.Synthesize[MergedDoc](docs, mergeOpts)
	if err != nil {
		log.Fatalf("Merge failed: %v", err)
	}

	fmt.Println("Merged Document:")
	fmt.Printf("  Title: %s\n", mergeResult.Synthesized.Title)
	fmt.Printf("  Authors: %v\n", mergeResult.Synthesized.Authors)
	fmt.Printf("  Sections: %v\n", mergeResult.Synthesized.Sections)

	fmt.Printf("\nSource Coverage:\n")
	for source, coverage := range mergeResult.SourceCoverage {
		fmt.Printf("  Source %d: %.0f%%\n", source, coverage*100)
	}

	fmt.Println("\n=== Synthesize Example Complete ===")
}
