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

	fmt.Println("=== Normalize Example ===")

	// Example 1: Normalize address data
	fmt.Println("--- Example 1: Address Normalization ---")
	type Address struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		State   string `json:"state"`
		Country string `json:"country"`
		Zip     string `json:"zip"`
	}

	address := Address{
		Street:  "123 Main St.",
		City:    "new york city",
		State:   "NY",
		Country: "USA",
		Zip:     "10001",
	}

	opts := ops.NewNormalizeOptions().
		WithRules(map[string]string{
			"street":  "Expand abbreviations (St. -> Street)",
			"city":    "Proper case, canonical name",
			"country": "Full country name",
		})

	result, err := ops.Normalize(address, opts)
	if err != nil {
		log.Fatalf("Normalization failed: %v", err)
	}

	fmt.Println("Original:")
	fmt.Printf("  %s, %s, %s, %s %s\n", address.Street, address.City, address.State, address.Country, address.Zip)
	fmt.Println("\nNormalized:")
	fmt.Printf("  %s, %s, %s, %s %s\n",
		result.Normalized.Street,
		result.Normalized.City,
		result.Normalized.State,
		result.Normalized.Country,
		result.Normalized.Zip)

	if len(result.Changes) > 0 {
		fmt.Println("\nChanges made:")
		for _, change := range result.Changes {
			fmt.Printf("  - %s: '%s' -> '%s' (%s)\n", change.Field, change.Original, change.Normalized, change.Reason)
		}
	}
	fmt.Println()

	// Example 2: Normalize dates
	fmt.Println("--- Example 2: Date Normalization ---")
	type Event struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	events := []Event{
		{Name: "Conference", StartDate: "Jan 15, 2024", EndDate: "January 17th, 2024"},
		{Name: "Workshop", StartDate: "2024/02/20", EndDate: "02-21-2024"},
		{Name: "Meetup", StartDate: "March 5th", EndDate: "3/5/24"},
	}

	dateOpts := ops.NewNormalizeOptions().
		WithStandard("ISO 8601").
		WithRules(map[string]string{
			"start_date": "YYYY-MM-DD format",
			"end_date":   "YYYY-MM-DD format",
		})

	results, err := ops.NormalizeBatch(events, dateOpts)
	if err != nil {
		log.Fatalf("Date normalization failed: %v", err)
	}

	fmt.Println("Normalized dates:")
	for i, r := range results {
		fmt.Printf("  %s: %s to %s\n",
			r.Normalized.Name,
			r.Normalized.StartDate,
			r.Normalized.EndDate)
		if i < len(events) {
			fmt.Printf("    (was: %s to %s)\n", events[i].StartDate, events[i].EndDate)
		}
	}
	fmt.Println()

	// Example 3: Normalize text with abbreviations
	fmt.Println("--- Example 3: Text Normalization ---")
	text := "The mtg is tmrw @ 3pm in conf rm 2B. Pls confirm ASAP. Thx!"

	textOpts := ops.NewNormalizeOptions().
		WithRules(map[string]string{
			"abbreviations": "Expand to full words",
			"punctuation":   "Standard punctuation",
		}).
		WithFixTypos(true)

	normalizedText, err := ops.NormalizeText(text, textOpts)
	if err != nil {
		log.Fatalf("Text normalization failed: %v", err)
	}

	fmt.Printf("Original: %s\n", text)
	fmt.Printf("Normalized: %s\n\n", normalizedText)

	// Example 4: Normalize with canonical mappings
	fmt.Println("--- Example 4: Canonical Mappings ---")
	type Company struct {
		Name     string `json:"name"`
		Industry string `json:"industry"`
		Country  string `json:"country"`
	}

	companies := []Company{
		{Name: "Microsoft Corp", Industry: "Tech", Country: "US"},
		{Name: "Alphabet Inc.", Industry: "Technology", Country: "United States"},
		{Name: "Meta Platforms", Industry: "tech/social", Country: "USA"},
	}

	canonicalOpts := ops.NewNormalizeOptions().
		WithCanonicalMappings(map[string]string{
			"US":          "United States",
			"USA":         "United States",
			"Tech":        "Technology",
			"tech":        "Technology",
			"tech/social": "Technology",
		}).
		WithRules(map[string]string{
			"industry": "Use standard industry classification",
		})

	companyResults, err := ops.NormalizeBatch(companies, canonicalOpts)
	if err != nil {
		log.Fatalf("Company normalization failed: %v", err)
	}

	fmt.Println("Normalized companies:")
	for _, r := range companyResults {
		fmt.Printf("  %s | %s | %s\n",
			r.Normalized.Name,
			r.Normalized.Industry,
			r.Normalized.Country)
	}

	fmt.Println("\n=== Normalize Example Complete ===")
}
