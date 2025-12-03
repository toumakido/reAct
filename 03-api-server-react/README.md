# 03-api-server-react

API server analysis using ReAct pattern with tree-based file exploration.

## Features

- Analyzes hierarchical Go API server code structure
- Tree-based directory listing with `ListFilesTree()`
- Hierarchical file path support in `ReadFile()`
- **Japanese response**: Final Answer is provided in Japanese
- **Reusable subagent**: Core logic extracted to `subagents/codeanalysis`

## Data Structure

```
data/
├── cmd/api/
│   ├── main.go              # Server entry point with routing
│   └── server.go            # Server struct with graceful shutdown
├── internal/
│   ├── handler/
│   │   ├── user.go          # User HTTP handlers
│   │   └── product.go       # Product HTTP handlers
│   ├── model/
│   │   ├── user.go          # User data models
│   │   └── product.go       # Product data models
│   └── service/
│       ├── user.go          # User business logic
│       └── product.go       # Product business logic
└── pkg/middleware/
    ├── auth.go              # Authentication middleware
    └── logging.go           # Request logging middleware
```

## API Endpoints

**Users:**
- `GET /api/users` - List all users
- `GET /api/users/{id}` - Get user by ID
- `POST /api/users` - Create new user

**Products:**
- `GET /api/products` - List all products
- `GET /api/products/{id}` - Get product by ID
- `POST /api/products` - Create new product
- `PATCH /api/products/{id}/stock` - Update product stock

## Usage

### As a CLI tool

```bash
go run . "What endpoints does this API server have?"
go run . "How is user authentication implemented?"
go run . "Explain the project structure"
go run . "How does the product management work?"
go run . "What middleware is used?"
```

### As a reusable subagent

```go
import (
    "context"
    "github.com/toumakido/reAct/lib/bedrock"
    "github.com/toumakido/reAct/subagents/codeanalysis"
)

func main() {
    ctx := context.Background()
    client, _ := bedrock.NewClient(ctx)

    config := codeanalysis.DefaultConfig()
    config.Verbose = true  // Enable detailed output
    config.MaxIterations = 20  // Customize max iterations

    answer, err := codeanalysis.RunAnalysis(ctx, client, "質問内容", config)
    if err != nil {
        // Handle error
    }

    // answer contains the Japanese response
    fmt.Println(answer)
}
```

## Architecture

### Two-Layer ReAct Structure

This implementation uses a two-layer ReAct architecture:

#### Layer 1: Orchestrator (main.go)
```
03-api-server-react/main.go
├── systemPrompt: Defines CallSubagent tool
└── runReActLoop()
    └── CallSubagent handler
        └── Delegates to subagents
```

**Role:** High-level task delegation
**Tool:** CallSubagent only
**Output:** Coordinates subagent execution

#### Layer 2: Subagent (subagents/codeanalysis)
```
subagents/codeanalysis/agent.go
├── systemPrompt: Defines ListFiles and ReadFile tools
└── RunAnalysis()
    ├── executeAction() - Executes ListFiles/ReadFile
    ├── parseAction() - Parses LLM output
    └── extractFinalAnswer() - Extracts Japanese answer
```

**Role:** File exploration and code analysis
**Tools:** ListFiles, ReadFile
**Output:** Japanese analysis results

### Benefits of This Architecture

1. **Separation of Concerns**: Orchestrator handles delegation, subagent handles file operations
2. **Reusability**: `subagents/codeanalysis` can be used independently in other projects
3. **Extensibility**: Easy to add new subagents (e.g., database analysis, API design)
4. **Clarity**: Each layer has a clear, focused responsibility

## Key Features

### 1. Japanese Final Answer
The system prompt instructs the LLM to provide the Final Answer in Japanese while keeping the reasoning process in English:

```
Final Answer: [Your complete and detailed answer - MUST be in Japanese]
```

### 2. Configurable Behavior
```go
type Config struct {
    MaxIterations int   // Maximum ReAct loop iterations
    Verbose       bool  // Enable/disable detailed output
}
```

### 3. Tool Execution
- `ListFiles`: Tree-structured directory listing
- `ReadFile`: Read specific Go source files

## Execution Flow

```
User Question
     ↓
[Orchestrator - main.go]
     │ systemPrompt: "You are a code analysis orchestrator"
     │ Tool: CallSubagent
     ↓
Thought: Need to analyze code, delegate to codeanalysis
Action: CallSubagent
Action Input: codeanalysis|ユーザーの質問
     ↓
[Subagent - codeanalysis]
     │ systemPrompt: "You are a code analysis assistant"
     │ Tools: ListFiles, ReadFile
     ↓
Thought: Check directory structure
Action: ListFiles → Observation: [tree]
     ↓
Thought: Read specific file
Action: ReadFile → Observation: [file content]
     ↓
Final Answer: [Japanese explanation]
     ↓
[Back to Orchestrator]
Observation: [Japanese explanation from subagent]
     ↓
Final Answer: [Pass through or summarize]
     ↓
User receives Japanese answer
```

## Changes from 02-code-react

- Introduced two-layer ReAct architecture
- Orchestrator layer only handles CallSubagent (no direct file access)
- Subagent layer handles all file operations
- Better separation of concerns for extensibility
- Can now easily add more specialized subagents
