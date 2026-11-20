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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ktlo/internal/common"
)

const (
	currentSchemaVersion  = 1
	errorInvalidReport    = "invalid test results report %q: %w"
	errorRequiredField    = "%s is required"
	pathResultsDirectory  = "results"
	reportFileName        = "workflow-test-analysis.json"
	reportFileMode        = 0o644
	reportDirectoryMode   = 0o755
	reportTempFilePattern = ".workflow-test-analysis-*.json"
)

const (
	errorFailureTextOnNonFailure = "failure text is only valid for a failed outcome"
	errorInvalidTest             = "invalid test observation %q: %w"
)

var (
	ErrInvalidLookbackCutoff  = errors.New("invalid workflow run cutoff")
	ErrInvalidReportIdentity  = errors.New("invalid test results report identity")
	ErrInvalidWorkflowRun     = errors.New("invalid workflow run")
	ErrMalformedReport        = errors.New("malformed test results report")
	ErrReportIdentityMismatch = errors.New("test results report identity does not match repository")
)

type Provider string

func (provider *Provider) UnmarshalJSON(data []byte) error {
	value, err := normalizeEnumJSON(data)
	if err != nil {
		return err
	}
	*provider = Provider(value)
	return nil
}

const (
	ProviderGitHub Provider = "github"
	ProviderGitLab Provider = "gitlab"
)

func (provider Provider) IsValid() bool {
	switch provider {
	case ProviderGitHub, ProviderGitLab:
		return true
	default:
		return false
	}
}

type AnalysisStatus string

func (status *AnalysisStatus) UnmarshalJSON(data []byte) error {
	value, err := normalizeEnumJSON(data)
	if err != nil {
		return err
	}
	*status = AnalysisStatus(value)
	return nil
}

const (
	AnalysisStatusComplete   AnalysisStatus = "complete"
	AnalysisStatusIncomplete AnalysisStatus = "incomplete"
)

func (status AnalysisStatus) IsValid() bool {
	switch status {
	case AnalysisStatusComplete, AnalysisStatusIncomplete:
		return true
	default:
		return false
	}
}

type Outcome string

func (outcome *Outcome) UnmarshalJSON(data []byte) error {
	value, err := normalizeEnumJSON(data)
	if err != nil {
		return err
	}
	*outcome = Outcome(value)
	return nil
}

const (
	OutcomeFailed       Outcome = "failed"
	OutcomeInconclusive Outcome = "inconclusive"
	OutcomePassed       Outcome = "passed"
	OutcomeSkipped      Outcome = "skipped"
)

func (outcome Outcome) IsValid() bool {
	switch outcome {
	case OutcomeFailed, OutcomeInconclusive, OutcomePassed, OutcomeSkipped:
		return true
	default:
		return false
	}
}

func (outcome *Outcome) Validate(providerConclusion string) error {
	if !outcome.IsValid() {
		return fmt.Errorf("invalid outcome %q", *outcome)
	}
	if *outcome == OutcomeInconclusive && strings.TrimSpace(providerConclusion) == "" {
		return errors.New("provider conclusion is required for an inconclusive outcome")
	}
	return nil
}

type TestKind string

func (kind *TestKind) UnmarshalJSON(data []byte) error {
	value, err := normalizeEnumJSON(data)
	if err != nil {
		return err
	}
	*kind = TestKind(value)
	return nil
}

const (
	TestKindPackage TestKind = "package"
	TestKindTest    TestKind = "test"
)

func (kind TestKind) IsValid() bool {
	switch kind {
	case TestKindPackage, TestKindTest:
		return true
	default:
		return false
	}
}

func normalizeEnumJSON(data []byte) (string, error) {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return "", err
	}
	return strings.ToLower(value), nil
}

type Repository = common.Repository
type PullRequestRef = common.PullRequestReference

type Report struct {
	SchemaVersion  int            `json:"schemaVersion"`
	Provider       Provider       `json:"provider"`
	Repository     Repository     `json:"repository"`
	LookbackStart  time.Time      `json:"lookbackStart,omitempty"`
	GeneratedAt    time.Time      `json:"generatedAt"`
	AnalysisStatus AnalysisStatus `json:"analysisStatus"`
	Workflows      []Workflow     `json:"workflows"`
}

