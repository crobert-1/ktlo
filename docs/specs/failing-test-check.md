# Failing Test Check

## Overview

The check will document tests that are failing or flaky for each configured repository. It is separate from PR status,
but may consume the local flaky-test analysis report described in
[flaky-test-analysis.md](./flaky-test-analysis.md). Enable `failingTestCheck` in the YAML config and run KTLO locally
to start analysis. No scheduled GitHub Actions workflow is planned at this stage.

## Current implementation boundary

The YAML key is `failingTestCheck`; its `enabled` setting defaults to `false`. The analyzer is being reorganized to consume a provider-neutral
workflow client interface that supplies repository identity, workflow run attempts, jobs, and step logs. Provider
adapters translate their native API responses into shared analyzer types; the GitHub adapter is the initial
implementation. The analyzer collects workflow run attempts, job and step conclusions, and test results within the
configured lookback, then merges them into the local report. The standalone check passes `failingTestCheck.lookback`
to the analyzer; it defaults to `7d` and accepts positive duration values such as `30m`, `12h`, and `7d`. Collection or parsing errors mark the report incomplete. The analyzer does
not currently use `failingTestCheck.maxArtifactAgeHours` during collection.
Test log parsing uses the `TestLogParser` interface; Go-specific parsing lives in `go_test_log_parser.go`, and additional
language parsers can be supplied when constructing the analyzer.

PR Status also runs the analyzer by default to refresh the shared report before checking PRs. Setting
`prStatus.updateFlakyAnalysisReport` to `false` disables that refresh; `failingTestCheck.enabled` independently controls
whether the analyzer runs as a standalone check.

## Future requirements

- Enforce `failingTestCheck.maxArtifactAgeHours` when consuming flaky-analysis reports; a value of `0` disables the
  age limit.
- Allow multiple analyzer instances to run in parallel. This is a future requirement; the initial implementation does
  not need to account for concurrent analyzer execution.

The config may set `failingTestCheck.lookback` to bound historical evidence; it defaults to `7d` and accepts positive
Go duration values such as `30m` and `12h`, plus `d` for days. It may set
`failingTestCheck.maxArtifactAgeHours` to define report freshness; it defaults to 48 hours, and zero disables the
freshness limit. If the configured lookback extends beyond available provider log retention, analyze available logs
and mark the report incomplete rather than silently shortening the window. For a failed PR run, the analyzer parses Go test output and compares identified failures with
outcomes from every workflow run on `main` within the lookback window. It classifies a test as flaky when a failed
attempt passes on retry or when the same test has both passing and failing outcomes across runs. Intermittent failures
on `main` count as evidence. Test identifiers include the full package and test/subtest path. Failures without a
stable test identifier remain unattributed and unclassified. Failed jobs without identifiable Go test output are
recorded separately with workflow name, job name, failed step when available, failure point or reason, rerun attempt,
and source URL. Retain every failed step in a failed job. Retain successful job outcomes as pass evidence for comparison.
Only jobs with a `failure` conclusion count as failed observations. Record other inconclusive conclusions, such as
timed-out, cancelled, or neutral, with their original conclusion and exclude them from pass/fail rates and flaky
classification. Each GitHub Actions rerun
attempt counts as a separate observation, and detailed observations retain
source-run metadata in the local consolidated JSON report. A non-test workflow/job/step failure is classified as
flaky when it fails and then passes on retry or has both failed and passed outcomes across workflow runs in the
lookback window. Match failures by workflow, job, and failed step names, while retaining each observation's reason
text independently as a concise error excerpt from the failed step's logs, so different reasons for the same failure
point remain visible. Extract error-like lines with nearby context, capped at 4 KiB per failed-step excerpt. If no
failed step name is available, identify the failure by workflow and job names. Do not store the full failed-step log.
If some relevant logs cannot be downloaded or parsed, retain the observations that were collected and mark the whole
repository report incomplete. Label all rates from it as partial and based on available observations. PR Status may
still use an incomplete report to classify a failure as flaky when collected observations
explicitly show both failed and passed outcomes for the same test or workflow/job/step. Otherwise it leaves the failure
unclassified rather than concluding that it is non-flaky.

Each run merges newly collected workflow runs into the existing local report. Before writing, remove only stored runs
that fall outside the currently configured lookback window; retain in-window runs. Identify runs by provider,
repository, run ID, and run attempt so a previously stored run is not added or counted again. Do not clear the full
report on each analysis.

The local report contains workflow runs, with jobs and steps nested beneath each workflow. Tests are nested under
`step.tests` when the source step is known; otherwise place them at the deepest known parent. Record passing and failing
outcomes at workflow, job, step, and test levels with source run metadata so consumers can correlate test outcomes with
their CI context and derive pass/fail rates. Count
each workflow, job, and step from its own conclusion; a failed step can count as failed even when its job passes, and a
failed job can count as failed even when its workflow passes.
Each observation retains its branch and commit SHA. Keep `main` observations separate from PR branch observations and
keep PR branch observations separate for each commit SHA; use the current head SHA for current PR classification.
Successful observations retain only the identity needed for grouping, outcome, run attempt, and source run URL. Do not
store logs or reason text for successful observations, and retain successful jobs and steps even if another part of the
workflow run failed. Store a failed test's message in its `failureMessage` field; omit that field for passing and
skipped test observations. Record tests marked `SKIP` as a separate outcome and exclude them from pass/fail rates and flaky
classification. Record workflow steps concluded as `skipped` separately and exclude them from pass/fail rates and flaky
classification. Record other inconclusive workflow/job/step conclusions with their original provider conclusion and
exclude them from rates and flaky classification. Available inconclusive outcomes do not make the overall report
incomplete; incompleteness is reserved for relevant logs or history that cannot be downloaded or parsed. When Go output
contains only a package-level `ok`, count matching identified tests in that package as
passed without reconstructing all passing test names. An explicit `SKIP` outcome overrides the package-level pass for
that test. When output reports package-level `FAIL` without an individual test name, record a package-level test
failure using the package identity rather than attributing it to every test in the package. Classify a package-level
failure as flaky when the same package has both failed and passed outcomes across retries or workflow runs in the
lookback window. Exclude package-level `ok (cached)` output from observations and pass/fail rates because it does not
represent a fresh test execution.
