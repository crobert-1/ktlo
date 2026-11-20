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
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"ktlo/internal/common"
	"ktlo/internal/config"
	"ktlo/internal/github"
	"ktlo/internal/output"
	"ktlo/internal/prstatus"
	"ktlo/internal/telemetry"
	"ktlo/internal/vulnerability"
)

const (
	errorRepositoryFailure                = "repository %q failed: %v"
	errorTelemetryShutdown                = "OpenTelemetry shutdown failed: %v"
	errorTelemetrySetup                   = "OpenTelemetry setup failed: %v"
	errorVulnerabilityCheck               = "Vulnerability Check failed: %v"
	messageSkipVulnerability              = "vulnerabilityCheck.id is not set; skipping Vulnerability Check"
	messageVulnerabilityReportDone        = "Markdown report written."
	messageVulnerabilityReportPath        = "Saved report: %s"
	messageVulnerabilityReportStart       = "Writing detailed Markdown reports to repository result files."
	messageVulnerabilityRepositoryFailure = "Vulnerability Check %s/%s %s failed: %s"
	exitCodeError                         = 1
	exitCodeUsage                         = 2
	githubTokenEnv                        = "GH_TOKEN"
	githubTokenFallbackEnv                = "GITHUB_TOKEN"
	instrumentationName                   = "ktlo/main"
	lineEnding                            = "\n"
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

func main() {
	// Skip the executable name so run receives only command-line arguments.
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) (exitCode int) {
	rootContext := context.Background()
	startedAt := time.Now()
	defer func() {
		_, _ = fmt.Fprintf(stderr, messageRunDuration+lineEnding, formatRunDuration(time.Since(startedAt)))
	}()

	loadedConfig, exitCode := loadConfig(args, stderr)
	if exitCode != 0 || loadedConfig == nil {
		return exitCode
	}
	if !loadedConfig.VulnerabilityCheck.HasID() {
		vulnerability.WriteProgress(stderr, messageSkipVulnerability)
	}
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

	if !loadedConfig.AnyCheckEnabled() {
		return 0
	}

	token := strings.TrimSpace(os.Getenv(githubTokenEnv))
	if token == emptyString {
		token = strings.TrimSpace(os.Getenv(githubTokenFallbackEnv))
	}
	if token == emptyString && loadedConfig.RequiresGitHubToken() {
		_, _ = fmt.Fprintln(stderr, github.ErrMissingToken)
		return exitCodeUsage
	}

	return runChecks(ctx, *loadedConfig, token, stdout, stderr)
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

func runChecks(ctx context.Context, config config.Config, token string, stdout io.Writer, stderr io.Writer) (exitCode int) {
	startedAt := time.Now()
	enabledRepositories := config.EnabledRepositories()
	ctx, span := telemetry.StartSpan(ctx, instrumentationName, spanRunChecks, attribute.Int(traceAttributeRepositoryCount, len(enabledRepositories)))
	defer func() {
		telemetry.FinishSpan(ctx, span, spanRunChecks, startedAt, exitCode != 0)
	}()
	if config.VulnerabilityCheck.HasID() {
		if err := runVulnerabilityCheck(ctx, config, token, stdout, stderr); err != nil {
			vulnerability.WriteProgress(stderr, errorVulnerabilityCheck, err)
			exitCode = exitCodeError
		}
	}

	for _, repository := range enabledRepositories {
		repositoryExitCode := runRepositoryChecks(ctx, repository, config, token, stdout, stderr)
		if repositoryExitCode == exitCodeUsage {
			return exitCodeUsage
		}
		if repositoryExitCode != 0 {
			exitCode = exitCodeError
		}
	}
	return exitCode
}

func runVulnerabilityCheck(ctx context.Context, loadedConfig config.Config, token string, stdout io.Writer, stderr io.Writer) error {
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
	analyzer.SetProgressWriter(stderr)
	report, err := analyzer.Analyze(ctx, loadedConfig.VulnerabilityCheck.ID, repositories)
	if err != nil {
		return err
	}
	vulnerability.WriteProgress(stderr, messageVulnerabilityReportStart)
	reportPaths, err := output.WriteVulnerabilityReportFiles(report)
	if err != nil {
		return err
	}
	for _, reportPath := range reportPaths {
		vulnerability.WriteProgress(stderr, messageVulnerabilityReportPath, reportPath)
	}
	if err := output.WriteVulnerabilitySummary(report, stdout); err != nil {
		return err
	}
	vulnerability.WriteProgress(stderr, messageVulnerabilityReportDone)
	for _, repository := range report.Repositories {
		for _, failure := range repository.Failures {
			vulnerability.WriteProgress(stderr, messageVulnerabilityRepositoryFailure, repository.Repository.Owner, repository.Repository.Name, failure.Operation, failure.Message)
		}
	}
	if report.HasFailures() {
		return fmt.Errorf("Vulnerability Check completed with repository failures")
	}
	return nil
}

func runRepositoryChecks(ctx context.Context, repositoryConfig config.RepositoryConfig, config config.Config, token string, stdout io.Writer, stderr io.Writer) (exitCode int) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, instrumentationName, spanRepository)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanRepository, startedAt, exitCode != 0)
	}()

	repositoryURL := repositoryConfig.URL
	var githubClient *github.Client
	var repository github.Repository
	var err error
	if config.RequiresGitHubToken() {
		githubClient, err = github.NewClient(ctx, strings.TrimSpace(repositoryURL), token, nil)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, errorRepositoryFailure+lineEnding, repositoryURL, err)
			if errors.Is(err, github.ErrInvalidRepositoryURL) {
				return exitCodeUsage
			}
			return exitCodeError
		}
		repository = githubClient.TargetRepository()
	} else {
		repository, err = github.ParseRepositoryURL(strings.TrimSpace(repositoryURL))
		if err != nil {
			_, _ = fmt.Fprintf(stderr, errorRepositoryFailure+lineEnding, repositoryURL, err)
			return exitCodeUsage
		}
	}
	span.SetAttributes(
		attribute.String(traceAttributeRepositoryName, repository.Name),
		attribute.String(traceAttributeRepositoryOwner, repository.Owner),
	)
	if config.FailingTestCheck.JSONReport.IsEnabled() {
		if err := runFailingTestJSONReport(ctx, githubClient, config.FailingTestCheck.JSONReport); err != nil {
			_, _ = fmt.Fprintf(stderr, errorRepositoryFailure+lineEnding, repositoryURL, err)
			exitCode = exitCodeError
		}
	}
	if config.FailingTestCheck.Readme.Enabled {
		if err := runFailingTestReadme(ctx, repository, stdout); err != nil {
			_, _ = fmt.Fprintf(stderr, errorRepositoryFailure+lineEnding, repositoryURL, err)
			exitCode = exitCodeError
		}
	}
	if config.PRStatus.IsEnabled() {
		if err := runPRStatus(ctx, githubClient, config.PRStatus, config.FailingTestCheck, stdout); err != nil {
			_, _ = fmt.Fprintf(stderr, errorRepositoryFailure+lineEnding, repositoryURL, err)
			exitCode = exitCodeError
		}
	}
	return exitCode
}

func runPRStatus(ctx context.Context, githubClient *github.Client, prStatusConfig config.PRStatusConfig, failingTestCheck config.FailingTestCheckConfig, stdout io.Writer) error {
	lookback, err := failingTestCheck.JSONReport.EffectiveLookback()
	if err != nil {
		return err
	}
	staleAge, err := prStatusConfig.EffectiveStaleAge()
	if err != nil {
		return err
	}
	runner, err := prstatus.NewRunner(githubClient, prStatusConfig.EffectiveDependencyUpdateAuthors(), lookback, staleAge, prStatusConfig.ParseWorkflowFailures)
	if err != nil {
		return err
	}
	return runner.Run(ctx, stdout)
}
