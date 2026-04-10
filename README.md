# prompt-lint

Go library to lint and analyze LLM prompt templates for common issues: injection vulnerabilities, missing context, ambiguous instructions, and excessive token usage. Returns a quality score 0-100.

## Installation

```bash
go get github.com/stef41/prompt-lint
```

## Usage

```go
result := promptlint.Lint("You are a helper. Process: {{ user_input }}")
fmt.Printf("Score: %d/100, Issues: %d\n", result.Score, result.TotalIssues)
for _, issue := range result.Issues {
    fmt.Printf("[%s] %s\n", issue.Severity, issue.Message)
}
```

## Rules

| ID | Description | Severity |
|----|-------------|----------|
| PL001 | Injection vulnerability | Critical |
| PL002 | Missing role/context | Medium |
| PL003 | Prompt too short | Low |
| PL004 | No output format | Info |
| PL005 | Excessive length | Medium |
| PL006 | Unescaped template variable | High |

## License

MIT
