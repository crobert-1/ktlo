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
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"ktlo/internal/common"
	"ktlo/internal/telemetry"
)

const (
	botLoginSuffix                        = "[bot]"
	errorMissingPullRequestCIData         = "GitHub response did not include CI/CD data for pull request #%d in %q"
	errorMissingPullRequestCIPage         = "GitHub response indicated another CI/CD context page without an end cursor for pull request #%d in %q"
	errorMissingPullRequestData           = "GitHub response did not include pull request data for %q"
	errorMissingPullRequestPage           = "GitHub response indicated another pull request page without an end cursor for %q"
	errorQueryCICDResults                 = "could not query CI/CD results for pull request #%d: %w"
	spanAnalyzePullRequest                = "github.analyze_pull_request"
	spanAnalyzePullRequests               = "github.analyze_pull_requests"
	spanListPullRequests                  = "github.list_pull_requests"
	traceAttributePullRequest             = "github.pull_request.number"
	traceAttributePullRequestCount        = "github.pull_request.count"
	traceAttributePullRequestContextCount = "github.pull_request.ci_context.count"
	traceAttributePullRequestPageCount    = "github.pull_request.page.count"
	traceAttributePullRequestResultCount  = "github.pull_request.result.count"
)

type PullRequest struct {
	Author  *common.PullRequestAuthor `json:"author"`
	HeadRef PullRequestHeadRef        `json:"headRef"`
	HeadSHA string                    `json:"headSha"`
	IsDraft bool                      `json:"isDraft"`
	Number  int                       `json:"number"`
	Title   string                    `json:"title"`
	URL     string                    `json:"url"`
}

type PullRequestHeadRef struct {
	Target PullRequestHeadCommit `json:"target"`
}

type PullRequestHeadCommit struct {
	CommittedDate time.Time `json:"committedDate"`
}

type PullRequestCIRollupState string

const (
	PullRequestCIRollupStateCancelled  PullRequestCIRollupState = "CANCELLED"
	PullRequestCIRollupStateError      PullRequestCIRollupState = "ERROR"
	PullRequestCIRollupStateExpected   PullRequestCIRollupState = "EXPECTED"
	PullRequestCIRollupStateFailure    PullRequestCIRollupState = "FAILURE"
	PullRequestCIRollupStateInProgress PullRequestCIRollupState = "IN_PROGRESS"
	PullRequestCIRollupStatePending    PullRequestCIRollupState = "PENDING"
	PullRequestCIRollupStateQueued     PullRequestCIRollupState = "QUEUED"
	PullRequestCIRollupStateRequested  PullRequestCIRollupState = "REQUESTED"
	PullRequestCIRollupStateSkipped    PullRequestCIRollupState = "SKIPPED"
	PullRequestCIRollupStateSuccess    PullRequestCIRollupState = "SUCCESS"
	PullRequestCIRollupStateUnknown    PullRequestCIRollupState = "UNKNOWN"
	PullRequestCIRollupStateWaiting    PullRequestCIRollupState = "WAITING"
)

func (state PullRequestCIRollupState) IsValid() bool {
	switch state {
	case PullRequestCIRollupStateCancelled, PullRequestCIRollupStateError, PullRequestCIRollupStateExpected,
		PullRequestCIRollupStateFailure, PullRequestCIRollupStateInProgress, PullRequestCIRollupStatePending,
		PullRequestCIRollupStateQueued, PullRequestCIRollupStateRequested, PullRequestCIRollupStateSkipped,
		PullRequestCIRollupStateSuccess, PullRequestCIRollupStateUnknown, PullRequestCIRollupStateWaiting:
		return true
	default:
		return false
	}
}

type PullRequestCIResults struct {
	Contexts    []PullRequestCIContext `json:"contexts"`
	Number      int                    `json:"number"`
	RollupState string                 `json:"rollupState,omitempty"`
}

