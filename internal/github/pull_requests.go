package github

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	botLoginSuffix                = "[bot]"
	errorMissingPullRequestCIData = "GitHub response did not include CI/CD data for pull request #%d in %q"
	errorMissingPullRequestCIPage = "GitHub response indicated another CI/CD context page without an end cursor for pull request #%d in %q"
	errorMissingPullRequestData   = "GitHub response did not include pull request data for %q"
	errorMissingPullRequestPage   = "GitHub response indicated another pull request page without an end cursor for %q"
	errorQueryCICDResults         = "could not query CI/CD results for pull request #%d: %w"
)

type PullRequest struct {
	Author  *PullRequestAuthor `json:"author"`
	IsDraft bool               `json:"isDraft"`
	Number  int                `json:"number"`
	Title   string             `json:"title"`
	URL     string             `json:"url"`
}

type PullRequestAuthor struct {
	Login string `json:"login"`
}

type PullRequestCIResults struct {
	Contexts    []PullRequestCIContext `json:"contexts"`
	Number      int                    `json:"number"`
	RollupState string                 `json:"rollupState,omitempty"`
}

type PullRequestCIContext struct {
	CompletedAt string `json:"completedAt,omitempty"`
	Conclusion  string `json:"conclusion,omitempty"`
	Context     string `json:"context,omitempty"`
	Description string `json:"description,omitempty"`
	DetailsURL  string `json:"detailsUrl,omitempty"`
	IsRequired  bool   `json:"isRequired"`
	Name        string `json:"name,omitempty"`
	StartedAt   string `json:"startedAt,omitempty"`
	State       string `json:"state,omitempty"`
	Status      string `json:"status,omitempty"`
	Summary     string `json:"summary,omitempty"`
	TargetURL   string `json:"targetUrl,omitempty"`
	Text        string `json:"text,omitempty"`
	Type        string `json:"__typename"`
}

type pullRequestVariables struct {
	After *string `json:"after,omitempty"`
	Name  string  `json:"name"`
	Owner string  `json:"owner"`
}

type pullRequestCICDVariables struct {
	After  *string `json:"after,omitempty"`
	Name   string  `json:"name"`
	Number int     `json:"number"`
	Owner  string  `json:"owner"`
}

type pullRequestsPage struct {
	Repository *pullRequestsRepository `json:"repository"`
}

type pullRequestCICDPage struct {
	Repository *pullRequestCICDRepository `json:"repository"`
}

type pullRequestCICDRepository struct {
	PullRequest *pullRequestCICDData `json:"pullRequest"`
}

type pullRequestCICDData struct {
	StatusCheckRollup *pullRequestCICDRollup `json:"statusCheckRollup"`
}

type pullRequestCICDRollup struct {
	Contexts pullRequestCICDContextConnection `json:"contexts"`
	State    string                           `json:"state"`
}

type pullRequestCICDContextConnection struct {
	Nodes    []*PullRequestCIContext `json:"nodes"`
	PageInfo pullRequestPageInfo     `json:"pageInfo"`
}

type pullRequestsRepository struct {
	PullRequests pullRequestConnection `json:"pullRequests"`
}

type pullRequestConnection struct {
	Nodes    []PullRequest       `json:"nodes"`
	PageInfo pullRequestPageInfo `json:"pageInfo"`
}

type pullRequestPageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

func (c *Client) OpenPullRequests(ctx context.Context) ([]PullRequest, error) {
	return c.openPullRequests(ctx, openPullRequestsQuery, nil)
}

func (c *Client) OpenDependencyBumpPullRequests(ctx context.Context, authors []PullRequestAuthor) ([]PullRequest, error) {
	if len(authors) == 0 {
		return []PullRequest{}, nil
	}
	return c.openPullRequests(ctx, openDependencyBumpPullRequestsQuery, authors)
}

