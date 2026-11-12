package output

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"ktlo/internal/github"
)

const (
	ciRollupError                    = "ERROR"
	ciRollupFailure                  = "FAILURE"
	ciRollupSuccess                  = "SUCCESS"
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
	formatMarkdownLink               = "[%s](%s)"
	formatMarkdownPullRequestText    = "%s (#%d)"
	formatMarkdownRow                = "| %s | %s | %s |\n"
	formatMarkdownSection            = "## %s\n\n%s"
	markdownBackslash                = "\\"
	markdownCarriageReturn           = "\r"
	markdownClosingBracket           = "]"
	markdownDocumentTitle            = "# PR Status Results\n\n"
	markdownEmptyLine                = "\n"
	markdownEscapedBackslash         = "\\\\"
	markdownEscapedClosingBracket    = "\\]"
	markdownEscapedOpeningBracket    = "\\["
	markdownHeader                   = "| Author | Pull Request | CI/CD |\n| --- | --- | --- |\n"
	markdownLineBreak                = "\n"
	markdownNewline                  = "<br>"
	markdownOpeningBracket           = "["
	markdownPipe                     = "\\|"
	markdownPipeCharacter            = "|"
	markdownWindowsNewline           = "\r\n"
	resultsDirectory                 = "results"
	resultsDirectoryMode             = 0o755
	resultsFileMode                  = 0o644
	resultsFileName                  = "RESULTS.md"
	resultsProviderGitHub            = "github"
	resultsTempFilePattern           = ".RESULTS-*.md"
	statusFailed                     = "Failed"
	statusNotReported                = "not reported"
	statusOthers                     = "Others"
	statusPassed                     = "Passed"
	successResultsFileWrittenMessage = "Markdown results written successfully to %s.\n"
)

func PrintPullRequestCIResults(stdout io.Writer, repository github.Repository, pullRequests []github.PullRequest, analysisResults []github.PullRequestCIResults) error {
	analysisResultsByNumber := make(map[int]github.PullRequestCIResults, len(analysisResults))
	for _, analysisResult := range analysisResults {
		analysisResultsByNumber[analysisResult.Number] = analysisResult
	}

	var content bytes.Buffer
	if _, err := io.WriteString(&content, markdownDocumentTitle); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	if err := writePullRequestGroup(&content, statusPassed, markdownHeader, classifyPullRequests(pullRequests, analysisResultsByNumber, ciRollupSuccess), analysisResultsByNumber); err != nil {
		return err
	}
	if err := writePullRequestGroup(&content, statusFailed, markdownHeader, classifyPullRequests(pullRequests, analysisResultsByNumber, ciRollupFailure, ciRollupError), analysisResultsByNumber); err != nil {
		return err
	}
	if err := writePullRequestGroup(&content, statusOthers, markdownHeader, classifyOtherPullRequests(pullRequests, analysisResultsByNumber), analysisResultsByNumber); err != nil {
		return err
	}

	resultsFilePath, err := writeResultsFile(content.Bytes(), repository)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, successResultsFileWrittenMessage, resultsFilePath); err != nil {
		return fmt.Errorf(errorWriteResultsMessage, err)
	}
	return nil
}

func classifyPullRequests(pullRequests []github.PullRequest, analysisResultsByNumber map[int]github.PullRequestCIResults, states ...string) []github.PullRequest {
	classifiedPullRequests := make([]github.PullRequest, 0)
	stateSet := make(map[string]struct{}, len(states))
	for _, state := range states {
		stateSet[state] = struct{}{}
	}
	for _, pullRequest := range pullRequests {
		analysisResult, found := analysisResultsByNumber[pullRequest.Number]
		if found {
			if _, matches := stateSet[analysisResult.RollupState]; matches {
				classifiedPullRequests = append(classifiedPullRequests, pullRequest)
			}
		}
	}
	return classifiedPullRequests
}

func classifyOtherPullRequests(pullRequests []github.PullRequest, analysisResultsByNumber map[int]github.PullRequestCIResults) []github.PullRequest {
	otherPullRequests := make([]github.PullRequest, 0)
	for _, pullRequest := range pullRequests {
		analysisResult, found := analysisResultsByNumber[pullRequest.Number]
		if !found || (analysisResult.RollupState != ciRollupSuccess && analysisResult.RollupState != ciRollupFailure && analysisResult.RollupState != ciRollupError) {
			otherPullRequests = append(otherPullRequests, pullRequest)
		}
	}
	return otherPullRequests
}

func writePullRequestGroup(writer io.Writer, groupName string, header string, pullRequests []github.PullRequest, analysisResultsByNumber map[int]github.PullRequestCIResults) error {
	if _, err := fmt.Fprintf(writer, formatMarkdownSection, groupName, header); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	for _, pullRequest := range pullRequests {
		var analysisResult github.PullRequestCIResults
		var authorLogin string
		var found bool
		var pullRequestLinkText string
		var rollupState string

		analysisResult, found = analysisResultsByNumber[pullRequest.Number]
		authorLogin = errorUnknownPullRequestAuthor
		if pullRequest.Author != nil && pullRequest.Author.Login != emptyString {
			authorLogin = pullRequest.Author.Login
		}
		rollupState = analysisResult.RollupState
		if !found || rollupState == emptyString {
			rollupState = statusNotReported
		}
		pullRequestLinkText = fmt.Sprintf(formatMarkdownPullRequestText, pullRequest.Title, pullRequest.Number)
		if _, err := fmt.Fprintf(writer, formatMarkdownRow, escapeMarkdown(authorLogin), fmt.Sprintf(formatMarkdownLink, escapeMarkdownLinkText(pullRequestLinkText), escapeMarkdown(pullRequest.URL)), escapeMarkdown(rollupState)); err != nil {
			return fmt.Errorf(errorWriteResultsFile, err)
		}
	}
	if _, err := io.WriteString(writer, markdownEmptyLine); err != nil {
		return fmt.Errorf(errorWriteResultsFile, err)
	}
	return nil
}

func writeResultsFile(content []byte, repository github.Repository) (string, error) {
	var directoryPath string
	var err error
	var file *os.File
	var resultsFilePath string
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
