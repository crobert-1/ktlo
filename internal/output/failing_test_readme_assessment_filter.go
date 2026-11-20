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
	"strings"

	"ktlo/internal/testresults"
)

const failingTestIDSubtestSeparator = "/"

type failingTestReadmeAssessmentFilter struct {
	failedObservationsByWorkflow map[string][]workflowTestFailureObservation
}

type workflowTestFailureObservation struct {
	scope []testresults.Test
	test  testresults.Test
}

func newFailingTestReadmeAssessmentFilter(report *testresults.Report) *failingTestReadmeAssessmentFilter {
	filter := &failingTestReadmeAssessmentFilter{
		failedObservationsByWorkflow: make(map[string][]workflowTestFailureObservation),
	}
	if report == nil {
		return filter
	}
	for workflowIndex := range report.Workflows {
		workflow := &report.Workflows[workflowIndex]
		if workflow.Branch != branchMain {
			continue
		}
		scopes := make([][]testresults.Test, 0)
		if len(workflow.Tests) > 0 {
			scopes = append(scopes, workflow.Tests)
		}
		for jobIndex := range workflow.Jobs {
			job := &workflow.Jobs[jobIndex]
			if len(job.Tests) > 0 {
				scopes = append(scopes, job.Tests)
			}
			for stepIndex := range job.Steps {
				step := &job.Steps[stepIndex]
				if len(step.Tests) > 0 {
					scopes = append(scopes, step.Tests)
				}
			}
		}
		for _, scope := range scopes {
			for _, test := range scope {
				if test.Kind != testresults.TestKindTest || test.Outcome != testresults.OutcomeFailed {
					continue
				}
				filter.failedObservationsByWorkflow[workflow.Name] = append(
					filter.failedObservationsByWorkflow[workflow.Name],
					workflowTestFailureObservation{scope: scope, test: test},
				)
			}
		}
	}
	return filter
}

func (filter *failingTestReadmeAssessmentFilter) suppressExplainedParentTests(assessments []testresults.FailureAssessment) []testresults.FailureAssessment {
	filteredAssessments := make([]testresults.FailureAssessment, 0, len(assessments))
	for assessmentIndex := range assessments {
		assessment := &assessments[assessmentIndex]
		if !filter.isExplainedByFailedDescendants(assessment) {
			filteredAssessments = append(filteredAssessments, *assessment)
		}
	}
	return filteredAssessments
}

func (filter *failingTestReadmeAssessmentFilter) isExplainedByFailedDescendants(assessment *testresults.FailureAssessment) bool {
	location := assessment.Location
	if location.TestKind != testresults.TestKindTest || location.TestID == emptyString {
		return false
	}
	failedObservations := filter.failedObservationsByWorkflow[location.WorkflowName]
	var parentFailureFound bool
	for _, observation := range failedObservations {
		if observation.test.ID != location.TestID {
			continue
		}
		parentFailureFound = true
		if !observation.hasExplainingFailedDescendant() {
			return false
		}
	}
	return parentFailureFound
}

func (observation *workflowTestFailureObservation) hasExplainingFailedDescendant() bool {
	parentFailureMessage := strings.TrimSpace(observation.test.FailureMessage)
	for _, candidate := range observation.scope {
		if candidate.Kind != testresults.TestKindTest || candidate.Outcome != testresults.OutcomeFailed ||
			!strings.HasPrefix(candidate.ID, observation.test.ID+failingTestIDSubtestSeparator) {
			continue
		}
		if parentFailureMessage == emptyString || strings.Contains(strings.TrimSpace(candidate.FailureMessage), parentFailureMessage) {
			return true
		}
	}
	return false
}
