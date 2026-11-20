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
	"fmt"
	"html"
	"io"
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
	ciRollupError                    = string(github.PullRequestCIRollupStateError)
	ciRollupFailure                  = string(github.PullRequestCIRollupStateFailure)
	ciRollupSkipped                  = string(github.PullRequestCIRollupStateSkipped)
	ciRollupSuccess                  = string(github.PullRequestCIRollupStateSuccess)
	cicdOutcomeCancelled             = "cancelled"
	cicdOutcomeFailed                = "failed"
	cicdOutcomePassed                = "passed"
	cicdOutcomeSkipped               = "skipped"
	emptyString                      = ""
	errorCloseResultsFile            = "could not close results file: %w"
	errorCreateResultsDirectory      = "could not create results directory %q: %w"
	errorCreateResultsFile           = "could not create results file in %q: %w"
	errorInvalidResultsPath          = "invalid repository path for results: owner %q, repository %q"
	errorRenameResultsFile           = "could not move results file into place: %w"
	errorSetResultsFileMode          = "could not set results file permissions: %w"
	errorUnknownPullRequestAuthor    = "unknown author"
	errorWriteResultsFile            = "could not write results file: %w"
	errorWriteResultsMessage         = "could not write results confirmation: %w"
	formatWorkflowFailureParseError  = "- Could not parse a failed workflow: %s\n"
	formatMarkdownLink               = "[%s](%s)"
	formatMarkdownPullRequestText    = "%s (#%d)"
	formatMarkdownSection            = "## %s\n\n%s"
	formatHTMLAnchor                 = "<a href=\"%s\">%s</a>"
	formatHTMLCell                   = "<td>%s</td>"
	formatCICDOutcomeCount           = "%d %s"
	formatCICDStatusCounts           = "%s — %s"
	formatFailureOutputCodeBlock     = "<pre><code>%s</code></pre>"
	formatPRDetailsAnchorTarget      = "#pr-%d"
	htmlFailureDetailsTableHeader    = "<table>\n<thead>\n<tr><th>Failure</th><th>Classification</th><th>Latest failure</th><th>Main history</th><th>PR commit history</th><th>Related CI history</th></tr>\n</thead>\n<tbody>\n"
	htmlPRStatusTableHeader          = "<table>\n<thead>\n<tr><th>Author</th><th>Pull Request</th><th>CI/CD</th></tr>\n</thead>\n<tbody>\n"
	htmlStalePRTableHeader           = "<table>\n<thead>\n<tr><th>Author</th><th>Pull Request</th></tr>\n</thead>\n<tbody>\n"
	htmlTableEnd                     = "</tbody>\n</table>\n\n"
	htmlTableRowEnd                  = "</tr>\n"
	htmlTableRowStart                = "<tr>"
	htmlTextLineBreak                = "<br>"
	htmlFragmentPrefix               = "#"
	httpsScheme                      = "https"
	markdownBackslash                = "\\"
	markdownCarriageReturn           = "\r"
	markdownClosingBracket           = "]"
	markdownDocumentTitle            = "# Dependency Bump PR Status\n\n"
	markdownEmptyLine                = "\n"
	markdownEscapedBackslash         = "\\\\"
	markdownEscapedClosingBracket    = "\\]"
	markdownEscapedOpeningBracket    = "\\["
	markdownLineBreak                = "\n"
	markdownNewline                  = "<br>"
	markdownOpeningBracket           = "["
	markdownPipe                     = "\\|"
	markdownPipeCharacter            = "|"
	markdownWindowsNewline           = "\r\n"
	resultsDirectory                 = "results"
	resultsDirectoryMode             = 0o755
	resultsFileMode                  = 0o644
	resultsFileName                  = "DEPENDENCY-BUMP-PRS.md"
	resultsProviderGitHub            = "github"
	resultsTempFilePattern           = ".DEPENDENCY-BUMP-PRS-*.md"
	statusFailed                     = "Failed"
	statusNotReported                = "not reported"
	statusOthers                     = "Others"
	statusPending                    = "Pending"
	statusPassed                     = "Passed"
	statusStale                      = "Stale"
	successResultsFileWrittenMessage = "Markdown results written successfully to %s.\n"
)

const (
	failureDetailsHeadingFormat      = "<a id=\"pr-%d\"></a>\n### %s\n\n"
	failureDetailLineFormat          = "- %s: %s\n"
	formatGitHubCommitURL            = "https://github.com/%s/%s/commit/%s"
	formatInferredWorkflowPasses     = "%s; includes %d inferred passes from successful workflows"
	formatInferredWorkflowPass       = "%s; includes %d inferred pass from a successful workflow"
	formatObservationCounts          = "%d/%d passed (%d%%), %d/%d failed (%d%%)"
	formatMainBaselineLabel          = "%s — last %s"
	formatPartialLabel               = "%s (%s)"
	formatRelatedCIHistoryLabel      = "%s — %s"
	separatorCICDOutcomeCounts       = ", "
	formatPackageIdentity            = "package %s"
	formatPRCommitLabel              = "%s %s"
	formatUnresolvedCISection        = "## CI checks not completed\n\n"
	formatUnresolvedPRHeading        = "<a id=\"pr-%d\"></a>\n### %s\n"
	formatTestIdentity               = "test %s"
	emptyTableCell                   = "—"
	formatUnavailableReasonNotice    = "Unknown — %s\n"
	incompleteReportNotice           = "Partial — the report is incomplete; findings use only available observations.\n"
	formatDay                        = "%d day"
	formatDays                       = "%d days"
	formatHour                       = "%d hour"
	formatHours                      = "%d hours"
	formatMinute                     = "%d minute"
	formatMinutes                    = "%d minutes"
	explanationHeadSHAMismatch       = "The report has observations for a different PR head SHA, so they cannot classify the current revision."
	explanationIncompleteLogCoverage = "Available observations do not establish whether this failure is broken or flaky."
	explanationInsufficientHistory   = "There is not enough matching history to classify this failure."
	explanationMalformedReport       = "The local flaky-analysis report is malformed, so its findings could not be used."
	explanationMissingReport         = "The local flaky-analysis report is missing, so no historical failure analysis is available."
	explanationReportReadError       = "The local flaky-analysis report could not be read, so its findings could not be used."
	explanationRepositoryMismatch    = "The local flaky-analysis report belongs to a different repository, so its findings could not be used."
	explanationUnidentifiableFailure = "The failed CI context could not be matched to a failure in the analysis report."
	explanationUsesInferredPasses    = "Includes inferred workflow passes."
	markdownIdentitySeparator        = " / "
	statusBroken                     = "Broken"
	statusFlaky                      = "Flaky"
	shortSHALength                   = 7
	labelCancelledCICheck            = "Cancelled check"
	labelMainBaseline                = "main"
	labelMainPass                    = "main pass"
	labelMainFailure                 = "main failure"
	labelPartial                     = "partial"
	labelPRCommitPass                = "%s pass"
	labelPRCommitFailure             = "%s failure"
	labelTestFailureRun              = "Failure run"
	labelUnknownBranch               = "unknown branch"
	markdownLineEnding               = "\n"
	unknownClassificationLabel       = "Unknown"
	unknownCIContextName             = "unknown CI check"
	unknownFailureIdentity           = "Failed CI rollup"
)

