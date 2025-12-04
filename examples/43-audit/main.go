package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// CustomerData represents customer information to audit
type CustomerData struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	SSN         string  `json:"ssn"`
	CreditCard  string  `json:"credit_card"`
	Password    string  `json:"password"`
	DateOfBirth string  `json:"date_of_birth"`
	Balance     float64 `json:"balance"`
	Status      string  `json:"status"`
	Notes       string  `json:"notes"`
}

// FinancialReport represents financial data to audit
type FinancialReport struct {
	Quarter      string  `json:"quarter"`
	Revenue      float64 `json:"revenue"`
	Expenses     float64 `json:"expenses"`
	NetIncome    float64 `json:"net_income"`
	Headcount    int     `json:"headcount"`
	AvgSalary    float64 `json:"avg_salary"`
	TotalPayroll float64 `json:"total_payroll"`
	GrowthRate   float64 `json:"growth_rate"`
	PrevQuarter  float64 `json:"prev_quarter_revenue"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Audit Example ===")

	// Example 1: Security and compliance audit
	fmt.Println("\n--- Example 1: Customer Data Audit ---")

	// Intentionally problematic customer data
	customer := CustomerData{
		ID:          "CUST-001",
		Name:        "John Doe",
		Email:       "john.doe@",                   // Invalid email
		Phone:       "555-1234",                    // Incomplete phone
		SSN:         "123-45-6789",                 // Plain text SSN!
		CreditCard:  "4532-1234-5678-9012",         // Plain text CC!
		Password:    "password123",                 // Weak password in plain text!
		DateOfBirth: "2030-01-15",                  // Future date!
		Balance:     -500.00,                       // Negative balance
		Status:      "ACTVE",                       // Typo
		Notes:       "Customer SSN is 123-45-6789", // SSN in notes!
	}

	result, err := schemaflow.Audit[CustomerData](customer, schemaflow.AuditOptions{
		Policies: []string{
			"PII data must not be stored in plain text",
			"Credit card numbers must be tokenized or encrypted",
			"Passwords must never be stored in plain text",
			"Email addresses must be valid format",
			"Phone numbers must be complete with area code",
			"Date of birth must be in the past",
			"Account balance should not be negative without reason",
		},
		Categories: []string{"security", "compliance", "quality"},
		Threshold:  0.3, // Report medium+ issues
		Deep:       true,
	})

	if err != nil {
		schemaflow.GetLogger().Error("Audit failed", "error", err)
		return
	}

	fmt.Printf("Audit Results for Customer %s:\n\n", customer.ID)

	fmt.Printf("Summary:\n")
	fmt.Printf("  Total Findings: %d\n", result.Summary.TotalFindings)
	fmt.Printf("  Critical Issues: %v\n", result.Summary.Critical)
	fmt.Printf("  Passes Audit: %v\n", result.Summary.PassesAudit)

	fmt.Printf("\nBy Severity:\n")
	for level, count := range result.Summary.BySeverity {
		fmt.Printf("  %s: %d\n", level, count)
	}

	fmt.Printf("\nBy Category:\n")
	for cat, count := range result.Summary.ByCategory {
		fmt.Printf("  %s: %d\n", cat, count)
	}

	fmt.Printf("\nDetailed Findings:\n")
	for _, finding := range result.Findings {
		severity := "INFO"
		switch {
		case finding.Severity >= 0.9:
			severity = "🔴 CRITICAL"
		case finding.Severity >= 0.7:
			severity = "🟠 HIGH"
		case finding.Severity >= 0.5:
			severity = "🟡 MEDIUM"
		case finding.Severity >= 0.3:
			severity = "🔵 LOW"
		}

		fmt.Printf("\n  %s [%s]\n", severity, finding.Category)
		fmt.Printf("  Field: %s\n", finding.Field)
		fmt.Printf("  Issue: %s\n", finding.Issue)
		if finding.Evidence != "" {
			fmt.Printf("  Evidence: %s\n", finding.Evidence)
		}
		if finding.Recommendation != "" {
			fmt.Printf("  Recommendation: %s\n", finding.Recommendation)
		}
	}

	// Example 2: Financial consistency audit
	fmt.Println("\n\n--- Example 2: Financial Report Audit ---")

	report := FinancialReport{
		Quarter:      "Q4 2023",
		Revenue:      10000000,
		Expenses:     8500000,
		NetIncome:    2000000, // Should be 1.5M (math error!)
		Headcount:    100,
		AvgSalary:    75000,
		TotalPayroll: 8000000, // 100 * 75000 = 7.5M (doesn't match!)
		GrowthRate:   0.25,    // 25%
		PrevQuarter:  8500000, // (10M - 8.5M) / 8.5M = 17.6% (not 25%!)
	}

	finResult, err := schemaflow.Audit[FinancialReport](report, schemaflow.AuditOptions{
		Policies: []string{
			"Net income must equal revenue minus expenses",
			"Total payroll should equal headcount times average salary",
			"Growth rate must match calculation from previous quarter",
			"All financial figures must be positive",
		},
		Categories: []string{"consistency", "quality"},
	})

	if err != nil {
		schemaflow.GetLogger().Error("Financial audit failed", "error", err)
		return
	}

	fmt.Printf("Financial Report Audit (%s):\n", report.Quarter)
	fmt.Printf("  Passes Audit: %v\n", finResult.Summary.PassesAudit)
	fmt.Printf("  Consistency Issues: %d\n", finResult.Summary.ByCategory["consistency"])

	fmt.Printf("\nFindings:\n")
	for _, f := range finResult.Findings {
		fmt.Printf("  - [%s] %s\n", f.Field, f.Issue)
		if f.Recommendation != "" {
			fmt.Printf("    Fix: %s\n", f.Recommendation)
		}
	}

	fmt.Println("\n=== Audit Example Complete ===")
}
