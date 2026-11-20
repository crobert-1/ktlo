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

package failingtest

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"ktlo/internal/common"
	"ktlo/internal/github"
	"ktlo/internal/telemetry"
	"ktlo/internal/testresults"
)

const githubArchiveJobNameColon = ":"
const githubArchiveJobNameColonReplacement = ""
const githubArchiveJobNamePathSeparator = "/"
const githubArchiveJobNameSeparator = "_"
const githubLogFileExtension = ".txt"
const githubLogStepNameSeparator = "_"
const githubSystemLogFileName = "system.txt"
const maxGitHubWorkflowAttemptLogArchiveBytes = 250 << 20
const maxGitHubWorkflowAttemptLogFileBytes = 10 << 20
const maxGitHubWorkflowAttemptLogTextBytes = 250 << 20
const workflowAttemptLogArchiveTempPattern = "ktlo-workflow-attempt-logs-*.zip"

const errorCloseTemporaryWorkflowLogArchive = "could not close temporary workflow log archive: %w"
const errorCloseWorkflowAttemptLogArchive = "could not close GitHub workflow attempt log archive for run %s attempt %d: %w"
const errorCloseWorkflowAttemptLogArchiveReader = "could not close GitHub workflow attempt log archive reader for run %s attempt %d: %w"
const errorInvalidGitHubWorkflowRunID = "invalid GitHub workflow run ID %q"
const errorRemoveTemporaryWorkflowLogArchive = "could not remove temporary workflow log archive %q: %w"

const (
	spanReadAndUnpackWorkflowAttemptLogs = "failingtest.read_and_unpack_workflow_attempt_logs"
	traceAttributeArchiveByteCount       = "ktlo.workflow_attempt_log.archive_byte_count"
	traceAttributeLogFileCount           = "ktlo.workflow_attempt_log.file.count"
	traceAttributeUncompressedLogBytes   = "ktlo.workflow_attempt_log.uncompressed_byte_count"
)

const spanDownloadWorkflowAttemptLogArchive = "failingtest.download_workflow_attempt_log_archive"

const spanUnpackWorkflowAttemptLogArchive = "failingtest.unpack_workflow_attempt_log_archive"

type GitHubWorkflowProvider struct {
	client *github.Client
}

type workflowAttemptLogArchiveMatcher struct {
	jobsByArchiveName map[string]WorkflowJob
}

type workflowAttemptLogArchiveEntry struct {
	jobID       string
	stepNumber  int
	isJobOutput bool
}

func newWorkflowAttemptLogArchiveMatcher(jobs []WorkflowJob) (*workflowAttemptLogArchiveMatcher, error) {
	matcher := &workflowAttemptLogArchiveMatcher{jobsByArchiveName: make(map[string]WorkflowJob, len(jobs))}
	for _, job := range jobs {
		archiveJobName := strings.ReplaceAll(job.Name, githubArchiveJobNamePathSeparator, githubArchiveJobNameSeparator)
		archiveJobName = strings.ReplaceAll(archiveJobName, githubArchiveJobNameColon, githubArchiveJobNameColonReplacement)
		if existingJob, exists := matcher.jobsByArchiveName[archiveJobName]; exists {
			return nil, fmt.Errorf("workflow jobs %q and %q have the same GitHub log archive name %q", existingJob.Name, job.Name, archiveJobName)
		}
		matcher.jobsByArchiveName[archiveJobName] = job
	}
	return matcher, nil
}