func (results *PullRequestCIResults) EffectiveRollupState() PullRequestCIRollupState {
	if results == nil {
		return PullRequestCIRollupStateUnknown
	}
	if len(results.Contexts) == 0 {
		if strings.TrimSpace(results.RollupState) == emptyString {
			return PullRequestCIRollupState(emptyString)
		}
		return results.providerRollupState()
	}

	var hasCancelledContext bool
	var hasErrorContext bool
	var hasFailureContext bool
	var hasPendingContext bool
	var hasSkippedContext bool
	var hasUnknownContext bool
	for contextIndex := range results.Contexts {
		contextState := results.Contexts[contextIndex].EffectiveState()
		switch contextState {
		case PullRequestCIRollupStateCancelled:
			hasCancelledContext = true
		case PullRequestCIRollupStateError:
			hasErrorContext = true
		case PullRequestCIRollupStateFailure:
			hasFailureContext = true
		case PullRequestCIRollupStateExpected, PullRequestCIRollupStateInProgress, PullRequestCIRollupStatePending,
			PullRequestCIRollupStateQueued, PullRequestCIRollupStateRequested, PullRequestCIRollupStateWaiting:
			hasPendingContext = true
		case PullRequestCIRollupStateSkipped:
			hasSkippedContext = true
		case PullRequestCIRollupStateSuccess:
		default:
			hasUnknownContext = true
		}
	}
	if hasErrorContext {
		return PullRequestCIRollupStateError
	}
	if hasFailureContext {
		return PullRequestCIRollupStateFailure
	}
	if hasPendingContext {
		return PullRequestCIRollupStatePending
	}
	if hasUnknownContext {
		return PullRequestCIRollupStateUnknown
	}
	if hasCancelledContext {
		return PullRequestCIRollupStateCancelled
	}
	if hasSkippedContext {
		return PullRequestCIRollupStateSkipped
	}
	return PullRequestCIRollupStateSuccess
}

func (results *PullRequestCIResults) HasPendingChecks() bool {
	if results == nil {
		return false
	}
	switch results.providerRollupState() {
	case PullRequestCIRollupStateExpected, PullRequestCIRollupStateInProgress, PullRequestCIRollupStatePending,
		PullRequestCIRollupStateQueued, PullRequestCIRollupStateRequested, PullRequestCIRollupStateWaiting:
		return true
	}
	for contextIndex := range results.Contexts {
		if results.Contexts[contextIndex].EffectiveState() == PullRequestCIRollupStatePending {
			return true
		}
	}
	return false
}

func (results *PullRequestCIResults) providerRollupState() PullRequestCIRollupState {
	if results == nil {
		return PullRequestCIRollupStateUnknown
	}
	state := PullRequestCIRollupState(strings.ToUpper(strings.TrimSpace(results.RollupState)))
	if !state.IsValid() {
		return PullRequestCIRollupStateUnknown
	}
	return state
}

