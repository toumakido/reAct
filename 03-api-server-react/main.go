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

const systemPrompt = `You are a code analysis assistant that can read Go source files to answer questions about API server implementations.

You must follow the ReAct (Reasoning and Acting) format strictly.

Your output format for each turn:

Thought: [Your reasoning about what to do next]
Action: [ActionName]
Action Input: [input for the action]

The system will then provide an Observation with the result.

Example of YOUR output:
Thought: I need to see the directory structure first to understand the API server layout.
Action: ListFiles
Action Input: ListFiles

After receiving the Observation from the system, continue:
Thought: Now I should read the main.go file to understand the server setup.
Action: ReadFile
Action Input: cmd/api/main.go

When you have gathered all necessary information, provide the final answer:

Final Answer: [Your complete answer to the user's question]

Available Actions:
- ListFiles: Lists all files and directories in tree format. No input required, just write "ListFiles" without any Action Input.
- ReadFile: Reads a Go source file from the data directory. Input should be the relative path from data directory (e.g., "cmd/api/main.go" or "internal/handler/user.go")

Important:
- YOU output: Thought, Action, Action Input
- SYSTEM provides: Observation
- Continue until you can provide the Final Answer
- The data directory contains a hierarchical Go API server project structure`

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
