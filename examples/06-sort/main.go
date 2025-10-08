package main

import (
	"fmt"
	"log"

	schemaflow "github.com/monstercameron/SchemaFlow"
	"github.com/monstercameron/SchemaFlow/ops"
)

// Task represents a work task
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Impact      string `json:"impact"`
	Effort      string `json:"effort"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		log.Fatalf("Failed to initialize SchemaFlow: %v", err)
	}

	// Unsorted tasks
	tasks := []Task{
		{
			ID:          1,
			Title:       "Update documentation",
			Description: "Update API documentation for v2.0 release",
			Deadline:    "Next week",
			Impact:      "Low - only affects developers",
			Effort:      "2 hours",
		},
		{
			ID:          2,
			Title:       "Fix critical security vulnerability",
			Description: "Patch SQL injection vulnerability in login system",
			Deadline:    "Today",
			Impact:      "Critical - affects all users, security risk",
			Effort:      "4 hours",
		},
		{
			ID:          3,
			Title:       "Add dark mode feature",
			Description: "Implement dark mode toggle in settings",
			Deadline:    "Next month",
			Impact:      "Medium - user requested feature",
			Effort:      "3 days",
		},
		{
			ID:          4,
			Title:       "Database backup failing",
			Description: "Automated backups haven't run in 3 days",
			Deadline:    "ASAP",
			Impact:      "High - data loss risk",
			Effort:      "1 hour",
		},
		{
			ID:          5,
			Title:       "Refactor payment module",
			Description: "Clean up legacy payment processing code",
			Deadline:    "Q2 2025",
			Impact:      "Low - technical debt",
			Effort:      "1 week",
		},
	}

	fmt.Println("📋 Sort Example - Task Prioritization")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println("\n📥 Unsorted Tasks:")
	for _, t := range tasks {
		fmt.Printf("  %d. %s (Deadline: %s)\n", t.ID, t.Title, t.Deadline)
	}

	// Sort tasks by priority (urgency + impact + effort)
	sortOpts := ops.NewSortOptions().WithCriteria("Priority by: 1) Urgency (deadline), 2) Business impact, 3) Effort (quick wins first)")
	sortOpts.OpOptions.Intelligence = schemaflow.Smart
	sortOpts.OpOptions.Steering = "Consider deadline urgency, business impact, and effort. Quick high-impact tasks should be prioritized."

	sortedTasks, err := schemaflow.Sort(tasks, sortOpts)

	if err != nil {
		log.Fatalf("Sorting failed: %v", err)
	}

	// Display sorted tasks
	fmt.Println("\n✅ Prioritized Tasks (Highest Priority First):")
	fmt.Println("---")
	for i, t := range sortedTasks {
		fmt.Printf("\n%d. 🎯 %s\n", i+1, t.Title)
		fmt.Printf("   Deadline: %s\n", t.Deadline)
		fmt.Printf("   Impact:   %s\n", t.Impact)
		fmt.Printf("   Effort:   %s\n", t.Effort)
		fmt.Printf("   Why:      ")
		if i == 0 {
			fmt.Println("Critical security issue with immediate deadline")
		} else if i == 1 {
			fmt.Println("High impact with minimal effort - quick win")
		} else if i == 2 {
			fmt.Println("Quick task before moving to longer projects")
		} else {
			fmt.Println("Can be scheduled for later")
		}
	}

	fmt.Println("\n✨ Success! Tasks intelligently prioritized")
}
