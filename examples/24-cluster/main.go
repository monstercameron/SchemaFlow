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

	fmt.Println("=== Cluster Example ===")

	// Example 1: Cluster articles by topic
	fmt.Println("--- Example 1: Article Clustering ---")
	articles := []string{
		"Python is a versatile programming language for data science",
		"Machine learning algorithms improve with more training data",
		"JavaScript frameworks like React are popular for web development",
		"Deep learning uses neural networks with multiple layers",
		"Vue.js provides reactive data binding for web applications",
		"Natural language processing enables text understanding",
		"TypeScript adds static typing to JavaScript",
		"Convolutional neural networks excel at image recognition",
		"Angular is a comprehensive framework for enterprise apps",
		"Transformer models have revolutionized NLP tasks",
	}

	opts := ops.NewClusterOptions().
		WithNumClusters(3).
		WithNamingStrategy("descriptive").
		WithIncludeOutliers(true)

	result, err := ops.Cluster(articles, opts)
	if err != nil {
		log.Fatalf("Clustering failed: %v", err)
	}

	fmt.Printf("Created %d clusters:\n", len(result.Clusters))
	for _, cluster := range result.Clusters {
		fmt.Printf("\n  Cluster: %s\n", cluster.Name)
		fmt.Printf("  Description: %s\n", cluster.Description)
		fmt.Printf("  Items (%d):\n", cluster.Size)
		for _, idx := range cluster.Indices {
			if idx < len(articles) {
				fmt.Printf("    - %s\n", articles[idx])
			}
		}
	}
	fmt.Println()

	// Example 2: Cluster support tickets
	fmt.Println("--- Example 2: Support Ticket Clustering ---")
	type Ticket struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	tickets := []Ticket{
		{ID: 1, Title: "Login issue", Content: "Cannot log into my account, password reset not working"},
		{ID: 2, Title: "Slow performance", Content: "The application is very slow when loading dashboard"},
		{ID: 3, Title: "Payment failed", Content: "My credit card payment was declined"},
		{ID: 4, Title: "Account locked", Content: "My account has been locked after multiple login attempts"},
		{ID: 5, Title: "Page load timeout", Content: "Getting timeout errors when accessing reports"},
		{ID: 6, Title: "Refund request", Content: "I need a refund for my subscription"},
		{ID: 7, Title: "Password recovery", Content: "The password reset email never arrives"},
		{ID: 8, Title: "App crashes", Content: "Application crashes when generating large reports"},
	}

	ticketOpts := ops.NewClusterOptions().
		WithNumClusters(0). // Auto-detect number of clusters
		WithNamingStrategy("keyword").
		WithSimilarityThreshold(0.6)

	ticketResult, err := ops.Cluster(tickets, ticketOpts)
	if err != nil {
		log.Fatalf("Ticket clustering failed: %v", err)
	}

	fmt.Printf("Auto-detected %d ticket clusters:\n", len(ticketResult.Clusters))
	for _, cluster := range ticketResult.Clusters {
		fmt.Printf("\n  Category: %s\n", cluster.Name)
		fmt.Printf("  Priority keywords: %v\n", cluster.Keywords)
		fmt.Printf("  Ticket count: %d\n", cluster.Size)
	}

	// Example 3: Clustering with outlier detection
	fmt.Println("\n--- Example 3: Outlier Detection ---")
	if len(ticketResult.Outliers) > 0 {
		fmt.Printf("Found %d outlier items that don't fit clusters:\n", len(ticketResult.Outliers))
		for _, idx := range ticketResult.OutlierIndices {
			if idx < len(tickets) {
				fmt.Printf("  - Ticket %d: %s\n", tickets[idx].ID, tickets[idx].Title)
			}
		}
	} else {
		fmt.Println("No outliers detected - all items fit into clusters")
	}

	fmt.Println("\n=== Cluster Example Complete ===")
}
