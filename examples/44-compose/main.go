package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// CompanyProfile is the target composed object
type CompanyProfile struct {
	Name         string   `json:"name"`
	LegalName    string   `json:"legal_name"`
	Industry     string   `json:"industry"`
	Founded      int      `json:"founded"`
	Employees    int      `json:"employees"`
	Revenue      float64  `json:"revenue"`
	Website      string   `json:"website"`
	Headquarters string   `json:"headquarters"`
	CEO          string   `json:"ceo"`
	Description  string   `json:"description"`
	Products     []string `json:"products"`
	Competitors  []string `json:"competitors"`
	StockTicker  string   `json:"stock_ticker"`
	MarketCap    float64  `json:"market_cap"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Compose Example ===")

	// Example 1: Compose company profile from multiple sources
	fmt.Println("\n--- Example 1: Company Profile from Multiple Sources ---")

	// Data from LinkedIn/website
	webSource := map[string]any{
		"name":         "Acme Technologies",
		"website":      "https://acme.tech",
		"description":  "Leading provider of innovative tech solutions",
		"employees":    "500-1000", // Approximate
		"industry":     "Technology",
		"headquarters": "San Francisco, CA",
	}

	// Data from SEC filings
	secSource := map[string]any{
		"legal_name":   "Acme Technologies, Inc.",
		"founded":      2015,
		"revenue":      125000000,
		"employees":    847, // Exact count
		"stock_ticker": "ACME",
		"market_cap":   2500000000,
		"ceo":          "Jane Smith",
	}

	// Data from product database
	productSource := map[string]any{
		"products": []string{
			"AcmeCloud Platform",
			"AcmeAI Assistant",
			"AcmeSecure",
			"AcmeAnalytics",
		},
	}

	// Data from market research
	marketSource := map[string]any{
		"competitors": []string{
			"TechCorp",
			"InnovateLabs",
			"FutureSoft",
		},
		"industry":  "Enterprise Software", // More specific
		"employees": "approximately 800",   // Different estimate
	}

	parts := []any{webSource, secSource, productSource, marketSource}

	result, err := schemaflow.Assemble[CompanyProfile](parts, schemaflow.ComposeOptions{
		MergeStrategy: "smart",
		FillGaps:      false,
		Validate:      true,
		Steering:      "Prefer SEC data for official figures. Use most specific industry classification.",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Composition failed", "error", err)
		return
	}

	fmt.Printf("Composed Company Profile:\n")
	fmt.Printf("  Name: %s\n", result.Composed.Name)
	fmt.Printf("  Legal Name: %s\n", result.Composed.LegalName)
	fmt.Printf("  Industry: %s\n", result.Composed.Industry)
	fmt.Printf("  Founded: %d\n", result.Composed.Founded)
	fmt.Printf("  Employees: %d\n", result.Composed.Employees)
	fmt.Printf("  Revenue: $%.0f\n", result.Composed.Revenue)
	fmt.Printf("  Website: %s\n", result.Composed.Website)
	fmt.Printf("  Headquarters: %s\n", result.Composed.Headquarters)
	fmt.Printf("  CEO: %s\n", result.Composed.CEO)
	fmt.Printf("  Stock Ticker: %s\n", result.Composed.StockTicker)
	fmt.Printf("  Market Cap: $%.0f\n", result.Composed.MarketCap)
	fmt.Printf("  Products: %v\n", result.Composed.Products)
	fmt.Printf("  Competitors: %v\n", result.Composed.Competitors)

	fmt.Printf("\nField Sources:\n")
	for _, fs := range result.FieldSources {
		method := fs.Method
		if fs.Conflicts {
			method += " (conflict resolved)"
		}
		fmt.Printf("  %s: sources %v (%s)\n", fs.Field, fs.Sources, method)
		if fs.Resolution != "" {
			fmt.Printf("    → %s\n", fs.Resolution)
		}
	}

	fmt.Printf("\nConflicts Resolved: %d\n", result.ConflictsResolved)
	fmt.Printf("Completeness: %.0f%%\n", result.Completeness*100)

	if len(result.UnusedParts) > 0 {
		fmt.Printf("Unused Part Indices: %v\n", result.UnusedParts)
	}

	// Example 2: Compose document from sections
	fmt.Println("\n\n--- Example 2: Document Composition ---")

	type Document struct {
		Title        string   `json:"title"`
		Abstract     string   `json:"abstract"`
		Introduction string   `json:"introduction"`
		Methods      string   `json:"methods"`
		Results      string   `json:"results"`
		Discussion   string   `json:"discussion"`
		Conclusion   string   `json:"conclusion"`
		References   []string `json:"references"`
		Authors      []string `json:"authors"`
	}

	section1 := map[string]any{
		"title":   "Impact of AI on Software Development",
		"authors": []string{"Dr. Alice Johnson", "Prof. Bob Williams"},
	}

	section2 := map[string]any{
		"abstract":     "This study examines how AI tools are transforming modern software development practices...",
		"introduction": "The software development industry has witnessed significant changes with the advent of AI-powered tools...",
	}

	section3 := map[string]any{
		"methods": "We conducted a survey of 500 developers across 50 companies...",
		"results": "78% of developers reported increased productivity when using AI coding assistants...",
	}

	section4 := map[string]any{
		"discussion": "Our findings suggest that AI tools are most effective for boilerplate code generation...",
		"conclusion": "AI-powered development tools represent a paradigm shift in how software is created...",
		"references": []string{
			"Smith et al., 2023, AI in Software Engineering",
			"Jones, 2024, The Future of Coding",
		},
	}

	docParts := []any{section1, section2, section3, section4}

	docResult, err := schemaflow.Assemble[Document](docParts, schemaflow.ComposeOptions{
		MergeStrategy: "combine",
		Template:      "Academic research paper structure",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Document composition failed", "error", err)
		return
	}

	fmt.Printf("Composed Document:\n")
	fmt.Printf("  Title: %s\n", docResult.Composed.Title)
	fmt.Printf("  Authors: %v\n", docResult.Composed.Authors)
	fmt.Printf("  Abstract: %s...\n", truncate(docResult.Composed.Abstract, 60))
	fmt.Printf("  Introduction: %s...\n", truncate(docResult.Composed.Introduction, 60))
	fmt.Printf("  Methods: %s...\n", truncate(docResult.Composed.Methods, 60))
	fmt.Printf("  Results: %s...\n", truncate(docResult.Composed.Results, 60))
	fmt.Printf("  Discussion: %s...\n", truncate(docResult.Composed.Discussion, 60))
	fmt.Printf("  Conclusion: %s...\n", truncate(docResult.Composed.Conclusion, 60))
	fmt.Printf("  References: %d items\n", len(docResult.Composed.References))
	fmt.Printf("\nCompleteness: %.0f%%\n", docResult.Completeness*100)

	fmt.Println("\n=== Compose Example Complete ===")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
