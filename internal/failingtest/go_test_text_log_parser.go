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
	"errors"
	"strings"

	"ktlo/internal/testresults"
)

const errorGoTestResultsMissingPackageResult = "go test results were not followed by a package result"
const goCachedResult = "(cached)"
const goLogTimestampSeparator = "Z "
const goRunLinePrefix = "=== RUN "
const goTestAssertionError = "error:"
const goTestDiagnosticFileSuffix = ".go:"
const goTestDiagnosticLineSuffix = ":"
const goTestErrorTrace = "error trace:"
const goTestFatalError = "fatal:"
const goTestLogsHeader = "logs:"
const goTestMessagesFieldPrefix = "messages:"
const goTestPanic = "panic:"
const goTestResultPrefix = "--- "
const goTestResultTokenSuffix = ":"
const goTestRunningTestsHeader = "running tests:"
const goTestTestFieldPrefix = "test:"
const goTestTimeoutPanicPrefix = "panic: test timed out after"
const packageTestIDSeparator = "."

type goPackageResult string

const (
	goPackageResultFailed goPackageResult = "FAIL"
	goPackageResultPassed goPackageResult = "OK"
)

func (result goPackageResult) IsValid() bool {
	switch result {
	case goPackageResultFailed, goPackageResultPassed:
		return true
	default:
		return false
	}
}

type goTestResult string

const (
	goTestResultFailed  goTestResult = "FAIL"
	goTestResultPassed  goTestResult = "PASS"
	goTestResultSkipped goTestResult = "SKIP"
)

func (result goTestResult) IsValid() bool {
	switch result {
	case goTestResultFailed, goTestResultPassed, goTestResultSkipped:
		return true
	default:
		return false
	}
}

// GoTestTextLogParser parses human-readable output from the Go test command.
type GoTestTextLogParser struct{}

func (GoTestTextLogParser) Specificity() int {
	return goTestTextParserSpecificity
}

func (parser GoTestTextLogParser) CanParse(stepLog string) bool {
	scanner := bufio.NewScanner(strings.NewReader(stepLog))
	for scanner.Scan() {
		line := parser.cleanGoTestLogLine(strings.TrimSpace(scanner.Text()))
		if strings.HasPrefix(line, goTestResultPrefix) && parser.parseTestResultLine(line) != nil {
			return true
		}
		if parser.parseGoPackageResultLine(line) != nil {
			return true
		}
	}
	return false
}

func (parser GoTestTextLogParser) Parse(stepLog string) ([]testresults.Test, error) {
	parseState := newGoTestLogParseState()
	scanner := bufio.NewScanner(strings.NewReader(stepLog))
	for scanner.Scan() {
		parseState.processLine(parser, parser.cleanGoTestLogLine(strings.TrimRight(scanner.Text(), "\r")))
	}
	return parseState.finalizeTestResults(scanner.Err())
}

func (GoTestTextLogParser) cleanGoTestLogLine(line string) string {
	if timestampIndex := strings.Index(line, goLogTimestampSeparator); timestampIndex > 0 && timestampIndex < 40 {
		// Skip the timestamp and its trailing space to leave the log payload.
		return line[timestampIndex+2:]
	}
	return line
}

func (parser GoTestTextLogParser) parseTestResultLine(line string) *testLogResult {
	// Expected fields: ---, outcome token, test name, optional duration.
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return nil
	}
	// The second field is the outcome token in --- PASS|FAIL|SKIP output.
	result := goTestResult(strings.ToUpper(strings.TrimSuffix(fields[1], goTestResultTokenSuffix)))
	if !result.IsValid() {
		return nil
	}
	outcome, valid := parser.outcomeFromString(string(result))
	if !valid {
		return nil
	}
	// The third field is the test name in --- PASS|FAIL|SKIP output.
	return &testLogResult{name: fields[2], outcome: outcome}
}

func (parser GoTestTextLogParser) parseGoPackageResultLine(line string) *goPackageLogResult {
	// Expected fields: package outcome, import path, optional cached marker.
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}
	result := goPackageResult(strings.ToUpper(fields[0]))
	if !result.IsValid() {
		return nil
	}
	outcome, valid := parser.outcomeFromString(string(result))
	if !valid {
		return nil
	}
	// The second field after the package result is the package import path.
	packageName := fields[1]
	// If present, the third field marks a cached package result.
	cached := len(fields) > 2 && fields[2] == goCachedResult
	return &goPackageLogResult{
		cached:  cached,
		name:    packageName,
		outcome: outcome,
	}
}

