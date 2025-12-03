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

	fmt.Println("=== Critique Example ===")

	// Example 1: Critique an essay
	fmt.Println("--- Example 1: Essay Critique ---")
	essay := `Climate change is a big problem. We should do something about it.
	Scientists say it's getting warmer. This is bad for polar bears and other animals.
	We need to use less energy and drive electric cars. The government should make laws.
	In conclusion, climate change is important and we need to act now.`

	opts := ops.NewCritiqueOptions().
		WithCriteria([]string{"argument_strength", "evidence", "clarity", "structure"}).
		WithIncludeSuggestions(true).
		WithIncludeFixes(true).
		WithIncludePositives(true).
		WithStyle("constructive")

	result, err := ops.Critique(essay, opts)
	if err != nil {
		log.Fatalf("Critique failed: %v", err)
	}

	fmt.Printf("Overall Score: %.2f/1.00\n\n", result.OverallScore)

	fmt.Println("Criteria Scores:")
	for criterion, score := range result.CriteriaScores {
		fmt.Printf("  %s: %.2f\n", criterion, score)
	}

	fmt.Println("\nIssues Found:")
	for i, issue := range result.Issues {
		fmt.Printf("\n  %d. [%s] %s\n", i+1, issue.Severity, issue.Criterion)
		fmt.Printf("     %s\n", issue.Description)
		if issue.Suggestion != "" {
			fmt.Printf("     Suggestion: %s\n", issue.Suggestion)
		}
		if issue.Fix != "" {
			fmt.Printf("     Fix: %s\n", issue.Fix)
		}
	}

	if len(result.Positives) > 0 {
		fmt.Println("\nPositive Feedback:")
		for _, pos := range result.Positives {
			fmt.Printf("  - %s: %s\n", pos.Criterion, pos.Description)
		}
	}

	fmt.Printf("\nSummary: %s\n\n", result.Summary)

	// Example 2: Code review critique
	fmt.Println("--- Example 2: Code Review ---")
	code := `
func processData(data []string) {
    for i := 0; i < len(data); i++ {
        item := data[i]
        if item != "" {
            result := doSomething(item)
            fmt.Println(result)
        }
    }
}
`

	codeOpts := ops.NewCritiqueOptions().
		WithDomain("software").
		WithRubric(map[string]string{
			"readability":    "Is the code easy to understand?",
			"efficiency":     "Are there performance issues?",
			"best_practices": "Does it follow Go idioms?",
			"error_handling": "Are errors handled properly?",
		}).
		WithStyle("balanced")

	codeResult, err := ops.Critique(code, codeOpts)
	if err != nil {
		log.Fatalf("Code critique failed: %v", err)
	}

	fmt.Printf("Code Quality Score: %.2f\n", codeResult.OverallScore)
	fmt.Println("\nCode Review Feedback:")
	for _, issue := range codeResult.Issues {
		fmt.Printf("  [%s] %s: %s\n", issue.Severity, issue.Criterion, issue.Description)
		if issue.Fix != "" {
			fmt.Printf("    Fix: %s\n", issue.Fix)
		}
	}
	fmt.Println()

	// Example 3: Product description critique
	fmt.Println("--- Example 3: Marketing Content Critique ---")
	productDesc := `
Our new SuperWidget 3000 is the best product ever made! 
It can do everything you need and more. 
Buy now and you won't regret it!
Limited time offer - act fast!
`

	marketingOpts := ops.NewCritiqueOptions().
		WithCriteria([]string{"persuasiveness", "clarity", "credibility", "call_to_action"}).
		WithAudience("general consumers").
		WithMaxIssues(5).
		WithSeverityFilter("all")

	marketingResult, err := ops.Critique(productDesc, marketingOpts)
	if err != nil {
		log.Fatalf("Marketing critique failed: %v", err)
	}

	fmt.Printf("Marketing Effectiveness: %.2f\n", marketingResult.OverallScore)
	fmt.Printf("Summary: %s\n", marketingResult.Summary)

	if len(marketingResult.TopPriorities) > 0 {
		fmt.Println("\nTop Priorities:")
		for i, priority := range marketingResult.TopPriorities {
			fmt.Printf("  %d. %s\n", i+1, priority)
		}
	}

	// Example 4: Harsh critique mode
	fmt.Println("\n--- Example 4: Harsh Critique Mode ---")
	presentation := "We did stuff this quarter. Sales were okay. Next quarter will be better probably."

	harshOpts := ops.NewCritiqueOptions().
		WithCriteria([]string{"professionalism", "specificity", "impact"}).
		WithStyle("harsh").
		WithIncludePositives(false)

	harshResult, err := ops.Critique(presentation, harshOpts)
	if err != nil {
		log.Fatalf("Harsh critique failed: %v", err)
	}

	fmt.Printf("Score: %.2f\n", harshResult.OverallScore)
	fmt.Println("Critical Issues:")
	for _, issue := range harshResult.Issues {
		if issue.Severity == "critical" || issue.Severity == "major" {
			fmt.Printf("  - %s\n", issue.Description)
		}
	}

	fmt.Println("\n=== Critique Example Complete ===")
}
