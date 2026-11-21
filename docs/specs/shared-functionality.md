# Shared functionality

The shared functionality is a set of common functions and utilities that can be used by the various KTLO scripts and tools.
The shared functionality is implemented in a separate package or module to keep the code organized and maintainable.

## Overview

### YAML configuration

Place configuration loading and validation in `internal/config`. The default config path is `config.yaml`, with
`--config` and `-c` available as command-line overrides. The selected command chooses which operation runs; configuration
contains shared repository definitions and reusable defaults, not command enablement switches.

Each repository entry is an object containing a required `url`, optional `enabled` (default `true`), and optional
`local-path`. URL-only entries fail config validation. Disabled entries are excluded from all commands. Vulnerability
Check uses `local-path` to scan the existing checkout as-is without changing repository files or Git metadata. Every
command requires at least one enabled repository.

`dependency-pr-status.update.dependency-update-authors` defaults to the configured dependency bot logins when omitted. `dependency-pr-status.update.stale-age`
defaults to `7d`, and the `dependency-pr-status update --stale-age` option can override it for one run. The flaky-test
JSON report lookback defaults to `7d`; `flaky-test-analysis --lookback` can override it for one run. The
`--overwrite` option forces a full report rebuild. Targeted workflow parsing is enabled by default for dependency PR
status; `--no-parse-workflow-failures` disables it for one run. These invocation controls are not stored as config
fields.

`telemetry.enabled` defaults to false. When enabled, non-empty `OTEL_EXPORTER_OTLP_ENDPOINT`,
`OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`, `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT`, and their matching `*_INSECURE` values
are applied to the corresponding process environment variables before telemetry setup; omitted or empty values leave
existing environment values intact. The exporters use gRPC directly, so no protocol selector is needed. Unknown YAML
keys are errors.

The config API owns parsing and reusable defaults:

```go
func (c DependencyPRStatusUpdateConfig) EffectiveStaleAge() (time.Duration, error)
func (c DependencyPRStatusUpdateConfig) EffectiveDependencyUpdateAuthors() []string
func (c FlakyTestAnalysisConfig) EffectiveLookback() (time.Duration, error)
func (c Config) EnabledRepositories() []RepositoryConfig
func (c Config) EffectiveLogLevel() logging.Level
func (c Config) Validate() error
```

### Common values

`internal/common` contains repository-independent values shared across packages and has no dependencies on other
repository packages. It owns the dependency-update PR author type:

```go
type PullRequestAuthor struct {
    Login string
}

type Repository struct {
    Owner string
    Name  string
}

type PullRequestReference struct {
    Number  int
    URL     string
    HeadSHA string
}

func (reference *PullRequestReference) Equals(other *PullRequestReference) bool

```

PR Status converts its configured author logins into `PullRequestAuthor` values and passes them to the GitHub client.
`Repository` and `PullRequestReference` are also shared from this package by GitHub and test-results report types.

`internal/common` also owns the provider-neutral workflow run, pull request association, job, and step types used by
the analyzer's provider interface and the GitHub client. Provider packages convert their native records into these
types at their boundaries.

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
    WorkflowFailureParseErrors  []string
    Lookback                    time.Duration
}

func PrintPullRequestCIResults(ctx context.Context, stdout io.Writer, repository github.Repository, results []PullRequestResult) error

type FlakyAnalysisResult struct {
    Repository        github.Repository
    Assessments       []testresults.FailureAssessment
    UnavailableReason testresults.UnknownReason
}

