# Flaky Test Analysis

## Overview

Run `ktlo flaky-test-analysis` (`ktlo f`) to analyze GitHub workflow logs and write a machine-readable report to the
local results directory for the PR status runner to consume. Do not add a scheduled GitHub Actions workflow at this
stage. The first implementation targets GitHub; the report layout and contract should also support GitLab.

## Retention and lookback

The analyzer's lookback window does not need to exceed the log retention period configured by GitHub or GitLab. The
analyzer should use only evidence available within that window. Missing history or incomplete analysis must be
represented as insufficient data, not as evidence that a failure is non-flaky. The window is set by
`failingTestCheck.jsonReport.lookback`, which defaults to `7d` and accepts positive Go duration values such as `30m` and `12h`,
plus `d` for days. The effective lookback must stay within the provider's
log retention period. If the configured window extends beyond available provider log retention, analyze the available
logs and mark the repository report incomplete rather than silently shortening the window. If some relevant logs cannot
be downloaded or parsed, write the observations that were collected and mark the repository report incomplete. Label
all rates from an incomplete report as partial and based on available observations.

For a failed PR run, parse the Go test output to identify failing tests. Compare those test identifiers with Go test
outcomes from every workflow run on `main` during the configured lookback period. Include both passed and failed main
workflow runs so intermittent failures on `main` can establish flakiness. Also record failed jobs that have no
identifiable Go test output. Capture the workflow name, job name, failed step when available, failure point or reason,
run attempt, and source URL so setup and other non-test failures remain traceable. Retain successful job outcomes as
pass evidence for comparison. Count only jobs with a `failure` conclusion as failed observations. Record other
inconclusive conclusions, such as timed-out, cancelled, or neutral, separately with their original conclusion; exclude
them from pass/fail rates and flaky classification. Retain every step with a `failure` conclusion in a failed job.
Collect PR-side observations only for open dependency-update PRs tracked by PR Status. Retain observations for each
commit on a tracked PR branch separately. Use only observations matching the current head SHA to classify the current
revision; older commit observations remain historical context.

Each analysis updates the existing local report incrementally. Load the current report, remove stored workflow-run
records whose run time falls outside the newly requested lookback window, retain records still inside that
window, and add newly collected runs. Do not clear the whole report before writing the new results. Give every stored
workflow run a stable provider run ID, its run attempt, and a run timestamp so the analyzer can identify the same run
across invocations and apply the lookback cutoff. Within a report, treat `(provider, repository, runId, runAttempt)` as
the unique run key. If that run key is already present, do not append it again or count its observations twice.
Update `generatedAt` after the merged report is written.

By default, reuse an existing report whose
`lookbackStart` covers the configured window and skip collection and log parsing for stored run attempts marked
`collectionComplete`. Retry incomplete or newly discovered attempts. If the report is missing `lookbackStart` or does
not cover the requested window, collect the full lookback to establish complete coverage. The `--overwrite` option
discards the existing report and rebuilds the configured lookback. Reject a mismatched
`schemaVersion` before collection unless `--overwrite` was supplied, and tell the user to use that option to replace
the old report. Do not automatically upgrade report schemas.

## Report artifact

Write a versioned JSON file named `workflow-test-analysis.json` locally at
`results/github/<owner>/<repo>/workflow-test-analysis.json`. The PR status runner reads this local file; do not upload it
as a GitHub Actions artifact or commit it to the repository. Keep the JSON separate from the KTLO script's generated
Markdown results. Do not include raw job logs; include test observations and links to source workflow runs or jobs.
Include a `generatedAt` timestamp and the requested `lookbackStart` so consumers can identify when the report was
generated and the period it covers. Each workflow attempt records `collectionComplete`; attempts with collection or
parsing errors are retried on the next incremental run. Older reports without this field are recollected once.

