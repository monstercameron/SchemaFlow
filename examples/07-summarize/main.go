package main

import (
	"fmt"
	"os"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		os.Exit(1)
	}

	// Long article text
	article := `
Artificial Intelligence Revolution in Healthcare

The healthcare industry is experiencing a profound transformation driven by artificial 
intelligence technologies. Machine learning algorithms are now capable of analyzing 
medical images with accuracy that rivals or exceeds human radiologists. In a recent 
study published in Nature Medicine, an AI system detected breast cancer in mammograms 
with 94.5% accuracy, compared to 88% for human radiologists.

Beyond diagnostics, AI is revolutionizing drug discovery. Traditional pharmaceutical 
research takes 10-15 years and costs over $2 billion per drug. AI-powered platforms 
can now screen millions of molecular combinations in weeks, identifying promising 
candidates for further testing. Atomwise, a San Francisco-based startup, used AI to 
discover two potential Ebola treatments in just one day - a process that would have 
taken years using conventional methods.

Patient care is also being transformed through AI-powered virtual health assistants 
and predictive analytics. These systems can monitor patient vital signs in real-time, 
predict potential health emergencies before they occur, and provide personalized 
treatment recommendations based on individual genetic profiles and medical histories.

However, challenges remain. Data privacy concerns, algorithmic bias, and the need for 
regulatory frameworks are critical issues that must be addressed. The FDA has approved 
only a handful of AI-based medical devices, and questions about liability when AI 
systems make errors remain unresolved.

Despite these challenges, experts predict that AI will become an integral part of 
healthcare within the next decade, potentially saving millions of lives and reducing 
healthcare costs by up to 30%. The key will be ensuring that these technologies are 
deployed ethically and equitably, benefiting all patients regardless of their 
socioeconomic status.
`

	fmt.Println("📰 Summarize Example - Article Condensation")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Printf("\n📄 Original Article Length: %d characters\n", len(article))
	fmt.Println("\n📥 Original Article:")
	fmt.Println("---")
	fmt.Println(article)
	fmt.Println("---")

	// Example 1: Simple string summary (original API)
	fmt.Println("\n🔹 Example 1: Simple Summary (string → string)")
	fmt.Println("-" + string(make([]byte, 40)))

	summaryOpts := schemaflow.NewSummarizeOptions()
	summaryOpts.TargetLength = 3 // 3 sentences
	summaryOpts.LengthUnit = "sentences"
	summaryOpts.OpOptions.Intelligence = schemaflow.Fast
	summaryOpts.OpOptions.Steering = "Create a concise summary capturing key points: AI in diagnostics, drug discovery, patient care, and challenges."

	summary, err := schemaflow.Summarize(article, summaryOpts)
	if err != nil {
		schemaflow.GetLogger().Error("Summarization failed", "error", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Summary:")
	fmt.Println("---")
	fmt.Println(summary)
	fmt.Println("---")

	// Example 2: Summary with metadata (new API)
	fmt.Println("\n🔹 Example 2: Summary with Metadata (WithMetadata API)")
	fmt.Println("-" + string(make([]byte, 40)))

	metadataOpts := schemaflow.NewSummarizeOptions()
	metadataOpts.TargetLength = 3
	metadataOpts.LengthUnit = "sentences"
	metadataOpts.OpOptions.Intelligence = schemaflow.Fast

	result, err := schemaflow.SummarizeWithMetadata(article, metadataOpts)
	if err != nil {
		schemaflow.GetLogger().Error("SummarizeWithMetadata failed", "error", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Summary with Metadata:")
	fmt.Println("---")
	fmt.Println(result.Text)
	fmt.Println("---")

	fmt.Printf("\n📊 Rich Metadata:\n")
	fmt.Printf("   Compression Ratio: %.1f%% of original\n", result.CompressionRatio*100)
	fmt.Printf("   Confidence:        %.0f%%\n", result.Confidence*100)

	if len(result.KeyPoints) > 0 {
		fmt.Println("\n📌 Key Points Extracted:")
		for i, point := range result.KeyPoints {
			fmt.Printf("   %d. %s\n", i+1, point)
		}
	}

	// Show comparison
	fmt.Println("\n📈 Summary Statistics:")
	fmt.Printf("   Original:    %d characters\n", len(article))
	fmt.Printf("   Summary:     %d characters\n", len(result.Text))
	fmt.Printf("   Reduction:   %.1f%% smaller\n", (1-result.CompressionRatio)*100)

	fmt.Println("\n✨ Success! SummarizeWithMetadata provides rich insights beyond just text")
}
