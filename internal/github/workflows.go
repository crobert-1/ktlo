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

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"ktlo/internal/common"
	"ktlo/internal/telemetry"
)

const (
	defaultRESTPageSize         = 100
	maxFilteredWorkflowRunCount = 1000
	restAPIEndpoint             = "https://api.github.com"
)

const maxConcurrentWorkflowRunMetadataRequests = 4

const (
	errorWorkflowAPIRequest        = "GitHub workflow API request failed"
	errorGetLatestWorkflowRun      = "could not get latest workflow run attempt for run %d: %w"
	errorInvalidLatestWorkflowRun  = "GitHub returned invalid workflow run ID %d or attempt %d for run %d"
	errorWorkflowRunAttemptID      = "workflow run ID and attempt must be positive"
	errorWorkflowRunSearchLimit    = "GitHub reports %d workflow runs in the lookback window, exceeding the API search limit of %d; analysis is incomplete"
	githubInstrumentationName      = "ktlo/internal/github"
	spanListWorkflowJobs           = "github.list_workflow_jobs"
	spanListWorkflowRunAttempts    = "github.list_workflow_run_attempts"
	spanLatestWorkflowRunAttempt   = "github.latest_workflow_run_attempt"
	spanWorkflowAttemptLogArchive  = "github.workflow_attempt_log_archive"
	traceAttributeJobCount         = "github.workflow_job.count"
	traceAttributeRunAttempt       = "github.workflow_run.attempt"
	traceAttributeRunID            = "github.workflow_run.id"
	traceAttributeRunStartedAfter  = "github.workflow_run.started_after"
	traceAttributeWorkflowRunCount = "github.workflow_run_attempt.count"
	workflowRunCreatedQueryParam   = "created"
)

type WorkflowConclusion string

func (conclusion *WorkflowConclusion) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*conclusion = WorkflowConclusion(strings.ToLower(value))
	return nil
}

const (
	WorkflowConclusionActionRequired WorkflowConclusion = "action_required"
	WorkflowConclusionCancelled      WorkflowConclusion = "cancelled"
	WorkflowConclusionFailure        WorkflowConclusion = "failure"
	WorkflowConclusionNeutral        WorkflowConclusion = "neutral"
	WorkflowConclusionSkipped        WorkflowConclusion = "skipped"
	WorkflowConclusionStale          WorkflowConclusion = "stale"
	WorkflowConclusionSuccess        WorkflowConclusion = "success"
	WorkflowConclusionTimedOut       WorkflowConclusion = "timed_out"
	WorkflowConclusionUnknown        WorkflowConclusion = "unknown"
)

func (conclusion WorkflowConclusion) IsValid() bool {
	switch conclusion {
	case WorkflowConclusionActionRequired, WorkflowConclusionCancelled, WorkflowConclusionFailure, WorkflowConclusionNeutral, WorkflowConclusionSkipped, WorkflowConclusionStale, WorkflowConclusionSuccess, WorkflowConclusionTimedOut, WorkflowConclusionUnknown:
		return true
	default:
		return false
	}
}

func (c *Client) ListAllWorkflowRunAttempts(ctx context.Context, startedAfter time.Time, knownAttempts map[common.WorkflowRunAttemptKey]struct{}) (workflowRunAttempts []common.WorkflowRunAttempt, resultErr error) {
	startedAt := time.Now()
	ctx, span := otel.Tracer(githubInstrumentationName).Start(ctx, spanListWorkflowRunAttempts)
	defer func() {
		telemetry.RecordDuration(ctx, spanListWorkflowRunAttempts, startedAt)
		if resultErr != nil {
			span.SetStatus(codes.Error, errorWorkflowAPIRequest)
		}
		span.SetAttributes(
			attribute.String(traceAttributeRunStartedAfter, startedAfter.Format(time.RFC3339Nano)),
			attribute.Int(traceAttributeWorkflowRunCount, len(workflowRunAttempts)),
		)
		span.End()
	}()
	workflowRunAttempts = make([]common.WorkflowRunAttempt, 0)
	queryErrors := make([]error, 0)
	for pageNumber := 1; ; pageNumber++ {
		if err := ctx.Err(); err != nil {
			queryErrors = append(queryErrors, err)
			return workflowRunAttempts, errors.Join(queryErrors...)
		}
		var page workflowRunsPage
		if err := c.getRESTJSON(ctx, c.workflowRunsEndpoint(pageNumber, startedAfter), &page); err != nil {
			queryErrors = append(queryErrors, err)
			break
		}
		if pageNumber == 1 && !startedAfter.IsZero() && page.TotalCount > maxFilteredWorkflowRunCount {
			queryErrors = append(queryErrors, fmt.Errorf(errorWorkflowRunSearchLimit, page.TotalCount, maxFilteredWorkflowRunCount))
		}
		candidateRuns := make([]workflowRun, 0, len(page.WorkflowRuns))
		for _, run := range page.WorkflowRuns {
			if run.ID < 1 || run.RunAttempt < 1 {
				queryErrors = append(queryErrors, fmt.Errorf("GitHub response included invalid workflow run ID %d or attempt %d", run.ID, run.RunAttempt))
				continue
			}
			if !startedAfter.IsZero() && !run.RunStartedAt.IsZero() && run.RunStartedAt.Before(startedAfter) {
				continue
			}
			candidateRuns = append(candidateRuns, run)
		}
		pageAttempts, attemptErrors := c.workflowRunAttemptsForRuns(ctx, candidateRuns, startedAfter, knownAttempts)
		workflowRunAttempts = append(workflowRunAttempts, pageAttempts...)
		queryErrors = append(queryErrors, attemptErrors...)
		resultLimit := page.TotalCount
		if !startedAfter.IsZero() && resultLimit > maxFilteredWorkflowRunCount {
			resultLimit = maxFilteredWorkflowRunCount
		}
		if len(page.WorkflowRuns) < defaultRESTPageSize || pageNumber*defaultRESTPageSize >= resultLimit {
			break
		}
	}
	return workflowRunAttempts, errors.Join(queryErrors...)
}

