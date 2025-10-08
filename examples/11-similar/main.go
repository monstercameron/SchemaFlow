package main

import (
	"fmt"
)

// NOTE: The Similar operation is defined in SchemaFlow but not yet implemented.
// This example demonstrates what the API would look like once implemented.
//
// Workaround: Use Compare operation with similarity focus, or Classify with categories.

func main() {
	fmt.Println("🔎 Similar Example - Duplicate Ticket Detection")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()
	fmt.Println("⚠️  STATUS: Similar operation not yet implemented")
	fmt.Println()
	fmt.Println("The Similar operation has options defined in ops/analysis.go")
	fmt.Println("but the implementation is pending.")
	fmt.Println()
	fmt.Println("📝 WORKAROUNDS:")
	fmt.Println("   1. Use Compare() with FocusOn=\"similarities\"")
	fmt.Println("   2. Use Classify() with categories")
	fmt.Println("   3. Use the Deduplicate() function in ops/extended.go")
	fmt.Println()
	fmt.Println("Example using Deduplicate (alternative):")
	fmt.Println("   result, err := ops.Deduplicate(tickets, 0.85)")
	fmt.Println()
	fmt.Println("See examples/10-compare for using Compare operation")
}
