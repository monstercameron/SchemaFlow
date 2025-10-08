package main

import (
	"fmt"
	"log"

	schemaflow "github.com/monstercameron/SchemaFlow"
	"github.com/monstercameron/SchemaFlow/ops"
)

// Product represents a product to compare
type Product struct {
	Name     string
	Price    float64
	Features []string
	Specs    map[string]string
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		log.Fatalf("Failed to initialize SchemaFlow: %v", err)
	}

	// Two products to compare
	productA := Product{
		Name:  "UltraPhone Pro Max",
		Price: 1299.99,
		Features: []string{
			"6.7-inch OLED display",
			"108MP triple camera",
			"5000mAh battery",
			"5G connectivity",
			"IP68 water resistance",
		},
		Specs: map[string]string{
			"Processor": "Snapdragon 8 Gen 3",
			"RAM":       "12GB",
			"Storage":   "256GB",
			"OS":        "Android 14",
		},
	}

	productB := Product{
		Name:  "SmartPhone Elite",
		Price: 1099.99,
		Features: []string{
			"6.5-inch AMOLED display",
			"64MP dual camera",
			"4500mAh battery",
			"5G connectivity",
			"Premium aluminum build",
		},
		Specs: map[string]string{
			"Processor": "A17 Bionic",
			"RAM":       "8GB",
			"Storage":   "128GB",
			"OS":        "iOS 17",
		},
	}

	fmt.Println("🔍 Compare Example - Product Comparison")
	fmt.Println("=" + string(make([]byte, 60)))

	fmt.Println("\n📱 Product A:")
	fmt.Printf("   Name:  %s ($%.2f)\n", productA.Name, productA.Price)
	fmt.Println("   Features:", productA.Features)

	fmt.Println("\n📱 Product B:")
	fmt.Printf("   Name:  %s ($%.2f)\n", productB.Name, productB.Price)
	fmt.Println("   Features:", productB.Features)

	// Compare the products
	compareOpts := ops.NewCompareOptions().
		WithComparisonAspects([]string{"camera", "battery", "display", "performance", "value"}).
		WithOutputFormat("structured")
	compareOpts.Depth = 7
	compareOpts.OpOptions.Intelligence = schemaflow.Smart

	comparison, err := ops.Compare(productA, productB, compareOpts)
	if err != nil {
		log.Fatalf("Comparison failed: %v", err)
	}

	// Display comparison
	fmt.Println("\n✅ Detailed Comparison:")
	fmt.Println("---")
	fmt.Println(comparison)
	fmt.Println("---")

	fmt.Println("\n📊 Quick Summary:")
	fmt.Println("   UltraPhone Pro Max:")
	fmt.Println("   ✓ Better camera (108MP vs 64MP)")
	fmt.Println("   ✓ Larger battery (5000mAh vs 4500mAh)")
	fmt.Println("   ✓ More RAM (12GB vs 8GB)")
	fmt.Println("   ✓ More storage (256GB vs 128GB)")
	fmt.Println()
	fmt.Println("   SmartPhone Elite:")
	fmt.Println("   ✓ Lower price ($1099 vs $1299)")
	fmt.Println("   ✓ Premium iOS ecosystem")
	fmt.Println("   ✓ Optimized hardware/software integration")

	fmt.Println("\n🎯 Recommendation:")
	fmt.Println("   Choose UltraPhone Pro Max for:")
	fmt.Println("   • Photography enthusiasts")
	fmt.Println("   • Heavy multitasking")
	fmt.Println("   • Android ecosystem preference")
	fmt.Println()
	fmt.Println("   Choose SmartPhone Elite for:")
	fmt.Println("   • Apple ecosystem users")
	fmt.Println("   • Better value for money")
	fmt.Println("   • Premium build quality")

	fmt.Println("\n✨ Success! Comprehensive product comparison complete")
}
