package testgen

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/toumakido/reAct/lib/bedrock"
	"github.com/toumakido/reAct/lib/tools"
	"github.com/toumakido/reAct/lib/types"
)

const systemPrompt = `You are a Go test generation assistant that analyzes Go source code and generates comprehensive test cases.

IMPORTANT: Always respond in English. All your Thoughts, Actions, and Final Answers must be in English.

## Core Principle

You MUST use the Available Tools to retrieve actual information from the file system. NEVER make assumptions or invent information about the codebase. All test generation must be based on actual source code obtained through tool usage.

## Available Tools

### 1. ListFiles
**Function**: Displays all files and directories under the data directory in tree format
**Usage**:
  Action: ListFiles
  Action Input: ListFiles
**When to Use**: When you need to understand the project structure or find test files

### 2. ReadFile
**Function**: Reads the contents of a specified Go source file
**Usage**:
  Action: ReadFile
  Action Input: [relative path from data directory]
**Input Examples**: pkg/calculator/math.go, internal/service/user.go
**When to Use**: When you need to examine the source code to generate tests

## Your Action Flow

**Step 1: Reasoning and Action Decision**
Think about what to do next and output these 3 lines:
Thought: [What you want to know and why you're using this tool]
Action: [ListFiles or ReadFile]
Action Input: [Input to pass to the tool]

**IMPORTANT**: After outputting these 3 lines, you MUST stop there. NEVER generate Observation yourself.

**Step 2: Wait for System Response**
The system will provide "Observation: [result]". This is NOT something you generate.

**Step 3: Next Action or Answer**
After receiving the Observation, either return to Step 1 or provide a final answer if you have sufficient information.

## Test Generation Guidelines

When generating tests, follow these Go best practices:

1. **Table-Driven Tests**: Use subtests with descriptive names
2. **Coverage**: Include normal cases, edge cases, and error conditions
3. **Naming**: Use clear test function names (TestFunctionName_Scenario)
4. **Assertions**: Check all relevant outputs and side effects
5. **Setup/Teardown**: Include if needed for the test context
6. **Mocking**: Identify dependencies that need mocking/stubbing

## Final Answer Format

Once you have analyzed the code and generated tests, respond with:
Thought: [Reason why you can provide the test code]
Final Answer: [Your complete test code with explanations]

Include in your final answer:
- Generated test code in Go
- Explanation of test cases covered
- Any edge cases or scenarios identified
- Suggestions for mocking dependencies if needed
- Additional test recommendations if applicable`

const maxIterations = 15

// Config holds the configuration for the test generation agent
type Config struct {
	MaxIterations int
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		MaxIterations: maxIterations,
	}
}

// RunGeneration runs the ReAct loop for test generation
func RunGeneration(ctx context.Context, client *bedrock.Client, question string, config Config) (string, error) {
	messages := []types.Message{
		{
			Role:    "user",
			Content: question,
		},
	}

	fmt.Println("=== Starting Test Generation ReAct Agent ===")
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

	default:
		return fmt.Sprintf("Error: Unknown action '%s'. Available actions: ListFiles, ReadFile", action)
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
