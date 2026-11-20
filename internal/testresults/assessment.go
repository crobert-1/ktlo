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

package testresults

import (
	"errors"
	"math"
	"strings"
	"time"
)

type Classification string

func (classification *Classification) UnmarshalJSON(data []byte) error {
	value, err := normalizeEnumJSON(data)
	if err != nil {
		return err
	}
	*classification = Classification(value)
	return nil
}

const (
	ClassificationBroken  Classification = "broken"
	ClassificationFlaky   Classification = "flaky"
	ClassificationUnknown Classification = "unknown"
)

func (classification Classification) IsValid() bool {
	switch classification {
	case ClassificationBroken, ClassificationFlaky, ClassificationUnknown:
		return true
	default:
		return false
	}
}

type UnknownReason string

func (reason *UnknownReason) UnmarshalJSON(data []byte) error {
	value, err := normalizeEnumJSON(data)
	if err != nil {
		return err
	}
	*reason = UnknownReason(value)
	return nil
}

const (
	UnknownReasonHeadSHAMismatch       UnknownReason = "pull_request_head_sha_mismatch"
	UnknownReasonIncompleteLogCoverage UnknownReason = "incomplete_log_coverage"
	UnknownReasonInsufficientHistory   UnknownReason = "insufficient_matching_history"
	UnknownReasonMalformedReport       UnknownReason = "malformed_report"
	UnknownReasonMissingReport         UnknownReason = "missing_report"
	UnknownReasonReportReadError       UnknownReason = "report_read_error"
	UnknownReasonRepositoryMismatch    UnknownReason = "repository_mismatch"
	UnknownReasonUnidentifiableFailure UnknownReason = "unidentifiable_failure"
)

func (reason UnknownReason) IsValid() bool {
	switch reason {
	case UnknownReasonHeadSHAMismatch, UnknownReasonIncompleteLogCoverage, UnknownReasonInsufficientHistory, UnknownReasonMalformedReport, UnknownReasonMissingReport, UnknownReasonReportReadError, UnknownReasonRepositoryMismatch, UnknownReasonUnidentifiableFailure:
		return true
	default:
		return false
	}
}

type CIFailureLocation struct {
	WorkflowName string
	JobName      string
	StepName     string
	TestID       string
	TestKind     TestKind
}

func (location *CIFailureLocation) specificity() int {
	if location.TestID != "" {
		return 4
	}
	if location.StepName != "" {
		return 3
	}
	if location.JobName != "" {
		return 2
	}
	if location.WorkflowName != "" {
		return 1
	}
	return 0
}

func (location *CIFailureLocation) testIdentity() CIFailureLocation {
	if location.TestID == "" {
		return *location
	}
	return CIFailureLocation{
		WorkflowName: location.WorkflowName,
		TestID:       location.TestID,
		TestKind:     location.TestKind,
	}
}

type CIContextFailure struct {
	Name           string
	DetailsURL     string
	TargetURL      string
	FailureMessage string
	RunStartedAt   time.Time
}

func (context *CIContextFailure) matchesName(name string) bool {
	return context != nil && name != "" && strings.EqualFold(strings.TrimSpace(context.Name), strings.TrimSpace(name))
}

func (context *CIContextFailure) matchesURL(url string) bool {
	if context == nil || url == "" {
		return false
	}
	url = strings.TrimSpace(url)
	return url == strings.TrimSpace(context.DetailsURL) || url == strings.TrimSpace(context.TargetURL)
}

func (context *CIContextFailure) mostSpecificLocation(locations map[CIFailureLocation]struct{}) (location CIFailureLocation, found bool) {
	var highestSpecificity int
	var highestSpecificityCount int
	for candidate := range locations {
		candidateSpecificity := candidate.specificity()
		if candidateSpecificity > highestSpecificity {
			location = candidate
			highestSpecificity = candidateSpecificity
			highestSpecificityCount = 1
			continue
		}
		if candidateSpecificity == highestSpecificity {
			highestSpecificityCount++
		}
	}
	return location, highestSpecificityCount == 1 && highestSpecificity > 0
}

