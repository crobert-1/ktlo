# PR Status

## Overview

The goal of this check is to review open PRs in the target repository and inform the user of their status,
especially which PRs are ready to merge, and which PRs require attention.

## Current implementation boundary

The top-level `run` method constructs one GitHub GraphQL client with the repository URL and passes it to
`internal/prstatus.Runner`. The client is shared across functionality during the invocation. The runner currently
fetches open dependency bump PRs from the configured author list through a paginated GitHub GraphQL query and writes
them under `results/github/<owner>/<repository>/`. The runner also queries each PR's paginated GraphQL status check rollup
and passes the PRs and results to `internal/output.PrintPullRequestCIResults`. The formatter writes successful rollups,
including PRs whose contexts are all passed or skipped, under `Passed`; missing results and PRs with in-progress checks
under `Pending`; failed or errored rollups without pending checks under `Failed`; and all remaining PRs under `Others` in
the repository-scoped results file, with a separate Markdown table for each group. After the file is written successfully, the runner reports
that on stdout. It includes whether each PR is a draft in the analysis data. Review approval, merge conflict, and CI
approval-to-run are not yet queried. PR Status reads the existing local flaky-analysis report for failure details when
available; it does not update that report. When standalone Failing Test Check is enabled in the same invocation, the
analyzer runs once before PR Status and PR Status consumes the updated report. Optionally, `prStatus.parseWorkflowFailures`
parses only failed GitHub Actions runs linked from failed CI contexts and merges those observations into an in-memory
report for matching. This targeted parsing uses the normal workflow/job/step/test builder and configured log parsers,
respects `failingTestCheck.jsonReport.lookback`, and never writes `workflow-test-analysis.json`. The report format is
described in [flaky-test-analysis.md](flaky-test-analysis.md). Client initialization also performs a repository
availability probe.

If the report is missing, stale, or incomplete, PR Status reports the corresponding unknown or partial analysis state;
targeted parsing is performed only when `prStatus.parseWorkflowFailures` is enabled and only for linked failed workflow
contexts. Only `failingTestCheck.jsonReport.enabled` runs the repository-wide lookback and writes or updates
`workflow-test-analysis.json`.

`prStatus.parseWorkflowFailures` defaults to `false`. When enabled, PR Status resolves workflow run IDs from failed
GitHub Actions check nodes through GraphQL and batches their associated check suites, workflow runs, jobs, and steps.
The selected run attempt is filtered by the shared JSON report lookback and PR head before logs are downloaded and
parsed. GitHub Actions runs with incomplete GraphQL metadata fall back to the REST run and jobs endpoints. Non-Actions
checks and failed contexts without a resolvable run retain the existing report-matching behavior. Parse/API errors are
reported in that PR's details and make the in-memory analysis partial; they do not stop output for other PRs. A
successful parsed workflow can add passing observations for its jobs, steps, and tests through the shared report
assessment rules. During one invocation, run metadata and parsed workflow attempts are cached across PRs. A complete
local report that already contains the exact run attempt is reused without downloading and parsing its logs again;
incomplete reports are reparsed for the linked failed run.

## Functionality

The PR Status check will perform the following actions:

### Dependency update PRs

Treat pull requests from these author logins as dependency update PRs:

- `dependabot`
- `renovate`
- `srv-gh-o11y-gdi`

When matching GitHub author logins, ignore a trailing `[bot]` suffix.

#### Implementation

The script checks for open PRs from the author logins in `prStatus.dependencyUpdateAuthors` in the target repository.
This list defaults to `dependabot`, `renovate`, and `srv-gh-o11y-gdi`. The GitHub
GraphQL repository pull request connection does not filter by author, so the query fetches open PRs and the client
filters the results against this list. It includes both draft and ready-for-review PRs, including each PR's current
head SHA for exact-revision report matching. The GitHub client's
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

- The current output has four CI status groups and a separate `Stale` table:
  - `Passed`: PRs whose reported CI contexts all succeeded or were skipped.
  - `Pending`: PRs with a missing or `not reported` CI result, or with at least one in-progress check. Pending takes
    precedence over all other CI groups, including `Failed`, when other checks have already failed.
  - `Failed`: PRs with at least one CI context in `FAILURE` or `ERROR` and no pending checks.
  - `Others`: PRs with cancelled, expected, unknown, or unavailable CI results that are not pending.
  - `Stale`: PRs whose latest commit is older than the configured `prStatus.staleAge`; this group has only author and
    pull-request columns, and stale PRs are omitted from the other tables and failure details.
