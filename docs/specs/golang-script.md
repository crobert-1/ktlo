# GoLang Script

## Overview

The `ktlo` executable runs one selected command against repositories in a shared YAML configuration file. A command name
selects the operation; config values provide reusable settings and repository definitions.

## Configuration

The default config path is `config.yaml`. Override it with `--config FILE` or `-c FILE`; the option may appear before or
after the command. The example configuration is:

```yaml
logLevel: default

otelGoInstrumentation:
  enabled: false
  OTEL_EXPORTER_OTLP_ENDPOINT: ""
  OTEL_EXPORTER_OTLP_INSECURE: ""
  OTEL_EXPORTER_OTLP_METRICS_ENDPOINT: ""
  OTEL_EXPORTER_OTLP_METRICS_INSECURE: ""
  OTEL_EXPORTER_OTLP_TRACES_ENDPOINT: ""
  OTEL_EXPORTER_OTLP_TRACES_INSECURE: ""

prStatus:
  dependencyUpdateAuthors:
    - dependabot
    - renovate
    - srv-gh-o11y-gdi
  staleAge: 7d

failingTestCheck:
  jsonReport:
    lookback: 7d

repositories:
  - url: https://github.com/owner/repository
    enabled: true
    # localPath: /Users/name/dev/repository
```

The config stores the progress log level, shared repositories, PR status author filters and default stale threshold,
flaky-test lookback, and optional OpenTelemetry settings. Repository entries require `url`; `enabled` defaults to
`true`, and `localPath` is optional. At least one repository must be enabled for a command. Unknown YAML fields are
errors.

`logLevel` defaults to `default` and accepts `none`, `default`, or `detailed`, case-insensitively. `none` suppresses
progress and runtime messages while errors and telemetry setup or shutdown failures remain visible. `default` reports
command and repository progress and preserves the existing vulnerability scan progress. `detailed` also reports
workflow collection and PR status operation stages.

`prStatus.dependencyUpdateAuthors` defaults to `dependabot`, `renovate`, and `srv-gh-o11y-gdi` when omitted.
`prStatus.staleAge` defaults to `7d`; `--stale-age` overrides it for one dependency PR status run.
`failingTestCheck.jsonReport.lookback` defaults to `7d`; `--lookback` overrides it for one flaky-test analysis run.
Both settings accept positive Go durations such as `30m`, `12h`, and `7d`.

## CLI contract

```text
Usage: ktlo [--config FILE] <command> [arguments]

Commands:
  vulnerability-check (v) [ID]
  flaky-test-analysis (f)
  generate-report (g)
  dependency-pr-status (d) <command>
  help [COMMAND]

Options:
  -c, --config FILE  Config file (default: config.yaml)
  -h, --help         Show help
```

Each command has its own help output, available with `<command> --help` or `ktlo help <command>`. Aliases are explicit;
other command prefixes are not inferred. The nested dependency PR action is `update` (`u`), available as
`ktlo dependency-pr-status update`, `ktlo d update`, or `ktlo d u`.

- `ktlo vulnerability-check [options] ID` (`ktlo v [options] ID`) requires one vulnerability identifier and checks
  all enabled repositories. `--release-version TAG` checks that exact release tag and prints one status line per
  repository. The identifier is a command argument, not a config field.
- `ktlo flaky-test-analysis` (`ktlo f`) collects workflow data and updates the local JSON report. `--lookback DURATION`
  overrides the config lookback for this run. `--overwrite` discards stored observations and rebuilds the full configured
  lookback.
- `ktlo generate-report` (`ktlo g`) renders the Markdown flaky-test report from the local JSON report. It does not
  contact GitHub and does not require a token.
- `ktlo dependency-pr-status update` (`ktlo d u`) writes dependency PR status. Targeted parsing of linked failed
  GitHub Actions runs is enabled by default; `--no-parse-workflow-failures` disables it for one invocation.
  `--stale-age DURATION` overrides the configured PR stale threshold.

Vulnerability checks, flaky-test analysis, and dependency PR status require a token from `GH_TOKEN`, falling back to
`GITHUB_TOKEN`. Credentials are not accepted in YAML or as command arguments. Commands attempt every enabled repository
and return a non-zero exit code if any repository operation fails.

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

A command execution reports total runtime to standard error at `default` and `detailed` levels. Help output does not emit
a runtime line. The duration is rounded to the nearest second and includes non-zero hours, minutes, and seconds.

Progress messages use a shared formatter in `internal/logging`. They are written to standard error with a timestamp and
a command-specific prefix: `vulnerability`, `flaky-test-analysis`, `generate-report`, or `dependency-pr-status`.
Repository processing messages identify the repository and report when processing starts and completes. Commands may
also report meaningful operation stages, such as workflow collection or report writing. Result documents and summaries
remain on standard output.

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
