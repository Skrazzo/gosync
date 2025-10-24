package main

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Check if config exists
	if !ConfigExists() {
		fmt.Println("No configuration file found. Starting setup...")

		config, err := ShowSetupForm(app)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during setup: %v\n", err)
			os.Exit(1)
		}

		if config == nil {
			fmt.Println("Setup cancelled.")
			os.Exit(0)
		}

		fmt.Printf("Configuration saved to %s\n", ConfigFileName)
		fmt.Println("Run gosync again to start synchronization.")
		os.Exit(0)
	}

	// Load existing config
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded configuration:\n")
	fmt.Printf("  Local:  %s\n", config.LocalDir)
	fmt.Printf("  Remote: %s@%s:%s\n", config.User, config.Host, config.RemoteDir)
	fmt.Printf("  Auth:   %s\n", config.AuthType)

	// TODO: Start the main TUI interface
	fmt.Println("\nMain TUI interface coming soon...")
}
