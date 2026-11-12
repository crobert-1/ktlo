# Shared functionality

The shared functionality is a set of common functions and utilities that can be used by the various KTLO scripts and tools.
The shared functionality is implemented in a separate package or module to keep the code organized and maintainable.

## Overview

### YAML configuration

Place configuration loading and validation in `internal/config`. The executable accepts a required `--config` path.
The YAML contains top-level camelCase blocks for functionality settings and a shared `repositories` list. `prStatus`
contains an `enabled` boolean that defaults to true when omitted; `failingTestCheck.enabled` defaults to false and its
`lookback` defaults to `7d`, accepting positive Go duration values such as `30m` and `12h` plus `d` for days.
`prStatus.updateFlakyAnalysisReport` defaults to true; when false, PR Status reads the
existing flaky-analysis report without updating it. Unknown keys are config errors. When any functionality is enabled, at least one
repository is required, and the script runs enabled functionality for each configured URL.

The defaulting behavior belongs to the functionality configuration:

```go
func (c PRStatusConfig) IsEnabled() bool
func (c PRStatusConfig) ShouldUpdateFlakyAnalysisReport() bool
func (c FailingTestCheckConfig) EffectiveLookback() (time.Duration, error)
func (c Config) AnyCheckEnabled() bool
```

### Common values

`internal/common` contains repository-independent values shared across packages and has no dependencies on other
repository packages. It owns the dependency-update PR author type and list:

```go
type PullRequestAuthor struct {
    Login string
}

func DependencyBumpPullRequestAuthors() []PullRequestAuthor
```

The function returns a fresh slice. PR Status and Failing Test Check pass this list to the GitHub client.

### GitHub GraphQL client

The initial API client supports GitHub.com only; GitLab and GitHub Enterprise support are deferred. Place it in `internal/github` and expose
a client method that sends a GraphQL query with optional variables and decodes the response data into a caller-provided
result value. Read credentials from `GH_TOKEN`, falling back to `GITHUB_TOKEN` when the first value is unset or blank,
and send the token in the HTTP `Authorization: Bearer` header. Handle transport failures,
non-2xx HTTP responses, GraphQL errors returned in successful HTTP responses, malformed JSON, and oversized responses
as errors. Use a 10-second default request timeout and cap response bodies at 1 MiB.

Never include the token in output or log messages; redact it from errors returned by the GitHub client.

### Provider query documents and callers

The initial implementation does not need a generic query builder. Store query contents as named constants in
`internal/github/queries.go`. Put provider caller methods in the appropriate GitHub source file; each caller passes its
named query constant to the GitHub client's query method and handles the typed result. Pass runtime values through
GraphQL variables rather than interpolating them into query text. Use the same pattern in a future `internal/gitlab`
package, with its own schema-specific query constants and caller methods.

### GraphQL transport client

Place provider-independent HTTP request and GraphQL response handling in `internal/graphql`. Its API accepts an
endpoint, an optional HTTP client, structured variables at initialization, a query, and a result target:

```go
func NewClient(endpoint string, httpClient *http.Client, variables any, headers http.Header) *Client
func (c *Client) Query(ctx context.Context, query string, result any) error
func (c *Client) QueryWithVariables(ctx context.Context, query string, variables any, result any) error
```

The transport client stores structured variables and HTTP headers supplied at initialization and includes them in each
`Query` request. `QueryWithVariables` uses the supplied structured variables for that request. Both methods handle JSON
encoding and decoding, HTTP status failures, GraphQL response errors, request timeouts, and response size limits.
The GitHub package supplies its endpoint and a named repository variables struct, then delegates query execution to
this client.

### Output formatting

Use `internal/output` for formatting script results. Its APIs are:

```go
type PullRequestResult struct {
    PullRequest                 github.PullRequest
    CIResults                   github.PullRequestCIResults
    FailureAssessments          []testresults.FailureAssessment
    AnalysisUnavailableReason   testresults.UnknownReason
}

func PrintPullRequestCIResults(stdout io.Writer, repository github.Repository, results []PullRequestResult) error

type FlakyAnalysisResult struct {
    Repository        github.Repository
    Assessments       []testresults.FailureAssessment
    UnavailableReason testresults.UnknownReason
}

func PrintFlakyAnalysisResults(stdout io.Writer, results []FlakyAnalysisResult) error
```

It writes a Markdown document titled `PR Status Results` to
`results/<provider>/<owner>/<repository>/RESULTS.md`; the current GitHub output uses `github` as the provider directory.
The document has separate `Passed`, `Failed`, and `Others` sections. Each section has a table with `Author`, `Pull
Request`, and `CI/CD` columns. The Pull Request cell links the formatted text
`title (#number)` to the PR URL. A `SUCCESS` rollup is Passed; `FAILURE` and `ERROR` are Failed; pending, expected,
missing, unknown, or unavailable results are Others. A pull request without an analysis result is also Others. Once the
file is completely written, it prints the path and a success confirmation to the supplied `stdout` writer. Writing uses
a temporary file and replaces the repository's results file only after the new document is complete.

