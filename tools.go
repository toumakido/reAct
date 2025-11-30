package main

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