func (r *Report) Validate() error {
	if err := r.ValidateRepo(); err != nil {
		return err
	}
	if r.SchemaVersion != currentSchemaVersion {
		return fmt.Errorf("unsupported test results schema version %d", r.SchemaVersion)
	}
	if r.GeneratedAt.IsZero() {
		return fmt.Errorf(errorRequiredField, "test results report generatedAt timestamp")
	}
	if r.AnalysisStatus == "" {
		return fmt.Errorf(errorRequiredField, "test results report analysis status")
	}
	if !r.AnalysisStatus.IsValid() {
		return fmt.Errorf("invalid analysis status %q", r.AnalysisStatus)
	}
	workflowRunKeys := make(map[workflowRunKey]struct{}, len(r.Workflows))
	for _, workflow := range r.Workflows {
		if err := workflow.Validate(); err != nil {
			return err
		}
		key := workflowRunKey{runID: workflow.RunID, runAttempt: workflow.RunAttempt}
		if _, exists := workflowRunKeys[key]; exists {
			return fmt.Errorf("duplicate workflow run ID %q attempt %d", workflow.RunID, workflow.RunAttempt)
		}
		workflowRunKeys[key] = struct{}{}
	}
	return nil
}

func (r *Report) ValidateRepo() error {
	if !r.Provider.IsValid() {
		return fmt.Errorf("%w: unsupported report provider %q", ErrInvalidReportIdentity, r.Provider)
	}
	if !validPathSegment(r.Repository.Owner) || !validPathSegment(r.Repository.Name) {
		return fmt.Errorf("%w: invalid repository identity: owner %q, name %q", ErrInvalidReportIdentity, r.Repository.Owner, r.Repository.Name)
	}
	return nil
}

type Workflow struct {
	Name               string          `json:"name"`
	RunID              string          `json:"runId"`
	RunAttempt         int             `json:"runAttempt"`
	RunStartedAt       time.Time       `json:"runStartedAt"`
	Branch             string          `json:"branch"`
	CommitSHA          string          `json:"commitSha"`
	Outcome            Outcome         `json:"outcome"`
	ProviderConclusion string          `json:"providerConclusion,omitempty"`
	WorkflowURL        string          `json:"runUrl"`
	PullRequest        *PullRequestRef `json:"pullRequest,omitempty"`
	CollectionComplete bool            `json:"collectionComplete"`
	Jobs               []Job           `json:"jobs,omitempty"`
	Tests              []Test          `json:"tests,omitempty"`
	ciFailureLocations []CIFailureLocation
}

func (workflow *Workflow) hasFailedTestInPackage(packageID string, jobName string, jobURL string) bool {
	if jobName == "" && jobURL == "" {
		for _, test := range workflow.Tests {
			if test.isFailedChildOfPackage(packageID) {
				return true
			}
		}
		return false
	}
	var matchingJob *Job
	for jobIndex := range workflow.Jobs {
		job := &workflow.Jobs[jobIndex]
		if jobURL != "" && job.JobURL != jobURL || jobURL == "" && job.Name != jobName {
			continue
		}
		if matchingJob != nil {
			return false
		}
		matchingJob = job
	}
	return matchingJob != nil && matchingJob.hasFailedTestInPackage(packageID)
}

func (workflow *Workflow) failureIsExplainedByChild(location CIFailureLocation, observation historyObservation) bool {
	if location.TestKind == TestKindPackage {
		return workflow.hasFailedTestInPackage(location.TestID, observation.jobName, observation.jobURL)
	}
	if location.TestID != "" || location.StepName != "" {
		return false
	}
	if location.JobName == "" {
		for _, job := range workflow.Jobs {
			if job.Outcome == OutcomeFailed {
				return true
			}
		}
		return false
	}
	for jobIndex := range workflow.Jobs {
		job := &workflow.Jobs[jobIndex]
		jobMatches := job.Name == location.JobName
		if observation.jobURL != "" {
			jobMatches = job.JobURL == observation.jobURL
		}
		if !jobMatches {
			continue
		}
		for _, step := range job.Steps {
			if step.Outcome == OutcomeFailed {
				return true
			}
		}
	}
	return false
}

