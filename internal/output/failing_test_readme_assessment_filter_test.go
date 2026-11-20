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

	"ktlo/internal/testresults"
)

const (
	failingTestFilterTestCaseEmptyParentMessage          = "empty parent message with failed descendant"
	failingTestFilterTestCaseParentMessageInChild        = "parent message repeated in failed descendant"
	failingTestFilterTestCaseParentHasIndependentMessage = "parent has independent failure message"
	failingTestFilterTestCaseDescendantInDifferentStep   = "descendant is in a different step"
	failingTestFilterTestCaseParentFailsInAnotherRun     = "parent has an unexplained failure in another run"
	failingTestFilterTestCaseDescendantPassed            = "descendant passed"
	failingTestFilterWorkflowName                        = "functional-tests"
	failingTestFilterParentID                            = "configuration_switching.Test_Functions"
	failingTestFilterChildID                             = "configuration_switching.Test_Functions/logs_and_metrics_attributes_verification"
	failingTestFilterRunID                               = "run-1"
	failingTestFilterSecondRunID                         = "run-2"
	failingTestFilterParentMessage                       = "parent setup failed"
	failingTestFilterChildMessage                        = "parent setup failed: condition never satisfied"
	failingTestFilterIndependentMessage                  = "teardown failed"
	failingTestFilterParentStepName                      = "parent-step"
	failingTestFilterChildStepName                       = "child-step"
)

func TestFailingTestReadmeAssessmentFilterSuppressesExplainedParents(t *testing.T) {
	parentAssessment := testresults.FailureAssessment{
		Location: testresults.CIFailureLocation{
			WorkflowName: failingTestFilterWorkflowName,
			TestID:       failingTestFilterParentID,
			TestKind:     testresults.TestKindTest,
		},
	}
	tests := []struct {
		name      string
		workflows []testresults.Workflow
		wantCount int
	}{
		{
			name: failingTestFilterTestCaseEmptyParentMessage,
			workflows: []testresults.Workflow{
				failingTestFilterWorkflow(failingTestFilterRunID, []testresults.Test{
					failingTestFilterTest(failingTestFilterParentID, testresults.OutcomeFailed, emptyString),
					failingTestFilterTest(failingTestFilterChildID, testresults.OutcomeFailed, failingTestFilterChildMessage),
				}),
			},
			wantCount: 0,
		},
		{
			name: failingTestFilterTestCaseParentMessageInChild,
			workflows: []testresults.Workflow{
				failingTestFilterWorkflow(failingTestFilterRunID, []testresults.Test{
					failingTestFilterTest(failingTestFilterParentID, testresults.OutcomeFailed, failingTestFilterParentMessage),
					failingTestFilterTest(failingTestFilterChildID, testresults.OutcomeFailed, failingTestFilterChildMessage),
				}),
			},
			wantCount: 0,
		},
		{
			name: failingTestFilterTestCaseParentHasIndependentMessage,
			workflows: []testresults.Workflow{
				failingTestFilterWorkflow(failingTestFilterRunID, []testresults.Test{
					failingTestFilterTest(failingTestFilterParentID, testresults.OutcomeFailed, failingTestFilterIndependentMessage),
					failingTestFilterTest(failingTestFilterChildID, testresults.OutcomeFailed, failingTestFilterChildMessage),
				}),
			},
			wantCount: 1,
		},
		{
			name: failingTestFilterTestCaseDescendantInDifferentStep,
			workflows: []testresults.Workflow{
				{
					Name:   failingTestFilterWorkflowName,
					Branch: branchMain,
					RunID:  failingTestFilterRunID,
					Jobs: []testresults.Job{{
						Name:    failingTestFilterWorkflowName,
						Outcome: testresults.OutcomeFailed,
						Steps: []testresults.Step{
							{
								Name:    failingTestFilterParentStepName,
								Outcome: testresults.OutcomeFailed,
								Tests:   []testresults.Test{failingTestFilterTest(failingTestFilterParentID, testresults.OutcomeFailed, emptyString)},
							},
							{
								Name:    failingTestFilterChildStepName,
								Outcome: testresults.OutcomeFailed,
								Tests:   []testresults.Test{failingTestFilterTest(failingTestFilterChildID, testresults.OutcomeFailed, failingTestFilterChildMessage)},
							},
						},
					}},
				},
			},
			wantCount: 1,
		},
		{
			name: failingTestFilterTestCaseParentFailsInAnotherRun,
			workflows: []testresults.Workflow{
				failingTestFilterWorkflow(failingTestFilterRunID, []testresults.Test{
					failingTestFilterTest(failingTestFilterParentID, testresults.OutcomeFailed, emptyString),
					failingTestFilterTest(failingTestFilterChildID, testresults.OutcomeFailed, failingTestFilterChildMessage),
				}),
				failingTestFilterWorkflow(failingTestFilterSecondRunID, []testresults.Test{
					failingTestFilterTest(failingTestFilterParentID, testresults.OutcomeFailed, failingTestFilterIndependentMessage),
				}),
			},
			wantCount: 1,
		},
		{
			name: failingTestFilterTestCaseDescendantPassed,
			workflows: []testresults.Workflow{
				failingTestFilterWorkflow(failingTestFilterRunID, []testresults.Test{
					failingTestFilterTest(failingTestFilterParentID, testresults.OutcomeFailed, emptyString),
					failingTestFilterTest(failingTestFilterChildID, testresults.OutcomePassed, emptyString),
				}),
			},
			wantCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			report := &testresults.Report{Workflows: test.workflows}
			filter := newFailingTestReadmeAssessmentFilter(report)
			got := filter.suppressExplainedParentTests([]testresults.FailureAssessment{parentAssessment})
			if len(got) != test.wantCount {
				t.Fatalf("got %d assessments, want %d", len(got), test.wantCount)
			}
		})
	}
}

func failingTestFilterWorkflow(runID string, tests []testresults.Test) testresults.Workflow {
	return testresults.Workflow{
		Name:   failingTestFilterWorkflowName,
		Branch: branchMain,
		RunID:  runID,
		Jobs: []testresults.Job{{
			Name:    failingTestFilterWorkflowName,
			Outcome: testresults.OutcomeFailed,
			Steps: []testresults.Step{{
				Name:    failingTestFilterWorkflowName,
				Outcome: testresults.OutcomeFailed,
				Tests:   tests,
			}},
		}},
	}
}

func failingTestFilterTest(testID string, outcome testresults.Outcome, failureMessage string) testresults.Test {
	return testresults.Test{
		ID:             testID,
		Kind:           testresults.TestKindTest,
		Outcome:        outcome,
		FailureMessage: failureMessage,
	}
}
