// 18-explain: Generate human-readable explanations of data structures
// Intelligence: Fast (Cerebras gpt-oss-120b)
// Expectations:
// - Explains user profile data for different audiences
// - Technical: focuses on struct design, serialization, patterns
// - Executive: focuses on business value, engagement metrics
// - Beginner: uses simple language and step-by-step format

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	schemaflow "github.com/monstercameron/SchemaFlow"
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

// UserProfile represents a user's profile data
type UserProfile struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Age        int      `json:"age"`
	Role       string   `json:"role"`
	IsVerified bool     `json:"is_verified"`
	Interests  []string `json:"interests"`
	JoinedAt   string   `json:"joined_at"`
}

func main() {
	loadEnv()

	fmt.Println("📚 Explain Example - Generate Human-Readable Explanations")
	fmt.Println("=" + string(make([]byte, 60)))

	// Initialize SchemaFlow with Fast intelligence (Cerebras)
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	// Sample user data to explain
	userData := UserProfile{
		ID:         12345,
		Name:       "Alice Johnson",
		Email:      "alice.johnson@example.com",
		Age:        28,
		Role:       "premium",
		IsVerified: true,
		Interests:  []string{"technology", "data science", "machine learning"},
		JoinedAt:   "2023-01-15",
	}

	fmt.Println("\n📊 Data to Explain:")
	fmt.Printf("   %+v\n", userData)

	// Example 1: Technical audience explanation
	fmt.Println("\n" + "─" + string(make([]byte, 60)))
	fmt.Println("1️⃣  Technical Audience Explanation")
	fmt.Println("─" + string(make([]byte, 60)))

	techResult, err := schemaflow.Explain(userData,
		schemaflow.NewExplainOptions().
			WithAudience("technical").
			WithDepth(2).
			WithFormat("structured").
			WithFocus("implementation").
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("\n   📝 Summary: %s\n", techResult.Summary)
		fmt.Printf("\n   📖 Explanation:\n   %s\n", techResult.Explanation)
		if len(techResult.KeyPoints) > 0 {
			fmt.Println("\n   🔑 Key Points:")
			for _, point := range techResult.KeyPoints {
				fmt.Printf("      • %s\n", point)
			}
		}
	}

	// Example 2: Executive audience explanation
	fmt.Println("\n" + "─" + string(make([]byte, 60)))
	fmt.Println("2️⃣  Executive Audience Explanation")
	fmt.Println("─" + string(make([]byte, 60)))

	execResult, err := schemaflow.Explain(userData,
		schemaflow.NewExplainOptions().
			WithAudience("executive").
			WithDepth(1).
			WithFormat("bullet-points").
			WithFocus("benefits").
			WithContext("Business intelligence and user analytics data").
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("\n   📝 Summary: %s\n", execResult.Summary)
		fmt.Printf("\n   📖 Explanation:\n   %s\n", execResult.Explanation)
		if len(execResult.KeyPoints) > 0 {
			fmt.Println("\n   🔑 Key Points:")
			for _, point := range execResult.KeyPoints {
				fmt.Printf("      • %s\n", point)
			}
		}
	}

	// Example 3: Beginner-friendly explanation
	fmt.Println("\n" + "─" + string(make([]byte, 60)))
	fmt.Println("3️⃣  Beginner-Friendly Explanation")
	fmt.Println("─" + string(make([]byte, 60)))

	beginnerResult, err := schemaflow.Explain(userData,
		schemaflow.NewExplainOptions().
			WithAudience("beginner").
			WithDepth(1).
			WithFormat("step-by-step").
			WithFocus("overview").
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("\n   📝 Summary: %s\n", beginnerResult.Summary)
		fmt.Printf("\n   📖 Explanation:\n   %s\n", beginnerResult.Explanation)
		if len(beginnerResult.KeyPoints) > 0 {
			fmt.Println("\n   🔑 Key Points:")
			for _, point := range beginnerResult.KeyPoints {
				fmt.Printf("      • %s\n", point)
			}
		}
	}

	fmt.Println()
	fmt.Println("✨ Success! Data explained for different audiences")
}
