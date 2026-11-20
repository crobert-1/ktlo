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
	"time"

	"ktlo/internal/config"
	"ktlo/internal/failingtest"
	"ktlo/internal/github"
	"ktlo/internal/output"
	"ktlo/internal/telemetry"
	"ktlo/internal/testresults"
)

const commandName = "ktlo"
const emptyString = ""
const errorMissingConfig = "required flag --config is missing"
const errorUnexpectedArgs = "unexpected positional arguments: %v"
const flagConfig = "config"
const helpConfig = "Path to the KTLO YAML config file"
const usageText = "Usage: ktlo --config <config.yaml>\n"

func loadConfig(args []string, stderr io.Writer) (*config.Config, int) {
	flags := flag.NewFlagSet(commandName, flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		_, _ = fmt.Fprint(stderr, usageText)
	}
	configPath := flags.String(flagConfig, emptyString, helpConfig)
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, 0
		}
		return nil, exitCodeUsage
	}
	if len(flags.Args()) > 0 {
		_, _ = fmt.Fprintf(stderr, errorUnexpectedArgs+lineEnding, flags.Args())
		return nil, exitCodeUsage
	}
	if *configPath == emptyString {
		_, _ = fmt.Fprintln(stderr, errorMissingConfig)
		_, _ = fmt.Fprint(stderr, usageText)
		return nil, exitCodeUsage
	}

	loadedConfig, err := config.Load(*configPath)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return nil, exitCodeUsage
	}
	return &loadedConfig, 0
}

func runFailingTestJSONReport(ctx context.Context, githubClient *github.Client, failingTestConfig config.FailingTestJSONReportConfig) (resultErr error) {
	startedAt := time.Now()
	ctx, span := telemetry.StartSpan(ctx, instrumentationName, spanFailingTestCheck)
	defer func() {
		telemetry.FinishSpan(ctx, span, spanFailingTestCheck, startedAt, resultErr != nil)
	}()

	workflowProvider, err := failingtest.NewGitHubWorkflowProvider(githubClient)
	if err != nil {
		return err
	}
	lookback, err := failingTestConfig.EffectiveLookback()
	if err != nil {
		return err
	}
	analyzer, err := failingtest.NewAnalyzerWithLookbackDurationAndOverwrite(workflowProvider, lookback, failingTestConfig.OverwriteExistingReport)
	if err != nil {
		return err
	}
	return analyzer.Run(ctx)
}

func runFailingTestReadme(ctx context.Context, repository github.Repository, stdout io.Writer) error {
	report := testresults.NewReport(testresults.ProviderGitHub, repository)
	if err := report.Load(); err != nil {
		return fmt.Errorf("could not load flaky-test report for README generation: %w", err)
	}
	return output.WriteFailingTestReadme(ctx, report, stdout)
}
