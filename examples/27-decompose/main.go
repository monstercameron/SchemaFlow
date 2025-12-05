// Example 27: Decompose Operation
// Breaks down complex items into smaller, manageable parts using LLM

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	schemaflow "github.com/monstercameron/SchemaFlow"
	"github.com/monstercameron/SchemaFlow/internal/ops"
	"github.com/monstercameron/SchemaFlow/internal/types"
)

// loadEnv loads environment variables from .env file
func loadEnv() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				log.Fatal("Error loading .env file")
			}
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	log.Fatal(".env file not found")
}

func main() {
	loadEnv()

	// Ensure environment is configured
	if os.Getenv("SCHEMAFLOW_API_KEY") == "" {
		log.Fatal("SCHEMAFLOW_API_KEY environment variable not set")
	}

	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		log.Fatalf("Failed to initialize SchemaFlow: %v", err)
	}

	fmt.Println("=== Decompose Example ===")
	fmt.Println("Breaks down complex items into smaller parts")
	fmt.Println()

	// Business Use Case: Sprint Planning - Break Epic into User Stories
	fmt.Println("--- Business Use Case: Epic → Sprint Backlog ---")

	type Epic struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Team        string `json:"team"`
	}

	epic := Epic{
		Title:       "Multi-Tenant Billing System",
		Description: "Build billing for SaaS platform with usage tracking and invoicing",
		Team:        "Platform Engineering",
	}

	fmt.Println("INPUT: Epic struct")
	fmt.Printf("  Title:       %q\n", epic.Title)
	fmt.Printf("  Description: %q\n", epic.Description)
	fmt.Printf("  Team:        %q\n\n", epic.Team)

	opts := ops.NewDecomposeOptions().
		WithStrategy("sequential").
		WithIncludeDependencies(true).
		WithTargetParts(5).
		WithIntelligence(types.Smart)

	result, err := ops.Decompose(epic, opts)
	if err != nil {
		log.Fatalf("Decomposition failed: %v", err)
	}

	fmt.Println("OUTPUT: DecomposeResult[Epic] → Sprint-Ready User Stories")
	fmt.Printf("  Total Stories: %d\n", len(result.Parts))
	fmt.Println("  ┌─────────────────────────────────────────────────────────────")
	for i, part := range result.Parts {
		fmt.Printf("  │ Story %d: %s\n", i+1, part.Name)
		fmt.Printf("  │   Description: %s\n", part.Description)
		if len(part.Dependencies) > 0 {
			fmt.Printf("  │   Blocked By: %v\n", part.Dependencies)
		}
		if i < len(result.Parts)-1 {
			fmt.Println("  │")
		}
	}
	fmt.Println("  └─────────────────────────────────────────────────────────────")

	fmt.Println("\n--- Business Use Case: Incident → Runbook Steps ---")

	incident := "P1: Checkout failures during peak traffic"

	fmt.Println("INPUT: string")
	fmt.Printf("  Incident: %q\n\n", incident)

	incidentOpts := ops.NewDecomposeOptions().
		WithStrategy("sequential").
		WithTargetParts(4).
		WithIntelligence(types.Smart)

	incidentResult, err := ops.Decompose(incident, incidentOpts)
	if err != nil {
		log.Fatalf("Incident decomposition failed: %v", err)
	}

	fmt.Println("OUTPUT: DecomposeResult[string] → Troubleshooting Runbook")
	fmt.Printf("  Total Steps: %d\n", len(incidentResult.Parts))
	fmt.Println("  ┌─────────────────────────────────────────────────────────────")
	for i, part := range incidentResult.Parts {
		fmt.Printf("  │ Step %d: %s\n", i+1, part.Name)
		fmt.Printf("  │   Action: %s\n", part.Description)
		if i < len(incidentResult.Parts)-1 {
			fmt.Println("  │")
		}
	}
	fmt.Println("  └─────────────────────────────────────────────────────────────")

	fmt.Println("\n=== Decompose Example Complete ===")
}
