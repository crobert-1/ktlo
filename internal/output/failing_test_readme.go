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

package output

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"ktlo/internal/github"
	"ktlo/internal/telemetry"
	"ktlo/internal/testresults"
)

const (
	branchMain                          = "main"
	errorCloseFailingReadme             = "could not close failing-test README: %w"
	errorCreateFailingReadmeDir         = "could not create failing-test README directory %q: %w"
	errorCreateFailingReadmeFile        = "could not create failing-test README file in %q: %w"
	errorNilTestResultsReport           = "test results report is nil"
	errorRenameFailingReadme            = "could not move failing-test README into place: %w"
	errorSetFailingReadmeMode           = "could not set failing-test README permissions: %w"
	errorWriteFailingReadme             = "could not write failing-test README: %w"
	failingReadmeAnalysisStatus         = "Analysis status: %s\n\n"
	failingReadmeBrokenWorkflowHeading  = "## Broken workflow and job failures\n\n"
	failingReadmeCodeBlockFormat        = "<pre><code>%s</code></pre>"
	failingReadmeContextTableHeader     = "<table>\n<thead>\n<tr><th>Workflow</th><th>Job</th><th>Step</th><th>Classification</th><th>Failures / observations</th><th>Failure rate</th><th>Recent failure</th></tr>\n</thead>\n<tbody>\n"
	failingReadmeCountFormat            = "%d / %d"
	failingReadmeDateFormat             = "2006-01-02"
	failingReadmeFailingHeading         = "## Failing tests\n\n"
	failingReadmeFileName               = "FLAKY-AND-FAILING-TESTS.md"
	failingReadmeFlakyHeading           = "## Flaky tests\n\n"
	failingReadmeFlakyWorkflowHeading   = "## Flaky workflow and job failures\n\n"
	failingReadmeGeneratedAt            = "Report generated: %s\n\n"
	failingReadmeHTMLCellFormat         = "<td>%s</td>"
	failingReadmeHTMLRowEnd             = "</tr>\n"
	failingReadmeHTMLRowStart           = "<tr>"
	failingReadmeHTMLTableEnd           = "</tbody>\n</table>\n\n"
	failingReadmeHTMLTableLineBreak     = "<br>"
	failingReadmeHTMLTestTableHeader    = "<table>\n<thead>\n<tr><th>Workflow</th><th>Job</th><th>Step</th><th>Test</th><th>Classification / reason</th><th>Failures / observations</th><th>Failure rate</th><th>Recent failure</th></tr>\n</thead>\n<tbody>\n"
	failingReadmeHTMLURLFormat          = "<a href=\"%s\">Failure run</a>"
	failingReadmeNoFailures             = "No failures observed on `main`.\n\n"
	failingReadmeObservationRange       = "Main branch observations: %s through %s\n\n"
	failingReadmePartialNotice          = "This report is partial; counts, rates, and classifications use the available observations. A test marked Broken may have successful runs that were not gathered.\n\n"
	failingReadmeRateFormat             = "%d%%"
	failingReadmeSuccessMessage         = "Failing-test README written successfully to %s.\n"
	failingReadmeTempFilePattern        = ".FLAKY-AND-FAILING-TESTS-*.md"
	failingReadmeTimestampFormat        = "2006-01-02 15:04:05 UTC"
	failingReadmeTitle                  = "# Flaky and Failing Tests\n\n"
	failingReadmeUnknownHeading         = "## Tests with unknown classification\n\n"
	failingReadmeUnknownWorkflowHeading = "## Workflow and job failures with unknown classification\n\n"
	failingReadmeUnknownDescription     = "These failures were observed, but available history is insufficient or incomplete to determine whether they are flaky or consistently failing.\n\n"
	failingReadmeUnknownReasonFormat    = "Unknown — %s"
	failingReadmeUnavailableRate        = "Unavailable"
	failingReadmeURLScheme              = "https"
	spanWriteFailingTestReadme          = "output.write_failing_test_readme"
)

