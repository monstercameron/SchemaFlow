package main

import (
	"fmt"
	"os"
	"time"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// Email represents a structured email
type Email struct {
	From    string    `json:"from"`
	To      string    `json:"to"`
	Subject string    `json:"subject"`
	Date    time.Time `json:"date"`
	Body    string    `json:"body"`
	Tags    []string  `json:"tags"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		os.Exit(1)
	}

	// Raw email text (unstructured)
	rawEmail := `
From: john.smith@example.com
To: sarah.jones@company.com
Subject: Project Update - Q4 Results

Hi Sarah,

I wanted to share the Q4 results with you. The project exceeded expectations 
with a 25% increase in productivity. The team did an excellent job meeting 
all deadlines.

Let me know if you need any additional details.

Best regards,
John Smith
Sent: December 15, 2024
`

	fmt.Println("📧 Email Parser Example")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println("\n📥 Raw Input:")
	fmt.Println(rawEmail)

	// Extract structured email from unstructured text
	email, err := schemaflow.Extract[Email](rawEmail, schemaflow.NewExtractOptions().
		WithIntelligence(schemaflow.Fast).
		WithSteering("Extract all email fields including metadata and categorize by tags"))

	if err != nil {
		schemaflow.GetLogger().Error("Extraction failed", "error", err)
		os.Exit(1)
	}

	// Display structured output
	fmt.Println("\n✅ Extracted Email:")
	fmt.Printf("  From:    %s\n", email.From)
	fmt.Printf("  To:      %s\n", email.To)
	fmt.Printf("  Subject: %s\n", email.Subject)
	fmt.Printf("  Date:    %s\n", email.Date.Format("2006-01-02"))
	fmt.Printf("  Tags:    %v\n", email.Tags)
	fmt.Printf("\n  Body:\n  %s\n", email.Body)

	fmt.Println("\n✨ Success! Unstructured text → Structured data")
}
