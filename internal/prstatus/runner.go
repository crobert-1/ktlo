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

package prstatus

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"ktlo/internal/common"
	"ktlo/internal/failingtest"
	"ktlo/internal/github"
	"ktlo/internal/output"
	"ktlo/internal/telemetry"
	"ktlo/internal/testresults"
)

const (
	errorInvalidReportMaxAge        = "PR status report max age must not be negative"
	errorInvalidReportLookback      = "PR status lookback must be greater than zero"
	errorInvalidStaleAge            = "PR status stale age must be greater than zero"
	errorMissingGitHubClient        = "GitHub client is required for the PR status check"
	errorMissingWorkflowRunResolver = "workflow provider cannot resolve workflow runs from failed CI links"
	prStatusInstrumentationName     = "ktlo/internal/prstatus"
	spanPRStatusRun                 = "prstatus.run"
	githubHostName                  = "github.com"
	githubActionsPathSegment        = "actions"
	githubRunsPathSegment           = "runs"
	pathSeparator                   = "/"
)

type Runner struct {
	dependencyUpdateAuthors []common.PullRequestAuthor
	githubClient            *github.Client
	maxReportAge            time.Duration
	lookback                time.Duration
	parseWorkflowFailures   bool
	staleAge                time.Duration
	workflowProvider        *failingtest.GitHubWorkflowProvider
}

func (r *Runner) isStale(pullRequest github.PullRequest, now time.Time) bool {
	committedAt := pullRequest.HeadRef.Target.CommittedDate
	return !committedAt.IsZero() && committedAt.Before(now.Add(-r.staleAge))
}

type reportAnalysis struct {
	report            *testresults.Report
	generatedAt       time.Time
	maxAge            time.Duration
	unavailableReason testresults.UnknownReason
}

func NewRunner(githubClient *github.Client, dependencyUpdateAuthors []string, maxReportAge time.Duration, lookback time.Duration, staleAge time.Duration, parseWorkflowFailures bool) (*Runner, error) {
	if githubClient == nil {
		return nil, errors.New(errorMissingGitHubClient)
	}
	if maxReportAge < 0 {
		return nil, errors.New(errorInvalidReportMaxAge)
	}
	if lookback <= 0 {
		return nil, errors.New(errorInvalidReportLookback)
	}
	if staleAge <= 0 {
		return nil, errors.New(errorInvalidStaleAge)
	}
	authors := make([]common.PullRequestAuthor, 0, len(dependencyUpdateAuthors))
	for _, login := range dependencyUpdateAuthors {
		authors = append(authors, common.PullRequestAuthor{Login: login})
	}
	workflowProvider, err := failingtest.NewGitHubWorkflowProvider(githubClient)
	if err != nil {
		return nil, err
	}
	return &Runner{
		dependencyUpdateAuthors: authors,
		githubClient:            githubClient,
		maxReportAge:            maxReportAge,
		lookback:                lookback,
		parseWorkflowFailures:   parseWorkflowFailures,
		staleAge:                staleAge,
		workflowProvider:        workflowProvider,
	}, nil
}

