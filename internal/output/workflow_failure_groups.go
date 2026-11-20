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
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"ktlo/internal/testresults"
)

const (
	workflowFailureGroupKeySeparator         = "\x00"
	workflowFailureMatrixDelimiter           = " ("
	workflowFailureMatrixSuffix              = ")"
	workflowFailureSignatureJob              = "job"
	workflowFailureSignatureWorkflow         = "workflow"
	workflowFailureStepSeparator             = "; "
	workflowMatrixVariantCountPluralFormat   = " (%d matrix variants)"
	workflowMatrixVariantCountSingularFormat = " (%d matrix variant)"
)

type workflowFailureGroup struct {
	affectedVariants map[string]struct{}
	classification   testresults.Classification
	counts           testresults.OutcomeCounts
	failureSignature string
	isMatrixFamily   bool
	jobFamily        string
	latestFailure    *testresults.FailureEvidence
	stepNames        []string
	workflowName     string
}

type workflowFailureGroupKey struct {
	jobFamily    string
	signature    string
	workflowName string
}

type workflowJobFamilyIndex map[string]map[string]map[string]struct{}

func buildWorkflowFailureGroups(report *testresults.Report) []workflowFailureGroup {
	familyIndex := buildWorkflowJobFamilyIndex(report)
	groups := make(map[workflowFailureGroupKey]*workflowFailureGroup)
	for workflowIndex := range report.Workflows {
		workflow := &report.Workflows[workflowIndex]
		if workflow.Branch != branchMain || !workflow.CollectionComplete {
			continue
		}
		failedJobFound := false
		for jobIndex := range workflow.Jobs {
			job := &workflow.Jobs[jobIndex]
			failedSteps := failedWorkflowSteps(job)
			if len(failedSteps) == 0 && job.Outcome != testresults.OutcomeFailed {
				continue
			}
			failedJobFound = true
			jobFamily, variant := workflowJobFamily(workflow.Name, job.Name, familyIndex)
			groupKey := workflowFailureGroupKey{
				workflowName: workflow.Name,
				jobFamily:    jobFamily,
				signature:    workflowJobFailureSignature(failedSteps),
			}
			group := groups[groupKey]
			if group == nil {
				group = &workflowFailureGroup{
					affectedVariants: make(map[string]struct{}),
					failureSignature: groupKey.signature,
					jobFamily:        jobFamily,
					isMatrixFamily:   len(familyIndex[workflow.Name][jobFamily]) > 1,
					stepNames:        workflowFailureStepNames(failedSteps),
					workflowName:     workflow.Name,
				}
				groups[groupKey] = group
			}
			group.affectedVariants[variant] = struct{}{}
		}
		if workflow.Outcome == testresults.OutcomeFailed && !failedJobFound {
			groupKey := workflowFailureGroupKey{workflowName: workflow.Name, signature: workflowFailureSignatureWorkflow}
			if groups[groupKey] == nil {
				groups[groupKey] = &workflowFailureGroup{workflowName: workflow.Name}
			}
		}
	}
	for groupKey, group := range groups {
		group.aggregateHistory(report, familyIndex, groupKey)
	}
	workflowFailures := make([]workflowFailureGroup, 0, len(groups))
	for _, group := range groups {
		group.classify()
		workflowFailures = append(workflowFailures, *group)
	}
	sort.Slice(workflowFailures, func(firstIndex int, secondIndex int) bool {
		firstGroup := workflowFailures[firstIndex]
		secondGroup := workflowFailures[secondIndex]
		if firstGroup.workflowName != secondGroup.workflowName {
			return firstGroup.workflowName < secondGroup.workflowName
		}
		if firstGroup.jobFamily != secondGroup.jobFamily {
			return firstGroup.jobFamily < secondGroup.jobFamily
		}
		firstSteps := strings.Join(firstGroup.stepNames, workflowFailureStepSeparator)
		secondSteps := strings.Join(secondGroup.stepNames, workflowFailureStepSeparator)
		if firstSteps != secondSteps {
			return firstSteps < secondSteps
		}
		return firstGroup.failureSignature < secondGroup.failureSignature
	})
	return workflowFailures
}

func buildWorkflowJobFamilyIndex(report *testresults.Report) workflowJobFamilyIndex {
	familyIndex := make(workflowJobFamilyIndex)
	for workflowIndex := range report.Workflows {
		workflow := &report.Workflows[workflowIndex]
		if !workflow.CollectionComplete {
			continue
		}
		if familyIndex[workflow.Name] == nil {
			familyIndex[workflow.Name] = make(map[string]map[string]struct{})
		}
		for jobIndex := range workflow.Jobs {
			jobName := workflow.Jobs[jobIndex].Name
			baseName, _, isMatrixName := splitWorkflowMatrixJobName(jobName)
			if !isMatrixName {
				continue
			}
			if familyIndex[workflow.Name][baseName] == nil {
				familyIndex[workflow.Name][baseName] = make(map[string]struct{})
			}
			familyIndex[workflow.Name][baseName][jobName] = struct{}{}
		}
	}
	return familyIndex
}