type PullRequestCIContext struct {
	CompletedAt string `json:"completedAt,omitempty"`
	Conclusion  string `json:"conclusion,omitempty"`
	Context     string `json:"context,omitempty"`
	Description string `json:"description,omitempty"`
	DetailsURL  string `json:"detailsUrl,omitempty"`
	NodeID      string `json:"id,omitempty"`
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

func (context *PullRequestCIContext) EffectiveState() PullRequestCIRollupState {
	if context == nil {
		return PullRequestCIRollupStateUnknown
	}
	state := strings.ToUpper(strings.TrimSpace(context.Conclusion))
	if state == emptyString {
		state = strings.ToUpper(strings.TrimSpace(context.State))
	}
	if state == emptyString {
		state = strings.ToUpper(strings.TrimSpace(context.Status))
	}
	switch state {
	case string(PullRequestCIRollupStateCancelled):
		return PullRequestCIRollupStateCancelled
	case string(PullRequestCIRollupStateError):
		return PullRequestCIRollupStateError
	case string(PullRequestCIRollupStateFailure):
		return PullRequestCIRollupStateFailure
	case string(PullRequestCIRollupStateExpected), string(PullRequestCIRollupStateInProgress),
		string(PullRequestCIRollupStatePending), string(PullRequestCIRollupStateQueued),
		string(PullRequestCIRollupStateRequested), string(PullRequestCIRollupStateWaiting):
		return PullRequestCIRollupStatePending
	case string(PullRequestCIRollupStateSkipped):
		return PullRequestCIRollupStateSkipped
	case string(PullRequestCIRollupStateSuccess):
		return PullRequestCIRollupStateSuccess
	default:
		return PullRequestCIRollupStateUnknown
	}
}

func (context *PullRequestCIContext) NameForDisplay() string {
	if context == nil {
		return ""
	}
	if context.Name != emptyString {
		return context.Name
	}
	return context.Context
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

func (c *Client) OpenDependencyBumpPullRequests(ctx context.Context, authors []common.PullRequestAuthor) ([]PullRequest, error) {
	if len(authors) == 0 {
		return []PullRequest{}, nil
	}
	return c.openPullRequests(ctx, openDependencyBumpPullRequestsQuery, authors)
}

func (c *Client) AnalyzePullRequests(ctx context.Context, pullRequests []PullRequest) (analysisResults []PullRequestCIResults, resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, githubInstrumentationName, spanAnalyzePullRequests, attribute.Int(traceAttributePullRequestCount, len(pullRequests)))
	defer func() {
		telemetry.RecordDuration(ctx, spanAnalyzePullRequests, startedAt)
		if resultErr != nil {
			span.SetStatus(codes.Error, emptyString)
		}
		span.SetAttributes(attribute.Int(traceAttributePullRequestResultCount, len(analysisResults)))
		span.End()
	}()
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	analysisResults = make([]PullRequestCIResults, 0, len(pullRequests))
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

func (c *Client) AnalyzePullRequest(ctx context.Context, pullRequest PullRequest) (analysisResults PullRequestCIResults, resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, githubInstrumentationName, spanAnalyzePullRequest, attribute.Int(traceAttributePullRequest, pullRequest.Number))
	defer func() {
		telemetry.RecordDuration(ctx, spanAnalyzePullRequest, startedAt)
		if resultErr != nil {
			span.SetStatus(codes.Error, emptyString)
		}
		span.SetAttributes(attribute.Int(traceAttributePullRequestContextCount, len(analysisResults.Contexts)))
		span.End()
	}()
	var after *string
	analysisResults = PullRequestCIResults{
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

func (c *Client) openPullRequests(ctx context.Context, query string, authors []common.PullRequestAuthor) (pullRequests []PullRequest, resultErr error) {
	pageCount := 0
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, githubInstrumentationName, spanListPullRequests)
	defer func() {
		telemetry.RecordDuration(ctx, spanListPullRequests, startedAt)
		if resultErr != nil {
			span.SetStatus(codes.Error, emptyString)
		}
		span.SetAttributes(
			attribute.Int(traceAttributePullRequestCount, len(pullRequests)),
			attribute.Int(traceAttributePullRequestPageCount, pageCount),
		)
		span.End()
	}()
	var after *string
	authorLogins := make(map[string]struct{}, len(authors))
	pullRequests = make([]PullRequest, 0)
	for _, author := range authors {
		authorLogins[normalizeAuthorLogin(author.Login)] = struct{}{}
	}

	for {
		pageCount++
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

func matchesAuthor(author *common.PullRequestAuthor, authorLogins map[string]struct{}) bool {
	if author == nil {
		return false
	}
	_, ok := authorLogins[normalizeAuthorLogin(author.Login)]
	return ok
}

func normalizeAuthorLogin(login string) string {
	return strings.TrimSuffix(strings.ToLower(login), botLoginSuffix)
}
