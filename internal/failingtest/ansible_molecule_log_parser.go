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
	"strings"

	"ktlo/internal/testresults"
)

const ansiColorCodeTerminator = "m"
const ansiEscapePrefix = "\x1b["
const ansibleMoleculeFailureAction = "idempotence"
const ansibleMoleculeFailureActionToken = "idempotence:"
const ansibleMoleculeFailureArrowToken = "➜"
const ansibleMoleculeFailureEndGroupMarker = "##[endgroup]"
const ansibleMoleculeFailureErrorPrefix = "ERROR "
const ansibleMoleculeFailureGroupMarker = "##[group]Molecule "
const ansibleMoleculeFailureGroupSuffix = " > idempotence"
const ansibleMoleculeFailureHeader = "Idempotence test failed because of the following tasks:"
const ansibleMoleculeFailureStatus = "Executed: Failed"
const ansibleMoleculeFailureStatusFieldSeparator = " "
const ansibleMoleculeFailureStatusSeparator = " ➜ "
const ansibleMoleculeFailureTaskPrefix = "* => "
const ansibleMoleculeFailureTaskSeparator = " : "
const ansibleMoleculeFailureTestIDSeparator = "/"

// AnsibleMoleculeLogParser parses failed Molecule idempotence tasks.
type AnsibleMoleculeLogParser struct{}

func (AnsibleMoleculeLogParser) Specificity() int {
	return ansibleMoleculeLogParserSpecificity
}

func (parser AnsibleMoleculeLogParser) CanParse(stepLog string) bool {
	return len(parser.failedTests(stepLog)) > 0
}

func (parser AnsibleMoleculeLogParser) Parse(stepLog string) ([]testresults.Test, error) {
	return parser.failedTests(stepLog), nil
}

func (AnsibleMoleculeLogParser) failedTests(stepLog string) []testresults.Test {
	lines := normalizedMoleculeLogLines(stepLog)
	failedScenarios := moleculeFailureScenarios(lines)
	failedTests := make([]testresults.Test, 0)
	seenTestIDs := make(map[string]struct{})
	activeScenario := ""
	for lineIndex, line := range lines {
		if scenario := moleculeIdempotenceGroupScenario(line); scenario != "" {
			activeScenario = scenario
			continue
		}
		if line == ansibleMoleculeFailureEndGroupMarker {
			activeScenario = ""
			continue
		}
		if activeScenario == "" || !failedScenarios[activeScenario] || !strings.Contains(line, ansibleMoleculeFailureHeader) {
			continue
		}
		// The task rows immediately follow this idempotence summary header.
		for _, failedTask := range moleculeFailureTasks(lines[lineIndex+1:]) {
			testID := activeScenario + ansibleMoleculeFailureTestIDSeparator + failedTask.role + ansibleMoleculeFailureTestIDSeparator + failedTask.name
			if _, duplicate := seenTestIDs[testID]; duplicate {
				continue
			}
			seenTestIDs[testID] = struct{}{}
			failureMessage := moleculeFailureMessage(activeScenario, failedTask.logLine)
			failedTests = append(failedTests, testresults.Test{
				ID:             testID,
				Kind:           testresults.TestKindTest,
				Outcome:        testresults.OutcomeFailed,
				FailureMessage: failureMessage,
			})
		}
	}
	return failedTests
}

type ansibleMoleculeFailedTask struct {
	logLine string
	name    string
	role    string
}

func normalizedMoleculeLogLines(stepLog string) []string {
	lines := strings.Split(stepLog, logLineSeparator)
	for lineIndex, line := range lines {
		lines[lineIndex] = strings.TrimSpace(stripANSIEscapeSequences(logLineContent(line)))
	}
	return lines
}

