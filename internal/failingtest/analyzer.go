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
	"context"
	"errors"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"ktlo/internal/common"
	"ktlo/internal/telemetry"
	"ktlo/internal/testresults"
)

const (
	defaultLookbackDuration               = 7 * 24 * time.Hour
	errorInvalidLookbackDuration          = "lookback duration must be greater than zero"
	spanAnalyzerRun                       = "failingtest.analyzer.run"
	spanCollectWorkflowRuns               = "failingtest.collect_workflow_runs"
	spanLoadAnalysisReport                = "failingtest.load_report"
	spanMergeAnalysisReport               = "failingtest.merge_report"
	spanPersistAnalysisReport             = "failingtest.persist_report"
	spanWriteAnalysisReport               = "failingtest.write_report"
	traceAttributeCutoff                  = "ktlo.analysis.cutoff"
	traceAttributeLookback                = "ktlo.analysis.lookback"
	traceAttributeWorkflowRunAttemptCount = "ktlo.workflow_run_attempt.count"
	traceAttributeWorkflowRecordCount     = "ktlo.workflow_record.count"
	analyzerInstrumentationName           = "ktlo/internal/failingtest"
)

const traceAttributeOverwriteReport = "ktlo.analysis.overwrite_report"

const traceAttributeKnownWorkflowRunAttemptCount = "ktlo.workflow_run_attempt.known_count"

type Analyzer struct {
	overwriteExistingReport bool
	workflowProvider        WorkflowProvider
	lookbackDuration        time.Duration
	testLogParsers          []TestLogParser
}

// NewAnalyzer uses the package's default parser set unless parsers are supplied.
func NewAnalyzer(workflowProvider WorkflowProvider, testLogParsers ...TestLogParser) (*Analyzer, error) {
	return newAnalyzer(workflowProvider, defaultLookbackDuration, false, testLogParsers...)
}

// NewAnalyzerWithLookbackDuration configures how far back workflow attempts are collected.
func NewAnalyzerWithLookbackDuration(workflowProvider WorkflowProvider, lookbackDuration time.Duration, testLogParsers ...TestLogParser) (*Analyzer, error) {
	return NewAnalyzerWithLookbackDurationAndOverwrite(workflowProvider, lookbackDuration, false, testLogParsers...)
}

// NewAnalyzerWithLookbackDurationAndOverwrite configures lookback and report reuse behavior.
func NewAnalyzerWithLookbackDurationAndOverwrite(workflowProvider WorkflowProvider, lookbackDuration time.Duration, overwriteExistingReport bool, testLogParsers ...TestLogParser) (*Analyzer, error) {
	if lookbackDuration <= 0 {
		return nil, errors.New(errorInvalidLookbackDuration)
	}
	return newAnalyzer(workflowProvider, lookbackDuration, overwriteExistingReport, testLogParsers...)
}

func newAnalyzer(workflowProvider WorkflowProvider, lookbackDuration time.Duration, overwriteExistingReport bool, testLogParsers ...TestLogParser) (*Analyzer, error) {
	if workflowProvider == nil {
		return nil, errors.New("workflow provider is required for the failing test check")
	}
	if len(testLogParsers) == 0 {
		testLogParsers = defaultTestLogParsers()
	}
	for _, parser := range testLogParsers {
		if parser == nil {
			return nil, errors.New("test log parser is required")
		}
	}
	testLogParsers = orderTestLogParsers(testLogParsers)
	return &Analyzer{
		overwriteExistingReport: overwriteExistingReport,
		workflowProvider:        workflowProvider,
		lookbackDuration:        lookbackDuration,
		testLogParsers:          append([]TestLogParser(nil), testLogParsers...),
	}, nil
}