func (c *Client) AnalyzePullRequests(ctx context.Context, pullRequests []PullRequest) ([]PullRequestCIResults, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	analysisResults := make([]PullRequestCIResults, 0, len(pullRequests))
	queryErrors := make([]error, 0)
	for _, pullRequest := range pullRequests {
		if err := ctx.Err(); err != nil {
			queryErrors = append(queryErrors, err)
			return analysisResults, errors.Join(queryErrors...)
		}
		analysisResult, err := c.AnalyzePullRequest(ctx, pullRequest)
		if err != nil {
			queryErrors = append(queryErrors, fmt.Errorf(errorQueryCICDResults, pullRequest.Number, err))
			continue
		}
		analysisResults = append(analysisResults, analysisResult)
	}
	return analysisResults, errors.Join(queryErrors...)
}

func (c *Client) AnalyzePullRequest(ctx context.Context, pullRequest PullRequest) (PullRequestCIResults, error) {
	var after *string
	analysisResults := PullRequestCIResults{
		Contexts: make([]PullRequestCIContext, 0),
		Number:   pullRequest.Number,
	}

	for {
		var page pullRequestCICDPage
		variables := pullRequestCICDVariables{
			After:  after,
			Name:   c.repositoryName,
			Number: pullRequest.Number,
			Owner:  c.owner,
		}
		if err := c.queryWithVariables(ctx, pullRequestCICDResultsQuery, variables, &page); err != nil {
			return analysisResults, err
		}
		if page.Repository == nil || page.Repository.PullRequest == nil {
			return analysisResults, fmt.Errorf(errorMissingPullRequestCIData, pullRequest.Number, c.repositoryURL)
		}

		pullRequestData := page.Repository.PullRequest
		if pullRequestData.StatusCheckRollup == nil {
			return analysisResults, nil
		}
		rollup := pullRequestData.StatusCheckRollup
		analysisResults.RollupState = rollup.State
		for _, checkContext := range rollup.Contexts.Nodes {
			if checkContext != nil {
				analysisResults.Contexts = append(analysisResults.Contexts, *checkContext)
			}
		}
		if !rollup.Contexts.PageInfo.HasNextPage {
			return analysisResults, nil
		}
		if rollup.Contexts.PageInfo.EndCursor == emptyString {
			return analysisResults, fmt.Errorf(errorMissingPullRequestCIPage, pullRequest.Number, c.repositoryURL)
		}
		after = &rollup.Contexts.PageInfo.EndCursor
	}
}

func (c *Client) openPullRequests(ctx context.Context, query string, authors []PullRequestAuthor) ([]PullRequest, error) {
	var after *string
	authorLogins := make(map[string]struct{}, len(authors))
	pullRequests := make([]PullRequest, 0)
	for _, author := range authors {
		authorLogins[normalizeAuthorLogin(author.Login)] = struct{}{}
	}

	for {
		var page pullRequestsPage
		variables := pullRequestVariables{
			After: after,
			Name:  c.repositoryName,
			Owner: c.owner,
		}
		if err := c.queryWithVariables(ctx, query, variables, &page); err != nil {
			return nil, err
		}
		if page.Repository == nil {
			return nil, fmt.Errorf(errorMissingPullRequestData, c.repositoryURL)
		}

		connection := page.Repository.PullRequests
		for _, pullRequest := range connection.Nodes {
			if len(authorLogins) == 0 || matchesAuthor(pullRequest.Author, authorLogins) {
				pullRequests = append(pullRequests, pullRequest)
			}
		}
		if !connection.PageInfo.HasNextPage {
			return pullRequests, nil
		}
		if connection.PageInfo.EndCursor == emptyString {
			return nil, fmt.Errorf(errorMissingPullRequestPage, c.repositoryURL)
		}
		after = &connection.PageInfo.EndCursor
	}
}

func matchesAuthor(author *PullRequestAuthor, authorLogins map[string]struct{}) bool {
	if author == nil {
		return false
	}
	_, ok := authorLogins[normalizeAuthorLogin(author.Login)]
	return ok
}

func normalizeAuthorLogin(login string) string {
	return strings.TrimSuffix(strings.ToLower(login), botLoginSuffix)
}