func (workflow *Workflow) observationsFor(location CIFailureLocation) []historyObservation {
	if location.WorkflowName != "" && location.WorkflowName != workflow.Name {
		return nil
	}
	observations := make([]historyObservation, 0)
	appendTestObservations := func(tests []Test, jobName string, jobURL string) {
		for _, test := range tests {
			if test.ID != location.TestID || (location.TestKind != "" && test.Kind != location.TestKind) {
				continue
			}
			observations = append(observations, historyObservation{
				outcome: test.Outcome,
				jobName: jobName,
				jobURL:  jobURL,
				message: test.FailureMessage,
			})
		}
	}
	if location.TestID != "" {
		if location.JobName == "" && location.StepName == "" {
			appendTestObservations(workflow.Tests, "", "")
		}
		// jobIndex is the current position in workflow.Jobs.
		for jobIndex := range workflow.Jobs {
			job := &workflow.Jobs[jobIndex]
			if location.JobName != "" && location.JobName != job.Name {
				continue
			}
			if location.StepName == "" {
				appendTestObservations(job.Tests, job.Name, job.JobURL)
			}
			// stepIndex is the current position in job.Steps.
			for stepIndex := range job.Steps {
				step := &job.Steps[stepIndex]
				if location.StepName == "" || step.Name == location.StepName {
					appendTestObservations(step.Tests, job.Name, job.JobURL)
				}
			}
		}
		return workflow.withSuccessfulWorkflowPass(location, observations)
	}
	if location.StepName != "" {
		// jobIndex is the current position in workflow.Jobs.
		for jobIndex := range workflow.Jobs {
			job := &workflow.Jobs[jobIndex]
			if location.JobName != "" && location.JobName != job.Name {
				continue
			}
			// stepIndex is the current position in job.Steps.
			for stepIndex := range job.Steps {
				step := &job.Steps[stepIndex]
				if step.Name == location.StepName {
					observations = append(observations, historyObservation{
						outcome: step.Outcome,
						jobName: job.Name,
						jobURL:  job.JobURL,
						message: step.ReasonExcerpt,
					})
				}
			}
		}
		return workflow.withSuccessfulWorkflowPass(location, observations)
	}
	if location.JobName != "" {
		// jobIndex is the current position in workflow.Jobs.
		for jobIndex := range workflow.Jobs {
			job := &workflow.Jobs[jobIndex]
			if job.Name == location.JobName {
				observations = append(observations, historyObservation{outcome: job.Outcome, jobName: job.Name, jobURL: job.JobURL})
			}
		}
		return workflow.withSuccessfulWorkflowPass(location, observations)
	}
	observations = append(observations, historyObservation{outcome: workflow.Outcome})
	return workflow.withSuccessfulWorkflowPass(location, observations)
}

func (workflow *Workflow) withSuccessfulWorkflowPass(location CIFailureLocation, observations []historyObservation) []historyObservation {
	if workflow.Outcome != OutcomePassed {
		return observations
	}
	for _, observation := range observations {
		if observation.outcome == OutcomePassed {
			return observations
		}
	}
	return append(observations, historyObservation{
		outcome:                     OutcomePassed,
		jobURL:                      workflow.jobURLFor(location),
		inferredFromWorkflowSuccess: true,
	})
}

func (workflow *Workflow) jobURLFor(location CIFailureLocation) string {
	for jobIndex := range workflow.Jobs {
		job := &workflow.Jobs[jobIndex]
		if location.JobName != "" && location.JobName == job.Name {
			return job.JobURL
		}
		if location.JobName == "" && location.StepName != "" {
			for stepIndex := range job.Steps {
				if job.Steps[stepIndex].Name == location.StepName {
					return job.JobURL
				}
			}
		}
	}
	return ""
}

