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
	"context"
	"strings"
	"testing"
	"time"

	"ktlo/internal/common"
	"ktlo/internal/testresults"
)

const workflowRecordBuilderTestFailureConclusion = "failure"
const workflowRecordBuilderTestPassConclusion = "success"
const workflowRecordBuilderTestStepName = "Unit tests"
const workflowRecordBuilderTestRunID = "43485897475"
const workflowRecordBuilderTestHeadSHA = "0123456789abcdef"
const workflowRecordBuilderTestRepositoryOwner = "signalfx"
const workflowRecordBuilderTestRepositoryName = "splunk-otel-collector"
const workflowRecordBuilderTestWorkflowName = "build-and-test"
const workflowRecordBuilderTestTestName = "TestTunnelMultipleTargetsConcurrently"
const workflowRecordBuilderTestPackage = "github.com/signalfx/splunk-otel-collector/internal/receiver/gnmireceiver"
const workflowRecordBuilderTestWindowsIISPackage = "github.com/signalfx/splunk-otel-collector/tests/zeroconfig/windows"
const workflowRecordBuilderTestCommonDimensionsPackage = "github.com/signalfx/splunk-otel-collector/internal/otelcol"
const workflowRecordBuilderTestCommonDimensionsTestName = "TestCommonDimensions"
const workflowRecordBuilderTestAssertionFailureText = "Error: Should be true"
const workflowRecordBuilderTestLinuxJobID = "job-linux"
const workflowRecordBuilderTestWindowsJobID = "job-windows"
const workflowRecordBuilderTestCommonDimensionsJobID = "job-common-dimensions"
const workflowRecordBuilderTestIISJobID = "job-iis"
const workflowRecordBuilderTestMainBranch = "main"
const workflowRecordBuilderTestJobURLPrefix = "https://github.com/signalfx/splunk-otel-collector/actions/runs/"
const workflowRecordBuilderTestJobURLSeparator = "/job/"
const workflowRecordBuilderTestSubtestSeparator = "."
const workflowRecordBuilderTestIISName = "TestWindowsIISInstrumentation"
const workflowRecordBuilderTestIISFailureMessage = "Failed to receive expected traces"
const workflowRecordBuilderTestLinuxJob = "test-ubuntu-24.04 / unit-test"
const workflowRecordBuilderTestWindowsJob = "test-windows-2025 / unit-test"
const workflowRecordBuilderTestCommonDimensionsJob = "test-ubuntu-24.04 / common-dimensions"
const workflowRecordBuilderTestIISJob = "test-windows-2025 / iis"
const workflowRecordBuilderTestMatrixFailureLog = `=== RUN   TestTunnelMultipleTargetsConcurrently
receiver_test.go:489:
    Error Trace: receiver_test.go:489
    Error: Should be true
    Test: TestTunnelMultipleTargetsConcurrently
--- FAIL: TestTunnelMultipleTargetsConcurrently (0.00s)
FAIL github.com/signalfx/splunk-otel-collector/internal/receiver/gnmireceiver 0.010s`
const workflowRecordBuilderTestMatrixPassLog = `=== RUN   TestTunnelMultipleTargetsConcurrently
--- PASS: TestTunnelMultipleTargetsConcurrently (0.00s)
ok github.com/signalfx/splunk-otel-collector/internal/receiver/gnmireceiver 0.010s`
const workflowRecordBuilderTestCommonDimensionsLog = `=== RUN   TestCommonDimensions
--- PASS: TestCommonDimensions (0.00s)
ok github.com/signalfx/splunk-otel-collector/internal/otelcol 0.010s`

type workflowRecordBuilderTestProvider struct {
	jobs []WorkflowJob
	logs WorkflowAttemptLogs
}

func (provider *workflowRecordBuilderTestProvider) Provider() testresults.Provider {
	return testresults.ProviderGitHub
}

func (provider *workflowRecordBuilderTestProvider) Repository() common.Repository {
	return common.Repository{Owner: workflowRecordBuilderTestRepositoryOwner, Name: workflowRecordBuilderTestRepositoryName}
}

