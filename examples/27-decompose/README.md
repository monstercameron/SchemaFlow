# Decompose Example

Demonstrates the `Decompose` operation for breaking complex items into parts.

## What is Decompose?

The `Decompose` operation breaks down complex tasks, systems, or concepts:
- **Sequential**: Ordered steps that depend on each other
- **Hierarchical**: Tree structure with parent-child relationships
- **Parallel**: Independent workstreams that can run concurrently
- **Functional**: Components based on responsibility

## Key Features

- Multiple decomposition strategies
- Dependency tracking
- Effort estimation
- Parallelization hints

## Usage

```go
opts := ops.NewDecomposeOptions().
    WithStrategy("sequential").
    WithIncludeDependencies(true).
    WithMaxParts(10)

result, err := ops.Decompose(complexTask, opts)
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithStrategy()` | Decomposition approach | "sequential" |
| `WithMaxParts()` | Maximum number of parts | 0 (unlimited) |
| `WithMinParts()` | Minimum number of parts | 2 |
| `WithIncludeDependencies()` | Track part dependencies | false |
| `WithMaxDepth()` | Maximum hierarchy depth | 1 |
| `WithEstimateEffort()` | Include effort estimates | false |

## Running

```bash
cd examples/27-decompose
go run main.go
```

## Example Output

```
Original task: Build a web application with user authentication...

Decomposed into 6 parts:

  1. Set up project infrastructure
     Description: Initialize repository, CI/CD, and development environment
     Depends on: []

  2. Design database schema
     Description: Create tables for users, sessions, and application data
     Depends on: [1]

  3. Implement authentication system
     Description: Build login, registration, and session management
     Depends on: [1, 2]

  4. Create core API endpoints
     Description: Implement REST API for main application features
     Depends on: [2, 3]

  5. Build frontend interface
     Description: Create React components and integrate with API
     Depends on: [4]

  6. Add real-time notifications
     Description: Implement WebSocket server and push notifications
     Depends on: [3, 4, 5]
```

## Use Cases

- **Project planning**: Break down initiatives into tasks
- **System design**: Identify components and interfaces
- **Learning paths**: Create curriculum from complex topics
- **Refactoring**: Plan migration steps
