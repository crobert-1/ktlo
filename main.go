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

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"ktlo/internal/common"
	"ktlo/internal/config"
	"ktlo/internal/github"
	"ktlo/internal/logging"
	"ktlo/internal/output"
	"ktlo/internal/prstatus"
	"ktlo/internal/telemetry"
	"ktlo/internal/vulnerability"
)

const (
	commandDependencyPRStatus             = "dependency-pr-status"
	commandFlakyTestAnalysis              = "flaky-test-analysis"
	commandGenerateReport                 = "generate-report"
	commandHelp                           = "help"
	commandName                           = "ktlo"
	commandUpdate                         = "update"
	commandVulnerabilityCheck             = "vulnerability-check"
	dashPrefix                            = "-"
	defaultConfigPath                     = "config.yaml"
	emptyString                           = ""
	equalsSign                            = "="
	errorConfigFlagValue                  = "option --config requires a file path"
	errorInvalidLookbackOverride          = "invalid --lookback value: %w"
	errorInvalidStaleAgeOverride          = "invalid --stale-age value: %w"
	errorMissingVulnerabilityID           = "vulnerability-check requires a vulnerability ID"
	errorRepositoryFailure                = "repository %q failed: %v"
	errorTelemetrySetup                   = "OpenTelemetry setup failed: %v"
	errorTelemetryShutdown                = "OpenTelemetry shutdown failed: %v"
	errorUnexpectedArguments              = "unexpected positional arguments: %v"
	errorUnknownCommand                   = "unknown command %q"
	errorVulnerabilityCheck               = "Vulnerability Check failed:\n  %v"
	exitCodeError                         = 1
	exitCodeUsage                         = 2
	flagConfigLong                        = "--config"
	flagConfigShort                       = "-c"
	flagHelpLong                          = "--help"
	flagHelpShort                         = "-h"
	flagLookback                          = "lookback"
	flagNoParseWorkflowFailures           = "no-parse-workflow-failures"
	flagOverwrite                         = "overwrite"
	flagParseWorkflowFailures             = "parse-workflow-failures"
	flagStaleAge                          = "stale-age"
	githubTokenEnv                        = "GH_TOKEN"
	githubTokenFallbackEnv                = "GITHUB_TOKEN"
	helpLookback                          = "Override the configured analysis window"
	helpNoParseWorkflowFailures           = "Skip targeted workflow log parsing for this run"
	helpOverwrite                         = "Rebuild the report instead of reusing stored data"
	helpParseWorkflowFailures             = "Parse failed workflow runs for more failure details (enabled by default)"
	helpStaleAge                          = "Override the configured stale PR threshold"
	instrumentationName                   = "ktlo/main"
	lineEnding                            = "\n"
	messageVulnerabilityReportDone        = "Markdown report written."
	messageVulnerabilityReportPath        = "Saved report: %s"
	messageVulnerabilityReportStart       = "Writing detailed Markdown reports to repository result files."
	messageVulnerabilityRepositoryFailure = "Vulnerability Check %s/%s %s failed: %s"
	shortcutDependencyPRStatus            = "d"
	shortcutFlakyTestAnalysis             = "f"
	shortcutGenerateReport                = "g"
	shortcutUpdate                        = "u"
	shortcutVulnerabilityCheck            = "v"
	telemetryShutdownTimeout              = 5 * time.Second
)

const durationParts = 3
const secondsPerHour = 3600
const secondsPerMinute = 60
const durationSeparator = ", "
const hoursUnit = "hour"
const messageRunDuration = "Total runtime: %s"
const minutesUnit = "minute"
const pluralDurationPart = "%d %s%s"
const pluralSuffix = "s"
const singularDurationPart = "%d %s"
const secondsUnit = "second"

const (
	spanFailingTestCheck = "ktlo.failing_test_check"
	spanRepository       = "ktlo.repository"
	spanRun              = "ktlo.run"
	spanRunChecks        = "ktlo.run_checks"
)

const traceAttributeRepositoryCount = "ktlo.repository.count"
const traceAttributeRepositoryName = "github.repository.name"
const traceAttributeRepositoryOwner = "github.repository.owner"