func (location *CIFailureLocation) classifyFailure(history OutcomeHistory, pullRequest *PullRequestRef, isAnalysisIncomplete bool) (Classification, UnknownReason) {
	mainCounts := history.MainBranch.Counts
	mainHasHistory := mainCounts.Passes+mainCounts.Failures > 0
	if mainCounts.hasMixedOutcomes() {
		return ClassificationFlaky, ""
	}

	var currentCommitCounts OutcomeCounts
	var currentCommitFound bool
	if pullRequest != nil {
		for _, commit := range history.PRCommits {
			if commit.CommitSHA == pullRequest.HeadSHA {
				currentCommitCounts = commit.Counts
				currentCommitFound = true
				if commit.Counts.hasMixedOutcomes() {
					return ClassificationFlaky, ""
				}
				break
			}
		}
		if !currentCommitFound && len(history.PRCommits) > 0 {
			return ClassificationUnknown, UnknownReasonHeadSHAMismatch
		}
	} else {
		currentCommitCounts = mainCounts
		currentCommitFound = mainHasHistory
	}

	if !mainHasHistory {
		return ClassificationUnknown, UnknownReasonInsufficientHistory
	}
	if !currentCommitFound || currentCommitCounts.Failures == 0 {
		if location.WorkflowName == "" && location.JobName == "" && location.StepName == "" && location.TestID == "" {
			return ClassificationUnknown, UnknownReasonUnidentifiableFailure
		}
		return ClassificationUnknown, UnknownReasonInsufficientHistory
	}
	if currentCommitCounts.Passes == 0 {
		return ClassificationBroken, ""
	}
	if isAnalysisIncomplete {
		return ClassificationUnknown, UnknownReasonIncompleteLogCoverage
	}
	return ClassificationUnknown, UnknownReasonInsufficientHistory
}

type FailureAssessment struct {
	PullRequest    *PullRequestRef
	Location       CIFailureLocation
	Classification Classification
	History        OutcomeHistory
	RelatedCI      []CIContextHistory
	Partial        bool
	UnknownReason  UnknownReason
}

func (assessment *FailureAssessment) LatestFailureEvidence() *FailureEvidence {
	return assessment.History.latestFailureFor(assessment.PullRequest)
}

type FailureEvidence struct {
	SourceRun RunReference
	Message   string
}

type CIContextHistory struct {
	Location CIFailureLocation
	History  OutcomeHistory
}

const mainBranchName = "main"

type OutcomeHistory struct {
	MainBranch OutcomeSummary
	PRCommits  []CommitOutcomeSummary
}

func (history *OutcomeHistory) InferredWorkflowPassCount() int {
	passCount := history.MainBranch.Counts.InferredPasses
	for _, commit := range history.PRCommits {
		passCount += commit.Counts.InferredPasses
	}
	return passCount
}

func (history *OutcomeHistory) addFailureForPullRequest(pullRequest *PullRequestRef, failure FailureEvidence) {
	if history == nil || pullRequest == nil {
		return
	}
	failure.SourceRun.CommitSHA = pullRequest.HeadSHA
	for commitIndex := range history.PRCommits {
		commit := &history.PRCommits[commitIndex]
		if commit.CommitSHA == pullRequest.HeadSHA {
			commit.addFailure(failure)
			return
		}
	}
	commit := CommitOutcomeSummary{CommitSHA: pullRequest.HeadSHA}
	commit.addFailure(failure)
	history.PRCommits = append(history.PRCommits, commit)
}

func (history *OutcomeHistory) latestFailureFor(pullRequest *PullRequestRef) *FailureEvidence {
	if pullRequest == nil {
		return history.MainBranch.LatestFailure
	}
	for commitIndex := range history.PRCommits {
		commit := &history.PRCommits[commitIndex]
		if commit.CommitSHA == pullRequest.HeadSHA {
			return commit.LatestFailure
		}
	}
	return nil
}

type historyObservation struct {
	outcome                     Outcome
	jobName                     string
	jobURL                      string
	message                     string
	inferredFromWorkflowSuccess bool
}

