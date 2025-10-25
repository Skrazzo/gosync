package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// addDirRecursively adds a directory and all its subdirectories to the watcher
func addDirRecursively(watcher *fsnotify.Watcher, path string) error {
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if err := watcher.Add(walkPath); err != nil {
				return fmt.Errorf("error adding directory %s: %v", walkPath, err)
			}
			fmt.Printf("  Watching: %s\n", walkPath)
		}
		return nil
	})
}

func StartFileWatcher() error {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating new watcher: %v", err)
	}
	defer watcher.Close()

	// Add current directory and all subdirectories recursively
	fmt.Println("Setting up file watcher...")
	err = addDirRecursively(watcher, ".")
	if err != nil {
		return fmt.Errorf("error setting up recursive watching: %v", err)
	}

	fmt.Println("\nFile watcher started. Watching all directories recursively for changes...")
	fmt.Println("---")

	// Start listening for events (blocking).
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// Log the raw event
			fmt.Printf("[EVENT] %s\n", event)

			// Log specific event types
			if event.Has(fsnotify.Write) {
				fmt.Printf("  ├─ Type: WRITE\n")
				fmt.Printf("  └─ File: %s\n", event.Name)
			}
			if event.Has(fsnotify.Create) {
				fmt.Printf("  ├─ Type: CREATE\n")
				fmt.Printf("  └─ File: %s\n", event.Name)

				// If a new directory was created, watch it too
				fileInfo, err := os.Stat(event.Name)
				if err == nil && fileInfo.IsDir() {
					fmt.Printf("  ├─ Info: New directory detected, adding to watcher\n")
					if err := addDirRecursively(watcher, event.Name); err != nil {
						fmt.Printf("  ├─ Warning: Could not add new directory to watcher: %v\n", err)
					}
				}
			}
			if event.Has(fsnotify.Remove) {
				fmt.Printf("  ├─ Type: REMOVE\n")
				fmt.Printf("  └─ File: %s\n", event.Name)
			}
			if event.Has(fsnotify.Rename) {
				fmt.Printf("  ├─ Type: RENAME\n")
				fmt.Printf("  └─ File: %s\n", event.Name)
			}
			if event.Has(fsnotify.Chmod) {
				fmt.Printf("  ├─ Type: CHMOD\n")
				fmt.Printf("  └─ File: %s\n", event.Name)
			}
			fmt.Println("---")

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Printf("[ERROR] %v\n", err)
		}
	}
}