`PrintPullRequestCIResults` adds failure details below the `Failed` table. Test failures include separate main and
per-PR-commit outcomes and related CI context; non-test failures include identity, reason, links, and classification
without pass/fail rates. `AnalysisUnavailableReason` explains missing or stale reports. `PrintFlakyAnalysisResults`
writes local lookup results to stdout and distinguishes unavailable reports from valid lookups with no matches. The
formatter maps `UnknownReason` values to human-readable explanations.

### GitHub API client

The initial implementation supports GitHub.com only. The `internal/github` client sends read-only GraphQL queries
using a token from `GH_TOKEN` or `GITHUB_TOKEN`. GitLab support is deferred. The public API is:

```go
func NewClient(ctx context.Context, repositoryURL string, token string, httpClient *http.Client) (*Client, error)
func (c *Client) Query(ctx context.Context, query string, result any) error
func (c *Client) OpenPullRequests(ctx context.Context) ([]PullRequest, error)
func (c *Client) OpenDependencyBumpPullRequests(ctx context.Context, authors []common.PullRequestAuthor) ([]PullRequest, error)
func (c *Client) AnalyzePullRequests(ctx context.Context, pullRequests []PullRequest) ([]PullRequestCIResults, error)
func (c *Client) AnalyzePullRequest(ctx context.Context, pullRequest PullRequest) (PullRequestCIResults, error)
func (c *Client) TargetRepository() Repository
func ParseRepositoryURL(rawURL string) (Repository, error)
func (c *Client) ListAllWorkflowRunAttempts(ctx context.Context, startedAfter time.Time) ([]WorkflowRunAttempt, error)
func (c *Client) ListWorkflowJobs(ctx context.Context, runID int64, attempt int) ([]WorkflowJob, error)
func (c *Client) JobLog(ctx context.Context, jobID int64) (io.ReadCloser, error)
var ErrInvalidRepositoryURL error
var ErrMissingToken error
```

The workflow result types preserve provider details at the GitHub boundary:

```go
type WorkflowRunAttempt struct {
    RunID              int64
    Attempt            int
    Name               string
    Branch             string
    HeadSHA             string
    RunStartedAt        time.Time
    Status              string
    ProviderConclusion  string
    WorkflowURL         string
    PullRequests        []AssociatedPullRequest
}

type AssociatedPullRequest struct {
    PRNumber int
    URL      string
    HeadSHA  string
}

type WorkflowJob struct {
    ID                  int64
    Name                string
    ProviderConclusion  string
    JobURL              string
    Steps               []WorkflowStep
}

type WorkflowStep struct {
    Number              int
    Name                string
    ProviderConclusion  string
}
```

For each supplied PR, `AnalyzePullRequests` calls `AnalyzePullRequest`, returning
successful results with any collected query errors. The per-PR query reads the status rollup and paginates the PR's
check runs and legacy status contexts. Each context includes its GraphQL type and available status, conclusion,
description, required flag, detail links, and check summary fields.

`NewClient` validates and stores the HTTPS GitHub.com repository URL using the client's private
`validateRepositoryURL` method. It configures `internal/graphql` with the stored owner and name as named fields in a
structured value, then uses `Query` to verify that the repository is available to the authenticated user before returning. The
probe requests the repository ID. `OpenPullRequests` fetches every open pull request, and
`OpenDependencyBumpPullRequests` fetches open PRs and filters the results against the supplied author list. Both methods
use named queries in `internal/github/queries.go`, include each PR's draft state, and follow GraphQL cursors until the
final page. The repository pull request connection has no author filter, so author matching happens in Go. Passing a
nil HTTP client selects the default 10-second timeout; callers may provide an HTTP client
for transport configuration. `Query` delegates HTTP and GraphQL handling to `internal/graphql` and redacts the token
from returned errors. Query documents use
the `$owner` and `$name` variables configured during initialization.

The workflow API returns individual attempts and preserves raw provider statuses and conclusions. Run and job listing
paginate internally and return records already collected with an error if a later page fails. The caller processes only
completed attempts; completed cancelled, timed-out, and other inconclusive outcomes remain available. `JobLog` returns
a plain-text stream for any completed job, including successful jobs. `ParseRepositoryURL` extracts and validates a
GitHub repository without network access; `NewClient` reuses it before performing the authenticated access check.
Workflow run data includes PR associations with `PRNumber` and head SHA when available.

The top-level `run` method creates one client per repository, including its repository preflight, and reuses it across
enabled functionality. The `internal/prstatus` runner API is:

```go
func NewRunner(githubClient *github.Client, maxReportAge time.Duration) (*Runner, error)
func (r *Runner) Run(ctx context.Context, stdout io.Writer) error
```

The runner loads the report once per repository, checks freshness, and obtains assessments matching each PR's exact
head SHA through `internal/testresults`. Missing or stale reports do not prevent PR Status output.

The runner retrieves open dependency update pull requests using the common author list, queries their CI status through
`AnalyzePullRequests`, and passes combined PR, CI, and assessment results to `output.PrintPullRequestCIResults`. The
formatter writes the three CI/CD groups to `results/github/<owner>/<repository>/RESULTS.md` and reports successful
file creation on stdout.

### Flaky test checker

