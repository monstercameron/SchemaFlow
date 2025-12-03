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

	fmt.Println("=== Compress Example ===")

	// Example 1: Compress a long article
	fmt.Println("--- Example 1: Article Compression ---")
	article := `Artificial intelligence has transformed the technology landscape in unprecedented ways.
	Machine learning algorithms now power everything from search engines to autonomous vehicles.
	Deep learning, a subset of machine learning, uses neural networks with multiple layers to
	process complex patterns in data. These networks can identify objects in images, understand
	natural language, and even generate creative content.
	
	The healthcare industry has particularly benefited from AI advancements. Diagnostic systems
	can now detect diseases from medical images with accuracy rivaling human experts. Drug
	discovery pipelines have been accelerated through AI-powered molecular analysis.
	
	However, the rapid advancement of AI also raises important ethical considerations. Issues
	around bias in training data, privacy concerns, and the potential displacement of workers
	require careful attention. Researchers and policymakers are working together to develop
	frameworks for responsible AI development.
	
	Looking forward, the integration of AI into everyday life will continue to deepen. Smart
	home devices, personalized education, and intelligent transportation systems are just a
	few areas where AI will play an increasingly central role.`

	opts := ops.NewCompressOptions().
		WithCompressionRatio(0.3). // Reduce to 30% of original
		WithStrategy("extractive").
		WithPreserveInfo([]string{"key facts", "conclusions"})

	result, err := ops.Compress(article, opts)
	if err != nil {
		log.Fatalf("Compression failed: %v", err)
	}

	fmt.Printf("Original length: %d characters\n", len(article))
	fmt.Printf("Compressed length: %d characters\n", result.CompressedSize)
	fmt.Printf("Actual compression ratio: %.2f\n", result.ActualRatio)
	fmt.Printf("\nCompressed text:\n%s\n\n", result.Compressed)

	// Example 2: Abstractive summarization
	fmt.Println("--- Example 2: Abstractive Compression ---")
	report := `Q3 Financial Report Summary:
	
	Revenue increased by 15% compared to Q2, reaching $4.2 million. This growth was driven
	primarily by strong performance in the enterprise segment, which saw a 25% increase.
	The consumer segment remained flat with 2% growth.
	
	Operating expenses grew by 8%, mainly due to increased R&D investment in our new AI
	product line. Marketing expenses were reduced by 5% through optimization of digital
	channels.
	
	Net profit margin improved from 12% to 14%, demonstrating improved operational efficiency.
	Cash flow from operations was positive at $800,000.
	
	Key achievements this quarter include the launch of three new features, signing 15 new
	enterprise clients, and expanding our team by 20 engineers.
	
	Looking ahead to Q4, we expect continued growth driven by holiday season demand and
	the full rollout of our AI features.`

	abstractOpts := ops.NewCompressOptions().
		WithCompressionRatio(0.4).
		WithStrategy("abstractive").
		WithPriority("revenue and profit metrics")

	abstractResult, err := ops.Compress(report, abstractOpts)
	if err != nil {
		log.Fatalf("Abstractive compression failed: %v", err)
	}

	fmt.Printf("Original report: %d chars\n", len(report))
	fmt.Printf("Compressed summary:\n%s\n\n", abstractResult.Compressed)

	// Example 3: Compress with key points preserved
	fmt.Println("--- Example 3: Key Points Extraction ---")
	meeting := `Team Meeting Notes - Project Alpha Update

	Attendees: John (PM), Sarah (Dev Lead), Mike (Designer), Lisa (QA)
	
	Status Updates:
	- Backend API is 80% complete, on track for Friday delivery
	- Frontend has some blockers with the authentication flow
	- Design mockups for v2 features are ready for review
	- QA found 3 critical bugs in the payment module
	
	Discussion Points:
	The team discussed the upcoming demo with the client scheduled for next Thursday.
	Sarah raised concerns about the authentication blockers potentially affecting the
	demo timeline. John suggested bringing in Alex from the platform team to help.
	
	Action Items:
	1. Sarah to document auth blockers and share by EOD
	2. John to schedule sync with Alex tomorrow morning
	3. Mike to present design mockups in Wednesday design review
	4. Lisa to prioritize payment bugs for sprint
	
	Next Meeting: Friday 10 AM to review demo readiness`

	keyPointsOpts := ops.NewCompressOptions().
		WithCompressionRatio(0.25).
		WithStrategy("hybrid").
		WithPriority("action items and blockers").
		WithPreserveInfo([]string{"deadlines", "attendees"})

	keyResult, err := ops.Compress(meeting, keyPointsOpts)
	if err != nil {
		log.Fatalf("Key points compression failed: %v", err)
	}

	fmt.Printf("Meeting notes compressed from %d to %d chars:\n", len(meeting), keyResult.CompressedSize)
	fmt.Printf("%s\n", keyResult.Compressed)

	// Example 4: Simple text compression
	fmt.Println("\n--- Example 4: Simple Text Compression ---")
	simpleText := "This is a very long text that contains many unnecessary words and phrases that could be shortened significantly while maintaining the core meaning."

	compressedText, err := ops.CompressText(simpleText, ops.NewCompressOptions().WithCompressionRatio(0.5))
	if err != nil {
		log.Fatalf("Simple compression failed: %v", err)
	}

	fmt.Printf("Original: %s\n", simpleText)
	fmt.Printf("Compressed: %s\n", compressedText)

	fmt.Println("\n=== Compress Example Complete ===")
}
