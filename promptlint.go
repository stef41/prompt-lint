// Package promptlint analyzes LLM prompt templates for common issues
// such as injection vulnerabilities, missing context, ambiguous instructions,
// and excessive token usage.
package promptlint

import (
	"regexp"
	"strings"
)

const (
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityMedium   = "medium"
	SeverityLow      = "low"
	SeverityInfo     = "info"
)

// Rule defines a lint rule for prompt analysis.
type Rule struct {
	ID          string
	Description string
	Severity    string
	Check       func(prompt string) []Issue
}

// Issue represents a lint finding.
type Issue struct {
	RuleID     string `json:"rule_id"`
	Severity   string `json:"severity"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
	Position   int    `json:"position"`
}

// Result holds the complete lint result.
type Result struct {
	TotalIssues int     `json:"total_issues"`
	Critical    int     `json:"critical"`
	High        int     `json:"high"`
	Medium      int     `json:"medium"`
	Low         int     `json:"low"`
	Info        int     `json:"info"`
	Issues      []Issue `json:"issues"`
	Score       int     `json:"score"`
}

var injectionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)ignore\s+(all\s+)?(previous|above)\s+(instructions|prompts)`),
	regexp.MustCompile(`(?i)disregard\s+(your|all)\s+(rules|instructions)`),
	regexp.MustCompile(`(?i)you\s+are\s+now\s+(a|an|in)\s+`),
	regexp.MustCompile(`(?i)pretend\s+(you|that|to\s+be)`),
	regexp.MustCompile(`(?i)reveal\s+(your|the)\s+(system|initial)\s+(prompt|instructions)`),
}

// DefaultRules returns the built-in set of lint rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			ID: "PL001", Description: "Injection vulnerability", Severity: SeverityCritical,
			Check: func(p string) []Issue {
				var issues []Issue
				for _, re := range injectionPatterns {
					if loc := re.FindStringIndex(p); loc != nil {
						issues = append(issues, Issue{
							RuleID: "PL001", Severity: SeverityCritical,
							Message:    "Injection pattern: " + p[loc[0]:loc[1]],
							Suggestion: "Remove or sanitize",
							Position:   loc[0],
						})
					}
				}
				return issues
			},
		},
		{
			ID: "PL002", Description: "Missing role/context", Severity: SeverityMedium,
			Check: func(p string) []Issue {
				lower := strings.ToLower(p)
				hasRole := strings.Contains(lower, "you are") || strings.Contains(lower, "your role") || strings.Contains(lower, "system:")
				if !hasRole && len(p) > 50 {
					return []Issue{{RuleID: "PL002", Severity: SeverityMedium, Message: "No role definition", Suggestion: "Add role prefix"}}
				}
				return nil
			},
		},
		{
			ID: "PL003", Description: "Prompt too short", Severity: SeverityLow,
			Check: func(p string) []Issue {
				if len(strings.TrimSpace(p)) < 20 {
					return []Issue{{RuleID: "PL003", Severity: SeverityLow, Message: "Prompt too short", Suggestion: "Add more context"}}
				}
				return nil
			},
		},
		{
			ID: "PL004", Description: "No output format", Severity: SeverityInfo,
			Check: func(p string) []Issue {
				lower := strings.ToLower(p)
				hasFormat := strings.Contains(lower, "json") || strings.Contains(lower, "format") || strings.Contains(lower, "respond with")
				if !hasFormat && len(p) > 100 {
					return []Issue{{RuleID: "PL004", Severity: SeverityInfo, Message: "No output format", Suggestion: "Specify output format"}}
				}
				return nil
			},
		},
		{
			ID: "PL005", Description: "Excessive length", Severity: SeverityMedium,
			Check: func(p string) []Issue {
				if len(strings.Fields(p)) > 2000 {
					return []Issue{{RuleID: "PL005", Severity: SeverityMedium, Message: "Over 2000 words", Suggestion: "Reduce verbosity"}}
				}
				return nil
			},
		},
		{
			ID: "PL006", Description: "Template variable", Severity: SeverityHigh,
			Check: func(p string) []Issue {
				re := regexp.MustCompile(`\{\{\s*\w+\s*\}\}`)
				matches := re.FindAllStringIndex(p, -1)
				var issues []Issue
				for _, m := range matches {
					issues = append(issues, Issue{
						RuleID: "PL006", Severity: SeverityHigh,
						Message:    "Template variable: " + p[m[0]:m[1]],
						Suggestion: "Sanitize user input",
						Position:   m[0],
					})
				}
				return issues
			},
		},
	}
}

// Lint analyzes a prompt with the default rules.
func Lint(prompt string) Result {
	return LintWithRules(prompt, DefaultRules())
}

// LintWithRules analyzes a prompt with custom rules.
func LintWithRules(prompt string, rules []Rule) Result {
	result := Result{}
	for _, rule := range rules {
		issues := rule.Check(prompt)
		for _, issue := range issues {
			result.Issues = append(result.Issues, issue)
			switch issue.Severity {
			case SeverityCritical:
				result.Critical++
			case SeverityHigh:
				result.High++
			case SeverityMedium:
				result.Medium++
			case SeverityLow:
				result.Low++
			case SeverityInfo:
				result.Info++
			}
		}
	}
	result.TotalIssues = len(result.Issues)
	result.Score = 100
	result.Score -= result.Critical * 25
	result.Score -= result.High * 15
	result.Score -= result.Medium * 10
	result.Score -= result.Low * 5
	result.Score -= result.Info * 2
	if result.Score < 0 {
		result.Score = 0
	}
	return result
}
