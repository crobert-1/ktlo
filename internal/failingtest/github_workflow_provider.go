package failingtest

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"time"

	"ktlo/internal/github"
	"ktlo/internal/testresults"
)

const maxGitHubJobLogArchiveBytes = 50 << 20
const maxGitHubJobLogFileBytes = 10 << 20
const githubLogStepNameSeparator = "_"

type GitHubWorkflowProvider struct {
	client *github.Client
}

var _ WorkflowProvider = (*GitHubWorkflowProvider)(nil)

func NewGitHubWorkflowProvider(client *github.Client) (*GitHubWorkflowProvider, error) {
	if client == nil {
		return nil, errors.New("GitHub client is required for the GitHub workflow provider")
	}
	return &GitHubWorkflowProvider{client: client}, nil
}

func (provider *GitHubWorkflowProvider) Provider() testresults.Provider {
	return testresults.ProviderGitHub
}

func (provider *GitHubWorkflowProvider) Repository() testresults.Repository {
	repository := provider.client.TargetRepository()
	return testresults.Repository{
		Owner: repository.Owner,
		Name:  repository.Name,
	}
}

func (provider *GitHubWorkflowProvider) ListWorkflowRunAttempts(ctx context.Context, startedAfter time.Time) ([]WorkflowRunAttempt, error) {
	githubAttempts, err := provider.client.ListAllWorkflowRunAttempts(ctx, startedAfter)
	workflowRunAttempts := make([]WorkflowRunAttempt, 0, len(githubAttempts))
	for _, githubAttempt := range githubAttempts {
		workflowRunAttempts = append(workflowRunAttempts, provider.workflowRunAttempt(githubAttempt))
	}
	return workflowRunAttempts, err
}

func (provider *GitHubWorkflowProvider) workflowRunAttempt(githubAttempt github.WorkflowRunAttempt) WorkflowRunAttempt {
	workflowRunAttempt := WorkflowRunAttempt{
		RunID:        strconv.FormatInt(githubAttempt.RunID, 10),
		Attempt:      githubAttempt.Attempt,
		Name:         githubAttempt.Name,
		Branch:       githubAttempt.Branch,
		HeadSHA:      githubAttempt.HeadSHA,
		RunStartedAt: githubAttempt.RunStartedAt,
		Status:       githubAttempt.Status,
		Conclusion:   string(githubAttempt.ProviderConclusion),
		WorkflowURL:  githubAttempt.WorkflowURL,
		PullRequests: make([]AssociatedPullRequest, 0, len(githubAttempt.PullRequests)),
	}
	for _, githubPullRequest := range githubAttempt.PullRequests {
		workflowRunAttempt.PullRequests = append(workflowRunAttempt.PullRequests, AssociatedPullRequest{
			Number:  githubPullRequest.PRNumber,
			URL:     githubPullRequest.URL,
			HeadSHA: githubPullRequest.HeadSHA,
		})
	}
	return workflowRunAttempt
}

func (provider *GitHubWorkflowProvider) ListWorkflowJobs(ctx context.Context, runID string, attempt int) ([]WorkflowJob, error) {
	parsedRunID, err := strconv.ParseInt(runID, 10, 64)
	if err != nil || parsedRunID < 1 {
		return nil, fmt.Errorf("invalid GitHub workflow run ID %q", runID)
	}
	githubJobs, err := provider.client.ListWorkflowJobs(ctx, parsedRunID, attempt)
	workflowJobs := make([]WorkflowJob, 0, len(githubJobs))
	for _, githubJob := range githubJobs {
		workflowJobs = append(workflowJobs, provider.workflowJob(githubJob))
	}
	return workflowJobs, err
}

func (provider *GitHubWorkflowProvider) workflowJob(githubJob github.WorkflowJob) WorkflowJob {
	workflowJob := WorkflowJob{
		ID:         strconv.FormatInt(githubJob.ID, 10),
		Name:       githubJob.Name,
		Conclusion: string(githubJob.ProviderConclusion),
		JobURL:     githubJob.JobURL,
		Steps:      make([]WorkflowStep, 0, len(githubJob.Steps)),
	}
	for _, githubStep := range githubJob.Steps {
		workflowJob.Steps = append(workflowJob.Steps, WorkflowStep{
			Number:     githubStep.Number,
			Name:       githubStep.Name,
			Conclusion: string(githubStep.ProviderConclusion),
		})
	}
	return workflowJob
}

func (provider *GitHubWorkflowProvider) StepLogs(ctx context.Context, jobID string) (map[int]string, error) {
	parsedJobID, err := strconv.ParseInt(jobID, 10, 64)
	if err != nil || parsedJobID < 1 {
		return nil, fmt.Errorf("invalid GitHub workflow job ID %q", jobID)
	}
	jobLog, err := provider.client.JobLog(ctx, parsedJobID)
	if err != nil {
		return nil, err
	}
	defer jobLog.Close()
	archiveData, err := io.ReadAll(io.LimitReader(jobLog, maxGitHubJobLogArchiveBytes+1))
	if err != nil {
		return nil, fmt.Errorf("could not read GitHub job log archive for job %q: %w", jobID, err)
	}
	if len(archiveData) > maxGitHubJobLogArchiveBytes {
		return nil, fmt.Errorf("GitHub job log archive for job %q exceeds %d bytes", jobID, maxGitHubJobLogArchiveBytes)
	}
	archiveReader, err := zip.NewReader(bytes.NewReader(archiveData), int64(len(archiveData)))
	if err != nil {
		return nil, fmt.Errorf("invalid GitHub job log archive for job %q: %w", jobID, err)
	}
	stepLogs := make(map[int]string)
	for _, archiveFile := range archiveReader.File {
		stepNumber, ok := provider.stepNumberFromLogFile(archiveFile.Name)
		if !ok {
			continue
		}
		if archiveFile.UncompressedSize64 > maxGitHubJobLogFileBytes {
			return nil, fmt.Errorf("GitHub job log file %q exceeds %d bytes", archiveFile.Name, maxGitHubJobLogFileBytes)
		}
		logFile, err := archiveFile.Open()
		if err != nil {
			return nil, fmt.Errorf("could not open GitHub job log file %q: %w", archiveFile.Name, err)
		}
		logData, readErr := io.ReadAll(io.LimitReader(logFile, maxGitHubJobLogFileBytes+1))
		closeErr := logFile.Close()
		if readErr != nil {
			return nil, fmt.Errorf("could not read GitHub job log file %q: %w", archiveFile.Name, readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("could not close GitHub job log file %q: %w", archiveFile.Name, closeErr)
		}
		if len(logData) > maxGitHubJobLogFileBytes {
			return nil, fmt.Errorf("GitHub job log file %q exceeds %d bytes", archiveFile.Name, maxGitHubJobLogFileBytes)
		}
		stepLogs[stepNumber] = string(logData)
	}
	return stepLogs, nil
}

func (provider *GitHubWorkflowProvider) stepNumberFromLogFile(fileName string) (int, bool) {
	baseName := path.Base(fileName)
	stepNumberText, _, found := strings.Cut(baseName, githubLogStepNameSeparator)
	if !found {
		return 0, false
	}
	stepNumber, err := strconv.Atoi(stepNumberText)
	return stepNumber, err == nil && stepNumber > 0
}
