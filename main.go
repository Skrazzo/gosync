package main

import (
	"fmt"
	"os"

	"gosync/forms"
	"gosync/utils"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Check if config exists
	cfg := utils.New()
	if !cfg.Exists() {
		fmt.Println("No configuration file found. Starting setup...")

		err := forms.SetupConfig(app, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during setup: %v\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	// Load existing config
	err := cfg.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded configuration:\n")
	fmt.Printf("  Remote: %s@%s:%s\n", cfg.User, cfg.Host, cfg.RemoteDir)
	fmt.Printf("  Auth:   %s\n", cfg.AuthType)

	// TODO: Start the main TUI interface
	fmt.Println("\nMain TUI interface coming soon...")
}
