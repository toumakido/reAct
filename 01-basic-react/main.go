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

const systemPrompt = `You are a helpful assistant that can read files to answer questions.

You must follow the ReAct (Reasoning and Acting) format strictly:

Thought: [Your reasoning about what to do next]
Action: ReadFile
Action Input: [filename]

After each action, you will receive an observation with the file content.
Then, continue with another Thought/Action/Observation cycle until you have enough information.

When you have gathered all necessary information and can provide the final answer, output:

Final Answer: [Your complete answer to the user's question]

Available Action:
- ReadFile: Reads a file from the data directory. Input should be just the filename (e.g., "start.txt")

Important:
- Always start by reading "start.txt" to begin your investigation
- Follow the clues in each file to find the next file to read
- Continue until you have enough information to answer the question
- Use "Final Answer:" only when you are ready to give the complete answer`

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

	fmt.Println("=== Starting ReAct Agent ===")
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

		if action != "ReadFile" {
			observation := fmt.Sprintf("Error: Unknown action '%s'. Only 'ReadFile' is supported.", action)
			fmt.Printf("Observation: %s\n\n", observation)

			messages = append(messages, types.Message{
				Role:    "user",
				Content: fmt.Sprintf("Observation: %s", observation),
			})
			continue
		}

		content, err := tools.ReadFile(actionInput)
		if err != nil {
			observation := fmt.Sprintf("Error reading file: %v", err)
			fmt.Printf("Observation: %s\n\n", observation)

			messages = append(messages, types.Message{
				Role:    "user",
				Content: fmt.Sprintf("Observation: %s", observation),
			})
			continue
		}

		observation := fmt.Sprintf("Observation: %s", content)
		fmt.Println(observation)
		fmt.Println()

		messages = append(messages, types.Message{
			Role:    "user",
			Content: observation,
		})
	}

	return fmt.Errorf("max iterations (%d) reached without final answer", maxIterations)
}

func parseAction(response string) (action string, actionInput string, found bool) {
	actionRegex := regexp.MustCompile(`(?i)Action:\s*(\w+)`)
	actionInputRegex := regexp.MustCompile(`(?i)Action Input:\s*(.+?)(?:\n|$)`)

	actionMatch := actionRegex.FindStringSubmatch(response)
	if len(actionMatch) < 2 {
		return "", "", false
	}
	action = strings.TrimSpace(actionMatch[1])

	actionInputMatch := actionInputRegex.FindStringSubmatch(response)
	if len(actionInputMatch) < 2 {
		return "", "", false
	}
	actionInput = strings.TrimSpace(actionInputMatch[1])

	return action, actionInput, true
}
