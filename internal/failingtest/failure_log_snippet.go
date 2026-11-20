// Copyright Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package failingtest

import (
	"bufio"
	"strings"
)

const logLineSeparator = "\n"
const maxFailureSnippetBytes = 4 << 10
const maxFallbackSnippetLines = 3
const failureDiffHeaderPrefix = "diff --git "
const failureDiffSummary = "One or more files were changed, this is not expected."

const (
	goModuleCommandLinePrefix      = "cd "
	goModuleDownloadFieldSeparator = " "
	goModuleDownloadLinePrefix     = "go: downloading "
	goModuleInstallLineMarker      = " && go install "
	goModuleMakeFailurePrefix      = "make: *** ["
	goModuleProxyInternalError     = "INTERNAL_ERROR"
	goModuleProxyStreamError       = "stream error: stream ID"
	goModuleProxyURLPrefix         = "https://proxy.golang.org/"
	goModuleProxyVersionMarker     = "/@v/"
	goModuleProxyZipSuffix         = ".zip"
)

// failureLineSelection marks which log lines to include in a failure snippet.
type failureLineSelection []bool

type failureIndicator string

const (
	failureIndicatorFail    failureIndicator = "fail"
	failureIndicatorError   failureIndicator = "error"
	failureIndicatorFailed  failureIndicator = "failed"
	failureIndicatorFailure failureIndicator = "failure"
	failureIndicatorFatal   failureIndicator = "fatal"
	failureIndicatorPanic   failureIndicator = "panic"
)

var failureIndicators = [...]failureIndicator{
	failureIndicatorError,
	failureIndicatorFail,
	failureIndicatorFailed,
	failureIndicatorFailure,
	failureIndicatorFatal,
	failureIndicatorPanic,
}

func extractFailureLogSnippet(stepLog string) string {
	if stepLog == "" {
		return ""
	}
	extractors := [...]failureLogSnippetExtractor{
		goTestJSONFailureSnippetExtractor{},
		goTestFailureSnippetExtractor{},
		failureDiffSnippetExtractor{},
		goModuleProxyFailureSnippetExtractor{},
		genericFailureSnippetExtractor{},
		fallbackFailureSnippetExtractor{},
	}
	for _, extractor := range extractors {
		if snippet := extractor.Extract(stepLog); snippet != "" {
			return snippet
		}
	}
	return ""
}

type failureLogSnippetExtractor interface {
	Extract(stepLog string) string
}

type genericFailureSnippetExtractor struct{}

func (genericFailureSnippetExtractor) Extract(stepLog string) string {
	lines := scanFailureLogLines(stepLog)
	selectedLines := selectedFailureContextLines(lines, failureContextLineSelection(lines))
	if len(selectedLines) == 0 {
		return ""
	}
	return limitFailureLogSnippet(strings.Join(selectedLines, logLineSeparator))
}

type fallbackFailureSnippetExtractor struct{}

func (fallbackFailureSnippetExtractor) Extract(stepLog string) string {
	lines := scanFailureLogLines(stepLog)
	selectedLines := fallbackFailureLogSnippetLines(lines)
	if len(selectedLines) == 0 {
		return ""
	}
	return limitFailureLogSnippet(strings.Join(selectedLines, logLineSeparator))
}

type goModuleProxyFailureSnippetExtractor struct{}

func (extractor goModuleProxyFailureSnippetExtractor) Extract(stepLog string) string {
	logLines := strings.Split(stepLog, logLineSeparator)
	for failureLineIndex, line := range logLines {
		failureLine := logLineContent(line)
		if !strings.Contains(failureLine, goModuleProxyStreamError) || !strings.Contains(failureLine, goModuleProxyInternalError) {
			continue
		}
		modulePath, moduleVersion, found := goModuleProxyFailureDetails(failureLine)
		if !found {
			continue
		}
		downloadLine := goModuleDownloadLinePrefix + modulePath + goModuleDownloadFieldSeparator + moduleVersion
		downloadLineIndex := -1
		for lineIndex := failureLineIndex - 1; lineIndex >= 0; lineIndex-- {
			if strings.TrimSpace(logLineContent(logLines[lineIndex])) == downloadLine {
				downloadLineIndex = lineIndex
				break
			}
		}
		return extractor.failureSnippet(logLines, modulePath, downloadLineIndex, failureLineIndex)
	}
	return ""
}

func (extractor goModuleProxyFailureSnippetExtractor) failureSnippet(logLines []string, modulePath string, downloadLineIndex int, failureLineIndex int) string {
	installCommand := extractor.goInstallCommand(logLines, downloadLineIndex, modulePath)
	snippetLines := make([]string, 0, 5)
	if installCommand != "" {
		snippetLines = append(snippetLines, installCommand)
	}
	if downloadLineIndex >= 0 {
		snippetLines = append(snippetLines, strings.TrimSpace(logLineContent(logLines[downloadLineIndex])))
	}
	snippetLines = append(snippetLines, strings.TrimSpace(logLineContent(logLines[failureLineIndex])))
	for lineIndex := failureLineIndex + 1; lineIndex < len(logLines); lineIndex++ {
		line := strings.TrimSpace(logLineContent(logLines[lineIndex]))
		// Expected fields: FAIL package path, or a standalone FAIL summary.
		packageFields := strings.Fields(line)
		if strings.HasPrefix(line, goTestResultPrefix+string(goTestResultFailed)+goTestResultTokenSuffix) ||
			len(packageFields) > 0 && strings.EqualFold(packageFields[0], string(goPackageResultFailed)) {
			snippetLines = append(snippetLines, line)
		}
		if strings.HasPrefix(line, goModuleMakeFailurePrefix) {
			snippetLines = append(snippetLines, line)
			break
		}
	}
	return limitFailureLogSnippetFromStart(strings.Join(snippetLines, logLineSeparator))
}

