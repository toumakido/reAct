# Code Analysis Subagent

## Name
codeanalysis

## Description
Performs comprehensive code analysis using autonomous ReAct loop with file exploration tools.

## Capabilities
- Explores directory structure (ListFiles tool)
- Reads Go source files (ReadFile tool)
- Analyzes code structure, relationships, and patterns
- Synthesizes information across multiple files
- Responds in any language (not limited to Japanese)

## When to Use
- Any question about the codebase structure
- Understanding API endpoints, handlers, or middleware
- Analyzing code relationships and architecture
- Explaining how specific features are implemented
- Any code-related query requiring file access

## Example Usage
```
Action: CallSubagent
Action Input: codeanalysis|What endpoints does this API server provide?
```

## Input Format
Questions must be in English.

## Output Format
Can respond in any language based on the orchestrator's requirements.