func (provider *workflowRecordBuilderTestProvider) ListWorkflowRunAttempts(context.Context, time.Time, map[common.WorkflowRunAttemptKey]struct{}) ([]WorkflowRunAttempt, error) {
	return nil, nil
}

func (provider *workflowRecordBuilderTestProvider) ListWorkflowJobs(context.Context, string, int) ([]WorkflowJob, error) {
	return provider.jobs, nil
}

func (provider *workflowRecordBuilderTestProvider) StepLogsForAttempt(context.Context, string, int, []WorkflowJob) (WorkflowAttemptLogs, error) {
	return provider.logs, nil
}

func TestParseWorkflowRunAttemptsKeepsMultipleJobResultsSeparate(t *testing.T) {
	testCases := []struct {
		jobs      []WorkflowJob
		logs      WorkflowAttemptLogs
		name      string
		wantByJob []expectedWorkflowJobTest
	}{
		{
			jobs: []WorkflowJob{
				workflowRecordBuilderTestJob(workflowRecordBuilderTestLinuxJobID, workflowRecordBuilderTestLinuxJob, workflowRecordBuilderTestPassConclusion),
				workflowRecordBuilderTestJob(workflowRecordBuilderTestWindowsJobID, workflowRecordBuilderTestWindowsJob, workflowRecordBuilderTestFailureConclusion),
			},
			logs: WorkflowAttemptLogs{Jobs: []WorkflowJobLogs{
				workflowRecordBuilderTestJobLog(workflowRecordBuilderTestLinuxJobID, workflowRecordBuilderTestMatrixPassLog),
				workflowRecordBuilderTestJobLog(workflowRecordBuilderTestWindowsJobID, workflowRecordBuilderTestMatrixFailureLog),
			}},
			name: "matching matrix test in successful and failing jobs",
			wantByJob: []expectedWorkflowJobTest{
				{jobName: workflowRecordBuilderTestLinuxJob, testID: workflowRecordBuilderTestPackage + workflowRecordBuilderTestSubtestSeparator + workflowRecordBuilderTestTestName, outcome: testresults.OutcomePassed},
				{jobName: workflowRecordBuilderTestWindowsJob, testID: workflowRecordBuilderTestPackage + workflowRecordBuilderTestSubtestSeparator + workflowRecordBuilderTestTestName, outcome: testresults.OutcomeFailed, failureMessage: workflowRecordBuilderTestAssertionFailureText},
			},
		},
		{
			jobs: []WorkflowJob{
				workflowRecordBuilderTestJob(workflowRecordBuilderTestCommonDimensionsJobID, workflowRecordBuilderTestCommonDimensionsJob, workflowRecordBuilderTestPassConclusion),
				workflowRecordBuilderTestJob(workflowRecordBuilderTestIISJobID, workflowRecordBuilderTestIISJob, workflowRecordBuilderTestFailureConclusion),
			},
			logs: WorkflowAttemptLogs{Jobs: []WorkflowJobLogs{
				workflowRecordBuilderTestJobLog(workflowRecordBuilderTestCommonDimensionsJobID, workflowRecordBuilderTestCommonDimensionsLog),
				workflowRecordBuilderTestJobLog(workflowRecordBuilderTestIISJobID, goTestTextLogParserWindowsIISFailureExample),
			}},
			name: "different tests in separate jobs",
			wantByJob: []expectedWorkflowJobTest{
				{jobName: workflowRecordBuilderTestCommonDimensionsJob, testID: workflowRecordBuilderTestCommonDimensionsPackage + workflowRecordBuilderTestSubtestSeparator + workflowRecordBuilderTestCommonDimensionsTestName, outcome: testresults.OutcomePassed},
				{jobName: workflowRecordBuilderTestIISJob, testID: workflowRecordBuilderTestWindowsIISPackage + workflowRecordBuilderTestSubtestSeparator + workflowRecordBuilderTestIISName, outcome: testresults.OutcomeFailed, failureMessage: workflowRecordBuilderTestIISFailureMessage},
			},
		},
		{
			jobs: []WorkflowJob{
				workflowRecordBuilderTestJob(workflowRecordBuilderTestWindowsJobID, workflowRecordBuilderTestWindowsJob, workflowRecordBuilderTestFailureConclusion),
			},
			logs: WorkflowAttemptLogs{Jobs: []WorkflowJobLogs{
				workflowRecordBuilderTestJobLog(workflowRecordBuilderTestWindowsJobID, goTestTextLogParserPackageUpgradeTimeoutExample),
			}},
			name: "timeout panic is attached to its active matrix subtest",
			wantByJob: []expectedWorkflowJobTest{
				{jobName: workflowRecordBuilderTestWindowsJob, testID: goTestTextLogParserWantPackageUpgradeTimeoutID, outcome: testresults.OutcomeFailed, failureMessage: goTestTextLogParserPackageUpgradeTimeoutText},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			provider := &workflowRecordBuilderTestProvider{jobs: testCase.jobs, logs: testCase.logs}
			workflowRuns, err := ParseWorkflowRunAttempts(context.Background(), provider, []WorkflowRunAttempt{{
				RunID:        workflowRecordBuilderTestRunID,
				Attempt:      1,
				Name:         workflowRecordBuilderTestWorkflowName,
				Branch:       workflowRecordBuilderTestMainBranch,
				HeadSHA:      workflowRecordBuilderTestHeadSHA,
				RunStartedAt: time.Date(2026, time.November, 13, 19, 0, 0, 0, time.UTC),
				Conclusion:   workflowRecordBuilderTestFailureConclusion,
			}})
			if err != nil {
				t.Fatalf("ParseWorkflowRunAttempts() error = %v", err)
			}
			if len(workflowRuns) != 1 {
				t.Fatalf("ParseWorkflowRunAttempts() returned %d workflows, want 1", len(workflowRuns))
			}
			workflow := &workflowRuns[0]
			if len(workflow.Jobs) != len(testCase.jobs) {
				t.Fatalf("workflow has %d jobs, want %d", len(workflow.Jobs), len(testCase.jobs))
			}
			for _, expected := range testCase.wantByJob {
				parsedTest, found := workflowRecordBuilderTestForJob(workflow, expected.jobName, expected.testID)
				if !found {
					t.Errorf("test %q was not found in job %q", expected.testID, expected.jobName)
					continue
				}
				if parsedTest.Outcome != expected.outcome {
					t.Errorf("test %q in job %q has outcome %q, want %q", expected.testID, expected.jobName, parsedTest.Outcome, expected.outcome)
				}
				if expected.failureMessage != "" && !strings.Contains(parsedTest.FailureMessage, expected.failureMessage) {
					t.Errorf("test %q in job %q has failure message %q, want it to contain %q", expected.testID, expected.jobName, parsedTest.FailureMessage, expected.failureMessage)
				}
			}
		})
	}
}

