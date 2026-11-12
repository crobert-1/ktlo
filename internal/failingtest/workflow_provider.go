package failingtest

import (
	"context"
	"time"

	"ktlo/internal/testresults"
)

type WorkflowProvider interface {
	Provider() testresults.Provider
	Repository() testresults.Repository
	ListWorkflowRunAttempts(ctx context.Context, startedAfter time.Time) ([]WorkflowRunAttempt, error)
	ListWorkflowJobs(ctx context.Context, runID string, attempt int) ([]WorkflowJob, error)
	StepLogs(ctx context.Context, jobID string) (map[int]string, error)
}

type WorkflowRunAttempt struct {
	RunID        string
	Attempt      int
	Name         string
	Branch       string
	HeadSHA      string
	RunStartedAt time.Time
	Status       string
	Conclusion   string
	WorkflowURL  string
	PullRequests []AssociatedPullRequest
}

type AssociatedPullRequest struct {
	Number  int
	URL     string
	HeadSHA string
}

type WorkflowJob struct {
	ID         string
	Name       string
	Conclusion string
	JobURL     string
	Steps      []WorkflowStep
}

type WorkflowStep struct {
	Number     int
	Name       string
	Conclusion string
}
