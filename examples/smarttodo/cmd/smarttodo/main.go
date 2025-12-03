package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/monstercameron/schemaflow"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/database"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/localization"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/models"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/tui"
)

func main() {
	// Disable schemaflow debug logging to prevent TUI corruption
	os.Setenv("SCHEMAFLOW_DEBUG", "false")

	// Check if API key exists before initializing
	needsAPIKey := false
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("SCHEMAFLOW_API_KEY")
	}

	if apiKey == "" {
		// No API key found, will need to set up
		needsAPIKey = true
		// Set a temporary key to allow the app to start
		os.Setenv("OPENAI_API_KEY", "sk-temp-key-will-be-replaced")
	}

	// Initialize SchemaFlow with environment variables
	if err := schemaflow.InitWithEnv(); err != nil && !needsAPIKey {
		fmt.Printf("❌ Error: %v\n", err)
		fmt.Println("Please ensure your .env file contains:")
		fmt.Println("  OPENAI_API_KEY=your-api-key")
		fmt.Println("Or set the environment variable:")
		fmt.Println("  export SCHEMAFLOW_API_KEY='your-api-key'")
		os.Exit(1)
	}

	// Ensure debug mode is off
	schemaflow.SetDebugMode(false)

	// Initialize localization system
	localization.InitLocalization()

	// Preload common strings for better performance
	go localization.PreloadCommonStrings()

	// Redirect log output to file to prevent TUI corruption
	// Note: Using structured logging instead of standard log
	logFile, err := os.OpenFile("/tmp/smarttodo.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		defer logFile.Close()
	}

	// Get database path
	dbPath := os.Getenv("SMARTTODO_DB")
	if dbPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			schemaflow.GetLogger().Error("Failed to get home directory", "error", err)
			os.Exit(1)
		}
		dbPath = filepath.Join(home, ".smarttodo.db")
	}

	// Initialize database
	db, err := database.NewDatabase(dbPath)
	if err != nil {
		schemaflow.GetLogger().Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Create the TUI program with API key flag
	model := tui.InitialModel(db)
	model.SetNeedsAPIKey(needsAPIKey)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Handle signals in a goroutine
	go func() {
		sig := <-sigChan
		schemaflow.GetLogger().Info("Received signal, initiating graceful shutdown", "signal", sig)

		// Send a message to start the closing animation
		p.Send(models.StartClosingMsg{})

		// After a short delay, force quit if animation doesn't complete
		go func() {
			time.Sleep(3 * time.Second)
			p.Send(tea.Quit())
		}()
	}()

	// Run the program
	if _, err := p.Run(); err != nil {
		schemaflow.GetLogger().Error("Error running program", "error", err)
		// Ensure database is closed
		if db != nil {
			db.Close()
		}
		os.Exit(1)
	}

	// Program exited normally
	schemaflow.GetLogger().Info("Smart Todo closed successfully")
	os.Exit(0)
}