func (GoTestTextLogParser) outcomeFromString(value string) (testresults.Outcome, bool) {
	switch strings.ToUpper(value) {
	case string(goPackageResultPassed), string(goTestResultPassed):
		return testresults.OutcomePassed, true
	case string(goPackageResultFailed):
		return testresults.OutcomeFailed, true
	case string(goTestResultSkipped):
		return testresults.OutcomeSkipped, true
	default:
		return "", false
	}
}

type testLogResult struct {
	assertionFailureDetails  string
	capturingAssertion       bool
	capturingPanic           bool
	failureTestFieldSeen     bool
	lastAssertionFailureLine string
	logsSectionStarted       bool
	name                     string
	outcome                  testresults.Outcome
	output                   string
	panicFailureDetails      string
	previousOutputLine       string
	runLine                  string
	timeoutSummaryComplete   bool
	timeoutSummaryStarted    bool
	timedOut                 bool
}

func (result *testLogResult) appendOutput(line string) {
	trimmedLine := strings.TrimSpace(line)
	if result.timeoutSummaryComplete {
		return
	}
	if strings.HasPrefix(trimmedLine, goTestTimeoutPanicPrefix) {
		result.timedOut = true
	}
	if result.timedOut && trimmedLine == goTestRunningTestsHeader {
		result.timeoutSummaryStarted = true
	}
	if result.timedOut && result.timeoutSummaryStarted && trimmedLine == "" {
		result.timeoutSummaryComplete = true
		return
	}
	result.appendAssertionFailureDetail(line)
	result.appendPanicFailureDetail(line)
	if result.output != "" {
		result.output += logLineSeparator
	}
	result.output += line
	result.output = limitFailureLogSnippet(result.output)
	result.previousOutputLine = line
}

func (result *testLogResult) appendAssertionFailureDetail(line string) {
	var lowerLine string
	trimmedLine := strings.TrimSpace(line)
	lowerLine = strings.ToLower(trimmedLine)
	if lowerLine == goTestLogsHeader {
		result.logsSectionStarted = true
		result.capturingAssertion = false
		result.failureTestFieldSeen = false
		return
	}
	if result.logsSectionStarted {
		return
	}
	if strings.HasPrefix(lowerLine, goTestErrorTrace) {
		if isGoTestDiagnosticLine(result.previousOutputLine) {
			result.appendAssertionFailureLine(result.previousOutputLine)
		}
		result.capturingAssertion = true
		result.failureTestFieldSeen = false
		result.appendAssertionFailureLine(line)
		return
	}
	if !result.capturingAssertion {
		return
	}
	if result.failureTestFieldSeen {
		if strings.HasPrefix(lowerLine, goTestMessagesFieldPrefix) {
			result.appendAssertionFailureLine(line)
			result.capturingAssertion = false
			result.failureTestFieldSeen = false
			return
		}
		if trimmedLine == "" {
			result.appendAssertionFailureLine(line)
			return
		}
		result.capturingAssertion = false
		result.failureTestFieldSeen = false
		return
	}
	result.appendAssertionFailureLine(line)
	if strings.HasPrefix(lowerLine, goTestTestFieldPrefix) {
		result.failureTestFieldSeen = true
	}
}

func (result *testLogResult) appendAssertionFailureLine(line string) {
	if result.lastAssertionFailureLine == line {
		return
	}
	if result.assertionFailureDetails != "" {
		result.assertionFailureDetails += logLineSeparator
	}
	result.assertionFailureDetails += line
	result.assertionFailureDetails = limitFailureLogTextFromStart(result.assertionFailureDetails, maxFailureSnippetBytes)
	result.lastAssertionFailureLine = line
}

func (result *testLogResult) appendPanicFailureDetail(line string) {
	trimmedLine := strings.TrimSpace(line)
	if strings.EqualFold(trimmedLine, goTestLogsHeader) {
		result.logsSectionStarted = true
		result.capturingPanic = false
		return
	}
	if result.logsSectionStarted {
		return
	}
	if strings.HasPrefix(strings.ToLower(trimmedLine), goTestTimeoutPanicPrefix) {
		return
	}
	if strings.HasPrefix(strings.ToLower(trimmedLine), goTestPanic) {
		result.capturingPanic = true
	}
	if !result.capturingPanic {
		return
	}
	if result.panicFailureDetails != "" {
		result.panicFailureDetails += logLineSeparator
	}
	result.panicFailureDetails += line
	result.panicFailureDetails = limitFailureLogTextFromStart(result.panicFailureDetails, maxFailureSnippetBytes)
}

