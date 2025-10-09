package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== SchemaFlow Explain Operation Examples ===")
	fmt.Println()

	// Note: For demonstration purposes, we'll show the API structure
	// without making actual LLM calls (requires valid API key)
	fmt.Println("📝 This example demonstrates the Explain operation API structure.")
	fmt.Println("🔑 To run with actual LLM calls, set SCHEMAFLOW_API_KEY or OPENAI_API_KEY environment variable.")
	fmt.Println()

	// Example 1: Explain for non-technical audience
	fmt.Println("1. Explanation for Non-Technical Audience:")
	fmt.Println("------------------------------------------")
	fmt.Println("Code that would be executed:")
	fmt.Println(`explanation, err := ops.Explain(sampleData,
    ops.NewExplainOptions().
        WithAudience("non-technical").
        WithDepth(2).
        WithFormat("paragraph").
        WithFocus("overview"))`)
	fmt.Println()
	fmt.Println("Mock Result:")
	fmt.Println("Summary: This data represents a user's profile and activity information")
	fmt.Println()
	fmt.Println("Full Explanation:")
	fmt.Println("This complex data structure contains information about a user named Alice Johnson. It includes her basic profile details like email and age, her preferences for using the application, a log of her recent activities, and some metadata about her account. The user appears to be an active premium member who enjoys technology and data science topics.")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  1. User has been active for over a year")
	fmt.Println("  2. Premium account with dark theme preference")
	fmt.Println("  3. Interested in technology and data science")
	fmt.Println("  4. Recent activity includes logging in and viewing analytics")
	fmt.Println()

	// Example 2: Technical explanation with implementation details
	fmt.Println("2. Technical Explanation (Implementation Focus):")
	fmt.Println("------------------------------------------------")
	fmt.Println("Code that would be executed:")
	fmt.Println(`explanation, err := ops.Explain(sampleData,
    ops.NewExplainOptions().
        WithAudience("technical").
        WithDepth(3).
        WithFormat("structured").
        WithFocus("implementation").
        WithContext("This is a user data structure from a web application"))`)
	fmt.Println()
	fmt.Println("Mock Result:")
	fmt.Println("Summary: Go struct representing user data with JSON serialization")
	fmt.Println()
	fmt.Println("Full Explanation:")
	fmt.Println("STRUCTURE OVERVIEW:")
	fmt.Println("The data is organized as a ComplexData struct with four main components:")
	fmt.Println("- UserProfile: Contains core user identification and verification fields")
	fmt.Println("- Preferences: User-configurable settings and interests")
	fmt.Println("- ActivityLog: Time-series array of user actions")
	fmt.Println("- Metadata: System-generated tracking and versioning information")
	fmt.Println()
	fmt.Println("DATA FLOW:")
	fmt.Println("1. UserProfile serves as the primary key with ID field")
	fmt.Println("2. Preferences array allows multiple interest categorization")
	fmt.Println("3. ActivityLog provides audit trail with timestamps and IP tracking")
	fmt.Println("4. Metadata enables versioning and custom field extension")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  1. Struct uses JSON tags for API serialization")
	fmt.Println("  2. ActivityLog implements time-series data pattern")
	fmt.Println("  3. Metadata map allows dynamic field extension")
	fmt.Println("  4. UserProfile includes verification and role-based access")
	fmt.Println()

	// Example 3: Executive summary in bullet points
	fmt.Println("3. Executive Summary (Bullet Points):")
	fmt.Println("-------------------------------------")
	fmt.Println("Code that would be executed:")
	fmt.Println(`explanation, err := ops.Explain(sampleData,
    ops.NewExplainOptions().
        WithAudience("executive").
        WithDepth(2).
        WithFormat("bullet-points").
        WithFocus("benefits").
        WithContext("Business intelligence and user analytics data"))`)
	fmt.Println()
	fmt.Println("Mock Result:")
	fmt.Println("Summary: Premium user profile with strong engagement metrics")
	fmt.Println()
	fmt.Println("Full Explanation:")
	fmt.Println("• High-value customer with premium subscription")
	fmt.Println("• Active user with recent login and feature usage")
	fmt.Println("• Technology-focused interests align with product roadmap")
	fmt.Println("• Complete profile data enables personalized marketing")
	fmt.Println("• Audit trail supports customer success initiatives")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  • 15+ months customer lifetime value")
	fmt.Println("  • Multiple product feature adoption")
	fmt.Println("  • Data completeness supports analytics")
	fmt.Println("  • Engagement patterns indicate loyalty")
	fmt.Println()

	// Example 4: Simple explanation for beginners
	fmt.Println("4. Beginner-Friendly Explanation:")
	fmt.Println("----------------------------------")
	fmt.Println("Code that would be executed:")
	fmt.Println(`explanation, err := ops.Explain(sampleData.UserProfile,
    ops.NewExplainOptions().
        WithAudience("beginner").
        WithDepth(1).
        WithFormat("step-by-step").
        WithFocus("overview"))`)
	fmt.Println()
	fmt.Println("Mock Result:")
	fmt.Println("Summary: Basic information about a person using our service")
	fmt.Println()
	fmt.Println("Full Explanation:")
	fmt.Println("Let me explain this user information step by step:")
	fmt.Println("1. This is data about a person named Alice Johnson")
	fmt.Println("2. She has an ID number to identify her in our system")
	fmt.Println("3. Her email address is alice.johnson@example.com")
	fmt.Println("4. She is 28 years old and joined in January 2023")
	fmt.Println("5. Her account is verified and she has premium access")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  1. Name and contact information")
	fmt.Println("  2. Age and join date")
	fmt.Println("  3. Account status and permissions")
	fmt.Println()

	// Example 5: Simple data structure explanation
	fmt.Println("5. Simple Data Structure Explanation:")
	fmt.Println("--------------------------------------")
	fmt.Println("Code that would be executed:")
	fmt.Println(`simpleData := map[string]interface{}{
    "product": "Laptop Computer",
    "price": 1299.99,
    "specs": map[string]interface{}{
        "cpu": "Intel i7",
        "ram": "16GB",
        "storage": "512GB SSD",
    },
    "in_stock": true,
}

explanation, err := ops.Explain(simpleData,
    ops.NewExplainOptions().
        WithAudience("children").
        WithDepth(2).
        WithFormat("paragraph").
        WithContext("This is information about a computer for sale"))`)
	fmt.Println()
	fmt.Println("Mock Result:")
	fmt.Println("Summary: Information about a cool computer you can buy")
	fmt.Println()
	fmt.Println("Full Explanation:")
	fmt.Println("Imagine a super smart computer that you can carry around! This laptop costs $1299.99 and has really good parts inside. It has a fast brain called an Intel i7, 16GB of memory to remember lots of things, and 512GB of super quick storage. The best part is that it's ready to buy right now - it's in stock and waiting for you!")
	fmt.Println()
	fmt.Println("Key Points:")
	fmt.Println("  1. Laptop computer for carrying around")
	fmt.Println("  2. Costs about thirteen hundred dollars")
	fmt.Println("  3. Has fast parts for good performance")
	fmt.Println("  4. Available to buy right now")
	fmt.Println()

	// Display metadata for the last explanation
	fmt.Println("6. Explanation Metadata:")
	fmt.Println("------------------------")
	fmt.Println("Audience: children")
	fmt.Println("Complexity: intermediate")
	fmt.Println("Metadata:")
	fmt.Println("  data_type: map[string]interface {}")
	fmt.Println("  field_count: 4")
	fmt.Println("  estimated_complexity: medium")
	fmt.Println("  explanation_depth: 2")
	fmt.Println("  focus_area: overview")
}