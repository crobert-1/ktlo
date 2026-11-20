# GoLang Script

## Overview

The GoLang script is the entry point for running enabled KTLO functionality against configured repositories.

## Input

The script accepts a required `--config` argument pointing to a YAML file. The YAML has top-level camelCase keys for
each functionality and a `repositories` list shared by enabled functionality. The initial configuration is:

```yaml
otelGoInstrumentation:
  enabled: false
  OTEL_EXPORTER_OTLP_ENDPOINT: ""
  OTEL_EXPORTER_OTLP_INSECURE: ""
  OTEL_EXPORTER_OTLP_METRICS_ENDPOINT: ""
  OTEL_EXPORTER_OTLP_METRICS_INSECURE: ""
  OTEL_EXPORTER_OTLP_TRACES_ENDPOINT: ""
  OTEL_EXPORTER_OTLP_TRACES_INSECURE: ""

prStatus:
  enabled: true
  dependencyUpdateAuthors:
    - dependabot
    - renovate
    - srv-gh-o11y-gdi
  parseWorkflowFailures: false
  staleAge: 7d

failingTestCheck:
  jsonReport:
    enabled: false
    lookback: 7d
    overwriteExistingReport: false
  readme:
    enabled: false
  maxArtifactAgeHours: 48

repositories:
  - url: https://github.com/owner/repository
    enabled: true
    # localPath: /Users/name/dev/repository
```

Each repository entry must be an object with a required `url`, an optional `enabled` boolean that defaults to `true`,
and an optional `localPath`. URL-only string entries fail config validation. Disabled repositories are excluded from
every check. Vulnerability Check scans a supplied `localPath` as-is without fetching, switching branches, or modifying
repository files or Git metadata. Its `origin` remote must match the configured URL. The scan may include uncommitted
changes, while report links point to the checkout's current `HEAD` commit. Relative paths are resolved from the process
working directory.

`prStatus.enabled` defaults to `true` if omitted. `dependencyUpdateAuthors` defaults to the existing dependency bot
logins shown above; set it to the author logins whose PRs should be included in dependency-update status. Setting
`prStatus.staleAge` defaults to `7d`; set it to a positive duration such as `60m`, `12h`, or `7d` to place PRs whose
latest commit is older than the duration in the `Stale` table. Setting `prStatus.enabled` to `false` skips PR status.
`prStatus.parseWorkflowFailures` defaults to `false`; when enabled, PR Status parses failed GitHub Actions runs linked
from failed PR checks, using the existing workflow and test-log parsing pipeline. It uses the JSON report lookback and
adds observations to an in-memory report only; it does not run the repository-wide analysis or write
`workflow-test-analysis.json`.
At least one enabled repository is required when any functionality is enabled. `failingTestCheck.jsonReport.enabled` and
`failingTestCheck.readme.enabled` default to `false`; JSON `lookback` defaults to `7d` and accepts positive Go duration
values such as `30m` and `12h`, plus `d` for days. Unknown configuration keys
are errors so misspelled settings are not silently ignored. Future functionalities will add their own top-level
camelCase blocks.
PR Status reads the existing local flaky-analysis report for failure details when one is available. Unless
`prStatus.parseWorkflowFailures` is enabled, it does not start analysis. Targeted parsing enriches an in-memory report
and never writes the JSON report. Only `failingTestCheck.jsonReport.enabled` runs repository-wide analysis and writes
the JSON report. `failingTestCheck.readme.enabled` renders the local report to `FLAKY-AND-FAILING-TESTS.md`; when both are enabled, JSON
collection runs first. README-only mode does not call GitHub APIs.
The top-level `run` delegates command-line parsing to `parseOptions` and YAML loading to a dedicated `loadConfig` helper.

Failing-test lookback uses GitHub's server-side workflow-run `created` filter to avoid scanning the repository's full
workflow history. ktlo then retains attempts whose `run_started_at` is at or after the configured cutoff. A workflow run
created before the cutoff but rerun during the lookback can be omitted because GitHub filters by the parent run's
creation time. GitHub returns at most 1,000 results for a filtered search; when the matching count exceeds that limit,
ktlo processes the returned results and marks the analysis incomplete.

For each workflow attempt containing successful or failed jobs, ktlo downloads one attempt-wide log archive, then
associates its step log files with the jobs returned by GitHub. Only logs for successful and failed jobs are parsed.
Archives are limited to 250 MiB compressed and 250 MiB of selected uncompressed log text; exceeding either limit marks
the analysis incomplete. If an archive cannot be associated with a selected job, the analysis is also marked incomplete.
The attempt archive can include logs for jobs that ktlo does not parse.

The CLI parsing API may return an optional flaky-analysis selector:

```go
type cliOptions struct {
    ConfigPath string
    Selector   *testresults.FailureSelector
}

func parseOptions(args []string, stderr io.Writer) (cliOptions, int)
func runLookup(config config.Config, selector testresults.FailureSelector, stdout io.Writer) error
```

