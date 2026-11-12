# Flaky Test Analysis

## Overview

Flaky test analysis runs through the KTLO script's `failingTestCheck` when enabled in the YAML config. Run the script
locally to analyze GitHub workflow logs and write a machine-readable report to the local results directory for the PR
status runner to consume. Do not add a scheduled GitHub Actions workflow at this stage. The first implementation
targets GitHub; the report layout and contract should also support GitLab.

## Retention and lookback

The analyzer's lookback window does not need to exceed the log retention period configured by GitHub or GitLab. The
analyzer should use only evidence available within that window. Missing history or incomplete analysis must be
represented as insufficient data, not as evidence that a failure is non-flaky. The window is set by
`failingTestCheck.lookback`, which defaults to `7d` and accepts positive Go duration values such as `30m` and `12h`,
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

## Report artifact

Write a versioned JSON file named `flaky-test-results.json` locally at
`results/github/<owner>/<repo>/flaky-test-results.json`. The PR status runner reads this local file; do not upload it
as a GitHub Actions artifact or commit it to the repository. Keep the JSON separate from the KTLO script's generated
Markdown results. Do not include raw job logs; include test observations and links to source workflow runs or jobs.
Include a `generatedAt` timestamp so the consumer can check freshness. The `failingTestCheck.maxArtifactAgeHours`
config defaults to 48 hours; `0` disables the freshness limit.
Normalize string-backed enum values to lowercase when decoding the report before validating them.

The report contains workflow runs, each with its own jobs, steps, and tests. A workflow contains its tests; when the
source step is known, nest a test under `step.tests`. Each level keeps its own outcome, so a failed step may be recorded
even when its job passes, or a failed job when its workflow passes. Every level retains branch, commit SHA, run URL,
and attempt metadata through its workflow parent; jobs also retain their job URL. Consumers derive pass/fail rates
from the nested records.
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
Classify a package-level failure as flaky when the same package has both failed and passed outcomes across retries or
workflow runs in the lookback window. Exclude package-level `ok (cached)` output from observations and pass/fail rates
because it does not represent a fresh test execution.
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
The report's `generatedAt` is checked against
`failingTestCheck.maxArtifactAgeHours`, which defaults to 48 hours; setting it to zero disables the freshness limit.

Classify a test as flaky if it fails and then passes on retry within the same workflow run, or if the same test has both
passing and failing outcomes across separate runs within the lookback window. Count each GitHub Actions rerun attempt
as a separate observation. Match observations using the full package and test/subtest path and retain each
observation's branch, outcome, run attempt, and source run URL. If a failure cannot be mapped to a test identifier,
report it as an unattributed failure and leave its test classification unknown. Non-test job failures must be retained
as separate workflow failure observations with their workflow/job/step identity and failure reason; do not discard them
when Go test output is absent. Apply the same flaky rule to these failures: a workflow/job/step failure is flaky when
it fails and then passes on retry, or when it has both failed and passed outcomes across workflow runs in the lookback
window. Match those failures by workflow, job, and failed step names; reason text does not change the failure identity.
If GitHub provides no failed step name, use the workflow and job names as the failure identity.
Preserve a concise error excerpt from the failed step's logs on every observation, including different excerpts for the
same workflow/job/step. Select error-like lines and nearby context, capped at 4 KiB per failed-step excerpt; do not
store the full failed-step log.

## PR status integration

The PR status runner reads the local report at `results/github/<owner>/<repo>/flaky-test-results.json` and matches
findings by provider, owner, repository name, PR number, and exact PR head SHA. Findings for a previous head SHA are
stale and must not be applied to the current PR revision. If the report is missing or older than the configured
freshness limit, the runner leaves the CI failure unclassified. An incomplete report may still classify a failure as
flaky when its available observations explicitly show both failed and passed outcomes for the same test or
workflow/job/step. Missing observations must not be treated as evidence that a failure is non-flaky; if the available
observations do not establish flakiness, leave the failure unclassified.

For a complete, fresh analysis, label a failed test, package, or workflow/job/step `Broken` when it has failures and no
passing observations in the lookback window. Show the number of passes and failures so the absence of passing evidence
is clear. A PR failure with no matching `main` observation is `Unknown` for insufficient baseline history, rather than
`Broken`; `Broken` requires at least one matching `main` observation. Apply the same rule to non-test workflow/job/step
failures: without a matching `main` observation, classify them as `Unknown` for insufficient baseline history.

Include observed pass/fail counts or rates in PR Status output only when explaining a test failure. Do not include those
rates when reporting a non-test workflow, job, or step failure. Calculate and report the counts or rates separately for
each failed test; do not aggregate them across multiple tests in a PR. When reporting a flaky non-test failure, include
its workflow/job/step identity, concise reason excerpt, and source run link. Include non-test failure details in PR Status
whether or not the available evidence confirms that the failure is flaky; leave the flaky classification unknown when
there is insufficient evidence. For each test failure, include separate pass/fail counts for `main` and for each
relevant PR branch commit, as well as the associated workflow, job, and step when those observations are available.
Never combine branch observations with `main` or combine results from different PR commits. Count rerun attempts
separately. Display each as counts and a percentage, such as `99/100 passed (99%)`. If the report is incomplete,
label rates as partial and clarify that they use only available observations. Round displayed percentages to the
nearest whole number.

Both PR Status and Failing Test Check must consume this report through one shared internal package. That package owns
the common report loading and decoding behavior so the two checks do not implement separate artifact readers.

The GitHub and GitLab integrations may use their own log APIs, but both must decode the same versioned report contract.
The provider's log retention bounds how much historical evidence the analyzer can use.
