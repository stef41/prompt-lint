package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	pl "github.com/stef41/prompt-lint"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: prompt-lint <prompt-text>\n")
		fmt.Fprintf(os.Stderr, "       prompt-lint --file <path>\n")
		os.Exit(1)
	}
	var prompt string
	jsonOutput := false
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--file":
			if i+1 < len(os.Args) {
				data, err := os.ReadFile(os.Args[i+1])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				prompt = string(data)
				i++
			}
		case "--json":
			jsonOutput = true
		default:
			if prompt == "" {
				prompt = strings.Join(os.Args[i:], " ")
				i = len(os.Args)
			}
		}
	}
	if prompt == "" {
		fmt.Fprintf(os.Stderr, "No prompt provided\n")
		os.Exit(1)
	}
	result := pl.Lint(prompt)
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(result)
	} else {
		fmt.Printf("Score: %d/100\n", result.Score)
		fmt.Printf("Issues: %d (critical=%d high=%d medium=%d low=%d info=%d)\n",
			result.TotalIssues, result.Critical, result.High, result.Medium, result.Low, result.Info)
		for _, issue := range result.Issues {
			fmt.Printf("  [%s] %s - %s\n", issue.Severity, issue.Message, issue.Suggestion)
		}
	}
	if result.Critical > 0 {
		os.Exit(2)
	}
}
