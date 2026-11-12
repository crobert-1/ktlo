package prstatus

import (
	"context"
	"errors"
	"io"

	"ktlo/internal/github"
	"ktlo/internal/output"
)

const (
	errorMissingGitHubClient = "GitHub client is required for the PR status check"
)

type Runner struct {
	githubClient *github.Client
}

func NewRunner(githubClient *github.Client) (*Runner, error) {
	if githubClient == nil {
		return nil, errors.New(errorMissingGitHubClient)
	}
	return &Runner{
		githubClient: githubClient,
	}, nil
}

func (r *Runner) Run(ctx context.Context, stdout io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	pullRequests, err := r.githubClient.OpenDependencyBumpPullRequests(ctx, dependencyUpdateAuthors[:])
	if err != nil {
		return err
	}
	analysisResults, err := r.githubClient.AnalyzePullRequests(ctx, pullRequests)
	if err != nil {
		return err
	}
	return output.PrintPullRequestCIResults(stdout, r.githubClient.TargetRepository(), pullRequests, analysisResults)
}
