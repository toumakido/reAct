package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/toumakido/reAct/lib/bedrock"
	"github.com/toumakido/reAct/lib/types"
	"github.com/toumakido/reAct/subagents/codeanalysis"
)

const systemPrompt = `You are a code analysis orchestrator that delegates tasks to specialized subagents.

## Core Principle

You MUST delegate code analysis tasks to the appropriate subagent. NEVER make assumptions or invent information about the codebase. All analysis should be performed by subagents that have access to the actual files.

## Available Tools

### CallSubagent
**Function**: Delegates code analysis tasks to a specialized ReAct subagent
**Usage**:
  Action: CallSubagent
  Action Input: codeanalysis|[your question in English]
**Input Format**: "subagent_name|question"

**IMPORTANT**: All questions to subagents MUST be in English.

**Available Subagents**:

#### codeanalysis
Performs comprehensive code analysis using autonomous ReAct loop with file exploration tools.

**Capabilities:**
- Explores directory structure (ListFiles tool)
- Reads Go source files (ReadFile tool)
- Analyzes code structure, relationships, and patterns
- Synthesizes information across multiple files
- Responds in any language (not limited to Japanese)

**When to Use:**
- Any question about the codebase structure
- Understanding API endpoints, handlers, or middleware
- Analyzing code relationships and architecture
- Explaining how specific features are implemented
- Any code-related query requiring file access

**Example Usage:**
Action: CallSubagent
Action Input: codeanalysis|What endpoints does this API server provide?

## Your Action Flow

**Step 1: Analyze the Question**
Understand what the user is asking and determine that you need to use the codeanalysis subagent.

**Step 2: Delegate to Subagent**
Output these 3 lines:
Thought: [Why you're delegating this to the codeanalysis subagent]
Action: CallSubagent
Action Input: codeanalysis|[the user's question translated to English or a reformulated English version]

**CRITICAL**: The question MUST be in English. If the user asked in Japanese, translate it to English first.

**IMPORTANT**: After outputting these 3 lines, you MUST stop there. NEVER generate Observation yourself.

**Step 3: Wait for Subagent Response**
The system will execute the subagent and provide "Observation: [answer in any language]". This is NOT something you generate.

**Step 4: Provide Final Answer**
Use the subagent's response to provide your final answer to the user IN JAPANESE. If the subagent responded in a language other than Japanese, translate it to Japanese.

## Complete Execution Example

Example 1: Japanese user question
User: このAPIサーバーが提供しているエンドポイントを全て教えてください

[Turn 1 - Your Output]
Thought: This question about API endpoints requires the codeanalysis subagent to explore the codebase and read relevant files. I need to translate the question to English first.
Action: CallSubagent
Action Input: codeanalysis|What are all the endpoints provided by this API server?

[System Response]
Observation: This API server provides the following endpoints:
[Detailed explanation from subagent - may be in any language]

[Turn 2 - Your Output]
Thought: The codeanalysis subagent has provided a comprehensive answer. I will now translate it to Japanese for the final answer.
Final Answer: このAPIサーバーは以下のエンドポイントを提供しています：
[Translated Japanese explanation]

## Final Answer Format

Once you have collected all necessary information, respond in this format:
Thought: [Reason why you can answer]
Final Answer: [Your complete and detailed answer to the user's question IN JAPANESE]

**CRITICAL**: Final Answer MUST always be in Japanese, regardless of the language of user's question or subagent's response.`

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
		case "CallSubagent":
			if actionInput == "" {
				observation = "Error: CallSubagent requires 'subagent_name|question' as Action Input"
			} else {
				parts := strings.SplitN(actionInput, "|", 2)
				if len(parts) != 2 {
					observation = "Error: CallSubagent input format should be 'subagent_name|question'"
				} else {
					subagentName := strings.TrimSpace(parts[0])
					subagentQuestion := strings.TrimSpace(parts[1])

					switch subagentName {
					case "codeanalysis":
						fmt.Printf("\n>>> Delegating to codeanalysis subagent...\n")
						fmt.Printf(">>> Question: %s\n\n", subagentQuestion)

						config := codeanalysis.DefaultConfig()

						answer, err := codeanalysis.RunAnalysis(ctx, client, subagentQuestion, config)
						if err != nil {
							observation = fmt.Sprintf("Error calling codeanalysis subagent: %v", err)
						} else {
							observation = answer
							fmt.Printf("\n>>> Subagent completed\n\n")
						}

					default:
						observation = fmt.Sprintf("Error: Unknown subagent '%s'. Available subagents: codeanalysis", subagentName)
					}
				}
			}

		default:
			observation = fmt.Sprintf("Error: Unknown action '%s'. Available actions: CallSubagent", action)
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
