package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultRESTPageSize = 100
	restAPIEndpoint     = "https://api.github.com"
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

type WorkflowRunAttempt struct {
	RunID              int64
	Attempt            int
	Name               string
	Branch             string
	HeadSHA            string
	RunStartedAt       time.Time
	Status             string
	ProviderConclusion WorkflowConclusion
	WorkflowURL        string
	PullRequests       []AssociatedPullRequest
}

type AssociatedPullRequest struct {
	PRNumber int
	URL      string
	HeadSHA  string
}

type WorkflowJob struct {
	ID                 int64
	Name               string
	ProviderConclusion WorkflowConclusion
	JobURL             string
	Steps              []WorkflowStep
}

type WorkflowStep struct {
	Number             int
	Name               string
	ProviderConclusion WorkflowConclusion
}

func (c *Client) ListAllWorkflowRunAttempts(ctx context.Context, startedAfter time.Time) ([]WorkflowRunAttempt, error) {
	workflowRunAttempts := make([]WorkflowRunAttempt, 0)
	queryErrors := make([]error, 0)
	for pageNumber := 1; ; pageNumber++ {
		if err := ctx.Err(); err != nil {
			queryErrors = append(queryErrors, err)
			return workflowRunAttempts, errors.Join(queryErrors...)
		}
		var page workflowRunsPage
		if err := c.getRESTJSON(ctx, c.workflowRunsEndpoint(pageNumber), &page); err != nil {
			queryErrors = append(queryErrors, err)
			break
		}
		for _, run := range page.WorkflowRuns {
			if run.ID < 1 || run.RunAttempt < 1 {
				queryErrors = append(queryErrors, fmt.Errorf("GitHub response included invalid workflow run ID %d or attempt %d", run.ID, run.RunAttempt))
				continue
			}
			for attemptNumber := 1; attemptNumber <= run.RunAttempt; attemptNumber++ {
				if err := ctx.Err(); err != nil {
					queryErrors = append(queryErrors, err)
					return workflowRunAttempts, errors.Join(queryErrors...)
				}
				var attempt workflowRun
				endpoint := c.workflowRunAttemptEndpoint(run.ID, attemptNumber)
				if err := c.getRESTJSON(ctx, endpoint, &attempt); err != nil {
					queryErrors = append(queryErrors, fmt.Errorf("could not get workflow run %d attempt %d: %w", run.ID, attemptNumber, err))
					continue
				}
				if !startedAfter.IsZero() && attempt.RunStartedAt.Before(startedAfter) {
					continue
				}
				workflowRunAttempts = append(workflowRunAttempts, attempt.toWorkflowRunAttempt())
			}
		}
		if len(page.WorkflowRuns) < defaultRESTPageSize || pageNumber*defaultRESTPageSize >= page.TotalCount {
			break
		}
	}
	return workflowRunAttempts, errors.Join(queryErrors...)
}

func (c *Client) ListWorkflowJobs(ctx context.Context, runID int64, attempt int) ([]WorkflowJob, error) {
	if runID < 1 || attempt < 1 {
		return nil, errors.New("workflow run ID and attempt must be positive")
	}
	workflowJobs := make([]WorkflowJob, 0)
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
			workflowJob := WorkflowJob{
				ID:                 job.ID,
				Name:               job.Name,
				ProviderConclusion: job.Conclusion,
				JobURL:             job.HTMLURL,
				Steps:              make([]WorkflowStep, 0, len(job.Steps)),
			}
			for _, step := range job.Steps {
				workflowJob.Steps = append(workflowJob.Steps, WorkflowStep{
					Number:             step.Number,
					Name:               step.Name,
					ProviderConclusion: step.Conclusion,
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

func (c *Client) JobLog(ctx context.Context, jobID int64) (io.ReadCloser, error) {
	if jobID < 1 {
		return nil, errors.New("workflow job ID must be positive")
	}
	var err error
	var response *http.Response
	response, err = c.getRESTResponse(ctx, c.jobLogEndpoint(jobID))
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (c *Client) workflowRunsEndpoint(pageNumber int) string {
	query := make(url.Values)
	query.Set("page", fmt.Sprint(pageNumber))
	query.Set("per_page", fmt.Sprint(defaultRESTPageSize))
	return fmt.Sprintf("%s/repos/%s/%s/actions/runs?%s", restAPIEndpoint, url.PathEscape(c.owner), url.PathEscape(c.repositoryName), query.Encode())
}

func (c *Client) workflowRunAttemptEndpoint(runID int64, attempt int) string {
	return fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/attempts/%d", restAPIEndpoint, url.PathEscape(c.owner), url.PathEscape(c.repositoryName), runID, attempt)
}

func (c *Client) workflowJobsEndpoint(runID int64, attempt int, pageNumber int) string {
	query := make(url.Values)
	query.Set("page", fmt.Sprint(pageNumber))
	query.Set("per_page", fmt.Sprint(defaultRESTPageSize))
	return fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/attempts/%d/jobs?%s", restAPIEndpoint, url.PathEscape(c.owner), url.PathEscape(c.repositoryName), runID, attempt, query.Encode())
}

func (c *Client) jobLogEndpoint(jobID int64) string {
	return fmt.Sprintf("%s/repos/%s/%s/actions/jobs/%d/logs", restAPIEndpoint, url.PathEscape(c.owner), url.PathEscape(c.repositoryName), jobID)
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

func (run *workflowRun) toWorkflowRunAttempt() WorkflowRunAttempt {
	workflowRunAttempt := WorkflowRunAttempt{
		RunID:              run.ID,
		Attempt:            run.RunAttempt,
		Name:               run.Name,
		Branch:             run.HeadBranch,
		HeadSHA:            run.HeadSHA,
		RunStartedAt:       run.RunStartedAt,
		Status:             run.Status,
		ProviderConclusion: run.Conclusion,
		WorkflowURL:        run.HTMLURL,
		PullRequests:       make([]AssociatedPullRequest, 0, len(run.PullRequests)),
	}
	for _, pullRequest := range run.PullRequests {
		workflowRunAttempt.PullRequests = append(workflowRunAttempt.PullRequests, AssociatedPullRequest{
			PRNumber: pullRequest.Number,
			URL:      pullRequest.URL,
			HeadSHA:  pullRequest.Head.SHA,
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
	Number     int                `json:"number"`
	Name       string             `json:"name"`
	Conclusion WorkflowConclusion `json:"conclusion"`
}