func (c *Client) workflowRunAttemptsForRuns(ctx context.Context, runs []workflowRun, startedAfter time.Time, knownAttempts map[common.WorkflowRunAttemptKey]struct{}) ([]common.WorkflowRunAttempt, []error) {
	tasks := make([]workflowRunAttemptRequest, 0)
	for _, run := range runs {
		for attemptNumber := 1; attemptNumber <= run.RunAttempt; attemptNumber++ {
			key := common.WorkflowRunAttemptKey{RunID: strconv.FormatInt(run.ID, 10), Attempt: attemptNumber}
			if _, known := knownAttempts[key]; known {
				continue
			}
			tasks = append(tasks, workflowRunAttemptRequest{attempt: attemptNumber, run: run})
		}
	}
	results := make([]workflowRunAttemptResult, len(tasks))
	workerCount := min(maxConcurrentWorkflowRunMetadataRequests, len(tasks))
	taskIndexes := make(chan int)
	var workers sync.WaitGroup
	for range workerCount {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for taskIndex := range taskIndexes {
				results[taskIndex] = c.workflowRunAttemptForRequest(ctx, tasks[taskIndex], startedAfter)
			}
		}()
	}
	for taskIndex := range tasks {
		taskIndexes <- taskIndex
	}
	close(taskIndexes)
	workers.Wait()
	attempts := make([]common.WorkflowRunAttempt, 0, len(results))
	queryErrors := make([]error, 0)
	for _, result := range results {
		if result.err != nil {
			queryErrors = append(queryErrors, result.err)
			continue
		}
		if result.inLookback {
			attempts = append(attempts, result.attempt)
		}
	}
	return attempts, queryErrors
}

func (c *Client) workflowRunAttemptForRequest(ctx context.Context, request workflowRunAttemptRequest, startedAfter time.Time) workflowRunAttemptResult {
	if err := ctx.Err(); err != nil {
		return workflowRunAttemptResult{err: err}
	}
	attempt := request.run
	if request.attempt != request.run.RunAttempt || request.run.RunStartedAt.IsZero() {
		attempt = workflowRun{}
		endpoint := c.workflowRunAttemptEndpoint(request.run.ID, request.attempt)
		if err := c.getRESTJSON(ctx, endpoint, &attempt); err != nil {
			return workflowRunAttemptResult{err: fmt.Errorf("could not get workflow run %d attempt %d: %w", request.run.ID, request.attempt, err)}
		}
	}
	if !startedAfter.IsZero() && attempt.RunStartedAt.Before(startedAfter) {
		return workflowRunAttemptResult{}
	}
	return workflowRunAttemptResult{attempt: attempt.toWorkflowRunAttempt(), inLookback: true}
}

