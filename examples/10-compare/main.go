package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
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
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
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

	// Compare the products using the new typed Compare
	compareOpts := schemaflow.NewCompareOptions().
		WithComparisonAspects([]string{"camera", "battery", "display", "performance", "value"}).
		WithFocusOn("both")
	compareOpts.Depth = 7
	compareOpts.OpOptions.Intelligence = schemaflow.Smart

	result, err := schemaflow.Compare[Product](productA, productB, compareOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Comparison failed", "error", err)
		return
	}

	// Display structured comparison results
	fmt.Println("\n✅ Comparison Results:")
	fmt.Println("---")
	fmt.Printf("📊 Overall Similarity: %.0f%%\n", result.SimilarityScore*100)
	fmt.Printf("📝 Verdict: %s\n", result.Verdict)

	// Show aspect scores if available
	if len(result.AspectScores) > 0 {
		fmt.Println("\n📈 Similarity by Aspect:")
		for aspect, score := range result.AspectScores {
			bar := ""
			filled := int(score * 10)
			for i := 0; i < 10; i++ {
				if i < filled {
					bar += "█"
				} else {
					bar += "░"
				}
			}
			fmt.Printf("   %s: %s %.0f%%\n", aspect, bar, score*100)
		}
	}

	// Show similarities
	if len(result.Similarities) > 0 {
		fmt.Println("\n✅ Similarities:")
		for _, sim := range result.Similarities {
			fmt.Printf("   • [%s] %s\n", sim.Aspect, sim.Description)
		}
	}

	// Show differences
	if len(result.Differences) > 0 {
		fmt.Println("\n❌ Differences:")
		for _, diff := range result.Differences {
			severity := ""
			switch diff.Severity {
			case "major":
				severity = "🔴"
			case "moderate":
				severity = "🟡"
			case "minor":
				severity = "🟢"
			default:
				severity = "⚪"
			}
			fmt.Printf("   %s [%s] %s\n", severity, diff.Aspect, diff.Description)
		}
	}

	fmt.Println("\n🎯 Recommendation based on comparison:")
	fmt.Printf("   Choose %s for:\n", productA.Name)
	fmt.Println("   • Photography enthusiasts")
	fmt.Println("   • Heavy multitasking")
	fmt.Println("   • Android ecosystem preference")
	fmt.Println()
	fmt.Printf("   Choose %s for:\n", productB.Name)
	fmt.Println("   • Apple ecosystem users")
	fmt.Println("   • Better value for money")
	fmt.Println("   • Premium build quality")

	fmt.Println("\n✨ Success! Typed product comparison with detailed analysis")
}
