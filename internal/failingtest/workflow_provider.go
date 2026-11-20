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
	"time"

	"ktlo/internal/common"
	"ktlo/internal/testresults"
)

const githubLogGroupPostPrefix = "##[group]Post "
const githubLogGroupPrefix = "##[group]"
const githubLogGroupRunPrefix = "##[group]Run "
const githubPostJobCleanupLine = "Post job cleanup."
const githubPostStepNamePrefix = "Post "
const workflowLogLineTimestampPrecision = time.Second

type WorkflowProvider interface {
	Provider() testresults.Provider
	Repository() common.Repository
	ListWorkflowRunAttempts(ctx context.Context, startedAfter time.Time, knownAttempts map[common.WorkflowRunAttemptKey]struct{}) ([]WorkflowRunAttempt, error)
	ListWorkflowJobs(ctx context.Context, runID string, attempt int) ([]WorkflowJob, error)
	StepLogsForAttempt(ctx context.Context, runID string, attempt int, jobs []WorkflowJob) (WorkflowAttemptLogs, error)
}

type LatestWorkflowRunAttemptProvider interface {
	LatestWorkflowRunAttempt(ctx context.Context, runID string) (WorkflowRunAttempt, error)
}

type WorkflowAttemptLogs struct {
	Jobs []WorkflowJobLogs
}

func (logs *WorkflowAttemptLogs) ForJobID(jobID string) (WorkflowJobLogs, bool) {
	for _, jobLogs := range logs.Jobs {
		if jobLogs.JobID == jobID {
			return jobLogs, true
		}
	}
	return WorkflowJobLogs{}, false
}

type WorkflowJobLogs struct {
	JobID     string
	JobOutput string
	Steps     []WorkflowStepLog
}

func (logs *WorkflowJobLogs) splitCombinedOutput(workflowJob WorkflowJob) {
	if logs.JobOutput == "" {
		return
	}
	existingStepNumbers := make(map[int]struct{}, len(logs.Steps))
	for _, stepLog := range logs.Steps {
		existingStepNumbers[stepLog.Number] = struct{}{}
	}
	stepLinesByNumber := make(map[int][]string)
	currentStepNumber := 0
	for _, line := range strings.Split(logs.JobOutput, logLineSeparator) {
		lineTimestamp, hasTimestamp := workflowLogLineTimestamp(line)
		if hasTimestamp {
			matchingStepNumbers := matchingWorkflowStepNumbers(lineTimestamp, workflowJob.Steps)
			currentStepNumber = workflowStepNumberForLogLine(line, workflowJob.Steps, matchingStepNumbers, currentStepNumber)
			if currentStepNumber == 0 && len(matchingStepNumbers) > 0 {
				// With second-precision metadata, prefer the earlier step for an ambiguous boundary line.
				currentStepNumber = matchingStepNumbers[0]
			}
		}
		if currentStepNumber != 0 {
			stepLinesByNumber[currentStepNumber] = append(stepLinesByNumber[currentStepNumber], line)
		}
	}
	for _, workflowStep := range workflowJob.Steps {
		if _, exists := existingStepNumbers[workflowStep.Number]; exists {
			continue
		}
		stepLines := stepLinesByNumber[workflowStep.Number]
		if len(stepLines) == 0 {
			continue
		}
		logs.Steps = append(logs.Steps, WorkflowStepLog{
			Number: workflowStep.Number,
			Text:   strings.Join(stepLines, logLineSeparator),
		})
	}
}

func matchingWorkflowStepNumbers(lineTimestamp time.Time, workflowSteps []WorkflowStep) []int {
	stepNumbers := make([]int, 0)
	for _, workflowStep := range workflowSteps {
		if workflowStep.StartedAt.IsZero() || workflowLogTimestampBeforeStep(lineTimestamp, workflowStep.StartedAt) ||
			(!workflowStep.CompletedAt.IsZero() && workflowLogTimestampAfterStep(lineTimestamp, workflowStep.CompletedAt)) {
			continue
		}
		stepNumbers = append(stepNumbers, workflowStep.Number)
	}
	return stepNumbers
}

func workflowLogTimestampBeforeStep(lineTimestamp time.Time, stepTimestamp time.Time) bool {
	if stepTimestamp.Nanosecond() == 0 {
		return lineTimestamp.Truncate(workflowLogLineTimestampPrecision).Before(stepTimestamp)
	}
	return lineTimestamp.Before(stepTimestamp)
}

func workflowLogTimestampAfterStep(lineTimestamp time.Time, stepTimestamp time.Time) bool {
	if stepTimestamp.Nanosecond() == 0 {
		return lineTimestamp.Truncate(workflowLogLineTimestampPrecision).After(stepTimestamp)
	}
	return lineTimestamp.After(stepTimestamp)
}

func workflowStepNumberForLogLine(line string, workflowSteps []WorkflowStep, matchingStepNumbers []int, currentStepNumber int) int {
	lineContent := workflowLogLineContent(line)
	if lineContent == githubPostJobCleanupLine {
		for _, workflowStep := range workflowSteps {
			if workflowStep.Number > currentStepNumber && strings.HasPrefix(workflowStep.Name, githubPostStepNamePrefix) &&
				containsWorkflowStepNumber(matchingStepNumbers, workflowStep.Number) {
				return workflowStep.Number
			}
		}
	}
	for _, groupPrefix := range []string{githubLogGroupRunPrefix, githubLogGroupPostPrefix} {
		if !strings.HasPrefix(lineContent, groupPrefix) {
			continue
		}
		stepLabel := strings.TrimPrefix(lineContent, githubLogGroupPrefix)
		for _, workflowStep := range workflowSteps {
			if strings.EqualFold(workflowStep.Name, stepLabel) && containsWorkflowStepNumber(matchingStepNumbers, workflowStep.Number) {
				return workflowStep.Number
			}
		}
	}
	if containsWorkflowStepNumber(matchingStepNumbers, currentStepNumber) {
		return currentStepNumber
	}
	return 0
}

func workflowLogLineContent(line string) string {
	if timestampEnd := strings.Index(line, workflowLogTimestampSeparator); timestampEnd >= 19 && timestampEnd <= 40 {
		return line[timestampEnd+2:]
	}
	return line
}

func containsWorkflowStepNumber(stepNumbers []int, stepNumber int) bool {
	for _, candidateStepNumber := range stepNumbers {
		if candidateStepNumber == stepNumber {
			return true
		}
	}
	return false
}

type WorkflowStepLog struct {
	Number int
	Text   string
}

func workflowLogLineTimestamp(line string) (time.Time, bool) {
	timestampEnd := strings.Index(line, workflowLogTimestampSeparator)
	if timestampEnd < 19 || timestampEnd > 40 {
		return time.Time{}, false
	}
	// The prefix through Z is the RFC3339 timestamp token.
	timestamp, err := time.Parse(time.RFC3339Nano, line[:timestampEnd+1])
	return timestamp, err == nil
}

const workflowLogTimestampSeparator = "Z "

type WorkflowRunAttempt = common.WorkflowRunAttempt
type WorkflowJob = common.WorkflowJob
type WorkflowStep = common.WorkflowStep