const (
	helpTextRoot = `Usage: ktlo [--config FILE] <command> [arguments]

Run a KTLO operation for repositories in the config file.

Commands:
  vulnerability-check (v) [ID]  Check repositories for a vulnerability
  flaky-test-analysis (f)       Analyze recent flaky-test data
  generate-report (g)           Generate the flaky-test Markdown report
  dependency-pr-status (d)       Manage dependency PR status
  help [COMMAND]                Show command help

Options:
  -c, --config FILE  Config file to use (default: config.yaml)
  -h, --help         Show help

Run "ktlo <command> --help" for command help.
`
	helpTextDependencyPRStatus = `Usage: ktlo dependency-pr-status <command> [options]
       ktlo d <command> [options]

Manage dependency pull request status for configured repositories.

Commands:
  update (u)  Update the dependency PR status report

Options:
  -c, --config FILE  Config file to use (default: config.yaml)
  -h, --help         Show this help

Run "ktlo dependency-pr-status update --help" for command help.
`
	helpTextDependencyPRStatusUpdate = `Usage: ktlo dependency-pr-status update [options]
       ktlo d update [options]
       ktlo d u [options]

Update the dependency PR status report for configured repositories.

Options:
  --parse-workflow-failures     Parse failed workflow runs (enabled by default)
  --no-parse-workflow-failures  Skip targeted workflow log parsing for this run
  --stale-age DURATION          Override the configured stale PR threshold
  -c, --config FILE             Config file to use (default: config.yaml)
  -h, --help                    Show this help

Examples:
  ktlo d update
  ktlo d u --config config.yaml
  ktlo d u --no-parse-workflow-failures
`
	helpTextFlakyTestAnalysis = `Usage: ktlo flaky-test-analysis [options]
       ktlo f [options]

Analyze recent workflow runs for configured repositories and update the
local flaky-test data.

Options:
  --lookback DURATION  Override the configured analysis window
  --overwrite          Rebuild the report instead of reusing stored data
  -c, --config FILE    Config file to use (default: config.yaml)
  -h, --help           Show this help

Examples:
  ktlo flaky-test-analysis
  ktlo f --lookback 30d --overwrite
`
	helpTextGenerateReport = `Usage: ktlo generate-report [options]
       ktlo g [options]

Generate the flaky-test Markdown report from local analysis data.

Options:
  -c, --config FILE  Config file to use (default: config.yaml)
  -h, --help         Show this help

Examples:
  ktlo generate-report
  ktlo g --config config.yaml
`
	helpTextVulnerabilityCheck = `Usage: ktlo vulnerability-check [options] <ID>
       ktlo v [options] <ID>

Check configured repositories for a vulnerability and write a report.

Arguments:
  ID  Vulnerability identifier, such as CVE-2026-1234 or GO-2026-1234

Options:
  -c, --config FILE  Config file to use (default: config.yaml)
  -h, --help         Show this help

Examples:
  ktlo v CVE-2026-1234
  ktlo --config config.yaml vulnerability-check GO-2026-1234
`
)

