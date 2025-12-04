package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// SalesRow represents individual sales records
type SalesRow struct {
	Product string  `json:"product"`
	Region  string  `json:"region"`
	Quarter string  `json:"quarter"`
	Revenue float64 `json:"revenue"`
	Units   int     `json:"units"`
}

// ProductSummary is pivoted view by product with quarters as columns
type ProductSummary struct {
	Product   string  `json:"product"`
	Q1Revenue float64 `json:"q1_revenue"`
	Q2Revenue float64 `json:"q2_revenue"`
	Q3Revenue float64 `json:"q3_revenue"`
	Q4Revenue float64 `json:"q4_revenue"`
	TotalYear float64 `json:"total_year"`
}

// NestedUser represents deeply nested user data
type NestedUser struct {
	ID      string `json:"id"`
	Profile struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Settings struct {
			Theme    string `json:"theme"`
			Language string `json:"language"`
		} `json:"settings"`
	} `json:"profile"`
	Address struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Country string `json:"country"`
	} `json:"address"`
	Subscription struct {
		Plan    string  `json:"plan"`
		Price   float64 `json:"price"`
		Renewal string  `json:"renewal"`
	} `json:"subscription"`
}

// FlatUser is the flattened version
type FlatUser struct {
	ID                  string  `json:"id"`
	ProfileName         string  `json:"profile_name"`
	ProfileEmail        string  `json:"profile_email"`
	SettingsTheme       string  `json:"settings_theme"`
	SettingsLanguage    string  `json:"settings_language"`
	AddressStreet       string  `json:"address_street"`
	AddressCity         string  `json:"address_city"`
	AddressCountry      string  `json:"address_country"`
	SubscriptionPlan    string  `json:"subscription_plan"`
	SubscriptionPrice   float64 `json:"subscription_price"`
	SubscriptionRenewal string  `json:"subscription_renewal"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Pivot Example ===")

	// Example 1: Pivot sales rows to quarterly columns
	fmt.Println("\n--- Example 1: Sales Data Pivot ---")

	salesData := []SalesRow{
		{Product: "Widget A", Region: "North", Quarter: "Q1", Revenue: 10000, Units: 100},
		{Product: "Widget A", Region: "North", Quarter: "Q2", Revenue: 12000, Units: 120},
		{Product: "Widget A", Region: "North", Quarter: "Q3", Revenue: 11000, Units: 110},
		{Product: "Widget A", Region: "North", Quarter: "Q4", Revenue: 15000, Units: 150},
		{Product: "Widget B", Region: "South", Quarter: "Q1", Revenue: 8000, Units: 80},
		{Product: "Widget B", Region: "South", Quarter: "Q2", Revenue: 9500, Units: 95},
		{Product: "Widget B", Region: "South", Quarter: "Q3", Revenue: 8500, Units: 85},
		{Product: "Widget B", Region: "South", Quarter: "Q4", Revenue: 12000, Units: 120},
		{Product: "Gadget X", Region: "East", Quarter: "Q1", Revenue: 20000, Units: 50},
		{Product: "Gadget X", Region: "East", Quarter: "Q2", Revenue: 22000, Units: 55},
		{Product: "Gadget X", Region: "East", Quarter: "Q3", Revenue: 25000, Units: 62},
		{Product: "Gadget X", Region: "East", Quarter: "Q4", Revenue: 28000, Units: 70},
	}

	fmt.Println("Original Data (sample):")
	for i := 0; i < 4 && i < len(salesData); i++ {
		s := salesData[i]
		fmt.Printf("  %s | %s | %s | $%.0f\n", s.Product, s.Region, s.Quarter, s.Revenue)
	}
	fmt.Println("  ...")

	result, err := schemaflow.Pivot[[]SalesRow, []ProductSummary](salesData, schemaflow.PivotOptions{
		PivotOn:   []string{"Quarter"},
		GroupBy:   []string{"Product"},
		Aggregate: "sum",
		Steering:  "Sum revenue by quarter for each product. Calculate total_year as sum of all quarters.",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Pivot failed", "error", err)
		return
	}

	fmt.Println("\nPivoted Data:")
	fmt.Println("  Product    | Q1      | Q2      | Q3      | Q4      | Total")
	fmt.Println("  -----------|---------|---------|---------|---------|----------")
	for _, p := range result.Pivoted {
		fmt.Printf("  %-10s | $%-6.0f | $%-6.0f | $%-6.0f | $%-6.0f | $%.0f\n",
			p.Product, p.Q1Revenue, p.Q2Revenue, p.Q3Revenue, p.Q4Revenue, p.TotalYear)
	}

	fmt.Printf("\nPivot Stats:\n")
	fmt.Printf("  Source Fields: %d\n", result.Stats.SourceFields)
	fmt.Printf("  Target Fields: %d\n", result.Stats.TargetFields)
	fmt.Printf("  Compressions: %d\n", result.Stats.Compressions)

	// Example 2: Flatten nested structure
	fmt.Println("\n\n--- Example 2: Flatten Nested Structure ---")

	nestedUser := NestedUser{
		ID: "USR-123",
	}
	nestedUser.Profile.Name = "John Doe"
	nestedUser.Profile.Email = "john@example.com"
	nestedUser.Profile.Settings.Theme = "dark"
	nestedUser.Profile.Settings.Language = "en-US"
	nestedUser.Address.Street = "123 Main St"
	nestedUser.Address.City = "Boston"
	nestedUser.Address.Country = "USA"
	nestedUser.Subscription.Plan = "Premium"
	nestedUser.Subscription.Price = 29.99
	nestedUser.Subscription.Renewal = "2024-12-01"

	fmt.Println("Nested Structure:")
	fmt.Println("  user:")
	fmt.Println("    id: USR-123")
	fmt.Println("    profile:")
	fmt.Println("      name: John Doe")
	fmt.Println("      email: john@example.com")
	fmt.Println("      settings:")
	fmt.Println("        theme: dark")
	fmt.Println("        language: en-US")
	fmt.Println("    address:")
	fmt.Println("      street: 123 Main St")
	fmt.Println("      city: Boston")
	fmt.Println("      ...")

	flatResult, err := schemaflow.Pivot[NestedUser, FlatUser](nestedUser, schemaflow.PivotOptions{
		Flatten: true,
	})

	if err != nil {
		schemaflow.GetLogger().Error("Flatten failed", "error", err)
		return
	}

	fmt.Println("\nFlattened Structure:")
	fmt.Printf("  id: %s\n", flatResult.Pivoted.ID)
	fmt.Printf("  profile_name: %s\n", flatResult.Pivoted.ProfileName)
	fmt.Printf("  profile_email: %s\n", flatResult.Pivoted.ProfileEmail)
	fmt.Printf("  settings_theme: %s\n", flatResult.Pivoted.SettingsTheme)
	fmt.Printf("  settings_language: %s\n", flatResult.Pivoted.SettingsLanguage)
	fmt.Printf("  address_street: %s\n", flatResult.Pivoted.AddressStreet)
	fmt.Printf("  address_city: %s\n", flatResult.Pivoted.AddressCity)
	fmt.Printf("  address_country: %s\n", flatResult.Pivoted.AddressCountry)
	fmt.Printf("  subscription_plan: %s\n", flatResult.Pivoted.SubscriptionPlan)
	fmt.Printf("  subscription_price: %.2f\n", flatResult.Pivoted.SubscriptionPrice)
	fmt.Printf("  subscription_renewal: %s\n", flatResult.Pivoted.SubscriptionRenewal)

	fmt.Printf("\nTransformation:\n")
	fmt.Printf("  Depth Change: %d (shallower)\n", flatResult.Stats.DepthChange)
	fmt.Printf("  Expansions: %d\n", flatResult.Stats.Expansions)

	if len(flatResult.DataLoss) > 0 {
		fmt.Printf("  Data Loss: %v\n", flatResult.DataLoss)
	} else {
		fmt.Println("  Data Loss: none")
	}

	fmt.Println("\n=== Pivot Example Complete ===")
}
