package main

import (
	"fmt"

	"github.com/monstercameron/SchemaFlow/ops"
)

func main() {
	fmt.Println("=== SchemaFlow Complete Operation Examples ===\n")

	// Example 1: Basic text completion
	fmt.Println("1. Basic Text Completion:")
	partial1 := "The weather today is"
	fmt.Printf("   Input: %q\n", partial1)

	// Note: In a real scenario, this would call the LLM
	// For demo purposes, we'll show the structure
	opts1 := ops.NewCompleteOptions().WithMaxLength(50)
	fmt.Printf("   Options: MaxLength=%d, Temperature=%.1f\n", opts1.MaxLength, opts1.Temperature)
	fmt.Printf("   Mock Completion: %q\n", "The weather today is beautiful and sunny with clear blue skies. Perfect for outdoor activities!")

	// Example 2: Completion with context
	fmt.Println("\n2. Completion with Context:")
	partial2 := "Please send me the"
	context2 := []string{
		"User: I need help with my order",
		"Assistant: I'd be happy to help! What seems to be the issue?",
		"User: I haven't received my package yet",
	}
	fmt.Printf("   Input: %q\n", partial2)
	fmt.Printf("   Context: %d messages\n", len(context2))
	for i, msg := range context2 {
		fmt.Printf("     %d. %s\n", i+1, msg)
	}

	opts2 := ops.NewCompleteOptions().
		WithContext(context2).
		WithMaxLength(100).
		WithTemperature(0.8)
	fmt.Printf("   Options: MaxLength=%d, Temperature=%.1f\n", opts2.MaxLength, opts2.Temperature)
	fmt.Printf("   Mock Completion: %q\n", "Please send me the tracking information for my order. I placed it last week and expected delivery by now.")

	// Example 3: Code completion with stop sequences
	fmt.Println("\n3. Code Completion with Stop Sequences:")
	partial3 := "function calculateTotal(items) {"
	stopSeq3 := []string{"}", "\n\n"}
	fmt.Printf("   Input: %q\n", partial3)
	fmt.Printf("   Stop Sequences: %v\n", stopSeq3)

	opts3 := ops.NewCompleteOptions().
		WithStopSequences(stopSeq3).
		WithMaxLength(200).
		WithTemperature(0.3) // Lower temperature for code
	fmt.Printf("   Options: MaxLength=%d, Temperature=%.1f\n", opts3.MaxLength, opts3.Temperature)
	fmt.Printf("   Mock Completion: %q\n", `function calculateTotal(items) {
  let total = 0;
  for (let item of items) {
    total += item.price * item.quantity;
  }
  return total;
}`)

	// Example 4: Creative writing completion
	fmt.Println("\n4. Creative Writing Completion:")
	partial4 := "Once upon a time, in a land far away,"
	fmt.Printf("   Input: %q\n", partial4)

	opts4 := ops.NewCompleteOptions().
		WithMaxLength(150).
		WithTemperature(1.2). // Higher temperature for creativity
		WithTopP(0.95)
	fmt.Printf("   Options: MaxLength=%d, Temperature=%.1f, TopP=%.2f\n",
		opts4.MaxLength, opts4.Temperature, opts4.TopP)
	fmt.Printf("   Mock Completion: %q\n", "Once upon a time, in a land far away, there lived a magical dragon who could speak every language in the world. His scales shimmered like diamonds in the sunlight, and his eyes glowed with ancient wisdom passed down through generations of mythical creatures.")

	// Example 5: Email completion
	fmt.Println("\n5. Email Completion:")
	partial5 := "Dear team,\n\nI hope this email finds you well."
	context5 := []string{
		"Subject: Project Update Meeting",
		"Previous email: Let's schedule a meeting to discuss the project progress",
	}
	fmt.Printf("   Input: %q\n", partial5)
	fmt.Printf("   Context: %v\n", context5)

	opts5 := ops.NewCompleteOptions().
		WithContext(context5).
		WithMaxLength(300).
		WithStopSequences([]string{"Best regards", "Sincerely"})
	fmt.Printf("   Options: MaxLength=%d\n", opts5.MaxLength)
	fmt.Printf("   Mock Completion: %q\n", `Dear team,

I hope this email finds you well. I wanted to follow up on our previous discussion about scheduling a project update meeting. Based on everyone's availability, I propose we meet this Friday at 2 PM in the main conference room. We'll review the current progress, discuss any blockers, and plan the next sprint objectives.

Please let me know if this time works for everyone, or if you'd prefer an alternative time.

Best regards,
[Your Name]`)

	// Example 6: Sentence completion with different modes
	fmt.Println("\n6. Sentence Completion (Different Modes):")
	partial6 := "The new feature will"
	fmt.Printf("   Input: %q\n", partial6)

	modes := []struct {
		name string
		opts ops.CompleteOptions
		mock string
	}{
		{"Conservative", ops.NewCompleteOptions().WithTemperature(0.1).WithMaxLength(80), "The new feature will improve user experience by providing faster load times and better error handling."},
		{"Balanced", ops.NewCompleteOptions().WithTemperature(0.7).WithMaxLength(80), "The new feature will help users manage their tasks more efficiently with an intuitive interface."},
		{"Creative", ops.NewCompleteOptions().WithTemperature(1.5).WithMaxLength(80), "The new feature will revolutionize how users interact with data through magical visualization portals."},
	}

	for _, mode := range modes {
		fmt.Printf("   %s Mode: Temperature=%.1f\n", mode.name, mode.opts.Temperature)
		fmt.Printf("     Mock Completion: %q\n", mode.mock)
	}

	// Example 7: Completion with length limits
	fmt.Println("\n7. Completion with Length Constraints:")
	partial7 := "In conclusion,"
	fmt.Printf("   Input: %q\n", partial7)

	opts7 := ops.NewCompleteOptions().
		WithMaxLength(50). // Very short completion
		WithStopSequences([]string{".", "!"})
	fmt.Printf("   Options: MaxLength=%d, StopSequences=%v\n",
		opts7.MaxLength, opts7.StopSequences)
	fmt.Printf("   Mock Completion: %q\n", "In conclusion, this approach provides the best balance of performance and maintainability.")

	// Example 8: API documentation completion
	fmt.Println("\n8. API Documentation Completion:")
	partial8 := "// GET /api/users - Retrieve a list of users"
	fmt.Printf("   Input: %q\n", partial8)

	opts8 := ops.NewCompleteOptions().
		WithMaxLength(250).
		WithTemperature(0.4). // Technical writing
		WithStopSequences([]string{"\n\n", "// POST"})
	fmt.Printf("   Options: MaxLength=%d, Temperature=%.1f\n",
		opts8.MaxLength, opts8.Temperature)
	fmt.Printf("   Mock Completion: %q\n", `// GET /api/users - Retrieve a list of users
//
// Query Parameters:
// - limit (optional): Maximum number of users to return (default: 20, max: 100)
// - offset (optional): Number of users to skip for pagination (default: 0)
// - search (optional): Search term to filter users by name or email
//
// Response: 200 OK
// Returns a JSON array of user objects with id, name, email, and created_at fields
//
// Example: GET /api/users?limit=10&search=john`)

	// Example 9: Chat message completion
	fmt.Println("\n9. Chat Message Completion:")
	partial9 := "That sounds like a great idea! I think we should"
	context9 := []string{
		"Alice: Hey, what do you think about organizing a team outing?",
		"Bob: That sounds fun! What kind of activities are you thinking?",
		"Alice: Maybe hiking or a picnic in the park",
	}
	fmt.Printf("   Input: %q\n", partial9)
	fmt.Printf("   Context: %d messages\n", len(context9))

	opts9 := ops.NewCompleteOptions().
		WithContext(context9).
		WithMaxLength(120).
		WithTemperature(0.9)
	fmt.Printf("   Options: MaxLength=%d, Temperature=%.1f\n",
		opts9.MaxLength, opts9.Temperature)
	fmt.Printf("   Mock Completion: %q\n", "That sounds like a great idea! I think we should plan a hiking trip to the nearby mountains. The trails there are perfect for our group size and skill level.")

	// Example 10: Error handling demonstration
	fmt.Println("\n10. Error Handling:")
	fmt.Println("   Testing various error conditions...")

	// Test empty input
	_, err1 := ops.Complete("", ops.NewCompleteOptions())
	if err1 != nil {
		fmt.Printf("   ✓ Empty input rejected: %v\n", err1)
	}

	// Test invalid options
	invalidOpts := ops.NewCompleteOptions().WithMaxLength(-1)
	_, err2 := ops.Complete("test", invalidOpts)
	if err2 != nil {
		fmt.Printf("   ✓ Invalid options rejected: %v\n", err2)
	}

	fmt.Println("\n=== Complete Operation Examples Complete ===")
	fmt.Println("\nNote: These examples show the API structure. In a real implementation,")
	fmt.Println("the Complete function would call an LLM to generate intelligent completions")
	fmt.Println("based on the partial text and context provided.")
}
