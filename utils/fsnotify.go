package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Function checks if string includes any of the patterns, and returns true if it does
func includesPattern(str string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(str, pattern) {
			return true
		}
	}
	return false
}

// addDirRecursively adds a directory and all its subdirectories to the watcher
func addDirRecursively(watcher *fsnotify.Watcher, path string, ignorePatterns []string) error {
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if includesPattern(walkPath, ignorePatterns) {
			return nil
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

func StartFileWatcher(cfg *Config, sftp *SFTP) error {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating new watcher: %v", err)
	}
	defer watcher.Close()

	// Add current directory and all subdirectories recursively
	fmt.Println("Setting up file watcher...")
	err = addDirRecursively(watcher, ".", cfg.Ignore)
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
			// fmt.Printf("[EVENT] %s\n", event)

			// Log specific event types
			if event.Has(fsnotify.Write) {
				// Add to upload queue
				sftp.Queue.Uploads = append(sftp.Queue.Uploads, event.Name)
			}

			if event.Has(fsnotify.Create) {
				// Add to upload queue
				sftp.Queue.Uploads = append(sftp.Queue.Uploads, event.Name)

				// If a new directory was created, watch it too
				fileInfo, err := os.Stat(event.Name)
				if err == nil && fileInfo.IsDir() {
					// TODO: Add to upload queue
					fmt.Printf("  ├─ Info: New directory detected, adding to watcher\n")
					if err := addDirRecursively(watcher, event.Name, cfg.Ignore); err != nil {
						// TODO: Add to error list
						fmt.Printf("  ├─ Warning: Could not add new directory to watcher: %v\n", err)
					}
				}
			}

			if event.Has(fsnotify.Remove) {
				// Add to delete queue
				sftp.Queue.Deletes = append(sftp.Queue.Deletes, event.Name)
			}

			if event.Has(fsnotify.Rename) {
				// Add to delete queue
				sftp.Queue.Deletes = append(sftp.Queue.Deletes, event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Printf("[ERROR] %v\n", err)
		}
	}
}
