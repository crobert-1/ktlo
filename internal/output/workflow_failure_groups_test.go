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

package output

import (
	"testing"
	"time"

	"ktlo/internal/testresults"
)

const (
	workflowFailureGroupsTestCaseIncompleteJobDoesNotAffectHistory  = "incomplete job does not affect history"
	workflowFailureGroupsTestCaseIncompleteMatrixVariantIsIgnored   = "incomplete matrix variant is ignored"
	workflowFailureGroupsTestCaseIncompleteWorkflowDoesNotSeedRow   = "incomplete workflow does not seed a workflow row"
	workflowFailureGroupsTestCaseIncompleteWorkflowOnlyRowIsIgnored = "incomplete workflow-only row is ignored"
	workflowFailureGroupsTestJobBuildLinux                          = "build (linux)"
	workflowFailureGroupsTestJobBuildMacOS                          = "build (macOS)"
	workflowFailureGroupsTestJobBuildMatrixDisplayName              = "build (2 matrix variants)"
	workflowFailureGroupsTestJobBuildWindows                        = "build (windows)"
	workflowFailureGroupsTestJobScan                                = "anchore-fs-scan"
	workflowFailureGroupsTestRunCompleteFailure                     = "run-complete-failure"
	workflowFailureGroupsTestRunCompletePass                        = "run-complete-pass"
	workflowFailureGroupsTestRunIncomplete                          = "run-incomplete"
	workflowFailureGroupsTestRunIncompletePass                      = "run-incomplete-pass"
	workflowFailureGroupsTestRunMatrix                              = "run-matrix"
	workflowFailureGroupsTestStepBuild                              = "Build"
	workflowFailureGroupsTestStepScan                               = "Run anchore scan"
	workflowFailureGroupsTestWorkflowName                           = "vulnerability scans"
)

