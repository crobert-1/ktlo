# PR Status

## Overview

The goal of this check is to review open PRs in the target repository and inform the user of their status,
especially which PRs are ready to merge, and which PRs require attention.

## Current implementation boundary

The top-level `run` method constructs one GitHub GraphQL client with the repository URL and passes it to
`internal/prstatus.Runner`. The client is shared across functionality during the invocation. The runner currently
fetches open dependency bump PRs from the configured author list through a paginated GitHub GraphQL query and writes
them under `results/github/<owner>/<repository>/`. The runner also queries each PR's paginated GraphQL status check rollup
and passes the PRs and results to `internal/output.PrintPullRequestCIResults`. The formatter writes successful rollups
under `Passed`, failed or errored rollups under `Failed`, and all remaining PRs under `Others` in the repository-scoped
results file, with a separate Markdown table for each group. After the file is written successfully, the runner reports
that on stdout. It includes whether each PR is a draft in the analysis data. Review approval, merge conflict, and CI
approval-to-run are not yet queried. By default, PR Status refreshes the local flaky-analysis report before reading it.
Set `prStatus.updateFlakyAnalysisReport` to `false` to check against the existing report without updating it. This
setting does not control writes to the PR Status `RESULTS.md` file. The report format is described in
[flaky-test-analysis.md](flaky-test-analysis.md). Client initialization also performs a repository availability probe.

When PR Status refreshes the report, it uses the configured `failingTestCheck.lookback` and the same incremental
merge behavior as the standalone Failing Test Check. If both functionality settings request an analysis during one
invocation, run the analyzer once per repository and share the refreshed report.

## Functionality

The PR Status check will perform the following actions:

### Dependency update PRs

Treat pull requests from these author logins as dependency update PRs:

- `dependabot`
- `renovate`
- `srv-gh-o11y-gdi`

When matching GitHub author logins, ignore a trailing `[bot]` suffix.

#### Implementation

The script checks for open PRs from the listed dependency update authors in the target repository. The GitHub
GraphQL repository pull request connection does not filter by author, so the query fetches open PRs and the client
filters the results against this list. It includes both draft and ready-for-review PRs. The GitHub client's
CI/CD analysis method accepts this PR list from its caller and queries each PR; the planned classifications are:

The GitHub client exposes the analysis method as:

```go
func (c *Client) AnalyzePullRequests(ctx context.Context, pullRequests []PullRequest) ([]PullRequestCIResults, error)
func (c *Client) AnalyzePullRequest(ctx context.Context, pullRequest PullRequest) (PullRequestCIResults, error)
```

`AnalyzePullRequests` does not fetch PRs. For each supplied PR, it calls `AnalyzePullRequest`, continues after
individual query errors, and returns successful results alongside
any collected errors. The per-PR method reads the overall status rollup and paginates through modern check runs and
legacy status contexts. The caller obtains the PR list through `github.Client.OpenDependencyBumpPullRequests`;
`Runner.Run` calls that client method directly. Current output groups PRs according to their CI/CD rollup; the fuller
statuses described below remain future work. The intended check will:

- Check if they are ready to merge:
  - Check if their CI checks have passed.
  - Check if they have been approved by a reviewer.
  - If the PRs have passed CI and have been approved, the script will inform the user that the PRs are ready to merge.

- If the PR's CI is blocked from running:
  - The script will inform the user that the PRs need their CI to be approved to run.

- If the PR has a merge conflict, the script will inform the user the PR needs to be re-generated to resolve the merge conflict.

- If the PR has failed CI, use the shared functionality to check the failure logs for known flaky test patterns,
  and determine if the PR's CI failures are due to flaky tests. If the failures are due to flaky tests, the script will
  inform the user that the PR's CI failures are due to flaky tests. Include observed pass/fail counts or rates only
  per failed test when explaining test failures; do not include them for non-test workflow, job, or step failures or
  aggregate rates across different tests. When reporting a flaky non-test failure, include its workflow/job/step
  identity, concise reason excerpt, and source run link. Include non-test failure details even when the evidence is
  insufficient to classify the failure as flaky; leave its classification unknown. For each failed test, also show
  pass/fail counts for the associated workflow, job, and step when those observations are available. Show `main` rates
  aggregated across its commits separately from PR branch rates, with PR branch rates separated by commit SHA. Never
  combine branch observations with `main` or combine results from different PR commits. Use only the current PR head
  SHA to classify that revision, and count rerun attempts separately. Display counts with percentages, such as
  `99/100 passed (99%)`. Label each PR commit rate with its branch name and a short SHA linked to the commit page;
  retain full SHAs in JSON. Label the aggregated baseline with `main` and the configured lookback period, such as
  `main — last 7 days`. When the report is incomplete, label rates as partial and state that they use only available
  observations. Round displayed percentages to the nearest whole number.

#### Output

- The current output has three CI rollup groups:
  - `Passed`: PRs whose rollup state is `SUCCESS`.
  - `Failed`: PRs whose rollup state is `FAILURE` or `ERROR`.
  - `Others`: PRs with pending, expected, missing, unknown, or unavailable CI results.
- Later status classification may add groups for ready to merge, CI blocked, merge conflict, flaky CI failures, and
  non-flaky CI failures.

- PRs should be listed with following information (and format):
  - PR author in its own column; a clickable `PR Title (#PRNumber)` link; and the current CI/CD rollup state.

The generated document title is `PR Status Results`. Each section's table has the columns `Author`, `Pull Request`, and `CI/CD`. The Pull Request cell is a clickable
Markdown link whose visible text is `title (#number)`. If GitHub reports no rollup state, display `not reported`. A
`SUCCESS` rollup belongs under `Passed`; `FAILURE` and `ERROR` belong under `Failed`. Pending, expected, missing,
unknown, or unavailable results belong under `Others`. A successful CI/CD rollup alone does not mean the PR is ready
to merge because approval status is not yet queried.

Below the `Failed` table, include per-PR failure details. For test failures, list each failed test separately with its
observed counts and percentages and the associated workflow/job/step counts when available. For non-test failures,
include the workflow/job/step identity, concise reason excerpt, source run link, and flaky classification when known.
Label failures without enough evidence to determine flakiness as `Unknown`, explain why the classification is unknown,
and include relevant links to the failure and source workflow run. Use standard reasons: missing report, stale report,
incomplete log coverage, insufficient matching history, unidentifiable test or workflow failure, or PR head SHA mismatch.
For a complete, fresh analysis, label a failed test, package, or workflow/job/step with no passing observations in the
lookback window as `Broken`, and show its pass and failure counts, such as `0 passes, 3 failures`.
A PR test failure with no matching `main` observation is `Unknown` for insufficient baseline history, rather than
`Broken`; `Broken` requires at least one matching `main` observation. Apply the same rule to non-test workflow/job/step
failures.
For test-rate details, link the PR failure and representative `main` pass/fail runs, plus representative PR branch
pass/fail runs for each commit. Keep every contributing run link in the local JSON report rather than listing them all
in Markdown. Choose the most recent run for each outcome and commit as the representative pass or fail link. Link
directly to the job when a job URL is available; otherwise link to the workflow run.

The generated file is `results/github/<owner>/<repository>/RESULTS.md` in the script's working directory. Standard output reports successful completion;
errors that prevent file creation or the check from completing are written to standard error and cause a non-zero exit
code.
