# 03-api-server-react

API server analysis using ReAct pattern with tree-based file exploration.

## Features

- Analyzes hierarchical Go API server code structure
- Tree-based directory listing with `ListFilesTree()`
- Hierarchical file path support in `ReadFile()`

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

```bash
go run . "What endpoints does this API server have?"
go run . "How is user authentication implemented?"
go run . "Explain the project structure"
go run . "How does the product management work?"
go run . "What middleware is used?"
```

## Changes from 02-code-react

- Added `ListFilesTree()` in `lib/tools/readfile.go` for tree representation
- Updated system prompt to support hierarchical paths
- Sample data is now a multi-layer API server structure