func (matcher *workflowAttemptLogArchiveMatcher) entryForArchiveFile(fileName string) (workflowAttemptLogArchiveEntry, bool) {
	baseName := path.Base(fileName)
	if path.Dir(fileName) == "." {
		// The root-level number identifies a job; nested log numbers identify steps.
		jobIndex, archiveJobName, found := strings.Cut(baseName, githubLogStepNameSeparator)
		if !found {
			return workflowAttemptLogArchiveEntry{}, false
		}
		jobIndexNumber, err := strconv.Atoi(jobIndex)
		if err != nil || jobIndexNumber < 0 {
			return workflowAttemptLogArchiveEntry{}, false
		}
		archiveJobName = strings.TrimSuffix(archiveJobName, githubLogFileExtension)
		if job, found := matcher.jobForArchiveName(archiveJobName); found {
			return workflowAttemptLogArchiveEntry{jobID: job.ID, isJobOutput: true}, true
		}
		return workflowAttemptLogArchiveEntry{}, false
	}
	if baseName == githubSystemLogFileName {
		return workflowAttemptLogArchiveEntry{}, false
	}
	archiveJobName := path.Base(path.Dir(fileName))
	job, found := matcher.jobForArchiveName(archiveJobName)
	if !found {
		return workflowAttemptLogArchiveEntry{}, false
	}
	stepNumber, found := logNumberFromArchiveFile(baseName)
	if !found {
		return workflowAttemptLogArchiveEntry{}, false
	}
	return workflowAttemptLogArchiveEntry{jobID: job.ID, stepNumber: stepNumber}, true
}

