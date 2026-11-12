package failingtest

import (
	"bufio"
	"strings"
)

const logLineSeparator = "\n"
const maxFailureSnippetBytes = 4 << 10
const maxFallbackSnippetLines = 3

// failureLineSelection marks which log lines to include in a failure snippet.
type failureLineSelection []bool

type failureIndicator string

const (
	failureIndicatorError   failureIndicator = "error"
	failureIndicatorFailed  failureIndicator = "failed"
	failureIndicatorFailure failureIndicator = "failure"
	failureIndicatorFatal   failureIndicator = "fatal"
	failureIndicatorPanic   failureIndicator = "panic"
)

var failureIndicators = [...]failureIndicator{
	failureIndicatorError,
	failureIndicatorFailed,
	failureIndicatorFailure,
	failureIndicatorFatal,
	failureIndicatorPanic,
}

func extractFailureLogSnippet(stepLog string) string {
	if stepLog == "" {
		return ""
	}
	lines := scanFailureLogLines(stepLog)
	selectedLines := selectedFailureContextLines(lines, failureContextLineSelection(lines))
	if len(selectedLines) == 0 {
		selectedLines = fallbackFailureLogSnippetLines(lines)
	}
	if len(selectedLines) == 0 {
		return ""
	}
	return limitFailureLogSnippet(strings.Join(selectedLines, logLineSeparator))
}

func scanFailureLogLines(stepLog string) []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(stepLog))
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines
}

func failureContextLineSelection(lines []string) failureLineSelection {
	lineSelection := make(failureLineSelection, len(lines))
	for lineIndex, line := range lines {
		if !containsFailureIndicator(line) {
			continue
		}
		for contextIndex := max(0, lineIndex-1); contextIndex <= min(len(lines)-1, lineIndex+1); contextIndex++ {
			lineSelection[contextIndex] = true
		}
	}
	return lineSelection
}

func containsFailureIndicator(line string) bool {
	lowerLine := strings.ToLower(line)
	for _, indicator := range failureIndicators {
		if strings.Contains(lowerLine, string(indicator)) {
			return true
		}
	}
	return false
}

func selectedFailureContextLines(lines []string, lineSelection failureLineSelection) []string {
	selectedLines := make([]string, 0)
	for lineIndex, line := range lines {
		if lineSelection[lineIndex] && line != "" {
			selectedLines = append(selectedLines, line)
		}
	}
	return selectedLines
}

func fallbackFailureLogSnippetLines(lines []string) []string {
	selectedLines := make([]string, 0, maxFallbackSnippetLines)
	for lineIndex := len(lines) - 1; lineIndex >= 0 && len(selectedLines) < maxFallbackSnippetLines; lineIndex-- {
		if lines[lineIndex] != "" {
			selectedLines = append([]string{lines[lineIndex]}, selectedLines...)
		}
	}
	return selectedLines
}

func limitFailureLogSnippet(snippet string) string {
	if len(snippet) > maxFailureSnippetBytes {
		// The byte offset caps persisted failure context at the configured maximum.
		snippet = snippet[:maxFailureSnippetBytes]
	}
	return snippet
}