type OutcomeSummary struct {
	Counts        OutcomeCounts
	LatestPass    *RunReference
	LatestFailure *FailureEvidence
}

type CommitOutcomeSummary struct {
	Branch        string
	CommitSHA     string
	Counts        OutcomeCounts
	LatestPass    *RunReference
	LatestFailure *FailureEvidence
}

func (summary *CommitOutcomeSummary) addFailure(failure FailureEvidence) {
	summary.Counts.Failures++
	if summary.LatestFailure == nil || failure.SourceRun.isLaterThan(&summary.LatestFailure.SourceRun) {
		summary.LatestFailure = &failure
	}
}

type OutcomeCounts struct {
	Passes         int
	Failures       int
	InferredPasses int
}

func (counts *OutcomeCounts) PassPercentage() (int, error) {
	total := counts.Passes + counts.Failures
	if total == 0 {
		return 0, errors.New("no pass or failure observations")
	}
	return int(math.Round(float64(counts.Passes) / float64(total) * 100)), nil
}

func (counts *OutcomeCounts) hasMixedOutcomes() bool {
	return counts.Passes > 0 && counts.Failures > 0
}

type RunReference struct {
	RunID        string
	RunAttempt   int
	Branch       string
	CommitSHA    string
	RunStartedAt time.Time
	WorkflowURL  string
	JobURL       string
}

func (reference *RunReference) isLaterThan(other *RunReference) bool {
	if !reference.RunStartedAt.Equal(other.RunStartedAt) {
		return reference.RunStartedAt.After(other.RunStartedAt)
	}
	if reference.RunAttempt != other.RunAttempt {
		return reference.RunAttempt > other.RunAttempt
	}
	return reference.RunID > other.RunID
}

type FailureSelector struct {
	TestID       string
	WorkflowName string
	JobName      string
}

func (selector *FailureSelector) Equals(location CIFailureLocation) bool {
	if selector == nil {
		return true
	}
	switch {
	case selector.TestID != "":
		return location.TestID == selector.TestID
	case selector.WorkflowName != "":
		return location.WorkflowName == selector.WorkflowName
	case selector.JobName != "":
		return location.JobName == selector.JobName
	default:
		return false
	}
}

func (selector *FailureSelector) Validate() error {
	if selector == nil {
		return errors.New("selector is required")
	}
	selectorCount := 0
	if selector.TestID != "" {
		selectorCount++
	}
	if selector.WorkflowName != "" {
		selectorCount++
	}
	if selector.JobName != "" {
		selectorCount++
	}
	if selectorCount != 1 {
		return errors.New("selector must specify exactly one of test ID, workflow name, or job name")
	}
	return nil
}

func (r *Report) AssessFailure(location CIFailureLocation, pullRequest *PullRequestRef) FailureAssessment {
	history := r.GetOutcomeHistory(location, pullRequest)
	return r.assessFailureWithHistory(location, pullRequest, history)
}

func (r *Report) assessFailureWithHistory(location CIFailureLocation, pullRequest *PullRequestRef, history OutcomeHistory) FailureAssessment {
	isAnalysisIncomplete := r.AnalysisStatus == AnalysisStatusIncomplete
	classification, unknownReason := location.classifyFailure(history, pullRequest, isAnalysisIncomplete)
	return FailureAssessment{
		PullRequest:    pullRequest,
		Location:       location,
		Classification: classification,
		History:        history,
		RelatedCI:      r.relatedCIHistories(location, pullRequest),
		Partial:        isAnalysisIncomplete,
		UnknownReason:  unknownReason,
	}
}