func (analyzer *Analyzer) Run(ctx context.Context) (resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanAnalyzerRun,
		attribute.String(traceAttributeLookback, analyzer.lookbackDuration.String()),
		attribute.Bool(traceAttributeOverwriteReport, analyzer.overwriteExistingReport),
	)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanAnalyzerRun, startedAt, resultErr != nil)
	}()
	if err := ctx.Err(); err != nil {
		return err
	}

	now := time.Now().UTC()
	cutoff := now.Add(-analyzer.lookbackDuration)
	report, err := analyzer.loadReport(ctx)
	if err != nil {
		return err
	}
	knownAttempts := make(map[common.WorkflowRunAttemptKey]struct{})
	if !analyzer.overwriteExistingReport && !report.LookbackStart.IsZero() && !report.LookbackStart.After(cutoff) {
		for _, workflow := range report.Workflows {
			if workflow.CollectionComplete {
				knownAttempts[common.WorkflowRunAttemptKey{RunID: workflow.RunID, Attempt: workflow.RunAttempt}] = struct{}{}
			}
		}
	}
	workflowRuns, collectionErr := analyzer.collectWorkflowRuns(ctx, cutoff, knownAttempts)
	return analyzer.persistReport(ctx, report, workflowRuns, cutoff, now, collectionErr)
}

func (analyzer *Analyzer) loadReport(ctx context.Context) (loadedReport *testresults.Report, resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanLoadAnalysisReport)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanLoadAnalysisReport, startedAt, resultErr != nil)
	}()

	loadedReport = testresults.NewReport(analyzer.workflowProvider.Provider(), analyzer.workflowProvider.Repository())
	if analyzer.overwriteExistingReport {
		return loadedReport, nil
	}
	if err := loadedReport.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return loadedReport, nil
}

func (analyzer *Analyzer) collectWorkflowRuns(ctx context.Context, cutoff time.Time, knownAttempts map[common.WorkflowRunAttemptKey]struct{}) (workflowRuns []testresults.Workflow, resultErr error) {
	var workflowRunAttempts []WorkflowRunAttempt
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanCollectWorkflowRuns,
		attribute.Int(traceAttributeKnownWorkflowRunAttemptCount, len(knownAttempts)),
		attribute.String(traceAttributeCutoff, cutoff.Format(time.RFC3339Nano)),
	)
	defer func() {
		span.SetAttributes(
			attribute.Int(traceAttributeWorkflowRunAttemptCount, len(workflowRunAttempts)),
			attribute.Int(traceAttributeWorkflowRecordCount, len(workflowRuns)),
		)
		telemetry.FinishSpan(ctx, span, spanCollectWorkflowRuns, startedAt, resultErr != nil)
	}()
	workflowRunAttempts, listErr := analyzer.workflowProvider.ListWorkflowRunAttempts(ctx, cutoff, knownAttempts)
	workflowRuns, resultErr = ParseWorkflowRunAttempts(ctx, analyzer.workflowProvider, workflowRunAttempts, analyzer.testLogParsers...)
	return workflowRuns, errors.Join(listErr, resultErr)
}

func (analyzer *Analyzer) persistReport(ctx context.Context, report *testresults.Report, workflowRuns []testresults.Workflow, cutoff time.Time, generatedAt time.Time, collectionErr error) (resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanPersistAnalysisReport)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanPersistAnalysisReport, startedAt, resultErr != nil)
	}()

	mergeStartedAt := time.Now()
	mergeContext, mergeSpan := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanMergeAnalysisReport)
	mergeErr := report.MergeWorkflowRuns(workflowRuns, cutoff)
	telemetry.FinishSpan(mergeContext, mergeSpan, spanMergeAnalysisReport, mergeStartedAt, mergeErr != nil)
	if mergeErr != nil {
		return mergeErr
	}
	if collectionErr != nil {
		report.AnalysisStatus = testresults.AnalysisStatusIncomplete
	} else {
		report.AnalysisStatus = testresults.AnalysisStatusComplete
	}
	report.GeneratedAt = generatedAt
	report.LookbackStart = cutoff
	writeStartedAt := time.Now()
	writeContext, writeSpan := telemetry.StartSpan(ctx, analyzerInstrumentationName, spanWriteAnalysisReport)
	writeErr := report.Write()
	telemetry.FinishSpan(writeContext, writeSpan, spanWriteAnalysisReport, writeStartedAt, writeErr != nil)
	return errors.Join(collectionErr, writeErr)
}
