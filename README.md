# ktlo

`ktlo` runs one selected command against repositories listed in a YAML configuration file. PR status results are written as
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

## Go maintenance targets

Run `make gotidy`, `make gofmt`, or `make golint` from the repository root to tidy modules, format Go files, or run
golangci-lint across all Go modules. Each Go module directory needs a `Makefile` that includes the shared
`Makefile.Common` using a relative path from that module directory. The lint target requires `golangci-lint` to be
installed and available on `PATH`. `make addlicense` adds the repository's Apache 2.0 header to Go files, and
`make checklicense` checks that each Go file has it; both targets require `addlicense` on `PATH`.

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
  configured repositories and writes detailed reports. Pass one CVE or Go vulnerability ID after the command.

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
operation runs. The example configuration is shown in [config.example.yaml](config.example.yaml):

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

- `repositories` (default: none): GitHub repositories to check. Each entry must be an object with a required `url`,
  optional `enabled` (default: `true`), and optional `localPath`; URL-only string entries fail config validation.
  Disabled repositories are skipped by every check. `localPath` supplies an existing local clone for Vulnerability Check;
  ktlo scans its current working tree as-is without fetching, switching branches, or modifying repository files or Git
  metadata. Its `origin` remote must match the configured URL. The scan may include uncommitted changes, while report
  links point to the checkout's current `HEAD` commit. Relative paths are resolved from the process working directory. At least
  one repository must be enabled for a command.
- `prStatus`
  - `dependencyUpdateAuthors` (default: `dependabot`, `renovate`, `srv-gh-o11y-gdi`): Author logins identifying
    dependency-update pull requests included in PR status results.
  - `staleAge` (default: `7d`): Moves PRs whose latest commit is older than this duration into the `Stale` table.
    Accepts positive durations such as `60m`, `12h`, and `7d`; `--stale-age` overrides it for one run.
- `failingTestCheck`
  - `jsonReport`
    - `lookback` (default: `7d`): How far back to search for workflow runs. GitHub filters by run creation time; ktlo then
      includes attempts started at or after the cutoff. A run created before the cutoff and rerun during the lookback may
      be missed. GitHub limits filtered searches to 1,000 results; exceeding that limit marks analysis incomplete.
      Accepts Go duration units such as `30m` and `12h`, plus `d` for days. Must be positive; `--lookback` overrides
      it for one analysis run. `--overwrite` rebuilds from the full lookback instead of reusing stored workflow
      attempts. An unsupported schema version stops the analysis and explains the overwrite option.

Unknown YAML keys are rejected. Commands operate on all enabled repositories unless the config disables an entry.
`--config` defaults to `config.yaml` and may appear before or after the command.

## License

This repository is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).