func (r *Report) AssessCIContextFailure(context CIContextFailure, pullRequest *PullRequestRef) (FailureAssessment, bool) {
	if r == nil || pullRequest == nil {
		return FailureAssessment{}, false
	}
	location, found := r.locationForCIContext(&context, pullRequest)
	if !found {
		return FailureAssessment{}, false
	}
	history := r.GetOutcomeHistory(location, pullRequest)
	runStartedAt := context.RunStartedAt
	if runStartedAt.IsZero() {
		runStartedAt = time.Now().UTC()
	}
	history.addFailureForPullRequest(pullRequest, FailureEvidence{
		SourceRun: RunReference{
			CommitSHA:    pullRequest.HeadSHA,
			RunStartedAt: runStartedAt,
			WorkflowURL:  context.TargetURL,
			JobURL:       context.DetailsURL,
		},
		Message: context.FailureMessage,
	})
	isAnalysisIncomplete := r.AnalysisStatus == AnalysisStatusIncomplete
	classification, unknownReason := location.classifyFailure(history, pullRequest, isAnalysisIncomplete)
	return FailureAssessment{
		PullRequest:    pullRequest,
		Location:       location,
		Classification: classification,
		History:        history,
		RelatedCI:      r.relatedCIHistories(location, pullRequest),
		Partial:        isAnalysisIncomplete,
		UnknownReason:  unknownReason,
	}, true
}

func (r *Report) locationForCIContext(context *CIContextFailure, pullRequest *PullRequestRef) (CIFailureLocation, bool) {
	urlMatches := make(map[CIFailureLocation]struct{})
	nameMatches := make(map[CIFailureLocation]struct{})
	for workflowIndex := range r.Workflows {
		workflow := &r.Workflows[workflowIndex]
		if workflow.Branch != mainBranchName && (workflow.PullRequest == nil || workflow.PullRequest.Number != pullRequest.Number) {
			continue
		}
		workflowLocation := CIFailureLocation{WorkflowName: workflow.Name}
		if context.matchesURL(workflow.WorkflowURL) {
			urlMatches[workflowLocation] = struct{}{}
		}
		if context.matchesName(workflow.Name) {
			nameMatches[workflowLocation] = struct{}{}
		}
		for jobIndex := range workflow.Jobs {
			job := &workflow.Jobs[jobIndex]
			jobLocation := CIFailureLocation{WorkflowName: workflow.Name, JobName: job.Name}
			if context.matchesURL(job.JobURL) {
				urlMatches[jobLocation] = struct{}{}
			}
			if context.matchesName(job.Name) {
				nameMatches[jobLocation] = struct{}{}
			}
			for stepIndex := range job.Steps {
				step := &job.Steps[stepIndex]
				stepLocation := CIFailureLocation{WorkflowName: workflow.Name, JobName: job.Name, StepName: step.Name}
				if context.matchesName(step.Name) {
					nameMatches[stepLocation] = struct{}{}
				}
				for _, test := range step.Tests {
					testLocation := CIFailureLocation{WorkflowName: workflow.Name, JobName: job.Name, StepName: step.Name, TestID: test.ID, TestKind: test.Kind}
					testLocation = (&testLocation).testIdentity()
					if context.matchesName(test.ID) {
						nameMatches[testLocation] = struct{}{}
					}
				}
			}
			for _, test := range job.Tests {
				testLocation := CIFailureLocation{WorkflowName: workflow.Name, JobName: job.Name, TestID: test.ID, TestKind: test.Kind}
				testLocation = (&testLocation).testIdentity()
				if context.matchesName(test.ID) {
					nameMatches[testLocation] = struct{}{}
				}
			}
		}
		for _, test := range workflow.Tests {
			testLocation := CIFailureLocation{WorkflowName: workflow.Name, TestID: test.ID, TestKind: test.Kind}
			testLocation = (&testLocation).testIdentity()
			if context.matchesName(test.ID) {
				nameMatches[testLocation] = struct{}{}
			}
		}
	}
	if location, found := context.mostSpecificLocation(urlMatches); found {
		return location, true
	}
	return context.mostSpecificLocation(nameMatches)
}

func (r *Report) FindAssessments(selector FailureSelector) ([]FailureAssessment, error) {
	if err := selector.Validate(); err != nil {
		return nil, err
	}
	return r.findAssessments(&selector, nil), nil
}

func (r *Report) AssessmentsForPullRequest(pullRequest PullRequestRef) []FailureAssessment {
	return r.findAssessments(nil, &pullRequest)
}

