package main

import (
	"fmt"
	"os"

	schemaflow "github.com/monstercameron/SchemaFlow"
	"github.com/monstercameron/SchemaFlow/core"
	"github.com/monstercameron/SchemaFlow/ops"
)

// Order represents an e-commerce order
type Order struct {
	ID              string
	Status          string
	Items           int
	Total           float64
	PaymentReceived bool
	ShippingAddress string
	TrackingNumber  string
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		core.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		os.Exit(1)
	}

	fmt.Println("🛡️ Guard Example - Order State Validation")
	fmt.Println("=" + string(make([]byte, 60)))

	// Test orders in different states
	testCases := []struct {
		name  string
		order Order
	}{
		{
			name: "Valid Order - Ready to Ship",
			order: Order{
				ID:              "ORD-1001",
				Status:          "processing",
				Items:           3,
				Total:           149.99,
				PaymentReceived: true,
				ShippingAddress: "123 Main St, City, State 12345",
				TrackingNumber:  "",
			},
		},
		{
			name: "Invalid - No Payment",
			order: Order{
				ID:              "ORD-1002",
				Status:          "processing",
				Items:           2,
				Total:           89.50,
				PaymentReceived: false,
				ShippingAddress: "456 Oak Ave, Town, State 54321",
				TrackingNumber:  "",
			},
		},
		{
			name: "Invalid - Missing Address",
			order: Order{
				ID:              "ORD-1003",
				Status:          "processing",
				Items:           1,
				Total:           29.99,
				PaymentReceived: true,
				ShippingAddress: "",
				TrackingNumber:  "",
			},
		},
		{
			name: "Invalid - Already Shipped",
			order: Order{
				ID:              "ORD-1004",
				Status:          "shipped",
				Items:           5,
				Total:           299.99,
				PaymentReceived: true,
				ShippingAddress: "789 Pine Rd, Village, State 98765",
				TrackingNumber:  "TRK123456789",
			},
		},
	}

	// Define guard checks for shipping
	checks := []func(Order) (bool, string){
		func(o Order) (bool, string) {
			if !o.PaymentReceived {
				return false, "Payment not received"
			}
			return true, "Payment confirmed"
		},
		func(o Order) (bool, string) {
			if o.ShippingAddress == "" {
				return false, "Shipping address missing"
			}
			return true, "Shipping address present"
		},
		func(o Order) (bool, string) {
			if o.Items <= 0 {
				return false, "No items in order"
			}
			return true, fmt.Sprintf("%d items ready", o.Items)
		},
		func(o Order) (bool, string) {
			if o.Status == "shipped" || o.Status == "delivered" {
				return false, "Order already shipped"
			}
			return true, "Order status valid for shipping"
		},
		func(o Order) (bool, string) {
			if o.Total <= 0 {
				return false, "Invalid order total"
			}
			return true, fmt.Sprintf("Order total: $%.2f", o.Total)
		},
	}

	// Test each order
	for i, tc := range testCases {
		fmt.Printf("\n%d. %s\n", i+1, tc.name)
		fmt.Println("---")
		fmt.Printf("   Order ID: %s\n", tc.order.ID)
		fmt.Printf("   Status: %s\n", tc.order.Status)
		fmt.Printf("   Items: %d\n", tc.order.Items)
		fmt.Printf("   Total: $%.2f\n", tc.order.Total)
		fmt.Printf("   Payment: %v\n", tc.order.PaymentReceived)
		fmt.Printf("   Address: %s\n", valueOrMissing(tc.order.ShippingAddress))

		fmt.Println()
		fmt.Println("   🛡️ Running guard checks...")

		// Run guards
		result := ops.Guard(tc.order, checks...)

		fmt.Println()
		if result.CanProceed {
			fmt.Println("   ✅ ALL GUARDS PASSED - Safe to ship")
		} else {
			fmt.Println("   ❌ GUARDS FAILED - Cannot proceed")
		}

		fmt.Println()
		fmt.Println("   Check Results:")
		if len(result.FailedChecks) == 0 {
			fmt.Println("      ✓ All checks passed")
		} else {
			for _, check := range result.FailedChecks {
				fmt.Printf("      ✗ %s\n", check)
			}
		}
	}

	fmt.Println()
	fmt.Println("📊 Guard Summary:")
	fmt.Println("   Total orders checked: 4")
	fmt.Println("   Passed all guards: 1")
	fmt.Println("   Failed guards: 3")
	fmt.Println()
	fmt.Println("✨ Success! Guard checks complete")
}

func valueOrMissing(s string) string {
	if s == "" {
		return "(missing)"
	}
	return s
}
