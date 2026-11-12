package failingtest

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ktlo/internal/testresults"
)

const workflowConclusionFailure = "failure"
const workflowConclusionSkipped = "skipped"
const workflowConclusionSuccess = "success"
const workflowConclusionUnknown = "unknown"

type workflowStepLog struct {
	stepNumber int
	text       string
}

type workflowJobBuildInput struct {
	job      WorkflowJob
	stepLogs []workflowStepLog
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
	workflowRun        *testresults.Workflow
	workflowRunAttempt WorkflowRunAttempt
}

func newWorkflowRecordBuilder(workflowProvider WorkflowProvider, testLogParsers []TestLogParser) *workflowRecordBuilder {
	return &workflowRecordBuilder{
		workflowProvider: workflowProvider,
		testLogParsers:   testLogParsers,
	}
}

func (builder *workflowRecordBuilder) build(ctx context.Context, workflowRunAttempt WorkflowRunAttempt) (*testresults.Workflow, error) {
	builder.collectionErrors = nil
	builder.workflowRun = generateReportWorkflowFromAttempt(workflowRunAttempt)
	builder.workflowRunAttempt = workflowRunAttempt
	builder.addPullRequest()
	jobBuildInput := builder.collectWorkflowJobLogs(ctx)
	testResultJobs, jobErrors := builder.buildTestResultJobs(jobBuildInput)
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

func (builder *workflowRecordBuilder) collectWorkflowJobLogs(ctx context.Context) testResultJobBuildInput {
	jobBuildInput := testResultJobBuildInput{jobs: make([]workflowJobBuildInput, 0)}
	workflowJobs, listErr := builder.workflowProvider.ListWorkflowJobs(ctx, builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt)
	if listErr != nil {
		jobBuildInput.collectionErrors = append(jobBuildInput.collectionErrors, fmt.Errorf("could not list jobs for workflow run %s attempt %d: %w", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, listErr))
	}
	for _, workflowJob := range workflowJobs {
		workflowJobInput := workflowJobBuildInput{job: workflowJob}
		if strings.EqualFold(workflowJob.Conclusion, workflowConclusionFailure) || strings.EqualFold(workflowJob.Conclusion, workflowConclusionSuccess) {
			stepLogs, err := builder.workflowProvider.StepLogs(ctx, workflowJob.ID)
			if err != nil {
				jobBuildInput.collectionErrors = append(jobBuildInput.collectionErrors, fmt.Errorf("could not read logs for workflow job %q: %w", workflowJob.ID, err))
			} else {
				for _, workflowStep := range workflowJob.Steps {
					if stepLog, found := stepLogs[workflowStep.Number]; found {
						workflowJobInput.stepLogs = append(workflowJobInput.stepLogs, workflowStepLog{
							stepNumber: workflowStep.Number,
							text:       stepLog,
						})
					}
				}
			}
		}
		jobBuildInput.jobs = append(jobBuildInput.jobs, workflowJobInput)
	}
	return jobBuildInput
}

func (builder *workflowRecordBuilder) buildTestResultJobs(jobBuildInput testResultJobBuildInput) ([]testresults.Job, []error) {
	jobs := make([]testresults.Job, 0, len(jobBuildInput.jobs))
	collectionErrors := make([]error, 0)
	for _, workflowJobInput := range jobBuildInput.jobs {
		job, jobErrors := builder.buildTestResultJob(workflowJobInput)
		jobs = append(jobs, job)
		collectionErrors = append(collectionErrors, jobErrors...)
	}
	return jobs, collectionErrors
}

func (builder *workflowRecordBuilder) buildTestResultJob(workflowJobInput workflowJobBuildInput) (testresults.Job, []error) {
	workflowJob := workflowJobInput.job
	job := testresults.Job{
		Name:   workflowJob.Name,
		JobURL: workflowJob.JobURL,
		Steps:  make([]testresults.Step, 0, len(workflowJob.Steps)),
	}
	job.Outcome, job.ProviderConclusion = resultOutcome(workflowJob.Conclusion, "")
	collectionErrors := make([]error, 0)
	for _, workflowStep := range workflowJob.Steps {
		step, err := builder.buildTestResultStep(workflowJob, workflowStep, workflowJobInput.logForStepNumber(workflowStep.Number))
		job.Steps = append(job.Steps, step)
		if err != nil {
			collectionErrors = append(collectionErrors, err)
		}
	}
	return job, collectionErrors
}

func (builder *workflowRecordBuilder) buildTestResultStep(workflowJob WorkflowJob, workflowStep WorkflowStep, stepLog string) (testresults.Step, error) {
	step := testresults.Step{Name: workflowStep.Name}
	step.Outcome, step.ProviderConclusion = resultOutcome(workflowStep.Conclusion, "")
	if strings.EqualFold(workflowStep.Conclusion, workflowConclusionFailure) {
		step.ReasonExcerpt = extractFailureLogSnippet(stepLog)
		if stepLog == "" {
			return step, fmt.Errorf("workflow run %s attempt %d job %q step %q has no downloadable log", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, workflowJob.Name, workflowStep.Name)
		}
	}
	if stepLog != "" {
		var err error
		step.Tests, err = builder.parseTestLog(stepLog)
		if err != nil {
			return step, fmt.Errorf("could not parse logs for workflow run %s attempt %d job %q step %q: %w", builder.workflowRunAttempt.RunID, builder.workflowRunAttempt.Attempt, workflowJob.Name, workflowStep.Name, err)
		}
	}
	return step, nil
}

func (builder *workflowRecordBuilder) parseTestLog(stepLog string) ([]testresults.Test, error) {
	for _, parser := range builder.testLogParsers {
		if parser.CanParse(stepLog) {
			return parser.Parse(stepLog)
		}
	}
	return nil, nil
}

func (builder *workflowRecordBuilder) addErrors(collectionErrors []error) {
	builder.collectionErrors = append(builder.collectionErrors, collectionErrors...)
}

func (builder *workflowRecordBuilder) validateWorkflowRecord() (*testresults.Workflow, error) {
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
