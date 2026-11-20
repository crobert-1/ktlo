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
	"testing"
	"time"
)

const workflowProviderTestGoPackageResult = "2026-11-17T16:39:30.0068713Z ok   github.com/signalfx/splunk-otel-collector-chart/functional_tests/logs 1481.294s"
const workflowProviderTestGoTestPass = "2026-11-17T16:39:30.0022023Z PASS"
const workflowProviderTestPostActionEnd = "2026-11-17T16:39:30.1210958Z ##[end-action id=__self.__run_6;outcome=skipped;conclusion=skipped;duration_ms=0]"
const workflowProviderTestPostCleanup = "2026-11-17T16:39:30.1358551Z Post job cleanup."
const workflowProviderTestPostDebugStart = "2026-11-17T16:39:30.1203734Z ##[start-action display=Collect Kubernetes Cluster debug info;id=__self.__run_6]"
const workflowProviderTestPostUploadStart = "2026-11-17T16:39:30.1212891Z ##[start-action display=Upload Kubernetes Cluster debug info;id=__self.__actions_upload-artifact]"
const workflowProviderTestStepEnd = "2026-11-17T16:39:30.1200934Z ##[end-action id=__self.run-functional-tests;outcome=success;conclusion=success;duration_ms=1487838]"
const workflowProviderTestStepStart = "2026-11-17T16:39:30.0015319Z --- PASS: Test_OTLPAnnotationPrecedence (636.62s)"
const workflowProviderTestTestName = "assign completion-second lines using log ordering and step markers"
const workflowProviderTestNameExactGroupStepBoundary = "an exact GitHub step group marker selects the matching adjacent step"
const workflowProviderTestNextStepGroup = "2026-11-17T16:39:30.1203734Z ##[group]Run go test ./..."
const workflowProviderTestNextStepOutput = "2026-11-17T16:39:30.1212891Z go test output belongs to the next step"
const workflowProviderTestRunStepName = "Run go test ./..."
const workflowProviderTestFunctionalStepName = "Run ./.github/actions/functional-test"
const workflowProviderTestPostStepName = "Post Run ./.github/actions/functional-test"
const workflowProviderTestStepLogMismatchFormat = "step %d log mismatch:\n got: %q\nwant: %q"

func TestWorkflowJobLogsSplitCombinedOutputAtSecondPrecisionStepBoundary(t *testing.T) {
	nextStep := []WorkflowStep{
		{Number: 3, Name: workflowProviderTestFunctionalStepName, StartedAt: time.Date(2026, time.November, 17, 16, 12, 28, 0, time.UTC), CompletedAt: time.Date(2026, time.November, 17, 16, 39, 30, 0, time.UTC)},
		{Number: 5, Name: workflowProviderTestRunStepName, StartedAt: time.Date(2026, time.November, 17, 16, 39, 30, 0, time.UTC), CompletedAt: time.Date(2026, time.November, 17, 16, 39, 34, 0, time.UTC)},
	}
	postCleanupStepStart := time.Date(2026, time.November, 17, 16, 39, 30, 0, time.UTC)
	stepEnd := time.Date(2026, time.November, 17, 16, 39, 30, 0, time.UTC)
	testStepStart := time.Date(2026, time.November, 17, 16, 12, 28, 0, time.UTC)
	stepTests := []WorkflowStep{
		{Number: 3, Name: workflowProviderTestFunctionalStepName, StartedAt: testStepStart, CompletedAt: stepEnd},
		{Number: 5, Name: workflowProviderTestPostStepName, StartedAt: postCleanupStepStart, CompletedAt: time.Date(2026, time.November, 17, 16, 39, 34, 0, time.UTC)},
	}
	testCases := []struct {
		name           string
		jobOutput      string
		parseAsGoTest  bool
		workflowSteps  []WorkflowStep
		wantStepOutput map[int]string
	}{
		{
			name: workflowProviderTestTestName,
			jobOutput: strings.Join([]string{
				workflowProviderTestStepStart,
				workflowProviderTestGoTestPass,
				workflowProviderTestGoPackageResult,
				workflowProviderTestStepEnd,
				workflowProviderTestPostDebugStart,
				workflowProviderTestPostActionEnd,
				workflowProviderTestPostUploadStart,
				workflowProviderTestPostCleanup,
			}, logLineSeparator),
			parseAsGoTest: true,
			workflowSteps: stepTests,
			wantStepOutput: map[int]string{
				3: strings.Join([]string{
					workflowProviderTestStepStart,
					workflowProviderTestGoTestPass,
					workflowProviderTestGoPackageResult,
					workflowProviderTestStepEnd,
					workflowProviderTestPostDebugStart,
					workflowProviderTestPostActionEnd,
					workflowProviderTestPostUploadStart,
				}, logLineSeparator),
				5: workflowProviderTestPostCleanup,
			},
		},
		{
			name:          workflowProviderTestNameExactGroupStepBoundary,
			jobOutput:     strings.Join([]string{workflowProviderTestStepStart, workflowProviderTestNextStepGroup, workflowProviderTestNextStepOutput}, logLineSeparator),
			workflowSteps: nextStep,
			wantStepOutput: map[int]string{
				3: workflowProviderTestStepStart,
				5: strings.Join([]string{workflowProviderTestNextStepGroup, workflowProviderTestNextStepOutput}, logLineSeparator),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			jobLogs := WorkflowJobLogs{JobOutput: testCase.jobOutput}
			jobLogs.splitCombinedOutput(WorkflowJob{Steps: testCase.workflowSteps})
			stepOutputByNumber := make(map[int]string, len(jobLogs.Steps))
			for _, stepLog := range jobLogs.Steps {
				stepOutputByNumber[stepLog.Number] = stepLog.Text
			}
			for stepNumber, wantOutput := range testCase.wantStepOutput {
				if gotOutput := stepOutputByNumber[stepNumber]; gotOutput != wantOutput {
					t.Errorf(workflowProviderTestStepLogMismatchFormat, stepNumber, gotOutput, wantOutput)
				}
			}
			if testCase.parseAsGoTest {
				if _, err := (GoTestTextLogParser{}).Parse(stepOutputByNumber[3]); err != nil {
					t.Errorf("could not parse the test step's Go output after splitting: %v", err)
				}
			}
		})
	}
}