- Derive the displayed CI status from the individual contexts rather than relying only on GitHub's aggregate rollup.
  Treat `SKIPPED` as a successful outcome. Treat `CANCELLED` as unresolved, not as a failure or pass. Add the
  non-zero passed, failed, skipped, and cancelled context counts to the CI/CD cell for every PR. Count `ERROR` contexts
  as failed. Keep the provider rollup only as a fallback when no individual contexts are available; omit counts when
  there are no contexts to count. If a PR has cancelled contexts but no failed contexts, show it under `Others` and
  identify the cancelled checks after the tables.
- Later status classification may add groups for ready to merge, CI blocked, merge conflict, flaky CI failures, and
  non-flaky CI failures.

- PRs should be listed with following information (and format):
  - PR author in its own column; a clickable `PR Title (#PRNumber)` link; and the current CI/CD rollup state.

The generated document title is `Dependency Bump PR Status`. Sort PRs by PR number in descending order within each status table,
and preserve that order in corresponding detail sections. Omit a section entirely when its table has no entries,
including its heading and header row. `Stale` has columns `Author` and `Pull Request`; other tables have columns
`Author`, `Pull Request`, and `CI/CD`. The Pull Request cell is a clickable
Markdown link whose visible text is `title (#number)`. Link the CI/CD status to the in-document details section for
that PR when one is generated: every PR with failed checks, including one assigned to `Pending`, and PRs in `Others` with cancelled checks. Link only the one-word
CI result summary; keep any outcome counts outside the link. Leave the status as plain text when the document has no
separate details section for that PR. If GitHub reports no rollup state, display `not reported`. The displayed CI/CD
state is derived from individual contexts. Any `FAILURE` or `ERROR` context places a PR under `Failed` unless it also
has a pending check, in which case it belongs under `Pending`. A PR is `Passed` when every reported context is
`SUCCESS` or `SKIPPED`. In every status group, append the
non-zero context totals for passed, failed, skipped, and cancelled outcomes to the CI/CD cell, such as
`FAILURE — 3 passed, 2 failed`. Omit zero totals. If no individual contexts are available, use the provider rollup as
a fallback and omit counts. A PR with no failed context but at least one `CANCELLED`, expected, unknown, or
unavailable context belongs under `Others` unless a check is pending. A missing result or a `PENDING`, `EXPECTED`,
`IN_PROGRESS`, `QUEUED`, `REQUESTED`, or `WAITING` check belongs under `Pending`. When the provider aggregate rollup
conflicts with complete context results, use the context-derived state, except a pending check takes precedence over
failures for table placement. A successful CI status alone does not mean the PR is ready to merge because
approval status is not yet queried.

The `prStatus.staleAge` setting defaults to `7d` and accepts positive durations such as `60m`, `12h`, or `7d`. Place
an open PR in `Stale` when its latest commit's `committedDate` is older than the current time minus this duration. Do
not fetch CI results or show failure details for stale PRs. If the latest commit timestamp is unavailable, classify
the PR using its CI results instead.

After the five tables and failure details, add a `CI checks not completed` section for PRs with cancelled contexts
and no failed contexts. Link cancelled checks when a details URL is available.

After the five tables, add one failure-details subsection for each PR with failed checks, preserving the order in its
status table. This includes PRs assigned to `Pending` because another check is still running. Head each subsection with a
link using the PR title and number, followed by one table containing all
distinct test and non-test failures, with one row per failure and columns for `Failure`, `Classification`, `Latest
failure`, `Main history`, `PR commit history`, and `Related CI history`. Show test identities in the `Failure` column,
and use the most-specific workflow, job, or step identity for non-test failures. Include an `Unknown` reason in the
classification cell unless the same report-wide issue is already explained once above the table. Render summary and
failure-detail tables as HTML tables. Format latest failure output as a full-width `<pre><code>` block inside its cell,
preserving newlines and escaping HTML-sensitive failure text. Keep the source run link outside the code block. Use
`—` for information that is unavailable or does not apply.