type expectedWorkflowJobTest struct {
	jobName        string
	testID         string
	outcome        testresults.Outcome
	failureMessage string
}

func workflowRecordBuilderTestJob(id string, name string, conclusion string) WorkflowJob {
	return WorkflowJob{
		ID:         id,
		Name:       name,
		Conclusion: conclusion,
		JobURL:     workflowRecordBuilderTestJobURLPrefix + workflowRecordBuilderTestRunID + workflowRecordBuilderTestJobURLSeparator + id,
		Steps: []WorkflowStep{{
			Number:     1,
			Name:       workflowRecordBuilderTestStepName,
			Conclusion: conclusion,
		}},
	}
}

func workflowRecordBuilderTestJobLog(jobID string, text string) WorkflowJobLogs {
	return WorkflowJobLogs{JobID: jobID, Steps: []WorkflowStepLog{{Number: 1, Text: text}}}
}

func workflowRecordBuilderTestForJob(workflow *testresults.Workflow, jobName string, testID string) (testresults.Test, bool) {
	for _, job := range workflow.Jobs {
		if job.Name != jobName {
			continue
		}
		for _, step := range job.Steps {
			for _, test := range step.Tests {
				if test.Kind == testresults.TestKindTest && test.ID == testID {
					return test, true
				}
			}
		}
	}
	return testresults.Test{}, false
}