func (r *Runner) Run(ctx context.Context, stdout io.Writer) (resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, prStatusInstrumentationName, spanPRStatusRun)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanPRStatusRun, startedAt, resultErr != nil)
	}()

	if err := ctx.Err(); err != nil {
		return err
	}
	pullRequests, err := r.githubClient.OpenDependencyBumpPullRequests(ctx, r.dependencyUpdateAuthors)
	if err != nil {
		return err
	}
	now := time.Now()
	activePullRequests := make([]github.PullRequest, 0, len(pullRequests))
	for _, pullRequest := range pullRequests {
		if r.isStale(pullRequest, now) {
			continue
		}
		activePullRequests = append(activePullRequests, pullRequest)
	}
	analysisResults, err := r.githubClient.AnalyzePullRequests(ctx, activePullRequests)
	if err != nil {
		return err
	}
	var reportState reportAnalysis
	if len(activePullRequests) > 0 {
		reportState = r.loadReportAnalysis()
	}
	results := make([]output.PullRequestResult, 0, len(pullRequests))
	for _, pullRequest := range pullRequests {
		isStale := r.isStale(pullRequest, now)
		pullRequestRef := testresults.PullRequestRef{
			Number:  pullRequest.Number,
			URL:     pullRequest.URL,
			HeadSHA: pullRequest.HeadSHA,
		}
		failureAssessments := []testresults.FailureAssessment(nil)
		unknownReason := reportState.unavailableReason
		if reportState.report != nil && !isStale {
			failureAssessments = reportState.report.AssessmentsForPullRequest(pullRequestRef)
			if unknownReason == "" && reportState.report.HasPullRequestHeadMismatch(&pullRequestRef) {
				unknownReason = testresults.UnknownReasonHeadSHAMismatch
			}
		}
		results = append(results, output.PullRequestResult{
			PullRequest:               pullRequest,
			FailureAssessments:        failureAssessments,
			IsStale:                   isStale,
			AnalysisUnavailableReason: unknownReason,
			ReportGeneratedAt:         reportState.generatedAt,
			MaxReportAge:              reportState.maxAge,
			Lookback:                  r.lookback,
		})
	}
	for _, analysisResult := range analysisResults {
		for resultIndex := range results {
			if results[resultIndex].PullRequest.Number == analysisResult.Number {
				results[resultIndex].CIResults = analysisResult
				break
			}
		}
	}
	if reportState.report != nil {
		if !r.parseWorkflowFailures {
			r.addReportAssessmentsForFailedCIContexts(reportState.report, results)
		}
	}
	if r.parseWorkflowFailures {
		parsedWorkflows, parseErrorsByPullRequest := r.parseWorkflowFailuresForPullRequests(ctx, results, reportState.report)
		if len(parsedWorkflows) > 0 {
			reportState.report = reportWithParsedWorkflowRuns(reportState.report, parsedWorkflows, r.githubClient.TargetRepository())
		}
		if len(parseErrorsByPullRequest) > 0 && reportState.report != nil {
			reportState.report.AnalysisStatus = testresults.AnalysisStatusIncomplete
		}
		for resultIndex := range results {
			result := &results[resultIndex]
			if result.IsStale {
				continue
			}
			if parseErrors := parseErrorsByPullRequest[result.PullRequest.Number]; len(parseErrors) > 0 {
				result.WorkflowFailureParseErrors = parseErrors
				if result.AnalysisUnavailableReason == "" {
					result.AnalysisUnavailableReason = testresults.UnknownReasonIncompleteLogCoverage
				}
			}
			if reportState.report != nil {
				pullRequest := testresults.PullRequestRef{
					Number:  result.PullRequest.Number,
					URL:     result.PullRequest.URL,
					HeadSHA: result.PullRequest.HeadSHA,
				}
				result.FailureAssessments = reportState.report.AssessmentsForPullRequest(pullRequest)
			}
		}
		if reportState.report != nil {
			r.addReportAssessmentsForFailedCIContexts(reportState.report, results)
		}
	}
	return output.PrintPullRequestCIResults(ctx, stdout, r.githubClient.TargetRepository(), results)
}

