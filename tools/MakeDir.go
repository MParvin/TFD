package tools

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// EnsureDirectories reads all directory paths from Viper config and creates them
func EnsureDirectories() error {
	// Define all directory keys to check
	dirKeys := []string{
		"directories.photo",
		"directories.video",
		"directories.music",
		"directories.voice",
		"directories.document",
		"directories.other",
	}

	for _, key := range dirKeys {
		dirPath := viper.GetString(key)
		if dirPath == "" {
			return fmt.Errorf("directory path for %s is not configured", key)
		}

		// Validate path is absolute
		if !filepath.IsAbs(dirPath) {
			return fmt.Errorf("directory path for %s must be absolute: %s", key, dirPath)
		}

		// Create directory if it doesn't exist
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}

		log.Printf("Ensured directory exists: %s", dirPath)
	}

	return nil
}
