package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// BasicPerson is the input type with minimal information
type BasicPerson struct {
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	City      string `json:"city"`
	JobTitle  string `json:"job_title"`
}

// EnrichedPerson has derived fields
type EnrichedPerson struct {
	// Original fields
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	City      string `json:"city"`
	JobTitle  string `json:"job_title"`

	// Derived fields
	Age            int    `json:"age"`
	AgeCategory    string `json:"age_category"`    // "youth", "adult", "senior"
	Generation     string `json:"generation"`      // "GenZ", "Millennial", etc.
	Region         string `json:"region"`          // Derived from city
	Timezone       string `json:"timezone"`        // Derived from city
	JobCategory    string `json:"job_category"`    // e.g., "tech", "healthcare"
	SeniorityLevel string `json:"seniority_level"` // "entry", "mid", "senior", "executive"
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Derive Example ===")

	// Example 1: Derive enriched person data
	fmt.Println("\n--- Example 1: Person Enrichment ---")

	person := BasicPerson{
		Name:      "Sarah Chen",
		BirthYear: 1992,
		City:      "San Francisco",
		JobTitle:  "Senior Software Engineer",
	}

	result, err := schemaflow.Derive[BasicPerson, EnrichedPerson](person, schemaflow.DeriveOptions{
		Fields:   []string{"age", "age_category", "generation", "region", "timezone", "job_category", "seniority_level"},
		Steering: "Use current year 2024 for age calculation",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Derivation failed", "error", err)
		return
	}

	fmt.Printf("Original Data:\n")
	fmt.Printf("  Name: %s\n", person.Name)
	fmt.Printf("  Birth Year: %d\n", person.BirthYear)
	fmt.Printf("  City: %s\n", person.City)
	fmt.Printf("  Job Title: %s\n", person.JobTitle)

	fmt.Printf("\nDerived Data:\n")
	fmt.Printf("  Age: %d\n", result.Derived.Age)
	fmt.Printf("  Age Category: %s\n", result.Derived.AgeCategory)
	fmt.Printf("  Generation: %s\n", result.Derived.Generation)
	fmt.Printf("  Region: %s\n", result.Derived.Region)
	fmt.Printf("  Timezone: %s\n", result.Derived.Timezone)
	fmt.Printf("  Job Category: %s\n", result.Derived.JobCategory)
	fmt.Printf("  Seniority Level: %s\n", result.Derived.SeniorityLevel)

	fmt.Printf("\nDerivation Details:\n")
	for _, d := range result.Derivations {
		fmt.Printf("  %s: %s (confidence: %.0f%%)\n", d.Field, d.Method, d.Confidence*100)
		if d.Reasoning != "" {
			fmt.Printf("    Reasoning: %s\n", d.Reasoning)
		}
	}

	// Example 2: Derive transaction risk features
	fmt.Println("\n--- Example 2: Transaction Risk Features ---")

	type RawTransaction struct {
		ID        string  `json:"id"`
		Amount    float64 `json:"amount"`
		Merchant  string  `json:"merchant"`
		Category  string  `json:"category"`
		Timestamp string  `json:"timestamp"`
		CardLast4 string  `json:"card_last4"`
		Location  string  `json:"location"`
		IsOnline  bool    `json:"is_online"`
	}

	type RiskEnrichedTransaction struct {
		// Original
		ID        string  `json:"id"`
		Amount    float64 `json:"amount"`
		Merchant  string  `json:"merchant"`
		Category  string  `json:"category"`
		Timestamp string  `json:"timestamp"`
		CardLast4 string  `json:"card_last4"`
		Location  string  `json:"location"`
		IsOnline  bool    `json:"is_online"`

		// Derived risk features
		AmountCategory string   `json:"amount_category"` // "low", "medium", "high", "very_high"
		TimeOfDay      string   `json:"time_of_day"`     // "morning", "afternoon", "evening", "night"
		IsWeekend      bool     `json:"is_weekend"`
		MerchantType   string   `json:"merchant_type"` // "known_retailer", "small_business", "unknown"
		RiskScore      float64  `json:"risk_score"`    // 0.0-1.0
		RiskFactors    []string `json:"risk_factors"`
	}

	transaction := RawTransaction{
		ID:        "TXN-123456",
		Amount:    2499.99,
		Merchant:  "ELECTROZONE ONLINE",
		Category:  "Electronics",
		Timestamp: "2024-01-15T02:30:00Z",
		CardLast4: "4532",
		Location:  "Lagos, Nigeria",
		IsOnline:  true,
	}

	txnResult, err := schemaflow.Derive[RawTransaction, RiskEnrichedTransaction](transaction, schemaflow.DeriveOptions{
		Steering: "Consider typical fraud patterns when deriving risk score and factors",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Transaction derivation failed", "error", err)
		return
	}

	fmt.Printf("Transaction: %s\n", transaction.ID)
	fmt.Printf("  Amount: $%.2f\n", transaction.Amount)
	fmt.Printf("  Merchant: %s\n", transaction.Merchant)

	fmt.Printf("\nDerived Risk Features:\n")
	fmt.Printf("  Amount Category: %s\n", txnResult.Derived.AmountCategory)
	fmt.Printf("  Time of Day: %s\n", txnResult.Derived.TimeOfDay)
	fmt.Printf("  Is Weekend: %v\n", txnResult.Derived.IsWeekend)
	fmt.Printf("  Merchant Type: %s\n", txnResult.Derived.MerchantType)
	fmt.Printf("  Risk Score: %.2f\n", txnResult.Derived.RiskScore)
	fmt.Printf("  Risk Factors:\n")
	for _, factor := range txnResult.Derived.RiskFactors {
		fmt.Printf("    - %s\n", factor)
	}

	fmt.Printf("\nOverall Confidence: %.0f%%\n", txnResult.OverallConfidence*100)

	fmt.Println("\n=== Derive Example Complete ===")
}
