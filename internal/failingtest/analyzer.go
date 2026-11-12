package failingtest

import (
	"context"
	"errors"
	"os"
	"time"

	"ktlo/internal/testresults"
)

const defaultLookbackDuration = 7 * 24 * time.Hour

const errorInvalidLookbackDuration = "lookback duration must be greater than zero"

// TestLogParser detects and parses a test log format into shared test observations.
type TestLogParser interface {
	CanParse(stepLog string) bool
	Parse(stepLog string) ([]testresults.Test, error)
}

type Analyzer struct {
	workflowProvider WorkflowProvider
	lookbackDuration time.Duration
	testLogParsers   []TestLogParser
}

// NewAnalyzer uses the package's default parser set unless parsers are supplied.
func NewAnalyzer(workflowProvider WorkflowProvider, testLogParsers ...TestLogParser) (*Analyzer, error) {
	return newAnalyzer(workflowProvider, defaultLookbackDuration, testLogParsers...)
}

// NewAnalyzerWithLookbackDuration configures how far back workflow attempts are collected.
func NewAnalyzerWithLookbackDuration(workflowProvider WorkflowProvider, lookbackDuration time.Duration, testLogParsers ...TestLogParser) (*Analyzer, error) {
	if lookbackDuration <= 0 {
		return nil, errors.New(errorInvalidLookbackDuration)
	}
	return newAnalyzer(workflowProvider, lookbackDuration, testLogParsers...)
}

func newAnalyzer(workflowProvider WorkflowProvider, lookbackDuration time.Duration, testLogParsers ...TestLogParser) (*Analyzer, error) {
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
	return &Analyzer{
		workflowProvider: workflowProvider,
		lookbackDuration: lookbackDuration,
		testLogParsers:   append([]TestLogParser(nil), testLogParsers...),
	}, nil
}

func (analyzer *Analyzer) Run(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	now := time.Now().UTC()
	cutoff := now.Add(-analyzer.lookbackDuration)
	report, err := analyzer.loadReport()
	if err != nil {
		return err
	}
	workflowRuns, collectionErr := analyzer.collectWorkflowRuns(ctx, cutoff)
	return analyzer.persistReport(report, workflowRuns, cutoff, now, collectionErr)
}

func (analyzer *Analyzer) loadReport() (*testresults.Report, error) {
	report := testresults.NewReport(analyzer.workflowProvider.Provider(), analyzer.workflowProvider.Repository())
	if err := report.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return report, nil
}

func (analyzer *Analyzer) collectWorkflowRuns(ctx context.Context, cutoff time.Time) ([]testresults.Workflow, error) {
	workflowRunAttempts, listErr := analyzer.workflowProvider.ListWorkflowRunAttempts(ctx, cutoff)
	collectionErrors := make([]error, 0)
	if listErr != nil {
		collectionErrors = append(collectionErrors, listErr)
	}
	builder := newWorkflowRecordBuilder(analyzer.workflowProvider, analyzer.testLogParsers)
	workflowRuns := make([]testresults.Workflow, 0, len(workflowRunAttempts))
	for _, workflowRunAttempt := range workflowRunAttempts {
		if err := ctx.Err(); err != nil {
			collectionErrors = append(collectionErrors, err)
			break
		}
		workflowRun, err := builder.build(ctx, workflowRunAttempt)
		if err != nil {
			collectionErrors = append(collectionErrors, err)
		}
		if workflowRun != nil {
			workflowRuns = append(workflowRuns, *workflowRun)
		}
	}
	return workflowRuns, errors.Join(collectionErrors...)
}

func (analyzer *Analyzer) persistReport(report *testresults.Report, workflowRuns []testresults.Workflow, cutoff time.Time, generatedAt time.Time, collectionErr error) error {
	if err := report.MergeWorkflowRuns(workflowRuns, cutoff); err != nil {
		return err
	}
	if collectionErr != nil {
		report.AnalysisStatus = testresults.AnalysisStatusIncomplete
	} else {
		report.AnalysisStatus = testresults.AnalysisStatusComplete
	}
	report.GeneratedAt = generatedAt
	return errors.Join(collectionErr, report.Write())
}