func (goModuleProxyFailureSnippetExtractor) goInstallCommand(logLines []string, downloadLineIndex int, modulePath string) string {
	if downloadLineIndex < 0 {
		return ""
	}
	for lineIndex := downloadLineIndex - 1; lineIndex >= 0; lineIndex-- {
		commandLine := strings.TrimSpace(logLineContent(logLines[lineIndex]))
		if !strings.Contains(commandLine, goModuleInstallLineMarker) {
			if strings.HasPrefix(commandLine, goModuleCommandLinePrefix) {
				break
			}
			continue
		}
		if strings.HasSuffix(commandLine, modulePath) {
			return commandLine
		}
		break
	}
	return ""
}

func goModuleProxyFailureDetails(failureLine string) (string, string, bool) {
	urlStart := strings.Index(failureLine, goModuleProxyURLPrefix)
	if urlStart < 0 {
		return "", "", false
	}
	// The URL starts at the proxy host and ends before its .zip suffix.
	urlEndOffset := strings.Index(failureLine[urlStart:], goModuleProxyZipSuffix)
	if urlEndOffset < 0 {
		return "", "", false
	}
	urlEnd := urlStart + urlEndOffset
	// The proxy URL path contains the module name followed by /@v/<version>.
	proxyPath := failureLine[urlStart+len(goModuleProxyURLPrefix) : urlEnd]
	versionMarkerIndex := strings.LastIndex(proxyPath, goModuleProxyVersionMarker)
	if versionMarkerIndex <= 0 {
		return "", "", false
	}
	// The marker separates the module path from its requested version.
	modulePath := proxyPath[:versionMarkerIndex]
	moduleVersion := proxyPath[versionMarkerIndex+len(goModuleProxyVersionMarker):]
	return modulePath, moduleVersion, modulePath != "" && moduleVersion != ""
}

type failureDiffSnippetExtractor struct{}

func (failureDiffSnippetExtractor) Extract(stepLog string) string {
	logLines := strings.Split(stepLog, logLineSeparator)
	diffStart := -1
	for lineIndex, line := range logLines {
		if strings.HasPrefix(strings.TrimSpace(logLineContent(line)), failureDiffHeaderPrefix) {
			diffStart = lineIndex
			// The preceding line may explain why this diff is a failure.
			if lineIndex > 0 && strings.Contains(logLineContent(logLines[lineIndex-1]), failureDiffSummary) {
				diffStart--
			}
			break
		}
	}
	if diffStart < 0 {
		return ""
	}
	snippetLines := make([]string, 0, len(logLines)-diffStart)
	for _, line := range logLines[diffStart:] {
		snippetLines = append(snippetLines, logLineContent(line))
	}
	return limitFailureLogSnippetFromStart(strings.Join(snippetLines, logLineSeparator))
}

func logLineContent(line string) string {
	timestampEnd := strings.Index(line, workflowLogTimestampSeparator)
	if timestampEnd >= 19 && timestampEnd <= 40 {
		// The payload begins after the timestamp's trailing Z and space.
		return line[timestampEnd+len(workflowLogTimestampSeparator):]
	}
	return line
}

func limitFailureLogSnippetFromStart(snippet string) string {
	return limitFailureLogTextFromStart(snippet, maxFailureSnippetBytes)
}

func limitFailureLogTextFromStart(snippet string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(snippet) <= maxBytes {
		return snippet
	}
	// Preserve the opening failure details within the byte limit.
	snippet = snippet[:maxBytes]
	if lineSeparatorIndex := strings.LastIndex(snippet, logLineSeparator); lineSeparatorIndex >= 0 {
		// Keep only complete lines before the byte cap.
		snippet = snippet[:lineSeparatorIndex]
	}
	return snippet
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
			// Select the failure line and its immediate neighbors.
			lineSelection[contextIndex] = true
		}
	}
	return lineSelection
}

func containsFailureIndicator(line string) bool {
	lowerLine := strings.ToLower(line)
	for _, indicator := range failureIndicators {
		indicatorText := string(indicator)
		searchStart := 0
		for searchStart < len(lowerLine) {
			indicatorIndex := strings.Index(lowerLine[searchStart:], indicatorText)
			if indicatorIndex < 0 {
				break
			}
			indicatorIndex += searchStart
			indicatorEnd := indicatorIndex + len(indicatorText)
			// Adjacent bytes determine whether the match is a whole word.
			if (indicatorIndex == 0 || !isFailureIndicatorWordByte(lowerLine[indicatorIndex-1])) &&
				(indicatorEnd == len(lowerLine) || !isFailureIndicatorWordByte(lowerLine[indicatorEnd])) {
				return true
			}
			searchStart = indicatorEnd
		}
	}
	return false
}

func isFailureIndicatorWordByte(value byte) bool {
	return value >= 'a' && value <= 'z' || value >= '0' && value <= '9' || value == '_'
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
	return limitFailureLogText(snippet, maxFailureSnippetBytes)
}

func limitFailureLogText(snippet string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(snippet) <= maxBytes {
		return snippet
	}
	// The suffix retains the most recent failure context within the byte limit.
	snippet = snippet[len(snippet)-maxBytes:]
	if lineSeparatorIndex := strings.Index(snippet, logLineSeparator); lineSeparatorIndex >= 0 {
		// Drop the first partial line introduced by the suffix boundary.
		snippet = snippet[lineSeparatorIndex+len(logLineSeparator):]
	}
	return snippet
}