type cliInvocation struct {
	command               string
	configPath            string
	lookback              string
	overwrite             bool
	parseWorkflowFailures bool
	staleAge              string
	vulnerabilityID       string
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) (exitCode int) {
	invocation, helpText, err := parseCLIInvocation(args)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		if helpText != emptyString {
			_, _ = fmt.Fprint(stderr, helpText)
		} else {
			_, _ = fmt.Fprint(stderr, helpTextRoot)
		}
		return exitCodeUsage
	}
	if helpText != emptyString {
		_, _ = fmt.Fprint(stdout, helpText)
		return 0
	}

	loadedConfig, err := config.Load(invocation.configPath)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return exitCodeUsage
	}
	if err := validateCommandConfig(invocation, loadedConfig); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return exitCodeUsage
	}

	rootContext := context.Background()
	startedAt := time.Now()
	logger := logging.New(stderr, commandLogPrefix(invocation.command), loadedConfig.EffectiveLogLevel())
	defer func() {
		logger.Progressf(messageRunDuration, formatRunDuration(time.Since(startedAt)))
	}()

	shutdownTelemetry, telemetryErr := telemetry.Setup(rootContext, loadedConfig.OTelGoInstrumentation)
	if telemetryErr != nil {
		_, _ = fmt.Fprintf(stderr, errorTelemetrySetup+lineEnding, telemetryErr)
	}
	if shutdownTelemetry != nil {
		defer func() {
			shutdownContext, cancel := context.WithTimeout(context.Background(), telemetryShutdownTimeout)
			defer cancel()
			if err := shutdownTelemetry(shutdownContext); err != nil {
				_, _ = fmt.Fprintf(stderr, errorTelemetryShutdown+lineEnding, err)
			}
		}()
	}

	ctx, span := telemetry.StartSpan(rootContext, instrumentationName, spanRun)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanRun, startedAt, exitCode != 0)
	}()

	token := strings.TrimSpace(os.Getenv(githubTokenEnv))
	if token == emptyString {
		token = strings.TrimSpace(os.Getenv(githubTokenFallbackEnv))
	}
	if token == emptyString && invocation.requiresGitHubToken() {
		_, _ = fmt.Fprintln(stderr, github.ErrMissingToken)
		return exitCodeUsage
	}

	return runCommand(ctx, invocation, loadedConfig, token, stdout, stderr)
}

func (invocation cliInvocation) requiresGitHubToken() bool {
	return invocation.command != commandGenerateReport
}

func validateCommandConfig(invocation cliInvocation, loadedConfig config.Config) error {
	if invocation.command == commandFlakyTestAnalysis {
		failingTestConfig := loadedConfig.FailingTestCheck.JSONReport
		if invocation.lookback != emptyString {
			failingTestConfig.Lookback = &invocation.lookback
			if _, err := failingTestConfig.EffectiveLookback(); err != nil {
				return fmt.Errorf(errorInvalidLookbackOverride, err)
			}
		}
	}
	if invocation.command == commandDependencyPRStatus && invocation.staleAge != emptyString {
		prStatusConfig := loadedConfig.PRStatus
		prStatusConfig.StaleAge = &invocation.staleAge
		if _, err := prStatusConfig.EffectiveStaleAge(); err != nil {
			return fmt.Errorf(errorInvalidStaleAgeOverride, err)
		}
	}
	return nil
}

func parseCLIInvocation(args []string) (cliInvocation, string, error) {
	filteredArgs, configPath, helpRequested, err := extractCLIOptions(args)
	if err != nil {
		return cliInvocation{}, emptyString, err
	}
	invocation := cliInvocation{configPath: configPath}
	if len(filteredArgs) == 0 {
		return invocation, helpTextRoot, nil
	}
	if filteredArgs[0] == commandHelp {
		helpRequested = true
		filteredArgs = filteredArgs[1:]
	}
	if len(filteredArgs) == 0 {
		return invocation, helpTextRoot, nil
	}

	commandName := filteredArgs[0]
	commandArgs := filteredArgs[1:]
	switch commandName {
	case commandVulnerabilityCheck, shortcutVulnerabilityCheck:
		invocation.command = commandVulnerabilityCheck
		if helpRequested {
			return invocation, helpTextVulnerabilityCheck, nil
		}
		if len(commandArgs) == 0 {
			return invocation, helpTextVulnerabilityCheck, fmt.Errorf(errorMissingVulnerabilityID)
		}
		if len(commandArgs) != 1 || strings.HasPrefix(commandArgs[0], dashPrefix) {
			return invocation, helpTextVulnerabilityCheck, fmt.Errorf(errorUnexpectedArguments, commandArgs)
		}
		invocation.vulnerabilityID = strings.TrimSpace(commandArgs[0])
		if invocation.vulnerabilityID == emptyString {
			return invocation, helpTextVulnerabilityCheck, fmt.Errorf(errorMissingVulnerabilityID)
		}
		return invocation, emptyString, nil
	case commandFlakyTestAnalysis, shortcutFlakyTestAnalysis:
		invocation.command = commandFlakyTestAnalysis
		if helpRequested {
			return invocation, helpTextFlakyTestAnalysis, nil
		}
		if err := parseFlakyTestAnalysisOptions(commandArgs, &invocation); err != nil {
			return invocation, helpTextFlakyTestAnalysis, err
		}
		return invocation, emptyString, nil
	case commandGenerateReport, shortcutGenerateReport:
		invocation.command = commandGenerateReport
		if helpRequested {
			return invocation, helpTextGenerateReport, nil
		}
		if len(commandArgs) != 0 {
			return invocation, helpTextGenerateReport, fmt.Errorf(errorUnexpectedArguments, commandArgs)
		}
		return invocation, emptyString, nil
	case commandDependencyPRStatus, shortcutDependencyPRStatus:
		if len(commandArgs) == 0 {
			return invocation, helpTextDependencyPRStatus, nil
		}
		if commandArgs[0] != commandUpdate && commandArgs[0] != shortcutUpdate {
			return invocation, helpTextDependencyPRStatus, fmt.Errorf(errorUnknownCommand, commandArgs[0])
		}
		invocation.command = commandDependencyPRStatus
		if helpRequested {
			return invocation, helpTextDependencyPRStatusUpdate, nil
		}
		if err := parseDependencyPRStatusOptions(commandArgs[1:], &invocation); err != nil {
			return invocation, helpTextDependencyPRStatusUpdate, err
		}
		return invocation, emptyString, nil
	default:
		return invocation, emptyString, fmt.Errorf(errorUnknownCommand, commandName)
	}
}

