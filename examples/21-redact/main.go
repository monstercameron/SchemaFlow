// 21-redact: LLM-powered sensitive data detection and character-level masking
// Intelligence: Fast (Cerebras gpt-oss-120b)
// Expectations:
// - LLM identifies sensitive data (emails, phones, SSNs, names, etc.)
// - Returns exact character positions (spans) for each detection
// - Masks at character level with configurable mask char
// - Supports partial reveal (show first N, last N characters)

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

func main() {
	loadEnv()

	fmt.Println("🔒 RedactLLM Example - AI-Powered Sensitive Data Masking")
	fmt.Println("=" + string(make([]byte, 60)))

	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		return
	}

	// Example 1: Basic sensitive data detection
	fmt.Println("\n1️⃣  Basic PII Detection")
	fmt.Println("─" + string(make([]byte, 60)))

	text1 := "Contact John Smith at john.smith@company.com or call 555-123-4567 for support."
	fmt.Printf("   Input:  %q\n", text1)

	result1, err := schemaflow.RedactLLM(text1,
		schemaflow.NewRedactLLMOptions().
			WithCategories([]string{"all"}).
			WithMaskChar('*').
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   Output: %q\n", result1.Text)
		fmt.Println("   Spans detected:")
		for _, span := range result1.Spans {
			fmt.Printf("      [%d:%d] %s → %q\n", span.Start, span.End, span.Category, span.Original)
		}
	}

	// Example 2: Partial reveal (show first/last characters)
	fmt.Println("\n2️⃣  Partial Reveal - Show First 2 and Last 2 Characters")
	fmt.Println("─" + string(make([]byte, 60)))

	text2 := "My email is alice.johnson@example.org and SSN is 123-45-6789."
	fmt.Printf("   Input:  %q\n", text2)

	result2, err := schemaflow.RedactLLM(text2,
		schemaflow.NewRedactLLMOptions().
			WithCategories([]string{"email", "ssn"}).
			WithMaskChar('*').
			WithShowFirst(2).
			WithShowLast(2).
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   Output: %q\n", result2.Text)
		fmt.Println("   Spans detected:")
		for _, span := range result2.Spans {
			fmt.Printf("      [%d:%d] %s → %q\n", span.Start, span.End, span.Category, span.Original)
		}
	}

	// Example 3: Custom mask character
	fmt.Println("\n3️⃣  Custom Mask Character (█)")
	fmt.Println("─" + string(make([]byte, 60)))

	text3 := "Credit card: 4532-1234-5678-9012, expires 12/25."
	fmt.Printf("   Input:  %q\n", text3)

	result3, err := schemaflow.RedactLLM(text3,
		schemaflow.NewRedactLLMOptions().
			WithCategories([]string{"credit_card"}).
			WithMaskChar('█').
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   Output: %q\n", result3.Text)
		fmt.Println("   Spans detected:")
		for _, span := range result3.Spans {
			fmt.Printf("      [%d:%d] %s → %q\n", span.Start, span.End, span.Category, span.Original)
		}
	}

	// Example 4: Detect secrets and API keys
	fmt.Println("\n4️⃣  Detect Secrets and API Keys")
	fmt.Println("─" + string(make([]byte, 60)))

	text4 := "Set API_KEY=sk-abc123xyz789secret and DB_PASSWORD=MyS3cretP@ss!"
	fmt.Printf("   Input:  %q\n", text4)

	result4, err := schemaflow.RedactLLM(text4,
		schemaflow.NewRedactLLMOptions().
			WithCategories([]string{"api_key", "password"}).
			WithMaskChar('#').
			WithShowFirst(3). // Show "sk-" prefix
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   Output: %q\n", result4.Text)
		fmt.Println("   Spans detected:")
		for _, span := range result4.Spans {
			fmt.Printf("      [%d:%d] %s → %q\n", span.Start, span.End, span.Category, span.Original)
		}
	}

	// Example 5: Mixed content with multiple categories
	fmt.Println("\n5️⃣  Complex Text with Multiple PII Types")
	fmt.Println("─" + string(make([]byte, 60)))

	text5 := `Customer: Bob Wilson
Email: bob.wilson@gmail.com
Phone: (555) 987-6543
Address: 123 Main Street, New York, NY 10001
DOB: 03/15/1985
SSN: 987-65-4321`

	fmt.Printf("   Input:\n%s\n\n", text5)

	result5, err := schemaflow.RedactLLM(text5,
		schemaflow.NewRedactLLMOptions().
			WithCategories([]string{"all"}).
			WithMaskChar('X').
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   Output:\n%s\n", result5.Text)
		fmt.Printf("\n   Categories found: %v\n", result5.Categories)
		fmt.Printf("   Total spans: %d\n", len(result5.Spans))
	}

	// Example 6: Specific category only
	fmt.Println("\n6️⃣  Detect Only Email Addresses")
	fmt.Println("─" + string(make([]byte, 60)))

	text6 := "Contacts: admin@company.com, support@help.org, John Smith (555-1234)"
	fmt.Printf("   Input:  %q\n", text6)

	result6, err := schemaflow.RedactLLM(text6,
		schemaflow.NewRedactLLMOptions().
			WithCategories([]string{"email"}). // Only emails, ignore phone
			WithMaskChar('*').
			WithIntelligence(schemaflow.Fast))

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   Output: %q\n", result6.Text)
		fmt.Println("   Spans detected:")
		for _, span := range result6.Spans {
			fmt.Printf("      [%d:%d] %s → %q\n", span.Start, span.End, span.Category, span.Original)
		}
	}

	fmt.Println()
	fmt.Println("✨ Success! Sensitive data detected and masked with LLM intelligence")
}
