# Dependency Tracer Subagent

## Name
deptrace

## Description
Traces function calls, type usage, and dependency relationships across Go codebase using autonomous ReAct loop with code analysis tools.

## Capabilities
- Searches for function and type definitions across the codebase
- Finds all references and usage locations of functions/types
- Traces call chains from entry points to target functions
- Analyzes import dependencies between packages
- Generates dependency diagrams in Mermaid format
- Identifies circular dependencies and potential issues
- Calculates impact radius for refactoring decisions

## When to Use
- Trace where a specific function is called from
- Find all usages of a type or interface
- Understand call chains and execution flow
- Analyze impact of changing a function signature
- Identify tightly coupled code that needs refactoring
- Generate architecture diagrams from code
- Find unused functions or dead code
- Map package dependencies

## Example Usage
```
Action: CallSubagent
Action Input: deptrace|Trace all call chains to the GetUser function in internal/service/user.go
```

## Input Format
Questions or commands must be in English. Should specify target function/type and optionally file path.

## Output Format
Can respond in any language based on the orchestrator's requirements. Returns call chains, dependency information, or Mermaid diagrams.