When a selector is supplied, lookup reads local reports and does not create a GitHub client or make network requests.
The configured repository URL is parsed without credentials or network access to locate each report. Without a
selector, the executable follows the normal enabled-check flow and requires the GitHub token.

The top-level script also provides a read-only flaky-analysis lookup mode. It accepts one test name, workflow name, or
job name as a selector and returns the matching flaky-analysis information from the local report. The selector is
looked up for the repositories listed in the config. This mode reads the existing report and does not fetch logs or
start a new analysis.

The current implementation supports GitHub.com repositories. GitLab support is deferred. The script requires a token
from `GH_TOKEN`, falling back to `GITHUB_TOKEN` when the first value is unset or blank. Credentials are not accepted in
the YAML file or command-line arguments.

## CLI contract

The executable is named `ktlo`:

```text
ktlo --config config.yaml
ktlo --config config.yaml --test package.TestName/subtest
ktlo --config config.yaml --workflow CI
ktlo --config config.yaml --job unit-tests
```

Without a lookup selector, the script runs enabled functionality for every configured repository. With a lookup
selector, it reads the local report for each configured repository and returns matching information. Select at most
one of `--test`, `--workflow`, or `--job` in a single invocation. Invalid configuration or repository URLs
are usage errors; connection failures, timeouts, unsuccessful HTTP responses, and GraphQL errors are operational
errors. Tokens must never appear in standard output, standard error, or other log messages, including when an API
response echoes one. The script attempts each configured repository even if a previous repository fails, then exits
non-zero if any configured check failed.

Every executable invocation reports its total wall-clock runtime to standard error as its final runtime line, including
usage-error and help invocations. The duration is rounded to the nearest second and includes only non-zero hours,
minutes, and seconds (for example, `Total runtime: 6 minutes, 5 seconds`).

Lookup output includes the matching test or workflow/job classification, pass/fail counts and percentages, `main`
history, PR branch and commit history, partial-data state, any `Unknown` reason, and representative evidence links.
The output preserves the distinction between `main` and each PR commit.

For each repository, `runChecks` creates one GitHub client and reuses it across checks that need provider access. If
`failingTestCheck.jsonReport.enabled` is
true, run the analyzer once before PR Status so PR Status can consume the updated report. If the standalone check is
disabled, PR Status may read the existing report but must not start analysis or modify it. Continue to PR Status after a
partial-analysis error if the report was written, and return a non-zero exit code for collected errors.

PR status results are written to `results/github/<owner>/<repository>/DEPENDENCY-BUMP-PRS.md` in the script's working directory,
grouped into `Passed`, `Failed`, and `Others` Markdown tables with a separate author column, a clickable `title
(#number)` pull request link, and the CI/CD rollup state. After each results file is completely written, the script
reports its path and successful completion on standard output. Review and merge status classification is future work.
Usage and operational errors are written to standard error. Invalid arguments, unsupported configuration keys, or a
check that cannot complete result in a non-zero exit code.

## Runtime instrumentation

OpenTelemetry traces and duration metrics are disabled by default. Set `otelGoInstrumentation.enabled` to `true` to
enable them. The `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`, and
`OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` YAML keys map to the matching environment variables consumed by the Go
OpenTelemetry SDK. `OTEL_EXPORTER_OTLP_INSECURE`, `OTEL_EXPORTER_OTLP_TRACES_INSECURE`, and
`OTEL_EXPORTER_OTLP_METRICS_INSECURE` configure plaintext gRPC connections. The common endpoint and insecure values
apply to both signals; signal-specific values override them. A configured non-empty value overrides the process
environment, and an omitted or empty value leaves an existing process environment value unchanged. The gRPC exporters
are selected directly by the application, so `OTEL_EXPORTER_OTLP_PROTOCOL` does not need to be configured. Export
setup or shutdown failures are written to standard error without changing the check's exit code. The CLI flushes
telemetry before exiting.

Traces cover the command and repository execution, PR status and failing-test checks, GitHub API methods and outbound
HTTP requests, workflow attempt archive download and parsing, report merge/write, and PR result output. GitHub REST JSON spans
include the HTTP method and full request URL to distinguish endpoints, IDs, and query parameters. The
`ktlo.operation.duration` histogram uses seconds and a bounded `ktlo.operation` attribute. HTTP client instrumentation
also exports `http.client.request.duration`. Metric attributes must not include repository, run, job, or PR identifiers.
Do not record credentials, request/response bodies, or log contents.

## Build

Run `make ktlo` from the repository root to build the Go executable and its dependencies. The output is
`./bin/ktlo`.

## Functionalities

The script delegates each functionality to its own package or module. Initial functionality must not create, update,
merge, or comment on repository resources.

## Future requirements

- Add progress log messages at appropriate log levels so users can follow repository processing and check progress.
- Rename the script's run modes to names that describe their behavior more clearly.