func (matcher *workflowAttemptLogArchiveMatcher) jobForArchiveName(archiveJobName string) (WorkflowJob, bool) {
	job, found := matcher.jobsByArchiveName[archiveJobName]
	return job, found
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

func (provider *GitHubWorkflowProvider) Repository() common.Repository {
	return provider.client.TargetRepository()
}

func (provider *GitHubWorkflowProvider) ListWorkflowRunAttempts(ctx context.Context, startedAfter time.Time, knownAttempts map[common.WorkflowRunAttemptKey]struct{}) ([]WorkflowRunAttempt, error) {
	return provider.client.ListAllWorkflowRunAttempts(ctx, startedAfter, knownAttempts)
}

func (provider *GitHubWorkflowProvider) LatestWorkflowRunAttempt(ctx context.Context, runID string) (WorkflowRunAttempt, error) {
	parsedRunID, err := strconv.ParseInt(runID, 10, 64)
	if err != nil || parsedRunID < 1 {
		return WorkflowRunAttempt{}, fmt.Errorf(errorInvalidGitHubWorkflowRunID, runID)
	}
	return provider.client.LatestWorkflowRunAttempt(ctx, parsedRunID)
}

func (provider *GitHubWorkflowProvider) ListWorkflowJobs(ctx context.Context, runID string, attempt int) ([]WorkflowJob, error) {
	parsedRunID, err := strconv.ParseInt(runID, 10, 64)
	if err != nil || parsedRunID < 1 {
		return nil, fmt.Errorf("invalid GitHub workflow run ID %q", runID)
	}
	return provider.client.ListWorkflowJobs(ctx, parsedRunID, attempt)
}

func (provider *GitHubWorkflowProvider) StepLogsForAttempt(ctx context.Context, runID string, attempt int, jobs []WorkflowJob) (workflowAttemptLogs WorkflowAttemptLogs, resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanReadAndUnpackWorkflowAttemptLogs,
		attribute.String(traceAttributeRunID, runID),
		attribute.Int(traceAttributeRunAttempt, attempt),
	)
	archiveByteCount := int64(0)
	logFileCount := 0
	stepLogByteCount := 0
	defer func() {
		span.SetAttributes(
			attribute.Int64(traceAttributeArchiveByteCount, archiveByteCount),
			attribute.Int(traceAttributeLogFileCount, logFileCount),
			attribute.Int(traceAttributeUncompressedLogBytes, stepLogByteCount),
		)
		telemetry.FinishSpan(ctx, span, spanReadAndUnpackWorkflowAttemptLogs, startedAt, resultErr != nil)
	}()

	parsedRunID, err := strconv.ParseInt(runID, 10, 64)
	if err != nil || parsedRunID < 1 || attempt < 1 {
		return WorkflowAttemptLogs{}, fmt.Errorf("invalid GitHub workflow run ID %q or attempt %d", runID, attempt)
	}
	archiveMatcher, err := newWorkflowAttemptLogArchiveMatcher(jobs)
	if err != nil {
		return WorkflowAttemptLogs{}, fmt.Errorf("could not match GitHub workflow log archive for run %s attempt %d: %w", runID, attempt, err)
	}
	jobLogArchive, err := provider.client.WorkflowAttemptLogArchive(ctx, parsedRunID, attempt)
	if err != nil {
		return WorkflowAttemptLogs{}, err
	}
	defer func() {
		if closeErr := jobLogArchive.Close(); closeErr != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf(errorCloseWorkflowAttemptLogArchive, runID, attempt, closeErr))
		}
	}()
	temporaryArchive, err := os.CreateTemp("", workflowAttemptLogArchiveTempPattern)
	if err != nil {
		return WorkflowAttemptLogs{}, fmt.Errorf("could not create temporary workflow log archive: %w", err)
	}
	defer func() {
		if removeErr := os.Remove(temporaryArchive.Name()); removeErr != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf(errorRemoveTemporaryWorkflowLogArchive, temporaryArchive.Name(), removeErr))
		}
	}()
	downloadStartedAt := time.Now()
	downloadContext, downloadSpan := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanDownloadWorkflowAttemptLogArchive,
		attribute.String(traceAttributeRunID, runID),
		attribute.Int(traceAttributeRunAttempt, attempt),
	)
	archiveByteCount, err = io.Copy(temporaryArchive, io.LimitReader(jobLogArchive, maxGitHubWorkflowAttemptLogArchiveBytes+1))
	downloadSpan.SetAttributes(attribute.Int64(traceAttributeArchiveByteCount, archiveByteCount))
	telemetry.FinishSpan(downloadContext, downloadSpan, spanDownloadWorkflowAttemptLogArchive, downloadStartedAt, err != nil || archiveByteCount > maxGitHubWorkflowAttemptLogArchiveBytes)
	if err != nil {
		_ = temporaryArchive.Close()
		return WorkflowAttemptLogs{}, fmt.Errorf("could not read GitHub workflow attempt log archive for run %s attempt %d: %w", runID, attempt, err)
	}
	if archiveByteCount > maxGitHubWorkflowAttemptLogArchiveBytes {
		_ = temporaryArchive.Close()
		return WorkflowAttemptLogs{}, fmt.Errorf("GitHub workflow attempt log archive for run %s attempt %d exceeds %d bytes", runID, attempt, maxGitHubWorkflowAttemptLogArchiveBytes)
	}
	if err := temporaryArchive.Close(); err != nil {
		return WorkflowAttemptLogs{}, fmt.Errorf(errorCloseTemporaryWorkflowLogArchive, err)
	}
	unpackStartedAt := time.Now()
	unpackContext, unpackSpan := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanUnpackWorkflowAttemptLogArchive,
		attribute.String(traceAttributeRunID, runID),
		attribute.Int(traceAttributeRunAttempt, attempt),
	)
	defer func() {
		unpackSpan.SetAttributes(
			attribute.Int(traceAttributeLogFileCount, logFileCount),
			attribute.Int(traceAttributeUncompressedLogBytes, stepLogByteCount),
		)
		telemetry.FinishSpan(unpackContext, unpackSpan, spanUnpackWorkflowAttemptLogArchive, unpackStartedAt, resultErr != nil)
	}()
	archiveReader, err := zip.OpenReader(temporaryArchive.Name())
	if err != nil {
		return WorkflowAttemptLogs{}, fmt.Errorf("invalid GitHub workflow attempt log archive for run %s attempt %d: %w", runID, attempt, err)
	}
	defer func() {
		if closeErr := archiveReader.Close(); closeErr != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf(errorCloseWorkflowAttemptLogArchiveReader, runID, attempt, closeErr))
		}
	}()
	workflowAttemptLogs = WorkflowAttemptLogs{Jobs: make([]WorkflowJobLogs, 0)}
	jobLogIndexesByID := make(map[string]int, len(jobs))
	for _, archiveFile := range archiveReader.File {
		if err := ctx.Err(); err != nil {
			return WorkflowAttemptLogs{}, err
		}
		archiveEntry, ok := archiveMatcher.entryForArchiveFile(archiveFile.Name)
		if !ok {
			continue
		}
		if archiveFile.UncompressedSize64 > maxGitHubWorkflowAttemptLogFileBytes {
			return WorkflowAttemptLogs{}, fmt.Errorf("GitHub workflow log file %q exceeds %d bytes", archiveFile.Name, maxGitHubWorkflowAttemptLogFileBytes)
		}
		logFile, err := archiveFile.Open()
		if err != nil {
			return WorkflowAttemptLogs{}, fmt.Errorf("could not open GitHub workflow log file %q: %w", archiveFile.Name, err)
		}
		logData, readErr := io.ReadAll(io.LimitReader(logFile, maxGitHubWorkflowAttemptLogFileBytes+1))
		closeErr := logFile.Close()
		if readErr != nil {
			return WorkflowAttemptLogs{}, fmt.Errorf("could not read GitHub workflow log file %q: %w", archiveFile.Name, readErr)
		}
		if closeErr != nil {
			return WorkflowAttemptLogs{}, fmt.Errorf("could not close GitHub workflow log file %q: %w", archiveFile.Name, closeErr)
		}
		if len(logData) > maxGitHubWorkflowAttemptLogFileBytes {
			return WorkflowAttemptLogs{}, fmt.Errorf("GitHub workflow log file %q exceeds %d bytes", archiveFile.Name, maxGitHubWorkflowAttemptLogFileBytes)
		}
		if stepLogByteCount+len(logData) > maxGitHubWorkflowAttemptLogTextBytes {
			return WorkflowAttemptLogs{}, fmt.Errorf("GitHub workflow attempt logs for run %s attempt %d exceed %d uncompressed bytes", runID, attempt, maxGitHubWorkflowAttemptLogTextBytes)
		}
		jobLogIndex, exists := jobLogIndexesByID[archiveEntry.jobID]
		if !exists {
			workflowAttemptLogs.Jobs = append(workflowAttemptLogs.Jobs, WorkflowJobLogs{JobID: archiveEntry.jobID, Steps: make([]WorkflowStepLog, 0)})
			jobLogIndex = len(workflowAttemptLogs.Jobs) - 1
			jobLogIndexesByID[archiveEntry.jobID] = jobLogIndex
		}
		if archiveEntry.isJobOutput {
			workflowAttemptLogs.Jobs[jobLogIndex].JobOutput = string(logData)
		} else {
			workflowAttemptLogs.Jobs[jobLogIndex].Steps = append(workflowAttemptLogs.Jobs[jobLogIndex].Steps, WorkflowStepLog{
				Number: archiveEntry.stepNumber,
				Text:   string(logData),
			})
		}
		logFileCount++
		stepLogByteCount += len(logData)
	}
	collectionErrors := make([]error, 0)
	for _, job := range jobs {
		jobLogIndex, found := jobLogIndexesByID[job.ID]
		if !found {
			collectionErrors = append(collectionErrors, fmt.Errorf("GitHub workflow attempt log archive for run %s attempt %d has no logs for job %q", runID, attempt, job.Name))
			continue
		}
		workflowAttemptLogs.Jobs[jobLogIndex].splitCombinedOutput(job)
	}
	return workflowAttemptLogs, errors.Join(collectionErrors...)
}

func logNumberFromArchiveFile(fileName string) (int, bool) {
	baseName := path.Base(fileName)
	stepNumberText, _, found := strings.Cut(baseName, githubLogStepNameSeparator)
	if !found {
		return 0, false
	}
	stepNumber, err := strconv.Atoi(stepNumberText)
	return stepNumber, err == nil && stepNumber >= 0
}
