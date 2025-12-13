package deptrace

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/toumakido/reAct/lib/bedrock"
	"github.com/toumakido/reAct/lib/tools"
	"github.com/toumakido/reAct/lib/types"
)

const systemPrompt = `You are a dependency tracing assistant that analyzes Go code to find function calls, type usage, and dependency relationships.

IMPORTANT: Always respond in English. All your Thoughts, Actions, and Final Answers must be in English.

## Core Principle

You MUST use the Available Tools to retrieve actual information from the file system. NEVER make assumptions or invent information about dependencies. All analysis must be based on actual source code obtained through tool usage.

## Available Tools

### 1. ListFiles
**Function**: Displays all files and directories under the data directory in tree format
**Usage**:
  Action: ListFiles
  Action Input: ListFiles
**When to Use**: When you need to understand the project structure or find all Go files

### 2. ReadFile
**Function**: Reads the contents of a specified Go source file
**Usage**:
  Action: ReadFile
  Action Input: [relative path from data directory]
**Input Examples**: internal/handler/user.go, pkg/service/auth.go
**When to Use**: When you need to examine code to find function calls, type usage, or imports

### 3. SearchPattern
**Function**: Searches for a pattern (function name, type name) across all Go files
**Usage**:
  Action: SearchPattern
  Action Input: [pattern to search]
**Input Examples**: GetUser, UserService, HandleRequest
**When to Use**: When you need to find all occurrences of a function or type name
**Note**: This is a simulated tool - in practice, you should read multiple files strategically

## Your Action Flow

**Step 1: Reasoning and Action Decision**
Think about what to do next and output these 3 lines:
Thought: [What you want to know and why you're using this tool]
Action: [ListFiles, ReadFile, or SearchPattern]
Action Input: [Input to pass to the tool]

**IMPORTANT**: After outputting these 3 lines, you MUST stop there. NEVER generate Observation yourself.

**Step 2: Wait for System Response**
The system will provide "Observation: [result]". This is NOT something you generate.

**Step 3: Next Action or Answer**
After receiving the Observation, either return to Step 1 or provide a final answer if you have sufficient information.

## Dependency Analysis Guidelines

When analyzing dependencies:

1. **Function Calls**: Look for direct function invocations, method calls
2. **Type Usage**: Find where structs, interfaces are instantiated or used
3. **Import Chains**: Track package imports and their relationships
4. **Call Chains**: Build paths from entry points (handlers, main) to target functions
5. **Visualization**: Generate Mermaid diagrams for complex relationships

## Mermaid Diagram Format

When generating diagrams, use this format:

` + "```mermaid" + `
graph TD
    A[Function A] --> B[Function B]
    B --> C[Function C]
    B --> D[Function D]
` + "```" + `

## Final Answer Format

Once you have traced the dependencies, respond with:
Thought: [Reason why you can provide the analysis]
Final Answer: [Your complete dependency analysis]

Include in your final answer:
- Call chains identified (from entry points to target)
- All locations where the function/type is used
- Dependency relationships
- Mermaid diagram if helpful for visualization
- Impact analysis for refactoring
- Any circular dependencies or issues found`

const maxIterations = 20

// Config holds the configuration for the dependency tracer agent
type Config struct {
	MaxIterations int
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		MaxIterations: maxIterations,
	}
}

// RunTrace runs the ReAct loop for dependency tracing
func RunTrace(ctx context.Context, client *bedrock.Client, question string, config Config) (string, error) {
	messages := []types.Message{
		{
			Role:    "user",
			Content: question,
		},
	}

	fmt.Println("=== Starting Dependency Tracer ReAct Agent ===")
	fmt.Printf("Question: %s\n\n", question)

	var finalAnswer string

	for i := 0; i < config.MaxIterations; i++ {
		fmt.Printf("--- Iteration %d ---\n", i+1)

		result, err := client.InvokeModel(ctx, systemPrompt, messages)
		if err != nil {
			return "", fmt.Errorf("failed to invoke model: %w", err)
		}

		fmt.Println(result.Text)
		fmt.Printf("\n[Token Usage] Input: %d, Output: %d, Total: %d\n\n",
			result.InputTokens, result.OutputTokens, result.InputTokens+result.OutputTokens)

		messages = append(messages, types.Message{
			Role:    "assistant",
			Content: result.Text,
		})

		if strings.Contains(result.Text, "Final Answer:") {
			finalAnswer = extractFinalAnswer(result.Text)
			fmt.Println("=== Agent Complete ===")
			return finalAnswer, nil
		}

		action, actionInput, found := parseAction(result.Text)
		if !found {
			continue
		}

		observation := executeAction(action, actionInput)

		messages = append(messages, types.Message{
			Role:    "user",
			Content: fmt.Sprintf("Observation: %s", observation),
		})
	}

	return "", fmt.Errorf("max iterations (%d) reached without final answer", config.MaxIterations)
}

func executeAction(action, actionInput string) string {
	switch action {
	case "ListFiles":
		result, err := tools.ListFilesTree()
		if err != nil {
			return fmt.Sprintf("Error listing files: %v", err)
		}
		return result

	case "ReadFile":
		if actionInput == "" {
			return "Error: ReadFile requires a filename as Action Input"
		}
		content, err := tools.ReadFile(actionInput)
		if err != nil {
			return fmt.Sprintf("Error reading file: %v", err)
		}
		return fmt.Sprintf("Content of %s:\n%s", actionInput, content)

	case "SearchPattern":
		if actionInput == "" {
			return "Error: SearchPattern requires a pattern as Action Input"
		}
		// SearchPattern is simulated - suggest reading files strategically
		return fmt.Sprintf("Note: To search for '%s', you should read relevant files based on project structure. Use ListFiles to see available files, then ReadFile to examine specific files where '%s' might be used.", actionInput, actionInput)

	default:
		return fmt.Sprintf("Error: Unknown action '%s'. Available actions: ListFiles, ReadFile, SearchPattern", action)
	}
}

func parseAction(response string) (action string, actionInput string, found bool) {
	actionRegex := regexp.MustCompile(`(?i)Action:\s*(\w+)`)
	actionMatch := actionRegex.FindStringSubmatch(response)
	if len(actionMatch) < 2 {
		return "", "", false
	}
	action = strings.TrimSpace(actionMatch[1])

	actionInputRegex := regexp.MustCompile(`(?i)Action Input:\s*(.+?)(?:\n|$)`)
	actionInputMatch := actionInputRegex.FindStringSubmatch(response)
	if len(actionInputMatch) >= 2 {
		actionInput = strings.TrimSpace(actionInputMatch[1])
	}

	return action, actionInput, true
}

func extractFinalAnswer(response string) string {
	lines := strings.Split(response, "\n")
	inFinalAnswer := false
	var answer []string

	for _, line := range lines {
		if strings.HasPrefix(line, "Final Answer:") {
			inFinalAnswer = true
			content := strings.TrimSpace(strings.TrimPrefix(line, "Final Answer:"))
			if content != "" {
				answer = append(answer, content)
			}
			continue
		}
		if inFinalAnswer {
			answer = append(answer, line)
		}
	}

	return strings.TrimSpace(strings.Join(answer, "\n"))
}