func PrintFlakyAnalysisResults(stdout io.Writer, results []FlakyAnalysisResult) error
```

It writes a Markdown document titled `Dependency Bump PR Status` to
`results/<provider>/<owner>/<repository>/DEPENDENCY-BUMP-PRS.md`; the current GitHub output uses `github` as the provider directory.
The document has separate `Passed`, `Failed`, and `Others` sections. Each section has a table with `Author`, `Pull
Request`, and `CI/CD` columns. The Pull Request cell links the formatted text
`title (#number)` to the PR URL. The formatter derives the effective state from individual CI contexts: all reported
contexts must succeed for `Passed`, and any `FAILURE` or `ERROR` context means `Failed`. Cancelled, skipped, pending,
expected, unknown, or unavailable contexts belong under `Others`; cancelled and skipped contexts are neither failures
nor passes. The provider rollup is a fallback when individual contexts are unavailable. A pull request without an
analysis result is also Others. Once the file is completely written, it prints the path and a success confirmation to
the supplied `stdout` writer. Writing uses a temporary file and replaces the repository's results file only after the
new document is complete.

`PrintPullRequestCIResults` adds one failure-details subsection per failed PR, in the same order as the `Failed` table.
It puts all distinct test and non-test failures for that PR in one table with failure, classification, latest failure,
`main` history, PR commit history, and related CI history columns. Unknown classifications include their reason;
failure messages link to their source run when available. Counts and rates are separated by `main` and PR commit,
with associated workflow/job/step history when available. Non-test rows and unavailable values use `—` for history
cells. Successful workflows contribute one inferred pass to each known child job, step, and test location, even when
that run has no explicit child result. Explicit child failures are retained beside the inferred pass, while an
explicit child pass is not double-counted. History cells show the number of inferred passes, and non-test rows
disclose when inferred passes informed classification. Failed CI contexts can use these observations at the same
most-specific location. Unmatched failed contexts appear as `Unknown` rows with their available name, description or
summary, and details link. Incomplete but usable evidence is labeled partial. Package-level fallback rows omit
aggregate counts. Parent failures explained by more specific child failures are suppressed. Report-wide issues appear
once per PR with a human-readable cause and explanation of the analysis impact. Individual `Unknown` classifications
explain their reason. `PrintFlakyAnalysisResults` writes local lookup results to stdout and
distinguishes unavailable reports from valid lookups with no matches. The
formatter maps `UnknownReason` values to human-readable explanations.
After failure details, it identifies cancelled checks with links when available and gives a count of skipped checks for
PRs that have no failed contexts.

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
func (c *Client) ListAllWorkflowRunAttempts(ctx context.Context, startedAfter time.Time) ([]common.WorkflowRunAttempt, error)
func (c *Client) ListWorkflowJobs(ctx context.Context, runID int64, attempt int) ([]common.WorkflowJob, error)
func (c *Client) WorkflowAttemptLogArchive(ctx context.Context, runID int64, attempt int) (io.ReadCloser, error)
var ErrInvalidRepositoryURL error
var ErrMissingToken error
```

The shared provider-neutral workflow types live in `internal/common`:

```go
type WorkflowRunAttempt struct {
    RunID        string
    Attempt      int
    Name         string
    Branch       string
    HeadSHA      string
    RunStartedAt time.Time
    Status       string
    Conclusion   string
    WorkflowURL  string
    PullRequests []PullRequestReference
}

type WorkflowJob struct {
    ID         string
    Name       string
    Conclusion string
    JobURL     string
    Steps      []WorkflowStep
}

type WorkflowStep struct {
    Number     int
    Name       string
    Conclusion string
}
```

GitHub's private REST response types retain GitHub-specific numeric IDs and conclusion enums, then convert to these
shared types. The failing-test provider interface and GitHub client use the same common workflow records, avoiding
duplicate provider-neutral struct definitions.

For each supplied PR, `AnalyzePullRequests` calls `AnalyzePullRequest`, returning
successful results with any collected query errors. The per-PR query reads the status rollup and paginates the PR's
check runs and legacy status contexts. Each context includes its GraphQL type and available status, conclusion,
description, required flag, detail links, and check summary fields.

`PullRequestCIResults.EffectiveRollupState` derives the effective state from those contexts. Explicit `FAILURE` or
`ERROR` contexts are failures; `CANCELLED` and `SKIPPED` contexts are unresolved and do not count as failures or passes.
The provider's aggregate rollup is used only when individual contexts are unavailable.

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

The workflow API returns individual attempts and preserves raw provider statuses and conclusions as strings in the
common types. Run and job listing
paginate internally and return records already collected with an error if a later page fails. The caller processes only
completed attempts; completed cancelled, timed-out, and other inconclusive outcomes remain available.
`WorkflowAttemptLogArchive` returns a ZIP stream for an individual attempt. The failing-test provider reads numbered
step log entries and job-level transcript entries such as `0_<job>.txt` separately; the latter is not a workflow step.
`ParseRepositoryURL` extracts and validates a GitHub repository without network access; `NewClient` reuses it before
performing the authenticated access check.
Workflow run data includes PR associations with `Number` and head SHA when available. The provider adapter converts the
common string run ID to GitHub's numeric ID only when calling GitHub job and log endpoints.

The top-level `run` method creates one client per repository, including its repository preflight, and reuses it across
the selected command. The `internal/prstatus` runner API is:

```go
func NewRunner(
    githubClient *github.Client,
    dependencyUpdateAuthors []string,
    lookback time.Duration,
    staleAge time.Duration,
    parseWorkflowFailures bool,
) (*Runner, error)
func (r *Runner) Run(ctx context.Context, stdout io.Writer) error
```

The runner loads the report once per repository and obtains assessments matching each PR's exact head SHA through
`internal/testresults`. It passes the configured lookback to the formatter so failure rates use the same setting.
Valid reports are used regardless of age; missing reports do not prevent PR Status output.

The runner retrieves open dependency update pull requests using the common author list, queries their CI status through
`AnalyzePullRequests`, and passes combined PR, CI, and assessment results to `output.PrintPullRequestCIResults`. The
formatter writes the CI/CD groups to `results/github/<owner>/<repository>/DEPENDENCY-BUMP-PRS.md` and reports successful
file creation on stdout.

### Flaky test checker

Flaky test analysis writes a versioned JSON report locally for PR Status and lookup mode. No scheduled GitHub Actions
workflow or remote artifact upload is planned. Its lookback is bounded by provider log retention. The report contract
and classification decisions are defined in [flaky-test-analysis.md](flaky-test-analysis.md).

The analyzer API is:

```go
func NewAnalyzer(workflowProvider WorkflowProvider, testLogParsers ...TestLogParser) (*Analyzer, error)
func NewAnalyzerWithLookbackDuration(workflowProvider WorkflowProvider, lookbackDuration time.Duration, testLogParsers ...TestLogParser) (*Analyzer, error)
func (a *Analyzer) Run(ctx context.Context) error
```

The provider contract is:

```go
type WorkflowProvider interface {
    Provider() testresults.Provider
    Repository() common.Repository
    ListWorkflowRunAttempts(ctx context.Context, startedAfter time.Time) ([]common.WorkflowRunAttempt, error)
    ListWorkflowJobs(ctx context.Context, runID string, attempt int) ([]common.WorkflowJob, error)
    StepLogsForAttempt(ctx context.Context, runID string, attempt int, jobs []common.WorkflowJob) (WorkflowAttemptLogs, error)
}
```

The provider supplies its identity, repository, run attempts, jobs, and attempt logs. It keeps job-level transcripts
separate from numbered step logs. When the archive provides only a combined job transcript, the GitHub provider uses
step start and completion timestamps to split it into step logs. The analyzer associates archive entries with selected
jobs and steps, passes step logs to the configured language parsers, and parses a job transcript into `job.tests` when
the source step cannot be identified. A combined transcript is not used as an individual failed step's excerpt when
the step cannot be isolated. It merges observations into the shared report and writes it. It continues after individual collection or parse failures,
retains available data, marks the report incomplete, writes it, and returns the collected errors. When standalone
Failing Test Check is enabled alongside PR Status, the top-level runner runs the analyzer once per repository before PR
Status. Otherwise PR Status reads the existing report without modifying it. Execution remains sequential; parallel
analyzers are a future requirement.

PR Status and Failing Test Check use `internal/testresults` for the shared report model, loading, writing, merging,
history, and assessment API. The report hierarchy is `Report -> Workflow -> Job -> Step -> Test`; tests sit at the
deepest known parent. Failed tests store `FailureMessage`, while passing or skipped tests omit it. `OutcomeCounts` is a
computed value containing pass, failure, and inferred pass counts; `PassPercentage() (int, error)` computes the rounded percentage
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

type CIContextFailure struct {
    Name           string
    DetailsURL     string
    TargetURL      string
    FailureMessage string
    RunStartedAt   time.Time
}

func (assessment *FailureAssessment) LatestFailureEvidence() *FailureEvidence

type FailureEvidence struct {
	SourceRun RunReference
	Message   string
}

type OutcomeHistory struct {
    MainBranch OutcomeSummary
    PRCommits  []CommitOutcomeSummary
}

type OutcomeSummary struct {
    Counts        OutcomeCounts
    LatestPass    *RunReference
    LatestFailure *FailureEvidence
}

type CommitOutcomeSummary struct {
    Branch        string
    CommitSHA     string
    Counts        OutcomeCounts
    LatestPass    *RunReference
    LatestFailure *FailureEvidence
}

type OutcomeCounts struct {
    Passes         int
    Failures       int
    InferredPasses int
}
```

`testresults.Repository` and `testresults.PullRequestRef` are aliases of the common types, preserving the report's JSON
field names while sharing repository and pull-request identity across packages. The persisted `Workflow` remains a
separate type because it adds normalized outcomes, parsed tests, and failure excerpts to a provider run attempt.
Each `OutcomeSummary` and `CommitOutcomeSummary` contains computed counts, a latest passing `RunReference`, and latest
failure `FailureEvidence` that keeps its source run and message together.
`OutcomeCounts` contains total `Passes`, `Failures`, and the subset of passes inferred from successful workflows in
`InferredPasses`; these are computed during history aggregation and are not stored in the JSON report. A successful
workflow contributes one inferred pass per known child location and workflow attempt, unless that location already
has an explicit pass in that attempt. Explicit child failures remain separate outcomes. A `CIFailureLocation`
identifies the deepest known workflow, job, step, or test and includes its known parent names. A `CIContextHistory`
records the separate workflow, job, or step history related to a test failure.

The report API includes `NewReport`, `Load`, `Write`, private `pathFor`, `MergeWorkflowRuns`, `GetOutcomeHistory`,
`AssessFailure`, `AssessCIContextFailure`, `FindAssessments`, and `AssessmentsForPullRequest`. `FindAssessments` accepts
exactly one exact test, workflow, or job selector. `AssessmentsForPullRequest` matches the PR number and exact current
head SHA. `AssessCIContextFailure` resolves a failed CI context against report URLs first and an unambiguous exact
name second, then assesses the most-specific matching location even if its report observations are only successful.
It adds the current context failure to that PR head's history. `GetOutcomeHistory` derives counts and latest evidence
links from explicit observations and inferred child passes, keeping `main` separate from each PR commit. Assessment
methods compute classification,
history, partial state, related CI histories, and any unknown reason.

The method signatures are:

```go
func NewReport(provider Provider, repository Repository) *Report
func (r *Report) Load() error
func (r *Report) Write() error
func (r *Report) MergeWorkflowRuns(runs []Workflow, cutoff time.Time) error
func (r *Report) pathFor() (string, error)
func (r *Report) GetOutcomeHistory(location CIFailureLocation, pullRequest *PullRequestRef) OutcomeHistory
func (history *OutcomeHistory) InferredWorkflowPassCount() int
func (r *Report) AssessFailure(location CIFailureLocation, pullRequest *PullRequestRef) FailureAssessment
func (r *Report) AssessCIContextFailure(context CIContextFailure, pullRequest *PullRequestRef) (FailureAssessment, bool)
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