func isGoTestDiagnosticLine(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return strings.HasSuffix(trimmedLine, goTestDiagnosticLineSuffix) && strings.Contains(trimmedLine, goTestDiagnosticFileSuffix)
}

func (result *testLogResult) failureMessage() string {
	if result.assertionFailureDetails != "" {
		return result.assertionFailureDetails
	}
	if result.panicFailureDetails != "" {
		return result.panicFailureDetails
	}
	return result.output
}

func (result *testLogResult) hasFailureOutput() bool {
	for _, line := range strings.Split(result.failureMessage(), logLineSeparator) {
		lowerLine := strings.ToLower(strings.TrimSpace(line))
		if strings.HasPrefix(lowerLine, goTestAssertionError) ||
			strings.HasPrefix(lowerLine, goTestErrorTrace) ||
			strings.HasPrefix(lowerLine, goTestFatalError) ||
			strings.HasPrefix(lowerLine, goTestPanic) {
			return true
		}
	}
	return false
}

func (result *testLogResult) hasFailedDescendant(packageTests []testLogResult) bool {
	for _, candidate := range packageTests {
		if candidate.outcome == testresults.OutcomeFailed && strings.HasPrefix(candidate.name, result.name+"/") {
			return true
		}
	}
	return false
}

func (result *testLogResult) hasIndependentFailureOutput(packageTests []testLogResult) bool {
	if !result.hasFailureOutput() {
		return false
	}
	parentOutput := strings.TrimSpace(result.failureMessage())
	for _, candidate := range packageTests {
		if candidate.outcome != testresults.OutcomeFailed || !strings.HasPrefix(candidate.name, result.name+"/") {
			continue
		}
		if strings.Contains(strings.TrimSpace(candidate.failureMessage()), parentOutput) {
			return false
		}
	}
	return true
}

func (result *testLogResult) failureSnippet() string {
	if result.runLine == "" {
		return limitFailureLogSnippet(result.failureMessage())
	}
	header := limitFailureLogSnippet(result.runLine) + logLineSeparator
	if result.timedOut {
		return header + result.output
	}
	if result.assertionFailureDetails != "" {
		output := limitFailureLogTextFromStart(result.assertionFailureDetails, maxFailureSnippetBytes-len(header))
		return header + output
	}
	if result.panicFailureDetails != "" {
		output := limitFailureLogTextFromStart(result.panicFailureDetails, maxFailureSnippetBytes-len(header))
		return header + output
	}
	output := limitFailureLogText(result.output, maxFailureSnippetBytes-len(header))
	if output == "" {
		return strings.TrimSuffix(header, logLineSeparator)
	}
	return header + output
}

func (result *testLogResult) toTestObservation(packageName string) testresults.Test {
	test := testresults.Test{
		ID:      packageName + packageTestIDSeparator + result.name,
		Kind:    testresults.TestKindTest,
		Outcome: result.outcome,
	}
	if result.outcome == testresults.OutcomeFailed {
		test.FailureMessage = result.failureMessage()
	}
	return test
}

type goPackageLogResult struct {
	cached  bool
	name    string
	outcome testresults.Outcome
}

type goTestLogParseState struct {
	currentTest     *testLogResult
	failureLogTest  *testLogResult
	packageTests    []testLogResult
	results         []testresults.Test
	runningTestName string
	// Preserve parent output until nested subtests emit their result lines.
	runningTestsByName map[string]*testLogResult
}

func newGoTestLogParseState() *goTestLogParseState {
	return &goTestLogParseState{
		packageTests:       make([]testLogResult, 0),
		results:            make([]testresults.Test, 0),
		runningTestsByName: make(map[string]*testLogResult),
	}
}

func (state *goTestLogParseState) omitAggregateTestFailures(packageTests []testLogResult) []testLogResult {
	retainedTests := make([]testLogResult, 0, len(packageTests))
	for _, test := range packageTests {
		if test.outcome == testresults.OutcomeFailed && test.hasFailedDescendant(packageTests) && !test.hasIndependentFailureOutput(packageTests) {
			continue
		}
		retainedTests = append(retainedTests, test)
	}
	return retainedTests
}

func (state *goTestLogParseState) latestFailedTest() *testLogResult {
	for testIndex := len(state.packageTests) - 1; testIndex >= 0; testIndex-- {
		if state.packageTests[testIndex].outcome != testresults.OutcomeFailed {
			continue
		}
		failedTest := state.packageTests[testIndex]
		return &failedTest
	}
	return nil
}

