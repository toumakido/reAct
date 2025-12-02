package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadFile reads a file from the data directory and returns its content
func ReadFile(filename string) (string, error) {
	path := filepath.Join("data", filename)

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return string(content), nil
}

// ListFiles lists all files in the data directory
func ListFiles() (string, error) {
	entries, err := os.ReadDir("data")
	if err != nil {
		return "", fmt.Errorf("failed to list files: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	if len(files) == 0 {
		return "No files found in data directory", nil
	}

	result := "Files in data directory:\n"
	for _, file := range files {
		result += fmt.Sprintf("- %s\n", file)
	}

	return result, nil
}
