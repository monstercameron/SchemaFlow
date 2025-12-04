package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Verify Example ===")

	// Example 1: Fact-check claims
	fmt.Println("--- Example 1: Fact Checking ---")
	claims := `
	1. The Earth is approximately 4.5 billion years old.
	2. Water boils at 100 degrees Celsius at sea level.
	3. The capital of Australia is Sydney.
	4. Humans have 206 bones in their body.
	5. The speed of light is about 300,000 km/s.
	`

	opts := schemaflow.NewVerifyOptions().
		WithCheckFacts(true).
		WithIncludeEvidence(true).
		WithExplainReasoning(true).
		WithStrictness("moderate")

	result, err := schemaflow.Verify(claims, opts)
	if err != nil {
		schemaflow.GetLogger().Error("Verification failed", "error", err)
		return
	}

	fmt.Printf("Overall Verdict: %s (Trust Score: %.2f)\n\n", result.OverallVerdict, result.TrustScore)

	fmt.Println("Individual Claims:")
	for _, claim := range result.Claims {
		verdict := claim.Verdict
		switch verdict {
		case "verified":
			verdict = "✓ " + verdict
		case "false":
			verdict = "✗ " + verdict
		case "partially_true":
			verdict = "◐ " + verdict
		default:
			verdict = "? " + verdict
		}

		fmt.Printf("\n  %s (%.0f%% confidence)\n", verdict, claim.Confidence*100)
		fmt.Printf("  Claim: %s\n", claim.Claim)
		if claim.Reasoning != "" {
			fmt.Printf("  Reasoning: %s\n", claim.Reasoning)
		}
		if claim.Corrections != "" {
			fmt.Printf("  Correction: %s\n", claim.Corrections)
		}
	}
	fmt.Println()

	// Example 2: Check logical consistency
	fmt.Println("--- Example 2: Logic Checking ---")
	argument := `
	Premise 1: All birds can fly.
	Premise 2: Penguins are birds.
	Conclusion: Therefore, penguins can fly.
	
	This means we can use penguins as aerial messengers.
	`

	logicOpts := schemaflow.NewVerifyOptions().
		WithCheckLogic(true).
		WithCheckFacts(true).
		WithExplainReasoning(true)

	logicResult, err := schemaflow.Verify(argument, logicOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Logic verification failed", "error", err)
		return
	}

	fmt.Printf("Argument Validity: %s\n", logicResult.OverallVerdict)
	fmt.Printf("Trust Score: %.2f\n", logicResult.TrustScore)

	if len(logicResult.LogicIssues) > 0 {
		fmt.Println("\nLogic Issues:")
		for _, issue := range logicResult.LogicIssues {
			fmt.Printf("  [%s] %s: %s\n", issue.Severity, issue.Type, issue.Description)
		}
	}

	fmt.Printf("\nSummary: %s\n\n", logicResult.Summary)

	// Example 3: Verify against sources
	fmt.Println("--- Example 3: Source-Based Verification ---")
	articleClaim := "The company reported $5 billion in Q3 revenue, a 20% increase from last year."

	sources := []any{
		map[string]any{
			"type":       "earnings_report",
			"q3_revenue": 4800000000, // Actual: $4.8B
			"yoy_growth": 0.18,       // Actual: 18%
		},
		map[string]any{
			"type":  "press_release",
			"quote": "We achieved strong growth this quarter with revenue approaching $5 billion.",
		},
	}

	sourceOpts := schemaflow.NewVerifyOptions().
		WithSources(sources).
		WithStrictness("strict").
		WithIncludeEvidence(true)

	sourceResult, err := schemaflow.Verify(articleClaim, sourceOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Source verification failed", "error", err)
		return
	}

	fmt.Printf("Claim: %s\n", articleClaim)
	fmt.Printf("Verdict: %s\n", sourceResult.OverallVerdict)

	for _, claim := range sourceResult.Claims {
		if len(claim.Evidence) > 0 {
			fmt.Println("Evidence:")
			for _, e := range claim.Evidence {
				fmt.Printf("  - %s\n", e)
			}
		}
		if claim.Corrections != "" {
			fmt.Printf("Correction: %s\n", claim.Corrections)
		}
	}
	fmt.Println()

	// Example 4: Check internal consistency
	fmt.Println("--- Example 4: Consistency Checking ---")
	document := `
	In the introduction, we state that the project started in 2020.
	The methodology section mentions data collection began in 2018.
	Results show analysis of 5 years of data from 2019-2023.
	The budget shows $1 million in Year 1 and $2 million in Year 2.
	The total budget is listed as $2.5 million.
	`

	consistencyOpts := schemaflow.NewVerifyOptions().
		WithCheckConsistency(true).
		WithCheckLogic(true)

	consistencyResult, err := schemaflow.Verify(document, consistencyOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Consistency check failed", "error", err)
		return
	}

	fmt.Printf("Document Consistency: %s\n", consistencyResult.OverallVerdict)

	if len(consistencyResult.ConsistencyIssues) > 0 {
		fmt.Println("\nInconsistencies Found:")
		for _, issue := range consistencyResult.ConsistencyIssues {
			fmt.Printf("  [%s] %s\n", issue.Type, issue.Description)
			if len(issue.Items) > 0 {
				fmt.Printf("    Conflicting items: %v\n", issue.Items)
			}
			if issue.Suggestion != "" {
				fmt.Printf("    Suggestion: %s\n", issue.Suggestion)
			}
		}
	}
	fmt.Println()

	// Example 5: Single claim verification
	fmt.Println("--- Example 5: Single Claim ---")
	singleClaim := "The Great Wall of China is visible from space with the naked eye."

	claimResult, err := schemaflow.VerifyClaim(singleClaim, schemaflow.NewVerifyOptions().WithExplainReasoning(true))
	if err != nil {
		schemaflow.GetLogger().Error("Claim verification failed", "error", err)
		return
	}

	fmt.Printf("Claim: %s\n", singleClaim)
	fmt.Printf("Verdict: %s (%.0f%% confidence)\n", claimResult.Verdict, claimResult.Confidence*100)
	fmt.Printf("Reasoning: %s\n", claimResult.Reasoning)

	fmt.Println("\n=== Verify Example Complete ===")
}