Flaky test analysis writes a versioned JSON report locally for PR Status and lookup mode. No scheduled GitHub Actions
workflow or remote artifact upload is planned. Its lookback is bounded by provider log retention. The report contract,
freshness rules, and classification decisions are defined in [flaky-test-analysis.md](flaky-test-analysis.md).

The analyzer API is:

```go
func NewAnalyzer(source workflowSource, settings config.FailingTestCheckConfig) (*Analyzer, error)
func (a *Analyzer) Run(ctx context.Context) error
```

The source contract is:

```go
type workflowSource interface {
    OpenDependencyBumpPullRequests(ctx context.Context, authors []common.PullRequestAuthor) ([]github.PullRequest, error)
    ListAllWorkflowRunAttempts(ctx context.Context, startedAfter time.Time) ([]github.WorkflowRunAttempt, error)
    ListWorkflowJobs(ctx context.Context, runID int64, attempt int) ([]github.WorkflowJob, error)
    JobLog(ctx context.Context, jobID int64) (io.ReadCloser, error)
}
```

The private parser has this boundary:

```go
func parseJobLog(job github.WorkflowJob, log io.Reader) (testresults.Job, error)
```

The source supplies dependency-update PRs, run attempts, jobs, and job-log streams. The analyzer parses logs privately,
merges observations into the shared report, and writes it. It continues after individual collection or parse failures,
retains available data, marks the report incomplete, writes it, and returns the collected errors. When standalone
Failing Test Check is enabled or PR Status has `updateFlakyAnalysisReport` enabled, the top-level runner runs the
analyzer once per repository before PR Status. Execution remains sequential; parallel analyzers are a future
requirement.

PR Status and Failing Test Check use `internal/testresults` for the shared report model, loading, writing, merging,
history, and assessment API. The report hierarchy is `Report -> Workflow -> Job -> Step -> Test`; tests sit at the
deepest known parent. Failed tests store `FailureMessage`, while passing or skipped tests omit it. `OutcomeCounts` is a
computed value containing pass and failure counts only; `PassPercentage() (int, error)` computes the rounded percentage
and errors when there are no samples. Classification and rate calculations remain separate from report data and file
I/O.

The principal report and assessment types are:

```go
type Report struct {
    SchemaVersion  int
    Provider       Provider
    Repository     Repository
    GeneratedAt    time.Time
    AnalysisStatus AnalysisStatus
    Workflows      []Workflow
}

type FailureAssessment struct {
    PullRequest    *PullRequestRef
    Location       CIFailureLocation
    Classification Classification
    History        OutcomeHistory
    RelatedCI      []CIContextHistory
    Partial        bool
    UnknownReason  UnknownReason
}

type OutcomeHistory struct {
    MainBranch OutcomeSummary
    PRCommits  []CommitOutcomeSummary
}
```

Each `OutcomeSummary` and `CommitOutcomeSummary` contains computed counts and latest pass/failure `RunReference` values.
`OutcomeCounts` contains `Passes` and `Failures`; no rates are stored in the JSON report. A `CIFailureLocation` identifies
the deepest known workflow, job, step, or test and includes its known parent names. A `CIContextHistory` records the
separate workflow, job, or step history related to a test failure.

The report API includes `NewReport`, `Load`, `Write`, private `pathFor`, `MergeWorkflowRuns`, `GetOutcomeHistory`,
`AssessFailure`, `FindAssessments`, and `AssessmentsForPullRequest`. `FindAssessments` accepts exactly one exact test,
workflow, or job selector. `AssessmentsForPullRequest` matches the PR number and exact current head SHA. `GetOutcomeHistory`
derives counts and latest evidence links from observations, keeping `main` separate from each PR commit. `AssessFailure`
computes classification, history, partial state, related CI histories, and any unknown reason.

The method signatures are:

```go
func NewReport(provider Provider, repository Repository) *Report
func (r *Report) Load() error
func (r *Report) Write() error
func (r *Report) MergeWorkflowRuns(runs []Workflow, cutoff time.Time) error
func (r *Report) pathFor() (string, error)
func (r *Report) GetOutcomeHistory(location CIFailureLocation, pullRequest *PullRequestRef) OutcomeHistory
func (r *Report) AssessFailure(location CIFailureLocation, pullRequest *PullRequestRef) FailureAssessment
func (r *Report) FindAssessments(selector FailureSelector) ([]FailureAssessment, error)
func (r *Report) AssessmentsForPullRequest(pullRequest PullRequestRef) []FailureAssessment
func (counts *OutcomeCounts) PassPercentage() (int, error)
```

`Load` preserves not-found errors and rejects malformed, unsupported-version, or wrong-repository reports. `Write`
validates and atomically replaces the report file. `MergeWorkflowRuns` removes old records before the cutoff, ignores
older incoming runs, and keeps an existing record when its `(runId, runAttempt)` key is already present. Incoming runs
require a stable run ID, positive attempt, and non-zero start time; invalid input leaves the report unchanged. Keep
stored report types and file I/O in `report.go`. Keep assessment types, classification, outcome history, and rate
calculations together in `assessment.go` because those APIs overlap.
