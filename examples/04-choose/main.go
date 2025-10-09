package main

import (
	"fmt"
	"os"

	schemaflow "github.com/monstercameron/SchemaFlow"
	"github.com/monstercameron/SchemaFlow/core"
	"github.com/monstercameron/SchemaFlow/ops"
)

// Product represents a product in the catalog
type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Price       float64  `json:"price"`
	Features    []string `json:"features"`
	Description string   `json:"description"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		core.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		os.Exit(1)
	}

	// Product catalog
	products := []Product{
		{
			ID:          1,
			Name:        "UltraBook Pro",
			Category:    "Laptop",
			Price:       1299.99,
			Features:    []string{"16GB RAM", "512GB SSD", "13-inch display", "10hr battery"},
			Description: "Lightweight laptop perfect for professionals on the go",
		},
		{
			ID:          2,
			Name:        "PowerStation Desktop",
			Category:    "Desktop",
			Price:       1899.99,
			Features:    []string{"32GB RAM", "1TB SSD", "RTX 4070", "27-inch monitor"},
			Description: "High-performance desktop for gaming and content creation",
		},
		{
			ID:          3,
			Name:        "BudgetBook",
			Category:    "Laptop",
			Price:       499.99,
			Features:    []string{"8GB RAM", "256GB SSD", "14-inch display", "8hr battery"},
			Description: "Affordable laptop for students and basic computing",
		},
		{
			ID:          4,
			Name:        "CreatorPro Workstation",
			Category:    "Desktop",
			Price:       2499.99,
			Features:    []string{"64GB RAM", "2TB SSD", "RTX 4090", "Dual 4K monitors"},
			Description: "Professional workstation for video editing and 3D rendering",
		},
	}

	// Customer query
	customerQuery := "I'm a college student studying computer science. I need something portable for coding and note-taking. My budget is around $500-600."

	fmt.Println("🛍️ Choose Example - Product Recommender")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println("\n💬 Customer Query:")
	fmt.Println(customerQuery)
	fmt.Println("\n📦 Available Products:")
	for _, p := range products {
		fmt.Printf("  %d. %s - $%.2f (%s)\n", p.ID, p.Name, p.Price, p.Category)
	}

	// Choose the best product for the customer
	chooseOpts := ops.NewChooseOptions().WithCriteria([]string{customerQuery})
	chooseOpts.OpOptions.Intelligence = schemaflow.Smart
	chooseOpts.OpOptions.Steering = "Consider budget, portability, and intended use. Prioritize value for students."

	chosen, err := schemaflow.Choose(products, chooseOpts)

	if err != nil {
		core.GetLogger().Error("Selection failed", "error", err)
		os.Exit(1)
	}

	// Display recommendation
	fmt.Println("\n✅ Recommended Product:")
	fmt.Println("---")
	fmt.Printf("🎯 %s\n", chosen.Name)
	fmt.Printf("💰 Price: $%.2f\n", chosen.Price)
	fmt.Printf("📱 Category: %s\n", chosen.Category)
	fmt.Printf("✨ Features: %v\n", chosen.Features)
	fmt.Printf("\n📝 Description:\n   %s\n", chosen.Description)

	fmt.Println("\n💡 Why this choice?")
	fmt.Println("   - Within budget ($500-600)")
	fmt.Println("   - Portable laptop (not desktop)")
	fmt.Println("   - Suitable for coding and note-taking")
	fmt.Println("   - Good value for students")

	fmt.Println("\n✨ Success! Intelligently selected best match")
}