func WriteFailingTestReadme(ctx context.Context, report *testresults.Report, stdout io.Writer) (resultErr error) {
	if report == nil {
		return fmt.Errorf(errorWriteFailingReadme, errors.New(errorNilTestResultsReport))
	}
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, outputInstrumentationName, spanWriteFailingTestReadme,
		attribute.String(traceAttributeRepositoryName, report.Repository.Name),
		attribute.String(traceAttributeRepositoryOwner, report.Repository.Owner),
	)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanWriteFailingTestReadme, startedAt, resultErr != nil)
	}()
	assessments := report.AssessmentsForBranch(branchMain)
	sort.Slice(assessments, func(firstIndex int, secondIndex int) bool {
		firstKey := assessmentSortKey(assessments[firstIndex])
		secondKey := assessmentSortKey(assessments[secondIndex])
		if firstKey != secondKey {
			return firstKey < secondKey
		}
		return assessments[firstIndex].Location.TestKind < assessments[secondIndex].Location.TestKind
	})
	assessments = newFailingTestReadmeAssessmentFilter(report).suppressExplainedParentTests(assessments)
	content, err := failingTestReadmeContent(report, assessments)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.Int(traceAttributeResultByteCount, len(content)))
	filePath, err := writeFailingTestReadmeFile(report.Repository, content)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, failingReadmeSuccessMessage, filePath); err != nil {
		return fmt.Errorf(errorWriteResultsMessage, err)
	}
	return nil
}

func failingTestReadmeContent(report *testresults.Report, assessments []testresults.FailureAssessment) ([]byte, error) {
	var content bytes.Buffer
	if _, err := io.WriteString(&content, failingReadmeTitle); err != nil {
		return nil, fmt.Errorf(errorWriteFailingReadme, err)
	}
	if _, err := fmt.Fprintf(&content, failingReadmeGeneratedAt, report.GeneratedAt.UTC().Format(failingReadmeTimestampFormat)); err != nil {
		return nil, fmt.Errorf(errorWriteFailingReadme, err)
	}
	if _, err := fmt.Fprintf(&content, failingReadmeAnalysisStatus, report.AnalysisStatus); err != nil {
		return nil, fmt.Errorf(errorWriteFailingReadme, err)
	}
	firstObservation, lastObservation := mainObservationRange(report)
	if !firstObservation.IsZero() {
		if _, err := fmt.Fprintf(&content, failingReadmeObservationRange, firstObservation.Format(failingReadmeDateFormat), lastObservation.Format(failingReadmeDateFormat)); err != nil {
			return nil, fmt.Errorf(errorWriteFailingReadme, err)
		}
	}
	if report.AnalysisStatus == testresults.AnalysisStatusIncomplete {
		if _, err := io.WriteString(&content, failingReadmePartialNotice); err != nil {
			return nil, fmt.Errorf(errorWriteFailingReadme, err)
		}
	}
	flakyTests := make([]testresults.FailureAssessment, 0)
	failingTests := make([]testresults.FailureAssessment, 0)
	unknownTests := make([]testresults.FailureAssessment, 0)
	workflowFailureGroups := buildWorkflowFailureGroups(report)
	brokenWorkflowFailures := make([]workflowFailureGroup, 0)
	flakyWorkflowFailures := make([]workflowFailureGroup, 0)
	unknownWorkflowFailures := make([]workflowFailureGroup, 0)
	for _, assessment := range assessments {
		if assessment.Location.TestID == emptyString {
			continue
		}
		switch assessment.Classification {
		case testresults.ClassificationFlaky:
			flakyTests = append(flakyTests, assessment)
		case testresults.ClassificationBroken:
			failingTests = append(failingTests, assessment)
		default:
			unknownTests = append(unknownTests, assessment)
		}
	}
	for _, group := range workflowFailureGroups {
		switch group.classification {
		case testresults.ClassificationFlaky:
			flakyWorkflowFailures = append(flakyWorkflowFailures, group)
		case testresults.ClassificationBroken:
			brokenWorkflowFailures = append(brokenWorkflowFailures, group)
		default:
			unknownWorkflowFailures = append(unknownWorkflowFailures, group)
		}
	}
	if len(assessments) == 0 && len(workflowFailureGroups) == 0 {
		if _, err := io.WriteString(&content, failingReadmeNoFailures); err != nil {
			return nil, fmt.Errorf(errorWriteFailingReadme, err)
		}
	}
	sections := []struct {
		assessments    []testresults.FailureAssessment
		description    string
		heading        string
		workflowGroups []workflowFailureGroup
		testRows       bool
		workflowRows   bool
	}{
		{assessments: flakyTests, heading: failingReadmeFlakyHeading, testRows: true},
		{workflowGroups: flakyWorkflowFailures, heading: failingReadmeFlakyWorkflowHeading, workflowRows: true},
		{assessments: failingTests, heading: failingReadmeFailingHeading, testRows: true},
		{workflowGroups: brokenWorkflowFailures, heading: failingReadmeBrokenWorkflowHeading, workflowRows: true},
		{assessments: unknownTests, description: failingReadmeUnknownDescription, heading: failingReadmeUnknownHeading, testRows: true},
		{workflowGroups: unknownWorkflowFailures, description: failingReadmeUnknownDescription, heading: failingReadmeUnknownWorkflowHeading, workflowRows: true},
	}
	for _, section := range sections {
		var err error
		if section.workflowRows {
			err = writeFailingReadmeWorkflowSection(&content, section.heading, section.description, section.workflowGroups)
		} else {
			err = writeFailingReadmeSection(&content, section.heading, section.description, section.assessments, section.testRows)
		}
		if err != nil {
			return nil, err
		}
	}
	return content.Bytes(), nil
}