func (c *Client) ListWorkflowJobs(ctx context.Context, runID int64, attempt int) (workflowJobs []common.WorkflowJob, resultErr error) {
	startedAt := time.Now()
	ctx, span := otel.Tracer(githubInstrumentationName).Start(ctx, spanListWorkflowJobs)
	defer func() {
		telemetry.RecordDuration(ctx, spanListWorkflowJobs, startedAt)
		if resultErr != nil {
			span.SetStatus(codes.Error, errorWorkflowAPIRequest)
		}
		span.SetAttributes(
			attribute.Int(traceAttributeRunAttempt, attempt),
			attribute.Int64(traceAttributeRunID, runID),
			attribute.Int(traceAttributeJobCount, len(workflowJobs)),
		)
		span.End()
	}()
	if runID < 1 || attempt < 1 {
		return nil, errors.New(errorWorkflowRunAttemptID)
	}
	workflowJobs = make([]common.WorkflowJob, 0)
	queryErrors := make([]error, 0)
	for pageNumber := 1; ; pageNumber++ {
		if err := ctx.Err(); err != nil {
			queryErrors = append(queryErrors, err)
			return workflowJobs, errors.Join(queryErrors...)
		}
		var page workflowJobsPage
		endpoint := c.workflowJobsEndpoint(runID, attempt, pageNumber)
		if err := c.getRESTJSON(ctx, endpoint, &page); err != nil {
			queryErrors = append(queryErrors, err)
			break
		}
		for _, job := range page.Jobs {
			workflowJob := common.WorkflowJob{
				ID:         strconv.FormatInt(job.ID, 10),
				Name:       job.Name,
				Conclusion: string(job.Conclusion),
				JobURL:     job.HTMLURL,
				Steps:      make([]common.WorkflowStep, 0, len(job.Steps)),
			}
			for _, step := range job.Steps {
				workflowJob.Steps = append(workflowJob.Steps, common.WorkflowStep{
					Number:      step.Number,
					Name:        step.Name,
					Conclusion:  string(step.Conclusion),
					StartedAt:   step.StartedAt,
					CompletedAt: step.CompletedAt,
				})
			}
			workflowJobs = append(workflowJobs, workflowJob)
		}
		if len(page.Jobs) < defaultRESTPageSize || pageNumber*defaultRESTPageSize >= page.TotalCount {
			break
		}
	}
	return workflowJobs, errors.Join(queryErrors...)
}

func (c *Client) LatestWorkflowRunAttempt(ctx context.Context, runID int64) (workflowRunAttempt common.WorkflowRunAttempt, resultErr error) {
	startedAt := time.Now()
	ctx, span := otel.Tracer(githubInstrumentationName).Start(ctx, spanLatestWorkflowRunAttempt)
	defer func() {
		telemetry.RecordDuration(ctx, spanLatestWorkflowRunAttempt, startedAt)
		if resultErr != nil {
			span.SetStatus(codes.Error, errorWorkflowAPIRequest)
		}
		span.SetAttributes(attribute.Int64(traceAttributeRunID, runID))
		span.End()
	}()
	if runID < 1 {
		return common.WorkflowRunAttempt{}, errors.New(errorWorkflowRunAttemptID)
	}
	var workflowRun workflowRun
	if err := c.getRESTJSON(ctx, c.workflowRunEndpoint(runID), &workflowRun); err != nil {
		return common.WorkflowRunAttempt{}, fmt.Errorf(errorGetLatestWorkflowRun, runID, err)
	}
	if workflowRun.ID != runID || workflowRun.RunAttempt < 1 {
		return common.WorkflowRunAttempt{}, fmt.Errorf(errorInvalidLatestWorkflowRun, workflowRun.ID, workflowRun.RunAttempt, runID)
	}
	workflowRunAttempt = workflowRun.toWorkflowRunAttempt()
	span.SetAttributes(attribute.Int(traceAttributeRunAttempt, workflowRunAttempt.Attempt))
	return workflowRunAttempt, nil
}

