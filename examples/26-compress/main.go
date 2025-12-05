// Example 26: Compress Operation
// Compresses text/data while preserving essential meaning using LLM

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	schemaflow "github.com/monstercameron/SchemaFlow"
	"github.com/monstercameron/SchemaFlow/internal/ops"
	"github.com/monstercameron/SchemaFlow/internal/types"
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

	// Ensure environment is configured
	if os.Getenv("SCHEMAFLOW_API_KEY") == "" {
		log.Fatal("SCHEMAFLOW_API_KEY environment variable not set")
	}

	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		log.Fatalf("Failed to initialize SchemaFlow: %v", err)
	}

	fmt.Println("=== Compress Example ===")
	fmt.Println("Compresses text while preserving essential meaning")
	fmt.Println()

	// Example 1: Compress a document struct
	fmt.Println("--- Example 1: Document Compression ---")

	type Document struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	doc := Document{
		Title:  "AI in Healthcare",
		Author: "Dr. Smith",
		Content: `Artificial intelligence has transformed the technology landscape in unprecedented ways.
Machine learning algorithms now power everything from search engines to autonomous vehicles.
Deep learning uses neural networks with multiple layers to process complex patterns in data.
The healthcare industry has particularly benefited from AI advancements. Diagnostic systems
can now detect diseases from medical images with accuracy rivaling human experts. Drug
discovery pipelines have been accelerated through AI-powered molecular analysis.
However, the rapid advancement of AI also raises important ethical considerations around
bias in training data, privacy concerns, and potential worker displacement.`,
	}

	fmt.Println("INPUT: Document struct")
	fmt.Printf("  Title: %q\n", doc.Title)
	fmt.Printf("  Author: %q\n", doc.Author)
	fmt.Printf("  Content: (%d chars) %q...\n\n", len(doc.Content), doc.Content[:80])

	opts := ops.NewCompressOptions().
		WithCompressionRatio(0.3).
		WithStrategy("semantic").
		WithPreserveInfo([]string{"key facts", "conclusions"}).
		WithIntelligence(types.Smart)

	result, err := ops.Compress(doc, opts)
	if err != nil {
		log.Fatalf("Compression failed: %v", err)
	}

	fmt.Println("OUTPUT: CompressResult[Document]")
	fmt.Printf("  OriginalSize: %d chars\n", result.OriginalSize)
	fmt.Printf("  CompressedSize: %d chars\n", result.CompressedSize)
	fmt.Printf("  ActualRatio: %.2f\n", result.ActualRatio)
	fmt.Printf("  Compressed: %q\n\n", result.Compressed)

	// Example 2: Compress meeting notes (as string)
	fmt.Println("--- Example 2: Meeting Notes Compression ---")

	meeting := `Project Alpha Update
Attendees: John (PM), Sarah (Dev Lead), Mike (Designer), Lisa (QA)

Status Updates:
- Backend API is 80% complete, on track for Friday delivery
- Frontend has blockers with authentication flow
- Design mockups for v2 features are ready for review
- QA found 3 critical bugs in payment module

Action Items:
1. Sarah to document auth blockers by EOD
2. John to schedule sync with Alex tomorrow
3. Mike to present design mockups Wednesday
4. Lisa to prioritize payment bugs

Next Meeting: Friday 10 AM`

	fmt.Println("INPUT: string")
	fmt.Printf("  Content: (%d chars)\n", len(meeting))
	fmt.Printf("  %q...\n\n", meeting[:100])

	meetingOpts := ops.NewCompressOptions().
		WithCompressionRatio(0.4).
		WithStrategy("lossy").
		WithPriority("actions").
		WithIntelligence(types.Smart)

	meetingResult, err := ops.Compress(meeting, meetingOpts)
	if err != nil {
		log.Fatalf("Meeting compression failed: %v", err)
	}

	fmt.Println("OUTPUT: CompressResult[string]")
	fmt.Printf("  OriginalSize: %d chars\n", meetingResult.OriginalSize)
	fmt.Printf("  CompressedSize: %d chars\n", meetingResult.CompressedSize)
	fmt.Printf("  ActualRatio: %.2f\n", meetingResult.ActualRatio)
	fmt.Printf("  Compressed: %q\n", meetingResult.Compressed)

	fmt.Println("\n=== Compress Example Complete ===")
}