func extractCLIOptions(args []string) ([]string, string, bool, error) {
	configPath := defaultConfigPath
	filteredArgs := make([]string, 0, len(args))
	helpRequested := false
	for index := 0; index < len(args); index++ {
		argument := args[index]
		switch {
		case argument == flagHelpShort || argument == flagHelpLong:
			helpRequested = true
		case argument == flagConfigShort || argument == flagConfigLong:
			if index+1 >= len(args) || args[index+1] == emptyString || strings.HasPrefix(args[index+1], dashPrefix) {
				return nil, emptyString, false, fmt.Errorf(errorConfigFlagValue)
			}
			configPath = args[index+1]
			index++
		case strings.HasPrefix(argument, flagConfigLong+equalsSign):
			configPath = strings.TrimPrefix(argument, flagConfigLong+equalsSign)
			if configPath == emptyString {
				return nil, emptyString, false, fmt.Errorf(errorConfigFlagValue)
			}
		case strings.HasPrefix(argument, flagConfigShort+equalsSign):
			configPath = strings.TrimPrefix(argument, flagConfigShort+equalsSign)
			if configPath == emptyString {
				return nil, emptyString, false, fmt.Errorf(errorConfigFlagValue)
			}
		default:
			filteredArgs = append(filteredArgs, argument)
		}
	}
	return filteredArgs, configPath, helpRequested, nil
}

