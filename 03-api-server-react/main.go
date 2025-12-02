package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/toumakido/reAct/lib/bedrock"
	"github.com/toumakido/reAct/lib/tools"
	"github.com/toumakido/reAct/lib/types"
)

const systemPrompt = `You are a code analysis assistant that reads Go source files and answers questions about API server implementations.

## Core Principle

You MUST use the Available Tools to retrieve actual information from the file system. NEVER make assumptions or invent information about the codebase. All your reasoning and answers must be based on information obtained through tool usage.

## Available Tools

### 1. ListFiles
**Function**: Displays all files and directories under the data directory in tree format
**Usage**:
  Action: ListFiles
  Action Input: ListFiles
**System Response**: Directory structure is returned in the format "Observation: [tree structure]"
**When to Use**: When you need to understand the project structure or check which files exist

### 2. ReadFile
**Function**: Reads the contents of a specified Go source file
**Usage**:
  Action: ReadFile
  Action Input: [relative path from data directory]
**Input Examples**: cmd/api/main.go, internal/handler/user.go, pkg/middleware/auth.go
**System Response**: Full file contents are returned in the format "Observation: Content of [filename]:\n[contents]"
**When to Use**: When you need to examine code in a specific file

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

## Complete Execution Example

[Turn 1 - Your Output]
Thought: I need to check the directory structure first to understand the project layout.
Action: ListFiles
Action Input: ListFiles

[System Response]
Observation: [The system will return the actual directory structure]

[Turn 2 - Your Output]
Thought: Based on the structure, I should read a specific file to get more details.
Action: ReadFile
Action Input: [path to relevant file]

[System Response]
Observation: Content of [filename]:
[The system will return the actual file contents]

[Turn 3 - Your Output]
Thought: I now have all the necessary information to answer the question.
Final Answer: [Your detailed answer based on the information gathered]

## Final Answer Format

Once you have collected all necessary information, respond in this format:
Thought: [Reason why you can answer]
Final Answer: [Your complete and detailed answer to the user's question]`

const maxIterations = 15

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run . \"Your question here\"")
	}

	question := os.Args[1]

	ctx := context.Background()

	client, err := bedrock.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Bedrock client: %v", err)
	}

	if err := runReActLoop(ctx, client, question); err != nil {
		log.Fatalf("Error during ReAct loop: %v", err)
	}
}

func runReActLoop(ctx context.Context, client *bedrock.Client, question string) error {
	messages := []types.Message{
		{
			Role:    "user",
			Content: question,
		},
	}

	fmt.Println("=== Starting API Server Analysis ReAct Agent ===")
	fmt.Printf("Question: %s\n\n", question)

	for i := 0; i < maxIterations; i++ {
		fmt.Printf("--- Iteration %d ---\n", i+1)

		result, err := client.InvokeModel(ctx, systemPrompt, messages)
		if err != nil {
			return fmt.Errorf("failed to invoke model: %w", err)
		}

		fmt.Println(result.Text)
		fmt.Printf("\n[Token Usage] Input: %d, Output: %d, Total: %d\n\n",
			result.InputTokens, result.OutputTokens, result.InputTokens+result.OutputTokens)

		messages = append(messages, types.Message{
			Role:    "assistant",
			Content: result.Text,
		})

		if strings.Contains(result.Text, "Final Answer:") {
			fmt.Println("=== Agent Complete ===")
			return nil
		}

		action, actionInput, found := parseAction(result.Text)
		if !found {
			continue
		}

		var observation string
		switch action {
		case "ListFiles":
			result, err := tools.ListFilesTree()
			if err != nil {
				observation = fmt.Sprintf("Error listing files: %v", err)
			} else {
				observation = result
			}

		case "ReadFile":
			if actionInput == "" {
				observation = "Error: ReadFile requires a filename as Action Input"
			} else {
				content, err := tools.ReadFile(actionInput)
				if err != nil {
					observation = fmt.Sprintf("Error reading file: %v", err)
				} else {
					observation = fmt.Sprintf("Content of %s:\n%s", actionInput, content)
				}
			}

		default:
			observation = fmt.Sprintf("Error: Unknown action '%s'. Available actions: ListFiles, ReadFile", action)
		}

		fmt.Printf("Observation: %s\n\n", observation)

		messages = append(messages, types.Message{
			Role:    "user",
			Content: fmt.Sprintf("Observation: %s", observation),
		})
	}

	return fmt.Errorf("max iterations (%d) reached without final answer", maxIterations)
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
