package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// InternalUser is the full internal representation
type InternalUser struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	PasswordHash  string   `json:"password_hash"`
	SSN           string   `json:"ssn"`
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	DateOfBirth   string   `json:"date_of_birth"`
	Phone         string   `json:"phone"`
	Address       Address  `json:"address"`
	Roles         []string `json:"roles"`
	CreatedAt     string   `json:"created_at"`
	LastLogin     string   `json:"last_login"`
	InternalNotes string   `json:"internal_notes"`
	CreditScore   int      `json:"credit_score"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

// PublicProfile is what external users/APIs see
type PublicProfile struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Location    string `json:"location"`
	MemberSince string `json:"member_since"`
	IsAdmin     bool   `json:"is_admin"`
	LastActive  string `json:"last_active"`
}

// AdminView is what admins see
type AdminView struct {
	ID        string   `json:"id"`
	FullName  string   `json:"full_name"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	Age       int      `json:"age"`
	Location  string   `json:"location"`
	Roles     []string `json:"roles"`
	CreatedAt string   `json:"created_at"`
	Notes     string   `json:"notes"`
	RiskLevel string   `json:"risk_level"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	fmt.Println("=== Project Example ===")

	internalUser := InternalUser{
		ID:           "USR-12345",
		Email:        "john.doe@example.com",
		PasswordHash: "$2a$10$xyz123...",
		SSN:          "123-45-6789",
		FirstName:    "John",
		LastName:     "Doe",
		DateOfBirth:  "1990-05-15",
		Phone:        "+1-555-123-4567",
		Address: Address{
			Street:  "123 Main St",
			City:    "Boston",
			State:   "MA",
			ZipCode: "02101",
			Country: "USA",
		},
		Roles:         []string{"user", "premium", "beta_tester"},
		CreatedAt:     "2022-01-15T10:30:00Z",
		LastLogin:     "2024-01-20T14:25:00Z",
		InternalNotes: "VIP customer, handle with care",
		CreditScore:   780,
	}

	// Example 1: Project to public profile (exclude sensitive data)
	fmt.Println("\n--- Example 1: Public Profile ---")

	publicResult, err := schemaflow.Project[InternalUser, PublicProfile](internalUser, schemaflow.ProjectOptions{
		Mappings: map[string]string{
			"id":         "user_id",
			"created_at": "member_since",
			"last_login": "last_active",
		},
		Exclude:      []string{"password_hash", "ssn", "credit_score", "internal_notes", "phone"},
		InferMissing: true,
		Steering:     "Combine first_name and last_name into display_name. Combine city and state into location. Check if 'admin' is in roles for is_admin.",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Public projection failed", "error", err)
		return
	}

	fmt.Printf("Public Profile:\n")
	fmt.Printf("  User ID: %s\n", publicResult.Projected.UserID)
	fmt.Printf("  Display Name: %s\n", publicResult.Projected.DisplayName)
	fmt.Printf("  Email: %s\n", publicResult.Projected.Email)
	fmt.Printf("  Location: %s\n", publicResult.Projected.Location)
	fmt.Printf("  Member Since: %s\n", publicResult.Projected.MemberSince)
	fmt.Printf("  Is Admin: %v\n", publicResult.Projected.IsAdmin)
	fmt.Printf("  Last Active: %s\n", publicResult.Projected.LastActive)

	fmt.Printf("\nField Mappings:\n")
	for _, m := range publicResult.Mappings {
		source := m.SourceField
		if source == "" {
			source = "[inferred]"
		}
		fmt.Printf("  %s → %s (%s)\n", source, m.TargetField, m.Method)
		if m.Transformation != "" {
			fmt.Printf("    Transform: %s\n", m.Transformation)
		}
	}

	if len(publicResult.Lost) > 0 {
		fmt.Printf("\nExcluded/Lost Fields: %v\n", publicResult.Lost)
	}
	if len(publicResult.Inferred) > 0 {
		fmt.Printf("Inferred Fields: %v\n", publicResult.Inferred)
	}
	fmt.Printf("Confidence: %.0f%%\n", publicResult.Confidence*100)

	// Example 2: Project to admin view (include more, derive some)
	fmt.Println("\n--- Example 2: Admin View ---")

	adminResult, err := schemaflow.Project[InternalUser, AdminView](internalUser, schemaflow.ProjectOptions{
		Exclude:      []string{"password_hash", "ssn"},
		InferMissing: true,
		Steering:     "Calculate age from date_of_birth (current year 2024). Determine risk_level based on credit_score (750+ = low, 650-749 = medium, below 650 = high). Combine city and state for location.",
	})

	if err != nil {
		schemaflow.GetLogger().Error("Admin projection failed", "error", err)
		return
	}

	fmt.Printf("Admin View:\n")
	fmt.Printf("  ID: %s\n", adminResult.Projected.ID)
	fmt.Printf("  Full Name: %s\n", adminResult.Projected.FullName)
	fmt.Printf("  Email: %s\n", adminResult.Projected.Email)
	fmt.Printf("  Phone: %s\n", adminResult.Projected.Phone)
	fmt.Printf("  Age: %d\n", adminResult.Projected.Age)
	fmt.Printf("  Location: %s\n", adminResult.Projected.Location)
	fmt.Printf("  Roles: %v\n", adminResult.Projected.Roles)
	fmt.Printf("  Created At: %s\n", adminResult.Projected.CreatedAt)
	fmt.Printf("  Notes: %s\n", adminResult.Projected.Notes)
	fmt.Printf("  Risk Level: %s\n", adminResult.Projected.RiskLevel)

	fmt.Printf("\nInferred Fields: %v\n", adminResult.Inferred)
	fmt.Printf("Confidence: %.0f%%\n", adminResult.Confidence*100)

	fmt.Println("\n=== Project Example Complete ===")
}
