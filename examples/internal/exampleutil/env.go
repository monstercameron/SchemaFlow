package exampleutil

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	schemaflow "github.com/monstercameron/schemaflow"
)

// Bootstrap initializes examples from the process environment and an optional nearby .env file.
// Missing .env files are ignored so examples still work in CI or when environment variables are already set.
func Bootstrap() error {
	loadNearestEnv()
	normalizeAPIKeyEnv()
	return schemaflow.InitWithEnv()
}

func normalizeAPIKeyEnv() {
	if os.Getenv("SCHEMAFLOW_API_KEY") == "" {
		if key := os.Getenv("OPENAI_API_KEY"); key != "" {
			_ = os.Setenv("SCHEMAFLOW_API_KEY", key)
		}
	}
}

func loadNearestEnv() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(envPath); statErr == nil {
			_ = godotenv.Load(envPath)
			return
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}
