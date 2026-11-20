# Failing Test Check

## Overview

Failing Test Check has two independent outputs: a machine-readable JSON report and a user-facing Markdown README.
Configure them under `failingTestCheck.jsonReport` and `failingTestCheck.readme`. Both are disabled by default and may
be enabled together. No scheduled GitHub Actions workflow is planned at this stage.

## Current implementation boundary

The JSON report mode's `enabled` setting defaults to `false`. The analyzer consumes a provider-neutral
workflow interface that supplies repository identity, workflow run attempts, jobs, and step logs. Providers use the
shared workflow records in `internal/common`; the GitHub adapter maps its native API records into them. The analyzer
collects workflow run attempts, job and step conclusions, and test results within the configured lookback, then merges
them into the local report. It passes `failingTestCheck.jsonReport.lookback` to the analyzer; it defaults to `7d` and
accepts positive duration values such as `30m`, `12h`, and `7d`. Collection or parsing errors mark the report
incomplete. With `overwriteExistingReport: false`, a report that covers the requested window is reused and attempts
marked complete are skipped; incomplete attempts are retried. Reports without coverage metadata or that do not cover
the requested window are collected again. With `overwriteExistingReport: true`, the existing report is discarded and
the full lookback is collected.
Test log parsing uses `TestLogParser` implementations selected by descending specificity. The structured Go test JSON
parser runs before the human-readable Go test parser; if multiple matching parsers share the highest specificity,
parsing reports an ambiguity instead of relying on registration order. Go test outcome parsing is separate from
failure excerpt extraction, which uses dedicated extractors for Go test output, diffs, Go module proxy failures, and
generic step errors. Parsers should require result-bearing format evidence rather than claiming logs from a run or
output marker alone. Additional language parsers can be supplied when constructing the analyzer.

The README mode reads the local JSON report and writes `results/github/<owner>/<repo>/FLAKY-AND-FAILING-TESTS.md`. It summarizes main
branch workflow, job, step, and test failures; test rows are separated into flaky, broken, and unknown classifications.
Rows include failure and observation counts, failure rates, and a recent failure example linked to its source run. If
both modes are enabled, JSON collection runs before README rendering. README-only mode does not call GitHub APIs and
requires an existing JSON report. De-duplicate parent CI failure rows when the same run has a failed child: omit a
workflow failure when a job failure identifies it, and omit a job failure when a step failure identifies it. Retain
parent failures that have no more specific failed child. PR Status independently reads the local report but does not
run analysis or update it.

## Future requirements

- Allow multiple analyzer instances to run in parallel. This is a future requirement; the initial implementation does
  not need to account for concurrent analyzer execution.

The config may set `failingTestCheck.jsonReport.lookback` to bound historical evidence; it defaults to `7d` and accepts
positive Go duration values such as `30m` and `12h`, plus `d` for days. If the configured lookback extends beyond
available provider log retention, analyze available logs and mark the report incomplete rather than silently shortening
the window. For a failed PR run, the analyzer parses Go test output and compares identified failures with
outcomes from every workflow run on `main` within the lookback window. It classifies a test as flaky when a failed
attempt passes on retry, when matching matrix jobs within a workflow run produce both outcomes, or when the same test
has both passing and failing outcomes across runs. Intermittent failures
on `main` count as evidence. Match test identity by workflow name, test kind, and full package and test/subtest path;
ignore matrix job names and operating systems when combining test outcomes. Keep each job and step as separate related
CI history so platform-specific evidence remains visible. Count each matching matrix job observation separately.
Test identifiers include the full package and test/subtest path. Failures without a
stable test identifier remain unattributed and unclassified. If a package fails while a Go test is still active and the
log has no test result line, attribute the failure to that test only when exactly one unclosed test has explicit
assertion, fatal, or panic output; otherwise keep the failure at package level. Failed jobs without identifiable Go test output are
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
point remain visible. When the log archive supplies only a combined job transcript, use step start and completion
timestamps to isolate each step before extracting its excerpt; do not reuse the full transcript as a failed step's
excerpt. Treat whole-second GitHub step timestamps as covering that full second so subsecond log lines at completion
are retained. If adjacent steps overlap at that precision, use exact GitHub group markers, post-job cleanup markers, and
log order to assign lines; do not discard ambiguous lines. If timestamps are unavailable or still do not identify a
step, leave the excerpt empty and retain the job link. Prefer a unified diff and its failure summary when present, strip log timestamps from that diff, and preserve
its opening file and hunk context within the 4 KiB excerpt limit. Otherwise select meaningful failure context rather
than substring matches inside identifiers, and keep the relevant failure block when enforcing the excerpt limit. If no
failed step name is available, identify the failure by workflow and job names. Do not store the full failed-step log.
In the generated README, merge matrix-expanded rows from the same workflow job family when they have the same set of failed step names. Failure excerpt text does not affect the grouping identity; retain it as supporting evidence in JSON and show the latest available excerpt in the README. Keep groups separate when the failed-step sets differ. Show the base job name, number of affected matrix variants, and combined failed-step names. Calculate grouped history by distinct workflow run attempt: an attempt fails if any variant has the grouped signature. Count a pass only when collection is complete, all variants that previously showed the grouped failure are present, every job in the family passed, and each grouped step passed. A partial matrix or a different failure does not count as a pass.
If some relevant logs cannot be downloaded or parsed, retain the observations that were collected and mark the whole
repository report incomplete. Label all rates from it as partial and based on available observations. PR Status may
still use an incomplete report to classify a failure as flaky when collected observations
explicitly show both failed and passed outcomes for the same test or workflow/job/step. Otherwise it leaves the failure
unclassified rather than concluding that it is non-flaky.

Each run merges newly collected workflow runs into the existing local report. Before writing, remove only stored runs
that fall outside the currently configured lookback window; retain in-window runs. Identify runs by provider,
repository, run ID, and run attempt so a previously stored run is not added or counted again. Do not clear the full
report on each analysis. Persist the requested window start as `lookbackStart` and a `collectionComplete` flag on each
attempt. A report may be reused incrementally only when `lookbackStart` is present and no later than the current cutoff.
Skip complete attempts before fetching attempt details, jobs, or logs; retry incomplete attempts and collect newly
discovered attempts. Older reports without the attempt flag are recollected once. For unknown attempts, use the
run-list data for the latest attempt when available, and request older attempt details only when they are not already
stored. Run-attempt detail requests and independent attempt parsing use small bounded worker pools.

Do not download logs for successful workflows. The successful workflow outcome supplies inferred pass observations for
known jobs, steps, and tests. For failed workflows, retain logs from successful matrix jobs as well as failed jobs so
test outcomes can be compared across matrix variants. Download one attempt archive per such workflow attempt rather
than one log archive per job. Record separate durations for archive download and ZIP unpacking, including archive and
uncompressed byte counts, so trace data distinguishes transfer time from local processing time.

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
