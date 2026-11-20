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
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"ktlo/internal/telemetry"
	"ktlo/internal/testresults"
)

const workflowConclusionFailure = "failure"
const workflowConclusionSkipped = "skipped"
const workflowConclusionSuccess = "success"
const workflowConclusionUnknown = "unknown"

const (
	spanBuildTestResultJobs    = "failingtest.build_test_result_jobs"
	spanBuildWorkflowRecord    = "failingtest.build_workflow_record"
	spanCollectWorkflowJobLogs = "failingtest.collect_workflow_job_logs"
	spanParseTestLog           = "failingtest.parse_test_log"
	traceAttributeJobCount     = "ktlo.workflow_job.count"
	traceAttributeLogByteCount = "ktlo.job_log.byte_count"
	traceAttributeRunAttempt   = "ktlo.workflow_run.attempt"
	traceAttributeRunID        = "ktlo.workflow_run.id"
	traceAttributeStepLogCount = "ktlo.step_log.count"
)

type workflowStepLog struct {
	stepNumber int
	text       string
}

type workflowJobBuildInput struct {
	job       WorkflowJob
	jobOutput string
	stepLogs  []workflowStepLog
}

func (jobInput *workflowJobBuildInput) logForStepNumber(stepNumber int) string {
	for _, stepLog := range jobInput.stepLogs {
		if stepLog.stepNumber == stepNumber {
			return stepLog.text
		}
	}
	return ""
}

type testResultJobBuildInput struct {
	collectionErrors []error
	jobs             []workflowJobBuildInput
}

type workflowRecordBuilder struct {
	collectionErrors   []error
	workflowProvider   WorkflowProvider
	testLogParsers     []TestLogParser
	testLogParserMutex *sync.Mutex
	workflowRun        *testresults.Workflow
	workflowRunAttempt WorkflowRunAttempt
}

func newWorkflowRecordBuilder(workflowProvider WorkflowProvider, testLogParsers []TestLogParser) *workflowRecordBuilder {
	return &workflowRecordBuilder{
		workflowProvider: workflowProvider,
		testLogParsers:   testLogParsers,
	}
}

func (builder *workflowRecordBuilder) build(ctx context.Context, workflowRunAttempt WorkflowRunAttempt) (workflow *testresults.Workflow, resultErr error) {
	return builder.buildWithWorkflowJobs(ctx, workflowRunAttempt, nil, false)
}

func (builder *workflowRecordBuilder) buildWithWorkflowJobs(ctx context.Context, workflowRunAttempt WorkflowRunAttempt, workflowJobs []WorkflowJob, useProvidedJobs bool) (workflow *testresults.Workflow, resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanBuildWorkflowRecord,
		attribute.String(traceAttributeRunID, workflowRunAttempt.RunID),
		attribute.Int(traceAttributeRunAttempt, workflowRunAttempt.Attempt),
	)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanBuildWorkflowRecord, startedAt, resultErr != nil)
	}()

	builder.collectionErrors = nil
	builder.workflowRun = generateReportWorkflowFromAttempt(workflowRunAttempt)
	builder.workflowRunAttempt = workflowRunAttempt
	builder.addPullRequest()
	jobBuildInput := builder.collectWorkflowJobLogs(ctx, workflowJobs, useProvidedJobs)
	testResultJobs, jobErrors := builder.buildTestResultJobs(ctx, jobBuildInput)
	builder.workflowRun.Jobs = testResultJobs
	builder.addErrors(jobBuildInput.collectionErrors)
	builder.addErrors(jobErrors)
	return builder.validateWorkflowRecord()
}

func generateReportWorkflowFromAttempt(workflowRunAttempt WorkflowRunAttempt) *testresults.Workflow {
	workflowRun := &testresults.Workflow{
		Name:         workflowRunAttempt.Name,
		RunID:        workflowRunAttempt.RunID,
		RunAttempt:   workflowRunAttempt.Attempt,
		RunStartedAt: workflowRunAttempt.RunStartedAt,
		Branch:       workflowRunAttempt.Branch,
		CommitSHA:    workflowRunAttempt.HeadSHA,
		WorkflowURL:  workflowRunAttempt.WorkflowURL,
	}
	workflowRun.Outcome, workflowRun.ProviderConclusion = resultOutcome(workflowRunAttempt.Conclusion, workflowRunAttempt.Status)
	return workflowRun
}