const (
	outputInstrumentationName      = "ktlo/internal/output"
	spanPrintPullRequestResults    = "output.print_pull_request_ci_results"
	spanWriteResultsFile           = "output.write_results_file"
	traceAttributePullRequestCount = "github.pull_request.count"
	traceAttributeRepositoryName   = "github.repository.name"
	traceAttributeRepositoryOwner  = "github.repository.owner"
	traceAttributeResultByteCount  = "ktlo.result.byte_count"
)

type PullRequestResult struct {
	PullRequest                github.PullRequest
	CIResults                  github.PullRequestCIResults
	IsStale                    bool
	FailureAssessments         []testresults.FailureAssessment
	AnalysisUnavailableReason  testresults.UnknownReason
	WorkflowFailureParseErrors []string
	Lookback                   time.Duration
}

type ciOutcomeCounts struct {
	cancelled int
	failed    int
	passed    int
	skipped   int
}

func (counts *ciOutcomeCounts) add(ciContext github.PullRequestCIContext) {
	switch ciContext.EffectiveState() {
	case github.PullRequestCIRollupStateCancelled:
		counts.cancelled++
	case github.PullRequestCIRollupStateError, github.PullRequestCIRollupStateFailure:
		counts.failed++
	case github.PullRequestCIRollupStateSkipped:
		counts.skipped++
	case github.PullRequestCIRollupStateSuccess:
		counts.passed++
	}
}

func (counts *ciOutcomeCounts) summary() string {
	outcomeCounts := make([]string, 0, 4)
	if counts.passed > 0 {
		outcomeCounts = append(outcomeCounts, fmt.Sprintf(formatCICDOutcomeCount, counts.passed, cicdOutcomePassed))
	}
	if counts.failed > 0 {
		outcomeCounts = append(outcomeCounts, fmt.Sprintf(formatCICDOutcomeCount, counts.failed, cicdOutcomeFailed))
	}
	if counts.skipped > 0 {
		outcomeCounts = append(outcomeCounts, fmt.Sprintf(formatCICDOutcomeCount, counts.skipped, cicdOutcomeSkipped))
	}
	if counts.cancelled > 0 {
		outcomeCounts = append(outcomeCounts, fmt.Sprintf(formatCICDOutcomeCount, counts.cancelled, cicdOutcomeCancelled))
	}
	return strings.Join(outcomeCounts, separatorCICDOutcomeCounts)
}

func PrintPullRequestCIResults(ctx context.Context, stdout io.Writer, repository github.Repository, results []PullRequestResult) (resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, outputInstrumentationName, spanPrintPullRequestResults,
		attribute.String(traceAttributeRepositoryName, repository.Name),
		attribute.String(traceAttributeRepositoryOwner, repository.Owner),
		attribute.Int(traceAttributePullRequestCount, len(results)),
	)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanPrintPullRequestResults, startedAt, resultErr != nil)
	}()

	sortedResults := append([]PullRequestResult(nil), results...)
	sort.SliceStable(sortedResults, func(firstIndex int, secondIndex int) bool {
		return sortedResults[firstIndex].PullRequest.Number > sortedResults[secondIndex].PullRequest.Number
	})
	analysisResultsByNumber := make(map[int]github.PullRequestCIResults, len(sortedResults))
	activeResults := make([]PullRequestResult, 0, len(sortedResults))
	pullRequests := make([]github.PullRequest, 0, len(sortedResults))
	stalePullRequests := make([]github.PullRequest, 0)
	for _, result := range sortedResults {
		if result.IsStale {
			stalePullRequests = append(stalePullRequests, result.PullRequest)
			continue
		}
		activeResults = append(activeResults, result)
		pullRequests = append(pullRequests, result.PullRequest)
		analysisResultsByNumber[result.CIResults.Number] = result.CIResults
	}

	var content bytes.Buffer
	if _, err := io.WriteString(&content, markdownDocumentTitle); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	if err := writePullRequestGroup(&content, statusPassed, htmlPRStatusTableHeader, classifyPullRequests(pullRequests, analysisResultsByNumber, ciRollupSuccess, ciRollupSkipped), analysisResultsByNumber); err != nil {
		return err
	}
	if err := writePullRequestGroup(&content, statusPending, htmlPRStatusTableHeader, classifyPendingPullRequests(pullRequests, analysisResultsByNumber), analysisResultsByNumber); err != nil {
		return err
	}
	if err := writePullRequestGroup(&content, statusFailed, htmlPRStatusTableHeader, classifyPullRequests(pullRequests, analysisResultsByNumber, ciRollupFailure, ciRollupError), analysisResultsByNumber); err != nil {
		return err
	}
	if err := writePullRequestGroup(&content, statusOthers, htmlPRStatusTableHeader, classifyOtherPullRequests(pullRequests, analysisResultsByNumber), analysisResultsByNumber); err != nil {
		return err
	}
	if err := writeStalePullRequestGroup(&content, stalePullRequests); err != nil {
		return err
	}
	if err := writeFailedPullRequestDetails(&content, repository, activeResults); err != nil {
		return err
	}
	if err := writeUnresolvedCICheckDetails(&content, activeResults); err != nil {
		return err
	}

	resultsFilePath, err := writeResultsFile(ctx, content.Bytes(), repository)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, successResultsFileWrittenMessage, resultsFilePath); err != nil {
		return fmt.Errorf(errorWriteResultsMessage, err)
	}
	return nil
}