In the `Main history` cell, include pass/failure counts and percentages for `main` and the configured lookback period.
In the `PR commit history` cell, include counts and percentages for each relevant PR commit. In the `Related CI
history` cell, include associated workflow, job, and step counts when available. Leave history cells as `—` for
non-test failures. Aggregate `main` observations across its commits, but keep PR observations separated by commit SHA.
Never combine branch observations with `main` or combine different PR commits. Count rerun attempts separately.
Display counts with percentages, such as `99/100 passed (99%)`. Label each PR commit with its branch name and a
short SHA linked to the commit page; retain full SHAs in JSON. When observations come from an incomplete report, mark
the counts and rates as partial and explain that they use only available observations. Round displayed percentages
to the nearest whole number. Include inferred successful-workflow passes in the history counts and identify them in
the cell. For non-test failures whose classification used inferred passes, note that in the classification cell.

Suppress a parent workflow, job, or step failure when a more specific failure row already explains it, while
retaining independent failures. If a package failure has no identifiable failed test, show one package-level fallback
row; omit it when individual test failures already explain the package failure. Add unmatched failed CI contexts as
`Unknown` rows with their available name, description or summary, and details link. If a failed rollup has no
identifiable contexts or report findings, add one fallback `Unknown` row for the rollup.

For non-test workflow, job, or step failures, show the most specific failed context, its identity, a concise reason
excerpt, a source run link, and its classification when known in the corresponding table row. Do not show pass/failure
counts or rates for these failures. A package-level fallback row shows its package identity, classification and
explanation, and source link, but no counts or rates because package counts can aggregate multiple tests.

Label any failure without enough evidence to determine its classification as `Unknown`, and explain the specific
reason and what analysis could not be completed. Link to the failure and relevant source workflow run. Standard
reasons include missing report, stale report, malformed report, repository mismatch, report read error, incomplete
log coverage, insufficient matching history, unidentifiable test or workflow failure, and PR head SHA mismatch. When
a report-level issue affects a PR, show one PR-level `Unknown` notice naming the cause and its impact; do not repeat
that same report issue in each failure block. For a stale report, include its generation time and the configured age
limit, and state that its findings were not used. Use `failingTestCheck.maxArtifactAgeHours` for this age limit even
when analysis is disabled; zero disables the limit. Malformed, repository-mismatched, or unreadable reports do not
prevent the CI results document from being written. For each failed GitHub CI context, try to resolve its details URL
to a report workflow or job first. If no URL matches, use an exact context name only when it identifies one
most-specific workflow, job, step, or test location in the report. Assess a matched location even when the report
contains only successful outcomes: count successful outcomes for that same location as passes and the current failed
context as a failure at the exact PR head SHA. For test history, match by workflow name, test kind, and full test
identifier across all matrix jobs, regardless of operating system. Report the workflow-level test counts together,
and retain per-job and per-step histories in related CI details. A successful workflow also contributes one inferred
pass for every known job, step, or test location in that workflow, even when that child has no explicit result. Keep an explicit
child failure alongside the inferred pass, but do not double-count an explicit child pass. If the context cannot be
matched unambiguously, show an `Unknown` failure block using the context's name,
available description or summary, and details link. For an incomplete but usable report, show available evidence as
partial. It may classify a failure as `Flaky` only when available observations show both passes and failures for that
same location; otherwise classify it as `Unknown` and explain the missing coverage.

For a complete, fresh analysis, classify a failed test, package-level fallback, or workflow/job/step as `Broken` when
it has failures and no passing observations in the lookback window and at least one matching `main` observation.
Without matching `main` evidence, classify the failure as `Unknown` for insufficient baseline history, rather than
`Broken`. Show counts only for failed tests; do not show counts for non-test failures or package-level fallback blocks.
For test-rate details, link the PR failure and representative `main` pass/fail runs, plus representative PR branch
pass/fail runs for each commit. Keep every contributing run link in the local JSON report rather than listing them all
in Markdown. Choose the most recent run for each outcome and commit as the representative pass or fail link. Link
directly to the job when a job URL is available; otherwise link to the workflow run.

The generated file is `results/github/<owner>/<repository>/DEPENDENCY-BUMP-PRS.md` in the script's working directory. Standard output reports successful completion;
errors that prevent file creation or the check from completing are written to standard error and cause a non-zero exit
code.
