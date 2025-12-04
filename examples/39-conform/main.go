package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// Address represents a mailing address
type Address struct {
	Name       string `json:"name"`
	Street1    string `json:"street1"`
	Street2    string `json:"street2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// ContactInfo with various formats
type ContactInfo struct {
	Phone     string `json:"phone"`
	AltPhone  string `json:"alt_phone"`
	Fax       string `json:"fax"`
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
	JoinDate  string `json:"join_date"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Conform Example ===")

	// Example 1: Conform address to USPS standard
	fmt.Println("\n--- Example 1: USPS Address Standardization ---")

	rawAddress := Address{
		Name:       "john smith",
		Street1:    "123 north main street, apt 4b",
		Street2:    "",
		City:       "los angeles",
		State:      "california",
		PostalCode: "90210",
		Country:    "usa",
	}

	fmt.Printf("Original Address:\n")
	fmt.Printf("  Name: %s\n", rawAddress.Name)
	fmt.Printf("  Street: %s\n", rawAddress.Street1)
	fmt.Printf("  City, State ZIP: %s, %s %s\n", rawAddress.City, rawAddress.State, rawAddress.PostalCode)
	fmt.Printf("  Country: %s\n", rawAddress.Country)

	result, err := schemaflow.Conform[Address](rawAddress, "USPS", schemaflow.ConformOptions{
		Strict: true,
	})

	if err != nil {
		schemaflow.GetLogger().Error("Conform failed", "error", err)
		return
	}

	fmt.Printf("\nUSPS Standardized:\n")
	fmt.Printf("  Name: %s\n", result.Conformed.Name)
	fmt.Printf("  Street1: %s\n", result.Conformed.Street1)
	if result.Conformed.Street2 != "" {
		fmt.Printf("  Street2: %s\n", result.Conformed.Street2)
	}
	fmt.Printf("  City, State ZIP: %s, %s %s\n", result.Conformed.City, result.Conformed.State, result.Conformed.PostalCode)
	fmt.Printf("  Country: %s\n", result.Conformed.Country)

	fmt.Printf("\nAdjustments Made:\n")
	for _, adj := range result.Adjustments {
		fmt.Printf("  %s: '%v' → '%v'\n", adj.Field, adj.OriginalValue, adj.ConformedValue)
		fmt.Printf("    Rule: %s\n", adj.Description)
	}

	fmt.Printf("\nCompliance: %.0f%%\n", result.Compliance*100)

	// Example 2: Conform contact info to E164 and ISO8601
	fmt.Println("\n--- Example 2: Phone (E164) & Date (ISO8601) ---")

	rawContact := ContactInfo{
		Phone:     "(555) 123-4567",
		AltPhone:  "1-800-FLOWERS",
		Fax:       "555.999.8888",
		Email:     "John.Smith@Email.COM",
		Birthdate: "March 15th, 1990",
		JoinDate:  "01/20/2024",
	}

	fmt.Printf("Original Contact:\n")
	fmt.Printf("  Phone: %s\n", rawContact.Phone)
	fmt.Printf("  Alt Phone: %s\n", rawContact.AltPhone)
	fmt.Printf("  Fax: %s\n", rawContact.Fax)
	fmt.Printf("  Email: %s\n", rawContact.Email)
	fmt.Printf("  Birthdate: %s\n", rawContact.Birthdate)
	fmt.Printf("  Join Date: %s\n", rawContact.JoinDate)

	// Conform to E164 for phones
	phoneResult, err := schemaflow.Conform[ContactInfo](rawContact, "E164+ISO8601", schemaflow.ConformOptions{
		Steering: "Assume US country code +1 for phone numbers. Convert dates to ISO8601 format.",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Phone conform failed", "error", err)
		return
	}

	fmt.Printf("\nStandardized Contact:\n")
	fmt.Printf("  Phone: %s\n", phoneResult.Conformed.Phone)
	fmt.Printf("  Alt Phone: %s\n", phoneResult.Conformed.AltPhone)
	fmt.Printf("  Fax: %s\n", phoneResult.Conformed.Fax)
	fmt.Printf("  Email: %s\n", phoneResult.Conformed.Email)
	fmt.Printf("  Birthdate: %s\n", phoneResult.Conformed.Birthdate)
	fmt.Printf("  Join Date: %s\n", phoneResult.Conformed.JoinDate)

	fmt.Printf("\nAdjustments:\n")
	for _, adj := range phoneResult.Adjustments {
		fmt.Printf("  %s: %s\n", adj.Field, adj.Description)
	}

	// Example 3: Custom standard conformance
	fmt.Println("\n--- Example 3: Custom Standard ---")

	type ProductCode struct {
		SKU         string `json:"sku"`
		Barcode     string `json:"barcode"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}

	rawProduct := ProductCode{
		SKU:         "electronics-laptop-dell-001",
		Barcode:     "1234567890",
		Category:    "LAPTOPS & COMPUTERS",
		Description: "Dell Laptop 15 inch",
	}

	customStandard := `
	SKU Format: CAT-SUB-BRAND-NUM (uppercase, max 20 chars)
	Barcode: EAN-13 format (13 digits with check digit)
	Category: Title case, no special characters
	Description: Sentence case, max 50 chars
	`

	customResult, err := schemaflow.Conform[ProductCode](rawProduct, customStandard, schemaflow.ConformOptions{
		Strict: false,
	})

	if err != nil {
		schemaflow.GetLogger().Error("Custom conform failed", "error", err)
		return
	}

	fmt.Printf("Original:\n")
	fmt.Printf("  SKU: %s\n", rawProduct.SKU)
	fmt.Printf("  Barcode: %s\n", rawProduct.Barcode)
	fmt.Printf("  Category: %s\n", rawProduct.Category)
	fmt.Printf("  Description: %s\n", rawProduct.Description)

	fmt.Printf("\nConformed:\n")
	fmt.Printf("  SKU: %s\n", customResult.Conformed.SKU)
	fmt.Printf("  Barcode: %s\n", customResult.Conformed.Barcode)
	fmt.Printf("  Category: %s\n", customResult.Conformed.Category)
	fmt.Printf("  Description: %s\n", customResult.Conformed.Description)

	fmt.Printf("\nCompliance: %.0f%%\n", customResult.Compliance*100)

	fmt.Println("\n=== Conform Example Complete ===")
}
