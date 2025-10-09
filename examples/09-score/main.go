package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
	"github.com/monstercameron/SchemaFlow/core"
	"github.com/monstercameron/SchemaFlow/ops"
)

// CodeSnippet represents a code submission
type CodeSnippet struct {
	ID       int
	Author   string
	Language string
	Code     string
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		core.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	// Code snippets to evaluate
	snippets := []CodeSnippet{
		{
			ID:       1,
			Author:   "Alice",
			Language: "Go",
			Code: `func calculateTotal(prices []float64) float64 {
    total := 0.0
    for _, price := range prices {
        total += price
    }
    return total
}`,
		},
		{
			ID:       2,
			Author:   "Bob",
			Language: "Go",
			Code: `func calculateTotal(p []float64) float64 {
    var t float64
    for i:=0;i<len(p);i++ {
        t=t+p[i]
    }
    return t
}`,
		},
		{
			ID:       3,
			Author:   "Carol",
			Language: "Go",
			Code: `// CalculateTotal computes the sum of all prices in the slice.
// It returns 0.0 for an empty slice.
// Time complexity: O(n)
func CalculateTotal(prices []float64) float64 {
    if len(prices) == 0 {
        return 0.0
    }
    
    total := 0.0
    for _, price := range prices {
        if price < 0 {
            continue // Skip negative prices
        }
        total += price
    }
    return total
}`,
		},
	}

	// Scoring criteria
	criteria := []string{
		"readability",
		"maintainability",
		"documentation",
		"error handling",
		"best practices",
	}

	fmt.Println("📊 Score Example - Code Quality Assessment")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println("🎯 Evaluation Criteria:", criteria)
	fmt.Println("📏 Scale: 1-10 (10 = excellent)")

	// Score each snippet
	type ScoredSnippet struct {
		Snippet CodeSnippet
		Score   float64
	}
	var scored []ScoredSnippet

	for _, snippet := range snippets {
		fmt.Printf("📝 Evaluating Submission #%d by %s\n", snippet.ID, snippet.Author)
		fmt.Println("---")
		fmt.Println(snippet.Code)
		fmt.Println("---")

		// Score the code
		scoreOpts := ops.NewScoreOptions().
			WithScaleMin(1).
			WithScaleMax(10).
			WithCriteria(criteria)
		scoreOpts.OpOptions.Intelligence = schemaflow.Smart

		score, err := ops.Score(snippet.Code, scoreOpts)
		if err != nil {
			core.GetLogger().Error("Failed to score snippet", "snippetID", snippet.ID, "error", err)
			continue
		}

		scored = append(scored, ScoredSnippet{snippet, score})

		// Display score with visualization
		stars := int(score)
		starBar := ""
		for i := 0; i < 10; i++ {
			if i < stars {
				starBar += "⭐"
			} else {
				starBar += "☆"
			}
		}

		fmt.Printf("\n✅ Score: %.1f/10 %s\n", score, starBar)

		// Provide feedback
		if score >= 8.5 {
			fmt.Println("💎 Excellent! Production-ready code with best practices.")
		} else if score >= 7.0 {
			fmt.Println("👍 Good code quality. Minor improvements possible.")
		} else if score >= 5.0 {
			fmt.Println("⚠️  Acceptable but needs improvement.")
		} else {
			fmt.Println("❌ Needs significant refactoring.")
		}
		fmt.Println()
	}

	// Show ranking
	fmt.Println("\n🏆 Final Rankings:")
	fmt.Println("---")
	for i := len(scored) - 1; i >= 0; i-- {
		for j := 0; j < len(scored)-i-1; j++ {
			if scored[j].Score < scored[j+1].Score {
				scored[j], scored[j+1] = scored[j+1], scored[j]
			}
		}
	}

	for i, s := range scored {
		medal := "🥉"
		if i == 0 {
			medal = "🥇"
		} else if i == 1 {
			medal = "🥈"
		}
		fmt.Printf("%s %d. %s - Score: %.1f/10\n", medal, i+1, s.Snippet.Author, s.Score)
	}

	fmt.Println("\n✨ Success! All code submissions evaluated")
}