func (r *Runner) parseWorkflowFailuresForPullRequests(ctx context.Context, results []output.PullRequestResult, report *testresults.Report) ([]testresults.Workflow, map[int][]string) {
	parsedWorkflows := make([]testresults.Workflow, 0)
	parseErrorsByPullRequest := make(map[int][]string)
	latestAttemptProvider, ok := any(r.workflowProvider).(failingtest.LatestWorkflowRunAttemptProvider)
	if !ok {
		for _, result := range results {
			if !result.IsStale {
				parseErrorsByPullRequest[result.PullRequest.Number] = append(parseErrorsByPullRequest[result.PullRequest.Number], errorMissingWorkflowRunResolver)
			}
		}
		return parsedWorkflows, parseErrorsByPullRequest
	}
	checkRunNodeIDs := make([]string, 0)
	for _, result := range results {
		if result.IsStale {
			continue
		}
		for _, ciContext := range result.CIResults.Contexts {
			if isFailedCIContext(&ciContext) && ciContext.NodeID != "" {
				checkRunNodeIDs = append(checkRunNodeIDs, ciContext.NodeID)
			}
		}
	}
	workflowDetailsByCheckRunID, _ := r.githubClient.WorkflowRunDetailsForCheckRuns(ctx, checkRunNodeIDs)
	workflowDetailsByRunID := make(map[string]github.WorkflowRunCheckSuiteDetails)
	for _, workflowDetails := range workflowDetailsByCheckRunID {
		workflowDetailsByRunID[workflowDetails.Attempt.RunID] = workflowDetails
	}
	cutoff := time.Now().Add(-r.lookback)
	existingWorkflowRuns := workflowRunsByAttempt(report)
	parsedAttempts := make(map[workflowRunAttemptKey]parsedWorkflowAttempt)
	resolvedAttempts := make(map[string]workflowRunAttemptResolution)
	parsedAttemptKeys := make(map[workflowRunAttemptKey]struct{})
	for _, result := range results {
		if result.IsStale {
			continue
		}
		seenRunIDs := make(map[string]struct{})
		for _, ciContext := range result.CIResults.Contexts {
			if !isFailedCIContext(&ciContext) {
				continue
			}
			workflowDetails, hasWorkflowDetails := workflowDetailsByCheckRunID[ciContext.NodeID]
			runID := workflowDetails.Attempt.RunID
			if !hasWorkflowDetails {
				var found bool
				runID, found = r.workflowRunIDForCIContext(ciContext)
				if !found {
					continue
				}
				workflowDetails, hasWorkflowDetails = workflowDetailsByRunID[runID]
			}
			if _, exists := seenRunIDs[runID]; exists {
				continue
			}
			seenRunIDs[runID] = struct{}{}
			resolution, cached := resolvedAttempts[runID]
			if !cached {
				if hasWorkflowDetails && workflowDetails.IsComplete {
					resolution.attempt = workflowDetails.Attempt
					resolution.fromGraphQL = true
				} else {
					resolution.attempt, resolution.err = latestAttemptProvider.LatestWorkflowRunAttempt(ctx, runID)
				}
				resolvedAttempts[runID] = resolution
			}
			if resolution.err != nil {
				parseErrorsByPullRequest[result.PullRequest.Number] = append(parseErrorsByPullRequest[result.PullRequest.Number], resolution.err.Error())
				continue
			}
			if resolution.attempt.RunStartedAt.IsZero() || resolution.attempt.RunStartedAt.Before(cutoff) {
				continue
			}
			if !workflowRunAttemptMatchesPullRequest(&resolution.attempt, &result.PullRequest) {
				if !resolution.fromGraphQL {
					continue
				}
				resolution.attempt, resolution.err = latestAttemptProvider.LatestWorkflowRunAttempt(ctx, runID)
				resolution.fromGraphQL = false
				resolvedAttempts[runID] = resolution
				if resolution.err != nil {
					parseErrorsByPullRequest[result.PullRequest.Number] = append(parseErrorsByPullRequest[result.PullRequest.Number], resolution.err.Error())
					continue
				}
				if resolution.attempt.RunStartedAt.IsZero() || resolution.attempt.RunStartedAt.Before(cutoff) || !workflowRunAttemptMatchesPullRequest(&resolution.attempt, &result.PullRequest) {
					continue
				}
			}
			attemptKey := workflowRunAttemptKey{runID: resolution.attempt.RunID, attempt: resolution.attempt.Attempt}
			if _, exists := existingWorkflowRuns[attemptKey]; exists {
				continue
			}
			if len(resolution.attempt.PullRequests) == 0 {
				resolution.attempt.PullRequests = []common.PullRequestReference{{
					Number:  result.PullRequest.Number,
					URL:     result.PullRequest.URL,
					HeadSHA: resolution.attempt.HeadSHA,
				}}
			}
			parsedAttempt, cached := parsedAttempts[attemptKey]
			if !cached {
				workflowDetailsKey := workflowRunAttemptKey{runID: workflowDetails.Attempt.RunID, attempt: workflowDetails.Attempt.Attempt}
				if hasWorkflowDetails && workflowDetails.IsComplete && workflowDetailsKey == attemptKey {
					workflowJobsByRunID := map[string][]failingtest.WorkflowJob{resolution.attempt.RunID: workflowDetails.Jobs}
					parsedAttempt.workflows, parsedAttempt.err = failingtest.ParseWorkflowRunAttemptsWithJobs(ctx, r.workflowProvider, []failingtest.WorkflowRunAttempt{resolution.attempt}, workflowJobsByRunID)
				} else {
					parsedAttempt.workflows, parsedAttempt.err = failingtest.ParseWorkflowRunAttempts(ctx, r.workflowProvider, []failingtest.WorkflowRunAttempt{resolution.attempt})
				}
				parsedAttempts[attemptKey] = parsedAttempt
			}
			if _, added := parsedAttemptKeys[attemptKey]; !added {
				parsedWorkflows = append(parsedWorkflows, parsedAttempt.workflows...)
				parsedAttemptKeys[attemptKey] = struct{}{}
			}
			if parsedAttempt.err != nil {
				parseErrorsByPullRequest[result.PullRequest.Number] = append(parseErrorsByPullRequest[result.PullRequest.Number], parsedAttempt.err.Error())
			}
		}
	}
	return parsedWorkflows, parseErrorsByPullRequest
}