func (workflow *Workflow) Validate() error {
	if strings.TrimSpace(workflow.Name) == "" {
		return fmt.Errorf("%w: %w", ErrInvalidWorkflowRun, fmt.Errorf(errorRequiredField, "workflow name"))
	}
	if strings.TrimSpace(workflow.RunID) == "" {
		return fmt.Errorf("%w: %w", ErrInvalidWorkflowRun, fmt.Errorf(errorRequiredField, "workflow run ID"))
	}
	if workflow.RunAttempt < 1 {
		return fmt.Errorf("%w: workflow run attempt must be positive", ErrInvalidWorkflowRun)
	}
	if workflow.RunStartedAt.IsZero() {
		return fmt.Errorf("%w: %w", ErrInvalidWorkflowRun, fmt.Errorf(errorRequiredField, "workflow run start time"))
	}
	if strings.TrimSpace(workflow.Branch) == "" {
		return fmt.Errorf("%w: %w", ErrInvalidWorkflowRun, fmt.Errorf(errorRequiredField, "workflow branch"))
	}
	if strings.TrimSpace(workflow.CommitSHA) == "" {
		return fmt.Errorf("%w: %w", ErrInvalidWorkflowRun, fmt.Errorf(errorRequiredField, "workflow commit SHA"))
	}
	if err := workflow.Outcome.Validate(workflow.ProviderConclusion); err != nil {
		return fmt.Errorf("%w: workflow %q: %w", ErrInvalidWorkflowRun, workflow.Name, err)
	}
	if workflow.PullRequest != nil {
		if workflow.PullRequest.Number < 1 {
			return fmt.Errorf("%w: workflow %q: pull request number must be positive", ErrInvalidWorkflowRun, workflow.Name)
		}
		if strings.TrimSpace(workflow.PullRequest.HeadSHA) == "" {
			return fmt.Errorf("%w: workflow %q: pull request head SHA is required", ErrInvalidWorkflowRun, workflow.Name)
		}
	}
	for _, job := range workflow.Jobs {
		if err := job.Validate(); err != nil {
			return fmt.Errorf("%w: workflow %q: %w", ErrInvalidWorkflowRun, workflow.Name, err)
		}
	}
	for _, test := range workflow.Tests {
		if err := test.Validate(); err != nil {
			return fmt.Errorf("%w: workflow %q: %w", ErrInvalidWorkflowRun, workflow.Name, err)
		}
	}
	return nil
}

func (workflow *Workflow) failureLocations() []CIFailureLocation {
	if workflow.ciFailureLocations != nil {
		return workflow.ciFailureLocations
	}
	locations := make([]CIFailureLocation, 0)
	testLocations := make(map[CIFailureLocation]struct{})
	appendFailedTest := func(test Test) {
		if test.Outcome != OutcomeFailed {
			return
		}
		location := (&CIFailureLocation{
			WorkflowName: workflow.Name,
			TestID:       test.ID,
			TestKind:     test.Kind,
		}).testIdentity()
		if _, exists := testLocations[location]; exists {
			return
		}
		testLocations[location] = struct{}{}
		locations = append(locations, location)
	}
	if workflow.Outcome == OutcomeFailed {
		locations = append(locations, CIFailureLocation{WorkflowName: workflow.Name})
	}
	for _, test := range workflow.Tests {
		appendFailedTest(test)
	}
	// jobIndex is the current position in workflow.Jobs.
	for jobIndex := range workflow.Jobs {
		job := &workflow.Jobs[jobIndex]
		if job.Outcome == OutcomeFailed {
			locations = append(locations, CIFailureLocation{WorkflowName: workflow.Name, JobName: job.Name})
		}
		for _, test := range job.Tests {
			appendFailedTest(test)
		}
		// stepIndex is the current position in job.Steps.
		for stepIndex := range job.Steps {
			step := &job.Steps[stepIndex]
			if step.Outcome == OutcomeFailed {
				locations = append(locations, CIFailureLocation{WorkflowName: workflow.Name, JobName: job.Name, StepName: step.Name})
			}
			for _, test := range step.Tests {
				appendFailedTest(test)
			}
		}
	}
	workflow.ciFailureLocations = locations
	return workflow.ciFailureLocations
}

type Job struct {
	Name               string  `json:"name"`
	Outcome            Outcome `json:"outcome"`
	ProviderConclusion string  `json:"providerConclusion,omitempty"`
	JobURL             string  `json:"jobUrl,omitempty"`
	Steps              []Step  `json:"steps,omitempty"`
	Tests              []Test  `json:"tests,omitempty"`
}

func (job *Job) hasFailedTestInPackage(packageID string) bool {
	for _, test := range job.Tests {
		if test.isFailedChildOfPackage(packageID) {
			return true
		}
	}
	for stepIndex := range job.Steps {
		step := &job.Steps[stepIndex]
		for _, test := range step.Tests {
			if test.isFailedChildOfPackage(packageID) {
				return true
			}
		}
	}
	return false
}

func (job *Job) hasTest(testID string, testKind TestKind) bool {
	for _, test := range job.Tests {
		if test.ID == testID && test.Kind == testKind {
			return true
		}
	}
	return false
}

func (job *Job) Validate() error {
	if strings.TrimSpace(job.Name) == "" {
		return fmt.Errorf(errorRequiredField, "job name")
	}
	if err := job.Outcome.Validate(job.ProviderConclusion); err != nil {
		return fmt.Errorf("job %q: %w", job.Name, err)
	}
	for _, step := range job.Steps {
		if err := step.Validate(); err != nil {
			return fmt.Errorf("job %q: %w", job.Name, err)
		}
	}
	for _, test := range job.Tests {
		if err := test.Validate(); err != nil {
			return fmt.Errorf("job %q: %w", job.Name, err)
		}
	}
	return nil
}