// WorkflowAttemptLogArchive returns the ZIP log archive for one workflow attempt.
func (c *Client) WorkflowAttemptLogArchive(ctx context.Context, runID int64, attempt int) (logArchive io.ReadCloser, resultErr error) {
	startedAt := time.Now()
	ctx, span := otel.Tracer(githubInstrumentationName).Start(ctx, spanWorkflowAttemptLogArchive)
	defer func() {
		telemetry.RecordDuration(ctx, spanWorkflowAttemptLogArchive, startedAt)
		if resultErr != nil {
			span.SetStatus(codes.Error, errorWorkflowAPIRequest)
		}
		span.SetAttributes(
			attribute.Int(traceAttributeRunAttempt, attempt),
			attribute.Int64(traceAttributeRunID, runID),
		)
		span.End()
	}()
	if runID < 1 || attempt < 1 {
		return nil, errors.New(errorWorkflowRunAttemptID)
	}
	var err error
	var response *http.Response
	response, err = c.getRESTResponse(ctx, c.workflowAttemptLogArchiveEndpoint(runID, attempt))
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (c *Client) workflowRunsEndpoint(pageNumber int, createdAfter time.Time) string {
	query := make(url.Values)
	query.Set("page", fmt.Sprint(pageNumber))
	query.Set("per_page", fmt.Sprint(defaultRESTPageSize))
	if !createdAfter.IsZero() {
		// This limits candidate workflow runs; attempts are still filtered by start time after fetching.
		createdAfter = createdAfter.UTC().Truncate(time.Second)
		query.Set(workflowRunCreatedQueryParam, ">="+createdAfter.Format(time.RFC3339))
	}
	return fmt.Sprintf("%s/repos/%s/%s/actions/runs?%s", restAPIEndpoint, url.PathEscape(c.owner), url.PathEscape(c.repositoryName), query.Encode())
}

func (c *Client) workflowRunAttemptEndpoint(runID int64, attempt int) string {
	return fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/attempts/%d", restAPIEndpoint, url.PathEscape(c.owner), url.PathEscape(c.repositoryName), runID, attempt)
}

func (c *Client) workflowRunEndpoint(runID int64) string {
	return fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d", restAPIEndpoint, url.PathEscape(c.owner), url.PathEscape(c.repositoryName), runID)
}

func (c *Client) workflowJobsEndpoint(runID int64, attempt int, pageNumber int) string {
	query := make(url.Values)
	query.Set("page", fmt.Sprint(pageNumber))
	query.Set("per_page", fmt.Sprint(defaultRESTPageSize))
	return fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/attempts/%d/jobs?%s", restAPIEndpoint, url.PathEscape(c.owner), url.PathEscape(c.repositoryName), runID, attempt, query.Encode())
}

func (c *Client) workflowAttemptLogArchiveEndpoint(runID int64, attempt int) string {
	return fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/attempts/%d/logs", restAPIEndpoint, url.PathEscape(c.owner), url.PathEscape(c.repositoryName), runID, attempt)
}

type workflowRunsPage struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []workflowRun `json:"workflow_runs"`
}

// workflowRun mirrors GitHub's response before mapping it to WorkflowRunAttempt.
type workflowRun struct {
	ID           int64                    `json:"id"`
	Name         string                   `json:"name"`
	HeadBranch   string                   `json:"head_branch"`
	HeadSHA      string                   `json:"head_sha"`
	Status       string                   `json:"status"`
	Conclusion   WorkflowConclusion       `json:"conclusion"`
	HTMLURL      string                   `json:"html_url"`
	RunAttempt   int                      `json:"run_attempt"`
	RunStartedAt time.Time                `json:"run_started_at"`
	PullRequests []workflowRunPullRequest `json:"pull_requests"`
}

type workflowRunAttemptRequest struct {
	attempt int
	run     workflowRun
}

type workflowRunAttemptResult struct {
	attempt    common.WorkflowRunAttempt
	err        error
	inLookback bool
}

func (run *workflowRun) toWorkflowRunAttempt() common.WorkflowRunAttempt {
	workflowRunAttempt := common.WorkflowRunAttempt{
		RunID:        strconv.FormatInt(run.ID, 10),
		Attempt:      run.RunAttempt,
		Name:         run.Name,
		Branch:       run.HeadBranch,
		HeadSHA:      run.HeadSHA,
		RunStartedAt: run.RunStartedAt,
		Status:       run.Status,
		Conclusion:   string(run.Conclusion),
		WorkflowURL:  run.HTMLURL,
		PullRequests: make([]common.PullRequestReference, 0, len(run.PullRequests)),
	}
	for _, pullRequest := range run.PullRequests {
		workflowRunAttempt.PullRequests = append(workflowRunAttempt.PullRequests, common.PullRequestReference{
			Number:  pullRequest.Number,
			URL:     pullRequest.URL,
			HeadSHA: pullRequest.Head.SHA,
		})
	}
	return workflowRunAttempt
}

type workflowRunPullRequest struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
	Head   struct {
		SHA string `json:"sha"`
	} `json:"head"`
}

type workflowJobsPage struct {
	TotalCount int           `json:"total_count"`
	Jobs       []workflowJob `json:"jobs"`
}

type workflowJob struct {
	ID         int64              `json:"id"`
	Name       string             `json:"name"`
	Conclusion WorkflowConclusion `json:"conclusion"`
	HTMLURL    string             `json:"html_url"`
	Steps      []workflowJobStep  `json:"steps"`
}

type workflowJobStep struct {
	Number      int                `json:"number"`
	Name        string             `json:"name"`
	Conclusion  WorkflowConclusion `json:"conclusion"`
	CompletedAt time.Time          `json:"completed_at"`
	StartedAt   time.Time          `json:"started_at"`
}
