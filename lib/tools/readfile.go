package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// ListFiles lists all files in the data directory (flat format)
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

// ListFilesTree lists all files and directories in the data directory in tree format
func ListFilesTree() (string, error) {
	var result string
	result += "data/\n"

	err := filepath.Walk("data", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == "data" {
			return nil
		}

		relPath, _ := filepath.Rel("data", path)
		depth := countDepth(relPath)

		prefix := buildTreePrefix(depth)
		name := filepath.Base(path)
		if info.IsDir() {
			name += "/"
		}

		result += fmt.Sprintf("%s%s\n", prefix, name)
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to list files: %w", err)
	}

	return result, nil
}

func countDepth(path string) int {
	if path == "." || path == "" {
		return 0
	}
	return strings.Count(filepath.ToSlash(path), "/") + 1
}

func buildTreePrefix(depth int) string {
	if depth == 0 {
		return ""
	}
	prefix := ""
	for i := 0; i < depth-1; i++ {
		prefix += "│   "
	}
	prefix += "├── "
	return prefix
}