type Step struct {
	Name               string  `json:"name"`
	Outcome            Outcome `json:"outcome"`
	ProviderConclusion string  `json:"providerConclusion,omitempty"`
	ReasonExcerpt      string  `json:"reasonExcerpt,omitempty"`
	Tests              []Test  `json:"tests,omitempty"`
}

func (step *Step) hasTest(testID string, testKind TestKind) bool {
	for _, test := range step.Tests {
		if test.ID == testID && test.Kind == testKind {
			return true
		}
	}
	return false
}

func (step *Step) Validate() error {
	if strings.TrimSpace(step.Name) == "" {
		return fmt.Errorf(errorRequiredField, "step name")
	}
	if err := step.Outcome.Validate(step.ProviderConclusion); err != nil {
		return fmt.Errorf("step %q: %w", step.Name, err)
	}
	if step.Outcome != OutcomeFailed && step.ReasonExcerpt != "" {
		return fmt.Errorf("step %q: %s", step.Name, errorFailureTextOnNonFailure)
	}
	for _, test := range step.Tests {
		if err := test.Validate(); err != nil {
			return fmt.Errorf("step %q: %w", step.Name, err)
		}
	}
	return nil
}

type Test struct {
	ID             string   `json:"id"`
	Kind           TestKind `json:"kind"`
	Outcome        Outcome  `json:"outcome"`
	FailureMessage string   `json:"failureMessage,omitempty"`
}

const testPackageIDSeparator = "."

func (test *Test) isFailedChildOfPackage(packageID string) bool {
	return test.Kind == TestKindTest && test.Outcome == OutcomeFailed && strings.HasPrefix(test.ID, packageID+testPackageIDSeparator)
}

func (test *Test) Validate() error {
	if strings.TrimSpace(test.ID) == "" {
		return fmt.Errorf(errorInvalidTest, test.ID, fmt.Errorf(errorRequiredField, "test ID"))
	}
	if !test.Kind.IsValid() {
		return fmt.Errorf(errorInvalidTest, test.ID, fmt.Errorf("invalid test kind %q", test.Kind))
	}
	if err := test.Outcome.Validate(""); err != nil {
		return fmt.Errorf(errorInvalidTest, test.ID, err)
	}
	if test.Outcome != OutcomeFailed && test.FailureMessage != "" {
		return fmt.Errorf(errorInvalidTest, test.ID, errors.New(errorFailureTextOnNonFailure))
	}
	return nil
}

func NewReport(provider Provider, repository Repository) *Report {
	return &Report{
		SchemaVersion: currentSchemaVersion,
		Provider:      provider,
		Repository:    repository,
		Workflows:     []Workflow{},
	}
}

func (r *Report) Load() error {
	if err := r.ValidateRepo(); err != nil {
		return err
	}

	reportPath, err := r.pathFor()
	if err != nil {
		return err
	}
	reportData, err := os.ReadFile(reportPath)
	if err != nil {
		return fmt.Errorf("could not load test results report %q: %w", reportPath, err)
	}

	var loadedReport Report
	if err := json.Unmarshal(reportData, &loadedReport); err != nil {
		return fmt.Errorf(errorInvalidReport, reportPath, fmt.Errorf("%w: %w", ErrMalformedReport, err))
	}
	if loadedReport.Provider != r.Provider {
		return fmt.Errorf(errorInvalidReport, reportPath, fmt.Errorf("%w: report provider does not match requested provider", ErrReportIdentityMismatch))
	}
	if loadedReport.Repository != r.Repository {
		return fmt.Errorf(errorInvalidReport, reportPath, fmt.Errorf("%w: report repository does not match requested repository", ErrReportIdentityMismatch))
	}
	if err := loadedReport.Validate(); err != nil {
		return fmt.Errorf(errorInvalidReport, reportPath, fmt.Errorf("%w: %w", ErrMalformedReport, err))
	}

	*r = loadedReport
	return nil
}