func (group *workflowFailureGroup) aggregateHistory(report *testresults.Report, familyIndex workflowJobFamilyIndex, groupKey workflowFailureGroupKey) {
	seenAttempts := make(map[string]struct{})
	for workflowIndex := range report.Workflows {
		workflow := &report.Workflows[workflowIndex]
		if !workflow.CollectionComplete || workflow.Branch != branchMain || workflow.Name != groupKey.workflowName {
			continue
		}
		attemptKey := workflowRunAttemptKey(workflow.RunID, workflow.RunAttempt)
		if _, seen := seenAttempts[attemptKey]; seen {
			continue
		}
		seenAttempts[attemptKey] = struct{}{}
		if groupKey.jobFamily == emptyString {
			group.addWorkflowAttempt(workflow)
			continue
		}
		familyJobs := workflowJobsInFamily(workflow, groupKey.jobFamily, familyIndex)
		if len(familyJobs) == 0 {
			continue
		}
		matchingFailureJobs := make([]*testresults.Job, 0)
		for jobIndex := range familyJobs {
			job := familyJobs[jobIndex]
			if jobMatchesWorkflowFailureGroup(job, groupKey.signature, group.stepNames) {
				matchingFailureJobs = append(matchingFailureJobs, job)
			}
		}
		if len(matchingFailureJobs) > 0 {
			group.counts.Failures++
			group.updateLatestFailure(workflow, matchingFailureJobs)
			continue
		}
		if workflow.CollectionComplete && group.familyPassed(familyJobs, familyIndex, workflow.Name) {
			group.counts.Passes++
		}
	}
}

func (group *workflowFailureGroup) addWorkflowAttempt(workflow *testresults.Workflow) {
	if workflow.Outcome == testresults.OutcomeFailed {
		group.counts.Failures++
		group.updateWorkflowFailure(workflow)
		return
	}
	if workflow.CollectionComplete && workflow.Outcome == testresults.OutcomePassed {
		group.counts.Passes++
	}
}

func (group *workflowFailureGroup) familyPassed(familyJobs []*testresults.Job, familyIndex workflowJobFamilyIndex, workflowName string) bool {
	observedVariants := make(map[string]struct{}, len(familyJobs))
	for _, job := range familyJobs {
		if job.Outcome != testresults.OutcomePassed {
			return false
		}
		observedVariants[workflowJobVariant(workflowName, job.Name, familyIndex)] = struct{}{}
		for _, expectedStepName := range group.stepNames {
			stepFound := false
			for _, step := range job.Steps {
				if step.Name != expectedStepName {
					continue
				}
				stepFound = step.Outcome == testresults.OutcomePassed
				break
			}
			if !stepFound {
				return false
			}
		}
	}
	for variant := range group.affectedVariants {
		if _, observed := observedVariants[variant]; !observed {
			return false
		}
	}
	return len(observedVariants) > 0
}

func (group *workflowFailureGroup) updateLatestFailure(workflow *testresults.Workflow, failedJobs []*testresults.Job) {
	for _, job := range failedJobs {
		message := workflowFailureMessage(job)
		reference := testresults.RunReference{
			RunID:        workflow.RunID,
			RunAttempt:   workflow.RunAttempt,
			Branch:       workflow.Branch,
			CommitSHA:    workflow.CommitSHA,
			RunStartedAt: workflow.RunStartedAt,
			WorkflowURL:  workflow.WorkflowURL,
			JobURL:       job.JobURL,
		}
		if group.latestFailure != nil && !workflowRunIsLater(&reference, &group.latestFailure.SourceRun) {
			continue
		}
		group.latestFailure = &testresults.FailureEvidence{SourceRun: reference, Message: message}
	}
}

func (group *workflowFailureGroup) updateWorkflowFailure(workflow *testresults.Workflow) {
	reference := testresults.RunReference{
		RunID:        workflow.RunID,
		RunAttempt:   workflow.RunAttempt,
		Branch:       workflow.Branch,
		CommitSHA:    workflow.CommitSHA,
		RunStartedAt: workflow.RunStartedAt,
		WorkflowURL:  workflow.WorkflowURL,
	}
	if group.latestFailure == nil || workflowRunIsLater(&reference, &group.latestFailure.SourceRun) {
		group.latestFailure = &testresults.FailureEvidence{SourceRun: reference}
	}
}

func (group *workflowFailureGroup) classify() {
	switch {
	case group.counts.Failures == 0:
		group.classification = testresults.ClassificationUnknown
	case group.counts.Passes > 0:
		group.classification = testresults.ClassificationFlaky
	default:
		group.classification = testresults.ClassificationBroken
	}
}

func (group *workflowFailureGroup) jobDisplayName() string {
	if !group.isMatrixFamily {
		return group.jobFamily
	}
	variantFormat := workflowMatrixVariantCountPluralFormat
	if len(group.affectedVariants) == 1 {
		variantFormat = workflowMatrixVariantCountSingularFormat
	}
	return group.jobFamily + fmt.Sprintf(variantFormat, len(group.affectedVariants))
}

