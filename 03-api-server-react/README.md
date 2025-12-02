# 03-api-server-react

API server analysis using ReAct pattern with tree-based file exploration.

## Features

- Analyzes hierarchical Go API server code structure
- Tree-based directory listing with `ListFilesTree()`
- Hierarchical file path support in `ReadFile()`

## Data Structure

```
data/
├── cmd/api/main.go          # Server entry point with routing
├── internal/
│   ├── handler/user.go      # HTTP handlers
│   ├── model/user.go        # Data models
│   └── service/user.go      # Business logic
└── pkg/middleware/auth.go   # Middleware
```

## Usage

```bash
go run . "What endpoints does this API server have?"
go run . "How is user authentication implemented?"
go run . "Explain the project structure"
```

## Changes from 02-code-react

- Added `ListFilesTree()` in `lib/tools/readfile.go` for tree representation
- Updated system prompt to support hierarchical paths
- Sample data is now a multi-layer API server structure