func TestBuildWorkflowFailureGroupsIgnoresIncompleteWorkflows(t *testing.T) {
	tests := []struct {
		name            string
		workflows       []testresults.Workflow
		wantGroups      int
		wantJobName     string
		wantStepName    string
		wantFailures    int
		wantPasses      int
		wantLatestRunID string
	}{
		{
			name: workflowFailureGroupsTestCaseIncompleteJobDoesNotAffectHistory,
			workflows: []testresults.Workflow{
				workflowFailureGroupsTestWorkflow(workflowFailureGroupsTestRunCompleteFailure, 1, true, testresults.OutcomeFailed,
					workflowFailureGroupsTestJob(workflowFailureGroupsTestJobScan, testresults.OutcomeFailed, workflowFailureGroupsTestStepScan, testresults.OutcomeFailed)),
				workflowFailureGroupsTestWorkflow(workflowFailureGroupsTestRunCompletePass, 2, true, testresults.OutcomePassed,
					workflowFailureGroupsTestJob(workflowFailureGroupsTestJobScan, testresults.OutcomePassed, workflowFailureGroupsTestStepScan, testresults.OutcomePassed)),
				workflowFailureGroupsTestWorkflow(workflowFailureGroupsTestRunIncomplete, 3, false, testresults.OutcomeFailed,
					workflowFailureGroupsTestJob(workflowFailureGroupsTestJobScan, testresults.OutcomeFailed, workflowFailureGroupsTestStepScan, testresults.OutcomeFailed)),
				workflowFailureGroupsTestWorkflow(workflowFailureGroupsTestRunIncompletePass, 4, false, testresults.OutcomePassed,
					workflowFailureGroupsTestJob(workflowFailureGroupsTestJobScan, testresults.OutcomePassed, workflowFailureGroupsTestStepScan, testresults.OutcomePassed)),
			},
			wantGroups:      1,
			wantJobName:     workflowFailureGroupsTestJobScan,
			wantStepName:    workflowFailureGroupsTestStepScan,
			wantFailures:    1,
			wantPasses:      1,
			wantLatestRunID: workflowFailureGroupsTestRunCompleteFailure,
		},
		{
			name: workflowFailureGroupsTestCaseIncompleteMatrixVariantIsIgnored,
			workflows: []testresults.Workflow{
				workflowFailureGroupsTestWorkflow(workflowFailureGroupsTestRunMatrix, 1, true, testresults.OutcomeFailed,
					workflowFailureGroupsTestJob(workflowFailureGroupsTestJobBuildLinux, testresults.OutcomeFailed, workflowFailureGroupsTestStepBuild, testresults.OutcomeFailed),
					workflowFailureGroupsTestJob(workflowFailureGroupsTestJobBuildWindows, testresults.OutcomeFailed, workflowFailureGroupsTestStepBuild, testresults.OutcomeFailed)),
				workflowFailureGroupsTestWorkflow(workflowFailureGroupsTestRunIncomplete, 2, false, testresults.OutcomeFailed,
					workflowFailureGroupsTestJob(workflowFailureGroupsTestJobBuildMacOS, testresults.OutcomeFailed, workflowFailureGroupsTestStepBuild, testresults.OutcomeFailed)),
			},
			wantGroups:      1,
			wantJobName:     workflowFailureGroupsTestJobBuildMatrixDisplayName,
			wantStepName:    workflowFailureGroupsTestStepBuild,
			wantFailures:    1,
			wantLatestRunID: workflowFailureGroupsTestRunMatrix,
		},
		{
			name: workflowFailureGroupsTestCaseIncompleteWorkflowDoesNotSeedRow,
			workflows: []testresults.Workflow{
				workflowFailureGroupsTestWorkflow(workflowFailureGroupsTestRunIncomplete, 1, false, testresults.OutcomeFailed,
					workflowFailureGroupsTestJob(workflowFailureGroupsTestJobScan, testresults.OutcomeFailed, workflowFailureGroupsTestStepScan, testresults.OutcomeFailed)),
			},
		},
		{
			name: workflowFailureGroupsTestCaseIncompleteWorkflowOnlyRowIsIgnored,
			workflows: []testresults.Workflow{
				workflowFailureGroupsTestWorkflow(workflowFailureGroupsTestRunCompleteFailure, 1, true, testresults.OutcomeFailed,
					workflowFailureGroupsTestJob(workflowFailureGroupsTestJobScan, testresults.OutcomeFailed, workflowFailureGroupsTestStepScan, testresults.OutcomeFailed)),
				workflowFailureGroupsTestWorkflow(workflowFailureGroupsTestRunIncomplete, 2, false, testresults.OutcomeFailed),
			},
			wantGroups:      1,
			wantJobName:     workflowFailureGroupsTestJobScan,
			wantStepName:    workflowFailureGroupsTestStepScan,
			wantFailures:    1,
			wantLatestRunID: workflowFailureGroupsTestRunCompleteFailure,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			report := &testresults.Report{Workflows: test.workflows}
			got := buildWorkflowFailureGroups(report)
			if len(got) != test.wantGroups {
				t.Fatalf("got %d workflow failure groups, want %d", len(got), test.wantGroups)
			}
			if test.wantGroups == 0 {
				return
			}
			if got[0].jobDisplayName() != test.wantJobName {
				t.Errorf("got job name %q, want %q", got[0].jobDisplayName(), test.wantJobName)
			}
			if got[0].stepDisplayName() != test.wantStepName {
				t.Errorf("got step name %q, want %q", got[0].stepDisplayName(), test.wantStepName)
			}
			if got[0].counts.Failures != test.wantFailures {
				t.Errorf("got %d failures, want %d", got[0].counts.Failures, test.wantFailures)
			}
			if got[0].counts.Passes != test.wantPasses {
				t.Errorf("got %d passes, want %d", got[0].counts.Passes, test.wantPasses)
			}
			if got[0].latestFailure == nil {
				t.Fatal("latest failure is nil")
			}
			if got[0].latestFailure.SourceRun.RunID != test.wantLatestRunID {
				t.Errorf("got latest failure run ID %q, want %q", got[0].latestFailure.SourceRun.RunID, test.wantLatestRunID)
			}
		})
	}
}

func workflowFailureGroupsTestWorkflow(runID string, day int, collectionComplete bool, outcome testresults.Outcome, jobs ...testresults.Job) testresults.Workflow {
	return testresults.Workflow{
		Name:               workflowFailureGroupsTestWorkflowName,
		RunID:              runID,
		RunAttempt:         1,
		RunStartedAt:       time.Date(2026, time.November, day, 12, 0, 0, 0, time.UTC),
		Branch:             branchMain,
		Outcome:            outcome,
		CollectionComplete: collectionComplete,
		Jobs:               jobs,
	}
}

func workflowFailureGroupsTestJob(name string, outcome testresults.Outcome, stepName string, stepOutcome testresults.Outcome) testresults.Job {
	return testresults.Job{
		Name:    name,
		Outcome: outcome,
		Steps: []testresults.Step{{
			Name:    stepName,
			Outcome: stepOutcome,
		}},
	}
}