func (builder *workflowRecordBuilder) addPullRequest() {
	if len(builder.workflowRunAttempt.PullRequests) > 0 {
		pullRequest := builder.workflowRunAttempt.PullRequests[0]
		builder.workflowRun.PullRequest = &testresults.PullRequestRef{
			Number:  pullRequest.Number,
			URL:     pullRequest.URL,
			HeadSHA: pullRequest.HeadSHA,
		}
		if builder.workflowRun.PullRequest.HeadSHA == "" {
			builder.workflowRun.PullRequest.HeadSHA = builder.workflowRunAttempt.HeadSHA
		}
	}
	if len(builder.workflowRunAttempt.PullRequests) > 1 {
		builder.collectionErrors = append(builder.collectionErrors, fmt.Errorf("workflow run %s attempt %d is associated with multiple pull requests; report can store only one", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt))
	}
}

func (builder *workflowRecordBuilder) collectWorkflowJobLogs(ctx context.Context, providedWorkflowJobs []WorkflowJob, useProvidedJobs bool) (jobBuildInput testResultJobBuildInput) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanCollectWorkflowJobLogs,
		attribute.String(traceAttributeRunID, builder.workflowRunAttempt.RunID),
		attribute.Int(traceAttributeRunAttempt, builder.workflowRunAttempt.Attempt),
	)
	defer func() {
		jobCount := len(jobBuildInput.jobs)
		stepLogCount := 0
		for _, workflowJob := range jobBuildInput.jobs {
			stepLogCount += len(workflowJob.stepLogs)
		}
		span.SetAttributes(
			attribute.Int(traceAttributeJobCount, jobCount),
			attribute.Int(traceAttributeStepLogCount, stepLogCount),
		)
		telemetry.FinishSpan(ctx, span, spanCollectWorkflowJobLogs, startedAt, len(jobBuildInput.collectionErrors) > 0)
	}()
	jobBuildInput = testResultJobBuildInput{jobs: make([]workflowJobBuildInput, 0)}
	workflowJobs := providedWorkflowJobs
	workflowSucceeded := strings.EqualFold(builder.workflowRunAttempt.Conclusion, workflowConclusionSuccess)
	if !useProvidedJobs {
		var listErr error
		workflowJobs, listErr = builder.workflowProvider.ListWorkflowJobs(ctx, builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt)
		if listErr != nil {
			jobBuildInput.collectionErrors = append(jobBuildInput.collectionErrors, fmt.Errorf("could not list jobs for workflow run %s attempt %d: %w", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, listErr))
		}
	}
	jobsNeedingLogs := make([]WorkflowJob, 0)
	for _, workflowJob := range workflowJobs {
		workflowJobInput := workflowJobBuildInput{job: workflowJob}
		if !workflowSucceeded && workflowJob.NeedsStepLogs() {
			jobsNeedingLogs = append(jobsNeedingLogs, workflowJob)
		}
		jobBuildInput.jobs = append(jobBuildInput.jobs, workflowJobInput)
	}
	if len(jobsNeedingLogs) == 0 {
		return jobBuildInput
	}
	attemptLogs, err := builder.workflowProvider.StepLogsForAttempt(ctx, builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, jobsNeedingLogs)
	if err != nil {
		jobBuildInput.collectionErrors = append(jobBuildInput.collectionErrors, fmt.Errorf("could not read logs for workflow run %s attempt %d: %w", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, err))
	}
	for jobIndex := range jobBuildInput.jobs {
		workflowJobInput := &jobBuildInput.jobs[jobIndex]
		if workflowSucceeded || !workflowJobInput.job.NeedsStepLogs() {
			continue
		}
		jobLogs, found := attemptLogs.ForJobID(workflowJobInput.job.ID)
		if !found {
			continue
		}
		workflowJobInput.jobOutput = jobLogs.JobOutput
		for _, workflowStep := range workflowJobInput.job.Steps {
			for _, stepLog := range jobLogs.Steps {
				if workflowStep.Number == stepLog.Number {
					workflowJobInput.stepLogs = append(workflowJobInput.stepLogs, workflowStepLog{
						stepNumber: stepLog.Number,
						text:       stepLog.Text,
					})
					break
				}
			}
		}
	}
	return jobBuildInput
}

func (builder *workflowRecordBuilder) buildTestResultJobs(ctx context.Context, jobBuildInput testResultJobBuildInput) (jobs []testresults.Job, collectionErrors []error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanBuildTestResultJobs, attribute.Int(traceAttributeJobCount, len(jobBuildInput.jobs)))
	defer func() {
		telemetry.FinishSpan(ctx, span, spanBuildTestResultJobs, startedAt, len(collectionErrors) > 0)
	}()
	jobs = make([]testresults.Job, 0, len(jobBuildInput.jobs))
	collectionErrors = make([]error, 0)
	for _, workflowJobInput := range jobBuildInput.jobs {
		job, jobErrors := builder.buildTestResultJob(ctx, workflowJobInput)
		jobs = append(jobs, job)
		collectionErrors = append(collectionErrors, jobErrors...)
	}
	return jobs, collectionErrors
}

