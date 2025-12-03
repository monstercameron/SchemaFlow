package main

import (
	"fmt"
	"strings"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// UserRegistration represents user registration data
type UserRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("✅ Validate Example - User Registration Validation")
	fmt.Println("=" + string(make([]byte, 60)))

	// Test cases
	testCases := []struct {
		name string
		data UserRegistration
	}{
		{
			name: "Valid Registration",
			data: UserRegistration{
				Username: "johndoe123",
				Email:    "john.doe@example.com",
				Password: "SecureP@ssw0rd!",
				Age:      25,
				Country:  "USA",
			},
		},
		{
			name: "Invalid Email",
			data: UserRegistration{
				Username: "janedoe",
				Email:    "not-an-email",
				Password: "GoodPassword123!",
				Age:      30,
				Country:  "Canada",
			},
		},
		{
			name: "Weak Password",
			data: UserRegistration{
				Username: "bobsmith",
				Email:    "bob@example.com",
				Password: "12345",
				Age:      22,
				Country:  "UK",
			},
		},
		{
			name: "Underage User",
			data: UserRegistration{
				Username: "younguser",
				Email:    "young@example.com",
				Password: "StrongPass123!",
				Age:      15,
				Country:  "Germany",
			},
		},
	}

	validationRules := `
Validation Rules:
1. Username: 3-20 characters, alphanumeric only
2. Email: Must be valid email format
3. Password: Minimum 8 characters, must include uppercase, lowercase, number, and special character
4. Age: Must be 18 or older
5. Country: Must be a valid country name
`

	for i, tc := range testCases {
		fmt.Printf("\n%d. %s\n", i+1, tc.name)
		fmt.Println("---")
		fmt.Printf("   Username: %s\n", tc.data.Username)
		fmt.Printf("   Email: %s\n", tc.data.Email)
		fmt.Printf("   Password: %s (length: %d)\n", maskPassword(tc.data.Password), len(tc.data.Password))
		fmt.Printf("   Age: %d\n", tc.data.Age)
		fmt.Printf("   Country: %s\n", tc.data.Country)

		// Validate
		result, err := schemaflow.Validate(tc.data, validationRules)
		if err != nil {
			schemaflow.GetLogger().Error("Validation error", "error", err)
			continue
		}

		if result.Valid {
			fmt.Println()
			fmt.Println("   ✅ VALID - Registration accepted")
		} else {
			fmt.Println()
			fmt.Println("   ❌ INVALID - Errors found:")
			for _, issue := range result.Issues {
				fmt.Printf("      • %s\n", issue)
			}
		}
	}

	fmt.Println()
	fmt.Println("📊 Validation Summary:")
	fmt.Println("   Total tested: 4 registrations")
	fmt.Println("   Expected: 1 valid, 3 invalid")
	fmt.Println()
	fmt.Println("✨ Success! Validation complete")
}

func maskPassword(password string) string {
	if len(password) <= 2 {
		return "***"
	}
	return password[:2] + strings.Repeat("*", len(password)-2)
}
