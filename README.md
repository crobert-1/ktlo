# ktlo

`ktlo` provides commands to check dependency-update pull requests, analyze flaky tests, generate test reports, and assess
vulnerabilities across enabled GitHub repositories listed in a YAML configuration file. Commands write Markdown or JSON
results under `results/github/<owner>/<repository>/`.

## Prerequisite

Set `GH_TOKEN` to a GitHub token that can read the target repository and its pull requests. `GITHUB_TOKEN` is also
accepted as a fallback.

```sh
export GH_TOKEN="your-github-token"
```

## Build and run

```sh
make ktlo
cp config.example.yaml config.yaml
./bin/ktlo --config config.yaml dependency-pr-status update
```

## CLI commands

```text
ktlo [--config FILE] <command> [arguments]

Commands:
  vulnerability-check (v) [ID]
  flaky-test-analysis (f)
  generate-report (g)
  dependency-pr-status (d) update (u)
  help [COMMAND]
```

Use `ktlo --help` for the command list and `ktlo <command> --help` for command options. The config option also accepts
the short form `-c` and may appear before or after the command. For example:

```sh
./bin/ktlo -c config.yaml v CVE-2026-1234
./bin/ktlo f --lookback 30d --overwrite
./bin/ktlo g
./bin/ktlo d u --stale-age 14d --no-parse-workflow-failures
```

## Run modes

The command selects the operation. Each command runs for every enabled repository in the config file.

- **Dependency PR status** (`dependency-pr-status update`, shortcut `d u`) checks open dependency-update pull requests
  and their CI status. It uses a valid local flaky-test report for failure details when available, but does not run
  analysis or update that report. It writes a Markdown report to
  `results/github/<owner>/<repository>/DEPENDENCY-BUMP-PRS.md`, with separate
  tables for passed, pending, failed, other, and stale PRs. Missing CI results and in-progress checks appear under
  Pending, which takes precedence when other checks have failed.
- **Flaky test analysis** (`flaky-test-analysis`, shortcut `f`) collects workflow runs created within the configured
  lookback, parses Go test logs, and writes
  `results/github/<owner>/<repository>/workflow-test-analysis.json`.
- **Generate report** (`generate-report`, shortcut `g`) reads the local JSON report and writes
  `results/github/<owner>/<repository>/FLAKY-AND-FAILING-TESTS.md`. It lists flaky tests, failing tests, and
  workflow/job/step failures observed on `main`, with failure counts, rates, and a recent failure example. This command
  does not call GitHub APIs.
- **Vulnerability check** (`vulnerability-check`, shortcut `v`) analyzes the supplied vulnerability ID across the
  configured repositories and writes detailed reports. Pass one supported vulnerability ID, such as a CVE, GHSA, or
  Go vulnerability ID, after the command.

  **Scope warning:** The check assesses dependencies at the latest commit on `main` only. If a dependency was removed,
  it may have been vulnerable in an earlier commit, but it is no longer present in the scanned snapshot and that past
  exposure is not reported.

`dependency-pr-status update` parses linked failed workflow runs by default. Add `--no-parse-workflow-failures` to skip
that work for a run. Use
`--lookback` and `--overwrite` with `flaky-test-analysis` to override the configured lookback or rebuild the report.
Use `--stale-age` with dependency PR status to override the configured stale PR threshold for one run.

Vulnerability checks, flaky-test analysis, and dependency PR status require a token from `GH_TOKEN` or its
`GITHUB_TOKEN` fallback. Report generation reads local data and does not require a token.

## Runtime instrumentation

Tracing and duration metrics are opt-in. To export both signals to an OTLP gRPC collector, set
`otelGoInstrumentation.enabled` to `true` in the YAML configuration. Set `OTEL_EXPORTER_OTLP_ENDPOINT` to configure
both signals, or use `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` and `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` to configure
traces and metrics separately. These YAML keys map to the corresponding environment variables used by the Go OpenTelemetry
SDK. Use the matching `*_INSECURE` keys for plaintext gRPC collectors. Configured values override matching process
environment variables; omitted values preserve them. For a local plaintext gRPC collector, use this configuration:

