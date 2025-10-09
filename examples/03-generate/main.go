package main

import (
	"encoding/json"
	"fmt"
	"os"

	schemaflow "github.com/monstercameron/SchemaFlow"
	"github.com/monstercameron/SchemaFlow/core"
	"github.com/monstercameron/SchemaFlow/ops"
)

// TestUser represents a user for testing
type TestUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
	Role     string `json:"role"`
	JoinDate string `json:"join_date"`
	Active   bool   `json:"active"`
}

// TestUserBatch represents multiple test users
type TestUserBatch struct {
	Users []TestUser `json:"users"`
	Count int        `json:"count"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		core.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		os.Exit(1)
	}

	fmt.Println("🧪 Generate Example - Test Data Generator")
	fmt.Println("=" + string(make([]byte, 50)))

	// Generate test users from a prompt
	prompt := `Generate 5 realistic test users for a social media application.
Include diverse:
- Names (from different cultures)
- Ages (18-65)
- Countries (worldwide distribution)
- Roles (user, moderator, admin)
- Join dates (within last 2 years)
- Mix of active and inactive users

Make the data realistic and varied.`

	fmt.Println("\n📝 Prompt:")
	fmt.Println(prompt)

	// Generate structured test data
	batch, err := schemaflow.Generate[TestUserBatch](
		prompt,
		ops.NewGenerateOptions().
			WithIntelligence(schemaflow.Fast).
			WithSteering("Create realistic, diverse test data with proper formatting"),
	)

	if err != nil {
		core.GetLogger().Error("Generation failed", "error", err)
		os.Exit(1)
	}

	// Display generated data
	fmt.Println("\n✅ Generated Test Users:")
	fmt.Println("---")

	for i, user := range batch.Users {
		fmt.Printf("\n👤 User %d:\n", i+1)
		fmt.Printf("   ID:        %d\n", user.ID)
		fmt.Printf("   Name:      %s\n", user.Name)
		fmt.Printf("   Email:     %s\n", user.Email)
		fmt.Printf("   Age:       %d\n", user.Age)
		fmt.Printf("   Country:   %s\n", user.Country)
		fmt.Printf("   Role:      %s\n", user.Role)
		fmt.Printf("   Join Date: %s\n", user.JoinDate)
		fmt.Printf("   Active:    %t\n", user.Active)
	}

	// Show as JSON for API usage
	jsonData, _ := json.MarshalIndent(batch, "", "  ")
	fmt.Println("\n📦 JSON Output (ready for API):")
	fmt.Println(string(jsonData))

	fmt.Printf("\n✨ Success! Generated %d realistic test users\n", batch.Count)
}
