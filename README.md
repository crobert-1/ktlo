# ktlo

`ktlo` runs enabled checks against repositories listed in a YAML configuration file. PR status results are written as
Markdown under `results/github/<owner>/<repository>/DEPENDENCY-BUMP-PRS.md`.

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
./bin/ktlo --config config.yaml
```

## Go maintenance targets

Run `make gotidy`, `make gofmt`, or `make golint` from the repository root to tidy modules, format Go files, or run
golangci-lint across all Go modules. Each Go module directory needs a `Makefile` that includes the shared
`Makefile.Common` using a relative path from that module directory. The lint target requires `golangci-lint` to be
installed and available on `PATH`. `make addlicense` adds the repository's Apache 2.0 header to Go files, and
`make checklicense` checks that each Go file has it; both targets require `addlicense` on `PATH`.

## Run modes

The configuration selects which checks run. Each enabled check runs for every configured repository, and both checks
can run in the same invocation.

- **PR status** (`prStatus.enabled`, enabled by default) checks open dependency-update pull requests and their CI status.
  It uses an existing local flaky-test report for failure details when available, but does not run analysis or update that
  report. It writes a Markdown report to `results/github/<owner>/<repository>/DEPENDENCY-BUMP-PRS.md`, with separate
  tables for passed, pending, failed, other, and stale PRs. Missing CI results and in-progress checks appear under
  Pending, which takes precedence when other checks have failed.
- **Flaky test JSON report** (`failingTestCheck.jsonReport.enabled`, disabled by default) collects workflow runs created
  within the configured lookback, parses Go test logs, and writes
  `results/github/<owner>/<repository>/workflow-test-analysis.json`.
- **Flaky test README** (`failingTestCheck.readme.enabled`, disabled by default) reads the local JSON report and writes
  `results/github/<owner>/<repository>/FLAKY-AND-FAILING-TESTS.md`. It lists flaky tests, failing tests, and workflow/job/step failures
  observed on `main`, with failure counts, rates, and a recent failure example. If both failing-test modes are enabled,
  JSON collection runs before the README is generated. README-only mode uses the existing JSON report and does not call
  GitHub APIs.

PR Status and JSON report generation require a token from `GH_TOKEN` or its `GITHUB_TOKEN` fallback. README-only mode
can run without a token when PR Status and JSON report generation are disabled.

## Runtime instrumentation

Tracing and duration metrics are opt-in. To export both signals to an OTLP gRPC collector, set
`otelGoInstrumentation.enabled` to `true` in the YAML configuration. Set `OTEL_EXPORTER_OTLP_ENDPOINT` to configure
both signals, or use `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` and `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` to configure
traces and metrics separately. These YAML keys map to the corresponding environment variables used by the Go OpenTelemetry
SDK. Use the matching `*_INSECURE` keys for plaintext gRPC collectors. Configured values override matching process
environment variables; omitted values preserve them. For a local plaintext gRPC collector, use this configuration:

```yaml
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

The YAML file passed with `--config` selects which checks to run, their settings, and the repositories to check. The
same options are shown in [config.example.yaml](config.example.yaml):

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

- `repositories` (default: none): GitHub repositories to check. Each entry must be an object with a required `url`,
  optional `enabled` (default: `true`), and optional `localPath`; URL-only string entries fail config validation.
  Disabled repositories are skipped by every check. `localPath` supplies an existing local clone for Vulnerability Check;
  ktlo scans its current working tree as-is without fetching, switching branches, or modifying repository files or Git
  metadata. Its `origin` remote must match the configured URL. The scan may include uncommitted changes, while report
  links point to the checkout's current `HEAD` commit. Relative paths are resolved from the process working directory. At least
  one repository must be enabled when a check is enabled.
- `prStatus`
  - `enabled` (default: `true`): Runs the PR status check.
  - `dependencyUpdateAuthors` (default: `dependabot`, `renovate`, `srv-gh-o11y-gdi`): Author logins identifying
    dependency-update pull requests included in PR status results.
  - `parseWorkflowFailures` (default: `false`): Fetches failed GitHub Actions runs linked from failed PR checks and
    parses their jobs, steps, and tests with the same parsers used by Failing Test Check. Parsed observations are used
    only in memory to improve failure matching; this option does not run the repository-wide lookback or write the
    flaky-test JSON report. It uses `failingTestCheck.jsonReport.lookback` to bound eligible workflow attempts.
  - `staleAge` (default: `7d`): Moves PRs whose latest commit is older than this duration into the `Stale` table.
    Accepts positive durations such as `60m`, `12h`, and `7d`.
- `failingTestCheck`
  - `jsonReport`
    - `enabled` (default: `false`): Collects workflow observations and writes the versioned flaky-test JSON report.
    - `lookback` (default: `7d`): How far back to search for workflow runs. GitHub filters by run creation time; ktlo then
      includes attempts started at or after the cutoff. A run created before the cutoff and rerun during the lookback may
      be missed. GitHub limits filtered searches to 1,000 results; exceeding that limit marks analysis incomplete.
      Accepts Go duration units such as `30m` and `12h`, plus `d` for days. Must be positive.
    - `overwriteExistingReport` (default: `false`): Rebuilds the report from the full lookback instead of reusing stored
      workflow attempts. When `false`, existing reports are incrementally updated; an unsupported schema version stops
      the script with instructions to enable overwriting.
  - `readme`
    - `enabled` (default: `false`): Writes a Markdown summary from the local JSON report. If JSON report generation is
      disabled, the report must already exist. The README reports observations on `main`; unknown classifications are
      shown with their reason. Its recent-failure column links to the source workflow or job.
  - `maxArtifactAgeHours` (default: `48`): Maximum age, in hours, for consuming a flaky-analysis report in PR status.
    `0` disables the freshness limit.

Unknown YAML keys are rejected. Both failing-test modes are opt-in; `prStatus` runs by default unless explicitly disabled.
To generate only the failing-test README from an existing local report without a GitHub token, also set
`prStatus.enabled: false` and `failingTestCheck.jsonReport.enabled: false`.

## License

This repository is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).
