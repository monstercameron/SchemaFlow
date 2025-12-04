package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// CustomerRecord represents a customer from different sources
type CustomerRecord struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	LastContact string `json:"last_contact"`
	Status      string `json:"status"`
	Source      string `json:"source"` // Track where this came from
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Resolve Example ===")

	// Example 1: Resolve conflicting customer records from multiple sources
	fmt.Println("\n--- Example 1: Customer Record Deduplication ---")

	// Same customer from different systems with conflicting data
	sources := []CustomerRecord{
		{
			ID:          "CUST-001",
			Name:        "John Smith",
			Email:       "john.smith@email.com",
			Phone:       "555-123-4567",
			Address:     "123 Main St, Boston, MA",
			LastContact: "2024-01-15",
			Status:      "active",
			Source:      "CRM",
		},
		{
			ID:          "CUST-001",
			Name:        "John A. Smith",
			Email:       "jsmith@work.com",
			Phone:       "(555) 123-4567",
			Address:     "123 Main Street, Boston, MA 02101",
			LastContact: "2024-01-20",
			Status:      "active",
			Source:      "Sales DB",
		},
		{
			ID:          "CUST-001",
			Name:        "John Smith",
			Email:       "john.smith@email.com",
			Phone:       "",
			Address:     "456 Oak Ave, Boston, MA",
			LastContact: "2023-12-01",
			Status:      "inactive",
			Source:      "Marketing",
		},
	}

	result, err := schemaflow.Resolve[CustomerRecord](sources, schemaflow.ResolveOptions{
		Strategy: "most-complete",
		Steering: "Prefer more recent data for contact info, most complete for address",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Resolution failed", "error", err)
		return
	}

	fmt.Printf("Resolved Record:\n")
	fmt.Printf("  ID: %s\n", result.Resolved.ID)
	fmt.Printf("  Name: %s\n", result.Resolved.Name)
	fmt.Printf("  Email: %s\n", result.Resolved.Email)
	fmt.Printf("  Phone: %s\n", result.Resolved.Phone)
	fmt.Printf("  Address: %s\n", result.Resolved.Address)
	fmt.Printf("  Last Contact: %s\n", result.Resolved.LastContact)
	fmt.Printf("  Status: %s\n", result.Resolved.Status)

	fmt.Printf("\nConflicts Resolved (%d):\n", len(result.Conflicts))
	for _, conflict := range result.Conflicts {
		fmt.Printf("  Field: %s\n", conflict.Field)
		fmt.Printf("    Values: %v\n", conflict.Values)
		fmt.Printf("    Resolution: %s\n", conflict.Resolution)
		fmt.Printf("    Chosen (source %d): %v\n", conflict.ChosenSource, conflict.ChosenValue)
	}

	fmt.Printf("\nSource Contributions:\n")
	for source, fields := range result.SourceContributions {
		fmt.Printf("  Source %d: %v\n", source, fields)
	}

	fmt.Printf("\nConfidence: %.0f%%\n", result.Confidence*100)

	// Example 2: Resolve product information from multiple vendors
	fmt.Println("\n--- Example 2: Product Information ---")

	type ProductInfo struct {
		SKU          string   `json:"sku"`
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		Price        float64  `json:"price"`
		Stock        int      `json:"stock"`
		Manufacturer string   `json:"manufacturer"`
		Categories   []string `json:"categories"`
		Source       string   `json:"source"`
	}

	productSources := []ProductInfo{
		{
			SKU:          "PROD-789",
			Name:         "Wireless Bluetooth Headphones",
			Description:  "High-quality wireless headphones",
			Price:        79.99,
			Stock:        150,
			Manufacturer: "AudioTech",
			Categories:   []string{"electronics", "audio"},
			Source:       "Warehouse A",
		},
		{
			SKU:          "PROD-789",
			Name:         "Wireless Bluetooth Headphones Pro",
			Description:  "Premium wireless headphones with noise cancellation and 30-hour battery",
			Price:        89.99,
			Stock:        75,
			Manufacturer: "AudioTech Inc.",
			Categories:   []string{"electronics", "audio", "wireless"},
			Source:       "Warehouse B",
		},
		{
			SKU:          "PROD-789",
			Name:         "BT Headphones",
			Description:  "",
			Price:        79.99,
			Stock:        200,
			Manufacturer: "AudioTech",
			Categories:   []string{"audio"},
			Source:       "Vendor Feed",
		},
	}

	productResult, err := schemaflow.Resolve[ProductInfo](productSources, schemaflow.ResolveOptions{
		Strategy: "most-complete",
		FieldPriorities: map[string]int{
			"price": 0, // Prefer Warehouse A for price
		},
	})

	if err != nil {
		schemaflow.GetLogger().Error("Product resolution failed", "error", err)
		return
	}

	fmt.Printf("Resolved Product:\n")
	fmt.Printf("  SKU: %s\n", productResult.Resolved.SKU)
	fmt.Printf("  Name: %s\n", productResult.Resolved.Name)
	fmt.Printf("  Description: %s\n", productResult.Resolved.Description)
	fmt.Printf("  Price: $%.2f\n", productResult.Resolved.Price)
	fmt.Printf("  Stock: %d\n", productResult.Resolved.Stock)
	fmt.Printf("  Manufacturer: %s\n", productResult.Resolved.Manufacturer)
	fmt.Printf("  Categories: %v\n", productResult.Resolved.Categories)

	fmt.Printf("\nConflicts: %d resolved\n", len(productResult.Conflicts))
	for _, c := range productResult.Conflicts {
		fmt.Printf("  - %s: %s\n", c.Field, c.Resolution)
	}

	fmt.Println("\n=== Resolve Example Complete ===")
}
