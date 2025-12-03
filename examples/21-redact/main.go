package main

import (
	"fmt"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

func main() {
	fmt.Println("=== SchemaFlow Redact Operation Examples ===")

	// Example 1: Basic text redaction
	fmt.Println("1. Basic Text Redaction:")
	text1 := "Contact john@example.com or call 555-123-4567 for support."
	fmt.Printf("   Original: %q\n", text1)

	result1, _ := schemaflow.Redact(text1, schemaflow.NewRedactOptions().WithCategories([]string{"PII"}))
	fmt.Printf("   Redacted: %q\n", result1)

	// Example 2: Different redaction strategies
	fmt.Println("\n2. Redaction Strategies:")
	email := "user@company.com"

	strategies := []struct {
		name string
		opts schemaflow.RedactOptions
	}{
		{"Mask", schemaflow.NewRedactOptions().WithCategories([]string{"PII"}).WithStrategy(schemaflow.RedactMask)},
		{"Nil", schemaflow.NewRedactOptions().WithCategories([]string{"PII"}).WithStrategy(schemaflow.RedactNil)},
		{"Jumble", schemaflow.NewRedactOptions().WithCategories([]string{"PII"}).WithStrategy(schemaflow.RedactJumble).WithJumbleSeed(42)},
	}

	for _, strategy := range strategies {
		result, _ := schemaflow.Redact(email, strategy.opts)
		fmt.Printf("   %s: %q → %q\n", strategy.name, email, result)
	}

	// Example 3: Struct redaction
	fmt.Println("\n3. Struct Field Redaction:")
	type User struct {
		Name     string `redact:"PII"`
		Email    string
		Password string
		Age      int
	}

	user := User{
		Name:     "John Smith",
		Email:    "john.smith@company.com",
		Password: "secret123",
		Age:      30,
	}

	fmt.Printf("   Original: %+v\n", user)
	redactedUser, _ := schemaflow.Redact(user, schemaflow.NewRedactOptions().WithCategories([]string{"PII"}))
	fmt.Printf("   Redacted: %+v\n", redactedUser)

	// Example 4: Struct field name detection
	fmt.Println("\n4. Field Name-Based Detection:")
	type Config struct {
		DatabasePassword string
		APIToken         string
		NormalField      string
	}

	config := Config{
		DatabasePassword: "dbpass123",
		APIToken:         "token456",
		NormalField:      "normal",
	}

	fmt.Printf("   Original: %+v\n", config)
	redactedConfig, _ := schemaflow.Redact(config, schemaflow.NewRedactOptions().WithCategories([]string{"secrets"}))
	fmt.Printf("   Redacted: %+v\n", redactedConfig)

	// Example 5: Map redaction
	fmt.Println("\n5. Map Value Redaction:")
	data := map[string]string{
		"email":       "contact@business.com",
		"password":    "password: mysecret",
		"full_name":   "Jane Doe",
		"description": "Regular text here",
	}

	fmt.Printf("   Original: %+v\n", data)
	redactedData, _ := schemaflow.Redact(data, schemaflow.NewRedactOptions().WithCategories([]string{"PII", "secrets"}))
	fmt.Printf("   Redacted: %+v\n", redactedData)

	// Example 6: Slice redaction
	fmt.Println("\n6. Slice Element Redaction:")
	emails := []string{
		"valid@email.com",
		"another@domain.org",
		"not an email",
	}

	fmt.Printf("   Original: %v\n", emails)
	redactedEmails, _ := schemaflow.Redact(emails, schemaflow.NewRedactOptions().WithCategories([]string{"PII"}))
	fmt.Printf("   Redacted: %v\n", redactedEmails)

	// Example 7: Jumble variations
	fmt.Println("\n7. Jumble Strategy Variations:")
	phone := "555-123-4567"
	name := "Alice Johnson"

	jumbleOpts := []struct {
		desc string
		opts schemaflow.RedactOptions
	}{
		{"Basic Jumble", schemaflow.NewRedactOptions().WithCategories([]string{"PII"}).WithStrategy(schemaflow.RedactJumble).WithJumbleMode(schemaflow.JumbleBasic).WithJumbleSeed(123)},
		{"Type-Aware Jumble", schemaflow.NewRedactOptions().WithCategories([]string{"PII"}).WithStrategy(schemaflow.RedactJumble).WithJumbleMode(schemaflow.JumbleTypeAware).WithJumbleSeed(123)},
	}

	for _, item := range []string{phone, name} {
		fmt.Printf("   Item: %q\n", item)
		for _, jumbleOpt := range jumbleOpts {
			result, _ := schemaflow.Redact(item, jumbleOpt.opts)
			fmt.Printf("     %s: %q\n", jumbleOpt.desc, result)
		}
	}

	// Example 8: Multiple categories
	fmt.Println("\n8. Multiple Categories:")
	mixedText := "User john@example.com has SSN 123-45-6789 and password: secret123"
	fmt.Printf("   Original: %q\n", mixedText)

	result8, _ := schemaflow.Redact(mixedText, schemaflow.NewRedactOptions().WithCategories([]string{"PII", "secrets"}))
	fmt.Printf("   Redacted: %q\n", result8)

	// Example 9: Custom patterns
	fmt.Println("\n9. Custom Patterns:")
	customText := "API key: sk-1234567890abcdef, DB: mysql://user:pass@host/db"
	fmt.Printf("   Original: %q\n", customText)

	result9, _ := schemaflow.Redact(customText, schemaflow.NewRedactOptions().
		WithCategories([]string{"secrets"}).
		WithCustomPatterns([]string{`sk-\w{20}`, `mysql://\S+`}))
	fmt.Printf("   Redacted: %q\n", result9)

	// Example 11: Custom mask configuration
	fmt.Println("\n11. Custom Mask Configuration:")
	sampleEmail := "user@company.com"

	masks := []struct {
		desc string
		opts schemaflow.RedactOptions
	}{
		{"Default mask", schemaflow.NewRedactOptions().WithCategories([]string{"PII"})},
		{"Custom text", schemaflow.NewRedactOptions().WithCategories([]string{"PII"}).WithMaskText("[PRIVATE]")},
		{"Hash symbols", schemaflow.NewRedactOptions().WithCategories([]string{"PII"}).WithMaskChar('#').WithMaskLength(5)},
		{"At symbols", schemaflow.NewRedactOptions().WithCategories([]string{"PII"}).WithMaskChar('@').WithMaskLength(8)},
		{"Original length X", schemaflow.NewRedactOptions().WithCategories([]string{"PII"}).WithMaskChar('X').WithMaskLength(-1)},
	}

	for _, mask := range masks {
		result, _ := schemaflow.Redact(sampleEmail, mask.opts)
		fmt.Printf("   %s: %q → %q\n", mask.desc, sampleEmail, result)
	}

	// Example 12: Error handling
	fmt.Println("\n12. Error Handling:")
	fmt.Println("   Testing various error conditions...")

	// Test empty categories
	_, err1 := schemaflow.Redact("test", schemaflow.NewRedactOptions().WithCategories([]string{}))
	if err1 != nil {
		fmt.Printf("   ✓ Empty categories rejected: %v\n", err1)
	}

	// Test invalid strategy
	_, err2 := schemaflow.Redact("test", schemaflow.RedactOptions{Categories: []string{"PII"}, Strategy: "invalid"})
	if err2 != nil {
		fmt.Printf("   ✓ Invalid strategy rejected: %v\n", err2)
	}

	fmt.Println("\n=== Redact Operation Examples Complete ===")
	fmt.Println("\nNote: The Redact operation creates new objects with sensitive data")
	fmt.Println("masked, scrambled, or removed according to the specified strategy.")
	fmt.Println("The original data remains unchanged.")
}
