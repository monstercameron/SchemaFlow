package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// SalesReport represents a quarterly sales report
type SalesReport struct {
	Quarter     string        `json:"quarter"`
	Year        int           `json:"year"`
	TotalSales  float64       `json:"total_sales"`
	UnitsSold   int           `json:"units_sold"`
	TopProducts []ProductSale `json:"top_products"`
	Regions     []RegionSale  `json:"regions"`
}

// ProductSale represents product sales data
type ProductSale struct {
	Name      string  `json:"name"`
	Revenue   float64 `json:"revenue"`
	UnitsSold int     `json:"units_sold"`
	Category  string  `json:"category"`
}

// RegionSale represents sales by region
type RegionSale struct {
	Region  string  `json:"region"`
	Revenue float64 `json:"revenue"`
	Growth  float64 `json:"growth"` // percentage
}

// KeyFindings represents typed findings from the report
type KeyFindings struct {
	TopPerformer      string `json:"top_performer"`
	BestRegion        string `json:"best_region"`
	GrowthTrend       string `json:"growth_trend"`
	OverallHealth     string `json:"overall_health"`
	RecommendedAction string `json:"recommended_action"`
}

// RiskAssessment represents typed risk analysis
type RiskAssessment struct {
	Risks []struct {
		Name       string  `json:"name"`
		Severity   string  `json:"severity"`
		Likelihood float64 `json:"likelihood"`
		Impact     string  `json:"impact"`
		Mitigation string  `json:"mitigation"`
	} `json:"risks"`
	OverallRiskLevel string `json:"overall_risk_level"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("❓ Question Example - Ask Questions About Data with Typed Answers")
	fmt.Println("=" + string(make([]byte, 65)))

	// Sample sales report
	report := SalesReport{
		Quarter:    "Q3",
		Year:       2024,
		TotalSales: 4850000,
		UnitsSold:  125000,
		TopProducts: []ProductSale{
			{Name: "Premium Widget", Revenue: 1200000, UnitsSold: 15000, Category: "Hardware"},
			{Name: "Smart Gadget", Revenue: 950000, UnitsSold: 12000, Category: "Electronics"},
			{Name: "Pro Service Plan", Revenue: 800000, UnitsSold: 8000, Category: "Services"},
			{Name: "Basic Widget", Revenue: 650000, UnitsSold: 45000, Category: "Hardware"},
			{Name: "Accessory Pack", Revenue: 450000, UnitsSold: 35000, Category: "Accessories"},
		},
		Regions: []RegionSale{
			{Region: "North America", Revenue: 2100000, Growth: 15.5},
			{Region: "Europe", Revenue: 1400000, Growth: 8.2},
			{Region: "Asia Pacific", Revenue: 950000, Growth: 22.3},
			{Region: "Latin America", Revenue: 400000, Growth: -3.1},
		},
	}

	fmt.Printf("\n📊 Sales Report: %s %d\n", report.Quarter, report.Year)
	fmt.Printf("   Total Sales: $%.2fM\n", report.TotalSales/1000000)
	fmt.Printf("   Units Sold: %d\n", report.UnitsSold)
	fmt.Println()

	// Example 1: Simple string answer
	fmt.Println("--- Example 1: Simple String Answer ---")
	simpleOpts := schemaflow.NewQuestionOptions("What is the best performing product by revenue?")
	simpleResult, err := schemaflow.Question[SalesReport, string](report, simpleOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Question failed", "error", err)
	} else {
		fmt.Printf("Q: What is the best performing product by revenue?\n")
		fmt.Printf("A: %s\n", simpleResult.Answer)
		fmt.Printf("   Confidence: %.0f%%\n", simpleResult.Confidence*100)
		if simpleResult.Reasoning != "" {
			fmt.Printf("   Reasoning: %s\n", simpleResult.Reasoning)
		}
	}
	fmt.Println()

	// Example 2: Boolean answer
	fmt.Println("--- Example 2: Boolean Answer ---")
	boolOpts := schemaflow.NewQuestionOptions("Are all regions showing positive growth?")
	boolResult, err := schemaflow.Question[SalesReport, bool](report, boolOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Question failed", "error", err)
	} else {
		fmt.Printf("Q: Are all regions showing positive growth?\n")
		if boolResult.Answer {
			fmt.Println("A: ✅ Yes")
		} else {
			fmt.Println("A: ❌ No")
		}
		fmt.Printf("   Confidence: %.0f%%\n", boolResult.Confidence*100)
		if len(boolResult.Evidence) > 0 {
			fmt.Println("   Evidence:")
			for _, e := range boolResult.Evidence {
				fmt.Printf("   - %s\n", e)
			}
		}
	}
	fmt.Println()

	// Example 3: Typed struct answer - Key Findings
	fmt.Println("--- Example 3: Typed Struct Answer (Key Findings) ---")
	findingsOpts := schemaflow.NewQuestionOptions("Summarize the key findings from this sales report including top performer, best region, growth trend, overall health assessment, and a recommended action")
	findingsResult, err := schemaflow.Question[SalesReport, KeyFindings](report, findingsOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Question failed", "error", err)
	} else {
		fmt.Println("Q: What are the key findings?")
		fmt.Println("A:")
		fmt.Printf("   🏆 Top Performer: %s\n", findingsResult.Answer.TopPerformer)
		fmt.Printf("   🌍 Best Region: %s\n", findingsResult.Answer.BestRegion)
		fmt.Printf("   📈 Growth Trend: %s\n", findingsResult.Answer.GrowthTrend)
		fmt.Printf("   ❤️ Overall Health: %s\n", findingsResult.Answer.OverallHealth)
		fmt.Printf("   💡 Recommended Action: %s\n", findingsResult.Answer.RecommendedAction)
		fmt.Printf("   Confidence: %.0f%%\n", findingsResult.Confidence*100)
	}
	fmt.Println()

	// Example 4: Typed struct answer - Risk Assessment
	fmt.Println("--- Example 4: Typed Struct Answer (Risk Assessment) ---")
	riskOpts := schemaflow.NewQuestionOptions("What are the business risks based on this sales data? Identify potential risks with their severity, likelihood, impact, and mitigation strategies.")
	riskResult, err := schemaflow.Question[SalesReport, RiskAssessment](report, riskOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Question failed", "error", err)
	} else {
		fmt.Println("Q: What are the business risks?")
		fmt.Printf("A: Overall Risk Level: %s\n", riskResult.Answer.OverallRiskLevel)
		fmt.Println("   Identified Risks:")
		for i, risk := range riskResult.Answer.Risks {
			fmt.Printf("\n   %d. %s\n", i+1, risk.Name)
			fmt.Printf("      Severity: %s, Likelihood: %.0f%%\n", risk.Severity, risk.Likelihood*100)
			fmt.Printf("      Impact: %s\n", risk.Impact)
			fmt.Printf("      Mitigation: %s\n", risk.Mitigation)
		}
		fmt.Printf("\n   Confidence: %.0f%%\n", riskResult.Confidence*100)
	}
	fmt.Println()

	// Example 5: Numeric answer
	fmt.Println("--- Example 5: Numeric Answer ---")
	numOpts := schemaflow.NewQuestionOptions("What is the average revenue per unit across all products?")
	numResult, err := schemaflow.Question[SalesReport, float64](report, numOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Question failed", "error", err)
	} else {
		fmt.Printf("Q: What is the average revenue per unit?\n")
		fmt.Printf("A: $%.2f per unit\n", numResult.Answer)
		fmt.Printf("   Confidence: %.0f%%\n", numResult.Confidence*100)
		if numResult.Reasoning != "" {
			fmt.Printf("   Reasoning: %s\n", numResult.Reasoning)
		}
	}
	fmt.Println()

	// Example 6: Using legacy interface for quick string answers
	fmt.Println("--- Example 6: Legacy String Interface ---")
	legacyAnswer, err := schemaflow.QuestionLegacy(report, "Which region should be prioritized for investment?")
	if err != nil {
		schemaflow.GetLogger().Error("Question failed", "error", err)
	} else {
		fmt.Printf("Q: Which region should be prioritized for investment?\n")
		fmt.Printf("A: %s\n", legacyAnswer)
	}

	fmt.Println()
	fmt.Println("✨ Success! Question examples complete")
	fmt.Println()
	fmt.Println("📝 Key Features Demonstrated:")
	fmt.Println("   - String answers with confidence and reasoning")
	fmt.Println("   - Boolean answers with evidence")
	fmt.Println("   - Typed struct answers for structured extraction")
	fmt.Println("   - Numeric answers with explanation")
	fmt.Println("   - Legacy interface for quick string questions")
}
