# Test Generator Subagent

## Name
testgen

## Description
Generates comprehensive Go test code using autonomous ReAct loop with code analysis and test generation tools.

## Capabilities
- Reads existing Go source files to understand implementation
- Analyzes function signatures, edge cases, and error conditions
- Generates table-driven tests following Go best practices
- Creates test cases for normal paths, edge cases, and error scenarios
- Suggests mock requirements for external dependencies
- Provides test coverage improvement recommendations

## When to Use
- Generate tests for existing functions without test coverage
- Add missing test cases to improve coverage
- Create table-driven tests from simple test functions
- Identify edge cases that need testing
- Generate boilerplate test code for new features
- Analyze testability and suggest refactoring

## Example Usage
```
Action: CallSubagent
Action Input: testgen|Generate tests for the CalculateTotal function in pkg/calculator/math.go
```

## Input Format
Questions or commands must be in English. Should specify target function/file path.

## Output Format
Can respond in any language based on the orchestrator's requirements. Returns generated test code or recommendations.
