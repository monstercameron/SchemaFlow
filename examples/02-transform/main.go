package main

import (
	"fmt"
	"os"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// Resume represents structured resume data
type Resume struct {
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Skills     []string `json:"skills"`
	Experience []struct {
		Company  string `json:"company"`
		Position string `json:"position"`
		Years    string `json:"years"`
	} `json:"experience"`
	Education string `json:"education"`
}

// MarkdownCV represents a formatted CV in markdown
type MarkdownCV struct {
	Content string `json:"content"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		os.Exit(1)
	}

	// Input: Structured resume data
	resume := Resume{
		Name:  "Jane Developer",
		Email: "jane.dev@email.com",
		Phone: "+1-555-0123",
		Skills: []string{
			"Go", "Python", "JavaScript",
			"Docker", "Kubernetes", "AWS",
		},
		Experience: []struct {
			Company  string `json:"company"`
			Position string `json:"position"`
			Years    string `json:"years"`
		}{
			{Company: "Tech Corp", Position: "Senior Engineer", Years: "2020-2024"},
			{Company: "StartupXYZ", Position: "Full Stack Developer", Years: "2018-2020"},
		},
		Education: "BS Computer Science, MIT, 2018",
	}

	fmt.Println("📄 Transform Example - Resume to CV")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println("\n📥 Input Resume (Structured):")
	fmt.Printf("  Name:   %s\n", resume.Name)
	fmt.Printf("  Email:  %s\n", resume.Email)
	fmt.Printf("  Skills: %v\n", resume.Skills)

	// Transform: Resume → Professional Markdown CV
	cv, err := schemaflow.Transform[Resume, MarkdownCV](
		resume,
		schemaflow.NewTransformOptions().
			WithIntelligence(schemaflow.Smart).
			WithSteering("Create a professional, well-formatted CV in markdown. Use headers, bullet points, and emphasis. Make it visually appealing."),
	)

	if err != nil {
		schemaflow.GetLogger().Error("Transformation failed", "error", err)
		os.Exit(1)
	}

	// Display transformed output
	fmt.Println("\n✅ Transformed CV (Markdown):")
	fmt.Println("---")
	fmt.Println(cv.Content)
	fmt.Println("---")

	fmt.Println("\n✨ Success! Structured data → Formatted document")
}