func (r *Report) Write() error {
	if err := r.Validate(); err != nil {
		return err
	}

	reportPath, err := r.pathFor()
	if err != nil {
		return err
	}
	reportDirectory := filepath.Dir(reportPath)
	if err := os.MkdirAll(reportDirectory, reportDirectoryMode); err != nil {
		return fmt.Errorf("could not create test results directory %q: %w", reportDirectory, err)
	}

	reportData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode test results report %q: %w", reportPath, err)
	}
	reportData = append(reportData, '\n')

	temporaryFile, err := os.CreateTemp(reportDirectory, reportTempFilePattern)
	if err != nil {
		return fmt.Errorf("could not create temporary test results report in %q: %w", reportDirectory, err)
	}
	temporaryPath := temporaryFile.Name()
	defer func() {
		_ = os.Remove(temporaryPath)
	}()

	if err := temporaryFile.Chmod(reportFileMode); err != nil {
		_ = temporaryFile.Close()
		return fmt.Errorf("could not set test results report permissions: %w", err)
	}
	if _, err := temporaryFile.Write(reportData); err != nil {
		_ = temporaryFile.Close()
		return fmt.Errorf("could not write test results report %q: %w", reportPath, err)
	}
	if err := temporaryFile.Close(); err != nil {
		return fmt.Errorf("could not close test results report %q: %w", temporaryPath, err)
	}
	if err := os.Rename(temporaryPath, reportPath); err != nil {
		return fmt.Errorf("could not move test results report into place: %w", err)
	}

	return nil
}

func (r *Report) MergeWorkflowRuns(runs []Workflow, cutoff time.Time) error {
	if err := r.ValidateRepo(); err != nil {
		return err
	}
	if cutoff.IsZero() {
		return ErrInvalidLookbackCutoff
	}
	for _, run := range runs {
		if err := run.Validate(); err != nil {
			return fmt.Errorf("cannot merge workflow run %q attempt %d: %w", run.RunID, run.RunAttempt, err)
		}
	}

	collectedRunsByKey := make(map[workflowRunKey]Workflow, len(runs))
	for _, run := range runs {
		if run.RunStartedAt.Before(cutoff) {
			continue
		}
		key := workflowRunKey{runID: run.RunID, runAttempt: run.RunAttempt}
		collectedRunsByKey[key] = run
	}
	mergedWorkflows := make([]Workflow, 0, len(r.Workflows)+len(runs))
	workflowRunKeys := make(map[workflowRunKey]struct{}, len(r.Workflows)+len(runs))
	for _, existingRun := range r.Workflows {
		if existingRun.RunStartedAt.Before(cutoff) {
			continue
		}
		key := workflowRunKey{runID: existingRun.RunID, runAttempt: existingRun.RunAttempt}
		if _, exists := workflowRunKeys[key]; exists {
			continue
		}
		workflowRunKeys[key] = struct{}{}
		if collectedRun, recollected := collectedRunsByKey[key]; recollected {
			mergedWorkflows = append(mergedWorkflows, collectedRun)
			continue
		}
		mergedWorkflows = append(mergedWorkflows, existingRun)
	}
	for _, run := range runs {
		if run.RunStartedAt.Before(cutoff) {
			continue
		}
		key := workflowRunKey{runID: run.RunID, runAttempt: run.RunAttempt}
		if _, exists := workflowRunKeys[key]; exists {
			continue
		}
		workflowRunKeys[key] = struct{}{}
		mergedWorkflows = append(mergedWorkflows, run)
	}

	r.Workflows = mergedWorkflows
	return nil
}

func (r *Report) HasPullRequestHeadMismatch(pullRequest *PullRequestRef) bool {
	if pullRequest == nil {
		return false
	}
	var hasDifferentHeadSHA bool
	for workflowIndex := range r.Workflows {
		workflow := &r.Workflows[workflowIndex]
		if workflow.PullRequest == nil || workflow.PullRequest.Number != pullRequest.Number {
			continue
		}
		if workflow.PullRequest.HeadSHA == pullRequest.HeadSHA {
			return false
		}
		hasDifferentHeadSHA = true
	}
	return hasDifferentHeadSHA
}

func (r *Report) pathFor() (string, error) {
	if err := r.ValidateRepo(); err != nil {
		return "", err
	}
	return filepath.Join(pathResultsDirectory, string(r.Provider), r.Repository.Owner, r.Repository.Name, reportFileName), nil
}

func validPathSegment(value string) bool {
	return value != "" && value != "." && value != ".." && !strings.ContainsAny(value, `/\\`)
}

type workflowRunKey struct {
	runID      string
	runAttempt int
}
