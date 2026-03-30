package tools

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// IsValidURL validates whether a string is a valid URL with http or https scheme
func IsValidURL(urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// DownloadFromURL downloads a file from a URL and saves it to the destination path
func DownloadFromURL(urlStr, destPath string, client *http.Client) error {
	if !IsValidURL(urlStr) {
		return fmt.Errorf("invalid URL: %s", urlStr)
	}

	if client == nil {
		client = http.DefaultClient
	}

	// Make HTTP GET request
	resp, err := client.Get(urlStr)
	if err != nil {
		return fmt.Errorf("failed to download from %s: %w", urlStr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d when downloading from %s", resp.StatusCode, urlStr)
	}

	// Extract filename from URL and sanitize it
	u, _ := url.Parse(urlStr)
	filename := filepath.Base(u.Path)
	if filename == "" || filename == "." {
		filename = "downloaded_file"
	}

	// Sanitize filename to prevent path traversal attacks
	filename = sanitizeFilename(filename)

	// Ensure destination directory exists
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", destPath, err)
	}

	// Create destination file
	fullPath := filepath.Join(destPath, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fullPath, err)
	}
	defer file.Close()

	// Download file content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	return nil
}

// sanitizeFilename removes dangerous characters from filenames to prevent path traversal
func sanitizeFilename(filename string) string {
	// Remove path separators and dangerous characters
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	filename = strings.ReplaceAll(filename, "\x00", "")

	// Trim whitespace
	filename = strings.TrimSpace(filename)

	// If filename becomes empty, use a default
	if filename == "" {
		filename = "file"
	}

	return filename
}