Run `ktlo generate-report` (`ktlo g`) to render the local report to
`results/github/<owner>/<repo>/FLAKY-AND-FAILING-TESTS.md`. The README summarizes failures observed on `main`, with separate test
sections for flaky, broken, and unknown classifications, and separate workflow/job/step sections for flaky, broken, and
unknown classifications. Each row includes
failure and observation counts, a failure rate, and the latest failure message with a link to its source run when
available. An incomplete report is labeled partial, and the notice clarifies that a test marked `Broken` may have
successful runs that were not gathered. Unknown classifications include their reason. The README
uses HTML tables so recent failure examples can render as multiline `<pre><code>` blocks within their cells. README-only
mode consumes an existing JSON report without making provider API calls.
In the README, omit a workflow-level failure observation when that run has a failed job, and omit a job-level failure
observation when that job has a failed step. Omit a package-level failure observation when the same workflow run
attempt and job contain an identified failed test in that package. Retain parent-level failures when no more specific
failed child identifies the cause. For nested test and subtest IDs, omit a failed parent test row when every observed
failure for it has a failed descendant in the same workflow attempt and job/step test scope, and the parent failure
message is empty or repeated in the descendant's message. Keep the parent row if any observed failure is not explained
this way. These README filters do not remove or rewrite observations in the JSON report.
Normalize string-backed enum values to lowercase when decoding the report before validating them.

The report contains workflow runs, each with its own jobs, steps, and tests. A workflow contains its tests; when the
source step is known, nest a test under `step.tests`. If a provider supplies only a job transcript and the parser
cannot identify the source step, store its test observations under `job.tests`. Each level keeps its own outcome, so a
failed step may be recorded even when its job passes, or a failed job when its workflow passes. Every level retains
branch, commit SHA, run URL, and attempt metadata through its workflow parent; jobs also retain their job URL.
Consumers derive pass/fail rates from the nested records.
Successful records retain only the identity needed for grouping, outcome, attempt, and source links. Do not store logs
or reason text for successful records. Keep successful jobs and steps even when another part of the workflow failed.
Record workflow steps concluded as `skipped` separately and exclude them from pass/fail rates and flaky classification.
Record other inconclusive workflow/job/step conclusions separately with their original provider conclusion and exclude
them from pass/fail rates and flaky classification.
These available inconclusive outcomes do not make the overall report incomplete; incompleteness is reserved for
relevant logs or history that cannot be downloaded or parsed.
Record Go tests marked `SKIP` as a separate `skipped` outcome and exclude them from pass/fail rates and flaky
classification. When output contains only a package-level `ok`, treat matching identified tests from that package as
passed for that run; do not require the output to list every passing test name. An explicit `SKIP` outcome overrides
that package-level pass for the skipped test. When output reports package-level `FAIL` without an individual test
name, record a package-level test failure using the package identity; do not attribute it to every test in the package.
Go's `--- FAIL` summaries for parent tests can repeat a failure from a nested test. When a failed parent has a failed
descendant and no independent failure output, omit that aggregate parent observation and retain the deepest identified
failure. Use that leaf test's output for its failure message and snippet; preserve a parent failure when its output
contains distinct failure evidence.
Classify a package-level failure as flaky when the same package has both failed and passed outcomes across retries or
workflow runs in the lookback window. Exclude package-level `ok (cached)` output from observations and pass/fail rates
because it does not represent a fresh test execution.
Choose among test log parsers by specificity, from structured format parsers to broader text parsers. When multiple
parsers match, use the highest-specificity parser; if multiple parsers match at that same specificity, record an
ambiguity error rather than choosing by registration order. Require result-bearing format evidence instead of
recognizing a log from a run or output marker alone. Keep failure excerpt extraction separate from test outcome parsing
so an excerpt match alone never creates a test failure.
For Go text logs, retain `Error Trace` assertion blocks as a failed test's failure message and excerpt, stopping before
subsequent `Logs:` output. For non-timeout panics, preserve the panic line and following goroutine stack. Use the bounded
raw test output only when no assertion block or panic stack is available; timeout failures retain the timeout summary.
Parse Molecule idempotence failures only when the log includes the matching scenario's idempotence group, its explicit
failed status, and the idempotence task summary. Record each listed failed task as a failed test observation whose
identity includes the Molecule scenario, role, and task name.
Each failed test observation may include a `failureMessage` containing its failure text. Omit `failureMessage` for
passing and skipped test observations so messages remain attached to the specific failed test and run attempt.

The initial report contract is:

```json
{
  "schemaVersion": 1,
  "provider": "github",
  "repository": {
    "owner": "owner",
    "name": "repository"
  },
  "generatedAt": "2026-11-05T20:00:00Z",
  "analysisStatus": "complete",
  "workflows": [
    {
      "name": "CI",
      "runId": "456",
      "runStartedAt": "2026-11-05T19:00:00Z",
      "branch": "main",
      "commitSha": "main-commit-sha-1",
      "runAttempt": 1,
      "outcome": "failed",
      "runUrl": "https://github.com/owner/repository/actions/runs/456",
      "jobs": [
        {
          "name": "unit-tests",
          "outcome": "failed",
          "jobUrl": "https://github.com/owner/repository/actions/runs/456/job/789",
          "steps": [
            {
              "name": "Run Go tests",
              "outcome": "failed",
              "reasonExcerpt": "--- FAIL: TestName (0.01s)",
              "tests": [
                {
                  "id": "package.TestName/subtest",
                  "kind": "test",
                  "outcome": "failed",
                  "failureMessage": "got unexpected result"
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "name": "CI",
      "runId": "457",
      "runStartedAt": "2026-11-05T19:30:00Z",
      "branch": "main",
      "commitSha": "main-commit-sha-2",
      "runAttempt": 1,
      "outcome": "passed",
      "runUrl": "https://github.com/owner/repository/actions/runs/457",
      "jobs": [
        {
          "name": "unit-tests",
          "outcome": "passed",
          "jobUrl": "https://github.com/owner/repository/actions/runs/457/job/790",
          "steps": [
            {
              "name": "Run Go tests",
              "outcome": "passed",
              "tests": [
                {
                  "id": "package.TestName/subtest",
                  "kind": "test",
                  "outcome": "passed"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

Each PR workflow includes `pullRequest.number` and the exact `pullRequest.headSha`. When the job is known but the source
step is not, place the test under `job.tests`; when neither is known, place it under `workflow.tests`. Keep
observations for each PR commit separate from
`main` observations and from other PR commits. Aggregate the `main` rate across its commits in the lookback window;
calculate PR branch rates separately for each commit SHA. PR Status uses only observations matching the current head
SHA to classify that revision; older commit observations remain historical context. A completed report with no matching flaky evidence means no
flaky outcome was identified in the available lookback window. Missing or incomplete reports must not be treated as a
completed analysis. When a scan is incomplete, retain collected observations and set `analysisStatus` to `incomplete`.
The report's `generatedAt` records when its contents were last updated. PR Status uses a valid local report regardless
of its age; run `ktlo flaky-test-analysis` to refresh the observations.

Classify a test as flaky if it fails and then passes on retry, if matching matrix jobs within a workflow run produce
both outcomes, or if the same test has both passing and failing outcomes across separate runs within the lookback
window. Match test identity by workflow name,
test kind, and full package and test/subtest path; ignore matrix job names and operating systems when combining test
outcomes. Keep each job and step as separate related CI history so the platform-specific evidence remains visible.
Count each GitHub Actions rerun attempt and each matching matrix job observation separately. Retain each
observation's branch, outcome, run attempt, and source run URL. If a failure cannot be mapped to a test identifier,
report it as an unattributed failure and leave its test classification unknown. Non-test job failures must be retained
as separate workflow failure observations with their workflow/job/step identity and failure reason; do not discard them
when Go test output is absent. Apply the same flaky rule to these failures: a workflow/job/step failure is flaky when
it fails and then passes on retry, or when it has both failed and passed outcomes across workflow runs in the lookback
window. Match those failures by workflow, job, and failed step names; reason text does not change the failure identity.
If GitHub provides no failed step name, use the workflow and job names as the failure identity.
For the `FLAKY-AND-FAILING-TESTS.md` workflow/job/step tables, merge matrix-expanded rows from the same
workflow job family when they have the same set of failed step names. Failure excerpt text does not affect the
grouping identity; retain it as supporting evidence in JSON and show the latest available excerpt in the README. Keep
groups separate when the failed-step sets differ. Show the base job name, number of affected matrix variants, and
combined failed-step names.

Calculate each grouped row's history by distinct, completely collected workflow run attempt. Exclude attempts with
`collectionComplete: false` from failure-group discovery, matrix-family detection, history counts, and latest-failure
evidence; their partial workflow, job, and step data is not reliable enough to classify. Among complete attempts, count
an attempt as failed if any variant has the grouped signature. Count it as passed only when all variants that
previously showed the grouped failure are present, every job in the family passed, and each grouped step passed. A
partial matrix or a different failure does not count as a pass. Keep the report-level incomplete status and notice
when collection errors occurred, even though incomplete attempts are excluded from these rows.
Preserve a concise error excerpt from the failed step's logs on every observation, including different excerpts for the
same workflow/job/step. When the log archive supplies only a combined job transcript, use step start and completion
timestamps to isolate each step before extracting its excerpt; do not reuse the full transcript as a failed step's
excerpt. Treat whole-second GitHub step timestamps as covering that full second so subsecond log lines at completion
are retained. If adjacent steps overlap at that precision, use exact GitHub group markers, post-job cleanup markers, and
log order to assign lines; do not discard ambiguous lines. If timestamps are unavailable or still do not identify a
step, leave the excerpt empty and retain the job link. Prefer a unified diff and its failure summary when present, strip log timestamps from that diff, and preserve
its opening file and hunk context within the 4 KiB per failed-step excerpt limit. Otherwise select meaningful failure
context rather than substring matches inside identifiers, and keep the relevant failure block within the limit; do not
store the full failed-step log. For Go module proxy download failures containing a stream `INTERNAL_ERROR`, include
the matching `go install` command and module download line when present, the proxy error, any immediately following
package failure summaries, and the following `make` failure line. If the matching download line is absent, include the
proxy error and failure summaries without unrelated downloads. Do not download logs for successful workflows;
successful workflow outcomes provide
inferred passes for known child test locations. For failed workflows, retain logs from successful matrix jobs so a
passing instance of a test can be compared with a failing instance. If a package fails while a Go test is still active and the log has no test result line, attribute the failure to
that test only when exactly one unclosed test has explicit assertion, fatal, or panic output; otherwise keep the
failure at package level.

## PR status integration

The PR status runner reads the local report at `results/github/<owner>/<repo>/workflow-test-analysis.json` and matches
findings by provider, owner, repository name, PR number, and exact PR head SHA. Findings for a previous head SHA are
stale and must not be applied to the current PR revision. If the report is missing, the runner leaves the CI failure
unclassified. A valid report is used regardless of its age. An incomplete report may classify a failure as `Flaky`
when its available observations explicitly show both failed and passed outcomes for the same test or workflow/job/step.
It may classify a failure as `Broken` when all observed outcomes at that location are failures and there is a matching
`main` observation. The partial report notice must clarify that an unobserved successful run may exist.

Label a failed test, package, or workflow/job/step `Broken` when it has failures and no passing observations in the
lookback window, including when analysis is incomplete. Show observed pass and failure counts, and label the report
partial so users know a successful run may not have been gathered. A PR failure with no matching `main` observation is
`Unknown` for insufficient baseline history, rather than `Broken`; `Broken` requires at least one matching `main`
observation. Apply the same rule to non-test workflow/job/step failures.

Include observed pass/fail counts or rates in PR Status output only when explaining a test failure. Do not include those
rates when reporting a non-test workflow, job, or step failure. Calculate and report the counts or rates separately for
each failed test; do not aggregate them across multiple tests in a PR. When reporting a flaky non-test failure, include
its workflow/job/step identity, concise reason excerpt, and source run link. Include non-test failure details in PR Status
whether or not the available evidence confirms that the failure is flaky; leave the flaky classification unknown when
there is insufficient evidence. For each test failure, include separate pass/fail counts for `main` and for each
relevant PR branch commit, as well as the associated workflow, job, and step when those observations are available.
Never combine branch observations with `main` or combine results from different PR commits. For a failed PR CI
context, resolve report history by its exact details URL or an unambiguous context name, and use successful
observations at that same workflow, job, step, or test location as passes. Each successful workflow also contributes
one inferred pass to every known job, step, and test location in that workflow, including locations without an
explicit result in that run. Keep explicit child failures alongside the inferred pass, and avoid duplicating an
explicit child pass. Mark inferred passes in PR Status output. Count rerun attempts separately. Display each as counts
and a percentage, such as
`99/100 passed (99%)`. In PR Status, show all distinct test and non-test failures for a PR in one table with failure,
classification, latest failure, `main` history, PR commit history, and related CI history columns; link the latest
failure to its source run when available. If the report is incomplete, label rates
as partial and clarify that they use only available observations. Round displayed percentages to the nearest whole
number.

Both PR Status and Failing Test Check must consume this report through one shared internal package. That package owns
the common report loading and decoding behavior so the two checks do not implement separate artifact readers.

The GitHub and GitLab integrations may use their own log APIs, but both must decode the same versioned report contract.
The provider's log retention bounds how much historical evidence the analyzer can use.

## TODO

- Add parsers for additional test log formats, ordered by specificity and backed by representative log fixtures.
- Add integration tests that save representative workflow result text and run the parser against known examples of
  different result kinds.
- Add a configuration option to treat a passed workflow as evidence that every job, step, and test passed, or instead
  inspect successful test results directly.
- Add a configuration option to analyze failed workflows only. For successful workflows, log that the workflow passed
  and skip job, step, and test parsing.