func (group *workflowFailureGroup) stepDisplayName() string {
	return strings.Join(group.stepNames, workflowFailureStepSeparator)
}

func failedWorkflowSteps(job *testresults.Job) []testresults.Step {
	failedSteps := make([]testresults.Step, 0)
	for _, step := range job.Steps {
		if step.Outcome == testresults.OutcomeFailed {
			failedSteps = append(failedSteps, step)
		}
	}
	sort.Slice(failedSteps, func(firstIndex int, secondIndex int) bool {
		return failedSteps[firstIndex].Name < failedSteps[secondIndex].Name
	})
	return failedSteps
}

func workflowJobFailureSignature(failedSteps []testresults.Step) string {
	if len(failedSteps) == 0 {
		return workflowFailureSignatureJob
	}
	return strings.Join(workflowFailureStepNames(failedSteps), workflowFailureGroupKeySeparator)
}

func jobMatchesWorkflowFailureGroup(job *testresults.Job, signature string, stepNames []string) bool {
	failedSteps := failedWorkflowSteps(job)
	if len(stepNames) == 0 {
		return len(failedSteps) == 0 && job.Outcome == testresults.OutcomeFailed
	}
	return len(failedSteps) > 0 && workflowJobFailureSignature(failedSteps) == signature
}

func workflowFailureStepNames(failedSteps []testresults.Step) []string {
	stepNames := make([]string, 0, len(failedSteps))
	for _, step := range failedSteps {
		if len(stepNames) == 0 || stepNames[len(stepNames)-1] != step.Name {
			stepNames = append(stepNames, step.Name)
		}
	}
	return stepNames
}

func workflowFailureMessage(job *testresults.Job) string {
	for _, step := range failedWorkflowSteps(job) {
		if strings.TrimSpace(step.ReasonExcerpt) != emptyString {
			return step.ReasonExcerpt
		}
	}
	return emptyString
}

func workflowJobsInFamily(workflow *testresults.Workflow, jobFamily string, familyIndex workflowJobFamilyIndex) []*testresults.Job {
	familyJobs := make([]*testresults.Job, 0)
	for jobIndex := range workflow.Jobs {
		job := &workflow.Jobs[jobIndex]
		candidateFamily, _ := workflowJobFamily(workflow.Name, job.Name, familyIndex)
		if candidateFamily == jobFamily {
			familyJobs = append(familyJobs, job)
		}
	}
	return familyJobs
}

func workflowJobFamily(workflowName string, jobName string, familyIndex workflowJobFamilyIndex) (string, string) {
	baseName, variantName, isMatrixName := splitWorkflowMatrixJobName(jobName)
	if !isMatrixName || len(familyIndex[workflowName][baseName]) < 2 {
		return jobName, jobName
	}
	return baseName, variantName
}

func workflowJobVariant(workflowName string, jobName string, familyIndex workflowJobFamilyIndex) string {
	_, variant := workflowJobFamily(workflowName, jobName, familyIndex)
	return variant
}

func splitWorkflowMatrixJobName(jobName string) (string, string, bool) {
	if !strings.HasSuffix(jobName, workflowFailureMatrixSuffix) {
		return jobName, emptyString, false
	}
	separatorIndex := strings.LastIndex(jobName, workflowFailureMatrixDelimiter)
	if separatorIndex <= 0 {
		return jobName, emptyString, false
	}
	// The prefix is the family name, and the parenthesized suffix is the matrix variant.
	baseName := strings.TrimSpace(jobName[:separatorIndex])
	variantName := strings.TrimSpace(jobName[separatorIndex+len(workflowFailureMatrixDelimiter) : len(jobName)-len(workflowFailureMatrixSuffix)])
	if baseName == emptyString || variantName == emptyString {
		return jobName, emptyString, false
	}
	return baseName, variantName, true
}

func workflowRunIsLater(firstReference *testresults.RunReference, secondReference *testresults.RunReference) bool {
	if !firstReference.RunStartedAt.Equal(secondReference.RunStartedAt) {
		return firstReference.RunStartedAt.After(secondReference.RunStartedAt)
	}
	if firstReference.RunAttempt != secondReference.RunAttempt {
		return firstReference.RunAttempt > secondReference.RunAttempt
	}
	return firstReference.RunID > secondReference.RunID
}

func workflowRunAttemptKey(runID string, attempt int) string {
	return runID + workflowFailureGroupKeySeparator + strconv.Itoa(attempt)
}

func (group *workflowFailureGroup) observedFailureRate() string {
	observations := group.counts.Passes + group.counts.Failures
	if observations == 0 {
		return failingReadmeUnavailableRate
	}
	return fmt.Sprintf(failingReadmeRateFormat, int(math.Round(float64(group.counts.Failures)/float64(observations)*100)))
}

func (group *workflowFailureGroup) observationCountLabel() string {
	return fmt.Sprintf(failingReadmeCountFormat, group.counts.Failures, group.counts.Passes+group.counts.Failures)
}

func (group *workflowFailureGroup) latestFailureEvidence() *testresults.FailureEvidence {
	return group.latestFailure
}
