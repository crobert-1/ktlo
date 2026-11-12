package testresults

import (
	"errors"
	"math"
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
	UnknownReasonMissingReport         UnknownReason = "missing_report"
	UnknownReasonStaleReport           UnknownReason = "stale_report"
	UnknownReasonUnidentifiableFailure UnknownReason = "unidentifiable_failure"
)

func (reason UnknownReason) IsValid() bool {
	switch reason {
	case UnknownReasonHeadSHAMismatch, UnknownReasonIncompleteLogCoverage, UnknownReasonInsufficientHistory, UnknownReasonMissingReport, UnknownReasonStaleReport, UnknownReasonUnidentifiableFailure:
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

	if isAnalysisIncomplete {
		return ClassificationUnknown, UnknownReasonIncompleteLogCoverage
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

type CIContextHistory struct {
	Location CIFailureLocation
	History  OutcomeHistory
}

const mainBranchName = "main"

type OutcomeHistory struct {
	MainBranch OutcomeSummary
	PRCommits  []CommitOutcomeSummary
}

type historyObservation struct {
	outcome Outcome
	jobURL  string
}

type OutcomeSummary struct {
	Counts        OutcomeCounts
	LatestPass    *RunReference
	LatestFailure *RunReference
}

type CommitOutcomeSummary struct {
	Branch        string
	CommitSHA     string
	Counts        OutcomeCounts
	LatestPass    *RunReference
	LatestFailure *RunReference
}

type OutcomeCounts struct {
	Passes   int
	Failures int
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

func (r *Report) FindAssessments(selector FailureSelector) ([]FailureAssessment, error) {
	if err := selector.Validate(); err != nil {
		return nil, err
	}
	return r.findAssessments(&selector, nil), nil
}

func (r *Report) AssessmentsForPullRequest(pullRequest PullRequestRef) []FailureAssessment {
	return r.findAssessments(nil, &pullRequest)
}

func (r *Report) GetOutcomeHistory(location CIFailureLocation, pullRequest *PullRequestRef) OutcomeHistory {
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
				mainAccumulator.addObservation(observation.outcome, reference)
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
			commit.outcomes.addObservation(observation.outcome, reference)
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
	ciHistories := make([]CIContextHistory, 0, 3)
	locations := []CIFailureLocation{{WorkflowName: location.WorkflowName}}
	if location.JobName != "" {
		locations = append(locations, CIFailureLocation{WorkflowName: location.WorkflowName, JobName: location.JobName})
	}
	if location.StepName != "" {
		locations = append(locations, CIFailureLocation{WorkflowName: location.WorkflowName, JobName: location.JobName, StepName: location.StepName})
	}
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
	latestFailure *RunReference
}

func (accumulator *outcomeAccumulator) addObservation(outcome Outcome, reference RunReference) {
	switch outcome {
	case OutcomePassed:
		accumulator.counts.Passes++
		if accumulator.latestPass == nil || reference.isLaterThan(accumulator.latestPass) {
			accumulator.latestPass = &reference
		}
	case OutcomeFailed:
		accumulator.counts.Failures++
		if accumulator.latestFailure == nil || reference.isLaterThan(accumulator.latestFailure) {
			accumulator.latestFailure = &reference
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