func parseFlakyTestAnalysisOptions(args []string, invocation *cliInvocation) error {
	flags := flag.NewFlagSet(commandFlakyTestAnalysis, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	lookback := flags.String(flagLookback, emptyString, helpLookback)
	overwrite := flags.Bool(flagOverwrite, false, helpOverwrite)
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() > 0 {
		return fmt.Errorf(errorUnexpectedArguments, flags.Args())
	}
	invocation.lookback = *lookback
	invocation.overwrite = *overwrite
	return nil
}

func parseDependencyPRStatusOptions(args []string, invocation *cliInvocation) error {
	flags := flag.NewFlagSet(commandDependencyPRStatus, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	parseWorkflowFailures := flags.Bool(flagParseWorkflowFailures, true, helpParseWorkflowFailures)
	flags.BoolVar(parseWorkflowFailures, flagNoParseWorkflowFailures, false, helpNoParseWorkflowFailures)
	staleAge := flags.String(flagStaleAge, emptyString, helpStaleAge)
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() > 0 {
		return fmt.Errorf(errorUnexpectedArguments, flags.Args())
	}
	invocation.parseWorkflowFailures = *parseWorkflowFailures
	invocation.staleAge = *staleAge
	return nil
}

func runCommand(ctx context.Context, invocation cliInvocation, loadedConfig config.Config, token string, stdout io.Writer, stderr io.Writer) (exitCode int) {
	startedAt := time.Now()
	enabledRepositories := loadedConfig.EnabledRepositories()
	logger := logging.New(stderr, commandLogPrefix(invocation.command), loadedConfig.EffectiveLogLevel())
	logger.Progressf("Starting %s for %d enabled repository/repositories.", invocation.command, len(enabledRepositories))
	defer func() {
		if exitCode == 0 {
			logger.Progressf("Completed %s successfully.", invocation.command)
			return
		}
		logger.Progressf("%s finished with exit code %d.", invocation.command, exitCode)
	}()
	ctx, span := telemetry.StartSpan(ctx, instrumentationName, spanRunChecks, attribute.Int(traceAttributeRepositoryCount, len(enabledRepositories)))
	defer func() {
		telemetry.FinishSpan(ctx, span, spanRunChecks, startedAt, exitCode != 0)
	}()
	if invocation.command == commandVulnerabilityCheck {
		if err := runVulnerabilityCheck(ctx, loadedConfig, invocation.vulnerabilityID, token, stdout, logger); err != nil {
			logger.Progressf("Command failed: %v", err)
			_, _ = fmt.Fprintf(stderr, errorVulnerabilityCheck+lineEnding, err)
			if errors.Is(err, github.ErrInvalidRepositoryURL) {
				return exitCodeUsage
			}
			return exitCodeError
		}
		return 0
	}

	for _, repository := range enabledRepositories {
		if err := runRepositoryCommand(ctx, repository, invocation, loadedConfig, token, stdout, logger); err != nil {
			logger.Progressf("Repository %q failed: %v", repository.URL, err)
			_, _ = fmt.Fprintf(stderr, errorRepositoryFailure+lineEnding, repository.URL, err)
			if errors.Is(err, github.ErrInvalidRepositoryURL) {
				return exitCodeUsage
			}
			exitCode = exitCodeError
		}
	}
	return exitCode
}

func runRepositoryCommand(ctx context.Context, repositoryConfig config.RepositoryConfig, invocation cliInvocation, loadedConfig config.Config, token string, stdout io.Writer, logger *logging.Logger) (resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, instrumentationName, spanRepository)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanRepository, startedAt, resultErr != nil)
	}()

	var githubClient *github.Client
	var repository github.Repository
	var err error
	if invocation.requiresGitHubToken() {
		githubClient, err = github.NewClient(ctx, strings.TrimSpace(repositoryConfig.URL), token, nil)
		if err == nil {
			repository = githubClient.TargetRepository()
		}
	} else {
		repository, err = github.ParseRepositoryURL(strings.TrimSpace(repositoryConfig.URL))
	}
	if err != nil {
		return err
	}
	logger.Progressf("Starting repository %s/%s.", repository.Owner, repository.Name)
	span.SetAttributes(
		attribute.String(traceAttributeRepositoryName, repository.Name),
		attribute.String(traceAttributeRepositoryOwner, repository.Owner),
	)

	switch invocation.command {
	case commandFlakyTestAnalysis:
		failingTestConfig := loadedConfig.FailingTestCheck.JSONReport
		if invocation.lookback != emptyString {
			failingTestConfig.Lookback = &invocation.lookback
		}
		resultErr = runFailingTestJSONReport(ctx, githubClient, failingTestConfig, invocation.overwrite, logger)
	case commandGenerateReport:
		resultErr = runFailingTestReadme(ctx, repository, stdout)
	case commandDependencyPRStatus:
		prStatusConfig := loadedConfig.PRStatus
		if invocation.staleAge != emptyString {
			prStatusConfig.StaleAge = &invocation.staleAge
		}
		resultErr = runPRStatus(ctx, githubClient, prStatusConfig, loadedConfig.FailingTestCheck.JSONReport, invocation.parseWorkflowFailures, stdout, logger)
	default:
		return fmt.Errorf(errorUnknownCommand, invocation.command)
	}
	if resultErr == nil {
		logger.Progressf("Finished repository %s/%s.", repository.Owner, repository.Name)
	}
	return resultErr
}

