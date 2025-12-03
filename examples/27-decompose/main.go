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

	fmt.Println("=== Decompose Example ===")

	// Example 1: Decompose a complex task
	fmt.Println("--- Example 1: Task Decomposition ---")
	task := "Build a web application with user authentication, database integration, and real-time notifications"

	opts := ops.NewDecomposeOptions().
		WithStrategy("sequential").
		WithIncludeDependencies(true).
		WithTargetParts(10)

	result, err := ops.Decompose(task, opts)
	if err != nil {
		log.Fatalf("Decomposition failed: %v", err)
	}

	fmt.Printf("Original task: %s\n\n", task)
	fmt.Printf("Decomposed into %d parts:\n", len(result.Parts))
	for i, part := range result.Parts {
		fmt.Printf("\n  %d. %s\n", i+1, part.Name)
		fmt.Printf("     Description: %s\n", part.Description)
		if len(part.Dependencies) > 0 {
			fmt.Printf("     Depends on: %v\n", part.Dependencies)
		}
	}
	fmt.Println()

	// Example 2: Hierarchical decomposition
	fmt.Println("--- Example 2: Hierarchical Decomposition ---")
	project := "Launch a successful e-commerce platform"

	hierOpts := ops.NewDecomposeOptions().
		WithStrategy("hierarchical").
		WithMaxDepth(2).
		WithIncludeDependencies(true)

	hierResult, err := ops.Decompose(project, hierOpts)
	if err != nil {
		log.Fatalf("Hierarchical decomposition failed: %v", err)
	}

	fmt.Printf("Project: %s\n\n", project)
	fmt.Printf("Hierarchical breakdown:\n")
	for _, part := range hierResult.Parts {
		indent := ""
		for i := 0; i < part.Depth; i++ {
			indent += "  "
		}
		fmt.Printf("\n%s[Depth %d] %s\n", indent, part.Depth, part.Name)
		fmt.Printf("%s  %s\n", indent, part.Description)
	}
	fmt.Println()

	// Example 3: Parallel decomposition
	fmt.Println("--- Example 3: Parallel Workstreams ---")
	initiative := "Migrate monolith to microservices architecture"

	parallelOpts := ops.NewDecomposeOptions().
		WithStrategy("parallel").
		WithIncludeEstimates(true)

	parallelResult, err := ops.Decompose(initiative, parallelOpts)
	if err != nil {
		log.Fatalf("Parallel decomposition failed: %v", err)
	}

	fmt.Printf("Initiative: %s\n\n", initiative)
	fmt.Printf("Parallel workstreams:\n")
	for i, part := range parallelResult.Parts {
		fmt.Printf("\n  Stream %d: %s\n", i+1, part.Name)
		if part.Estimate != "" {
			fmt.Printf("  Effort: %s\n", part.Estimate)
		}
		if len(part.Dependencies) == 0 {
			fmt.Println("  [Can run in parallel - no dependencies]")
		}
	}
	fmt.Println()

	// Example 4: Functional decomposition
	fmt.Println("--- Example 4: Functional Components ---")
	type SystemSpec struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Features    []string `json:"features"`
	}

	system := SystemSpec{
		Name:        "Order Processing System",
		Description: "Handle customer orders from placement to fulfillment",
		Features:    []string{"order creation", "payment processing", "inventory check", "shipping"},
	}

	type Component struct {
		Name           string   `json:"name"`
		Responsibility string   `json:"responsibility"`
		Inputs         []string `json:"inputs"`
		Outputs        []string `json:"outputs"`
	}

	funcOpts := ops.NewDecomposeOptions().WithStrategy("functional")

	funcResult, err := ops.DecomposeToSlice[SystemSpec, Component](system, funcOpts)
	if err != nil {
		log.Fatalf("Functional decomposition failed: %v", err)
	}

	fmt.Printf("System: %s\n\n", system.Name)
	fmt.Printf("Functional components:\n")
	for _, component := range funcResult {
		fmt.Printf("\n  Component: %s\n", component.Name)
		fmt.Printf("  Responsibility: %s\n", component.Responsibility)
		if len(component.Inputs) > 0 {
			fmt.Printf("  Inputs: %v\n", component.Inputs)
		}
		if len(component.Outputs) > 0 {
			fmt.Printf("  Outputs: %v\n", component.Outputs)
		}
	}

	fmt.Println("\n=== Decompose Example Complete ===")
}
