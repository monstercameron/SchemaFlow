package main

import (
	"fmt"
	"log"

	"github.com/monstercameron/SchemaFlow/ops"
)

func main() {
	fmt.Println("=== SchemaFlow Suggest Operation Examples ===")

	// Example 1: Basic suggestions for data processing workflow
	fmt.Println("\n1. Data Processing Workflow Suggestions:")
	currentState := map[string]any{
		"task":          "ETL pipeline optimization",
		"data_volume":   "100GB daily",
		"issues":        []string{"slow processing", "memory usage", "error rates"},
		"current_stack": []string{"Python", "Pandas", "PostgreSQL"},
	}

	suggestions, err := ops.Suggest[string](currentState,
		ops.NewSuggestOptions().
			WithTopN(5).
			WithDomain("data-engineering"))
	if err != nil {
		log.Printf("Example 1 failed: %v", err)
	} else {
		fmt.Println("Current state:", currentState["task"])
		fmt.Println("Suggestions:")
		for i, suggestion := range suggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion)
		}
	}

	// Example 2: API endpoint suggestions
	fmt.Println("\n2. API Design Suggestions:")
	apiContext := map[string]any{
		"resource":    "user profiles",
		"operations":  []string{"create", "read", "update", "delete"},
		"constraints": []string{"RESTful", "versioned", "authenticated"},
		"data_fields": []string{"id", "name", "email", "preferences"},
	}

	apiSuggestions, err := ops.Suggest[string](apiContext,
		ops.NewSuggestOptions().
			WithStrategy(ops.SuggestPattern).
			WithTopN(4).
			WithDomain("api-design"))
	if err != nil {
		log.Printf("Example 2 failed: %v", err)
	} else {
		fmt.Println("Resource:", apiContext["resource"])
		fmt.Println("API endpoint suggestions:")
		for i, suggestion := range apiSuggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion)
		}
	}

	// Example 3: Configuration optimization suggestions
	fmt.Println("\n3. Configuration Optimization:")
	configContext := map[string]any{
		"system":       "web application",
		"current_load": "1000 req/sec",
		"issues":       []string{"high latency", "memory leaks"},
		"environment":  "production",
		"tech_stack":   []string{"Go", "PostgreSQL", "Redis"},
	}

	configSuggestions, err := ops.Suggest[string](configContext,
		ops.NewSuggestOptions().
			WithStrategy(ops.SuggestHybrid).
			WithTopN(3).
			WithIncludeReasons(true))
	if err != nil {
		log.Printf("Example 3 failed: %v", err)
	} else {
		fmt.Println("System:", configContext["system"])
		fmt.Println("Optimization suggestions:")
		for i, suggestion := range configSuggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion)
		}
	}

	// Example 4: Custom types - Action suggestions
	fmt.Println("\n4. Workflow Action Suggestions:")

	type Action struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		Category    string `json:"category"`
	}

	workflowContext := map[string]any{
		"phase":    "data validation",
		"dataset":  "customer_records",
		"issues":   []string{"inconsistent formats", "missing required fields"},
		"deadline": "end of week",
	}

	actionSuggestions, err := ops.Suggest[Action](workflowContext,
		ops.NewSuggestOptions().
			WithTopN(3).
			WithDomain("data-quality"))
	if err != nil {
		log.Printf("Example 4 failed: %v", err)
	} else {
		fmt.Println("Workflow phase:", workflowContext["phase"])
		fmt.Println("Action suggestions:")
		for i, action := range actionSuggestions {
			fmt.Printf("  %d. %s (%s) - %s\n", i+1, action.Name, action.Priority, action.Description)
		}
	}

	// Example 5: Error handling and constraints
	fmt.Println("\n5. Error Recovery Suggestions:")
	errorContext := map[string]any{
		"error":       "database connection timeout",
		"component":   "user authentication service",
		"frequency":   "5 times per hour",
		"impact":      "user login failures",
		"constraints": []string{"zero downtime", "preserve user sessions"},
	}

	recoverySuggestions, err := ops.Suggest[string](errorContext,
		ops.NewSuggestOptions().
			WithConstraints([]string{"zero downtime", "preserve user sessions"}).
			WithCategories([]string{"reliability", "monitoring"}).
			WithTopN(3))
	if err != nil {
		log.Printf("Example 5 failed: %v", err)
	} else {
		fmt.Println("Error:", errorContext["error"])
		fmt.Println("Recovery suggestions:")
		for i, suggestion := range recoverySuggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion)
		}
	}

	fmt.Println("\n=== Suggest Operation Examples Complete ===")
	fmt.Println("\nNote: The Suggest operation generates context-aware recommendations")
	fmt.Println("to help optimize workflows, solve problems, and guide decision-making.")
	fmt.Println("Suggestions are generated using LLM intelligence and can be customized")
	fmt.Println("with domains, constraints, categories, and ranking preferences.")
}