func commandLogPrefix(command string) string {
	switch command {
	case commandDependencyPRStatus:
		return logging.PrefixDependencyPRStatus
	case commandFlakyTestAnalysis:
		return logging.PrefixFlakyTestAnalysis
	case commandGenerateReport:
		return logging.PrefixGenerateReport
	default:
		return logging.PrefixVulnerability
	}
}

func formatRunDuration(duration time.Duration) string {
	totalSeconds := int64(duration.Round(time.Second) / time.Second)
	hours := totalSeconds / secondsPerHour
	minutes := totalSeconds % secondsPerHour / secondsPerMinute
	seconds := totalSeconds % secondsPerMinute
	parts := make([]string, 0, durationParts)
	if hours > 0 {
		parts = append(parts, formatDurationPart(hours, hoursUnit))
	}
	if minutes > 0 {
		parts = append(parts, formatDurationPart(minutes, minutesUnit))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, formatDurationPart(seconds, secondsUnit))
	}
	return strings.Join(parts, durationSeparator)
}

func formatDurationPart(value int64, unit string) string {
	if value == 1 {
		return fmt.Sprintf(singularDurationPart, value, unit)
	}
	return fmt.Sprintf(pluralDurationPart, value, unit, pluralSuffix)
}

func runVulnerabilityCheck(ctx context.Context, loadedConfig config.Config, vulnerabilityID string, token string, stdout io.Writer, logger *logging.Logger) error {
	enabledRepositories := loadedConfig.EnabledRepositories()
	repositories := make([]common.Repository, 0, len(enabledRepositories))
	for _, repositoryConfig := range enabledRepositories {
		repository, err := github.ParseRepositoryURL(strings.TrimSpace(repositoryConfig.URL))
		if err != nil {
			return fmt.Errorf("could not parse configured repository %q for Vulnerability Check: %w", repositoryConfig.URL, err)
		}
		repository.LocalPath = repositoryConfig.LocalPath
		repositories = append(repositories, repository)
	}
	analyzer, err := vulnerability.NewAnalyzer(token)
	if err != nil {
		return err
	}
	analyzer.SetLogger(logger)
	report, err := analyzer.Analyze(ctx, vulnerabilityID, repositories)
	if err != nil {
		return err
	}
	logger.Progressf(messageVulnerabilityReportStart)
	reportPaths, err := output.WriteVulnerabilityReportFiles(report)
	if err != nil {
		return err
	}
	for _, reportPath := range reportPaths {
		logger.Progressf(messageVulnerabilityReportPath, reportPath)
	}
	if err := output.WriteVulnerabilitySummary(report, stdout); err != nil {
		return err
	}
	logger.Progressf(messageVulnerabilityReportDone)
	for _, repository := range report.Repositories {
		for _, failure := range repository.Failures {
			logger.Progressf(messageVulnerabilityRepositoryFailure, repository.Repository.Owner, repository.Repository.Name, failure.Operation, failure.Message)
		}
	}
	if report.HasFailures() {
		return fmt.Errorf("Vulnerability Check completed with repository failures")
	}
	return nil
}

func runPRStatus(ctx context.Context, githubClient *github.Client, prStatusConfig config.PRStatusConfig, failingTestConfig config.FailingTestJSONReportConfig, parseWorkflowFailures bool, stdout io.Writer, logger *logging.Logger) error {
	lookback, err := failingTestConfig.EffectiveLookback()
	if err != nil {
		return err
	}
	staleAge, err := prStatusConfig.EffectiveStaleAge()
	if err != nil {
		return err
	}
	runner, err := prstatus.NewRunner(githubClient, prStatusConfig.EffectiveDependencyUpdateAuthors(), lookback, staleAge, parseWorkflowFailures)
	if err != nil {
		return err
	}
	runner.SetLogger(logger)
	return runner.Run(ctx, stdout)
}