func (state *goTestLogParseState) finalizeFailureLogTest() {
	state.queueCurrentTestForPackage()
	state.packageTests = state.omitAggregateTestFailures(state.packageTests)
	state.failureLogTest = state.latestFailedTest()
}

func (state *goTestLogParseState) processLine(parser GoTestTextLogParser, line string) {
	parsedLine := strings.TrimSpace(line)
	if strings.HasPrefix(parsedLine, goRunLinePrefix) {
		state.setRunningTestName(parsedLine)
		return
	}
	if strings.HasPrefix(parsedLine, goTestResultPrefix) {
		state.setTestResult(parser, parsedLine)
		return
	}
	if packageResult := parser.parseGoPackageResultLine(parsedLine); packageResult != nil {
		state.setPackageResult(packageResult)
		return
	}
	state.appendLineToActiveTestOutput(line)
}

func (state *goTestLogParseState) setRunningTestName(line string) {
	// Expected fields: ===, RUN, test or subtest name.
	runningFields := strings.Fields(line)
	if len(runningFields) >= 3 {
		// The third field is the test or subtest name in === RUN output.
		state.runningTestName = runningFields[2]
		if state.runningTestsByName[state.runningTestName] == nil {
			state.runningTestsByName[state.runningTestName] = &testLogResult{
				name:    state.runningTestName,
				runLine: line,
			}
		}
	}
}

func (state *goTestLogParseState) setTestResult(parser GoTestTextLogParser, line string) {
	state.queueCurrentTestForPackage()
	parsedTest := parser.parseTestResultLine(line)
	if parsedTest == nil {
		state.currentTest = nil
	} else {
		state.currentTest = state.runningTestsByName[parsedTest.name]
		if state.currentTest == nil {
			state.currentTest = &testLogResult{name: parsedTest.name}
		}
		state.currentTest.outcome = parsedTest.outcome
		delete(state.runningTestsByName, parsedTest.name)
	}
	state.runningTestName = ""
}

func (state *goTestLogParseState) setPackageResult(packageResult *goPackageLogResult) {
	state.queueCurrentTestForPackage()
	if packageResult.outcome == testresults.OutcomeFailed {
		state.queueUnclosedFailedTest()
	}
	state.finalizeFailureLogTest()
	if !packageResult.cached {
		state.results = append(state.results, testresults.Test{ID: packageResult.name, Kind: testresults.TestKindPackage, Outcome: packageResult.outcome})
	}
	for _, packageTest := range state.packageTests {
		state.results = append(state.results, packageTest.toTestObservation(packageResult.name))
	}
	state.packageTests = state.packageTests[:0]
	state.runningTestName = ""
	state.runningTestsByName = make(map[string]*testLogResult)
}

func (state *goTestLogParseState) queueUnclosedFailedTest() {
	var failureTest *testLogResult
	for _, runningTest := range state.runningTestsByName {
		if !runningTest.hasFailureOutput() {
			continue
		}
		if failureTest != nil {
			return
		}
		failureTest = runningTest
	}
	if failureTest == nil {
		return
	}
	failureTest.outcome = testresults.OutcomeFailed
	state.packageTests = append(state.packageTests, *failureTest)
	delete(state.runningTestsByName, failureTest.name)
	if state.runningTestName == failureTest.name {
		state.runningTestName = ""
	}
}

func (state *goTestLogParseState) queueCurrentTestForPackage() {
	if state.currentTest == nil {
		return
	}
	state.packageTests = append(state.packageTests, *state.currentTest)
	state.currentTest = nil
}

func (state *goTestLogParseState) appendLineToActiveTestOutput(line string) {
	if runningTest := state.runningTestsByName[state.runningTestName]; runningTest != nil {
		runningTest.appendOutput(line)
	}
	if state.currentTest != nil && state.runningTestName == "" && state.currentTest.outcome == testresults.OutcomeFailed {
		state.currentTest.appendOutput(line)
	}
}

func (state *goTestLogParseState) finalizeTestResults(scanErr error) ([]testresults.Test, error) {
	state.queueCurrentTestForPackage()
	if scanErr != nil {
		return state.results, scanErr
	}
	if len(state.packageTests) > 0 {
		return state.results, errors.New(errorGoTestResultsMissingPackageResult)
	}
	return state.results, nil
}
