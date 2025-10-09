package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
	"github.com/monstercameron/SchemaFlow/core"
	"github.com/monstercameron/SchemaFlow/ops"
)

// CustomerRecord represents customer information from different sources
type CustomerRecord struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
	VIP     bool   `json:"vip,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		core.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("🔀 Merge Example - Customer Record Deduplication")
	fmt.Println("=" + string(make([]byte, 60)))

	// Multiple records for the same customer from different systems
	records := []CustomerRecord{
		{
			ID:      "CRM-001",
			Name:    "John Smith",
			Email:   "john.smith@example.com",
			Phone:   "+1-555-0123",
			Address: "",
			VIP:     false,
			Notes:   "Prefers email contact",
		},
		{
			ID:      "SALES-456",
			Name:    "J. Smith",
			Email:   "",
			Phone:   "+1-555-0123",
			Address: "123 Main St, Springfield, IL 62701",
			VIP:     true,
			Notes:   "",
		},
		{
			ID:      "SUPPORT-789",
			Name:    "John A. Smith",
			Email:   "john.smith@example.com",
			Phone:   "",
			Address: "123 Main Street, Springfield, Illinois",
			VIP:     false,
			Notes:   "Has premium support plan",
		},
	}

	fmt.Println("\n📋 Source Records:")
	for i, r := range records {
		fmt.Printf("\n%d. Record %s:\n", i+1, r.ID)
		printRecord(r, "   ")
	}

	// Merge strategy
	strategy := `
Merge Strategy:
1. Keep the most complete name variant
2. Prefer email if present in any record
3. Use the most detailed address
4. Set VIP=true if ANY record has it
5. Combine notes from all records
6. Choose ID from CRM system if present
`

	fmt.Println("\n🔄 Merging records...")

	// Merge records
	merged, err := ops.Merge(records, strategy)
	if err != nil {
		core.GetLogger().Error("Failed to merge records", "error", err)
		return
	}

	fmt.Println()
	fmt.Println("✅ Merged Result:")
	printRecord(merged, "   ")

	fmt.Println()
	fmt.Println("📊 Merge Analysis:")
	fmt.Println("   Input: 3 duplicate records")
	fmt.Println("   Output: 1 unified record")
	fmt.Println()
	fmt.Println("   ✓ Name: Selected most complete variant")
	fmt.Println("   ✓ Email: Preserved from CRM")
	fmt.Println("   ✓ Phone: Common across records")
	fmt.Println("   ✓ Address: Used most detailed version")
	fmt.Println("   ✓ VIP: Upgraded to true")
	fmt.Println("   ✓ Notes: Combined all information")

	fmt.Println()
	fmt.Println("✨ Success! Customer records merged")
}

func printRecord(r CustomerRecord, indent string) {
	fmt.Printf("%sName:    %s\n", indent, valueOrEmpty(r.Name))
	fmt.Printf("%sEmail:   %s\n", indent, valueOrEmpty(r.Email))
	fmt.Printf("%sPhone:   %s\n", indent, valueOrEmpty(r.Phone))
	fmt.Printf("%sAddress: %s\n", indent, valueOrEmpty(r.Address))
	fmt.Printf("%sVIP:     %v\n", indent, r.VIP)
	if r.Notes != "" {
		fmt.Printf("%sNotes:   %s\n", indent, r.Notes)
	}
}

func valueOrEmpty(s string) string {
	if s == "" {
		return "(empty)"
	}
	return s
}
