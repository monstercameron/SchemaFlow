package main

import (
	"fmt"
	"log"
	"os"

	"github.com/monstercameron/SchemaFlow/internal/ops"
)

func main() {
	// Ensure environment is configured
	if os.Getenv("SCHEMAFLOW_API_KEY") == "" {
		log.Fatal("SCHEMAFLOW_API_KEY environment variable not set")
	}

	fmt.Println("=== Predict Example ===")

	// Example 1: Predict sales for next quarter
	fmt.Println("--- Example 1: Sales Forecast ---")
	type SalesForecast struct {
		Revenue    float64 `json:"revenue"`
		GrowthRate float64 `json:"growth_rate"`
		Units      int     `json:"units"`
	}

	historicalSales := []map[string]any{
		{"quarter": "Q1 2023", "revenue": 1000000, "units": 5000},
		{"quarter": "Q2 2023", "revenue": 1150000, "units": 5500},
		{"quarter": "Q3 2023", "revenue": 1100000, "units": 5200},
		{"quarter": "Q4 2023", "revenue": 1400000, "units": 7000},
		{"quarter": "Q1 2024", "revenue": 1250000, "units": 6000},
		{"quarter": "Q2 2024", "revenue": 1450000, "units": 7200},
	}

	opts := ops.NewPredictOptions().
		WithHorizon("Q3 2024").
		WithIncludeConfidenceInterval(true).
		WithConfidenceLevel(0.8).
		WithIncludeReasoning(true)

	result, err := ops.Predict[SalesForecast](historicalSales, opts)
	if err != nil {
		log.Fatalf("Prediction failed: %v", err)
	}

	fmt.Println("Q3 2024 Sales Forecast:")
	fmt.Printf("  Predicted Revenue: $%.2f\n", result.Prediction.Revenue)
	fmt.Printf("  Predicted Units: %d\n", result.Prediction.Units)
	fmt.Printf("  Growth Rate: %.1f%%\n", result.Prediction.GrowthRate*100)
	fmt.Printf("  Confidence: %.0f%%\n\n", result.Confidence*100)

	if result.Interval != nil {
		fmt.Printf("  %.0f%% Confidence Interval:\n", result.Interval.ConfidenceLevel*100)
		fmt.Printf("    Revenue: $%.2f - $%.2f\n", result.Interval.Lower, result.Interval.Upper)
	}

	if result.Reasoning != "" {
		fmt.Printf("\n  Reasoning: %s\n", result.Reasoning)
	}
	fmt.Println()

	// Example 2: Predict with scenarios
	fmt.Println("--- Example 2: Scenario-Based Prediction ---")
	type MarketPrediction struct {
		Price  float64 `json:"price"`
		Volume int     `json:"volume"`
		Trend  string  `json:"trend"`
	}

	marketData := []map[string]any{
		{"date": "2024-01", "price": 100, "volume": 10000},
		{"date": "2024-02", "price": 105, "volume": 12000},
		{"date": "2024-03", "price": 102, "volume": 11000},
		{"date": "2024-04", "price": 110, "volume": 15000},
		{"date": "2024-05", "price": 115, "volume": 14000},
	}

	scenarioOpts := ops.NewPredictOptions().
		WithHorizon("next 3 months").
		WithIncludeScenarios(true).
		WithNumScenarios(3).
		WithFactors([]string{"market_sentiment", "seasonality", "competition"})

	scenarioResult, err := ops.Predict[MarketPrediction](marketData, scenarioOpts)
	if err != nil {
		log.Fatalf("Scenario prediction failed: %v", err)
	}

	fmt.Println("Market Prediction (Base Case):")
	fmt.Printf("  Price: $%.2f\n", scenarioResult.Prediction.Price)
	fmt.Printf("  Volume: %d\n", scenarioResult.Prediction.Volume)
	fmt.Printf("  Trend: %s\n\n", scenarioResult.Prediction.Trend)

	fmt.Println("Alternative Scenarios:")
	for _, scenario := range scenarioResult.Scenarios {
		fmt.Printf("\n  %s (%.0f%% probability)\n", scenario.Name, scenario.Probability*100)
		fmt.Printf("  Description: %s\n", scenario.Description)
		if len(scenario.Conditions) > 0 {
			fmt.Printf("  Conditions: %v\n", scenario.Conditions)
		}
	}
	fmt.Println()

	// Example 3: Predict with assumptions
	fmt.Println("--- Example 3: Prediction with Assumptions ---")
	type ResourceForecast struct {
		Headcount int     `json:"headcount"`
		Budget    float64 `json:"budget"`
		Projects  int     `json:"projects"`
	}

	teamHistory := []map[string]any{
		{"year": 2021, "headcount": 10, "budget": 500000, "projects": 5},
		{"year": 2022, "headcount": 15, "budget": 750000, "projects": 8},
		{"year": 2023, "headcount": 22, "budget": 1100000, "projects": 12},
	}

	assumptionOpts := ops.NewPredictOptions().
		WithHorizon("2025").
		WithAssumptions([]string{
			"20% annual growth target",
			"No major market disruptions",
			"Continued remote work policy",
		}).
		WithIncludeReasoning(true)

	resourceResult, err := ops.Predict[ResourceForecast](teamHistory, assumptionOpts)
	if err != nil {
		log.Fatalf("Resource prediction failed: %v", err)
	}

	fmt.Println("2025 Resource Forecast:")
	fmt.Printf("  Projected Headcount: %d\n", resourceResult.Prediction.Headcount)
	fmt.Printf("  Projected Budget: $%.0f\n", resourceResult.Prediction.Budget)
	fmt.Printf("  Expected Projects: %d\n", resourceResult.Prediction.Projects)

	if len(resourceResult.Assumptions) > 0 {
		fmt.Println("\n  Based on assumptions:")
		for _, a := range resourceResult.Assumptions {
			fmt.Printf("    - %s\n", a)
		}
	}

	if len(resourceResult.Risks) > 0 {
		fmt.Println("\n  Identified Risks:")
		for _, r := range resourceResult.Risks {
			fmt.Printf("    ⚠ %s\n", r)
		}
	}
	fmt.Println()

	// Example 4: Trend-based prediction
	fmt.Println("--- Example 4: Trend Analysis ---")
	type TrendPrediction struct {
		Value     float64 `json:"value"`
		Direction string  `json:"direction"`
		Velocity  float64 `json:"velocity"`
	}

	trendData := []float64{10, 12, 15, 14, 18, 22, 21, 25, 28}

	trendOpts := ops.NewPredictOptions().
		WithHorizon("next 3 periods").
		WithMethod("trend").
		WithIncludeReasoning(true)

	trendResult, err := ops.Predict[TrendPrediction](trendData, trendOpts)
	if err != nil {
		log.Fatalf("Trend prediction failed: %v", err)
	}

	fmt.Println("Trend Analysis:")
	fmt.Printf("  Next Value: %.2f\n", trendResult.Prediction.Value)
	fmt.Printf("  Direction: %s\n", trendResult.Prediction.Direction)
	fmt.Printf("  Velocity: %.2f per period\n", trendResult.Prediction.Velocity)

	if len(trendResult.Factors) > 0 {
		fmt.Println("\n  Key Factors:")
		for _, f := range trendResult.Factors {
			fmt.Printf("    %s [%s]: weight %.2f\n", f.Name, f.Impact, f.Weight)
		}
	}

	fmt.Println("\n=== Predict Example Complete ===")
}
