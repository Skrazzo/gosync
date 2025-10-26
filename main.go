package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gosync/forms"
	"gosync/utils"

	"github.com/rivo/tview"
)

func main() {
	// Parse command-line flags
	consoleMode := flag.Bool("console", false, "Run in console mode (no dashboard, just file watcher logs)")
	flag.Parse()

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

	// Connect to remote server
	sftp := utils.NewSftp()
	// if err := sftp.Connect(cfg); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error connecting to remote server: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer sftp.Close()

	// fmt.Printf("Connected to %s@%s\n", cfg.User, cfg.Host)

	// Run in console mode or dashboard mode
	if *consoleMode {
		fmt.Println("\n=== Console Mode ===")
		fmt.Println("Running file watcher in console mode. Press Ctrl+C to exit.\n")

		// Go routine to start file watcher
		go func() {
			if err := utils.StartFileWatcher(cfg, sftp); err != nil {
				fmt.Fprintf(os.Stderr, "File watcher error: %v\n", err)
				os.Exit(1)
			}
		}()

		for {
			utils.ClearScreen()

			fmt.Println("Upload queue:")

			for _, file := range sftp.Queue.Uploads {
				fmt.Printf("  %s\n", file)
			}

			fmt.Println("\nDelete queue:")

			for _, file := range sftp.Queue.Deletes {
				fmt.Printf("  %s\n", file)
			}

			time.Sleep(time.Second)
		}
	} else {
		ShowView()
	}

}