func (r *Report) AssessmentsForBranch(branch string) []FailureAssessment {
	if r == nil {
		return nil
	}
	locations := make(map[CIFailureLocation]struct{})
	for workflowIndex := range r.Workflows {
		workflow := &r.Workflows[workflowIndex]
		if workflow.Branch != branch {
			continue
		}
		for _, location := range workflow.failureLocations() {
			location = (&location).testIdentity()
			locations[location] = struct{}{}
		}
	}
	assessments := make([]FailureAssessment, 0, len(locations))
	for location := range locations {
		history := r.getOutcomeHistory(location, nil, true)
		if (location.TestKind == TestKindPackage || location.TestID == "") && history.MainBranch.Counts.Failures == 0 {
			continue
		}
		assessments = append(assessments, r.assessFailureWithHistory(location, nil, history))
	}
	return assessments
}

func (r *Report) GetOutcomeHistory(location CIFailureLocation, pullRequest *PullRequestRef) OutcomeHistory {
	return r.getOutcomeHistory(location, pullRequest, false)
}

func (r *Report) getOutcomeHistory(location CIFailureLocation, pullRequest *PullRequestRef, omitExplainedFailures bool) OutcomeHistory {
	if location.TestID != "" {
		location = (&location).testIdentity()
	}
	commitAccumulators := make(map[string]*commitAccumulator)
	commitOrder := make([]string, 0)
	mainAccumulator := outcomeAccumulator{}
	// workflowIndex is the current position in r.Workflows.
	for workflowIndex := range r.Workflows {
		workflow := &r.Workflows[workflowIndex]
		observations := workflow.observationsFor(location)
		if len(observations) == 0 {
			continue
		}
		for _, observation := range observations {
			if omitExplainedFailures && observation.outcome == OutcomeFailed && workflow.failureIsExplainedByChild(location, observation) {
				continue
			}
			reference := RunReference{
				RunID:        workflow.RunID,
				RunAttempt:   workflow.RunAttempt,
				Branch:       workflow.Branch,
				CommitSHA:    workflow.CommitSHA,
				RunStartedAt: workflow.RunStartedAt,
				WorkflowURL:  workflow.WorkflowURL,
				JobURL:       observation.jobURL,
			}
			if workflow.Branch == mainBranchName {
				mainAccumulator.addObservation(observation.outcome, reference, observation.message, observation.inferredFromWorkflowSuccess)
				continue
			}
			if pullRequest == nil || workflow.PullRequest == nil || workflow.PullRequest.Number != pullRequest.Number {
				continue
			}
			commit, exists := commitAccumulators[workflow.CommitSHA]
			if !exists {
				commit = &commitAccumulator{
					branch:    workflow.Branch,
					commitSHA: workflow.CommitSHA,
				}
				commitAccumulators[workflow.CommitSHA] = commit
				commitOrder = append(commitOrder, workflow.CommitSHA)
			}
			commit.outcomes.addObservation(observation.outcome, reference, observation.message, observation.inferredFromWorkflowSuccess)
		}
	}

	history := OutcomeHistory{
		MainBranch: mainAccumulator.summary(),
		PRCommits:  make([]CommitOutcomeSummary, 0, len(commitOrder)),
	}
	for _, commitSHA := range commitOrder {
		history.PRCommits = append(history.PRCommits, commitAccumulators[commitSHA].summary())
	}
	return history
}

func (r *Report) findAssessments(selector *FailureSelector, pullRequest *PullRequestRef) []FailureAssessment {
	assessments := make([]FailureAssessment, 0)
	assessmentKeys := make(map[assessmentKey]struct{})
	// workflowIndex is the current position in r.Workflows.
	for workflowIndex := range r.Workflows {
		workflow := &r.Workflows[workflowIndex]
		if pullRequest != nil && !workflow.PullRequest.Equals(pullRequest) {
			continue
		}
		for _, location := range workflow.failureLocations() {
			if !selector.Equals(location) {
				continue
			}
			location = (&location).testIdentity()
			key := assessmentKey{location: location}
			if workflow.PullRequest != nil {
				key.pullRequestNumber = workflow.PullRequest.Number
				key.pullRequestHeadSHA = workflow.PullRequest.HeadSHA
			}
			if _, exists := assessmentKeys[key]; exists {
				continue
			}
			assessmentKeys[key] = struct{}{}
			assessments = append(assessments, r.AssessFailure(location, workflow.PullRequest))
		}
	}
	return assessments
}

