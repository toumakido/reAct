package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SubagentMetadata represents metadata for a subagent
type SubagentMetadata struct {
	Name         string
	Description  string
	Capabilities []string
	WhenToUse    []string
	Example      string
	InputFormat  string
	OutputFormat string
}

// LoadSubagentMetadata loads all subagent metadata from the subagents directory
func LoadSubagentMetadata(subagentsDir string) ([]SubagentMetadata, error) {
	entries, err := os.ReadDir(subagentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read subagents directory: %w", err)
	}

	var metadataList []SubagentMetadata
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		metadataPath := filepath.Join(subagentsDir, entry.Name(), "metadata.md")
		if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
			continue
		}

		content, err := os.ReadFile(metadataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read metadata file %s: %w", metadataPath, err)
		}

		metadata, err := parseMetadata(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse metadata for %s: %w", entry.Name(), err)
		}

		metadataList = append(metadataList, metadata)
	}

	return metadataList, nil
}

// parseMetadata parses markdown metadata into SubagentMetadata struct
func parseMetadata(content string) (SubagentMetadata, error) {
	var metadata SubagentMetadata
	lines := strings.Split(content, "\n")

	var currentSection string
	var listItems []string
	var exampleLines []string
	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Code block detection for examples
		if strings.HasPrefix(trimmed, "```") {
			if !inCodeBlock {
				inCodeBlock = true
				continue
			} else {
				inCodeBlock = false
				metadata.Example = strings.Join(exampleLines, "\n")
				exampleLines = nil
				continue
			}
		}

		if inCodeBlock {
			exampleLines = append(exampleLines, line)
			continue
		}

		// Section headers
		if strings.HasPrefix(trimmed, "## ") {
			// Save previous section's list items
			if currentSection != "" && len(listItems) > 0 {
				saveListItems(&metadata, currentSection, listItems)
				listItems = nil
			}
			currentSection = strings.TrimPrefix(trimmed, "## ")
			continue
		}

		// List items
		if strings.HasPrefix(trimmed, "- ") {
			listItems = append(listItems, strings.TrimPrefix(trimmed, "- "))
			continue
		}

		// Handle non-list content based on current section
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			switch currentSection {
			case "Name":
				metadata.Name = trimmed
			case "Description":
				metadata.Description = trimmed
			case "Input Format":
				metadata.InputFormat = trimmed
			case "Output Format":
				metadata.OutputFormat = trimmed
			}
		}
	}

	// Save last section's list items
	if len(listItems) > 0 {
		saveListItems(&metadata, currentSection, listItems)
	}

	return metadata, nil
}

func saveListItems(metadata *SubagentMetadata, section string, items []string) {
	switch section {
	case "Capabilities":
		metadata.Capabilities = items
	case "When to Use":
		metadata.WhenToUse = items
	}
}

// FormatSubagentInfo formats subagent metadata into a string for system prompt
func FormatSubagentInfo(metadata SubagentMetadata) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("#### %s\n", metadata.Name))
	builder.WriteString(fmt.Sprintf("%s\n\n", metadata.Description))

	if len(metadata.Capabilities) > 0 {
		builder.WriteString("**Capabilities:**\n")
		for _, cap := range metadata.Capabilities {
			builder.WriteString(fmt.Sprintf("- %s\n", cap))
		}
		builder.WriteString("\n")
	}

	if len(metadata.WhenToUse) > 0 {
		builder.WriteString("**When to Use:**\n")
		for _, usage := range metadata.WhenToUse {
			builder.WriteString(fmt.Sprintf("- %s\n", usage))
		}
		builder.WriteString("\n")
	}

	if metadata.Example != "" {
		builder.WriteString("**Example Usage:**\n")
		builder.WriteString(metadata.Example)
		builder.WriteString("\n")
	}

	return builder.String()
}