type workflowRunAttemptResolution struct {
	attempt     failingtest.WorkflowRunAttempt
	err         error
	fromGraphQL bool
}

type parsedWorkflowAttempt struct {
	workflows []testresults.Workflow
	err       error
}

func workflowRunsByAttempt(report *testresults.Report) map[workflowRunAttemptKey]testresults.Workflow {
	workflowRuns := make(map[workflowRunAttemptKey]testresults.Workflow)
	if report == nil || report.AnalysisStatus != testresults.AnalysisStatusComplete {
		return workflowRuns
	}
	for _, workflowRun := range report.Workflows {
		workflowRuns[workflowRunAttemptKey{runID: workflowRun.RunID, attempt: workflowRun.RunAttempt}] = workflowRun
	}
	return workflowRuns
}

func workflowRunAttemptMatchesPullRequest(attempt *failingtest.WorkflowRunAttempt, pullRequest *github.PullRequest) bool {
	if attempt == nil || pullRequest == nil {
		return false
	}
	for _, associatedPullRequest := range attempt.PullRequests {
		if associatedPullRequest.Number == pullRequest.Number && (associatedPullRequest.HeadSHA == "" || associatedPullRequest.HeadSHA == pullRequest.HeadSHA) {
			return true
		}
	}
	return len(attempt.PullRequests) == 0 && attempt.HeadSHA == pullRequest.HeadSHA
}

func (r *Runner) workflowRunIDForCIContext(ciContext github.PullRequestCIContext) (string, bool) {
	for _, contextURL := range []string{ciContext.DetailsURL, ciContext.TargetURL} {
		parsedURL, err := url.Parse(contextURL)
		if err != nil || !strings.EqualFold(parsedURL.Hostname(), githubHostName) {
			continue
		}
		pathSegments := strings.Split(strings.Trim(parsedURL.Path, pathSeparator), pathSeparator)
		if len(pathSegments) < 5 {
			continue
		}
		repository := r.githubClient.TargetRepository()
		if pathSegments[0] != repository.Owner || pathSegments[1] != repository.Name || pathSegments[2] != githubActionsPathSegment || pathSegments[3] != githubRunsPathSegment {
			continue
		}
		// The run ID is the fifth path segment in /owner/repository/actions/runs/{runID}.
		runID := pathSegments[4]
		parsedRunID, err := strconv.ParseInt(runID, 10, 64)
		if err == nil && parsedRunID > 0 {
			return runID, true
		}
	}
	return "", false
}

func reportWithParsedWorkflowRuns(report *testresults.Report, workflowRuns []testresults.Workflow, repository common.Repository) *testresults.Report {
	if report == nil {
		report = testresults.NewReport(testresults.ProviderGitHub, repository)
		report.AnalysisStatus = testresults.AnalysisStatusComplete
		report.GeneratedAt = time.Now().UTC()
	} else {
		reportCopy := *report
		reportCopy.Workflows = append([]testresults.Workflow(nil), report.Workflows...)
		report = &reportCopy
	}
	workflowRunsToReplace := make(map[workflowRunAttemptKey]struct{}, len(workflowRuns))
	for _, workflowRun := range workflowRuns {
		workflowRunsToReplace[workflowRunAttemptKey{runID: workflowRun.RunID, attempt: workflowRun.RunAttempt}] = struct{}{}
	}
	retainedWorkflowRuns := make([]testresults.Workflow, 0, len(report.Workflows)+len(workflowRuns))
	for _, workflowRun := range report.Workflows {
		if _, replace := workflowRunsToReplace[workflowRunAttemptKey{runID: workflowRun.RunID, attempt: workflowRun.RunAttempt}]; replace {
			continue
		}
		retainedWorkflowRuns = append(retainedWorkflowRuns, workflowRun)
	}
	report.Workflows = append(retainedWorkflowRuns, workflowRuns...)
	return report
}