func (r *Report) relatedCIHistories(location CIFailureLocation, pullRequest *PullRequestRef) []CIContextHistory {
	if location.TestID == "" {
		return nil
	}
	locations := []CIFailureLocation{{WorkflowName: location.WorkflowName}}
	seenLocations := map[CIFailureLocation]struct{}{locations[0]: {}}
	for workflowIndex := range r.Workflows {
		workflow := &r.Workflows[workflowIndex]
		if workflow.Name != location.WorkflowName {
			continue
		}
		for jobIndex := range workflow.Jobs {
			job := &workflow.Jobs[jobIndex]
			hasTestObservation := job.hasTest(location.TestID, location.TestKind)
			if hasTestObservation {
				jobLocation := CIFailureLocation{WorkflowName: location.WorkflowName, JobName: job.Name}
				if _, exists := seenLocations[jobLocation]; !exists {
					seenLocations[jobLocation] = struct{}{}
					locations = append(locations, jobLocation)
				}
			}
			for stepIndex := range job.Steps {
				step := &job.Steps[stepIndex]
				if !step.hasTest(location.TestID, location.TestKind) {
					continue
				}
				jobLocation := CIFailureLocation{WorkflowName: location.WorkflowName, JobName: job.Name}
				if _, exists := seenLocations[jobLocation]; !exists {
					seenLocations[jobLocation] = struct{}{}
					locations = append(locations, jobLocation)
				}
				stepLocation := CIFailureLocation{WorkflowName: location.WorkflowName, JobName: job.Name, StepName: step.Name}
				if _, exists := seenLocations[stepLocation]; !exists {
					seenLocations[stepLocation] = struct{}{}
					locations = append(locations, stepLocation)
				}
			}
		}
	}
	ciHistories := make([]CIContextHistory, 0, len(locations))
	for _, relatedLocation := range locations {
		ciHistories = append(ciHistories, CIContextHistory{
			Location: relatedLocation,
			History:  r.GetOutcomeHistory(relatedLocation, pullRequest),
		})
	}
	return ciHistories
}

type assessmentKey struct {
	location           CIFailureLocation
	pullRequestNumber  int
	pullRequestHeadSHA string
}

type outcomeAccumulator struct {
	counts        OutcomeCounts
	latestPass    *RunReference
	latestFailure *FailureEvidence
}

func (accumulator *outcomeAccumulator) addObservation(outcome Outcome, reference RunReference, message string, inferredFromWorkflowSuccess bool) {
	switch outcome {
	case OutcomePassed:
		accumulator.counts.Passes++
		if inferredFromWorkflowSuccess {
			accumulator.counts.InferredPasses++
		}
		if accumulator.latestPass == nil || reference.isLaterThan(accumulator.latestPass) {
			accumulator.latestPass = &reference
		}
	case OutcomeFailed:
		accumulator.counts.Failures++
		if accumulator.latestFailure == nil || reference.isLaterThan(&accumulator.latestFailure.SourceRun) {
			accumulator.latestFailure = &FailureEvidence{SourceRun: reference, Message: message}
		}
	}
}

func (accumulator *outcomeAccumulator) summary() OutcomeSummary {
	return OutcomeSummary{
		Counts:        accumulator.counts,
		LatestPass:    accumulator.latestPass,
		LatestFailure: accumulator.latestFailure,
	}
}

type commitAccumulator struct {
	branch    string
	commitSHA string
	outcomes  outcomeAccumulator
}

func (accumulator *commitAccumulator) summary() CommitOutcomeSummary {
	return CommitOutcomeSummary{
		Branch:        accumulator.branch,
		CommitSHA:     accumulator.commitSHA,
		Counts:        accumulator.outcomes.counts,
		LatestPass:    accumulator.outcomes.latestPass,
		LatestFailure: accumulator.outcomes.latestFailure,
	}
}