func writeStalePullRequestGroup(writer io.Writer, pullRequests []github.PullRequest) error {
	if len(pullRequests) == 0 {
		return nil
	}
	if _, err := fmt.Fprintf(writer, formatMarkdownSection, statusStale, htmlStalePRTableHeader); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	for _, pullRequest := range pullRequests {
		authorLogin := errorUnknownPullRequestAuthor
		if pullRequest.Author != nil && pullRequest.Author.Login != emptyString {
			authorLogin = pullRequest.Author.Login
		}
		pullRequestLinkText := fmt.Sprintf(formatMarkdownPullRequestText, pullRequest.Title, pullRequest.Number)
		pullRequestLink := htmlAnchor(pullRequestLinkText, pullRequest.URL)
		if err := writeHTMLTableRow(writer, escapeHTMLTableText(authorLogin), pullRequestLink); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(writer, htmlTableEnd); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	return nil
}

func writeFailedPullRequestDetails(writer io.Writer, repository github.Repository, results []PullRequestResult) error {
	for _, result := range results {
		if !isFailedRollup(result.CIResults.EffectiveRollupState()) {
			continue
		}
		if err := writePullRequestFailureDetails(writer, repository, result); err != nil {
			return err
		}
	}
	return nil
}

func writeUnresolvedCICheckDetails(writer io.Writer, results []PullRequestResult) error {
	var sectionWritten bool
	for _, result := range results {
		if isFailedRollup(result.CIResults.EffectiveRollupState()) {
			continue
		}
		var cancelledChecks []github.PullRequestCIContext
		for _, ciContext := range result.CIResults.Contexts {
			if ciContext.EffectiveState() == github.PullRequestCIRollupStateCancelled {
				cancelledChecks = append(cancelledChecks, ciContext)
			}
		}
		if len(cancelledChecks) == 0 {
			continue
		}
		if !sectionWritten {
			if _, err := io.WriteString(writer, formatUnresolvedCISection); err != nil {
				return fmt.Errorf(errorWriteResultsFile, err)
			}
			sectionWritten = true
		}

		var pullRequestLink string
		pullRequestLinkText := fmt.Sprintf(formatMarkdownPullRequestText, result.PullRequest.Title, result.PullRequest.Number)
		pullRequestLink = fmt.Sprintf(formatMarkdownLink, escapeMarkdownLinkText(pullRequestLinkText), escapeMarkdown(result.PullRequest.URL))
		if _, err := fmt.Fprintf(writer, formatUnresolvedPRHeading, result.PullRequest.Number, pullRequestLink); err != nil {
			return fmt.Errorf(errorWriteResultsFile, err)
		}
		for _, cancelledCheck := range cancelledChecks {
			checkName := ciContextName(cancelledCheck)
			checkURL := cancelledCheck.DetailsURL
			if checkURL == emptyString {
				checkURL = cancelledCheck.TargetURL
			}
			if checkURL != emptyString {
				checkName = fmt.Sprintf(formatMarkdownLink, escapeMarkdownLinkText(checkName), escapeMarkdown(checkURL))
			}
			if err := writeFailureDetailLine(writer, labelCancelledCICheck, checkName); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(writer, markdownEmptyLine); err != nil {
			return fmt.Errorf(errorWriteResultsFile, err)
		}
	}
	return nil
}

func ciContextName(ciContext github.PullRequestCIContext) string {
	name := ciContext.NameForDisplay()
	if name != emptyString {
		return name
	}
	return unknownCIContextName
}

func writePullRequestFailureDetails(writer io.Writer, repository github.Repository, result PullRequestResult) error {
	var pullRequestLink string
	pullRequestLinkText := fmt.Sprintf(formatMarkdownPullRequestText, result.PullRequest.Title, result.PullRequest.Number)
	pullRequestLink = fmt.Sprintf(formatMarkdownLink, escapeMarkdownLinkText(pullRequestLinkText), escapeMarkdown(result.PullRequest.URL))
	if _, err := fmt.Fprintf(writer, failureDetailsHeadingFormat, result.PullRequest.Number, pullRequestLink); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	if err := writeAnalysisUnavailableNotice(writer, result); err != nil {
		return err
	}
	for _, parseError := range result.WorkflowFailureParseErrors {
		if _, err := fmt.Fprintf(writer, formatWorkflowFailureParseError, escapeMarkdown(parseError)); err != nil {
			return fmt.Errorf(errorWriteResultsFile, err)
		}
	}
	if result.AnalysisUnavailableReason != "" {
		if _, err := io.WriteString(writer, markdownEmptyLine); err != nil {
			return fmt.Errorf(errorWriteResultsFile, err)
		}
	}

	assessments := append([]testresults.FailureAssessment(nil), result.FailureAssessments...)
	sort.SliceStable(assessments, func(firstIndex int, secondIndex int) bool {
		return assessmentSortKey(assessments[firstIndex]) < assessmentSortKey(assessments[secondIndex])
	})
	writtenFailure := false
	visibleAssessments := make([]testresults.FailureAssessment, 0, len(assessments))
	for _, assessment := range assessments {
		if hasMoreSpecificAssessment(assessment, assessments) || isExplainedByFailedTest(assessment, assessments) {
			continue
		}
		visibleAssessments = append(visibleAssessments, assessment)
	}
	unmatchedContexts := make([]github.PullRequestCIContext, 0)
	for _, failureContext := range result.CIResults.Contexts {
		if !isFailureContext(failureContext) || matchesAssessment(failureContext, assessments) {
			continue
		}
		unmatchedContexts = append(unmatchedContexts, failureContext)
	}
	if len(visibleAssessments) > 0 || len(unmatchedContexts) > 0 {
		if err := writeFailureDetailsTable(writer, repository, result, visibleAssessments, unmatchedContexts); err != nil {
			return err
		}
		writtenFailure = true
	}
	if !writtenFailure && result.AnalysisUnavailableReason == "" {
		if err := writeFailureDetailsTable(writer, repository, result, nil, nil); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(writer, markdownEmptyLine); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	return nil
}

func writeAnalysisUnavailableNotice(writer io.Writer, result PullRequestResult) error {
	if result.AnalysisUnavailableReason == "" {
		return nil
	}
	if result.AnalysisUnavailableReason == testresults.UnknownReasonIncompleteLogCoverage {
		if _, err := io.WriteString(writer, incompleteReportNotice); err != nil {
			return fmt.Errorf(errorWriteResultsFile, err)
		}
		return nil
	}
	reason := unknownReasonExplanation(result.AnalysisUnavailableReason)
	if _, err := fmt.Fprintf(writer, formatUnavailableReasonNotice, escapeMarkdown(reason)); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	return nil
}

func writeFailureDetailsTable(writer io.Writer, repository github.Repository, result PullRequestResult, assessments []testresults.FailureAssessment, unmatchedContexts []github.PullRequestCIContext) error {
	if err := writeFailureDetailsTableHeader(writer); err != nil {
		return err
	}
	for _, assessment := range assessments {
		row := newFailureDetailsTableRow(repository, result, assessment)
		if err := row.write(writer); err != nil {
			return err
		}
	}
	for _, failureContext := range unmatchedContexts {
		if err := writeUnmatchedFailureContextRow(writer, failureContext); err != nil {
			return err
		}
	}
	if len(assessments) == 0 && len(unmatchedContexts) == 0 {
		if err := writeUnidentifiableRollupRow(writer); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(writer, htmlTableEnd); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	return nil
}

func writeFailureDetailsTableHeader(writer io.Writer) error {
	if _, err := io.WriteString(writer, htmlFailureDetailsTableHeader); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	return nil
}

type failureDetailsTableRow struct {
	failure          string
	classification   string
	latestFailure    string
	mainHistory      string
	prCommitHistory  string
	relatedCIHistory string
}

func newFailureDetailsTableRow(repository github.Repository, result PullRequestResult, assessment testresults.FailureAssessment) failureDetailsTableRow {
	failureIdentity := assessmentIdentity(assessment.Location)
	if isFailedTestAssessment(assessment) {
		failureIdentity = fmt.Sprintf(formatTestIdentity, assessment.Location.TestID)
	}
	row := failureDetailsTableRow{
		failure:          failureIdentity,
		classification:   classificationCell(assessment, result),
		latestFailure:    latestFailureHTMLCell(assessment.LatestFailureEvidence()),
		mainHistory:      mainHistoryCell(result.Lookback, assessment.History.MainBranch, assessment.Partial),
		prCommitHistory:  prCommitHistoryCell(repository, assessment.History.PRCommits, assessment.Partial),
		relatedCIHistory: relatedCIHistoryCell(repository, result.Lookback, assessment.RelatedCI, assessment.Partial),
	}
	if !isFailedTestAssessment(assessment) {
		row.mainHistory = emptyTableCell
		row.prCommitHistory = emptyTableCell
		row.relatedCIHistory = emptyTableCell
	}
	return row
}

func (row *failureDetailsTableRow) write(writer io.Writer) error {
	return writeHTMLTableRow(writer,
		escapeHTMLTableText(row.failure),
		row.classification,
		tableCell(row.latestFailure),
		tableCell(row.mainHistory),
		tableCell(row.prCommitHistory),
		tableCell(row.relatedCIHistory),
	)
}

func isFailedTestAssessment(assessment testresults.FailureAssessment) bool {
	return assessment.Location.TestID != emptyString && assessment.Location.TestKind != testresults.TestKindPackage
}

func classificationCell(assessment testresults.FailureAssessment, result PullRequestResult) string {
	classificationValues := []string{classificationLabel(assessment.Classification)}
	if assessment.Classification == testresults.ClassificationUnknown && (assessment.UnknownReason != result.AnalysisUnavailableReason || result.AnalysisUnavailableReason == testresults.UnknownReasonIncompleteLogCoverage) {
		classificationValues = append(classificationValues, unknownReasonExplanation(assessment.UnknownReason))
	}
	if !isFailedTestAssessment(assessment) && assessment.History.InferredWorkflowPassCount() > 0 {
		classificationValues = append(classificationValues, explanationUsesInferredPasses)
	}
	for valueIndex := range classificationValues {
		classificationValues[valueIndex] = escapeHTMLTableText(classificationValues[valueIndex])
	}
	return strings.Join(classificationValues, markdownNewline)
}

func latestFailureHTMLCell(evidence *testresults.FailureEvidence) string {
	if evidence == nil {
		return emptyTableCell
	}
	cellValues := make([]string, 0, 2)
	if evidence.Message != emptyString {
		cellValues = append(cellValues, formatFailureOutput(evidence.Message))
	}
	runURL := evidence.SourceRun.JobURL
	if runURL == emptyString {
		runURL = evidence.SourceRun.WorkflowURL
	}
	if runURL != emptyString {
		failureRunLink := htmlAnchor(labelTestFailureRun, runURL)
		cellValues = append(cellValues, failureRunLink)
	}
	if len(cellValues) == 0 {
		return emptyTableCell
	}
	return strings.Join(cellValues, markdownNewline)
}

func formatFailureOutput(message string) string {
	message = strings.ReplaceAll(message, markdownWindowsNewline, markdownLineBreak)
	message = strings.ReplaceAll(message, markdownCarriageReturn, markdownLineBreak)
	return fmt.Sprintf(formatFailureOutputCodeBlock, html.EscapeString(message))
}

func tableCell(value string) string {
	if value == emptyString {
		return emptyTableCell
	}
	return value
}

func mainHistoryCell(lookback time.Duration, outcome testresults.OutcomeSummary, partial bool) string {
	counts := outcomeCountsCell(outcome.Counts, partial)
	if counts == emptyString {
		return emptyString
	}
	baselineLabel := labelMainBaseline
	if lookback > 0 {
		baselineLabel = fmt.Sprintf(formatMainBaselineLabel, labelMainBaseline, formatDuration(lookback))
	}
	historyValues := []string{escapeHTMLTableText(counts)}
	if passLink := runReferenceHTMLLink(labelMainPass, outcome.LatestPass); passLink != emptyString {
		historyValues = append(historyValues, passLink)
	}
	if outcome.LatestFailure != nil {
		if failureLink := runReferenceHTMLLink(labelMainFailure, &outcome.LatestFailure.SourceRun); failureLink != emptyString {
			historyValues = append(historyValues, failureLink)
		}
	}
	return fmt.Sprintf(formatRelatedCIHistoryLabel, escapeHTMLTableText(baselineLabel), strings.Join(historyValues, markdownNewline))
}

func prCommitHistoryCell(repository github.Repository, commits []testresults.CommitOutcomeSummary, partial bool) string {
	commitHistories := make([]string, 0, len(commits))
	for _, commit := range commits {
		counts := outcomeCountsCell(commit.Counts, partial)
		if counts == emptyString {
			continue
		}
		commitLabel := prCommitLabel(commit.Branch, commit.CommitSHA)
		commitURL := fmt.Sprintf(formatGitHubCommitURL, repository.Owner, repository.Name, commit.CommitSHA)
		commitLink := htmlAnchor(commitLabel, commitURL)
		historyValues := []string{escapeHTMLTableText(counts)}
		if passLink := runReferenceHTMLLink(fmt.Sprintf(labelPRCommitPass, commitLabel), commit.LatestPass); passLink != emptyString {
			historyValues = append(historyValues, passLink)
		}
		if commit.LatestFailure != nil {
			if failureLink := runReferenceHTMLLink(fmt.Sprintf(labelPRCommitFailure, commitLabel), &commit.LatestFailure.SourceRun); failureLink != emptyString {
				historyValues = append(historyValues, failureLink)
			}
		}
		commitHistories = append(commitHistories, fmt.Sprintf(formatRelatedCIHistoryLabel, commitLink, strings.Join(historyValues, markdownNewline)))
	}
	return strings.Join(commitHistories, markdownNewline)
}

func relatedCIHistoryCell(repository github.Repository, lookback time.Duration, histories []testresults.CIContextHistory, partial bool) string {
	contextHistories := make([]string, 0)
	for _, history := range histories {
		contextName := assessmentIdentity(history.Location)
		mainCounts := outcomeCountsCell(history.History.MainBranch.Counts, partial)
		if mainCounts != emptyString {
			mainLabel := labelMainBaseline
			if lookback > 0 {
				mainLabel = fmt.Sprintf(formatMainBaselineLabel, labelMainBaseline, formatDuration(lookback))
			}
			mainHistoryValues := []string{escapeHTMLTableText(mainCounts)}
			if passLink := runReferenceHTMLLink(fmt.Sprintf(formatRelatedCIHistoryLabel, contextName, labelMainPass), history.History.MainBranch.LatestPass); passLink != emptyString {
				mainHistoryValues = append(mainHistoryValues, passLink)
			}
			if history.History.MainBranch.LatestFailure != nil {
				failureLabel := fmt.Sprintf(formatRelatedCIHistoryLabel, contextName, labelMainFailure)
				if failureLink := runReferenceHTMLLink(failureLabel, &history.History.MainBranch.LatestFailure.SourceRun); failureLink != emptyString {
					mainHistoryValues = append(mainHistoryValues, failureLink)
				}
			}
			contextHistories = append(contextHistories, fmt.Sprintf(formatRelatedCIHistoryLabel, escapeHTMLTableText(contextName), fmt.Sprintf(formatRelatedCIHistoryLabel, escapeHTMLTableText(mainLabel), strings.Join(mainHistoryValues, markdownNewline))))
		}
		for _, commit := range history.History.PRCommits {
			counts := outcomeCountsCell(commit.Counts, partial)
			if counts == emptyString {
				continue
			}
			commitLabel := prCommitLabel(commit.Branch, commit.CommitSHA)
			commitURL := fmt.Sprintf(formatGitHubCommitURL, repository.Owner, repository.Name, commit.CommitSHA)
			commitLink := htmlAnchor(commitLabel, commitURL)
			commitHistoryValues := []string{escapeHTMLTableText(counts)}
			if passLink := runReferenceHTMLLink(fmt.Sprintf(labelPRCommitPass, commitLabel), commit.LatestPass); passLink != emptyString {
				commitHistoryValues = append(commitHistoryValues, passLink)
			}
			if commit.LatestFailure != nil {
				if failureLink := runReferenceHTMLLink(fmt.Sprintf(labelPRCommitFailure, commitLabel), &commit.LatestFailure.SourceRun); failureLink != emptyString {
					commitHistoryValues = append(commitHistoryValues, failureLink)
				}
			}
			contextHistories = append(contextHistories, fmt.Sprintf(formatRelatedCIHistoryLabel, escapeHTMLTableText(contextName), fmt.Sprintf(formatRelatedCIHistoryLabel, commitLink, strings.Join(commitHistoryValues, markdownNewline))))
		}
	}
	return strings.Join(contextHistories, markdownNewline)
}

func outcomeCountsCell(counts testresults.OutcomeCounts, partial bool) string {
	total := counts.Passes + counts.Failures
	if total == 0 {
		return emptyString
	}
	passPercentage, err := counts.PassPercentage()
	if err != nil {
		return emptyString
	}
	countsText := fmt.Sprintf(formatObservationCounts, counts.Passes, total, passPercentage, counts.Failures, total, 100-passPercentage)
	if counts.InferredPasses == 1 {
		countsText = fmt.Sprintf(formatInferredWorkflowPass, countsText, counts.InferredPasses)
	} else if counts.InferredPasses > 1 {
		countsText = fmt.Sprintf(formatInferredWorkflowPasses, countsText, counts.InferredPasses)
	}
	if partial {
		countsText = fmt.Sprintf(formatPartialLabel, countsText, labelPartial)
	}
	return countsText
}

func runReferenceHTMLLink(label string, reference *testresults.RunReference) string {
	if reference == nil {
		return emptyString
	}
	runURL := reference.JobURL
	if runURL == "" {
		runURL = reference.WorkflowURL
	}
	if runURL == "" {
		return emptyString
	}
	return htmlAnchor(label, runURL)
}

func writeUnmatchedFailureContextRow(writer io.Writer, failureContext github.PullRequestCIContext) error {
	identity := ciContextName(failureContext)
	message := failureContext.Description
	if message == "" {
		message = failureContext.Summary
	}
	if message == "" {
		message = failureContext.Text
	}
	contextURL := failureContext.DetailsURL
	if contextURL == "" {
		contextURL = failureContext.TargetURL
	}
	latestFailure := emptyTableCell
	if message != "" {
		latestFailure = formatFailureOutput(message)
	}
	if contextURL != "" {
		failureRunLink := htmlAnchor(labelTestFailureRun, contextURL)
		if latestFailure == emptyTableCell {
			latestFailure = failureRunLink
		} else {
			latestFailure = strings.Join([]string{latestFailure, failureRunLink}, markdownNewline)
		}
	}
	row := failureDetailsTableRow{
		failure:          identity,
		classification:   unknownClassificationCell(testresults.UnknownReasonUnidentifiableFailure),
		latestFailure:    latestFailure,
		mainHistory:      emptyTableCell,
		prCommitHistory:  emptyTableCell,
		relatedCIHistory: emptyTableCell,
	}
	return row.write(writer)
}

func writeUnidentifiableRollupRow(writer io.Writer) error {
	row := failureDetailsTableRow{
		failure:          unknownFailureIdentity,
		classification:   unknownClassificationCell(testresults.UnknownReasonUnidentifiableFailure),
		latestFailure:    emptyTableCell,
		mainHistory:      emptyTableCell,
		prCommitHistory:  emptyTableCell,
		relatedCIHistory: emptyTableCell,
	}
	return row.write(writer)
}

func writeFailureDetailLine(writer io.Writer, label string, value string) error {
	if _, err := fmt.Fprintf(writer, failureDetailLineFormat, escapeMarkdown(label), escapeMarkdown(value)); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	return nil
}

func isFailedRollup(rollupState github.PullRequestCIRollupState) bool {
	return rollupState == github.PullRequestCIRollupStateFailure || rollupState == github.PullRequestCIRollupStateError
}

func unknownClassificationCell(reason testresults.UnknownReason) string {
	return strings.Join([]string{
		escapeHTMLTableText(unknownClassificationLabel),
		escapeHTMLTableText(unknownReasonExplanation(reason)),
	}, markdownNewline)
}

func isFailureContext(failureContext github.PullRequestCIContext) bool {
	state := failureContext.EffectiveState()
	return state == github.PullRequestCIRollupStateFailure || state == github.PullRequestCIRollupStateError
}

func matchesAssessment(failureContext github.PullRequestCIContext, assessments []testresults.FailureAssessment) bool {
	contextName := failureContext.Name
	if contextName == "" {
		contextName = failureContext.Context
	}
	contextURLs := []string{failureContext.DetailsURL, failureContext.TargetURL}
	for _, assessment := range assessments {
		location := assessment.Location
		if contextName != "" && (strings.EqualFold(contextName, location.WorkflowName) || strings.EqualFold(contextName, location.JobName) || strings.EqualFold(contextName, location.StepName) || strings.EqualFold(contextName, location.TestID)) {
			return true
		}
		evidence := assessment.LatestFailureEvidence()
		if evidence == nil {
			continue
		}
		for _, contextURL := range contextURLs {
			if contextURL != "" && (contextURL == evidence.SourceRun.JobURL || contextURL == evidence.SourceRun.WorkflowURL) {
				return true
			}
		}
	}
	return false
}

func hasMoreSpecificAssessment(assessment testresults.FailureAssessment, assessments []testresults.FailureAssessment) bool {
	parent := assessment.Location
	if parent.TestID != "" {
		return false
	}
	for _, candidate := range assessments {
		child := candidate.Location
		if child.WorkflowName != parent.WorkflowName || child.TestID == "" && child.JobName == parent.JobName && child.StepName == parent.StepName {
			continue
		}
		if parent.JobName == "" || parent.JobName == child.JobName && (parent.StepName == "" || parent.StepName == child.StepName) {
			return true
		}
	}
	return false
}

func isExplainedByFailedTest(assessment testresults.FailureAssessment, assessments []testresults.FailureAssessment) bool {
	packageLocation := assessment.Location
	if packageLocation.TestKind != testresults.TestKindPackage {
		return false
	}
	for _, candidate := range assessments {
		testLocation := candidate.Location
		if testLocation.TestKind == testresults.TestKindTest && testLocation.WorkflowName == packageLocation.WorkflowName && testLocation.JobName == packageLocation.JobName && testLocation.StepName == packageLocation.StepName {
			return true
		}
	}
	return false
}

func assessmentSortKey(assessment testresults.FailureAssessment) string {
	location := assessment.Location
	return strings.Join([]string{location.WorkflowName, location.JobName, location.StepName, location.TestID}, markdownLineEnding)
}

func assessmentIdentity(location testresults.CIFailureLocation) string {
	parts := make([]string, 0, 4)
	if location.WorkflowName != "" {
		parts = append(parts, location.WorkflowName)
	}
	if location.JobName != "" {
		parts = append(parts, location.JobName)
	}
	if location.StepName != "" {
		parts = append(parts, location.StepName)
	}
	if location.TestID != "" {
		if location.TestKind == testresults.TestKindPackage {
			parts = append(parts, fmt.Sprintf(formatPackageIdentity, location.TestID))
		} else {
			parts = append(parts, fmt.Sprintf(formatTestIdentity, location.TestID))
		}
	}
	return strings.Join(parts, markdownIdentitySeparator)
}

func classificationLabel(classification testresults.Classification) string {
	switch classification {
	case testresults.ClassificationBroken:
		return statusBroken
	case testresults.ClassificationFlaky:
		return statusFlaky
	default:
		return unknownClassificationLabel
	}
}

func unknownReasonExplanation(reason testresults.UnknownReason) string {
	switch reason {
	case testresults.UnknownReasonHeadSHAMismatch:
		return explanationHeadSHAMismatch
	case testresults.UnknownReasonIncompleteLogCoverage:
		return explanationIncompleteLogCoverage
	case testresults.UnknownReasonInsufficientHistory:
		return explanationInsufficientHistory
	case testresults.UnknownReasonMalformedReport:
		return explanationMalformedReport
	case testresults.UnknownReasonMissingReport:
		return explanationMissingReport
	case testresults.UnknownReasonReportReadError:
		return explanationReportReadError
	case testresults.UnknownReasonRepositoryMismatch:
		return explanationRepositoryMismatch
	case testresults.UnknownReasonUnidentifiableFailure:
		return explanationUnidentifiableFailure
	default:
		return explanationUnidentifiableFailure
	}
}

func formatDuration(duration time.Duration) string {
	if duration%(24*time.Hour) == 0 {
		days := int(duration / (24 * time.Hour))
		if days == 1 {
			return fmt.Sprintf(formatDay, days)
		}
		return fmt.Sprintf(formatDays, days)
	}
	if duration%time.Hour == 0 {
		hours := int(duration / time.Hour)
		if hours == 1 {
			return fmt.Sprintf(formatHour, hours)
		}
		return fmt.Sprintf(formatHours, hours)
	}
	if duration%time.Minute == 0 {
		minutes := int(duration / time.Minute)
		if minutes == 1 {
			return fmt.Sprintf(formatMinute, minutes)
		}
		return fmt.Sprintf(formatMinutes, minutes)
	}
	return duration.String()
}

func abbreviatedSHA(sha string) string {
	if len(sha) <= shortSHALength {
		return sha
	}
	// Seven leading SHA characters are used as the short commit identifier.
	return sha[:shortSHALength]
}

func prCommitLabel(branch string, sha string) string {
	if branch == emptyString {
		branch = labelUnknownBranch
	}
	return fmt.Sprintf(formatPRCommitLabel, branch, abbreviatedSHA(sha))
}

func classifyPullRequests(pullRequests []github.PullRequest, analysisResultsByNumber map[int]github.PullRequestCIResults, states ...string) []github.PullRequest {
	classifiedPullRequests := make([]github.PullRequest, 0)
	stateSet := make(map[string]struct{}, len(states))
	for _, state := range states {
		stateSet[state] = struct{}{}
	}
	for _, pullRequest := range pullRequests {
		analysisResult, found := analysisResultsByNumber[pullRequest.Number]
		if isPendingPullRequest(analysisResult, found) {
			continue
		}
		if found {
			if _, matches := stateSet[string(analysisResult.EffectiveRollupState())]; matches {
				classifiedPullRequests = append(classifiedPullRequests, pullRequest)
			}
		}
	}
	return classifiedPullRequests
}

func classifyPendingPullRequests(pullRequests []github.PullRequest, analysisResultsByNumber map[int]github.PullRequestCIResults) []github.PullRequest {
	pendingPullRequests := make([]github.PullRequest, 0)
	for _, pullRequest := range pullRequests {
		analysisResult, found := analysisResultsByNumber[pullRequest.Number]
		if isPendingPullRequest(analysisResult, found) {
			pendingPullRequests = append(pendingPullRequests, pullRequest)
		}
	}
	return pendingPullRequests
}

func classifyOtherPullRequests(pullRequests []github.PullRequest, analysisResultsByNumber map[int]github.PullRequestCIResults) []github.PullRequest {
	otherPullRequests := make([]github.PullRequest, 0)
	for _, pullRequest := range pullRequests {
		analysisResult, found := analysisResultsByNumber[pullRequest.Number]
		if isPendingPullRequest(analysisResult, found) {
			continue
		}
		rollupState := analysisResult.EffectiveRollupState()
		if !found || (rollupState != github.PullRequestCIRollupStateSuccess && rollupState != github.PullRequestCIRollupStateSkipped && rollupState != github.PullRequestCIRollupStateFailure && rollupState != github.PullRequestCIRollupStateError) {
			otherPullRequests = append(otherPullRequests, pullRequest)
		}
	}
	return otherPullRequests
}

func isPendingPullRequest(analysisResult github.PullRequestCIResults, found bool) bool {
	if !found {
		return true
	}
	if analysisResult.EffectiveRollupState() == github.PullRequestCIRollupState(emptyString) {
		return true
	}
	return analysisResult.HasPendingChecks()
}

func hasPullRequestDetails(analysisResult github.PullRequestCIResults) bool {
	if isFailedRollup(analysisResult.EffectiveRollupState()) {
		return true
	}
	for _, ciContext := range analysisResult.Contexts {
		switch ciContext.EffectiveState() {
		case github.PullRequestCIRollupStateCancelled:
			return true
		}
	}
	return false
}

func writePullRequestGroup(writer io.Writer, groupName string, header string, pullRequests []github.PullRequest, analysisResultsByNumber map[int]github.PullRequestCIResults) error {
	if len(pullRequests) == 0 {
		return nil
	}
	if _, err := fmt.Fprintf(writer, formatMarkdownSection, groupName, header); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	for _, pullRequest := range pullRequests {
		var analysisResult github.PullRequestCIResults
		var authorLogin string
		var found bool
		var pullRequestLinkText string
		var resultSummary string
		var rollupState string

		analysisResult, found = analysisResultsByNumber[pullRequest.Number]
		authorLogin = errorUnknownPullRequestAuthor
		if pullRequest.Author != nil && pullRequest.Author.Login != emptyString {
			authorLogin = pullRequest.Author.Login
		}
		rollupState = string(analysisResult.EffectiveRollupState())
		if !found || rollupState == emptyString {
			rollupState = statusNotReported
		}
		if rollupState == ciRollupSkipped {
			rollupState = ciRollupSuccess
		}
		resultSummary = escapeHTMLTableText(rollupState)
		if hasPullRequestDetails(analysisResult) {
			resultSummary = htmlAnchor(rollupState, fmt.Sprintf(formatPRDetailsAnchorTarget, pullRequest.Number))
		}
		var outcomeCounts ciOutcomeCounts
		for _, ciContext := range analysisResult.Contexts {
			outcomeCounts.add(ciContext)
		}
		if outcomeCountSummary := outcomeCounts.summary(); outcomeCountSummary != emptyString {
			resultSummary = fmt.Sprintf(formatCICDStatusCounts, resultSummary, escapeHTMLTableText(outcomeCountSummary))
		}
		pullRequestLinkText = fmt.Sprintf(formatMarkdownPullRequestText, pullRequest.Title, pullRequest.Number)
		if err := writeHTMLTableRow(writer,
			escapeHTMLTableText(authorLogin),
			htmlAnchor(pullRequestLinkText, pullRequest.URL),
			resultSummary,
		); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(writer, htmlTableEnd); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	return nil
}

func writeResultsFile(ctx context.Context, content []byte, repository github.Repository) (resultsFilePath string, resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, outputInstrumentationName, spanWriteResultsFile,
		attribute.String(traceAttributeRepositoryName, repository.Name),
		attribute.String(traceAttributeRepositoryOwner, repository.Owner),
		attribute.Int(traceAttributeResultByteCount, len(content)),
	)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanWriteResultsFile, startedAt, resultErr != nil)
	}()

	var directoryPath string
	var err error
	var file *os.File
	directoryPath, err = resultsDirectoryPath(repository)
	if err != nil {
		return emptyString, err
	}
	if err := os.MkdirAll(directoryPath, resultsDirectoryMode); err != nil {
		return emptyString, fmt.Errorf(errorCreateResultsDirectory, directoryPath, err)
	}
	resultsFilePath = filepath.Join(directoryPath, resultsFileName)
	file, err = os.CreateTemp(directoryPath, resultsTempFilePattern)
	if err != nil {
		return emptyString, fmt.Errorf(errorCreateResultsFile, directoryPath, err)
	}
	defer func() {
		_ = os.Remove(file.Name())
	}()

	if _, err := io.Copy(file, bytes.NewReader(content)); err != nil {
		_ = file.Close()
		return emptyString, fmt.Errorf(errorWriteResultsFile, err)
	}
	if err := file.Chmod(resultsFileMode); err != nil {
		_ = file.Close()
		return emptyString, fmt.Errorf(errorSetResultsFileMode, err)
	}
	if err := file.Close(); err != nil {
		return emptyString, fmt.Errorf(errorCloseResultsFile, err)
	}
	if err := os.Rename(file.Name(), resultsFilePath); err != nil {
		return emptyString, fmt.Errorf(errorRenameResultsFile, err)
	}
	return resultsFilePath, nil
}

func resultsDirectoryPath(repository github.Repository) (string, error) {
	if !filepath.IsLocal(repository.Owner) || repository.Owner == "." || repository.Owner == ".." || !filepath.IsLocal(repository.Name) || repository.Name == "." || repository.Name == ".." {
		return emptyString, fmt.Errorf(errorInvalidResultsPath, repository.Owner, repository.Name)
	}
	return filepath.Join(resultsDirectory, resultsProviderGitHub, repository.Owner, repository.Name), nil
}

func writeHTMLTableRow(writer io.Writer, values ...string) error {
	if _, err := io.WriteString(writer, htmlTableRowStart); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	for _, value := range values {
		if _, err := fmt.Fprintf(writer, formatHTMLCell, value); err != nil {
			return fmt.Errorf(errorWriteResultsFile, err)
		}
	}
	if _, err := io.WriteString(writer, htmlTableRowEnd); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	return nil
}

func htmlAnchor(label string, target string) string {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return escapeHTMLTableText(label)
	}
	if parsedURL.Scheme != emptyString {
		if !strings.EqualFold(parsedURL.Scheme, httpsScheme) || parsedURL.Host == emptyString {
			return escapeHTMLTableText(label)
		}
	} else if !strings.HasPrefix(target, htmlFragmentPrefix) || parsedURL.Path != emptyString {
		return escapeHTMLTableText(label)
	}
	return fmt.Sprintf(formatHTMLAnchor, html.EscapeString(target), html.EscapeString(label))
}

func escapeHTMLTableText(value string) string {
	value = strings.ReplaceAll(value, markdownWindowsNewline, markdownLineBreak)
	value = strings.ReplaceAll(value, markdownCarriageReturn, markdownLineBreak)
	return strings.ReplaceAll(html.EscapeString(value), markdownLineBreak, htmlTextLineBreak)
}

func escapeMarkdown(value string) string {
	value = strings.ReplaceAll(value, markdownWindowsNewline, markdownNewline)
	value = strings.ReplaceAll(value, markdownLineBreak, markdownNewline)
	value = strings.ReplaceAll(value, markdownCarriageReturn, markdownNewline)
	return strings.ReplaceAll(value, markdownPipeCharacter, markdownPipe)
}

func escapeMarkdownLinkText(value string) string {
	value = strings.ReplaceAll(value, markdownBackslash, markdownEscapedBackslash)
	value = strings.ReplaceAll(value, markdownOpeningBracket, markdownEscapedOpeningBracket)
	value = strings.ReplaceAll(value, markdownClosingBracket, markdownEscapedClosingBracket)
	return escapeMarkdown(value)
}