type workflowRunAttemptKey struct {
	runID   string
	attempt int
}

func (r *Runner) addReportAssessmentsForFailedCIContexts(report *testresults.Report, results []output.PullRequestResult) {
	for resultIndex := range results {
		result := &results[resultIndex]
		if result.IsStale {
			continue
		}
		pullRequest := testresults.PullRequestRef{
			Number:  result.PullRequest.Number,
			URL:     result.PullRequest.URL,
			HeadSHA: result.PullRequest.HeadSHA,
		}
		assessedLocations := make(map[testresults.CIFailureLocation]struct{}, len(result.FailureAssessments))
		for _, assessment := range result.FailureAssessments {
			assessedLocations[assessment.Location] = struct{}{}
		}
		for contextIndex := range result.CIResults.Contexts {
			ciContext := &result.CIResults.Contexts[contextIndex]
			if !isFailedCIContext(ciContext) {
				continue
			}
			assessment, found := report.AssessCIContextFailure(ciFailureContextReference(ciContext), &pullRequest)
			if !found {
				continue
			}
			if _, exists := assessedLocations[assessment.Location]; exists {
				continue
			}
			assessedLocations[assessment.Location] = struct{}{}
			result.FailureAssessments = append(result.FailureAssessments, assessment)
		}
	}
}

func isFailedCIContext(ciContext *github.PullRequestCIContext) bool {
	if ciContext == nil {
		return false
	}
	state := ciContext.EffectiveState()
	return state == github.PullRequestCIRollupStateFailure || state == github.PullRequestCIRollupStateError
}

func ciFailureContextReference(ciContext *github.PullRequestCIContext) testresults.CIContextFailure {
	failureMessage := ciContext.Description
	if failureMessage == "" {
		failureMessage = ciContext.Summary
	}
	if failureMessage == "" {
		failureMessage = ciContext.Text
	}
	runStartedAt, _ := time.Parse(time.RFC3339Nano, ciContext.StartedAt)
	if runStartedAt.IsZero() {
		runStartedAt, _ = time.Parse(time.RFC3339Nano, ciContext.CompletedAt)
	}
	return testresults.CIContextFailure{
		Name:           ciContext.NameForDisplay(),
		DetailsURL:     ciContext.DetailsURL,
		TargetURL:      ciContext.TargetURL,
		FailureMessage: failureMessage,
		RunStartedAt:   runStartedAt,
	}
}

func (r *Runner) loadReportAnalysis() reportAnalysis {
	analysis := reportAnalysis{maxAge: r.maxReportAge}
	report := testresults.NewReport(testresults.ProviderGitHub, r.githubClient.TargetRepository())
	if err := report.Load(); err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			analysis.unavailableReason = testresults.UnknownReasonMissingReport
		case errors.Is(err, testresults.ErrReportIdentityMismatch):
			analysis.unavailableReason = testresults.UnknownReasonRepositoryMismatch
		case errors.Is(err, testresults.ErrMalformedReport):
			analysis.unavailableReason = testresults.UnknownReasonMalformedReport
		default:
			analysis.unavailableReason = testresults.UnknownReasonReportReadError
		}
		return analysis
	}
	analysis.generatedAt = report.GeneratedAt
	if analysis.maxAge > 0 && time.Since(report.GeneratedAt) > analysis.maxAge {
		analysis.unavailableReason = testresults.UnknownReasonStaleReport
		return analysis
	}
	analysis.report = report
	if report.AnalysisStatus == testresults.AnalysisStatusIncomplete {
		analysis.unavailableReason = testresults.UnknownReasonIncompleteLogCoverage
	}
	return analysis
}