```yaml
logLevel: default

otelGoInstrumentation:
  enabled: true
  OTEL_EXPORTER_OTLP_TRACES_ENDPOINT: dns://localhost:4317
  OTEL_EXPORTER_OTLP_TRACES_INSECURE: "true"
  OTEL_EXPORTER_OTLP_METRICS_ENDPOINT: dns://localhost:4317
  OTEL_EXPORTER_OTLP_INSECURE: "true"
```

Traces include spans for CLI stages, GitHub API methods and HTTP requests, workflow attempt log archive collection and parsing, and
report persistence. `github.get_rest_json` spans include the HTTP method and full request URL, including path IDs and
query parameters such as pagination and lookback filters. The `ktlo.operation.duration` histogram records method durations in seconds with a bounded
`ktlo.operation` attribute. The HTTP instrumentation also reports `http.client.request.duration`. IDs and repository
names appear only on trace spans, not as metric dimensions. The CLI flushes telemetry before exiting.

## Configuration

The YAML file passed with `--config` stores reusable settings and the repositories to check. The command selects which
operation runs. `--config` defaults to `config.yaml`. Unknown YAML keys, obsolete URL-only repository entries, and
configurations without at least one enabled repository fail validation.

### Options

- `otelGoInstrumentation` (optional; default: disabled): Configures OpenTelemetry export.
  - `enabled` (optional; default: `false`): Enables tracing and duration metrics.
  - `OTEL_EXPORTER_OTLP_ENDPOINT` (optional; default: unset): Base OTLP gRPC endpoint for traces and metrics.
  - `OTEL_EXPORTER_OTLP_INSECURE` (optional; default: unset): Configures insecure transport for the base endpoint.
  - `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` (optional; default: unset): Overrides the base endpoint for metrics.
  - `OTEL_EXPORTER_OTLP_METRICS_INSECURE` (optional; default: unset): Configures insecure transport for the metrics endpoint.
  - `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` (optional; default: unset): Overrides the base endpoint for traces.
  - `OTEL_EXPORTER_OTLP_TRACES_INSECURE` (optional; default: unset): Configures insecure transport for the traces endpoint.
  - Empty or omitted exporter values leave corresponding process environment variables unchanged; non-empty configured
    values override them.
- `prStatus` (optional): Configures Dependency PR Status.
  - `dependencyUpdateAuthors` (optional; default: `dependabot`, `renovate`, `srv-gh-o11y-gdi`): GitHub author logins
    used to identify dependency-update PRs.
  - `staleAge` (optional; default: `7d`): Marks PRs whose latest commit is older than this duration as stale. Accepts
    positive durations such as `60m`, `12h`, and `7d`; `--stale-age` overrides it for one run.
- `failingTestCheck` (optional): Configures flaky-test analysis.
  - `jsonReport` (optional): Configures the stored workflow test report.
    - `lookback` (optional; default: `7d`): How far back to search for workflow runs. Accepts positive durations such
      as `30m`, `12h`, and `7d`; `--lookback` overrides it for one run.
- `logLevel` (optional; default: `default`): Controls progress output written to standard error. `none` suppresses
  progress messages, `default` shows command and repository progress, and `detailed` also shows operation-stage updates.
  Errors remain visible at every level.
- `repositories` (required; no default): Repositories checked by commands. Supply at least one entry and enable at least
  one repository.
  - Each entry is an object with these options:
    - `url` (required; no default): GitHub repository URL.
    - `enabled` (optional; default: `true`): Whether checks include this repository.
    - `localPath` (optional; no default): Path to an existing local checkout for Vulnerability Check. Its `origin` must
      match `url`. KTLO scans the working tree as-is without fetching, switching branches, or modifying files or Git
      metadata. Uncommitted changes may be included while links point to the checkout's current `HEAD`. Relative paths
      are resolved from the process working directory.

### Example

See [config.example.yaml](config.example.yaml) for a complete example:

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

Commands operate on all enabled repositories unless the config disables an entry. `--config` may appear before or after
the command.

## License

This repository is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).
