package failingtest

import (
	"bufio"
	"errors"
	"strings"

	"ktlo/internal/testresults"
)

const (
	goCachedResult          = "(cached)"
	goLogTimestampSeparator = "Z "
	goRunLinePrefix         = "=== RUN "
	goTestResultPrefix      = "--- "
	goTestResultTokenSuffix = ":"
	packageTestIDSeparator  = "."
)

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

// GoTestLogParser parses output from the Go test command.
type GoTestLogParser struct{}

func defaultTestLogParsers() []TestLogParser {
	return []TestLogParser{GoTestLogParser{}}
}

func (parser GoTestLogParser) CanParse(stepLog string) bool {
	scanner := bufio.NewScanner(strings.NewReader(stepLog))
	for scanner.Scan() {
		line := parser.cleanGoTestLogLine(strings.TrimSpace(scanner.Text()))
		if strings.HasPrefix(line, goRunLinePrefix) {
			return true
		}
		if strings.HasPrefix(line, goTestResultPrefix) && parser.parseTestResultLine(line) != nil {
			return true
		}
		if parser.parseGoPackageResultLine(line) != nil {
			return true
		}
	}
	return false
}

func (parser GoTestLogParser) Parse(stepLog string) ([]testresults.Test, error) {
	parseState := newGoTestLogParseState()
	scanner := bufio.NewScanner(strings.NewReader(stepLog))
	for scanner.Scan() {
		parseState.processLine(parser, parser.cleanGoTestLogLine(strings.TrimSpace(scanner.Text())))
	}
	return parseState.finalizeTestResults(scanner.Err())
}

func (GoTestLogParser) cleanGoTestLogLine(line string) string {
	if timestampIndex := strings.Index(line, goLogTimestampSeparator); timestampIndex > 0 && timestampIndex < 40 {
		// Skip the timestamp and its trailing space to leave the log payload.
		return line[timestampIndex+2:]
	}
	return line
}

func (parser GoTestLogParser) parseTestResultLine(line string) *testLogResult {
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

func (parser GoTestLogParser) parseGoPackageResultLine(line string) *goPackageLogResult {
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

func (GoTestLogParser) outcomeFromString(value string) (testresults.Outcome, bool) {
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
	name    string
	outcome testresults.Outcome
	output  string
}

func (result *testLogResult) appendOutput(line string) {
	if len(result.output)+len(line)+len(logLineSeparator) > maxFailureSnippetBytes {
		return
	}
	if result.output != "" {
		result.output += logLineSeparator
	}
	result.output += line
}

func (result *testLogResult) toTestObservation(packageName string) testresults.Test {
	test := testresults.Test{
		ID:      packageName + packageTestIDSeparator + result.name,
		Kind:    testresults.TestKindTest,
		Outcome: result.outcome,
	}
	if result.outcome == testresults.OutcomeFailed {
		test.FailureMessage = result.output
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

func (state *goTestLogParseState) processLine(parser GoTestLogParser, line string) {
	if strings.HasPrefix(line, goRunLinePrefix) {
		state.setRunningTestName(line)
		return
	}
	if strings.HasPrefix(line, goTestResultPrefix) {
		state.setTestResult(parser, line)
		return
	}
	if packageResult := parser.parseGoPackageResultLine(line); packageResult != nil {
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
			state.runningTestsByName[state.runningTestName] = &testLogResult{name: state.runningTestName}
		}
	}
}

func (state *goTestLogParseState) setTestResult(parser GoTestLogParser, line string) {
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
	if !packageResult.cached {
		state.results = append(state.results, testresults.Test{ID: packageResult.name, Kind: testresults.TestKindPackage, Outcome: packageResult.outcome})
	}
	for _, packageTest := range state.packageTests {
		state.results = append(state.results, packageTest.toTestObservation(packageResult.name))
	}
	state.packageTests = state.packageTests[:0]
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
		return state.results, errors.New("Go test results were not followed by a package result")
	}
	return state.results, nil
}
