package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// JobCandidate represents a job applicant
type JobCandidate struct {
	Name             string   `json:"name"`
	YearsExp         int      `json:"years_experience"`
	Education        string   `json:"education"`
	Skills           []string `json:"skills"`
	Location         string   `json:"location"`
	SalaryExpect     int      `json:"salary_expectation"`
	NoticePeriod     int      `json:"notice_period_days"`
	RemotePreference string   `json:"remote_preference"`
}

// LoanApplication represents a loan request
type LoanApplication struct {
	ApplicantName   string  `json:"applicant_name"`
	Amount          float64 `json:"amount"`
	Purpose         string  `json:"purpose"`
	Income          float64 `json:"annual_income"`
	CreditScore     int     `json:"credit_score"`
	DebtToIncome    float64 `json:"debt_to_income"`
	EmploymentYears int     `json:"employment_years"`
	Collateral      string  `json:"collateral"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Arbitrate Example ===")

	// Example 1: Job candidate selection
	fmt.Println("\n--- Example 1: Candidate Selection ---")

	candidates := []JobCandidate{
		{
			Name:             "Alice Chen",
			YearsExp:         5,
			Education:        "Master's CS",
			Skills:           []string{"Go", "Python", "Kubernetes", "AWS"},
			Location:         "San Francisco",
			SalaryExpect:     180000,
			NoticePeriod:     14,
			RemotePreference: "hybrid",
		},
		{
			Name:             "Bob Smith",
			YearsExp:         8,
			Education:        "Bachelor's CS",
			Skills:           []string{"Java", "Spring", "Docker", "GCP"},
			Location:         "Austin",
			SalaryExpect:     165000,
			NoticePeriod:     30,
			RemotePreference: "remote",
		},
		{
			Name:             "Carol Johnson",
			YearsExp:         3,
			Education:        "PhD CS",
			Skills:           []string{"Go", "Rust", "ML", "AWS", "Kubernetes"},
			Location:         "New York",
			SalaryExpect:     200000,
			NoticePeriod:     60,
			RemotePreference: "office",
		},
	}

	rules := []string{
		"Must have at least 3 years of experience",
		"Must know Go or Python",
		"Prefer candidates with Kubernetes experience",
		"Salary expectation must be under $190,000",
		"Notice period should be under 30 days",
		"Remote or hybrid preference is a plus",
	}

	result, err := schemaflow.Arbitrate[JobCandidate](candidates, schemaflow.ArbitrateOptions{
		Rules:           rules,
		RequireAllRules: false, // Best effort, not all rules must pass
		Steering:        "We need someone who can start quickly and knows our stack (Go, K8s)",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Arbitration failed", "error", err)
		return
	}

	fmt.Printf("Winner: %s\n", result.Winner.Name)
	fmt.Printf("Winner Score: %.0f%%\n\n", result.Scores[result.WinnerIndex]*100)

	fmt.Println("Evaluation Summary:")
	for _, eval := range result.Evaluations {
		candidate := candidates[eval.Index]
		fmt.Printf("\n  %s (Score: %.0f%%)\n", candidate.Name, eval.TotalScore*100)

		passed := 0
		failed := 0
		for _, r := range eval.RuleResults {
			if r.Passed {
				passed++
			} else {
				failed++
			}
		}
		fmt.Printf("  Rules: %d passed, %d failed\n", passed, failed)

		fmt.Println("  Rule Details:")
		for _, rule := range eval.RuleResults {
			status := "✓"
			if !rule.Passed {
				status = "✗"
			}
			fmt.Printf("    %s %s\n", status, rule.Rule)
			if rule.Reasoning != "" {
				fmt.Printf("      → %s\n", rule.Reasoning)
			}
		}
	}

	fmt.Printf("\nReasoning: %s\n", result.Reasoning)

	// Example 2: Loan approval with strict rules
	fmt.Println("\n--- Example 2: Loan Approval (Strict) ---")

	loanApplications := []LoanApplication{
		{
			ApplicantName:   "John Doe",
			Amount:          50000,
			Purpose:         "Home Improvement",
			Income:          95000,
			CreditScore:     720,
			DebtToIncome:    0.28,
			EmploymentYears: 5,
			Collateral:      "Home equity",
		},
		{
			ApplicantName:   "Jane Wilson",
			Amount:          75000,
			Purpose:         "Business Expansion",
			Income:          120000,
			CreditScore:     680,
			DebtToIncome:    0.42,
			EmploymentYears: 3,
			Collateral:      "None",
		},
		{
			ApplicantName:   "Mike Brown",
			Amount:          30000,
			Purpose:         "Debt Consolidation",
			Income:          65000,
			CreditScore:     750,
			DebtToIncome:    0.35,
			EmploymentYears: 8,
			Collateral:      "Vehicle",
		},
	}

	loanRules := []string{
		"Credit score must be at least 650",
		"Debt-to-income ratio must be under 0.40",
		"Loan amount must not exceed 60% of annual income",
		"Must have at least 2 years of employment",
		"Collateral required for loans over $50,000",
	}

	loanResult, err := schemaflow.Arbitrate[LoanApplication](loanApplications, schemaflow.ArbitrateOptions{
		Rules:           loanRules,
		RequireAllRules: true, // All rules must pass
	})

	if err != nil {
		schemaflow.GetLogger().Error("Loan arbitration failed", "error", err)
		return
	}

	fmt.Println("Loan Evaluation Results:")
	for _, eval := range loanResult.Evaluations {
		app := loanApplications[eval.Index]
		status := "APPROVED"
		if eval.Disqualified {
			status = "DENIED"
		}
		fmt.Printf("\n  %s - %s (Score: %.0f%%)\n",
			app.ApplicantName, status, eval.TotalScore*100)
		fmt.Printf("  Loan: $%.0f for %s\n", app.Amount, app.Purpose)

		for _, rule := range eval.RuleResults {
			if !rule.Passed {
				fmt.Printf("    ✗ FAILED: %s\n", rule.Rule)
				fmt.Printf("      Reason: %s\n", rule.Reasoning)
			}
		}
	}

	if loanResult.Winner.ApplicantName != "" {
		fmt.Printf("\nBest Candidate: %s\n", loanResult.Winner.ApplicantName)
	} else {
		fmt.Println("\nNo candidates passed all requirements")
	}

	fmt.Println("\n=== Arbitrate Example Complete ===")
}
