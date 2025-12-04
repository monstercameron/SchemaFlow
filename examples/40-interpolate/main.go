package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// SalesMetric represents a daily sales record
type SalesMetric struct {
	Date      string  `json:"date"`
	Revenue   float64 `json:"revenue"`
	Units     int     `json:"units"`
	Customers int     `json:"customers"`
	AvgOrder  float64 `json:"avg_order"`
}

// EmployeeRecord represents attendance/performance data
type EmployeeRecord struct {
	Week          int     `json:"week"`
	HoursWorked   float64 `json:"hours_worked"`
	TasksComplete int     `json:"tasks_complete"`
	Rating        float64 `json:"rating"`
	Present       bool    `json:"present"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Interpolate Example ===")

	// Example 1: Fill gaps in sales time series
	fmt.Println("\n--- Example 1: Sales Data with Gaps ---")

	// Note: Some days have missing data (zero values)
	salesData := []SalesMetric{
		{Date: "2024-01-01", Revenue: 5200.50, Units: 42, Customers: 38, AvgOrder: 136.86},
		{Date: "2024-01-02", Revenue: 4800.00, Units: 38, Customers: 35, AvgOrder: 137.14},
		{Date: "2024-01-03", Revenue: 0, Units: 0, Customers: 0, AvgOrder: 0}, // Missing!
		{Date: "2024-01-04", Revenue: 6100.00, Units: 48, Customers: 42, AvgOrder: 145.24},
		{Date: "2024-01-05", Revenue: 0, Units: 0, Customers: 0, AvgOrder: 0}, // Missing!
		{Date: "2024-01-06", Revenue: 0, Units: 0, Customers: 0, AvgOrder: 0}, // Missing!
		{Date: "2024-01-07", Revenue: 7500.00, Units: 55, Customers: 48, AvgOrder: 156.25},
	}

	fmt.Println("Original Data (0 = missing):")
	for _, s := range salesData {
		status := ""
		if s.Revenue == 0 {
			status = " [MISSING]"
		}
		fmt.Printf("  %s: $%.2f, %d units%s\n", s.Date, s.Revenue, s.Units, status)
	}

	result, err := schemaflow.Interpolate[SalesMetric](salesData, schemaflow.InterpolateOptions{
		Method:   "contextual",
		Steering: "Consider day-of-week patterns. Weekends might have different patterns. Zero values indicate missing data.",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Interpolation failed", "error", err)
		return
	}

	fmt.Println("\nInterpolated Data:")
	for _, s := range result.Complete {
		status := ""
		for _, filled := range result.Filled {
			if filled.Index == indexOf(result.Complete, s) {
				status = fmt.Sprintf(" [FILLED, %.0f%% conf]", filled.Confidence*100)
				break
			}
		}
		fmt.Printf("  %s: $%.2f, %d units%s\n", s.Date, s.Revenue, s.Units, status)
	}

	fmt.Printf("\nFilled Items: %d\n", len(result.Filled))
	for _, f := range result.Filled {
		fmt.Printf("  Index %d: %s\n", f.Index, f.Method)
		if f.Reasoning != "" {
			fmt.Printf("    Reasoning: %s\n", f.Reasoning)
		}
	}

	fmt.Printf("\nGaps Detected: %d\n", result.GapCount)
	fmt.Printf("Average Confidence: %.0f%%\n", result.AverageConfidence*100)

	// Example 2: Employee performance with missing weeks
	fmt.Println("\n--- Example 2: Employee Performance Data ---")

	performanceData := []EmployeeRecord{
		{Week: 1, HoursWorked: 40.0, TasksComplete: 12, Rating: 4.2, Present: true},
		{Week: 2, HoursWorked: 38.5, TasksComplete: 11, Rating: 4.0, Present: true},
		{Week: 3, HoursWorked: 0, TasksComplete: 0, Rating: 0, Present: false}, // Vacation?
		{Week: 4, HoursWorked: 0, TasksComplete: 0, Rating: 0, Present: false}, // Still out?
		{Week: 5, HoursWorked: 42.0, TasksComplete: 14, Rating: 4.5, Present: true},
		{Week: 6, HoursWorked: 40.0, TasksComplete: 13, Rating: 4.3, Present: true},
		{Week: 7, HoursWorked: 0, TasksComplete: 0, Rating: 0, Present: false}, // Sick?
		{Week: 8, HoursWorked: 41.5, TasksComplete: 15, Rating: 4.6, Present: true},
	}

	fmt.Println("Original Performance Data:")
	for _, p := range performanceData {
		if p.Present {
			fmt.Printf("  Week %d: %.1f hrs, %d tasks, rating %.1f\n", p.Week, p.HoursWorked, p.TasksComplete, p.Rating)
		} else {
			fmt.Printf("  Week %d: [ABSENT/MISSING]\n", p.Week)
		}
	}

	perfResult, err := schemaflow.Interpolate[EmployeeRecord](performanceData, schemaflow.InterpolateOptions{
		Method:   "pattern",
		Steering: "Weeks 3-4 were likely planned vacation. Week 7 might be sick day. Interpolate what their performance WOULD have been if present.",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Performance interpolation failed", "error", err)
		return
	}

	fmt.Println("\nInterpolated Performance:")
	for _, p := range perfResult.Complete {
		marker := ""
		for _, f := range perfResult.Filled {
			if f.Index == p.Week-1 {
				marker = fmt.Sprintf(" [INTERPOLATED: %s]", f.Method)
			}
		}
		fmt.Printf("  Week %d: %.1f hrs, %d tasks, rating %.1f%s\n",
			p.Week, p.HoursWorked, p.TasksComplete, p.Rating, marker)
	}

	fmt.Println("\n=== Interpolate Example Complete ===")
}

// Helper function to find index
func indexOf[T any](slice []T, item T) int {
	for i := range slice {
		if fmt.Sprintf("%v", slice[i]) == fmt.Sprintf("%v", item) {
			return i
		}
	}
	return -1
}