func writeFailingReadmeSection(writer io.Writer, heading string, description string, assessments []testresults.FailureAssessment, testRows bool) error {
	if len(assessments) == 0 {
		return nil
	}
	if _, err := io.WriteString(writer, heading); err != nil {
		return fmt.Errorf(errorWriteFailingReadme, err)
	}
	if description != emptyString {
		if _, err := io.WriteString(writer, description); err != nil {
			return fmt.Errorf(errorWriteFailingReadme, err)
		}
	}
	header := failingReadmeContextTableHeader
	if testRows {
		header = failingReadmeHTMLTestTableHeader
	}
	if _, err := io.WriteString(writer, header); err != nil {
		return fmt.Errorf(errorWriteFailingReadme, err)
	}
	for _, assessment := range assessments {
		if err := writeFailingReadmeRow(writer, assessment, testRows); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(writer, failingReadmeHTMLTableEnd); err != nil {
		return fmt.Errorf(errorWriteFailingReadme, err)
	}
	return nil
}

func writeFailingReadmeRow(writer io.Writer, assessment testresults.FailureAssessment, testRows bool) error {
	location := assessment.Location
	counts := assessment.History.MainBranch.Counts
	totalObservations := counts.Passes + counts.Failures
	countLabel := fmt.Sprintf(failingReadmeCountFormat, counts.Failures, totalObservations)
	rateLabel := failingReadmeUnavailableRate
	if totalObservations > 0 {
		rateLabel = fmt.Sprintf(failingReadmeRateFormat, int(math.Round(float64(counts.Failures)/float64(totalObservations)*100)))
	}
	classification := classificationLabel(assessment.Classification)
	if assessment.Classification == testresults.ClassificationUnknown {
		classification = fmt.Sprintf(failingReadmeUnknownReasonFormat, unknownReasonExplanation(assessment.UnknownReason))
	}
	failure := failingReadmeFailureCell(assessment.LatestFailureEvidence())
	if testRows {
		values := []string{
			escapeFailingReadmeHTMLCell(tableCell(location.WorkflowName)),
			escapeFailingReadmeHTMLCell(tableCell(location.JobName)),
			escapeFailingReadmeHTMLCell(tableCell(location.StepName)),
			escapeFailingReadmeHTMLCell(assessmentIdentity(testresults.CIFailureLocation{TestID: location.TestID, TestKind: location.TestKind})),
			escapeFailingReadmeHTMLCell(classification),
			escapeFailingReadmeHTMLCell(countLabel),
			escapeFailingReadmeHTMLCell(rateLabel),
			failure,
		}
		return writeFailingReadmeHTMLRow(writer, values)
	}
	values := []string{
		escapeFailingReadmeHTMLCell(tableCell(location.WorkflowName)),
		escapeFailingReadmeHTMLCell(tableCell(location.JobName)),
		escapeFailingReadmeHTMLCell(tableCell(location.StepName)),
		escapeFailingReadmeHTMLCell(classification),
		escapeFailingReadmeHTMLCell(countLabel),
		escapeFailingReadmeHTMLCell(rateLabel),
		failure,
	}
	return writeFailingReadmeHTMLRow(writer, values)
}

func writeFailingReadmeWorkflowSection(writer io.Writer, heading string, description string, groups []workflowFailureGroup) error {
	if len(groups) == 0 {
		return nil
	}
	if _, err := io.WriteString(writer, heading); err != nil {
		return fmt.Errorf(errorWriteFailingReadme, err)
	}
	if description != emptyString {
		if _, err := io.WriteString(writer, description); err != nil {
			return fmt.Errorf(errorWriteFailingReadme, err)
		}
	}
	if _, err := io.WriteString(writer, failingReadmeContextTableHeader); err != nil {
		return fmt.Errorf(errorWriteFailingReadme, err)
	}
	for groupIndex := range groups {
		if err := writeFailingReadmeWorkflowRow(writer, &groups[groupIndex]); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(writer, failingReadmeHTMLTableEnd); err != nil {
		return fmt.Errorf(errorWriteFailingReadme, err)
	}
	return nil
}

func writeFailingReadmeWorkflowRow(writer io.Writer, group *workflowFailureGroup) error {
	values := []string{
		escapeFailingReadmeHTMLCell(group.workflowName),
		escapeFailingReadmeHTMLCell(group.jobDisplayName()),
		escapeFailingReadmeHTMLCell(group.stepDisplayName()),
		escapeFailingReadmeHTMLCell(classificationLabel(group.classification)),
		escapeFailingReadmeHTMLCell(group.observationCountLabel()),
		escapeFailingReadmeHTMLCell(group.observedFailureRate()),
		failingReadmeFailureCell(group.latestFailureEvidence()),
	}
	return writeFailingReadmeHTMLRow(writer, values)
}

func writeFailingReadmeHTMLRow(writer io.Writer, values []string) error {
	if _, err := io.WriteString(writer, failingReadmeHTMLRowStart); err != nil {
		return fmt.Errorf(errorWriteFailingReadme, err)
	}
	for _, value := range values {
		if _, err := fmt.Fprintf(writer, failingReadmeHTMLCellFormat, value); err != nil {
			return fmt.Errorf(errorWriteFailingReadme, err)
		}
	}
	if _, err := io.WriteString(writer, failingReadmeHTMLRowEnd); err != nil {
		return fmt.Errorf(errorWriteFailingReadme, err)
	}
	return nil
}

func failingReadmeFailureCell(evidence *testresults.FailureEvidence) string {
	if evidence == nil {
		return emptyTableCell
	}
	cellValues := make([]string, 0, 2)
	if evidence.Message != emptyString {
		message := strings.ReplaceAll(evidence.Message, markdownWindowsNewline, markdownLineBreak)
		message = strings.ReplaceAll(message, markdownCarriageReturn, markdownLineBreak)
		cellValues = append(cellValues, fmt.Sprintf(failingReadmeCodeBlockFormat, html.EscapeString(message)))
	}
	runURL := evidence.SourceRun.JobURL
	if runURL == emptyString {
		runURL = evidence.SourceRun.WorkflowURL
	}
	parsedURL, err := url.Parse(runURL)
	if err == nil && parsedURL.Scheme == failingReadmeURLScheme && parsedURL.Host != emptyString {
		cellValues = append(cellValues, fmt.Sprintf(failingReadmeHTMLURLFormat, html.EscapeString(runURL)))
	}
	if len(cellValues) == 0 {
		return emptyTableCell
	}
	return strings.Join(cellValues, markdownNewline)
}

func escapeFailingReadmeHTMLCell(value string) string {
	value = strings.ReplaceAll(value, markdownWindowsNewline, markdownLineBreak)
	value = strings.ReplaceAll(value, markdownCarriageReturn, markdownLineBreak)
	return strings.ReplaceAll(html.EscapeString(value), markdownLineBreak, failingReadmeHTMLTableLineBreak)
}

func mainObservationRange(report *testresults.Report) (firstObservation time.Time, lastObservation time.Time) {
	for _, workflow := range report.Workflows {
		if workflow.Branch != branchMain || workflow.RunStartedAt.IsZero() {
			continue
		}
		if firstObservation.IsZero() || workflow.RunStartedAt.Before(firstObservation) {
			firstObservation = workflow.RunStartedAt
		}
		if lastObservation.IsZero() || workflow.RunStartedAt.After(lastObservation) {
			lastObservation = workflow.RunStartedAt
		}
	}
	return firstObservation, lastObservation
}

func writeFailingTestReadmeFile(repository github.Repository, content []byte) (filePath string, resultErr error) {
	directoryPath, err := resultsDirectoryPath(repository)
	if err != nil {
		return emptyString, err
	}
	if err := os.MkdirAll(directoryPath, resultsDirectoryMode); err != nil {
		return emptyString, fmt.Errorf(errorCreateFailingReadmeDir, directoryPath, err)
	}
	filePath = filepath.Join(directoryPath, failingReadmeFileName)
	temporaryFile, err := os.CreateTemp(directoryPath, failingReadmeTempFilePattern)
	if err != nil {
		return emptyString, fmt.Errorf(errorCreateFailingReadmeFile, directoryPath, err)
	}
	defer func() {
		_ = os.Remove(temporaryFile.Name())
	}()
	if _, err := temporaryFile.Write(content); err != nil {
		_ = temporaryFile.Close()
		return emptyString, fmt.Errorf(errorWriteFailingReadme, err)
	}
	if err := temporaryFile.Chmod(resultsFileMode); err != nil {
		_ = temporaryFile.Close()
		return emptyString, fmt.Errorf(errorSetFailingReadmeMode, err)
	}
	if err := temporaryFile.Close(); err != nil {
		return emptyString, fmt.Errorf(errorCloseFailingReadme, err)
	}
	if err := os.Rename(temporaryFile.Name(), filePath); err != nil {
		return emptyString, fmt.Errorf(errorRenameFailingReadme, err)
	}
	return filePath, nil
}