func (builder *workflowRecordBuilder) buildTestResultJob(ctx context.Context, workflowJobInput workflowJobBuildInput) (testresults.Job, []error) {
	workflowJob := workflowJobInput.job
	job := testresults.Job{
		Name:   workflowJob.Name,
		JobURL: workflowJob.JobURL,
		Steps:  make([]testresults.Step, 0, len(workflowJob.Steps)),
	}
	job.Outcome, job.ProviderConclusion = resultOutcome(workflowJob.Conclusion, "")
	collectionErrors := make([]error, 0)
	if len(workflowJobInput.stepLogs) == 0 && workflowJobInput.jobOutput != "" {
		var err error
		job.Tests, err = builder.parseTestLog(ctx, workflowJobInput.jobOutput)
		if err != nil {
			collectionErrors = append(collectionErrors, fmt.Errorf("could not parse logs for workflow run %s attempt %d job %q: %w", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, workflowJob.Name, err))
		}
	}
	for _, workflowStep := range workflowJob.Steps {
		stepLog := workflowJobInput.logForStepNumber(workflowStep.Number)
		step, err := builder.buildTestResultStep(ctx, workflowJob, workflowStep, stepLog, workflowJobInput.jobOutput)
		job.Steps = append(job.Steps, step)
		if err != nil {
			collectionErrors = append(collectionErrors, err)
		}
	}
	return job, collectionErrors
}

func (builder *workflowRecordBuilder) buildTestResultStep(ctx context.Context, workflowJob WorkflowJob, workflowStep WorkflowStep, stepLog string, jobOutput string) (testresults.Step, error) {
	step := testresults.Step{Name: workflowStep.Name}
	step.Outcome, step.ProviderConclusion = resultOutcome(workflowStep.Conclusion, "")
	if strings.EqualFold(workflowStep.Conclusion, workflowConclusionFailure) {
		step.ReasonExcerpt = extractFailureLogSnippet(stepLog)
		if stepLog == "" && jobOutput == "" {
			return step, fmt.Errorf("workflow run %s attempt %d job %q step %q has no downloadable log", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, workflowJob.Name, workflowStep.Name)
		}
	}
	if stepLog != "" {
		var err error
		step.Tests, err = builder.parseTestLog(ctx, stepLog)
		if err != nil {
			return step, fmt.Errorf("could not parse logs for workflow run %s attempt %d job %q step %q: %w", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, workflowJob.Name, workflowStep.Name, err)
		}
	}
	return step, nil
}

func (builder *workflowRecordBuilder) parseTestLog(ctx context.Context, stepLog string) (tests []testresults.Test, resultErr error) {
	if builder.testLogParserMutex != nil {
		builder.testLogParserMutex.Lock()
		defer builder.testLogParserMutex.Unlock()
	}
	parser, err := selectTestLogParser(stepLog, builder.testLogParsers)
	if err != nil || parser == nil {
		return nil, err
	}
	startedAt := time.Now()
	parseContext, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanParseTestLog,
		attribute.Int(traceAttributeLogByteCount, len(stepLog)),
		attribute.Int(traceAttributeTestLogParserSpecificity, parser.Specificity()),
	)
	tests, resultErr = parser.Parse(stepLog)
	telemetry.FinishSpan(parseContext, span, spanParseTestLog, startedAt, resultErr != nil)
	return tests, resultErr
}

func (builder *workflowRecordBuilder) addErrors(collectionErrors []error) {
	builder.collectionErrors = append(builder.collectionErrors, collectionErrors...)
}

func (builder *workflowRecordBuilder) validateWorkflowRecord() (*testresults.Workflow, error) {
	builder.workflowRun.CollectionComplete = len(builder.collectionErrors) == 0
	if err := builder.workflowRun.Validate(); err != nil {
		builder.collectionErrors = append(builder.collectionErrors, fmt.Errorf("workflow run %s attempt %d is invalid: %w", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, err))
		return nil, errors.Join(builder.collectionErrors...)
	}
	return builder.workflowRun, errors.Join(builder.collectionErrors...)
}

func resultOutcome(conclusion string, status string) (testresults.Outcome, string) {
	switch strings.ToLower(conclusion) {
	case workflowConclusionSuccess:
		return testresults.OutcomePassed, ""
	case workflowConclusionFailure:
		return testresults.OutcomeFailed, ""
	case workflowConclusionSkipped:
		return testresults.OutcomeSkipped, ""
	default:
		if conclusion != "" {
			return testresults.OutcomeInconclusive, string(conclusion)
		}
		if status != "" {
			return testresults.OutcomeInconclusive, status
		}
		return testresults.OutcomeInconclusive, workflowConclusionUnknown
	}
}
