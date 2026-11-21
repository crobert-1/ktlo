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
	"fmt"
	"io"
	"time"

	"ktlo/internal/config"
	"ktlo/internal/failingtest"
	"ktlo/internal/github"
	"ktlo/internal/logging"
	"ktlo/internal/output"
	"ktlo/internal/telemetry"
	"ktlo/internal/testresults"
)

func runFailingTestJSONReport(ctx context.Context, githubClient *github.Client, failingTestConfig config.FailingTestJSONReportConfig, overwriteExistingReport bool, logger *logging.Logger) (resultErr error) {
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
	analyzer, err := failingtest.NewAnalyzerWithLookbackDurationAndOverwrite(workflowProvider, lookback, overwriteExistingReport)
	if err != nil {
		return err
	}
	analyzer.SetLogger(logger)
	return analyzer.Run(ctx)
}

func runFailingTestReadme(ctx context.Context, repository github.Repository, stdout io.Writer) error {
	report := testresults.NewReport(testresults.ProviderGitHub, repository)
	if err := report.Load(); err != nil {
		return fmt.Errorf("could not load flaky-test report for README generation: %w", err)
	}
	return output.WriteFailingTestReadme(ctx, report, stdout)
}