func stripANSIEscapeSequences(line string) string {
	var strippedLine strings.Builder
	for lineIndex := 0; lineIndex < len(line); {
		if strings.HasPrefix(line[lineIndex:], ansiEscapePrefix) {
			sequenceEndOffset := strings.Index(line[lineIndex+len(ansiEscapePrefix):], ansiColorCodeTerminator)
			if sequenceEndOffset >= 0 {
				lineIndex += len(ansiEscapePrefix) + sequenceEndOffset + len(ansiColorCodeTerminator)
				continue
			}
		}
		strippedLine.WriteByte(line[lineIndex])
		lineIndex++
	}
	return strippedLine.String()
}

func moleculeFailureScenarios(lines []string) map[string]bool {
	failedScenarios := make(map[string]bool)
	for _, line := range lines {
		statusIndex := strings.Index(line, ansibleMoleculeFailureStatusSeparator)
		if statusIndex < 0 || !strings.HasPrefix(line, ansibleMoleculeFailureErrorPrefix) {
			continue
		}
		// Expected fields: ERROR, scenario name, arrow, idempotence, status text.
		fields := strings.Fields(line)
		if len(fields) < 6 || fields[0] != strings.TrimSpace(ansibleMoleculeFailureErrorPrefix) ||
			fields[2] != ansibleMoleculeFailureArrowToken || fields[3] != ansibleMoleculeFailureActionToken ||
			fields[4]+ansibleMoleculeFailureStatusFieldSeparator+fields[5] != ansibleMoleculeFailureStatus {
			continue
		}
		failedScenarios[fields[1]] = true
	}
	return failedScenarios
}

func moleculeIdempotenceGroupScenario(line string) string {
	if !strings.HasPrefix(line, ansibleMoleculeFailureGroupMarker) || !strings.HasSuffix(line, ansibleMoleculeFailureGroupSuffix) {
		return ""
	}
	return strings.TrimSuffix(strings.TrimPrefix(line, ansibleMoleculeFailureGroupMarker), ansibleMoleculeFailureGroupSuffix)
}

func moleculeFailureTasks(lines []string) []ansibleMoleculeFailedTask {
	failedTasks := make([]ansibleMoleculeFailedTask, 0)
	seenTasks := make(map[string]struct{})
	insideFailureSummary := false
	for _, line := range lines {
		if strings.Contains(line, ansibleMoleculeFailureHeader) {
			insideFailureSummary = true
			continue
		}
		if !insideFailureSummary {
			continue
		}
		if !strings.HasPrefix(line, ansibleMoleculeFailureTaskPrefix) {
			if line == "" {
				continue
			}
			insideFailureSummary = false
			continue
		}
		failedTask, valid := parseMoleculeFailureTask(line)
		if !valid {
			continue
		}
		if _, duplicate := seenTasks[failedTask.logLine]; duplicate {
			continue
		}
		seenTasks[failedTask.logLine] = struct{}{}
		failedTasks = append(failedTasks, failedTask)
	}
	return failedTasks
}

func parseMoleculeFailureTask(line string) (ansibleMoleculeFailedTask, bool) {
	taskDetails := strings.TrimPrefix(line, ansibleMoleculeFailureTaskPrefix)
	taskSeparatorIndex := strings.Index(taskDetails, ansibleMoleculeFailureTaskSeparator)
	if taskSeparatorIndex <= 0 {
		return ansibleMoleculeFailedTask{}, false
	}
	// The separator divides the role name from the failed task name.
	role := strings.TrimSpace(taskDetails[:taskSeparatorIndex])
	// The task name starts after the role/task separator.
	taskName := strings.TrimSpace(taskDetails[taskSeparatorIndex+len(ansibleMoleculeFailureTaskSeparator):])
	if role == "" || taskName == "" {
		return ansibleMoleculeFailedTask{}, false
	}
	return ansibleMoleculeFailedTask{logLine: line, name: taskName, role: role}, true
}

func moleculeFailureMessage(scenario string, taskLine string) string {
	return scenario + ansibleMoleculeFailureStatusSeparator + ansibleMoleculeFailureAction + ": " + ansibleMoleculeFailureStatus +
		logLineSeparator + ansibleMoleculeFailureHeader + logLineSeparator + taskLine
}
