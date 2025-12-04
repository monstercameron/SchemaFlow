package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// ProjectPlan represents the negotiated output
type ProjectPlan struct {
	Duration      int      `json:"duration"`       // Weeks
	Budget        int      `json:"budget"`         // Dollars
	TeamSize      int      `json:"team_size"`      // People
	Features      []string `json:"features"`       // Features to include
	Quality       string   `json:"quality"`        // "high", "medium", "basic"
	Deadline      string   `json:"deadline"`       // Target date
	DeliverySplit []string `json:"delivery_split"` // Phased delivery
}

// ProjectConstraints represents competing requirements
type ProjectConstraints struct {
	MaxBudget          int      `json:"max_budget"`
	MinFeatures        []string `json:"min_features"`
	DesiredFeatures    []string `json:"desired_features"`
	HardDeadline       string   `json:"hard_deadline"`
	PreferredDeadline  string   `json:"preferred_deadline"`
	QualityRequirement string   `json:"quality_requirement"`
	TeamAvailability   int      `json:"team_availability"`
	RiskTolerance      string   `json:"risk_tolerance"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Negotiate Example ===")

	// Example 1: Project planning with competing constraints
	fmt.Println("\n--- Example 1: Project Planning ---")

	constraints := ProjectConstraints{
		MaxBudget:          100000,
		MinFeatures:        []string{"user_auth", "dashboard", "api"},
		DesiredFeatures:    []string{"analytics", "mobile_app", "integrations", "ai_features"},
		HardDeadline:       "2024-06-01",
		PreferredDeadline:  "2024-04-01",
		QualityRequirement: "high",
		TeamAvailability:   4, // 4 developers available
		RiskTolerance:      "low",
	}

	result, err := schemaflow.Negotiate[ProjectPlan](constraints, schemaflow.NegotiateOptions{
		Strategy: "balanced",
		Steering: "Prioritize meeting the hard deadline over including all desired features",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Negotiation failed", "error", err)
		return
	}

	fmt.Printf("Negotiated Plan:\n")
	fmt.Printf("  Duration: %d weeks\n", result.Solution.Duration)
	fmt.Printf("  Budget: $%d\n", result.Solution.Budget)
	fmt.Printf("  Team Size: %d\n", result.Solution.TeamSize)
	fmt.Printf("  Features: %v\n", result.Solution.Features)
	fmt.Printf("  Quality: %s\n", result.Solution.Quality)
	fmt.Printf("  Deadline: %s\n", result.Solution.Deadline)

	if len(result.Solution.DeliverySplit) > 0 {
		fmt.Printf("  Phased Delivery:\n")
		for _, phase := range result.Solution.DeliverySplit {
			fmt.Printf("    - %s\n", phase)
		}
	}

	fmt.Printf("\nTradeoffs Made:\n")
	for _, tradeoff := range result.Tradeoffs {
		fmt.Printf("  - Sacrificed: %s, Gained: %s\n", tradeoff.Sacrificed, tradeoff.Gained)
		fmt.Printf("    Impact: %s\n", tradeoff.Impact)
		if tradeoff.Reasoning != "" {
			fmt.Printf("    Reason: %s\n", tradeoff.Reasoning)
		}
	}

	fmt.Printf("\nConstraint Satisfaction:\n")
	for constraint, score := range result.Satisfaction {
		fmt.Printf("  %s: %.0f%%\n", constraint, score*100)
	}

	fmt.Printf("\nOverall Confidence: %.0f%%\n", result.Confidence*100)

	// Example 2: Salary negotiation
	fmt.Println("\n--- Example 2: Salary Negotiation ---")

	type SalaryOffer struct {
		BaseSalary   int    `json:"base_salary"`
		Bonus        int    `json:"bonus"`
		Equity       string `json:"equity"`
		RemoteDays   int    `json:"remote_days"`
		VacationDays int    `json:"vacation_days"`
		StartDate    string `json:"start_date"`
		SigningBonus int    `json:"signing_bonus"`
	}

	salaryConstraints := map[string]any{
		"candidate_minimum_salary": 150000,
		"candidate_wants_remote":   "3-5 days per week",
		"candidate_wants_equity":   "at least 0.1%",
		"candidate_ideal_start":    "2024-02-01",
		"company_max_budget":       170000,
		"company_equity_pool":      "0.05-0.15%",
		"company_remote_policy":    "2-3 days per week",
		"company_needs_start_by":   "2024-03-01",
		"industry_average":         145000,
	}

	salaryResult, err := schemaflow.Negotiate[SalaryOffer](salaryConstraints, schemaflow.NegotiateOptions{
		Strategy: "balanced",
		Priorities: map[string]float64{
			"base_salary":   0.4,
			"remote_days":   0.25,
			"equity":        0.2,
			"vacation_days": 0.15,
		},
	})

	if err != nil {
		schemaflow.GetLogger().Error("Salary negotiation failed", "error", err)
		return
	}

	fmt.Printf("Negotiated Offer:\n")
	fmt.Printf("  Base Salary: $%d\n", salaryResult.Solution.BaseSalary)
	fmt.Printf("  Bonus: $%d\n", salaryResult.Solution.Bonus)
	fmt.Printf("  Equity: %s\n", salaryResult.Solution.Equity)
	fmt.Printf("  Remote Days: %d/week\n", salaryResult.Solution.RemoteDays)
	fmt.Printf("  Vacation Days: %d\n", salaryResult.Solution.VacationDays)
	fmt.Printf("  Start Date: %s\n", salaryResult.Solution.StartDate)
	if salaryResult.Solution.SigningBonus > 0 {
		fmt.Printf("  Signing Bonus: $%d\n", salaryResult.Solution.SigningBonus)
	}

	fmt.Printf("\nKey Tradeoffs:\n")
	for _, tradeoff := range salaryResult.Tradeoffs {
		fmt.Printf("  - Sacrificed: %s → Gained: %s\n", tradeoff.Sacrificed, tradeoff.Gained)
	}

	fmt.Println("\n=== Negotiate Example Complete ===")
}
